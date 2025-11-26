package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute VvpInstance. >>> Resource test cases, automatically generated.
// Case 4636
func TestAccAliCloudRealtimeComputeVvpInstance_basic4636(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4636)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-realtimecomputevvpinstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage": []map[string]interface{}{
						{
							"oss": []map[string]interface{}{
								{
									"bucket": "${alicloud_oss_bucket.defaultOSS.bucket}",
								},
							},
						},
					},
					"vvp_instance_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":           "cn-hangzhou-i",
					"vswitch_ids": []string{
						"${data.alicloud_vswitches.default.ids.0}"},
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "zone_id"},
			},
		},
	})
}

var AlicloudRealtimeComputeVvpInstanceMap4636 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRealtimeComputeVvpInstanceBasicDependence4636(name string) string {
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
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_oss_bucket" "defaultOSS" {
  bucket = var.name
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

`, name)
}

// Case 4594
func TestAccAliCloudRealtimeComputeVvpInstance_basic4594(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4594)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage": []map[string]interface{}{
						{
							"oss": []map[string]interface{}{
								{
									"bucket": "${alicloud_oss_bucket.defaultOSS.bucket}",
								},
							},
						},
					},
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "8",
						},
					},
					"vvp_instance_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_ids": []string{
						"${data.alicloud_vswitches.default.ids.0}"},
					"zone_id":       "cn-hangzhou-i",
					"payment_type":  "Subscription",
					"pricing_cycle": "Month",
					"duration":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "Subscription",
						"pricing_cycle":     "Month",
						"duration":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "4",
							"memory_gb": "16",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "zone_id"},
			},
		},
	})
}

var AlicloudRealtimeComputeVvpInstanceMap4594 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRealtimeComputeVvpInstanceBasicDependence4594(name string) string {
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
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_oss_bucket" "defaultOSS" {
  bucket = var.name
}


`, name)
}

// Case 4636  twin
func TestAccAliCloudRealtimeComputeVvpInstance_basic4636_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4636)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage": []map[string]interface{}{
						{
							"oss": []map[string]interface{}{
								{
									"bucket": "${alicloud_oss_bucket.defaultOSS.bucket}",
								},
							},
						},
					},
					"vvp_instance_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":           "cn-hangzhou-i",
					"vswitch_ids": []string{
						"${data.alicloud_vswitches.default.ids.0}"},
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "zone_id"},
			},
		},
	})
}

// Case 4594  twin
func TestAccAliCloudRealtimeComputeVvpInstance_basic4594_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4594)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeVvpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{14})
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage": []map[string]interface{}{
						{
							"oss": []map[string]interface{}{
								{
									"bucket": "${alicloud_oss_bucket.defaultOSS.bucket}",
								},
							},
						},
					},
					"vvp_instance_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_ids": []string{
						"${data.alicloud_vswitches.default.ids.0}"},
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "4",
							"memory_gb": "16",
						},
					},
					"zone_id":       "cn-hangzhou-i",
					"payment_type":  "Subscription",
					"pricing_cycle": "Month",
					"duration":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "Subscription",
						"pricing_cycle":     "Month",
						"duration":          "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "zone_id"},
			},
		},
	})
}

// Test RealtimeCompute VvpInstance. <<< Resource test cases, automatically generated.
