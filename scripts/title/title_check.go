// title_check.go enforces the pull request title conventions against the
// changed files: every registered entity a changed flat source file
// implements must be announced in the title, docs-only and testcase-only
// changes must say so, and no title may truncate with an ellipsis.
//
// A flat filename is not a reliable entity name — shared plumbing files
// declare no constructor, and stock files can declare constructors
// registered under a different TypeName — so the entities a flat file
// implements are resolved through the provider.go registration join
// (scripts/internal/flatlayout), the same source of truth the consistency
// and testing coverage gates use.
package main

import (
	"flag"
	"os"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/flatlayout"
	log "github.com/sirupsen/logrus"
	"github.com/waigani/diffparser"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.DisableTimestamp = false
	customFormatter.DisableColors = false
	customFormatter.ForceColors = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var (
	fileNames = flag.String("fileNames", "", "the files to check diff")
	title     = flag.String("title", "", "the pull request title to check")
)

// entity is one registered provider entity a changed flat source file
// implements.
type entity struct {
	typeName   string // e.g. "alicloud_cloudauth_face_configs"
	dataSource bool   // registered in DataSourcesMap rather than ResourcesMap
}

// titleFile is what the title policy needs to know about one diff file.
type titleFile struct {
	name string
	// added marks a newly added file: its entities demand the
	// "New Resource/New Data Source" title entries instead of the
	// "resource/<name>/data-source/<name>" ones.
	added bool
	// testcase, docs and service mirror the file classification the
	// previous title gates applied; flat marks an alicloud flat source
	// file.
	testcase, docs, service, flat bool
	// entities are the registered entities a flat file implements,
	// resolved through the provider.go registration join; an empty list
	// means the title gate has nothing to demand for it.
	entities []entity
}

func main() {
	flag.Parse()
	log.Infof("==> PR Title: %s", *title)
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := os.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	// The provider.go map literals are the only flat registration source;
	// the constructor → registration table is shared by every flat file
	// in the diff.
	providerSrc, err := os.ReadFile("alicloud/provider.go")
	if err != nil {
		log.Errorf("read alicloud/provider.go failed. Error: %s", err)
		os.Exit(1)
	}
	files := resolveTitleFiles(diff.Files, flatlayout.ProviderRegistrations(providerSrc))
	if !checkTitle(*title, files) {
		log.Errorf("Pull request title validation failed.")
		os.Exit(1)
	}
}

// resolveTitleFiles classifies the diff files and joins each flat source
// file to the entities it implements. Deleted files impose no title
// requirements — a removed entity has nothing to announce — so they are
// dropped here, matching the status filters the previous gates applied.
func resolveTitleFiles(files []*diffparser.DiffFile, flatRegistrations map[string][]flatlayout.Registration) []titleFile {
	var resolved []titleFile
	for _, file := range files {
		if file.Mode == diffparser.DELETED {
			continue
		}
		tf := titleFile{name: file.NewName, added: file.Mode == diffparser.NEW}
		tf.testcase, tf.docs, tf.service, tf.flat = classify(file.NewName)
		if tf.flat {
			tf.entities = flatEntities(file.NewName, flatRegistrations)
		}
		resolved = append(resolved, tf)
	}
	return resolved
}

// classify mirrors the classification order of the previous title gates:
// test files first, then docs, service sources and flat sources.
func classify(name string) (testcase, docs, service, flat bool) {
	if strings.HasSuffix(name, "_test.go") {
		return true, false, false, false
	}
	if strings.HasPrefix(name, "website/docs") {
		return false, true, false, false
	}
	if strings.HasPrefix(name, "alicloud/service/") && strings.HasSuffix(name, ".go") {
		return false, false, true, false
	}
	if strings.HasPrefix(name, "alicloud/resource_alicloud") || strings.HasPrefix(name, "alicloud/data_source_alicloud") {
		return false, false, false, true
	}
	return false, false, false, false
}

// flatEntities resolves the registered entities one changed flat source
// file implements: the constructors it declares joined against the
// provider.go registrations, one entity per TypeName so an alias pair
// backed by a single constructor yields both. A shared plumbing file
// implements none; a file the join cannot read falls back to the
// filename-derived name the previous title gates required.
func flatEntities(name string, flatRegistrations map[string][]flatlayout.Registration) []entity {
	joined, err := flatlayout.JoinWith(name, flatRegistrations)
	if err == nil {
		entities := make([]entity, 0, len(joined.Registrations))
		for _, r := range joined.Registrations {
			entities = append(entities, entity{typeName: r.TypeName, dataSource: r.DataSource})
		}
		return entities
	}
	log.Warnf("fail to join flat-layout file %s to its registrations, falling back to the filename-derived name. Error: %s", name, err)
	derived := entity{}
	base := strings.TrimSuffix(name[strings.LastIndex(name, "/")+1:], ".go")
	derived.dataSource = strings.HasPrefix(base, "data_source_")
	derived.typeName = strings.TrimPrefix(strings.TrimPrefix(base, "resource_"), "data_source_")
	return []entity{derived}
}

// checkTitle enforces the title conventions: every registered entity a
// changed flat source implements must be announced in the title, a
// docs-only change must start with the docs marker, a testcase-only one
// must carry the testcase marker, and no title may truncate with an
// ellipsis.
func checkTitle(prTitle string, files []titleFile) bool {
	failed := false
	reportError := func(format string, args ...interface{}) {
		log.Errorf(format, args...)
		failed = true
	}

	if strings.Contains(prTitle, "…") || strings.HasSuffix(prTitle, "...") {
		reportError("The PR title should not omit important infos about added or modified files. xxx... is not recommended.")
	}

	docsOnly, testcaseOnly := true, true
	docsChanged, testcaseChanged := false, false
	for _, file := range files {
		switch {
		case file.testcase:
			testcaseChanged = true
		case file.docs:
			docsChanged = true
			testcaseOnly = false
		case file.service:
			docsOnly = false
			testcaseOnly = false
		case file.flat:
			docsOnly = false
			testcaseOnly = false
			for _, e := range file.entities {
				if file.added {
					if e.dataSource {
						if !strings.Contains(prTitle, "New Data Source: "+e.typeName) {
							reportError("The PR title should contains new datasource info New Data Source: %s", e.typeName)
						}
					} else if !strings.Contains(prTitle, "New Resource: "+e.typeName) {
						reportError("The PR title should contains new resource info New Resource: %s", e.typeName)
					}
					continue
				}
				titleStr := "resource/"
				if e.dataSource {
					titleStr = "data-source/"
				}
				titleStr += e.typeName + ": "
				if !strings.Contains(prTitle, titleStr) {
					reportError("The PR title should contains modified info like %sxxx", titleStr)
				}
			}
		}
	}

	if docsChanged && docsOnly {
		if !strings.HasPrefix(prTitle, "docs: ") {
			reportError("The PR title should contains docs modified info like docs: xxx")
		}
	} else if testcaseChanged && testcaseOnly {
		if !strings.Contains(prTitle, "testcase: ") {
			reportError("The PR title should contains testcase modified info like testcase: xxx")
		}
	}
	return !failed
}
