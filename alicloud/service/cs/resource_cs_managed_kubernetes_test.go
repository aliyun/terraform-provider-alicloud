package cs_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCSManagedKubernetes_basic(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csManagedKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependence)

	clusterCaCertFile, clientCertFile, clientKeyFile, err := CreateTempFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clientCertFile.Name())
	defer os.Remove(clientKeyFile.Name())
	defer os.Remove(clusterCaCertFile.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":                     acctest.CHECKSET,
						"deletion_protection":                   "false",
						"enable_rrsa":                           "false",
						"timezone":                              "Asia/Shanghai",
						"proxy_mode":                            "ipvs",
						"new_nat_gateway":                       "true",
						"nat_gateway_id":                        acctest.CHECKSET,
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
				// Addon config/version upgrades and the delete-time retain list
				// cannot be exercised against a live cluster in this suite;
				// removing the only allowed vswitch selector stops the plan
				// before the placeholders reach the API.
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]interface{}{
						{"name": "terway-eniip", "config": "securityGroups", "version": "1.3.0", "disabled": "false"},
					},
					"retain_resources":   []string{"ecs"},
					"worker_vswitch_ids": acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("one of `vswitch_ids,worker_vswitch_ids,zone_ids` must be\\s+specified"),
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
						"resource_group_id": acctest.CHECKSET,
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
						"security_group_id": acctest.CHECKSET,
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

func TestAccAliCloudCSManagedKubernetes_encryption(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csManagedKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-encryption-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"worker_vswitch_ids":   []string{"${local.vswitch_id}"},
					"pod_cidr":             "10.93.0.0/16",
					"service_cidr":         "172.21.0.0/16",
					"slb_internet_enabled": "true",
					"cluster_spec":         "ack.pro.small",
					"deletion_protection":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"pod_cidr":             "10.93.0.0/16",
						"service_cidr":         "172.21.0.0/16",
						"slb_internet_enabled": "true",
						"cluster_spec":         "ack.pro.small",
						"deletion_protection":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encryption_provider_key": "${data.alicloud_kms_keys.default.keys[0].key_id}",
					"disable_encryption":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encryption_provider_key": acctest.CHECKSET,
						"disable_encryption":      "false",
					}),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Minute) },
				Config: testAccConfig(map[string]interface{}{
					"encryption_provider_key": "",
					"disable_encryption":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disable_encryption": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "user_ca", "name_prefix", "slb_internet_enabled", "api_audiences", "service_account_issuer", "load_balancer_spec", "cluster_ca_cert", "client_key", "client_cert", "worker_vswitch_ids"},
			},
		},
	})
}

func TestAccAliCloudCSManagedKubernetes_essd_migrate_upgrade(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csManagedKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependence_essd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"version":                        acctest.CHECKSET,
						"pod_cidr":                       "10.94.0.0/16",
						"service_cidr":                   "172.22.0.0/16",
						"deletion_protection":            "false",
						"cluster_spec":                   "ack.standard",
						"new_nat_gateway":                "true",
						"nat_gateway_id":                 acctest.CHECKSET,
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
						"version": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled"},
			},
			{
				// no-op apply that only sleeps before the framework destroys the resources,
				// giving terway pod/node ENIs time to be released so the dependent vSwitch
				// can be deleted without hitting DependencyViolation.NetworkInterface.
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_sleep": "100",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSManagedKubernetes_controlPlanLog(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfig)

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
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"ip_stack":                     "ipv4",
					"vswitch_ids":                  []string{"${local.vswitch_id}", "${local.vswitch_id_1}"},
					"pod_vswitch_ids":              []string{"${local.vswitch_id}"},
					"control_plane_log_ttl":        "30",
					"control_plane_log_components": []string{"apiserver", "kcm", "scheduler"},
					"control_plane_log_project":    "",
					"audit_log_config": []map[string]string{{
						"enabled":          "true",
						"sls_project_name": "${alicloud_log_project.log.0.name}",
					}},
					"user_ca":                        tmpCAFile.Name(),
					"addons":                         []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                acctest.CHECKSET,
						"cluster_spec":                        "ack.pro.small",
						"deletion_protection":                 "false",
						"new_nat_gateway":                     "true",
						"nat_gateway_id":                      acctest.CHECKSET,
						"service_cidr":                        "172.23.0.0/16",
						"proxy_mode":                          "ipvs",
						"ip_stack":                            "ipv4",
						"vswitch_ids.#":                       "2",
						"control_plane_log_ttl":               "30",
						"control_plane_log_components.0":      "apiserver",
						"control_plane_log_components.1":      "kcm",
						"control_plane_log_components.2":      "scheduler",
						"control_plane_log_project":           acctest.CHECKSET,
						"audit_log_config.0.enabled":          "true",
						"audit_log_config.0.sls_project_name": acctest.CHECKSET,
						"skip_set_certificate_authority":      "false",
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
					"control_plane_log_project": "${alicloud_log_project.log.0.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_log_project": fmt.Sprintf("%s-0", name),
					})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_log_config": []map[string]string{{
						"enabled":          "true",
						"sls_project_name": "${alicloud_log_project.log.1.name}",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_log_config.0.sls_project_name": fmt.Sprintf("%s-1", name),
					})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_log_config": []map[string]string{{
						"enabled": "false",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_log_config.0.enabled": "false",
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

func TestAccAliCloudCSManagedKubernetes_auto(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependenceAuto)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"zone_ids":             []string{"${data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id}"},
					"cluster_spec":         "ack.pro.small",
					"new_nat_gateway":      "false",
					"slb_internet_enabled": "false",
					"deletion_protection":  "false",
					"node_cidr_mask":       "26",
					"service_cidr":         "172.23.0.0/16",
					"proxy_mode":           "ipvs",
					"ip_stack":             "ipv4",
					"timezone":             "Asia/Shanghai",
					"addons":               []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"auto_mode": []map[string]string{{
						"enabled": "true",
					}},
					"maintenance_window": []map[string]string{
						{
							"enable":           "true",
							"maintenance_time": "2024-10-15T12:31:00.000+08:00",
							"duration":         "3h",
							"weekly_period":    "Thursday",
						},
					},
					"skip_set_certificate_authority": "false",
					"control_plane_log_ttl":          "30",
					"control_plane_log_components":   []string{"apiserver", "kcm", "scheduler"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                  name,
						"vswitch_ids.#":                         "1",
						"cluster_spec":                          "ack.pro.small",
						"vpc_id":                                acctest.CHECKSET,
						"deletion_protection":                   "false",
						"new_nat_gateway":                       "false",
						"slb_internet_enabled":                  "false",
						"ip_stack":                              "ipv4",
						"timezone":                              "Asia/Shanghai",
						"resource_group_id":                     acctest.CHECKSET,
						"skip_set_certificate_authority":        "false",
						"auto_mode.0.enabled":                   "true",
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "true",
						"maintenance_window.0.maintenance_time": "2024-10-15T12:31:00.000+08:00",
						"maintenance_window.0.duration":         "3h",
						"maintenance_window.0.weekly_period":    "Thursday",
						"control_plane_log_ttl":                 "30",
						"control_plane_log_components.0":        "apiserver",
						"control_plane_log_components.1":        "kcm",
						"control_plane_log_components.2":        "scheduler",
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

func TestAccAliCloudCSManagedKubernetes_controlPlaneEndpoints(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependenceNlb)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// create with a custom NLB endpoint (load_balancers_config)
				Config: testAccConfig(map[string]interface{}{
					"name":                         name,
					"cluster_spec":                 "ack.pro.small",
					"vswitch_ids":                  []string{"${alicloud_vswitch.default_1.id}"},
					"new_nat_gateway":              "false",
					"deletion_protection":          "false",
					"service_cidr":                 "172.23.0.0/16",
					"pod_cidr":                     "172.20.0.0/16",
					"is_enterprise_security_group": "true",
					"addons":                       []map[string]string{{"name": "flannel", "config": "", "version": "", "disabled": "false"}},
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.default.id}",
									"endpoint_type":    "private",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                             name,
						"cluster_spec":                     "ack.pro.small",
						"vpc_id":                           acctest.CHECKSET,
						"control_plane_endpoints_config.#": "1",
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
					}),
				),
			},
			{
				// change load_balancers_config: switch to another NLB instance
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]interface{}{
						{
							"disabled": "true",
							"name":     "terway-eniip",
						},
					},
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
					}),
				),
			},
			{
				// add a second endpoint (A -> AB), both private
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.default.id}",
									"endpoint_type":    "private",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.load_balancers_config.#":                  "2",
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.1.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint":         acctest.CHECKSET,
					}),
				),
			},
			{
				// remove the second endpoint (AB -> A)
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.load_balancers_config.#":                  "1",
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
						// the second endpoint is removed; clear its accumulated expectations
						"control_plane_endpoints_config.0.load_balancers_config.1.load_balancer_id": acctest.REMOVEKEY,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint_type":    acctest.REMOVEKEY,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint":         acctest.REMOVEKEY,
					}),
				),
			},
			{
				// add a public endpoint (Internet NLB) alongside the private one; backend requires at least one private NLB retained
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.public.id}",
									"endpoint_type":    "public",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.load_balancers_config.#":                  "2",
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.1.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint_type":    "public",
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint":         acctest.CHECKSET,
					}),
				),
			},
			{
				// remove the public endpoint, keep only the private one
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.load_balancers_config.#":                  "1",
						"control_plane_endpoints_config.0.load_balancers_config.0.load_balancer_id": acctest.CHECKSET,
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint_type":    "private",
						"control_plane_endpoints_config.0.load_balancers_config.0.endpoint":         acctest.CHECKSET,
						// the public endpoint is removed; clear its accumulated expectations
						"control_plane_endpoints_config.0.load_balancers_config.1.load_balancer_id": acctest.REMOVEKEY,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint_type":    acctest.REMOVEKEY,
						"control_plane_endpoints_config.0.load_balancers_config.1.endpoint":         acctest.REMOVEKEY,
					}),
				),
			},
			{
				// enable cluster internal domain name access (internal_dns_config) with bind_vpcs
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
							},
							"internal_dns_config": []map[string]interface{}{
								{
									"enabled":   "true",
									"bind_vpcs": []string{"${alicloud_vpc.default.id}"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.internal_dns_config.0.enabled":     "true",
						"control_plane_endpoints_config.0.internal_dns_config.0.bind_vpcs.#": "1",
					}),
				),
			},
			{
				// change internal_dns_config.bind_vpcs: bind an additional VPC
				Config: testAccConfig(map[string]interface{}{
					"control_plane_endpoints_config": []map[string]interface{}{
						{
							"load_balancers_config": []map[string]interface{}{
								{
									"load_balancer_id": "${alicloud_nlb_load_balancer.update.id}",
									"endpoint_type":    "private",
								},
							},
							"internal_dns_config": []map[string]interface{}{
								{
									"enabled":   "true",
									"bind_vpcs": []string{"${alicloud_vpc.default.id}", "${alicloud_vpc.bind.id}"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"control_plane_endpoints_config.0.internal_dns_config.0.enabled":     "true",
						"control_plane_endpoints_config.0.internal_dns_config.0.bind_vpcs.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "addons", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled", "zone_ids", "is_enterprise_security_group"},
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
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfigDependenceAuto)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"version":                        acctest.CHECKSET,
						"vswitch_ids.#":                  "1",
						"cluster_spec":                   "ack.pro.small",
						"profile":                        "Acs",
						"vpc_id":                         acctest.CHECKSET,
						"deletion_protection":            "false",
						"new_nat_gateway":                "false",
						"slb_internet_enabled":           "false",
						"ip_stack":                       "ipv4",
						"timezone":                       "Asia/Shanghai",
						"resource_group_id":              acctest.CHECKSET,
						"skip_set_certificate_authority": "false",
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
  count       = 2
  name        = format("%%s-%%d", var.name, count.index)
  description = "created by terraform for managedkubernetes cluster"
  lifecycle {
    ignore_changes = [
      policy
    ]
  }
}

locals {
  vswitch_id   = alicloud_vswitch.vswitches.0.id
  vswitch_id_1 = alicloud_vswitch.vswitches.1.id
}
`, name)
}

func resourceCSManagedKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

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
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitches.0.id
}

resource "alicloud_security_group" "default" {
  count  = 2
  vpc_id = alicloud_vpc.vpc.0.id
}

data "alicloud_kms_keys" "default" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]},{\"Key\":\"KeySpec\",\"Values\":[\"Aliyun_AES_256\"]},{\"Key\":\"KeyUsage\",\"Values\":[\"ENCRYPT/DECRYPT\"]},{\"Key\":\"CreatorType\",\"Values\":[\"User\"]}]"
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
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  instance_charge_type          = "PostPaid"
  system_disk_performance_level = "PL0"
  desired_size                  = 1
  instance_patterns {
    instance_family_level = "EnterpriseLevel"
    cpu_architectures     = ["X86"]
    instance_categories   = ["General-purpose"]
    min_cpu_cores         = 4
    max_cpu_cores         = 8
    min_memory_size       = 8
    max_memory_size       = 16
  }
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

func resourceCSManagedKubernetesConfigDependenceNlb(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vpc" "bind" {
  vpc_name   = "${var.name}-bind"
  cidr_block = "10.2.0.0/21"
}

resource "alicloud_vswitch" "default_1" {
  vswitch_name = "${var.name}-1"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.1.1.0/24"
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default_2" {
  vswitch_name = "${var.name}-2"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.1.2.0/24"
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
  security_group_type = "enterprise"
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = "${var.name}-1"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_type = "Network"
  address_type       = "Intranet"
  vpc_id             = alicloud_vpc.default.id
  security_group_ids = [alicloud_security_group.default.id]
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_2.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_nlb_load_balancer" "update" {
  load_balancer_name = "${var.name}-2"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_type = "Network"
  address_type       = "Intranet"
  vpc_id             = alicloud_vpc.default.id
  security_group_ids = [alicloud_security_group.default.id]
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_2.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_nlb_load_balancer" "public" {
  load_balancer_name = "${var.name}-3"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_type = "Network"
  address_type       = "Internet"
  vpc_id             = alicloud_vpc.default.id
  security_group_ids = [alicloud_security_group.default.id]
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_2.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}
`, name)
}

var csManagedKubernetesBasicMap = map[string]string{
	"new_nat_gateway":                      "true",
	"slb_internet_enabled":                 "true",
	"name":                                 acctest.CHECKSET,
	"security_group_id":                    acctest.CHECKSET,
	"version":                              acctest.CHECKSET,
	"certificate_authority.0.cluster_cert": acctest.CHECKSET,
	"certificate_authority.0.client_cert":  acctest.CHECKSET,
	"certificate_authority.0.client_key":   acctest.CHECKSET,
	"connections.0.api_server_internet":    acctest.CHECKSET,
	"connections.0.api_server_intranet":    acctest.CHECKSET,
	"connections.0.master_public_ip":       acctest.CHECKSET,
	"connections.0.service_domain":         acctest.CHECKSET,
	"worker_ram_role_name":                 acctest.CHECKSET,
	"vpc_id":                               acctest.CHECKSET,
	"resource_group_id":                    acctest.CHECKSET,
	"slb_internet":                         acctest.CHECKSET,
	"slb_intranet":                         acctest.CHECKSET,
	"cluster_spec":                         acctest.CHECKSET,
	"slb_id":                               acctest.CHECKSET,
}

// Test Ack Cluster. >>> Resource test cases, automatically generated.
// Case 集群测试-升级 12341
func TestAccAliCloudCSManagedKubernetes_basic12341(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_managed_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAckClusterBasicDependence12341)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name,
					"service_cidr": "${var.service_cidr}",
					"addons": []map[string]interface{}{
						{
							"name": "terway-eniip",
						},
					},
					"is_enterprise_security_group": "true",
					"zone_ids": []string{
						"${data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id}"},
					"cluster_spec":        "ack.pro.small",
					"deletion_protection": "false",
					"version":             "1.32.7-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"service_cidr":                 acctest.CHECKSET,
						"addons.#":                     "1",
						"is_enterprise_security_group": "true",
						"zone_ids.#":                   "1",
						"cluster_spec":                 "ack.pro.small",
						"deletion_protection":          "false",
						"version":                      "1.32.7-aliyun.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone": "Asia/Singapore",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone": "Asia/Singapore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "1.33.3-aliyun.1",
					"upgrade_policy": []map[string]interface{}{
						{
							"control_plane_only": true,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.33.3-aliyun.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "addons", "new_nat_gateway", "user_ca", "name_prefix", "load_balancer_spec", "slb_internet_enabled", "zone_ids", "is_enterprise_security_group", "upgrade_policy"},
			},
		},
	})
}

func AlicloudAckClusterBasicDependence12341(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

variable "cluster_name" {
  default = "test-create-cluster"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

variable "cluster_name_update" {
  default = "test-create-cluster-update"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "container_cidr" {
  default = "10.53.0.0/16"
}

variable "cluster_type" {
  default = "ManagedKubernetes"
}`, name)
}

func TestAccAliCloudCSManagedKubernetes_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_cs_managed_kubernetes.default"
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
				Config: testAccCSManagedKubernetesStateMigrationConfigV0(name),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
				Config:                   testAccCSManagedKubernetesStateMigrationConfigV1(name),
			},
		},
	})
}

func testAccCSManagedKubernetesStateMigrationDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
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
}`, name)
}

func testAccCSManagedKubernetesStateMigrationConfigV0(name string) string {
	return testAccCSManagedKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_managed_kubernetes" "default" {
  name                           = var.name
  cluster_spec                   = "ack.pro.small"
  vswitch_ids                    = [alicloud_vswitch.default.id]
  new_nat_gateway                = true
  pod_cidr                       = "10.77.0.0/16"
  service_cidr                   = "172.27.0.0/16"
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
}

output "alicloud_cs_managed_kubernetes" {
  value = alicloud_cs_managed_kubernetes.default.certificate_authority["cluster_cert"]
}

output "connections" {
  value = alicloud_cs_managed_kubernetes.default.connections["api_server_internet"]
}
`
}

func testAccCSManagedKubernetesStateMigrationConfigV1(name string) string {
	return testAccCSManagedKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_managed_kubernetes" "default" {
  name                           = var.name
  cluster_spec                   = "ack.pro.small"
  vswitch_ids                    = [alicloud_vswitch.default.id]
  new_nat_gateway                = true
  pod_cidr                       = "10.77.0.0/16"
  service_cidr                   = "172.27.0.0/16"
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
}

output "alicloud_cs_managed_kubernetes" {
  value = alicloud_cs_managed_kubernetes.default.certificate_authority.0.cluster_cert
}

output "connections" {
  value = alicloud_cs_managed_kubernetes.default.connections.0.api_server_internet
}
`
}

// Test Ack Cluster. <<< Resource test cases, automatically generated.

func TestAccAliCloudCSManagedKubernetes_collectionOrder(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_managed_kubernetes.collectionorder"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedkubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfig)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                         name,
					"cluster_spec":                 "ack.standard",
					"vswitch_ids":                  []string{"${local.vswitch_id}", "${local.vswitch_id_1}"},
					"pod_vswitch_ids":              []string{"${local.vswitch_id}", "${local.vswitch_id_1}"},
					"new_nat_gateway":              "true",
					"pod_cidr":                     "10.77.0.0/16",
					"service_cidr":                 "172.27.0.0/16",
					"is_enterprise_security_group": "true",
					"addons": []map[string]string{
						{"name": "metrics-server", "disabled": "false"},
						{"name": "arms-prometheus", "disabled": "false"},
					},
					"api_audiences":                []string{"https://kubernetes.default.svc", "https://kubernetes.default.svc.cluster.local"},
					"control_plane_log_project":    "${alicloud_log_project.log.0.name}",
					"control_plane_log_components": []string{"kubernetes-audit", "kubernetes-controller-manager"},
					"retain_resources":             []string{"${local.vswitch_id}", "${local.vswitch_id_1}"},
					"delete_options": []map[string]interface{}{
						{"delete_mode": "delete", "resource_type": "SLB"},
						{"delete_mode": "delete", "resource_type": "ALB"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// addons: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{"name": "arms-prometheus", "disabled": "false"},
						{"name": "metrics-server", "disabled": "false"},
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
			// api_audiences: verify reorder convergence (ForceNew)
			{
				Config: testAccConfig(map[string]interface{}{
					"api_audiences": []string{"https://kubernetes.default.svc.cluster.local", "https://kubernetes.default.svc"},
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
			// control_plane_log_components: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"control_plane_log_components": []string{"kubernetes-controller-manager", "kubernetes-audit"},
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
					"retain_resources": []string{"${local.vswitch_id_1}", "${local.vswitch_id}"},
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
			// delete_options: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]interface{}{
						{"delete_mode": "delete", "resource_type": "ALB"},
						{"delete_mode": "delete", "resource_type": "SLB"},
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
			// pod_vswitch_ids: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"pod_vswitch_ids": []string{"${local.vswitch_id_1}", "${local.vswitch_id}"},
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
			// vswitch_ids: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{"${local.vswitch_id_1}", "${local.vswitch_id}"},
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
		},
	})
}

func TestAccAliCloudCSManagedKubernetes_collectionOrder_zoneIds(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail
	resourceId := "alicloud_cs_managed_kubernetes.collectionorderzone"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccmanagedk8szone-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSManagedKubernetesConfig)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                         name,
					"cluster_spec":                 "ack.standard",
					"zone_ids":                     []string{"${data.alicloud_zones.default.zones.0.id}", "${data.alicloud_zones.default.zones.1.id}"},
					"new_nat_gateway":              "true",
					"pod_cidr":                     "10.77.0.0/16",
					"service_cidr":                 "172.27.0.0/16",
					"is_enterprise_security_group": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// zone_ids: verify reorder produces a non-empty plan
				Config: testAccConfig(map[string]interface{}{
					"zone_ids": []string{"${data.alicloud_zones.default.zones.1.id}", "${data.alicloud_zones.default.zones.0.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Apply the reordered config and verify convergence (empty plan)
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}
