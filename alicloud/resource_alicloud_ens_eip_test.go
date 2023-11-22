package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens Eip. >>> Resource test cases, automatically generated.
// Case 5131
func TestAccAliCloudEnsEip_basic5131(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipMap5131)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipBasicDependence5131)
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
					"payment_type":         "PayAsYouGo",
					"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
					"internet_charge_type": "95BandwidthByMonth",
					"eip_name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":         "PayAsYouGo",
						"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
						"internet_charge_type": "95BandwidthByMonth",
						"eip_name":             name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "EipDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "EipDescription_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "EipDescription_UPDATE_autost",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "EipDescription_UPDATE_autost",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "EipDescription_autotest",
					"bandwidth":            "5",
					"isp":                  "cmcc",
					"payment_type":         "PayAsYouGo",
					"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
					"eip_name":             name + "_update",
					"internet_charge_type": "95BandwidthByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "EipDescription_autotest",
						"bandwidth":            "5",
						"isp":                  "cmcc",
						"payment_type":         "PayAsYouGo",
						"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
						"eip_name":             name + "_update",
						"internet_charge_type": "95BandwidthByMonth",
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

var AlicloudEnsEipMap5131 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"bandwidth":   "5",
}

func AlicloudEnsEipBasicDependence5131(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5131  twin
func TestAccAliCloudEnsEip_basic5131_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipMap5131)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipBasicDependence5131)
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
					"description":          "EipDescription_UPDATE_autost",
					"bandwidth":            "6",
					"isp":                  "cmcc",
					"payment_type":         "PayAsYouGo",
					"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
					"eip_name":             name,
					"internet_charge_type": "95BandwidthByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "EipDescription_UPDATE_autost",
						"bandwidth":            "6",
						"isp":                  "cmcc",
						"payment_type":         "PayAsYouGo",
						"ens_region_id":        "cn-chenzhou-telecom_unicom_cmcc",
						"eip_name":             name,
						"internet_charge_type": "95BandwidthByMonth",
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

// Test Ens Eip. <<< Resource test cases, automatically generated.
