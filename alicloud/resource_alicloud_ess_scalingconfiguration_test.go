package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccEssScalingConfigurationConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^ubuntu_14")),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"key_name",
						regexp.MustCompile("^tf-testAccEssScalingConfigurationConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"user_data",
						"#!/bin/bash\necho \"hello\"\n"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.#",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.delete_with_instance",
						"false"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.size",
						"20"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_SystemDisk(t *testing.T) {
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
			{
				Config: testAccEssScalingConfigurationConfig_SystemDisk(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^ubuntu_14")),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"key_name",
						regexp.MustCompile("^tf-testAccEssScalingConfigurationConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"user_data",
						"#!/bin/bash\necho \"hello\"\n"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"system_disk_size",
						"40"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.#",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.delete_with_instance",
						"false"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.size",
						"20"),
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
			{
				Config: testAccEssScalingConfiguration_multiConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
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
						regexp.MustCompile("^ubuntu_14")),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"key_name",
						regexp.MustCompile("^tf-testAccEssScalingConfiguration-multi-*")),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.bar",
						"role_name",
						regexp.MustCompile("^tf-testAccEssScalingConfiguration-multi-*")),
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
			{
				Config: testAccEssScalingConfiguration_active(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
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
						regexp.MustCompile("^ubuntu_14")),
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
			{
				Config: testAccEssScalingConfiguration_disable(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
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
						regexp.MustCompile("^ubuntu_14")),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_modify(t *testing.T) {
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
			{
				Config: testAccEssScalingConfigurationConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
				),
			},
			{
				Config: testAccScalingConfigurationModify,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"image_id",
						regexp.MustCompile("^centos")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"scaling_configuration_name",
						"tf-testAccEssConfiguration-modify"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"internet_charge_type",
						"PayByBandwidth"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"internet_max_bandwidth_out",
						"5"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"system_disk_category",
						"cloud_ssd"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"system_disk_size",
						"50"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"role_name",
						"tf-testAccEssConfiguration-modify"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"key_name",
						"tf-testAccEssConfiguration-modify"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"instance_name",
						"tf-testAccEssConfiguration-modify"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"tags.name",
						"tf-test"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.#",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.size",
						"20"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.category",
						"cloud_ssd"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_configuration.foo",
						"data_disk.0.delete_with_instance",
						"false"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_multi_sgs(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EssScalingConfigurationMultiSgSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingConfigurationConfig_multi_sg(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr("alicloud_ess_scaling_configuration.foo",
						"security_group_ids.#",
						"2"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfiguration_modify_single_2_multi_sg(t *testing.T) {
	var sc ess.ScalingConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EssScalingConfigurationMultiSgSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_configuration.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScalingConfiguration_multi_sg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttr("alicloud_ess_scaling_configuration.foo",
						"security_group_ids.#",
						"2"),
				),
			},
			{
				Config: testAccScalingConfiguration_single_sg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingConfigurationExists(
						"alicloud_ess_scaling_configuration.foo", &sc),
					resource.TestCheckResourceAttrSet("alicloud_ess_scaling_configuration.foo",
						"security_group_id"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		attr, err := essService.DescribeScalingConfigurationById(rs.Primary.ID)
		log.Printf("[DEBUG] check scaling configuration %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}
		_, err := essService.DescribeScalingConfigurationById(rs.Primary.ID)

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

func testAccEssScalingConfigurationConfig_multi_sg(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingConfigurationConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_security_group" "scaling_conf_default" {
	  name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_ids = ["${alicloud_security_group.default.id}","${alicloud_security_group.scaling_conf_default.id}"]
		force_delete = true
	}
	`, common, rand)
}

func testAccEssScalingConfigurationConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingConfigurationConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		key_name = "${alicloud_key_pair.key.id}"
		force_delete = true
		data_disk = [
		{
			size = 20
			category = "cloud_efficiency"
			delete_with_instance = false
		}
	]
		user_data = <<EOF
#!/bin/bash
echo "hello"
EOF
	}
	
	resource "alicloud_key_pair" "key" {
	  key_name = "${var.name}"
	}
	`, common, rand)
}

func testAccEssScalingConfigurationConfig_SystemDisk(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingConfigurationConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		key_name = "${alicloud_key_pair.key.id}"
		system_disk_size = 40
		force_delete = true
		data_disk = [
		{
			size = 20
			category = "cloud_efficiency"
			delete_with_instance = false
		}
	]
		user_data = <<EOF
#!/bin/bash
echo "hello"
EOF
	}

	resource "alicloud_key_pair" "key" {
	  key_name = "${var.name}"
	}
	`, common, rand)
}

func testAccEssScalingConfiguration_multiConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingConfiguration-multi-%d"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]

	}
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		key_name = "${alicloud_key_pair.key.id}"
		role_name = "${alicloud_ram_role.role.id}"
		force_delete = true
	}

	resource "alicloud_ess_scaling_configuration" "bar" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

		image_id = "${data.alicloud_images.default.images.0.id}"
		  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
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
	`, common, rand)
}

func testAccEssScalingConfiguration_active(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssScalingConfiguration_active-%d"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		active = true
	
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = true
	}
	`, common, rand)
}

func testAccEssScalingConfiguration_disable(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssConfiguration_disable-%d"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		enable = false
	
		image_id = "${data.alicloud_images.default.images.0.id}"
		  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = true
	}
	`, common, rand)
}

const testAccScalingConfigurationModify = EcsInstanceCommonTestCase + `

variable "name" {
		default = "tf-testAccEssConfiguration-modify"
	}

resource "alicloud_ess_scaling_group" "scaling_conf_foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

data "alicloud_images" "scaling_conf_foo" {
  name_regex  = "^centos.*_64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_security_group" "scaling_conf_foo" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_key_pair" "scaling_conf_foo" {
  key_name = "${var.name}"
}

resource "alicloud_ram_role" "scaling_conf_foo" {
	  name = "${var.name}"
	  services = ["ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}

resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.scaling_conf_foo.id}"
		image_id = "${data.alicloud_images.scaling_conf_foo.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.1.id}"
		security_group_id = "${alicloud_security_group.scaling_conf_foo.id}"
		scaling_configuration_name = "${var.name}"
		internet_charge_type = "PayByBandwidth"
		internet_max_bandwidth_out = 5
		system_disk_category = "cloud_ssd"
		system_disk_size = "50"
		role_name = "${alicloud_ram_role.scaling_conf_foo.id}"
		key_name = "${alicloud_key_pair.scaling_conf_foo.id}"
		instance_name = "${var.name}"
		tags = {
			name = "tf-test"
		    }
		force_delete = true
		data_disk = [
		{
			size = 20
			category = "cloud_ssd"
			delete_with_instance = false
		}
	]
		user_data = <<EOF
#!/bin/bash
echo "world"
EOF
	}

`

const testAccScalingConfiguration_multi_sg = EcsInstanceCommonTestCase + `

variable "name" {
		default = "tf-testAccEssConfiguration_multi_sg"
	}

resource "alicloud_ess_scaling_group" "scaling_conf_foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

resource "alicloud_security_group" "scaling_conf_foo" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.scaling_conf_foo.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_ids = ["${alicloud_security_group.scaling_conf_foo.id}","${alicloud_security_group.default.id}"]
		scaling_configuration_name = "${var.name}"
		force_delete = true
	}
`

const testAccScalingConfiguration_single_sg = EcsInstanceCommonTestCase + `

variable "name" {
		default = "tf-testAccEssConfiguration_multi_sg"
	}

resource "alicloud_ess_scaling_group" "scaling_conf_foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.scaling_conf_foo.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		scaling_configuration_name = "${var.name}"
		force_delete = true
	}
`
