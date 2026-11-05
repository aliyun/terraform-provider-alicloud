package cs_test

import (
	"fmt"
	"strings"
	"testing"

	vpcSDK "github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	denverdinoCS "github.com/denverdino/aliyungo/cs"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudCSKubernetesNodes_basic(t *testing.T) {
	name := fmt.Sprintf("tf-testAcc-cs-nodes-%d", sdkacctest.RandInt())
	poolName := name + "-pool"
	secondaryPoolName := name + "-pool-b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// ACK is exercised in cn-hangzhou like the node pool resource tests;
			// the zone, instance types and VPC layout are still resolved by data
			// sources, so nothing in the configuration pins a zone or a spec.
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCsKubernetesNodesDestroy,
		Steps: []resource.TestStep{
			{
				// The node pools have to exist before the nodes can be queried: a data
				// source created in the same apply would read the cluster before its
				// nodes register.
				Config: testAccAliCloudCSKubernetesNodesConfig(name),
			},
			{
				Config: testAccAliCloudCSKubernetesNodesConfig(name) + testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "node_names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_nodes.default", "node_names.0"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "nodes.0.instance_role", "Worker"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodepool_id",
						"alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodes.0.nodepool_id",
						"alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
				),
			},
			{
				// The name filter resolves the pool ID inside the data source, so a pool
				// that references this data source does not form a dependency cycle.
				Config: testAccAliCloudCSKubernetesNodesConfig(name) + testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolName(poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "node_names.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodepool_id",
						"alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodes.0.nodepool_id",
						"alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
				),
			},
			{
				// Identity check: the same scope holds two node pools, and the name is
				// switched from the first to the second. A data source that kept the
				// nodepool_id resolved on the previous step would still answer with the
				// first pool's node, so both the pool ID and the node's own pool ID are
				// asserted against the second pool.
				Config: testAccAliCloudCSKubernetesNodesConfig(name) + testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolName(secondaryPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "node_names.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodepool_id",
						"alicloud_cs_kubernetes_node_pool.secondary", "node_pool_id"),
					resource.TestCheckResourceAttrPair("data.alicloud_cs_kubernetes_nodes.default", "nodes.0.nodepool_id",
						"alicloud_cs_kubernetes_node_pool.secondary", "node_pool_id"),
				),
			},
			{
				// A name that matches no pool yields an empty result rather than every node
				// of the cluster.
				Config: testAccAliCloudCSKubernetesNodesConfig(name) + testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolName("tf-no-such-node-pool"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "node_names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "ids.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cs_kubernetes_nodes.default", "nodepool_id"),
				),
			},
		},
	})
}

// testAccCheckCsKubernetesNodesDestroy verifies that every object the test
// creates itself is gone: the cluster, both node pools, the VPC and the vSwitch.
// CsService and VpcService keep their client field unexported, so the check talks
// to the same SDK clients through the configured provider meta instead.
func testAccCheckCsKubernetesNodesDestroy(s *terraform.State) error {
	client, ok := acctest.Provider.Meta().(*connectivity.AliyunClient)
	if !ok || client == nil {
		return fmt.Errorf("the alicloud provider is not configured")
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "alicloud_cs_managed_kubernetes":
			raw, err := client.WithCsClient(func(csClient *denverdinoCS.Client) (interface{}, error) {
				return csClient.DescribeKubernetesClusterDetail(rs.Primary.ID)
			})
			if err != nil {
				if isNotFound(err) {
					continue
				}
				return fmt.Errorf("checking cluster %s: %s", rs.Primary.ID, err)
			}
			if cluster, ok := raw.(*denverdinoCS.KubernetesClusterDetail); ok && cluster.ClusterId != rs.Primary.ID {
				continue
			}
			return fmt.Errorf("cluster %s still exists", rs.Primary.ID)
		case "alicloud_cs_kubernetes_node_pool":
			parts := strings.SplitN(rs.Primary.ID, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("unexpected node pool id %q", rs.Primary.ID)
			}
			raw, err := client.WithCsClient(func(csClient *denverdinoCS.Client) (interface{}, error) {
				return csClient.DescribeNodePoolDetail(parts[0], parts[1])
			})
			if err != nil {
				if isNotFound(err) {
					continue
				}
				return fmt.Errorf("checking node pool %s: %s", rs.Primary.ID, err)
			}
			if nodePool, ok := raw.(*denverdinoCS.NodePoolDetail); ok && nodePool.NodePoolId != parts[1] {
				continue
			}
			return fmt.Errorf("node pool %s still exists", rs.Primary.ID)
		case "alicloud_vpc":
			raw, err := client.WithVpcClient(func(vpcClient *vpcSDK.Client) (interface{}, error) {
				request := vpcSDK.CreateDescribeVpcsRequest()
				request.RegionId = client.RegionId
				request.VpcId = rs.Primary.ID
				return vpcClient.DescribeVpcs(request)
			})
			if err != nil {
				if isNotFound(err) {
					continue
				}
				return fmt.Errorf("checking vpc %s: %s", rs.Primary.ID, err)
			}
			for _, v := range raw.(*vpcSDK.DescribeVpcsResponse).Vpcs.Vpc {
				if v.VpcId == rs.Primary.ID {
					return fmt.Errorf("vpc %s still exists", rs.Primary.ID)
				}
			}
		case "alicloud_vswitch":
			raw, err := client.WithVpcClient(func(vpcClient *vpcSDK.Client) (interface{}, error) {
				request := vpcSDK.CreateDescribeVSwitchAttributesRequest()
				request.RegionId = client.RegionId
				request.VSwitchId = rs.Primary.ID
				return vpcClient.DescribeVSwitchAttributes(request)
			})
			if err != nil {
				if isNotFound(err) {
					continue
				}
				return fmt.Errorf("checking vswitch %s: %s", rs.Primary.ID, err)
			}
			if raw.(*vpcSDK.DescribeVSwitchAttributesResponse).VSwitchId == rs.Primary.ID {
				return fmt.Errorf("vswitch %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	message := strings.ToLower(err.Error())
	for _, code := range []string{
		"errorclusternotfound",
		"errornodepoolnotfound",
		"invalidvswitchid.notfound",
		"forbidden.vpcnotfound",
		"invalidvpcid.notfound",
	} {
		if strings.Contains(message, code) {
			return true
		}
	}
	return false
}

func testAccAliCloudCSKubernetesNodesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = local.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
  eni_amount           = 2
  instance_type_family = "ecs.c7"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = local.zone_id
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name         = var.name
  cluster_spec = "ack.pro.small"
  vswitch_ids  = [alicloud_vswitch.default.id]
  // The zone is resolved at run time and not every zone supports a NAT
  // gateway, so the cluster relies on the VPC's internal network instead.
  new_nat_gateway = false
  pod_cidr        = "172.20.0.0/16"
  service_cidr    = "172.21.0.0/20"
}

locals {
  # zones.0 is not a stable pick: the data source order changes between reads.
  # The alphabetically first resolved zone is used for both the vSwitch and the
  # instance type lookup, so they always agree.
  zone_id = element(sort(data.alicloud_zones.default.zones[*].id), 0)
  # The family filter keeps the candidate set small; sorting its IDs hides the
  # price order, which changes between reads of the unfiltered data source.
  node_instance_types = [element(sort(data.alicloud_instance_types.default.instance_types[*].id), 0)]
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "${var.name}-pool"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = local.node_instance_types
  desired_size         = "1"
  system_disk_category = "cloud_essd"
  system_disk_size     = "40"
  image_type           = "AliyunLinux3ContainerOptimized"
}

resource "alicloud_cs_kubernetes_node_pool" "secondary" {
  node_pool_name       = "${var.name}-pool-b"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = local.node_instance_types
  desired_size         = "1"
  system_disk_category = "cloud_essd"
  system_disk_size     = "40"
  image_type           = "AliyunLinux3ContainerOptimized"
}
`, name)
}

func testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolID() string {
	return `
data "alicloud_cs_kubernetes_nodes" "default" {
  cluster_id  = alicloud_cs_managed_kubernetes.default.id
  nodepool_id = alicloud_cs_kubernetes_node_pool.default.node_pool_id
}
`
}

func testAccAliCloudCSKubernetesNodesDataSourceConfigByPoolName(poolName string) string {
	return fmt.Sprintf(`
data "alicloud_cs_kubernetes_nodes" "default" {
  cluster_id    = alicloud_cs_managed_kubernetes.default.id
  nodepool_name = "%s"
}
`, poolName)
}
