package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudADBResourcePool_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_resource_pool.default"
	checkoutSupportedRegions(t, true, connectivity.ADBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudADBResourcePoolMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourcePool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TF-TESTACCADBRESOURCEPOOL%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudADBResourcePoolBasicDependence0)
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
					"pool_name":     name,
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
					"query_type":    "batch",
					"node_num":      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pool_name":     name,
						"db_cluster_id": CHECKSET,
						"query_type":    "batch",
						"node_num":      "2",
					}),
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

func TestAccAlicloudADBResourcePool_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_resource_pool.default"
	checkoutSupportedRegions(t, true, connectivity.ADBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudADBResourcePoolMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourcePool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TF-TESTACCADBRESOURCEPOOL%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudADBResourcePoolBasicDependence0)
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
					"pool_name":     name,
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pool_name":     name,
						"db_cluster_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_type": "batch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_type": "batch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_num": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_type": "interactive",
					"node_num":   "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_type": "interactive",
						"node_num":   "2",
					}),
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

var AlicloudADBResourcePoolMap0 = map[string]string{}

func AlicloudADBResourcePoolBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "ADB"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "32Core128GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = data.alicloud_vswitches.default.ids[0]
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF"
    For     = "acceptance-test-update"
  }
}
`, name)
}
