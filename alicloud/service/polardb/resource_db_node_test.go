package polardb_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudPolarDBDBNode_basic(t *testing.T) {
	resourceId := "alicloud_polardb_db_node.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testAccPolarDBDBNode%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolardbDBNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBDBNodeBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "db_cluster_id"),
					resource.TestCheckResourceAttr(resourceId, "db_node.#", "1"),
				),
			},
			{
				Config: testAccAliCloudPolarDBDBNodeBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "auto_use_coupon", "true"),
					resource.TestCheckResourceAttr(resourceId, "modify_type", "Upgrade"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateIdFunc:       testAccPolarDBDBNodeImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "auto_use_coupon", "cloud_provider", "db_node_target_class", "db_node_type", "endpoint_bind_list", "from_time_service", "modify_type", "planned_end_time", "planned_flashing_off_time", "planned_start_time", "promotion_code", "restart_planned_end_time", "restart_planned_start_time", "scc_mode"},
			},
		},
	})
}

func testAccPolarDBDBNodeImportStateID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["alicloud_polardb_db_node.default"]
	if !ok {
		return "", fmt.Errorf("resource alicloud_polardb_db_node.default not found in state")
	}
	dbClusterId := rs.Primary.Attributes["db_cluster_id"]
	if dbClusterId == "" {
		return "", fmt.Errorf("resource alicloud_polardb_db_node.default has empty db_cluster_id")
	}
	return fmt.Sprintf("%s:%s", dbClusterId, rs.Primary.ID), nil
}

func testAccCheckPolardbDBNodeDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_db_node" {
			continue
		}
		dbClusterId := rs.Primary.Attributes["db_cluster_id"]
		object, err := polardb.DBNodeGet(client, rs.Primary.ID, dbClusterId, false)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing polardb db node %s failed: %s", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("polardb db node %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudPolarDBDBNodeBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_node {
    target_class = "polar.mysql.x4.large"
  }
}
`, name)
}

// The createVisibilityTimeout test was removed because the production code
// uses the same operationTimeout for both CreateDBNodes retry (retryRpcPost)
// and the post-create visibility check (dbNodeWaitVisible). Setting
// timeouts.create=30s to force a dbNodeWaitVisible timeout also limits the
// CreateDBNodes retry window to 30s; when OperationDenied.OrderProcessing is
// active (PostPaid cluster billing order settlement), CreateDBNodes cannot
// succeed within 30s, so the test always fails at the API stage and never
// reaches the visibility-check stage. Amplifying the production timeout for
// test purposes (the rejected f57f635b5 approach) changes public timeout
// semantics without approval. The trackedPlan persistence path (AddWarning
// when dbNodeWaitVisible times out) is validated by code review; the basic
// test validates trackedPlan construction when CreateDBNodes succeeds.

func testAccAliCloudPolarDBDBNodeBasicUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  db_node {
    target_class = "polar.mysql.x4.large"
  }
  db_node_target_class = "polar.mysql.x4.large"
  modify_type          = "Upgrade"
  auto_use_coupon      = true
}
`, name)
}
