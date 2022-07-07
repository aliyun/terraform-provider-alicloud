package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	packageIds := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func TestAccAlicloudCommonBandwidthPackage_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence0)
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.change.ids.0}",
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
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
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

func TestAccAlicloudCommonBandwidthPackage_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
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
					"bandwidth":            `10`,
					"isp":                  "BGP",
					"internet_charge_type": "PayByBandwidth",
					"ratio":                `100`,
					"name":                 name,
					"description":          name,
					"zone":                 "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":            "10",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",
						"description":          name,
						"ratio":                "100",
						"name":                 name,
						"zone":                 CHECKSET,
						"resource_group_id":    CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "zone"},
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackage_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
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
					"bandwidth":              `10`,
					"isp":                    "BGP",
					"internet_charge_type":   "PayByBandwidth",
					"ratio":                  `100`,
					"bandwidth_package_name": name,
					"description":            name,
					"zone":                   "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":            "10",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",

						"ratio":                  "100",
						"bandwidth_package_name": name,
						"description":            name,
						"zone":                   CHECKSET,
						"resource_group_id":      CHECKSET,
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

func TestAccAlicloudCommonBandwidthPackage_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudCommonBandwidthPackageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scommonbandwidthpackage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCommonBandwidthPackageBasicDependence1)
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
					"bandwidth":              `10`,
					"isp":                    "BGP",
					"internet_charge_type":   "PayByDominantTraffic",
					"ratio":                  `100`,
					"bandwidth_package_name": name,
					"description":            name,
					"zone":                   "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "10",
						"isp":                    "BGP",
						"internet_charge_type":   "PayByDominantTraffic",
						"ratio":                  "100",
						"bandwidth_package_name": name,
						"description":            name,
						"zone":                   CHECKSET,
						"resource_group_id":      CHECKSET,
					}),
				),
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
  name_regex = "default"
}
data "alicloud_resource_manager_resource_groups" "change" {
  name_regex = "terraformci"
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
