package dataworks_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudDataWorksNodeOnBaselineDataSource exercises the singular
// framework-native data source backed by the dataworks-public 2020-05-18
// GetNodeOnBaseline action. The test account must have a pre-created DataWorks
// baseline and authorization on GetNodeOnBaseline; otherwise the test skips
// cleanly via ALICLOUD_DATA_WORKS_BASELINE_ID rather than failing as a code
// defect. No depends_on is used on the data block; the implicit dependency is
// established via filter-parameter interpolation.
func TestAccAlicloudDataWorksNodeOnBaselineDataSource(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID"); v == "" {
		t.Skip("Skipping because ALICLOUD_DATA_WORKS_BASELINE_ID is not set; the test account must have a pre-created DataWorks baseline and authorization on GetNodeOnBaseline.")
	}
	baselineID := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID")
	resourceId := "data.alicloud_data_works_node_on_baseline.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksNodeOnBaselineDataSourceConfig(baselineID),
				Check: resource.ComposeTestCheckFunc(
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
