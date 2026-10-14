package main

import (
	"os"
	"path/filepath"
	"reflect"
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
