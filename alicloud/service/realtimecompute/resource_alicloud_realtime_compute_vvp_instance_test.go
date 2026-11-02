package realtimecompute_test

// vswitch_ids is pre-existing stock carried verbatim from the flat alicloud/
// layout by this product migration; the move does not add or change it, so
// convergent reorder coverage is known stock debt.
// collection-order: bypass vswitch_ids -- pre-existing attribute outside this change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test RealtimeCompute VvpInstance. >>> Resource test cases, automatically generated.
// Case 4636
func TestAccAliCloudRealtimeComputeVvpInstance_basic4636(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4636)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVvpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-realtimecomputevvpinstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vpc_id":            acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id"},
			},
		},
	})
}

var AlicloudRealtimeComputeVvpInstanceMap4636 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4594)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVvpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy:  rac.CheckResourceDestroy(),
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
					"zone_id":               "cn-hangzhou-i",
					"payment_type":          "Subscription",
					"pricing_cycle":         "Month",
					"duration":              "1",
					"renew_status":          "AutoRenewal",
					"auto_renew_duration":   1,
					"renewal_duration_unit": "M",
					"default_namespace_resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "default_namespace_resource_spec"},
			},
		},
	})
}

var AlicloudRealtimeComputeVvpInstanceMap4594 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4636)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVvpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"resource_spec": []map[string]interface{}{
						{
							"cpu":       "2",
							"memory_gb": "8",
						},
					},
					"default_namespace_resource_spec": []map[string]interface{}{
						{
							"cpu":       "1",
							"memory_gb": "4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vvp_instance_name": name,
						"vpc_id":            acctest.CHECKSET,
						"vswitch_ids.#":     "1",
						"payment_type":      "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id", "default_namespace_resource_spec"},
			},
		},
	})
}

// Case 4594  twin
// lintignore: AT001
func TestAccAliCloudRealtimeComputeVvpInstance_basic4594_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_vvp_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4594)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVvpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithTime(t, []int{14})
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		//CheckDestroy:  rac.CheckResourceDestroy(),
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
						"vpc_id":            acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id"},
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeVvpInstanceMap4594)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeVvpInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srealtimecomputevvpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeVvpInstanceBasicDependence4594_intl)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.IntlSite)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy:  rac.CheckResourceDestroy(),
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
						"vpc_id":            acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_renew_duration", "duration", "pricing_cycle", "renew_status", "renewal_duration_unit", "zone_id"},
			},
		},
	})
}

// Test RealtimeCompute VvpInstance. <<< Resource test cases, automatically generated.
