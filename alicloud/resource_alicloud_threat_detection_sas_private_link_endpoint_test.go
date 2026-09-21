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
							"zone_id":     "cn-shanghai-f",
						},
						{
							"v_switch_id": "${alicloud_vswitch.default2.id}",
							"zone_id":     "cn-shanghai-b",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name":         CHECKSET,
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"zones.#":           "2",
						"status":            CHECKSET,
						"update_domain":     CHECKSET,
						"jsrv_domain":       CHECKSET,
						"region_id":         CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"node_name":         "${var.name}_update",
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"zones": []map[string]interface{}{
						{
							"v_switch_id": "${alicloud_vswitch.default.id}",
							"zone_id":     "cn-shanghai-f",
						},
						{
							"v_switch_id": "${alicloud_vswitch.default2.id}",
							"zone_id":     "cn-shanghai-b",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name": CHECKSET,
						"zones.#":   "2",
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
					"region_id":         "cn-shanghai",
					"zones": []map[string]interface{}{
						{
							"v_switch_id": "${alicloud_vswitch.default.id}",
							"zone_id":     "cn-shanghai-f",
						},
						{
							"v_switch_id": "${alicloud_vswitch.default2.id}",
							"zone_id":     "cn-shanghai-b",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name":         CHECKSET,
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"region_id":         "cn-shanghai",
						"zones.#":           "2",
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
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/21"
  zone_id    = "cn-shanghai-f"
}

resource "alicloud_vswitch" "default2" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.8.0/21"
  zone_id    = "cn-shanghai-b"
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

`, name)
}
