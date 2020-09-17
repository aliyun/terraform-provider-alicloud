package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPolarDBDatabase_update(t *testing.T) {
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
			testAccPreCheckWithNoDefaultVpc(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":  "${alicloud_polardb_cluster.instance.id}",
					"db_name":        "tftestdatabase",
					"db_description": "test",
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

func resourcePolarDBDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	variable "instancechargetype" {
		default = "PostPaid"
	}

	variable "engine" {
		default = "MySQL"
	}

	variable "engineversion" {
		default = "8.0"
	}

	variable "instanceclass" {
		default = "polar.mysql.x4.large"
	}

	resource "alicloud_polardb_cluster" "instance" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		pay_type = "${var.instancechargetype}"
		db_node_class = "${var.instanceclass}"
		vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
		description = "${var.name}"
	}`, PolarDBCommonTestCase, name)
}
