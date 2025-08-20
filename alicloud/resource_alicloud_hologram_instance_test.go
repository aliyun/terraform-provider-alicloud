package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Hologram Instance. >>> Resource test cases, automatically generated.
// Case 3920
func TestAccAliCloudHologramInstance_basic3920(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap3920)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence3920)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name": name,
					"payment_type":  "PayAsYouGo",
					"instance_type": "Warehouse",
					"pricing_cycle": "Hour",
					"cpu":           "32",
					"gateway_count": "2",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"instance_name": name,
						"payment_type":  "PayAsYouGo",
						"instance_type": "Warehouse",
						"endpoints.#":   "2",
						"cpu":           "32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "64",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "64",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_count": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Suspended",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Suspended",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_count": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle": "Hour",
					"cpu":           "64",
					"duration":      "1",
					"auto_pay":      "true",
					"instance_name": name + "_update",
					"gateway_count": "2",
					"payment_type":  "PayAsYouGo",
					"instance_type": "Warehouse",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"pricing_cycle": "Hour",
						"cpu":           "64",
						"duration":      "1",
						"auto_pay":      "true",
						"instance_name": name + "_update",
						"gateway_count": "2",
						"payment_type":  "PayAsYouGo",
						"instance_type": "Warehouse",
						"endpoints.#":   "2",
						"status":        "Running",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap3920 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence3920(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}


`, name)
}

// Case 4132
func TestAccAliCloudHologramInstance_basic4132(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4132)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4132)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "Standard",
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name": name,
					"payment_type":  "PayAsYouGo",
					"pricing_cycle": "Hour",
					"cpu":           "32",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "Standard",
						"cpu":           "32",
						"zone_id":       CHECKSET,
						"instance_name": name,
						"payment_type":  "PayAsYouGo",
						"endpoints.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "64",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "64",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"duration":      "1",
					"instance_type": "Standard",
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle": "Hour",
					"cpu":           "64",
					"instance_name": name + "_update",
					"payment_type":  "PayAsYouGo",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration":      "1",
						"instance_type": "Standard",
						"zone_id":       CHECKSET,
						"pricing_cycle": "Hour",
						"cpu":           "64",
						"instance_name": name + "_update",
						"payment_type":  "PayAsYouGo",
						"endpoints.#":   "2",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap4132 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence4132(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaulVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaulVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = alicloud_vpc.defaulVpc.id
  resource_group_name = var.name

}


`, name)
}

// Case 4785
func TestAccAliCloudHologramInstance_basic4785(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4785)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4785)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSharedSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name": name,
					"payment_type":  "PayAsYouGo",
					"instance_type": "Shared",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"instance_name": name,
						"payment_type":  "PayAsYouGo",
						"instance_type": "Shared",
						"endpoints.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "32",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name + "_update",
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Shared",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":            "Running",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
						"zone_id":           CHECKSET,
						"pricing_cycle":     "Hour",
						"cpu":               "32",
						"storage_size":      "0",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name + "_update",
						"payment_type":      "PayAsYouGo",
						"instance_type":     "Shared",
						"endpoints.#":       "2",
						"status":            "Running",
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap4785 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence4785(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

data "alicloud_resource_manager_resource_groups" "default"{}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-shanghai-e"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}


`, name)
}

// Case 3916
func TestAccAliCloudHologramInstance_basic3916(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap3916)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence3916)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{22})
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_name": name,
					"payment_type":  "Subscription",
					"instance_type": "Standard",
					"pricing_cycle": "Month",
					"duration":      "1",
					"auto_pay":      "true",
					"cpu":           "32",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
							"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"instance_name": name,
						"payment_type":  "Subscription",
						"instance_type": "Standard",
						"endpoints.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "32",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "32",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
							"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap3916 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence3916(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = "cn-hangzhou-j"
}

`, name)
}

// Case 4858
func TestAccAliCloudHologramInstance_basic4858(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4858)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4858)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name": name,
					"payment_type":  "PayAsYouGo",
					"instance_type": "Standard",
					"pricing_cycle": "Hour",
					"cpu":           "32",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"instance_name": name,
						"cpu":           "32",
						"payment_type":  "PayAsYouGo",
						"instance_type": "Standard",
						"endpoints.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "64",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "64",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch2.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch2.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "32",
					"scale_type": "DOWNGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "32",
						"scale_type": "DOWNGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "64",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "64",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "64",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name + "_update",
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Standard",
					"endpoints": []map[string]interface{}{
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
						{
							"type": "Intranet",
						},
					},
					"status":            "Running",
					"initial_databases": "abcd, 123, _underline_db",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
						"zone_id":           CHECKSET,
						"pricing_cycle":     "Hour",
						"cpu":               "64",
						"storage_size":      "0",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name + "_update",
						"payment_type":      "PayAsYouGo",
						"instance_type":     "Standard",
						"endpoints.#":       "2",
						"status":            "Running",
						"initial_databases": "abcd, 123, _underline_db",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap4858 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence4858(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}

resource "alicloud_vpc" "defaultVPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultVSwitch2" {
  vpc_id       = alicloud_vpc.defaultVPC2.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}


`, name)
}

// Case 4783 Deprecated Follower
func SkipTestAccAliCloudHologramInstance_basic4783(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4783)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4783)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":            "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name":      name,
					"payment_type":       "PayAsYouGo",
					"pricing_cycle":      "Hour",
					"cpu":                "32",
					"instance_type":      "Follower",
					"leader_instance_id": "${alicloud_hologram_instance.leaderInstance.id}",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"instance_name": name,
						"payment_type":  "PayAsYouGo",
						"instance_type": "Follower",
						"endpoints.#":   "2",
						"cpu":           "32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":        "64",
					"scale_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":        "64",
						"scale_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type": "Internet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "64",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name + "_update",
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Follower",
					"endpoints": []map[string]interface{}{
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":             "Running",
					"leader_instance_id": "${alicloud_hologram_instance.leaderInstance.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size":  "0",
						"zone_id":            CHECKSET,
						"pricing_cycle":      "Hour",
						"cpu":                "64",
						"storage_size":       "0",
						"duration":           "1",
						"auto_pay":           "true",
						"instance_name":      name + "_update",
						"payment_type":       "PayAsYouGo",
						"instance_type":      "Follower",
						"endpoints.#":        "1",
						"status":             "Running",
						"leader_instance_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

var AlicloudHologramInstanceMap4783 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudHologramInstanceBasicDependence4783(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name

}

resource "alicloud_hologram_instance" "leaderInstance" {
  zone_id       = "cn-hangzhou-j"
  pricing_cycle = "Hour"
  cpu           = "32"
  duration      = "1"
  instance_name = var.name

  endpoints {
    type       = "VPCSingleTunnel"
    vswitch_id = alicloud_vswitch.defaultVSwitch.id
    vpc_id     = alicloud_vswitch.defaultVSwitch.vpc_id
  }
  payment_type  = "PayAsYouGo"
  instance_type = "Standard"
}


`, name)
}

// Case 3920  twin
func TestAccAliCloudHologramInstance_basic3920_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap3920)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence3920)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle": "Hour",
					"cpu":           "32",
					"duration":      "1",
					"auto_pay":      "true",
					"instance_name": name,
					"gateway_count": "4",
					"payment_type":  "PayAsYouGo",
					"instance_type": "Warehouse",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":     "Running",
					"scale_type": "UPGRADE",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"pricing_cycle": "Hour",
						"cpu":           "32",
						"duration":      "1",
						"auto_pay":      "true",
						"instance_name": name,
						"gateway_count": "4",
						"payment_type":  "PayAsYouGo",
						"instance_type": "Warehouse",
						"endpoints.#":   "2",
						"status":        "Running",
						"scale_type":    "UPGRADE",
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Case 4132  twin
func TestAccAliCloudHologramInstance_basic4132_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4132)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4132)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"duration":      "1",
					"instance_type": "Standard",
					"zone_id":       "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle": "Hour",
					"cpu":           "32",
					"instance_name": name,
					"payment_type":  "PayAsYouGo",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
						{
							"type": "Internet",
						},
					},
					"scale_type": "UPGRADE",
					"status":     "Running",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration":      "1",
						"instance_type": "Standard",
						"zone_id":       "cn-hangzhou-j",
						"pricing_cycle": "Hour",
						"cpu":           "32",
						"instance_name": name,
						"payment_type":  "PayAsYouGo",
						"endpoints.#":   "3",
						"scale_type":    "UPGRADE",
						"status":        "Running",
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Case 4785  twin
func TestAccAliCloudHologramInstance_basic4785_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4785)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4785)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSharedSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "32",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name,
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Shared",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":            "Running",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"scale_type":        "UPGRADE",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
						"zone_id":           CHECKSET,
						"pricing_cycle":     "Hour",
						"cpu":               "32",
						"storage_size":      "0",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name,
						"payment_type":      "PayAsYouGo",
						"instance_type":     "Shared",
						"endpoints.#":       "2",
						"status":            "Running",
						"resource_group_id": CHECKSET,
						"scale_type":        "UPGRADE",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Case 3916  twin
func SkipTestAccAliCloudHologramInstance_basic3916_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap3916)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence3916)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "cn-shenzhen-f",
					"pricing_cycle":     "Month",
					"cpu":               "32",
					"storage_size":      "200",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name,
					"payment_type":      "Subscription",
					"instance_type":     "Standard",
					"cold_storage_size": "200",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
							"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
						},
					},
					"scale_type": "UPGRADE",
					"status":     "Suspended",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           "cn-shenzhen-f",
						"pricing_cycle":     "Month",
						"cpu":               "32",
						"storage_size":      "200",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name,
						"payment_type":      "Subscription",
						"instance_type":     "Standard",
						"cold_storage_size": "200",
						"endpoints.#":       "2",
						"scale_type":        "UPGRADE",
						"status":            "Suspended",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Case 4858  twin
func TestAccAliCloudHologramInstance_basic4858_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4858)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4858)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "32",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name,
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Standard",
					"endpoints": []map[string]interface{}{
						{
							"type": "Intranet",
						},
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":            "Running",
					"initial_databases": "abcd, 123, _underline_db",
					"scale_type":        "UPGRADE",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "0",
						"zone_id":           CHECKSET,
						"pricing_cycle":     "Hour",
						"cpu":               "32",
						"storage_size":      "0",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name,
						"payment_type":      "PayAsYouGo",
						"instance_type":     "Standard",
						"endpoints.#":       "2",
						"status":            "Running",
						"initial_databases": "abcd, 123, _underline_db",
						"scale_type":        "UPGRADE",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Case 4783  twin Deprecated Follower
func SkipTestAccAliCloudHologramInstance_basic4783_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramInstanceMap4783)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sholograminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramInstanceBasicDependence4783)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
					"cpu":               "32",
					"storage_size":      "0",
					"duration":          "1",
					"auto_pay":          "true",
					"instance_name":     name,
					"payment_type":      "PayAsYouGo",
					"instance_type":     "Follower",
					"endpoints": []map[string]interface{}{
						{
							"type":       "VPCSingleTunnel",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
							"vpc_id":     "${alicloud_vswitch.defaultVSwitch.vpc_id}",
						},
					},
					"status":             "Running",
					"leader_instance_id": "${alicloud_hologram_instance.leaderInstance.id}",
					"scale_type":         "UPGRADE",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size":  "0",
						"zone_id":            CHECKSET,
						"pricing_cycle":      "Hour",
						"cpu":                "32",
						"storage_size":       "0",
						"duration":           "1",
						"auto_pay":           "true",
						"instance_name":      name,
						"payment_type":       "PayAsYouGo",
						"instance_type":      "Follower",
						"endpoints.#":        "1",
						"status":             "Running",
						"leader_instance_id": CHECKSET,
						"scale_type":         "UPGRADE",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "initial_databases", "pricing_cycle", "scale_type"},
			},
		},
	})
}

// Test Hologram Instance. <<< Resource test cases, automatically generated.
