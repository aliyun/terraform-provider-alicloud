package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens CommonBandwidthPackage. >>> Resource test cases.
func TestAccAliCloudEnsCommonBandwidthPackage_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsCommonBandwidthPackageMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsCommonBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senscbwp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsCommonBandwidthPackageBasicDependence)
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
					"bandwidth":     "10",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"description":   name,
					"name":          name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":     "10",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"description":   name,
						"name":          name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("%s-update", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s-update", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("%s-update", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("%s-update", name),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEnsCommonBandwidthPackageMap = map[string]string{
	"status":        CHECKSET,
	"creation_time": CHECKSET,
}

func AliCloudEnsCommonBandwidthPackageBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Ens CommonBandwidthPackage. <<< Resource test cases.
