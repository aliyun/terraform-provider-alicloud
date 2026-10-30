package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRdsDBReadWriteSplittingConnectionMssql_create(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_read_write_splitting_connection.create"

	var dbRwSplitConnMap = map[string]string{
		"instance_id":       acctest.CHECKSET,
		"distribution_type": "Standard",
	}

	ra := acctest.ResourceAttrInit(resourceId, dbRwSplitConnMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBInstanceNetInfo")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := "tf-testAccDBRwSplitConnMssql_create"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDBReadWriteSplittingConfigDependenceMssql)
	// connection_prefix is Optional+ForceNew (Schema string len 1..31). Generate one unique
	// alphanumeric prefix per run and reuse the same value across both steps to avoid a
	// ForceNew rebuild between steps; Read parses the value back from the real ConnectionString.
	rwPrefix := fmt.Sprintf("tf%d", tfacctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"distribution_type": "Standard",
					"connection_prefix": rwPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": rwPrefix,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"distribution_type": "Custom",
					"connection_prefix": rwPrefix,
					"weight": map[string]interface{}{
						"\"${alicloud_db_readonly_instance.default.id}\"": "200",
						"master": "200",
						"slave":  "400",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testAccCheck(map[string]string{
						"weight.%":          "3",
						"distribution_type": "Custom",
					}),
				),
			},
		},
	})

}

func resourceDBReadWriteSplittingConfigDependenceMssql(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine                   = "SQLServer"
	engine_version           = "2017_ent"
	instance_charge_type     = "PostPaid"
 	db_instance_storage_type = "cloud_essd"
	category                 = "AlwaysOn"
}

data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vswitches.default.vswitches[0].vpc_id
}

resource "alicloud_db_instance" "default" {
    engine                   = "SQLServer"
	engine_version           = "2017_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_type            = "mssql.x4.medium.e2"
	instance_storage         = 20
	category                 = "AlwaysOn"
	vswitch_id               = data.alicloud_vswitches.default.vswitches[0].vswitch_id
	instance_name            = var.name
}

resource "alicloud_db_readonly_instance" "default" {
	zone_id               = alicloud_db_instance.default.zone_id
	master_db_instance_id = alicloud_db_instance.default.id
	engine_version        = alicloud_db_instance.default.engine_version
	instance_storage      = alicloud_db_instance.default.instance_storage
	instance_type         = "rds.mssql.s2.large"
	instance_name         = "${var.name}_ro"
	vswitch_id            = data.alicloud_vswitches.default.vswitches[0].vswitch_id
}

`, name)
}
