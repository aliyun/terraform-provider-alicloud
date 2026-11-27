//nolint:all
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/waigani/diffparser"
)

// The same-named framework pair the gate mishandled: one type name registered
// as a managed resource and as a data source, each with its own factory,
// implementation file, schema and docs page.
const sameNameEntity = "alicloud_sae_load_balancer_nlb"

const sameNameRegisterResourceFirst = `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
	}
}
`

const sameNameRegisterDataSourceFirst = `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBResource,
			},
		},
	}
}
`

// The two implementations declare different schemas: passing the gate means
// each entity was compared against its own schema and its own docs page, not
// one schema against both pages.
const sameNameResourceImplSrc = `package sae

func NewLoadBalancerNLBResource() *nlbResource {
	return &nlbResource{}
}

func nlbResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"listener_port":   schema.Int64Attribute{Required: true},
		"server_group_id": schema.StringAttribute{Optional: true},
	}
}
`

const sameNameDataSourceImplSrc = `package sae

func NewLoadBalancerNLBDataSource() *nlbDataSource {
	return &nlbDataSource{}
}

func nlbDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"load_balancer_id": schema.StringAttribute{Optional: true},
		"output_file":      schema.StringAttribute{Optional: true},
	}
}
`

// sameNameSharedImplSrc declares both factories of the pair in one file: the
// layout the package tolerates but the registration can only resolve
// file-grained, so both (kind, name) pairs resolve to this same ImplPath.
const sameNameSharedImplSrc = `package sae

func NewLoadBalancerNLBResource() *nlbResource {
	return &nlbResource{}
}

func nlbResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"listener_port":   schema.Int64Attribute{Required: true},
		"server_group_id": schema.StringAttribute{Optional: true},
	}
}

func NewLoadBalancerNLBDataSource() *nlbDataSource {
	return &nlbDataSource{}
}

func nlbDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"load_balancer_id": schema.StringAttribute{Optional: true},
		"output_file":      schema.StringAttribute{Optional: true},
	}
}
`

const sameNameDocsR = `# alicloud_sae_load_balancer_nlb

Provides an SAE NLB resource.

## Argument Reference

The following arguments are supported:

* ` + "`listener_port`" + ` - (Required) The port of the listener.
* ` + "`server_group_id`" + ` - (Optional) The id of the server group.
`

// sameNameDocsRMissingField drops one attribute bullet: a real documentation
// defect of the resource only.
const sameNameDocsRMissingField = `# alicloud_sae_load_balancer_nlb

Provides an SAE NLB resource.

## Argument Reference

The following arguments are supported:

* ` + "`listener_port`" + ` - (Required) The port of the listener.
`

const sameNameDocsD = `# alicloud_sae_load_balancer_nlb

This data source provides the SAE NLB.

## Argument Reference

* ` + "`load_balancer_id`" + ` - (Optional) The id of the load balancer.
* ` + "`output_file`" + ` - (Optional) The output file path.
`

const (
	sameNameDocsRPath = "website/docs/r/sae_load_balancer_nlb.html.markdown"
	sameNameDocsDPath = "website/docs/d/sae_load_balancer_nlb.html.markdown"

	sameNameResourceImplPath = "alicloud/service/sae/resource_load_balancer_nlb.go"
	sameNameDataSourcePath   = "alicloud/service/sae/data_source_load_balancer_nlb.go"
	sameNameSharedImplPath   = "alicloud/service/sae/load_balancer_nlb.go"
)

func sameNameFrameworkNames() map[splitlayout.Kind]map[string]struct{} {
	return map[splitlayout.Kind]map[string]struct{}{
		splitlayout.KindResource:   {sameNameEntity: {}},
		splitlayout.KindDataSource: {sameNameEntity: {}},
	}
}

// writeSameNameConsistencyTree materializes the source and docs tree for the
// same-named pair. An empty docsR or docsD leaves that page out, which is how
// the single-entity defect cases keep the other page intact.
func writeSameNameConsistencyTree(t *testing.T, registerSrc, docsR, docsD string) {
	t.Helper()
	serviceDir := filepath.Join("alicloud", "service", "sae")
	files := map[string]string{
		filepath.Join(serviceDir, "register.go"):                      registerSrc,
		filepath.Join(serviceDir, "resource_load_balancer_nlb.go"):    sameNameResourceImplSrc,
		filepath.Join(serviceDir, "data_source_load_balancer_nlb.go"): sameNameDataSourceImplSrc,
	}
	if docsR != "" {
		files[sameNameDocsRPath] = docsR
	}
	if docsD != "" {
		files[sameNameDocsDPath] = docsD
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// modifiedDiff builds the parsed diff runDiffCheck consumes. The checker only
// reads file paths and modes — never hunks — so a file list is a faithful
// stand-in for a real diff.
func modifiedDiff(paths ...string) *diffparser.Diff {
	files := make([]*diffparser.DiffFile, 0, len(paths))
	for _, p := range paths {
		files = append(files, &diffparser.DiffFile{Mode: diffparser.MODIFIED, OrigName: p, NewName: p})
	}
	return &diffparser.Diff{Files: files}
}

// TestUnitRunDiffCheckSameNameDocsOnly covers a docs-only diff that touches
// both pages of the pair. Before the (kind, name) identity, the second page
// was deduplicated away by the name-only map and both entities shared — or
// missed — a single comparison. Register order and diff order are flipped
// between subtests: the outcome must not depend on either.
func TestUnitRunDiffCheckSameNameDocsOnly(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		register string
		pages    []string
	}{
		{"resource first in register, r page first in diff", sameNameRegisterResourceFirst, []string{sameNameDocsRPath, sameNameDocsDPath}},
		{"data source first in register, d page first in diff", sameNameRegisterDataSourceFirst, []string{sameNameDocsDPath, sameNameDocsRPath}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeSameNameConsistencyTree(t, tc.register, sameNameDocsR, sameNameDocsD)

			report, exitCode := runDiffCheck(modifiedDiff(tc.pages...), sameNameFrameworkNames(), &schema.Provider{})

			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0 (both entities documented consistently)", exitCode)
			}
			if report.entities != 2 {
				t.Errorf("entities = %d, want 2: the same-named pair is two entities", report.entities)
			}
			if report.frameworkChecked != 2 {
				t.Errorf("frameworkChecked = %d, want 2: each entity must be compared once", report.frameworkChecked)
			}
			if report.compared != 2 || report.passed != 2 {
				t.Errorf("compared = %d, passed = %d, want 2 and 2", report.compared, report.passed)
			}
			if len(report.failedNames) != 0 {
				t.Errorf("failedNames = %v, want none", report.failedNames)
			}
		})
	}
}

// TestUnitRunDiffCheckSameNameSourceOnly covers a source-only diff that
// touches both implementation files of the pair: the split-layout resolution
// yields one registration per file, each carrying its kind.
func TestUnitRunDiffCheckSameNameSourceOnly(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		files []string
	}{
		{"resource impl first", []string{sameNameResourceImplPath, sameNameDataSourcePath}},
		{"data source impl first", []string{sameNameDataSourcePath, sameNameResourceImplPath}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)

			report, exitCode := runDiffCheck(modifiedDiff(tc.files...), sameNameFrameworkNames(), &schema.Provider{})

			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0 (both entities documented consistently)", exitCode)
			}
			if report.entities != 2 || report.frameworkChecked != 2 || report.compared != 2 || report.passed != 2 {
				t.Errorf("entities=%d frameworkChecked=%d compared=%d passed=%d, want 2/2/2/2",
					report.entities, report.frameworkChecked, report.compared, report.passed)
			}
			if len(report.failedNames) != 0 {
				t.Errorf("failedNames = %v, want none", report.failedNames)
			}
		})
	}
}

// TestUnitRunDiffCheckSameNameSingleEntityDefect gives only the resource a
// real documentation defect — a page the parser cannot open, then a page
// missing one attribute — and requires the data source of the pair to keep
// passing: the failure must be scoped to the entity whose docs are actually
// wrong, not smeared across the shared name.
func TestUnitRunDiffCheckSameNameSingleEntityDefect(t *testing.T) {
	t.Run("resource docs page missing", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, "", sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsRPath, sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (the resource page is missing)", exitCode)
		}
		if report.entities != 2 {
			t.Errorf("entities = %d, want 2", report.entities)
		}
		if report.passed != 1 || report.docsUnreadable != 1 {
			t.Errorf("passed = %d, docsUnreadable = %d, want the data source passed (1) and the resource page unreadable (1)",
				report.passed, report.docsUnreadable)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sameNameEntity)
		}
	})
	t.Run("resource docs page missing one attribute", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, sameNameDocsRMissingField, sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsRPath, sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (server_group_id is undocumented)", exitCode)
		}
		if report.entities != 2 {
			t.Errorf("entities = %d, want 2", report.entities)
		}
		if report.compared != 2 || report.passed != 1 {
			t.Errorf("compared = %d, passed = %d, want both compared (2) with only the data source passing (1)",
				report.compared, report.passed)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sameNameEntity)
		}
	})
}

// writeSameNameSharedImplTree materializes the pair with both factories
// declared in one implementation file. The registration is the same as the
// separated layout; only the file layout changes, so a differing gate outcome
// between the two trees is attributable to the shared file alone.
func writeSameNameSharedImplTree(t *testing.T, registerSrc, docsR, docsD string) {
	t.Helper()
	serviceDir := filepath.Join("alicloud", "service", "sae")
	files := map[string]string{
		filepath.Join(serviceDir, "register.go"):          registerSrc,
		filepath.Join(serviceDir, "load_balancer_nlb.go"): sameNameSharedImplSrc,
	}
	if docsR != "" {
		files[sameNameDocsRPath] = docsR
	}
	if docsD != "" {
		files[sameNameDocsDPath] = docsD
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestUnitRunDiffCheckSameNameFactoriesShareFile pins the file-granularity
// guard. When both factories of the pair live in one implementation file, the
// whole-file schema parse cannot isolate the entities — each arrival would be
// handed the union of their attributes — so the gate fails the entity instead
// of comparing a merged schema against its docs. The separated layout under
// the same registration and the same docs-only diff keeps passing, so the
// failure comes from the shared file, not from the pair itself.
func TestUnitRunDiffCheckSameNameFactoriesShareFile(t *testing.T) {
	t.Run("both factories in one file fail the entity instead of comparing the merged schema", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameSharedImplTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (the shared implementation file cannot isolate the pair)", exitCode)
		}
		if report.entities != 1 || report.frameworkChecked != 1 {
			t.Errorf("entities = %d, frameworkChecked = %d, want 1 and 1 (the d/ page routes to the data source alone)", report.entities, report.frameworkChecked)
		}
		if report.compared != 0 || report.passed != 0 || report.schemaUnreadable != 1 {
			t.Errorf("compared = %d, passed = %d, schemaUnreadable = %d, want 0, 0 and 1 (the guard fails before the comparison)", report.compared, report.passed, report.schemaUnreadable)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sameNameEntity)
		}
	})
	t.Run("one factory per file keeps the same diff passing", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (the data source is documented consistently)", exitCode)
		}
		if report.entities != 1 || report.frameworkChecked != 1 || report.compared != 1 || report.passed != 1 {
			t.Errorf("entities=%d frameworkChecked=%d compared=%d passed=%d, want 1/1/1/1",
				report.entities, report.frameworkChecked, report.compared, report.passed)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
}

// The alias fixture for the alias-registration guard test: one factory
// (NewLoadBalancerNLBResource) registered under two type names — the
// alicloud_disk/alicloud_ecs_disk and sslcertificateservice legacy-alias
// shape — while the data source of the pair keeps its own separate factory.
const aliasEntity = sameNameEntity + "_alias"

const sameNameAliasRegister = `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBResource,
			},
			{
				TypeName: "alicloud_sae_load_balancer_nlb_alias",
				Factory:  NewLoadBalancerNLBResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
	}
}
`

// The alias entity's page documents the same schema the shared factory
// declares: an alias reuses the factory, so it reuses its schema too.
const sameNameDocsRAlias = `# alicloud_sae_load_balancer_nlb_alias

Provides an SAE NLB resource under its legacy alias.

## Argument Reference

The following arguments are supported:

* ` + "`listener_port`" + ` - (Required) The port of the listener.
* ` + "`server_group_id`" + ` - (Optional) The id of the server group.
`

const sameNameDocsRAliasMissingField = `# alicloud_sae_load_balancer_nlb_alias

Provides an SAE NLB resource under its legacy alias.

## Argument Reference

The following arguments are supported:

* ` + "`listener_port`" + ` - (Required) The port of the listener.
`

const sameNameDocsRAliasPath = "website/docs/r/sae_load_balancer_nlb_alias.html.markdown"

func aliasFrameworkNames() map[splitlayout.Kind]map[string]struct{} {
	return map[splitlayout.Kind]map[string]struct{}{
		splitlayout.KindResource:   {sameNameEntity: {}, aliasEntity: {}},
		splitlayout.KindDataSource: {sameNameEntity: {}},
	}
}

// writeSameNameAliasTree materializes the aliased pair: the register reuses
// one resource factory for two type names, each entity keeps its own docs
// page, and the implementation layout is one factory per file.
func writeSameNameAliasTree(t *testing.T, docsR, docsRAlias, docsD string) {
	t.Helper()
	serviceDir := filepath.Join("alicloud", "service", "sae")
	files := map[string]string{
		filepath.Join(serviceDir, "register.go"):                      sameNameAliasRegister,
		filepath.Join(serviceDir, "resource_load_balancer_nlb.go"):    sameNameResourceImplSrc,
		filepath.Join(serviceDir, "data_source_load_balancer_nlb.go"): sameNameDataSourceImplSrc,
	}
	if docsR != "" {
		files[sameNameDocsRPath] = docsR
	}
	if docsRAlias != "" {
		files[sameNameDocsRAliasPath] = docsRAlias
	}
	if docsD != "" {
		files[sameNameDocsDPath] = docsD
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestUnitRunDiffCheckAliasRegistrationsShareFactory pins the alias side of
// the file-granularity guard. Alias registrations reuse one factory for
// several type names — a layout the repo documents explicitly (the
// sslcertificateservice legacy aliases, alicloud_disk/alicloud_ecs_disk) —
// and the aliased file still declares exactly one constructor whose
// whole-file parse is that schema alone. The guard must count distinct
// factories, not registrations: counting entries would fail the entity the
// moment the framework side adopts the alias migration, the same class of
// false block this gate exists to eliminate. The guard must keep its real
// work too: a documentation defect of the aliased entity still fails that
// entity alone.
func TestUnitRunDiffCheckAliasRegistrationsShareFactory(t *testing.T) {
	t.Run("an aliased implementation file is compared once per type name", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameAliasTree(t, sameNameDocsR, sameNameDocsRAlias, sameNameDocsD)

		// The implementation file declares the factories of two registrations
		// (the entity and its alias) while the data source file declares one.
		report, exitCode := runDiffCheck(modifiedDiff(sameNameResourceImplPath, sameNameDataSourcePath), aliasFrameworkNames(), &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (one factory per file compares each type name against its own page)", exitCode)
		}
		if report.entities != 3 || report.frameworkChecked != 3 {
			t.Errorf("entities = %d, frameworkChecked = %d, want 3 and 3 (resource, alias and data source)", report.entities, report.frameworkChecked)
		}
		if report.compared != 3 || report.passed != 3 {
			t.Errorf("compared = %d, passed = %d, want 3 and 3", report.compared, report.passed)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
	t.Run("a docs-only diff checks the aliased entity with its own factory's schema", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameAliasTree(t, sameNameDocsR, sameNameDocsRAlias, sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsRAliasPath), aliasFrameworkNames(), &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (the alias page documents the shared factory's schema)", exitCode)
		}
		if report.entities != 1 || report.frameworkChecked != 1 || report.compared != 1 || report.passed != 1 {
			t.Errorf("entities=%d frameworkChecked=%d compared=%d passed=%d, want 1/1/1/1",
				report.entities, report.frameworkChecked, report.compared, report.passed)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
	t.Run("an aliased entity's own docs defect fails that entity alone", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameAliasTree(t, sameNameDocsR, sameNameDocsRAliasMissingField, sameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsRPath, sameNameDocsRAliasPath), aliasFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (server_group_id is undocumented on the alias page)", exitCode)
		}
		if report.entities != 2 || report.frameworkChecked != 2 {
			t.Errorf("entities = %d, frameworkChecked = %d, want 2 and 2 (both pages of the shared factory are in the diff)", report.entities, report.frameworkChecked)
		}
		if report.compared != 2 || report.passed != 1 {
			t.Errorf("compared = %d, passed = %d, want both compared (2) with only the resource passing (1)",
				report.compared, report.passed)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != aliasEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, aliasEntity)
		}
	})
}

// The ecs fixture for the cross-product factory-name test: its register
// references the same constructor name the sae package declares. Go function
// names are only package-unique, so a product that compiles a same-named
// factory declares it in its own package — which is what ecs's own
// implementation file does here.
const crossProductFactoryRegister = `package ecs

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecs",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ecs_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
	}
}
`

const crossProductFactoryImplSrc = `package ecs

func NewLoadBalancerNLBDataSource() *nlbDataSource {
	return &nlbDataSource{}
}
`

// writeEcsFactoryNameCoincidence materializes the ecs product whose register
// references a factory name the sae package also declares — the name
// coincidence the guard must not mistake for a shared file.
func writeEcsFactoryNameCoincidence(t *testing.T) {
	t.Helper()
	ecsDir := filepath.Join("alicloud", "service", "ecs")
	files := map[string]string{
		filepath.Join(ecsDir, "register.go"):                      crossProductFactoryRegister,
		filepath.Join(ecsDir, "data_source_load_balancer_nlb.go"): crossProductFactoryImplSrc,
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestUnitRunDiffCheckFactoryNameCoincidenceAcrossProducts pins the product
// scope of the file-granularity guard. Go function names are unique only
// inside a package and the framework factory names carry no product prefix,
// so another product's register can reference a factory whose name this file
// also declares — its own same-named function in its own package. The guard
// counts a file's factories from the file's own product register only:
// matching names across every register instead would report the ecs entry
// next to the sae one, trip the "several entities" guard on a perfectly
// isolated one-factory-per-file layout, and put another product's entity into
// the failure message — the same wrong-entity class the gate exists to
// eliminate. The within-product firing is pinned too, so the scoping cannot
// quietly lose the guard's real work.
func TestUnitRunDiffCheckFactoryNameCoincidenceAcrossProducts(t *testing.T) {
	t.Run("data source checks alone despite the ecs name coincidence", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)
		writeEcsFactoryNameCoincidence(t)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (the data source file declares one factory and is compared normally)", exitCode)
		}
		if report.entities != 1 || report.frameworkChecked != 1 || report.compared != 1 || report.passed != 1 {
			t.Errorf("entities=%d frameworkChecked=%d compared=%d passed=%d, want 1/1/1/1",
				report.entities, report.frameworkChecked, report.compared, report.passed)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
	t.Run("the guard still fires within the file's own product", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameSharedImplTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)
		writeEcsFactoryNameCoincidence(t)

		report, exitCode := runDiffCheck(modifiedDiff(sameNameDocsDPath), sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (the sae file shares the pair's factories and cannot isolate them)", exitCode)
		}
		if report.entities != 1 || report.frameworkChecked != 1 {
			t.Errorf("entities = %d, frameworkChecked = %d, want 1 and 1 (the d/ page routes to the data source alone)", report.entities, report.frameworkChecked)
		}
		if report.schemaUnreadable != 1 || report.compared != 0 || report.passed != 0 {
			t.Errorf("schemaUnreadable = %d, compared = %d, passed = %d, want 1, 0 and 0 (the guard fails before the comparison)",
				report.schemaUnreadable, report.compared, report.passed)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sameNameEntity)
		}
	})
}

// TestUnitRunDiffCheckSDKv2FlatLayout pins the single-type SDKv2 path: a flat
// implementation file dispatches through the provider maps exactly as it did
// before the kind was threaded through, including the data_source_ prefix
// fallback.
func TestUnitRunDiffCheckSDKv2FlatLayout(t *testing.T) {
	t.Chdir(t.TempDir())
	docs := map[string]string{
		"website/docs/r/widget.html.markdown":  "# alicloud_widget\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `name` - (Required) The name of the widget.\n",
		"website/docs/d/widgets.html.markdown": "# alicloud_widgets\n\n## Argument Reference\n\n* `name_regex` - (Optional) A regex string to filter widgets by name.\n",
	}
	for path, src := range docs {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
	provider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_widget": {Schema: map[string]*schema.Schema{
				"name": {Type: schema.TypeString, Required: true},
			}},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_widgets": {Schema: map[string]*schema.Schema{
				"name_regex": {Type: schema.TypeString, Optional: true},
			}},
		},
	}

	report, exitCode := runDiffCheck(
		modifiedDiff("alicloud/resource_alicloud_widget.go", "alicloud/data_source_alicloud_widgets.go"),
		map[splitlayout.Kind]map[string]struct{}{},
		provider,
	)

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0 (single-type SDKv2 entities documented consistently)", exitCode)
	}
	if report.entities != 2 || report.compared != 2 || report.passed != 2 {
		t.Errorf("entities=%d compared=%d passed=%d, want 2/2/2", report.entities, report.compared, report.passed)
	}
	if report.frameworkChecked != 0 {
		t.Errorf("frameworkChecked = %d, want 0: neither name is framework-native", report.frameworkChecked)
	}
}

// TestUnitRunDiffCheckOtherFrameworkKinds keeps the framework categories that
// never collide by name behaving as before: an ephemeral page resolves through
// its own directory, and a list resource — which has no directory of its own
// and documents under r/ — still reaches its framework check through the
// single-registration fallback.
func TestUnitRunDiffCheckOtherFrameworkKinds(t *testing.T) {
	t.Run("ephemeral page in its own directory", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeFrameworkKindTree(t, "kms", "ephemeral_plaintext", splitlayout.KindEphemeral,
			"NewPlaintextEphemeralResource", "payload", "ephemeral-resources")

		frameworkNames := map[splitlayout.Kind]map[string]struct{}{
			splitlayout.KindEphemeral: {"alicloud_kms_plaintext": {}},
		}
		report, exitCode := runDiffCheck(modifiedDiff("website/docs/ephemeral-resources/kms_plaintext.html.markdown"), frameworkNames, &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (ephemeral entity documented consistently)", exitCode)
		}
		if report.frameworkChecked != 1 || report.passed != 1 {
			t.Errorf("frameworkChecked=%d passed=%d, want 1/1", report.frameworkChecked, report.passed)
		}
	})
	t.Run("list resource page under r", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeFrameworkKindTree(t, "cas", "list_certificate", splitlayout.KindListResource,
			"NewCertificateListResource", "ids", "r")

		frameworkNames := map[splitlayout.Kind]map[string]struct{}{
			splitlayout.KindListResource: {"alicloud_cas_certificate": {}},
		}
		report, exitCode := runDiffCheck(modifiedDiff("website/docs/r/cas_certificate.html.markdown"), frameworkNames, &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (list resource documented consistently)", exitCode)
		}
		if report.frameworkChecked != 1 || report.passed != 1 {
			t.Errorf("frameworkChecked=%d passed=%d, want 1/1: the page under r/ must still reach the list-resource check",
				report.frameworkChecked, report.passed)
		}
	})
}

// TestUnitRunDiffCheckListResourceImplAndDocsSameDiff pins the entity-level
// dedup contract for a list resource whose implementation file and r/ docs
// page land in the same diff. The two arrivals state different kinds — the
// split-layout resolution carries the registered (list) kind, the r/
// directory says "resource" — so a dedup key built from the diff's own kind
// lets both through and the one entity is compared twice with every count
// inflated. The key must be built after the dispatch normalized the kind;
// the diff order of the two files must not matter.
func TestUnitRunDiffCheckListResourceImplAndDocsSameDiff(t *testing.T) {
	implPath := "alicloud/service/cas/list_certificate.go"
	docsPath := "website/docs/r/cas_certificate.html.markdown"
	frameworkNames := map[splitlayout.Kind]map[string]struct{}{
		splitlayout.KindListResource: {"alicloud_cas_certificate": {}},
	}
	for _, tc := range []struct {
		desc string
		diff []string
	}{
		{"implementation file first", []string{implPath, docsPath}},
		{"docs page first", []string{docsPath, implPath}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeFrameworkKindTree(t, "cas", "list_certificate", splitlayout.KindListResource,
				"NewCertificateListResource", "ids", "r")

			report, exitCode := runDiffCheck(modifiedDiff(tc.diff...), frameworkNames, &schema.Provider{})

			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0 (the list resource is documented consistently)", exitCode)
			}
			if report.entities != 1 {
				t.Errorf("entities = %d, want 1: the implementation and the docs page are one entity",
					report.entities)
			}
			if report.frameworkChecked != 1 {
				t.Errorf("frameworkChecked = %d, want 1: the entity must be compared once, not once per arrival kind",
					report.frameworkChecked)
			}
			if report.compared != 1 || report.passed != 1 {
				t.Errorf("compared = %d, passed = %d, want 1 and 1", report.compared, report.passed)
			}
			if len(report.failedNames) != 0 {
				t.Errorf("failedNames = %v, want none", report.failedNames)
			}
		})
	}
	// The same single-check contract holds when the page is actually wrong:
	// the one entity fails once, not once per arrival kind.
	t.Run("docs page missing the attribute", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeFrameworkKindTree(t, "cas", "list_certificate", splitlayout.KindListResource,
			"NewCertificateListResource", "ids", "r")
		if err := os.WriteFile(docsPath, []byte("# alicloud_cas_certificate\n\n## Argument Reference\n"), 0644); err != nil {
			t.Fatal(err)
		}

		report, exitCode := runDiffCheck(modifiedDiff(implPath, docsPath), frameworkNames, &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (the page documents no attribute)", exitCode)
		}
		if report.entities != 1 || report.compared != 1 || report.passed != 0 {
			t.Errorf("entities=%d compared=%d passed=%d, want 1/1/0: the defect is reported once",
				report.entities, report.compared, report.passed)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != "alicloud_cas_certificate" {
			t.Errorf("failedNames = %v, want [alicloud_cas_certificate]", report.failedNames)
		}
	})
}

// TestUnitRunDiffCheckSDKv2FlatDataSourceImplAndDocsSameDiff pins the
// entity-level dedup contract for an SDKv2 data source whose flat
// implementation file and d/ docs page land in the same diff. The flat file
// names the entity with its data_source_ prefix while the docs page names it
// without, so a key built from the raw name lets both through and the one
// entity is compared twice. The key normalizes the prefix; the diff order of
// the two files must not matter.
func TestUnitRunDiffCheckSDKv2FlatDataSourceImplAndDocsSameDiff(t *testing.T) {
	docsPath := "website/docs/d/widgets.html.markdown"
	for _, tc := range []struct {
		desc string
		diff []string
	}{
		{"implementation file first", []string{"alicloud/data_source_alicloud_widgets.go", docsPath}},
		{"docs page first", []string{docsPath, "alicloud/data_source_alicloud_widgets.go"}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			if err := os.MkdirAll(filepath.Dir(docsPath), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(docsPath, []byte("# alicloud_widgets\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `name_regex` - (Optional) A regex string to filter widgets by name.\n"), 0644); err != nil {
				t.Fatal(err)
			}
			provider := &schema.Provider{
				DataSourcesMap: map[string]*schema.Resource{
					"alicloud_widgets": {Schema: map[string]*schema.Schema{
						"name_regex": {Type: schema.TypeString, Optional: true},
					}},
				},
			}

			report, exitCode := runDiffCheck(modifiedDiff(tc.diff...), map[splitlayout.Kind]map[string]struct{}{}, provider)

			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0 (the data source is documented consistently)", exitCode)
			}
			if report.entities != 1 {
				t.Errorf("entities = %d, want 1: the flat file and the docs page are one entity however they name it",
					report.entities)
			}
			if report.compared != 1 || report.passed != 1 {
				t.Errorf("compared = %d, passed = %d, want 1 and 1", report.compared, report.passed)
			}
			if report.frameworkChecked != 0 {
				t.Errorf("frameworkChecked = %d, want 0: the entity is SDKv2", report.frameworkChecked)
			}
			if len(report.failedNames) != 0 {
				t.Errorf("failedNames = %v, want none", report.failedNames)
			}
		})
	}
}

// The SDKv2 same-named pair the review asked this regression for:
// alicloud_actiontrail_global_events_storage_region is registered in the real
// actiontrail service package both as an SDKv2 resource and as an SDKv2 data
// source, so the provider carries one entity of that name in each map. Before
// the kind picked the map, the SDKv2 lookup probed ResourcesMap first and the
// resource claimed the data source's docs arrival: a diff touching only the
// d/ page compared the resource's schema against the r/ page, and the touched
// d/ page slipped through unchecked.
const sdkv2SameNameEntity = "alicloud_actiontrail_global_events_storage_region"

const (
	sdkv2SameNameDocsRPath = "website/docs/r/actiontrail_global_events_storage_region.html.markdown"
	sdkv2SameNameDocsDPath = "website/docs/d/actiontrail_global_events_storage_region.html.markdown"
)

// The two schemas share only storage_region: passing or failing must come
// from the entity the kind selected, never from a schema that fits both
// pages.
func sdkv2SameNameProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			sdkv2SameNameEntity: {Schema: map[string]*schema.Schema{
				"actiontrail_name": {Type: schema.TypeString, Required: true},
				"storage_region":   {Type: schema.TypeString, Required: true},
			}},
		},
		DataSourcesMap: map[string]*schema.Resource{
			sdkv2SameNameEntity: {Schema: map[string]*schema.Schema{
				"storage_region": {Type: schema.TypeString, Optional: true},
				"output_file":    {Type: schema.TypeString, Optional: true},
			}},
		},
	}
}

const sdkv2SameNameDocsR = "# alicloud_actiontrail_global_events_storage_region\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `actiontrail_name` - (Required) The name of the actiontrail.\n* `storage_region` - (Required) The region storing the events.\n"

// sdkv2SameNameDocsRMissingField drops storage_region: a defect of the r/
// page only, used to prove a d/-only diff never reads that page.
const sdkv2SameNameDocsRMissingField = "# alicloud_actiontrail_global_events_storage_region\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `actiontrail_name` - (Required) The name of the actiontrail.\n"

const sdkv2SameNameDocsD = "# alicloud_actiontrail_global_events_storage_region\n\nThis data source provides the global events storage region.\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `storage_region` - (Optional) The region storing the events.\n* `output_file` - (Optional) The output file path.\n"

// sdkv2SameNameDocsDWrongField documents the resource's actiontrail_name on
// the data source's page and drops the data source's output_file: a defect of
// the touched d/ page only. The r/ page stays correct, so a checker that
// still compared the resource's schema against the r/ page would pass and
// the wrong field would go unnoticed.
const sdkv2SameNameDocsDWrongField = "# alicloud_actiontrail_global_events_storage_region\n\nThis data source provides the global events storage region.\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `storage_region` - (Optional) The region storing the events.\n* `actiontrail_name` - (Required) The name of the actiontrail.\n"

// writeSDKv2SameNameDocs materializes whichever pages the subtest needs; an
// empty string leaves that page out entirely.
func writeSDKv2SameNameDocs(t *testing.T, docsR, docsD string) {
	t.Helper()
	files := map[string]string{}
	if docsR != "" {
		files[sdkv2SameNameDocsRPath] = docsR
	}
	if docsD != "" {
		files[sdkv2SameNameDocsDPath] = docsD
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestUnitRunDiffCheckSDKv2SameNameKindDispatch pins the SDKv2 side of the
// (kind, name) identity: the kind the diff already identified — the docs
// directory here — selects the provider map and the docs page, so a
// same-named SDKv2 resource and data source are each compared against their
// own schema and their own page.
func TestUnitRunDiffCheckSDKv2SameNameKindDispatch(t *testing.T) {
	t.Run("touched d page consistent", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSDKv2SameNameDocs(t, sdkv2SameNameDocsR, sdkv2SameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sdkv2SameNameDocsDPath), map[splitlayout.Kind]map[string]struct{}{}, sdkv2SameNameProvider())

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (the touched d/ page matches the data source)", exitCode)
		}
		if report.entities != 1 || report.compared != 1 || report.passed != 1 {
			t.Errorf("entities=%d compared=%d passed=%d, want 1/1/1: the d/ page checks the data source alone",
				report.entities, report.compared, report.passed)
		}
		if report.frameworkChecked != 0 {
			t.Errorf("frameworkChecked = %d, want 0: the pair is SDKv2", report.frameworkChecked)
		}
	})
	t.Run("touched d page with a wrong field fails", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSDKv2SameNameDocs(t, sdkv2SameNameDocsR, sdkv2SameNameDocsDWrongField)

		report, exitCode := runDiffCheck(modifiedDiff(sdkv2SameNameDocsDPath), map[splitlayout.Kind]map[string]struct{}{}, sdkv2SameNameProvider())

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (actiontrail_name is not an argument of the data source)", exitCode)
		}
		if report.entities != 1 || report.compared != 1 || report.passed != 0 {
			t.Errorf("entities=%d compared=%d passed=%d, want 1/1/0: the d/ page must be the page compared",
				report.entities, report.compared, report.passed)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sdkv2SameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sdkv2SameNameEntity)
		}
	})
	t.Run("untouched r page defect does not bleed into the d check", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSDKv2SameNameDocs(t, sdkv2SameNameDocsRMissingField, sdkv2SameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sdkv2SameNameDocsDPath), map[splitlayout.Kind]map[string]struct{}{}, sdkv2SameNameProvider())

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (only the d/ page is in the diff and it is consistent)", exitCode)
		}
		if report.compared != 1 || report.passed != 1 {
			t.Errorf("compared = %d, passed = %d, want 1 and 1: the check reads the d/ page, not the r/ one",
				report.compared, report.passed)
		}
	})
	t.Run("touched r page with a wrong field fails", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSDKv2SameNameDocs(t, sdkv2SameNameDocsRMissingField, sdkv2SameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sdkv2SameNameDocsRPath), map[splitlayout.Kind]map[string]struct{}{}, sdkv2SameNameProvider())

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (storage_region is undocumented on the r/ page)", exitCode)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sdkv2SameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sdkv2SameNameEntity)
		}
	})
	t.Run("both pages in the diff", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSDKv2SameNameDocs(t, sdkv2SameNameDocsR, sdkv2SameNameDocsD)

		report, exitCode := runDiffCheck(modifiedDiff(sdkv2SameNameDocsRPath, sdkv2SameNameDocsDPath), map[splitlayout.Kind]map[string]struct{}{}, sdkv2SameNameProvider())

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (both entities documented consistently)", exitCode)
		}
		if report.entities != 2 || report.compared != 2 || report.passed != 2 {
			t.Errorf("entities=%d compared=%d passed=%d, want 2/2/2: the same-named pair is two entities",
				report.entities, report.compared, report.passed)
		}
	})
}

// TestUnitRunDiffCheckSDKv2DataSourceOnlyBothPages pins the SDKv2 dedup key
// for a name registered only in DataSourcesMap whose r/ and d/ pages are both
// in the diff. The historical resource→data-source fallback rewrites the r/
// arrival's class to the data source, so the key must be built from the
// settled class: keying on the arrival's own kind would give one entity two
// keys and compare it twice with every count inflated, whichever order the
// diff lists the pages in. The misplaced r/ page itself stays un-compared,
// exactly as before this change — that redirect is preserved base behavior,
// and making it fail instead is a deliberate tightening that waits for a
// maintainer's call, so the r/ page here documents a field the data source
// does not have: the moment the gate starts comparing it, this test fails
// and the tightening has landed.
func TestUnitRunDiffCheckSDKv2DataSourceOnlyBothPages(t *testing.T) {
	provider := &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_widgets": {Schema: map[string]*schema.Schema{
				"name_regex": {Type: schema.TypeString, Optional: true},
			}},
		},
	}
	docsD := "website/docs/d/widgets.html.markdown"
	docsR := "website/docs/r/widgets.html.markdown"
	pages := map[string]string{
		docsD: "# alicloud_widgets\n\nThis data source provides the widgets.\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `name_regex` - (Optional) A regex string to filter widgets by name.\n",
		docsR: "# alicloud_widgets\n\n## Argument Reference\n\nThe following arguments are supported:\n\n* `bogus_field` - (Optional) Not an argument of the data source.\n",
	}
	for _, tc := range []struct {
		desc string
		diff []string
	}{
		{"r page first", []string{docsR, docsD}},
		{"d page first", []string{docsD, docsR}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			for path, src := range pages {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(src), 0644); err != nil {
					t.Fatal(err)
				}
			}

			report, exitCode := runDiffCheck(modifiedDiff(tc.diff...), map[splitlayout.Kind]map[string]struct{}{}, provider)

			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0 (the entity's own d/ page is consistent)", exitCode)
			}
			if report.entities != 1 || report.compared != 1 || report.passed != 1 {
				t.Errorf("entities=%d compared=%d passed=%d, want 1/1/1: the fallback settles both arrivals on one data source identity",
					report.entities, report.compared, report.passed)
			}
			if report.frameworkChecked != 0 {
				t.Errorf("frameworkChecked = %d, want 0: the entity is SDKv2", report.frameworkChecked)
			}
			if len(report.failedNames) != 0 {
				t.Errorf("failedNames = %v, want none", report.failedNames)
			}
		})
	}
}

// TestUnitRunDiffCheckDocsPageUnclassifiableFails keeps a page that matches
// the docs pattern but lives in no directory the checker maps to a kind from
// passing silently. The pattern is not anchored, so a deeper page still
// matches while its directory cannot be read; the gate must fail on what it
// cannot compare instead of warning and skipping.
func TestUnitRunDiffCheckDocsPageUnclassifiableFails(t *testing.T) {
	t.Chdir(t.TempDir())
	deepPath := "vendor/website/docs/r/orphan.html.markdown"

	report, exitCode := runDiffCheck(modifiedDiff(deepPath), map[splitlayout.Kind]map[string]struct{}{}, &schema.Provider{})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1 (the page is in the diff and cannot be classified)", exitCode)
	}
	if report.entities != 0 {
		t.Errorf("entities = %d, want 0: an unclassifiable page is not a checked entity", report.entities)
	}
	if len(report.failedNames) != 1 || report.failedNames[0] != deepPath {
		t.Errorf("failedNames = %v, want [%s]", report.failedNames, deepPath)
	}
}

// TestUnitRunDiffCheckFrameworkFlatFilePrefersRegistration pins which path is
// the authority for a framework entity's implementation: a diff-supplied
// flat-layout path that disagrees with the registration must lose to the
// registered implementation, because the registration is what pairs the type
// name with the factory that declares its schema. The flat path here does not
// even exist on disk — trusting it would fail on a missing file instead of
// comparing the entity the dispatch actually claimed.
func TestUnitRunDiffCheckFrameworkFlatFilePrefersRegistration(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFrameworkKindTree(t, "widget", "resource_gadget", splitlayout.KindResource,
		"NewGadgetResource", "name", "r")
	frameworkNames := map[splitlayout.Kind]map[string]struct{}{
		splitlayout.KindResource: {"alicloud_widget_gadget": {}},
	}

	report, exitCode := runDiffCheck(modifiedDiff("alicloud/resource_alicloud_widget_gadget.go"), frameworkNames, &schema.Provider{})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0 (the entity compares against its registered implementation)", exitCode)
	}
	if report.frameworkChecked != 1 || report.compared != 1 || report.passed != 1 {
		t.Errorf("frameworkChecked=%d compared=%d passed=%d, want 1/1/1",
			report.frameworkChecked, report.compared, report.passed)
	}
	if len(report.failedNames) != 0 {
		t.Errorf("failedNames = %v, want none", report.failedNames)
	}
}

// TestUnitRunDiffCheckDispatchFallbackOnlyForListResources pins the one pair
// the single-registration fallback may substitute: a resource-kind arrival of
// a name registered only as a list resource. Any other cross-kind
// substitution is the wrong-entity comparison this gate exists to prevent —
// a name the SDKv2 provider registers under the diff's kind must stay on the
// SDKv2 path (alicloud_kms_secret is both an SDKv2 resource and a framework
// ephemeral), and a docs page naming a kind nothing registered must fail
// instead of silently checking some other entity's page.
func TestUnitRunDiffCheckDispatchFallbackOnlyForListResources(t *testing.T) {
	t.Run("cross-dialect same name stays SDKv2", func(t *testing.T) {
		t.Chdir(t.TempDir())
		// One register declaring the same name as an SDKv2 resource and as a
		// framework ephemeral, the way service/kms declares alicloud_kms_secret.
		serviceDir := filepath.Join("alicloud", "service", "gadget")
		registerSrc := "package gadget\n\nfunc ServicePackage() conns.ServicePackage {\n	return conns.ServicePackage{\n		SDKResources: []conns.SDKResource{\n			{\n				TypeName: \"alicloud_gadget\",\n				Factory:  resourceGadget,\n			},\n		},\n		EphemeralResources: []conns.EphemeralResource{\n			{\n				TypeName: \"alicloud_gadget\",\n				Factory:  NewGadgetEphemeralResource,\n			},\n		},\n	}\n}\n"
		files := map[string]string{
			filepath.Join(serviceDir, "register.go"):                      registerSrc,
			filepath.Join(serviceDir, "resource_gadget.go"):               "package gadget\n\nfunc resourceGadget() *schema.Resource { return nil }\n",
			filepath.Join(serviceDir, "ephemeral_gadget.go"):              "package gadget\n\nfunc NewGadgetEphemeralResource() *entity { return nil }\n\nfunc schemaDecl() map[string]schema.Attribute {\n	return map[string]schema.Attribute{\n		\"payload\": schema.StringAttribute{Optional: true},\n	}\n}\n",
			filepath.Join("website", "docs", "r", "gadget.html.markdown"): "# alicloud_gadget\n\n## Argument Reference\n\n* `name` - (Required) The name of the gadget.\n",
		}
		for path, src := range files {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(src), 0644); err != nil {
				t.Fatal(err)
			}
		}
		provider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"alicloud_gadget": {Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Required: true},
				}},
			},
		}
		frameworkNames := map[splitlayout.Kind]map[string]struct{}{
			splitlayout.KindEphemeral: {"alicloud_gadget": {}},
		}

		report, exitCode := runDiffCheck(modifiedDiff("alicloud/service/gadget/resource_gadget.go"), frameworkNames, provider)

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (the SDKv2 resource is documented consistently)", exitCode)
		}
		if report.entities != 1 || report.compared != 1 || report.passed != 1 {
			t.Errorf("entities=%d compared=%d passed=%d, want 1/1/1: the SDKv2 entity is checked once",
				report.entities, report.compared, report.passed)
		}
		if report.frameworkChecked != 0 {
			t.Errorf("frameworkChecked = %d, want 0: the ephemeral registration must not claim the SDKv2 resource's arrival",
				report.frameworkChecked)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
	t.Run("docs page naming an unregistered kind fails", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeFrameworkKindTree(t, "gizmo", "resource_thing", splitlayout.KindResource,
			"NewThingResource", "name", "r")
		frameworkNames := map[splitlayout.Kind]map[string]struct{}{
			splitlayout.KindResource: {"alicloud_gizmo_thing": {}},
		}

		// The d/ page names a data source that nothing registers; substituting
		// the resource registration would check the untouched r/ page instead
		// and let the dangling d/ page pass.
		report, exitCode := runDiffCheck(modifiedDiff("website/docs/d/gizmo_thing.html.markdown"), frameworkNames, &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (the d/ page documents no data source)", exitCode)
		}
		if report.unresolved != 1 {
			t.Errorf("unresolved = %d, want 1: the name must reach the SDKv2 path and be reported missing", report.unresolved)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != "alicloud_gizmo_thing" {
			t.Errorf("failedNames = %v, want [alicloud_gizmo_thing]", report.failedNames)
		}
	})
}

// writeFrameworkKindTree materializes one framework entity of a given kind:
// its register entry, an implementation file declaring the factory and a
// schema, and a docs page under the given docs directory.
func writeFrameworkKindTree(t *testing.T, product, stem string, kind splitlayout.Kind, factory, attribute, docsDir string) {
	t.Helper()
	typeName := "alicloud_" + product + "_" + strings_TrimPrefixList(stem)
	// The docs page is named after the type name minus the alicloud_ prefix,
	// which embeds the product segment — not after the implementation stem.
	typeBase := product + "_" + strings_TrimPrefixList(stem)
	serviceDir := filepath.Join("alicloud", "service", product)
	registerSrc := "package " + product + "\n\nfunc ServicePackage() conns.ServicePackage {\n	return conns.ServicePackage{\n		Name: " + `"` + title(product) + `"` + ",\n		" + kindField(kind) + ": []conns." + kindInterface(kind) + "{\n			{\n				TypeName: " + `"` + typeName + `"` + ",\n				Factory:  " + factory + ",\n			},\n		},\n	}\n}\n"
	implSrc := "package " + product + "\n\nfunc " + factory + "() *entity { return nil }\n\nfunc schemaDecl() map[string]schema.Attribute {\n	return map[string]schema.Attribute{\n		" + `"` + attribute + `"` + ": schema.StringAttribute{Optional: true},\n	}\n}\n"
	docsSrc := "# " + typeName + "\n\n## Argument Reference\n\n* `" + attribute + "` - (Optional) The " + attribute + " of the " + product + ".\n"
	files := map[string]string{
		filepath.Join(serviceDir, "register.go"):                             registerSrc,
		filepath.Join(serviceDir, stem+".go"):                                implSrc,
		filepath.Join("website", "docs", docsDir, typeBase+".html.markdown"): docsSrc,
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func strings_TrimPrefixList(stem string) string {
	// stems look like "ephemeral_plaintext" or "list_certificate"; the type
	// name drops the kind prefix.
	for _, prefix := range []string{"ephemeral_", "list_", "action_", "resource_", "data_source_"} {
		if len(stem) > len(prefix) && stem[:len(prefix)] == prefix {
			return stem[len(prefix):]
		}
	}
	return stem
}

func title(product string) string {
	if product == "" {
		return product
	}
	return string(rune(product[0]-32)) + product[1:]
}

func kindField(kind splitlayout.Kind) string {
	switch kind {
	case splitlayout.KindResource:
		return "Resources"
	case splitlayout.KindDataSource:
		return "DataSources"
	case splitlayout.KindEphemeral:
		return "EphemeralResources"
	case splitlayout.KindListResource:
		return "ListResources"
	case splitlayout.KindAction:
		return "Actions"
	}
	return ""
}

func kindInterface(kind splitlayout.Kind) string {
	switch kind {
	case splitlayout.KindResource:
		return "Resource"
	case splitlayout.KindDataSource:
		return "DataSource"
	case splitlayout.KindEphemeral:
		return "EphemeralResource"
	case splitlayout.KindListResource:
		return "ListResource"
	case splitlayout.KindAction:
		return "Action"
	}
	return ""
}

// TestUnitRunDiffCheckDeletedDocsPage keeps the deletion gate: removing a page
// of a still-registered entity fails, removing page and entity together does
// not.
func TestUnitRunDiffCheckDeletedDocsPage(t *testing.T) {
	t.Run("entity still registered", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameConsistencyTree(t, sameNameRegisterResourceFirst, sameNameDocsR, sameNameDocsD)

		diff := &diffparser.Diff{Files: []*diffparser.DiffFile{{
			Mode:     diffparser.DELETED,
			OrigName: sameNameDocsRPath,
		}}}
		report, exitCode := runDiffCheck(diff, sameNameFrameworkNames(), &schema.Provider{})

		if exitCode != 1 {
			t.Fatalf("exitCode = %d, want 1 (page deleted while the entity is registered)", exitCode)
		}
		if len(report.failedNames) != 1 || report.failedNames[0] != sameNameEntity {
			t.Errorf("failedNames = %v, want [%s]", report.failedNames, sameNameEntity)
		}
	})
	t.Run("entity removed with its page", func(t *testing.T) {
		t.Chdir(t.TempDir())

		diff := &diffparser.Diff{Files: []*diffparser.DiffFile{{
			Mode:     diffparser.DELETED,
			OrigName: "website/docs/r/gone_entity.html.markdown",
		}}}
		report, exitCode := runDiffCheck(diff, map[splitlayout.Kind]map[string]struct{}{}, &schema.Provider{})

		if exitCode != 0 {
			t.Fatalf("exitCode = %d, want 0 (entity and page are both gone)", exitCode)
		}
		if len(report.failedNames) != 0 {
			t.Errorf("failedNames = %v, want none", report.failedNames)
		}
	})
}
