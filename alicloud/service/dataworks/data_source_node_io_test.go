package dataworks_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudDataWorksNodeIoDataSource exercises the Framework-native
// alicloud_data_works_node_io data source against a real DataWorks scheduling node.
//
// ListNodeIO is a read-only query over an existing node, so the test runs only when
// a real node fixture is supplied via ALICLOUD_DATA_WORKS_NODE_ID. The fixture is a
// pre-existing, read-only node supplied by the test account/maintainer; this test
// does not create or own a long-lived DataWorks node. When the env var is absent the
// case is skipped rather than reported as PASS.
func TestAccAliCloudDataWorksNodeIoDataSource(t *testing.T) {
	acctest.PreCheck(t)

	nodeId := os.Getenv("ALICLOUD_DATA_WORKS_NODE_ID")
	if nodeId == "" {
		t.Skip("env ALICLOUD_DATA_WORKS_NODE_ID (a real DataWorks node id) is required for this test")
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
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksNodeIoDataSourceConfig(nodeId, projectEnv, ioType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_node_io.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "node_id", nodeId),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "project_env", projectEnv),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "io_type", ioType),
					resource.TestCheckResourceAttr("data.alicloud_data_works_node_io.default", "output_file", "/tmp/node_io_result.txt"),
				),
			},
		},
	})
}

func testAccAliCloudDataWorksNodeIoDataSourceConfig(nodeId, projectEnv, ioType string) string {
	return fmt.Sprintf(`
data "alicloud_data_works_node_io" "default" {
  node_id     = "%[1]s"
  project_env = "%[2]s"
  io_type     = "%[3]s"
  ids         = []
  output_file = "/tmp/node_io_result.txt"
}
`, nodeId, projectEnv, ioType)
}
