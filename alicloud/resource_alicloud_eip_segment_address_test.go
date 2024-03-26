package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Eip SegmentAddress. >>> Resource test cases, automatically generated.
// Case 3419
func TestAccAlicloudEipSegmentAddress_basic3419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_segment_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEipSegmentAddressMap3419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipSegmentAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipsegmentaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEipSegmentAddressBasicDependence3419)
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
					"eip_mask": "28",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_mask": "28",
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
				ImportStateVerifyIgnore: []string{"bandwidth", "eip_mask", "internet_charge_type", "isp", "netmode"},
			},
		},
	})
}

var AlicloudEipSegmentAddressMap3419 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEipSegmentAddressBasicDependence3419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3419  twin
func TestAccAlicloudEipSegmentAddress_basic3419_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_segment_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEipSegmentAddressMap3419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipSegmentAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipsegmentaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEipSegmentAddressBasicDependence3419)
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
				ImportStateVerifyIgnore: []string{"bandwidth", "eip_mask", "internet_charge_type", "isp", "netmode"},
			},
		},
	})
}

// Test Eip SegmentAddress. <<< Resource test cases, automatically generated.
