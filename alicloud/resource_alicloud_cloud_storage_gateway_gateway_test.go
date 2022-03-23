package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayGateway_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayMap0)
	var rand = acctest.RandIntRange(10000, 99999)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%scloudstoragegatewaygateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayBasicDependence0)
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
					"storage_bundle_id": "${alicloud_cloud_storage_gateway_storage_bundle.default.id}",
					"vswitch_id":        "${local.vswitch_id}",
					"type":              "Iscsi",
					"payment_type":      "PayAsYouGo",
					"location":          "Cloud",
					"description":       "Description",
					"gateway_class":     "Basic",
					"gateway_name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":          "Iscsi",
						"payment_type":  "PayAsYouGo",
						"location":      "Cloud",
						"description":   "Description",
						"gateway_class": "Basic",
						"gateway_name":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Enhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Enhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Advanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Advanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "File",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "File",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"public_network_bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_network_bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "DescriptionAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "DescriptionAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": "gateway_name_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": "gateway_name_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":              "DescriptionAll",
					"public_network_bandwidth": "20",
					"location":                 "Cloud",
					"gateway_name":             "gateway_nameAll",
					"gateway_class":            "Basic",
					"type":                     "Iscsi",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "DescriptionAll",
						"public_network_bandwidth": "20",
						"gateway_name":             "gateway_nameAll",
						"gateway_class":            "Basic",
						"type":                     "Iscsi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "payment_type", "reason_detail", "release_after_expiration"},
			},
		},
	})
}

func TestAccAlicloudCloudStorageGatewayGateway_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayMap1)
	var rand = acctest.RandIntRange(10000, 99999)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%scloudstoragegatewaygateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayBasicDependence1)
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
					"storage_bundle_id": "${alicloud_cloud_storage_gateway_storage_bundle.default.id}",
					"type":              "Iscsi",
					"payment_type":      "PayAsYouGo",
					"location":          "On_Premise",
					"gateway_name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":         "Iscsi",
						"payment_type": "PayAsYouGo",
						"location":     "On_Premise",
						"gateway_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "File",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "File",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "DescriptionAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "DescriptionAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": "gateway_name_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": "gateway_name_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "DescriptionAll",
					"gateway_name": "gateway_nameAll",
					"type":         "Iscsi",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "DescriptionAll",
						"gateway_name": "gateway_nameAll",
						"type":         "Iscsi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "payment_type", "reason_detail", "release_after_expiration"},
			},
		},
	})
}

var AlicloudCloudStorageGatewayGatewayMap0 = map[string]string{
	"public_network_bandwidth": CHECKSET,
	"reason_type":              NOSET,
	"status":                   CHECKSET,
	"type":                     CHECKSET,
	"location":                 CHECKSET,
	"storage_bundle_id":        CHECKSET,
	"vswitch_id":               CHECKSET,
}

func AlicloudCloudStorageGatewayGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}


data "alicloud_zones" "default"{
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones[1].id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones[1].id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}
`, name)
}

func AlicloudCloudStorageGatewayGatewayBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}
`, name)
}

var AlicloudCloudStorageGatewayGatewayMap1 = map[string]string{
	"reason_type":       NOSET,
	"status":            CHECKSET,
	"type":              CHECKSET,
	"location":          CHECKSET,
	"storage_bundle_id": CHECKSET,
}
