package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
		entity{typeName: "alicloud_cloudauth_face_config", existing: true})}
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

func TestServiceEntitiesJoinsTheRegister(t *testing.T) {
	root := writeFixture(t)
	chdir(t, root)
	empty := &baseIndex{names: map[string]bool{}}
	cases := map[string][]entity{
		"alicloud/service/sdkp/resource_thing.go": {
			{typeName: "alicloud_sdkp_thing"},
		},
		"alicloud/service/sdkp/data_source_things.go": {
			{typeName: "alicloud_sdkp_things", dataSource: true},
		},
		"alicloud/service/fwp/resource_user.go": {
			{typeName: "alicloud_fwp_user"},
		},
		"alicloud/service/fwp/data_source_users.go": {
			{typeName: "alicloud_fwp_users", dataSource: true},
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

func TestServiceEntitiesMarksExistingEntities(t *testing.T) {
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
	})
	got := serviceEntities("alicloud/service/sdkp/resource_thing.go", index)
	want := []entity{{typeName: "alicloud_sdkp_thing", existing: true}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceEntities with a base index = %+v, want %+v", got, want)
	}
}

func TestServiceEntitiesSkipsNonEntityServiceFiles(t *testing.T) {
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

func TestBaseIndexFromSources(t *testing.T) {
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
	)
	for _, name := range []string{"alicloud_vpc", "alicloud_vpcs", "alicloud_sdkp_thing", "alicloud_fwp_user", "alicloud_fwp_token"} {
		if !index.has(name) {
			t.Errorf("baseIndexFromSources lost %s", name)
		}
	}
	if index.has("alicloud_missing") {
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
	index := baseIndexFromSources(nil, map[string][]byte{dir + "/register.go": register})
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
