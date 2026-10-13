package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			testAccPreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
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

var AlicloudGpdbDbResourceGroupMap6919 = map[string]string{}

func AlicloudGpdbDbResourceGroupBasicDependence6919(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_gpdb_account" "default1" {
  account_name        = "tftestacc121"
  db_instance_id      = alicloud_gpdb_instance.default.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}

resource "alicloud_gpdb_account" "default2" {
  account_name        = "tftestacc122"
  db_instance_id      = alicloud_gpdb_instance.default.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}

resource "alicloud_gpdb_account" "default3" {
  account_name        = "tftestacc123"
  db_instance_id      = alicloud_gpdb_instance.default.id
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
			testAccPreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_config": "{\\\"CpuRateLimit\\\":10,\\\"MemoryLimit\\\":10,\\\"MemorySharedQuota\\\":80,\\\"MemorySpillRatio\\\":0,\\\"Concurrency\\\":10}",
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
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
