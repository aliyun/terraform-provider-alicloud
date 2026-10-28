package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudRdsRCClusterNodePool_basic covers the create, update, import
// and delete lifecycle of alicloud_rds_rc_cluster_node_pool.
//
// The node pool is a sub-resource of an RDS Custom container cluster. The test
// provisions the whole dependency chain itself: ACK managed cluster, the
// alicloud_rds_rc_cluster import, then the node pool, and destroys it in
// reverse order through the resource graph.
func TestAccAliCloudRdsRCClusterNodePool_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-rc-nodepool-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRdsRCClusterNodePoolConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster_node_pool.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster_node_pool.default", "node_pool_id"),
					resource.TestCheckResourceAttrPair("alicloud_rds_rc_cluster_node_pool.default", "cluster_id", "alicloud_rds_rc_cluster.default", "id"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "name", name),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "scaling_group.#", "1"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "scaling_group.0.desired_size", "0"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "kubernetes_config.0.labels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "kubernetes_config.0.labels.0.key", "env"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "kubernetes_config.0.labels.0.value", "prod"),
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster_node_pool.default", "status.0.state"),
				),
			},
			{
				Config: testAccAliCloudRdsRCClusterNodePoolConfig(name + "-upd"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "name", name+"-upd"),
					resource.TestCheckResourceAttr("alicloud_rds_rc_cluster_node_pool.default", "scaling_group.0.desired_size", "0"),
				),
			},
			{
				ResourceName:      "alicloud_rds_rc_cluster_node_pool.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudRdsRCClusterNodePoolConfig(name string) string {
	return testAccRCClusterAckClusterConfig(name) + fmt.Sprintf(`
resource "alicloud_rds_rc_cluster" "default" {
  cluster_id         = alicloud_cs_managed_kubernetes.default.id
  vpc_id             = alicloud_vpc.default.id
  v_switch_id        = alicloud_vswitch.default.id
  client_key         = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_key
  client_certificate = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_cert
}

resource "alicloud_rds_rc_cluster_node_pool" "default" {
  cluster_id = alicloud_rds_rc_cluster.default.id
  name       = "%s"

  scaling_group = [
    {
      v_switch_ids         = [alicloud_vswitch.default.id]
      instance_types       = ["mysql.x4.xlarge.6cm"]
      instance_charge_type = "PostPaid"
      period               = 1
      period_unit          = "Month"
      auto_renew           = false
      security_group_ids   = [alicloud_cs_managed_kubernetes.default.security_group_id]
      desired_size         = 0
    }
  ]

  kubernetes_config = [
    {
      labels = [
        {
          key   = "env"
          value = "prod"
        }
      ]
      taints = []
    }
  ]
}
`, name)
}
