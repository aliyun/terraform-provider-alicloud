package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ecs SavingPlan. >>> Resource test cases, automatically generated.
// Case 3970
func TestAccAlicloudEcsSavingPlan_basic3970(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_saving_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsSavingPlanMap3970)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSavingPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssavingplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsSavingPlanBasicDependence3970)
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
					"committed_amount": "0.01",
					"purchase_method":  "group",
					"offering_type":    "total",
					"plan_type":        "ecs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "group",
						"offering_type":    "total",
						"plan_type":        "ecs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"committed_amount": "0.01",
					"purchase_method":  "group",
					"offering_type":    "total",
					"plan_type":        "ecs",
					"period":           "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "group",
						"offering_type":    "total",
						"plan_type":        "ecs",
						"period":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_family_set", "period_unit", "purchase_method", "saving_plan_name"},
			},
		},
	})
}

var AlicloudEcsSavingPlanMap3970 = map[string]string{
	"create_time":  CHECKSET,
	"payment_type": CHECKSET,
}

func AlicloudEcsSavingPlanBasicDependence3970(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3682
func TestAccAlicloudEcsSavingPlan_basic3682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_saving_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsSavingPlanMap3682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSavingPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssavingplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsSavingPlanBasicDependence3682)
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
					"committed_amount": "0.01",
					"purchase_method":  "family",
					"offering_type":    "total",
					"plan_type":        "ecs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "family",
						"offering_type":    "total",
						"plan_type":        "ecs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"committed_amount": "0.01",
					"purchase_method":  "family",
					"offering_type":    "total",
					"instance_family":  "ecs.e3",
					"plan_type":        "ecs",
					"period":           "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "family",
						"offering_type":    "total",
						"instance_family":  "ecs.e3",
						"plan_type":        "ecs",
						"period":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_family_set", "period_unit", "purchase_method", "saving_plan_name"},
			},
		},
	})
}

var AlicloudEcsSavingPlanMap3682 = map[string]string{
	"create_time":  CHECKSET,
	"payment_type": CHECKSET,
}

func AlicloudEcsSavingPlanBasicDependence3682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3970  twin
func TestAccAlicloudEcsSavingPlan_basic3970_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_saving_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsSavingPlanMap3970)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSavingPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssavingplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsSavingPlanBasicDependence3970)
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
					"committed_amount": "0.01",
					"purchase_method":  "group",
					"offering_type":    "total",
					"plan_type":        "ecs",
					"period":           "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "group",
						"offering_type":    "total",
						"plan_type":        "ecs",
						"period":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_family_set", "period_unit", "purchase_method", "saving_plan_name"},
			},
		},
	})
}

// Case 3682  twin
func TestAccAlicloudEcsSavingPlan_basic3682_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_saving_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsSavingPlanMap3682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSavingPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssavingplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsSavingPlanBasicDependence3682)
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
					"committed_amount": "0.01",
					"purchase_method":  "family",
					"offering_type":    "total",
					"instance_family":  "ecs.e3",
					"plan_type":        "ecs",
					"period":           "1",
					"period_unit":      "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"committed_amount": "0.01",
						"purchase_method":  "family",
						"offering_type":    "total",
						"instance_family":  "ecs.e3",
						"plan_type":        "ecs",
						"period":           "1",
						"period_unit":      "Year",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_family_set", "period_unit", "purchase_method", "saving_plan_name"},
			},
		},
	})
}

// Test Ecs SavingPlan. <<< Resource test cases, automatically generated.
