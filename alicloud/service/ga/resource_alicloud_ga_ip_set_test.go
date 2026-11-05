package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaIpSet_basic(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_ip_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaIpSetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaIpSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", AliCloudGaIpSetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"accelerate_region_id": acctest.DefaultRegionToTest,
					"bandwidth":            "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"accelerate_region_id": acctest.DefaultRegionToTest,
						"bandwidth":            "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": `10`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaIpSet_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_ip_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaIpSetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaIpSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", AliCloudGaIpSetBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${alicloud_ga_accelerator.default.id}",
					"accelerate_region_id": acctest.DefaultRegionToTest,
					"bandwidth":            "5",
					"ip_version":           "IPv4",
					"isp_type":             "BGP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"accelerate_region_id": acctest.DefaultRegionToTest,
						"bandwidth":            "5",
						"ip_version":           "IPv4",
						"isp_type":             "BGP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": `10`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudGaIpSetMap = map[string]string{
	"ip_address_list.#": acctest.CHECKSET,
	"ip_version":        "IPv4",
	"status":            "active",
}

func AliCloudGaIpSetBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Basic"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
  		auto_pay               = true
  		auto_use_coupon        = true
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}
`, name)
}

func AliCloudGaIpSetBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ga_accelerator" "default" {
  		accelerator_name       = var.name
  		description            = var.name
  		spec                   = "1"
  		auto_use_coupon        = true
  		duration               = 1
  		bandwidth_billing_type = "CDT"
	}
`, name)
}
