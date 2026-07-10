package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Cms AggTaskGroup. >>> Resource test cases, automatically generated.
// Case aggTaskGroup 8019
func TestAccAliCloudCmsAggTaskGroup_basic8019(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_agg_task_group.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAggTaskGroupMap8019)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAggTaskGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	toTime := time.Now().AddDate(0, 0, 2).Unix()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAggTaskGroupBasicDependence8019)
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
					"source_prometheus_id":  "${alicloud_cms_prometheus_instance.default.0.id}",
					"target_prometheus_id":  "${alicloud_cms_prometheus_instance.default.1.id}",
					"agg_task_group_name":   name,
					"agg_task_group_config": `groups:\n- name: \"node.rules\"\n  interval: \"60s\"\n  rules:\n  - record: \"node_namespace_pod:kube_pod_info:\"\n    expr: \"max(label_replace(kube_pod_info{job=\\\"kubernetes-pods-kube-state-metrics\\\" }, \\\"pod\\\", \\\"$1\\\", \\\"pod\\\", \\\"(.*)\\\")) by (node, namespace, pod, cluster)\"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_prometheus_id":  CHECKSET,
						"target_prometheus_id":  CHECKSET,
						"agg_task_group_name":   name,
						"agg_task_group_config": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agg_task_group_config": `groups:\n- name: \"node.rules\"\n  interval: \"80s\"\n  rules:\n  - record: \"node_namespace_pod:kube_pod_info:\"\n    expr: \"max(label_replace(kube_pod_info{job=\\\"kubernetes-pods-kube-state-metrics\\\" }, \\\"pod\\\", \\\"$1\\\", \\\"pod\\\", \\\"(.*)\\\")) by (node, namespace, pod, cluster)\"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agg_task_group_config": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agg_task_group_config_type": "RecordingRuleYaml",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agg_task_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agg_task_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_retries": "18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_retries": "18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_run_time_in_seconds": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_run_time_in_seconds": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"precheck_string": `{\"policy\":\"skip\",\"prometheusId\":\"` + "${alicloud_cms_prometheus_instance.default.0.id}" + `\",\"query\":\"noPrecheck\",\"threshold\":0.5,\"timeout\":20,\"type\":\"none\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"precheck_string": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_mode": "Cron",
					"cron_expr":     "0/1 * * * *",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_mode": "Cron",
						"cron_expr":     "0/1 * * * *",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cron_expr": "0/2 * * * *",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cron_expr": "0/2 * * * *",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_time_expr": "@s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_time_expr": "@s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"to_time": toTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"to_time": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agg_task_group_config_type", "override_if_exists"},
			},
		},
	})
}

func TestAccAliCloudCmsAggTaskGroup_basic8019_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_agg_task_group.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAggTaskGroupMap8019)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAggTaskGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	toTime := time.Now().AddDate(0, 0, 2).Unix()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAggTaskGroupBasicDependence8019)
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
					"source_prometheus_id":       "${alicloud_cms_prometheus_instance.default.0.id}",
					"target_prometheus_id":       "${alicloud_cms_prometheus_instance.default.1.id}",
					"agg_task_group_name":        name,
					"agg_task_group_config":      `groups:\n- name: \"node.rules\"\n  interval: \"60s\"\n  rules:\n  - record: \"node_namespace_pod:kube_pod_info:\"\n    expr: \"max(label_replace(kube_pod_info{job=\\\"kubernetes-pods-kube-state-metrics\\\" }, \\\"pod\\\", \\\"$1\\\", \\\"pod\\\", \\\"(.*)\\\")) by (node, namespace, pod, cluster)\"`,
					"agg_task_group_config_type": "RecordingRuleYaml",
					"cron_expr":                  "0/1 * * * *",
					"delay":                      "60",
					"description":                name,
					"max_retries":                "18",
					"max_run_time_in_seconds":    "200",
					"precheck_string":            `{\"policy\":\"skip\",\"prometheusId\":\"` + "${alicloud_cms_prometheus_instance.default.0.id}" + `\",\"query\":\"noPrecheck\",\"threshold\":0.5,\"timeout\":20,\"type\":\"none\"}`,
					"schedule_mode":              "Cron",
					"schedule_time_expr":         "@s",
					"status":                     "Running",
					"to_time":                    toTime,
					"override_if_exists":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_prometheus_id":    CHECKSET,
						"target_prometheus_id":    CHECKSET,
						"agg_task_group_name":     name,
						"agg_task_group_config":   CHECKSET,
						"cron_expr":               "0/1 * * * *",
						"delay":                   "60",
						"description":             name,
						"max_retries":             "18",
						"max_run_time_in_seconds": "200",
						"precheck_string":         CHECKSET,
						"schedule_mode":           "Cron",
						"schedule_time_expr":      "@s",
						"status":                  "Running",
						"to_time":                 CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agg_task_group_config_type", "override_if_exists"},
			},
		},
	})
}

var AliCloudCmsAggTaskGroupMap8019 = map[string]string{
	"agg_task_group_id":       CHECKSET,
	"delay":                   CHECKSET,
	"max_retries":             CHECKSET,
	"max_run_time_in_seconds": CHECKSET,
	"precheck_string":         CHECKSET,
	"region_id":               CHECKSET,
	"schedule_mode":           CHECKSET,
	"schedule_time_expr":      CHECKSET,
	"status":                  CHECKSET,
}

func AliCloudCmsAggTaskGroupBasicDependence8019(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_prometheus_instance" "default" {
  count                    = 2
  prometheus_instance_name = "${var.name}_${count.index}"
  workspace                = alicloud_cms_workspace.default.id
}
`, name)
}

// Test Cms AggTaskGroup. <<< Resource test cases, automatically generated.
