package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc HaVip. >>> Resource test cases, automatically generated.
// Case 2535
func TestAccAlicloudVpcHaVip_basic2535(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ha_vip.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcHaVipMap2535)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcHaVip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcHaVipBasicDependence2535)
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
					"vswitch_id":  "${alicloud_vswitch.defaultVswitch.id}",
					"ha_vip_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":  CHECKSET,
						"ha_vip_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_vip_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_vip_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "test",
					"vswitch_id":        "${alicloud_vswitch.defaultVswitch.id}",
					"ha_vip_name":       name + "_update",
					"ip_address":        "192.168.1.101",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test",
						"vswitch_id":        CHECKSET,
						"ha_vip_name":       name + "_update",
						"ip_address":        "192.168.1.101",
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcHaVipMap2535 = map[string]string{}

func AlicloudVpcHaVipBasicDependence2535(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-test-acc-vpc"
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVswitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "192.168.0.0/21"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
  description  = "tf-testacc-vswitch"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 2535  twin
func TestAccAlicloudVpcHaVip_basic2535_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ha_vip.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcHaVipMap2535)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcHaVip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcHaVipBasicDependence2535)
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
					"description":       "test-update",
					"vswitch_id":        "${alicloud_vswitch.defaultVswitch.id}",
					"ha_vip_name":       name,
					"ip_address":        "192.168.1.101",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test-update",
						"vswitch_id":        CHECKSET,
						"ha_vip_name":       name,
						"ip_address":        "192.168.1.101",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Test Vpc HaVip. <<< Resource test cases, automatically generated.
