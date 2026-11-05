package coverage

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSelection(t *testing.T) {
	src := []byte(`package alicloud; func resourceExample() {}`)
	old := map[string]string{"alicloud_example": "resourceExample"}
	current := map[string]string{"alicloud_example": "resourceExample"}
	sources := map[string][]byte{"alicloud/resource_alicloud_example.go": src}
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }
	cases := []struct {
		name          string
		paths         []string
		before, after map[string]string
		sources       map[string][]byte
		want          map[string]string
		fail          bool
	}{
		{"implementation", []string{"alicloud/resource_alicloud_example.go"}, old, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"test only or deleted test", []string{"alicloud/resource_alicloud_example_test.go"}, old, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"no resources", []string{"README.md"}, old, current, sources, map[string]string{}, false},
		{"registration", []string{"alicloud/provider.go"}, map[string]string{}, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"unresolved registration", []string{"alicloud/provider.go"}, map[string]string{}, map[string]string{"alicloud_unknown": "unknown"}, sources, nil, true},
		{"unregistered changed file", []string{"alicloud/resource_alicloud_unknown.go"}, old, current, map[string][]byte{"alicloud/resource_alicloud_unknown.go": []byte(`package alicloud;func unknown(){}`)}, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Select(tc.paths, tc.before, tc.after, tc.sources, missing)
			if (err != nil) != tc.fail || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v fail=%v", got, err, tc.want, tc.fail)
			}
		})
	}
	t.Run("rename", func(t *testing.T) {
		got, err := Select([]string{"alicloud/resource_alicloud_example.go", "alicloud/resource_alicloud_renamed.go"}, old, current, map[string][]byte{"alicloud/resource_alicloud_renamed.go": src}, func(string) ([]byte, error) { return src, nil })
		if err != nil || got["alicloud_example"] != "alicloud/resource_alicloud_renamed_test.go" {
			t.Fatal(got, err)
		}
	})
	t.Run("deleted resource", func(t *testing.T) {
		got, err := Select([]string{"alicloud/resource_alicloud_example.go"}, old, map[string]string{}, map[string][]byte{}, func(string) ([]byte, error) { return src, nil })
		if err != nil || len(got) != 0 {
			t.Fatal(got, err)
		}
	})
}

func TestRegistry(t *testing.T) {
	got, err := Registry([]byte(`package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":resourceExample()}}`))
	if err != nil || got["alicloud_example"] != "resourceExample" {
		t.Fatal(got, err)
	}
	for _, src := range []string{`package alicloud;var p = Provider{ResourcesMap: dynamic()}`, `package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":dynamic(name)}}`} {
		if _, err := Registry([]byte(src)); err == nil {
			t.Fatal("unresolved registry passed")
		}
	}
}

func TestSplitSelection(t *testing.T) {
	impl := []byte(`package cloudauth; func ResourceFaceConfig() {}`)
	implPath := "alicloud/service/cloudauth/resource_face_config.go"
	testPath := "alicloud/service/cloudauth/resource_face_config_test.go"
	imsPath := "alicloud/service/ims/resource_user.go"
	after := map[string]string{"alicloud_cloudauth_face_config": "cloudauth.ResourceFaceConfig"}
	before := map[string]string{"alicloud_cloudauth_face_config": "resourceAlicloudCloudauthFaceConfig"}
	sources := map[string][]byte{
		implPath: impl,
		imsPath:  []byte(`package ims; func resourceUser() {}`),
	}
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }
	want := map[string]string{"alicloud_cloudauth_face_config": testPath}
	cases := []struct {
		name   string
		paths  []string
		before map[string]string
		want   map[string]string
	}{
		{"split implementation", []string{implPath}, before, want},
		{"split test file", []string{testPath}, before, want},
		{"register is a registration source", []string{"alicloud/service/cloudauth/register.go"}, map[string]string{}, want},
		{"framework file without registration", []string{imsPath}, before, map[string]string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Select(tc.paths, tc.before, after, sources, missing)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	t.Run("cross-package function names do not collide", func(t *testing.T) {
		shared := map[string][]byte{
			"alicloud/service/a/resource_x.go": []byte(`package a; func Shared() {}`),
			"alicloud/service/b/resource_x.go": []byte(`package b; func Shared() {}`),
		}
		got, err := Select([]string{"alicloud/service/a/register.go"}, map[string]string{}, map[string]string{"alicloud_a_x": "a.Shared"}, shared, missing)
		if err != nil || got["alicloud_a_x"] != "alicloud/service/a/resource_x_test.go" {
			t.Fatal(got, err)
		}
	})
	t.Run("within-package duplicate fails", func(t *testing.T) {
		dup := map[string][]byte{
			"alicloud/service/a/resource_x.go": []byte(`package a; func Dup() {}`),
			"alicloud/service/a/resource_y.go": []byte(`package a; func Dup() {}`),
		}
		if _, err := Select(nil, map[string]string{}, map[string]string{}, dup, missing); err == nil {
			t.Fatal("expected duplicate function failure")
		}
	})
	t.Run("deleted split implementation with kept registration", func(t *testing.T) {
		if _, err := Select([]string{implPath}, before, after, map[string][]byte{}, func(string) ([]byte, error) { return impl, nil }); err == nil {
			t.Fatal("expected missing-implementation failure")
		}
	})
	t.Run("deleted split resource", func(t *testing.T) {
		got, err := Select([]string{implPath}, before, map[string]string{}, map[string][]byte{}, func(string) ([]byte, error) { return impl, nil })
		if err != nil || len(got) != 0 {
			t.Fatal(got, err)
		}
	})
}

func TestRegistryUnion(t *testing.T) {
	providerSrc := []byte(`package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":resourceExample()}}`)
	register := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cloudauth",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_config", Factory: ResourceFaceConfig},
		},
		SDKDataSources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_configs", Factory: DataSourceCloudauthFaceConfigs},
		},
		Resources: []conns.FrameworkResource{
			{TypeName: "alicloud_cloudauth_framework", Factory: NewFramework},
		},
	}
}`)
	got, err := RegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": register})
	if err != nil {
		t.Fatal(err)
	}
	if got["alicloud_example"] != "resourceExample" {
		t.Fatal("flat entry lost:", got)
	}
	if got["alicloud_cloudauth_face_config"] != "cloudauth.ResourceFaceConfig" {
		t.Fatal("SDK registration missing:", got)
	}
	if _, ok := got["alicloud_cloudauth_face_configs"]; ok {
		t.Fatal("data source leaked into the resource registry:", got)
	}
	if _, ok := got["alicloud_cloudauth_framework"]; ok {
		t.Fatal("framework registration leaked into the resource registry:", got)
	}
	dup := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{SDKResources: []conns.SDKResource{{TypeName: "alicloud_example", Factory: resourceExample}}}
}`)
	if _, err := RegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": dup}); err == nil {
		t.Fatal("duplicate declaration passed")
	}
}

func TestUnitFullySplitRegistry(t *testing.T) {
	empty := []byte(`package alicloud; var p = Provider{ResourcesMap: map[string]*Resource{}}`)
	register := []byte(`package example
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
 SDKResources: []conns.SDKResource{{TypeName: "alicloud_example", Factory: ResourceExample}},
} }`)
	got, err := RegistryUnion(empty, map[string][]byte{"alicloud/service/example/register.go": register})
	if err != nil || got["alicloud_example"] != "example.ResourceExample" {
		t.Fatalf("fully split registry: %v, %v", got, err)
	}
	for _, src := range []string{
		`package alicloud; var p = Provider{}`,
		`package alicloud; var p = Provider{ResourcesMap: dynamic()}`,
		`package alicloud; var p = Provider{ResourcesMap: map[string]*Resource{}, ResourcesMap: map[string]*Resource{}}`,
	} {
		if _, err := Registry([]byte(src)); err == nil {
			t.Fatalf("invalid root registration passed: %s", src)
		}
	}
}

// TestUnitHelperTest covers same-package helper test files (validation, unit, ...)
// whose mechanically trimmed implementation path never existed. Select must
// associate them with the registered resource whose real implementation is
// already in fileFunctions and run that resource's main coverage check, never
// calling oldSource on the nonexistent derived path. Real coverage gaps still
// surface because the resource stays selected, and genuinely unresolvable or
// deleted-but-registered cases still fail closed. The TestUnit prefix makes
// ci-check-quick's changed-unit-tests stage execute these regressions.
func TestUnitHelperTest(t *testing.T) {
	flatOld := map[string]string{"alicloud_example": "resourceExample"}
	flatAfter := map[string]string{"alicloud_example": "resourceExample"}
	flatSources := map[string][]byte{
		"alicloud/resource_alicloud_example.go": []byte(`package alicloud; func resourceExample() {}`),
	}
	splitBefore := map[string]string{"alicloud_cloudauth_face_config": "resourceAlicloudCloudauthFaceConfig"}
	splitAfter := map[string]string{"alicloud_cloudauth_face_config": "cloudauth.ResourceFaceConfig"}
	splitSources := map[string][]byte{
		"alicloud/service/cloudauth/resource_face_config.go": []byte(`package cloudauth; func ResourceFaceConfig() {}`),
	}
	// missing fails on every call, so any case that reaches oldSource errors.
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }

	cases := []struct {
		name    string
		paths   []string
		before  map[string]string
		after   map[string]string
		sources map[string][]byte
		want    map[string]string
		fail    bool
	}{
		// AC1: legal helper naming must not trigger a git show on a file that
		// never existed; the resource's main test is selected instead.
		{"flat validation helper", []string{"alicloud/resource_alicloud_example_validation_test.go"}, flatOld, flatAfter, flatSources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"flat unit helper", []string{"alicloud/resource_alicloud_example_unit_test.go"}, flatOld, flatAfter, flatSources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"split validation helper", []string{"alicloud/service/cloudauth/resource_face_config_validation_test.go"}, splitBefore, splitAfter, splitSources, map[string]string{"alicloud_cloudauth_face_config": "alicloud/service/cloudauth/resource_face_config_test.go"}, false},
		{"split unit helper", []string{"alicloud/service/cloudauth/resource_face_config_unit_test.go"}, splitBefore, splitAfter, splitSources, map[string]string{"alicloud_cloudauth_face_config": "alicloud/service/cloudauth/resource_face_config_test.go"}, false},
		// AC2: existing main test, new registration and existing-resource
		// helper changes all select the resource without error.
		{"flat main test", []string{"alicloud/resource_alicloud_example_test.go"}, flatOld, flatAfter, flatSources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"flat new registration", []string{"alicloud/provider.go"}, map[string]string{}, flatAfter, flatSources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"split new registration", []string{"alicloud/service/cloudauth/register.go"}, map[string]string{}, splitAfter, splitSources, map[string]string{"alicloud_cloudauth_face_config": "alicloud/service/cloudauth/resource_face_config_test.go"}, false},
		// AC3: renaming a main test to a helper name cannot bypass the gate;
		// the resource is still selected so its coverage check runs.
		{"renamed main test selects resource", []string{"alicloud/resource_alicloud_example_helper_test.go"}, flatOld, flatAfter, flatSources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		// AC4: a test with no matching implementation in the package is
		// genuinely unresolvable and fails closed.
		{"unresolvable test blocks", []string{"alicloud/service/cloudauth/resource_orphan_test.go"}, splitBefore, splitAfter, splitSources, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Select(tc.paths, tc.before, tc.after, tc.sources, missing)
			if (err != nil) != tc.fail || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v fail=%v", got, err, tc.want, tc.fail)
			}
		})
	}
	// AC1: assert oldSource is never called for a helper test.
	t.Run("helper does not call oldSource", func(t *testing.T) {
		called := false
		oldSource := func(string) ([]byte, error) { called = true; return nil, fmt.Errorf("missing") }
		if _, err := Select([]string{"alicloud/service/cloudauth/resource_face_config_validation_test.go"}, splitBefore, splitAfter, splitSources, oldSource); err != nil {
			t.Fatalf("helper test failed: %v", err)
		}
		if called {
			t.Fatal("oldSource was called for a helper test")
		}
	})
	// AC4: implementation deleted but still registered blocks; oldSource
	// succeeds on the base, add fails because the implementation is gone.
	t.Run("deleted implementation with kept registration blocks", func(t *testing.T) {
		impl := []byte(`package cloudauth; func ResourceFaceConfig() {}`)
		if _, err := Select([]string{"alicloud/service/cloudauth/resource_face_config.go"}, splitBefore, splitAfter, map[string][]byte{}, func(string) ([]byte, error) { return impl, nil }); err == nil {
			t.Fatal("expected missing-implementation failure")
		}
	})
	// AC2/AC3: the longest implementation stem wins, so an extra-resource
	// helper selects the extra resource and not its shorter neighbor.
	t.Run("longest stem match", func(t *testing.T) {
		sources := map[string][]byte{
			"alicloud/service/cloudauth/resource_face_config.go":       []byte(`package cloudauth; func ResourceFaceConfig() {}`),
			"alicloud/service/cloudauth/resource_face_config_extra.go": []byte(`package cloudauth; func ResourceFaceConfigExtra() {}`),
		}
		after := map[string]string{
			"alicloud_cloudauth_face_config":       "cloudauth.ResourceFaceConfig",
			"alicloud_cloudauth_face_config_extra": "cloudauth.ResourceFaceConfigExtra",
		}
		got, err := Select([]string{"alicloud/service/cloudauth/resource_face_config_extra_validation_test.go"}, splitBefore, after, sources, missing)
		if err != nil {
			t.Fatalf("longest stem match failed: %v", err)
		}
		if got["alicloud_cloudauth_face_config_extra"] != "alicloud/service/cloudauth/resource_face_config_extra_test.go" {
			t.Fatalf("extra resource not selected: %v", got)
		}
		if _, ok := got["alicloud_cloudauth_face_config"]; ok {
			t.Fatalf("shorter resource should not be selected by the extra helper: %v", got)
		}
	})
}

// TestUnitDataSourceSelection covers data-source selection through SelectScope
// with DataSourceScope: flat and split data-source implementations, tests and
// registrations all select the registered data source, and the resource and
// data-source scopes never steal each other's files (a data_source_ file is
// invisible to ResourceScope and a resource_ file is invisible to
// DataSourceScope). The TestUnit prefix makes ci-check-quick's
// changed-unit-tests stage execute these regressions.
func TestUnitDataSourceSelection(t *testing.T) {
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }

	// Flat data source: alicloud/data_source_alicloud_example.go.
	flatSrc := []byte(`package alicloud; func dataSourceExample() {}`)
	flatPath := "alicloud/data_source_alicloud_example.go"
	flatOld := map[string]string{"alicloud_example": "dataSourceExample"}
	flatAfter := map[string]string{"alicloud_example": "dataSourceExample"}
	flatSources := map[string][]byte{flatPath: flatSrc}
	flatWant := map[string]string{"alicloud_example": "alicloud/data_source_alicloud_example_test.go"}

	// Split data source: alicloud/service/cloudauth/data_source_face_configs.go.
	splitImpl := []byte(`package cloudauth; func DataSourceCloudauthFaceConfigs() {}`)
	splitPath := "alicloud/service/cloudauth/data_source_face_configs.go"
	splitTest := "alicloud/service/cloudauth/data_source_face_configs_test.go"
	splitOld := map[string]string{"alicloud_cloudauth_face_configs": "dataSourceAlicloudCloudauthFaceConfigs"}
	splitAfter := map[string]string{"alicloud_cloudauth_face_configs": "cloudauth.DataSourceCloudauthFaceConfigs"}
	splitSources := map[string][]byte{splitPath: splitImpl}
	splitWant := map[string]string{"alicloud_cloudauth_face_configs": splitTest}

	cases := []struct {
		name    string
		paths   []string
		before  map[string]string
		after   map[string]string
		sources map[string][]byte
		want    map[string]string
	}{
		{"flat implementation", []string{flatPath}, flatOld, flatAfter, flatSources, flatWant},
		{"flat test file", []string{"alicloud/data_source_alicloud_example_test.go"}, flatOld, flatAfter, flatSources, flatWant},
		{"flat registration", []string{"alicloud/provider.go"}, map[string]string{}, flatAfter, flatSources, flatWant},
		{"split implementation", []string{splitPath}, splitOld, splitAfter, splitSources, splitWant},
		{"split test file", []string{splitTest}, splitOld, splitAfter, splitSources, splitWant},
		{"split registration", []string{"alicloud/service/cloudauth/register.go"}, map[string]string{}, splitAfter, splitSources, splitWant},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SelectScope(tc.paths, tc.before, tc.after, tc.sources, missing, DataSourceScope)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	// Cross-scope isolation: a data source's Select never matches a resource_
	// file and a resource's Select never matches a data_source_ file. Each
	// scope's after registry here carries a matching constructor, so a leak
	// would surface as a selected entry; the correct scope filter selects none.
	t.Run("resource scope skips data source files", func(t *testing.T) {
		got, err := SelectScope([]string{splitPath}, splitOld, splitAfter, splitSources, missing, ResourceScope)
		if err != nil || len(got) != 0 {
			t.Fatalf("resource scope selected a data source file: %v, %v", got, err)
		}
	})
	t.Run("data source scope skips resource files", func(t *testing.T) {
		resImpl := []byte(`package cloudauth; func ResourceFaceConfig() {}`)
		resPath := "alicloud/service/cloudauth/resource_face_config.go"
		resOld := map[string]string{"alicloud_cloudauth_face_config": "resourceAlicloudCloudauthFaceConfig"}
		resAfter := map[string]string{"alicloud_cloudauth_face_config": "cloudauth.ResourceFaceConfig"}
		resSources := map[string][]byte{resPath: resImpl}
		got, err := SelectScope([]string{resPath}, resOld, resAfter, resSources, missing, DataSourceScope)
		if err != nil || len(got) != 0 {
			t.Fatalf("data source scope selected a resource file: %v, %v", got, err)
		}
	})
}

// TestUnitDataSourcesRegistry covers DataSourcesRegistry and
// DataSourcesRegistryUnion: the DataSourcesMap literal is parsed the same way
// ResourcesMap is, the split SDKDataSources fold in under a product qualifier,
// resources and framework registrations never leak into the data-source
// registry, duplicates fail, and a resource and data source sharing a type name
// each resolve in their own registry. The TestUnit prefix makes ci-check-quick's
// changed-unit-tests stage execute these regressions.
func TestUnitDataSourcesRegistry(t *testing.T) {
	got, err := DataSourcesRegistry([]byte(`package alicloud;var p = Provider{DataSourcesMap:map[string]*Resource{"alicloud_example":dataSourceExample()}}`))
	if err != nil || got["alicloud_example"] != "dataSourceExample" {
		t.Fatalf("DataSourcesRegistry: %v, %v", got, err)
	}
	for _, src := range []string{
		`package alicloud;var p = Provider{DataSourcesMap: dynamic()}`,
		`package alicloud;var p = Provider{DataSourcesMap:map[string]*Resource{"alicloud_example":dynamic(name)}}`,
		`package alicloud;var p = Provider{}`,
		`package alicloud;var p = Provider{DataSourcesMap: map[string]*Resource{}, DataSourcesMap: map[string]*Resource{}}`,
	} {
		if _, err := DataSourcesRegistry([]byte(src)); err == nil {
			t.Fatalf("invalid data-source registration passed: %s", src)
		}
	}
	providerSrc := []byte(`package alicloud;var p = Provider{DataSourcesMap:map[string]*Resource{"alicloud_example":dataSourceExample()}}`)
	register := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cloudauth",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_config", Factory: ResourceFaceConfig},
		},
		SDKDataSources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_configs", Factory: DataSourceCloudauthFaceConfigs},
		},
		Resources: []conns.FrameworkResource{
			{TypeName: "alicloud_cloudauth_framework", Factory: NewFramework},
		},
	}
}`)
	got, err = DataSourcesRegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": register})
	if err != nil {
		t.Fatal(err)
	}
	if got["alicloud_example"] != "dataSourceExample" {
		t.Fatal("flat data-source entry lost:", got)
	}
	if got["alicloud_cloudauth_face_configs"] != "cloudauth.DataSourceCloudauthFaceConfigs" {
		t.Fatal("SDK data-source registration missing:", got)
	}
	if _, ok := got["alicloud_cloudauth_face_config"]; ok {
		t.Fatal("resource leaked into the data-source registry:", got)
	}
	if _, ok := got["alicloud_cloudauth_framework"]; ok {
		t.Fatal("framework resource leaked into the data-source registry:", got)
	}
	dup := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{SDKDataSources: []conns.SDKResource{{TypeName: "alicloud_example", Factory: dataSourceExample}}}
}`)
	if _, err := DataSourcesRegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": dup}); err == nil {
		t.Fatal("duplicate data-source declaration passed")
	}
	// A resource and data source sharing a type name stay in separate
	// registries: DataSourcesRegistryUnion resolves the data source and
	// RegistryUnion resolves the resource, with no cross-contamination.
	shared := []byte(`package shared
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		SDKResources: []conns.SDKResource{{TypeName: "alicloud_shared", Factory: ResourceShared}},
		SDKDataSources: []conns.SDKResource{{TypeName: "alicloud_shared", Factory: DataSourceShared}},
	}
}`)
	dsGot, err := DataSourcesRegistryUnion([]byte(`package alicloud;var p = Provider{DataSourcesMap:map[string]*Resource{}}`), map[string][]byte{"alicloud/service/shared/register.go": shared})
	if err != nil || dsGot["alicloud_shared"] != "shared.DataSourceShared" {
		t.Fatalf("same-named data source not resolved: %v, %v", dsGot, err)
	}
	resGot, err := RegistryUnion([]byte(`package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{}}`), map[string][]byte{"alicloud/service/shared/register.go": shared})
	if err != nil || resGot["alicloud_shared"] != "shared.ResourceShared" {
		t.Fatalf("same-named resource not resolved: %v, %v", resGot, err)
	}
}
