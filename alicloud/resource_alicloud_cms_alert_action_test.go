package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccAliCloudCmsAlertAction_basic verifies the create→read→update→import→delete lifecycle.
func TestAccAliCloudCmsAlertAction_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_action.basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, "tf-test-cms-alert-action", func(name string) string {
		return ""
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlertActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_action_name": "tf-test-alert-action-basic",
					"type":              "WEBHOOK",
					"webhook_param": []map[string]interface{}{
						{
							"method":       "GET",
							"url":          "https://httpbin.org/get",
							"content_type": "JSON",
							"headers": map[string]interface{}{
								"X-Custom-Header": "test-value",
							},
						},
					},
					"mns_param": []map[string]interface{}{
						{
							"mns_type":  "queue",
							"name":      "test-mns-queue",
							"region_id": "cn-hangzhou",
						},
					},
					"sls_param": []map[string]interface{}{
						{
							"logstore":  "test-logstore",
							"project":   "test-project",
							"region_id": "cn-hangzhou",
						},
					},
					"ess_param": []map[string]interface{}{
						{
							"ess_group_id": "test-ess-group",
							"ess_rule_id":  "test-ess-rule",
							"region_id":    "cn-hangzhou",
						},
					},
					"fc_param": []map[string]interface{}{
						{
							"function":  "test-fc-function",
							"region_id": "cn-hangzhou",
							"service":   "test-fc-service",
						},
					},
					"pager_duty_param": []map[string]interface{}{
						{
							"key": "test-integration-key",
							"url": "https://events.pagerduty.com/test",
						},
					},
					"fc3_param": []map[string]interface{}{
						{
							"region_id": "cn-hangzhou",
							"function":  "test-fc3-function",
							"qualifier": "test-qualifier",
						},
					},
					"eb_param": []map[string]interface{}{
						{
							"region_id":      "cn-hangzhou",
							"event_bus_name": "test-event-bus",
							"subject":        "test-subject",
							"eb_source":      "test-source",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlertActionExists("alicloud_cms_alert_action.basic", &v),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "alert_action_name", "tf-test-alert-action-basic"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "type", "WEBHOOK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "webhook_param.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "webhook_param.0.url", "https://httpbin.org/get"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "webhook_param.0.content_type", "JSON"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_action_name": "tf-test-alert-action-updated",
					"type":              "WEBHOOK",
					"webhook_param": []map[string]interface{}{
						{
							"method":       "POST",
							"url":          "https://httpbin.org/post",
							"content_type": "FORM",
							"headers": map[string]interface{}{
								"X-Custom-Header": "test-value-updated",
							},
						},
					},
					"mns_param": []map[string]interface{}{
						{
							"mns_type":  "topic",
							"name":      "test-mns-queue-updated",
							"region_id": "cn-beijing",
						},
					},
					"sls_param": []map[string]interface{}{
						{
							"logstore":  "test-logstore-updated",
							"project":   "test-project-updated",
							"region_id": "cn-beijing",
						},
					},
					"ess_param": []map[string]interface{}{
						{
							"ess_group_id": "test-ess-group-updated",
							"ess_rule_id":  "test-ess-rule-updated",
							"region_id":    "cn-beijing",
						},
					},
					"fc_param": []map[string]interface{}{
						{
							"function":  "test-fc-function-updated",
							"region_id": "cn-beijing",
							"service":   "test-fc-service-updated",
						},
					},
					"pager_duty_param": []map[string]interface{}{
						{
							"key": "test-integration-key-updated",
							"url": "https://events.pagerduty.com/test-updated",
						},
					},
					"fc3_param": []map[string]interface{}{
						{
							"region_id": "cn-beijing",
							"function":  "test-fc3-function-updated",
							"qualifier": "test-qualifier-updated",
						},
					},
					"eb_param": []map[string]interface{}{
						{
							"region_id":      "cn-beijing",
							"event_bus_name": "test-event-bus-updated",
							"subject":        "test-subject-updated",
							"eb_source":      "test-source-updated",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlertActionExists("alicloud_cms_alert_action.basic", &v),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "alert_action_name", "tf-test-alert-action-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_action.basic", "webhook_param.0.method", "POST"),
				),
			},
			{
				ResourceName:      "alicloud_cms_alert_action.basic",
				ImportState:       true,
				ImportStateVerify: true,
				// The API only returns the param block matching the action type (WEBHOOK here);
				// the other param blocks are write-only for this type and cannot be read back.
				ImportStateVerifyIgnore: []string{"mns_param", "sls_param", "ess_param", "fc_param", "pager_duty_param", "fc3_param", "eb_param"},
			},
		},
	})
}

func testAccCheckCmsAlertActionExists(n string, v *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsServiceV2 := CmsServiceV2{client}
		object, err := cmsServiceV2.DescribeCmsAlertAction(rs.Primary.ID)
		if err != nil {
			return err
		}
		*v = object
		return nil
	}
}

func testAccCheckCmsAlertActionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_alert_action" {
			continue
		}
		_, err := cmsServiceV2.DescribeCmsAlertAction(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Cms AlertAction %s still exists", rs.Primary.ID)
	}
	return nil
}
