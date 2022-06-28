package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudCenRouteMap_basic_child_instance_same_region(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default"
	ra := resourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenRouteMapChildInstanceSameRegionConfigDependence)

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
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      defaultRegionToTest,
					"map_result":         "Permit",
					"priority":           "3",
					"transmit_direction": "RegionIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":             CHECKSET,
						"cen_region_id":      defaultRegionToTest,
						"map_result":         "Permit",
						"priority":           "3",
						"transmit_direction": "RegionIn",
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
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"map_result": "Deny",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"map_result": "Deny",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":   "1",
					"map_result": "Permit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":   "1",
						"map_result": "Permit",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_region_ids":                      []string{"${var.child_region}"},
					"source_instance_ids":                    []string{"${data.alicloud_vpcs.vpc.ids.0}"},
					"source_route_table_ids":                 []string{"${data.alicloud_vpcs.vpc.vpcs.0.route_table_id}"},
					"destination_instance_ids":               []string{"${data.alicloud_vpcs.vpc.ids.0}"},
					"destination_route_table_ids":            []string{"${data.alicloud_vpcs.vpc.vpcs.0.route_table_id}"},
					"destination_cidr_blocks":                []string{"${data.alicloud_vpcs.vpc.vpcs.0.cidr_block}"},
					"match_community_set":                    []string{"65501:1"},
					"match_asns":                             []string{"65501"},
					"operate_community_set":                  []string{"65501:1"},
					"next_priority":                          "3",
					"source_instance_ids_reverse_match":      "false",
					"destination_instance_ids_reverse_match": "false",
					"cidr_match_mode":                        "Include",
					"as_path_match_mode":                     "Include",
					"community_match_mode":                   "Include",
					"community_operate_mode":                 "Additive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_region_ids.#":                    "1",
						"source_instance_ids.#":                  "1",
						"source_route_table_ids.#":               "1",
						"destination_instance_ids.#":             "1",
						"destination_route_table_ids.#":          "1",
						"destination_cidr_blocks.#":              "1",
						"match_community_set.#":                  "1",
						"match_asns.#":                           "1",
						"operate_community_set.#":                "1",
						"next_priority":                          "3",
						"source_instance_ids_reverse_match":      "false",
						"destination_instance_ids_reverse_match": "false",
						"cidr_match_mode":                        "Include",
						"as_path_match_mode":                     "Include",
						"community_match_mode":                   "Include",
						"community_operate_mode":                 "Additive",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudCenRouteMap_basic_transit_router_route_table_id(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default"
	ra := resourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenRouteMapTransitRouterRouteTableIdConfigDependence)

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
					"cen_id":                        "${alicloud_cen_instance.default.id}",
					"cen_region_id":                 defaultRegionToTest,
					"map_result":                    "Permit",
					"priority":                      "3",
					"transmit_direction":            "RegionIn",
					"transit_router_route_table_id": "${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                        CHECKSET,
						"cen_region_id":                 defaultRegionToTest,
						"map_result":                    "Permit",
						"priority":                      "3",
						"transmit_direction":            "RegionIn",
						"transit_router_route_table_id": CHECKSET,
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

func TestAccAlicloudCenRouteMap_basic_child_instance_different_region(t *testing.T) {
	resourceId := "alicloud_cen_route_map.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	ra := resourceAttrInit(resourceId, cenRouteMapBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenRouteMapChildInstanceDifferentRegionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenRouteMapAttachmentDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      defaultRegionToTest,
					"map_result":         "Permit",
					"priority":           "3",
					"transmit_direction": "RegionIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":             CHECKSET,
						"cen_region_id":      defaultRegionToTest,
						"map_result":         "Permit",
						"priority":           "3",
						"transmit_direction": "RegionIn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"map_result": "Deny",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"map_result": "Deny",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":   "1",
					"map_result": "Permit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":   "1",
						"map_result": "Permit",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_region_ids":                      []string{"${var.vpc_region_01}"},
					"source_instance_ids":                    []string{"${data.alicloud_vpcs.vpc01.ids.0}"},
					"source_route_table_ids":                 []string{"${data.alicloud_vpcs.vpc01.vpcs.0.route_table_id}"},
					"destination_instance_ids":               []string{"${data.alicloud_vpcs.vpc00.ids.0}"},
					"destination_route_table_ids":            []string{"${data.alicloud_vpcs.vpc00.vpcs.0.route_table_id}"},
					"destination_cidr_blocks":                []string{"${data.alicloud_vpcs.vpc00.vpcs.0.cidr_block}"},
					"match_community_set":                    []string{"65501:1"},
					"match_asns":                             []string{"65501"},
					"operate_community_set":                  []string{"65501:1"},
					"next_priority":                          "3",
					"source_instance_ids_reverse_match":      "false",
					"destination_instance_ids_reverse_match": "false",
					"cidr_match_mode":                        "Include",
					"as_path_match_mode":                     "Include",
					"community_match_mode":                   "Include",
					"community_operate_mode":                 "Additive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_region_ids.#":                    "1",
						"source_instance_ids.#":                  "1",
						"source_route_table_ids.#":               "1",
						"destination_instance_ids.#":             "1",
						"destination_route_table_ids.#":          "1",
						"destination_cidr_blocks.#":              "1",
						"match_community_set.#":                  "1",
						"match_asns.#":                           "1",
						"operate_community_set.#":                "1",
						"next_priority":                          "3",
						"source_instance_ids_reverse_match":      "false",
						"destination_instance_ids_reverse_match": "false",
						"cidr_match_mode":                        "Include",
						"as_path_match_mode":                     "Include",
						"community_match_mode":                   "Include",
						"community_operate_mode":                 "Additive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_region_ids":                      []string{"${var.vpc_region_00}"},
					"source_instance_ids":                    []string{"${data.alicloud_vpcs.vpc00.ids.0}"},
					"source_route_table_ids":                 []string{"${data.alicloud_vpcs.vpc00.vpcs.0.route_table_id}"},
					"destination_instance_ids":               []string{"${data.alicloud_vpcs.vpc01.ids.0}"},
					"destination_route_table_ids":            []string{"${data.alicloud_vpcs.vpc01.vpcs.0.route_table_id}"},
					"destination_cidr_blocks":                []string{"${data.alicloud_vpcs.vpc01.vpcs.0.cidr_block}"},
					"match_community_set":                    []string{"65500:1"},
					"match_asns":                             []string{"65500"},
					"operate_community_set":                  []string{"65500:1"},
					"next_priority":                          "5",
					"source_instance_ids_reverse_match":      "true",
					"destination_instance_ids_reverse_match": "true",
					"cidr_match_mode":                        "Complete",
					"as_path_match_mode":                     "Complete",
					"community_match_mode":                   "Complete",
					"community_operate_mode":                 "Replace",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_region_ids.#":                    "1",
						"source_instance_ids.#":                  "1",
						"source_route_table_ids.#":               "1",
						"destination_instance_ids.#":             "1",
						"destination_route_table_ids.#":          "1",
						"destination_cidr_blocks.#":              "1",
						"match_community_set.#":                  "1",
						"match_asns.#":                           "1",
						"operate_community_set.#":                "1",
						"next_priority":                          "5",
						"source_instance_ids_reverse_match":      "true",
						"destination_instance_ids_reverse_match": "true",
						"cidr_match_mode":                        "Complete",
						"as_path_match_mode":                     "Complete",
						"community_match_mode":                   "Complete",
						"community_operate_mode":                 "Replace",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudCenRouteMap_multi(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default.4"
	ra := resourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenRouteMapConfigMultiDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      defaultRegionToTest,
					"count":              "5",
					"priority":           "${count.index+1}",
					"transmit_direction": "RegionIn",
					"map_result":         "Permit",
					"description":        "${var.name}-${count.index}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

var cenRouteMapBasicMap = map[string]string{
	"cen_id":       CHECKSET,
	"route_map_id": CHECKSET,
}

func resourceCenRouteMapTransitRouterRouteTableIdConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "child_region" {
    default = "%s"
}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
	transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	transit_router_route_table_name =  var.name
	transit_router_route_table_description = "description"
}

`, name, defaultRegionToTest)
}

func resourceCenRouteMapChildInstanceSameRegionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "child_region" {
    default = "%s"
}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = "${var.name}"
}

data "alicloud_vpcs" "vpc" {
	name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance_attachment" "default00" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${data.alicloud_vpcs.vpc.ids.0}"
	child_instance_type = "VPC"
	child_instance_region_id = "${var.child_region}"
}

resource "alicloud_cen_instance_attachment" "default01" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${data.alicloud_vpcs.vpc.ids.0}"
	child_instance_type = "VPC"
	child_instance_region_id = "${var.child_region}"
}

`, name, defaultRegionToTest)
}

func resourceCenRouteMapChildInstanceDifferentRegionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vpc_region_00" {
    default = "cn-hangzhou"
}

variable "vpc_region_01" {
    default = "cn-shanghai"
}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = "${var.name}"
}

provider "alicloud" {
	alias = "vpc00_region"
	region = "${var.vpc_region_00}"
}

provider "alicloud" {
	alias = "vpc01_region"
	region = "${var.vpc_region_01}"
}

data "alicloud_vpcs" "vpc00" {
	provider = "alicloud.vpc00_region"
	name_regex = "default-NODELETING"
}

data "alicloud_vpcs" "vpc01" {
	provider = "alicloud.vpc01_region"
	name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance_attachment" "default00" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${data.alicloud_vpcs.vpc00.ids.0}"
	child_instance_type = "VPC"
	child_instance_region_id = "${var.vpc_region_00}"
}

resource "alicloud_cen_instance_attachment" "default01" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${data.alicloud_vpcs.vpc01.ids.0}"
	child_instance_type = "VPC"
	child_instance_region_id = "${var.vpc_region_01}"
}

`, name)
}

func resourceCenRouteMapConfigMultiDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cen_instance" "default" {
	name = "${var.name}"
}
`, name)
}

func testAccCheckCenRouteMapAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenRouteMapAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenRouteMapAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cenService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}
