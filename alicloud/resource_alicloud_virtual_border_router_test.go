package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVirtualBorderRouter_basic(t *testing.T) {
	var v vpc.VirtualBorderRouterType
	resourceId := "alicloud_virtual_border_router.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccVbrBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVirtualBorderRouterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithPhysicalConnectionSetting(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id": os.Getenv("ALICLOUD_PHYSICAL_CONNECTION_ID"),
					"vlan_id":                "2500",
					"local_gateway_ip":       "10.0.0.1",
					"peer_gateway_ip":        "10.0.0.2",
					"peering_subnet_mask":    "255.255.255.0",
					"name":                   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"vlan_interface_id": CHECKSET,
						"route_table_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id": "2501",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": "2501",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_gateway_ip": "10.0.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_gateway_ip": "10.0.0.4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_subnet_mask": "255.255.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_subnet_mask": "255.255.0.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_description", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_description", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id":             "2502",
					"local_gateway_ip":    "10.0.0.4",
					"peer_gateway_ip":     "10.0.0.5",
					"peering_subnet_mask": "255.255.252.0",
					"name":                name,
					"description":         fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_description", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id":             "2502",
						"local_gateway_ip":    "10.0.0.4",
						"peer_gateway_ip":     "10.0.0.5",
						"peering_subnet_mask": "255.255.252.0",
						"name":                name,
						"description":         fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_description", rand),
						"vlan_interface_id":   CHECKSET,
						"route_table_id":      CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVirtualBorderRouter_multi(t *testing.T) {
	var v vpc.VirtualBorderRouterType
	resourceId := "alicloud_virtual_border_router.default.1"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccVbrMulti%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVirtualBorderRouterConfigDependenceForMulti)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                  "2",
					"physical_connection_id": os.Getenv("ALICLOUD_PHYSICAL_CONNECTION_ID"),
					"vlan_id":                "${element(var.vlan_id_list,count.index)}",
					"local_gateway_ip":       "${element(var.local_gateway_ip_list,count.index)}",
					"peer_gateway_ip":        "${element(var.peer_gateway_ip_list,count.index)}",
					"peering_subnet_mask":    "255.255.255.0",
					"name":                   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id":             "2504",
						"local_gateway_ip":    "10.0.0.3",
						"peer_gateway_ip":     "10.0.0.4",
						"peering_subnet_mask": "255.255.255.0",
						"name":                name,
						"vlan_interface_id":   CHECKSET,
						"route_table_id":      CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceVirtualBorderRouterConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	`, name)
}

func resourceVirtualBorderRouterConfigDependenceForMulti(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	variable "vlan_id_list" {
		type = "list"
		default = [ "2503", "2504" ]
	}
	variable "local_gateway_ip_list" {
		type = "list"
		default = [ "10.0.0.1", "10.0.0.3" ]
	}
	variable "peer_gateway_ip_list" {
		type = "list"
		default = [ "10.0.0.2", "10.0.0.4" ]
	}
	`, name)
}
