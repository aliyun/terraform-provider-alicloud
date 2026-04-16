package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCENTransitRouterGrantAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_grant_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterGrantAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterGrantAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%scentransitroutergrantattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterGrantAttachmentBasicDependence0)
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
					"order_type":    "PayByCenOwner",
					"instance_id":   "${data.alicloud_vpcs.default.ids.0}",
					"cen_owner_id":  "${data.alicloud_account.default.id}",
					"cen_id":        "${alicloud_cen_instance.default.id}",
					"instance_type": "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type":    "PayByCenOwner",
						"instance_id":   CHECKSET,
						"cen_owner_id":  CHECKSET,
						"cen_id":        CHECKSET,
						"instance_type": "VPC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_type": "PayByResourceOwner",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByResourceOwner",
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

// TestAccAliCloudCENTransitRouterGrantAttachment_crossAccountOrderType covers
// the cross-account variant where the Grant is owned by account B (the VPC
// owner) and authorizes account A's CEN. It exercises order_type updates on
// the Grant itself — the authorizing account's knob that the consuming
// attachment inherits.
func TestAccAliCloudCENTransitRouterGrantAttachment_crossAccountOrderType(t *testing.T) {
	testAccPreCheckCENCrossAccount(t)
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_grant_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterGrantAttachmentMap0)
	providerFactories, factoryProviders := cenCrossAccountProviderFactories()
	// The Grant resource lives in account B (TerraformTest): describe with B.
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		client := cenCrossAccountClientByAK(*factoryProviders, sharedCENCrossAccountCreds.testAK)
		if client == nil {
			return &CbnService{}
		}
		return &CbnService{client}
	}, "DescribeCenTransitRouterGrantAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%scentrgrantx%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCENTransitRouterGrantAttachmentCrossAccountConfig(name, "PayByCenOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type":    "PayByCenOwner",
						"instance_type": "VPC",
					}),
				),
			},
			{
				Config: testAccCENTransitRouterGrantAttachmentCrossAccountConfig(name, "PayByResourceOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByResourceOwner",
					}),
				),
			},
			// ImportState intentionally omitted (see VPN cross-account test).
		},
	})
}

func testAccCENTransitRouterGrantAttachmentCrossAccountConfig(name, orderType string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
%s

data "alicloud_account" "a" {
  provider = alicloud.a
}

# --- Account B: VPC to be granted ---
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

# --- Account A: CEN instance receiving the grant ---
resource "alicloud_cen_instance" "default" {
  provider          = alicloud.a
  cen_instance_name = var.name
  description       = "cross-account grant attachment test"
}

# --- Account B: the grant itself (resource under test) ---
resource "alicloud_cen_transit_router_grant_attachment" "default" {
  instance_type = "VPC"
  instance_id   = alicloud_vpc.default.id
  cen_owner_id  = data.alicloud_account.a.id
  cen_id        = alicloud_cen_instance.default.id
  order_type    = %q
}
`, name, cenCrossAccountProviderBlocks(), orderType)
}

var AlicloudCenTransitRouterGrantAttachmentMap0 = map[string]string{}

func AlicloudCenTransitRouterGrantAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = "test for transit router grant attachment"
}

`, name)
}
