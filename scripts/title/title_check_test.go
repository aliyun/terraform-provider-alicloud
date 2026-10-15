package main

import (
	"reflect"
	"testing"

	"github.com/waigani/diffparser"
)

// flatFile builds a flat titleFile with explicit entities, the way
// resolveTitleFiles produces them from a successful join.
func flatFile(name string, added bool, entities ...entity) titleFile {
	return titleFile{name: name, added: added, flat: true, entities: entities}
}

func TestCheckTitleSharedPlumbingImposesNoRequirement(t *testing.T) {
	files := []titleFile{flatFile("alicloud/data_source_alicloud_common.go", false)}
	if !checkTitle("refactor: migrate alicloud_cloudauth_face_configs data source", files) {
		t.Error("checkTitle(shared plumbing) failed, want pass: the file implements no registered entity")
	}
}

func TestCheckTitleRequiresModifiedEntityInTitle(t *testing.T) {
	files := []titleFile{flatFile("alicloud/resource_alicloud_instance.go", false, entity{typeName: "alicloud_instance"})}
	if checkTitle("fix: something else", files) {
		t.Error("checkTitle(modified entity missing from title) passed, want fail")
	}
	if !checkTitle("fix: resource/alicloud_instance: fix the instance bootstrap", files) {
		t.Error("checkTitle(modified entity in title) failed, want pass")
	}
}

func TestCheckTitleModeBRequiresTheRegisteredName(t *testing.T) {
	files := []titleFile{flatFile("alicloud/resource_alicloud_edas_application_package_attachment.go", false,
		entity{typeName: "alicloud_edas_application_deployment"})}
	if checkTitle("fix: resource/edas_application_package_attachment: wrong name", files) {
		t.Error("checkTitle(mode B filename name) passed, want fail: the registered TypeName is required")
	}
	if !checkTitle("fix: resource/alicloud_edas_application_deployment: right name", files) {
		t.Error("checkTitle(mode B registered name) failed, want pass")
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

func TestCheckTitleMultiEntityFileRequiresEveryEntity(t *testing.T) {
	files := []titleFile{flatFile("alicloud/data_source_alicloud_instance_types.go", false,
		entity{typeName: "alicloud_instance_types", dataSource: true},
		entity{typeName: "alicloud_subnets", dataSource: true})}
	if checkTitle("fix: data-source/alicloud_instance_types: one of two", files) {
		t.Error("checkTitle(one of two entities) passed, want fail")
	}
	if !checkTitle("fix: data-source/alicloud_instance_types: x and data-source/alicloud_subnets: y", files) {
		t.Error("checkTitle(both entities) failed, want pass")
	}
}

func TestCheckTitleUnregisteredFlatFileImposesNoRequirement(t *testing.T) {
	// e.g. resource_alicloud_havip.go: a stock file whose constructor is
	// registered nowhere.
	files := []titleFile{flatFile("alicloud/resource_alicloud_havip.go", false)}
	if !checkTitle("fix: anything", files) {
		t.Error("checkTitle(unregistered constructor file) failed, want pass")
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

func TestCheckTitleServiceSourcesImposeNoRequirement(t *testing.T) {
	service := titleFile{name: "alicloud/service/cloudauth/data_source_face_configs.go", service: true}
	testcase := titleFile{name: "alicloud/service/cloudauth/data_source_face_configs_test.go", testcase: true}
	if !checkTitle("refactor: migrate the cloudauth face configs data source", []titleFile{service, testcase}) {
		t.Error("checkTitle(service sources) failed, want pass: they impose no title entries")
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

func TestResolveTitleFilesSkipsDeletedAndClassifies(t *testing.T) {
	files := []*diffparser.DiffFile{
		{Mode: diffparser.DELETED, NewName: "alicloud/data_source_alicloud_cloudauth_face_configs.go"},
		{Mode: diffparser.MODIFIED, NewName: "alicloud/data_source_alicloud_common.go"},
		{Mode: diffparser.NEW, NewName: "alicloud/service/cloudauth/data_source_face_configs_test.go"},
		{Mode: diffparser.MODIFIED, NewName: "alicloud/service/cloudauth/register.go"},
		{Mode: diffparser.MODIFIED, NewName: "website/docs/d/cloudauth_face_configs.html.markdown"},
		{Mode: diffparser.MODIFIED, NewName: "alicloud/provider.go"},
		{Mode: diffparser.MODIFIED, NewName: "scripts/internal/flatlayout/flatlayout.go"},
	}
	resolved := resolveTitleFiles(files, nil)
	if len(resolved) != len(files)-1 {
		t.Fatalf("resolveTitleFiles returned %d files, want %d (deleted dropped)", len(resolved), len(files)-1)
	}
	// testcase, docs, service, flat per file.
	wantFlags := map[string][4]bool{
		"alicloud/data_source_alicloud_common.go":                     {false, false, false, true},
		"alicloud/service/cloudauth/data_source_face_configs_test.go": {true, false, false, false},
		"alicloud/service/cloudauth/register.go":                      {false, false, true, false},
		"website/docs/d/cloudauth_face_configs.html.markdown":         {false, true, false, false},
		"alicloud/provider.go":                                        {false, false, false, false},
		"scripts/internal/flatlayout/flatlayout.go":                   {false, false, false, false},
	}
	addedSeen := false
	for _, f := range resolved {
		want, ok := wantFlags[f.name]
		if !ok {
			t.Fatalf("resolveTitleFiles returned unexpected file %s", f.name)
		}
		got := [4]bool{f.testcase, f.docs, f.service, f.flat}
		if got != want {
			t.Errorf("resolveTitleFiles classified %s as %+v, want %+v", f.name, got, want)
		}
		if f.name == "alicloud/service/cloudauth/data_source_face_configs_test.go" {
			addedSeen = f.added
		}
	}
	if !addedSeen {
		t.Error("resolveTitleFiles did not mark the NEW service test file as added")
	}
}

func TestFlatEntitiesFallsBackToFilenameDerivation(t *testing.T) {
	// The join reads repository files from the working directory; a
	// package test never runs at the repo root, so the read fails and
	// the fallback must yield the name the previous gates required.
	cases := map[string][]entity{
		"alicloud/data_source_alicloud_common.go": {{typeName: "alicloud_common", dataSource: true}},
		"alicloud/resource_alicloud_instance.go":  {{typeName: "alicloud_instance"}},
	}
	for path, want := range cases {
		if got := flatEntities(path, nil); !reflect.DeepEqual(got, want) {
			t.Errorf("flatEntities(%s) = %+v, want %+v", path, got, want)
		}
	}
}
