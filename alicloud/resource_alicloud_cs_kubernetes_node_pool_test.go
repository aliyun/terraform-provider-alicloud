package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCSKubernetesNodePool_basic(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Rds)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"cluster_id":  "${local.cluster_id}",
					"vswitch_ids": []string{"${local.vswitch_id}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.default.instance_types.0.id}",
						"${data.alicloud_instance_types.default.instance_types.1.id}",
					},
					"desired_size":          "1",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks": []map[string]string{
						{
							"size":     "100",
							"category": "cloud_ssd",
							"name":     name,
						},
					},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "surge_percentage": "10", "max_unavailable": "0"}},
					"image_type":            "AliyunLinux3",
					"security_hardening_os": "true",
					"cpu_policy":            "none",
					"spot_strategy":         "NoSpot",
					"rds_instances":         []string{"${alicloud_db_instance.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"cluster_id":                    CHECKSET,
						"vswitch_ids.#":                 "1",
						"instance_types.#":              "2",
						"desired_size":                  "1",
						"key_name":                      CHECKSET,
						"system_disk_category":          "cloud_efficiency",
						"system_disk_size":              "40",
						"install_cloud_monitor":         "false",
						"data_disks.#":                  "1",
						"data_disks.0.size":             "100",
						"data_disks.0.category":         "cloud_ssd",
						"data_disks.0.name":             name,
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.Foo":                      "Bar",
						"management.#":                  "1",
						"management.0.auto_repair":      "true",
						"management.0.auto_upgrade":     "true",
						"management.0.surge":            "0",
						"management.0.surge_percentage": "10",
						"management.0.max_unavailable":  "0",
						"security_group_ids.#":          "1",
						"image_type":                    "AliyunLinux3",
						"security_hardening_os":         "true",
						"cpu_policy":                    "none",
						"spot_strategy":                 "NoSpot",
						"rds_instances.#":               "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name": "${alicloud_key_pair.default2.key_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rds_instances": []string{"${alicloud_db_instance.default.0.id}", "${alicloud_db_instance.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rds_instances.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "2",
					"management":   []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "surge_percentage": "20", "max_unavailable": "1"}},
					"data_disks": []map[string]string{
						{
							"size":     "40",
							"category": "cloud_ssd",
							"name":     name + "_update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size":                  "2",
						"data_disks.#":                  "1",
						"data_disks.0.size":             "40",
						"data_disks.0.category":         "cloud_ssd",
						"data_disks.0.name":             name + "_update",
						"management.#":                  "1",
						"management.0.auto_repair":      "true",
						"management.0.auto_upgrade":     "true",
						"management.0.surge":            "1",
						"management.0.surge_percentage": "20",
						"management.0.max_unavailable":  "1",
					}),
				),
			},
			// check: remove nodes
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "1",
					}),
				),
			},
			// check: kubelet config
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{{
						"registry_pull_qps":     "0",
						"registry_burst":        "0",
						"event_record_qps":      "0",
						"event_burst":           "0",
						"serialize_image_pulls": "false",
						"cpu_manager_policy":    "none",
						"feature_gates": map[string]interface{}{
							"\"APIListChunking\"":                "true",
							"\"CustomCPUCFSQuotaPeriod\"":        "true",
							"\"MemoryManager\"":                  "true",
							"\"RotateKubeletServerCertificate\"": "true",
						},
					}},
					"rolling_policy": []map[string]interface{}{{
						"max_parallelism": "1",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kubelet_configuration.#":                       "1",
						"kubelet_configuration.0.registry_pull_qps":     "0",
						"kubelet_configuration.0.registry_burst":        "0",
						"kubelet_configuration.0.event_record_qps":      "0",
						"kubelet_configuration.0.event_burst":           "0",
						"kubelet_configuration.0.serialize_image_pulls": "false",
						"kubelet_configuration.0.cpu_manager_policy":    "none",
						"rolling_policy.#":                              "1",
						"rolling_policy.0.max_parallelism":              "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePoolWithNodeCount_basic(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.with_node_count"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "${local.cluster_id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"node_count":            "2",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks":            []map[string]string{{"size": "100", "category": "cloud_ssd"}},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
					"image_type":            "AliyunLinux3",
					"cpu_policy":            "none",
					"spot_strategy":         "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"node_count":                   "2",
						"key_name":                     CHECKSET,
						"system_disk_category":         "cloud_efficiency",
						"system_disk_size":             "40",
						"install_cloud_monitor":        "false",
						"data_disks.#":                 "1",
						"data_disks.0.size":            "100",
						"data_disks.0.category":        "cloud_ssd",
						"tags.%":                       "2",
						"tags.Created":                 "TF",
						"tags.Foo":                     "Bar",
						"management.#":                 "1",
						"management.0.auto_repair":     "true",
						"management.0.auto_upgrade":    "true",
						"management.0.surge":           "0",
						"management.0.max_unavailable": "0",
						"security_group_ids.#":         "1",
						"image_type":                   "AliyunLinux3",
						"cpu_policy":                   "none",
						"spot_strategy":                "NoSpot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "cpu_policy"},
			},
			// check: scale out
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count":       "2",
					"system_disk_size": "80",
					"data_disks":       []map[string]string{{"size": "40", "category": "cloud"}},
					"management":       []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "max_unavailable": "0"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count":                   "2",
						"system_disk_size":             "80",
						"data_disks.#":                 "1",
						"data_disks.0.size":            "40",
						"data_disks.0.category":        "cloud",
						"management.#":                 "1",
						"management.0.auto_repair":     "true",
						"management.0.auto_upgrade":    "true",
						"management.0.surge":           "1",
						"management.0.max_unavailable": "0",
					}),
				),
			},
			// check: remove nodes
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count": "1",
					}),
				),
			},
			// check: change node_count to desire_size
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count":   "#REMOVEKEY",
					"desired_size": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_autoScaling(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.autocaling"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                   name,
					"cluster_id":             "${local.cluster_id}",
					"vswitch_ids":            []string{"${local.vswitch_id}"},
					"instance_types":         []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"key_name":               "${alicloud_key_pair.default.key_name}",
					"system_disk_categories": []string{"cloud_efficiency", "cloud_essd"},
					"system_disk_size":       "40",
					"install_cloud_monitor":  "false",
					"platform":               "AliyunLinux",
					"scaling_policy":         "release",
					"scaling_config":         []map[string]string{{"enable": "true", "min_size": "1", "max_size": "10", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
					"cpu_policy":             "none",
					"spot_strategy":          "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                      name,
						"cluster_id":                                CHECKSET,
						"vswitch_ids.#":                             "1",
						"instance_types.#":                          "1",
						"key_name":                                  CHECKSET,
						"system_disk_categories.#":                  "2",
						"system_disk_size":                          "40",
						"install_cloud_monitor":                     "false",
						"platform":                                  "AliyunLinux",
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "true",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "10",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
						"cpu_policy":                                "none",
						"spot_strategy":                             "NoSpot",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                   name + "_update",
					"system_disk_categories": []string{"cloud_efficiency"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                     name + "_update",
						"system_disk_categories.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_type": "AliyunLinux3",
					"platform":   "AliyunLinux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_type": "AliyunLinux3",
						"platform":   "AliyunLinux",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// check: update config
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_policy": "release",
					"scaling_config": []map[string]string{{"enable": "true", "min_size": "2", "max_size": "20", "type": "spot", "is_bond_eip": "true", "eip_internet_charge_type": "PayByTraffic", "eip_bandwidth": "100"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "true",
						"scaling_config.0.min_size":                 "2",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "spot",
						"scaling_config.0.is_bond_eip":              "true",
						"scaling_config.0.eip_internet_charge_type": "PayByTraffic",
						"scaling_config.0.eip_bandwidth":            "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]string{{"enable": "true", "min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "false", "eip_internet_charge_type": "PayByTraffic", "eip_bandwidth": "100"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "true",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "false",
						"scaling_config.0.eip_internet_charge_type": "PayByTraffic",
						"scaling_config.0.eip_bandwidth":            "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]string{{"enable": "false", "min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "false", "eip_internet_charge_type": "PayByTraffic", "eip_bandwidth": "100"}},
					"desired_size":   "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "false",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "false",
						"scaling_config.0.eip_internet_charge_type": "PayByTraffic",
						"scaling_config.0.eip_bandwidth":            "100",
						"desired_size":                              "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_PrePaid(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.pre_paid_nodepool"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "${local.cluster_id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"password":              "Terraform1234",
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "120",
					"install_cloud_monitor": "false",
					"instance_charge_type":  "PrePaid",
					"period":                "1",
					"period_unit":           "Month",
					"auto_renew_period":     "1",
					"cpu_policy":            "none",
					"spot_strategy":         "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"cluster_id":            CHECKSET,
						"password":              CHECKSET,
						"vswitch_ids.#":         "1",
						"instance_types.#":      "1",
						"security_group_id":     CHECKSET,
						"system_disk_category":  "cloud_efficiency",
						"system_disk_size":      "120",
						"instance_charge_type":  "PrePaid",
						"install_cloud_monitor": "false",
						"period":                "1",
						"period_unit":           "Month",
						"auto_renew_period":     "1",
						"cpu_policy":            "none",
						"spot_strategy":         "NoSpot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "cpu_policy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":  "PrePaid",
					"auto_renew_period":     "2",
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":  "PrePaid",
						"auto_renew_period":     "2",
						"install_cloud_monitor": "true",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_Spot(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.spot_nodepool"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                       name,
					"cluster_id":                 "${local.cluster_id}",
					"vswitch_ids":                []string{"${local.vswitch_id}"},
					"instance_types":             []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"image_id":                   "aliyun_2_1903_x64_20G_alibase_20231008.vhd",
					"system_disk_category":       "cloud_efficiency",
					"system_disk_size":           "120",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"password":                   "Terraform1234",
					"desired_size":               "1",
					"install_cloud_monitor":      "false",
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "5",
					"spot_strategy":              "SpotWithPriceLimit",
					"spot_price_limit": []map[string]string{
						{
							"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
							"price_limit":   "0.57",
						},
					},
					"cpu_policy": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                             name,
						"cluster_id":                       CHECKSET,
						"vswitch_ids.#":                    "1",
						"instance_types.#":                 "1",
						"image_id":                         "aliyun_2_1903_x64_20G_alibase_20231008.vhd",
						"system_disk_category":             "cloud_efficiency",
						"system_disk_size":                 "120",
						"resource_group_id":                CHECKSET,
						"password":                         CHECKSET,
						"desired_size":                     "1",
						"install_cloud_monitor":            "false",
						"internet_charge_type":             "PayByTraffic",
						"internet_max_bandwidth_out":       "5",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit.#":               "1",
						"spot_price_limit.0.instance_type": CHECKSET,
						"spot_price_limit.0.price_limit":   "0.57",
						"cpu_policy":                       "none",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "cpu_policy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "aliyun_3_x64_20G_alibase_20230110.vhd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": "aliyun_3_x64_20G_alibase_20230110.vhd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "10",
					"spot_price_limit": []map[string]string{
						{
							"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
							"price_limit":   "0.60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type":             "PayByTraffic",
						"internet_max_bandwidth_out":       "10",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit.#":               "1",
						"spot_price_limit.0.instance_type": CHECKSET,
						"spot_price_limit.0.price_limit":   "0.60",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_KMS(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_BYOK)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"cluster_id":           "${local.cluster_id}",
					"vswitch_ids":          []string{"${local.vswitch_id}"},
					"instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"desired_size":         "1",
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "40",
					"data_disks": []map[string]string{
						{
							"kms_key_id": "${data.alicloud_kms_keys.default.ids.0}",
							"encrypted":  "true",
							"size":       "100",
							"category":   "cloud_essd",
						},
					},
					"tags":                          map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"image_type":                    "AliyunLinux3",
					"system_disk_encrypted":         "true",
					"system_disk_kms_key":           "${data.alicloud_kms_keys.default.ids.0}",
					"system_disk_encrypt_algorithm": "aes-256",
					"cpu_policy":                    "none",
					"spot_strategy":                 "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"cluster_id":                    CHECKSET,
						"vswitch_ids.#":                 "1",
						"instance_types.#":              "1",
						"desired_size":                  "1",
						"key_name":                      CHECKSET,
						"system_disk_category":          "cloud_essd",
						"system_disk_size":              "40",
						"data_disks.#":                  "1",
						"data_disks.0.size":             "100",
						"data_disks.0.category":         "cloud_essd",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.Foo":                      "Bar",
						"security_group_ids.#":          "1",
						"image_type":                    "AliyunLinux3",
						"system_disk_encrypted":         "true",
						"system_disk_kms_key":           CHECKSET,
						"system_disk_encrypt_algorithm": "aes-256",
						"cpu_policy":                    "none",
						"spot_strategy":                 "NoSpot",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_kms_key": "${data.alicloud_kms_keys.default.ids.1}",
					"data_disks": []map[string]string{
						{
							"size":       "100",
							"category":   "cloud_essd",
							"device":     "/dev/xvdc",
							"kms_key_id": "${data.alicloud_kms_keys.default.ids.1}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_kms_key": CHECKSET,
						"data_disks.#":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_encrypted":         "false",
					"system_disk_encrypt_algorithm": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_encrypted":         "false",
						"system_disk_encrypt_algorithm": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_BYOK(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_BYOK)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"cluster_id":           "${local.cluster_id}",
					"vswitch_ids":          []string{"${local.vswitch_id}"},
					"instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"desired_size":         "1",
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "40",
					"data_disks": []map[string]string{
						{
							"kms_key_id": "${data.alicloud_kms_keys.default.ids.0}",
							"encrypted":  "true",
							"size":       "100",
							"category":   "cloud_essd",
						},
						{
							"size":        "100",
							"category":    "cloud_essd",
							"device":      "/dev/xvdb",
							"snapshot_id": "${alicloud_ecs_snapshot.default.0.id}",
						},
					},
					"tags":                          map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"image_type":                    "AliyunLinux3",
					"system_disk_encrypted":         "true",
					"system_disk_kms_key":           "${data.alicloud_kms_keys.default.ids.0}",
					"system_disk_encrypt_algorithm": "aes-256",
					"cpu_policy":                    "none",
					"spot_strategy":                 "NoSpot",
					"ram_role_name":                 "${alicloud_ram_role.role.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"cluster_id":                    CHECKSET,
						"vswitch_ids.#":                 "1",
						"instance_types.#":              "1",
						"desired_size":                  "1",
						"key_name":                      CHECKSET,
						"system_disk_category":          "cloud_essd",
						"system_disk_size":              "40",
						"data_disks.#":                  "2",
						"data_disks.0.size":             "100",
						"data_disks.0.category":         "cloud_essd",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.Foo":                      "Bar",
						"security_group_ids.#":          "1",
						"image_type":                    "AliyunLinux3",
						"system_disk_encrypted":         "true",
						"system_disk_kms_key":           CHECKSET,
						"system_disk_encrypt_algorithm": "aes-256",
						"cpu_policy":                    "none",
						"spot_strategy":                 "NoSpot",
						"ram_role_name":                 CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_kms_key": "${data.alicloud_kms_keys.default.ids.1}",
					"data_disks": []map[string]string{
						{
							"size":        "100",
							"category":    "cloud_essd",
							"device":      "/dev/xvdc",
							"snapshot_id": "${alicloud_ecs_snapshot.default.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"system_disk_kms_key": "${data.alicloud_kms_keys.default.ids.1}",
						"data_disks.#": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_encrypted":         "false",
					"system_disk_encrypt_algorithm": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_encrypted":         "false",
						"system_disk_encrypt_algorithm": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_DeploymentSet(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_DeploymentSet)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "${local.cluster_id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"desired_size":          "2",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks": []map[string]string{
						{
							"size":         "100",
							"category":     "cloud_ssd",
							"auto_format":  "true",
							"file_system":  "ext4",
							"mount_target": "/var/lib/kubelet,/var/lib/containerd",
						},
						{
							"size":         "100",
							"category":     "cloud_ssd",
							"auto_format":  "true",
							"file_system":  "xfs",
							"mount_target": "/mnt/path3",
						},
					},
					"tags":              map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":        []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
					"image_type":        "AliyunLinux3",
					"deployment_set_id": "${alicloud_ecs_deployment_set.default.id}",
					"cpu_policy":        "none",
					"spot_strategy":     "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"desired_size":                 "2",
						"key_name":                     CHECKSET,
						"system_disk_category":         "cloud_efficiency",
						"system_disk_size":             "40",
						"install_cloud_monitor":        "false",
						"data_disks.#":                 CHECKSET,
						"tags.%":                       "2",
						"tags.Created":                 "TF",
						"tags.Foo":                     "Bar",
						"management.#":                 "1",
						"management.0.auto_repair":     "true",
						"management.0.auto_upgrade":    "true",
						"management.0.surge":           "0",
						"management.0.max_unavailable": "0",
						"image_type":                   "AliyunLinux3",
						"deployment_set_id":            CHECKSET,
						"cpu_policy":                   "none",
						"spot_strategy":                "NoSpot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "force_delete"},
			},
			// auto_format, mount_target cannot be modified if the custom mount path feature is enabled
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]string{
						{
							"size":         "100",
							"category":     "cloud_ssd",
							"auto_format":  "true",
							"file_system":  "ext4",
							"mount_target": "/var/lib/kubelet,/var/lib/containerd",
						},
						{
							"size":         "120",
							"category":     "cloud_ssd",
							"auto_format":  "true",
							"file_system":  "xfs",
							"mount_target": "/mnt/path3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#": CHECKSET,
					}),
				),
			},
			// check: remove nodes
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesNodePool_AttachInstances(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_AttachInstances)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"cluster_id":           "${local.cluster_id}",
					"vswitch_ids":          []string{"${local.vswitch_id}"},
					"instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_efficiency",
					"system_disk_size":     "40",
					"image_type":           "AliyunLinux3",
					"instances":            []string{"${alicloud_instance.default.0.id}"},
					"format_disk":          false,
					"keep_instance_name":   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "1",
						"key_name":             CHECKSET,
						"system_disk_category": "cloud_efficiency",
						"system_disk_size":     "40",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "force_delete", "instances", "format_disk", "keep_instance_name"},
			},
			// change, attach 1 and remove 0
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"instances":          []string{"${alicloud_instance.default.1.id}"},
			//		"format_disk":        true,
			//		"keep_instance_name": false,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			// attach one more instance
			{
				Config: testAccConfig(map[string]interface{}{
					"instances":          []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
					"format_disk":        true,
					"keep_instance_name": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// remove instance
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"instances": []string{"${alicloud_instance.default.1.id}"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
		},
	})
}

// auto_scaling has concurrent config conflict
func SkipTestAccAliCloudCSKubernetesNodePool_ScalingConflict(t *testing.T) {
	var v *cs.NodePoolDetail

	resourceId := "alicloud_cs_kubernetes_node_pool.autoscaling"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_ScalingConflict)

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
					"cluster_id":             "${local.cluster_id}",
					"vswitch_ids":            []string{"${local.vswitch_id}"},
					"instance_types":         []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"key_name":               "${alicloud_key_pair.default.key_name}",
					"system_disk_categories": []string{"cloud_efficiency", "cloud_essd"},
					"system_disk_size":       "40",
					"install_cloud_monitor":  "false",
					"image_type":             "AliyunLinux3",
					"scaling_policy":         "release",
					"scaling_config":         []map[string]string{{"enable": "true", "min_size": "1", "max_size": "10", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
					"cpu_policy":             "none",
					"spot_strategy":          "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                      name,
						"cluster_id":                                CHECKSET,
						"vswitch_ids.#":                             "1",
						"instance_types.#":                          "1",
						"key_name":                                  CHECKSET,
						"system_disk_categories.#":                  "2",
						"system_disk_size":                          "40",
						"install_cloud_monitor":                     "false",
						"image_type":                                "AliyunLinux3",
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "true",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "10",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
						"cpu_policy":                                "none",
						"spot_strategy":                             "NoSpot",
					}),
				),
			},
		},
	})
}

var csdKubernetesNodePoolBasicMap = map[string]string{
	"system_disk_size":     "40",
	"system_disk_category": "cloud_efficiency",
}

func resourceCSNodePoolConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	kubernetes_node_role       = "Worker"
    system_disk_category       = "cloud_efficiency"
    eni_amount           = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
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

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}
`, name)
}

func resourceCSNodePoolConfigDependence_Auto(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "YourPw123456"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

data "alicloud_instance_types" "cloud_essd" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
}

data "alicloud_instance_types" "cloud_auto" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_auto"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING-ACK$"
}

data "alicloud_vswitches" "vsw1" {
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
}

data "alicloud_vswitches" "vsw2" {
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 10)
}

data "alicloud_vswitches" "vsw3" {
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 11)
}

resource "alicloud_vswitch" "vsw1" {
  count      = length(data.alicloud_vswitches.vsw1.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vsw2" {
  count      = length(data.alicloud_vswitches.vsw2.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 10)
  zone_id    = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_vswitch" "vsw3" {
  count      = length(data.alicloud_vswitches.vsw3.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 11)
  zone_id    = data.alicloud_zones.default.zones.2.id
}

resource "alicloud_key_pair" "default" {
  key_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vsw1]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultrt8z7K" {
  time_points     = ["1", "22", "23"]
  repeat_weekdays = ["1", "2", "3"]
  name            = var.name

  retention_days = "-1"
}

resource "alicloud_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
}

resource "alicloud_ecs_capacity_reservation" "default" {
  count = 2
  description               = "terraform-example"
  platform                  = "linux"
  capacity_reservation_name = "terraform-example"
  end_time_type             = "Unlimited"
  instance_amount           = 1
  instance_type             = data.alicloud_instance_types.cloud_essd.instance_types.0.id
  match_criteria            = "Open"
  zone_ids = ["${data.alicloud_zones.default.zones.0.id}"]
}

data "alicloud_cs_kubernetes_version" "default" {
  cluster_type       = "ManagedKubernetes"
  profile            = "Default"
}

locals {
  cluster_id           = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
  vsw1                 = length(data.alicloud_vswitches.vsw1.ids) > 0 ? data.alicloud_vswitches.vsw1.ids[0] : concat(alicloud_vswitch.vsw1.*.id, [""])[0]
  vsw2                 = length(data.alicloud_vswitches.vsw2.ids) > 0 ? data.alicloud_vswitches.vsw2.ids[0] : concat(alicloud_vswitch.vsw2.*.id, [""])[0]
  vsw3                 = length(data.alicloud_vswitches.vsw3.ids) > 0 ? data.alicloud_vswitches.vsw3.ids[0] : concat(alicloud_vswitch.vsw3.*.id, [""])[0]
  capacity_reservations = alicloud_ecs_capacity_reservation.default.*.id
}
`, name)
}

func resourceCSNodePoolConfigDependence_Auto_Tee(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "password" {
    default = "YourPw123456"
}
data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
}

data "alicloud_vswitches" "vsw1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
}

resource "alicloud_vswitch" "vsw1" {
  count      = length(data.alicloud_vswitches.vsw1.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vsw1]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
  vsw1 = length(data.alicloud_vswitches.vsw1.ids) > 0 ? data.alicloud_vswitches.vsw1.ids[0] : concat(alicloud_vswitch.vsw1.*.id, [""])[0]
}
`, name)
}

// system disk encryption only support region HongKong zones B/C
func resourceCSNodePoolConfigDependence_BYOK(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_kms_keys" "default" {
  status  = "Enabled"
  filters = "[{\"Key\":\"CreatorType\", \"Values\":[\"User\"]}]"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING-ACK$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  security_group_type = "enterprise"
}

resource "alicloud_security_group" "group1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  security_group_type = "enterprise"
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default.*"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_ram_role" "role" {
  name        = var.name
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a role test."
}

resource "alicloud_ecs_disk" "default" {
  count     = 2
  disk_name = var.name
  zone_id   = data.alicloud_zones.default.zones.0.id
  category  = "cloud_efficiency"
  size      = "20"
}


data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
  instance_type  = data.alicloud_instance_types.default.instance_types.0.id
}

resource "alicloud_instance" "default" {
  count             = 2
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = "terraform-example"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.group.id]
  vswitch_id        = local.vswitch_id
}

resource "alicloud_ecs_disk_attachment" "default" {
  count       = 2
  disk_id     = element(alicloud_ecs_disk.default.*.id, count.index)
  instance_id = element(alicloud_instance.default.*.id, count.index)
  timeouts {
    delete = "5m"
  }
}

resource "alicloud_ecs_snapshot" "default" {
  count          = 2
  force          = "true"
  category       = "standard"
  description    = "terraform-example"
  disk_id        = element(alicloud_ecs_disk_attachment.default.*.disk_id, count.index)
  retention_days = "1"
  snapshot_name  = "terraform-example"
  tags = {
    Created = "TF"
    For     = "example"
  }
}
`, name)
}

func resourceCSNodePoolConfigDependence_Rds(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	kubernetes_node_role       = "Worker"
    system_disk_category       = "cloud_efficiency"
    eni_amount           = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
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

resource "alicloud_db_instance" "default" {
 count = 2
 engine               = "MySQL"
 engine_version       = "5.6"
 instance_type        = "rds.mysql.s2.large"
 instance_storage     = "30"
 instance_charge_type = "Postpaid"
 instance_name        = "tf-testacckubernetes"
 vswitch_id           = local.vswitch_id
 monitoring_period    = "60"
}

resource "alicloud_key_pair" "default" {
	key_name = var.name
}


resource "alicloud_key_pair" "default2" {
	key_name_prefix = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  is_enterprise_security_group = true
  slb_internet_enabled = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}


`, name)
}

func resourceCSNodePoolConfigDependence_DeploymentSet(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	kubernetes_node_role       = "Worker"
    system_disk_category       = "cloud_efficiency"
    eni_amount           = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
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

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

resource "alicloud_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
 name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
 count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
 name                 = var.name
 cluster_spec         = "ack.pro.small"
 worker_vswitch_ids   = [local.vswitch_id]
 new_nat_gateway      = false
 pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
 service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
 slb_internet_enabled = true
 is_enterprise_security_group = true
}

locals {
 vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
 cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}
`, name)
}

func resourceCSNodePoolConfigDependence_AttachInstances(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  kubernetes_node_role       = "Worker"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
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

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  security_group_type = "enterprise"
}

data "alicloud_images" "default" {
  name_regex = "^aliyun_3_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_instance" "default" {
  count             = 2
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = "terraform-example"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.group.id]
  vswitch_id        = local.vswitch_id
  system_disk_category       = "cloud_efficiency"
  lifecycle {
    ignore_changes = [user_data, instance_name, image_id, security_groups, tags]
  }
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
 name_regex = "^Default"
 enable_details = true
}

resource "alicloud_cs_managed_kubernetes" "default" {
 count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
 name                 = var.name
 cluster_spec         = "ack.pro.small"
 worker_vswitch_ids   = [local.vswitch_id]
 new_nat_gateway      = false
 pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
 service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
 slb_internet_enabled = true
 is_enterprise_security_group = true
}

locals {
 vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
 cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
 cluster_sg_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.security_group_id : alicloud_cs_managed_kubernetes.default.0.security_group_id
}
`, name)
}

func resourceCSNodePoolConfigDependence_ScalingConflict(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING-ACK$"
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

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  deletion_protection  = false
  is_enterprise_security_group = true
}

resource "alicloud_cs_kubernetes_node_pool" "spot_auto_scaling" {
  name                 = "spot_auto_scaling"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [local.vswitch_id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # spot price config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.default.instance_types.0.id
    price_limit   = "0.70"
  }
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}
`, name)
}

// Case 节点池测试_spot_instance 5288
func TestAccAliCloudAckNodepool_basic5288(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5288)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name,
					"cluster_id":     "${local.cluster_id}",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}", "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}"},
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":       name,
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "2",
						"system_disk_category": "cloud_essd",
						"system_disk_size":     "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_az_policy": "COST_OPTIMIZED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_az_policy": "COST_OPTIMIZED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_policy": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_policy": "none",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_instance_remedy": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"on_demand_base_capacity":                  "1",
					"on_demand_percentage_above_base_capacity": "20",
					"compensate_with_on_demand":                "true",
					"spot_instance_pools":                      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity":                  "1",
						"on_demand_percentage_above_base_capacity": "20",
						"compensate_with_on_demand":                "true",
						"spot_instance_pools":                      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "SpotWithPriceLimit",
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}",
							"price_limit":   "0.96",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":      "SpotWithPriceLimit",
						"spot_price_limit.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy":    "SpotAsPriceGo",
					"spot_price_limit": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":      "SpotAsPriceGo",
						"spot_price_limit.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "PayByTraffic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByTraffic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${local.vsw1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}", "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"on_demand_base_capacity": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "SpotAsPriceGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy": "SpotAsPriceGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compensate_with_on_demand": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compensate_with_on_demand": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"on_demand_percentage_above_base_capacity": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_percentage_above_base_capacity": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_instance_pools": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_pools": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "PayByBandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByBandwidth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_instance_remedy": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "false",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "COST_OPTIMIZED",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "0",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}", "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}"},
					"spot_instance_remedy":    "true",
					"on_demand_base_capacity": "1",
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}",
							"price_limit":   "0.96",
						},
					},
					"spot_strategy":                            "SpotWithPriceLimit",
					"compensate_with_on_demand":                "true",
					"on_demand_percentage_above_base_capacity": "20",
					"spot_instance_pools":                      "2",
					"internet_charge_type":                     "PayByTraffic",
					"internet_max_bandwidth_out":               "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                           name + "_update",
						"cluster_id":                               CHECKSET,
						"instance_charge_type":                     "PostPaid",
						"system_disk_category":                     "cloud_essd",
						"vswitch_ids.#":                            "1",
						"install_cloud_monitor":                    "false",
						"login_as_non_root":                        "true",
						"system_disk_size":                         "120",
						"multi_az_policy":                          "COST_OPTIMIZED",
						"cpu_policy":                               "none",
						"runtime_version":                          CHECKSET,
						"desired_size":                             "0",
						"instance_types.#":                         "2",
						"spot_instance_remedy":                     "true",
						"on_demand_base_capacity":                  "1",
						"spot_price_limit.#":                       "2",
						"spot_strategy":                            "SpotWithPriceLimit",
						"compensate_with_on_demand":                "true",
						"on_demand_percentage_above_base_capacity": "20",
						"spot_instance_pools":                      "2",
						"internet_charge_type":                     "PayByTraffic",
						"internet_max_bandwidth_out":               "5",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5288 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence5288(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster" {
  default = "c0196d207b11d4a25ae4cad2a6f029a38"
}

variable "password" {
  default = "spot_instance"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_resource_manager_resource_group" "default" {
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.9.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vsw2" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.10.0/24"
  zone_id    = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_vswitch" "vsw3" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.11.0/24"
  zone_id    = data.alicloud_zones.default.zones.2.id
}


`, name)
}

// Case 节点池测试_kubelet 5291
func TestAccAliCloudAckNodepool_basic5291_old(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5291)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name,
					"cluster_id":     "${local.cluster_id}",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"system_disk_category": "cloud_auto",
					"system_disk_size":     "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":       name,
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "1",
						"system_disk_category": "cloud_auto",
						"system_disk_size":     "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_az_policy": "PRIORITY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_az_policy": "PRIORITY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_policy": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_policy": "none",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
								// eviction check failed
								//"\"RotateKubeletServerCertificate\"": "true",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"cpu_manager_policy": "none",
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":           "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"registry_pull_qps": "30",
							"registry_burst":    "10",
							"event_record_qps":  "40",
							"kube_api_burst":    "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "Target",
							"private_pool_options_id":             "${local.capacity_reservations[1]}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_provisioned_iops": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_provisioned_iops": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_bursting_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// Clusters with version greater than or equal to 1.34 cannot set 'unschedulable' to true
					"unschedulable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unschedulable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"read_only_port":          "10000",
							"max_pods":                "10",
							"container_log_max_size":  "2Mi",
							"container_log_max_files": "10",
							"feature_gates": map[string]interface{}{
								"\"APIListChunking\"":                "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*", "net.ipv4.route.min_pmtu", "kernel.shm"},
							"registry_pull_qps": "10",
							"registry_burst":    "20",
							"event_record_qps":  "10",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "20%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "20%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "2Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "2m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "20Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":               "0.5",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"event_burst":           "40",
							"kube_api_qps":          "22",
							"serialize_image_pulls": "false",
							"cpu_manager_policy":    "static",
							"kube_api_burst":        "25",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"private_pool_options": []map[string]interface{}{
			//			{
			//				"private_pool_options_match_criteria": "None",
			//				"private_pool_options_id":             "${local.capacity_reservations[0]}",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_provisioned_iops": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_provisioned_iops": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_bursting_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_bursting_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unschedulable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unschedulable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"cis_enabled":          "true",
					"system_disk_category": "cloud_auto",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "1",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"cpu_manager_policy": "none",
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":           "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"registry_pull_qps": "30",
							"registry_burst":    "10",
							"event_record_qps":  "40",
							"kube_api_burst":    "20",
						},
					},
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "Target",
							"private_pool_options_id":             "${local.capacity_reservations[1]}",
						},
					},
					"system_disk_provisioned_iops": "100",
					"system_disk_bursting_enabled": "true",
					"unschedulable":                "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":               name + "_update",
						"cluster_id":                   CHECKSET,
						"instance_charge_type":         "PostPaid",
						"cis_enabled":                  "true",
						"system_disk_category":         "cloud_auto",
						"vswitch_ids.#":                "1",
						"install_cloud_monitor":        "true",
						"system_disk_size":             "120",
						"multi_az_policy":              "PRIORITY",
						"cpu_policy":                   "none",
						"runtime_version":              CHECKSET,
						"desired_size":                 "1",
						"instance_types.#":             "1",
						"system_disk_provisioned_iops": "100",
						"system_disk_bursting_enabled": "true",
						"unschedulable":                "false",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5291 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

// Case 节点池测试-PrePaid 5266
func TestAccAliCloudAckNodepool_basic5266(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5266)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name,
					"cluster_id":     "${local.cluster_id}",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "120",
					"desired_size":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":       name,
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "1",
						"system_disk_category": "cloud_essd",
						"system_disk_size":     "120",
						"desired_size":         "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
					"period":               "1",
					"auto_renew_period":    "1",
					"auto_renew":           "true",
					"period_unit":          "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"period":               "1",
						"auto_renew_period":    "1",
						"auto_renew":           "true",
						"period_unit":          "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_az_policy": "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_az_policy": "BALANCE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_policy": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_policy": "none",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_performance_level": "PL0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_performance_level": "PL0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_encrypted": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_encrypted": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PrePaid",
					"auto_renew_period":    "1",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "BALANCE",
					"cpu_policy":            "none",
					"period":                "1",
					"tee_config": []map[string]interface{}{
						{
							"tee_enable": "false",
						},
					},
					"system_disk_performance_level": "PL0",
					"runtime_version":               "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "test",
						},
					},
					"period_unit":           "Month",
					"desired_size":          "0",
					"auto_renew":            "false",
					"system_disk_encrypted": "false",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                name + "_update",
						"cluster_id":                    CHECKSET,
						"instance_charge_type":          "PrePaid",
						"auto_renew_period":             "1",
						"system_disk_category":          "cloud_essd",
						"vswitch_ids.#":                 "1",
						"install_cloud_monitor":         "true",
						"login_as_non_root":             "true",
						"system_disk_size":              "120",
						"multi_az_policy":               "BALANCE",
						"cpu_policy":                    "none",
						"period":                        "1",
						"system_disk_performance_level": "PL0",
						"runtime_version":               CHECKSET,
						"labels.#":                      "1",
						"period_unit":                   "Month",
						"desired_size":                  "0",
						"auto_renew":                    "false",
						"system_disk_encrypted":         "false",
						"instance_types.#":              "1",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5266 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence5266(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster" {
  default = "c0196d207b11d4a25ae4cad2a6f029a38"
}

variable "password" {
  default = "tf-example123456"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.4.4.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}


`, name)
}

// Case 节点池测试_basic 5172
func TestAccAliCloudAckNodepool_basic5172(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name,
					"cluster_id":     "${local.cluster_id}",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "120",
					"desired_size":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":       name,
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "1",
						"system_disk_category": "cloud_essd",
						"system_disk_size":     "120",
						"desired_size":         "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
									"vul_level":    "asap",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_az_policy": "PRIORITY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_az_policy": "PRIORITY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_policy": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_policy": "none",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]interface{}{
						{
							"category":                "cloud_ssd",
							"encrypted":               "true",
							"size":                    "40",
							"auto_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
						},
						{
							"category":         "cloud_auto",
							"bursting_enabled": "true",
							"provisioned_iops": "100",
							"size":             "100",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "example",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"taints": []map[string]interface{}{
						{
							"key":    "tf",
							"effect": "NoSchedule",
							"value":  "example",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"taints.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${var.password}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_policy": "release",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_policy": "release",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unschedulable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unschedulable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pre_user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pre_user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_snapshot_policy_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${local.vsw1}", "${local.vsw2}", "${local.vsw3}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
						"${data.alicloud_instance_types.cloud_essd.instance_types.1.id}",
						"${data.alicloud_instance_types.cloud_essd.instance_types.2.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL1",
							"size":              "40",
						},
						{
							"category":         "cloud_auto",
							"bursting_enabled": "false",
							"provisioned_iops": "2000",
							"size":             "40",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "label1",
							"value": "value1",
						},
						{
							"key":   "label2",
							"value": "value2",
						},
						{
							"key":   "label3",
							"value": "value3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"taints": []map[string]interface{}{
						{
							"key":    "taint1",
							"effect": "NoSchedule",
						},
						{
							"key":    "taint2",
							"effect": "NoSchedule",
						},
						{
							"key":    "taint3",
							"effect": "NoSchedule",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"taints.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${local.vsw1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data":     "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
					"pre_user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":     "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
						"pre_user_data": "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${local.cluster_id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
									"vul_level":    "asap",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "80",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "0",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"data_disks": []map[string]interface{}{
						{
							"category":                "cloud_ssd",
							"encrypted":               "true",
							"size":                    "40",
							"auto_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
						},
						{
							"category":         "cloud_auto",
							"bursting_enabled": "true",
							"provisioned_iops": "100",
							"size":             "100",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
						},
					},
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "example",
						},
					},
					"taints": []map[string]interface{}{
						{
							"key":    "tf",
							"effect": "NoSchedule",
							"value":  "example",
						},
					},
					"password":                       "${var.password}",
					"scaling_policy":                 "release",
					"deployment_set_id":              "${alicloud_ecs_deployment_set.default.id}",
					"node_name_mode":                 "customized,aliyun,ip,com",
					"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
					"unschedulable":                  "false",
					"pre_user_data":                  "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                 name + "_update",
						"cluster_id":                     CHECKSET,
						"resource_group_id":              CHECKSET,
						"instance_charge_type":           "PostPaid",
						"system_disk_category":           "cloud_essd",
						"vswitch_ids.#":                  "1",
						"install_cloud_monitor":          "true",
						"login_as_non_root":              "true",
						"system_disk_size":               "80",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "none",
						"runtime_version":                CHECKSET,
						"desired_size":                   "0",
						"instance_types.#":               "1",
						"data_disks.#":                   "3",
						"labels.#":                       "1",
						"taints.#":                       "1",
						"password":                       CHECKSET,
						"scaling_policy":                 "release",
						"deployment_set_id":              CHECKSET,
						"node_name_mode":                 "customized,aliyun,ip,com",
						"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
						"unschedulable":                  "false",
						"pre_user_data":                  "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgcHJlLXVzZXItZGF0YSBhZnRlciBtb2RpZmllZCI=",
						"system_disk_snapshot_policy_id": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5172 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence5172(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster" {
  default = "c0196d207b11d4a25ae4cad2a6f029a38"
}

variable "password" {
  default = "tf-example123456"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_resource_manager_resource_group" "default" {
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.9.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vsw2" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.10.0/24"
  zone_id    = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_vswitch" "vsw3" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.11.0/24"
  zone_id    = data.alicloud_zones.default.zones.2.id
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultrt8z7K" {
  time_points               = ["1", "22", "23"]
  repeat_weekdays           = ["1", "2", "3"]
  auto_snapshot_policy_name = var.name
  retention_days            = "-1"
}

resource "alicloud_ecs_deployment_set" "default" {
  group_count         = "3"
  strategy            = "Availability"
  deployment_set_name = var.name
}


`, name)
}

// Case 节点池测试_soc 5401
func TestAccAliCloudAckNodepool_basic5401(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5401)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"node_pool_name":       name,
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "120",
					"desired_size":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":           CHECKSET,
						"vswitch_ids.#":        "1",
						"instance_types.#":     "1",
						"node_pool_name":       name,
						"system_disk_category": "cloud_essd",
						"system_disk_size":     "120",
						"desired_size":         "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${var.password}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_efficiency",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_efficiency",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_name": "containerd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_name": "containerd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"soc_enabled":          "true",
					"instance_charge_type": "PostPaid",
					"login_as_non_root":    "false",
					"desired_size":         "0",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"system_disk_size": "60",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"password":             "${var.password}",
					"system_disk_category": "cloud_efficiency",
					"runtime_name":         "containerd",
					"runtime_version":      "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"node_pool_name":       name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"soc_enabled":          "true",
						"instance_charge_type": "PostPaid",
						"login_as_non_root":    "false",
						"desired_size":         "0",
						"vswitch_ids.#":        "1",
						"system_disk_size":     "60",
						"instance_types.#":     "1",
						"password":             CHECKSET,
						"system_disk_category": "cloud_efficiency",
						"runtime_name":         "containerd",
						"runtime_version":      CHECKSET,
						"node_pool_name":       name + "_update",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5401 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence5401(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster" {
  default = "c0196d207b11d4a25ae4cad2a6f029a38"
}

variable "password" {
  default = "tf-example123456"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.9.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}


`, name)
}

// Case 节点池测试_teeconfig 5628
func SkipTestAccAliCloudAckNodepool_basic5628(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5628)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto_Tee)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name": name,
					"instance_types": []string{
						"ecs.c7t.xlarge", "ecs.g7t.xlarge", "ecs.r7t.xlarge"},
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"cluster_id":           "${local.cluster_id}",
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "120",
					"image_type":           "AliyunLinuxSecurity",
					"desired_size":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":       name,
						"instance_types.#":     "3",
						"vswitch_ids.#":        "1",
						"cluster_id":           CHECKSET,
						"system_disk_category": "cloud_essd",
						"system_disk_size":     "120",
						"image_type":           "AliyunLinuxSecurity",
						"desired_size":         "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_name": "containerd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_name": "containerd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
					"runtime_name":         "containerd",
					"tee_config": []map[string]interface{}{
						{
							"tee_enable": "true",
						},
					},
					"node_pool_name": name + "_update",
					"desired_size":   "0",
					"instance_types": []string{
						"ecs.c7t.xlarge", "ecs.g7t.xlarge", "ecs.r7t.xlarge"},
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"cluster_id":           "${local.cluster_id}",
					"system_disk_size":     "40",
					"system_disk_category": "cloud_essd",
					"image_type":           "AliyunLinuxSecurity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"runtime_name":         "containerd",
						"node_pool_name":       name + "_update",
						"desired_size":         "0",
						"instance_types.#":     "3",
						"vswitch_ids.#":        "1",
						"cluster_id":           CHECKSET,
						"system_disk_size":     "40",
						"system_disk_category": "cloud_essd",
						"image_type":           "AliyunLinuxSecurity",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap5628 = map[string]string{
	"platform":                 CHECKSET,
	"instance_charge_type":     "PostPaid",
	"runtime_name":             CHECKSET,
	"image_type":               CHECKSET,
	"tee_config.#":             CHECKSET,
	"node_name_mode":           CHECKSET,
	"image_id":                 CHECKSET,
	"multi_az_policy":          CHECKSET,
	"cpu_policy":               CHECKSET,
	"runtime_version":          CHECKSET,
	"security_group_ids.#":     CHECKSET,
	"node_pool_id":             CHECKSET,
	"system_disk_categories.#": CHECKSET,
	"spot_strategy":            CHECKSET,
	"scaling_policy":           CHECKSET,
	"scaling_config.#":         CHECKSET,
	"security_group_id":        CHECKSET,
	"management.#":             CHECKSET,
	"system_disk_category":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence5628(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster" {
  default = "c0196d207b11d4a25ae4cad2a6f029a38"
}

variable "password" {
  default = "tf-example123456"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id     = "vpc-bp1b444zex5kv0jwh0je4"
  cidr_block = "10.0.9.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}


`, name)
}

// Case 节点池测试_spot_instance 5288  twin
func TestAccAliCloudAckNodepool_basic5288_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5288)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"auto_renew_period":    "0",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "false",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "COST_OPTIMIZED",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "0",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
						"${data.alicloud_instance_types.cloud_essd.instance_types.1.id}"},
					"spot_instance_remedy":    "false",
					"on_demand_base_capacity": "2",
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "${data.alicloud_instance_types.cloud_essd.instance_types.1.id}",
							"price_limit":   "0.96",
						},
					},
					"spot_strategy":                            "SpotAsPriceGo",
					"compensate_with_on_demand":                "false",
					"on_demand_percentage_above_base_capacity": "30",
					"spot_instance_pools":                      "1",
					"internet_charge_type":                     "PayByBandwidth",
					"internet_max_bandwidth_out":               "10",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                           name,
						"cluster_id":                               CHECKSET,
						"instance_charge_type":                     "PostPaid",
						"system_disk_category":                     "cloud_essd",
						"vswitch_ids.#":                            "1",
						"install_cloud_monitor":                    "false",
						"login_as_non_root":                        "true",
						"system_disk_size":                         "120",
						"multi_az_policy":                          "COST_OPTIMIZED",
						"cpu_policy":                               "none",
						"runtime_version":                          CHECKSET,
						"desired_size":                             "0",
						"instance_types.#":                         "2",
						"spot_instance_remedy":                     "false",
						"on_demand_base_capacity":                  "2",
						"spot_price_limit.#":                       "2",
						"spot_strategy":                            "SpotAsPriceGo",
						"compensate_with_on_demand":                "false",
						"on_demand_percentage_above_base_capacity": "30",
						"spot_instance_pools":                      "1",
						"internet_charge_type":                     "PayByBandwidth",
						"internet_max_bandwidth_out":               "10",
						"tags.%":                                   "2",
						"tags.Created":                             "TF",
						"tags.For":                                 "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试_kubelet 5291  twin
func TestAccAliCloudAckNodepool_basic5291_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5291)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"cis_enabled":          "true",
					"system_disk_category": "cloud_auto",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "1",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
					},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "40",
							"kube_api_qps":          "22",
							"serialize_image_pulls": "false",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "20%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "20%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "20Gi",
							},
							"cpu_manager_policy": "static",
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "2Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "2m30s",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":               "0.5",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"read_only_port":          "10000",
							"max_pods":                "10",
							"container_log_max_size":  "2Mi",
							"container_log_max_files": "10",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"": "true",
								"\"APIListChunking\"":      "true",
							},
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*", "net.ipv4.route.min_pmtu", "kernel.shm"},
							"registry_pull_qps": "10",
							"registry_burst":    "20",
							"event_record_qps":  "10",
							"kube_api_burst":    "25",
						},
					},
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "Target",
							"private_pool_options_id":             "${local.capacity_reservations[0]}",
						},
					},
					"system_disk_provisioned_iops": "200",
					"unschedulable":                "false",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":               name,
						"cluster_id":                   CHECKSET,
						"instance_charge_type":         "PostPaid",
						"cis_enabled":                  "true",
						"system_disk_category":         "cloud_auto",
						"vswitch_ids.#":                "1",
						"install_cloud_monitor":        "true",
						"login_as_non_root":            "true",
						"system_disk_size":             "120",
						"multi_az_policy":              "PRIORITY",
						"cpu_policy":                   "none",
						"runtime_version":              CHECKSET,
						"desired_size":                 "1",
						"instance_types.#":             "1",
						"system_disk_provisioned_iops": "200",
						"unschedulable":                "false",
						"tags.%":                       "2",
						"tags.Created":                 "TF",
						"tags.For":                     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试-PrePaid 5266  twin
func TestAccAliCloudAckNodepool_basic5266_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5266)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PrePaid",
					"auto_renew_period":    "2",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "BALANCE",
					"cpu_policy":            "none",
					"period":                "6",
					"tee_config": []map[string]interface{}{
						{
							"tee_enable": "false",
						},
					},
					"system_disk_performance_level": "PL0",
					"runtime_version":               "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "test",
						},
					},
					"period_unit":           "Month",
					"desired_size":          "0",
					"auto_renew":            "true",
					"system_disk_encrypted": "false",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                name,
						"cluster_id":                    CHECKSET,
						"instance_charge_type":          "PrePaid",
						"auto_renew_period":             "2",
						"system_disk_category":          "cloud_essd",
						"vswitch_ids.#":                 "1",
						"install_cloud_monitor":         "true",
						"login_as_non_root":             "true",
						"system_disk_size":              "120",
						"multi_az_policy":               "BALANCE",
						"cpu_policy":                    "none",
						"period":                        "6",
						"system_disk_performance_level": "PL0",
						"runtime_version":               CHECKSET,
						"labels.#":                      "1",
						"period_unit":                   "Month",
						"desired_size":                  "0",
						"auto_renew":                    "true",
						"system_disk_encrypted":         "false",
						"instance_types.#":              "1",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试_basic 5172  twin
func TestAccAliCloudAckNodepool_basic5172_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
									"vul_level":    "asap",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vsw1}", "${local.vsw2}", "${local.vsw3}"},
					"install_cloud_monitor": "false",
					"login_as_non_root":     "true",
					"system_disk_size":      "100",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "0",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
						"${data.alicloud_instance_types.cloud_essd.instance_types.1.id}",
						"${data.alicloud_instance_types.cloud_essd.instance_types.2.id}",
					},
					"data_disks": []map[string]interface{}{
						{
							"category":                "cloud_essd",
							"encrypted":               "true",
							"size":                    "40",
							"auto_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
							"performance_level":       "PL1",
						},
						{
							"category":         "cloud_auto",
							"bursting_enabled": "false",
							"provisioned_iops": "2000",
							"size":             "40",
						},
						{
							"category":          "cloud_essd",
							"performance_level": "PL0",
							"size":              "40",
						},
					},
					"labels": []map[string]interface{}{
						{
							"key":   "label1",
							"value": "value1",
						},
						{
							"key":   "label2",
							"value": "value2",
						},
						{
							"key":   "label3",
							"value": "value3",
						},
					},
					"taints": []map[string]interface{}{
						{
							"key":    "taint1",
							"effect": "NoSchedule",
							"value":  "example",
						},
						{
							"key":    "taint2",
							"effect": "NoSchedule",
						},
						{
							"key":    "taint3",
							"effect": "NoSchedule",
						},
					},
					"password":                       "${var.password}",
					"scaling_policy":                 "release",
					"deployment_set_id":              "${alicloud_ecs_deployment_set.default.id}",
					"node_name_mode":                 "customized,aliyun,ip,com",
					"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
					"unschedulable":                  "false",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
					"system_disk_performance_level":  "PL0",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                 name,
						"cluster_id":                     CHECKSET,
						"resource_group_id":              CHECKSET,
						"instance_charge_type":           "PostPaid",
						"system_disk_category":           "cloud_essd",
						"vswitch_ids.#":                  "3",
						"install_cloud_monitor":          "false",
						"login_as_non_root":              "true",
						"system_disk_size":               "100",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "none",
						"runtime_version":                CHECKSET,
						"desired_size":                   "0",
						"instance_types.#":               "3",
						"data_disks.#":                   "3",
						"labels.#":                       "3",
						"taints.#":                       "3",
						"password":                       CHECKSET,
						"scaling_policy":                 "release",
						"deployment_set_id":              CHECKSET,
						"node_name_mode":                 "customized,aliyun,ip,com",
						"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
						"unschedulable":                  "false",
						"system_disk_snapshot_policy_id": CHECKSET,
						"system_disk_performance_level":  "PL0",
						"tags.%":                         "2",
						"tags.Created":                   "TF",
						"tags.For":                       "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试_soc 5401  twin
func TestAccAliCloudAckNodepool_basic5401_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5401)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":           "${local.cluster_id}",
					"soc_enabled":          "true",
					"instance_charge_type": "PostPaid",
					"login_as_non_root":    "false",
					"desired_size":         "0",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"system_disk_size": "60",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_efficiency.instance_types.0.id}"},
					"password":             "${var.password}",
					"system_disk_category": "cloud_efficiency",
					"runtime_name":         "containerd",
					"runtime_version":      "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"node_pool_name":       name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":           CHECKSET,
						"soc_enabled":          "true",
						"instance_charge_type": "PostPaid",
						"login_as_non_root":    "false",
						"desired_size":         "0",
						"vswitch_ids.#":        "1",
						"system_disk_size":     "60",
						"instance_types.#":     "1",
						"password":             CHECKSET,
						"system_disk_category": "cloud_efficiency",
						"runtime_name":         "containerd",
						"runtime_version":      CHECKSET,
						"node_pool_name":       name,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试_teeconfig 5628  twin
func SkipTestAccAliCloudAckNodepool_basic5628_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5628)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto_Tee)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
					"runtime_name":         "containerd",
					"tee_config": []map[string]interface{}{
						{
							"tee_enable": "true",
						},
					},
					"node_pool_name": name,
					"desired_size":   "0",
					"instance_types": []string{
						"ecs.c7t.xlarge", "ecs.g7t.xlarge", "ecs.r7t.xlarge"},
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"cluster_id":           "${local.cluster_id}",
					"system_disk_size":     "40",
					"system_disk_category": "cloud_essd",
					"image_type":           "AliyunLinuxSecurity",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"runtime_name":         "containerd",
						"node_pool_name":       name,
						"desired_size":         "0",
						"instance_types.#":     "3",
						"vswitch_ids.#":        "1",
						"cluster_id":           CHECKSET,
						"system_disk_size":     "40",
						"system_disk_category": "cloud_essd",
						"image_type":           "AliyunLinuxSecurity",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Case 节点池测试_kubelet 5291
func TestAccAliCloudAckNodepool_basic5291(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap5291)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sacknodepool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Auto)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"management": []map[string]interface{}{
						{
							"enable":          "true",
							"auto_repair":     "true",
							"auto_upgrade":    "true",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_vul_fix": "true",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "true",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "true",
								},
							},
						},
					},
					"system_disk_category": "cloud_auto",
					"vswitch_ids": []string{
						"${local.vsw1}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "${data.alicloud_cs_kubernetes_version.default.metadata.0.runtime.0.version}",
					"desired_size":          "1",
					"instance_types": []string{
						"${data.alicloud_instance_types.cloud_essd.instance_types.0.id}",
					},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"cpu_manager_policy": "none",
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":    "true",
								"\"CustomCPUCFSQuotaPeriod\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"registry_pull_qps":              "30",
							"registry_burst":                 "10",
							"event_record_qps":               "40",
							"kube_api_burst":                 "20",
							"container_log_monitor_interval": "30s",
							"container_log_max_workers":      "5",
							"tracing": []map[string]interface{}{
								{
									"sampling_rate_per_million": "1000",
									"endpoint":                  "localhost:4318",
								},
							},
							"topology_manager_policy":         "best-effort",
							"pod_pids_limit":                  "1024",
							"memory_manager_policy":           "None",
							"image_gc_low_threshold_percent":  "10",
							"image_gc_high_threshold_percent": "80",
							"cpu_cfs_quota_period":            "100ms",
							"cpu_cfs_quota":                   "true",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
							},
							"cluster_dns": []string{
								"8.8.8.8"},
						},
					},
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "Target",
							"private_pool_options_id":             "${local.capacity_reservations[1]}",
						},
					},
					"system_disk_provisioned_iops": "100",
					"system_disk_bursting_enabled": "true",
					"unschedulable":                "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":               name,
						"cluster_id":                   CHECKSET,
						"instance_charge_type":         "PostPaid",
						"system_disk_category":         "cloud_auto",
						"vswitch_ids.#":                "1",
						"install_cloud_monitor":        "true",
						"login_as_non_root":            "true",
						"system_disk_size":             "120",
						"multi_az_policy":              "PRIORITY",
						"cpu_policy":                   "none",
						"desired_size":                 "1",
						"instance_types.#":             "1",
						"system_disk_provisioned_iops": "100",
						"system_disk_bursting_enabled": "true",
						"unschedulable":                "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"read_only_port":          "10000",
							"max_pods":                "10",
							"container_log_max_size":  "2Mi",
							"container_log_max_files": "10",
							"feature_gates": map[string]interface{}{
								"\"APIListChunking\"":         "true",
								"\"CustomCPUCFSQuotaPeriod\"": "true",
								"\"MemoryManager\"":           "true",
							},
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*", "net.ipv4.route.min_pmtu", "kernel.shm"},
							"registry_pull_qps": "10",
							"registry_burst":    "20",
							"event_record_qps":  "10",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "20%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "20%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "2Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "2m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "20Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":               "0.5",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"event_burst":                    "40",
							"kube_api_qps":                   "22",
							"serialize_image_pulls":          "false",
							"cpu_manager_policy":             "static",
							"kube_api_burst":                 "25",
							"container_log_monitor_interval": "10s",
							"container_log_max_workers":      "3",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4317",
									"sampling_rate_per_million": "10000",
								},
							},
							"topology_manager_policy": "restricted",
							"pod_pids_limit":          "800",
							"cluster_dns": []string{
								"8.8.8.8", "4.4.4.4"},
							"image_gc_low_threshold_percent":  "15",
							"image_gc_high_threshold_percent": "75",
							"cpu_cfs_quota_period":            "50ms",
							"cpu_cfs_quota":                   "false",
							"memory_manager_policy":           "Static",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "3Gi",
									},
								},
							},
						},
					},
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "Target",
							"private_pool_options_id":             "${local.capacity_reservations[1]}",
						},
					},
					"system_disk_provisioned_iops": "200",
					"unschedulable":                "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_provisioned_iops": "200",
						"unschedulable":                "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*"},
							"cpu_manager_policy": "none",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

// Test Ack Nodepool. >>> Resource test cases, automatically generated.
// Case np-instancepattern 11802
func TestAccAliCloudAckNodepool_basic11802(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap11802)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence11802)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":          "PostPaid",
					"runtime_name":                  "containerd",
					"image_type":                    "AliyunLinux3",
					"system_disk_encrypt_algorithm": "aes-256",
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "None",
						},
					},
					"system_disk_performance_level": "PL0",
					"node_name_mode":                "customized,aliyun,ip,com",
					"image_id":                      "aliyun_3_x64_20G_alibase_20240819.vhd",
					"install_cloud_monitor":         "false",
					"multi_az_policy":               "PRIORITY",
					"cpu_policy":                    "static",
					"node_pool_name":                name,
					"runtime_version":               "1.6.20",
					"spot_instance_pools":           "1",
					"labels": []map[string]interface{}{
						{
							"key":   "test",
							"value": "test",
						},
					},
					"taints": []map[string]interface{}{
						{
							"key":    "test_taint_key",
							"effect": "NoSchedule",
							"value":  "test_taint_val",
						},
					},
					"internet_charge_type":       "PayByBandwidth",
					"internet_max_bandwidth_out": "1",
					"login_as_non_root":          "false",
					"desired_size":               "0",
					"system_disk_size":           "40",
					"spot_strategy":              "NoSpot",
					"password":                   "Zsyh1234!",
					"scaling_policy":             "release",
					"vswitch_ids": []string{
						"${alicloud_vswitch.defaultziRRat.id}", "${alicloud_vswitch.defaultT8D8ss.id}", "${alicloud_vswitch.defaultFsk7cj.id}"},
					"instance_types":               []string{},
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"compensate_with_on_demand":    "false",
					"spot_instance_remedy":         "false",
					"soc_enabled":                  "false",
					"system_disk_provisioned_iops": "3000",
					"system_disk_encrypted":        "false",
					"security_hardening_os":        "false",
					"system_disk_categories": []string{
						"cloud_essd", "cloud_auto"},
					"system_disk_bursting_enabled": "false",
					"user_data":                    "${var.user_data}",
					"pre_user_data":                "${var.user_data}",
					"scaling_config": []map[string]interface{}{
						{
							"enable":   "false",
							"type":     "cpu",
							"max_size": "10",
							"min_size": "1",
						},
					},
					"unschedulable":     "false",
					"deployment_set_id": "${alicloud_ecs_deployment_set.defaultACvZwq.id}",
					"management": []map[string]interface{}{
						{
							"enable":          "false",
							"auto_repair":     "false",
							"auto_upgrade":    "false",
							"surge":           "0",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_vul_fix": "false",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "false",
								},
							},
						},
					},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"cpu_manager_policy":    "none",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":    "true",
								"\"CustomCPUCFSQuotaPeriod\"": "true",
							},
							"registry_pull_qps":              "30",
							"registry_burst":                 "10",
							"event_record_qps":               "40",
							"kube_api_burst":                 "20",
							"container_log_monitor_interval": "30s",
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"container_log_max_workers": "5",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4318",
									"sampling_rate_per_million": "1000",
								},
							},
							"topology_manager_policy": "best-effort",
							"pod_pids_limit":          "1024",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "1",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "2",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
							},
							"memory_manager_policy": "None",
							"cluster_dns": []string{
								"8.8.8.8", "10.96.0.10", "192.168.0.100"},
							"image_gc_low_threshold_percent":  "10",
							"image_gc_high_threshold_percent": "80",
							"cpu_cfs_quota_period":            "100ms",
							"cpu_cfs_quota":                   "true",
						},
					},
					"instance_patterns": []map[string]interface{}{
						{
							"instance_family_level": "EnterpriseLevel",
							"excluded_instance_types": []string{
								"ecs.c6.*"},
							"min_cpu_cores":   "4",
							"max_cpu_cores":   "8",
							"min_memory_size": "8",
							"max_memory_size": "16",
							"instance_categories": []string{
								"General-purpose"},
							"cpu_architectures": []string{
								"X86"},
						},
					},
					"cluster_id":                     "${alicloud_cs_managed_kubernetes.defaultC02XDz.id}",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default4L8w9Z.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":           "PostPaid",
						"runtime_name":                   "containerd",
						"image_type":                     "AliyunLinux3",
						"system_disk_encrypt_algorithm":  "aes-256",
						"system_disk_performance_level":  "PL0",
						"node_name_mode":                 "customized,aliyun,ip,com",
						"image_id":                       "aliyun_3_x64_20G_alibase_20240819.vhd",
						"install_cloud_monitor":          "false",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "static",
						"node_pool_name":                 name,
						"runtime_version":                "1.6.20",
						"spot_instance_pools":            "1",
						"labels.#":                       "1",
						"taints.#":                       "1",
						"internet_charge_type":           "PayByBandwidth",
						"internet_max_bandwidth_out":     "1",
						"login_as_non_root":              "false",
						"desired_size":                   "0",
						"system_disk_size":               "40",
						"spot_strategy":                  "NoSpot",
						"password":                       "Zsyh1234!",
						"scaling_policy":                 "release",
						"vswitch_ids.#":                  "3",
						"instance_types.#":               "0",
						"resource_group_id":              CHECKSET,
						"compensate_with_on_demand":      "false",
						"spot_instance_remedy":           "false",
						"soc_enabled":                    "false",
						"system_disk_provisioned_iops":   "3000",
						"system_disk_encrypted":          "false",
						"security_hardening_os":          "false",
						"system_disk_categories.#":       "2",
						"system_disk_bursting_enabled":   "false",
						"user_data":                      CHECKSET,
						"pre_user_data":                  CHECKSET,
						"unschedulable":                  "false",
						"deployment_set_id":              CHECKSET,
						"instance_patterns.#":            "1",
						"cluster_id":                     CHECKSET,
						"system_disk_snapshot_policy_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{
						"ecs.g6.xlarge"},
					"instance_patterns": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#":    "1",
						"instance_patterns.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{},
					"instance_patterns": []map[string]interface{}{
						{
							"cores":                 "4",
							"memory":                "16",
							"instance_family_level": "EnterpriseLevel",
							"excluded_instance_types": []string{
								"ecs.c6.*"},
							"instance_categories": []string{
								"General-purpose"},
							"cpu_architectures": []string{
								"X86"},
							"instance_type_families": []string{
								"ecs.g6"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#":    "0",
						"instance_patterns.#": "1",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eflo_node_group", "password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap11802 = map[string]string{
	"scaling_group_id": CHECKSET,
	"node_pool_id":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence11802(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "cluster_name" {
  default = "tf-test-cluster-1"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw4_cidr" {
  default = "10.1.3.0/24"
}

variable "rg_name_1" {
  default = "tf-test-resource-group-1"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "rg_name_2" {
  default = "tf-test-resource-group-2"
}

variable "container_cidr" {
  default = "172.16.222.0/24"
}

variable "user_data" {
  default = "I18vYmluL3No"
}

variable "service_cidr" {
  default = "172.16.216.0/24"
}

variable "vsw3_cidr" {
  default = "10.1.2.0/24"
}

variable "kubernetes_version" {
  default = "1.32.1-aliyun.1"
}

variable "user_data_1" {
  default = "IyEvYmluL3NoIGVjaG8gIkhlbGxvIFdvcmxkLiBUaGUgdGltZSBpcyBudWcgJChkYXRlIC1SKSkhfCB0ZWUgL3Jvb3QvdXNlcmRhdGFfdGVzdC50eHQ="
}

variable "zone_3" {
  default = "cn-hangzhou-i"
}

variable "zone_4" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_auto_snapshot_policy" "default4L8w9Z" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot"
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultcu8Fq3" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot-2"
}

resource "alicloud_vpc" "defaultqe0KHK" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_ecs_deployment_set" "defaultACvZwq" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds"
}

resource "alicloud_ecs_deployment_set" "default2UoIBA" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds-1"
}

resource "alicloud_vswitch" "defaultVTblQn" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "defaultziRRat" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_vswitch" "defaultT8D8ss" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_3
  cidr_block = var.vsw3_cidr
}

resource "alicloud_vswitch" "defaultFsk7cj" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_4
  cidr_block = var.vsw4_cidr
}

resource "alicloud_cs_managed_kubernetes" "defaultC02XDz" {
  pod_cidr     = var.container_cidr
  vswitch_ids  = ["${alicloud_vswitch.defaultVTblQn.id}"]
  service_cidr = var.service_cidr
  cluster_spec = "ack.pro.small"
  deletion_protection = false
  new_nat_gateway     = false
  is_enterprise_security_group = true
}


`, name)
}

// Case np-kubelet 10193
func TestAccAliCloudAckNodepool_basic10193(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap10193)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence10193)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":          "PostPaid",
					"runtime_name":                  "containerd",
					"image_type":                    "AliyunLinux3",
					"system_disk_encrypt_algorithm": "aes-256",
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "None",
						},
					},
					"system_disk_performance_level": "PL0",
					"node_name_mode":                "customized,aliyun,ip,com",
					"image_id":                      "aliyun_3_x64_20G_alibase_20240819.vhd",
					"install_cloud_monitor":         "false",
					"multi_az_policy":               "PRIORITY",
					"cpu_policy":                    "static",
					"node_pool_name":                name,
					"runtime_version":               "1.6.20",
					"spot_instance_pools":           "1",
					"labels": []map[string]interface{}{
						{
							"key":   "test",
							"value": "test",
						},
					},
					"taints": []map[string]interface{}{
						{
							"key":    "test_taint_key",
							"effect": "NoSchedule",
							"value":  "test_taint_val",
						},
					},
					"internet_charge_type":       "PayByBandwidth",
					"internet_max_bandwidth_out": "1",
					"login_as_non_root":          "false",
					"desired_size":               "0",
					"system_disk_size":           "40",
					"spot_strategy":              "NoSpot",
					"password":                   "Zsyh1234!",
					"scaling_policy":             "release",
					"vswitch_ids": []string{
						"${alicloud_vswitch.defaultziRRat.id}", "${alicloud_vswitch.defaultT8D8ss.id}", "${alicloud_vswitch.defaultFsk7cj.id}"},
					"instance_types": []string{
						"ecs.c7.xlarge"},
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"compensate_with_on_demand":      "false",
					"spot_instance_remedy":           "false",
					"soc_enabled":                    "false",
					"system_disk_provisioned_iops":   "3000",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default4L8w9Z.id}",
					"system_disk_encrypted":          "false",
					"security_hardening_os":          "false",
					"system_disk_categories": []string{
						"cloud_essd", "cloud_auto"},
					"system_disk_bursting_enabled": "false",
					"user_data":                    "${var.user_data}",
					"pre_user_data":                "${var.user_data}",
					"scaling_config": []map[string]interface{}{
						{
							"enable":   "false",
							"type":     "cpu",
							"max_size": "10",
							"min_size": "1",
						},
					},
					"unschedulable":     "false",
					"deployment_set_id": "${alicloud_ecs_deployment_set.defaultACvZwq.id}",
					"management": []map[string]interface{}{
						{
							"enable":          "false",
							"auto_repair":     "false",
							"auto_upgrade":    "false",
							"surge":           "0",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_vul_fix": "false",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "false",
								},
							},
						},
					},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"cpu_manager_policy":    "none",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":    "true",
								"\"CustomCPUCFSQuotaPeriod\"": "true",
							},
							"registry_pull_qps":              "30",
							"registry_burst":                 "10",
							"event_record_qps":               "40",
							"kube_api_burst":                 "20",
							"container_log_monitor_interval": "30s",
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"container_log_max_workers": "5",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4318",
									"sampling_rate_per_million": "1000",
								},
							},
							"topology_manager_policy": "best-effort",
							"pod_pids_limit":          "1024",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "1",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "2",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
							},
							"memory_manager_policy": "None",
							"cluster_dns": []string{
								"8.8.8.8", "10.96.0.10", "192.168.0.100"},
							"image_gc_low_threshold_percent":  "10",
							"image_gc_high_threshold_percent": "80",
							"cpu_cfs_quota_period":            "100ms",
							"cpu_cfs_quota":                   "true",
						},
					},
					"cluster_id": "${alicloud_cs_managed_kubernetes.defaultC02XDz.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":           "PostPaid",
						"runtime_name":                   "containerd",
						"image_type":                     "AliyunLinux3",
						"system_disk_encrypt_algorithm":  "aes-256",
						"system_disk_performance_level":  "PL0",
						"node_name_mode":                 "customized,aliyun,ip,com",
						"image_id":                       "aliyun_3_x64_20G_alibase_20240819.vhd",
						"install_cloud_monitor":          "false",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "static",
						"node_pool_name":                 name,
						"runtime_version":                "1.6.20",
						"spot_instance_pools":            "1",
						"labels.#":                       "1",
						"taints.#":                       "1",
						"internet_charge_type":           "PayByBandwidth",
						"internet_max_bandwidth_out":     "1",
						"login_as_non_root":              "false",
						"desired_size":                   "0",
						"system_disk_size":               "40",
						"spot_strategy":                  "NoSpot",
						"password":                       "Zsyh1234!",
						"scaling_policy":                 "release",
						"vswitch_ids.#":                  "3",
						"instance_types.#":               "1",
						"resource_group_id":              CHECKSET,
						"compensate_with_on_demand":      "false",
						"spot_instance_remedy":           "false",
						"soc_enabled":                    "false",
						"system_disk_provisioned_iops":   "3000",
						"system_disk_snapshot_policy_id": CHECKSET,
						"system_disk_encrypted":          "false",
						"security_hardening_os":          "false",
						"system_disk_categories.#":       "2",
						"system_disk_bursting_enabled":   "false",
						"user_data":                      CHECKSET,
						"pre_user_data":                  CHECKSET,
						"unschedulable":                  "false",
						"deployment_set_id":              CHECKSET,
						"cluster_id":                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":                    "40",
							"kube_api_qps":                   "22",
							"cpu_manager_policy":             "static",
							"container_log_max_files":        "10",
							"registry_pull_qps":              "10",
							"registry_burst":                 "20",
							"event_record_qps":               "10",
							"kube_api_burst":                 "25",
							"container_log_monitor_interval": "10s",
							"container_log_max_workers":      "3",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4317",
									"sampling_rate_per_million": "10000",
								},
							},
							"topology_manager_policy": "restricted",
							"pod_pids_limit":          "800",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "3Gi",
									},
								},
							},
							"memory_manager_policy":           "Static",
							"image_gc_low_threshold_percent":  "15",
							"image_gc_high_threshold_percent": "75",
							"cpu_cfs_quota_period":            "50ms",
							"serialize_image_pulls":           "false",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "20%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "20%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "2Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "2m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "20Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":               "0.5",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"read_only_port":         "10000",
							"max_pods":               "10",
							"container_log_max_size": "2Mi",
							"feature_gates": map[string]interface{}{
								"\"APIListChunking\"":                "true",
								"\"CustomCPUCFSQuotaPeriod\"":        "true",
								"\"MemoryManager\"":                  "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"kernel.msgmax", "net.ipv4.route.min_pmtu", "kernel.shmmax"},
							"cluster_dns": []string{
								"8.8.8.8", "4.4.4.4"},
							"cpu_cfs_quota": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
							"cluster_dns":            []string{},
							"feature_gates": map[string]interface{}{
								"\"APIListChunking\"":                "true",
								"\"CustomCPUCFSQuotaPeriod\"":        "true",
								"\"MemoryManager\"":                  "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eflo_node_group", "password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap10193 = map[string]string{
	"scaling_group_id": CHECKSET,
	"node_pool_id":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence10193(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "cluster_name" {
  default = "tf-test-cluster-1"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw4_cidr" {
  default = "10.1.3.0/24"
}

variable "rg_name_1" {
  default = "tf-test-resource-group-1"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "rg_name_2" {
  default = "tf-test-resource-group-2"
}

variable "container_cidr" {
  default = "172.16.220.0/24"
}

variable "user_data" {
  default = "I18vYmluL3No"
}

variable "service_cidr" {
  default = "172.16.212.0/24"
}

variable "vsw3_cidr" {
  default = "10.1.2.0/24"
}

variable "kubernetes_version" {
  default = "1.32.1-aliyun.1"
}

variable "user_data_1" {
  default = "IyEvYmluL3NoIGVjaG8gIkhlbGxvIFdvcmxkLiBUaGUgdGltZSBpcyBudWcgJChkYXRlIC1SKSkhfCB0ZWUgL3Jvb3QvdXNlcmRhdGFfdGVzdC50eHQ="
}

variable "zone_3" {
  default = "cn-hangzhou-i"
}

variable "zone_4" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_auto_snapshot_policy" "default4L8w9Z" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot"
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultcu8Fq3" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot-2"
}

resource "alicloud_vpc" "defaultqe0KHK" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_ecs_deployment_set" "defaultACvZwq" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds"
}

resource "alicloud_ecs_deployment_set" "default2UoIBA" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds-1"
}

resource "alicloud_vswitch" "defaultVTblQn" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "defaultziRRat" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_vswitch" "defaultT8D8ss" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_3
  cidr_block = var.vsw3_cidr
}

resource "alicloud_vswitch" "defaultFsk7cj" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_4
  cidr_block = var.vsw4_cidr
}

resource "alicloud_cs_managed_kubernetes" "defaultC02XDz" {
  vswitch_ids     = ["${alicloud_vswitch.defaultVTblQn.id}"]
  service_cidr    = var.service_cidr
  pod_vswitch_ids = []
  cluster_spec    = "ack.pro.small"
  pod_cidr        = var.container_cidr
  new_nat_gateway = false
  deletion_protection  = false
  is_enterprise_security_group = true
}


`, name)
}

// Case np-automode 11829
func TestAccAliCloudAckNodepool_basic11829(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap11829)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence11829)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudAckNodepoolBasicDependence11829(name) + fmt.Sprintf(`
resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name = "%s"
  cluster_id     = alicloud_cs_managed_kubernetes.defaultC02XDz.id
  vswitch_ids    = ["${alicloud_vswitch.defaultziRRat.id}", "${alicloud_vswitch.defaultT8D8ss.id}", "${alicloud_vswitch.defaultFsk7cj.id}"]

  auto_mode {
    enabled = true
  }

  scaling_config {
    max_size = 50
    min_size = 0
  }

  instance_patterns {
    min_cpu_cores           = 4
    max_cpu_cores           = 8
    min_memory_size         = 8
    max_memory_size         = 16
    instance_family_level   = "EnterpriseLevel"
    excluded_instance_types = ["ecs.c6.*"]
    instance_categories     = ["General-purpose"]
    cpu_architectures       = ["X86"]
  }

  data_disks {
    size     = 120
    encrypted = "false"
    category = "cloud_essd"
  }

  labels {
    key   = "test"
    value = "test"
  }

  taints {
    key    = "test_taint_key"
    effect = "NoSchedule"
    value  = "test_taint_val"
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0

  lifecycle {
    ignore_changes = [
      management,
      install_cloud_monitor,
      cpu_policy,
      node_name_mode,
      runtime_name,
      runtime_version,
      unschedulable,
      user_data,
      pre_user_data,
      auto_renew,
      auto_renew_period,
      cis_enabled,
      compensate_with_on_demand,
      deployment_set_id,
      image_id,
      image_type,
      instance_charge_type,
      instance_metadata_options,
      internet_charge_type,
      internet_max_bandwidth_out,
      key_name,
      login_as_non_root,
      password,
      multi_az_policy,
      on_demand_base_capacity,
      on_demand_percentage_above_base_capacity,
      period,
      period_unit,
      platform,
      private_pool_options,
      ram_role_name,
      rds_instances,
      scaling_policy,
      security_group_id,
      security_group_ids,
      security_hardening_os,
      soc_enabled,
      spot_instance_pools,
      spot_instance_remedy,
      spot_price_limit,
      spot_strategy,
      system_disk_category,
      system_disk_categories,
      system_disk_size,
      system_disk_bursting_enabled,
      system_disk_performance_level,
      system_disk_encrypted,
      system_disk_kms_key,
      system_disk_snapshot_policy_id,
      system_disk_encrypt_algorithm,
      system_disk_provisioned_iops,
      tee_config,
    ]
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":      name,
						"labels.#":            "1",
						"taints.#":            "1",
						"vswitch_ids.#":       "3",
						"resource_group_id":   CHECKSET,
						"instance_patterns.#": "1",
						"data_disks.#":        "1",
						"cluster_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_patterns": []map[string]interface{}{
						{
							"instance_family_level": "EnterpriseLevel",
							"min_cpu_cores":         "4",
							"max_cpu_cores":         "8",
							"min_memory_size":       "4",
							"max_memory_size":       "16",
							"instance_categories": []string{
								"General-purpose"},
							"cpu_architectures": []string{
								"X86"},
						},
					},
					"data_disks": []map[string]interface{}{
						{
							"size":     "130",
							"category": "cloud_essd",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_patterns.#": "1",
						"data_disks.#":        "1",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eflo_node_group", "password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap11829 = map[string]string{
	"scaling_group_id": CHECKSET,
	"node_pool_id":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence11829(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "cluster_name" {
  default = "tf-test-cluster-1"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw4_cidr" {
  default = "10.1.3.0/24"
}

variable "rg_name_1" {
  default = "tf-test-resource-group-1"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "rg_name_2" {
  default = "tf-test-resource-group-2"
}

variable "container_cidr" {
  default = "172.16.214.0/24"
}

variable "user_data" {
  default = "I18vYmluL3No"
}

variable "service_cidr" {
  default = "172.16.215.0/24"
}

variable "vsw3_cidr" {
  default = "10.1.2.0/24"
}

variable "kubernetes_version" {
  default = "1.32.1-aliyun.1"
}

variable "user_data_1" {
  default = "IyEvYmluL3NoIGVjaG8gIkhlbGxvIFdvcmxkLiBUaGUgdGltZSBpcyBudWcgJChkYXRlIC1SKSkhfCB0ZWUgL3Jvb3QvdXNlcmRhdGFfdGVzdC50eHQ="
}

variable "zone_3" {
  default = "cn-hangzhou-i"
}

variable "zone_4" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_auto_snapshot_policy" "default4L8w9Z" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot"
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultcu8Fq3" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot-2"
}

resource "alicloud_vpc" "defaultqe0KHK" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_ecs_deployment_set" "defaultACvZwq" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds"
}

resource "alicloud_ecs_deployment_set" "default2UoIBA" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds-1"
}

resource "alicloud_vswitch" "defaultVTblQn" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "defaultziRRat" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_vswitch" "defaultT8D8ss" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_3
  cidr_block = var.vsw3_cidr
}

resource "alicloud_vswitch" "defaultFsk7cj" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_4
  cidr_block = var.vsw4_cidr
}

resource "alicloud_cs_managed_kubernetes" "defaultC02XDz" {
  auto_mode {
    enabled = true
  }
  vswitch_ids     = ["${alicloud_vswitch.defaultVTblQn.id}"]
  service_cidr    = var.service_cidr
  pod_vswitch_ids = ["${alicloud_vswitch.defaultVTblQn.id}"]
  cluster_spec    = "ack.pro.small"
  new_nat_gateway      = false
  deletion_protection  = false
  is_enterprise_security_group = true
  control_plane_log_components = ["apiserver","kcm","ccm"]
}


`, name)
}

// Case np-kubelet_副本1763091494524 11832
func TestAccAliCloudAckNodepool_basic11832(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudAckNodepoolMap11832)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckNodepool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence11832)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":          "PostPaid",
					"runtime_name":                  "containerd",
					"image_type":                    "AliyunLinux3",
					"system_disk_encrypt_algorithm": "aes-256",
					"on_demand_base_capacity":       "0",
					"private_pool_options": []map[string]interface{}{
						{
							"private_pool_options_match_criteria": "None",
						},
					},
					"system_disk_performance_level": "PL0",
					"node_name_mode":                "customized,aliyun,ip,com",
					"image_id":                      "aliyun_3_x64_20G_alibase_20240819.vhd",
					"install_cloud_monitor":         "false",
					"multi_az_policy":               "PRIORITY",
					"cpu_policy":                    "static",
					"node_pool_name":                name,
					"runtime_version":               "1.6.20",
					"spot_instance_pools":           "1",
					"labels": []map[string]interface{}{
						{
							"key":   "test",
							"value": "test",
						},
					},
					"taints": []map[string]interface{}{
						{
							"key":    "test_taint_key",
							"effect": "NoSchedule",
							"value":  "test_taint_val",
						},
					},
					"internet_charge_type":       "PayByBandwidth",
					"internet_max_bandwidth_out": "1",
					"login_as_non_root":          "false",
					"desired_size":               "0",
					"system_disk_size":           "40",
					"spot_strategy":              "NoSpot",
					"password":                   "Zsyh1234!",
					"scaling_policy":             "release",
					"vswitch_ids": []string{
						"${alicloud_vswitch.defaultziRRat.id}", "${alicloud_vswitch.defaultT8D8ss.id}", "${alicloud_vswitch.defaultFsk7cj.id}"},
					"instance_types": []string{
						"ecs.g6.xlarge"},
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"compensate_with_on_demand":      "false",
					"spot_instance_remedy":           "false",
					"soc_enabled":                    "false",
					"system_disk_provisioned_iops":   "3000",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default4L8w9Z.id}",
					"system_disk_encrypted":          "false",
					"security_hardening_os":          "false",
					"system_disk_categories": []string{
						"cloud_essd", "cloud_auto"},
					"system_disk_bursting_enabled": "false",
					"user_data":                    "${var.user_data}",
					"pre_user_data":                "${var.user_data}",
					"scaling_config": []map[string]interface{}{
						{
							"enable":   "false",
							"type":     "cpu",
							"max_size": "10",
							"min_size": "1",
						},
					},
					"unschedulable":     "false",
					"deployment_set_id": "${alicloud_ecs_deployment_set.defaultACvZwq.id}",
					"management": []map[string]interface{}{
						{
							"enable":          "false",
							"auto_repair":     "false",
							"auto_upgrade":    "false",
							"surge":           "0",
							"max_unavailable": "1",
							"auto_repair_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_vul_fix": "false",
							"auto_vul_fix_policy": []map[string]interface{}{
								{
									"restart_node": "false",
								},
							},
							"auto_upgrade_policy": []map[string]interface{}{
								{
									"auto_upgrade_kubelet": "false",
								},
							},
						},
					},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":           "50",
							"kube_api_qps":          "20",
							"serialize_image_pulls": "true",
							"cpu_manager_policy":    "none",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "10%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "10%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "1.5Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "1m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":    "500m",
								"\"memory\"": "1Gi",
							},
							"read_only_port":          "0",
							"max_pods":                "200",
							"container_log_max_size":  "10Mi",
							"container_log_max_files": "15",
							"feature_gates": map[string]interface{}{
								"\"GracefulNodeShutdown\"":    "true",
								"\"CustomCPUCFSQuotaPeriod\"": "true",
							},
							"registry_pull_qps":              "30",
							"registry_burst":                 "10",
							"event_record_qps":               "40",
							"kube_api_burst":                 "20",
							"container_log_monitor_interval": "30s",
							"allowed_unsafe_sysctls": []string{
								"net.ipv4.route.min_pmtu"},
							"container_log_max_workers": "5",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4318",
									"sampling_rate_per_million": "1000",
								},
							},
							"topology_manager_policy": "best-effort",
							"pod_pids_limit":          "1024",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "1",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
								{
									"numa_node": "2",
									"limits": map[string]interface{}{
										"\"memory\"": "1Gi",
									},
								},
							},
							"memory_manager_policy": "None",
							"cluster_dns": []string{
								"8.8.8.8", "10.96.0.10", "192.168.0.100"},
							"image_gc_low_threshold_percent":  "10",
							"image_gc_high_threshold_percent": "80",
							"cpu_cfs_quota_period":            "100ms",
							"cpu_cfs_quota":                   "true",
						},
					},
					"cluster_id": "${alicloud_cs_managed_kubernetes.defaultC02XDz.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":           "PostPaid",
						"runtime_name":                   "containerd",
						"image_type":                     "AliyunLinux3",
						"system_disk_encrypt_algorithm":  "aes-256",
						"on_demand_base_capacity":        "0",
						"system_disk_performance_level":  "PL0",
						"node_name_mode":                 "customized,aliyun,ip,com",
						"image_id":                       "aliyun_3_x64_20G_alibase_20240819.vhd",
						"install_cloud_monitor":          "false",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "static",
						"node_pool_name":                 name,
						"runtime_version":                "1.6.20",
						"spot_instance_pools":            "1",
						"labels.#":                       "1",
						"taints.#":                       "1",
						"internet_charge_type":           "PayByBandwidth",
						"internet_max_bandwidth_out":     "1",
						"login_as_non_root":              "false",
						"desired_size":                   "0",
						"system_disk_size":               "40",
						"spot_strategy":                  "NoSpot",
						"password":                       "Zsyh1234!",
						"scaling_policy":                 "release",
						"vswitch_ids.#":                  "3",
						"instance_types.#":               "1",
						"resource_group_id":              CHECKSET,
						"compensate_with_on_demand":      "false",
						"spot_instance_remedy":           "false",
						"soc_enabled":                    "false",
						"system_disk_provisioned_iops":   "3000",
						"system_disk_snapshot_policy_id": CHECKSET,
						"system_disk_encrypted":          "false",
						"security_hardening_os":          "false",
						"system_disk_categories.#":       "2",
						"system_disk_bursting_enabled":   "false",
						"user_data":                      CHECKSET,
						"pre_user_data":                  CHECKSET,
						"unschedulable":                  "false",
						"deployment_set_id":              CHECKSET,
						"cluster_id":                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":                    "40",
							"kube_api_qps":                   "22",
							"cpu_manager_policy":             "static",
							"container_log_max_files":        "10",
							"registry_pull_qps":              "10",
							"registry_burst":                 "20",
							"event_record_qps":               "10",
							"kube_api_burst":                 "25",
							"container_log_monitor_interval": "10s",
							"container_log_max_workers":      "3",
							"tracing": []map[string]interface{}{
								{
									"endpoint":                  "localhost:4317",
									"sampling_rate_per_million": "10000",
								},
							},
							"topology_manager_policy": "restricted",
							"pod_pids_limit":          "800",
							"reserved_memory": []map[string]interface{}{
								{
									"numa_node": "0",
									"limits": map[string]interface{}{
										"\"memory\"": "3Gi",
									},
								},
							},
							"memory_manager_policy":           "Static",
							"image_gc_low_threshold_percent":  "15",
							"image_gc_high_threshold_percent": "75",
							"cpu_cfs_quota_period":            "50ms",
							"serialize_image_pulls":           "false",
							"eviction_hard": map[string]interface{}{
								"\"memory.available\"":            "1024Mi",
								"\"nodefs.available\"":            "20%",
								"\"nodefs.inodesFree\"":           "1000",
								"\"imagefs.available\"":           "20%",
								"\"imagefs.inodesFree\"":          "1000",
								"\"allocatableMemory.available\"": "2048",
								"\"pid.available\"":               "1000",
							},
							"eviction_soft": map[string]interface{}{
								"\"memory.available\"": "2Gi",
							},
							"eviction_soft_grace_period": map[string]interface{}{
								"\"memory.available\"": "2m30s",
							},
							"system_reserved": map[string]interface{}{
								"\"cpu\"":               "1",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "20Gi",
							},
							"kube_reserved": map[string]interface{}{
								"\"cpu\"":               "0.5",
								"\"memory\"":            "1Gi",
								"\"ephemeral-storage\"": "10Gi",
							},
							"read_only_port":         "10000",
							"max_pods":               "10",
							"container_log_max_size": "2Mi",
							"feature_gates": map[string]interface{}{
								"\"APIListChunking\"":                "true",
								"\"CustomCPUCFSQuotaPeriod\"":        "true",
								"\"MemoryManager\"":                  "true",
								"\"RotateKubeletServerCertificate\"": "true",
							},
							"allowed_unsafe_sysctls": []string{
								"kernel.msgmax", "net.ipv4.route.min_pmtu", "kernel.shmmax"},
							"cluster_dns": []string{
								"8.8.8.8", "4.4.4.4"},
							"cpu_cfs_quota": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
							"cluster_dns":            []string{},
							"feature_gates": map[string]interface{}{
								"\"RotateKubeletServerCertificate\"": "true",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eflo_node_group", "password", "rolling_policy", "update_nodes"},
			},
		},
	})
}

var AlicloudAckNodepoolMap11832 = map[string]string{
	"scaling_group_id": CHECKSET,
	"node_pool_id":     CHECKSET,
}

func AlicloudAckNodepoolBasicDependence11832(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "cluster_name" {
  default = "tf-test-cluster-1"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw4_cidr" {
  default = "10.1.3.0/24"
}

variable "rg_name_1" {
  default = "tf-test-resource-group-1"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "rg_name_2" {
  default = "tf-test-resource-group-2"
}

variable "container_cidr" {
  default = "172.16.211.0/24"
}

variable "user_data" {
  default = "I18vYmluL3No"
}

variable "service_cidr" {
  default = "172.16.212.0/24"
}

variable "vsw3_cidr" {
  default = "10.1.2.0/24"
}

variable "kubernetes_version" {
  default = "1.32.1-aliyun.1"
}

variable "user_data_1" {
  default = "IyEvYmluL3NoIGVjaG8gIkhlbGxvIFdvcmxkLiBUaGUgdGltZSBpcyBudWcgJChkYXRlIC1SKSkhfCB0ZWUgL3Jvb3QvdXNlcmRhdGFfdGVzdC50eHQ="
}

variable "zone_3" {
  default = "cn-hangzhou-i"
}

variable "zone_4" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_auto_snapshot_policy" "default4L8w9Z" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot"
}

resource "alicloud_ecs_auto_snapshot_policy" "defaultcu8Fq3" {
  time_points     = ["1"]
  retention_days  = "1"
  repeat_weekdays = ["1"]
  name            = "tf-test-autosnapshot-2"
}

resource "alicloud_vpc" "defaultqe0KHK" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_ecs_deployment_set" "defaultACvZwq" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds"
}

resource "alicloud_ecs_deployment_set" "default2UoIBA" {
  strategy            = "Availability"
  deployment_set_name = "tf-test-ds-1"
}

resource "alicloud_vswitch" "defaultVTblQn" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "defaultziRRat" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_vswitch" "defaultT8D8ss" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_3
  cidr_block = var.vsw3_cidr
}

resource "alicloud_vswitch" "defaultFsk7cj" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_4
  cidr_block = var.vsw4_cidr
}

resource "alicloud_cs_managed_kubernetes" "defaultC02XDz" {
  vswitch_ids     = ["${alicloud_vswitch.defaultVTblQn.id}"]
  service_cidr    = var.service_cidr
  pod_vswitch_ids = []
  cluster_spec    = "ack.pro.small"
  pod_cidr        = var.container_cidr
  deletion_protection = false
  is_enterprise_security_group = true
}


`, name)
}

// Test Ack Nodepool. <<< Resource test cases, automatically generated.
