package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudRdsRCCluster_basic covers the import, read, import-state and
// delete lifecycle of alicloud_rds_rc_cluster. The cluster imports an ACK
// managed cluster that the test provisions itself, so the lifecycle no longer
// depends on a pre-imported cluster supplied through the environment.
func TestAccAliCloudRdsRCCluster_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-rc-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRCClusterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster.default", "id"),
					resource.TestCheckResourceAttrPair("alicloud_rds_rc_cluster.default", "cluster_id", "alicloud_cs_managed_kubernetes.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster.default", "vcluster_status"),
					resource.TestCheckResourceAttrPair("alicloud_rds_rc_cluster.default", "vpc_id", "alicloud_vpc.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_rds_rc_cluster.default", "region_id"),
				),
			},
			{
				ResourceName:      "alicloud_rds_rc_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
				// The import endpoint only carries the cluster ID, and the API
				// does not return the kubeconfig credentials or the vSwitch.
				ImportStateVerifyIgnore: []string{"client_key", "client_certificate", "v_switch_id", "support_disk_category"},
			},
		},
	})
}

func testAccRCClusterConfig(name string) string {
	return testAccRCClusterAckClusterConfig(name) + `
resource "alicloud_rds_rc_cluster" "default" {
  cluster_id         = alicloud_cs_managed_kubernetes.default.id
  vpc_id             = alicloud_vpc.default.id
  v_switch_id        = alicloud_vswitch.default.id
  client_key         = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_key
  client_certificate = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_cert
}
`
}

// testAccRCClusterAckClusterConfig provisions the ACK managed cluster that a
// RDS Custom container cluster imports, together with its VPC, vSwitch and
// kubeconfig credentials. Shared by the container cluster and node pool tests.
func testAccRCClusterAckClusterConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 37)
  service_cidr         = cidrsubnet("172.18.0.0/16", 4, 7)
  slb_internet_enabled = true
}

data "alicloud_cs_cluster_credential" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}
`, name)
}
