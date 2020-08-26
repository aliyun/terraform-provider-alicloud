package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSManagedKubernetes_basic(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := resourceAttrInit(resourceId, csManagedKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                        name,
					"worker_vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_number":               "1",
					"password":                    "Test12345",
					"pod_cidr":                    "172.20.0.0/16",
					"service_cidr":                "172.21.0.0/20",
					"worker_disk_size":            "50",
					"worker_disk_category":        "cloud_ssd",
					"worker_data_disk_size":       "20",
					"worker_data_disk_category":   "cloud_ssd",
					"worker_instance_charge_type": "PostPaid",
					"slb_internet_enabled":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"worker_number":             "1",
						"password":                  "Test12345",
						"pod_cidr":                  "172.20.0.0/16",
						"service_cidr":              "172.21.0.0/20",
						"worker_disk_size":          "50",
						"worker_disk_category":      "cloud_ssd",
						"worker_data_disk_size":     "20",
						"worker_data_disk_category": "cloud_ssd",
						"slb_internet_enabled":      "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor", "user_ca", "force_update",
					"node_cidr_mask", "slb_internet_enabled", "vswitch_ids", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_instance_types", "log_config",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes",
					"cpu_policy", "proxy_mode"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_nat_gateway": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_nat_gateway": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_number": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_number": "5",
					}),
				),
			},
		},
	})
}

func resourceCSManagedKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_log_project" "log" {
  name        = "${var.name}"
  description = "created by terraform for managedkubernetes cluster"
}
`, name)
}

func resourceCSManagedKubernetesConfigDependence_multiAZ(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_instance_types" "default1" {
	availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_instance_types" "default2" {
	availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-2], "id")}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "default1" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.2.0/24"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
}

resource "alicloud_vswitch" "default2" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.3.0/24"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-2], "id")}"
}

resource "alicloud_nat_gateway" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  specification = "Small"
}

resource "alicloud_snat_entry" "default" {
  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.default.id}"
  snat_ip           = "${alicloud_eip.default.ip_address}"
}

resource "alicloud_snat_entry" "default1" {
  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.default1.id}"
  snat_ip           = "${alicloud_eip.default.ip_address}"
}

resource "alicloud_snat_entry" "default2" {
  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.default2.id}"
  snat_ip           = "${alicloud_eip.default.ip_address}"
}

resource "alicloud_eip" "default" {
  name = "${var.name}"
  bandwidth = "100"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${alicloud_eip.default.id}"
  instance_id   = "${alicloud_nat_gateway.default.id}"
}
`, name)
}

func TestAccAlicloudCSManagedKubernetes_upgrade(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := resourceAttrInit(resourceId, csManagedKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                        name,
					"worker_vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_number":               "3",
					"password":                    "Test12345",
					"pod_cidr":                    "172.20.0.0/16",
					"service_cidr":                "172.21.0.0/20",
					"worker_disk_size":            "50",
					"worker_disk_category":        "cloud_ssd",
					"worker_data_disk_size":       "20",
					"worker_data_disk_category":   "cloud_ssd",
					"worker_instance_charge_type": "PostPaid",
					"slb_internet_enabled":        "true",
					"version":                     "1.14.8-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"password":                  "Test12345",
						"pod_cidr":                  "172.20.0.0/16",
						"service_cidr":              "172.21.0.0/20",
						"worker_disk_size":          "50",
						"worker_disk_category":      "cloud_ssd",
						"worker_data_disk_size":     "20",
						"worker_data_disk_category": "cloud_ssd",
						"slb_internet_enabled":      "true",
						"version":                   "1.14.8-aliyun.1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor", "user_ca", "force_update",
					"node_cidr_mask", "slb_internet_enabled", "vswitch_ids", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_instance_types", "log_config",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes",
					"cpu_policy", "proxy_mode"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "1.16.9-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.16.9-aliyun.1",
					}),
				),
			},
		},
	})
}

var csManagedKubernetesBasicMap = map[string]string{
	"new_nat_gateway":             "true",
	"worker_number":               "3",
	"worker_instance_types.0":     CHECKSET,
	"worker_disk_size":            "40",
	"worker_disk_category":        "cloud_efficiency",
	"worker_data_disk_size":       "40",
	"worker_instance_charge_type": "PostPaid",
	"slb_internet_enabled":        "true",
	"install_cloud_monitor":       "true",
	"force_update":                "false",
}
