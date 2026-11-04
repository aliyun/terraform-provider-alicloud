package rds_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass data_disk -- pre-existing TypeList; this whole-product package migration moves files without changing schema, CRUD or list order semantics
// collection-order: bypass security_group_ids -- pre-existing TypeList; this whole-product package migration moves files without changing schema, CRUD or list order semantics
// Test Rds Custom. >>> Resource test cases, automatically generated.
// Case rdscustom_run_ins_extra_param 10836
func TestAccAliCloudRdsCustom_basic10836(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomMap10836)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustom")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomBasicDependence10836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"image_id":                      "",
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
						"vswitch_id":                    acctest.CHECKSET,
						"auto_renew":                    "false",
						"period":                        "1",
						"auto_pay":                      "true",
						"instance_type":                 "mysql.x2.xlarge.6cm",
						"data_disk.#":                   "1",
						"status":                        "Running",
						"security_group_ids.#":          "1",
						"io_optimized":                  "optimized",
						"description":                   "ran_test_ram_code",
						"key_pair_name":                 acctest.CHECKSET,
						"zone_id":                       acctest.CHECKSET,
						"instance_charge_type":          "Prepaid",
						"internet_charge_type":          "PayByTraffic",
						"internet_max_bandwidth_out":    "0",
						"image_id":                      "",
						"security_enhancement_strategy": "Active",
						"period_unit":                   "Month",
						"password":                      "jingyiTEST@123",
						"resource_group_id":             acctest.CHECKSET,
						"host_name":                     acctest.CHECKSET,
						"create_mode":                   acctest.CHECKSET,
						"create_extra_param":            acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
	"region_id": acctest.CHECKSET,
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

data "alicloud_vswitches" "default" {
  vpc_id  = "%s"
  status  = "Available"
  zone_id = var.test_zone_id
}

data "alicloud_security_groups" "default" {
  vpc_id = "%s"
}

resource "alicloud_ecs_key_pair" "KeyPairName" {
  key_pair_name = var.name
}`, name, rdsCustomFixtureVpcID, rdsCustomFixtureVpcID)
}

// Test Rds Custom. <<< Resource test cases, automatically generated.

// Case resourceCase_20260415_01_clone_1_clone_0 12788
func TestAccAliCloudRdsCustom_basic12788(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomMap12788)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustom")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomBasicDependence12788)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"vswitch_id":           acctest.CHECKSET,
						"dry_run":              "false",
						"deployment_set_id":    acctest.CHECKSET,
						"auto_pay":             "true",
						"force":                "false",
						"support_case":         "eni",
						"security_group_ids.#": "1",
						"instance_name":        "namecreate",
						"data_disk.#":          "1",
						"create_mode":          acctest.CHECKSET,
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
	"system_disk_id": acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AlicloudRdsCustomBasicDependence12788(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vswitches" "default" {
  vpc_id  = "%s"
  status  = "Available"
  zone_id = "cn-beijing-i"
}

data "alicloud_security_groups" "default" {
  vpc_id = "%s"
}

resource "alicloud_rds_custom_deployment_set" "deploymentSet" {
  custom_deployment_set_name            = var.name
  description                           = var.name
  group_count                           = 3
  on_unable_to_redeploy_failed_instance = "CancelMembershipAndStart"
  strategy                              = "Availability"
}`, name, rdsCustomFixtureVpcID, rdsCustomFixtureVpcID)
}

// Case resourceCase_20260415_01 12772
func TestAccAliCloudRdsCustom_basic12772(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomMap12772)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustom")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	privateIpAddress := fmt.Sprintf("172.16.3.%d", tfacctest.RandIntRange(10, 240))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomBasicDependence12772)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
						"vswitch_id":           acctest.CHECKSET,
						"dry_run":              "false",
						"auto_pay":             "true",
						"security_group_ids.#": "1",
						"data_disk.#":          "1",
						"create_mode":          acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"system_disk_id": acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AlicloudRdsCustomBasicDependence12772(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vswitches" "default" {
  vpc_id     = "%s"
  cidr_block = "172.16.3.0/24"
  status     = "Available"
  zone_id    = "cn-beijing-i"
}

data "alicloud_security_groups" "default" {
  vpc_id = "%s"
}`, name, rdsCustomFixtureVpcID, rdsCustomFixtureVpcID)
}

// collection-order: bypass security_group_ids -- pre-existing stock attribute outside this change; the returned member order is defined by the VPC security-group echo rather than by the configuration

// Case property_coverage 20260817: user_data, user_data_in_base64, auto_use_coupon, promotion_code,
// enable_jumbo_frame, deletion_protection, reboot, reboot_time, reboot_when_finished,
// data_disk.{device,encrypted,snapshot_id} and system_disk.performance_level.
func TestAccAliCloudRdsCustom_newAttributes(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomNewAttributesMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustom")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	rebootTime := time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02T15:04Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomBasicDependenceNewAttributes)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "tf-testacc-rds-custom-attributes",
					"instance_charge_type": "PostPaid",
					"amount":               "1",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"security_group_ids": []string{
						"${data.alicloud_security_groups.default.ids.0}"},
					"auto_pay":             "true",
					"force":                "true",
					"instance_type":        "mysql.x2.large.9cm",
					"create_mode":          "0",
					"spot_strategy":        "NoSpot",
					"host_name":            "tf-testacc-rds-custom",
					"period_unit":          "Month",
					"password":             "@MY7yxqc9YXQC",
					"auto_use_coupon":      "true",
					"enable_jumbo_frame":   "true",
					"promotion_code":       "tfacc-test-promo",
					"reboot":               "false",
					"reboot_time":          rebootTime,
					"reboot_when_finished": "true",
					"user_data":            "echo Hello Terraform",
					"create_extra_param":   "{}",
					"system_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "40",
							"performance_level": "PL1",
						},
					},
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
							"encrypted":         "true",
							"device":            "",
							"snapshot_id":       "",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "50",
							"encrypted":         "true",
							"device":            "",
							"snapshot_id":       "",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "60",
							"encrypted":         "false",
							"device":            "",
							"snapshot_id":       "${var.snapshot_id}",
						},
					},
					"timeouts": []map[string]interface{}{
						{
							"create": "40m",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_use_coupon":                 "true",
						"enable_jumbo_frame":              "true",
						"promotion_code":                  "tfacc-test-promo",
						"reboot":                          "false",
						"reboot_when_finished":            "true",
						"user_data":                       "echo Hello Terraform",
						"create_extra_param":              "{}",
						"data_disk.#":                     "3",
						"data_disk.0.encrypted":           "true",
						"data_disk.2.size":                "60",
						"system_disk.0.performance_level": "PL1",
					}),
					resource.TestCheckResourceAttr(resourceId, "data_disk.2.snapshot_id", rdsCustomSnapshotFixtureID),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "50",
							"encrypted":         "true",
							"device":            "",
							"snapshot_id":       "",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
							"encrypted":         "true",
							"device":            "",
							"snapshot_id":       "",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "60",
							"encrypted":         "false",
							"device":            "",
							"snapshot_id":       "${var.snapshot_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":      "3",
						"data_disk.0.size": "50",
					}),
					resource.TestCheckResourceAttr(resourceId, "data_disk.2.snapshot_id", rdsCustomSnapshotFixtureID),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_jumbo_frame": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_jumbo_frame": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reboot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reboot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "mysql.x4.large.9cm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "mysql.x4.large.9cm",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_pay", "auto_renew", "auto_use_coupon", "create_extra_param", "create_mode", "direction", "dry_run", "force", "force_stop", "host_name", "image_id", "instance_charge_type", "internet_charge_type", "internet_max_bandwidth_out", "io_optimized", "key_pair_name", "password", "period", "period_unit", "promotion_code", "reboot", "reboot_time", "reboot_when_finished", "security_enhancement_strategy", "spot_strategy", "support_case", "user_data_in_base64"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data":           "IyEvYmluL3NoCmVjaG8gSGVsbG8gVGVycmFmb3Jt",
					"user_data_in_base64": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":           "IyEvYmluL3NoCmVjaG8gSGVsbG8gVGVycmFmb3Jt",
						"user_data_in_base64": "true",
					}),
				),
			},
		},
	})
}

var AlicloudRdsCustomNewAttributesMap = map[string]string{
	"region_id":      acctest.CHECKSET,
	"system_disk_id": acctest.CHECKSET,
}

// rdsCustomFixtureVpcID pins the account VPC used by the generated RDS Custom
// cases so the fixtures do not depend on account-wide data source ordering
// (which changed when unrelated test security groups were created).
const rdsCustomFixtureVpcID = "vpc-2ze059yng967hqhwmx9h9"

// rdsCustomSnapshotFixtureID is an existing accomplished snapshot of a data disk
// in the test account (cn-beijing). ECS rejects CreateSnapshot on a disk that was
// never attached to an instance, so the case consumes this account snapshot for
// data_disk.snapshot_id instead of creating one inside the fixture.
const rdsCustomSnapshotFixtureID = "s-2ze5a4vzfo21qnyxjreg"

func AlicloudRdsCustomBasicDependenceNewAttributes(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "snapshot_id" {
  default = "%s"
}

data "alicloud_vswitches" "default" {
  vpc_id  = "%s"
  status  = "Available"
  zone_id = "cn-beijing-i"
}

data "alicloud_security_groups" "default" {
  vpc_id = "%s"
}`, name, rdsCustomSnapshotFixtureID, rdsCustomFixtureVpcID, rdsCustomFixtureVpcID)
}
