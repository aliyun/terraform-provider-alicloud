package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPolarDBAccountPrivilege_update(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.PolarDBSupportRegions)
	var v *polardb.DBAccount
	name := "tf-testAccPolarDBAccountPrivilege_update"
	resourceId := "alicloud_polardb_account_privilege.default"
	var basicMap = map[string]string{
		"db_cluster_id":     CHECKSET,
		"account_name":      "tftestprivilege",
		"account_privilege": "ReadOnly",
		"db_names.#":        "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBAccountPrivilegeConfigDependence)

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
					"db_cluster_id":     "${alicloud_polardb_cluster.default.id}",
					"account_name":      "${alicloud_polardb_account.default.account_name}",
					"account_privilege": "ReadOnly",
					"db_names":          "${alicloud_polardb_database.default.*.db_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_cluster.default.id}",
					"account_name":      "${alicloud_polardb_account.default.account_name}",
					"account_privilege": "ReadOnly",
					"db_names":          []string{"${alicloud_polardb_database.default.0.db_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_cluster.default.id}",
					"account_name":      "${alicloud_polardb_account.default.account_name}",
					"account_privilege": "ReadOnly",
					"db_names":          "${alicloud_polardb_database.default.*.db_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func resourcePolarDBAccountPrivilegeConfigDependence(name string) string {
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
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
        db_node_class = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
		vswitch_id = local.vswitch_id
		description = "${var.name}"
	}
	resource "alicloud_polardb_database" "default" {
	  count = 2
	  db_cluster_id = "${alicloud_polardb_cluster.default.id}"
	  db_name = "tfaccountpri_${count.index}"
	  db_description = "from terraform"
	}

	resource "alicloud_polardb_account" "default" {
	  db_cluster_id = "${alicloud_polardb_cluster.default.id}"
	  account_name = "tftestprivilege"
	  account_type = "Normal"
	  account_password = "Test12345"
	  account_description = "from terraform"
	}
`, PolarDBCommonTestCase, name)
}
