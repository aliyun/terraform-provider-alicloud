package dataworks_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The DataQualityAlertRule Target references a DataQualityScan, which has no
// Terraform resource and whose CreateDataQualityScan call requires a deeply
// nested Spec/Parameters/RuntimeResource graph (monitor tables, data sources,
// triggers, hooks) that belongs to a separate governed resource. Rather than
// fabricate a scan inside this test, the ACC borrows a pre-existing scan and
// its owning project from the test account via env vars and skips when they
// are not configured. The scan must belong to the borrowed project.
//
//	ALICLOUD_DATA_WORKS_DATA_QUALITY_PROJECT_ID — project id owning the scan
//	ALICLOUD_DATA_WORKS_DATA_QUALITY_SCAN_ID     — DataQualityScan id (in that project)
//
// The receiver_values/extension fields are exercised through an optional
// Dingding webhook URL; when unset the test still covers the receiver_type
// in-place update but cannot prove the value-bearing receiver types.
//
//	ALICLOUD_DATA_WORKS_DATA_QUALITY_RECEIVER_WEBHOOK_URL — a Dingtalk robot
//	    webhook URL used as receiver_values for receiver_type=DingdingUrl
func testAccPreCheckDataQualityAlertRule(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
	acctest.PreCheck(t)
	if v := os.Getenv("ALICLOUD_DATA_WORKS_DATA_QUALITY_PROJECT_ID"); v == "" {
		t.Skip("ALICLOUD_DATA_WORKS_DATA_QUALITY_PROJECT_ID must be set to a project owning a DataQualityScan to run this test")
	}
	if v := os.Getenv("ALICLOUD_DATA_WORKS_DATA_QUALITY_SCAN_ID"); v == "" {
		t.Skip("ALICLOUD_DATA_WORKS_DATA_QUALITY_SCAN_ID must be set to a DataQualityScan id (in the borrowed project) to run this test")
	}
}

func TestAccAliCloudDataWorksDataQualityAlertRule_basic(t *testing.T) {
	testAccPreCheckDataQualityAlertRule(t)

	resourceId := "alicloud_data_works_data_quality_alert_rule.default"
	projectId := os.Getenv("ALICLOUD_DATA_WORKS_DATA_QUALITY_PROJECT_ID")
	scanId := os.Getenv("ALICLOUD_DATA_WORKS_DATA_QUALITY_SCAN_ID")
	webhookURL := os.Getenv("ALICLOUD_DATA_WORKS_DATA_QUALITY_RECEIVER_WEBHOOK_URL")

	steps := []resource.TestStep{
		{
			Config: testAccDataWorksDataQualityAlertRuleBasic(projectId, scanId),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceId, "id"),
				resource.TestCheckResourceAttr(resourceId, "project_id", projectId),
				resource.TestCheckResourceAttr(resourceId, "condition", "results.any { r -> r.status == 'fail' && r.rule.severity == 'High' }"),
				resource.TestCheckResourceAttr(resourceId, "target.type", "DataQualityScan"),
				resource.TestCheckResourceAttr(resourceId, "target.ids.#", "1"),
				resource.TestCheckResourceAttr(resourceId, "notification.channels.#", "1"),
				resource.TestCheckResourceAttr(resourceId, "notification.channels.0", "Mail"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.#", "1"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.receiver_type", "TaskOwner"),
			),
		},
		{
			// In-place update: condition, channels and receiver_type change
			// (TaskOwner -> DataQualityScanOwner) without touching target.ids
			// (R7.6 update coverage of receiver_type).
			Config: testAccDataWorksDataQualityAlertRuleUpdate(projectId, scanId),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceId, "id"),
				resource.TestCheckResourceAttr(resourceId, "condition", "results.any { r -> r.status == 'fail' }"),
				resource.TestCheckResourceAttr(resourceId, "notification.channels.0", "Sms"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.receiver_type", "DataQualityScanOwner"),
			),
		},
	}

	if webhookURL != "" {
		// receiver_values and extension are value-bearing fields that the
		// TaskOwner/DataQualityScanOwner receiver types never set; this step
		// switches to a DingdingUrl receiver carrying both receiver_values
		// (the webhook URL list) and extension ({"atAll":true}), exercising
		// the two previously zero-covered fields (R7.13). It also uses a
		// multi-channel [Mail, Sms] notification that exercises the Read
		// order normalization (R4.19).
		steps = append(steps, resource.TestStep{
			Config: testAccDataWorksDataQualityAlertRuleReceiverValues(projectId, scanId, webhookURL),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceId, "id"),
				resource.TestCheckResourceAttr(resourceId, "notification.channels.#", "2"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.#", "1"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.receiver_type", "DingdingUrl"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.receiver_values.#", "1"),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.receiver_values.0", webhookURL),
				resource.TestCheckResourceAttr(resourceId, "notification.receivers.0.extension", `{"atAll":true}`),
			),
		})
	}

	steps = append(steps, resource.TestStep{
		ResourceName:            resourceId,
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDataQualityAlertRule(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksDataQualityAlertRuleDestroy,
		Steps:                    steps,
	})
}

func testAccCheckDataWorksDataQualityAlertRuleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_data_works_data_quality_alert_rule" {
			continue
		}
		_, err := dataworks.NewService(client).DescribeDataQualityAlertRule(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing data quality alert rule %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("data quality alert rule %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccDataWorksDataQualityAlertRuleBasic(projectId, scanId string) string {
	return fmt.Sprintf(`
resource "alicloud_data_works_data_quality_alert_rule" "default" {
  project_id = %[1]s
  condition  = "results.any { r -> r.status == 'fail' && r.rule.severity == 'High' }"
  target {
    type = "DataQualityScan"
    ids  = [%[2]s]
  }
  notification {
    channels = ["Mail"]
    receivers {
      receiver_type = "TaskOwner"
    }
  }
}
`, projectId, scanId)
}

func testAccDataWorksDataQualityAlertRuleUpdate(projectId, scanId string) string {
	return fmt.Sprintf(`
resource "alicloud_data_works_data_quality_alert_rule" "default" {
  project_id = %[1]s
  condition  = "results.any { r -> r.status == 'fail' }"
  target {
    type = "DataQualityScan"
    ids  = [%[2]s]
  }
  notification {
    channels = ["Sms"]
    receivers {
      receiver_type = "DataQualityScanOwner"
    }
  }
}
`, projectId, scanId)
}

func testAccDataWorksDataQualityAlertRuleReceiverValues(projectId, scanId, webhookURL string) string {
	return fmt.Sprintf(`
resource "alicloud_data_works_data_quality_alert_rule" "default" {
  project_id = %[1]s
  condition  = "results.any { r -> r.status == 'fail' }"
  target {
    type = "DataQualityScan"
    ids  = [%[2]s]
  }
  notification {
    channels = ["Mail", "Sms"]
    receivers {
      receiver_type   = "DingdingUrl"
      receiver_values = ["%[3]s"]
      extension       = "{\"atAll\":true}"
    }
  }
}
`, projectId, scanId, webhookURL)
}
