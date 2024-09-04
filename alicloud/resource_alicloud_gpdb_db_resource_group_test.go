package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudGpdbDbResourceGroup_basic6919(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbDbResourceGroupMap6919)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbresourcegroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbDbResourceGroupBasicDependence6919)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.defaultJXWSlW.id}",
					"resource_group_name":   "yb_test_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":10,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
						"db_instance_id":        CHECKSET,
						"resource_group_name":   "yb_test_group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":20,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":20,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.defaultJXWSlW.id}",
					"resource_group_name":   "yb_test_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":10,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
						"db_instance_id":        CHECKSET,
						"resource_group_name":   "yb_test_group",
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

var AlicloudGpdbDbResourceGroupMap6919 = map[string]string{}

func AlicloudGpdbDbResourceGroupBasicDependence6919(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultZc8RD9" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultRv5UXt" {
  vpc_id     = alicloud_vpc.defaultZc8RD9.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultJXWSlW" {
  instance_spec              = "2C8G"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  engine                     = "gpdb"
  // resource_management_mode   = "resourceGroup"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = data.alicloud_zones.default.zones.0.id
  vswitch_id                 = alicloud_vswitch.defaultRv5UXt.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultZc8RD9.id
  db_instance_mode           = "StorageElastic"
  description                = "创建资源组依赖实例_01"
}
`, name)
}

func AlicloudGpdbDbResourceGroupBasicDependence6919Raw(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultZc8RD9" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultRv5UXt" {
  vpc_id     = alicloud_vpc.defaultZc8RD9.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultJXWSlW" {
  instance_spec              = "2C8G"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  engine                     = "gpdb"
  // resource_management_mode   = "resourceGroup"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = data.alicloud_zones.default.zones.0.id
  vswitch_id                 = alicloud_vswitch.defaultRv5UXt.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultZc8RD9.id
  db_instance_mode           = "StorageElastic"
  description                = "创建资源组依赖实例_01"
}

resource "alicloud_gpdb_account" "default1" {
  account_name        = format("%%s%%s", var.name, "1")
  db_instance_id      = alicloud_gpdb_instance.defaultJXWSlW.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}

resource "alicloud_gpdb_account" "default2" {
  account_name        = format("%%s%%s", var.name, "2")
  db_instance_id      = alicloud_gpdb_instance.defaultJXWSlW.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}

resource "alicloud_gpdb_account" "default3" {
  account_name        = format("%%s%%s", var.name, "3")
  db_instance_id      = alicloud_gpdb_instance.defaultJXWSlW.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}

`, name)
}

// Case 资源组测试_依赖资源_GPDB_14 6919  twin
func TestAccAliCloudGpdbDbResourceGroup_basic6919_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbDbResourceGroupMap6919)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbresourcegroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbDbResourceGroupBasicDependence6919)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.defaultJXWSlW.id}",
					"resource_group_name":   "yb_test_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":10,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
						"db_instance_id":        CHECKSET,
						"resource_group_name":   "yb_test_group",
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

// Case 资源组测试_依赖资源_GPDB_14 6919  raw
func TestAccAliCloudGpdbDbResourceGroup_basic6919_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbDbResourceGroupMap6919)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tftestb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbDbResourceGroupBasicDependence6919Raw)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.defaultJXWSlW.id}",
					"resource_group_name":   "yb_test_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":10,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
						"db_instance_id":        CHECKSET,
						"resource_group_name":   "yb_test_group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":20,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_config": "{\"CpuRateLimit\":20,\"MemoryLimit\":10,\"MemorySharedQuota\":80,\"MemorySpillRatio\":0,\"Concurrency\":10}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_list": []string{"${alicloud_gpdb_account.default1.account_name}", "${alicloud_gpdb_account.default2.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_list": []string{"${alicloud_gpdb_account.default2.account_name}", "${alicloud_gpdb_account.default3.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_list": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_list.#": "0",
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
