package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksNodeIoDataSource(t *testing.T) {
	testAccPreCheck(t)
	nodeId := os.Getenv("ALICLOUD_DATA_WORKS_NODE_ID")
	if nodeId == "" {
		t.Skip("env ALICLOUD_DATA_WORKS_NODE_ID is required for this test")
	}
	ioType := os.Getenv("ALICLOUD_DATA_WORKS_IO_TYPE")
	if ioType == "" {
		ioType = "input"
	}
	projectEnv := os.Getenv("ALICLOUD_DATA_WORKS_PROJECT_ENV")
	if projectEnv == "" {
		projectEnv = "PROD"
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "alicloud_data_works_node_io" "default" {
  node_id     = "%s"
  project_env = "%s"
  io_type     = "%s"
  ids         = []
  output_file = "/tmp/node_io_result.txt"
}
`, nodeId, projectEnv, ioType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_node_io.default"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "node_id", nodeId),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "project_env", projectEnv),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "io_type", ioType),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "output_file", "/tmp/node_io_result.txt"),
				),
			},
		},
	})
}
