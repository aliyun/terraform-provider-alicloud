package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func unitFixture(t *testing.T, files map[string]string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err = os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(old) })
	for p, s := range files {
		if err = os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(p, []byte(s), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

const unitSample = "package p\nimport \"testing\"\nfunc TestUnitA(t *testing.T) {}\nfunc TestAccCloud(t *testing.T) {}\n// func TestUnitComment(t *testing.T) {}\nvar text = `func TestUnitString(t *testing.T) {}`\n"

func TestUnitChangedSelectorCompanionAndEditedTests(t *testing.T) {
	unitFixture(t, map[string]string{"p/resource_alicloud_foo_test.go": unitSample, "p/other_test.go": "package p\nfunc TestUnitB(t *testing.T) {}"})
	cases, broad, err := selectCases([]string{"p/resource_alicloud_foo.go", "p/resource_alicloud_foo_test.go"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 || cases[0].Name != "TestUnitA" || len(broad) != 0 {
		t.Fatalf("%+v %v", cases, broad)
	}
}
func TestUnitChangedSelectorHelperAndServiceBroadenPackage(t *testing.T) {
	for _, name := range []string{"common.go", "service_cs.go", "cs_service.go", "diff_suppress_funcs.go"} {
		t.Run(name, func(t *testing.T) {
			unitFixture(t, map[string]string{"p/foo_test.go": unitSample, "p/bar_test.go": "package p\nfunc TestUnitB(t *testing.T) {}"})
			cases, broad, err := selectCases([]string{"p/" + name})
			if err != nil {
				t.Fatal(err)
			}
			if len(cases) != 2 || !reflect.DeepEqual(broad, []string{"./p"}) {
				t.Fatalf("%+v %v", cases, broad)
			}
		})
	}
}
func TestUnitChangedSelectorDeletedAndRenamedTests(t *testing.T) {
	unitFixture(t, map[string]string{"p/new_test.go": unitSample})
	cases, broad, err := selectCases([]string{"p/old_test.go", "p/new_test.go"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 || len(broad) != 1 {
		t.Fatalf("%+v %v", cases, broad)
	}
}
func TestUnitChangedSelectorBuildConstraints(t *testing.T) {
	unitFixture(t, map[string]string{"p/resource_x_test.go": "//go:build impossible_ci_unit_tag\n\n" + unitSample})
	cases, _, err := selectCases([]string{"p/resource_x.go"})
	if err != nil || len(cases) != 0 {
		t.Fatalf("%+v %v", cases, err)
	}
}
func TestUnitChangedSelectorDocsAndAbsentTests(t *testing.T) {
	unitFixture(t, map[string]string{})
	cases, _, err := selectCases([]string{"README.md", "p/resource_x.go"})
	if err != nil || len(cases) != 0 {
		t.Fatalf("%+v %v", cases, err)
	}
}
func TestUnitChangedSelectorModuleChanges(t *testing.T) {
	unitFixture(t, map[string]string{"a/foo_test.go": unitSample, "b/foo_test.go": unitSample})
	cases, broad, err := selectCases([]string{"go.sum"})
	if err != nil || len(cases) != 2 || len(broad) != 2 {
		t.Fatalf("%+v %v %v", cases, broad, err)
	}
}
func TestUnitChangedSelectorMalformedTestFails(t *testing.T) {
	unitFixture(t, map[string]string{"p/resource_x_test.go": "package p\nfunc TestUnitBroken("})
	_, _, err := selectCases([]string{"p/resource_x.go"})
	if err == nil {
		t.Fatal("invalid Go must not become no applicable tests")
	}
}

func TestUnitChangedRunnerUsesExactPackageAndNoAcc(t *testing.T) {
	unitFixture(t, map[string]string{"bin/go": "#!/bin/sh\ntest -z \"$TF_ACC\" || exit 7\ncase \"$*\" in *'-count=1 -p=1 ./p -run=^(TestUnitFixture)$ -timeout=10m'*) ;; *) exit 8 ;; esac\necho '{\"Action\":\"pass\",\"Test\":\"TestUnitFixture\"}'\n"})
	os.Chmod("bin/go", 0755)
	root, _ := os.Getwd()
	t.Setenv("PATH", filepath.Join(root, "bin")+":"+os.Getenv("PATH"))
	t.Setenv("TF_ACC", "1")
	cases := []unitCase{{"./p", "TestUnitFixture", "MISSING"}}
	if err := runCases(cases); err != nil || cases[0].Result != "PASS" {
		t.Fatalf("%+v %v", cases, err)
	}
}
func TestUnitChangedRunnerPreservesFailureSkipAndMissing(t *testing.T) {
	for _, action := range []string{"fail", "skip", "", "subtest-skip"} {
		t.Run(action, func(t *testing.T) {
			text := "#!/bin/sh\necho '{\"Action\":\"" + action + "\",\"Test\":\"TestUnitFixture\"}'\n"
			if action == "fail" {
				text += "exit 1\n"
			}
			if action == "subtest-skip" {
				text = "#!/bin/sh\necho '{\"Action\":\"skip\",\"Test\":\"TestUnitFixture/optional\"}'\necho '{\"Action\":\"pass\",\"Test\":\"TestUnitFixture\"}'\n"
			}
			unitFixture(t, map[string]string{"bin/go": text})
			os.Chmod("bin/go", 0755)
			root, _ := os.Getwd()
			t.Setenv("PATH", filepath.Join(root, "bin")+":"+os.Getenv("PATH"))
			cases := []unitCase{{"./p", "TestUnitFixture", "MISSING"}}
			err := runCases(cases)
			want := map[string]string{"fail": "FAIL", "skip": "SKIP", "": "MISSING", "subtest-skip": "SKIP"}[action]
			if cases[0].Result != want || (action == "fail" && err == nil) {
				t.Fatalf("%+v %v", cases, err)
			}
		})
	}
}

func TestUnitChangedSummaryAllowsSkipWithoutHidingErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		results []string
		err     error
		want    string
		skipped int
	}{
		{"pass", []string{"PASS"}, nil, "passed", 0},
		{"mixed", []string{"PASS", "SKIP"}, nil, "passed", 1},
		{"all skipped", []string{"SKIP", "SKIP"}, nil, "passed", 2},
		{"failure wins", []string{"SKIP", "FAIL"}, nil, "failed", 1},
		{"missing remains invalid", []string{"PASS", "MISSING"}, nil, "failed", 0},
		{"execution error wins", []string{"SKIP"}, errors.New("compiler failed"), "failed", 1},
		{"unknown outcome", []string{"unexpected"}, nil, "failed", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := unitSummary{Status: "failed"}
			for _, result := range tc.results {
				s.Cases = append(s.Cases, unitCase{Package: "./p", Name: "TestUnitFixture", Result: result})
			}
			finishUnitSummary(&s, tc.err)
			if s.Status != tc.want || s.Skipped != tc.skipped {
				t.Fatalf("unexpected summary: %+v", s)
			}
			for i, result := range tc.results {
				if s.Cases[i].Result != result {
					t.Fatal("case outcome rewritten")
				}
			}
		})
	}
}

const registrationFixture = `package alicloud
func Provider() *schema.Provider {
 return &schema.Provider{
  ResourcesMap: map[string]*schema.Resource{RESOURCE_ENTRIES},
  DataSourcesMap: map[string]*schema.Resource{DATA_ENTRIES},
  ConfigureFunc: configure,
 }
}
`

func registrationSource(resources, data string) string {
	return strings.NewReplacer("RESOURCE_ENTRIES", resources, "DATA_ENTRIES", data).Replace(registrationFixture)
}

func TestUnitChangedRegistrationProof(t *testing.T) {
	base := registrationSource(`"alicloud_old": resourceOld(),`, `"alicloud_data": dataSourceOld(),`)
	for _, tc := range []struct {
		name, after string
		ok          bool
		count       int
	}{
		{"add resource", registrationSource(`"alicloud_old": resourceOld(), "alicloud_new": resourceNew(),`, `"alicloud_data": dataSourceOld(),`), true, 1},
		{"add data source", registrationSource(`"alicloud_old": resourceOld(),`, `"alicloud_data": dataSourceOld(), "alicloud_new": dataSourceNew(),`), true, 1},
		{"remove both", registrationSource("", ""), true, 2},
		{"rename alias", registrationSource(`"alicloud_alias": resourceOld(),`, `"alicloud_data": dataSourceOld(),`), true, 2},
		{"replace factory", strings.Replace(base, "resourceOld()", "resourceNew()", 1), false, 0},
		{"change configure", strings.Replace(base, "ConfigureFunc: configure", "ConfigureFunc: configureNew", 1), false, 0},
		{"other declaration", base + "\nvar shared = 1\n", false, 0},
		{"import change", strings.Replace(base, "package alicloud", "package alicloud\nimport x \"example.com/x\"", 1), false, 0},
		{"dynamic map", strings.Replace(base, `map[string]*schema.Resource{"alicloud_old": resourceOld(),}`, "resources()", 1), false, 0},
		{"duplicate key", registrationSource(`"alicloud_old": resourceOld(), "alicloud_old": resourceOld(),`, `"alicloud_data": dataSourceOld(),`), false, 0},
		{"dynamic key", registrationSource(`key: resourceOld(),`, `"alicloud_data": dataSourceOld(),`), false, 0},
		{"missing map", strings.Replace(base, "DataSourcesMap:", "SomethingElse:", 1), false, 0},
		{"second provider literal", strings.Replace(base, "return &schema.Provider", "_ = schema.Provider{}; return &schema.Provider", 1), false, 0},
		{"build directive", "//go:build custom\n\n" + base, false, 0},
		{"syntax error", "package alicloud\nfunc Provider(", false, 0},
		{"format only", strings.ReplaceAll(base, ": ", ":    "), true, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changed, ok := registrationChanges([]byte(base), []byte(tc.after))
			if ok != tc.ok || (ok && len(changed) != tc.count) {
				t.Fatalf("ok=%v changes=%d", ok, len(changed))
			}
		})
	}
	if _, ok := registrationChanges([]byte("invalid"), []byte(base)); ok {
		t.Fatal("invalid base narrowed")
	}
}

func registrationGitBase(t *testing.T) string {
	t.Helper()
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=Unit Fixture", "-c", "user.email=fixture@example.invalid", "-c", "core.hooksPath=/dev/null", "commit", "-qm", "fixture"}} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git: %v %s", err, out)
		}
	}
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func TestUnitChangedRegistrationSelection(t *testing.T) {
	before := registrationSource(`"alicloud_old": resourceOld(),`, "")
	for _, tc := range []struct {
		name, after string
		extra       []string
		remove      string
		broad       bool
		names       []string
	}{
		{name: "add resource and data source", after: registrationSource(`"alicloud_old": resourceOld(), "alicloud_new": resourceNew(),`, `"alicloud_data": dataSourceNew(),`), names: []string{"TestUnitData", "TestUnitNew"}},
		{name: "remove registration", after: registrationSource("", ""), names: []string{"TestUnitOld"}},
		{name: "alias uses actual factory", after: registrationSource(`"alicloud_old": resourceOld(), "alicloud_alias": resourceNew(),`, ""), names: []string{"TestUnitNew"}},
		{name: "helper still broadens", after: registrationSource("", ""), extra: []string{"alicloud/common.go"}, broad: true},
		{name: "module still broadens", after: registrationSource("", ""), extra: []string{"go.mod"}, broad: true},
		{name: "replacement broadens", after: registrationSource(`"alicloud_old": resourceNew(),`, ""), broad: true},
		{name: "unknown factory broadens", after: registrationSource(`"alicloud_new": missing(),`, ""), broad: true},
		{name: "qualified factory broadens", after: registrationSource(`"alicloud_new": ecs.Resource(),`, ""), broad: true},
		{name: "arguments broaden", after: registrationSource(`"alicloud_new": resourceNew(1),`, ""), broad: true},
		{name: "missing companion broadens", after: registrationSource("", ""), remove: "alicloud/resource_alicloud_old_test.go", broad: true},
		{name: "deleted implementation broadens", after: registrationSource("", ""), remove: "alicloud/resource_alicloud_old.go", broad: true},
		{name: "no unit companion is explicit no match", after: registrationSource(`"alicloud_old": resourceOld(),`, `"alicloud_acc": dataSourceAcc(),`), names: []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unitFixture(t, map[string]string{
				"alicloud/provider.go":                       before,
				"alicloud/resource_alicloud_old.go":          "package alicloud\nfunc resourceOld() {}",
				"alicloud/resource_alicloud_old_test.go":     "package alicloud\nfunc TestUnitOld(t *testing.T) {}",
				"alicloud/resource_alicloud_new.go":          "package alicloud\nfunc resourceNew() {}",
				"alicloud/resource_alicloud_new_test.go":     "package alicloud\nfunc TestUnitNew(t *testing.T) {}",
				"alicloud/data_source_alicloud_data.go":      "package alicloud\nfunc dataSourceNew() {}",
				"alicloud/data_source_alicloud_data_test.go": "package alicloud\nfunc TestUnitData(t *testing.T) {}",
				"alicloud/data_source_alicloud_acc.go":       "package alicloud\nfunc dataSourceAcc() {}",
				"alicloud/data_source_alicloud_acc_test.go":  "package alicloud\nfunc TestAccOnly(t *testing.T) {}",
				"alicloud/common_test.go":                    "package alicloud\nfunc TestUnitShared(t *testing.T) {}",
			})
			base := registrationGitBase(t)
			if err := os.WriteFile("alicloud/provider.go", []byte(tc.after), 0644); err != nil {
				t.Fatal(err)
			}
			if tc.remove != "" {
				if err := os.Remove(tc.remove); err != nil {
					t.Fatal(err)
				}
			}
			files := append([]string{"alicloud/provider.go"}, tc.extra...)
			cases, broad, err := selectCases(files, base)
			if err != nil {
				t.Fatal(err)
			}
			if tc.broad {
				if !reflect.DeepEqual(broad, []string{"./alicloud"}) {
					t.Fatalf("missing fallback: %v", broad)
				}
				return
			}
			names := []string{}
			for _, c := range cases {
				names = append(names, c.Name)
			}
			if len(broad) != 0 || !reflect.DeepEqual(names, tc.names) {
				t.Fatalf("cases=%v broad=%v", names, broad)
			}
		})
	}
}

func TestUnitChangedRegistrationMissingBaseBroadens(t *testing.T) {
	unitFixture(t, map[string]string{"alicloud/provider.go": registrationSource("", ""), "alicloud/common_test.go": unitSample})
	for _, base := range []string{"", "nonexistent"} {
		cases, broad, err := selectCases([]string{"alicloud/provider.go"}, base)
		if err != nil || len(cases) != 1 || len(broad) != 1 {
			t.Fatalf("%v %v %v", cases, broad, err)
		}
	}
}
