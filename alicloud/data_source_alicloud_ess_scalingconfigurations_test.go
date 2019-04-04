package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssScalingconfigurationsDataSource_ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingconfigurationsIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_ids", "configurations.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingconfigurationsIdsNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_ids", "configurations.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingconfigurationsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingconfigurationsNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_name_regex", "configurations.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingconfigurationsNameRegexNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_name_regex", "configurations.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingconfigurationsDataSource_scaling_group_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingconfigurationsScalingGroupId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_scaling_group_id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_scaling_group_id", "configurations.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingconfigurationsScalingGroupIdNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_scaling_group_id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_scaling_group_id", "configurations.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingconfigurationsDataSource_combined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingconfigurationsCombined,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.name", "tf-testAccDataSourceScalingConfiguration1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.image_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.security_group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.system_disk_category"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.system_disk_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.internet_max_bandwidth_in"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.internet_max_bandwidth_out"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.internet_charge_type"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.0.data_disks.#", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingconfigurationsCombinedNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_configurations.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_configurations.foo_combined", "configurations.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudScalingconfigurationsCombined = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_combined"{
	scaling_group_id = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_group_id}"
	name_regex = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_configuration_name}"
	ids = ["${alicloud_ess_scaling_configuration.scaling_configuration1.id}"]
}
`

const testAccCheckAlicloudScalingconfigurationsCombinedNotFound = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_combined"{
	scaling_group_id = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_group_id}"
	name_regex = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_configuration_name}-fake"
	ids = ["${alicloud_ess_scaling_configuration.scaling_configuration1.id}"]
}
`

const testAccCheckAlicloudScalingconfigurationsScalingGroupId = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_scaling_group_id"{
	scaling_group_id = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_group_id}"
}
`
const testAccCheckAlicloudScalingconfigurationsScalingGroupIdNotFound = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_scaling_group_id"{
	scaling_group_id = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_group_id}-fake"
}
`

const testAccCheckAlicloudScalingconfigurationsNameRegex = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_name_regex"{
	name_regex = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_configuration_name}"
}
`
const testAccCheckAlicloudScalingconfigurationsNameRegexNotFound = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_name_regex"{
	name_regex = "${alicloud_ess_scaling_configuration.scaling_configuration1.scaling_configuration_name}-fake"
}
`

const testAccCheckAlicloudScalingconfigurationsIds = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_ids"{
	ids = ["${alicloud_ess_scaling_configuration.scaling_configuration1.id}"]
}
`
const testAccCheckAlicloudScalingconfigurationsIdsNotFound = testAccCheckAlicloudScalingconfigurationsBasicConfig + `
data "alicloud_ess_scaling_configurations" "foo_ids"{
	ids = ["${alicloud_ess_scaling_configuration.scaling_configuration1.id}-fake"]
}
`

const testAccCheckAlicloudScalingconfigurationsBasicConfig = EcsInstanceCommonTestCase + `

variable "name" {
	default = "tf-testAccDataSourceEssScalingRules"
}

resource "alicloud_ess_scaling_group" "scalinggroup_foo1" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "tf-testAccDataSourceScalingRuleEssScalingGroup1"
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

resource "alicloud_ess_scaling_configuration" "scaling_configuration1"{
	scaling_group_id = "${alicloud_ess_scaling_group.scalinggroup_foo1.id}"
	scaling_configuration_name = "tf-testAccDataSourceScalingConfiguration1"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.default.id}"
	force_delete = true
}
`
