package polardb_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccNodeClassesRegion = "cn-hangzhou"

func TestAccAlicloudPolarDBNodeClasses(t *testing.T) {
	resourceId := "data.alicloud_polardb_node_classes.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// PolarDB is Hangzhou-only (connectivity.PolarDBSupportRegions); the
			// flat suite pinned the region through checkoutSupportedRegions, so
			// pin it here.
			os.Setenv("ALICLOUD_REGION", testAccNodeClassesRegion)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// pay_type filter, prepaid
			{
				Config: testAccNodeClassesConfig(`pay_type = "PrePaid"`),
				Check:  checkNodeClassesExist(),
			},
			// pay_type filter, postpaid
			{
				Config: testAccNodeClassesConfig(`pay_type = "PostPaid"`),
				Check:  checkNodeClassesExist(),
			},
			// engine and version filter, MySQL 5.6
			{
				Config: testAccNodeClassesConfig(`
  pay_type   = "PostPaid"
  db_type    = "Mysql"
  db_version = "5.6"`),
				Check: checkNodeClassesExist(),
			},
			// engine and version filter, PostgreSQL 11
			{
				Config: testAccNodeClassesConfig(`
  pay_type   = "PostPaid"
  db_type    = "PostgreSQL"
  db_version = "11"`),
				Check: checkNodeClassesExist(),
			},
			// engine and version filter, PostgreSQL 14. The API returns PostgreSQL 14
			// as the compact engine identifier "pg14", which the previous substring
			// matching against "PostgreSQL" filtered out.
			{
				Config: testAccNodeClassesConfig(`
  pay_type   = "PostPaid"
  db_type    = "PostgreSQL"
  db_version = "14"`),
				Check: checkNodeClassesExist(),
			},
			// engine and version filter, Oracle 14. The API returns Oracle 14 as the
			// compact engine identifier "pgo14"; normalizeEngineType must map the
			// "pgo" prefix to "oracle" for the db_type filter to keep the resources.
			{
				Config: testAccNodeClassesConfig(`
  pay_type   = "PostPaid"
  db_type    = "Oracle"
  db_version = "14"`),
				Check: resource.ComposeTestCheckFunc(
					checkNodeClassesExist(),
					resource.TestCheckResourceAttr(resourceId, "classes.0.supported_engines.0.engine", "pgo14"),
				),
			},
			// db_node_class filter, matching
			{
				Config: testAccNodeClassesConfig(`
  pay_type      = "PostPaid"
  db_node_class = "polar.pg.x4.large"`),
				Check: checkNodeClassesExist(),
			},
			// db_node_class filter, no match
			{
				Config: testAccNodeClassesConfig(`
  pay_type      = "PostPaid"
  db_node_class = "fake"`),
				Check: resource.TestCheckResourceAttr(resourceId, "classes.#", "0"),
			},
			// engine, version and category filter, MySQL 8.0
			{
				Config: testAccNodeClassesConfig(`
  pay_type   = "PostPaid"
  db_type    = "Mysql"
  db_version = "8.0"
  category   = "Normal"`),
				Check: checkNodeClassesExist(),
			},
			// region_id filter
			{
				Config: testAccNodeClassesConfig(fmt.Sprintf(`
  pay_type  = "PostPaid"
  region_id = "%s"`, testAccNodeClassesRegion)),
				Check: checkNodeClassesExist(),
			},
		},
	})
}

func checkNodeClassesExist() resource.TestCheckFunc {
	resourceId := "data.alicloud_polardb_node_classes.default"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceId, "classes.#"),
		resource.TestCheckResourceAttrSet(resourceId, "classes.0.supported_engines.0.available_resources.0.db_node_class"),
		resource.TestCheckResourceAttrSet(resourceId, "classes.0.supported_engines.0.engine"),
		resource.TestCheckResourceAttrSet(resourceId, "classes.0.zone_id"),
	)
}

func testAccNodeClassesConfig(attrs string) string {
	return fmt.Sprintf(`
data "alicloud_polardb_zones" "resources" {}

data "alicloud_polardb_node_classes" "default" {
%s
}
`, attrs)
}
