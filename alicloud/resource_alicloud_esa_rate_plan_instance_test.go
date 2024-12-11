package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa RatePlanInstance. >>> Resource test cases, automatically generated.
// Case 套餐_2.0 8489
func TestAccAliCloudEsaRatePlanInstance_basic8489(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rate_plan_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaRatePlanInstanceMap8489)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRatePlanInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sesarateplaninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaRatePlanInstanceBasicDependence8489)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":         "NS",
					"auto_renew":   "true",
					"period":       "1",
					"payment_type": "Subscription",
					"coverage":     "overseas",
					"plan_name":    "basic",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":         "NS",
						"auto_renew":   "true",
						"period":       "1",
						"payment_type": "Subscription",
						"coverage":     "overseas",
						"plan_name":    "basic",
						"auto_pay":     "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": "medium",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "coverage", "period", "type"},
			},
		},
	})
}

var AlicloudEsaRatePlanInstanceMap8489 = map[string]string{
	"status":          CHECKSET,
	"create_time":     CHECKSET,
	"instance_status": CHECKSET,
}

func AlicloudEsaRatePlanInstanceBasicDependence8489(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Esa RatePlanInstance. <<< Resource test cases, automatically generated.

func TestAccAliCloudEsaRatePlanInstance_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rate_plan_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaRatePlanInstanceMap8489)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRatePlanInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sesarateplaninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaRatePlanInstanceBasicDependence8489)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":         "NS",
					"auto_renew":   "true",
					"period":       "1",
					"payment_type": "Subscription",
					"coverage":     "overseas",
					"plan_name":    "entranceplan_intl",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":         "NS",
						"auto_renew":   "true",
						"period":       "1",
						"payment_type": "Subscription",
						"coverage":     "overseas",
						"plan_name":    "entranceplan_intl",
						"auto_pay":     "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "basicplan_intl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": "basicplan_intl",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "coverage", "period", "type"},
			},
		},
	})
}
