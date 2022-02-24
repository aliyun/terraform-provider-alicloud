package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECPInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_instance.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPInstanceBasicDependence0)
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
					"instance_name":     name,
					"description":       name,
					"key_pair_name":     "${alicloud_ecp_key_pair.default.key_pair_name}",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"payment_type":      "PayAsYouGo",
					"vnc_password":      "Cp1234",
					"force":             "true",
					"depends_on":        []string{"alicloud_ecp_key_pair.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"description":       name,
						"key_pair_name":     CHECKSET,
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"payment_type":      CHECKSET,
						"instance_type":     CHECKSET,
						"resolution":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "auto_renew", "force", "vnc_password", "auto_pay", "payment_type", "eip_bandwidth", "period_unit"},
			},
		},
	})
}

var AlicloudECPInstanceMap0 = map[string]string{}

func AlicloudECPInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`, name)
}
func TestAccAlicloudECPInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_instance.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPInstanceMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPInstanceBasicDependence1)
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
					"instance_name":     name,
					"description":       name,
					"force":             "true",
					"payment_type":      "PayAsYouGo",
					"key_pair_name":     "${alicloud_ecp_key_pair.default.0.key_pair_name}",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"status":            "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"description":       name,
						"payment_type":      "PayAsYouGo",
						"key_pair_name":     CHECKSET,
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"status":            "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${alicloud_ecp_key_pair.default.1.key_pair_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vnc_password": "Cp1232",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vnc_password": "Cp1232",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "update",
					"description":       name + "update",
					"key_pair_name":     "${alicloud_ecp_key_pair.default.0.key_pair_name}",
					"vnc_password":      "Cp1234",
					"payment_type":      "PayAsYouGo",
					"force":             "true",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"status":            "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "update",
						"description":       name + "update",
						"key_pair_name":     CHECKSET,
						"vnc_password":      "Cp1234",
						"security_group_id": CHECKSET,
						"payment_type":      "PayAsYouGo",
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"status":            "Running",
					}),
				),
			},
		},
	})
}

var AlicloudECPInstanceMap1 = map[string]string{
	"period_unit":   NOSET,
	"period":        NOSET,
	"status":        CHECKSET,
	"auto_renew":    NOSET,
	"force":         CHECKSET,
	"eip_bandwidth": NOSET,
	"amount":        NOSET,
	"auto_pay":      NOSET,
}

func AlicloudECPInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  count           = 2
  key_pair_name   = join("", [var.name, count.index])
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`, name)
}
