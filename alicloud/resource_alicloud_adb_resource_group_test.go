package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudAdsResourceGroup_basic2004(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudAdsResourceGroupMap2004)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.ADBResourceGroupSupportRegions)
	name := fmt.Sprintf("tf_testacc_AdbRG%d", rand)
	name = strings.ToUpper(name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdsResourceGroupBasicDependence2004)
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
					"group_name":    "${var.name}",
					"group_type":    "batch",
					"node_num":      "1",
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":    CHECKSET,
						"group_type":    "batch",
						"node_num":      "1",
						"db_cluster_id": CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"group_name":    "${var.name}",
					"node_num":      "2",
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":    CHECKSET,
						"node_num":      "2",
						"db_cluster_id": CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAdsResourceGroupMap2004 = map[string]string{}

func AlicloudAdsResourceGroupBasicDependence2004(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "creation" {
  default = "ADB"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_adb_zones" "zones_ids" {}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_adb_zones.zones_ids.zones.0.id
  vswitch_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_adb_db_cluster" "default" {
  compute_resource    = "48Core192GBNEW"
  db_cluster_category = "MixedStorage"
  db_node_class       = "E32"
  db_node_count       = 1
  db_node_storage     = 100
  description         = var.name
  elastic_io_resource = 2
  maintain_time       = "04:00Z-05:00Z"
  mode                = "flexible"
  payment_type        = "PayAsYouGo"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_ips        = ["10.168.1.12", "10.168.1.11"]
  vpc_id              = data.alicloud_vpcs.default.ids.0
  vswitch_id          = alicloud_vswitch.vswitch.id
  zone_id             = data.alicloud_adb_zones.zones_ids.zones[0].id
  tags = {
    Created = "TF",
    For     = "example",
  }
}

`, name)
}
