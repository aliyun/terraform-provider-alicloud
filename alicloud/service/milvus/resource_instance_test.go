// Package milvus_test. This file is generated automatically. Please do not modify it manually, thank you!
package milvus_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/milvus"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					"storage_class":     "alicloud-disk-essd-pl2",
					"size":              "65536",
					"performance_level": "PL2",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMap11774)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependence11774)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"zone_id":               acctest.CHECKSET,
						"vswitch_ids.#":         "1",
						"db_admin_password":     password,
						"components.#":          "5",
						"instance_name":         name,
						"db_version":            acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no", "running_time"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11774 = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMap11770)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependence11770)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
							"cu_num":         "4",
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
							"cu_num":         "16",
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
						"zone_id":             acctest.CHECKSET,
						"vswitch_ids.#":       "1",
						"db_admin_password":   password,
						"components.#":        "5",
						"instance_name":       name,
						"db_version":          acctest.CHECKSET,
						"vpc_id":              acctest.CHECKSET,
						"ha":                  "false",
						"payment_type":        "PayAsYouGo",
						"multi_zone_mode":     "Single",
						"kms_key_id":          "key-test-milvus",
						"encrypted":           "false",
						"resource_group_id":   acctest.CHECKSET,
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
							"cu_num":         "8",
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
							"cu_num":         "32",
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
							"cu_num":         "16",
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
						"resource_group_id": acctest.CHECKSET,
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
							"cu_num":         "32",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no", "running_time"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11770 = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMap11772)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependence11772)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"zone_id":               acctest.CHECKSET,
						"vswitch_ids.#":         "1",
						"db_admin_password":     password,
						"components.#":          "5",
						"instance_name":         name,
						"db_version":            acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no", "running_time"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11772 = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMap11771)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependence11771)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"zone_id":           acctest.CHECKSET,
						"vswitch_ids.#":     "1",
						"db_admin_password": password,
						"components.#":      "5",
						"instance_name":     name,
						"db_version":        acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no", "running_time"},
			},
		},
	})
}

var AliCloudMilvusInstanceMap11771 = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMapAutoPay)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependence11771)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":             acctest.CHECKSET,
						"vswitch_ids.#":       "1",
						"db_admin_password":   password,
						"components.#":        "5",
						"instance_name":       name,
						"db_version":          acctest.CHECKSET,
						"vpc_id":              acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "db_admin_password", "is_multi_az_storage", "load_replicas", "payment_duration", "payment_duration_unit", "promotion_no", "running_time"},
			},
		},
	})
}

var AliCloudMilvusInstanceMapAutoPay = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
}

// Case instance-双可用区vswitch_ids换序覆盖 (TypeList order coverage for vswitch_ids).
// cn-hangzhou supports Milvus availability zones J and I (see the region-and-zone
// document); multi_zone_mode="Two" is the documented dual-availability-zone mode.
// vswitch_ids is a ForceNew TypeList, so reordering members triggers replacement,
// which the order-coverage checker permits. The convergent second apply reorders
// the two distinct vswitch members and requires an empty post-apply plan.
func TestAccAliCloudMilvusInstance_vswitchIdsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_milvus_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMilvusInstanceMapMultiZone)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return milvus.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmilvus%d", rand)
	password := testAccMilvusInstancePassword(rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMilvusInstanceBasicDependenceMultiZone)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.primary_zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.primaryVswitch.id}",
							"zone_id": "${alicloud_vswitch.primaryVswitch.zone_id}",
						},
						{
							"vsw_id":  "${alicloud_vswitch.secondaryVswitch.id}",
							"zone_id": "${alicloud_vswitch.secondaryVswitch.zone_id}",
						},
					},
					"db_admin_password": password,
					"components":        testAccMilvusInstanceStandardComponents(),
					"instance_name":     name,
					"db_version":        "2.4",
					"vpc_id":            "${alicloud_vpc.defaultILXuit.id}",
					"ha":                "false",
					"payment_type":      "PayAsYouGo",
					"multi_zone_mode":   "Two",
					"auto_pay":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           acctest.CHECKSET,
						"vswitch_ids.#":     "2",
						"db_admin_password": password,
						"components.#":      "5",
						"instance_name":     name,
						"db_version":        acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
						"ha":                "false",
						"payment_type":      "PayAsYouGo",
						"multi_zone_mode":   "Two",
						"auto_pay":          "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${var.primary_zone_id}",
					"vswitch_ids": []map[string]interface{}{
						{
							"vsw_id":  "${alicloud_vswitch.secondaryVswitch.id}",
							"zone_id": "${alicloud_vswitch.secondaryVswitch.zone_id}",
						},
						{
							"vsw_id":  "${alicloud_vswitch.primaryVswitch.id}",
							"zone_id": "${alicloud_vswitch.primaryVswitch.zone_id}",
						},
					},
					"db_admin_password": password,
					"components":        testAccMilvusInstanceStandardComponents(),
					"instance_name":     name,
					"db_version":        "2.4",
					"vpc_id":            "${alicloud_vpc.defaultILXuit.id}",
					"ha":                "false",
					"payment_type":      "PayAsYouGo",
					"multi_zone_mode":   "Two",
					"auto_pay":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           acctest.CHECKSET,
						"vswitch_ids.#":     "2",
						"db_admin_password": password,
						"components.#":      "5",
						"instance_name":     name,
						"db_version":        acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
						"ha":                "false",
						"payment_type":      "PayAsYouGo",
						"multi_zone_mode":   "Two",
						"auto_pay":          "true",
					}),
				),
			},
		},
	})
}

var AliCloudMilvusInstanceMapMultiZone = map[string]string{
	"status":               acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"region_id":            acctest.CHECKSET,
	"order_id":             acctest.CHECKSET,
	"security_group_ids.#": acctest.CHECKSET,
	"expire_time":          acctest.CHECKSET,
	"running_time":         acctest.CHECKSET,
}

func AliCloudMilvusInstanceBasicDependenceMultiZone(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "primary_zone_id" {
  default = "cn-hangzhou-j"
}

variable "secondary_zone_id" {
  default = "cn-hangzhou-i"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "primaryVswitch" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.primary_zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-primary"
}

resource "alicloud_vswitch" "secondaryVswitch" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.secondary_zone_id
  cidr_block   = "172.16.2.0/24"
  vswitch_name = "milvus-secondary"
}

`, name)
}

// Test Milvus Instance. <<< Resource test cases, automatically generated.
