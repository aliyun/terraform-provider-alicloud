package mnsopen_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	mnsopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/mnsopen"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// MessageService AccountLogging is an account-level singleton (SetAccountAttributes / GetAccountAttributes);
// destroy is a no-op, so CheckDestroy cannot be asserted. When log_enabled is false the API returns no
// project_name / log_store_name, so those are cleared together with disabling logging.
// lintignore: AT001
func TestAccAliCloudMessageServiceAccountLogging_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_account_logging.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudMessageServiceAccountLoggingMap0)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceAccountLogging")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-example-mns-log-%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudMessageServiceAccountLoggingBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		Providers: tfacctest.SDKProviders,
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
						"project_name":          tfacctest.CHECKSET,
						"log_store_name":        tfacctest.CHECKSET,
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
						"project_name":          tfacctest.CHECKSET,
						"log_store_name":        tfacctest.CHECKSET,
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
