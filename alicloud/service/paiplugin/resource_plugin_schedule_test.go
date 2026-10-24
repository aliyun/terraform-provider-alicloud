package paiplugin_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Creating a PAI Plugin Schedule requires pre-existing fixtures that the test
// account cannot provision on its own: a crowd (group_id), and optionally an
// approved SMS signature/template. The schedule also triggers real outreach, so
// the config schedules execution far in the future and deletes before it fires,
// avoiding any actual send. The test is skipped unless ALICLOUD_PAI_PLUGIN_GROUP_ID
// is provided, which is the honest "not run" state when the fixture is absent.
func TestAccAliCloudPaiPluginSchedule(t *testing.T) {
	if v := os.Getenv("ALICLOUD_PAI_PLUGIN_GROUP_ID"); v == "" {
		t.Skip("ALICLOUD_PAI_PLUGIN_GROUP_ID must be set to a real crowd ID to run this acceptance test")
	}

	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiPluginScheduleConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_schedule.default", "id"),
					resource.TestCheckResourceAttr("alicloud_pai_plugin_schedule.default", "name", fmt.Sprintf("tf-testacc%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_plugin_schedule.default", "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_schedule.default", "group_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_schedule.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_schedule.default", "create_time"),
				),
			},
			{
				ResourceName:      "alicloud_pai_plugin_schedule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// AISendStartDate / AISendEndDate are write-only: accepted by
				// CreateSchedule but never returned by GetSchedule, so they cannot
				// be populated on import. This is the verified write-only ForceNew
				// exception, not a ForceNew field that the API echoes back.
				ImportStateVerifyIgnore: []string{"ai_send_start_date", "ai_send_end_date"},
			},
		},
	})
}

func testAccAliCloudPaiPluginScheduleConfig(r int) string {
	groupId := os.Getenv("ALICLOUD_PAI_PLUGIN_GROUP_ID")
	signatureId := os.Getenv("ALICLOUD_PAI_PLUGIN_SIGNATURE_ID")
	templateId := os.Getenv("ALICLOUD_PAI_PLUGIN_TEMPLATE_ID")
	templateCode := os.Getenv("ALICLOUD_PAI_PLUGIN_TEMPLATE_CODE")
	signName := os.Getenv("ALICLOUD_PAI_PLUGIN_SIGN_NAME")
	if signName == "" {
		signName = "阿里云"
	}

	return fmt.Sprintf(`
resource "alicloud_pai_plugin_schedule" "default" {
  name              = "tf-testacc%[1]d"
  group_id          = "%[2]s"
  signature_id      = "%[3]s"
  template_id       = "%[4]s"
  template_code     = "%[5]s"
  sign_name         = "%[6]s"
  payment_type      = "PayAsYouGo"
  execute_time      = "2099-01-01 12:00:00"
  end_time          = "2099-12-31"
  ai_send_start_date = "2099-01-01 12:00:00"
  ai_send_end_date   = "2099-02-01 12:00:00"
  repeat_cycle      = 1
  repeat_cycle_unit = 0
  repeat_times      = 1
}
`, r, groupId, signatureId, templateId, templateCode, signName)
}
