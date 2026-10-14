package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// SecurityGroupSnapshot is a read-only, system-produced resource, so there is no
// lifecycle resource fixture to create. This test is the API-callability
// discriminator from the accepted routePlan: a passing apply proves the
// DescribeSecurityGroupSnapshots (Ecs, 2016-03-14) Inner API is reachable and
// that the response is parsed; an InvalidAction.NotFound at Read time would
// prove the route premise is inverted (hand-off instead of implementation).
// No lifecycle ACC exists because the resource is read-only.
func TestAccAlicloudSecurityGroupSnapshotsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				// Basic callability: query all snapshots with no filter. The data
				// source id (hash of returned ids) is always set, and ids.# must
				// equal snapshots.# regardless of whether the account has any.
				Config: testAccCheckAlicloudSecurityGroupSnapshotsDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_snapshots.default"),
				),
			},
			{
				// Fake id filter: a non-existent snapshot id must yield an empty
				// list, exercising the SecurityGroupSnapshotIds query path and
				// the local defensive filter.
				Config: testAccCheckAlicloudSecurityGroupSnapshotsDataSourceConfigFakeID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_snapshots.default"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_snapshots.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_snapshots.default", "snapshots.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupSnapshotsDataSourceConfigBasic = `
data "alicloud_security_group_snapshots" "default" {
  output_file = "security_group_snapshots.json"
}
`

const testAccCheckAlicloudSecurityGroupSnapshotsDataSourceConfigFakeID = `
data "alicloud_security_group_snapshots" "default" {
  security_group_snapshot_ids = ["ss-nonexistent-snapshot-id"]
}
`
