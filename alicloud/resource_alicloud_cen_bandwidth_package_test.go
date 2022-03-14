package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_bandwidth_package", &resource.Sweeper{
		Name: "alicloud_cen_bandwidth_package",
		F:    testSweepCenBandwidthPackage,
		Dependencies: []string{
			"alicloud_cen_bandwidth_limit",
		},
	})
}

func testSweepCenBandwidthPackage(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []cbn.CenBandwidthPackage
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		var raw interface{}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenBandwidthPackages(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieving CEN Bandwidth Package: %s", err)
			break
		}
		response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)
		if len(response.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			break
		}
		insts = append(insts, response.CenBandwidthPackages.CenBandwidthPackage...)

		if len(response.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.Name
		id := v.CenBandwidthPackageId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CEN bandwidth package: %s (%s)", name, id)
			continue
		}
		sweeped = true
		if v.Status == string(InUse) {
			log.Printf("[INFO] Deleting CEN bandwidth package attachment: %s (%s)", name, id)
			request := cbn.CreateUnassociateCenBandwidthPackageRequest()
			request.CenId = v.CenIds.CenId[0]
			request.CenBandwidthPackageId = id
			_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.UnassociateCenBandwidthPackage(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to Unassociate CEN bandwidth package (%s (%s)): %s", name, id, err)
			}
		}
		log.Printf("[INFO] Deleting CEN bandwidth package: %s (%s)", name, id)
		request := cbn.CreateDeleteCenBandwidthPackageRequest()
		request.CenBandwidthPackageId = id
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenBandwidthPackage(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN bandwidth package (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 5 seconds to eusure these instances have been deleted.
		time.Sleep(5 * time.Second)
	}
	return nil
}

// Skip this testcase because of the account cannot purchase non-internal products.
func SkipTestAccAlicloudCenBandwidthPackage_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, cenBandwidthPackageBasicMap)

	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cenBwp, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":             "5",
					"geographic_region_ids": []string{"China", "China"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":               "5",
						"geographic_region_ids.#": "1",
					}),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "China"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"geographic_region_ids": []string{"China", "Asia-Pacific"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"geographic_region_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":             "5",
					"geographic_region_ids": []string{"China", "Asia-Pacific"},
					"name":                  fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(cenBandwidthPackageBasicMap),
				),
			},
		},
	})
}

// Skip this testcase because of the account cannot purchase non-internal products.
func SkipTestAccAlicloudCenBandwidthPackage_multi(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, cenBandwidthPackageBasicMap)

	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cenBwp, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_multi)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":             "5",
					"geographic_region_ids": []string{"China", "Asia-Pacific"},
					"name":                  fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(cenBandwidthPackageBasicMap),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_upgrade(t *testing.T) {
	var v cbn.CenBandwidthPackage
	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_upgrade)
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
					"bandwidth":              "5",
					"geographic_region_a_id": "China",
					"geographic_region_b_id": "China",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "5",
						"geographic_region_a_id": "China",
						"geographic_region_b_id": "China",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_bandwidth_package_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_bandwidth_package_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              "5",
					"geographic_region_a_id": "China",
					"geographic_region_b_id": "China",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "5",
						"geographic_region_a_id": "China",
						"geographic_region_b_id": "China",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_basic1(t *testing.T) {
	var v cbn.CenBandwidthPackage
	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_upgrade)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":                  "5",
					"geographic_region_ids":      []string{"China", "China"},
					"cen_bandwidth_package_name": "${var.name}",
					"payment_type":               "PrePaid",
					"period":                     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":                  "5",
						"geographic_region_ids.#":    "1",
						"cen_bandwidth_package_name": name,
						"payment_type":               "PrePaid",
						"period":                     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_basic2(t *testing.T) {
	var v cbn.CenBandwidthPackage
	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_upgrade)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":             "5",
					"geographic_region_ids": []string{"China", "China"},
					"name":                  "${var.name}",
					"charge_type":           "PrePaid",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":               "5",
						"geographic_region_ids.#": "1",
						"name":                    name,
						"charge_type":             "PrePaid",
						"period":                  "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_basic3(t *testing.T) {
	var v cbn.CenBandwidthPackage
	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_upgrade)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":                  "5",
					"geographic_region_ids":      []string{"China", "China"},
					"cen_bandwidth_package_name": "${var.name}",
					"payment_type":               "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":                  "5",
						"geographic_region_ids.#":    "1",
						"cen_bandwidth_package_name": name,
						"payment_type":               "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_basic4(t *testing.T) {
	var v cbn.CenBandwidthPackage
	resourceId := "alicloud_cen_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCen%sBandwidthPackage-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageConfigDependence_upgrade)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":             "5",
					"geographic_region_ids": []string{"China", "China"},
					"name":                  "${var.name}",
					"charge_type":           "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":               "5",
						"geographic_region_ids.#": "1",
						"name":                    name,
						"charge_type":             "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

var cenBandwidthPackageBasicMap = map[string]string{
	"bandwidth": "5",
}

func resourceCenBandwidthPackageConfigDependence(name string) string {
	return ""
}

func resourceCenBandwidthPackageConfigDependence_multi(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cen_bandwidth_package" "default1" {
    name = "%s-multi"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"China"]
}
`, name)
}

func testAccCheckCenBandwidthPackageRegionId(cenBwp *cbn.CenBandwidthPackage, regionAId string, regionBId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		responseRegionAId := convertGeographicRegionAIdResponse(cenBwp.GeographicRegionAId)
		responseRegionBId := convertGeographicRegionBIdResponse(cenBwp.GeographicRegionBId)
		if (responseRegionAId == regionAId && responseRegionBId == regionBId) ||
			(responseRegionAId == regionBId && responseRegionBId == regionAId) {
			return nil
		} else {
			return fmt.Errorf("CEN Bandwidth Package %s geographic region ID error", cenBwp.CenBandwidthPackageId)
		}
	}
}

func resourceCenBandwidthPackageConfigDependence_upgrade(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}`, name)
}
