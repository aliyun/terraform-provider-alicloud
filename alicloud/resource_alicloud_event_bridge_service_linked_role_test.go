package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEventBridgeServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeEventServiceLinkedRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudEventBridgeServiceLinkedRoleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": "AliyunServiceRoleForEventBridgeSourceRocketMQ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": "AliyunServiceRoleForEventBridgeSourceRocketMQ",
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

func TestAccAlicloudEventBridgeServiceLinkedRole_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeEventServiceLinkedRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudEventBridgeServiceLinkedRoleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": "AliyunServiceRoleForEventBridgeSendToMNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": "AliyunServiceRoleForEventBridgeSendToMNS",
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

var AlicloudEventBridgeEventServiceLinkedRoleMap0 = map[string]string{
	"product_name": "AliyunServiceRoleForEventBridgeSourceRocketMQ",
}

func testAccCheckAlicloudEventBridgeServiceLinkedRoleDependence(name string) string {
	return ""
}
