package dbs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain delegates to resource.TestMain so the AddTestSweepers
// registrations declared in this dbs_test package (e.g.
// alicloud_dbs_backup_plan) are processed for the leaf package's
// test binary, mirroring the root alicloud_sweeper_test.go entry.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
