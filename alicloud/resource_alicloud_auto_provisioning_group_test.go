package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAutoProvisioningGroup(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	var basicMap = map[string]string{
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"spot_target_capacity":                       "2",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":            "${alicloud_launch_template.template.id}",
					"total_target_capacity":         "4",
					"pay_as_you_go_target_capacity": "1",
					"spot_target_capacity":          "2",
					"launch_template_config": []map[string]string{{
						"vswitch_id":        "${alicloud_vswitch.default.id}",
						"weighted_capacity": "1",
						"max_price":         "2",
						"instance_type":     "ecs.sn1ne.large",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "auto_provisioning_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "auto_provisioning_group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_provisioning_group_name": "auto_provisioning_group_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_provisioning_group_name": "auto_provisioning_group_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_spot_price": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_spot_price": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"excess_capacity_termination_policy": "termination",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"excess_capacity_termination_policy": "termination",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_target_capacity_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_target_capacity_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"terminate_instances_with_expiration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"terminate_instances_with_expiration": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"total_target_capacity": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_target_capacity": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pay_as_you_go_target_capacity": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_as_you_go_target_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_target_capacity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_target_capacity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_target_capacity":                "2",
					"pay_as_you_go_target_capacity":       "1",
					"total_target_capacity":               "4",
					"terminate_instances_with_expiration": "false",
					"default_target_capacity_type":        "Spot",
					"excess_capacity_termination_policy":  "no-termination",
					"auto_provisioning_group_name":        "auto_provisioning_group",
					"max_spot_price":                      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_target_capacity":                "2",
						"pay_as_you_go_target_capacity":       "1",
						"total_target_capacity":               "4",
						"terminate_instances_with_expiration": "false",
						"default_target_capacity_type":        "Spot",
						"excess_capacity_termination_policy":  "no-termination",
						"auto_provisioning_group_name":        "auto_provisioning_group",
						"max_spot_price":                      "2",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudAutoProvisioningGroup_valid(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	validFrom := time.Now().AddDate(0, 0, 1).UTC().Format("2006-01-02T15:04:05Z")
	validUntil := time.Now().AddDate(0, 0, 3).UTC().Format("2006-01-02T15:04:05Z")
	var basicMap = map[string]string{
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"spot_target_capacity":                       "2",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
		"valid_from":                                 validFrom,
		"valid_until":                                validUntil,
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":            "${alicloud_launch_template.template.id}",
					"total_target_capacity":         "4",
					"pay_as_you_go_target_capacity": "1",
					"spot_target_capacity":          "2",
					"valid_from":                    validFrom,
					"valid_until":                   validUntil,
					"launch_template_config": []map[string]string{{
						"vswitch_id":        "${alicloud_vswitch.default.id}",
						"weighted_capacity": "1",
						"max_price":         "2",
						"instance_type":     "ecs.sn1ne.large",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}
func resourceAutoProvisioningGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s

variable "name" {
  default = "%s"
}
resource "alicloud_launch_template" "template" {
  name              = "${var.name}"
  image_id          = "${data.alicloud_images.default.images.0.id}"
  instance_type     = "ecs.sn1ne.large"
  security_group_id = "${alicloud_security_group.default.id}"
}`, EcsInstanceCommonTestCase, name)
}

func TestAccAlicloudAutoProvisioningGroup_launchConfiguration(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccapglc%d", rand)
	resourceId := "alicloud_auto_provisioning_group.default"
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, resourceAttrInit(resourceId, map[string]string{
		"total_target_capacity":         "1",
		"pay_as_you_go_target_capacity": "1",
	}))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoProvisioningGroupLaunchConfigurationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					resource.TestCheckResourceAttr(resourceId, "launch_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "launch_configuration.0.image_id", "${data.alicloud_images.default.images.0.id}"),
					resource.TestCheckResourceAttr(resourceId, "launch_configuration.0.system_disk_size", "40"),
					resource.TestCheckResourceAttr(resourceId, "launch_configuration.0.security_group_id", "${alicloud_security_group.default.id}"),
				),
			},
		},
	})
}

func testAccAutoProvisioningGroupLaunchConfigurationConfig(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
  default = "%s"
}

resource "alicloud_auto_provisioning_group" "default" {
  total_target_capacity         = "1"
  pay_as_you_go_target_capacity = "1"

  launch_template_config {
    instance_type     = "ecs.sn1ne.large"
    vswitch_id        = alicloud_vswitch.default.id
    weighted_capacity = "1"
    max_price         = "2"
  }

  launch_configuration {
    image_id                       = data.alicloud_images.default.images.0.id
    image_family                   = ""
    instance_name                  = "tf-testAcc-apg-instance"
    instance_description           = "test acc instance description"
    host_name                      = "tf-test-host"
    host_names                     = []
    credit_specification           = "standard"
    deployment_set_id              = ""
    auto_release_time              = ""
    io_optimized                   = "optimized"
    security_group_id              = alicloud_security_group.default.id
    security_group_ids             = []
    security_enhancement_strategy  = "DeletionProtection"
    internet_max_bandwidth_in      = 5
    internet_max_bandwidth_out     = 5
    internet_charge_type           = "PayByTraffic"
    password                       = "TestAcc1234!"
    password_inherit               = false
    key_pair_name                  = ""
    ram_role_name                  = ""
    user_data                      = ""
    resource_group_id              = ""
    system_disk_category           = "cloud_efficiency"
    system_disk_size               = 40
    system_disk_performance_level  = "PL1"
    system_disk_name               = "tf-test-sys-disk"
    system_disk_description        = "test system disk"

    system_disk {
      encrypted          = "false"
      kms_key_id         = ""
      encrypt_algorithm  = ""
      provisioned_iops   = ""
      bursting_enabled   = ""
    }

    data_disk {
      category             = "cloud_efficiency"
      disk_name            = "tf-test-data-disk"
      size                 = "40"
      device               = "/dev/xvdb"
      snapshot_id          = ""
      description          = "test data disk"
      delete_with_instance = "true"
      encrypted            = "false"
      kms_key_id           = ""
      encrypt_algorithm    = ""
      performance_level    = "PL1"
      provisioned_iops     = ""
      bursting_enabled     = ""
    }

    network_interface {
      security_group_id   = alicloud_security_group.default.id
      security_group_ids  = []
      instance_type       = "ecs.sn1ne.large"
    }

    tag {
      key   = "test-key"
      value = "test-value"
    }

    arn {
      rolearn         = ""
      role_type       = ""
      assume_role_for = ""
    }

    period            = ""
    period_unit       = ""
    auto_renew        = ""
    auto_renew_period = ""

    additional_info {
      pvd_config = ""
    }
  }
}
`, EcsInstanceCommonTestCase, name)
}
