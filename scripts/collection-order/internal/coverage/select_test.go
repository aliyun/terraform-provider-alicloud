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
	// Helper tests (resource_foo_crud_test.go, resource_foo_internal_test.go)
	// derive a resource_foo_crud.go impl that does not exist; they must resolve to
	// the registered resource's real impl (resource_foo.go) so the collection-order
	// check still applies — coverage is preserved, not blanket-ignored. An orphan
	// helper test (no matching impl, or an impl that is not registered) fails closed.
	t.Run("split crud helper test resolves to registered impl", func(t *testing.T) {
		got, err := Select([]string{"alicloud/service/cloudauth/resource_face_config_crud_test.go"}, before, after, sources, missing)
		if err != nil || got["alicloud_cloudauth_face_config"] != "alicloud/service/cloudauth/resource_face_config_test.go" {
			t.Fatalf("crud helper test must resolve to the registered resource's standard test; got %v, %v", got, err)
		}
	})
	t.Run("split internal helper test resolves to registered impl", func(t *testing.T) {
		got, err := Select([]string{"alicloud/service/cloudauth/resource_face_config_internal_test.go"}, before, after, sources, missing)
		if err != nil || got["alicloud_cloudauth_face_config"] != "alicloud/service/cloudauth/resource_face_config_test.go" {
			t.Fatalf("internal helper test must resolve to the registered resource's standard test; got %v, %v", got, err)
		}
	})
	t.Run("orphan helper test with no impl fails closed", func(t *testing.T) {
		if _, err := Select([]string{"alicloud/service/cloudauth/resource_orphan_crud_test.go"}, before, after, sources, missing); err == nil {
			t.Fatal("an orphan helper test with no matching impl must fail closed, not be silently skipped")
		}
	})
	t.Run("helper test resolving to an unregistered impl fails closed", func(t *testing.T) {
		unregisteredSources := map[string][]byte{
			implPath: impl,
			"alicloud/service/cloudauth/resource_other.go": []byte(`package cloudauth; func ResourceOther() {}`),
		}
		if _, err := Select([]string{"alicloud/service/cloudauth/resource_other_crud_test.go"}, before, after, unregisteredSources, missing); err == nil {
			t.Fatal("a helper test resolving to an unregistered impl must fail closed")
		}
	})
	// A framework-native resource may keep its ACC test in the flat alicloud/
	// directory (package alicloud_test, so `go test ./alicloud` discovers it)
	// while its real implementation lives in alicloud/service/<product>/. The
	// flat impl path the standard pairing derives does not exist at either
	// revision; the check must resolve it to the split-layout impl so the
	// collection-order check still applies — coverage is preserved, not
	// blanket-ignored. A flat test whose split impl is not registered fails
	// closed, and a flat test with no matching split impl also fails closed.
	t.Run("flat test for split-layout framework resource resolves to split impl", func(t *testing.T) {
		splitImpl := "alicloud/service/ecs/resource_security_group_snapshot_policy.go"
		splitTest := "alicloud/service/ecs/resource_security_group_snapshot_policy_test.go"
		flatTest := "alicloud/resource_alicloud_security_group_snapshot_policy_test.go"
		fwAfter := map[string]string{"alicloud_security_group_snapshot_policy": "ecs.NewSecurityGroupSnapshotPolicy"}
		fwSources := map[string][]byte{
			splitImpl: []byte(`package ecs; func NewSecurityGroupSnapshotPolicy() {}`),
		}
		got, err := Select([]string{flatTest}, map[string]string{}, fwAfter, fwSources, missing)
		if err != nil {
			t.Fatalf("flat test for a split-layout resource must resolve; got %v", err)
		}
		if got["alicloud_security_group_snapshot_policy"] != splitTest {
			t.Fatalf("flat test must resolve to the split-layout standard test; got %v, want %v", got, map[string]string{"alicloud_security_group_snapshot_policy": splitTest})
		}
	})
	t.Run("flat test for unregistered split impl fails closed", func(t *testing.T) {
		splitImpl := "alicloud/service/ecs/resource_security_group_snapshot_policy.go"
		flatTest := "alicloud/resource_alicloud_security_group_snapshot_policy_test.go"
		fwSources := map[string][]byte{
			splitImpl: []byte(`package ecs; func NewSecurityGroupSnapshotPolicy() {}`),
		}
		if _, err := Select([]string{flatTest}, map[string]string{}, map[string]string{}, fwSources, missing); err == nil {
			t.Fatal("a flat test whose split impl is not registered must fail closed")
		}
	})
	t.Run("flat test with no matching split impl fails closed", func(t *testing.T) {
		flatTest := "alicloud/resource_alicloud_security_group_snapshot_policy_test.go"
		// sources only have a flat impl for a different resource
		flatSources := map[string][]byte{
			"alicloud/resource_alicloud_example.go": []byte(`package alicloud; func resourceExample() {}`),
		}
		if _, err := Select([]string{flatTest}, map[string]string{}, map[string]string{"alicloud_security_group_snapshot_policy": "ecs.NewSecurityGroupSnapshotPolicy"}, flatSources, missing); err == nil {
			t.Fatal("a flat test with no matching split impl must fail closed")
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
