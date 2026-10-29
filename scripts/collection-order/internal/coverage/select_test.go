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
		// A flat file registered in before whose type name is no longer in
		// after is a real public registration deletion, not a relocation: it
		// fails closed so a deletion cannot be waived as a move.
		if _, err := Select([]string{"alicloud/resource_alicloud_example.go"}, old, map[string]string{}, map[string][]byte{}, func(string) ([]byte, error) { return src, nil }); err == nil {
			t.Fatal("expected registered-resource deletion failure")
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

// TestUnitFlatDeletionClassification covers the flat-layout deletion cases the
// collection-order gate must distinguish after a product split moves flat
// files into service packages. The type name — not the factory name, which a
// split move renames and package-qualifies — is the stable link between a
// deleted flat file and the registry. The TestUnit prefix makes
// ci-check-quick's changed-unit-tests stage execute these regressions.
func TestUnitFlatDeletionClassification(t *testing.T) {
	// #11402 model: alicloud_havip and alicloud_vpc_ha_vip are public aliases
	// that at the flat layout both point to resourceAliCloudVpcHaVip (declared
	// in resource_alicloud_vpc_ha_vip.go). A dead, never-registered
	// resourceAlicloudHavip (typoed "Alicloud", not "AliCloud") lives in the
	// resource_alicloud_havip.go file whose name matches the alias — exactly
	// the file a type-name path guess would pick.
	before := map[string]string{
		"alicloud_havip":      "resourceAliCloudVpcHaVip",
		"alicloud_vpc_ha_vip": "resourceAliCloudVpcHaVip",
	}
	// After the split both aliases are still registered, under the renamed,
	// package-qualified factory.
	after := map[string]string{
		"alicloud_havip":      "vpc.ResourceAliCloudVpcHaVip",
		"alicloud_vpc_ha_vip": "vpc.ResourceAliCloudVpcHaVip",
	}
	deadFactory := []byte(`package alicloud; func resourceAlicloudHavip() {}`)
	realFactory := []byte(`package alicloud; func resourceAliCloudVpcHaVip() {}`)
	deadFile := "alicloud/resource_alicloud_havip.go"
	realFile := "alicloud/resource_alicloud_vpc_ha_vip.go"

	// AC5: the dead factory's file (declares resourceAlicloudHavip, never a
	// before value) moved with the split must not be flagged as a
	// registered-resource deletion.
	t.Run("dead factory relocation is not a deletion", func(t *testing.T) {
		got, err := Select([]string{deadFile}, before, after, map[string][]byte{}, func(string) ([]byte, error) { return deadFactory, nil })
		if err != nil || len(got) != 0 {
			t.Fatalf("expected no error and no selection; got %v, %v", got, err)
		}
	})

	// A registered implementation moved with an export rename: the deleted
	// flat file declares the bare baseline factory, whose type name is still
	// in after under the qualified renamed factory. The moved (split) file
	// carries selection; the deleted flat file is benign.
	t.Run("factory rename move is not a deletion", func(t *testing.T) {
		got, err := Select([]string{realFile}, before, after, map[string][]byte{}, func(string) ([]byte, error) { return realFactory, nil })
		if err != nil || len(got) != 0 {
			t.Fatalf("expected no error and no selection; got %v, %v", got, err)
		}
	})

	// AC6: a constructor registered in before whose type name is no longer in
	// after is a real public registration deletion, even though the
	// implementation file is also gone — it cannot be waived as a relocation.
	t.Run("real registration deletion fails", func(t *testing.T) {
		removedAfter := map[string]string{"alicloud_vpc_ha_vip": "vpc.ResourceAliCloudVpcHaVip"}
		if _, err := Select([]string{realFile}, before, removedAfter, map[string][]byte{}, func(string) ([]byte, error) { return realFactory, nil }); err == nil {
			t.Fatal("expected registered-resource deletion failure")
		}
	})

	// AC7: a deleted flat file whose base content cannot be read still fails
	// closed; the oldSource error is not swallowed.
	t.Run("missing old source fails", func(t *testing.T) {
		missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }
		if _, err := Select([]string{realFile}, before, after, map[string][]byte{}, missing); err == nil {
			t.Fatal("expected oldSource resolution failure")
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
