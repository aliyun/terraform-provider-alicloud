package splitlayout

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnitCoreLayoutResolution(t *testing.T) {
	t.Chdir(t.TempDir())
	files := map[string]string{
		"alicloud/provider/core/register.go": `package core
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
 SDKDataSources: []conns.SDKDataSource{
  {TypeName: "alicloud_account", Factory: DataSourceAccount},
  {TypeName: "alicloud_account_alias", Factory: DataSourceAccount},
 },
 DataSources: []conns.DataSource{{TypeName: "alicloud_core_native", Factory: NewNative}},
} }`,
		"alicloud/provider/core/data_source_account.go": "package core; func DataSourceAccount() {}",
		"alicloud/provider/core/data_source_native.go":  "package core; func NewNative() {}",
		"alicloud/service/core/register.go": `package core
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
 SDKDataSources: []conns.SDKDataSource{{TypeName: "alicloud_product_core", Factory: DataSourceAccount}},
} }`,
		"alicloud/service/core/data_source_account.go": "package core; func DataSourceAccount() {}",
	}
	for path, src := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
	key, stem, isTest, ok := IsSplitSource("alicloud/provider/core/data_source_account_test.go")
	if !ok || key != "provider/core" || stem != "data_source_account" || !isTest {
		t.Fatalf("core path was not recognized: %q %q %v %v", key, stem, isTest, ok)
	}
	resolved, err := Resolve(key, stem)
	if err != nil || len(resolved) != 2 {
		t.Fatalf("core aliases: %+v, %v", resolved, err)
	}
	for _, r := range resolved {
		if r.Kind != KindDataSource || r.ImplPath != "alicloud/provider/core/data_source_account.go" || r.TestPath != "alicloud/provider/core/data_source_account_test.go" {
			t.Fatalf("wrong core binding: %+v", r)
		}
	}
	for name, dir := range map[string]string{"alicloud_account": "alicloud/provider/core", "alicloud_product_core": "alicloud/service/core"} {
		r := FindByTypeName(name)
		if r == nil || filepath.Dir(r.ImplPath) != dir {
			t.Fatalf("%s resolved to %+v", name, r)
		}
	}
	if !FrameworkRegistered("alicloud_core_native") {
		t.Fatal("core framework registration was lost")
	}
	r := FrameworkFindByTypeName("alicloud_core_native")
	if r == nil || r.ImplPath != "alicloud/provider/core/data_source_native.go" {
		t.Fatalf("framework core: %+v", r)
	}
	fw, err := FrameworkResolve("provider/core", "data_source_native")
	if err != nil || len(fw) != 1 {
		t.Fatalf("framework resolve: %+v, %v", fw, err)
	}
	paths, err := RegisterPaths()
	if err != nil || len(paths) != 2 {
		t.Fatalf("registration paths: %v, %v", paths, err)
	}
}

func TestUnitCorePathsExcludeProviderInternals(t *testing.T) {
	for _, path := range []string{
		"alicloud/provider/conns/data_source_account.go",
		"alicloud/provider/registry/register.go",
		"alicloud/provider/core/nested/data_source_account.go",
		"alicloud/provider/core/../conns/data_source_account.go",
	} {
		if _, _, _, ok := IsSplitSource(path); ok {
			t.Errorf("accepted unsupported path %s", path)
		}
		if _, ok := ServiceProduct(path); ok {
			t.Errorf("accepted unsupported package %s", path)
		}
	}
}
