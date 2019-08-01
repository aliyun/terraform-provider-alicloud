package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSManagedKubernetes_basic(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions) },

		IDRefreshName: "alicloud_cs_managed_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckManagedKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKubernetes_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_managed_kubernetes.k8s", &k8s),
					resource.TestMatchResourceAttr("alicloud_cs_managed_kubernetes.k8s", "name", regexp.MustCompile("^tf-testAccManagedKubernetes-basic*")),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_size", "50"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_types.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "vswitch_ids.#", "1"),

					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "vpc_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "security_group_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "availability_zone"),
				),
			},
			{
				ResourceName:      "alicloud_cs_managed_kubernetes.k8s",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name_prefix", "new_nat_gateway", "pod_cidr",
					"service_cidr", "password", "install_cloud_monitor", "slb_internet_enabled",
					"vswitch_ids", "worker_instance_types", "worker_numbers", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_number", "force_update"},
			},
		},
	})
}

func TestAccAlicloudCSManagedKubernetes_multiaz(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions) },

		IDRefreshName: "alicloud_cs_managed_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckManagedKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKubernetes_multiaz,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_managed_kubernetes.k8s", &k8s),
					resource.TestMatchResourceAttr("alicloud_cs_managed_kubernetes.k8s", "name", regexp.MustCompile("^tf-testAccManagedKubernetes-multiaz*")),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "vpc_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "security_group_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "availability_zone"),
				),
			},
		},
	})
}

func TestAccAlicloudCSManagedKubernetes_update(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions) },

		IDRefreshName: "alicloud_cs_managed_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckManagedKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKubernetes_update_before,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_managed_kubernetes.k8s", &k8s),
					resource.TestMatchResourceAttr("alicloud_cs_managed_kubernetes.k8s", "name", regexp.MustCompile("^tf-testAccManagedKubernetes-update*")),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_types.#", "1"),

					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "vpc_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "security_group_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "availability_zone"),
				),
			},
			{
				Config: testAccManagedKubernetes_update_after,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_managed_kubernetes.k8s", &k8s),
					resource.TestMatchResourceAttr("alicloud_cs_managed_kubernetes.k8s", "name", regexp.MustCompile("^tf-testAccManagedKubernetes-update*")),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_number", "4"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_cs_managed_kubernetes.k8s", "worker_instance_types.#", "1"),

					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "vpc_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "security_group_id"),
					resource.TestCheckResourceAttrSet("alicloud_cs_managed_kubernetes.k8s", "availability_zone"),
				),
			},
		},
	})
}

func testAccCheckManagedKubernetesClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_managed_kubernetes" {
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

const testAccManagedKubernetes_basic = `
variable "name" {
	default = "tf-testAccManagedKubernetes-basic"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
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

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  vswitch_ids = ["${alicloud_vswitch.foo.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 2
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  worker_disk_size = 50
}
`

const testAccManagedKubernetes_update_before = `
variable "name" {
	default = "tf-testAccManagedKubernetes-update"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
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

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  vswitch_ids = ["${alicloud_vswitch.foo.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 2
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  worker_disk_category  = "cloud_efficiency"
}
`

const testAccManagedKubernetes_update_after = `
variable "name" {
	default = "tf-testAccManagedKubernetes-update"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
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

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  vswitch_ids = ["${alicloud_vswitch.foo.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 4
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  worker_disk_category  = "cloud_efficiency"
}
`

const testAccManagedKubernetes_multiaz = `
variable "name" {
	default = "tf-testAccManagedKubernetes-multiaz"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "instance_types_1_worker" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_instance_types" "instance_types_2_worker" {
	availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones)-1)%length(data.alicloud_zones.main.zones)], "id")}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_instance_types" "instance_types_3_worker" {
	availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones)-2)%length(data.alicloud_zones.main.zones)], "id")}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "vsw1" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_vswitch" "vsw2" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.2.0/24"
  availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones)-1)%length(data.alicloud_zones.main.zones)], "id")}"
}

resource "alicloud_vswitch" "vsw3" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.3.0/24"
  availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones)-2)%length(data.alicloud_zones.main.zones)], "id")}"
}

resource "alicloud_nat_gateway" "nat_gateway" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  specification = "Small"
}

resource "alicloud_snat_entry" "snat_entry_1" {
  snat_table_id     = "${alicloud_nat_gateway.nat_gateway.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.vsw1.id}"
  snat_ip           = "${alicloud_eip.eip.ip_address}"
}

resource "alicloud_snat_entry" "snat_entry_2" {
  snat_table_id     = "${alicloud_nat_gateway.nat_gateway.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.vsw2.id}"
  snat_ip           = "${alicloud_eip.eip.ip_address}"
}

resource "alicloud_snat_entry" "snat_entry_3" {
  snat_table_id     = "${alicloud_nat_gateway.nat_gateway.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.vsw3.id}"
  snat_ip           = "${alicloud_eip.eip.ip_address}"
}

resource "alicloud_eip" "eip" {
  name = "${var.name}"
  bandwidth = "100"
}

resource "alicloud_eip_association" "eip_asso" {
  allocation_id = "${alicloud_eip.eip.id}"
  instance_id   = "${alicloud_nat_gateway.nat_gateway.id}"
}

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  vswitch_ids = ["${alicloud_vswitch.vsw1.id}", "${alicloud_vswitch.vsw2.id}", "${alicloud_vswitch.vsw3.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.instance_types_1_worker.instance_types.0.id}", "${data.alicloud_instance_types.instance_types_2_worker.instance_types.0.id}", "${data.alicloud_instance_types.instance_types_3_worker.instance_types.0.id}"]
  worker_number = 2
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
}
`
