package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsClassDetailsDataSource(t *testing.T) {
	resourceId := "data.alicloud_rds_class_details.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRdsClassDetailsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":                       CHECKSET,
						"commodity_code":           "bards",
						"class_code":               "mysql.n4.medium.2c",
						"engine_version":           "8.0",
						"engine":                   "MySQL",
						"max_iombps":               CHECKSET,
						"max_connections":          CHECKSET,
						"class_group":              CHECKSET,
						"cpu":                      CHECKSET,
						"memory_class":             CHECKSET,
						"max_iops":                 CHECKSET,
						"reference_price":          CHECKSET,
						"category":                 CHECKSET,
						"db_instance_storage_type": CHECKSET,
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
