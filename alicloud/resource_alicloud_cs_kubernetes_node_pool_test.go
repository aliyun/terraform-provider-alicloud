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

func TestAccAlicloudCSKubernetesNodePool_basic(t *testing.T) {
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
					"cluster_id":            "${alicloud_cs_managed_kubernetes.default.0.id}",
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
					"runtime_name":          "containerd",
					"runtime_version":       "1.4.8",
					"image_type":            "AliyunLinux",
					"deployment_set_id":     "${alicloud_ecs_deployment_set.default.id}",
					"cis_enabled":           "true",
					"rds_instances":         []string{"${alicloud_db_instance.default.id}"},
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
						"runtime_name":                 "containerd",
						"runtime_version":              "1.4.8",
						"image_type":                   "AliyunLinux",
						"deployment_set_id":            CHECKSET,
						"cis_enabled":                  "true",
						"rds_instances.#":              "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// check: scale out
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_size":     "2",
					"system_disk_size": "80",
					"data_disks":       []map[string]string{{"size": "40", "category": "cloud"}},
					"management":       []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "max_unavailable": "1"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size":                 "2",
						"system_disk_size":             "80",
						"data_disks.#":                 "1",
						"data_disks.0.size":            "40",
						"data_disks.0.category":        "cloud",
						"management.#":                 "1",
						"management.0.auto_repair":     "true",
						"management.0.auto_upgrade":    "true",
						"management.0.surge":           "1",
						"management.0.max_unavailable": "1",
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
		},
	})
}

func TestAccAlicloudCSKubernetesNodePoolWithNodeCount_basic(t *testing.T) {
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
					"cluster_id":            "${alicloud_cs_managed_kubernetes.default.0.id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"node_count":            "1",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks":            []map[string]string{{"size": "100", "category": "cloud_ssd"}},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
					"security_group_ids":    []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"runtime_name":          "containerd",
					"runtime_version":       "1.4.8",
					"image_type":            "CentOS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"node_count":                   "1",
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
						"runtime_name":                 "containerd",
						"runtime_version":              "1.4.8",
						"image_type":                   "CentOS",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// check: scale out
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count":       "2",
					"system_disk_size": "80",
					"data_disks":       []map[string]string{{"size": "40", "category": "cloud"}},
					"management":       []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "max_unavailable": "1"}},
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
						"management.0.max_unavailable": "1",
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

func TestAccAlicloudCSKubernetesNodePool_autoScaling(t *testing.T) {
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
					"name":                  name,
					"cluster_id":            "${alicloud_cs_managed_kubernetes.default.0.id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"platform":              "AliyunLinux",
					"scaling_policy":        "release",
					"scaling_config":        []map[string]string{{"min_size": "1", "max_size": "10", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"key_name":                     CHECKSET,
						"system_disk_category":         "cloud_efficiency",
						"system_disk_size":             "40",
						"install_cloud_monitor":        "false",
						"platform":                     "AliyunLinux",
						"scaling_policy":               "release",
						"scaling_config.#":             "1",
						"scaling_config.0.min_size":    "1",
						"scaling_config.0.max_size":    "10",
						"scaling_config.0.type":        "cpu",
						"scaling_config.0.is_bond_eip": "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "node_count", "desired_size"},
			},
			// check: update config
			{
				Config: testAccConfig(map[string]interface{}{
					"platform":       "AliyunLinux",
					"scaling_policy": "release",
					"scaling_config": []map[string]string{{"min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform":                                  "AliyunLinux",
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]string{{"min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "false", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_config.#":                          "1",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "false",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCSKubernetesNodePool_PrePaid(t *testing.T) {
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
					"cluster_id":            "${alicloud_cs_managed_kubernetes.default.0.id}",
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

func TestAccAlicloudCSKubernetesNodePool_Spot(t *testing.T) {
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
					"cluster_id":                 "${alicloud_cs_managed_kubernetes.default.0.id}",
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

func TestAccAlicloudCSKubernetesNodePool_BYOK(t *testing.T) {
	// modify to supported region
	prevRegion := os.Getenv("ALICLOUD_REGION")
	os.Setenv("ALICLOUD_REGION", "cn-hongkong")
	defer os.Setenv("ALICLOUD_REGION", prevRegion)

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
			testAccPreCheckWithRegions(t, true, connectivity.ACKSystemDiskEncryptionSupportRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"cluster_id":                    "${alicloud_cs_managed_kubernetes.default.0.id}",
					"vswitch_ids":                   []string{"${local.vswitch_id}"},
					"instance_types":                []string{"ecs.c7.xlarge"},
					"desired_size":                  "1",
					"key_name":                      "${alicloud_key_pair.default.key_name}",
					"system_disk_category":          "cloud_essd",
					"system_disk_size":              "40",
					"data_disks":                    []map[string]string{{"size": "100", "category": "cloud_essd"}},
					"tags":                          map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"security_group_ids":            []string{"${alicloud_security_group.group.id}", "${alicloud_security_group.group1.id}"},
					"image_type":                    "CentOS",
					"system_disk_encrypted":         "true",
					"system_disk_kms_key":           "${data.alicloud_kms_keys.default.ids.0}",
					"system_disk_encrypt_algorithm": "aes-256",
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

func TestAccAlicloudCSKubernetesNodePool_AliyunLinux3(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

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
					"name":                  name,
					"cluster_id":            "${alicloud_cs_managed_kubernetes.default.0.id}",
					"vswitch_ids":           []string{"${local.vswitch_id}"},
					"instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"desired_size":          "1",
					"key_name":              "${alicloud_key_pair.default.key_name}",
					"system_disk_category":  "cloud_efficiency",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks":            []map[string]string{{"size": "100", "category": "cloud_ssd"}},
					"image_type":            "AliyunLinux3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"instance_types.#": "1",
						"desired_size":     "1",
						"key_name":         CHECKSET,
						"image_type":       "AliyunLinux3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}"},
					"desired_size":   "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_size":     "2",
						"instance_types.#": "2",
						"image_type":       "AliyunLinux3",
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
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "group" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "group1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}



data "alicloud_db_instances" "default" {
	engine = "MySQL"
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

resource "alicloud_key_pair" "default" {
	key_name = var.name
}

resource "alicloud_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = data.alicloud_db_instances.default.instances.0.engine
  engine_version   = data.alicloud_db_instances.default.instances.0.engine_version
  instance_type    = data.alicloud_db_instances.default.instances.0.instance_type
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [local.vswitch_id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
  
  maintenance_window {
    enable            = true
    maintenance_time  = "03:00:00Z"
    duration          = "3h"
    weekly_period     = "Thursday"
  }
}
`, name)
}

// system disk encryption only support region HongKong zones B/C
func resourceCSNodePoolConfigDependence_BYOK(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "alicloud_zone" {
    default = "cn-hongkong-c"
}

data "alicloud_kms_keys" "default" {
    status = "Enabled"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = var.alicloud_zone
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = var.alicloud_zone
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
  zone_id           = var.alicloud_zone
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_key_pair" "default" {
	key_pair_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [local.vswitch_id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
  
  maintenance_window {
    enable            = true
    maintenance_time  = "03:00:00Z"
    duration          = "3h"
    weekly_period     = "Thursday"
  }
}

`, name)
}
