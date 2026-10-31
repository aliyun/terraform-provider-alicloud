package main

import "testing"

// TestUnitResolveRootEntity covers the root-layout type-name derivation that
// list_resources uses for alicloud/resource_alicloud_* and
// alicloud/data_source_alicloud_* files. It is a TestUnit* so ci-check-quick
// compiles the helper (including the splitlayout.Resolve path) and verifies
// the pure root derivation without requiring the repository root as the
// working directory that splitlayout.Resolve requires.
func TestUnitResolveRootEntity(t *testing.T) {
	cases := []struct {
		path string
		want string
		ok   bool
	}{
		{"alicloud/resource_alicloud_foo.go", "alicloud_foo", true},
		{"alicloud/data_source_alicloud_bar.go", "alicloud_bar", true},
		{"alicloud/resource_alicloud_foo_test.go", "", false},
		{"alicloud/helper.go", "", false},
		{"alicloud/provider.go", "", false},
	}
	for _, c := range cases {
		got, _, ok := resolveRootEntity(c.path)
		if ok != c.ok || got != c.want {
			t.Errorf("resolveRootEntity(%q) = (%q,%v), want (%q,%v)", c.path, got, ok, c.want, c.ok)
		}
	}
}
