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
func TestAccAliCloudExpressConnectRouterInterface_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressconnectRouterInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sExpressconnectRouterInterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressconnectRouterInterfaceBasicDependence)
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
					"description":           name,
					"opposite_region_id":    os.Getenv("ALICLOUD_REGION"),
					"router_id":             "${alicloud_vpc.default.router_id}",
					"role":                  "InitiatingSide",
					"router_type":           "VRouter",
					"payment_type":          "PayAsYouGo",
					"router_interface_name": name,
					"spec":                  "Mini.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"opposite_region_id":    os.Getenv("ALICLOUD_REGION"),
						"router_id":             CHECKSET,
						"role":                  "InitiatingSide",
						"router_type":           "VRouter",
						"payment_type":          "PayAsYouGo",
						"router_interface_name": name,
						"spec":                  "Mini.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"opposite_router_type": "VRouter",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"opposite_router_type": "VRouter",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"router_interface_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"router_interface_name": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "Mini.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "Mini.5",
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

func TestAccAliCloudExpressConnectRouterInterface_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressconnectRouterInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sExpressconnectRouterInterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressconnectRouterInterfaceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           name,
					"opposite_region_id":    os.Getenv("ALICLOUD_REGION"),
					"router_id":             "${alicloud_vpc.default.router_id}",
					"role":                  "InitiatingSide",
					"router_type":           "VRouter",
					"payment_type":          "Subscription",
					"router_interface_name": name,
					"spec":                  "Mini.2",
					"period":                "1",
					"pricing_cycle":         "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"opposite_region_id":    os.Getenv("ALICLOUD_REGION"),
						"router_id":             CHECKSET,
						"role":                  "InitiatingSide",
						"router_type":           "VRouter",
						"payment_type":          "Subscription",
						"router_interface_name": name,
						"spec":                  "Mini.2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "pricing_cycle"},
			},
		},
	})
}

var AlicloudExpressconnectRouterInterfaceMap = map[string]string{}

func AlicloudExpressconnectRouterInterfaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

`, name)
}
