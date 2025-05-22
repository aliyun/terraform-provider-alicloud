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
					"name":                   name,
					"worker_vswitch_ids":     []string{"${local.vswitch_id}"},
					"pod_cidr":               "10.93.0.0/16",
					"service_cidr":           "172.21.0.0/16",
					"slb_internet_enabled":   "true",
					"cluster_spec":           "ack.pro.small",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"security_group_id":      "${alicloud_security_group.default.0.id}",
					"deletion_protection":    "false",
					"enable_rrsa":            "false",
					"timezone":               "Asia/Shanghai",
					"proxy_mode":             "ipvs",
					"new_nat_gateway":        "true",
					"api_audiences":          []string{"https://kubernetes.default.svc"},
					"service_account_issuer": "https://kubernetes.default.svc",
					"cluster_domain":         "cluster.test",
					"custom_san":             "www.terraform.io",
					"maintenance_window": []map[string]string{
						{
							"enable":           "true",
							"maintenance_time": "2024-10-15T12:31:00.000+08:00",
							"duration":         "3h",
							"weekly_period":    "Thursday",
						},
					},
					"operation_policy": []map[string]interface{}{
						{
							"cluster_auto_upgrade": []map[string]interface{}{
								{
									"enabled": "true",
									"channel": "patch",
								},
							},
						},
					},
					"cluster_ca_cert":                clusterCaCertFile.Name(),
					"client_key":                     clientKeyFile.Name(),
					"client_cert":                    clientCertFile.Name(),
					"skip_set_certificate_authority": "false",
					"depends_on":                     []string{"alicloud_security_group.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
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
						"cluster_domain":                        "cluster.test",
						"custom_san":                            "www.terraform.io",
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "true",
						"maintenance_window.0.maintenance_time": "2024-10-15T12:31:00.000+08:00",
						"maintenance_window.0.duration":         "3h",
						"maintenance_window.0.weekly_period":    "Thursday",
						"operation_policy.#":                    "1",
						"operation_policy.0.cluster_auto_upgrade.#":         "1",
						"operation_policy.0.cluster_auto_upgrade.0.enabled": "true",
						"operation_policy.0.cluster_auto_upgrade.0.channel": "patch",
						"skip_set_certificate_authority":                    "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "user_ca", "name_prefix", "slb_internet_enabled", "api_audiences", "service_account_issuer", "load_balancer_spec", "encryption_provider_key", "cluster_ca_cert", "client_key", "client_cert", "worker_vswitch_ids"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name + "_update",
					"custom_san":          "www.terraform.io,terraform.test",
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name + "_update",
						"custom_san":          "www.terraform.io,terraform.test",
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
					"security_group_id": "${alicloud_security_group.default.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintenance_window": []map[string]string{{"enable": "false", "maintenance_time": "2024-10-15T11:31:00.000+08:00", "duration": "5h", "weekly_period": "Monday,Thursday"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "false",
						"maintenance_window.0.maintenance_time": "2024-10-15T11:31:00.000+08:00",
						"maintenance_window.0.duration":         "5h",
						"maintenance_window.0.weekly_period":    "Monday,Thursday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_policy": []map[string]interface{}{
						{
							"cluster_auto_upgrade": []map[string]interface{}{
								{
									"enabled": "false",
									"channel": "rapid",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_policy.#":                                "1",
						"operation_policy.0.cluster_auto_upgrade.#":         "1",
						"operation_policy.0.cluster_auto_upgrade.0.enabled": "false",
						"operation_policy.0.cluster_auto_upgrade.0.channel": "rapid",
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
					"version":             "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.2.version}",
					"pod_cidr":            "10.94.0.0/16",
					"service_cidr":        "172.22.0.0/16",
					"deletion_protection": "false",
					"cluster_spec":        "ack.standard",
					"new_nat_gateway":     "true",
					"proxy_mode":          "ipvs",
					"vswitch_ids":         []string{"${local.vswitch_id}"},
					"tags": map[string]string{
						"Platform": "TF",
					},
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// cluster args
						"name":                           name,
						"version":                        CHECKSET,
						"pod_cidr":                       "10.94.0.0/16",
						"service_cidr":                   "172.22.0.0/16",
						"deletion_protection":            "false",
						"cluster_spec":                   "ack.standard",
						"new_nat_gateway":                "true",
						"nat_gateway_id":                 CHECKSET,
						"proxy_mode":                     "ipvs",
						"vswitch_ids.#":                  "1",
						"tags.%":                         "1",
						"tags.Platform":                  "TF",
						"skip_set_certificate_authority": "false",
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
					"version": "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.1.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled"},
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
					"name_prefix":                    "tf-testaccmanagedkubernetes",
					"cluster_spec":                   "ack.pro.small",
					"is_enterprise_security_group":   "true",
					"deletion_protection":            "false",
					"new_nat_gateway":                "true",
					"node_cidr_mask":                 "26",
					"service_cidr":                   "172.23.0.0/16",
					"proxy_mode":                     "ipvs",
					"ip_stack":                       "ipv4",
					"vswitch_ids":                    []string{"${local.vswitch_id}", "${local.vswitch_id_1}"},
					"pod_vswitch_ids":                []string{"${local.vswitch_id}"},
					"control_plane_log_ttl":          "30",
					"control_plane_log_components":   []string{"apiserver", "kcm", "scheduler"},
					"control_plane_log_project":      "",
					"user_ca":                        tmpCAFile.Name(),
					"addons":                         []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"skip_set_certificate_authority": "false",
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
						"ip_stack":                       "ipv4",
						"vswitch_ids.#":                  "2",
						"control_plane_log_ttl":          "30",
						"control_plane_log_components.0": "apiserver",
						"control_plane_log_components.1": "kcm",
						"control_plane_log_components.2": "scheduler",
						"control_plane_log_project":      CHECKSET,
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{"${local.vswitch_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					})),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "user_ca", "name_prefix", "addons", "is_enterprise_security_group", "pod_vswitch_ids", "slb_internet_enabled", "load_balancer_spec"},
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

func TestAccAliCloudCSManagedKubernetesAuto(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfigDependenceAuto)

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
					"name":                           name,
					"zone_ids":                       []string{"${data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id}"},
					"cluster_spec":                   "ack.pro.small",
					"new_nat_gateway":                "false",
					"slb_internet_enabled":           "false",
					"deletion_protection":            "false",
					"node_cidr_mask":                 "26",
					"service_cidr":                   "172.23.0.0/16",
					"proxy_mode":                     "ipvs",
					"ip_stack":                       "ipv4",
					"timezone":                       "Asia/Shanghai",
					"addons":                         []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"vswitch_ids.#":                  "1",
						"cluster_spec":                   "ack.pro.small",
						"vpc_id":                         CHECKSET,
						"deletion_protection":            "false",
						"new_nat_gateway":                "false",
						"slb_internet_enabled":           "false",
						"ip_stack":                       "ipv4",
						"timezone":                       "Asia/Shanghai",
						"resource_group_id":              CHECKSET,
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "addons", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled", "zone_ids"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone": "Europe/London",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone": "Europe/London",
					})),
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
						{
							"delete_mode":   "delete",
							"resource_type": "PrivateZone",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_Data",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_ControlPlane",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{})),
			},
		},
	})
}

func TestAccAliCloudCSManagedKubernetesForProfile(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSManagedKubernetesConfigDependenceAuto)

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
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.1.version}",
					"zone_ids":                       []string{"${data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id}"},
					"profile":                        "Acs",
					"cluster_spec":                   "ack.pro.small",
					"new_nat_gateway":                "false",
					"slb_internet_enabled":           "false",
					"deletion_protection":            "false",
					"service_cidr":                   "172.23.0.0/16",
					"ip_stack":                       "ipv4",
					"timezone":                       "Asia/Shanghai",
					"skip_set_certificate_authority": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"version":                        CHECKSET,
						"vswitch_ids.#":                  "1",
						"cluster_spec":                   "ack.pro.small",
						"profile":                        "Acs",
						"vpc_id":                         CHECKSET,
						"deletion_protection":            "false",
						"new_nat_gateway":                "false",
						"slb_internet_enabled":           "false",
						"ip_stack":                       "ipv4",
						"timezone":                       "Asia/Shanghai",
						"resource_group_id":              CHECKSET,
						"skip_set_certificate_authority": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "certificate_authority", "addons", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled", "zone_ids"},
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
						{
							"delete_mode":   "delete",
							"resource_type": "PrivateZone",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_Data",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_ControlPlane",
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

resource "alicloud_vpc" "vpc" {
  count      = 1
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  count        = 2
  vpc_id       = alicloud_vpc.vpc.0.id
  cidr_block   = format("192.168.%%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.default.zones[count.index].id
  vswitch_name = var.name
}

resource "alicloud_log_project" "log" {
  name        = var.name
  description = "created by terraform for managedkubernetes cluster"
  lifecycle {
    ignore_changes = [
      policy
    ]
  }
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.0.id
  vswitch_id_1 = alicloud_vswitch.vswitches.1.id
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

data "alicloud_cs_kubernetes_version" "kubernetes_versions" {
  cluster_type = "ManagedKubernetes"
  profile      = "Default"
}
	
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  count      = 1
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  count        = 1
  vpc_id       = alicloud_vpc.vpc.0.id
  cidr_block   = format("192.168.%%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.default.zones[count.index].id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.0.id
}

resource "alicloud_security_group" "default" {
  count  = 2
  vpc_id = alicloud_vpc.vpc.0.id
}
`, name)
}

func resourceCSManagedKubernetesConfigDependence_essd(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "instance_type" {
  default = "ecs.c6.xlarge"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
  available_instance_type     = var.instance_type
}

data "alicloud_cs_kubernetes_version" "kubernetes_versions" {
  cluster_type = "ManagedKubernetes"
  profile      = "Default"
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
  zone_id      = data.alicloud_zones.default.zones[count.index].id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.0.id
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_managed_kubernetes.default.id
  node_pool_name                = var.name
  vswitch_ids                   = [local.vswitch_id]
  instance_types                = [var.instance_type]
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  instance_charge_type          = "PostPaid"
  system_disk_performance_level = "PL0"
  desired_size                  = 1
}
`, name)
}

func resourceCSManagedKubernetesConfigDependenceAuto(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_cs_kubernetes_version" "kubernetes_versions" {
  cluster_type = "ManagedKubernetes"
  profile      = "Acs"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
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
