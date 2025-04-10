package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eflo ExperimentPlan. >>> Resource test cases, automatically generated.
// Case 测试计划_V2_预发_测试 10572
func TestAccAliCloudEfloExperimentPlan_basic10572(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_experiment_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloExperimentPlanMap10572)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloExperimentPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloExperimentPlanBasicDependence10572)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id": "36",
					"template_id": "${alicloud_eflo_experiment_plan_template.defaultpSZN7t.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_id": CHECKSET,
						"template_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_params"},
			},
		},
	})
}

func TestAccAliCloudEfloExperimentPlan_basic10572_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_experiment_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloExperimentPlanMap10572)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloExperimentPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloExperimentPlanBasicDependence10572)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"resource_id":       "36",
					"plan_name":         name,
					"template_id":       "${alicloud_eflo_experiment_plan_template.defaultpSZN7t.id}",
					"external_params": map[string]interface{}{
						"node": "test",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_id":       CHECKSET,
						"template_id":       CHECKSET,
						"plan_name":         name,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_params"},
			},
		},
	})
}

var AliCloudEfloExperimentPlanMap10572 = map[string]string{
	"create_time":       CHECKSET,
	"plan_name":         CHECKSET,
	"resource_group_id": CHECKSET,
}

func AliCloudEfloExperimentPlanBasicDependence10572(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_experiment_plan_template" "defaultpSZN7t" {
 template_pipeline {
   workload_id   = "2"
   workload_name = "MatMul"
   env_params {
     cpu_per_worker     = "90"
     gpu_per_worker     = "8"
     memory_per_worker  = "500"
     share_memory       = "500"
     worker_num         = "1"
     py_torch_version   = "1"
     gpu_driver_version = "1"
     cuda_version       = "1"
     nccl_version       = "1"
   }
   pipeline_order = "1"
   scene          = "baseline"
 }
 privacy_level        = "private"
 template_name        = var.name
 template_description = var.name
}

`, name)
}

// Test Eflo ExperimentPlan. <<< Resource test cases, automatically generated.
