package governance_test

import "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

// resourceTestAccConfigFunc is the package-local counterpart of the flat
// alicloud test helper (alicloud/service_alicloud_common_test.go), which is
// not reachable from a service package test. It delegates to the shared
// acctest.ResourceTestAccConfig builder and preserves the bare name so the
// scripts/collection-order TypeList order coverage checker recognizes the
// config builder by name when tracking cumulative step deltas: the checker
// matches the bare identifier "resourceTestAccConfigFunc" and follows the
// accumulated attribute map across apply steps to detect convergent reorders.
func resourceTestAccConfigFunc(resourceId string, name string, configDependence func(name string) string) acctest.ResourceTestAccConfigFunc {
	return acctest.ResourceTestAccConfig(resourceId, name, configDependence)
}
