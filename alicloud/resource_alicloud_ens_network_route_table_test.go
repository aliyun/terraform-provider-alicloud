package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
	name := fmt.Sprintf("tf-testacc%sensnetroutetable%d", defaultRegionToTest, rand)
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
					"associate_type":                 "Gateway",
					"description":                    "tf-test-route-table-description",
					"is_default_gateway_route_table": true,
					"network_id":                     "${alicloud_ens_network.default.id}",
					"route_table_name":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"associate_type":                 "Gateway",
						"description":                    "tf-test-route-table-description",
						"is_default_gateway_route_table": "true",
						"route_table_name":               name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"associate_type":                 "Gateway",
					"description":                    "tf-test-route-table-description-updated",
					"is_default_gateway_route_table": true,
					"network_id":                     "${alicloud_ens_network.default.id}",
					"route_table_name":               name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"associate_type":                 "Gateway",
						"description":                    "tf-test-route-table-description-updated",
						"is_default_gateway_route_table": "true",
						"route_table_name":               name + "_update",
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

var AlicloudEnsNetworkRouteTableMap0 = map[string]string{
	"associate_type":                 CHECKSET,
	"description":                    CHECKSET,
	"is_default_gateway_route_table": CHECKSET,
	"network_id":                     CHECKSET,
	"route_table_name":               CHECKSET,
	"route_table_type":               CHECKSET,
	"status":                         CHECKSET,
	"create_time":                    CHECKSET,
}

func AlicloudEnsNetworkRouteTableBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}
`, name)
}
