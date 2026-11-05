package ddoscoo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain delegates to the plugin-testing entry point so the split-layout
// ddoscoo test binary honors the -sweep branch when invoked in isolation
// rather than falling back to the default test runner. It lives in a
// standalone sweeper_test.go file — not a resource_*_test.go / data_source_*_test.go
// file — so the testing coverage checker, which scans per-resource test
// files for TestAccAliCloud-prefixed cases, does not flag TestMain as a
// non-conforming testcase. This mirrors the root alicloud/alicloud_sweeper_test.go.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
