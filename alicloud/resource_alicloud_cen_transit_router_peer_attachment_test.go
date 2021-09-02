package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterPeerAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_peer_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterPeerAttachmentMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterPeerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterPeerAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":                       "alicloud.cn",
					"cen_id":                         "${alicloud_cen_instance.default.id}",
					"transit_router_id":              "${alicloud_cen_transit_router.default_0.transit_router_id}",
					"peer_transit_router_id":         "${alicloud_cen_transit_router.default_1.transit_router_id}",
					"peer_transit_router_region_id":  "cn-beijing",
					"transit_router_attachment_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_id":                         CHECKSET,
						"peer_transit_router_id":         CHECKSET,
						"transit_router_id":              CHECKSET,
						"peer_transit_router_region_id":  "cn-beijing",
						"transit_router_attachment_name": name,
					}),
				),
			},
			// This step can not work in the multi region.
			//{
			//	ResourceName:            resourceId,
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"dry_run", "route_table_association_enabled", "route_table_propagation_enabled"},
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "tf-testaccdescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "tf-testaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"auto_publish_route_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_bandwidth_package_id": "${alicloud_cen_bandwidth_package.default.id}",
					"bandwidth":                `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"cen_bandwidth_package_id": CHECKSET,
						"bandwidth":                "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled":            `true`,
					"bandwidth":                             `5`,
					"transit_router_attachment_description": "desp",
					"transit_router_attachment_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"auto_publish_route_enabled":            "true",
						"bandwidth":                             "5",
						"transit_router_attachment_description": "desp",
						"transit_router_attachment_name":        name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterPeerAttachmentMap = map[string]string{
	"auto_publish_route_enabled":            CHECKSET,
	"bandwidth":                             CHECKSET,
	"cen_bandwidth_package_id":              "",
	"cen_id":                                CHECKSET,
	"dry_run":                               NOSET,
	"peer_transit_router_id":                CHECKSET,
	"peer_transit_router_region_id":         "cn-beijing",
	"resource_type":                         "TR",
	"route_table_association_enabled":       NOSET,
	"route_table_propagation_enabled":       NOSET,
	"status":                                "Attached",
	"transit_router_attachment_description": "",
	"transit_router_attachment_name":        CHECKSET,
	"transit_router_id":                     CHECKSET,
}

func testAccCheckCenTransitRouterPeerAttachmentExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_cen_transit_router_peer_attachment ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}

			resp, err := cbnService.DescribeCenTransitRouterPeerAttachment(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_transit_router_peer_attachment not found")
	}
}

func testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenTransitRouterPeerAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenTransitRouterPeerAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {

	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_transit_router_peer_attachment" {
			continue
		}
		resp, err := cbnService.DescribeCenTransitRouterPeerAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Transit Router Attachment still exist,  ID %s ", fmt.Sprint(resp["TransitRouterAttachmentId"]))
		}
	}

	return nil
}

func AlicloudCenTransitRouterPeerAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "cn"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider = alicloud.cn
  name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider = alicloud.cn
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider = alicloud.cn
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default_0" {
  provider = alicloud.cn
  cen_id = alicloud_cen_bandwidth_package_attachment.default.instance_id
  transit_router_name = "${var.name}-00"
}

resource "alicloud_cen_transit_router" "default_1" {
  provider = alicloud.bj
  cen_id = alicloud_cen_transit_router.default_0.cen_id
  transit_router_name = "${var.name}-01"
}

`, name)
}
