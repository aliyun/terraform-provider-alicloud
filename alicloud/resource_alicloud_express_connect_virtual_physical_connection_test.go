package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudExpressConnectVirtualPhysicalConnection_basic2033(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_virtual_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectVirtualPhysicalConnectionMap2033)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectVirtualPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sExpressConnectVirtualPhysicalConnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectVirtualPhysicalConnectionBasicDependence2033)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithExpressConnectUidSetting(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_physical_connection_name": "amp_resource_test",
					"description":                      "amp_resource_test",
					"order_mode":                       "PayByPhysicalConnectionOwner",
					"parent_physical_connection_id":    "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"spec":                             "50M",
					"vlan_id":                          "789",
					"vpconn_ali_uid":                   "${var.vpconn_ali_uid}",
					"expect_spec":                      "100M",
					"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_physical_connection_name": "amp_resource_test",
						"description":                      "amp_resource_test",
						"order_mode":                       "PayByPhysicalConnectionOwner",
						"parent_physical_connection_id":    CHECKSET,
						"spec":                             "50M",
						"vlan_id":                          "789",
						"vpconn_ali_uid":                   CHECKSET,
						"expect_spec":                      "100M",
						"resource_group_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expect_spec": "200M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expect_spec": "200M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id": "1124",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": "1124",
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

var AlicloudExpressConnectVirtualPhysicalConnectionMap2033 = map[string]string{}

func AlicloudExpressConnectVirtualPhysicalConnectionBasicDependence2033(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vpconn_ali_uid" {
	default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

`, name, os.Getenv("ALICLOUD_EXPRESS_CONNECT_UID"))
}
