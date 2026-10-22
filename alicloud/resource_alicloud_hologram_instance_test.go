package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithTime(t, []int{22})
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"leader_instance_id": "${alicloud_hologram_instance.leaderInstance.id}",
					"zone_id":            "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_name":      name,
					"payment_type":       "PayAsYouGo",
					"instance_type":      "Standard",
					"pricing_cycle":      "Hour",
					"cpu":                "32",
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"endpoints.#":   "2",
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
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSharedSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "0",
					"zone_id":           "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"pricing_cycle":     "Hour",
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
						"storage_size":      "0",
						"duration":          "1",
						"auto_pay":          "true",
						"instance_name":     name,
						"payment_type":      "PayAsYouGo",
						"instance_type":     "Shared",
						"endpoints.#":       "2",
						"status":            "Running",
						"resource_group_id": CHECKSET,
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"enable_ssl":        "true",
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
						"enable_ssl":        "true",
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

type vpcHologramAPI struct {
	region         string
	list           func() (map[string]interface{}, error)
	get            func(string) (map[string]interface{}, error)
	remove         func(string) (map[string]interface{}, error)
	now            func() time.Time
	sleep          func(time.Duration)
	releaseTimeout time.Duration
	maxPolls       int
}

func nativeVpcHologramAPI(client *connectivity.AliyunClient) vpcHologramAPI {
	return vpcHologramAPI{region: client.RegionId,
		// Official ListInstances has a single optional opaque body and a complete
		// InstanceList response: no page/size/next-token fields exist in its schema.
		list: func() (map[string]interface{}, error) {
			return client.RoaPost("Hologram", "2022-06-01", "/api/v1/instances", map[string]*string{}, nil, map[string]interface{}{}, true)
		},
		get: func(id string) (map[string]interface{}, error) {
			return client.RoaGet("Hologram", "2022-06-01", "/api/v1/instances/"+id, map[string]*string{}, nil, map[string]interface{}{"instanceId": id})
		},
		remove: func(id string) (map[string]interface{}, error) {
			return client.RoaPost("Hologram", "2022-06-01", "/api/v1/instances/"+id+"/delete", map[string]*string{}, nil, map[string]interface{}{"instanceId": id, "RegionId": client.RegionId}, true)
		},
		now: time.Now, sleep: time.Sleep, releaseTimeout: 30 * time.Minute, maxPolls: 180,
	}
}

func vpcHologramSweeper(region, vpcID string, api vpcHologramAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_hologram_instance", F: func(actual string) error {
		if region == "" || vpcID == "" || actual != region || api.region != region {
			return fmt.Errorf("Hologram cleanup requires exact VPC and region")
		}
		if api.list == nil || api.get == nil || api.remove == nil || api.now == nil || api.sleep == nil || api.maxPolls < 1 || api.releaseTimeout <= 0 {
			return fmt.Errorf("Hologram cleanup requires bounded injected API")
		}
		rows, err := hologramCompleteInstanceList(api, region)
		if err != nil {
			return err
		}
		var failures []string
		for _, row := range rows {
			selected, err := hologramExclusiveVpc(row, vpcID)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if !selected {
				continue
			}
			id := row["InstanceId"].(string)
			response, err := api.get(id)
			if err != nil {
				failures = append(failures, fmt.Sprintf("Hologram detail %s: %v", id, err))
				continue
			}
			detail, err := hologramVerifiedDetail(response, region, id)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			selected, err = hologramExclusiveVpc(detail, vpcID)
			if err != nil || !selected {
				failures = append(failures, fmt.Sprintf("Hologram %s complete endpoint ownership changed or shared: %v", id, err))
				continue
			}
			if detail["InstanceChargeType"] != "PostPaid" {
				failures = append(failures, fmt.Sprintf("Hologram %s billing requires manual release; not PostPaid", id))
				continue
			}
			response, err = api.remove(id)
			if err != nil && !hologramExactNotFound(err) {
				failures = append(failures, fmt.Sprintf("Hologram delete %s: %v", id, err))
				continue
			}
			if err == nil {
				if err = hologramSuccessfulResponse(response); err != nil {
					failures = append(failures, err.Error())
					continue
				}
			}
			if err = hologramWaitReleased(api, region, vpcID, id); err != nil {
				failures = append(failures, err.Error())
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("Hologram owner cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}

func hologramSuccessfulResponse(response map[string]interface{}) error {
	success, ok := response["Success"].(bool)
	if !ok || !success {
		return fmt.Errorf("Hologram response did not explicitly succeed")
	}
	if status, exists := response["HttpStatusCode"]; exists && fmt.Sprint(status) != "200" {
		return fmt.Errorf("Hologram unexpected response status")
	}
	return nil
}

func hologramCompleteInstanceList(api vpcHologramAPI, region string) ([]map[string]interface{}, error) {
	response, err := api.list()
	if err != nil {
		return nil, err
	}
	if err = hologramSuccessfulResponse(response); err != nil {
		return nil, err
	}
	// A changed/partial response must not masquerade as the documented full list.
	for _, key := range []string{"NextToken", "nextToken", "HasMore", "hasMore", "PageNumber", "pageNumber", "PageSize", "pageSize"} {
		if value, exists := response[key]; exists && value != nil {
			return nil, fmt.Errorf("Hologram unexpected pagination field %s", key)
		}
	}
	raw, ok := response["InstanceList"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Hologram missing complete InstanceList")
	}
	rows := make([]map[string]interface{}, 0, len(raw))
	seen := map[string]bool{}
	for _, value := range raw {
		row, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Hologram malformed instance list")
		}
		id, ok := row["InstanceId"].(string)
		if !ok || !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`).MatchString(id) || seen[id] {
			return nil, fmt.Errorf("Hologram missing/duplicate/invalid instance identity")
		}
		if row["RegionId"] != region {
			return nil, fmt.Errorf("Hologram list region missing or mismatched")
		}
		seen[id] = true
		rows = append(rows, row)
	}
	return rows, nil
}

func hologramVerifiedDetail(response map[string]interface{}, region, id string) (map[string]interface{}, error) {
	if err := hologramSuccessfulResponse(response); err != nil {
		return nil, err
	}
	detail, ok := response["Instance"].(map[string]interface{})
	if !ok || detail["InstanceId"] != id || detail["RegionId"] != region {
		return nil, fmt.Errorf("Hologram detail identity/region mismatched")
	}
	return detail, nil
}

func hologramExclusiveVpc(instance map[string]interface{}, vpcID string) (bool, error) {
	endpoints, ok := instance["Endpoints"].([]interface{})
	if !ok || len(endpoints) == 0 {
		return false, fmt.Errorf("Hologram complete Endpoints missing")
	}
	selected, foreign := false, false
	for _, raw := range endpoints {
		endpoint, ok := raw.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("Hologram malformed endpoint")
		}
		kind, ok := endpoint["Type"].(string)
		if !ok {
			return false, fmt.Errorf("Hologram endpoint Type missing")
		}
		switch kind {
		case "Intranet", "Internet":
			if v, exists := endpoint["VpcId"]; exists && v != nil && v != "" {
				return false, fmt.Errorf("Hologram non-VPC endpoint has unexpected VPC binding")
			}
		case "VPCSingleTunnel":
			id, ok := endpoint["VpcId"].(string)
			if !ok || id == "" {
				return false, fmt.Errorf("Hologram VPC endpoint missing VpcId")
			}
			if id == vpcID {
				selected = true
			} else {
				foreign = true
			}
		default:
			return false, fmt.Errorf("Hologram unknown endpoint Type; refusing release")
		}
	}
	if selected && foreign {
		return false, fmt.Errorf("Hologram shared VPC endpoints; instance retained")
	}
	return selected, nil
}

func hologramExactNotFound(err error) bool {
	for {
		if wrapped, ok := err.(*ComplexError); ok {
			err = wrapped.Cause
			continue
		}
		break
	}
	switch e := err.(type) {
	case *tea.SDKError:
		return tea.StringValue(e.Code) == "resource not exists failed" || tea.StringValue(e.Message) == "resource not exists failed"
	case *sdkerrors.ServerError:
		return e.ErrorCode() == "resource not exists failed" || e.Message() == "resource not exists failed"
	}
	return false
}

func hologramWaitReleased(api vpcHologramAPI, region, vpcID, id string) error {
	deadline := api.now().Add(api.releaseTimeout)
	for attempt := 0; attempt < api.maxPolls && api.now().Before(deadline); attempt++ {
		response, err := api.get(id)
		if err != nil {
			if !hologramExactNotFound(err) {
				return fmt.Errorf("Hologram release verification: %w", err)
			}
			rows, listErr := hologramCompleteInstanceList(api, region)
			if listErr != nil {
				return listErr
			}
			present := false
			for _, row := range rows {
				if row["InstanceId"] == id {
					present = true
					break
				}
			}
			if !present {
				return nil
			}
		} else {
			detail, verifyErr := hologramVerifiedDetail(response, region, id)
			if verifyErr != nil {
				return verifyErr
			}
			selected, scopeErr := hologramExclusiveVpc(detail, vpcID)
			if scopeErr != nil || !selected {
				return fmt.Errorf("Hologram ownership changed during release verification")
			}
		}
		if attempt+1 < api.maxPolls && api.now().Before(deadline) {
			api.sleep(10 * time.Second)
		}
	}
	return fmt.Errorf("Hologram exact instance still present after bounded release verification")
}

func hologramUnitInstance(vpc, charge string) map[string]interface{} {
	return map[string]interface{}{"InstanceId": "hg-instance", "RegionId": "cn-hangzhou", "InstanceChargeType": charge, "InstanceStatus": "Running", "Endpoints": []interface{}{map[string]interface{}{"Type": "Intranet"}, map[string]interface{}{"Type": "VPCSingleTunnel", "VpcId": vpc}}}
}
func hologramUnitList(row map[string]interface{}) map[string]interface{} {
	rows := []interface{}{}
	if row != nil {
		rows = append(rows, row)
	}
	return map[string]interface{}{"Success": true, "HttpStatusCode": 200, "InstanceList": rows}
}
func hologramUnitDetail(row map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"Success": true, "HttpStatusCode": 200, "Instance": row}
}
func hologramUnitAPI() vpcHologramAPI {
	now := time.Unix(10, 0)
	return vpcHologramAPI{region: "cn-hangzhou", now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }, releaseTimeout: time.Minute, maxPolls: 3,
		list: func() (map[string]interface{}, error) {
			return hologramUnitList(hologramUnitInstance("target", "PostPaid")), nil
		},
		get: func(string) (map[string]interface{}, error) {
			return hologramUnitDetail(hologramUnitInstance("target", "PostPaid")), nil
		},
		remove: func(string) (map[string]interface{}, error) { return map[string]interface{}{"Success": true}, nil }}
}

func TestUnitVpcHologramExclusivePostpaidOwnerAndAsyncConfirmation(t *testing.T) {
	api := hologramUnitAPI()
	writes, polls, lists := 0, 0, 0
	api.list = func() (map[string]interface{}, error) {
		lists++
		if polls >= 2 {
			return hologramUnitList(nil), nil
		}
		return hologramUnitList(hologramUnitInstance("target", "PostPaid")), nil
	}
	api.get = func(id string) (map[string]interface{}, error) {
		if id != "hg-instance" {
			t.Fatal("wrong exact ID")
		}
		if writes > 0 {
			polls++
			if polls >= 2 {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "resource not exists failed", "message": "resource not exists failed"})
			}
		}
		return hologramUnitDetail(hologramUnitInstance("target", "PostPaid")), nil
	}
	api.remove = func(id string) (map[string]interface{}, error) {
		writes++
		if id != "hg-instance" {
			t.Fatal("wrong delete ID")
		}
		return map[string]interface{}{"Success": true}, nil
	}
	if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err != nil {
		t.Fatal(err)
	}
	if writes != 1 || polls != 2 || lists != 2 {
		t.Fatalf("writes=%d polls=%d lists=%d", writes, polls, lists)
	}
}

func TestUnitVpcHologramSharedVpcAndBillingNeverRelease(t *testing.T) {
	for _, scenario := range []string{"shared", "protected-only", "PrePaid", "unknown-billing", "unknown-type", "unexpected-internet-vpc", "missing-vpc", "missing-endpoints", "changed-detail"} {
		api := hologramUnitAPI()
		writes := 0
		row := hologramUnitInstance("target", "PostPaid")
		switch scenario {
		case "shared":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "VPCSingleTunnel", "VpcId": "protected"})
		case "protected-only":
			row = hologramUnitInstance("protected", "PostPaid")
		case "PrePaid":
			row["InstanceChargeType"] = "PrePaid"
		case "unknown-billing":
			delete(row, "InstanceChargeType")
		case "unknown-type":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "Unknown"})
		case "unexpected-internet-vpc":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "Internet", "VpcId": "protected"})
		case "missing-vpc":
			row["Endpoints"] = []interface{}{map[string]interface{}{"Type": "VPCSingleTunnel"}}
		case "missing-endpoints":
			delete(row, "Endpoints")
		}
		api.list = func() (map[string]interface{}, error) { return hologramUnitList(row), nil }
		api.get = func(string) (map[string]interface{}, error) {
			if scenario == "changed-detail" {
				return hologramUnitDetail(hologramUnitInstance("protected", "PostPaid")), nil
			}
			return hologramUnitDetail(row), nil
		}
		api.remove = func(string) (map[string]interface{}, error) {
			writes++
			return map[string]interface{}{"Success": true}, nil
		}
		err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou")
		if writes != 0 {
			t.Fatal("unsafe release", scenario)
		}
		if scenario != "protected-only" && err == nil {
			t.Fatal("unsafe boundary accepted", scenario)
		}
	}
}

func TestUnitVpcHologramMalformedRegionalListFailsBeforeDeletion(t *testing.T) {
	for _, scenario := range []string{"wrong-region", "missing-region", "duplicate", "missing-list", "success-false", "partial-token", "bad-id"} {
		api := hologramUnitAPI()
		writes := 0
		row := hologramUnitInstance("target", "PostPaid")
		response := hologramUnitList(row)
		switch scenario {
		case "wrong-region":
			row["RegionId"] = "cn-beijing"
		case "missing-region":
			delete(row, "RegionId")
		case "duplicate":
			response["InstanceList"] = []interface{}{row, row}
		case "missing-list":
			delete(response, "InstanceList")
		case "success-false":
			response["Success"] = false
		case "partial-token":
			response["NextToken"] = "not-complete"
		case "bad-id":
			row["InstanceId"] = "../wrong-path"
		}
		api.list = func() (map[string]interface{}, error) { return response, nil }
		api.remove = func(string) (map[string]interface{}, error) { writes++; return nil, nil }
		if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err == nil || writes != 0 {
			t.Fatalf("scenario=%s writes=%d err=%v", scenario, writes, err)
		}
	}
}

func TestUnitVpcHologramOpaqueErrorsAndReleaseWaitAreBounded(t *testing.T) {
	for _, scenario := range []string{"list-denied", "get-denied", "delete-denied", "delete-unknown-success", "poll-denied", "poll-unknown-notfound", "timeout", "notfound-but-listed", "ownership-change"} {
		api := hologramUnitAPI()
		writes, gets := 0, 0
		originalGet := api.get
		api.get = func(id string) (map[string]interface{}, error) {
			gets++
			if scenario == "get-denied" || writes > 0 && scenario == "poll-denied" {
				return nil, fmt.Errorf("permission denied")
			}
			if writes > 0 && scenario == "poll-unknown-notfound" {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "Forbidden", "message": "not found or inaccessible"})
			}
			if writes > 0 && scenario == "notfound-but-listed" {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "resource not exists failed"})
			}
			if writes > 0 && scenario == "ownership-change" {
				return hologramUnitDetail(hologramUnitInstance("protected", "PostPaid")), nil
			}
			return originalGet(id)
		}
		if scenario == "list-denied" {
			api.list = func() (map[string]interface{}, error) { return nil, fmt.Errorf("list denied") }
		}
		api.remove = func(string) (map[string]interface{}, error) {
			writes++
			if scenario == "delete-denied" {
				return nil, fmt.Errorf("delete denied")
			}
			if scenario == "delete-unknown-success" {
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{"Success": true}, nil
		}
		if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err == nil {
			t.Fatal("false success", scenario)
		}
		if writes > 1 || gets > 4 {
			t.Fatalf("unbounded scenario=%s writes=%d gets=%d", scenario, writes, gets)
		}
	}
	api := hologramUnitAPI()
	if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-beijing"); err == nil {
		t.Fatal("region mismatch accepted")
	}
}
