package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudPolarDBDBNode_basic(t *testing.T) {
	r := tfacctest.RandStringFromCharSet(8, tfacctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolardbDBNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBDBNodeConfig(r, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "id"),
					resource.TestCheckResourceAttrPair("alicloud_polardb_db_node.default", "db_cluster_id", "alicloud_polardb_cluster.default", "id"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node_type", "RO"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node.#", "1"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node.0.target_class", "polar.mysql.x4.medium"),
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "db_node.0.zone_id"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node_ids.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_polardb_db_node.default", "db_node_ids.0", "alicloud_polardb_db_node.default", "id"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "imci_switch", "OFF"),
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "resource_group_id"),
				),
			},
			{
				ResourceName:            "alicloud_polardb_db_node.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_use_coupon", "db_node_type"},
			},
			{
				// ModifyDBNodeClass: the resize pair must change together; the
				// coupon flag rides along with the resize request.
				Config: testAccAliCloudPolarDBDBNodeConfig(r, `
  auto_use_coupon      = false
  db_node_target_class = "polar.mysql.x4.large"
  modify_type          = "Upgrade"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "auto_use_coupon", "false"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node_target_class", "polar.mysql.x4.large"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "modify_type", "Upgrade"),
				),
			},
			{
				// ModifyDBNodeSccMode: SCC mode is a write-only update parameter.
				Config: testAccAliCloudPolarDBDBNodeConfig(r, `
  auto_use_coupon      = false
  db_node_target_class = "polar.mysql.x4.large"
  modify_type          = "Upgrade"
  scc_mode             = "ON"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "scc_mode", "ON"),
				),
			},
			{
				// RestartDBNode runs last: FromTimeService=true restarts the node
				// immediately, and no further write should race the restart.
				Config: testAccAliCloudPolarDBDBNodeConfig(r, `
  auto_use_coupon      = false
  db_node_target_class = "polar.mysql.x4.large"
  modify_type          = "Upgrade"
  scc_mode             = "ON"
  from_time_service    = "true"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "from_time_service", "true"),
				),
			},
		},
	})
}

func testAccCheckPolardbDBNodeDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_db_node" {
			continue
		}
		id := rs.Primary.ID
		response, err := client.RpcPost("polardb", "2017-08-01", "DescribeDBClusters", nil, map[string]interface{}{
			"RegionId":  client.RegionId,
			"DBNodeIds": id,
		}, true)
		if err != nil {
			return fmt.Errorf("describing clusters while checking the destroyed DB node %s: %w", id, err)
		}
		if dbNodeIDFound(response, id) {
			return fmt.Errorf("the DB node %s still exists after destroy", id)
		}
	}
	return nil
}

func dbNodeIDFound(response map[string]interface{}, id string) bool {
	items, _ := response["Items"].(map[string]interface{})
	clusters, _ := items["DBCluster"].([]interface{})
	for _, item := range clusters {
		cluster, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		nodes, _ := cluster["DBNodes"].(map[string]interface{})
		children, _ := nodes["DBNode"].([]interface{})
		for _, child := range children {
			node, ok := child.(map[string]interface{})
			if ok && fmt.Sprint(node["DBNodeId"]) == id {
				return true
			}
		}
	}
	return false
}

func testAccAliCloudPolarDBDBNodeConfig(r, nodeArgs string) string {
	// The zone comes from the class data source; the classes stay literal because
	// the resize step needs a specific larger spec than the create step and the
	// data source does not expose a stable "next size up" relation.
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccPolarDBDBNode-%[1]s"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  zone_id    = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

resource "alicloud_polardb_cluster" "default" {
  db_type            = "MySQL"
  db_version         = "8.0"
  pay_type           = "PostPaid"
  db_node_class      = "polar.mysql.x4.medium"
  vswitch_id         = alicloud_vswitch.default.id
  description        = var.name
  encrypt_new_tables = "OFF"

  # ModifyDBNodeClass changes the cluster-level node class for every node, so the
  # node resource owns db_node_class here; the cluster must not fight that change.
  lifecycle {
    ignore_changes = [db_node_class]
  }
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_node_type  = "RO"
  imci_switch   = "OFF"

  db_node = [{
    target_class = "polar.mysql.x4.medium"
    zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  }]
%[2]s
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`, r, nodeArgs)
}
