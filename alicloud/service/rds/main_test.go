package rds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain delegates to the plugin-testing resource runner so the package's
// registered AddTestSweepers expose the standard -sweep entry for the test binary.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
