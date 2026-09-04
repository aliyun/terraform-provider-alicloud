// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Milvus Instance. >>> Resource test cases, automatically generated.
func testAccMilvusInstancePassword(rand int) string {
	return fmt.Sprintf("TfAccMilvus%d@Aa", rand)
}

func testAccMilvusInstanceStandardComponents() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":           "data",
			"cu_num":         "2",
			"replica":        "1",
			"cu_type":        "general",
			"disk_size_type": "Normal",
			"data_disk": []map[string]interface{}{
				{
					"storage_class":     "alicloud-disk-essd-pl1",
					"size":              "100",
					"performance_level": "PL1",
					"enabled":           "true",
				},
			},
		},
		{
			"type":           "index",
			"cu_num":         "4",
			"replica":        "2",
			"cu_type":        "general",
			"disk_size_type": "Normal",
		},
		{
			"type":           "query",
			"cu_num":         "4",
			"replica":        "2",
			"cu_type":        "general",
			"disk_size_type": "Normal",
		},
		{
			"type":           "proxy",
			"cu_num":         "2",
			"replica":        "2",
			"cu_type":        "general",
			"disk_size_type": "Normal",
		},
		{
			"type":           "mix_coordinator",
			"cu_num":         "4",
			"replica":        "2",
			"cu_type":        "general",
			"disk_size_type": "Normal",
		},
	}
}

// Case instance_包年包月-年_张家口 11774
func TestAccAliCloudMilvusInstance_basic11774(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudMilvusInstanceMap11774)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MilvusServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMilvusInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMilvusInstanceBasicDependence11774)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.defaultN80M7S.id}",
							"zone_id": "${alicloud_vswitch.defaultN80M7S.zone_id}",
						},
					},
					"db_admin_password":     password,
					"components":            testAccMilvusInstanceStandardComponents(),
					"instance_name":         name,
					"db_version":            "2.4",
					"vpc_id":                "${alicloud_vpc.defaultILXuit.id}",
					"ha":                    "false",
					"payment_type":          "Subscription",
					"multi_zone_mode":       "Single",
					"payment_duration_unit": "year",
					"payment_duration":      "1",
					"auto_renew":            "true",
					"auto_pay":              "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":               CHECKSET,
						"vswitch_ids.#":         "1",
						"db_admin_password":     password,
						"components.#":          "5",
						"instance_name":         name,
						"db_version":            CHECKSET,
						"vpc_id":                CHECKSET,
						"ha":                    "false",
						"payment_type":          "Subscription",
						"multi_zone_mode":       "Single",
						"payment_duration_unit": "year",
						"payment_duration":      "1",
						"auto_renew":            "true",
						"auto_pay":              "true",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "backup_restore_info", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11774 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"region_id":            CHECKSET,
	"order_id":             CHECKSET,
	"security_group_ids.#": CHECKSET,
	"expire_time":          CHECKSET,
	"running_time":         CHECKSET,
}

func AliCloudMilvusInstanceBasicDependence11774(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-test"
}


`, name)
}

// Case instance-按量更新_张家口 11770
func TestAccAliCloudMilvusInstance_basic11770(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudMilvusInstanceMap11770)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MilvusServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMilvusInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMilvusInstanceBasicDependence11770)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.defaultN80M7S.id}",
							"zone_id": "${alicloud_vswitch.defaultN80M7S.zone_id}",
						},
					},
					"db_admin_password": password,
					"components": []map[string]interface{}{
						{
							"type":           "data",
							"cu_num":         "2",
							"replica":        "1",
							"disk_size_type": "Normal",
						},
						{
							"type":           "index",
							"cu_num":         "4",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "query",
							"cu_num":         "8",
							"replica":        "2",
							"disk_size_type": "Large",
						},
						{
							"type":           "proxy",
							"cu_num":         "2",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "mix_coordinator",
							"cu_num":         "4",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
					},
					"instance_name":       name,
					"db_version":          "2.4",
					"vpc_id":              "${alicloud_vpc.defaultILXuit.id}",
					"ha":                  "false",
					"payment_type":        "PayAsYouGo",
					"multi_zone_mode":     "Single",
					"kms_key_id":          "key-test-milvus",
					"encrypted":           "false",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"auto_backup":         "false",
					"is_multi_az_storage": "true",
					"configuration":       "rootCoord:\\n    maxDatabaseNum: 64 # Maximum number of database",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":             CHECKSET,
						"vswitch_ids.#":       "1",
						"db_admin_password":   password,
						"components.#":        "5",
						"instance_name":       name,
						"db_version":          CHECKSET,
						"vpc_id":              CHECKSET,
						"ha":                  "false",
						"payment_type":        "PayAsYouGo",
						"multi_zone_mode":     "Single",
						"kms_key_id":          "key-test-milvus",
						"encrypted":           "false",
						"resource_group_id":   CHECKSET,
						"auto_backup":         "false",
						"is_multi_az_storage": "true",
						"configuration":       "rootCoord:\n    maxDatabaseNum: 64 # Maximum number of database",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_开启备份",
					"auto_backup":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_开启备份",
						"auto_backup":   "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"components": []map[string]interface{}{
						{
							"type":           "mix_coordinator",
							"cu_num":         "4",
							"replica":        "1",
							"disk_size_type": "Normal",
							"cu_type":        "general",
						},
						{
							"type":           "index",
							"cu_type":        "general",
							"cu_num":         "4",
							"replica":        "1",
							"disk_size_type": "Normal",
						},
						{
							"type":           "data",
							"cu_num":         "2",
							"replica":        "1",
							"cu_type":        "general",
							"disk_size_type": "Normal",
							"data_disk": []map[string]interface{}{
								{
									"storage_class":     "alicloud-disk-essd-pl1",
									"size":              "100",
									"performance_level": "PL1",
									"enabled":           "true",
								},
							},
						},
						{
							"type":           "query",
							"cu_num":         "8",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Large",
						},
						{
							"type":           "proxy",
							"cu_num":         "2",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Normal",
						},
					},
					"instance_name": name + "_降配",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"components.#":  "5",
						"instance_name": name + "_降配",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"components": []map[string]interface{}{
						{
							"type":           "mix_coordinator",
							"cu_type":        "general",
							"cu_num":         "4",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "index",
							"cu_type":        "general",
							"cu_num":         "8",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "data",
							"cu_num":         "4",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Normal",
						},
						{
							"type":           "query",
							"cu_num":         "8",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Large",
						},
						{
							"type":           "proxy",
							"cu_num":         "2",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Normal",
						},
					},
					"instance_name": name + "_升配",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"components.#":  "5",
						"instance_name": name + "_升配",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_更新配置项",
					"configuration": "rootCoord:\\n    maxDatabaseNum: 64 # Maximum number of database\\n    maxPartitionNum: 4096 ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_更新配置项",
						"configuration": "rootCoord:\n    maxDatabaseNum: 64 # Maximum number of database\n    maxPartitionNum: 4096 ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "_更新资源组",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_更新资源组",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"components": []map[string]interface{}{
						{
							"type":           "mix_coordinator",
							"cu_type":        "general",
							"cu_num":         "4",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "index",
							"cu_type":        "general",
							"cu_num":         "8",
							"replica":        "2",
							"disk_size_type": "Normal",
						},
						{
							"type":           "data",
							"cu_num":         "8",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Normal",
						},
						{
							"type":           "query",
							"cu_num":         "8",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Large",
						},
						{
							"type":           "proxy",
							"cu_num":         "2",
							"replica":        "2",
							"cu_type":        "general",
							"disk_size_type": "Normal",
						},
					},
					"instance_name": name + "_开启高可用",
					"ha":            "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"components.#":  "5",
						"instance_name": name + "_开启高可用",
						"ha":            "true",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "backup_restore_info", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11770 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"region_id":            CHECKSET,
	"order_id":             CHECKSET,
	"security_group_ids.#": CHECKSET,
	"expire_time":          CHECKSET,
	"running_time":         CHECKSET,
}

func AliCloudMilvusInstanceBasicDependence11770(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-test"
}


`, name)
}

// Case instance包年包月-月_张家口 11772
func TestAccAliCloudMilvusInstance_basic11772(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudMilvusInstanceMap11772)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MilvusServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMilvusInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMilvusInstanceBasicDependence11772)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.defaultN80M7S.id}",
							"zone_id": "${alicloud_vswitch.defaultN80M7S.zone_id}",
						},
					},
					"db_admin_password":     password,
					"components":            testAccMilvusInstanceStandardComponents(),
					"instance_name":         name,
					"db_version":            "2.4",
					"vpc_id":                "${alicloud_vpc.defaultILXuit.id}",
					"ha":                    "false",
					"payment_type":          "Subscription",
					"multi_zone_mode":       "Single",
					"payment_duration_unit": "month",
					"payment_duration":      "1",
					"auto_renew":            "true",
					"is_multi_az_storage":   "true",
					"auto_pay":              "true",
					"load_replicas":         "2",
					"promotion_no":          "youhuiquan_promotion_option_id_for_blank",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":               CHECKSET,
						"vswitch_ids.#":         "1",
						"db_admin_password":     password,
						"components.#":          "5",
						"instance_name":         name,
						"db_version":            CHECKSET,
						"vpc_id":                CHECKSET,
						"ha":                    "false",
						"payment_type":          "Subscription",
						"multi_zone_mode":       "Single",
						"payment_duration_unit": "month",
						"payment_duration":      "1",
						"auto_renew":            "true",
						"is_multi_az_storage":   "true",
						"auto_pay":              "true",
						"load_replicas":         "2",
						"promotion_no":          "youhuiquan_promotion_option_id_for_blank",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "backup_restore_info", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11772 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"region_id":            CHECKSET,
	"order_id":             CHECKSET,
	"security_group_ids.#": CHECKSET,
	"expire_time":          CHECKSET,
	"running_time":         CHECKSET,
}

func AliCloudMilvusInstanceBasicDependence11772(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-test"
}


`, name)
}

// Case instance-按量更新_tag_张家口 11771
func TestAccAliCloudMilvusInstance_basic11771(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudMilvusInstanceMap11771)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MilvusServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMilvusInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMilvusInstanceBasicDependence11771)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.defaultN80M7S.id}",
							"zone_id": "${alicloud_vswitch.defaultN80M7S.zone_id}",
						},
					},
					"db_admin_password": password,
					"components":        testAccMilvusInstanceStandardComponents(),
					"instance_name":     name,
					"db_version":        "2.4",
					"vpc_id":            "${alicloud_vpc.defaultILXuit.id}",
					"ha":                "false",
					"payment_type":      "PayAsYouGo",
					"multi_zone_mode":   "Single",
					"kms_key_id":        "k-test",
					"encrypted":         "false",
					"load_replicas":     "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_ids.#":     "1",
						"db_admin_password": password,
						"components.#":      "5",
						"instance_name":     name,
						"db_version":        CHECKSET,
						"vpc_id":            CHECKSET,
						"ha":                "false",
						"payment_type":      "PayAsYouGo",
						"multi_zone_mode":   "Single",
						"kms_key_id":        "k-test",
						"encrypted":         "false",
						"load_replicas":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "backup_restore_info", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11771 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"region_id":            CHECKSET,
	"order_id":             CHECKSET,
	"security_group_ids.#": CHECKSET,
	"expire_time":          CHECKSET,
	"running_time":         CHECKSET,
}

func AliCloudMilvusInstanceBasicDependence11771(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-test"
}


`, name)
}

// Case instance-按量-auto_pay
func TestAccAliCloudMilvusInstance_autoPay(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudMilvusInstanceMapAutoPay)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MilvusServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMilvusInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMilvusInstanceBasicDependence11771)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.defaultN80M7S.id}",
							"zone_id": "${alicloud_vswitch.defaultN80M7S.zone_id}",
						},
					},
					"db_admin_password":   password,
					"components":          testAccMilvusInstanceStandardComponents(),
					"instance_name":       name,
					"db_version":          "2.4",
					"vpc_id":              "${alicloud_vpc.defaultILXuit.id}",
					"ha":                  "false",
					"payment_type":        "PayAsYouGo",
					"multi_zone_mode":     "Single",
					"kms_key_id":          "k-test",
					"encrypted":           "false",
					"auto_pay":            "false",
					"is_multi_az_storage": "false",
					"load_replicas":       "2",
					"promotion_no":        "youhuiquan_promotion_option_id_for_blank",
					"backup_restore_info": []map[string]interface{}{
						{
							"backup_id":         "bt-69b8fdaff88db73f",
							"source_cluster_id": "c-2d95c862a7142aec",
							"backup_name":       "auto_backup_2026_08_25_03_00_00_902610000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":             CHECKSET,
						"vswitch_ids.#":       "1",
						"db_admin_password":   password,
						"components.#":        "5",
						"instance_name":       name,
						"db_version":          CHECKSET,
						"vpc_id":              CHECKSET,
						"ha":                  "false",
						"payment_type":        "PayAsYouGo",
						"multi_zone_mode":     "Single",
						"kms_key_id":          "k-test",
						"encrypted":           "false",
						"auto_pay":            "false",
						"is_multi_az_storage": "false",
						"load_replicas":       "2",
						"promotion_no":        "youhuiquan_promotion_option_id_for_blank",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "backup_restore_info", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no"},
			},
		},
	})
}

var AliCloudMilvusInstanceMapAutoPay = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"region_id":            CHECKSET,
	"order_id":             CHECKSET,
	"security_group_ids.#": CHECKSET,
	"expire_time":          CHECKSET,
	"running_time":         CHECKSET,
}

// Test Milvus Instance. <<< Resource test cases, automatically generated.
