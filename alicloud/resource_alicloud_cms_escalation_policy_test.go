package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms EscalationPolicy. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsEscalationPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_escalation_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEscalationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEscalationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccescalation%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEscalationPolicyBasicDependence)
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
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"name":        name,
					"enable":      "true",
					"description": "Terraform managed escalation policy",
					"escalation_stage_list": []map[string]interface{}{
						{
							"index":                 "1",
							"cycle_notify_interval": "5",
							"cycle_notify_count":    "3",
							"trigger_delay":         "10",
							"target_incident_state": "acknowledged",
							"notify_channels": []map[string]interface{}{
								{
									"channel_type":         "DING",
									"receivers":            []string{"tf-test-group"},
									"enabled_sub_channels": []string{},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"time_zone":            "Asia/Shanghai",
									"start_time_in_minute": "0",
									"end_time_in_minute":   "1439",
									"day_in_week":          []string{"1", "2", "3", "4", "5"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":   CHECKSET,
						"name":        name,
						"enable":      "true",
						"description": "Terraform managed escalation policy",
						"uuid":        CHECKSET,
						"create_time": CHECKSET,
						"update_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"name":        fmt.Sprintf("%s-updated", name),
					"enable":      "false",
					"description": fmt.Sprintf("%s-updated", name),
					"escalation_stage_list": []map[string]interface{}{
						{
							"index":                 "2",
							"cycle_notify_interval": "10",
							"cycle_notify_count":    "5",
							"trigger_delay":         "20",
							"target_incident_state": "resolved",
							"notify_channels": []map[string]interface{}{
								{
									"channel_type":         "SMS",
									"receivers":            []string{"tf-test-group-2"},
									"enabled_sub_channels": []string{"DING"},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"time_zone":            "UTC",
									"start_time_in_minute": "60",
									"end_time_in_minute":   "1380",
									"day_in_week":          []string{"6", "7"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable":      "false",
						"description": fmt.Sprintf("%s-updated", name),
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

func TestAccAliCloudCmsEscalationPolicy_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_escalation_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEscalationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEscalationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccescalationtwin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEscalationPolicyBasicDependence)
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
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"name":        name,
					"enable":      "true",
					"description": "Terraform twin escalation policy",
					"escalation_stage_list": []map[string]interface{}{
						{
							"index":                 "1",
							"cycle_notify_interval": "5",
							"cycle_notify_count":    "3",
							"trigger_delay":         "10",
							"target_incident_state": "acknowledged",
							"notify_channels": []map[string]interface{}{
								{
									"channel_type":         "DING",
									"receivers":            []string{"tf-test-group"},
									"enabled_sub_channels": []string{},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"time_zone":            "Asia/Shanghai",
									"start_time_in_minute": "0",
									"end_time_in_minute":   "1439",
									"day_in_week":          []string{"1", "2", "3", "4", "5"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace": CHECKSET,
						"name":      name,
						"enable":    "true",
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

func TestAccAliCloudCmsEscalationPolicy_disappear(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_escalation_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEscalationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEscalationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccescalationdis%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEscalationPolicyBasicDependence)
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
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"name":        name,
					"enable":      "true",
					"description": "Terraform disappear escalation policy",
					"escalation_stage_list": []map[string]interface{}{
						{
							"index":                 "1",
							"cycle_notify_interval": "5",
							"cycle_notify_count":    "3",
							"trigger_delay":         "10",
							"target_incident_state": "acknowledged",
							"notify_channels": []map[string]interface{}{
								{
									"channel_type":         "DING",
									"receivers":            []string{"tf-test-group"},
									"enabled_sub_channels": []string{},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"time_zone":            "Asia/Shanghai",
									"start_time_in_minute": "0",
									"end_time_in_minute":   "1439",
									"day_in_week":          []string{"1", "2", "3", "4", "5"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace": CHECKSET,
						"name":      name,
						"enable":    "true",
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

var AliCloudCmsEscalationPolicyMap = map[string]string{
	"create_time": CHECKSET,
	"uuid":        CHECKSET,
	"update_time": CHECKSET,
}

func AliCloudCmsEscalationPolicyBasicDependence(name string) string {
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
