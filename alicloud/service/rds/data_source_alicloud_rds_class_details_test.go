package rds_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRdsClassDetailsDataSource(t *testing.T) {
	resourceId := "data.alicloud_rds_class_details.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRdsClassDetailsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":                       acctest.CHECKSET,
						"commodity_code":           "bards",
						"class_code":               "mysql.n4.medium.2c",
						"engine_version":           "8.0",
						"engine":                   "MySQL",
						"max_iombps":               acctest.CHECKSET,
						"max_connections":          acctest.CHECKSET,
						"class_group":              acctest.CHECKSET,
						"cpu":                      acctest.CHECKSET,
						"memory_class":             acctest.CHECKSET,
						"max_iops":                 acctest.CHECKSET,
						"reference_price":          acctest.CHECKSET,
						"category":                 acctest.CHECKSET,
						"db_instance_storage_type": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccCheckAlicloudRdsClassDetailsDataSource2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":                       acctest.CHECKSET,
						"commodity_code":           "rds",
						"class_code":               "rds.mssql.s2.large",
						"engine_version":           "2017_ent",
						"engine":                   "SQLServer",
						"max_iombps":               acctest.CHECKSET,
						"max_connections":          acctest.CHECKSET,
						"class_group":              acctest.CHECKSET,
						"cpu":                      acctest.CHECKSET,
						"memory_class":             acctest.CHECKSET,
						"max_iops":                 acctest.CHECKSET,
						"reference_price":          acctest.CHECKSET,
						"category":                 acctest.CHECKSET,
						"db_instance_storage_type": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudRdsClassDetailsDataSource = `
data "alicloud_rds_class_details" "default" {
  commodity_code = "bards"
  class_code     = "mysql.n4.medium.2c"
  engine_version = "8.0"
  engine   		 = "MySQL"
}
`

const testAccCheckAlicloudRdsClassDetailsDataSource2 = `
data "alicloud_rds_class_details" "default" {
  commodity_code = "rds"
  class_code     = "rds.mssql.s2.large"
  engine_version = "2017_ent"
  engine         = "SQLServer"
}
`
