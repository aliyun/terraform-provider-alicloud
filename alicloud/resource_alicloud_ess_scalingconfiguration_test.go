package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssScalingConfiguration_basic(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"key_name",
						"testAccEssScalingConfigurationConfig"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_multiConfig(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.bar",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingConfiguration_multiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.bar", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"active",
						"false"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"image_id",
						regexp.MustCompile("^centos_6")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"key_name",
						"testAccEssScalingConfiguration-multi"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"role_name",
						"testAccEssScalingConfiguration-multi"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_active(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingConfiguration_active,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"active",
						"true"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_disable(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingConfiguration_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"enable",
						"false"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
				),
			},
		},
	})
}

func testAccCheckEssScalingConfigurationExists(n string, d *ess.ScalingConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Scaling Configuration ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeScalingConfigurationById(rs.Primary.ID)
		log.Printf("[DEBUG] check scaling configuration %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}
		_, err := client.DescribeScalingConfigurationById(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling configuration %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAccEssScalingConfigurationConfig = `
provider "alicloud" {
  region = "cn-qingdao"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_ssd"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingConfigurationConfig"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	force_delete = true
}

resource "alicloud_key_pair" "key" {
  key_name = "${var.name}"
}
`

const testAccEssScalingConfiguration_multiConfig = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingConfiguration-multi"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "test-for-ess"
}

resource "alicloud_security_group" "tf_test_foo" {
  	name = "${var.name}"
	vpc_id = "${alicloud_vpc.vpc.id}"
	description = "foo"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]

}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	role_name = "${alicloud_ram_role.role.id}"
	force_delete = true
}

resource "alicloud_ess_scaling_configuration" "bar" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	role_name = "${alicloud_ram_role.role.id}"
	force_delete = true
}
resource "alicloud_key_pair" "key" {
  key_name = "${var.name}"
}

resource "alicloud_ram_role" "role" {
  name = "${var.name}"
  services = ["ecs.aliyuncs.com"]
  description = "Test role for ECS and access to OSS."
  force = true
}

resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Allow"
      action = ["oss:Get", "oss:List", "ecs:*"]
      resource = [ "*" ]
    }
  ]
  description = "Test role policy for ECS and access to OSS."
  force = true
}
resource "alicloud_ram_role_policy_attachment" "role-policy" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
`

const testAccEssScalingConfiguration_active = `
provider "alicloud" {
  region = "eu-central-1"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingConfiguration_active"
}
resource "alicloud_vpc" "vpc" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
	vpc_id = "${alicloud_vpc.vpc.id}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "test-for-ess"
}

resource "alicloud_security_group" "tf_test_foo" {
	vpc_id = "${alicloud_vpc.vpc.id}"
  	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	active = true

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`

const testAccEssScalingConfiguration_disable = `
provider "alicloud" {
  region = "cn-huhehaote"
}
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingConfiguration_disable"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "test-for-ess"
}

resource "alicloud_security_group" "tf_test_foo" {
  	name = "${var.name}"
	vpc_id = "${alicloud_vpc.vpc.id}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.vswitch.id}"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = false

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`
