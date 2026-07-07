package alicloud

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const EdgeKubernetesConfigTpl = `
variable "name" {
  default = "%s"
}

variable "instance_type" {
  default = "ecs.c6.xlarge"
}

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

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_vpc" "vpc" {
  count      = 1
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  count        = 1
  vpc_id       = alicloud_vpc.vpc.0.id
  cidr_block   = format("192.168.%%d.0/24", count.index + 1)
  zone_id      = data.alicloud_db_zones.default.zones[count.index].id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.0.id
}

resource "alicloud_log_project" "log" {
  project_name = var.name
  description  = "created by terraform for edgekubernetes cluster"
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.0.min
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.vswitches.0.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_snapshot_policy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["1", "2", "3"]
  retention_days            = -1
  time_points               = ["1", "22", "23"]
}
`

var edgeCheckMap = map[string]string{
	"new_nat_gateway":       "true",
	"worker_number":         "1",
	"slb_internet_enabled":  "true",
	"install_cloud_monitor": "true",
}

func TestAccAliCloudEdgeKubernetes_basic(t *testing.T) {
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      resourceAttrCheck.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                        name,
					"worker_vswitch_ids":          []string{"${local.vswitch_id}"},
					"worker_instance_types":       []string{"${var.instance_type}"},
					"version":                     "1.20.11-aliyunedge.1",
					"worker_number":               "2",
					"password":                    "Test12345",
					"pod_cidr":                    "10.100.0.0/16",
					"service_cidr":                "172.30.0.0/16",
					"worker_instance_charge_type": "PostPaid",
					"new_nat_gateway":             "false",
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
					"is_enterprise_security_group":   "true",
					"deletion_protection":            "true",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"rds_instances":                  []string{"${alicloud_db_instance.default.id}"},
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"version":                        "1.20.11-aliyunedge.1",
						"worker_number":                  "2",
						"password":                       "Test12345",
						"pod_cidr":                       "10.100.0.0/16",
						"service_cidr":                   "172.30.0.0/16",
						"new_nat_gateway":                "false",
						"slb_internet_enabled":           "true",
						"deletion_protection":            "true",
						"resource_group_id":              CHECKSET,
						"rds_instances.#":                "1",
						"worker_data_disks.#":            "1",
						"worker_data_disks.0.category":   "cloud_ssd",
						"worker_data_disks.0.size":       "200",
						"worker_data_disks.0.encrypted":  "false",
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password", "install_cloud_monitor", "node_cidr_mask", "slb_internet_enabled", "worker_disk_category", "worker_disk_size", "worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number", "worker_vswitch_ids", "proxy_mode", "is_enterprise_security_group", "rds_instances", "worker_data_disks"},
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
					"version": "1.22.15-aliyunedge.1",
					"name":    name + "_update_again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.22.15-aliyunedge.1",
						"name":    name + "_update_again",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEdgeKubernetes_essd(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.EssdSupportRegions)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      resourceAttrCheck.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// global args
					"name":                         name,
					"version":                      "1.26.3-aliyun.1",
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"node_cidr_mask":               "24",
					"install_cloud_monitor":        "true",
					"slb_internet_enabled":         "true",
					"new_nat_gateway":              "true",
					"is_enterprise_security_group": "true",
					"deletion_protection":          "true",
					"pod_cidr":                     "10.101.0.0/16",
					"service_cidr":                 "172.30.0.0/16",
					// worker args
					"password":                       "Test12345",
					"worker_number":                  "1",
					"worker_vswitch_ids":             []string{"${local.vswitch_id}"},
					"worker_instance_types":          []string{"${var.instance_type}"},
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
						"version":              "1.26.3-aliyun.1",
						"resource_group_id":    CHECKSET,
						"slb_internet_enabled": "true",
						"deletion_protection":  "true",
						"pod_cidr":             "10.101.0.0/16",
						"service_cidr":         "172.30.0.0/16",
						// check worker args
						"password": "Test12345",
						//"worker_number":                  "1",
						"worker_disk_snapshot_policy_id": CHECKSET,
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
					"worker_number":                  "2",
					"worker_vswitch_ids":             []string{"${local.vswitch_id}"},
					"worker_instance_types":          []string{"${var.instance_type}"},
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
						"worker_number":                  "1",
						"worker_vswitch_ids.#":           "1",
						"worker_instance_types.#":        "1",
						"worker_disk_category":           "cloud_essd",
						"worker_disk_snapshot_policy_id": CHECKSET,
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      resourceAttrCheck.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                        name,
					"cluster_spec":                "ack.pro.small",
					"worker_vswitch_ids":          []string{"${local.vswitch_id}"},
					"worker_instance_types":       []string{"${var.instance_type}"},
					"version":                     "1.22.15-aliyunedge.1",
					"worker_number":               "1",
					"password":                    "Test12345",
					"pod_cidr":                    "10.102.0.0/16",
					"service_cidr":                "172.30.0.0/16",
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
					"runtime": []map[string]interface{}{
						{
							"name":    "containerd",
							"version": "1.6.28",
						},
					},
					"load_balancer_spec":             "slb.s2.small",
					"is_enterprise_security_group":   "true",
					"deletion_protection":            "false",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"cluster_spec":                   "ack.pro.small",
						"version":                        "1.22.15-aliyunedge.1",
						"worker_number":                  "1",
						"password":                       "Test12345",
						"pod_cidr":                       "10.102.0.0/16",
						"service_cidr":                   "172.30.0.0/16",
						"slb_internet_enabled":           "true",
						"deletion_protection":            "false",
						"resource_group_id":              CHECKSET,
						"worker_data_disks.#":            "1",
						"worker_data_disks.0.category":   "cloud_ssd",
						"worker_data_disks.0.size":       "200",
						"worker_data_disks.0.encrypted":  "false",
						"runtime.0.name":                 "containerd",
						"runtime.0.version":              "1.6.28",
						"load_balancer_spec":             "slb.s2.small",
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "name_prefix", "new_nat_gateway", "pod_cidr", "service_cidr", "password", "install_cloud_monitor", "node_cidr_mask", "slb_internet_enabled", "worker_disk_category", "worker_disk_size", "worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number", "worker_vswitch_ids", "proxy_mode", "is_enterprise_security_group", "rds_instances", "worker_data_disks", "load_balancer_spec", "runtime"},
			},
		},
	})
}

func edgeKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(EdgeKubernetesConfigTpl, name)
}

func TestCSEdgeKubernetesStateUpgradeV0(t *testing.T) {
	fields := []string{"runtime", "certificate_authority", "connections"}
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string][]interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id": "c-edge-123",
				"runtime": map[string]interface{}{
					"name":    "containerd",
					"version": "1.6.20",
				},
				"certificate_authority": map[string]interface{}{
					"cluster_cert": "cluster-cert",
					"client_cert":  "client-cert",
					"client_key":   "client-key",
				},
				"connections": map[string]interface{}{
					"api_server_internet": "https://1.2.3.4:6443",
					"api_server_intranet": "https://10.0.0.1:6443",
					"master_public_ip":    "1.2.3.4",
					"service_domain":      "*.c-edge-123.cs.local",
				},
			},
			expected: map[string][]interface{}{
				"runtime": {
					map[string]interface{}{
						"name":    "containerd",
						"version": "1.6.20",
					},
				},
				"certificate_authority": {
					map[string]interface{}{
						"cluster_cert": "cluster-cert",
						"client_cert":  "client-cert",
						"client_key":   "client-key",
					},
				},
				"connections": {
					map[string]interface{}{
						"api_server_internet": "https://1.2.3.4:6443",
						"api_server_intranet": "https://10.0.0.1:6443",
						"master_public_ip":    "1.2.3.4",
						"service_domain":      "*.c-edge-123.cs.local",
					},
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":                    "c-edge-456",
				"runtime":               map[string]interface{}{},
				"certificate_authority": map[string]interface{}{},
				"connections":           map[string]interface{}{},
			},
			expected: map[string][]interface{}{
				"runtime":               {},
				"certificate_authority": {},
				"connections":           {},
			},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":                    "c-edge-789",
				"runtime":               nil,
				"certificate_authority": nil,
				"connections":           nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "c-edge-000",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceAlicloudCSEdgeKubernetesStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, field := range fields {
				got := result[field]
				if tc.expected == nil {
					if got != nil {
						t.Errorf("%s: expected nil, got %v", field, got)
					}
					continue
				}
				exp := tc.expected[field]
				gotList, ok := got.([]interface{})
				if !ok {
					t.Fatalf("%s: expected []interface{}, got %T", field, got)
				}
				if len(gotList) != len(exp) {
					t.Fatalf("%s: expected length %d, got %d", field, len(exp), len(gotList))
				}
				for i, item := range gotList {
					gotMap := item.(map[string]interface{})
					expMap := exp[i].(map[string]interface{})
					for k, v := range expMap {
						if gotMap[k] != v {
							t.Errorf("%s[%d].%s: expected %v, got %v", field, i, k, v, gotMap[k])
						}
					}
				}
			}
		})
	}
}

func TestCSEdgeKubernetesSchemaVersionV0ToV1(t *testing.T) {
	r := resourceAlicloudCSEdgeKubernetes()
	if r.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion 1, got %d", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 StateUpgrader, got %d", len(r.StateUpgraders))
	}
	if r.StateUpgraders[0].Version != 0 {
		t.Errorf("expected StateUpgrader version 0, got %d", r.StateUpgraders[0].Version)
	}
}

func TestAccAlicloudCSEdgeKubernetes_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_cs_edge_kubernetes.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCsManagedKubernetes")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccCSEdgeKubernetesStateMigrationConfigV0(name),
			},
			{
				ProviderFactories: testAccProviderFactory,
				Config:            testAccCSEdgeKubernetesStateMigrationConfigV1(name),
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
  available_disk_category     = "cloud_efficiency"
  available_instance_type     = "ecs.c6.xlarge"
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
  worker_instance_types          = ["ecs.c6.xlarge"]
  worker_number                  = 1
  password                       = "Test12345"
  pod_cidr                       = "10.76.0.0/16"
  service_cidr                   = "172.26.0.0/16"
  worker_disk_category           = "cloud_efficiency"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
  runtime = {
    name    = "containerd"
    version = "1.6.20"
  }
}

output "runtime" {
  value = alicloud_cs_edge_kubernetes.default.runtime["name"]
}

output "alicloud_cs_edge_kubernetes" {
  value = alicloud_cs_edge_kubernetes.default.certificate_authority["cluster_cert"]
}

output "connections" {
  value = alicloud_cs_edge_kubernetes.default.connections["api_server_internet"]
}
`
}

func testAccCSEdgeKubernetesStateMigrationConfigV1(name string) string {
	return testAccCSEdgeKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_edge_kubernetes" "default" {
  name                           = var.name
  worker_vswitch_ids             = [alicloud_vswitch.default.id]
  worker_instance_types          = ["ecs.c6.xlarge"]
  worker_number                  = 1
  password                       = "Test12345"
  pod_cidr                       = "10.76.0.0/16"
  service_cidr                   = "172.26.0.0/16"
  worker_disk_category           = "cloud_efficiency"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
  runtime {
    name    = "containerd"
    version = "1.6.20"
  }
}

output "runtime" {
  value = alicloud_cs_edge_kubernetes.default.runtime.0.name
}

output "alicloud_cs_edge_kubernetes" {
  value = alicloud_cs_edge_kubernetes.default.certificate_authority.0.cluster_cert
}

output "connections" {
  value = alicloud_cs_edge_kubernetes.default.connections.0.api_server_internet
}
`
}
