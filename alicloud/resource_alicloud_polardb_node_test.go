package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPolarDBNode(t *testing.T) {
	var v *polardb.DBNode
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testPolardbnode-%d", rand)
	var polardbNodeMap = map[string]string{
		"db_cluster_id": CHECKSET,
		"db_node_class": CHECKSET,
	}
	resourceId := "alicloud_polardb_node.default"
	ra := resourceAttrInit(resourceId, polardbNodeMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBNodeConfigDependence)
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
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"db_node_class": "${alicloud_polardb_cluster.cluster.db_node_class}",
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
					"db_node_class": "polar.mysql.x4.large",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"db_node_class": CHECKSET}),
				),
			},
		},
	})
}

func resourcePolarDBNodeConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "vpcs_ds"{
		is_default = "true"
	}

	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
		db_node_class = "polar.mysql.x4.medium"
		vswitch_id = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.vswitch_ids.0}"
		description = "${var.name}"
	}`, name)
}
