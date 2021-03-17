package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaBandwidthPackage_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudGaBandwidthPackageMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaBandwidthPackageBasicDependence)
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
					"bandwidth":      `100`,
					"type":           "Basic",
					"bandwidth_type": "Basic",
					"billing_type":   "PayBy95",
					"payment_type":   "PayAsYouGo",
					"ratio":          "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":      "100",
						"type":           "Basic",
						"bandwidth_type": "Basic",
						"billing_type":   "PayBy95",
						"payment_type":   "PayAsYouGo",
						"ratio":          "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_type", "payment_type", "ratio"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_type": "Enhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_type": "Enhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "bandwidthpackageDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "bandwidthpackageDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": "${var.name}",
					"description":            "bandwidthpackage",
					"bandwidth":              "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name,
						"description":            "bandwidthpackage",
						"bandwidth":              "50",
					}),
				),
			},
		},
	})
}

var AlicloudGaBandwidthPackageMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaBandwidthPackageBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
