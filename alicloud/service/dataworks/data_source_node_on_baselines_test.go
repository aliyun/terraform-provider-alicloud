package dataworks_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudDataWorksNodeOnBaselinesDataSource exercises the plural
// framework-native data source backed by the dataworks-public 2020-05-18
// ListNodesByBaseline action. The test account must have a pre-created
// DataWorks baseline and authorization on ListNodesByBaseline; otherwise the
// test skips cleanly via ALICLOUD_DATA_WORKS_BASELINE_ID rather than failing
// as a code defect. The second step narrows the plural data source via an ids
// filter that references the first node returned by the baseline lookup, using
// filter-parameter interpolation as the implicit dependency mechanism
// (depends_on on data blocks is not supported).
func TestAccAlicloudDataWorksNodeOnBaselinesDataSource(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID"); v == "" {
		t.Skip("Skipping because ALICLOUD_DATA_WORKS_BASELINE_ID is not set; the test account must have a pre-created DataWorks baseline and authorization on ListNodesByBaseline.")
	}
	baselineID := os.Getenv("ALICLOUD_DATA_WORKS_BASELINE_ID")
	resourceId := "data.alicloud_data_works_node_on_baselines.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksNodeOnBaselinesDataSourceConfig(baselineID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "baseline_id"),
					resource.TestCheckResourceAttrSet(resourceId, "nodes.#"),
				),
			},
			{
				Config: testAccCheckAlicloudDataWorksNodeOnBaselinesDataSourceConfigWithIDFilter(baselineID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "baseline_id"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDataWorksNodeOnBaselinesDataSourceConfig(baselineID string) string {
	return `
variable "baseline_id" {
  default = "` + baselineID + `"
}

data "alicloud_data_works_node_on_baselines" "default" {
  baseline_id = var.baseline_id
}

output "nodes_count" {
  value = data.alicloud_data_works_node_on_baselines.default.nodes.#
}
`
}

func testAccCheckAlicloudDataWorksNodeOnBaselinesDataSourceConfigWithIDFilter(baselineID string) string {
	return `
variable "baseline_id" {
  default = "` + baselineID + `"
}

data "alicloud_data_works_node_on_baselines" "all" {
  baseline_id = var.baseline_id
}

data "alicloud_data_works_node_on_baselines" "default" {
  baseline_id = var.baseline_id
  ids         = [data.alicloud_data_works_node_on_baselines.all.nodes.0.id]
}
`
}
