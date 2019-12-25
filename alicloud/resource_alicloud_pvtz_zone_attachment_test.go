package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPvtzZoneAttachment_basic(t *testing.T) {
	var v pvtz.DescribeZoneInfoResponse

	resourceId := "alicloud_pvtz_zone_attachment.default"
	ra := resourceAttrInit(resourceId, pvtzZoneAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneAttachmentConfigDependence)

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
					"zone_id": alicloud_pvtz_zone.default.id,
					"vpc_ids": []string{alicloud_vpc.default.id},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_ids.#": "1",
						"vpcs.#":    "1",
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
					"vpc_ids": []string{alicloud_vpc.default1.id},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_ids.#": "1",
						"vpcs.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_ids": "${alicloud_vpc.defaults.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_ids.#": "2",
						"vpcs.#":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_ids": REMOVEKEY,
					"vpcs": []map[string]interface{}{
						{
							"vpc_id": alicloud_vpc.default.id,
						},
						{
							"vpc_id": alicloud_vpc.default1.id,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpcs.#":    "2",
						"vpc_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_client_ip": "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpcs":           REMOVEKEY,
					"vpc_ids":        []string{alicloud_vpc.default1.id},
					"lang":           "zh",
					"user_client_ip": "192.168.1.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_ids.#":      "1",
						"vpcs.#":         "1",
						"lang":           "zh",
						"user_client_ip": "192.168.1.2",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudPvtzZoneAttachment_multi(t *testing.T) {
	var v pvtz.DescribeZoneInfoResponse

	resourceId := "alicloud_pvtz_zone_attachment.default.4"
	ra := resourceAttrInit(resourceId, pvtzZoneAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneAttachmentConfigDependence)

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
					"zone_id": alicloud_pvtz_zone.default.id,
					"vpc_ids": "${alicloud_vpc.defaults.*.id}",
					"count":   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourcePvtzZoneAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "number" {
		  default = "2"
	}
	resource "alicloud_vpc" "default" {
		name = "tf-testaccPvtzZoneAttachmentConfig"
		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vpc" "default1" {
		name = "tf-testaccPvtzZoneAttachmentConfigUpdate"
		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vpc" "defaults" {
		count = var.number
		cidr_block = "172.16.0.0/12"
		name = "tf-testaccPvtzZoneAttachmentConfigMulti"
	}

	resource "alicloud_pvtz_zone" "default" {
		name = "%s"
	}
	`, name)
}

var pvtzZoneAttachmentBasicMap = map[string]string{
	"zone_id":   CHECKSET,
	"vpc_ids.#": CHECKSET,
}
