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
					"worker_vswitch_ids":          []string{"${local.vswitch_id}"},
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
					"enable_rrsa":                 "true",
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
					"worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number",
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
					"name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
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
					"enable_rrsa": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_rrsa": "null",
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
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// cluster args
					"name":                name,
					"password":            "Test12345",
					"pod_cidr":            "172.20.0.0/16",
					"version":             "1.20.11-aliyun.1",
					"service_cidr":        "172.21.0.0/20",
					"deletion_protection": "true",
					"cluster_spec":        "ack.standard",
					// worker args
					"worker_number":                  "2",
					"worker_vswitch_ids":             []string{"${local.vswitch_id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_instance_charge_type":    "PostPaid",
					"worker_data_disk_category":      "cloud_ssd",
					"worker_data_disk_size":          "20",
					"worker_disk_category":           "cloud_essd",
					"worker_disk_size":               "100",
					"worker_disk_performance_level":  "PL0",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_data_disks": []map[string]string{
						{
							"category":                "cloud_essd",
							"size":                    "100",
							"auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
							"performance_level":       "PL0",
						},
					},
					"tags": map[string]string{
						"Platform": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// cluster args
						"name":                name,
						"version":             "1.20.11-aliyun.1",
						"password":            "Test12345",
						"pod_cidr":            "172.20.0.0/16",
						"service_cidr":        "172.21.0.0/20",
						"deletion_protection": "true",
						"cluster_spec":        "ack.standard",
						// worker args
						"worker_number":                  "2",
						"worker_data_disk_category":      "cloud_ssd",
						"worker_data_disk_size":          "20",
						"worker_instance_charge_type":    "PostPaid",
						"worker_disk_size":               "100",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_performance_level":  "PL0",
						"worker_disk_snapshot_policy_id": CHECKSET,
						"worker_data_disks.#":            "1",
						"tags.%":                         "1",
						"tags.Platform":                  "TF",
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
					"worker_instance_charge_type", "worker_instance_types", "log_config", "tags", "worker_data_disk_category", "worker_data_disk_size",
					"master_vswitch_ids", "worker_number", "worker_vswitch_ids", "exclude_autoscaler_nodes", "cpu_policy", "proxy_mode", "cluster_domain",
					"custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "worker_disk_snapshot_policy_id", "worker_disk_performance_level",
					"taints", "encryption_provider_key", "worker_data_disks", "rds_instances", "load_balancer_spec", "worker_number"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":        "2",
						"tags.Platform": "TF",
						"tags.Env":      "Pre",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					// cluster args
					"deletion_protection": "false",
					"cluster_spec":        "ack.pro.small", // migrate cluster
					// worker args
					"worker_number":                  "3",
					"worker_vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_disk_category":           "cloud_essd",
					"worker_disk_size":               "120",
					"worker_disk_performance_level":  "PL1",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_data_disks": []map[string]string{
						{
							"category":                "cloud_essd",
							"size":                    "120",
							"auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
							"performance_level":       "PL1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// global args
						"deletion_protection": "false",
						"cluster_spec":        "ack.pro.small",
						// worker args
						"worker_number":                  "3",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_size":               "120",
						"worker_disk_performance_level":  "PL1",
						"worker_disk_snapshot_policy_id": CHECKSET,
						"worker_data_disks.#":            "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCSManagedKubernetes_controlPlanLog(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfig)

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
					"name":                         name,
					"cluster_spec":                 "ack.pro.small",
					"is_enterprise_security_group": "true",
					"deletion_protection":          "false",
					"node_cidr_mask":               "26",
					"pod_cidr":                     "172.20.0.0/16",
					"service_cidr":                 "172.21.0.0/20",
					"os_type":                      "Linux",
					"platform":                     "AliyunLinux",
					"password":                     "Test12345",
					"worker_number":                "0",
					"worker_vswitch_ids":           []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_disk_size":             "50",
					"worker_disk_category":         "cloud_ssd",
					"control_plane_log_ttl":        "30",
					"control_plane_log_components": []string{"apiserver", "kcm", "scheduler"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"cluster_spec":         "ack.pro.small",
						"deletion_protection":  "false",
						"pod_cidr":             "172.20.0.0/16",
						"service_cidr":         "172.21.0.0/20",
						"os_type":              "Linux",
						"platform":             "AliyunLinux",
						"password":             "Test12345",
						"worker_number":        "0",
						"worker_disk_size":     "50",
						"worker_disk_category": "cloud_ssd",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr", "service_cidr", "control_plane_log_ttl",
					"node_cidr_mask", "vswitch_ids", "worker_disk_category", "worker_disk_size", "control_plane_log_components",
					"worker_instance_charge_type", "worker_instance_types", "os_type", "platform", "timezone", "password",
					"exclude_autoscaler_nodes", "install_cloud_monitor", "proxy_mode", "slb_internet_enabled", "worker_vswitch_ids",
					"cpu_policy", "enable_ssh", "is_enterprise_security_group", "worker_number",
				},
			},
		},
	})
}

func resourceCSManagedKubernetesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 4
	memory_size = 8
	kubernetes_node_role = "Worker"
}
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
`, name)
}

func resourceCSManagedKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 4
	memory_size = 8
	kubernetes_node_role = "Worker"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_keys" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
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
  vswitch_id           = "${local.vswitch_id}"
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
					"worker_vswitch_ids":          []string{"${local.vswitch_id}"},
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
					"version":                     "1.20.11-aliyun.1",
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
						"version":                   "1.20.11-aliyun.1",
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
					"version": "1.22.3-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.22.3-aliyun.1",
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
