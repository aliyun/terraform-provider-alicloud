package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Config: testAccEssAlarm_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists("alicloud_ess_alarm.foo", &alarm),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "name", "tf-testAccEssAlarm_basic"),
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
				Config: testAccEssAlarm_with_dimension,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists("alicloud_ess_alarm.foo", &alarm),
					resource.TestCheckResourceAttr("alicloud_ess_alarm.foo", "name", "tf-testAccEssAlarm_with_dimension"),
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
				Config: testAccEssAlarm,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists(
						"alicloud_ess_alarm.foo", &alarm),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"name",
						"tf-testAccEssAlarm_update"),
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
				Config: testAccEssAlarm_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAlarmExists(
						"alicloud_ess_alarm.foo", &alarm),
					resource.TestCheckResourceAttr(
						"alicloud_ess_alarm.foo",
						"name",
						"tf-testAccEssAlarm_update_new"),
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

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeEssAlarmById(rs.Primary.ID)
		log.Printf("[DEBUG] check ess alarm %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssAlarmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_alarm" {
			continue
		}
		if _, err := client.DescribeEssAlarmById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Ess alarm %s still exists.", rs.Primary.ID)
	}
	return nil
}

const testAccEssAlarm_basic = `

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  	name = "tf-testAccEssAlarm_basic"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
	name = "tf-testAccEssAlarm_basic_foo"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
	name = "tf-testAccEssAlarm_basic_bar"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "tf-testAccEssAlarm_basic"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_rule_name = "tf-testAccEssAlarm_basic"
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 2
	cooldown = 60
}

resource "alicloud_ess_alarm" "foo" {
	name = "tf-testAccEssAlarm_basic"
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
`

const testAccEssAlarm_with_dimension = `

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  	name = "tf-testAccEssAlarm_with_dimension"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
	name = "tf-testAccEssAlarm_with_dimension_foo"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
	name = "tf-testAccEssAlarm_with_dimension_bar"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "tf-testAccEssAlarm_with_dimension"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_rule_name = "tf-testAccEssAlarm_with_dimension"
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 2
	cooldown = 60
}

resource "alicloud_ess_alarm" "foo" {
	name = "tf-testAccEssAlarm_with_dimension"
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
`

const testAccEssAlarm = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  	name = "tf-testAccEssAlarm_update"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
	name = "tf-testAccEssAlarm_update_foo"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
	name = "tf-testAccEssAlarm_update_bar"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "tf-testAccEssAlarm_update"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_rule_name = "tf-testAccEssAlarm_update"
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 2
	cooldown = 60
}

resource "alicloud_ess_alarm" "foo" {
	name = "tf-testAccEssAlarm_update"
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
`
const testAccEssAlarm_update = `

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  	name = "tf-testAccEssAlarm_update"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
	name = "tf-testAccEssAlarm_update_foo"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
	name = "tf-testAccEssAlarm_update_bar"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "tf-testAccEssAlarm_update"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_rule_name = "tf-testAccEssAlarm_update"
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 2
	cooldown = 60
}

resource "alicloud_ess_alarm" "foo" {
	name = "tf-testAccEssAlarm_update_new"
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
`
