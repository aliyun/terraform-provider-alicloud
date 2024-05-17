package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCSManagedKubernetes_basic(t *testing.T) {
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

	clusterCaCertFile, clientCertFile, clientKeyFile, err := CreateTempFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clientCertFile.Name())
	defer os.Remove(clientKeyFile.Name())
	defer os.Remove(clusterCaCertFile.Name())

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
					"name":                    name,
					"version":                 "1.24.6-aliyun.1",
					"worker_vswitch_ids":      []string{"${local.vswitch_id}"},
					"pod_cidr":                "10.93.0.0/16",
					"service_cidr":            "172.21.0.0/16",
					"slb_internet_enabled":    "true",
					"load_balancer_spec":      "slb.s2.small",
					"cluster_spec":            "ack.pro.small",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"security_group_id":       "${alicloud_security_group.default.id}",
					"deletion_protection":     "false",
					"enable_rrsa":             "false",
					"timezone":                "Asia/Shanghai",
					"proxy_mode":              "ipvs",
					"new_nat_gateway":         "true",
					"api_audiences":           []string{"https://kubernetes.default.svc"},
					"service_account_issuer":  "https://kubernetes.default.svc",
					"cluster_domain":          "cluster.local",
					"custom_san":              "www.terraform.io",
					"encryption_provider_key": "${data.alicloud_kms_keys.default.keys.0.id}",
					"maintenance_window":      []map[string]string{{"enable": "true", "maintenance_time": "03:00:00Z", "duration": "3h", "weekly_period": "Thursday"}},
					"cluster_ca_cert":         clusterCaCertFile.Name(),
					"client_key":              clientKeyFile.Name(),
					"client_cert":             clientCertFile.Name(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
						"version":                               "1.24.6-aliyun.1",
						"pod_cidr":                              "10.93.0.0/16",
						"service_cidr":                          "172.21.0.0/16",
						"slb_internet_enabled":                  "true",
						"cluster_spec":                          "ack.pro.small",
						"resource_group_id":                     CHECKSET,
						"deletion_protection":                   "false",
						"enable_rrsa":                           "false",
						"timezone":                              "Asia/Shanghai",
						"proxy_mode":                            "ipvs",
						"new_nat_gateway":                       "true",
						"nat_gateway_id":                        CHECKSET,
						"cluster_domain":                        "cluster.local",
						"custom_san":                            "www.terraform.io",
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
				ImportStateVerifyIgnore: []string{"new_nat_gateway", "user_ca", "timezone", "name_prefix", "api_audiences",
					"service_account_issuer", "load_balancer_spec", "encryption_provider_key", "cluster_ca_cert", "client_key", "client_cert",
				},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name + "_update",
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name + "_update",
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
					"enable_rrsa": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_rrsa":             "true",
						"rrsa_metadata.0.enabled": "true",
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
					"maintenance_window": []map[string]string{{"enable": "false", "maintenance_time": "", "duration": "", "weekly_period": ""}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "false",
						"maintenance_window.0.maintenance_time": "",
						"maintenance_window.0.duration":         "",
						"maintenance_window.0.weekly_period":    "",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSManagedKubernetes_essd_migrate_upgrade(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfigDependence_essd)

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
					"version":             "1.24.6-aliyun.1",
					"pod_cidr":            "10.94.0.0/16",
					"service_cidr":        "172.22.0.0/16",
					"deletion_protection": "false",
					"cluster_spec":        "ack.standard",
					"new_nat_gateway":     "true",
					"proxy_mode":          "ipvs",
					"worker_vswitch_ids":  []string{"${local.vswitch_id}"},
					"tags": map[string]string{
						"Platform": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// cluster args
						"name":                name,
						"version":             "1.24.6-aliyun.1",
						"pod_cidr":            "10.94.0.0/16",
						"service_cidr":        "172.22.0.0/16",
						"deletion_protection": "false",
						"cluster_spec":        "ack.standard",
						"new_nat_gateway":     "true",
						"nat_gateway_id":      CHECKSET,
						"proxy_mode":          "ipvs",
						"tags.%":              "1",
						"tags.Platform":       "TF",
					}),
				),
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
					// migrate cluster
					"cluster_spec": "ack.pro.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_spec": "ack.pro.small",
					}),
				),
			},
			{
				// upgrade
				Config: testAccConfig(map[string]interface{}{
					"version": "1.26.15-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.26.15-aliyun.1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"new_nat_gateway", "user_ca", "timezone", "name_prefix"},
			},
		},
	})
}

func TestAccAliCloudCSManagedKubernetes_controlPlanLog(t *testing.T) {
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

	tmpCAFile, err := os.CreateTemp("", "tf-acc-alicloud-cs-managed-kubernetes-userca")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpCAFile.Name())
	err = os.WriteFile(tmpCAFile.Name(), []byte(caCert), 0644)
	if err != nil {
		t.Fatal(err)
	}

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
					"name_prefix":                  "tf-testaccmanagedkubernetes",
					"cluster_spec":                 "ack.pro.small",
					"is_enterprise_security_group": "true",
					"deletion_protection":          "false",
					"new_nat_gateway":              "true",
					"node_cidr_mask":               "26",
					"service_cidr":                 "172.23.0.0/16",
					"proxy_mode":                   "ipvs",
					"worker_vswitch_ids":           []string{"${local.vswitch_id}"},
					"pod_vswitch_ids":              []string{"${local.vswitch_id}"},
					"control_plane_log_ttl":        "30",
					"control_plane_log_components": []string{"apiserver", "kcm", "scheduler"},
					"control_plane_log_project":    "",
					"user_ca":                      tmpCAFile.Name(),
					"addons":                       []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           CHECKSET,
						"cluster_spec":                   "ack.pro.small",
						"deletion_protection":            "false",
						"new_nat_gateway":                "true",
						"nat_gateway_id":                 CHECKSET,
						"service_cidr":                   "172.23.0.0/16",
						"proxy_mode":                     "ipvs",
						"control_plane_log_ttl":          "30",
						"control_plane_log_components.0": "apiserver",
						"control_plane_log_components.1": "kcm",
						"control_plane_log_components.2": "scheduler",
						"control_plane_log_project":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"control_plane_log_ttl":        "90",
					"control_plane_log_components": []string{"apiserver", "kcm", "scheduler", "ccm"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_log_ttl":          "90",
						"control_plane_log_components.0": "apiserver",
						"control_plane_log_components.1": "kcm",
						"control_plane_log_components.2": "scheduler",
						"control_plane_log_components.3": "ccm",
					})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"control_plane_log_project": "${alicloud_log_project.log.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_log_project": name,
					})),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"new_nat_gateway", "user_ca", "timezone", "name_prefix", "addons",
					"is_enterprise_security_group"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]interface{}{
						{
							"delete_mode":   "delete",
							"resource_type": "SLB",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "ALB",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{})),
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_log_project" "log" {
  name        = var.name
  description = "created by terraform for managedkubernetes cluster"
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
`, name)
}

func resourceCSManagedKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_keys" "default" {
  status  = "Enabled"
  filters = "[{\"Key\":\"CreatorType\", \"Values\":[\"User\"]}]"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func resourceCSManagedKubernetesConfigDependence_essd(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
  kubernetes_node_role = "Worker"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_managed_kubernetes.default.id
  name                          = var.name
  vswitch_ids                   = [local.vswitch_id]
  instance_types                = [data.alicloud_instance_types.default.instance_types.0.id]
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  system_disk_performance_level = "PL0"
  desired_size                  = 1
}
`, name)
}

var csManagedKubernetesBasicMap = map[string]string{
	"new_nat_gateway":                    "true",
	"slb_internet_enabled":               "true",
	"name":                               CHECKSET,
	"security_group_id":                  CHECKSET,
	"version":                            CHECKSET,
	"certificate_authority.cluster_cert": CHECKSET,
	"certificate_authority.client_cert":  CHECKSET,
	"certificate_authority.client_key":   CHECKSET,
	"connections.api_server_internet":    CHECKSET,
	"connections.api_server_intranet":    CHECKSET,
	"connections.master_public_ip":       CHECKSET,
	"connections.service_domain":         CHECKSET,
	"worker_ram_role_name":               CHECKSET,
	"vpc_id":                             CHECKSET,
	"resource_group_id":                  CHECKSET,
	"slb_internet":                       CHECKSET,
	"slb_intranet":                       CHECKSET,
	"cluster_spec":                       CHECKSET,
	"slb_id":                             CHECKSET,
}
