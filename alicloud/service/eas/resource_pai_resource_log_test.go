package eas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiResourceLog_basic exercises the enable/read/disable
// lifecycle of alicloud_pai_resource_log.
//
// The log delivery config is a singleton on a parent EAS resource group, and
// CreateResourceLog requires a pre-existing SLS project and log store that the
// API will not create. Provisioning an EAS resource group with real compute
// resources is out of scope for the acceptance test account (billing/quota), so
// the parent cluster_id, resource_id and the SLS project_name/log_store are
// supplied through environment variables. The test is skipped when any of them
// is unset, which keeps the suite green in environments without a PAI-EAS
// resource group while still validating the full lifecycle where one exists.
func TestAccAliCloudPaiResourceLog_basic(t *testing.T) {
	clusterId := os.Getenv("ALICLOUD_PAI_RESOURCE_LOG_CLUSTER_ID")
	resourceId := os.Getenv("ALICLOUD_PAI_RESOURCE_LOG_RESOURCE_ID")
	projectName := os.Getenv("ALICLOUD_PAI_RESOURCE_LOG_PROJECT_NAME")
	logStore := os.Getenv("ALICLOUD_PAI_RESOURCE_LOG_LOG_STORE")
	if clusterId == "" || resourceId == "" || projectName == "" || logStore == "" {
		t.Skip("set ALICLOUD_PAI_RESOURCE_LOG_{CLUSTER_ID,RESOURCE_ID,PROJECT_NAME,LOG_STORE} to a pre-existing PAI-EAS resource group and SLS project/logstore to run this acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiResourceLogConfig(clusterId, resourceId, projectName, logStore),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_resource_log.default", "cluster_id", clusterId),
					resource.TestCheckResourceAttr("alicloud_pai_resource_log.default", "resource_id", resourceId),
					resource.TestCheckResourceAttr("alicloud_pai_resource_log.default", "project_name", projectName),
					resource.TestCheckResourceAttr("alicloud_pai_resource_log.default", "log_store", logStore),
					resource.TestCheckResourceAttrSet("alicloud_pai_resource_log.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_resource_log.default", "status"),
				),
			},
			{
				ResourceName:      "alicloud_pai_resource_log.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudPaiResourceLogConfig(clusterId, resourceId, projectName, logStore string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_resource_log" "default" {
  cluster_id   = "%[1]s"
  resource_id  = "%[2]s"
  project_name = "%[3]s"
  log_store    = "%[4]s"
}
`, clusterId, resourceId, projectName, logStore)
}
