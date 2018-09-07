package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_scalinggroup", &resource.Sweeper{
		Name: "alicloud_ess_scalinggroup",
		F:    testSweepEssGroups,
	})
}

func testSweepEssGroups(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var groups []ess.ScalingGroup
	req := ess.CreateDescribeScalingGroupsRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.essconn.DescribeScalingGroups(req)
		if err != nil {
			return fmt.Errorf("Error retrieving Scaling groups: %s", err)
		}
		if resp == nil || len(resp.ScalingGroups.ScalingGroup) < 1 {
			break
		}
		groups = append(groups, resp.ScalingGroups.ScalingGroup...)

		if len(resp.ScalingGroups.ScalingGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range groups {
		name := v.ScalingGroupName
		id := v.ScalingGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Scaling Group: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Scaling Group: %s (%s)", name, id)
		req := ess.CreateDeleteScalingGroupRequest()
		req.ScalingGroupId = id
		req.ForceDelete = requests.NewBoolean(true)
		if _, err := conn.essconn.DeleteScalingGroup(req); err != nil {
			log.Printf("[ERROR] Failed to delete Scaling Group (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(2 * time.Minute)
	}
	return nil
}

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
				Config: testAccEssScalingGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "tf-testAccEssScalingGroup"),
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
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "tf-testAccEssScalingGroup_update"),
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
						"alicloud_ess_scaling_group.foo", "scaling_group_name", "tf-testAccEssScalingGroup_vpc"),
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
	var slb slb.DescribeLoadBalancerAttributeResponse

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

func TestAccAlicloudEssScalingGroup_slbempty(t *testing.T) {
	var sg ess.ScalingGroup

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
				Config: testAccEssScalingGroup_slbempty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.scaling", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "loadbalancer_ids.#", "0"),
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

		if _, err := client.DescribeScalingGroupById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling group %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAccEssScalingGroup = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "tf-testAccEssScalingGroup"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}-bar"
}

resource "alicloud_security_group" "tf_test_foo" {
 	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
	vswitch_ids = ["${alicloud_vswitch.foo.id}", "${alicloud_vswitch.bar.id}"]
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = true
	active = true
	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 10
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}

`

const testAccEssScalingGroup_update = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "tf-testAccEssScalingGroup"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}-bar"
}

resource "alicloud_security_group" "tf_test_foo" {
 	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 2
	max_size = 2
	scaling_group_name = "${var.name}_update"
	default_cooldown = 20
	vswitch_ids = ["${alicloud_vswitch.foo.id}", "${alicloud_vswitch.bar.id}"]
	removal_policies = ["OldestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = true
	active = true
	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 10
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
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

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "tf-testAccEssScalingGroup_vpc"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
        name = "${var.name}-foo"
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}-bar"
}

resource "alicloud_security_group" "tf_test_foo" {
 	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 2
	scaling_group_name = "${var.name}"
	default_cooldown = 20
	vswitch_ids = ["${alicloud_vswitch.foo.id}", "${alicloud_vswitch.bar.id}"]
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = true
	active = true
	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 10
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}
`

const testAccEssScalingGroup_slb = `
variable "name" {
	default = "tf-testAccEssScalingGroup_slb"
}
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "sg" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_ess_scaling_group" "scaling" {
  min_size = "1"
  max_size = "1"
  scaling_group_name = "${var.name}"
  removal_policies = ["OldestInstance", "NewestInstance"]
  vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]
  loadbalancer_ids = ["${alicloud_slb.instance.0.id}","${alicloud_slb.instance.1.id}"]
  depends_on = ["alicloud_slb_listener.tcp"]
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  active = true
  enable = true
  image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_group_id = "${alicloud_security_group.sg.id}"
  force_delete = "true"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_slb" "instance" {
  count=2
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
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

const testAccEssScalingGroup_slbempty = `
variable "name" {
	default = "tf-testAccEssScalingGroup_slbempty"
}
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "sg" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_ess_scaling_group" "scaling" {
  min_size = "1"
  max_size = "1"
  scaling_group_name = "${var.name}"
  removal_policies = ["OldestInstance", "NewestInstance"]
  vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]
  loadbalancer_ids = []
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  active = true
  enable = true
  image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_group_id = "${alicloud_security_group.sg.id}"
  force_delete = "true"
  internet_charge_type = "PayByTraffic"
}
`
