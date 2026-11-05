package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudThreatDetectionClientFileProtectEvents_basic queries the data
// source with a non-existent rule name. File protect client events are
// system-generated, so a freshly created rule produces no events; the test
// therefore validates the data source registration, the ListFileProtectClientEvent
// call, the response parsing and the rule_name filter by asserting an empty
// events list and a computed id.
func TestAccAliCloudThreatDetectionClientFileProtectEvents_basic(t *testing.T) {
	rand := tFacctest.RandIntRange(10000, 99999)
	ruleName := fmt.Sprintf("tf-testacc-cfp-events-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientFileProtectEventsBasicConfig(ruleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_client_file_protect_events.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_client_file_protect_events.default", "events.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_client_file_protect_events.default", "ids.#", "0"),
				),
			},
		},
	})
}

// TestAccAliCloudThreatDetectionClientFileProtectEvents_alertLevels queries the
// data source with the alert_levels repeatList filter alongside rule_name to
// validate the repeatList request encoding for ListFileProtectClientEvent.
func TestAccAliCloudThreatDetectionClientFileProtectEvents_alertLevels(t *testing.T) {
	rand := tFacctest.RandIntRange(10000, 99999)
	ruleName := fmt.Sprintf("tf-testacc-cfp-events-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientFileProtectEventsAlertLevelsConfig(ruleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_client_file_protect_events.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_client_file_protect_events.default", "events.#", "0"),
				),
			},
		},
	})
}

// TestAccAliCloudThreatDetectionClientFileProtectEvents_idsFilter queries the
// data source with the ids provider-side filter; a non-matching id yields an
// empty events list even when the rule_name filter is absent.
func TestAccAliCloudThreatDetectionClientFileProtectEvents_idsFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientFileProtectEventsIdsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_client_file_protect_events.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_client_file_protect_events.default", "events.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_client_file_protect_events.default", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccClientFileProtectEventsBasicConfig(ruleName string) string {
	return fmt.Sprintf(`
data "alicloud_threat_detection_client_file_protect_events" "default" {
  rule_name = "%[1]s"
}
`, ruleName)
}

func testAccClientFileProtectEventsAlertLevelsConfig(ruleName string) string {
	return fmt.Sprintf(`
data "alicloud_threat_detection_client_file_protect_events" "default" {
  rule_name    = "%[1]s"
  alert_levels = [2, 3]
}
`, ruleName)
}

func testAccClientFileProtectEventsIdsConfig() string {
	return `
data "alicloud_threat_detection_client_file_protect_events" "default" {
  ids = ["999999999999"]
}
`
}
