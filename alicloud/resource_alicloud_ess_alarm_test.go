package alicloud

import (
	"fmt"
	"log"
	"testing"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEssAlarm_basic(t *testing.T) {
	var alarm ess.Alarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_alarm.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAlarm_basic(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists("alicloud_ess_alarm.foo", &alarm),
					resource.TestMatchResourceAttr("alicloud_ess_alarm.foo", "name", regexp.MustCompile("^tf-testAccEssAlarm_basic-*")),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "metric_type", "system"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "metric_name", "CpuUtilization"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "period", "300"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "statistics", "Average"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "comparison_operator", ">="),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "evaluation_count", "2"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "threshold", "200.3"),
				),
			},
		},
	})
}

func TestAccAlicloudEssAlarm_with_dimension(t *testing.T) {
	var alarm ess.Alarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_alarm.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAlarm_with_dimension(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists("alicloud_ess_alarm.foo", &alarm),
					resource.TestMatchResourceAttr("alicloud_ess_alarm.foo", "name", regexp.MustCompile("^tf-testAccEssAlarm_with_dimension-*")),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "metric_type", "system"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "metric_name", "PackagesNetIn"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "period", "300"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "statistics", "Average"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "comparison_operator", ">="),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "evaluation_count", "2"),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "threshold", "200.3"),
				),
			},
		},
	})
}

func TestAccAlicloudEssAlarm_update(t *testing.T) {
	var alarm ess.Alarm

	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_alarm.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAlarm(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists(
						"alicloud_ess_alarm.foo", &alarm),
					resource.TestMatchResourceAttr(
						"alicloud_ess_alarm.foo",
						"name",
						regexp.MustCompile("^tf-testAccEssAlarm_update-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"metric_type",
						"system"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"metric_name",
						"CpuUtilization"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"period",
						"300"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"statistics",
						"Average"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"comparison_operator",
						">="),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"evaluation_count",
						"2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"threshold",
						"200.3"),
				),
			},

			{
				Config: testAccEssAlarm_update(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists(
						"alicloud_ess_alarm.foo", &alarm),
					resource.TestMatchResourceAttr(
						"alicloud_ess_alarm.foo",
						"name",
						regexp.MustCompile("^tf-testAccEssAlarm_update_new-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"description",
						"Acc alarm test update"),
				),
			},
		},
	})
}

func testAccCheckEssAlarmExists(n string, d *ess.Alarm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Alarm ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		attr, err := essService.DescribeEssAlarmById(rs.Primary.ID)
		log.Printf("[DEBUG] check ess alarm %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssAlarmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_alarm" {
			continue
		}
		if _, err := essService.DescribeEssAlarmById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Ess alarm %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssAlarm_basic(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAlarm_basic-%d"
	}
	resource "alicloud_vswitch" "bar" {
		name = "${var.name}_bar"
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}

	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_rule_name = "${var.name}"
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

	resource "alicloud_ess_alarm" "foo" {
	    name = "${var.name}"
	    description = "Acc alarm test"
	    alarm_actions = ["${alicloud_ess_scaling_rule.foo.ari}"]
	    scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	    metric_type = "system"
	    metric_name = "CpuUtilization"
	    period = 300
	    statistics = "Average"
	    threshold = 200.3
	    comparison_operator = ">="
		evaluation_count = 2
	}
	`, common, rand)
}

func testAccEssAlarm_with_dimension(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAlarm_with_dimension-%d"
	}
	resource "alicloud_vswitch" "bar" {
		name = "${var.name}"
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}

	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_rule_name = "${var.name}"
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

	resource "alicloud_ess_alarm" "foo" {
	    name = "${var.name}"
	    description = "Acc alarm test"
	    alarm_actions = ["${alicloud_ess_scaling_rule.foo.ari}"]
	    scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	    metric_type = "system"
	    metric_name = "PackagesNetIn"
	    period = 300
	    statistics = "Average"
	    threshold = 200.3
	    comparison_operator = ">="
		evaluation_count = 2
		dimensions = {
			device = "eth0"
		    }
	}
	`, common, rand)
}

func testAccEssAlarm(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAlarm_update-%d"
	}
	resource "alicloud_vswitch" "bar" {
		name = "${var.name}"
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}

	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_rule_name = "${var.name}"
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

	resource "alicloud_ess_alarm" "foo" {
	    name = "${var.name}"
	    description = "Acc alarm test"
	    alarm_actions = ["${alicloud_ess_scaling_rule.foo.ari}"]
	    scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	    metric_type = "system"
	    metric_name = "CpuUtilization"
	    period = 300
	    statistics = "Average"
	    threshold = 200.3
	    comparison_operator = ">="
		    evaluation_count = 2
	}
	`, common, rand)
}
func testAccEssAlarm_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAlarm_update_new-%d"
	}

	resource "alicloud_vswitch" "bar" {
		name = "${var.name}_bar"
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}

	resource "alicloud_ess_scaling_rule" "foo" {
		scaling_rule_name = "${var.name}"
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

	resource "alicloud_ess_alarm" "foo" {
	    name = "${var.name}"
	    description = "Acc alarm test update"
	    alarm_actions = ["${alicloud_ess_scaling_rule.foo.ari}"]
	    scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	    metric_type = "system"
	    metric_name = "CpuUtilization"
	    period = 300
	    statistics = "Average"
	    threshold = 200.3
	    comparison_operator = ">="
	    evaluation_count = 2
	}
	`, common, rand)
}
