package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cen TransitRouterEcrAttachment. >>> Resource test cases, automatically generated.
// Case ECR Attachment 5366
func TestAccAliCloudCenTransitRouterEcrAttachment_basic5366(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_ecr_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterEcrAttachmentMap5366)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterEcrAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterecrattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterEcrAttachmentBasicDependence5366)
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
					"ecr_id":                                "${alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id}",
					"cen_id":                                "${alicloud_express_connect_router_tr_association.defaultedPu6c.cen_id}",
					"transit_router_ecr_attachment_name":    name,
					"transit_router_attachment_description": "ecr attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultQa94Y1.transit_router_id}",
					"ecr_owner_id":                          "${data.alicloud_account.current.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_id":                             CHECKSET,
						"transit_router_ecr_attachment_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "ecr attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "ecr attachment",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_ecr_attachment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_ecr_attachment_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_id":                                "${alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id}",
					"cen_id":                                "${alicloud_express_connect_router_tr_association.defaultedPu6c.cen_id}",
					"transit_router_ecr_attachment_name":    name + "_update",
					"transit_router_attachment_description": "ecr attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultQa94Y1.transit_router_id}",
					"ecr_owner_id":                          "${data.alicloud_account.current.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_id":                                CHECKSET,
						"cen_id":                                CHECKSET,
						"transit_router_ecr_attachment_name":    name + "_update",
						"transit_router_attachment_description": "ecr attachment",
						"transit_router_id":                     CHECKSET,
						"ecr_owner_id":                          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudCenTransitRouterEcrAttachmentMap5366 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenTransitRouterEcrAttachmentBasicDependence5366(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
    default = "%s"
}

variable "asn" {
  default = "4200000667"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultO8Hcfx" {
  alibaba_side_asn = var.asn
  ecr_name         = "resource-test-123456"
}

resource "alicloud_cen_instance" "defaultQKBiay" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultQa94Y1" {
  cen_id              = alicloud_cen_instance.defaultQKBiay.id
  transit_router_name = var.name
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_tr_association" "defaultedPu6c" {
  association_region_id   = var.region
  ecr_id                  = alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id
  cen_id                  = alicloud_cen_instance.defaultQKBiay.id
  transit_router_id       = alicloud_cen_transit_router.defaultQa94Y1.transit_router_id
  transit_router_owner_id = data.alicloud_account.current.id
}


`, name, defaultRegionToTest)
}

// Case ECR Attachment 5366  twin
func TestAccAliCloudCenTransitRouterEcrAttachment_basic5366_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_ecr_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterEcrAttachmentMap5366)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterEcrAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterecrattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterEcrAttachmentBasicDependence5366)
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
					"ecr_id":                                "${alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id}",
					"cen_id":                                "${alicloud_express_connect_router_tr_association.defaultedPu6c.cen_id}",
					"transit_router_ecr_attachment_name":    name,
					"transit_router_attachment_description": "ecr attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultQa94Y1.transit_router_id}",
					"ecr_owner_id":                          "${data.alicloud_account.current.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_id":                                CHECKSET,
						"cen_id":                                CHECKSET,
						"transit_router_ecr_attachment_name":    name,
						"transit_router_attachment_description": "ecr attachment",
						"transit_router_id":                     CHECKSET,
						"ecr_owner_id":                          CHECKSET,
						"tags.%":                                "2",
						"tags.Created":                          "TF",
						"tags.For":                              "Test",
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

// Case ECR Attachment 5366  raw
func TestAccAliCloudCenTransitRouterEcrAttachment_basic5366_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_ecr_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterEcrAttachmentMap5366)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterEcrAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterecrattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterEcrAttachmentBasicDependence5366)
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
					"ecr_id":                                "${alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id}",
					"cen_id":                                "${alicloud_express_connect_router_tr_association.defaultedPu6c.cen_id}",
					"transit_router_ecr_attachment_name":    name,
					"transit_router_attachment_description": "ecr attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultQa94Y1.transit_router_id}",
					"ecr_owner_id":                          "${data.alicloud_account.current.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_id":                                CHECKSET,
						"cen_id":                                CHECKSET,
						"transit_router_ecr_attachment_name":    name,
						"transit_router_attachment_description": "ecr attachment",
						"transit_router_id":                     CHECKSET,
						"ecr_owner_id":                          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_ecr_attachment_name":    name + "_update",
					"transit_router_attachment_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_ecr_attachment_name":    name + "_update",
						"transit_router_attachment_description": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

// Test Cen TransitRouterEcrAttachment. <<< Resource test cases, automatically generated.
