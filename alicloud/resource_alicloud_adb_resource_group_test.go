package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAdbResourceGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-AdbResourceGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence0)
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
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
					"group_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"group_name":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_type": "batch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_type": "batch",
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
					"users": []string{"${alicloud_adb_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"users.#": "1",
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

func TestAccAliCloudAdbResourceGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-AdbResourceGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence0)
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
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
					"group_name":    name,
					"group_type":    "batch",
					"node_num":      "1",
					"users":         []string{"${alicloud_adb_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"group_name":    CHECKSET,
						"group_type":    "batch",
						"node_num":      "1",
						"users.#":       "1",
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

var AliCloudAdbResourceGroupMap0 = map[string]string{
	"group_type":  CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
}

func AliCloudAdbResourceGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_adb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.0
	}

	resource "alicloud_adb_db_cluster" "default" {
  		compute_resource    = "32Core128GBNEW"
  		db_cluster_category = "MixedStorage"
  		description         = var.name
  		elastic_io_resource = 1
  		mode                = "flexible"
  		payment_type        = "PayAsYouGo"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		zone_id             = data.alicloud_adb_zones.default.zones.0.id
	}

	resource "alicloud_adb_account" "default" {
  		db_cluster_id    = alicloud_adb_db_cluster.default.id
  		account_name     = "tf_account_name"
  		account_password = "YourPassword123!"
	}
`, name)
}
