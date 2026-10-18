package dataworks_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudDataWorksTableLineageDataSource verifies the
// alicloud_data_works_table_lineage data source against a real table that has
// lineage relationships in the test account.
//
// The data source queries pre-existing table lineage that cannot be
// fixture-created inside an ACC run, so the test is environment-guarded: it
// runs only when ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME (and optionally
// ALICLOUD_DATA_WORKS_LINEAGE_DATABASE_NAME) are set, and skips otherwise. The
// account must also hold DataWorks Data Map authorization for
// GetMetaTableLineage; without it the API returns ResourceNotAuthorized.Api.
// This mirrors how other pure-query data sources handle pre-existing
// infrastructure.
func TestAccAlicloudDataWorksTableLineageDataSource(t *testing.T) {
	tableName := os.Getenv("ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME")
	databaseName := os.Getenv("ALICLOUD_DATA_WORKS_LINEAGE_DATABASE_NAME")
	if tableName == "" {
		t.Skip("Skipping test: ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME is not set; this data source queries pre-existing table lineage and cannot fixture its own test data.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, "up"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_table_lineage.up", "id"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_table_lineage.up", "direction", "up"),
				),
			},
			{
				Config: testAccAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, "down"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_table_lineage.down", "id"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_table_lineage.down", "direction", "down"),
				),
			},
		},
	})
}

func testAccAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, direction string) string {
	config := `
data "alicloud_data_works_table_lineage" "` + direction + `" {
  direction  = "` + direction + `"
  table_name = "` + tableName + `"
`
	if databaseName != "" {
		config += `  database_name = "` + databaseName + `"
`
	}
	config += "}\n"
	return config
}
