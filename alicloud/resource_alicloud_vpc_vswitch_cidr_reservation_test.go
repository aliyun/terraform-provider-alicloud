package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc VswitchCidrReservation. >>> Resource test cases, automatically generated.
// Case 3117
func TestAccAlicloudVpcVswitchCidrReservation_basic3117(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_vswitch_cidr_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVswitchCidrReservationMap3117)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVswitchCidrReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcVswitchCidrReservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVswitchCidrReservationBasicDependence3117)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.VpcVSwitchCidrReservationSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_cidr":         "10.0.10.0/24",
					"vswitch_cidr_reservation_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_cidr":         "10.0.10.0/24",
						"vswitch_cidr_reservation_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "testupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version":                    "IPv4",
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_description":  "test",
					"cidr_reservation_cidr":         "10.0.10.0/24",
					"vswitch_cidr_reservation_name": name + "_update",
					"cidr_reservation_type":         "Prefix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":                    "IPv4",
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_description":  "test",
						"cidr_reservation_cidr":         "10.0.10.0/24",
						"vswitch_cidr_reservation_name": name + "_update",
						"cidr_reservation_type":         "Prefix",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_reservation_mask"},
			},
		},
	})
}

var AlicloudVpcVswitchCidrReservationMap3117 = map[string]string{}

func AlicloudVpcVswitchCidrReservationBasicDependence3117(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "10.0.0.0/20"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}


`, name)
}

// Case 3138
func TestAccAlicloudVpcVswitchCidrReservation_basic3138(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_vswitch_cidr_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVswitchCidrReservationMap3138)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVswitchCidrReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcVswitchCidrReservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVswitchCidrReservationBasicDependence3138)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.VpcVSwitchCidrReservationSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_mask":         "24",
					"vswitch_cidr_reservation_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_mask":         "24",
						"vswitch_cidr_reservation_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "testupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_reservation_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_reservation_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_cidr_reservation_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_cidr_reservation_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version":                    "IPv4",
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_description":  "test",
					"vswitch_cidr_reservation_name": name + "_update",
					"cidr_reservation_mask":         "24",
					"cidr_reservation_type":         "Prefix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":                    "IPv4",
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_description":  "test",
						"vswitch_cidr_reservation_name": name + "_update",
						"cidr_reservation_mask":         "24",
						"cidr_reservation_type":         "Prefix",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_reservation_mask"},
			},
		},
	})
}

var AlicloudVpcVswitchCidrReservationMap3138 = map[string]string{}

func AlicloudVpcVswitchCidrReservationBasicDependence3138(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "10.0.0.0/20"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}


`, name)
}

// Case 3117  twin
func TestAccAlicloudVpcVswitchCidrReservation_basic3117_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_vswitch_cidr_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVswitchCidrReservationMap3117)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVswitchCidrReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcVswitchCidrReservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVswitchCidrReservationBasicDependence3117)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.VpcVSwitchCidrReservationSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version":                    "IPv4",
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_description":  "test",
					"cidr_reservation_cidr":         "10.0.10.0/24",
					"vswitch_cidr_reservation_name": name,
					"cidr_reservation_type":         "Prefix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":                    "IPv4",
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_description":  "test",
						"cidr_reservation_cidr":         "10.0.10.0/24",
						"vswitch_cidr_reservation_name": name,
						"cidr_reservation_type":         "Prefix",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_reservation_mask"},
			},
		},
	})
}

// Case 3138  twin
func TestAccAlicloudVpcVswitchCidrReservation_basic3138_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_vswitch_cidr_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVswitchCidrReservationMap3138)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVswitchCidrReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcVswitchCidrReservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVswitchCidrReservationBasicDependence3138)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.VpcVSwitchCidrReservationSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version":                    "IPv4",
					"vswitch_id":                    "${alicloud_vswitch.defaultVSwitch.id}",
					"cidr_reservation_description":  "test",
					"vswitch_cidr_reservation_name": name,
					"cidr_reservation_mask":         "24",
					"cidr_reservation_type":         "Prefix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":                    "IPv4",
						"vswitch_id":                    CHECKSET,
						"cidr_reservation_description":  "test",
						"vswitch_cidr_reservation_name": name,
						"cidr_reservation_mask":         "24",
						"cidr_reservation_type":         "Prefix",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_reservation_mask"},
			},
		},
	})
}

// Test Vpc VswitchCidrReservation. <<< Resource test cases, automatically generated.
