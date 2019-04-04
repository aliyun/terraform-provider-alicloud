package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// Skip this test, because alicloud_ess_scaling_rule.scaling_rule1.id is now
// composed of scaling_group_id and scaling_rule_id.
func SkipTestAccAlicloudEssScalingrulesDataSource_ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_ids", "rules.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceIdsNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_ids"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_ids", "rules.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingrulesDataSource_scaling_group_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceScalingGroupId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_scaling_group_id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_scaling_group_id", "rules.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceScalingGroupIdNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_scaling_group_id"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_scaling_group_id", "rules.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingrulesDataSource_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_type"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_type", "rules.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceTypeNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_type"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_type", "rules.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingrulesDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_scaling_group_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_scaling_group_name_regex", "rules.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceNameRegexNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_scaling_group_name_regex"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_scaling_group_name_regex", "rules.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingrulesDataSource_combined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceCombined,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.scaling_rule_ari"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.name", "tf-testAccDataSourceScalingRule1"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.type", "SimpleScalingRule"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.adjustment_value", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.0.adjustment_type", "QuantityChangeInCapacity"),
				),
			},
			{
				Config: testAccCheckAlicloudScalingrulesDataSourceCombinedNotFound,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ess_scaling_rules.foo_combined"),
					resource.TestCheckResourceAttr("data.alicloud_ess_scaling_rules.foo_combined", "rules.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudScalingrulesDataSourceType = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_type" {
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}"
	type = "SimpleScalingRule"
}
`

const testAccCheckAlicloudScalingrulesDataSourceTypeNotFound = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_type" {
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}"
	type = "TargetTrackingScalingRule"
}
`

const testAccCheckAlicloudScalingrulesDataSourceIds = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_ids" {
	ids = ["${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_id}"]
}
`

const testAccCheckAlicloudScalingrulesDataSourceIdsNotFound = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_ids" {
	ids = ["${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_id}-fake"]
}
`

const testAccCheckAlicloudScalingrulesDataSourceScalingGroupId = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_scaling_group_id" {
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}"
}
`
const testAccCheckAlicloudScalingrulesDataSourceScalingGroupIdNotFound = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_scaling_group_id" {
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}-fake"
}
`

const testAccCheckAlicloudScalingrulesDataSourceNameRegex = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_scaling_group_name_regex" {
	name_regex = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_name}"
}
`

const testAccCheckAlicloudScalingrulesDataSourceNameRegexNotFound = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_scaling_group_name_regex" {
	name_regex = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_name}-fake"
}
`

const testAccCheckAlicloudScalingrulesDataSourceCombined = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_combined"{
	name_regex = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_name}"
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}"
	type = "SimpleScalingRule"
}
`
const testAccCheckAlicloudScalingrulesDataSourceCombinedNotFound = testAccCheckAlicloudScalingrulesBasicConfig + `
data "alicloud_ess_scaling_rules" "foo_combined"{
	name_regex = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_rule_name}-fake"
	scaling_group_id = "${alicloud_ess_scaling_rule.scaling_rule1.scaling_group_id}"
	type = "SimpleScalingRule"
}
`

const testAccCheckAlicloudScalingrulesBasicConfig = EcsInstanceCommonTestCase + `

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

resource "alicloud_ess_scaling_rule" "scaling_rule1"{
	scaling_group_id = "${alicloud_ess_scaling_group.scalinggroup_foo1.id}"
	scaling_rule_name = "tf-testAccDataSourceScalingRule1"
	adjustment_type = "QuantityChangeInCapacity"
	adjustment_value = 1
}

`
