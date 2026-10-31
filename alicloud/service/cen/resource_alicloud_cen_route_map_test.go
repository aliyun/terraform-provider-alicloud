package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudCenRouteMap_basic_child_instance_same_region(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default"
	ra := acctest.ResourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenRouteMapChildInstanceSameRegionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_cen_instance_attachment.default00", "alicloud_cen_instance_attachment.default01"},
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      "${var.child_region}",
					"map_result":         "Permit",
					"priority":           "3",
					"transmit_direction": "RegionIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":             acctest.CHECKSET,
						"cen_region_id":      acctest.CHECKSET,
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
			{
				// Cover the remaining must-set + modify-set route map attrs:
				// destination_child_instance_types, source_child_instance_types,
				// route_types, preference, prepend_as_path. They first appear at
				// this step (index > 0) -> counted as modified. Valid in RegionIn
				// because the prior step sets source_region_ids = cen_region_id and
				// source_instance_ids (the PrependAsPath RegionIn constraint:
				// source region must equal the route map's region).
				Config: testAccConfig(map[string]interface{}{
					"destination_child_instance_types": []string{"VPC"},
					"source_child_instance_types":      []string{"VPC"},
					"route_types":                      []string{"System"},
					"preference":                       "50",
					"prepend_as_path":                  []string{"65501"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_child_instance_types.#": "1",
						"source_child_instance_types.#":      "1",
						"route_types.#":                      "1",
						"preference":                         "50",
						"prepend_as_path.#":                  "1",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudCenRouteMap_basic_transit_router_route_table_id(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default"
	ra := acctest.ResourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenRouteMapTransitRouterRouteTableIdConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":                    []string{"alicloud_cen_transit_router.default"},
					"cen_id":                        "${alicloud_cen_instance.default.id}",
					"cen_region_id":                 acctest.DefaultRegionToTest,
					"map_result":                    "Permit",
					"priority":                      "3",
					"transmit_direction":            "RegionIn",
					"transit_router_route_table_id": "${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                        acctest.CHECKSET,
						"cen_region_id":                 acctest.DefaultRegionToTest,
						"map_result":                    "Permit",
						"priority":                      "3",
						"transmit_direction":            "RegionIn",
						"transit_router_route_table_id": acctest.CHECKSET,
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

func TestAccAliCloudCenRouteMap_basic_child_instance_different_region(t *testing.T) {
	resourceId := "alicloud_cen_route_map.default"

	ra := acctest.ResourceAttrInit(resourceId, cenRouteMapBasicMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenRouteMapChildInstanceDifferentRegionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCenRouteMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_cen_instance_attachment.default00", "alicloud_cen_instance_attachment.default01"},
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      "${var.vpc_region_00}",
					"map_result":         "Permit",
					"priority":           "3",
					"transmit_direction": "RegionIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":             acctest.CHECKSET,
						"cen_region_id":      acctest.CHECKSET,
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
			{
				// Modify cen_region_id (vpc_region_00 -> vpc_region_01) to cover
				// the modify-set for this Required, non-ForceNew field. Per the
				// resource schema/Read/Update, cen_region_id is updatable
				// (Update sends request.CenRegionId on HasChange; Read d.Set's
				// it back). The different_region dependence provisions VPCs in
				// both cn-hangzhou and cn-shanghai, so a real second region is
				// used (not a fabricated id).
				Config: testAccConfig(map[string]interface{}{
					"cen_region_id": "${var.vpc_region_01}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_region_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})

}

func TestAccAliCloudCenRouteMap_multi(t *testing.T) {
	var routeMap cbn.RouteMap
	resourceId := "alicloud_cen_route_map.default.4"
	ra := acctest.ResourceAttrInit(resourceId, cenRouteMapBasicMap)
	serviceFunc := func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &routeMap, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccresourceAlicloudCenRouteMap%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenRouteMapConfigMultiDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_cen_transit_router.default"},
					"cen_id":             "${alicloud_cen_instance.default.id}",
					"cen_region_id":      acctest.DefaultRegionToTest,
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
	"cen_id":       acctest.CHECKSET,
	"route_map_id": acctest.CHECKSET,
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

`, name, acctest.DefaultRegionToTest)
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

`, name, acctest.DefaultRegionToTest)
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

provider "alicloudalt" {
	region = "${var.vpc_region_01}"
}

data "alicloud_vpcs" "vpc00" {
	name_regex = "default-NODELETING"
}

data "alicloud_vpcs" "vpc01" {
	provider = alicloudalt
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

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}
`, name)
}

func testAccCheckCenRouteMapDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	cbnService := cen.NewCbnService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_route_map" {
			continue
		}

		_, err := cbnService.DescribeCenRouteMap(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("CEN Route Map still exists: %s", rs.Primary.ID)
	}

	return nil
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
	cenService := cen.NewCenService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cenService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}
