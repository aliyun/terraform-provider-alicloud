package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc RouteEntry. >>> Resource test cases, automatically generated.
// Case RouteEntry单跳预发 10324
func TestAccAliCloudVpcRouteEntry_basic10324(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap10324)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence10324)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"nexthop_type":           "Ipv4Gateway",
					"route_table_id":         "${alicloud_vpc.defaultVpc.route_table_id}",
					"route_entry_name":       name,
					"nexthop_id":             "${alicloud_vpc_ipv4_gateway.defaultIpv4Gateway.id}",
					"destination_cidr_block": "1.1.1.1/32",
					"depends_on":             []string{"alicloud_vpc_ipv4_gateway.defaultIpv4Gateway"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":           "Ipv4Gateway",
						"route_table_id":         CHECKSET,
						"route_entry_name":       name,
						"nexthop_id":             CHECKSET,
						"destination_cidr_block": "1.1.1.1/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nexthop_type":     "HaVip",
					"route_entry_name": name + "_update",
					"nexthop_id":       "${alicloud_havip.defaultHavip.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":     "HaVip",
						"route_entry_name": name + "_update",
						"nexthop_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
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

var AlicloudVpcRouteEntryMap10324 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence10324(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = "TFRouteEntry1"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVswitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "TFRouteEntry1"
}

resource "alicloud_vpc_ipv4_gateway" "defaultIpv4Gateway" {
  ipv4_gateway_name = "TFRouteEntry"
  vpc_id            = alicloud_vpc.defaultVpc.id
  enabled           = true
}

resource "alicloud_havip" "defaultHavip" {
  vswitch_id  = alicloud_vswitch.defaultVswitch.id
  ha_vip_name = "TFRouteEntry1"
}


`, name)
}

func TestAccAliCloudVpcRouteEntry_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence1)
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
					"route_table_id": "vtb-bp11pzfej88uvbu2cb9gr",
					"description":    "1.1.1.1/32",
					"next_hops": []map[string]interface{}{
						{
							"nexthop_type": "RouterInterface",
							"nexthop_id":   "ri-bp13h9zi93400onfxmnmc",
							"weight":       "100",
						},
						{
							"nexthop_type": "RouterInterface",
							"nexthop_id":   "ri-bp1q6pci7f128etm4rzg9",
							"weight":       "100",
						},
					},
					"route_entry_name":       name,
					"destination_cidr_block": "10.1.1.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":         CHECKSET,
						"description":            "1.1.1.1/32",
						"next_hops.#":            "2",
						"route_entry_name":       name,
						"destination_cidr_block": "10.1.1.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_publish_targets": []map[string]interface{}{
						{
							"target_type":        "ECR",
							"target_instance_id": "ecr-jvp19bb8mrhdnc7fep",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_publish_targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_publish_targets": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_publish_targets.#": "0",
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

func TestAccAliCloudVpcRouteEntry_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence1)
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
					"nexthop_type":           "Ipv4Gateway",
					"route_table_id":         "vtb-bp11pzfej88uvbu2cb9gr",
					"route_entry_name":       name,
					"nexthop_id":             "ipv4gw-bp1wvf0vp9zyr1ziel0nw",
					"destination_cidr_block": "1.1.1.1/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":           "Ipv4Gateway",
						"route_table_id":         CHECKSET,
						"route_entry_name":       name,
						"nexthop_id":             CHECKSET,
						"destination_cidr_block": "1.1.1.1/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_publish_targets": []map[string]interface{}{
						{
							"target_type":        "ECR",
							"target_instance_id": "ecr-jvp19bb8mrhdnc7fep",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_publish_targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_publish_targets": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_publish_targets.#": "0",
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

var AlicloudVpcRouteEntryMap1 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
