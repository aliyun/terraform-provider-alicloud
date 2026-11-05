package eip_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eip"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eip SegmentAddress. >>> Resource test cases, automatically generated.
// Case 3419
func TestAccAliCloudEipSegmentAddress_basic3419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_segment_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipSegmentAddressMap3419)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipSegmentAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipsegmentaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipSegmentAddressBasicDependence3419)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_mask":          "28",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_mask":          "28",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_mask":             "28",
					"bandwidth":            "5",
					"isp":                  "BGP",
					"internet_charge_type": "PayByBandwidth",
					"netmode":              "public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_mask":             "28",
						"bandwidth":            "5",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",
						"netmode":              "public",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bandwidth", "eip_mask", "internet_charge_type", "isp", "netmode", "resource_group_id"},
			},
		},
	})
}

var AlicloudEipSegmentAddressMap3419 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEipSegmentAddressBasicDependence3419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}

// Case 3419  twin
func TestAccAliCloudEipSegmentAddress_basic3419_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_segment_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipSegmentAddressMap3419)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipSegmentAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipsegmentaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipSegmentAddressBasicDependence3419)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.VpcDhcpOptionsSetSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_mask":             "28",
					"bandwidth":            "5",
					"isp":                  "BGP",
					"internet_charge_type": "PayByBandwidth",
					"netmode":              "public",
					"zone":                 "cn-hangzhou-i",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_mask":             "28",
						"bandwidth":            "5",
						"isp":                  "BGP",
						"internet_charge_type": "PayByBandwidth",
						"netmode":              "public",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bandwidth", "eip_mask", "internet_charge_type", "isp", "netmode"},
			},
		},
	})
}

// Test Eip SegmentAddress. <<< Resource test cases, automatically generated.
