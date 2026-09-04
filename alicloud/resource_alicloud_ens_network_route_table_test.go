package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Ens NetworkRouteTable. >>> Resource test cases, automatically generated.
func TestAccAliCloudEnsNetworkRouteTable_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_network_route_table.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNetworkRouteTableMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNetworkRouteTable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnetworkroutetable%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNetworkRouteTableBasicDependence0)
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
					"network_id":    "${alicloud_ens_network.defaultObbrL7.id}",
					"associate_type": "Gateway",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_id":    CHECKSET,
						"associate_type": "Gateway",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_id":    "${alicloud_ens_network.defaultObbrL7.id}",
					"associate_type": "Gateway",
					"route_table_name": name,
					"description": "test description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_id":    CHECKSET,
						"associate_type": "Gateway",
						"route_table_name": name,
						"description": "test description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_id":    "${alicloud_ens_network.defaultObbrL7.id}",
					"associate_type": "Gateway",
					"route_table_name": name,
					"description": "test description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_id":    CHECKSET,
						"associate_type": "Gateway",
						"route_table_name": name,
						"description": "test description",
						"is_default_gateway_route_table": "true",
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

var AlicloudEnsNetworkRouteTableMap0 = map[string]string{
	"create_time": CHECKSET,
	"status":      CHECKSET,
	"route_table_type": CHECKSET,
}

func AlicloudEnsNetworkRouteTableBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-31"
}

resource "alicloud_ens_network" "defaultObbrL7" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}


`, name)
}

// Test Ens NetworkRouteTable. <<< Resource test cases, automatically generated.
