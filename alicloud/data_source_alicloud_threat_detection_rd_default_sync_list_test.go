package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test ThreatDetection RdDefaultSyncList. >>> Data source test cases.
func TestAccAliCloudThreatDetectionRdDefaultSyncListDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionrddefaultsynclist%d", defaultRegionToTest, rand)
	dataSourceId := "data.alicloud_threat_detection_rd_default_sync_list.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Set up the managed list first, so that the data source reads a
				// deterministic, non-empty folder list in the next step.
				Config: testAccThreatDetectionRdDefaultSyncListDataSourceConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_rd_default_sync_list.default", "folder_ids.#", "1"),
				),
			},
			{
				// The data source has no input arguments: it reads the current
				// account's default synchronization list. The managed resource
				// already exists from the previous step, so the plan-time read
				// returns the same folder list (no depends_on needed).
				Config: testAccThreatDetectionRdDefaultSyncListDataSourceConfig(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "folder_ids.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceId, "folder_ids.0", "alicloud_resource_manager_folder.default", "id"),
				),
			},
		},
	})
}

func testAccThreatDetectionRdDefaultSyncListDataSourceConfig(name string, withDataSource bool) string {
	config := fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_resource_manager_folder" "default" {
  folder_name = "tf-testacc-rd-sync-folder-${var.name}"
}

resource "alicloud_threat_detection_rd_default_sync_list" "default" {
  folder_ids = [alicloud_resource_manager_folder.default.id]
}
`, name)
	if withDataSource {
		config += `
data "alicloud_threat_detection_rd_default_sync_list" "default" {
}
`
	}
	return config
}

// Test ThreatDetection RdDefaultSyncList. <<< Data source test cases.
