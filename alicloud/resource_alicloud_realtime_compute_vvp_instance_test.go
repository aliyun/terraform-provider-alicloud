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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "ha", "ha_zone_id", "ha_vswitch_ids", "ha_resource_spec", "namespace_resource_specs"},
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
// lintignore: AT001
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
					"ha_resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "4",
						},
					},
					"ha":                true,
					"ha_vswitch_ids":    []string{"${data.alicloud_vswitches.default.ids.0}"},
					"vvp_instance_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_ids": []string{
						"${data.alicloud_vswitches.default.ids.0}"},
					"zone_id":               "cn-hangzhou-i",
					"payment_type":          "Subscription",
					"pricing_cycle":         "Month",
					"duration":              "1",
					"renew_status":          "AutoRenewal",
					"auto_renew_duration":   1,
					"renewal_duration_unit": "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "Subscription",
						"pricing_cycle":     "Month",
						"duration":          "1",
						"ha":                "true",
						"ha_vswitch_ids.#":  "1",
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
					"ha_resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "8",
						},
					},
					"ha_zone_id":          "cn-hangzhou-i",
					"auto_renew_duration": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_zone_id":             "cn-hangzhou-i",
						"ha_resource_spec.0.cpu": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "ha", "ha_zone_id", "ha_vswitch_ids", "ha_resource_spec", "namespace_resource_specs"},
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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "ha", "ha_zone_id", "ha_vswitch_ids", "ha_resource_spec", "namespace_resource_specs"},
			},
		},
	})
}

// Case 4594  twin
// lintignore: AT001
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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "ha", "ha_zone_id", "ha_vswitch_ids", "ha_resource_spec", "namespace_resource_specs"},
			},
		},
	})
}

func AlicloudRealtimeComputeVvpInstanceBasicDependence4594_intl(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = "ap-southeast-1a"
}

resource "alicloud_oss_bucket" "defaultOSS" {
  bucket = var.name
}

`, name)
}

// TestAccAliCloudRealtimeComputeVvpInstance_basic4594_intl tests the renewal
// trio (renew_status / auto_renew_duration / renewal_duration_unit) on an
// international-site subscription instance. Skipped on domestic accounts.
// lintignore: AT001
func TestAccAliCloudRealtimeComputeVvpInstance_basic4594_intl(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594_intl)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
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
					"vvp_instance_name":     name,
					"vpc_id":                "${alicloud_vpc.default.id}",
					"vswitch_ids":           []string{"${alicloud_vswitch.default.id}"},
					"zone_id":               "ap-southeast-1a",
					"payment_type":          "Subscription",
					"pricing_cycle":         "Month",
					"duration":              "1",
					"renew_status":          "AutoRenewal",
					"auto_renew_duration":   1,
					"renewal_duration_unit": "M",
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
					"auto_renew_duration": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "ha", "ha_zone_id", "ha_vswitch_ids", "ha_resource_spec", "namespace_resource_specs"},
			},
		},
	})
}

// Case 4636 order
// TestAccAliCloudRealtimeComputeVvpInstance_vswitchIDsOrder verifies that
// reordering the configurable vswitch_ids TypeList (ForceNew) yields a
// non-empty plan and converges after applying the reordered configuration.
// See scripts/collection-order/README.md for the required three-step sequence.
func TestAccAliCloudRealtimeComputeVvpInstance_vswitchIDsOrder(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceOrderDependence)
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
					"vpc_id":            "${alicloud_vpc.default.id}",
					"zone_id":           "cn-hangzhou-i",
					"vswitch_ids": []string{
						"${alicloud_vswitch.first.id}",
						"${alicloud_vswitch.second.id}",
					},
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${alicloud_vswitch.second.id}",
						"${alicloud_vswitch.first.id}",
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${alicloud_vswitch.second.id}",
						"${alicloud_vswitch.first.id}",
					},
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Case 4594 order
// lintignore: AT001
// TestAccAliCloudRealtimeComputeVvpInstance_haVswitchIDsOrder verifies that
// reordering the configurable ha_vswitch_ids TypeList yields a non-empty plan
// and converges after applying the reordered configuration. See
// scripts/collection-order/README.md for the required three-step sequence.
func TestAccAliCloudRealtimeComputeVvpInstance_haVswitchIDsOrder(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeVvpInstanceOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
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
					"ha_resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "4",
						},
					},
					"ha":         true,
					"ha_zone_id": "cn-hangzhou-i",
					"ha_vswitch_ids": []string{
						"${alicloud_vswitch.first.id}",
						"${alicloud_vswitch.second.id}",
					},
					"vvp_instance_name": name,
					"vpc_id":            "${alicloud_vpc.default.id}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.first.id}",
					},
					"zone_id":       "cn-hangzhou-i",
					"payment_type":  "Subscription",
					"pricing_cycle": "Month",
					"duration":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vswitch_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_vswitch_ids": []string{
						"${alicloud_vswitch.second.id}",
						"${alicloud_vswitch.first.id}",
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_vswitch_ids": []string{
						"${alicloud_vswitch.second.id}",
						"${alicloud_vswitch.first.id}",
					},
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func AlicloudRealtimeComputeVvpInstanceOrderDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "first" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = "cn-hangzhou-i"
}

resource "alicloud_vswitch" "second" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.2.0/24"
  zone_id    = "cn-hangzhou-i"
}

resource "alicloud_oss_bucket" "defaultOSS" {
  bucket = var.name
}

`, name)
}

// Test RealtimeCompute VvpInstance. <<< Resource test cases, automatically generated.
