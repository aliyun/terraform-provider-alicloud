package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssScalinggroupsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_name_regex", "groups.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceNameRegexNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_name_regex", "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalinggroupsDataSource_ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_ids", "groups.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceIdsNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_ids", "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalinggroupsDataSource_combined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceCombined,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.name", "tf-testAccEssScalingGroup1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.min_size", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.max_size", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.cooldown_time", "20"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.removal_policies.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.removal_policies.0", "OldestInstance"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.removal_policies.1", "NewestInstance"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.load_balancer_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.db_instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.vswitch_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.total_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.active_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.pending_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.removing_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_groups.foo_combined", "groups.0.lifecycle_state"),
				),
			},
			{
				Config: testAccCheckAlicloudScalinggroupsDataSourceCombinedNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_groups.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_groups.foo_combined", "groups.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudScalinggroupsDataSourceNameRegex = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_name_regex" {
	name_regex = "${alicloud_ess_scaling_group.scalinggroup_foo1.scaling_group_name}"
}
`

const testAccCheckAlicloudScalinggroupsDataSourceNameRegexNotFound = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_name_regex" {
	name_regex = "${alicloud_ess_scaling_group.scalinggroup_foo1.scaling_group_name}-fake"
}
`

const testAccCheckAlicloudScalinggroupsDataSourceIds = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_ids" {
	ids = ["${alicloud_ess_scaling_group.scalinggroup_foo1.id}"]
}
`

const testAccCheckAlicloudScalinggroupsDataSourceIdsNotFound = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_ids" {
	ids = ["${alicloud_ess_scaling_group.scalinggroup_foo1.id}-fake"]
}

`

const testAccCheckAlicloudScalinggroupsDataSourceCombined = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_combined" {
	ids = ["${alicloud_ess_scaling_group.scalinggroup_foo1.id}"]
	name_regex = "${alicloud_ess_scaling_group.scalinggroup_foo1.scaling_group_name}"
}
`
const testAccCheckAlicloudScalinggroupsDataSourceCombinedNotFound = testAccCheckAlicloudScalinggroupsBasicConfig + `
data "alicloud_ess_scaling_groups" "foo_combined" {
	ids = ["${alicloud_ess_scaling_group.scalinggroup_foo1.id}-fake"]
	name_regex = "${alicloud_ess_scaling_group.scalinggroup_foo1.scaling_group_name}"
}
`

const testAccCheckAlicloudScalinggroupsBasicConfig = EcsInstanceCommonTestCase + `

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroups"
}

resource "alicloud_ess_scaling_group" "scalinggroup_foo1" {
	min_size = 0
	max_size = 2
	scaling_group_name = "tf-testAccEssScalingGroup1"
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

`
