package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionSasPrivateLinkEndpoint_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_sas_private_link_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionSasPrivateLinkEndpointMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionSasPrivateLinkEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sSasPrivateLinkEndpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionSasPrivateLinkEndpointBasicDependence)
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
					"node_name":         "${var.name}",
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"zones": []map[string]interface{}{
						{
							"v_switch_id": "${alicloud_vswitch.default.id}",
							"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name":         CHECKSET,
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"zones.#":           "1",
						"status":            CHECKSET,
						"update_domain":     CHECKSET,
						"jsrv_domain":       CHECKSET,
						"region_id":         CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"node_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name": CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudThreatDetectionSasPrivateLinkEndpoint_allFields(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_sas_private_link_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionSasPrivateLinkEndpointMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionSasPrivateLinkEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sSasPrivateLinkEndpointAll%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionSasPrivateLinkEndpointBasicDependence)
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
					"node_name":         "${var.name}",
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"region_id":         "${data.alicloud_regions.default.regions.0.id}",
					"zones": []map[string]interface{}{
						{
							"v_switch_id": "${alicloud_vswitch.default.id}",
							"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name":         CHECKSET,
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"region_id":         CHECKSET,
						"zones.#":           "1",
						"status":            CHECKSET,
						"update_domain":     CHECKSET,
						"jsrv_domain":       CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudThreatDetectionSasPrivateLinkEndpointMap = map[string]string{}

func AlicloudThreatDetectionSasPrivateLinkEndpointBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/21"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

`, name)
}
