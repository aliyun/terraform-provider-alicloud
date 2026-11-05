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
// A whole-product migration can name too many entities for one title —
// the per-entity fragments ran the largest packages (RAM, ResourceManager,
// CR, CloudMonitorService) past the SCM title length and were rejected
// with SCM_PR_CONTENT_INVALID. The concise form migrate/<product>:
// waives the per-entity fragments for one product's existing migration
// only: every changed service resource or data source must belong to that
// single product, every one must already exist at the base (a migration,
// never a new entity), and the marker's product must match. New entities,
// flat entities and a cross-product mix keep demanding their declarations,
// so the marker is a scoped shortcut, not a "migrate" bypass. The base
// index also reads the root alicloud/provider.go so a public alias
// registered there (Quotas #11288) is announced as modified when it moves
// into a service register, not as new.
//
// The base index is kind-aware: a Resource and a Data Source may share a
// type name — it is legal and common in the Provider — so the index keeps
// the two kinds apart. A migration that moves a root DataSource alias into
// a service Resource register, or vice versa, must announce the new kind,
// not waive it behind migrate/<product>:. The root alicloud/provider.go
// DataSourcesMap yields data sources and ResourcesMap yields resources, so
// ProviderRegisteredNames returns the kind each map belongs to.
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
	"regexp"
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
	// product is the split-layout package key a service entity belongs to
	// ("" for flat entities). The whole-product migration aggregate is
	// scoped to one product, so the title gate needs each entity's product
	// to refuse a cross-product change hiding behind a single marker.
	product string
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
	if _, split := splitlayout.ServiceProduct(name); split && strings.HasSuffix(name, ".go") {
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
// exist. The index is kind-aware, so a base DataSource does not satisfy a
// same-named service Resource: the two kinds are kept apart.
func serviceEntities(name string, index *baseIndex) []entity {
	product, _, _, ok := splitlayout.IsSplitSource(name)
	if !ok {
		return nil
	}
	var entities []entity
	for _, r := range joinRegistrations(name) {
		switch r.Kind {
		case splitlayout.KindResource:
			entities = append(entities, entity{typeName: r.TypeName, existing: index.has(r.TypeName, false), product: product})
		case splitlayout.KindDataSource:
			entities = append(entities, entity{typeName: r.TypeName, dataSource: true, existing: index.has(r.TypeName, true), product: product})
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

// baseIndex is the set of provider entities that existed at the base
// commit, in either layout: the flat source files directly, the
// registrations of every product register.go, and the root
// alicloud/provider.go public aliases. It decides whether a service entity
// is announced as new or as existing. The index is kind-aware: resources
// and data sources are kept in separate maps, because a Resource and a
// Data Source may share a type name, and a same-name collision must not
// let one kind waive the other's declaration. Ephemeral, list and action
// registrations have no house title form, so they never satisfy a resource
// or data-source lookup; they are retained in a separate others map rather
// than folded into resources, so a base ephemeral (or list/action) does
// not let a same-named head Resource be waived behind migrate/<product>:.
type baseIndex struct {
	resources   map[string]bool
	dataSources map[string]bool
	others      map[string]bool
}

// has reports whether typeName existed at the indexed base commit as the
// given kind: dataSource consults the data-source map, otherwise the
// resource map. A nil index or a missing kind reports not existing.
func (b *baseIndex) has(typeName string, dataSource bool) bool {
	if b == nil {
		return false
	}
	if dataSource {
		return b.dataSources[typeName]
	}
	return b.resources[typeName]
}

// loadBaseIndex builds the entity index of one base commit from git: the
// flat sources directly under alicloud/, every register.go under
// alicloud/service/ and alicloud/provider/core, and the root
// alicloud/provider.go whose DataSourcesMap and ResourcesMap register the
// public aliases the flat layout keeps at the provider root. A base without
// the split layout lists no products and simply carries the flat and root
// names.
func loadBaseIndex(base string) (*baseIndex, error) {
	flat, err := gitLines("ls-tree", "--name-only", base, "--", "alicloud/")
	if err != nil {
		return nil, err
	}
	sources := make(map[string][]byte)
	if serviceFiles, err := gitLines("ls-tree", "-r", "--name-only", base, "--", "alicloud/service", "alicloud/provider/core"); err != nil {
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
	// The root provider.go registers aliases (Quotas #11288) that no flat
	// filename stem and no service register.go carries; without them a
	// migration that moves such an alias into a service register is
	// misannounced as new. A base predating provider.go simply has none.
	providerSrc, err := exec.Command("git", "show", base+":alicloud/provider.go").Output()
	if err != nil {
		providerSrc = nil
	}
	return baseIndexFromSources(flat, sources, providerSrc), nil
}

// baseIndexFromSources builds the entity index from the flat source paths,
// the register.go sources of one tree and the root provider.go source; the
// git plumbing in loadBaseIndex only gathers those. Flat stems, register
// entries and root aliases are all filed by kind, so a same-name
// Resource/DataSource collision never lets one kind satisfy the other.
func baseIndexFromSources(flatFiles []string, registers map[string][]byte, providerSrc []byte) *baseIndex {
	index := &baseIndex{resources: make(map[string]bool), dataSources: make(map[string]bool), others: make(map[string]bool)}
	for _, file := range flatFiles {
		if !strings.HasSuffix(file, ".go") || strings.HasSuffix(file, "_test.go") {
			continue
		}
		stem := strings.TrimSuffix(file[strings.LastIndex(file, "/")+1:], ".go")
		// Only a real flat implementation file — resource_alicloud_* or
		// data_source_alicloud_* — names a provider entity; the _test.go
		// files were skipped above, and a helper or provider.go stem with no
		// resource_/data_source_ prefix names no entity. A type name that
		// legitimately ends in _test is still proven by a real register.go
		// or root provider.go registration, not by a test-file stem.
		var kind splitlayout.Kind
		var typeName string
		switch {
		case strings.HasPrefix(stem, "resource_"):
			kind = splitlayout.KindResource
			typeName = strings.TrimPrefix(stem, "resource_")
		case strings.HasPrefix(stem, "data_source_"):
			kind = splitlayout.KindDataSource
			typeName = strings.TrimPrefix(stem, "data_source_")
		default:
			continue
		}
		index.add(typeName, kind)
	}
	for _, src := range registers {
		entries := append(splitlayout.SDKRegisterEntries(src), splitlayout.FrameworkRegisterEntries(src)...)
		for _, e := range entries {
			index.add(e.TypeName, e.Kind)
		}
	}
	for _, e := range splitlayout.ProviderRegisteredNames(providerSrc) {
		index.add(e.TypeName, e.Kind)
	}
	return index
}

// add records typeName at the base commit as the given kind. A Data Source
// goes into the data-source map and a Resource into the resource map; any
// other kind (ephemeral, action, list) goes into the separate others map,
// so it is retained without satisfying a resource or data-source lookup.
// The title gate only resolves Resource and DataSource entities, so those
// kinds never reach a title demand, but folding them into resources would
// let a base ephemeral satisfy a same-named head Resource and be waived
// behind migrate/<product>: — keeping them apart closes that gap.
func (b *baseIndex) add(typeName string, kind splitlayout.Kind) {
	switch kind {
	case splitlayout.KindDataSource:
		b.dataSources[typeName] = true
	case splitlayout.KindResource:
		b.resources[typeName] = true
	default:
		b.others[typeName] = true
	}
}

// aggregateMarkerRe matches the concise whole-product migration form
// migrate/<product>:. The product segment mirrors the split-layout package
// key charset, so a path-like description or the bare word "migrate" never
// parses as a marker, and resource/<type>:/data-source/<type>:/ fragments
// stay distinct from it.
var aggregateMarkerRe = regexp.MustCompile(`migrate/([0-9A-Za-z_-]+):`)

// aggregate reports whether the title's migrate/<product>: marker waives the
// per-entity resource/<type>:/data-source/<type>:/ fragments for one
// product's existing migration. The waiver applies only when every changed
// service resource or data source belongs to a single product, every one of
// them already existed at the base (a migration, never a new entity) and the
// product the marker names is that same product. A title that merely
// contains the word "migrate", or a change that mixes products, adds new
// entities or names the wrong product, gets no waiver: the per-entity rules
// keep applying, and a title that leaned on the waiver then fails. New
// entities and flat entities are never waived, so the marker cannot hide an
// undeclared addition or a flat alias left in the root layout.
type aggregate struct {
	product string
	ok      bool
}

// covers reports whether the aggregate waives the modified-form requirement
// for one existing service entity.
func (a aggregate) covers(e entity) bool {
	return a.ok && e.product != "" && e.product == a.product
}

// parseAggregate decides whether the title's migrate/<product>: marker
// waives the per-entity fragments for one product's existing migration. It
// is a pure waiver decision: when the change does not qualify it simply
// returns ok=false, and the per-entity loop enforces the ordinary rules —
// a new entity still demands its New Resource/New Data Source form, an
// existing entity still demands its resource/<type>:/data-source/<type>:/
// fragment. So a cross-product mix, a new entity or a wrong product is
// identifiable through the very declarations the waiver would have hidden,
// and a title that declares every entity beside the marker still passes.
func parseAggregate(prTitle string, files []titleFile) aggregate {
	product := aggregateMarker(prTitle)
	if product == "" {
		return aggregate{}
	}
	var products []string
	newEntity := false
	for _, file := range files {
		for _, e := range file.entities {
			if e.product == "" {
				continue // flat entities are not part of a product migration
			}
			if !e.existing {
				newEntity = true
			}
			products = append(products, e.product)
		}
	}
	if len(products) == 0 {
		// No service entity to migrate: the marker is ordinary text, and
		// the docs-only, testcase-only and per-entity rules still judge it.
		return aggregate{product: product}
	}
	first := products[0]
	single := true
	for _, p := range products {
		if p != first {
			single = false
			break
		}
	}
	if single && !newEntity && first == product {
		return aggregate{product: product, ok: true}
	}
	return aggregate{product: product}
}

// aggregateMarker returns the product named by a migrate/<product>: marker,
// or "" when the title uses no such marker.
func aggregateMarker(prTitle string) string {
	m := aggregateMarkerRe.FindStringSubmatch(prTitle)
	if m == nil {
		return ""
	}
	return m[1]
}

// checkTitle enforces the title conventions: every registered entity a
// changed source file implements must be announced in the title — the
// new-entity form for entities the base did not know, the modified form
// otherwise, unless the concise whole-product migration form
// migrate/<product>: waives the per-entity fragments for one product's
// existing migration — a docs-only change must start with the docs marker, a
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

	agg := parseAggregate(prTitle, files)

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
			if agg.covers(e) {
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
