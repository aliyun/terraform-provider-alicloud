package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
	"time"
)

func TestAccAliCloudEssInstanceRefresh_image(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_instance_refresh.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	var v map[string]interface{}
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssInstanceRefresh-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssInstanceRefreshImage)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssInstanceRefreshDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":                     []string{"alicloud_ess_scaling_configuration.default"},
					"scaling_group_id":               "${alicloud_ess_scaling_group.default.id}",
					"desired_configuration_image_id": "${data.alicloud_images.default2.images[0].id}",
					"min_healthy_percentage":         "90",
					"max_healthy_percentage":         "150",
					"checkpoint_pause_time":          "1",
					"skip_matching":                  "false",
					"checkpoints": []map[string]string{
						{
							"percentage": "50",
						},
						{
							"percentage": "70",
						},
						{
							"percentage": "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Cancelled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Cancelled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEssInstanceRefresh_template(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_instance_refresh.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	var v map[string]interface{}
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssInstanceRefresh-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssInstanceRefreshTemplate)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssInstanceRefreshDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"desired_configuration_launch_template_id":      "${alicloud_ecs_launch_template.default.id}",
					"desired_configuration_launch_template_version": "Default",
					"desired_configuration_launch_template_overrides": []map[string]string{
						{
							"instance_type": "${data.alicloud_instance_types.default1.instance_types.0.id}",
						},
						{
							"instance_type": "${data.alicloud_instance_types.default1.instance_types.1.id}",
						},
					},
					"min_healthy_percentage": "90",
					"max_healthy_percentage": "150",
					"checkpoint_pause_time":  "1",
					"skip_matching":          "false",
					"checkpoints": []map[string]string{
						{
							"percentage": "50",
						},
						{
							"percentage": "60",
						},
						{
							"percentage": "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEssInstanceRefresh_container(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_instance_refresh.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	var v map[string]interface{}
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssInstanceRefresh-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssInstanceRefreshContainer)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssInstanceRefreshDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":             []string{"alicloud_ess_eci_scaling_configuration.default"},
					"scaling_group_id":       "${alicloud_ess_scaling_group.default.id}",
					"min_healthy_percentage": "90",
					"max_healthy_percentage": "150",
					"checkpoint_pause_time":  "1",
					"skip_matching":          "false",
					"checkpoints": []map[string]string{
						{
							"percentage": "50",
						},
						{
							"percentage": "60",
						},
						{
							"percentage": "100",
						},
					},
					"desired_configuration_containers": []map[string]interface{}{
						{
							"args":     []string{"999999", "8888"},
							"commands": []string{"sleep", "echo 1"},
							"name":     "container-1",
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/tomcat:8",
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "PATH",
									"value":                "/usr/local/bin",
									"field_ref_field_path": "fieldPath",
								},
								{
									"key":                  "PATH1",
									"value":                "/usr/local/bin1",
									"field_ref_field_path": "fieldPath1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RollbackSuccessful",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RollbackSuccessful",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccCheckEssInstanceRefreshDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_instance_refresh" {
			continue
		}
		wait := incrementalWait(1*time.Second, 2*time.Second)

		var raw map[string]interface{}
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = essService.DescribeEssInstanceRefresh(rs.Primary.ID)
			if err != nil {
				wait()
				return resource.RetryableError(err)
			}
			if raw["Status"] != "cancel" && raw["Status"] != "suspend" && raw["Status"] != "resume" && raw["Status"] != "rollback" {
				wait()
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	return nil
}

func testAccEssInstanceRefreshContainer(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "10"
      desired_capacity = "1"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
      group_type = "ECI"
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
    data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_instance_types" "sn1" {
	 instance_type_family = "ecs.sn1"
     availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
    data "alicloud_images" "default2" {
		name_regex  = "^centos"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_eci_scaling_configuration" "default" {
          depends_on = [
        "alicloud_ess_scaling_group.default"]
		scaling_group_id = alicloud_ess_scaling_group.default.id
        security_group_id = alicloud_security_group.default.id
        container_group_name = "test"
        cpu = "2"
        memory = "4"
        containers {
			name = "container-1"
            cpu = "2"
            memory = "4"
	        image = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/redis:5.0.1-alpine"
		}
		force_delete = true
		active = true
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssInstanceRefreshTemplate(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "10"
      desired_capacity = "1"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
    data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^ubu"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_instance_types" "sn1" {
	 instance_type_family = "ecs.sn1"
     availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
    data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
          depends_on = [
        "alicloud_ess_scaling_group.default"]
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default1.images[0].id
		instance_type = data.alicloud_instance_types.default1.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}
    resource "alicloud_ecs_launch_template" "default" {
		launch_template_name = "tf-test12145"
		image_id             =  data.alicloud_images.default2.images.0.id
		instance_charge_type = "PrePaid"
		instance_type        =  data.alicloud_instance_types.sn1.instance_types.0.id
		internet_charge_type          = "PayByBandwidth"
		internet_max_bandwidth_in     = "5"
		internet_max_bandwidth_out    = "0"
		io_optimized                  = "optimized"
		network_type                  = "vpc"
		security_enhancement_strategy = "Active"
		spot_price_limit              = "5"
        spot_strategy                 = "SpotWithPriceLimit"
		security_group_id             = alicloud_security_group.default.id
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssInstanceRefreshImage(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "10"
      desired_capacity = "1"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
    data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^ubu"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
          depends_on = [
        "alicloud_ess_scaling_group.default"]
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default1.images[0].id
		instance_type = data.alicloud_instance_types.default1.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}
	`, EcsInstanceCommonTestCase, name)
}
