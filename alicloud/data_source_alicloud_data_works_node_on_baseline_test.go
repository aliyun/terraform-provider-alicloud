package alicloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksNodeOnBaselineDataSource(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID"); v == "" {
		t.Skip("Skipping because ALICLOUD_DATA_WORKS_BASELINE_ID is not set; the test account must have a pre-created DataWorks baseline and authorization on GetNodeOnBaseline.")
	}
	baselineID := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID")
	resourceId := "data.alicloud_data_works_node_on_baseline.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksNodeOnBaselineDataSourceConfig(baselineID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "baseline_id"),
					resource.TestCheckResourceAttrSet(resourceId, "node_id"),
					resource.TestCheckResourceAttrSet(resourceId, "node_name"),
					resource.TestCheckResourceAttrSet(resourceId, "owner"),
					resource.TestCheckResourceAttrSet(resourceId, "project_id"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
				),
			},
		},
	})
}

// testAccCheckAlicloudDataWorksNodeOnBaselineDataSourceConfig builds a singular
// data source config that uses filter-parameter interpolation for implicit
// dependencies. Per the SDK v1.17.2 / 0.12.7-sdk semantics, data blocks must
// not declare depends_on; the baseline_id is sourced from an env-var-driven
// local so the test stays runnable only when a pre-created baseline exists.
func testAccCheckAlicloudDataWorksNodeOnBaselineDataSourceConfig(baselineID string) string {
	return `
variable "baseline_id" {
  default = "` + baselineID + `"
}

data "alicloud_data_works_node_on_baseline" "default" {
  baseline_id = var.baseline_id
}

output "node_id" {
  value = data.alicloud_data_works_node_on_baseline.default.node_id
}
`
}
