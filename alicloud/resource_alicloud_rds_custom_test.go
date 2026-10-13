package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Rds Custom. >>> Resource test cases, automatically generated.
// Case rdscustom_run_ins_extra_param 10836
func TestAccAliCloudRdsCustom_basic10836(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomMap10836)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomBasicDependence10836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"amount":        "1",
					"vswitch_id":    "${data.alicloud_vswitches.default.ids.0}",
					"auto_renew":    "false",
					"period":        "1",
					"auto_pay":      "true",
					"instance_type": "mysql.x2.xlarge.6cm",
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "50",
							"performance_level": "PL1",
						},
					},
					"status": "Running",
					"security_group_ids": []string{
						"${data.alicloud_security_groups.default.ids.0}"},
					"io_optimized":                  "optimized",
					"description":                   "ran_test_ram_code",
					"key_pair_name":                 "${alicloud_ecs_key_pair.KeyPairName.id}",
					"zone_id":                       "${var.test_zone_id}",
					"instance_charge_type":          "Prepaid",
					"internet_charge_type":          "PayByTraffic",
					"internet_max_bandwidth_out":    "0",
					"image_id":                      "aliyun_3_x64_20G_alibase_20260625.vhd",
					"security_enhancement_strategy": "Active",
					"period_unit":                   "Month",
					"password":                      "jingyiTEST@123",
					"system_disk": []map[string]interface{}{
						{
							"size":     "40",
							"category": "cloud_essd",
						},
					},
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"host_name":          "1750263391",
					"create_mode":        "0",
					"create_extra_param": "{}",
					"spot_strategy":      "NoSpot",
					"timeouts": []map[string]interface{}{
						{
							"create": "20m",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"amount":                        "1",
						"vswitch_id":                    CHECKSET,
						"auto_renew":                    "false",
						"period":                        "1",
						"auto_pay":                      "true",
						"instance_type":                 "mysql.x2.xlarge.6cm",
						"data_disk.#":                   "1",
						"status":                        "Running",
						"security_group_ids.#":          "1",
						"io_optimized":                  "optimized",
						"description":                   "ran_test_ram_code",
						"key_pair_name":                 CHECKSET,
						"zone_id":                       CHECKSET,
						"instance_charge_type":          "Prepaid",
						"internet_charge_type":          "PayByTraffic",
						"internet_max_bandwidth_out":    "0",
						"image_id":                      "aliyun_3_x64_20G_alibase_20260625.vhd",
						"security_enhancement_strategy": "Active",
						"period_unit":                   "Month",
						"password":                      "jingyiTEST@123",
						"resource_group_id":             CHECKSET,
						"host_name":                     CHECKSET,
						"create_mode":                   CHECKSET,
						"create_extra_param":            CHECKSET,
						"spot_strategy":                 "NoSpot",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "mysql.x4.xlarge.6cm",
					"status":            "Stopped",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"force_stop":        "false",
					"direction":         "Up",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "mysql.x4.xlarge.6cm",
						"status":            "Stopped",
						"resource_group_id": CHECKSET,
						"force_stop":        "false",
						"direction":         "Up",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "Running",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "Running",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "Stopped",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"force_stop":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "Stopped",
						"resource_group_id": CHECKSET,
						"force_stop":        "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_pay", "auto_renew", "create_extra_param", "create_mode", "direction", "dry_run", "force_stop", "host_name", "image_id", "instance_charge_type", "internet_charge_type", "internet_max_bandwidth_out", "io_optimized", "key_pair_name", "password", "period", "period_unit", "security_enhancement_strategy", "spot_strategy", "support_case"},
			},
		},
	})
}

var AlicloudRdsCustomMap10836 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudRdsCustomBasicDependence10836(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "test_region_id" {
  default = "cn-beijing"
}

variable "test_zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  ids = [data.alicloud_vswitches.default.vswitches.0.vpc_id]
}

data "alicloud_vswitches" "default" {
  status  = "Available"
  zone_id = var.test_zone_id
}

data "alicloud_security_groups" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecs_key_pair" "KeyPairName" {
  key_pair_name = var.name
}`, name)
}

// Test Rds Custom. <<< Resource test cases, automatically generated.

// Case resourceCase_20260415_01_clone_1_clone_0 12788
func TestAccAliCloudRdsCustom_basic12788(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomMap12788)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomBasicDependence12788)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "ran资源用名称050801",
					"instance_charge_type": "PostPaid",
					"auto_renew":           "false",
					"amount":               "1",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"dry_run":              "false",
					"deployment_set_id":    "${alicloud_rds_custom_deployment_set.deploymentSet.id}",
					"auto_pay":             "true",
					"force":                "false",
					"support_case":         "eni",
					"security_group_ids": []string{
						"${data.alicloud_security_groups.default.ids.0}"},
					"system_disk": []map[string]interface{}{
						{
							"category": "cloud_essd",
							"size":     "40",
						},
					},
					"instance_name": "namecreate",
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
						},
					},
					"create_mode":   "0",
					"instance_type": "mysql.x2.large.9cm",
					"spot_strategy": "NoSpot",
					"host_name":     "testhostNameran",
					"period_unit":   "Month",
					"password":      "@MY7yxqc9YXQC",
					"timeouts": []map[string]interface{}{
						{
							"create": "20m",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "ran资源用名称050801",
						"instance_charge_type": "PostPaid",
						"auto_renew":           "false",
						"amount":               "1",
						"vswitch_id":           CHECKSET,
						"dry_run":              "false",
						"deployment_set_id":    CHECKSET,
						"auto_pay":             "true",
						"force":                "false",
						"support_case":         "eni",
						"security_group_ids.#": "1",
						"instance_name":        "namecreate",
						"data_disk.#":          "1",
						"create_mode":          CHECKSET,
						"instance_type":        "mysql.x2.large.9cm",
						"spot_strategy":        "NoSpot",
						"host_name":            "testhostNameran",
						"period_unit":          "Month",
						"password":             "@MY7yxqc9YXQC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "修改描述",
					"force_stop":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "修改描述",
						"force_stop":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "update名称",
					"host_name":     "hostnameupdate",
					"password":      "@MY7yxqc9YXQCUPDATE",
					"force":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "update名称",
						"host_name":     "hostnameupdate",
						"password":      "@MY7yxqc9YXQCUPDATE",
						"force":         "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":     "Stopped",
					"force_stop": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":     "Stopped",
						"force_stop": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_pay", "auto_renew", "create_mode", "direction", "dry_run", "force", "force_stop", "host_name", "image_id", "instance_charge_type", "internet_charge_type", "internet_max_bandwidth_out", "io_optimized", "key_pair_name", "password", "period", "period_unit", "security_enhancement_strategy", "spot_strategy", "support_case"},
			},
		},
	})
}

var AlicloudRdsCustomMap12788 = map[string]string{
	"system_disk_id": CHECKSET,
	"region_id":      CHECKSET,
}

func AlicloudRdsCustomBasicDependence12788(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  ids = [data.alicloud_vswitches.default.vswitches.0.vpc_id]
}

data "alicloud_vswitches" "default" {
  status  = "Available"
  zone_id = "cn-beijing-i"
}

data "alicloud_security_groups" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_rds_custom_deployment_set" "deploymentSet" {
  custom_deployment_set_name            = var.name
  description                           = var.name
  group_count                           = 3
  on_unable_to_redeploy_failed_instance = "CancelMembershipAndStart"
  strategy                              = "Availability"
}`, name)
}

// Case resourceCase_20260415_01 12772
func TestAccAliCloudRdsCustom_basic12772(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomMap12772)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	privateIpAddress := fmt.Sprintf("172.16.3.%d", acctest.RandIntRange(10, 240))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomBasicDependence12772)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "ran资源用例描述042001",
					"private_ip_address":   privateIpAddress,
					"instance_charge_type": "PostPaid",
					"auto_renew":           "false",
					"amount":               "1",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"dry_run":              "false",
					"auto_pay":             "true",
					"security_group_ids": []string{
						"${data.alicloud_security_groups.default.ids.0}"},
					"system_disk": []map[string]interface{}{
						{
							"category": "cloud_essd",
							"size":     "40",
						},
					},
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
						},
					},
					"create_mode":   "0",
					"instance_type": "mysql.x2.large.9cm",
					"spot_strategy": "NoSpot",
					"host_name":     "testhostNameran",
					"period_unit":   "Month",
					"password":      "@MY7yxqc9YXQC",
					"timeouts": []map[string]interface{}{
						{
							"create": "20m",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "ran资源用例描述042001",
						"private_ip_address":   privateIpAddress,
						"instance_charge_type": "PostPaid",
						"auto_renew":           "false",
						"amount":               "1",
						"vswitch_id":           CHECKSET,
						"dry_run":              "false",
						"auto_pay":             "true",
						"security_group_ids.#": "1",
						"data_disk.#":          "1",
						"create_mode":          CHECKSET,
						"instance_type":        "mysql.x2.large.9cm",
						"spot_strategy":        "NoSpot",
						"host_name":            "testhostNameran",
						"period_unit":          "Month",
						"password":             "@MY7yxqc9YXQC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "修改描述",
					"force_stop":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "修改描述",
						"force_stop":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password":      "@MY7yxqc9YXQCUPDATE",
					"instance_name": "update名称",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password":      "@MY7yxqc9YXQCUPDATE",
						"instance_name": "update名称",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":     "Stopped",
					"force_stop": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":     "Stopped",
						"force_stop": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_pay", "auto_renew", "create_mode", "direction", "dry_run", "force_stop", "host_name", "image_id", "instance_charge_type", "internet_charge_type", "internet_max_bandwidth_out", "io_optimized", "key_pair_name", "password", "period", "period_unit", "security_enhancement_strategy", "spot_strategy", "support_case"},
			},
		},
	})
}

var AlicloudRdsCustomMap12772 = map[string]string{
	"system_disk_id": CHECKSET,
	"region_id":      CHECKSET,
}

func AlicloudRdsCustomBasicDependence12772(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  ids = [data.alicloud_vswitches.default.vswitches.0.vpc_id]
}

data "alicloud_vswitches" "default" {
  status  = "Available"
  zone_id = "cn-beijing-i"
}

data "alicloud_security_groups" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}`, name)
}
