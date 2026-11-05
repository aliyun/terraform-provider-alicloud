package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSweepPackageUsesLiteralRegisteredOwner(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "alicloud", "service", "example")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	source := `package example_test
 import r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
 // r.AddTestSweepers("alicloud_fake", nil)
 func init(){r.AddTestSweepers("alicloud_example",nil)}
 `
	if err := os.WriteFile(filepath.Join(dir, "resource_test.go"), []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := sweepPackages(root, "alicloud_example")
	if err != nil || !reflect.DeepEqual(got, []string{"./alicloud/service/example"}) {
		t.Fatalf("%v %v", got, err)
	}
	got, err = sweepPackages(root, "alicloud_fake")
	if err != nil || len(got) != 0 {
		t.Fatalf("comment treated as registration: %v %v", got, err)
	}
}
func TestSweepPackageRejectsMissingOrAmbiguousOwner(t *testing.T) {
	for _, ps := range [][]string{nil, {"./alicloud/sweep", "./alicloud/service/example"}} {
		if _, err := selectSweepPackage(ps, ""); err == nil {
			t.Fatalf("accepted ambiguous or missing owner: %v", ps)
		}
	}
	ps := []string{"./alicloud/sweep", "./alicloud/service/example"}
	if got, err := selectSweepPackage(ps, "./alicloud/sweep"); err != nil || got != "./alicloud/sweep" {
		t.Fatalf("%s %v", got, err)
	}
	if _, err := selectSweepPackage(ps, "./alicloud"); err == nil {
		t.Fatal("accepted unrelated override")
	}
}
