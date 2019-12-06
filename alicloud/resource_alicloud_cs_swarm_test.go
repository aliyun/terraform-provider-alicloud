package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cs_cluster", &resource.Sweeper{
		Name: "alicloud_cs_cluster",
		F:    testSweepCSSwarms,
	})
}

func testSweepCSSwarms(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeClusters("")
	})
	if err != nil {
		return fmt.Errorf("Error retrieving CS Clusters: %s", err)
	}
	clusters, _ := raw.([]cs.ClusterType)
	sweeped := false

	var vpcIds, vswIds, groupIds, slbIds []string
	for _, v := range clusters {
		name := v.Name
		id := v.ClusterID
		if string(v.RegionID) != region {
			log.Printf("The current cluster is in the region %s. Skipped.", string(v.RegionID))
			continue
		}
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CS Clusters: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting CS Clusters: %s (%s)", name, id)
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.DeleteCluster(id)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CS Clusters (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
		vpcIds = append(vpcIds, v.VPCID)
		vswIds = append(vswIds, v.VSwitchID)
		groupIds = append(groupIds, v.SecurityGroupID)
		slbIds = append(slbIds, v.ExternalLoadbalancerID)
	}
	if sweeped {
		// Waiting 2 minutes to eusure these swarms have been deleted.
		time.Sleep(2 * time.Minute)
	}
	// Currently, the CS will retain some resources after the cluster is deleted.
	slbS := SlbService{client}
	for _, id := range slbIds {
		if err := slbS.sweepSlb(id); err != nil {
			fmt.Printf("[ERROR] Failed to deleting slb %s: %s", id, WrapError(err))
		}
	}
	ecsS := EcsService{client}
	for _, id := range groupIds {
		if err := ecsS.sweepSecurityGroup(id); err != nil {
			fmt.Printf("[ERROR] Failed to deleting SG %s: %s", id, WrapError(err))
		}
	}
	vpcS := VpcService{client}
	for _, id := range vswIds {
		if err := vpcS.sweepVSwitch(id); err != nil {
			fmt.Printf("[ERROR] Failed to deleting VSW %s: %s", id, WrapError(err))
		}
	}
	for _, id := range vpcIds {
		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = id
		if raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		}); err != nil {
			fmt.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		} else {
			response, _ := raw.(vpc.DescribeNatGatewaysResponse)
			for _, nat := range response.NatGateways.NatGateway {
				if err := vpcS.sweepNatGateway(nat.NatGatewayId); err != nil {
					fmt.Printf("[ERROR] Failed to delete nat gateway %s: %s", nat.Name, err)
				}
			}
		}
		if err := vpcS.sweepVpc(id); err != nil {
			fmt.Printf("[ERROR] Failed to deleting VPC %s: %s", id, WrapError(err))
		}
	}
	return nil
}

func SkipTestAccAlicloudCSSwarm_vpc(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.0.eip", ""),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSSwarm_vpc_zero_node(t *testing.T) {
	var container cs.ClusterType
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic_zero_node(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "0"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-basic-zero-node-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "0"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", ""),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "0"),
				),
			},

			{
				Config: testAccCSSwarm_basic_zero_node_update(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-basic-zero-node-update-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSSwarm_vpc_noslb(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic_noslb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.0.eip", ""),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "slb_id", ""),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSSwarm_update(t *testing.T) {
	var container cs.ClusterType
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_update(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-update-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
				),
			},

			{
				Config: testAccCSSwarm_updateAfter(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-updateafter-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "3"),
				),
			},
		},
	})
}

func testAccCheckContainerClusterExists(n string, d *cs.ClusterType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}

		if cluster.Primary.ID == "" {
			return fmt.Errorf("No Container cluster ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeCluster(cluster.Primary.ID)
		})
		log.Printf("[DEBUG] check cluster %s attribute %#v", cluster.Primary.ID, raw)

		if err != nil {
			return err
		}
		attr, _ := raw.(cs.ClusterType)
		if attr.ClusterID == "" {
			return fmt.Errorf("Container cluster not found")
		}

		*d = attr
		return nil
	}
}

func testAccCheckSwarmClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_swarm" {
			continue
		}

		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeCluster(rs.Primary.ID)
		})

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				continue
			}
			return err
		}
		cluster, _ := raw.(cs.ClusterType)
		if cluster.ClusterID != "" {
			return fmt.Errorf("Error container cluster %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCSSwarm_basic = `
variable "name" {
	default = "tf-testAccCSSwarm-basic"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  release_eip = "true"
}
`

func testAccCSSwarm_basic_zero_node(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-basic-zero-node-%d"
	}
	data "alicloud_images" main {
		most_recent = true
		name_regex = "^centos_6\\w{1,5}[64].*"
	}

	data "alicloud_zones" main {
		  available_resource_creation = "VSwitch"
	}
	data "alicloud_instance_types" "default" {
		 availability_zone = "${data.alicloud_zones.main.zones.0.id}"
		cpu_core_count = 1
		memory_size = 2
	}

	resource "alicloud_vpc" "foo" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 0
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	  release_eip = "true"
	}
	`, rand)
}

func testAccCSSwarm_basic_zero_node_update(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-basic-zero-node-update-%d"
	}
	data "alicloud_images" main {
		most_recent = true
		name_regex = "^centos_6\\w{1,5}[64].*"
	}

	data "alicloud_zones" main {
		  available_resource_creation = "VSwitch"
	}
	data "alicloud_instance_types" "default" {
		 availability_zone = "${data.alicloud_zones.main.zones.0.id}"
		cpu_core_count = 1
		memory_size = 2
	}

	resource "alicloud_vpc" "foo" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 2
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}

const testAccCSSwarm_basic_noslb = `
variable "name" {
	default = "tf-testAccCSSwarm-basic"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  release_eip = "true"
  need_slb = "false"
}
`

func testAccCSSwarm_update(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-update-%d"
	}
	data "alicloud_images" main {
		most_recent = true
		name_regex = "^centos_6\\w{1,5}[64].*"
	}

	data "alicloud_zones" main {
		  available_resource_creation = "VSwitch"
	}
	data "alicloud_instance_types" "default" {
		 availability_zone = "${data.alicloud_zones.main.zones.0.id}"
		cpu_core_count = 1
		memory_size = 2
	}

	resource "alicloud_vpc" "foo" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 2
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}

func testAccCSSwarm_updateAfter(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-updateafter-%d"
	}
	data "alicloud_images" main {
		most_recent = true
		name_regex = "^centos_6\\w{1,5}[64].*"
	}

	data "alicloud_zones" main {
		  available_resource_creation = "VSwitch"
	}
	data "alicloud_instance_types" "default" {
		 availability_zone = "${data.alicloud_zones.main.zones.0.id}"
		cpu_core_count = 1
		memory_size = 2
	}

	resource "alicloud_vpc" "foo" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 3
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}
