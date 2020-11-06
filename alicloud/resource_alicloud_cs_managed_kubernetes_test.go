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
					"worker_number":               "2",
					"password":                    "Test12345",
					"pod_cidr":                    "172.20.0.0/16",
					"service_cidr":                "172.21.0.0/20",
					"worker_disk_size":            "50",
					"worker_disk_category":        "cloud_ssd",
					"worker_data_disk_size":       "20",
					"worker_data_disk_category":   "cloud_ssd",
					"worker_instance_charge_type": "PostPaid",
					"slb_internet_enabled":        "true",
					"cluster_spec":                "ack.pro.small",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":         "false",
					"timezone":                    "Asia/Shanghai",
					"os_type":                     "Linux",
					"platform":                    "CentOS",
					"node_port_range":             "30000-32767",
					"cluster_domain":              "cluster.local",
					"custom_san":                  "www.terraform.io",
					"encryption_provider_key":     "${data.alicloud_kms_secrets.default.secrets.0.id}",
					"runtime":                     map[string]interface{}{"Name": "docker", "Version": "19.03.5"},
					"rds_instances":               []string{"${alicloud_db_instance.default.id}"},
					"taints":                      []map[string]string{{"key": "tf-key1", "value": "tf-value1", "effect": "NoSchedule"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"worker_number":             "2",
						"password":                  "Test12345",
						"pod_cidr":                  "172.20.0.0/16",
						"service_cidr":              "172.21.0.0/20",
						"worker_disk_size":          "50",
						"worker_disk_category":      "cloud_ssd",
						"worker_data_disk_size":     "20",
						"worker_data_disk_category": "cloud_ssd",
						"slb_internet_enabled":      "true",
						"cluster_spec":              "ack.pro.small",
						"resource_group_id":         CHECKSET,
						"deletion_protection":       "false",
						"timezone":                  "Asia/Shanghai",
						"os_type":                   "Linux",
						"platform":                  "CentOS",
						"node_port_range":           "30000-32767",
						"cluster_domain":            "cluster.local",
						"custom_san":                "www.terraform.io",
						"rds_instances.#":           "1",
						"taints.#":                  "1",
						"taints.0.key":              "tf-key1",
						"taints.0.value":            "tf-value1",
						"taints.0.effect":           "NoSchedule",
						"runtime.Name":              "docker",
						"runtime.Version":           "19.03.5",
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
					"cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "encryption_provider_key", "rds_instances"},
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

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_secrets" "default" {}

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

resource "alicloud_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alicloud_vswitch.default.id}"
  monitoring_period    = "60"
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
					"cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "encryption_provider_key", "rds_instances"},
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
