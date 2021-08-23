package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostUser_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_user.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostUserBasicDependence0)
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
					"user_name":   "tf-testAccBastionHostUser-12345",
					"source":      "Local",
					"instance_id": "${alicloud_bastionhost_instance.default.id}",
					"password":    "tf-testAcc-oAupFqRaH24MdOSrsIKsu3qw",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   "tf-testAccBastionHostUser-12345",
						"source":      "Local",
						"instance_id": "${alicloud_bastionhost_instance.default.id}",
						"password":    "tf-testAcc-oAupFqRaH24MdOSrsIKsu3qw",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAcc-mrPFCPi3MuIloLzTvVzQbUbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAcc-mrPFCPi3MuIloLzTvVzQbUbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "tf-testAcc-5V8AgQKKw389irWIePb47aOq@163.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "tf-testAcc-5V8AgQKKw389irWIePb47aOq@163.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "tf-testAcc-RZEdvPXF9A3w3ArhFwuAfUoY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "tf-testAcc-RZEdvPXF9A3w3ArhFwuAfUoY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "HK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "HK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "TW",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "TW",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "RU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "RU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SG",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "ID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "ID",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "DE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "DE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "US",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "US",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "JP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "JP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "IN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "IN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "KR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "KR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "PH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "PH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "702345672",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "702345672",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Frozen",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Frozen",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-Li6bvnYmD9ryuLUt2Wsdn4gy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-Li6bvnYmD9ryuLUt2Wsdn4gy",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":             "tf-testAcc-X23IfHiv8DnMoYChjEnb6X2h",
					"email":               "tf-testAcc-Rw5hfV8W1mkMO44chYBC07sq@163.com",
					"display_name":        "tf-testAcc-yAwB1akRJGW9RVMaTEdOHOHS",
					"mobile_country_code": "CN",
					"mobile":              "13312345678",
					"password":            "tf-testAcc-lBdrpSbUJ4Ddw9oSGCeI2u2p",
					"status":              "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":             "tf-testAcc-X23IfHiv8DnMoYChjEnb6X2h",
						"email":               "tf-testAcc-Rw5hfV8W1mkMO44chYBC07sq@163.com",
						"display_name":        "tf-testAcc-yAwB1akRJGW9RVMaTEdOHOHS",
						"mobile_country_code": "CN",
						"mobile":              "13312345678",
						"password":            "tf-testAcc-lBdrpSbUJ4Ddw9oSGCeI2u2p",
						"status":              "Normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudBastionhostUserMap0 = map[string]string{
	"display_name":        CHECKSET,
	"status":              CHECKSET,
	"instance_id":         CHECKSET,
	"mobile_country_code": "",
	"user_id":             CHECKSET,
}

func AlicloudBastionhostUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
 zone_id = local.zone_id
 vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
 count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id       = data.alicloud_vpcs.default.ids.0
 zone_id      = data.alicloud_zones.default.ids.0
 cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 name   = var.name
}
locals {
 vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
 zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
 description        = var.name
 license_code       = "bhah_ent_50_asset"
 period             = "1"
 vswitch_id         = local.vswitch_id
 security_group_ids = [alicloud_security_group.default.id]
}
`, name)
}
func TestAccAlicloudBastionhostUser_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_user.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostUserMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostUserBasicDependence1)
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
					"user_name":      "tf-testAccBastionhostUserRam-123456",
					"source":         "Ram",
					"instance_id":    "${alicloud_bastionhost_instance.default.id}",
					"source_user_id": "247823888127488180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":      "tf-testAccBastionhostUserRam-123456",
						"source":         "Ram",
						"instance_id":    "${alicloud_bastionhost_instance.default.id}",
						"source_user_id": "247823888127488180",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAccBastionhostUserRam-123456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAccBastionhostUserRam-123456",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "tf-testAcc-LmwD6dS7fyO93I@163.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "tf-testAcc-LmwD6dS7fyO93I@163.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "tf-testAccBastionhostUserRam-456789",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "tf-testAccBastionhostUserRam-456789",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "HK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "HK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "TW",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "TW",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "RU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "RU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SG",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "ID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "ID",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "DE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "DE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "US",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "US",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "JP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "JP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "IN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "IN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "KR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "KR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "PH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "PH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "702345672",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "702345672",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Frozen",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Frozen",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":             "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"email":               "tf-testAcc-75MYawy06OnL4zTD4xdi6n4T@163.com",
					"display_name":        "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"mobile_country_code": "CN",
					"mobile":              "13312345678",
					"password":            "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"status":              "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":             "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"email":               "tf-testAcc-75MYawy06OnL4zTD4xdi6n4T@163.com",
						"display_name":        "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"mobile_country_code": "CN",
						"mobile":              "13312345678",
						"password":            "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"status":              "Normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudBastionhostUserMap1 = map[string]string{
	"user_id":             CHECKSET,
	"display_name":        CHECKSET,
	"status":              CHECKSET,
	"instance_id":         CHECKSET,
	"mobile_country_code": "",
}

func AlicloudBastionhostUserBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
 zone_id = local.zone_id
 vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
 count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id       = data.alicloud_vpcs.default.ids.0
 zone_id      = data.alicloud_zones.default.ids.0
 cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 name   = var.name
}
locals {
 vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
 zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
 description        = var.name
 license_code       = "bhah_ent_50_asset"
 period             = "1"
 vswitch_id         = local.vswitch_id
 security_group_ids = [alicloud_security_group.default.id]
}
`, name)
}
