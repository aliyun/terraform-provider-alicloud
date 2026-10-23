package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBDBNode_basic(t *testing.T) {
	r := tfacctest.RandStringFromCharSet(8, tfacctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBDBNodeConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "id"),
					resource.TestCheckResourceAttrPair("alicloud_polardb_db_node.default", "db_cluster_id", "alicloud_polardb_cluster.default", "id"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node_type", "RO"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node.#", "1"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node.0.target_class", "polar.mysql.x4.medium"),
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "db_node.0.zone_id"),
					resource.TestCheckResourceAttr("alicloud_polardb_db_node.default", "db_node_ids.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_polardb_db_node.default", "db_node_ids.0", "alicloud_polardb_db_node.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_polardb_db_node.default", "resource_group_id"),
				),
			},
			{
				ResourceName:            "alicloud_polardb_db_node.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_use_coupon", "db_node_type"},
			},
		},
	})
}

func testAccAliCloudPolarDBDBNodeConfig(r string) string {
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
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_node_type  = "RO"
  db_node = [{
    target_class = "polar.mysql.x4.medium"
    zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  }]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`, r)
}
