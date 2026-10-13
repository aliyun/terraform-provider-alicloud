package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// MessageService AccountLogging is an account-level singleton (SetAccountAttributes / GetAccountAttributes);
// destroy is a no-op, so CheckDestroy cannot be asserted. When log_enabled is false the API returns no
// project_name / log_store_name, so those are cleared together with disabling logging.
// lintignore: AT001
func TestAccAliCloudMessageServiceAccountLogging_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_account_logging.default"
	ra := resourceAttrInit(resourceId, AlicloudMessageServiceAccountLoggingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceAccountLogging")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-example-mns-log-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMessageServiceAccountLoggingBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"log_enabled":           "true",
					"message_trace_enabled": "false",
					"project_name":          "${alicloud_log_project.default.project_name}",
					"log_store_name":        "${alicloud_log_store.default.logstore_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_enabled":           "true",
						"message_trace_enabled": "false",
						"project_name":          CHECKSET,
						"log_store_name":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"message_trace_enabled": "true",
					"project_name":          "${alicloud_log_project.update.project_name}",
					"log_store_name":        "${alicloud_log_store.update.logstore_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message_trace_enabled": "true",
						"project_name":          CHECKSET,
						"log_store_name":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_enabled":    "false",
					"project_name":   "",
					"log_store_name": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_enabled":    "false",
						"project_name":   "",
						"log_store_name": "",
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

var AlicloudMessageServiceAccountLoggingMap0 = map[string]string{}

func AlicloudMessageServiceAccountLoggingBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "tf-testacc message service account logging"
}

resource "alicloud_log_store" "default" {
  project_name     = alicloud_log_project.default.project_name
  logstore_name    = var.name
  retention_period = 30
  shard_count      = 2
}

resource "alicloud_log_project" "update" {
  project_name = "${var.name}-update"
  description  = "tf-testacc message service account logging update"
}

resource "alicloud_log_store" "update" {
  project_name     = alicloud_log_project.update.project_name
  logstore_name    = "${var.name}-update"
  retention_period = 30
  shard_count      = 2
}
`, name)
}
