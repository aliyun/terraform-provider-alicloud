package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEssScalingRule_basic(t *testing.T) {
	var sc ess.ScalingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_rule.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"1"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingRule_update(t *testing.T) {
	var sc ess.ScalingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_rule.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRule(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"1"),
				),
			},

			{
				Config: testAccEssScalingRule_update(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"2"),
				),
			},
		},
	})
}

func testAccCheckEssScalingRuleExists(n string, d *ess.ScalingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Scaling Rule ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		ids := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		attr, err := essService.DescribeScalingRuleById(ids[0], ids[1])
		log.Printf("[DEBUG] check scaling rule %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_rule" {
			continue
		}
		ids := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		_, err := essService.DescribeScalingRuleById(ids[0], ids[1])

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling rule %s still exists.", ids[1])
	}

	return nil
}

func testAccEssScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "bar" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"

		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
		cooldown = 120
	}
	`, common, rand)
}

func testAccEssScalingRule(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "bar" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
	
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
		cooldown = 120
	}
	`, common, rand)
}

func testAccEssScalingRule_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "bar" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
	
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	`, common, rand)
}
