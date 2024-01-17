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
					"desired_size":          "1",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks":            []map[string]string{{"size": "100", "category": "cloud_ssd"}},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
					"security_group_ids":    []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"image_type":            "AliyunLinux",
					"cis_enabled":           "true",
					"cpu_policy":            "none",
					"spot_strategy":         "NoSpot",
					"rds_instances":         []string{"pgm-gw88d7qc6b2z6ksw", "pgm-gw8i4i431o53i7sp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"desired_size":                 "1",
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
						"security_group_ids.#":         "2",
						"image_type":                   "AliyunLinux",
						"cis_enabled":                  "true",
						"cpu_policy":                   "none",
						"spot_strategy":                "NoSpot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "cpu_policy", "kubelet_configuration", "rds_instances"},
			},
			// check: scale out
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"desired_size":     "2",
			//		"system_disk_size": "80",
			//		"data_disks":       []map[string]string{{"size": "40", "category": "cloud"}},
			//		"management":       []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "max_unavailable": "1"}},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"desired_size":                 "2",
			//			"system_disk_size":             "80",
			//			"data_disks.#":                 "1",
			//			"data_disks.0.size":            "40",
			//			"data_disks.0.category":        "cloud",
			//			"management.#":                 "1",
			//			"management.0.auto_repair":     "true",
			//			"management.0.auto_upgrade":    "true",
			//			"management.0.surge":           "1",
			//			"management.0.max_unavailable": "1",
			//		}),
			//	),
			//},
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
					"security_group_ids":    []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"image_type":            "CentOS",
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
						"security_group_ids.#":         "2",
						"image_type":                   "CentOS",
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

func TestAccAliCloudCSKubernetesNodePoo_attach(t *testing.T) {
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
					"instances":            []string{"${local.vswitch_id}"},
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_efficiency",
					"system_disk_size":     "40",
					"tags":                 map[string]interface{}{"Created": "TF", "Foo": "Bar"},

					"image_type":    "CentOS",
					"cpu_policy":    "none",
					"spot_strategy": "NoSpot",
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
						"security_group_ids.#":         "2",
						"image_type":                   "CentOS",
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
					"image_type": "CentOS",
					"platform":   "CentOS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_type": "CentOS",
						"platform":   "CentOS",
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
					"scaling_config": []map[string]string{{"enable": "true", "min_size": "2", "max_size": "20", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByTraffic", "eip_bandwidth": "100"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.enable":                   "true",
						"scaling_config.0.min_size":                 "2",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
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
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
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
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
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
			//testAccPreCheckWithRegions(t, true, connectivity.ACKSystemDiskEncryptionSupportRegions)
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
					"instance_types":       []string{"ecs.c7.xlarge"},
					"desired_size":         "1",
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "40",
					"data_disks": []map[string]string{
						{
							"kms_key_id": "${alicloud_kms_key.default0.id}",
							"encrypted":  "true",
							"size":       "100",
							"category":   "cloud_essd",
						},
					},
					"tags":                          map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"security_group_ids":            []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"image_type":                    "CentOS",
					"system_disk_encrypted":         "true",
					"system_disk_kms_key":           "${alicloud_kms_key.default0.id}",
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
						"security_group_ids.#":          "2",
						"image_type":                    "CentOS",
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
					"system_disk_kms_key": "${alicloud_kms_key.default1.id}",
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
			//testAccPreCheckWithRegions(t, true, connectivity.ACKSystemDiskEncryptionSupportRegions)
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
					"instance_types":       []string{"ecs.c7.xlarge"},
					"desired_size":         "1",
					"key_name":             "${alicloud_key_pair.default.key_name}",
					"system_disk_category": "cloud_essd",
					"system_disk_size":     "40",
					"data_disks": []map[string]string{
						//{"kms_key_id": "${data.alicloud_kms_keys.default.ids.0}", "encrypted": "true", "size": "100", "category": "cloud_essd"},
						{
							"size":        "100",
							"category":    "cloud_essd",
							"device":      "/dev/xvdb",
							"snapshot_id": "${alicloud_ecs_snapshot.default.0.id}",
						},
					},
					"tags":               map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"security_group_ids": []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"image_type":         "CentOS",
					//"system_disk_encrypted":         "true",
					//"system_disk_kms_key":           "${data.alicloud_kms_keys.default.ids.0}",
					//"system_disk_encrypt_algorithm": "aes-256",
					"cpu_policy":    "none",
					"spot_strategy": "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"cluster_id":            CHECKSET,
						"vswitch_ids.#":         "1",
						"instance_types.#":      "1",
						"desired_size":          "1",
						"key_name":              CHECKSET,
						"system_disk_category":  "cloud_essd",
						"system_disk_size":      "40",
						"data_disks.#":          "1",
						"data_disks.0.size":     "100",
						"data_disks.0.category": "cloud_essd",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.Foo":              "Bar",
						"security_group_ids.#":  "2",
						"image_type":            "CentOS",
						//"system_disk_encrypted":         "true",
						//"system_disk_kms_key":           CHECKSET,
						//"system_disk_encrypt_algorithm": "aes-256",
						"cpu_policy":    "none",
						"spot_strategy": "NoSpot",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					//"system_disk_kms_key": "${data.alicloud_kms_keys.default.ids.1}",
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
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"system_disk_encrypted":         "false",
			//		"system_disk_encrypt_algorithm": REMOVEKEY,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"system_disk_encrypted":         "false",
			//			"system_disk_encrypt_algorithm": REMOVEKEY,
			//		}),
			//	),
			//},
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
					"data_disks":            []map[string]string{{"size": "100", "category": "cloud_ssd"}},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
					"image_type":            "AliyunLinux",
					"deployment_set_id":     "${alicloud_ecs_deployment_set.default.id}",
					"cpu_policy":            "none",
					"spot_strategy":         "NoSpot",
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
						"image_type":                   "AliyunLinux",
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
				ImportStateVerifyIgnore: []string{"password", "cpu_policy"},
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
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
    instance_type_family       = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
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

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "group1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}
`, name)
}

// system disk encryption only support region HongKong zones B/C
func resourceCSNodePoolConfigDependence_BYOK(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

resource "alicloud_kms_key" "default0" {
 description = "Default"
 status = "Enabled"
 pending_window_in_days = 7
}

resource "alicloud_kms_key" "default1" {
 description = "Default"
 status = "Enabled"
 pending_window_in_days = 7
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
  instance_type_family       = "ecs.sn1ne"
  system_disk_category       = "cloud_essd"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "group1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

resource "alicloud_key_pair" "default" {
	key_pair_name = "tf-testacc-1"
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default.*"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name_prefix          = "Default"
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}


resource "alicloud_ecs_disk" "default" {
  count = 2
  disk_name = var.name
  zone_id   = data.alicloud_zones.default.zones.0.id
  category  = "cloud_efficiency"
  size      = "20"
}


data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_instance" "default" {
  count = 2
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = "terraform-example"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.group.id]
  vswitch_id        = local.vswitch_id
}

resource "alicloud_ecs_disk_attachment" "default" {
  count = 2
  disk_id     = element(alicloud_ecs_disk.default.*.id, count.index)
  instance_id = element(alicloud_instance.default.*.id, count.index)
}


resource "alicloud_ecs_snapshot" "default" {
  count = 2
  force = "true"
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
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
    instance_type_family       = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
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

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "group1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

//resource "alicloud_ecs_elasticity_assurance" "default0yDzRS" {
//  instance_amount                     = "1"
//  zone_ids                            = ["${data.alicloud_zones.default.zones.0.id}"]
//  period                              = "1"
//  private_pool_options_match_criteria = "Open"
//  assurance_times                     = "Unlimited"
//  period_unit                         = "Month"
//  instance_type                       = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
//}
//
//resource "alicloud_ecs_elasticity_assurance" "default4NZTlr" {
//  instance_amount                     = "1"
//  zone_ids                            = ["${data.alicloud_zones.default.zones.0.id}"]
//  period                              = "1"
//  private_pool_options_match_criteria = "Open"
//  assurance_times                     = "Unlimited"
//  period_unit                         = "Month"
//  instance_type                       = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
//}


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
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
    instance_type_family       = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
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
  name_regex = "Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  cluster_id =  length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}
`, name)
}

// Test Ack Nodepool. >>> Resource test cases, automatically generated.
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5288)
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
					"cluster_id":     "${var.cluster}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"instance_types": []string{
						"ecs.c8a.xlarge", "ecs.hfc7.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":   name,
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"instance_types.#": "2",
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
					"auto_renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "1",
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "1.6.20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": "1.6.20",
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
					"on_demand_base_capacity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "ecs.c8a.xlarge",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "ecs.hfc7.xlarge",
							"price_limit":   "0.96",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_price_limit.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "SpotWithPriceLimit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy": "SpotWithPriceLimit",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compensate_with_on_demand": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compensate_with_on_demand": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"on_demand_percentage_above_base_capacity": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_percentage_above_base_capacity": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_instance_pools": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_pools": "2",
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
						"${alicloud_vswitch.vsw1.id}"},
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
						"ecs.c8a.xlarge", "ecs.hfc7.xlarge"},
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${var.cluster}",
					"instance_charge_type": "PostPaid",
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
						"${alicloud_vswitch.vsw1.id}"},
					"login_as_non_root": "true",
					"system_disk_size":  "120",
					"multi_az_policy":   "COST_OPTIMIZED",
					"cpu_policy":        "none",
					"period":            "1",
					"runtime_version":   "1.6.20",
					"desired_size":      "0",
					"instance_types": []string{
						"ecs.c8a.xlarge", "ecs.hfc7.xlarge"},
					"spot_instance_remedy":    "true",
					"on_demand_base_capacity": "1",
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "ecs.c8a.xlarge",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "ecs.hfc7.xlarge",
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
						"auto_renew_period":                        "1",
						"system_disk_category":                     "cloud_essd",
						"vswitch_ids.#":                            "1",
						"login_as_non_root":                        "true",
						"system_disk_size":                         "120",
						"multi_az_policy":                          "COST_OPTIMIZED",
						"cpu_policy":                               "none",
						"period":                                   "1",
						"runtime_version":                          "1.6.20",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5291)
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
					"cluster_id":     "${var.cluster}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"instance_types": []string{
						"ecs.g5.xlarge", "ecs.c5.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":   name,
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"instance_types.#": "2",
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
					"runtime_version": "1.6.20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": "1.6.20",
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
								"\"GracefulNodeShutdown\"": "true",
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
							"private_pool_options_id":             "依赖资源引用：default0yDzRS.ElasticityAssuranceId",
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
					"unschedulable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unschedulable": "true",
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
						"${alicloud_vswitch.vsw1.id}"},
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
						"ecs.g5.xlarge", "ecs.c5.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":  "40",
							"kube_api_qps": "22",
							"eviction_hard": map[string]interface{}{
								"\"nodefs.available\"":  "20%",
								"\"imagefs.available\"": "20%",
							},
							"system_reserved": map[string]interface{}{
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
								"\"cpu\"": "0.5",
							},
							"read_only_port":          "10000",
							"max_pods":                "10",
							"container_log_max_size":  "2Mi",
							"container_log_max_files": "10",
							"feature_gates":           map[string]interface{}{},
							"allowed_unsafe_sysctls": []string{
								"kernel.msg*"},
							"registry_pull_qps": "10",
							"registry_burst":    "20",
							"event_record_qps":  "10",
							"kube_api_burst":    "25",
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
							"private_pool_options_match_criteria": "None",
							"private_pool_options_id":             "${alicloud_ecs_elasticity_assurance.default4NZTlr.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
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
					"cluster_id":           "${var.cluster}",
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
						"${alicloud_vswitch.vsw1.id}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "1.6.20",
					"desired_size":          "1",
					"instance_types": []string{
						"ecs.g5.xlarge", "ecs.c5.xlarge"},
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
								"\"GracefulNodeShutdown\"": "true",
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
							"private_pool_options_id":             "依赖资源引用：default0yDzRS.ElasticityAssuranceId",
						},
					},
					"system_disk_provisioned_iops": "100",
					"system_disk_bursting_enabled": "true",
					"unschedulable":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":               name + "_update",
						"cluster_id":                   CHECKSET,
						"instance_charge_type":         "PostPaid",
						"system_disk_category":         "cloud_essd",
						"vswitch_ids.#":                "1",
						"install_cloud_monitor":        "true",
						"login_as_non_root":            "true",
						"system_disk_size":             "120",
						"multi_az_policy":              "PRIORITY",
						"cpu_policy":                   "none",
						"runtime_version":              "1.6.20",
						"desired_size":                 "1",
						"instance_types.#":             "2",
						"system_disk_provisioned_iops": "100",
						"system_disk_bursting_enabled": "true",
						"unschedulable":                "true",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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

func AlicloudAckNodepoolBasicDependence5291(name string) string {
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

resource "alicloud_ecs_elasticity_assurance" "default0yDzRS" {
  instance_amount                     = "1"
  zone_id                             = ["${alicloud_vswitch.vsw1.zone_id}"]
  period                              = "1"
  private_pool_options_match_criteria = "Open"
  assurance_times                     = "Unlimited"
  period_unit                         = "Month"
}

resource "alicloud_ecs_elasticity_assurance" "default4NZTlr" {
  instance_amount                     = "1"
  zone_id                             = ["${alicloud_vswitch.vsw1.zone_id}"]
  period                              = "1"
  private_pool_options_match_criteria = "Open"
  assurance_times                     = "Unlimited"
  period_unit                         = "Month"
}


`, name)
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5266)
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
					"cluster_id":     "${var.cluster}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"instance_types": []string{
						"ecs.c6.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":   name,
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"instance_types.#": "1",
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
					"auto_renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "1",
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
					"runtime_version": "1.6.20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": "1.6.20",
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
					"period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period_unit": "Month",
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
						"${alicloud_vswitch.vsw1.id}"},
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
						"ecs.c6.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "1",
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
					"node_pool_name":       name + "_update",
					"cluster_id":           "${var.cluster}",
					"instance_charge_type": "PrePaid",
					"auto_renew_period":    "1",
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
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
						"${alicloud_vswitch.vsw1.id}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "BALANCE",
					"cpu_policy":            "none",
					"period":                "1",
					"tee_config": []map[string]interface{}{
						{},
					},
					"system_disk_performance_level": "PL0",
					"runtime_version":               "1.6.20",
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "test",
						},
					},
					"period_unit":  "Month",
					"desired_size": "0",
					"instance_types": []string{
						"ecs.c6.xlarge"},
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
						"runtime_version":               "1.6.20",
						"labels.#":                      "1",
						"period_unit":                   "Month",
						"desired_size":                  "0",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5172)
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
					"cluster_id":     "${var.cluster}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"instance_types": []string{
						"ecs.g5.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":   name,
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"instance_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
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
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime_version": "1.6.20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": "1.6.20",
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
					"key_name": "ack",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": "ack",
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
						"${alicloud_vswitch.vsw1.id}"},
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
						"ecs.g5.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "1",
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
					"vswitch_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"instance_types": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]interface{}{
						{
							"category": "cloud_essd",
						},
						{
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
							"key": "taint1",
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
					"vswitch_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "0",
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
					"user_data": "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_pool_name":       name + "_update",
					"cluster_id":           "${var.cluster}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_charge_type": "PostPaid",
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
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
						"${alicloud_vswitch.vsw1.id}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "80",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"period":                "1",
					"runtime_version":       "1.6.20",
					"desired_size":          "0",
					"instance_types": []string{
						"ecs.g5.xlarge"},
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
					"key_name":                       "ack",
					"deployment_set_id":              "${alicloud_ecs_deployment_set.default.id}",
					"node_name_mode":                 "customized,aliyun,ip,com",
					"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
					"system_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.defaultrt8z7K.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                 name + "_update",
						"cluster_id":                     CHECKSET,
						"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
						"instance_charge_type":           "PostPaid",
						"system_disk_category":           "cloud_essd",
						"vswitch_ids.#":                  "1",
						"install_cloud_monitor":          "true",
						"login_as_non_root":              "true",
						"system_disk_size":               "80",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "none",
						"period":                         "1",
						"runtime_version":                "1.6.20",
						"desired_size":                   "0",
						"instance_types.#":               "1",
						"data_disks.#":                   "3",
						"labels.#":                       "1",
						"taints.#":                       "1",
						"password":                       CHECKSET,
						"scaling_policy":                 "release",
						"key_name":                       "ack",
						"deployment_set_id":              CHECKSET,
						"node_name_mode":                 "customized,aliyun,ip,com",
						"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gInRoaXMgaXMgYSBleGFtcGxlIg==",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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

  retention_days = "-1"
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5401)
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
					"cluster_id": "c0196d207b11d4a25ae4cad2a6f029a38",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"instance_types": []string{
						"ecs.u1-c1m1.xlarge"},
					"node_pool_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":       "c0196d207b11d4a25ae4cad2a6f029a38",
						"vswitch_ids.#":    "1",
						"instance_types.#": "1",
						"node_pool_name":   name,
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
					"runtime_version": "1.6.20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime_version": "1.6.20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
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
						"ecs.u1-c1m1.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "1",
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
					"cluster_id": "c0196d207b11d4a25ae4cad2a6f029a38",
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
						},
					},
					"soc_enabled":          "true",
					"instance_charge_type": "PostPaid",
					"desired_size":         "0",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"system_disk_size": "60",
					"instance_types": []string{
						"ecs.u1-c1m1.xlarge"},
					"password":             "${var.password}",
					"system_disk_category": "cloud_efficiency",
					"runtime_name":         "containerd",
					"runtime_version":      "1.6.20",
					"node_pool_name":       name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":           "c0196d207b11d4a25ae4cad2a6f029a38",
						"soc_enabled":          "true",
						"instance_charge_type": "PostPaid",
						"desired_size":         "0",
						"vswitch_ids.#":        "1",
						"system_disk_size":     "60",
						"instance_types.#":     "1",
						"password":             CHECKSET,
						"system_disk_category": "cloud_efficiency",
						"runtime_name":         "containerd",
						"runtime_version":      "1.6.20",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
func TestAccAliCloudAckNodepool_basic5628(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5628)
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
						"${alicloud_vswitch.vsw1.id}"},
					"cluster_id": "${var.cluster}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":   name,
						"instance_types.#": "3",
						"vswitch_ids.#":    "1",
						"cluster_id":       CHECKSET,
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
					"image_type": "AliyunLinuxSecurity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_type": "AliyunLinuxSecurity",
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
					"instance_types": []string{
						"ecs.c7t.xlarge", "ecs.g7t.xlarge", "ecs.r7t.xlarge"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
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
					"instance_charge_type": "PrePaid",
					"runtime_name":         "containerd",
					"tee_config": []map[string]interface{}{
						{
							"tee_enable": "true",
						},
					},
					"node_pool_name": name + "_update",
					"desired_size":   "1",
					"instance_types": []string{
						"ecs.c7t.xlarge", "ecs.g7t.xlarge", "ecs.r7t.xlarge"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"cluster_id":           "${var.cluster}",
					"system_disk_size":     "40",
					"system_disk_category": "cloud_essd",
					"image_type":           "AliyunLinuxSecurity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"runtime_name":         "containerd",
						"node_pool_name":       name + "_update",
						"desired_size":         "1",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Rds)
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
						"${local.vswitch_id}"},
					"login_as_non_root": "true",
					"system_disk_size":  "120",
					"multi_az_policy":   "COST_OPTIMIZED",
					"cpu_policy":        "none",
					"period":            "1",
					"runtime_version":   "1.6.20",
					"desired_size":      "0",
					"instance_types": []string{
						"ecs.c8a.xlarge", "ecs.hfc7.xlarge"},
					"on_demand_base_capacity": "2",
					"spot_price_limit": []map[string]interface{}{
						{
							"instance_type": "ecs.c8a.xlarge",
							"price_limit":   "0.96",
						},
						{
							"instance_type": "ecs.hfc7.xlarge",
							"price_limit":   "0.96",
						},
					},
					"spot_strategy": "SpotWithPriceLimit",
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
						"node_pool_name":          name,
						"cluster_id":              CHECKSET,
						"instance_charge_type":    "PostPaid",
						"auto_renew_period":       "0",
						"system_disk_category":    "cloud_essd",
						"vswitch_ids.#":           "1",
						"login_as_non_root":       "true",
						"system_disk_size":        "120",
						"multi_az_policy":         "COST_OPTIMIZED",
						"cpu_policy":              "none",
						"period":                  "1",
						"runtime_version":         "1.6.20",
						"desired_size":            "0",
						"instance_types.#":        "2",
						"on_demand_base_capacity": "2",
						"spot_price_limit.#":      "2",
						"spot_strategy":           "SpotWithPriceLimit",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Rds)
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
					"system_disk_category": "cloud_essd",
					"vswitch_ids": []string{
						"${local.vswitch_id}"},
					"install_cloud_monitor": "true",
					"login_as_non_root":     "true",
					"system_disk_size":      "120",
					"multi_az_policy":       "PRIORITY",
					"cpu_policy":            "none",
					"runtime_version":       "1.6.20",
					"desired_size":          "2",
					"instance_types": []string{
						"ecs.g5.xlarge", "ecs.c5.xlarge"},
					"kubelet_configuration": []map[string]interface{}{
						{
							"event_burst":  "40",
							"kube_api_qps": "22",
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
							"private_pool_options_match_criteria": "None",
							"private_pool_options_id":             "${alicloud_ecs_elasticity_assurance.default4NZTlr.id}",
						},
					},
					"system_disk_provisioned_iops": "200",
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
						"system_disk_category":         "cloud_essd",
						"vswitch_ids.#":                "1",
						"install_cloud_monitor":        "true",
						"login_as_non_root":            "true",
						"system_disk_size":             "120",
						"multi_az_policy":              "PRIORITY",
						"cpu_policy":                   "none",
						"runtime_version":              "1.6.20",
						"desired_size":                 "2",
						"instance_types.#":             "2",
						"system_disk_provisioned_iops": "200",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence_Rds)
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
					"node_pool_name":       name,
					"cluster_id":           "${local.cluster_id}",
					"instance_charge_type": "PostPaid",
					"auto_renew_period":    "0",
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
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
						"${local.vswitch_id}"},
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
					"runtime_version":               "1.6.20",
					"labels": []map[string]interface{}{
						{
							"key":   "tf",
							"value": "test",
						},
					},
					"period_unit":  "",
					"desired_size": "0",
					"auto_renew":   "false",
					"instance_types": []string{
						"ecs.c6.xlarge"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_pool_name":                name,
						"cluster_id":                    CHECKSET,
						"instance_charge_type":          "PostPaid",
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
						"runtime_version":               "1.6.20",
						"labels.#":                      "1",
						"period_unit":                   "Month",
						"desired_size":                  "0",
						"auto_renew":                    "true",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5172)
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
					"node_pool_name":       name,
					"cluster_id":           "${var.cluster}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_charge_type": "PostPaid",
					"scaling_config": []map[string]interface{}{
						{
							"type": "cpu",
						},
					},
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
						"${alicloud_vswitch.vsw1.id}", "${alicloud_vswitch.vsw2.id}", "${alicloud_vswitch.vsw3.id}"},
					"login_as_non_root": "true",
					"system_disk_size":  "100",
					"multi_az_policy":   "PRIORITY",
					"cpu_policy":        "none",
					"period":            "1",
					"runtime_version":   "1.6.20",
					"desired_size":      "0",
					"instance_types": []string{
						"ecs.g5.xlarge", "ecs.c5.xlarge", "ecs.ic5.xlarge"},
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
					"key_name":                       "ack",
					"deployment_set_id":              "${alicloud_ecs_deployment_set.default.id}",
					"node_name_mode":                 "customized,aliyun,ip,com",
					"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
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
						"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
						"instance_charge_type":           "PostPaid",
						"system_disk_category":           "cloud_essd",
						"vswitch_ids.#":                  "3",
						"login_as_non_root":              "true",
						"system_disk_size":               "100",
						"multi_az_policy":                "PRIORITY",
						"cpu_policy":                     "none",
						"period":                         "1",
						"runtime_version":                "1.6.20",
						"desired_size":                   "0",
						"instance_types.#":               "3",
						"data_disks.#":                   "3",
						"labels.#":                       "3",
						"taints.#":                       "3",
						"password":                       CHECKSET,
						"scaling_policy":                 "release",
						"key_name":                       "ack",
						"deployment_set_id":              CHECKSET,
						"node_name_mode":                 "customized,aliyun,ip,com",
						"user_data":                      "IyEvYmluL2Jhc2gKCmVjaG8gIlRoaXMgaXMgYSBleGFtcGxlIg==",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5401)
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
					"cluster_id": "c0196d207b11d4a25ae4cad2a6f029a38",
					"kubelet_configuration": []map[string]interface{}{
						{
							"allowed_unsafe_sysctls": []string{},
						},
					},
					"soc_enabled":          "true",
					"instance_charge_type": "PostPaid",
					"desired_size":         "0",
					"vswitch_ids": []string{
						"${alicloud_vswitch.vsw1.id}"},
					"system_disk_size": "60",
					"instance_types": []string{
						"ecs.u1-c1m1.xlarge"},
					"password":             "${var.password}",
					"system_disk_category": "cloud_efficiency",
					"runtime_name":         "containerd",
					"runtime_version":      "1.6.20",
					"node_pool_name":       name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":           "c0196d207b11d4a25ae4cad2a6f029a38",
						"soc_enabled":          "true",
						"instance_charge_type": "PostPaid",
						"desired_size":         "0",
						"vswitch_ids.#":        "1",
						"system_disk_size":     "60",
						"instance_types.#":     "1",
						"password":             CHECKSET,
						"system_disk_category": "cloud_efficiency",
						"runtime_name":         "containerd",
						"runtime_version":      "1.6.20",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
			},
		},
	})
}

// Case 节点池测试_teeconfig 5628  twin
func TestAccAliCloudAckNodepool_basic5628_twin(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckNodepoolBasicDependence5628)
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
						"${alicloud_vswitch.vsw1.id}"},
					"cluster_id":           "${var.cluster}",
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
				ImportStateVerifyIgnore: []string{"rolling_policy"},
			},
		},
	})
}

// Test Ack Nodepool. <<< Resource test cases, automatically generated.
