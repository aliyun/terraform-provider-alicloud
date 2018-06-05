package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssScalingGroup_basic(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "sg-for-test-config"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingGroup_update(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "sg-for-test"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "2"),
				),
			},

			resource.TestStep{
				Config: testAccEssScalingGroup_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "update"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_vpc(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroup_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "sg-for-test-vpc"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "vswitch_ids.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_slb(t *testing.T) {
	var sg ess.ScalingGroup
	var slb slb.LoadBalancerType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.scaling",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroup_slb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.scaling", &sg),
					testAccCheckSlbExists(
						"alicloud_slb.instance.0", &slb),
					testAccCheckSlbExists(
						"alicloud_slb.instance.1", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "loadbalancer_ids.#", "2"),
				),
			},
		},
	})

}

func testAccCheckEssScalingGroupExists(n string, d *ess.ScalingGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Scaling Group ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeScalingGroupById(rs.Primary.ID)
		log.Printf("[DEBUG] check scaling group %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		_, err := client.DescribeScalingGroupById(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			// Verify the error is what we want
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccEssScalingGroupConfig = `
resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "sg-for-test-config"
	removal_policies = ["OldestInstance", "NewestInstance"]
}
`

const testAccEssScalingGroup = `
resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "sg-for-test"
	removal_policies = ["OldestInstance", "NewestInstance"]
}
`

const testAccEssScalingGroup_update = `
resource "alicloud_ess_scaling_group" "foo" {
	min_size = 2
	max_size = 2
	scaling_group_name = "update"
	removal_policies = ["OldestInstance"]
}
`
const testAccEssScalingGroup_vpc = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  	name = "tf_test_foo"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 2
	scaling_group_name = "sg-for-test-vpc"
	default_cooldown = 20
	vswitch_ids = ["${alicloud_vswitch.foo.id}", "${alicloud_vswitch.bar.id}"]
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = true
	active = true
	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 10
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}
`

const testAccEssScalingGroup_slb = `
provider "alicloud" {
  region = "cn-hangzhou"
}
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "sg" {
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_ess_scaling_group" "scaling" {
  min_size = "1"
  max_size = "1"
  scaling_group_name = "sg-for-test-slb"
  removal_policies = ["OldestInstance", "NewestInstance"]
  vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]
  loadbalancer_ids = ["${alicloud_slb.instance.0.id}","${alicloud_slb.instance.1.id}"]
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  active = true
  enable = true
  image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  instance_type = "ecs.n4.small"
  security_group_id = "${alicloud_security_group.sg.id}"
  force_delete = "true"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_slb" "instance" {
  count=2
  name = "slb-for-ess"
  internet_charge_type = "paybytraffic"
  internet = false
}
resource "alicloud_slb_listener" "tcp" {
  count = 2
  load_balancer_id = "${element(alicloud_slb.instance.*.id, count.index)}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
}
`
