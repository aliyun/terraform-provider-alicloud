package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA RoutineRelatedRecord. >>> Resource test cases, automatically generated.
// Case resource_routineRelatedRecord_test
func TestAccAliCloudESARoutineRelatedRecordresource_routineRelatedRecord_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_related_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineRelatedRecordresource_routineRelatedRecord_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRelatedRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARoutineRelatedRecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineRelatedRecordresource_routineRelatedRecord_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "tftestacc.com",
					"site_id":     "618651327383200",
					"name":        "${alicloud_esa_routine.resource_Routine_Record_test.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "tftestacc.com",
					"site_id":     "618651327383200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AliCloudESARoutineRelatedRecordresource_routineRelatedRecord_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARoutineRelatedRecordresource_routineRelatedRecord_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_routine" "resource_Routine_Record_test" {
  description   = "test-routine2"
  name          = "test-routine2"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

`, name)
}

// Test ESA RoutineRelatedRecord. <<< Resource test cases, automatically generated.
// Test Esa RoutineRelatedRecord. >>> Resource test cases, automatically generated.
// Case resource_routineRelatedRecord_test 11878
func TestAccAliCloudEsaRoutineRelatedRecord_basic11878(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_related_record.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaRoutineRelatedRecordMap11878)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRelatedRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaRoutineRelatedRecordBasicDependence11878)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "chenxin0116.site",
					"site_id":     "${alicloud_esa_site.resource_Site_routineRelatedRecord_test.id}",
					"name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_name": "chenxin0116.site",
						"site_id":     CHECKSET,
						"name":        name,
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

var AlicloudEsaRoutineRelatedRecordMap11878 = map[string]string{
	"record_id": CHECKSET,
}

func AlicloudEsaRoutineRelatedRecordBasicDependence11878(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_routineRelatedRecord_test" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "rongbeitest2"
}

resource "alicloud_esa_routine" "resource_Routine_Record_test" {
  description = "test-routine2"
  name        = "test-routine2"
}

resource "alicloud_esa_site" "resource_Site_routineRelatedRecord_test" {
  site_name   = "chenxin0116.site"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_routineRelatedRecord_test.id
  coverage    = "overseas"
  access_type = "NS"
}


`, name)
}

// Test Esa RoutineRelatedRecord. <<< Resource test cases, automatically generated.
