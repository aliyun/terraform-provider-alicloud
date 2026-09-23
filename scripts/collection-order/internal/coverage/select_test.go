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
