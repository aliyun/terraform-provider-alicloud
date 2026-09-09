package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms EvaluationTask. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsEvaluationTask_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_evaluation_task.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEvaluationTaskMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEvaluationTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccetaskbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEvaluationTaskBasicDependence)
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
					"workspace":      "${alicloud_cms_workspace.default.workspace_name}",
					"task_name":      name,
					"task_mode":      "Manual",
					"data_type":      "Metric",
					"data_filter":    name,
					"channel":        name,
					"description":    name,
					"run_strategies": name,
					"status":         "Running",
					"config":         `{\"key\":\"value\"}`,
					"tags":           `{\"env\":\"test\"}`,
					"evaluators": []interface{}{
						map[string]interface{}{
							"name":             name,
							"result_name":      name,
							"result_type":      "Metric",
							"config":           `{\"scope\":\"all\"}`,
							"filters":          `{\"region\":\"cn-hangzhou\"}`,
							"variable_mapping": `{\"metric\":\"cpu\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":      CHECKSET,
						"task_name":      name,
						"task_mode":      "Manual",
						"data_type":      "Metric",
						"data_filter":    name,
						"channel":        name,
						"description":    name,
						"run_strategies": name,
						"status":         "Running",
						"config":         CHECKSET,
						"tags":           CHECKSET,
						"task_id":        CHECKSET,
						"create_time":    CHECKSET,
						"region_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    fmt.Sprintf("%s-updated", name),
					"data_filter":    fmt.Sprintf("%s-updated", name),
					"channel":        fmt.Sprintf("%s-updated", name),
					"run_strategies": fmt.Sprintf("%s-updated", name),
					"status":         "Completed",
					"config":         `{\"key\":\"updated\"}`,
					"tags":           `{\"env\":\"prod\"}`,
					"evaluators": []interface{}{
						map[string]interface{}{
							"name":             fmt.Sprintf("%s-updated", name),
							"result_name":      fmt.Sprintf("%s-updated", name),
							"result_type":      "Log",
							"config":           `{\"scope\":\"updated\"}`,
							"filters":          `{\"region\":\"cn-shanghai\"}`,
							"variable_mapping": `{\"metric\":\"memory\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    fmt.Sprintf("%s-updated", name),
						"data_filter":    fmt.Sprintf("%s-updated", name),
						"channel":        fmt.Sprintf("%s-updated", name),
						"run_strategies": fmt.Sprintf("%s-updated", name),
						"status":         "Completed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "",
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

func TestAccAliCloudCmsEvaluationTask_disappear(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_evaluation_task.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEvaluationTaskMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEvaluationTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccetaskdisp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEvaluationTaskBasicDependence)
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
					"workspace": "${alicloud_cms_workspace.default.workspace_name}",
					"task_name": name,
					"task_mode": "Manual",
					"data_type": "Metric",
					"evaluators": []interface{}{
						map[string]interface{}{
							"name":             name,
							"result_name":      name,
							"variable_mapping": `{}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace": CHECKSET,
						"task_name": name,
						"task_mode": "Manual",
						"data_type": "Metric",
						"task_id":   CHECKSET,
					}),
				),
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				ResourceName:       resourceId,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

var AliCloudCmsEvaluationTaskMap = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"task_id":     CHECKSET,
}

func AliCloudCmsEvaluationTaskBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}
`, name)
}
