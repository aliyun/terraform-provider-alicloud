package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_common_bandwidth_package", &resource.Sweeper{
		Name: "alicloud_common_bandwidth_package",
		F:    testSweepCommonBandwidthPackage,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_common_bandwidth_package_attachment",
		},
	})
}

func testSweepCommonBandwidthPackage(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeCommonBandwidthPackages"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	packageIds := make([]string, 0)
	for {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return fmt.Errorf("Error retrieving CommonBandwidthPackages: %s", err)
		}

		resp, err := jsonpath.Get("$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			id := fmt.Sprint(item["BandwidthPackageId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Common Bandwidth Package: %s (%s)", name, id)
				continue
			}
			packageIds = append(packageIds, id)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, id := range packageIds {
		log.Printf("[INFO] Deleting Common Bandwidth Package: (%s)", id)
		request := map[string]interface{}{
			"BandwidthPackageId": id,
		}
		action := "DeleteCommonBandwidthPackage"
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Common Bandwidth Package (%s): %v", id, err)
		}
	}
	return nil
}

func TestAccAliCloudCommonBandwidthPackage_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":            `10`,
					"isp":                  "BGP",
					"internet_charge_type": "PayByBandwidth",
					"ratio":                `100`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":            "10",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",
						"ratio":                "100",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name + "1",
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
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              `10`,
					"description":            name,
					"bandwidth_package_name": "${var.name}",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"deletion_protection":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "10",
						"description":            name,
						"bandwidth_package_name": name,
						"resource_group_id":      CHECKSET,
						"deletion_protection":    "false",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackage_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":            `10`,
					"isp":                  "BGP",
					"internet_charge_type": "PayByBandwidth",
					"ratio":                `100`,
					"name":                 name,
					"description":          name,
					//"zone":                 "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":            "10",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",
						"description":          name,
						"ratio":                "100",
						"name":                 name,
						//"zone":                 CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "1",
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
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackage_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              `10`,
					"isp":                    "BGP",
					"internet_charge_type":   "PayByBandwidth",
					"ratio":                  `100`,
					"bandwidth_package_name": name,
					"description":            name,
					// "zone":                   "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "10",
						"isp":                    "BGP",
						"internet_charge_type":   "PayByBandwidth",
						"ratio":                  "100",
						"bandwidth_package_name": name,
						"description":            name,
						// "zone":                   CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackage_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              `10`,
					"isp":                    "BGP",
					"internet_charge_type":   "PayByDominantTraffic",
					"ratio":                  `100`,
					"bandwidth_package_name": name,
					"description":            name,
					// "zone":                   "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "10",
						"isp":                    "BGP",
						"internet_charge_type":   "PayByDominantTraffic",
						"ratio":                  "100",
						"bandwidth_package_name": name,
						"description":            name,
						// "zone":                   CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackage_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":                 `2`,
					"isp":                       "BGP",
					"internet_charge_type":      "PayBy95",
					"ratio":                     `20`,
					"bandwidth_package_name":    name,
					"description":               name,
					"security_protection_types": []string{"AntiDDoS_Enhanced"},
					// "zone":                      "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":                   "2",
						"isp":                         "BGP",
						"internet_charge_type":        "PayBy95",
						"ratio":                       "20",
						"bandwidth_package_name":      name,
						"description":                 name,
						"security_protection_types.#": "1",
						// "zone":                        CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
		},
	})
}

var AlicloudCommonBandwidthPackageMap0 = map[string]string{
	"isp":                  "BGP",
	"internet_charge_type": "PayByBandwidth",
	"ratio":                "100",
	"deletion_protection":  "false",
}

func AlicloudCommonBandwidthPackageBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

func AlicloudCommonBandwidthPackageBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_resource_manager_resource_groups" "change" {
  name_regex = "terraformci"
}
data "alicloud_zones" "default" {}
`, name)
}

// Test Cbwp CommonBandwidthPackage. >>> Resource test cases, automatically generated.
// Case 3426
func TestAccAliCloudCbwpCommonBandwidthPackage_basic3426(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCbwpCommonBandwidthPackageMap3426)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scbwpcp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCbwpCommonBandwidthPackageBasicDependence3426)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1, 2, 3, 4, 5, 6, 7})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": "tf-testacc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "testupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.change.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": "tf-testacc-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": "tf-testacc-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "test",
					"isp":                  "BGP",
					"bandwidth":            "1000",
					"ratio":                "100",
					"internet_charge_type": "PayByBandwidth",
					"resource_group_id":    "${alicloud_resource_manager_resource_group.default.id}",
					"zone":                 "${data.alicloud_zones.default.zones.0.id}",
					"security_protection_types": []string{
						"AntiDDoS_Enhanced"},
					"bandwidth_package_name": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "test",
						"isp":                         "BGP",
						"bandwidth":                   "1000",
						"ratio":                       "100",
						"internet_charge_type":        "PayByBandwidth",
						"resource_group_id":           CHECKSET,
						"zone":                        CHECKSET,
						"security_protection_types.#": "1",
						"bandwidth_package_name":      "tf-testacc",
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
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

var AlicloudCbwpCommonBandwidthPackageMap3426 = map[string]string{
	"payment_type":         CHECKSET,
	"ratio":                "100",
	"status":               CHECKSET,
	"isp":                  "BGP",
	"create_time":          CHECKSET,
	"deletion_protection":  CHECKSET,
	"internet_charge_type": "PayByTraffic",
}

func AlicloudCbwpCommonBandwidthPackageBasicDependence3426(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {}

resource "alicloud_resource_manager_resource_group" "default" {
  display_name        = "test03"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "change" {
  display_name        = "test04"
  resource_group_name = "${var.name}1"
}


`, name)
}

// Case 3426  twin
func TestAccAliCloudCbwpCommonBandwidthPackage_basic3426_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCbwpCommonBandwidthPackageMap3426)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCbwpCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scbwpcp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCbwpCommonBandwidthPackageBasicDependence3426)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1, 2, 3, 4, 5, 6, 7})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "testupdate",
					"isp":                  "BGP",
					"bandwidth":            "1001",
					"ratio":                "100",
					"internet_charge_type": "PayByBandwidth",
					"resource_group_id":    "${alicloud_resource_manager_resource_group.change.id}",
					"zone":                 "${data.alicloud_zones.default.zones.0.id}",
					"security_protection_types": []string{
						"AntiDDoS_Enhanced"},
					"bandwidth_package_name": "tf-testacc-update",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "testupdate",
						"isp":                         "BGP",
						"bandwidth":                   "1001",
						"ratio":                       "100",
						"internet_charge_type":        "PayByBandwidth",
						"resource_group_id":           CHECKSET,
						"zone":                        CHECKSET,
						"security_protection_types.#": "1",
						"bandwidth_package_name":      "tf-testacc-update",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

// Test Cbwp CommonBandwidthPackage. <<< Resource test cases, automatically generated.
