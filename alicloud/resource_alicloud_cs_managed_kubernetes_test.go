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
					"version":                     "1.18.8-aliyun.1",
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
					"load_balancer_spec":          "slb.s2.small",
					"cluster_spec":                "ack.pro.small",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":         "true",
					"timezone":                    "Asia/Shanghai",
					"os_type":                     "Linux",
					"platform":                    "CentOS",
					"node_port_range":             "30000-32767",
					"cluster_domain":              "cluster.local",
					"custom_san":                  "www.terraform.io",
					"encryption_provider_key":     "${data.alicloud_kms_keys.default.keys.0.id}",
					"runtime":                     map[string]interface{}{"Name": "docker", "Version": "19.03.5"},
					"rds_instances":               []string{"${alicloud_db_instance.default.id}"},
					"taints":                      []map[string]string{{"key": "tf-key1", "value": "tf-value1", "effect": "NoSchedule"}},
					"maintenance_window":          []map[string]string{{"enable": "true", "maintenance_time": "03:00:00Z", "duration": "3h", "weekly_period": "Thursday"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
						"version":                               "1.18.8-aliyun.1",
						"worker_number":                         "2",
						"password":                              "Test12345",
						"pod_cidr":                              "172.20.0.0/16",
						"service_cidr":                          "172.21.0.0/20",
						"worker_disk_size":                      "50",
						"worker_disk_category":                  "cloud_ssd",
						"worker_data_disk_size":                 "20",
						"worker_data_disk_category":             "cloud_ssd",
						"slb_internet_enabled":                  "true",
						"cluster_spec":                          "ack.pro.small",
						"resource_group_id":                     CHECKSET,
						"deletion_protection":                   "true",
						"timezone":                              "Asia/Shanghai",
						"os_type":                               "Linux",
						"platform":                              "CentOS",
						"node_port_range":                       "30000-32767",
						"cluster_domain":                        "cluster.local",
						"custom_san":                            "www.terraform.io",
						"rds_instances.#":                       "1",
						"taints.#":                              "1",
						"taints.0.key":                          "tf-key1",
						"taints.0.value":                        "tf-value1",
						"taints.0.effect":                       "NoSchedule",
						"runtime.Name":                          "docker",
						"runtime.Version":                       "19.03.5",
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "true",
						"maintenance_window.0.maintenance_time": "03:00:00Z",
						"maintenance_window.0.duration":         "3h",
						"maintenance_window.0.weekly_period":    "Thursday",
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
					"worker_instance_charge_type", "worker_instance_types", "maintenance_window", "log_config",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes",
					"cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "encryption_provider_key", "rds_instances", "load_balancer_spec"},
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
					"name": "tf-managed-k8s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-managed-k8s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintenance_window": []map[string]string{{"enable": "true", "maintenance_time": "05:00:00Z", "duration": "5h", "weekly_period": "Monday,Thursday"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "true",
						"maintenance_window.0.maintenance_time": "05:00:00Z",
						"maintenance_window.0.duration":         "5h",
						"maintenance_window.0.weekly_period":    "Monday,Thursday",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_number": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_number": "3",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCSManagedKubernetes_essd(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.EssdSupportRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"password":                       "Test12345",
					"pod_cidr":                       "172.20.0.0/16",
					"version":                        "1.18.8-aliyun.1",
					"service_cidr":                   "172.21.0.0/20",
					"deletion_protection":            "true",
					"worker_number":                  "3",
					"worker_data_disk_category":      "cloud_ssd",
					"worker_data_disk_size":          "20",
					"worker_instance_charge_type":    "PostPaid",
					"cluster_spec":                   "ack.pro.small",
					"worker_vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_disk_category":           "cloud_essd",
					"worker_disk_size":               "120",
					"worker_disk_performance_level":  "PL0",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_data_disks":              []map[string]string{{"category": "cloud_essd", "size": "120", "auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}", "performance_level": "PL0"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"version":                        "1.18.8-aliyun.1",
						"password":                       "Test12345",
						"pod_cidr":                       "172.20.0.0/16",
						"service_cidr":                   "172.21.0.0/20",
						"deletion_protection":            "true",
						"worker_number":                  "3",
						"worker_data_disk_category":      "cloud_ssd",
						"worker_data_disk_size":          "20",
						"worker_instance_charge_type":    "PostPaid",
						"cluster_spec":                   "ack.pro.small",
						"worker_disk_size":               "120",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_performance_level":  "PL0",
						"worker_disk_snapshot_policy_id": CHECKSET,
						"worker_data_disks.#":            "1",
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
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_number", "worker_vswitch_ids", "exclude_autoscaler_nodes",
					"cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime",
					"worker_disk_snapshot_policy_id", "worker_disk_performance_level", "taints", "encryption_provider_key", "worker_data_disks", "rds_instances", "load_balancer_spec"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// worker args. scale out
					"worker_number":                  "4",
					"worker_vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_disk_category":           "cloud_essd",
					"worker_disk_size":               "100",
					"worker_disk_performance_level":  "PL1",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_data_disks":              []map[string]string{{"category": "cloud_essd", "size": "100", "auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}", "performance_level": "PL1"}},
					// global args
					"new_nat_gateway":     "true",
					"name":                "tf-managed-k8s",
					"deletion_protection": "false",
					"maintenance_window":  []map[string]string{{"enable": "true", "maintenance_time": "03:00:00Z", "duration": "3h", "weekly_period": "Thursday"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// worker args
						"worker_number":                  "3",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_size":               "100",
						"worker_disk_performance_level":  "PL1",
						"worker_disk_snapshot_policy_id": CHECKSET,
						"worker_data_disks.#":            "1",
						// global args
						"new_nat_gateway":      "true",
						"name":                 "tf-managed-k8s",
						"deletion_protection":  "false",
						"maintenance_window.#": "1",
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
	instance_type_family = "ecs.c6"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_keys" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}"
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

resource "alicloud_snapshot_policy" "default" {
	name            = "${var.name}"
	repeat_weekdays = ["1", "2", "3"]
	retention_days  = -1
	time_points     = ["1", "22", "23"]
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
