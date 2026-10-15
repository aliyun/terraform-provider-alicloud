package pai_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// PAI Studio experiments (the PaiStudio/2021-02-02 visual-modeling product) are
// not manageable through any Terraform resource in this provider, so the test
// cannot create its fixture inline. Point it at a pre-existing experiment via
// these environment variables; when they are absent the test is skipped rather
// than fabricated against a non-existent object.
func testAccPaiNodeOutputPreCheck(t *testing.T) (experimentId, nodeId, outputId string) {
	experimentId = os.Getenv("ALICLOUD_PAI_NODE_OUTPUT_EXPERIMENT_ID")
	nodeId = os.Getenv("ALICLOUD_PAI_NODE_OUTPUT_NODE_ID")
	outputId = os.Getenv("ALICLOUD_PAI_NODE_OUTPUT_OUTPUT_ID")
	if experimentId == "" || nodeId == "" || outputId == "" {
		t.Skip("env ALICLOUD_PAI_NODE_OUTPUT_EXPERIMENT_ID/NODE_ID/OUTPUT_ID must point at a pre-existing PAI Studio experiment, node and output")
	}
	return
}

func TestAccAliCloudPaiNodeOutputDataSource(t *testing.T) {
	experimentId, nodeId, outputId := testAccPaiNodeOutputPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiNodeOutputDataSourceConfig(experimentId, nodeId, outputId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_pai_node_output.default", "id", fmt.Sprintf("%s:%s:%s", experimentId, nodeId, outputId)),
					resource.TestCheckResourceAttr("data.alicloud_pai_node_output.default", "experiment_id", experimentId),
					resource.TestCheckResourceAttr("data.alicloud_pai_node_output.default", "node_id", nodeId),
					resource.TestCheckResourceAttr("data.alicloud_pai_node_output.default", "output_id", outputId),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_node_output.default", "type"),
				),
			},
		},
	})
}

func testAccPaiNodeOutputDataSourceConfig(experimentId, nodeId, outputId string) string {
	return fmt.Sprintf(`
data "alicloud_pai_node_output" "default" {
  experiment_id = "%s"
  node_id       = "%s"
  output_id     = "%s"
}
`, experimentId, nodeId, outputId)
}
