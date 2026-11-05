package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

// flatFile builds a flat titleFile with explicit entities, the way
// resolveEntities produces them from the filename derivation.
func flatFile(name string, added bool, entities ...entity) titleFile {
	return titleFile{name: name, added: added, flat: true, entities: entities}
}

// serviceFile builds a service titleFile with explicit entities, the way
// resolveEntities produces them from the register.go join.
func serviceFile(name string, entities ...entity) titleFile {
	return titleFile{name: name, service: true, entities: entities}
}

// ent builds one entity carrying the kind, existing state and product the
// kind-aware serviceEntities derivation would produce. The aggregate waiver
// tests use it to exercise parseAggregate/checkTitle logic; the kind-collision
// e2e tests below drive the real base index instead of hand-setting existing.
func ent(name string, dataSource, existing bool, product string) entity {
	return entity{typeName: name, dataSource: dataSource, existing: existing, product: product}
}

// svcEntities builds one service titleFile whose entities carry an explicit
// product, the way serviceEntities derives them after a real register.go join
// against a base index.
func svcEntities(product string, es ...entity) titleFile {
	return titleFile{name: "alicloud/service/" + product + "/resource.go", service: true, entities: es}
}

func TestCheckTitleRequiresModifiedEntityInTitle(t *testing.T) {
	files := []titleFile{flatFile("alicloud/resource_alicloud_instance.go", false,
		entity{typeName: "alicloud_instance", existing: true})}
	if checkTitle("fix: something else", files) {
		t.Error("checkTitle(modified entity missing from title) passed, want fail")
	}
	if !checkTitle("fix: resource/alicloud_instance: fix the instance bootstrap", files) {
		t.Error("checkTitle(modified entity in title) failed, want pass")
	}
}

func TestCheckTitleRequiresAddedEntityInTitle(t *testing.T) {
	files := []titleFile{flatFile("alicloud/data_source_alicloud_instance_types.go", true,
		entity{typeName: "alicloud_instance_types", dataSource: true})}
	if checkTitle("feat: add something", files) {
		t.Error("checkTitle(added data source missing) passed, want fail")
	}
	if !checkTitle("feat: New Data Source: alicloud_instance_types", files) {
		t.Error("checkTitle(added data source announced) failed, want pass")
	}
}

func TestCheckTitleFlatFilenameDrivesTheDemand(t *testing.T) {
	// The flat side stayed filename-derived after the #10839 revert, so a
	// shared plumbing file still demands its filename stem — including
	// data_source_alicloud_common.go.
	files := []titleFile{flatFile("alicloud/data_source_alicloud_common.go", false,
		entity{typeName: "alicloud_common", dataSource: true, existing: true})}
	if checkTitle("refactor: migrate the cloudauth face configs data source", files) {
		t.Error("checkTitle(shared plumbing without its stem) passed, want fail")
	}
	if !checkTitle("refactor: data-source/alicloud_common: migrate the face configs data source", files) {
		t.Error("checkTitle(shared plumbing with its stem) failed, want pass")
	}
}

func TestCheckTitleMultiEntityFileRequiresEveryEntity(t *testing.T) {
	files := []titleFile{serviceFile("alicloud/service/ims/data_source_users.go",
		entity{typeName: "alicloud_ims_users", dataSource: true, existing: true},
		entity{typeName: "alicloud_ims_subnets", dataSource: true, existing: true})}
	if checkTitle("fix: data-source/alicloud_ims_users: one of two", files) {
		t.Error("checkTitle(one of two entities) passed, want fail")
	}
	if !checkTitle("fix: data-source/alicloud_ims_users: x and data-source/alicloud_ims_subnets: y", files) {
		t.Error("checkTitle(both entities) failed, want pass")
	}
}

func TestCheckTitleDocsAndTestcaseOnlyRules(t *testing.T) {
	docs := titleFile{name: "website/docs/d/instance.html.markdown", docs: true}
	testcase := titleFile{name: "alicloud/resource_alicloud_instance_test.go", testcase: true}
	if checkTitle("docs: update instance docs", []titleFile{testcase}) {
		t.Error("checkTitle(testcase-only without marker) passed, want fail")
	}
	if !checkTitle("testcase: stabilize instance test", []titleFile{testcase}) {
		t.Error("checkTitle(testcase-only with marker) failed, want pass")
	}
	if checkTitle("testcase: docs change", []titleFile{docs}) {
		t.Error("checkTitle(docs-only without prefix) passed, want fail")
	}
	if !checkTitle("docs: update instance docs", []titleFile{docs}) {
		t.Error("checkTitle(docs-only with prefix) failed, want pass")
	}
	// A docs+testcase change keeps the docs-only branch: the previous
	// gates announced it as docs.
	if !checkTitle("docs: update instance docs and test", []titleFile{docs, testcase}) {
		t.Error("checkTitle(docs and testcase) failed, want pass with the docs prefix")
	}
}

func TestCheckTitleServiceSourcesWithoutEntitiesImposeNoRequirement(t *testing.T) {
	// register.go, shared helpers and the ephemeral/action/list stems
	// clear the docs-only and testcase-only states only.
	files := []titleFile{
		{name: "alicloud/service/ims/register.go", service: true},
		{name: "alicloud/service/ims/helpers.go", service: true},
		{name: "alicloud/service/sts/ephemeral_token.go", service: true},
	}
	if !checkTitle("refactor: reorganize the ims registrations", files) {
		t.Error("checkTitle(service sources without entities) failed, want pass")
	}
}

func TestCheckTitleRequiresServiceEntityInTitle(t *testing.T) {
	files := []titleFile{serviceFile("alicloud/service/ims/resource_user.go",
		entity{typeName: "alicloud_ims_user"})}
	if checkTitle("feat: add the ims user", files) {
		t.Error("checkTitle(new service resource missing) passed, want fail")
	}
	if !checkTitle("feat: New Resource: alicloud_ims_user", files) {
		t.Error("checkTitle(new service resource announced) failed, want pass")
	}
	modified := []titleFile{serviceFile("alicloud/service/ims/resource_user.go",
		entity{typeName: "alicloud_ims_user", existing: true})}
	if !checkTitle("resource/alicloud_ims_user: add a phone attribute", modified) {
		t.Error("checkTitle(modified service resource) failed, want pass")
	}
	dataSource := []titleFile{serviceFile("alicloud/service/ims/data_source_users.go",
		entity{typeName: "alicloud_ims_users", dataSource: true, existing: true})}
	if checkTitle("resource/alicloud_ims_users: wrong map", dataSource) {
		t.Error("checkTitle(data source under resource form) passed, want fail")
	}
	if !checkTitle("data-source/alicloud_ims_users: add a filter", dataSource) {
		t.Error("checkTitle(data source under data-source form) failed, want pass")
	}
}

func TestCheckTitleMigrationAnnouncesTheModifiedForm(t *testing.T) {
	// PR #10817: the implementation files are added, but the entity
	// already existed at the base — the title announces the move with the
	// modified form, never with the new-entity one.
	files := []titleFile{serviceFile("alicloud/service/cloudauth/resource_face_config.go",
		entity{typeName: "alicloud_cloudauth_face_config", existing: true, product: "cloudauth"})}
	if checkTitle("New Resource: alicloud_cloudauth_face_config migrate to split layout", files) {
		t.Error("checkTitle(migration under the new-entity form) passed, want fail")
	}
	if !checkTitle("resource/alicloud_cloudauth_face_config: migrate to alicloud/service/cloudauth split layout", files) {
		t.Error("checkTitle(migration under the modified form) failed, want pass")
	}
}

func TestCheckTitleRejectsEllipsisTitles(t *testing.T) {
	if checkTitle("fix: resource/alicloud_instance: something...", nil) {
		t.Error(`checkTitle("..." suffix) passed, want fail`)
	}
	if checkTitle("fix: resource/alicloud_instance: something…", nil) {
		t.Error(`checkTitle("…" ellipsis) passed, want fail`)
	}
}

func TestParseNameStatus(t *testing.T) {
	out := strings.Join([]string{
		"A\talicloud/service/ims/resource_user.go",
		"M\talicloud/resource_alicloud_vpc.go",
		"R100\talicloud/resource_alicloud_old.go\talicloud/resource_alicloud_new.go",
		"D\talicloud/resource_alicloud_gone.go",
		"",
	}, "\n")
	files := parseNameStatus(out)
	if len(files) != 3 {
		t.Fatalf("parseNameStatus returned %d files, want 3 (deleted dropped)", len(files))
	}
	if files[0].name != "alicloud/service/ims/resource_user.go" || !files[0].added || !files[0].service {
		t.Errorf("parseNameStatus classified the added service file as %+v", files[0])
	}
	if files[1].name != "alicloud/resource_alicloud_vpc.go" || files[1].added || !files[1].flat {
		t.Errorf("parseNameStatus classified the modified flat file as %+v", files[1])
	}
	// A rename classifies by its new name with modified semantics.
	if files[2].name != "alicloud/resource_alicloud_new.go" || files[2].added || !files[2].flat {
		t.Errorf("parseNameStatus classified the renamed file as %+v", files[2])
	}
}

func TestClassify(t *testing.T) {
	cases := map[string][4]bool{
		"alicloud/data_source_alicloud_common.go":                     {false, false, false, true},
		"alicloud/service/cloudauth/data_source_face_configs_test.go": {true, false, false, false},
		"alicloud/service/cloudauth/register.go":                      {false, false, true, false},
		"alicloud/service/cloudauth/data_source_face_configs.go":      {false, false, true, false},
		"website/docs/d/cloudauth_face_configs.html.markdown":         {false, true, false, false},
		"alicloud/provider.go":                                        {false, false, false, false},
		"scripts/internal/splitlayout/splitlayout.go":                 {false, false, false, false},
	}
	for name, want := range cases {
		got := [4]bool{}
		got[0], got[1], got[2], got[3] = classify(name)
		if got != want {
			t.Errorf("classify(%s) = %+v, want %+v", name, got, want)
		}
	}
}

func TestFlatEntities(t *testing.T) {
	cases := map[string][]entity{
		"alicloud/data_source_alicloud_common.go": {{typeName: "alicloud_common", dataSource: true, existing: true}},
		"alicloud/resource_alicloud_instance.go":  {{typeName: "alicloud_instance", existing: true}},
	}
	for path, want := range cases {
		if got := flatEntities(path, false); !reflect.DeepEqual(got, want) {
			t.Errorf("flatEntities(%s, modified) = %+v, want %+v", path, got, want)
		}
	}
	if got := flatEntities("alicloud/resource_alicloud_instance.go", true); !reflect.DeepEqual(got, []entity{{typeName: "alicloud_instance"}}) {
		t.Errorf("flatEntities(added) = %+v, want the non-existing form", got)
	}
}

func TestResolveEntitiesDerivesFlatEntities(t *testing.T) {
	out := "A\talicloud/resource_alicloud_vpc.go\nM\talicloud/data_source_alicloud_common.go\n"
	files := parseNameStatus(out)
	if err := resolveEntities(files, ""); err != nil {
		t.Fatalf("resolveEntities failed: %s", err)
	}
	want := [][]entity{
		{{typeName: "alicloud_vpc"}},
		{{typeName: "alicloud_common", dataSource: true, existing: true}},
	}
	for i, w := range want {
		if !reflect.DeepEqual(files[i].entities, w) {
			t.Errorf("resolveEntities derived %+v for %s, want %+v", files[i].entities, files[i].name, w)
		}
	}
}

// writeFixture builds a minimal split-layout tree — one SDKv2-dialect
// product and one framework-dialect product — and returns its root, so
// the register.go join can be exercised against real parseable sources.
func writeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"alicloud/service/sdkp/register.go": `package sdkp

var ServicePackage = c.ServicePackage{
	SDKResources: []c.SDKResource{
		{TypeName: "alicloud_sdkp_thing", Factory: ResourceThing},
	},
	SDKDataSources: []c.SDKResource{
		{TypeName: "alicloud_sdkp_things", Factory: DataSourceThings},
	},
}
`,
		"alicloud/service/sdkp/resource_thing.go":     "package sdkp\n\nfunc ResourceThing() *r.Resource { return nil }\n",
		"alicloud/service/sdkp/data_source_things.go": "package sdkp\n\nfunc DataSourceThings() *s.DataSource { return nil }\n",
		"alicloud/service/fwp/register.go": `package fwp

var ServicePackage = c.ServicePackage{
	Resources: []c.Resource{
		{TypeName: "alicloud_fwp_user", Factory: NewUserResource},
	},
	DataSources: []c.DataSource{
		{TypeName: "alicloud_fwp_users", Factory: NewUsersDataSource},
	},
	EphemeralResources: []c.EphemeralResource{
		{TypeName: "alicloud_fwp_token", Factory: NewTokenEphemeral},
	},
}
`,
		"alicloud/service/fwp/resource_user.go":     "package fwp\n\nfunc NewUserResource() *x.Resource { return nil }\n",
		"alicloud/service/fwp/data_source_users.go": "package fwp\n\nfunc NewUsersDataSource() *x.DataSource { return nil }\n",
		"alicloud/service/fwp/ephemeral_token.go":   "package fwp\n\nfunc NewTokenEphemeral() *x.Ephemeral { return nil }\n",
		"alicloud/service/fwp/helpers.go":           "package fwp\n\nfunc helper() {}\n",
	}
	for path, src := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// chdir switches the test into dir and restores the previous directory
// afterwards; splitlayout resolves register.go paths against the working
// directory. Tests that chdir must not run in parallel.
func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	})
}

func TestUnitServiceEntitiesJoinsTheRegister(t *testing.T) {
	root := writeFixture(t)
	chdir(t, root)
	empty := &baseIndex{resources: map[string]bool{}, dataSources: map[string]bool{}}
	cases := map[string][]entity{
		"alicloud/service/sdkp/resource_thing.go": {
			{typeName: "alicloud_sdkp_thing", product: "sdkp"},
		},
		"alicloud/service/sdkp/data_source_things.go": {
			{typeName: "alicloud_sdkp_things", dataSource: true, product: "sdkp"},
		},
		"alicloud/service/fwp/resource_user.go": {
			{typeName: "alicloud_fwp_user", product: "fwp"},
		},
		"alicloud/service/fwp/data_source_users.go": {
			{typeName: "alicloud_fwp_users", dataSource: true, product: "fwp"},
		},
	}
	for path, want := range cases {
		if got := serviceEntities(path, empty); !reflect.DeepEqual(got, want) {
			t.Errorf("serviceEntities(%s) = %+v, want %+v", path, got, want)
		}
	}
	// Ephemeral registrations have no title form: they only clear the
	// docs-only and testcase-only states.
	if got := serviceEntities("alicloud/service/fwp/ephemeral_token.go", empty); len(got) != 0 {
		t.Errorf("serviceEntities(ephemeral) = %+v, want none", got)
	}
}

func TestUnitServiceEntitiesMarksExistingEntities(t *testing.T) {
	root := writeFixture(t)
	chdir(t, root)
	index := baseIndexFromSources(nil, map[string][]byte{
		"alicloud/service/sdkp/register.go": []byte(`package sdkp

var ServicePackage = c.ServicePackage{
	SDKResources: []c.SDKResource{
		{TypeName: "alicloud_sdkp_thing", Factory: ResourceThing},
	},
}
`),
	}, nil)
	got := serviceEntities("alicloud/service/sdkp/resource_thing.go", index)
	want := []entity{{typeName: "alicloud_sdkp_thing", existing: true, product: "sdkp"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceEntities with a base index = %+v, want %+v", got, want)
	}
}

func TestUnitServiceEntitiesSkipsNonEntityServiceFiles(t *testing.T) {
	root := writeFixture(t)
	chdir(t, root)
	for _, path := range []string{
		"alicloud/service/sdkp/register.go",
		"alicloud/service/fwp/helpers.go",
		"alicloud/service/fwp/resource_user_test.go",
		"alicloud/service/fwp/unknown_kind.go",
	} {
		if got := serviceEntities(path, &baseIndex{}); len(got) != 0 {
			t.Errorf("serviceEntities(%s) = %+v, want none", path, got)
		}
	}
}

func TestResolveLines(t *testing.T) {
	root := writeFixture(t)
	chdir(t, root)
	got := resolveLines("alicloud/service/fwp/resource_user.go")
	want := []string{"alicloud/service/fwp/resource_user.go\talicloud_fwp_user\tresource"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("resolveLines = %q, want %q", got, want)
	}
	if got := resolveLines("alicloud/service/fwp/helpers.go"); len(got) != 0 {
		t.Errorf("resolveLines(helpers) = %q, want none", got)
	}
}

// quotasProviderSrc returns a root alicloud/provider.go source that registers
// the public Quotas alias alicloud_quotas_application_info in DataSourcesMap
// (Quotas #11288), the form splitlayout.ProviderRegisteredNames reads.
func quotasProviderSrc() []byte {
	return []byte(`package main

var Provider = schema.Provider{
	DataSourcesMap: map[string]*schema.Resource{
		"alicloud_quotas_application_info": QuotasApplicationInfo(),
	},
}
`)
}

func TestUnitBaseIndexFromSources(t *testing.T) {
	index := baseIndexFromSources(
		[]string{
			"alicloud/resource_alicloud_vpc.go",
			"alicloud/data_source_alicloud_vpcs.go",
			"alicloud/provider.go",
			"alicloud/service",
		},
		map[string][]byte{
			"alicloud/service/sdkp/register.go": []byte(`package sdkp

var ServicePackage = c.ServicePackage{
	SDKResources: []c.SDKResource{
		{TypeName: "alicloud_sdkp_thing", Factory: ResourceThing},
	},
}
`),
			"alicloud/service/fwp/register.go": []byte(`package fwp

var ServicePackage = c.ServicePackage{
	Resources: []c.Resource{
		{TypeName: "alicloud_fwp_user", Factory: NewUserResource},
	},
	EphemeralResources: []c.EphemeralResource{
		{TypeName: "alicloud_fwp_token", Factory: NewTokenEphemeral},
	},
}
`),
		},
		quotasProviderSrc(),
	)
	cases := []struct {
		name string
		ds   bool
	}{
		{"alicloud_vpc", false},
		{"alicloud_vpcs", true},
		{"alicloud_sdkp_thing", false},
		{"alicloud_fwp_user", false},
		{"alicloud_quotas_application_info", true},
	}
	for _, c := range cases {
		if !index.has(c.name, c.ds) {
			t.Errorf("baseIndexFromSources lost %s (dataSource=%v)", c.name, c.ds)
		}
	}
	// Kind-aware: the root Quotas alias is a DataSource, so it must not
	// satisfy a same-named Resource lookup — the collision the fix targets.
	if index.has("alicloud_quotas_application_info", false) {
		t.Error("baseIndexFromSources filed the Quotas DataSource alias as a resource")
	}
	// Ephemeral registrations are no longer folded into resources: they
	// live in the others map, so they neither satisfy a resource lookup
	// nor a data-source one, and a same-named head Resource is new.
	if index.has("alicloud_fwp_token", false) {
		t.Error("baseIndexFromSources filed the ephemeral alicloud_fwp_token as a resource")
	}
	if index.has("alicloud_fwp_token", true) {
		t.Error("baseIndexFromSources filed the ephemeral alicloud_fwp_token as a data source")
	}
	if !index.others["alicloud_fwp_token"] {
		t.Error("baseIndexFromSources lost the ephemeral alicloud_fwp_token from others")
	}
	if index.has("alicloud_missing", false) {
		t.Error("baseIndexFromSources invented alicloud_missing")
	}
}

func TestUnitCoreTitleResolution(t *testing.T) {
	t.Chdir(t.TempDir())
	dir := "alicloud/provider/core"
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	register := []byte(`package core
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
 SDKDataSources: []conns.SDKDataSource{{TypeName: "alicloud_account", Factory: DataSourceAccount}},
} }`)
	if err := os.WriteFile(dir+"/register.go", register, 0644); err != nil {
		t.Fatal(err)
	}
	path := dir + "/data_source_account.go"
	if err := os.WriteFile(path, []byte("package core; func DataSourceAccount() {}"), 0644); err != nil {
		t.Fatal(err)
	}
	_, _, split, flat := classify(path)
	if !split || flat {
		t.Fatal("core implementation bypassed entity title checks")
	}
	index := baseIndexFromSources(nil, map[string][]byte{dir + "/register.go": register}, nil)
	entities := serviceEntities(path, index)
	if len(entities) != 1 || entities[0].typeName != "alicloud_account" || !entities[0].dataSource || !entities[0].existing {
		t.Fatalf("core entity classified incorrectly: %+v", entities)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("testcase: unrelated change", files) {
		t.Fatal("core type announcement was not required")
	}
	if !checkTitle("data-source/alicloud_account: migrate to core package", files) {
		t.Fatal("valid core title rejected")
	}
	if got := resolveLines(path); len(got) != 1 || got[0] != path+"\talicloud_account\tdata_source" {
		t.Fatalf("core documentation mapping missing: %v", got)
	}
}

func TestUnitCheckTitleAggregateWholeProductMigration(t *testing.T) {
	// A whole-product migration names too many entities for one title; the
	// concise migrate/<product>: waives every per-entity fragment when all
	// are existing entities of that one product.
	var es []entity
	for i := 0; i < 47; i++ {
		es = append(es, ent(fmt.Sprintf("alicloud_ram_entity_%d", i), false, true, "ram"))
	}
	files := []titleFile{svcEntities("ram", es...)}
	if !checkTitle("migrate/ram: move the whole RAM package to split layout", files) {
		t.Error("checkTitle(whole-product migration) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateWaivesResourceAndDataSource(t *testing.T) {
	files := []titleFile{svcEntities("cbn",
		ent("alicloud_cbn_resource", false, true, "cbn"),
		ent("alicloud_cbn_data", true, true, "cbn"))}
	if !checkTitle("migrate/cbn: move the resource and data source", files) {
		t.Error("checkTitle(resource+datasource migration) failed, want pass")
	}
	// Without the marker both fragments are demanded.
	if checkTitle("refactor: move the resource and data source", files) {
		t.Error("checkTitle(no marker, missing fragments) passed, want fail")
	}
	if !checkTitle("refactor: resource/alicloud_cbn_resource: x and data-source/alicloud_cbn_data: y", files) {
		t.Error("checkTitle(both fragments without marker) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateRootAliasesMigration(t *testing.T) {
	// A root alias (Quotas #11288) moved into a service register is an
	// existing entity of its product, so migrate/<product>: waives it.
	files := []titleFile{svcEntities("quotas",
		ent("alicloud_quotas_application_info", true, true, "quotas"))}
	if !checkTitle("migrate/quotas: move the root alias into a service register", files) {
		t.Error("checkTitle(root alias migration) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateRejectsNewEntity(t *testing.T) {
	// A new entity is never waived: the marker cannot hide an undeclared
	// addition, so the title must still announce New Resource:.
	files := []titleFile{svcEntities("ram",
		ent("alicloud_ram_new", false, false, "ram"))}
	if checkTitle("migrate/ram: move the package", files) {
		t.Error("checkTitle(marker hiding a new entity) passed, want fail")
	}
	if !checkTitle("migrate/ram: move the package New Resource: alicloud_ram_new", files) {
		t.Error("checkTitle(new entity announced beside the marker) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateRejectsCrossProduct(t *testing.T) {
	// A cross-product mix is not one product's migration, so the waiver does
	// not apply and every entity keeps demanding its own fragment — the
	// outsider the marker would have hidden stays declared.
	files := []titleFile{
		svcEntities("ram", ent("alicloud_ram_thing", false, true, "ram")),
		svcEntities("cbn", ent("alicloud_cbn_thing", false, true, "cbn")),
	}
	if checkTitle("migrate/ram: move things", files) {
		t.Error("checkTitle(cross-product behind one marker) passed, want fail")
	}
	if checkTitle("migrate/ram: move things resource/alicloud_cbn_thing: x", files) {
		t.Error("checkTitle(cross-product with only the outsider fragment) passed, want fail")
	}
	if !checkTitle("migrate/ram: move things resource/alicloud_ram_thing: x resource/alicloud_cbn_thing: y", files) {
		t.Error("checkTitle(cross-product with every fragment) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateRejectsWrongProduct(t *testing.T) {
	files := []titleFile{svcEntities("ram", ent("alicloud_ram_thing", false, true, "ram"))}
	if checkTitle("migrate/cbn: wrong product", files) {
		t.Error("checkTitle(wrong product marker) passed, want fail")
	}
	if !checkTitle("migrate/cbn: wrong product resource/alicloud_ram_thing: x", files) {
		t.Error("checkTitle(wrong product with the fragment) failed, want pass")
	}
}

func TestUnitCheckTitleAggregateDoesNotWaiveFlatEntities(t *testing.T) {
	// Flat entities carry no product, so the marker never waives them.
	files := []titleFile{
		svcEntities("ram", ent("alicloud_ram_thing", false, true, "ram")),
		flatFile("alicloud/resource_alicloud_flat.go", false,
			entity{typeName: "alicloud_flat", existing: true}),
	}
	if checkTitle("migrate/ram: move things", files) {
		t.Error("checkTitle(flat entity hidden behind the marker) passed, want fail")
	}
	if !checkTitle("migrate/ram: move things resource/alicloud_flat: x", files) {
		t.Error("checkTitle(flat entity announced beside the marker) failed, want pass")
	}
}

func TestUnitCheckTitleMigrateWordIsNotABypass(t *testing.T) {
	// The bare word "migrate" without the /<product>: marker is ordinary
	// text; the per-entity rules keep applying.
	files := []titleFile{svcEntities("ram", ent("alicloud_ram_thing", false, true, "ram"))}
	if checkTitle("migrate ram: move things", files) {
		t.Error("checkTitle(bare migrate word bypass) passed, want fail")
	}
	if !checkTitle("migrate ram: move things resource/alicloud_ram_thing: x", files) {
		t.Error("checkTitle(bare migrate word with fragment) failed, want pass")
	}
}

func TestUnitAggregateMarkerParsing(t *testing.T) {
	cases := map[string]string{
		"migrate/ram: x":           "ram",
		"migrate/ram-resource: x":  "ram-resource",
		"feat: migrate/cbn: move":  "cbn",
		"migrate ram: x":           "", // space, not a marker
		"resource/alicloud_ram: x": "", // not the migrate form
		"":                         "",
		"migrate/: x":              "", // empty product segment
	}
	for title, want := range cases {
		if got := aggregateMarker(title); got != want {
			t.Errorf("aggregateMarker(%q) = %q, want %q", title, got, want)
		}
	}
}

// writeDemoFixture builds a one-entity demo split-layout tree. When
// headResource is true the head registers alicloud_demo_thing as an SDK
// Resource (resource_thing.go + ResourceThing); otherwise as a Data Source
// (data_source_thing.go + DataSourceThing). The base commit owned no demo
// register, so the base index is built separately from the root provider.go
// alias the caller supplies — the e2e tests drive the real baseIndexFromSources
// + serviceEntities + resolveLines + checkTitle, never hand-setting existing.
// writeDemoTypeFixture builds a one-entity demo split-layout tree that
// registers typeName. When headResource is true the head registers typeName as
// an SDK Resource (resource_thing.go + ResourceThing); otherwise as a Data
// Source (data_source_thing.go + DataSourceThing). The base commit owned no
// demo register, so the base index is built separately from the flat files
// and root provider.go alias the caller supplies — the e2e tests drive the real
// baseIndexFromSources + serviceEntities + resolveLines + checkTitle, never
// hand-setting existing.
func writeDemoTypeFixture(t *testing.T, typeName string, headResource bool) string {
	t.Helper()
	root := t.TempDir()
	var reg, impl, implPath string
	if headResource {
		reg = fmt.Sprintf(`package demo

var ServicePackage = c.ServicePackage{
	SDKResources: []c.SDKResource{
		{TypeName: %q, Factory: ResourceThing},
	},
}
`, typeName)
		impl = "package demo\n\nfunc ResourceThing() *r.Resource { return nil }\n"
		implPath = "alicloud/service/demo/resource_thing.go"
	} else {
		reg = fmt.Sprintf(`package demo

var ServicePackage = c.ServicePackage{
	SDKDataSources: []c.SDKResource{
		{TypeName: %q, Factory: DataSourceThing},
	},
}
`, typeName)
		impl = "package demo\n\nfunc DataSourceThing() *s.DataSource { return nil }\n"
		implPath = "alicloud/service/demo/data_source_thing.go"
	}
	files := map[string]string{
		"alicloud/service/demo/register.go": reg,
		implPath:                            impl,
	}
	for path, src := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func writeDemoFixture(t *testing.T, headResource bool) string {
	return writeDemoTypeFixture(t, "alicloud_demo_thing", headResource)
}

// demoProviderSrc returns a root alicloud/provider.go source that registers
// alicloud_demo_thing in DataSourcesMap (dataSource=true) or ResourcesMap,
// the public alias form splitlayout.ProviderRegisteredNames reads.
func demoProviderSrc(dataSource bool) []byte {
	field := "ResourcesMap"
	if dataSource {
		field = "DataSourcesMap"
	}
	return []byte(fmt.Sprintf(`package main

var Provider = schema.Provider{
	%s: map[string]*schema.Resource{
		"alicloud_demo_thing": Thing(),
	},
}
`, field))
}

// TestUnitCheckTitleKindCollisionForward is the regression the fix targets: the
// base root provider.go registers alicloud_demo_thing as a DataSource, and the
// head adds a same-named demo Resource. A name-only index would mark the new
// Resource existing and let migrate/demo: waive it; the kind-aware index keeps
// the base DataSource from satisfying the Resource, so migrate/demo: fails and
// New Resource: is demanded.
func TestUnitCheckTitleKindCollisionForward(t *testing.T) {
	root := writeDemoFixture(t, true)
	chdir(t, root)
	base := baseIndexFromSources(nil, map[string][]byte{}, demoProviderSrc(true))
	path := "alicloud/service/demo/resource_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_thing", existing: false, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (base DataSource must not satisfy the new Resource)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("migrate/demo: move the demo thing to split layout", files) {
		t.Error("checkTitle(migrate/demo hiding a new same-named Resource) passed, want fail")
	}
	if !checkTitle("feat: New Resource: alicloud_demo_thing", files) {
		t.Error("checkTitle(New Resource for the kind collision) failed, want pass")
	}
	if got := resolveLines(path); len(got) != 1 || got[0] != path+"\talicloud_demo_thing\tresource" {
		t.Errorf("resolveLines = %v, want the resource mapping", got)
	}
}

// TestUnitCheckTitleKindCollisionReverse is the mirror: the base root
// provider.go registers alicloud_demo_thing as a Resource, and the head adds a
// same-named demo Data Source. migrate/demo: fails and New Data Source: is
// demanded.
func TestUnitCheckTitleKindCollisionReverse(t *testing.T) {
	root := writeDemoFixture(t, false)
	chdir(t, root)
	base := baseIndexFromSources(nil, map[string][]byte{}, demoProviderSrc(false))
	path := "alicloud/service/demo/data_source_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_thing", dataSource: true, existing: false, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (base Resource must not satisfy the new Data Source)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("migrate/demo: move the demo thing to split layout", files) {
		t.Error("checkTitle(migrate/demo hiding a new same-named Data Source) passed, want fail")
	}
	if !checkTitle("feat: New Data Source: alicloud_demo_thing", files) {
		t.Error("checkTitle(New Data Source for the kind collision) failed, want pass")
	}
}

// TestUnitCheckTitleSameKindAliasMigrationIsWaived proves the fix does not
// over-block: when the base DataSource alias moves into the demo register as
// the same kind, the entity is existing and migrate/demo: waives it.
func TestUnitCheckTitleSameKindAliasMigrationIsWaived(t *testing.T) {
	root := writeDemoFixture(t, false)
	chdir(t, root)
	base := baseIndexFromSources(nil, map[string][]byte{}, demoProviderSrc(true))
	path := "alicloud/service/demo/data_source_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_thing", dataSource: true, existing: true, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (same-kind alias must be existing)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if !checkTitle("migrate/demo: move the demo data source alias into the service register", files) {
		t.Error("checkTitle(same-kind alias migration) failed, want pass")
	}
}

// demoBaseRegisterSrc returns a base demo register.go source that registers
// alicloud_demo_thing as the given framework-native kind — ephemeral, list
// or action — so the kind-collision fixture can drive a base registration
// that has no house title form through the real baseIndexFromSources path.
func demoBaseRegisterSrc(kind splitlayout.Kind) []byte {
	var field, elem string
	switch kind {
	case splitlayout.KindEphemeral:
		field, elem = "EphemeralResources", "c.EphemeralResource"
	case splitlayout.KindListResource:
		field, elem = "ListResources", "c.ListResource"
	case splitlayout.KindAction:
		field, elem = "Actions", "c.Action"
	default:
		panic(fmt.Sprintf("demoBaseRegisterSrc: unsupported kind %s", kind))
	}
	return []byte(fmt.Sprintf(`package demo

var ServicePackage = c.ServicePackage{
	%s: []%s{
		{TypeName: "alicloud_demo_thing", Factory: DemoThing},
	},
}
`, field, elem))
}

// TestUnitCheckTitleNonTitleKindBaseDoesNotSatisfyResource is the regression the
// kind-aware others map targets: the base register.go registers
// alicloud_demo_thing only as an ephemeral, list resource or action — kinds
// with no house title form — and the head adds a same-named SDK Resource. A
// name-only (or resource-folding) index would mark the new Resource existing
// and let migrate/demo: waive it; the kind-aware index keeps the base entity
// in others, so it never satisfies the Resource lookup, migrate/demo: fails
// and New Resource: is demanded. The entities are derived through the real
// baseIndexFromSources + serviceEntities + resolveLines + checkTitle path,
// never hand-setting existing.
func TestUnitCheckTitleNonTitleKindBaseDoesNotSatisfyResource(t *testing.T) {
	for _, kind := range []splitlayout.Kind{splitlayout.KindEphemeral, splitlayout.KindListResource, splitlayout.KindAction} {
		t.Run(string(kind), func(t *testing.T) {
			root := writeDemoFixture(t, true)
			chdir(t, root)
			base := baseIndexFromSources(nil, map[string][]byte{
				"alicloud/service/demo/register.go": demoBaseRegisterSrc(kind),
			}, nil)
			if base.has("alicloud_demo_thing", false) {
				t.Fatalf("base %s alicloud_demo_thing satisfies a resource lookup", kind)
			}
			if !base.others["alicloud_demo_thing"] {
				t.Fatalf("base %s alicloud_demo_thing not retained in others", kind)
			}
			path := "alicloud/service/demo/resource_thing.go"
			entities := serviceEntities(path, base)
			want := []entity{{typeName: "alicloud_demo_thing", existing: false, product: "demo"}}
			if !reflect.DeepEqual(entities, want) {
				t.Fatalf("kind=%s: serviceEntities = %+v, want %+v (base non-title kind must not satisfy the new Resource)", kind, entities, want)
			}
			files := []titleFile{serviceFile(path, entities...)}
			if checkTitle("migrate/demo: move the demo thing to split layout", files) {
				t.Errorf("kind=%s: migrate/demo waived a new same-named Resource", kind)
			}
			if !checkTitle("feat: New Resource: alicloud_demo_thing", files) {
				t.Errorf("kind=%s: New Resource for the non-title-kind collision rejected", kind)
			}
			if got := resolveLines(path); len(got) != 1 || got[0] != path+"\talicloud_demo_thing\tresource" {
				t.Errorf("kind=%s: resolveLines = %v, want the resource mapping", kind, got)
			}
		})
	}
}

// TestUnitCheckTitleFlatTestFileDoesNotSatisfyResource is the regression the
// flat _test.go filter targets: the base flat tree lists
// resource_alicloud_demo_test.go, a test file that names no provider entity,
// and the head registers alicloud_demo_test (a legitimate _test-suffixed type
// name) as a demo Resource. A stem-filing index would record
// alicloud_demo_test as a base resource and let migrate/demo: waive the new
// entity; the fix skips _test.go stems, so the base knows nothing, the entity
// is new, migrate/demo: fails and New Resource: is demanded. The entities are
// derived through the real baseIndexFromSources + serviceEntities +
// resolveLines + checkTitle path, never hand-setting existing.
func TestUnitCheckTitleFlatTestFileDoesNotSatisfyResource(t *testing.T) {
	root := writeDemoTypeFixture(t, "alicloud_demo_test", true)
	chdir(t, root)
	base := baseIndexFromSources(
		[]string{"alicloud/resource_alicloud_demo_test.go"},
		map[string][]byte{}, nil)
	if base.has("alicloud_demo_test", false) {
		t.Fatal("base flat _test.go satisfied a resource lookup")
	}
	path := "alicloud/service/demo/resource_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_test", existing: false, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (base _test.go must not satisfy the new Resource)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("migrate/demo: move the demo thing to split layout", files) {
		t.Error("checkTitle(migrate/demo hiding a new Resource behind a base _test.go) passed, want fail")
	}
	if !checkTitle("feat: New Resource: alicloud_demo_test", files) {
		t.Error("checkTitle(New Resource for the base _test.go collision) failed, want pass")
	}
	if got := resolveLines(path); len(got) != 1 || got[0] != path+"\talicloud_demo_test\tresource" {
		t.Errorf("resolveLines = %v, want the resource mapping", got)
	}
}

// TestUnitCheckTitleFlatTestFileDoesNotSatisfyDataSource is the data-source
// mirror: the base flat tree lists data_source_alicloud_demo_test.go and the
// head registers alicloud_demo_test as a demo Data Source. A stem-filing index
// would record it as a base data source and waive it; the fix skips _test.go,
// so the entity is new, migrate/demo: fails and New Data Source: is demanded.
func TestUnitCheckTitleFlatTestFileDoesNotSatisfyDataSource(t *testing.T) {
	root := writeDemoTypeFixture(t, "alicloud_demo_test", false)
	chdir(t, root)
	base := baseIndexFromSources(
		[]string{"alicloud/data_source_alicloud_demo_test.go"},
		map[string][]byte{}, nil)
	if base.has("alicloud_demo_test", true) {
		t.Fatal("base flat _test.go satisfied a data-source lookup")
	}
	path := "alicloud/service/demo/data_source_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_test", dataSource: true, existing: false, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (base _test.go must not satisfy the new Data Source)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("migrate/demo: move the demo thing to split layout", files) {
		t.Error("checkTitle(migrate/demo hiding a new Data Source behind a base _test.go) passed, want fail")
	}
	if !checkTitle("feat: New Data Source: alicloud_demo_test", files) {
		t.Error("checkTitle(New Data Source for the base _test.go collision) failed, want pass")
	}
}

// TestUnitCheckTitleFlatHelperDoesNotSatisfyResource is the regression the
// non-implementation stem filter targets: the base flat tree lists
// alicloud_demo_thing.go, a shared helper with no resource_/data_source_
// prefix that names no provider entity, and the head registers
// alicloud_demo_thing as a demo Resource. A stem-filing index would record
// alicloud_demo_thing as a base resource and let migrate/demo: waive the new
// entity; the fix skips non-implementation stems, so the entity is new,
// migrate/demo: fails and New Resource: is demanded.
func TestUnitCheckTitleFlatHelperDoesNotSatisfyResource(t *testing.T) {
	root := writeDemoTypeFixture(t, "alicloud_demo_thing", true)
	chdir(t, root)
	base := baseIndexFromSources(
		[]string{"alicloud/alicloud_demo_thing.go"},
		map[string][]byte{}, nil)
	if base.has("alicloud_demo_thing", false) {
		t.Fatal("base flat helper satisfied a resource lookup")
	}
	path := "alicloud/service/demo/resource_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_thing", existing: false, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (base helper must not satisfy the new Resource)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if checkTitle("migrate/demo: move the demo thing to split layout", files) {
		t.Error("checkTitle(migrate/demo hiding a new Resource behind a base helper) passed, want fail")
	}
	if !checkTitle("feat: New Resource: alicloud_demo_thing", files) {
		t.Error("checkTitle(New Resource for the base helper collision) failed, want pass")
	}
}

// TestUnitCheckTitleFlatImplementationStillSatisfies proves the fix does not
// over-block: a real flat implementation file resource_alicloud_demo_thing.go
// keeps filing alicloud_demo_thing as a base resource, so when the head moves
// it into the demo register as the same kind, the entity is existing and
// migrate/demo: waives the per-entity fragment. The flat implementation prefix
// is admitted exactly as before; only _test.go and non-implementation stems
// are dropped.
func TestUnitCheckTitleFlatImplementationStillSatisfies(t *testing.T) {
	root := writeDemoTypeFixture(t, "alicloud_demo_thing", true)
	chdir(t, root)
	base := baseIndexFromSources(
		[]string{"alicloud/resource_alicloud_demo_thing.go"},
		map[string][]byte{}, nil)
	if !base.has("alicloud_demo_thing", false) {
		t.Fatal("base flat implementation lost the real resource")
	}
	path := "alicloud/service/demo/resource_thing.go"
	entities := serviceEntities(path, base)
	want := []entity{{typeName: "alicloud_demo_thing", existing: true, product: "demo"}}
	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("serviceEntities = %+v, want %+v (real flat implementation must satisfy the migration)", entities, want)
	}
	files := []titleFile{serviceFile(path, entities...)}
	if !checkTitle("migrate/demo: move the demo thing into the service register", files) {
		t.Error("checkTitle(real flat implementation migration) failed, want pass")
	}
	if checkTitle("feat: New Resource: alicloud_demo_thing", files) {
		t.Error("checkTitle(existing flat implementation demanded New Resource) passed, want fail")
	}
}
