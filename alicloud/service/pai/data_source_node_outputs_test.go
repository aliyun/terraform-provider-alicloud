package pai_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ListNodeOutputs enumerates the outputs of a node inside a PAI Studio
// experiment. Such experiments are not manageable through a Terraform resource,
// so this test reads a pre-existing node pointed at by environment variables and
// skips when they are not provided.
func testAccPaiNodeOutputsPreCheck(t *testing.T) (experimentId, nodeId string) {
	experimentId = os.Getenv("ALICLOUD_PAI_NODE_OUTPUTS_EXPERIMENT_ID")
	nodeId = os.Getenv("ALICLOUD_PAI_NODE_OUTPUTS_NODE_ID")
	if experimentId == "" || nodeId == "" {
		t.Skip("env ALICLOUD_PAI_NODE_OUTPUTS_EXPERIMENT_ID/NODE_ID must point at a pre-existing PAI Studio experiment and node")
	}
	return
}

func TestAccAliCloudPaiNodeOutputsDataSource(t *testing.T) {
	experimentId, nodeId := testAccPaiNodeOutputsPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiNodeOutputsDataSourceConfig(experimentId, nodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_pai_node_outputs.default", "id", fmt.Sprintf("%s:%s", experimentId, nodeId)),
					resource.TestCheckResourceAttr("data.alicloud_pai_node_outputs.default", "experiment_id", experimentId),
					resource.TestCheckResourceAttr("data.alicloud_pai_node_outputs.default", "node_id", nodeId),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_node_outputs.default", "outputs.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_node_outputs.default", "ids.#"),
				),
			},
		},
	})
}

func testAccPaiNodeOutputsDataSourceConfig(experimentId, nodeId string) string {
	return fmt.Sprintf(`
data "alicloud_pai_node_outputs" "default" {
  experiment_id = "%s"
  node_id       = "%s"
}
`, experimentId, nodeId)
}
