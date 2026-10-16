package main

import "testing"

// TestUnitFrameworkNativeSkip exercises the pure logic that decides whether a
// resource or data source absent from the SDKv2 provider maps is registered
// framework-native and should be skipped. It stubs frameworkRegistered so no
// filesystem register.go is read; the production default reads the service
// layout, but the decision logic is what matters here.
func TestUnitFrameworkNativeSkip(t *testing.T) {
	orig := frameworkRegistered
	defer func() { frameworkRegistered = orig }()

	t.Run("framework-native resource skips", func(t *testing.T) {
		frameworkRegistered = func(name string) bool { return name == "alicloud_ecs_saving_plan" }
		if !frameworkNativeSkip("ecs_saving_plan") {
			t.Fatal("framework-native resource should be skipped")
		}
	})
	t.Run("framework-native data source skips", func(t *testing.T) {
		frameworkRegistered = func(name string) bool { return name == "alicloud_ecs_saving_plans" }
		if !frameworkNativeSkip("ecs_saving_plans") {
			t.Fatal("framework-native data source should be skipped")
		}
	})
	t.Run("unregistered name does not skip", func(t *testing.T) {
		frameworkRegistered = func(string) bool { return false }
		if frameworkNativeSkip("nonexistent_xyz") {
			t.Fatal("unregistered name should not be skipped")
		}
	})
	t.Run("prefix is prepended exactly once", func(t *testing.T) {
		// frameworkNativeSkip prepends "alicloud_"; a caller passing the
		// already-prefixed name must not match, proving no double-prefix.
		frameworkRegistered = func(name string) bool { return name == "alicloud_ecs_saving_plan" }
		if frameworkNativeSkip("alicloud_ecs_saving_plan") {
			t.Fatal("frameworkNativeSkip should prepend alicloud_, not double-prepend")
		}
	})
}

// TestUnitCheckResourceByNameFrameworkNativeSkips verifies the end-to-end path:
// when a resource is absent from the SDKv2 ResourcesMap (Framework-native
// resources live in a Framework provider server, never the SDKv2 map) and is
// registered framework-native, checkResourceByName returns true rather than
// failing. "zzz_nonexistent" is guaranteed absent; the stubbed
// frameworkRegistered makes it appear framework-native. This proves the
// checker skips framework resources instead of false-failing coverage.
func TestUnitCheckResourceByNameFrameworkNativeSkips(t *testing.T) {
	orig := frameworkRegistered
	frameworkRegistered = func(string) bool { return true }
	defer func() { frameworkRegistered = orig }()
	if !checkResourceByName("zzz_nonexistent", true, "resource", "", "") {
		t.Fatal("framework-native resource should be skipped by checkResourceByName")
	}
}

// TestUnitCheckResourceByNameUnregisteredFails is the negative counterpart: a name
// absent from the SDKv2 maps and not framework-native must fail (return false),
// so a genuinely unregistered resource is reported rather than silently passed.
func TestUnitCheckResourceByNameUnregisteredFails(t *testing.T) {
	orig := frameworkRegistered
	frameworkRegistered = func(string) bool { return false }
	defer func() { frameworkRegistered = orig }()
	if checkResourceByName("zzz_nonexistent", true, "resource", "", "") {
		t.Fatal("unregistered resource should fail coverage check")
	}
}

// TestUnitCheckResourceByNameFrameworkNativeDataSourceSkips verifies the data
// source branch of the same fail-open path.
func TestUnitCheckResourceByNameFrameworkNativeDataSourceSkips(t *testing.T) {
	orig := frameworkRegistered
	frameworkRegistered = func(string) bool { return true }
	defer func() { frameworkRegistered = orig }()
	if !checkResourceByName("zzz_nonexistent", false, "data source", "", "") {
		t.Fatal("framework-native data source should be skipped by checkResourceByName")
	}
}
