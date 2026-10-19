package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDataWorksConnectionsDataSource(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATAWORKS_PROJECT_ID"); v == "" {
		t.Skip("Skipping: ALICLOUD_DATAWORKS_PROJECT_ID env var is required to run this test.")
	}
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksconnection%d", defaultRegionToTest, rand)
	resourceId := "data.alicloud_data_works_connections.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksConnectionsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id":      "${alicloud_data_works_connection.default.project_id}",
			"connection_type": "mysql",
			"name":            name,
			"ids":             []string{"${alicloud_data_works_connection.default.connection_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id":      "${alicloud_data_works_connection.default.project_id}",
			"connection_type": "mysql",
			"name":            fmt.Sprintf("%s_fake", name),
			"ids":             []string{"${alicloud_data_works_connection.default.connection_id}_fake"},
		}),
	}

	var existDataWorksConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"connections.#":                 "1",
			"connections.0.project_id":      CHECKSET,
			"connections.0.connection_id":   CHECKSET,
			"connections.0.connection_name": name,
			"connections.0.connection_type": "mysql",
		}
	}

	var fakeDataWorksConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"connections.#": "0",
			"ids.#":         "0",
		}
	}

	var DataWorksConnectionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksConnectionsMapFunc,
		fakeMapFunc:  fakeDataWorksConnectionsMapFunc,
	}

	DataWorksConnectionsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataWorksConnectionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "project_id" {
  default = "%s"
}

resource "alicloud_data_works_connection" "default" {
  project_id      = "${var.project_id}"
  connection_name = var.name
  connection_type = "mysql"
  sub_type        = "mysql"
  env_type        = 0
  content         = "{\"database\":\"tf_example\",\"host\":\"127.0.0.1\",\"password\":\"tf_example_pw\",\"port\":\"3306\",\"username\":\"tf_example_user\"}"
  description     = "tf-example-connection-desc-1"
}

data "alicloud_data_works_connections" "default" {
  project_id      = alicloud_data_works_connection.default.project_id
  connection_type = "mysql"
  name            = var.name
  ids             = [alicloud_data_works_connection.default.connection_id]
}
`, name, os.Getenv("ALICLOUD_DATAWORKS_PROJECT_ID"))
}
