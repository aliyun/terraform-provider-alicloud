package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ExpressConnectRouter ExpressConnectRouterVbrChildInstance. >>> Resource test cases, automatically generated.
// Case 初始版本测试用例 6368
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstance_basic6368(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vbr_child_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceMap6368)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervbrchildinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceBasicDependence6368)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"child_instance_type":      "VBR",
					"child_instance_id":        "${alicloud_express_connect_virtual_border_router.defaultydbbk3.id}",
					"child_instance_region_id": "${data.alicloud_regions.default.regions.0.id}",
					"ecr_id":                   "${alicloud_express_connect_router_express_connect_router.defaultAAlhUy.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"child_instance_type":      "VBR",
						"child_instance_id":        CHECKSET,
						"child_instance_region_id": CHECKSET,
						"ecr_id":                   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated Description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated Description",
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

func TestAccAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstance_basic6368_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vbr_child_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceMap6368)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervbrchildinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceBasicDependence6368)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"child_instance_type":      "VBR",
					"child_instance_id":        "${alicloud_express_connect_virtual_border_router.defaultydbbk3.id}",
					"child_instance_owner_id":  "${data.alicloud_account.default.id}",
					"child_instance_region_id": "${data.alicloud_regions.default.regions.0.id}",
					"ecr_id":                   "${alicloud_express_connect_router_express_connect_router.defaultAAlhUy.id}",
					"description":              "Description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"child_instance_type":      "VBR",
						"child_instance_id":        CHECKSET,
						"child_instance_owner_id":  CHECKSET,
						"child_instance_region_id": CHECKSET,
						"ecr_id":                   CHECKSET,
						"description":              "Description",
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

var AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceMap6368 = map[string]string{
	"status":                  CHECKSET,
	"create_time":             CHECKSET,
	"child_instance_owner_id": CHECKSET,
}

func AliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceBasicDependence6368(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "defaultydbbk3" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1000"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultAAlhUy" {
  alibaba_side_asn = "65532"
}


`, name)
}

// Test ExpressConnectRouter ExpressConnectRouterVbrChildInstance. <<< Resource test cases, automatically generated.
