package cs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const EdgeKubernetesCommonConfigTpl = `
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_cs_kubernetes_version" "edge" {
  cluster_type = "ManagedKubernetes"
  profile      = "Edge"
}
`

const EdgeKubernetesBasicConfigTpl = `
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_instance_types" "default" {
  availability_zone     = data.alicloud_db_zones.default.zones.0.id
  instance_charge_type  = "PostPaid"
  kubernetes_node_role  = "Worker"
  system_disk_category  = "cloud_essd"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.0.min
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.vswitches.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.vpc.id
  nat_gateway_name = var.name
  vswitch_id       = alicloud_vswitch.vswitches.id
  nat_type         = "Enhanced"
}

resource "alicloud_eip_address" "default" {
  address_name         = var.name
  payment_type         = "PayAsYouGo"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_snat_entry" "default" {
  snat_table_id    = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vswitches.id
  snat_ip           = alicloud_eip_address.default.ip_address
  snat_entry_name   = var.name

  depends_on = [alicloud_eip_association.default]
}

locals {
  # Make cluster creation wait until outbound access is ready while testing
  # new_nat_gateway = false.
  vswitch_id = alicloud_snat_entry.default.source_vswitch_id
}
`

const EdgeKubernetesEssdConfigTpl = `
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
  available_instance_type     = "ecs.c6.xlarge"
}

data "alicloud_instance_types" "default" {
  availability_zone     = data.alicloud_zones.default.zones.0.id
  instance_charge_type  = "PostPaid"
  kubernetes_node_role  = "Worker"
  system_disk_category  = "cloud_essd"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.id
}

resource "alicloud_snapshot_policy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["1", "2", "3"]
  retention_days            = -1
  time_points               = ["1", "22", "23"]
}
`

const EdgeKubernetesProConfigTpl = `
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
  available_instance_type     = "ecs.c6.xlarge"
}

data "alicloud_instance_types" "default" {
  availability_zone     = data.alicloud_zones.default.zones.0.id
  instance_charge_type  = "PostPaid"
  kubernetes_node_role  = "Worker"
  system_disk_category  = "cloud_essd"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.id
}
`

var edgeCheckMap = map[string]string{
	"new_nat_gateway":       "true",
	"worker_number":         "2",
	"slb_internet_enabled":  "true",
	"install_cloud_monitor": "true",
}

func TestAccAliCloudEdgeKubernetes_basic(t *testing.T) {
	var cluster *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_edge_kubernetes.default"
	resourceAttr := acctest.ResourceAttrInit(resourceId, edgeCheckMap)
	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	resourceCheck := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &cluster, serviceFunc, "DescribeCsManagedKubernetes")
	resourceAttrCheck := acctest.ResourceAttrCheckInit(resourceCheck, resourceAttr)

	testAccCheck := resourceAttrCheck.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccedgekubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, edgeKubernetesBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             resourceAttrCheck.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"worker_vswitch_ids":            []string{"${local.vswitch_id}"},
					"worker_instance_types":         []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"version":                       "${data.alicloud_cs_kubernetes_version.edge.metadata.0.version}",
					"worker_number":                 "2",
					"password":                      "Test12345",
					"pod_cidr":                      "10.100.0.0/16",
					"service_cidr":                  "172.30.0.0/16",
					"worker_instance_charge_type":   "PostPaid",
					"worker_disk_category":          "cloud_essd",
					"worker_disk_performance_level": "PL0",
					"new_nat_gateway":               "false",
					"node_cidr_mask":                "24",
					"install_cloud_monitor":         "true",
					"slb_internet_enabled":          "true",
					"worker_data_disks": []map[string]string{
						{
							"category":          "cloud_essd",
							"size":              "200",
							"encrypted":         "true",
							"performance_level": "PL0",
						},
					},
					"is_enterprise_security_group":   "false",
					"deletion_protection":            "false",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"rds_instances":                  []string{"${alicloud_db_instance.default.id}"},
					"skip_set_certificate_authority": "false",
					"log_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "tfacc-log-project",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
						"version":                               acctest.CHECKSET,
						"worker_number":                         "2",
						"worker_disk_category":                  "cloud_essd",
						"worker_disk_performance_level":         "PL0",
						"password":                              "Test12345",
						"pod_cidr":                              "10.100.0.0/16",
						"service_cidr":                          "172.30.0.0/16",
						"new_nat_gateway":                       "false",
						"slb_internet_enabled":                  "true",
						"deletion_protection":                   "false",
						"resource_group_id":                     acctest.CHECKSET,
						"rds_instances.#":                       "1",
						"worker_data_disks.#":                   "1",
						"worker_data_disks.0.category":          "cloud_essd",
						"worker_data_disks.0.size":              "200",
						"worker_data_disks.0.encrypted":         "true",
						"worker_data_disks.0.performance_level": "PL0",
						"skip_set_certificate_authority":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]interface{}{
						{
							"name":     "tfacc-addon",
							"config":   "{\"enable\":\"true\"}",
							"version":  "1.0.0",
							"disabled": "true",
						},
					},
					"availability_zone": "cn-hangzhou-g",
					"proxy_mode":        "iptables",
					"retain_resources":  []string{"tfacc-retain"},
					"user_data":         "I2RhdGE6IHRmYWNjCg==",
					"security_group_id": "sg-tfacc-placeholder",
					"resource_group_id": "rg-tfacc-placeholder",
					"cluster_spec":      "ack.standard",
					"runtime": []map[string]string{
						{
							"name":    "containerd",
							"version": "1.6.20",
						},
					},
					"version":     "1.28.3-aliyun-1",
					"password":    acctest.REMOVEKEY,
					"key_name":    "tfacc-key-pair",
					"name":        acctest.REMOVEKEY,
					"name_prefix": "tfacc-cluster-",
					"worker_data_disks": []map[string]string{
						{
							"device":                  "/dev/vdc",
							"kms_key_id":              "key-tfacc-placeholder",
							"name":                    "tfacc-disk",
							"snapshot_id":             "s-tfacc-placeholder",
							"auto_snapshot_policy_id": "sp-tfacc-placeholder",
							"category":                "cloud_ssd",
							"encrypted":               "false",
						},
					},
					"worker_number": acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password", "install_cloud_monitor", "node_cidr_mask", "slb_internet_enabled", "worker_disk_category", "worker_disk_size", "worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number", "worker_vswitch_ids", "proxy_mode", "worker_disk_performance_level", "is_enterprise_security_group", "rds_instances", "worker_data_disks"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_update",
					"log_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "tfacc-log-project-updated",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
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
					"name": name + "_update_again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": acctest.CHECKSET,
						"name":    name + "_update_again",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEdgeKubernetes_essd(t *testing.T) {
	var cluster *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_edge_kubernetes.default"
	resourceAttr := acctest.ResourceAttrInit(resourceId, edgeCheckMap)
	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	resourceCheck := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &cluster, serviceFunc, "DescribeCsManagedKubernetes")
	resourceAttrCheck := acctest.ResourceAttrCheckInit(resourceCheck, resourceAttr)

	testAccCheck := resourceAttrCheck.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccedgekubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, edgeKubernetesEssdConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             resourceAttrCheck.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// global args
					"name":                         name,
					"version":                      "${data.alicloud_cs_kubernetes_version.edge.metadata.0.version}",
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"node_cidr_mask":               "24",
					"install_cloud_monitor":        "true",
					"slb_internet_enabled":         "true",
					"new_nat_gateway":              "true",
					"is_enterprise_security_group": "false",
					"deletion_protection":          "false",
					"pod_cidr":                     "10.101.0.0/16",
					"service_cidr":                 "172.30.0.0/16",
					// worker args
					"password":                       "Test12345",
					"worker_number":                  "2",
					"worker_vswitch_ids":             []string{"${local.vswitch_id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_instance_charge_type":    "PostPaid",
					"worker_disk_category":           "cloud_essd",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_disk_size":               "120",
					"worker_disk_performance_level":  "PL0",
					"worker_data_disks": []map[string]string{
						{
							"category":                "cloud_essd",
							"auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
							"size":                    "100",
							"performance_level":       "PL0",
						},
					},
					"tags": map[string]string{
						"Platform": "TF",
					},
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// check global args
						"name":                 name,
						"version":              acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"slb_internet_enabled": "true",
						"deletion_protection":  "false",
						"pod_cidr":             "10.101.0.0/16",
						"service_cidr":         "172.30.0.0/16",
						// check worker args
						"password":                       "Test12345",
						"worker_number":                  "2",
						"worker_disk_snapshot_policy_id": acctest.CHECKSET,
						"worker_disk_size":               "120",
						"worker_disk_performance_level":  "PL0",
						"worker_data_disks.#":            "1",
						"tags.%":                         "1",
						"tags.Platform":                  "TF",
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password", "worker_number", "install_cloud_monitor", "node_cidr_mask", "slb_internet_enabled", "tags", "worker_disk_category", "worker_disk_size", "worker_instance_charge_type", "worker_disk_snapshot_policy_id", "worker_instance_types", "log_config", "worker_vswitch_ids", "proxy_mode", "worker_disk_performance_level", "is_enterprise_security_group", "rds_instances", "worker_data_disks"},
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
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
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
					"deletion_protection":            "false",
					"worker_number":                  "3",
					"worker_vswitch_ids":             []string{"${local.vswitch_id}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_disk_category":           "cloud_essd",
					"worker_disk_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
					"worker_disk_size":               "120",
					"worker_disk_performance_level":  "PL1",
					"worker_data_disks": []map[string]string{
						{
							"category":                "cloud_essd",
							"auto_snapshot_policy_id": "${alicloud_snapshot_policy.default.id}",
							"size":                    "120",
							"performance_level":       "PL1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":            "false",
						"worker_number":                  "3",
						"worker_vswitch_ids.#":           "1",
						"worker_instance_types.#":        "1",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_snapshot_policy_id": acctest.CHECKSET,
						"worker_disk_size":               "120",
						"worker_disk_performance_level":  "PL1",
						"worker_data_disks.#":            "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEdgeKubernetes_pro(t *testing.T) {
	var cluster *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_edge_kubernetes.default"
	resourceAttr := acctest.ResourceAttrInit(resourceId, edgeCheckMap)
	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	resourceCheck := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &cluster, serviceFunc, "DescribeCsManagedKubernetes")
	resourceAttrCheck := acctest.ResourceAttrCheckInit(resourceCheck, resourceAttr)

	testAccCheck := resourceAttrCheck.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccedgekubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, edgeKubernetesProConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             resourceAttrCheck.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"cluster_spec":                  "ack.pro.small",
					"worker_vswitch_ids":            []string{"${local.vswitch_id}"},
					"worker_instance_types":         []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"version":                       "${data.alicloud_cs_kubernetes_version.edge.metadata.0.version}",
					"worker_number":                 "2",
					"password":                      "Test12345",
					"pod_cidr":                      "10.102.0.0/16",
					"service_cidr":                  "172.30.0.0/16",
					"worker_instance_charge_type":   "PostPaid",
					"worker_disk_category":          "cloud_essd",
					"worker_disk_performance_level": "PL0",
					"new_nat_gateway":               "true",
					"node_cidr_mask":                "24",
					"install_cloud_monitor":         "true",
					"slb_internet_enabled":          "true",
					"worker_data_disks": []map[string]string{
						{
							"category":          "cloud_essd",
							"size":              "200",
							"encrypted":         "true",
							"performance_level": "PL0",
						},
					},
					"runtime": []map[string]interface{}{
						{
							"name":    "${data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.name}",
							"version": "${data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.version}",
						},
					},
					"load_balancer_spec":             "slb.s2.small",
					"is_enterprise_security_group":   "false",
					"deletion_protection":            "false",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
						"cluster_spec":                          "ack.pro.small",
						"version":                               acctest.CHECKSET,
						"worker_number":                         "2",
						"worker_disk_category":                  "cloud_essd",
						"worker_disk_performance_level":         "PL0",
						"password":                              "Test12345",
						"pod_cidr":                              "10.102.0.0/16",
						"service_cidr":                          "172.30.0.0/16",
						"slb_internet_enabled":                  "true",
						"deletion_protection":                   "false",
						"resource_group_id":                     acctest.CHECKSET,
						"worker_data_disks.#":                   "1",
						"worker_data_disks.0.category":          "cloud_essd",
						"worker_data_disks.0.size":              "200",
						"worker_data_disks.0.encrypted":         "true",
						"worker_data_disks.0.performance_level": "PL0",
						"runtime.0.name":                        acctest.CHECKSET,
						"runtime.0.version":                     acctest.CHECKSET,
						"load_balancer_spec":                    "slb.s2.small",
						"skip_set_certificate_authority":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "tfacc-log-project-pro",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password", "install_cloud_monitor", "node_cidr_mask", "slb_internet_enabled", "worker_disk_category", "worker_disk_size", "worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number", "worker_vswitch_ids", "proxy_mode", "worker_disk_performance_level", "is_enterprise_security_group", "rds_instances", "worker_data_disks", "load_balancer_spec", "runtime"},
			},
		},
	})
}

func edgeKubernetesBasicConfigDependence(name string) string {
	return fmt.Sprintf(EdgeKubernetesCommonConfigTpl+EdgeKubernetesBasicConfigTpl, name)
}

func edgeKubernetesEssdConfigDependence(name string) string {
	return fmt.Sprintf(EdgeKubernetesCommonConfigTpl+EdgeKubernetesEssdConfigTpl, name)
}

func edgeKubernetesProConfigDependence(name string) string {
	return fmt.Sprintf(EdgeKubernetesCommonConfigTpl+EdgeKubernetesProConfigTpl, name)
}

func TestAccAliCloudEdgeKubernetes_StateMigrationV0ToV1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_edge_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCsManagedKubernetes")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccCSEdgeKubernetesStateMigrationConfigV0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "worker_number", "2"),
					resource.TestCheckResourceAttr(resourceId, "worker_disk_category", "cloud_essd"),
					resource.TestCheckResourceAttr(resourceId, "worker_disk_performance_level", "PL0"),
					resource.TestCheckResourceAttrSet(resourceId, "runtime.name"),
					resource.TestCheckResourceAttrSet(resourceId, "runtime.version"),
					resource.TestCheckResourceAttrSet(resourceId, "certificate_authority.cluster_cert"),
					resource.TestCheckResourceAttrSet(resourceId, "connections.api_server_internet"),
				),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
				Config:                   testAccCSEdgeKubernetesStateMigrationConfigV1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "worker_number", "2"),
					resource.TestCheckResourceAttr(resourceId, "worker_disk_category", "cloud_essd"),
					resource.TestCheckResourceAttr(resourceId, "worker_disk_performance_level", "PL0"),
					resource.TestCheckResourceAttr(resourceId, "runtime.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "runtime.0.name"),
					resource.TestCheckResourceAttrSet(resourceId, "runtime.0.version"),
					resource.TestCheckResourceAttr(resourceId, "certificate_authority.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "certificate_authority.0.cluster_cert"),
					resource.TestCheckResourceAttr(resourceId, "connections.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "connections.0.api_server_internet"),
				),
			},
		},
	})
}

func testAccCSEdgeKubernetesStateMigrationDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
  available_instance_type     = "ecs.c6.xlarge"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_charge_type = "PostPaid"
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
}

data "alicloud_cs_kubernetes_version" "edge" {
  cluster_type = "ManagedKubernetes"
  profile      = "Edge"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
`, name)
}

func testAccCSEdgeKubernetesStateMigrationConfigV0(name string) string {
	return testAccCSEdgeKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_edge_kubernetes" "default" {
  name                           = var.name
  worker_vswitch_ids             = [alicloud_vswitch.default.id]
  worker_instance_types          = [data.alicloud_instance_types.default.instance_types.0.id]
  worker_number                  = 2
  password                       = "Test12345"
  pod_cidr                       = "10.76.0.0/16"
  service_cidr                   = "172.26.0.0/16"
  version                        = data.alicloud_cs_kubernetes_version.edge.metadata.0.version
  worker_disk_category           = "cloud_essd"
  worker_disk_performance_level  = "PL0"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  is_enterprise_security_group   = false
  deletion_protection            = false
  skip_set_certificate_authority = false
  runtime = {
    name    = data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.name
    version = data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.version
  }
}

output "runtime_name" {
  value = alicloud_cs_edge_kubernetes.default.runtime["name"]
}

output "certificate_authority_cluster_cert" {
  value = alicloud_cs_edge_kubernetes.default.certificate_authority["cluster_cert"]
}

output "api_server_internet" {
  value = alicloud_cs_edge_kubernetes.default.connections["api_server_internet"]
}
`
}

func testAccCSEdgeKubernetesStateMigrationConfigV1(name string) string {
	return testAccCSEdgeKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_edge_kubernetes" "default" {
  name                           = var.name
  worker_vswitch_ids             = [alicloud_vswitch.default.id]
  worker_instance_types          = [data.alicloud_instance_types.default.instance_types.0.id]
  worker_number                  = 2
  password                       = "Test12345"
  pod_cidr                       = "10.76.0.0/16"
  service_cidr                   = "172.26.0.0/16"
  version                        = data.alicloud_cs_kubernetes_version.edge.metadata.0.version
  worker_disk_category           = "cloud_essd"
  worker_disk_performance_level  = "PL0"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  is_enterprise_security_group   = false
  deletion_protection            = false
  skip_set_certificate_authority = false
  runtime {
    name    = data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.name
    version = data.alicloud_cs_kubernetes_version.edge.metadata.0.runtime.0.version
  }
}

output "runtime_name" {
  value = alicloud_cs_edge_kubernetes.default.runtime.0.name
}

output "certificate_authority_cluster_cert" {
  value = alicloud_cs_edge_kubernetes.default.certificate_authority.0.cluster_cert
}

output "api_server_internet" {
  value = alicloud_cs_edge_kubernetes.default.connections.0.api_server_internet
}
`
}

func edgeKubernetesCollectionOrderConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_cs_kubernetes_version" "edge" {
  cluster_type = "ManagedKubernetes"
  profile      = "Edge"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
}

data "alicloud_db_zones" "rds" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  db_instance_storage_type = "cloud_essd"
  category                 = "Basic"
}

data "alicloud_db_instance_classes" "rds" {
  zone_id                  = data.alicloud_db_zones.rds.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  db_instance_storage_type = "cloud_essd"
  category                 = "Basic"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  count        = 2
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = format("192.168.%%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.default.zones[count.index].id
  vswitch_name = var.name
}

resource "alicloud_vswitch" "rds" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.30.0/24"
  zone_id      = data.alicloud_db_zones.rds.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id_0       = alicloud_vswitch.vswitches.0.id
  vswitch_id_1       = alicloud_vswitch.vswitches.1.id
  rds_instance_class = [for c in data.alicloud_db_instance_classes.rds.instance_classes : c.instance_class if startswith(c.instance_class, "mysql")][0]
}

resource "alicloud_security_group" "group" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.vpc.id
}

resource "alicloud_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = local.rds_instance_class
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "tf-testaccedgekubernetes"
  vswitch_id           = alicloud_vswitch.rds.id
  monitoring_period    = "60"
}

resource "alicloud_db_instance" "extra" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = local.rds_instance_class
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "tf-testaccedgekubernetes-extra"
  vswitch_id           = alicloud_vswitch.rds.id
  monitoring_period    = "60"
}
`, name)
}

// TestAccAliCloudEdgeKubernetes_collectionOrder verifies that reordering
// TypeList members converges to an empty plan for the 7 configurable
// multi-member TypeList attributes of cs_edge_kubernetes:
// worker_vswitch_ids, worker_instance_types, worker_data_disks,
// rds_instances, retain_resources, addons and runtime.
//
// lintignore: AT001
func TestAccAliCloudEdgeKubernetes_collectionOrder(t *testing.T) {
	var cluster *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_edge_kubernetes.default"
	resourceAttr := acctest.ResourceAttrInit(resourceId, edgeCheckMap)
	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	resourceCheck := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &cluster, serviceFunc, "DescribeCsManagedKubernetes")
	resourceAttrCheck := acctest.ResourceAttrCheckInit(resourceCheck, resourceAttr)

	testAccCheck := resourceAttrCheck.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccedgekubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, edgeKubernetesCollectionOrderConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             resourceAttrCheck.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"worker_vswitch_ids":             []string{"${local.vswitch_id_0}", "${local.vswitch_id_1}"},
					"worker_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}"},
					"version":                        "${data.alicloud_cs_kubernetes_version.edge.metadata.0.version}",
					"worker_number":                  "2",
					"password":                       "Test12345",
					"pod_cidr":                       "10.103.0.0/16",
					"service_cidr":                   "172.30.0.0/16",
					"worker_instance_charge_type":    "PostPaid",
					"worker_disk_category":           "cloud_essd",
					"worker_disk_performance_level":  "PL0",
					"new_nat_gateway":                "true",
					"node_cidr_mask":                 "24",
					"install_cloud_monitor":          "true",
					"slb_internet_enabled":           "true",
					"is_enterprise_security_group":   "false",
					"deletion_protection":            "false",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"skip_set_certificate_authority": "false",
					"worker_data_disks": []map[string]string{
						{
							"category":          "cloud_essd",
							"size":              "200",
							"performance_level": "PL0",
						},
						{
							"category":          "cloud_essd",
							"size":              "100",
							"performance_level": "PL1",
						},
					},
					"addons": []map[string]string{
						{"name": "terway-eniip", "disabled": "false"},
						{"name": "metrics-server", "disabled": "false"},
					},
					"rds_instances":    []string{"${alicloud_db_instance.default.id}", "${alicloud_db_instance.extra.id}"},
					"retain_resources": []string{"${alicloud_security_group.group.id}", "${alicloud_vpc.vpc.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_vswitch_ids.#":    "2",
						"worker_instance_types.#": "2",
						"worker_data_disks.#":     "2",
						"addons.#":                "2",
						"rds_instances.#":         "2",
						"retain_resources.#":      "2",
					}),
				),
			},
			// worker_vswitch_ids: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_vswitch_ids": []string{"${local.vswitch_id_1}", "${local.vswitch_id_0}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// worker_instance_types: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_instance_types": []string{"${data.alicloud_instance_types.default.instance_types.1.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// worker_data_disks: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_data_disks": []map[string]string{
						{
							"category":          "cloud_essd",
							"size":              "100",
							"performance_level": "PL1",
						},
						{
							"category":          "cloud_essd",
							"size":              "200",
							"performance_level": "PL0",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// rds_instances: verify reorder convergence
			// rds_instances assertions observe retained state, not cloud
			// updates; stock Read does not independently refresh this field.
			{
				Config: testAccConfig(map[string]interface{}{
					"rds_instances": []string{"${alicloud_db_instance.extra.id}", "${alicloud_db_instance.default.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// retain_resources: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{"${alicloud_vpc.vpc.id}", "${alicloud_security_group.group.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// addons: verify reorder convergence
			// addons assertions observe retained state, not cloud updates;
			// stock Read does not independently refresh this field.
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{"name": "metrics-server", "disabled": "false"},
						{"name": "terway-eniip", "disabled": "false"},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// runtime: verify reorder convergence
			// runtime is Optional and not set in step 0; set 2 members first,
			// then reorder and verify convergence.  Runtime assertions observe
			// retained state, not cloud updates; stock Read does not
			// independently refresh this field.
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "docker", "version": "19.03.15"},
						{"name": "containerd", "version": "1.6.20"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "runtime.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "1.6.20"},
						{"name": "docker", "version": "19.03.15"},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "runtime.#", "2"),
				),
			},
		},
	})
}
