package eas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiLogDownloadTask_basic covers the Create/Read/Delete lifecycle
// of the EAS log download task. The log download task API
// (/api/v2/log-download-tasks) is a regional, public-cluster EAS v2 endpoint whose
// ClusterId is the provider region, so the task references a public EAS service by
// name. The fixture therefore builds a public echo_processor EAS service (mirroring
// the proven alicloud_pai_service basic3213 pattern: develop=false, no workspace_id,
// cn-hangzhou) and pins the provider to cn-hangzhou where EAS is available. The task
// is first created with the minimal required argument, then recreated with the
// optional instance_name/container filters and a distinct alias_name, and finally
// re-imported. from_time/to_time are omitted in both steps because the EAS Create
// API enforces a server-side seven-day window on the log collection range and
// applies sensible defaults (FromTime=ToTime-7d, ToTime=now) when they are absent.
// All settable attributes are ForceNew, so the optional-field step exercises the
// recreate path.
func TestAccAliCloudPaiLogDownloadTask_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiLogDownloadTaskBasic(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "task_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "service_name"),
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "region_id"),
					resource.TestCheckResourceAttr("alicloud_pai_log_download_task.default", "alias_name", fmt.Sprintf("tftestaccldtalias%d", r)),
				),
			},
			{
				Config: testAccAliCloudPaiLogDownloadTaskFull(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_log_download_task.default", "id"),
					resource.TestCheckResourceAttr("alicloud_pai_log_download_task.default", "alias_name", fmt.Sprintf("tftestaccldtaliasfull%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_log_download_task.default", "instance_name", fmt.Sprintf("tftestaccldtinst%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_log_download_task.default", "container", "main"),
				),
			},
			{
				ResourceName:      "alicloud_pai_log_download_task.default",
				ImportState:       true,
				ImportStateVerify: true,
				// The EAS log download task completes asynchronously (Running ->
				// Succeeded), so status/download_url/finish_time drift between the
				// recreate step and the import re-read; from_time/to_time are Computed
				// API defaults that the import does not reproduce exactly.
				ImportStateVerifyIgnore: []string{"from_time", "to_time", "status", "download_url", "finish_time"},
			},
		},
	})
}

// testAccAliCloudPaiLogDownloadTaskBase builds the public EAS service the log
// download task depends on: an echo_processor service pinned to cn-hangzhou with
// develop=false and no workspace_id, matching the alicloud_pai_service basic3213
// fixture. The service name uses no separator so EAS accepts it.
func testAccAliCloudPaiLogDownloadTaskBase(r int) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_service" "default" {
  develop        = "false"
  service_config = jsonencode({
    metadata = {
      cpu      = 1
      gpu      = 0
      instance = 1
      memory   = 2000
      name     = "tftestaccldt%[1]d"
      rpc      = { keepalive = 70000 }
    }
    model_path      = "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz"
    processor_entry = "libecho.so"
    processor_path  = "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz"
    processor_type  = "cpp"
  })
}
`, r)
}

func testAccAliCloudPaiLogDownloadTaskBasic(r int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_log_download_task" "default" {
  service_name = alicloud_pai_service.default.id
  alias_name   = "tftestaccldtalias%[2]d"
}
`, testAccAliCloudPaiLogDownloadTaskBase(r), r)
}

// testAccAliCloudPaiLogDownloadTaskFull recreates the task with the optional
// instance_name and container filters set, plus a distinct alias_name. from_time
// and to_time are intentionally omitted: the EAS Create API applies a server-side
// seven-day window on the log collection range (FromTime defaults to ToTime-7d,
// ToTime to now) and rejects explicit ranges that fall outside it, so exercising
// the API defaults keeps the fixture stable across runs while still driving the
// ForceNew recreate path for alias_name/instance_name/container.
func testAccAliCloudPaiLogDownloadTaskFull(r int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_log_download_task" "default" {
  service_name  = alicloud_pai_service.default.id
  alias_name    = "tftestaccldtaliasfull%[2]d"
  instance_name = "tftestaccldtinst%[2]d"
  container     = "main"
}
`, testAccAliCloudPaiLogDownloadTaskBase(r), r)
}
