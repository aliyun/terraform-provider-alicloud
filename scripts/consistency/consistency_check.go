//nolint:all
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	providerregistry "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/frameworkschema"
	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
	set "github.com/deckarep/golang-set"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	resourceNames = flag.String("resourceNames", "", "the names of the terraform resources to diff")
	fileNames     = flag.String("fileNames", "", "the files to check diff")
	filterList    = map[string][]string{
		"alicloud_amqp_instance":            {"logistics"},
		"alicloud_cms_alarm":                {"notify_type"},
		"alicloud_cs_serverless_kubernetes": {"private_zone", "create_v2_cluster"},
		"alicloud_slb_listener":             {"lb_protocol", "instance_port", "lb_port"},
		"alicloud_kvstore_instance":         {"connection_string"},
		"alicloud_instance":                 {"subnet_id"},
		"alicloud_hbr_ots_backup_plan":      {"vault_id"},
		"alicloud_nat_gateway":              {"vswitch_id"},
		"alicloud_ecs_disk":                 {"advanced_features", "encrypt_algorithm", "dedicated_block_storage_cluster_id"},
	}
	skippedSchemaKeys = map[string]struct{}{
		"page_size":   {},
		"page_number": {},
		"total_count": {},
		"max_results": {},
	}
	// Attribute-phase modifiers are matched only inside a CLOSED parenthesized
	// modifier group whose members are modifier tokens — "(Optional)",
	// "(Required, ForceNew)" — not free-form prose. attributeFieldRegex
	// captures the whole description tail in words[2], so a plain Contains
	// would mistake prose like "Optional metadata returned by the service."
	// or "Required for subsequent API calls." for a modifier, and the prior
	// \([^)]*\b<modifier>\b would also match an unclosed "(Optional" or a
	// prose parenthesis like "(Required for subsequent API calls)" /
	// "(Optional metadata returned by the service)". The group regex requires
	// a closing ")" and only modifier tokens (Optional/Required/ForceNew/
	// Deprecated, comma-separated) inside, so prose and unclosed parens are
	// rejected; multi-token groups and Deprecated's legacy Contains behavior
	// are preserved.
	attrModifierGroupRe = regexp.MustCompile(`\(\s*(?:Optional|Required|ForceNew|Deprecated)(?:\s*,\s*(?:Optional|Required|ForceNew|Deprecated))*\s*\)`)
)

// checkReport aggregates the per-entity outcome. A diff can touch more than a
// thousand docs pages, so the run must end with one line that says how much was
// actually scanned: "failed" alone cannot tell a gate that checked two entries
// apart from one that checked all of them.
type checkReport struct {
	entities         int
	compared         int
	passed           int
	docsUnreadable   int
	schemaUnreadable int
	unresolved       int
	frameworkChecked int
	duplicates       int
	failedNames      []string
}

func (r *checkReport) fail(resourceName string) {
	for _, name := range r.failedNames {
		if name == resourceName {
			return
		}
	}
	r.failedNames = append(r.failedNames, resourceName)
}

func (r *checkReport) print() {
	failed := len(r.failedNames)
	log.Infof("================ Consistency summary ================")
	log.Infof("entities=%d compared=%d passed=%d failed=%d framework=%d docs_unreadable=%d schema_unreadable=%d unresolved=%d duplicates=%d",
		r.entities, r.compared, r.passed, failed, r.frameworkChecked, r.docsUnreadable, r.schemaUnreadable, r.unresolved, r.duplicates)
	if r.duplicates > 0 {
		log.Warnf("%d attribute(s) are documented in both the `## Argument Reference` and the `## Attributes Reference`; remove them from the latter", r.duplicates)
	}
	if failed > 0 {
		log.Errorf("failed entities: %s", strings.Join(r.failedNames, ", "))
	}
}

type ResourceAttribute struct {
	Name        string
	Type        string
	Optional    string
	Required    string
	ForceNew    bool
	Default     string
	ElemType    string
	Deprecated  string
	DocsLineNum int
	Removed     string
}

func main() {
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, err := ioutil.ReadFile(*fileNames)
	if err != nil {
		log.Errorf("fail to read the diff file %s: %s", *fileNames, err)
		os.Exit(1)
	}
	diff, err := diffparser.Parse(string(byt))
	// A discarded parse error used to leave an empty file list behind and the
	// gate passed without comparing anything.
	if err != nil {
		log.Errorf("fail to parse the diff file %s: %s", *fileNames, err)
		os.Exit(1)
	}
	report, exitCode := runDiffCheck(diff, frameworkRegisteredNames(), sdkv2.Provider())
	report.print()
	if exitCode > 0 {
		os.Exit(exitCode)
	}
	return
}

// kindFromDocsDir maps the docs directory segment a page lives in to the kind
// of entity it documents. The directory is the authority the checker has for
// a docs-only path: a page under r/ documents the resource of that name, one
// under d/ the data source, and the same page name can exist under both for a
// name registered as each kind.
func kindFromDocsDir(path string) (splitlayout.Kind, bool) {
	segments := strings.Split(path, "/")
	if len(segments) != 4 {
		return "", false
	}
	switch segments[2] {
	case "r":
		return splitlayout.KindResource, true
	case "d":
		return splitlayout.KindDataSource, true
	case "ephemeral-resources":
		return splitlayout.KindEphemeral, true
	case "actions":
		return splitlayout.KindAction, true
	}
	return "", false
}

// frameworkDispatch resolves the framework kind to check one name under. The
// diff's own kind is the authority when the name is registered under it, and
// the one exception is the pair the docs directory cannot express: a list
// resource documents its page under r/, where the directory alone says
// "resource", so a resource-kind arrival of a name registered only as a list
// resource still dispatches to that registration. Every other kind is named
// by its own directory (d/, ephemeral-resources/, actions/) or by its split
// registration, so any other single-registration fallback would compare the
// diff's page against a different entity and leave the touched page
// unchecked — most damagingly for a name the SDKv2 provider registers under
// the diff's kind (alicloud_kms_secret is both an SDKv2 resource and a
// framework ephemeral): the hijack would run the framework check on the
// SDKv2 entity's arrival and its docs would never be compared. Such a name
// stays unclaimed and falls to the SDKv2 path below.
func frameworkDispatch(frameworkNames map[splitlayout.Kind]map[string]struct{}, kind splitlayout.Kind, name string) (splitlayout.Kind, bool) {
	if _, ok := frameworkNames[kind][name]; ok {
		return kind, true
	}
	var registered []splitlayout.Kind
	for k, names := range frameworkNames {
		if _, ok := names[name]; ok {
			registered = append(registered, k)
		}
	}
	if len(registered) == 1 && kind == splitlayout.KindResource && registered[0] == splitlayout.KindListResource {
		return registered[0], true
	}
	return "", false
}

// runDiffCheck walks the parsed diff against the registered entities and
// returns the aggregated report plus the process exit code. main supplies the
// production registrations — the framework names from the runtime registry and
// the SDKv2 provider — while unit tests inject fixtures and a fake source tree
// so the dispatch can be exercised without reading the real repo.
func runDiffCheck(diff *diffparser.Diff, frameworkNames map[splitlayout.Kind]map[string]struct{}, sdkProvider *schema.Provider) (*checkReport, int) {
	report := &checkReport{}
	exitCode := 0
	fileRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*_test.go")
	fileDocsRegex := regexp.MustCompile("website/docs/(r|d|ephemeral-resources|actions)/[0-9a-zA-Z_]*.html.markdown")
	resourceNameMap := make(map[string]struct{})
	// frameworkDocsDir maps a registration kind to the docs directory that
	// documents it. The repo has no list-resource directory, so a list resource
	// falls back to the resource directory and its page is reported missing when
	// it has none.
	frameworkDocsDir := func(kind splitlayout.Kind) string {
		switch kind {
		case splitlayout.KindDataSource:
			return "website/docs/d/"
		case splitlayout.KindEphemeral:
			return "website/docs/ephemeral-resources/"
		case splitlayout.KindAction:
			return "website/docs/actions/"
		default:
			return "website/docs/r/"
		}
	}
	// finish records the outcome of one comparison. Both schema dialects end
	// here so their counts cannot drift apart.
	finish := func(resourceName string, consistent bool) {
		report.compared++
		if consistent {
			log.Infof("--- PASS!\n\n")
			report.passed++
			return
		}
		log.Errorf("--- Failed!\n\n")
		report.fail(resourceName)
		exitCode = 1
	}
	// checkFrameworkName checks a framework-native entity. Its schema lives
	// outside the SDKv2 runtime maps, so it is parsed from the declaration in
	// the implementation file. The registration reverse-lookup supplies that
	// file: a docs-only diff carries no source path and the lookup fills it,
	// and a diff-supplied path that disagrees with the registration loses to
	// it — the registration pairs the type name with the factory that declares
	// its schema, so trusting the diff path instead would compare one entity's
	// schema against another's docs when the dispatch just switched the kind.
	// A disagreement is surfaced rather than swallowed. The kind comes from the
	// dispatch — the diff path is authoritative — so a same-named resource and
	// data source are each compared against their own schema and docs page. The
	// check is no weaker than the SDKv2 one — the same parse feeds the same
	// comparison — and an entity whose schema cannot be located fails the gate
	// instead of passing unexamined. The registration's file granularity is
	// guarded too: a file declaring several entities' factories has no
	// per-entity schema to parse, and the entity fails rather than being
	// compared with the union of the entities' attributes.
	checkFrameworkName := func(kind splitlayout.Kind, resourceName, sourcePath string) {
		resolved := splitlayout.FrameworkFindByKindAndTypeName(kind, resourceName)
		if resolved != nil {
			if sourcePath != "" && sourcePath != resolved.ImplPath {
				log.Warnf("the diff supplied %s as the implementation of %s (%s) but the registration resolves to %s; using the registered implementation", sourcePath, resourceName, kind, resolved.ImplPath)
			}
			sourcePath = resolved.ImplPath
		}
		report.frameworkChecked++
		log.Infof("==> Checking framework-native %s (%s) attributes consistency...", resourceName, kind)
		if sourcePath == "" {
			log.Errorf("cannot locate the implementation of the framework entity %s, so its docs cannot be compared", resourceName)
			report.unresolved++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			log.Errorf("reading the framework entity %s implementation %s failed: %s", resourceName, sourcePath, err)
			report.schemaUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		// The registration resolves file-grained: the package tolerates one
		// implementation file declaring several entities' factories, and then
		// the whole-file parse below would hand each of those entities the
		// union of their attributes — the same wrong-entity comparison the
		// (kind, name) dispatch exists to eliminate. A file that shared
		// factories cannot isolate the entities, so the gate fails the
		// entity instead of comparing a merged schema against its docs. A
		// guard that cannot prove either way — the source or the product's
		// register failed to parse or read — fails the entity too: it must
		// not be conflated with "no conflict".
		declared, err := splitlayout.FrameworkFactoriesDeclaredIn(sourcePath, source)
		if err != nil {
			log.Errorf("verifying which framework entities %s declares failed: %s; the docs of %s cannot be compared", sourcePath, err, resourceName)
			report.schemaUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		// The guard counts distinct factories, not registrations. Alias
		// registrations legitimately reuse one factory for several type names
		// (alicloud_disk/alicloud_ecs_disk, the sslcertificateservice legacy
		// aliases kept for backward compatibility), and a file declaring one
		// constructor exposes exactly that one schema no matter how many type
		// names point at it — the whole-file parse below is then the entity's
		// own schema, not a union. Counting entries instead would fail a
		// perfectly isolated aliased entity the moment the framework side
		// adopts the alias migration, the same class of false block this gate
		// exists to eliminate. A single factory must still be the checked
		// entity's own: when the only registered factory declared in the file
		// belongs to a different entity, comparing the file's schema against
		// this entity's docs would be a wrong-entity comparison, so the gate
		// fails the entity instead of guessing.
		declaredFactories := make(map[string]struct{}, len(declared))
		entityDeclared := false
		declaredNames := make([]string, 0, len(declared))
		for _, e := range declared {
			declaredFactories[e.Factory] = struct{}{}
			declaredNames = append(declaredNames, string(e.Kind)+" "+e.TypeName)
			if e.Kind == kind && e.TypeName == resourceName {
				entityDeclared = true
			}
		}
		if len(declaredFactories) > 1 {
			log.Errorf("the implementation %s declares the factories of several framework entities (%s); their schemas share one file and cannot be isolated, so the docs of %s cannot be compared", sourcePath, strings.Join(declaredNames, ", "), resourceName)
			report.schemaUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		if len(declaredFactories) == 1 && !entityDeclared {
			log.Errorf("the only framework factory declared in %s belongs to other entities (%s), not to %s (%s); the file's schema cannot stand in for this entity, so its docs cannot be compared", sourcePath, strings.Join(declaredNames, ", "), resourceName, kind)
			report.schemaUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		parsed := frameworkschema.Parse(string(source))
		if len(parsed) == 0 {
			log.Errorf("%s declares no framework schema attribute in %s: the declaration is absent or not a static map, so its docs cannot be compared", resourceName, sourcePath)
			report.schemaUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		schemaAttributes := frameworkSchemaAttributes(parsed)
		docsAttributes := make(map[string]ResourceAttribute)
		if duplicates, err := parseResourceDocs(resourceName, frameworkDocsDir(kind), kind == splitlayout.KindResource, docsAttributes, frameworkSchemaKeySet(schemaAttributes)); err != nil {
			log.Errorf("parsing the resource %s docs failed: %s", resourceName, err)
			report.docsUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		} else {
			report.duplicates += len(duplicates)
		}
		finish(resourceName, consistencyCheckAttributes(resourceName, docsAttributes, schemaAttributes))
	}
	// checkName deduplicates and runs the consistency pipeline for one
	// resource or data source touched by the diff. The dedup key is
	// (kind, name): a same-named resource and data source are two entities and
	// each must be compared once with its own schema and docs page, while the
	// implementation and docs page of one entity still collapse into a single
	// check. sourcePath carries the implementation file when the diff includes
	// it; a docs-only diff leaves it empty.
	//
	// The dispatch normalizes the kind before the key is built, because the
	// key must carry the registration's identity, not the diff's claim about
	// it. A list resource's page lives under r/, so its implementation file
	// arrives as (list, name) and its docs page as (resource, name): keying on
	// the diff's own kind would let both through — two keys, one entity — and
	// the entity would be compared twice with every count inflated. After the
	// dispatch both calls carry the registered kind and collapse into one key.
	// A name the dispatch cannot claim keeps the diff's kind. The key also
	// normalizes the flat-layout data_source_ prefix away — the flat
	// implementation file names an SDKv2 data source with it while its docs
	// page names the entity without — so those two arrivals collapse too;
	// framework type names never carry the prefix, so the trim is a no-op
	// there.
	//
	// On the framework path the key is built right after the dispatch. On the
	// SDKv2 path it is built after the map lookup below settles the entity's
	// real class: the historical resource→data-source fallback rewrites that
	// class, so a name registered only as an SDKv2 data source whose r/ and
	// d/ pages are both in the diff would otherwise hold one key per arrival
	// for one entity and be compared twice, whichever order the diff lists
	// the pages in. Keying on the settled class collapses both arrivals into
	// the data source's single comparison.
	checkName := func(kind splitlayout.Kind, resourceName, sourcePath string) {
		fwKind, dispatched := frameworkDispatch(frameworkNames, kind, resourceName)
		if dispatched {
			kind = fwKind
		}
		if dispatched {
			dedupKey := string(kind) + "\x00" + strings.TrimPrefix(resourceName, "data_source_")
			if _, ok := resourceNameMap[dedupKey]; ok {
				return
			}
			resourceNameMap[dedupKey] = struct{}{}
			report.entities++
			checkFrameworkName(kind, resourceName, sourcePath)
			return
		}

		// The kind the diff already identified — from the docs directory, the
		// flat filename prefix or the split registration — picks the SDKv2 map
		// and the docs directory here, the same authority the dispatch and the
		// dedup key below follow. A type name can be registered in both maps
		// (alicloud_actiontrail_global_events_storage_region is one SDKv2
		// resource and one SDKv2 data source): a resource-first probe would let
		// the resource claim the data source's arrival, compare the touched d/
		// page against the resource's schema and the r/ page, and leave the
		// d/ page unchecked. The data-source kind therefore goes straight to
		// DataSourcesMap and d/; the resource kind — and framework-only kinds
		// the dispatch could not claim, which fall through to this path —
		// keeps the historical resource-first lookup with the data-source
		// fallback, so a single-type arrival behaves as it did before.
		isResource := kind != splitlayout.KindDataSource
		docsPath := "website/docs/r/"
		if !isResource {
			docsPath = "website/docs/d/"
		}
		var resource *schema.Resource
		var ok bool
		if isResource {
			resource, ok = sdkProvider.ResourcesMap[resourceName]
			if !ok || resource == nil {
				name := strings.TrimPrefix(resourceName, "data_source_")
				resource, ok = sdkProvider.DataSourcesMap[name]
				if ok && resource != nil {
					resourceName = name
					isResource = false
					docsPath = "website/docs/d/"
				}
			}
		} else {
			name := strings.TrimPrefix(resourceName, "data_source_")
			resource, ok = sdkProvider.DataSourcesMap[name]
			if ok && resource != nil {
				resourceName = name
			}
		}
		// The key carries the class the lookup settled on rather than the
		// arrival's: once the fallback above has rewritten a resource-kind
		// arrival into its data-source identity, the entity's d/ page and the
		// r/ page the fallback redirects both key on (data_source, name) and
		// collapse into one check of one schema against one page.
		effectiveKind := kind
		if !isResource {
			effectiveKind = splitlayout.KindDataSource
		}
		dedupKey := string(effectiveKind) + "\x00" + strings.TrimPrefix(resourceName, "data_source_")
		if _, ok := resourceNameMap[dedupKey]; ok {
			return
		}
		resourceNameMap[dedupKey] = struct{}{}
		report.entities++
		log.Infof("==> Checking resource or data-source %s attributes consistency...", resourceName)
		if !ok || resource == nil {
			log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resourceName)
			report.unresolved++
			report.fail(resourceName)
			exitCode = 1
			return
		}
		resourceSchema := resource.Schema
		resourceSchemaFromDocs := make(map[string]ResourceAttribute)
		// A docs page the parser cannot read or parse leaves the entity entirely
		// unchecked. That used to be logged and then skipped with a green exit,
		// which turned every unreadable page into an invisible coverage hole; a
		// missing page or a malformed list is a documentation defect of its own
		// and must fail the gate.
		if duplicates, err := parseResourceDocs(resourceName, docsPath, isResource, resourceSchemaFromDocs, schemaKeySet(resourceSchema)); err != nil {
			log.Errorf("parsing the resource %s docs failed: %s", resourceName, err)
			report.docsUnreadable++
			report.fail(resourceName)
			exitCode = 1
			return
		} else {
			report.duplicates += len(duplicates)
		}

		finish(resourceName, consistencyCheck(resourceName, resourceSchemaFromDocs, resourceSchema))
	}
	// One split stem resolves once: the implementation and its test share it.
	splitSeen := make(map[string]bool)
	for _, file := range diff.Files {
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			name := strings.TrimPrefix(strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go"), "resource_")
			// The flat-layout filename states the kind through its prefix the
			// same way the docs directory does for a page.
			kind := splitlayout.KindResource
			if strings.HasPrefix(name, "data_source_") {
				kind = splitlayout.KindDataSource
			}
			checkName(kind, name, file.NewName)
		} else if fileDocsRegex.MatchString(file.NewName) {
			if kind, ok := kindFromDocsDir(file.NewName); ok {
				checkName(kind, "alicloud_"+strings.TrimSuffix(strings.Split(file.NewName, "/")[3], ".html.markdown"), "")
			} else {
				// fileDocsRegex is not anchored, so a page deeper than the repo
				// root still matches it while kindFromDocsDir cannot read a
				// kind off its directory. The page is in the diff and the gate
				// has no way to compare it, so the gate must fail rather than
				// warn and skip — an unmapped page was an invisible coverage
				// hole, the same hole an unreadable page used to be.
				log.Errorf("docs page %s matches the docs pattern but lives in no directory this checker maps to a kind", file.NewName)
				report.fail(file.NewName)
				exitCode = 1
			}
		} else if product, stem, _, ok := splitlayout.IsSplitSource(file.NewName); ok {
			// The implementation and its test share the stem; the first of
			// the pair to appear in the diff carries both paths.
			if splitSeen[product+"/"+stem] {
				continue
			}
			splitSeen[product+"/"+stem] = true
			resolved, err := splitlayout.Resolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve split-layout file %s. Error: %s", file.NewName, err)
				continue
			}
			// One file can register more than one resource, so check every
			// matched registration.
			for _, r := range resolved {
				checkName(r.Kind, r.TypeName, r.ImplPath)
			}
			// Framework-native registrations — managed resources, data
			// sources, ephemerals, list resources and actions — live outside
			// the provider's SDKv2 maps, so their schema is parsed from the
			// implementation file rather than from the runtime map. Each
			// registration carries its kind, so a same-named pair splits into
			// two checks here.
			fwResolved, err := splitlayout.FrameworkResolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve framework registrations of %s. Error: %s", file.NewName, err)
				continue
			}
			for _, r := range fwResolved {
				checkName(r.Kind, r.TypeName, r.ImplPath)
			}
			if len(resolved) == 0 && len(fwResolved) == 0 {
				log.Warnf("==> %s registers no resource, data source, ephemeral, action or list resource", file.NewName)
			}
		}
	}
	// A removed file keeps only OrigName: diffparser leaves NewName empty for a
	// DELETED file, so the loop above never sees the deletion. Without this pass a
	// docs page can disappear while the entity it documents stays registered, and
	// the gate used to stay green. Only the zero NewName is trusted as a deletion
	// marker: DELETED is iota zero, so the mode alone cannot distinguish an unset
	// field from a real removal.
	//
	// The boundary is deliberate: entity removal in the split layout is expressed
	// by dropping a registration entry, not by deleting a file, so a registration
	// removal that leaves its page behind is not visible here. This pass covers
	// the direction that does surface in a diff — the page going away first.
	//
	// registered is deliberately name-level across all kinds: removing a page
	// must be blocked by ANY entity of that name still being registered, which
	// is what the name-only check always did.
	registered := func(entity string) bool {
		for _, names := range frameworkNames {
			if _, ok := names[entity]; ok {
				return true
			}
		}
		if resource, ok := sdkProvider.ResourcesMap[entity]; ok && resource != nil {
			return true
		}
		if dataSource, ok := sdkProvider.DataSourcesMap[entity]; ok && dataSource != nil {
			return true
		}
		return false
	}
	deletedPages := make(map[string]struct{})
	for _, file := range diff.Files {
		if file.Mode != diffparser.DELETED || file.NewName != "" {
			continue
		}
		if !fileDocsRegex.MatchString(file.OrigName) {
			continue
		}
		entity := entityFromDocsPage(file.OrigName)
		if _, ok := deletedPages[entity]; ok {
			continue
		}
		deletedPages[entity] = struct{}{}
		if !registered(entity) {
			log.Infof("==> %s and its docs page %s are both removed, nothing to compare", entity, file.OrigName)
			continue
		}
		log.Errorf("docs page %s is deleted but %s is still registered", file.OrigName, entity)
		report.fail(entity)
		exitCode = 1
	}
	return report, exitCode
}

// frameworkRegisteredNames returns the framework-native registrations keyed by
// the kind of entity they name. The kind is part of the identity: a product
// can register one type name as a managed resource and as a data source
// (alicloud_sae_load_balancer_nlb), and the name-only set the checker used
// before collapsed that pair into one entity whose single comparison could
// pick the wrong schema and docs page. The runtime registry carries the kind
// through the field each registration is declared in.
func frameworkRegisteredNames() map[splitlayout.Kind]map[string]struct{} {
	names := make(map[splitlayout.Kind]map[string]struct{})
	for _, kind := range []splitlayout.Kind{
		splitlayout.KindResource,
		splitlayout.KindDataSource,
		splitlayout.KindEphemeral,
		splitlayout.KindListResource,
		splitlayout.KindAction,
	} {
		names[kind] = make(map[string]struct{})
	}
	for _, sp := range providerregistry.ServicePackages() {
		for _, d := range sp.Resources {
			names[splitlayout.KindResource][d.TypeName] = struct{}{}
		}
		for _, d := range sp.DataSources {
			names[splitlayout.KindDataSource][d.TypeName] = struct{}{}
		}
		// Ephemerals, list resources and actions are registered entities
		// too: their names must be recognized, not reported as missing from
		// the SDKv2 provider maps.
		for _, d := range sp.EphemeralResources {
			names[splitlayout.KindEphemeral][d.TypeName] = struct{}{}
		}
		for _, d := range sp.ListResources {
			names[splitlayout.KindListResource][d.TypeName] = struct{}{}
		}
		for _, d := range sp.Actions {
			names[splitlayout.KindAction][d.TypeName] = struct{}{}
		}
	}
	return names
}

// entityFromDocsPage derives the entity a docs page documents. The page name is
// the only link the gate has from a removed markdown file back to a resource or
// data source.
func entityFromDocsPage(path string) string {
	return "alicloud_" + strings.TrimSuffix(filepath.Base(path), ".html.markdown")
}

func schemaKeySet(resourceSchema map[string]*schema.Schema) map[string]struct{} {
	flattened := make(map[string]ResourceAttribute)
	getResourceAttributes("", flattened, resourceSchema, "")
	keys := make(map[string]struct{}, len(flattened))
	for key := range flattened {
		keys[key] = struct{}{}
	}
	return keys
}

// headingOwnsScope reports whether a `### `name“ heading names a schema
// attribute (or an attribute path prefix) and therefore prefixes the bullets
// below it. Documentation-only grouping headings — a per-protocol section on a
// resource whose schema is flat, for example — must leave the prefix empty,
// otherwise every field under the heading is looked up as `group.field` and
// reported missing. A nil key set means the schema is unknown (the parser is
// being unit tested in isolation), where the historical assumption that every
// documented block heading is an attribute is kept.
func headingOwnsScope(name string, schemaKeys map[string]struct{}) bool {
	if schemaKeys == nil {
		return true
	}
	if name == "" {
		return false
	}
	if _, ok := schemaKeys[name]; ok {
		return true
	}
	prefix := name + "."
	for key := range schemaKeys {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// attributeNameRe is the attribute-name contract shared by both reference
// sections. A bullet whose backticked text carries anything else — a space, a
// slash, a dot — is a documentation grouping/see-also line, not an attribute
// entry, and must stay out of the docs map in BOTH comparison directions.
var attributeNameRe = regexp.MustCompile(`^[a-zA-Z_0-9]+$`)

// parseResourceDocs reads one docs page into resourceAttributes. It returns the
// attributes documented in both reference sections: those are reported as
// warnings, not failures, because a repo-wide census on the current tree finds a
// handful of long-standing pages that list an input again as an output, and
// failing the gate on them would block every unrelated change to those pages.
// Returning them keeps the warning out of the silent "logged an error, exited 0"
// bucket: the summary carries the count.
func parseResourceDocs(resourceName, docsPath string, isResource bool, resourceAttributes map[string]ResourceAttribute, schemaKeys map[string]struct{}) ([]string, error) {
	splitRes := strings.Split(resourceName, "alicloud_")
	if len(splitRes) < 2 {
		log.Errorf("parsing resource name %s failed.", resourceName)
		return nil, fmt.Errorf("parsing resource name %s failed", resourceName)
	}

	filePath := strings.Join([]string{docsPath, splitRes[1], ".html.markdown"}, "")

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open resource %s docs failed. Error: %s", filePath, err)
		return nil, err
	}
	defer file.Close()

	argsRegex := regexp.MustCompile("## Argument Reference")
	attribRegex := regexp.MustCompile("## Attributes Reference")
	secondLevelRegex := regexp.MustCompile("^### `([a-zA-Z_\\.0-9-]*)`")
	// The name capture is deliberately wider than the attribute-name contract so
	// that a see-also/grouping bullet is recognized here and rejected by
	// parseMatchLine, instead of silently falling through both regexes.
	argumentsFieldRegex := regexp.MustCompile("^\\* `([^`]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")
	attributeFieldRegex := regexp.MustCompile("^\\s*\\* `([^`]*)`[ ]*-?(.*)")

	name := filepath.Base(filePath)
	re := regexp.MustCompile("[a-z0-9A-Z_]*")
	resourceName = "alicloud_" + re.FindString(name)

	scanner := bufio.NewScanner(file)
	phase := "Argument"
	record := false
	subAttributeName := ""
	line := 0
	rootPrefixLen := 0
	rootName := ""
	// indentStep is the per-list indentation increment detected from the
	// first descent in the Attribute phase. The project's markdownlint
	// (MD007) mandates a two-space step, while legacy docs use four; both
	// must be parsed. The step is detected dynamically so the back-up
	// arithmetic is not hard-coded to a single indent width.
	indentStep := 0
	extraResourceAttribute := map[string]struct{}{}
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		if strings.HasPrefix(text, "#") && (strings.HasSuffix(text, "Timeouts") || strings.HasSuffix(text, "Import")) {
			break
		}
		if argsRegex.MatchString(text) {
			record = true
			phase = "Argument"
			continue
		}
		if attribRegex.MatchString(text) {
			record = true
			phase = "Attribute"
			subAttributeName = ""
			continue
		}
		if secondLevelRegex.MatchString(strings.TrimSpace(text)) {
			record = true
			parts := strings.Split(strings.TrimSpace(text), " ")
			block := strings.Replace(strings.Trim(parts[len(parts)-1], "`"), "-", ".", -1)
			// Only a heading that names a schema attribute opens a nested path
			// scope. This is what fixes the false positives: a block heading that is
			// not in the schema is a documentation grouping and must not prefix the
			// bullets below it.
			if headingOwnsScope(block, schemaKeys) {
				subAttributeName = block
			} else {
				subAttributeName = ""
			}
			// A `### `block`` heading opens a fresh attribute scope. The
			// indentation context of the previous list (depth, current field
			// and detected step) must not bleed into the new scope, otherwise
			// a block following a nested list would mis-parse its first field.
			rootPrefixLen = 0
			rootName = ""
			indentStep = 0
			continue
		}
		if record {
			var matched [][]string
			if phase == "Argument" {
				matched = argumentsFieldRegex.FindAllStringSubmatch(text, 1)
			} else if phase == "Attribute" {
				matched = attributeFieldRegex.FindAllStringSubmatch(text, 1)
			}

			for _, m := range matched {
				if phase == "Attribute" {
					thisLen := len(strings.Split(m[0], "*")[0])

					if rootPrefixLen < thisLen {
						// Descending one indentation level. The step is
						// established by the first descent and every
						// subsequent descent must reuse the same step; a
						// different delta means the list is inconsistently
						// indented (e.g. a two-space list that suddenly widens
						// by four, or a skip-level child), which this gate
						// rejects.
						diff := thisLen - rootPrefixLen
						if indentStep == 0 {
							// With a "* " marker, a child list starts two to
							// five columns after its parent. One column is
							// still a sibling; six or more is paragraph/code
							// content, not a nested attribute list. MD007
							// separately enforces the project's indent style.
							if diff < 2 || diff > 5 {
								return nil, fmt.Errorf("resource %s docs %s: indentation delta %d is not a nested attribute list", resourceName, docsPath, diff)
							}
							indentStep = diff
						} else if diff != indentStep {
							return nil, fmt.Errorf("resource %s docs %s have not been formatted: inconsistent indentation step", resourceName, docsPath)
						}
						if subAttributeName != "" {
							subAttributeName += "." + rootName
						} else {
							subAttributeName = rootName
						}
						rootName = m[1]
					} else if rootPrefixLen > thisLen {
						parts := strings.Split(subAttributeName, ".")
						backIndex := rootPrefixLen - thisLen
						// Ascending back toward a sibling or ancestor. The
						// indent delta must be a whole multiple of the
						// detected step; otherwise the list is malformed. The
						// number of levels to pop is delta / step, which works
						// for both two-space and four-space indentation.
						if indentStep == 0 || backIndex%indentStep != 0 {
							return nil, fmt.Errorf("resource %s docs %s have not been formatted: indent step %d does not divide indent delta %d", resourceName, docsPath, indentStep, backIndex)
						}
						backIndex /= indentStep
						if len(parts) > 0 {
							for backIndex > 0 {
								backIndex--
								if strings.Contains(subAttributeName, ".") {
									subAttributeName = strings.TrimSuffix(strings.TrimSuffix(subAttributeName, parts[len(parts)-1]), ".")
									parts = parts[:len(parts)-1]
								} else {
									subAttributeName = ""
								}
								rootName = m[1]
							}
						}
					} else {
						rootName = m[1]
					}
					rootPrefixLen = thisLen
				}
				attribute := parseMatchLine(m, phase, subAttributeName)
				if attribute == nil {
					continue
				}
				attribute.DocsLineNum = line
				// A repeat is only the Argument/Attributes conflict this check
				// exists for when it happens in the Attributes Reference. A field
				// may legitimately be listed once per grouping section inside the
				// Argument Reference — a per-protocol block on a resource whose
				// schema is flat repeats frontend_port, health_check and friends
				// in every protocol section — and that is not an output being
				// re-documented as a configuration argument.
				if _, ok := resourceAttributes[attribute.Name]; !ok {
					resourceAttributes[attribute.Name] = *attribute
				} else if isResource && phase == "Attribute" {
					extraResourceAttribute[attribute.Name] = struct{}{}
				}
				if phase == "Attribute" {
					for key := range extraResourceAttribute {
						if strings.HasPrefix(attribute.Name, key+".") {
							delete(extraResourceAttribute, key)
						}
					}
				}
			}
		}
	}

	// A duplicate is counted, not failed: the run must stay green so an unrelated
	// edit to one of the long-standing pages still passes. Logging it at error
	// level contradicted that — the job exited 0 while the log said ERRO — so it
	// is a warning, matching the summary line below.
	for key := range extraResourceAttribute {
		log.Warnf("'%v' has been set in the `## Argument Reference` and it should be removed from `## Attributes Reference`", key)
	}
	duplicates := make([]string, 0, len(extraResourceAttribute))
	for key := range extraResourceAttribute {
		duplicates = append(duplicates, key)
	}
	sort.Strings(duplicates)

	return duplicates, nil
}

func parseMatchLine(words []string, phase, rootName string) *ResourceAttribute {
	// A backticked bullet that does not name an attribute — "HTTP Listener",
	// "foo.bar", a path — is a documentation grouping or see-also line. Keeping
	// it out of the docs map makes both comparison directions agree: it is not
	// reported as a missing schema field, and it cannot mask a real attribute.
	if len(words) < 2 || !attributeNameRe.MatchString(words[1]) {
		return nil
	}
	result := ResourceAttribute{}
	if phase == "Argument" && len(words) >= 4 {
		if rootName != "" {
			result.Name = rootName + "." + words[1]
		} else {
			result.Name = words[1]
		}
		//result["Description"] = words[3]
		if strings.Contains(words[2], "Optional") {
			result.Optional = "true"
		}
		if strings.Contains(words[2], "Required") {
			result.Required = "true"
		}
		if strings.Contains(words[2], "ForceNew") {
			result.ForceNew = true
		}
		if strings.Contains(words[2], "Deprecated") {
			result.Deprecated = "Deprecated since"
		}
		return &result
	}
	if phase == "Attribute" && len(words) >= 3 {
		if rootName != "" {
			result.Name = rootName + "." + words[1]
		} else {
			result.Name = words[1]
		}
		// Attribute-phase docs embed the modifier inside a CLOSED parenthesized
		// modifier group (e.g. "- (Optional) The tag ..."). attributeFieldRegex
		// captures the whole description tail in words[2], unlike
		// argumentsFieldRegex whose second capture group is the modifier group
		// itself. A plain Contains on words[2] would therefore treat free-form
		// prose like "Optional metadata returned by the service." or "Required
		// for subsequent API calls." as modifiers, and an unclosed "(Optional"
		// or a prose parenthesis like "(Required for subsequent API calls)" /
		// "(Optional metadata returned by the service)" would also be wrong.
		// Match only a closed modifier group whose members are modifier tokens;
		// Deprecated keeps its legacy Contains behavior.
		for _, group := range attrModifierGroupRe.FindAllString(words[2], -1) {
			if strings.Contains(group, "Optional") {
				result.Optional = "true"
			}
			if strings.Contains(group, "Required") {
				result.Required = "true"
			}
			if strings.Contains(group, "ForceNew") {
				result.ForceNew = true
			}
		}
		if strings.Contains(words[2], "Deprecated") {
			result.Deprecated = "Deprecated since"
		}
		//result["Description"] = words[2]
		return &result
	}
	return nil
}

func consistencyCheck(resourceName string, resourceAttributeFromDocs map[string]ResourceAttribute, resourceSchemaDefined map[string]*schema.Schema) bool {
	// the number of the schema field + 1(id) should equal to the number defined in document
	resourceAttributes := make(map[string]ResourceAttribute)
	getResourceAttributes("", resourceAttributes, resourceSchemaDefined, "")
	return consistencyCheckAttributes(resourceName, resourceAttributeFromDocs, resourceAttributes)
}

// consistencyCheckAttributes compares one flattened schema against the docs.
// It takes the flattened form rather than *schema.Schema so that the framework
// dialect — whose schema never enters the SDKv2 runtime maps — runs this same
// comparison instead of a parallel reimplementation.
func consistencyCheckAttributes(resourceName string, resourceAttributeFromDocs, resourceAttributes map[string]ResourceAttribute) bool {
	isConsistent := true
	filteredList := set.NewSet()
	if val, ok := filterList[resourceName]; ok {
		for _, v := range val {
			filteredList.Add(v)
		}
	}

	for attributeKey, attributeValue := range resourceAttributes {
		if _, ok := skippedSchemaKeys[attributeKey]; ok {
			continue
		}
		// `id` is implicit in both dialects and optional in the docs: the
		// reverse loop below tolerates a documented id that no schema declares,
		// and this direction must tolerate an undeclared-in-docs id for the same
		// reason. Only a framework entity ever reaches here with one, because an
		// SDKv2 schema never declares id.
		if attributeKey == "id" {
			continue
		}
		// The `## Timeouts` block is a standard docs section, not an attribute
		// reference, and the parser reads only the two reference sections. An
		// SDKv2 resource keeps timeouts outside its Schema map so it never
		// reached this loop; a framework resource declares them as a plain
		// attribute, and comparing that against a section the parser does not
		// read would report every such resource as undocumented.
		if attributeKey == "timeouts" || strings.HasPrefix(attributeKey, "timeouts.") {
			continue
		}
		attributeDocsValue, ok := resourceAttributeFromDocs[attributeKey]
		if attributeValue.Removed != "" {
			continue
		}
		if !ok {
			isConsistent = false
			log.Errorf("'%v' is not found in the docs", attributeKey)
		}
		if attributeValue.Deprecated != "" {
			if attributeDocsValue.Deprecated == "" {
				isConsistent = false
				log.Errorf("'%v' should be marked as Deprecated in the document description", attributeKey)
			}
			continue
		}
		if attributeValue.Optional == "true" && attributeDocsValue.Optional != attributeValue.Optional {
			isConsistent = false
			log.Errorf("'%v' should be marked as Optional in the document", attributeKey)
		}
		if attributeValue.Required == "true" && attributeDocsValue.Required != attributeValue.Required {
			isConsistent = false
			log.Errorf("'%v' should be marked as Required in the document", attributeKey)
		}
		if attributeValue.ForceNew && !attributeDocsValue.ForceNew {
			isConsistent = false
			log.Errorf("'%v' should be marked as ForceNew in the document description", attributeKey)
		}
		// Reverse consistency: the forward branches above only verify that a
		// schema-declared Optional/Required attribute is documented as such.
		// They stay silent in the opposite direction — a Computed-only
		// attribute (neither Optional nor Required in the schema) advertised
		// as configurable in the docs slipped through, e.g. ssl_vpn_client_
		// cert's status/ca_cert/client_cert/client_key/client_config, whose
		// Argument Reference lines wrongly carried "(Optional)". Reject that
		// direction too, so the docs cannot grant configuration capabilities
		// the schema does not have.
		if attributeValue.Optional != "true" && attributeValue.Required != "true" {
			if attributeDocsValue.Optional == "true" {
				isConsistent = false
				log.Errorf("'%v' is not configurable in the schema and should not be marked as Optional in the document", attributeKey)
			}
			if attributeDocsValue.Required == "true" {
				isConsistent = false
				log.Errorf("'%v' is not configurable in the schema and should not be marked as Required in the document", attributeKey)
			}
		}
	}
	for attributeKey, _ := range resourceAttributeFromDocs {
		if _, ok := resourceAttributes[attributeKey]; !ok && attributeKey != "id" {
			isConsistent = false
			log.Errorf("'%v' which described in the docs not found in the resource schema", attributeKey)
		}
	}
	return isConsistent
}

func getResourceAttributes(rootName string, resourceAttributeMap map[string]ResourceAttribute, resourceSchema map[string]*schema.Schema, rootRemoved string) {
	for key, value := range resourceSchema {
		if rootName != "" {
			key = rootName + "." + key
		}

		var thisRemoved = ""
		if len(thisRemoved) == 0 && len(rootRemoved) != 0 {
			thisRemoved = rootRemoved
		}

		if _, ok := resourceAttributeMap[key]; !ok {
			resourceAttributeMap[key] = ResourceAttribute{
				Name:       key,
				Type:       value.Type.String(),
				Optional:   fmt.Sprint(value.Optional),
				Required:   fmt.Sprint(value.Required),
				ForceNew:   value.ForceNew,
				Default:    fmt.Sprint(value.Default),
				Deprecated: value.Deprecated,
				Removed:    thisRemoved,
			}
		}
		if value.Type == schema.TypeSet || value.Type == schema.TypeList {
			if v, ok := value.Elem.(*schema.Schema); ok {
				vv := resourceAttributeMap[key]
				vv.ElemType = v.Type.String()
				resourceAttributeMap[key] = vv
			} else {
				vv := resourceAttributeMap[key]
				vv.ElemType = "Object"
				resourceAttributeMap[key] = vv
				getResourceAttributes(key, resourceAttributeMap, value.Elem.(*schema.Resource).Schema, thisRemoved)
			}
		}
		if value.Type == schema.TypeMap {
			if _, ok := value.Elem.(*schema.Resource); ok {
				vv := resourceAttributeMap[key]
				vv.ElemType = "Object"
				resourceAttributeMap[key] = vv
				getResourceAttributes(key, resourceAttributeMap, value.Elem.(*schema.Resource).Schema, thisRemoved)

			}
		}
	}
}

// frameworkSchemaAttributes flattens a parsed framework declaration into the
// same shape getResourceAttributes produces from an SDKv2 schema. Optional and
// Required default to "false" because consistencyCheckAttributes compares
// their string form against what the docs parser extracted, and an unset
// framework field means neither.
func frameworkSchemaAttributes(parsed map[string]map[string]interface{}) map[string]ResourceAttribute {
	attributes := make(map[string]ResourceAttribute, len(parsed))
	for path, props := range parsed {
		attribute := ResourceAttribute{Name: path, Optional: "false", Required: "false"}
		if value, ok := props["Type"].(string); ok {
			attribute.Type = value
		}
		if value, ok := props["Optional"].(bool); ok {
			attribute.Optional = fmt.Sprint(value)
		}
		if value, ok := props["Required"].(bool); ok {
			attribute.Required = fmt.Sprint(value)
		}
		if value, ok := props["ForceNew"].(bool); ok {
			attribute.ForceNew = value
		}
		if value, ok := props["Deprecated"].(string); ok {
			attribute.Deprecated = value
		}
		attributes[path] = attribute
	}
	return attributes
}

// frameworkSchemaKeySet mirrors schemaKeySet for a flattened Framework schema.
// Keep every path so headingOwnsScope can recognize nested block headings.
func frameworkSchemaKeySet(attributes map[string]ResourceAttribute) map[string]struct{} {
	keys := make(map[string]struct{}, len(attributes))
	for name := range attributes {
		keys[name] = struct{}{}
	}
	return keys
}
