package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
)

const EdgeKubernetesConfigTpl = `
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
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.default.id
  monitoring_period    = "60"
}

`

var edgeCheckMap = map[string]string{
	"new_nat_gateway":              "true",
	"worker_number":                "1",
	"slb_internet_enabled":         "true",
	"install_cloud_monitor":        "true",
	"is_enterprise_security_group": "true",
}

func TestAccAlicloudEdgeKubernetes(t *testing.T) {
	var cluster *cs.KubernetesClusterDetail
	resourceId := "alicloud_cs_edge_kubernetes.default"
	resourceAttr := resourceAttrInit(resourceId, edgeCheckMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	resourceCheck := resourceCheckInitWithDescribeMethod(resourceId, &cluster, serviceFunc, "DescribeCsManagedKubernetes")
	resourceAttrCheck := resourceAttrCheckInit(resourceCheck, resourceAttr)

	testAccCheck := resourceAttrCheck.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccedgekubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, edgeKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  resourceAttrCheck.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                        name,
					"worker_vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
					"worker_instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"version":                     "1.12.6-aliyunedge.2",
					"worker_number":               "1",
					"password":                    "Test12345",
					"pod_cidr":                    "172.20.0.0/16",
					"service_cidr":                "172.21.0.0/20",
					"worker_instance_charge_type": "PostPaid",
					"new_nat_gateway":             "true",
					"node_cidr_mask":              "24",
					"install_cloud_monitor":       "true",
					"slb_internet_enabled":        "true",
					"worker_data_disks": []map[string]string{
						{
							"category":  "cloud_ssd",
							"size":      "200",
							"encrypted": "false",
						},
					},
					"is_enterprise_security_group": "true",
					"deletion_protection":          "true",
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"rds_instances":                []string{"${alicloud_db_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"version":                       "1.12.6-aliyunedge.2",
						"worker_number":                 "1",
						"password":                      "Test12345",
						"pod_cidr":                      "172.20.0.0/16",
						"service_cidr":                  "172.21.0.0/20",
						"slb_internet_enabled":          "true",
						"is_enterprise_security_group":  "true",
						"deletion_protection":           "true",
						"resource_group_id":             CHECKSET,
						"rds_instances.#":               "1",
						"worker_data_disks.#":           "1",
						"worker_data_disks.0.category":  "cloud_ssd",
						"worker_data_disks.0.size":      "200",
						"worker_data_disks.0.encrypted": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password",
					"install_cloud_monitor", "force_update", "node_cidr_mask", "slb_internet_enabled", "worker_disk_category",
					"worker_disk_size", "worker_instance_charge_type", "worker_instance_types", "log_config",
					"worker_vswitch_ids", "proxy_mode", "is_enterprise_security_group", "rds_instances", "worker_data_disks"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "modified-edge-cluster",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "modified-edge-cluster",
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
					"version":       "1.14.8-aliyunedge.1",
					"name":          "modified-edge-cluster-again",
					"worker_number": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":       "1.14.8-aliyunedge.1",
						"name":          "modified-edge-cluster-again",
						"worker_number": "3",
					}),
				),
			},
		},
	})
}

func edgeKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(EdgeKubernetesConfigTpl, name)
}
