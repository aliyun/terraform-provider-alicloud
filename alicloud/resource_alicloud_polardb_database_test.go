package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBDatabase_update(t *testing.T) {
	var database *polardb.Database
	resourceId := "alicloud_polardb_database.default"

	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDatabaseConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":      "${alicloud_polardb_cluster.instance.id}",
					"db_name":            "tftestdatabase",
					"character_set_name": "utf8",
					"db_description":     "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"db_description": "from terraform"}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBDatabasePostgreSQL_update(t *testing.T) {
	var database *polardb.Database
	resourceId := "alicloud_polardb_database.default"

	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDatabaseConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":      "${alicloud_polardb_cluster.instance.id}",
					"db_name":            "tftestdatabase",
					"character_set_name": "utf8",
					"db_description":     "test",
					"account_name":       "${alicloud_polardb_account.acc1.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourcePolarDBDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  zone_id    = local.zone_id
      category   = "Normal"
	}
	resource "alicloud_polardb_cluster" "instance" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
		db_node_class = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
		vswitch_id = local.vswitch_id
		description = "${var.name}"
    }`, PolarDBCommonTestCase, name)
}

func resourcePolarDBDatabaseConfigPostgreSQLDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "pg"
	  db_version = "14"
      pay_type   = "PostPaid"
	  zone_id    = local.zone_id
      category   = "Normal"
	}
	resource "alicloud_polardb_cluster" "instance" {
		db_type = "PostgreSQL"
		db_version = "14"
		pay_type = "PostPaid"
		db_node_class = "polar.pg.x4.medium"
		vswitch_id = local.vswitch_id
		description = "${var.name}"
	}
    resource "alicloud_polardb_account" "acc1" {
       db_cluster_id       = alicloud_polardb_cluster.instance.id
       account_name        = "terraformtest"
       account_password    = "Example1234"
       account_description = "terraformTest"
    }`, PolarDBPostgreSQLCommonTestCase, name)
}
