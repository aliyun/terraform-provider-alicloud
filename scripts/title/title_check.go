// title_check.go enforces the pull request title conventions against the
// changed files: every registered entity a changed provider source file
// implements must be announced in the title, docs-only and testcase-only
// changes must say so, and no title may truncate with an ellipsis.
//
// The changed files and their statuses come from `git diff --name-status`
// between the -base and -head commits, so the gate needs no diff file
// hand-off. Flat filenames yield their entity names directly, the way the
// previous title gates derived them. Split-layout filenames cannot: the
// product slug only coincidentally matches the resource-name prefix, and
// a stem like ims/resource_user.go names no entity on its own, so each
// changed service source is joined to the registrations of its product's
// register.go through scripts/internal/splitlayout — the same source of
// truth the consistency and testing coverage gates use.
//
// A service entity is announced as new only when it did not exist at the
// base commit, in either layout: a migration that moves a flat resource
// into the split layout adds implementation files but still announces the
// move with the `resource/<name>:` form, the way PR #10817 and #10831
// did. Ephemeral, action and list registrations have no house title form
// yet; they only clear the docs-only and testcase-only states, and so do
// register.go and shared helper files that implement no registration.
//
// The gate runs as the base-restored copy of itself in CI — the workflows
// check out the pull request head, then restore scripts/title and
// scripts/internal/splitlayout from the base so master's code judges the
// pull request data — and it must run there without network access, so it
// imports only the standard library and the in-module splitlayout
// package. -resolve turns the same register.go join into a lookup for the
// docs example gate: it prints the entities service implementation files
// resolve to, because a split-layout filename cannot yield the docs path.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

var (
	title   = flag.String("title", "", "the pull request title to check")
	base    = flag.String("base", "", "the base commit of the change range")
	head    = flag.String("head", "", "the head commit of the change range")
	resolve = flag.Bool("resolve", false, "resolve service implementation files to their registered entities and exit")
)

// entity is one registered provider entity a changed source file
// implements.
type entity struct {
	typeName   string // e.g. "alicloud_cloudauth_face_configs"
	dataSource bool   // a data source rather than a managed resource
	// existing reports that the entity already existed before the change:
	// an entity that moved from the flat layout into the split one is
	// announced with the modified form, not the new-entity one.
	existing bool
}

// titleFile is what the title policy needs to know about one diff file.
type titleFile struct {
	name string
	// added marks a newly added file: flat entities derived from it are
	// announced as new. Service entities decide through the base commit's
	// entity index instead, so added carries no meaning for them.
	added bool
	// testcase, docs and service mirror the file classification the
	// previous title gates applied; flat marks an alicloud flat source
	// file.
	testcase, docs, service, flat bool
	// entities are the registered entities the file implements; an empty
	// list means the title gate has nothing to demand for it.
	entities []entity
}

func main() {
	flag.Parse()
	if *resolve {
		resolveMode(flag.Args())
		return
	}
	if *base == "" || *head == "" {
		fmt.Fprintln(os.Stderr, "ERROR: both -base and -head are required")
		os.Exit(1)
	}
	fmt.Printf("==> PR Title: %s\n", *title)
	files, err := changedFiles(*base, *head)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("==> the change range is empty, shipped!")
		return
	}
	if err := resolveEntities(files, *base); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if !checkTitle(*title, files) {
		fmt.Fprintln(os.Stderr, "ERROR: Pull request title validation failed.")
		os.Exit(1)
	}
}

// changedFiles lists the added, modified and renamed files of the change
// range with their classification. Renamed files classify by their new
// name, the way the previous title gates did; deleted files impose no
// title requirements — a removed entity has nothing to announce — so the
// status filter drops them before they ever reach the parser.
func changedFiles(base, head string) ([]titleFile, error) {
	out, err := exec.Command("git", "diff", "--name-status", "-M", "--diff-filter=AMR", base, head).Output()
	if err != nil {
		return nil, fmt.Errorf("git diff --name-status %s %s failed: %w", base, head, err)
	}
	return parseNameStatus(string(out)), nil
}

// parseNameStatus parses `git diff --name-status` output: one record per
// line, the status letter, then the path — two paths for a rename, of
// which the new one classifies the file.
func parseNameStatus(out string) []titleFile {
	var files []titleFile
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 2 || fields[0] == "D" {
			continue
		}
		tf := titleFile{name: fields[len(fields)-1], added: fields[0] == "A"}
		tf.testcase, tf.docs, tf.service, tf.flat = classify(tf.name)
		files = append(files, tf)
	}
	return files
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

// resolveEntities joins each classified source file to the entities it
// implements: flat files derive theirs from the filename, service source
// files through the register.go join of both dialects, with the
// new-versus-existing form decided by the base commit's entity index.
func resolveEntities(files []titleFile, base string) error {
	var index *baseIndex // built lazily: only service entities consult it
	for i := range files {
		switch {
		case files[i].flat:
			files[i].entities = flatEntities(files[i].name, files[i].added)
		case files[i].service:
			if index == nil {
				var err error
				index, err = loadBaseIndex(base)
				if err != nil {
					return err
				}
			}
			files[i].entities = serviceEntities(files[i].name, index)
		}
	}
	return nil
}

// flatEntities derives the entity a changed flat source file announces
// from its filename, the way the previous title gates did: the stem after
// the layout prefix is the entity name, an added file announces a new
// entity and a modified one an existing one.
func flatEntities(name string, added bool) []entity {
	stem := strings.TrimSuffix(name[strings.LastIndex(name, "/")+1:], ".go")
	return []entity{{
		typeName:   strings.TrimPrefix(strings.TrimPrefix(stem, "resource_"), "data_source_"),
		dataSource: strings.HasPrefix(stem, "data_source_"),
		existing:   !added,
	}}
}

// serviceEntities joins one changed service source file to the registered
// entities it implements. Resource and data-source registrations demand
// their title entries; ephemeral, action and list registrations have no
// house title form yet, and register.go or shared helper files implement
// no registration at all, so all of those only clear the docs-only and
// testcase-only states. Whether an entity counts as new is decided by the
// base commit's index, not by the file's diff status: a migration into
// the split layout adds implementation files for entities that already
// exist.
func serviceEntities(name string, index *baseIndex) []entity {
	var entities []entity
	for _, r := range joinRegistrations(name) {
		switch r.Kind {
		case splitlayout.KindResource:
			entities = append(entities, entity{typeName: r.TypeName, existing: index.has(r.TypeName)})
		case splitlayout.KindDataSource:
			entities = append(entities, entity{typeName: r.TypeName, dataSource: true, existing: index.has(r.TypeName)})
		}
	}
	return entities
}

// joinRegistrations resolves one service source path through the
// register.go join of both dialects. Files that are not implementation
// stems, and files whose registrations cannot be read, resolve to
// nothing: the gate stays quiet about them — flagging missing
// registrations is the consistency gate's business.
func joinRegistrations(name string) []splitlayout.Resolved {
	product, stem, isTest, ok := splitlayout.IsSplitSource(name)
	if !ok || isTest {
		return nil
	}
	joined, err := splitlayout.Resolve(product, stem)
	if err != nil {
		fmt.Printf("INFO: %s resolves no registration: %s\n", name, err)
	}
	frameworkJoined, err := splitlayout.FrameworkResolve(product, stem)
	if err != nil {
		fmt.Printf("INFO: %s resolves no framework registration: %s\n", name, err)
	}
	joined = append(joined, frameworkJoined...)
	if len(joined) == 0 {
		fmt.Printf("INFO: no registration resolves %s; no title entries demanded\n", name)
	}
	return joined
}

// resolveMode prints the registered entities the given service
// implementation files resolve to, one `<path>\t<TypeName>\t<Kind>` line
// per registration, for the docs example gate to derive documentation
// paths from — a split-layout filename cannot yield them itself. Files
// that resolve to nothing print nothing; the exit status stays 0.
func resolveMode(paths []string) {
	for _, path := range paths {
		for _, line := range resolveLines(path) {
			fmt.Println(line)
		}
	}
}

// resolveLines formats the resolve-mode output for one path.
func resolveLines(path string) []string {
	var lines []string
	for _, r := range joinRegistrations(path) {
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s", path, r.TypeName, r.Kind))
	}
	return lines
}

// baseIndex is the set of provider entity names that existed at the base
// commit, in either layout: the flat source files directly, and the
// registrations of every product register.go. It decides whether a
// service entity is announced as new or as existing.
type baseIndex struct {
	names map[string]bool
}

// has reports whether typeName existed at the indexed base commit.
func (b *baseIndex) has(typeName string) bool {
	return b != nil && b.names[typeName]
}

// loadBaseIndex builds the entity index of one base commit from git: the
// flat sources directly under alicloud/, and every register.go under
// alicloud/service/. A base without the split layout lists no products
// and simply carries the flat names.
func loadBaseIndex(base string) (*baseIndex, error) {
	flat, err := gitLines("ls-tree", "--name-only", base, "--", "alicloud/")
	if err != nil {
		return nil, err
	}
	sources := make(map[string][]byte)
	if serviceFiles, err := gitLines("ls-tree", "-r", "--name-only", base, "--", "alicloud/service"); err != nil {
		return nil, err
	} else {
		for _, path := range serviceFiles {
			if !strings.HasSuffix(path, "/register.go") {
				continue
			}
			src, err := exec.Command("git", "show", base+":"+path).Output()
			if err != nil {
				return nil, fmt.Errorf("git show %s:%s failed: %w", base, path, err)
			}
			sources[path] = src
		}
	}
	return baseIndexFromSources(flat, sources), nil
}

// baseIndexFromSources builds the entity index from the flat source paths
// and the register.go sources of one tree; the git plumbing in
// loadBaseIndex only gathers those.
func baseIndexFromSources(flatFiles []string, registers map[string][]byte) *baseIndex {
	index := &baseIndex{names: make(map[string]bool)}
	for _, file := range flatFiles {
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		stem := strings.TrimSuffix(file[strings.LastIndex(file, "/")+1:], ".go")
		index.names[strings.TrimPrefix(strings.TrimPrefix(stem, "resource_"), "data_source_")] = true
	}
	for _, src := range registers {
		entries := append(splitlayout.SDKRegisterEntries(src), splitlayout.FrameworkRegisterEntries(src)...)
		for _, e := range entries {
			index.names[e.TypeName] = true
		}
	}
	return index
}

// checkTitle enforces the title conventions: every registered entity a
// changed source file implements must be announced in the title — the
// new-entity form for entities the base did not know, the modified form
// otherwise — a docs-only change must start with the docs marker, a
// testcase-only one must carry the testcase marker, and no title may
// truncate with an ellipsis.
func checkTitle(prTitle string, files []titleFile) bool {
	failed := false
	reportError := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
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
		case file.service, file.flat:
			docsOnly = false
			testcaseOnly = false
		}
		for _, e := range file.entities {
			if !e.existing {
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

// gitLines runs one git command and returns its non-empty output lines.
func gitLines(args ...string) ([]string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}
	var lines []string
	for _, line := range strings.Split(string(out), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
