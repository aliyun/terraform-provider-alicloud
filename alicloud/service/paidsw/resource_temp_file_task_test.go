package paidsw_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The temp file task requires an existing DSW instance (instance_id). The
// provider does not yet manage a DSW instance resource, so the ACC test
// references a pre-existing instance supplied through ALICLOUD_PAIDSW_INSTANCE_ID.
// The test is skipped when the variable is absent rather than fabricating a
// pass; QA validates against an account with a real DSW instance.
func testAccPreCheckPaiDswTempFileTask(t *testing.T) {
	acctest.PreCheck(t)
	if os.Getenv("ALICLOUD_PAIDSW_INSTANCE_ID") == "" {
		t.Skip("env ALICLOUD_PAIDSW_INSTANCE_ID is not set; skipping PAI DSW temp file task ACC test")
	}
}

func TestAccAliCloudPaiDswTempFileTask(t *testing.T) {
	instanceId := os.Getenv("ALICLOUD_PAIDSW_INSTANCE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPaiDswTempFileTask(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiDswTempFileTaskConfig(instanceId, "2099-12-31 23:59:59"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_temp_file_task.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_temp_file_task.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_temp_file_task.default", "gmt_expired_time", "2099-12-31 23:59:59"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_temp_file_task.default", "region_id"),
				),
			},
			{
				Config: testAccAliCloudPaiDswTempFileTaskConfig(instanceId, "2099-06-30 00:00:00"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_dsw_temp_file_task.default", "gmt_expired_time", "2099-06-30 00:00:00"),
				),
			},
			{
				ResourceName:      "alicloud_pai_dsw_temp_file_task.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudPaiDswTempFileTaskConfig(instanceId, gmtExpiredTime string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_temp_file_task" "default" {
  instance_id     = "%[1]s"
  gmt_expired_time = "%[2]s"
}
`, instanceId, gmtExpiredTime)
}
