package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"regexp"
	"testing"
)

func TestAccAlicloudEssScalingConfiguration_basic(t *testing.T) {
	var sc ess.ScalingConfigurationItemType

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
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_type",
						"ecs.n4.large"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"key_name",
						"ess_testcase_for_scaling_configuration"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_multiConfig(t *testing.T) {
	var sc ess.ScalingConfigurationItemType

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
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"instance_type",
						"ecs.n4.large"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"image_id",
						regexp.MustCompile("^centos_6")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"key_name",
						"ess_testcase_for_scaling_configuration"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"role_name",
						"EcsRamRoleForEssTest"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_active(t *testing.T) {
	var sc ess.ScalingConfigurationItemType

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
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_type",
						"ecs.n4.large"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_inactive(t *testing.T) {
	var sc ess.ScalingConfigurationItemType

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
				Config: testAccEssScalingConfiguration_inActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"active",
						"false"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_type",
						"ecs.n4.large"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
				),
			},
		},
	})
}

// Skip test enable security group reseult
func TestAccAlicloudEssScalingConfiguration_enable(t *testing.T) {
	var sc ess.ScalingConfigurationItemType

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
				Config: testAccEssScalingConfiguration_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"enable",
						"true"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_type",
						"ecs.n4.large"),
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
	var sc ess.ScalingConfigurationItemType

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
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_type",
						"ecs.n4.large"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos_6")),
				),
			},
		},
	})
}

func testAccCheckEssScalingConfigurationExists(n string, d *ess.ScalingConfigurationItemType) resource.TestCheckFunc {
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

		if attr == nil {
			return fmt.Errorf("Scaling Configuration not found")
		}

		*d = *attr
		return nil
	}
}

func testAccCheckEssScalingConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}
		ins, err := client.DescribeScalingConfigurationById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error ESS scaling configuration still exist")
		}

		// Verify the error is what we want
		if err != nil {
			// Verify the error is what we want
			e, _ := err.(*common.Error)
			if e.Code == InstanceNotFound {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccEssScalingConfigurationConfig = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_security_group" "tf_test_foo" {
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
	scaling_group_name = "test-scaling-configuration"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	force_delete = true
}

resource "alicloud_key_pair" "key" {
  key_name = "ess_testcase_for_scaling_configuration"
}
`

const testAccEssScalingConfiguration_multiConfig = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_instance_type = "ecs.n4.large"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  name = "test-for-ess"
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
	vpc_id = "${alicloud_vpc.vpc.id}"
	description = "foo"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "test-scaling-configuration-multi"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_id = "${alicloud_vswitch.vswitch.id}"

}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	role_name = "${alicloud_ram_role.role.id}"
	force_delete = true
}

resource "alicloud_ess_scaling_configuration" "bar" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	key_name = "${alicloud_key_pair.key.id}"
	role_name = "${alicloud_ram_role.role.id}"
	force_delete = true
}
resource "alicloud_key_pair" "key" {
  key_name = "ess_testcase_for_scaling_configuration"
}

resource "alicloud_ram_role" "role" {
  name = "EcsRamRoleForEssTest"
  services = ["ecs.aliyuncs.com"]
  description = "Test role for ECS and access to OSS."
  force = true
}

resource "alicloud_ram_policy" "policy" {
  name = "EcsRamRolePolicyTest"
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
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_security_group" "tf_test_foo" {
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
	scaling_group_name = "test-scaling-configuration-active"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	active = true

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`

const testAccEssScalingConfiguration_inActive = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_security_group" "tf_test_foo" {
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
	scaling_group_name = "test-scaling-configuration-inactive"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	active = false

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	substitute = "${alicloud_ess_scaling_configuration.bar.id}"
}
resource "alicloud_ess_scaling_configuration" "bar" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`

const testAccEssScalingConfiguration_enable = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_security_group" "tf_test_foo" {
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
	scaling_group_name = "test-scaling-configuration-able"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	active = true
	enable = true

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`

const testAccEssScalingConfiguration_disable = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_security_group" "tf_test_foo" {
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
	scaling_group_name = "test-scaling-configuration-disable"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	enable = false

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.large"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
}
`
