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
					"deletion_protection":          "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"worker_number":                "1",
						"password":                     "Test12345",
						"pod_cidr":                     "172.20.0.0/16",
						"service_cidr":                 "172.21.0.0/20",
						"slb_internet_enabled":         "true",
						"is_enterprise_security_group": "true",
						"deletion_protection":          "false",
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
					"worker_vswitch_ids", "proxy_mode", "worker_data_disks", "is_enterprise_security_group"},
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
