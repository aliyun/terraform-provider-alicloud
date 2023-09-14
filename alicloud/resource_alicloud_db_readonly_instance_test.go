package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRdsDBReadonlyInstance_update(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "20",
		"engine_version":        "8.0",
		"engine":                "MySQL",
		"port":                  "3306",
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"vswitch_id":            CHECKSET,
		"connection_string":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.read.instance_classes.20.instance_class}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_ips":          "${alicloud_db_instance.default.security_ips}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "effective_time"},
			},

			// upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.read.instance_classes.20.storage_range.step}",
					"effective_time":   "Immediate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": CHECKSET}),
				),
			},
			// upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.read.instance_classes.21.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_type": CHECKSET}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_ro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_ro",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			//UpgradeDBInstanceKernelVersion
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_db_instance_kernel_version": "false",
					"upgrade_time":                       "Immediate",
					"switch_time":                        "2020-01-15T00:00:00Z",
					"target_minor_version":               "rds_20201031",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "rds_20201031",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_restart": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_restart": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "delayed_insert_timeout",
							"value": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			//
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.read.instance_classes.20.instance_class}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage + 2*data.alicloud_db_instance_classes.read.instance_classes.20.storage_range.step}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":    name,
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.13", "100.69.7.113"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_type": "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_mode": "Cover",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modify_mode": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_attribute": "hidden",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_attribute": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"whitelist_network_type": "MIX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"whitelist_network_type": CHECKSET,
					}),
				),
			},
		},
	})

}

func TestAccAliCloudRdsDBReadonlyInstancePostgreSQL_update(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "20",
		"engine_version":        "13.0",
		"engine":                "PostgreSQL",
		"port":                  CHECKSET,
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"vswitch_id":            CHECKSET,
		"connection_string":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.ro.instance_classes.0.instance_class}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
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
			//upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": CHECKSET}),
				),
			},
			//upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.ro.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_type": CHECKSET}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_ro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_ro",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "1",
					"ca_type":     "aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "1",
						"acl":             CHECKSET,
						"replication_acl": CHECKSET,
						"server_cert":     CHECKSET,
						"server_key":      CHECKSET,
						"ca_type":         "aliyun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_enabled": "1",
					"client_ca_cert":    client_ca_cert,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_enabled": "1",
						"client_ca_cert":    client_ca_cert2,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": client_cert_revocation_list2,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl":             "cert",
					"replication_acl": "cert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl":             "cert",
						"replication_acl": "cert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id":       "${alicloud_db_instance.default.id}",
					"zone_id":                     "${alicloud_db_instance.default.zone_id}",
					"engine_version":              "${alicloud_db_instance.default.engine_version}",
					"instance_type":               "${data.alicloud_db_instance_classes.ro.instance_classes.0.instance_class}",
					"instance_storage":            "${alicloud_db_instance.default.instance_storage + 2*data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "${data.alicloud_vswitches.default.ids.0}",
					"ssl_enabled":                 "1",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
					"deletion_protection":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":               name,
						"instance_storage":            CHECKSET,
						"ssl_enabled":                 "1",
						"ca_type":                     "aliyun",
						"client_ca_enabled":           "1",
						"client_ca_cert":              client_ca_cert2,
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": client_cert_revocation_list2,
						"acl":                         "cert",
						"replication_acl":             "cert",
						"server_cert":                 CHECKSET,
						"server_key":                  CHECKSET,
						"deletion_protection":         "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_attribute": "hidden",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_attribute": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"whitelist_network_type": "MIX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"whitelist_network_type": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBReadonlyInstanceMySQL_updatePayType(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_mysql_%d", rand)
	var DBReadonlyMap = map[string]string{}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigMySQLDependence_ro)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.read.instance_classes.0.instance_class}",
					"instance_storage":      "${data.alicloud_db_instance_classes.read.instance_classes.0.storage_range.min}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type":  "Prepaid",
					"period":                "1",
					"auto_renew":            "true",
					"auto_renew_period":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"engine_version":        CHECKSET,
						"instance_type":         CHECKSET,
						"instance_storage":      CHECKSET,
						"instance_name":         name,
						"vswitch_id":            CHECKSET,
						"instance_charge_type":  CHECKSET,
						"period":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":        "true",
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_essd2",
					"instance_storage":         "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": CHECKSET,
						"instance_storage":         CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBReadonlyInstanceMySQL_downgrade(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_mysql_%d", rand)
	var DBReadonlyMap = map[string]string{}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigMySQLDependence_ro_downgrade)
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
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "mysqlro.n4.medium.1c",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type":  "Prepaid",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"engine_version":        CHECKSET,
						"instance_type":         CHECKSET,
						"instance_storage":      CHECKSET,
						"instance_name":         name,
						"vswitch_id":            CHECKSET,
						"instance_charge_type":  CHECKSET,
						"period":                "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "mysqlro.n2.medium.1c",
					"direction":     "Down",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "mysqlro.n2.medium.1c",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "auto_renew_period", "auto_renew", "direction"},
			},
		},
	})
}

func TestAccAliCloudRdsDBReadonlyInstancePostgreSQL_updateDBInstanceSSL(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "20",
		"engine_version":        "13.0",
		"engine":                "PostgreSQL",
		"port":                  CHECKSET,
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"vswitch_id":            CHECKSET,
		"connection_string":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.ro.instance_classes.0.instance_class}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"server_cert":           "-----BEGIN CERTIFICATE-----\\nMIIDdjCCAl4CCQCcm+erkcKN7DANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQGEwJj\\nbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmFsaXl1\\nbjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0LmNvbTEaMBgGCSqGSIb3\\nDQEJARYLMTIzQDEyMy5jb20wHhcNMTkwNDI2MDM0ODAxWhcNMjQwNDI1MDM0ODAx\\nWjB9MQswCQYDVQQGEwJjbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcx\\nDzANBgNVBAoMBmFsaXl1bjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0\\nLmNvbTEaMBgGCSqGSIb3DQEJARYLMTIzQDEyMy5jb20wggEiMA0GCSqGSIb3DQEB\\nAQUAA4IBDwAwggEKAoIBAQDKMKF5qmN/uoMjdH3D8aPRcUOA0s8rZpYhG8zbkF1j\\n8gHYoB/FDvM7G7dfVsyjbMwLOxKvAhWvHHSpEz/t7gB+QdwrAMiMJwGmtCnXrh2E\\nWiXgalMe1y4a/T5R7q+m4T1zFATf+kbnHWfkSGF4W7b6UBoaH+9StQ95CnqzNf/2\\np/Of7+S0XzCxFXw8GIVzZk0xFe6lHJzaq06f3mvzrD+4rpO56tTUvrgTY/n61gsF\\nZP7f0CJ2JQh6eNRFOEUSfxKu/Dy/+IsQxorCJY2Q59ZAf3rXrqDN104jw9PlwnLl\\nqfZz3RMODN6BWjxE8rvRtT0qMfuAfv1gjBdWZN0hUYBRAgMBAAEwDQYJKoZIhvcN\\nAQELBQADggEBAABzo82TxGp5poVkd5pLWj5ACgcBv8Cs6oH9D+4Jz9BmyuBUsQXh\\n2aG0hQAe1mU61C9konsl/GTW8umJQ4M4lYEztXXwMf5PlBMGwebM0ZbSGg6jKtZg\\nWCgJ3eP/FMmyXGL5Jji5+e09eObhUDVle4tdi0On97zBoz85W02rgWFAqZJwiEAP\\nt+c7jX7uOSBq2/38iGStlrX5yB1at/gJXXiA5CL5OtlR3Okvb0/QH37efO1Nu39m\\nlFi0ODPAVyXjVypAiLguDxPn6AtDTdk9Iw9B19OD4NrzNRWgSSX5vuxo/VcRcgWk\\n3gEe9Ca0ZKN20q9XgthAiFFjl1S9ZgdA6Zc=\\n-----END CERTIFICATE-----",
					"server_key":            "-----BEGIN RSA PRIVATE KEY-----\\nMIIEowIBAAKCAQEAyjCheapjf7qDI3R9w/Gj0XFDgNLPK2aWIRvM25BdY/IB2KAf\\nxQ7zOxu3X1bMo2zMCzsSrwIVrxx0qRM/7e4AfkHcKwDIjCcBprQp164dhFol4GpT\\nHtcuGv0+Ue6vpuE9cxQE3/pG5x1n5EhheFu2+lAaGh/vUrUPeQp6szX/9qfzn+/k\\ntF8wsRV8PBiFc2ZNMRXupRyc2qtOn95r86w/uK6TuerU1L64E2P5+tYLBWT+39Ai\\ndiUIenjURThFEn8Srvw8v/iLEMaKwiWNkOfWQH96166gzddOI8PT5cJy5an2c90T\\nDgzegVo8RPK70bU9KjH7gH79YIwXVmTdIVGAUQIDAQABAoIBAE1J4a/8biR5S3/W\\nG+03BYQeY8tuyjqw8FqfoeOcf9agwAvqybouSNQjeCk9qOQfxq/UWQQFK/zQR9gJ\\nv7pX7GBXFK5rkj3g+0SaQhRsPmRFgY0Tl8qGPt2aSKRRNVv5ZeADmwlzRn86QmiF\\nMp0rkfqFfDTYWEepZszCML0ouzuxsW/9tq7rvtSjsgATNt31B3vFa3D3JBi31jUh\\n5nfR9A3bATze7mQw3byEDiVl5ASRDgYyur403P1fDnMy9DBHZ8NaPOsFF6OBpJal\\nBJsG5z00hll5PFN2jfmBQKlvAeU7wfwqdaSnGHOfqf2DeTTaFjIQ4gUhRn/m6pLo\\n6kXttLECgYEA9sng0Qz/TcPFfM4tQ1gyvB1cKnnGIwg1FP8sfUjbbEgjaHhA224S\\nk3BxtX2Kq6fhTXuwusAFc6OVMAZ76FgrQ5K4Ci7+DTsrF28z4b8td+p+lO/DxgP9\\nlTgN+ddsiTOV4fUef9Z3yY0Zr0CnBUMbQYRaV2UIbCdiB0G4V/bt9TsCgYEA0bya\\nOo9wGI0RJV0bYP7qwO74Ra1/i1viWbRlS7jU37Q+AZstrlKcQ5CTPzOjKFKMiUzl\\n4miWacZ0/q2n+Mvd7NbXGXTLijahnyOYKaHJYyh4oBymfkgAifRstE0Ki9gdvArb\\n/I+emC0GvLSyfGN8UUeDJs4NmqdEXGqjo2JOV+MCgYALFv1MR5o9Y1u/hQBRs2fs\\nPiGDIx+9OUQxYloccyaxEfjNXAIGGkcpavchIbgWiJ++PJ2vdquIC8TLeK8evL+M\\n9M3iX0Q5UfxYvD2HmnCvn9D6Xl/cyRcfGnq+TGjrLW9BzSMGuZt+aiHKV0xqFx7l\\nbc4leTvMqGRmURS4lzcQOwKBgQCDzA/i4sYfN25h21tcHXSpnsG3D2rJyQi5NCo/\\nZjunA92/JqOTGuiFcLGHEszhhtY3ZXJET1LNz18vtzKJnpqrvOnYXlOVW/U+SqDQ\\n8JDb1c/PVZGuY1KrXkR9HLiW3kz5IJ3S3PFdUVYdeTN8BQxXCyg4V12nJJtJs912\\ny0zN3wKBgGDS6YttCN6aI4EOABYE8fI1EYQ7vhfiYsaWGWSR1l6bQey7KR6M1ACz\\nZzMASNyytVt12yXE4/Emv6/pYqigbDLfL1zQJSLJ3EHJYTh2RxjR+AaGDudYFG/T\\nliQ9YXhV5Iu2x1pNwrtFnssDdaaGpfA7l3xC00BL7Z+SAJyI4QKA\\n-----END RSA PRIVATE KEY-----",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "ssl_enabled", "server_cert", "server_key"},
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
			//upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": CHECKSET}),
				),
			},
			//upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.ro.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_type": CHECKSET}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_ro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_ro",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "0",
					"ca_type":     "custom",
					"server_cert": "",
					"server_key":  "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_enabled": "1",
					"client_ca_cert":    client_ca_cert,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl":             "cert",
					"replication_acl": "cert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id":       "${alicloud_db_instance.default.id}",
					"zone_id":                     "${alicloud_db_instance.default.zone_id}",
					"engine_version":              "${alicloud_db_instance.default.engine_version}",
					"instance_type":               "${data.alicloud_db_instance_classes.ro.instance_classes.0.instance_class}",
					"instance_storage":            "${alicloud_db_instance.default.instance_storage + 2*data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "${data.alicloud_vswitches.default.ids.0}",
					"ssl_enabled":                 "1",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
					"deletion_protection":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":               name,
						"instance_storage":            CHECKSET,
						"ssl_enabled":                 "1",
						"ca_type":                     "aliyun",
						"client_ca_enabled":           "1",
						"client_ca_cert":              client_ca_cert2,
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": client_cert_revocation_list2,
						"acl":                         "cert",
						"replication_acl":             "cert",
						"server_cert":                 CHECKSET,
						"server_key":                  CHECKSET,
						"deletion_protection":         "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array_attribute": "hidden",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_attribute": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"whitelist_network_type": "MIX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"whitelist_network_type": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBReadonlyInstancePostgreSQL_updatePayType(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_pgsql_%d", rand)
	var DBReadonlyMap = map[string]string{}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigPostgreSQLDependence_ro)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.read.instance_classes.0.instance_class}",
					"instance_storage":      "${data.alicloud_db_instance_classes.read.instance_classes.0.storage_range.min}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type":  "Prepaid",
					"period":                "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"engine_version":        CHECKSET,
						"instance_type":         CHECKSET,
						"instance_storage":      CHECKSET,
						"instance_name":         name,
						"vswitch_id":            CHECKSET,
						"instance_charge_type":  CHECKSET,
						"period":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":        "true",
					"auto_renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_essd2",
					"instance_storage":         "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": CHECKSET,
						"instance_storage":         CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBReadonlyInstanceSQLServer_updatePayType(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_sqlserver_%d", rand)
	var DBReadonlyMap = map[string]string{}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigSQLServerDependence_ro)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id":    "${alicloud_db_instance.default.id}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"engine_version":           "${alicloud_db_instance.default.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.read.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.read.instance_classes.0.storage_range.min}",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"engine_version":        CHECKSET,
						"instance_type":         CHECKSET,
						"instance_storage":      CHECKSET,
						"instance_name":         name,
						"vswitch_id":            CHECKSET,
						"instance_charge_type":  CHECKSET,
						"period":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":        "true",
					"auto_renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceDBReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.28.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.28.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	security_ips = ["10.168.1.12", "100.69.7.112"]
	target_minor_version =  "rds_20201031"
}

data "alicloud_db_instance_classes" "read" {
    db_instance_id = alicloud_db_instance.default.id 
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
    commodity_code = "rords"
}
`, name)
}

func resourceDBReadonlyInstanceConfigPostgreSQLDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "13.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "13.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "PostgreSQL"
	engine_version = "13.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	security_ips = ["10.168.1.12", "100.69.7.112"]
}

data "alicloud_db_instance_classes" "ro" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "13.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
    commodity_code = "rords"
	db_instance_id = alicloud_db_instance.default.id
}
`, name)
}

func resourceDBReadonlyInstanceConfigMySQLDependence_ro(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PrePaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}

data "alicloud_db_instance_classes" "read" {
    db_instance_id = alicloud_db_instance.default.id 
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
    commodity_code = "rds_rordspre_public_cn"
}
`, name)
}

func resourceDBReadonlyInstanceConfigMySQLDependence_ro_downgrade(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PrePaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}
`, name)
}

func resourceDBReadonlyInstanceConfigPostgreSQLDependence_ro(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "10.0"
	instance_charge_type = "PrePaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "10.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
    engine = "PostgreSQL"
	engine_version = "10.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}

data "alicloud_db_instance_classes" "read" {
    db_instance_id = alicloud_db_instance.default.id 
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "10.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
    commodity_code = "rds_rordspre_public_cn"
}
`, name)
}

func resourceDBReadonlyInstanceConfigSQLServerDependence_ro(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "SQLServer"
	engine_version = "2019_ent"
	instance_charge_type = "PostPaid"
	category = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "SQLServer"
	engine_version = "2019_ent"
    category = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
    engine = "SQLServer"
	engine_version = "2019_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	category = "AlwaysOn"
}

data "alicloud_db_instance_classes" "read" {
    db_instance_id = alicloud_db_instance.default.id 
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "SQLServer"
	engine_version = "2019_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
    commodity_code = "rds_rordspre_public_cn"
}
`, name)
}
