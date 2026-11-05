package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Esa RatePlanInstance. >>> Resource test cases, automatically generated.
// Case 套餐_副本 8110
func TestAccAliCloudEsaRatePlanInstance_basic8110(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rate_plan_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRatePlanInstanceMap8110)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRatePlanInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRatePlanInstanceBasicDependence8110)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"period":       "1",
					"plan_name":    "basic",
					"auto_pay":     "true",
					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":       "1",
						"plan_name":    "basic",
						"auto_pay":     "true",
						"payment_type": "Subscription",
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

func TestAccAliCloudEsaRatePlanInstance_basic8110_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rate_plan_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRatePlanInstanceMap8110)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRatePlanInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRatePlanInstanceBasicDependence8110)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "coverage", "period", "type"},
			},
		},
	})
}

var AliCloudEsaRatePlanInstanceMap8110 = map[string]string{
	"status":          acctest.CHECKSET,
	"create_time":     acctest.CHECKSET,
	"instance_status": acctest.CHECKSET,
}

func AliCloudEsaRatePlanInstanceBasicDependence8110(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

var AliCloudEsaRatePlanInstanceMap8489 = map[string]string{
	"status":          acctest.CHECKSET,
	"create_time":     acctest.CHECKSET,
	"instance_status": acctest.CHECKSET,
}

func AliCloudEsaRatePlanInstanceBasicDependence8489(name string) string {
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRatePlanInstanceMap8489)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRatePlanInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sesarateplaninstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRatePlanInstanceBasicDependence8489)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithAccountSiteType(t, acctest.IntlSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"period":       "1",
					"payment_type": "Subscription",
					"plan_name":    "entranceplan_intl",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":       "1",
						"payment_type": "Subscription",
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

func TestAccAliCloudEsaRatePlanInstance_intl_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rate_plan_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRatePlanInstanceMap8489)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRatePlanInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sesarateplaninstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRatePlanInstanceBasicDependence8489)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithAccountSiteType(t, acctest.IntlSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":         "NS",
					"auto_renew":   "true",
					"period":       "1",
					"payment_type": "Subscription",
					"coverage":     "overseas",
					"plan_name":    "basicplan_intl",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":         "NS",
						"auto_renew":   "true",
						"period":       "1",
						"payment_type": "Subscription",
						"coverage":     "overseas",
						"plan_name":    "basicplan_intl",
						"auto_pay":     "true",
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
