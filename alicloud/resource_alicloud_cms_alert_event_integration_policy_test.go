package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms AlertEventIntegrationPolicy. >>> Resource test cases.

func TestAccAliCloudCmsAlertEventIntegrationPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_event_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlertEventIntegrationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertEventIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertEventIntegrationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_event_integration_policy_name": name,
					"workspace":                           "${alicloud_cms_workspace.default.id}",
					"type":                                "SYS_EVENT",
					"description":                         "tf-test-description",
					"integration_setting":                 "test-integration-setting",
					"enable":                              true,
					"filter_setting": []map[string]interface{}{{
						"conditions": []map[string]interface{}{
							{"field": "name", "value": "test", "op": "eq"},
						},
						"expression": "test-expr",
						"relation":   "and",
					}},
					"transformer_setting": []map[string]interface{}{{
						"source":    "name",
						"target":    "name",
						"type":      "replace",
						"value":     "test",
						"variable":  "name",
						"label_key": "key",
						"reg_exp":   ".*",
						"mapping":   map[string]interface{}{"key": "value"},
						"filter_setting": []map[string]interface{}{{
							"conditions": []map[string]interface{}{
								{"field": "name", "value": "test", "op": "eq"},
							},
							"expression": "test-expr",
							"relation":   "and",
						}},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_event_integration_policy_name": name,
						"type":                                "SYS_EVENT",
						"description":                         "tf-test-description",
						"enable":                              "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_event_integration_policy_name": name + "_update",
					"type":                                "CMS_ALERT",
					"description":                         "tf-test-description-updated",
					"integration_setting":                 "test-integration-setting-updated",
					"enable":                              false,
					"filter_setting": []map[string]interface{}{{
						"conditions": []map[string]interface{}{
							{"field": "updated_field", "value": "updated", "op": "ne"},
						},
						"expression": "updated-expr",
						"relation":   "or",
					}},
					"transformer_setting": []map[string]interface{}{{
						"source":    "updated_source",
						"target":    "updated",
						"type":      "extract",
						"value":     "updated",
						"variable":  "updated",
						"label_key": "key2",
						"reg_exp":   "^.*$",
						"mapping":   map[string]interface{}{"key2": "value2"},
						"filter_setting": []map[string]interface{}{{
							"conditions": []map[string]interface{}{
								{"field": "updated_field", "value": "updated", "op": "ne"},
							},
							"expression": "updated-expr",
							"relation":   "or",
						}},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_event_integration_policy_name": name + "_update",
						"type":                                "CMS_ALERT",
						"description":                         "tf-test-description-updated",
						"enable":                              "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "update_time"},
			},
		},
	})
}

func TestAccAliCloudCmsAlertEventIntegrationPolicy_minimal(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_event_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlertEventIntegrationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertEventIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertEventIntegrationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_event_integration_policy_name": name,
					"workspace":                           "${alicloud_cms_workspace.default.id}",
					"type":                                "CUSTOM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_event_integration_policy_name": name,
						"type":                                "CUSTOM",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "update_time"},
			},
		},
	})
}

var AliCloudCmsAlertEventIntegrationPolicyMap = map[string]string{
	"region_id":   CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
	"user_id":     CHECKSET,
}

func AliCloudCmsAlertEventIntegrationPolicyBasicDependence(name string) string {
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
`, name)
}

// Test Cms AlertEventIntegrationPolicy. <<< Resource test cases.
