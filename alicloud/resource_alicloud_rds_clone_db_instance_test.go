package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_rds_clone_db_instance", &resource.Sweeper{
		Name: "alicloud_rds_clone_db_instance",
		F:    testSweepDBInstances,
	})
}

func resourceCloneDBInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = "pg.n2.2c.2m"
  instance_storage         = "30"
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}
`, name)
}

var cloneInstanceBasicMap = map[string]string{}
var cloneInstanceClusterMap = map[string]string{
	"zone_id_slave_a": CHECKSET,
	"zone_id_slave_b": CHECKSET,
}

func TestAccAliCloudRdsCloneDBInstancePostgreSQLSSL(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccpgsslclone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_db_instance_id":    CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"payment_type":             "PayAsYouGo",
						"backup_id":                CHECKSET,
						"engine_version":           "13.0",
						"db_instance_class":        CHECKSET,
						"db_instance_storage":      CHECKSET,
						"zone_id":                  CHECKSET,
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
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
					"db_instance_description": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tcp_connection_type": "SHORT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_connection_type": "SHORT",
					}),
				),
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
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "pg.n2.4c.2m",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": CHECKSET,
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
					"port":                     "3333",
					"connection_string_prefix": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "1",
						"ca_type":         "aliyun",
						"acl":             "prefer",
						"replication_acl": "prefer",
						"server_cert":     CHECKSET,
						"server_key":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":                 "1",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description":     "tf-testAccDBInstance_instance_name",
					"security_ips":                []string{"10.168.1.12", "100.69.7.112"},
					"port":                        "3333",
					"connection_string_prefix":    "${var.name}",
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
					"backup_type":              "FullBackup",
					"certificate":              "-----BEGIN CERTIFICATE-----\ntfacc\n-----END CERTIFICATE-----",
					"db_name":                  "tfaccdb",
					"db_names":                 "[\"tfaccdb\"]",
					"dedicated_host_group_id":  "dhg-tfacc-placeholder",
					"direction":                "Up",
					"effective_time":           "2026-01-01T12:05:05Z",
					"encryption_key":           "key-tfacc-placeholder",
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"force_restart":            "true",
					"ha_mode":                  "RPO",
					"password":                 "YourPassword_123",
					"period":                   "Year",
					"private_ip_address":       "172.16.0.10",
					"private_key":              "-----BEGIN PRIVATE KEY-----\ntfacc\n-----END PRIVATE KEY-----",
					"released_keep_policy":     "Delete",
					"resource_group_id":        "rg-tfacc-placeholder",
					"restore_table":            "false",
					"restore_time":             "2026-01-01T11:05:05Z",
					"role_arn":                 "acs:ram::1234567890123456:role/tfacc-placeholder",
					"server_cert":              "-----BEGIN CERTIFICATE-----\ntfacc-server\n-----END CERTIFICATE-----",
					"server_key":               "-----BEGIN PRIVATE KEY-----\ntfacc-server\n-----END PRIVATE KEY-----",
					"source_biz":               "1",
					"switch_time":              "2026-01-01T13:05:05Z",
					"sync_mode":                "Sync",
					"table_meta":               "[{\"name\":\"tfacc\",\"type\":\"int\"}]",
					"tde_status":               "Enabled",
					"used_time":                "12",
					"db_instance_storage_type": "cloud_essd2",
					"payment_type":             "Subscription",
					"serverless_config": []map[string]interface{}{
						{
							"auto_pause":   "true",
							"switch_force": "true",
						},
					},
					"parameters": []map[string]interface{}{
						{
							"name":  "tfacc_param",
							"value": "1",
						},
					},
					"pg_hba_conf": []map[string]interface{}{
						{
							"type":        "host",
							"database":    "all",
							"user":        "all",
							"address":     "0.0.0.0/0",
							"mask":        "255.255.255.255",
							"method":      "md5",
							"option":      "sameuser",
							"priority_id": "0",
						},
					},
					"source_db_instance_id": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id", "client_ca_enabled", "client_crl_enabled", "connection_string_prefix", "ssl_enabled", "pg_hba_conf"},
			},
		},
	})
}

// SSL function and pg_hba_conf function incompatible, so add this test case for pg_hba_conf without ssl function.
func TestAccAliCloudRdsCloneDBInstancePostgreSQL_PG_HBA_CONF(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccpghbaclone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_db_instance_id":    CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"payment_type":             "PayAsYouGo",
						"backup_id":                CHECKSET,
						"engine_version":           "13.0",
						"db_instance_class":        CHECKSET,
						"db_instance_storage":      CHECKSET,
						"zone_id":                  CHECKSET,
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccDBInstance_instance_name",
					}),
				),
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
			{
				Config: testAccConfig(map[string]interface{}{
					"pg_hba_conf": []interface{}{
						map[string]interface{}{
							"type":        "host",
							"user":        "all",
							"address":     "0.0.0.0/0",
							"database":    "all",
							"method":      "md5",
							"priority_id": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pg_hba_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "pg.n2.4c.2m",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": CHECKSET,
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
					"port":                     "3333",
					"connection_string_prefix": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description":  "tf-testAccDBInstance_instance_name",
					"security_ips":             []string{"10.168.1.12", "100.69.7.112"},
					"port":                     "3333",
					"connection_string_prefix": "${var.name}",
					"deletion_protection":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id", "connection_string_prefix", "pg_hba_conf"},
			},
		},
	})
}

func TestAccAliCloudRdsCloneDBInstanceMySQL_ServerlessBasic(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_CloneMySQLServerlessBasic_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceMySQLServerlessBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_class":        "${alicloud_db_instance.default.instance_type}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"instance_network_type":    "VPC",
					"vpc_id":                   "${alicloud_db_instance.default.vpc_id}",
					"vswitch_id":               "${alicloud_db_instance.default.vswitch_id}",
					"category":                 "serverless_basic",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "Serverless",
					"db_instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "0.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":                CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_network_type":            "VPC",
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"category":                         "serverless_basic",
						"db_instance_storage_type":         "cloud_essd",
						"payment_type":                     "Serverless",
						"db_instance_storage":              CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "0.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "7",
							"min_capacity": "1.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "7",
						"serverless_config.0.min_capacity": "1.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id"},
			},
		},
	})
}

func TestAccAliCloudRdsCloneDBInstanceMySQL_Cluster(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence_MySQLCluster)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"zone_id_slave_a":          "${alicloud_db_instance.default.zone_id_slave_a}",
					"zone_id_slave_b":          "${alicloud_db_instance.default.zone_id_slave_b}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_db_instance_id":    CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"payment_type":             "PayAsYouGo",
						"backup_id":                CHECKSET,
						"engine_version":           CHECKSET,
						"db_instance_class":        CHECKSET,
						"db_instance_storage":      CHECKSET,
						"zone_id":                  CHECKSET,
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_minor_version": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_minor_version": "Auto",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min + data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccDBInstance_instance_name",
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
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
					"tcp_connection_type": "SHORT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_connection_type": "SHORT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "22:00Z-02:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "22:00Z-02:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3333",
					"connection_string_prefix": "rm-ccccccc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": "rm-ccccccc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id", "connection_string_prefix", "pg_hba_conf"},
			},
		},
	})
}

func resourceCloneDBInstanceMySQLServerlessBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "MySQL"
    engine_version = "8.0"
    instance_charge_type = "Serverless"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "MySQL"
    engine_version = "8.0"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_charge_type     = "Serverless"
  instance_name            = var.name
  zone_id                  = data.alicloud_db_zones.default.ids.1
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  db_instance_storage_type = "cloud_essd"
  category                 = "serverless_basic"
  serverless_config {
    max_capacity = 8
    min_capacity = 0.5
    auto_pause   = false
    switch_force = false
  }
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

`, name)
}

func resourceCloneDBInstanceConfigDependence_MySQLCluster(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id 		 = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = "mysql.n2e.small.xc"
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
  zone_id 				   = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.0
  zone_id_slave_b          = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

`, name)
}

func TestAccAliCloudRdsCloneDBInstanceMySQL_ServerlessStandard(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_CloneMySQLServerlessStandard_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceMySQLServerlessStandardConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_class":        "${alicloud_db_instance.default.instance_type}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"zone_id_slave_a":          "${alicloud_db_instance.default.zone_id_slave_a}",
					"instance_network_type":    "VPC",
					"vpc_id":                   "${alicloud_db_instance.default.vpc_id}",
					"vswitch_id":               "${join(\",\", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])}",
					"category":                 "serverless_standard",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "Serverless",
					"db_instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "0.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":                CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_network_type":            "VPC",
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"category":                         "serverless_standard",
						"db_instance_storage_type":         "cloud_essd",
						"payment_type":                     "Serverless",
						"db_instance_storage":              CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "0.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "7",
							"min_capacity": "1.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "7",
						"serverless_config.0.min_capacity": "1.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id"},
			},
		},
	})
}

func resourceCloneDBInstanceMySQLServerlessStandardConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "MySQL"
    engine_version = "8.0"
    instance_charge_type = "Serverless"
    category = "serverless_standard"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "MySQL"
    engine_version = "8.0"
    category = "serverless_standard"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "vswitche1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_vswitches" "vswitche2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_charge_type     = "Serverless"
  instance_name            = var.name
  zone_id                  = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.1
  vswitch_id               = join(",", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])
  db_instance_storage_type = "cloud_essd"
  category                 = "serverless_standard"
  serverless_config {
    max_capacity = 8
    min_capacity = 0.5
    auto_pause   = false
    switch_force = false
  }
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

`, name)
}

func TestAccAliCloudRdsCloneDBInstancePostgreSQL_ServerlessBasic(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_ClonePostgreSQLServerlessBasic_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstancePostgreSQLServerlessBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_class":        "${alicloud_db_instance.default.instance_type}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"instance_network_type":    "VPC",
					"vpc_id":                   "${alicloud_db_instance.default.vpc_id}",
					"vswitch_id":               "${alicloud_db_instance.default.vswitch_id}",
					"category":                 "serverless_basic",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "Serverless",
					"db_instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "12",
							"min_capacity": "0.5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":                CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_network_type":            "VPC",
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"category":                         "serverless_basic",
						"db_instance_storage_type":         "cloud_essd",
						"payment_type":                     "Serverless",
						"db_instance_storage":              CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "12",
						"serverless_config.0.min_capacity": "0.5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "11",
							"min_capacity": "1.5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "11",
						"serverless_config.0.min_capacity": "1.5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id"},
			},
		},
	})
}

func resourceCloneDBInstancePostgreSQLServerlessBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "PostgreSQL"
    engine_version = "14.0"
    instance_charge_type = "Serverless"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "PostgreSQL"
    engine_version = "14.0"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "14.0"
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_charge_type     = "Serverless"
  instance_name            = var.name
  zone_id                  = data.alicloud_db_zones.default.ids.1
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  db_instance_storage_type = "cloud_essd"
  category                 = "serverless_basic"
  serverless_config {
    max_capacity = 12
    min_capacity = 0.5
  }
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

`, name)
}

func TestAccAliCloudRdsCloneDBInstancegeneral_essd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccpgsslclone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence_general_essd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_storage_type": "general_essd",
					"payment_type":             "PayAsYouGo",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_db_instance_id":    CHECKSET,
						"db_instance_storage_type": "general_essd",
						"payment_type":             "PayAsYouGo",
						"backup_id":                CHECKSET,
						"engine_version":           "8.0",
						"db_instance_class":        CHECKSET,
						"db_instance_storage":      CHECKSET,
						"zone_id":                  CHECKSET,
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id", "client_ca_enabled", "client_crl_enabled", "connection_string_prefix", "ssl_enabled", "pg_hba_conf"},
			},
		},
	})
}
func resourceCloneDBInstanceConfigDependence_general_essd(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id 		 = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
  zone_id 				   = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

`, name)
}

// TestAccAliCloudRdsCloneDBInstanceMySQL_Tags verifies that tags configured at
// creation time are carried atomically by CloneDBInstance (AC1) and that the
// subsequent update/remove of tags still works through setInstanceTags (AC2).
func TestAccAliCloudRdsCloneDBInstanceMySQL_Tags(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccmysqltagsclone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence_general_essd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				// AC1: tags set in the create config are carried atomically by the cloned instance on creation (via CloneDBInstance).
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_storage_type": "general_essd",
					"payment_type":             "PayAsYouGo",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_db_instance_id":    CHECKSET,
						"db_instance_storage_type": "general_essd",
						"payment_type":             "PayAsYouGo",
						"backup_id":                CHECKSET,
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "acceptance Test",
					}),
				),
			},
			{
				// AC2: removing all tags exercises the existing setInstanceTags (UnTagResources) path on the cloned instance. REMOVEKEY clears every key Step 1 set, so the accumulating check helper does not retain a stale expectation for a removed tag.
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
		},
	})
}
func TestAccAliCloudRdsCloneDBInstanceSQLServer_ServerlessHA(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_CloneMssqlServerlessHA_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceMssqlServerlessHAConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_db_instance_id":    "${alicloud_db_instance.default.id}",
					"db_instance_class":        "${alicloud_db_instance.default.instance_type}",
					"zone_id":                  "${alicloud_db_instance.default.zone_id}",
					"zone_id_slave_a":          "${alicloud_db_instance.default.zone_id_slave_a}",
					"instance_network_type":    "VPC",
					"vpc_id":                   "${alicloud_db_instance.default.vpc_id}",
					"vswitch_id":               "${join(\",\", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])}",
					"category":                 "serverless_ha",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "Serverless",
					"db_instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"backup_id":                "${alicloud_rds_backup.default.backup_id}",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":                CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_network_type":            "VPC",
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"category":                         "serverless_ha",
						"db_instance_storage_type":         "cloud_essd",
						"payment_type":                     "Serverless",
						"db_instance_storage":              CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "6",
							"min_capacity": "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "6",
						"serverless_config.0.min_capacity": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id"},
			},
		},
	})
}

func resourceCloneDBInstanceMssqlServerlessHAConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "SQLServer"
    engine_version = "2019_std_sl"
    instance_charge_type = "Serverless"
    category = "serverless_ha"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "SQLServer"
    engine_version = "2019_std_sl"
    category = "serverless_ha"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "vswitche1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}
data "alicloud_vswitches" "vswitche2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.2
}

resource "alicloud_db_instance" "default" {
  engine                   = "SQLServer"
  engine_version           = "2019_std_sl"
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_charge_type     = "Serverless"
  instance_name            = var.name
  zone_id                  = data.alicloud_db_zones.default.ids.1
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.2
  vswitch_id               = join(",", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])
  db_instance_storage_type = "cloud_essd"
  category                 = "serverless_ha"
  serverless_config {
    max_capacity = 8
    min_capacity = 2
  }
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  backup_method     = "Physical"
  backup_type       = "FullBackup"
  remove_from_state = "true"
}

`, name)
}

// TestUnitBuildCloneDBInstanceTags verifies that the Create path injects the
// configured tags atomically into the CloneDBInstance request as Tag.N.Key /
// Tag.N.Value pairs at request-build time. Unlike the existing ACC, which
// only checks the final state tags (set later by Update -> setInstanceTags),
// this request-level unit test fails if the atomic tag injection is removed.
func TestUnitBuildCloneDBInstanceTags(t *testing.T) {
	p := Provider().ResourcesMap
	r, ok := p["alicloud_rds_clone_db_instance"]
	if !ok {
		t.Fatalf("resource alicloud_rds_clone_db_instance not registered in provider")
	}
	d, err := schema.InternalMap(r.Schema).Data(nil, nil)
	if err != nil {
		t.Fatalf("failed to build ResourceData: %v", err)
	}

	tags := map[string]interface{}{
		"env":  "prod",
		"team": "rds",
	}
	if err := d.Set("tags", tags); err != nil {
		t.Fatalf("failed to set tags: %v", err)
	}

	request := make(map[string]interface{})
	buildCloneDBInstanceTags(d, request)

	// At least two configured tags must be present as paired Tag.N.Key /
	// Tag.N.Value entries in the CloneDBInstance request.
	assert.Contains(t, request, "Tag.1.Key", "CloneDBInstance request must contain Tag.1.Key")
	assert.Contains(t, request, "Tag.1.Value", "CloneDBInstance request must contain Tag.1.Value")
	assert.Contains(t, request, "Tag.2.Key", "CloneDBInstance request must contain Tag.2.Key")
	assert.Contains(t, request, "Tag.2.Value", "CloneDBInstance request must contain Tag.2.Value")

	// Each Tag.N.Key must be paired with its matching Tag.N.Value, and every
	// configured tag must appear with its expected value (order is irrelevant
	// since tags come from a map).
	paired := make(map[string]string, len(tags))
	for i := 1; i <= len(tags); i++ {
		key := request[fmt.Sprintf("Tag.%d.Key", i)]
		value := request[fmt.Sprintf("Tag.%d.Value", i)]
		if key == nil || value == nil {
			t.Fatalf("Tag.%d pair incomplete in CloneDBInstance request: key=%v value=%v", i, key, value)
		}
		paired[key.(string)] = value.(string)
	}
	for k, v := range tags {
		got, ok := paired[k]
		if !ok {
			t.Errorf("configured tag %q missing from CloneDBInstance request", k)
			continue
		}
		assert.Equal(t, v, got, "tag %q value mismatch in CloneDBInstance request", k)
	}
}

// TestUnitBuildCloneDBInstanceTags_NoTags verifies that no Tag.N.* entries are
// injected when no tags are configured, guarding against stray tag encoding.
func TestUnitBuildCloneDBInstanceTags_NoTags(t *testing.T) {
	p := Provider().ResourcesMap
	r, ok := p["alicloud_rds_clone_db_instance"]
	if !ok {
		t.Fatalf("resource alicloud_rds_clone_db_instance not registered in provider")
	}
	d, err := schema.InternalMap(r.Schema).Data(nil, nil)
	if err != nil {
		t.Fatalf("failed to build ResourceData: %v", err)
	}

	request := make(map[string]interface{})
	buildCloneDBInstanceTags(d, request)
	for k := range request {
		if strings.HasPrefix(k, "Tag.") {
			t.Errorf("unexpected Tag.* entry %q in CloneDBInstance request when no tags configured", k)
		}
	}
}

// TestUnitCloneDBInstanceCreateRequestInjectsTags asserts that the production
// CloneDBInstance request constructor — buildCloneDBInstanceCreateRequest, to
// which resourceAlicloudRdsCloneDbInstanceCreate delegates all request building
// and which returns the exact map sent to client.RpcPost for the non-serverless
// CloneDBInstance call — injects the configured tags atomically as paired
// Tag.N.Key / Tag.N.Value entries.
//
// This is the request-level mutation guard the original ACC lacked. The ACC only
// checks final-state tags (set later by Update -> setInstanceTags), so removing
// the atomic tag injection still passed. This test instead inspects the actual
// creation request: removing ONLY the buildCloneDBInstanceTags(d, request) call
// from buildCloneDBInstanceCreateRequest (while leaving the leaf helper and its
// direct unit tests, TestUnitBuildCloneDBInstanceTags / _NoTags, untouched)
// leaves those helper tests green but makes THIS test fail, because the returned
// CloneDBInstance request no longer carries Tag.N.Key / Tag.N.Value.
//
// No cloud credentials or network access are required: buildCloneDBInstanceCreateRequest
// only reads the exported client.SourceIp / client.RegionId fields, so a zero-value
// client is sufficient and no gomonkey patching of RPC methods (which Go 1.25
// register-ABI makes unreliable on concrete-type pointer-receiver methods such as
// AliyunClient.RpcPost / RdsService.DescribeDBInstance) is needed.
func TestUnitCloneDBInstanceCreateRequestInjectsTags(t *testing.T) {
	p := Provider().ResourcesMap
	r, ok := p["alicloud_rds_clone_db_instance"]
	if !ok {
		t.Fatalf("resource alicloud_rds_clone_db_instance not registered in provider")
	}
	d, err := schema.InternalMap(r.Schema).Data(nil, nil)
	if err != nil {
		t.Fatalf("failed to build ResourceData: %v", err)
	}

	tags := map[string]interface{}{
		"env":  "prod",
		"team": "rds",
	}
	sourceInstance := "rm-source-mock"
	// Explicit literal-key Set calls populate each configured attribute on
	// ResourceData exactly as terraform apply would, so the production request
	// constructor reads them on the real create path. Literal keys are required
	// by tflint R001 (variable-key d.Set is flagged); the six input values and
	// the two-tag pairing assertions below are unchanged.
	if e := d.Set("source_db_instance_id", sourceInstance); e != nil {
		t.Fatalf("failed to set source_db_instance_id: %v", e)
	}
	if e := d.Set("payment_type", "PayAsYouGo"); e != nil {
		t.Fatalf("failed to set payment_type: %v", e)
	}
	if e := d.Set("db_instance_class", "mysql.n2.medium.25"); e != nil {
		t.Fatalf("failed to set db_instance_class: %v", e)
	}
	if e := d.Set("db_instance_storage", 50); e != nil {
		t.Fatalf("failed to set db_instance_storage: %v", e)
	}
	if e := d.Set("db_instance_storage_type", "cloud_essd"); e != nil {
		t.Fatalf("failed to set db_instance_storage_type: %v", e)
	}
	if e := d.Set("tags", tags); e != nil {
		t.Fatalf("failed to set tags: %v", e)
	}

	// Zero-value client: only exported SourceIp/RegionId are read during request
	// building. No credentials, endpoints or network access are touched.
	client := &connectivity.AliyunClient{}

	// The production request constructor Create delegates to; the returned map is
	// the exact CloneDBInstance request sent to RpcPost (non-serverless path).
	request := buildCloneDBInstanceCreateRequest(d, client)

	// Sanity: this is the real CloneDBInstance request, not a detached helper —
	// the source instance id and storage type from schema are present, proving
	// the production request-construction path ran.
	assert.Equal(t, sourceInstance, request["DBInstanceId"],
		"buildCloneDBInstanceCreateRequest must populate DBInstanceId from source_db_instance_id, proving it is the production request path")
	assert.Equal(t, "cloud_essd", request["DBInstanceStorageType"],
		"buildCloneDBInstanceCreateRequest must populate DBInstanceStorageType, proving it is the production request path")

	// At least two configured tags must be present as paired Tag.N.Key /
	// Tag.N.Value entries in the actual CloneDBInstance request.
	assert.Contains(t, request, "Tag.1.Key", "CloneDBInstance request must contain Tag.1.Key")
	assert.Contains(t, request, "Tag.1.Value", "CloneDBInstance request must contain Tag.1.Value")
	assert.Contains(t, request, "Tag.2.Key", "CloneDBInstance request must contain Tag.2.Key")
	assert.Contains(t, request, "Tag.2.Value", "CloneDBInstance request must contain Tag.2.Value")

	paired := make(map[string]string, len(tags))
	for i := 1; i <= len(tags); i++ {
		key := request[fmt.Sprintf("Tag.%d.Key", i)]
		value := request[fmt.Sprintf("Tag.%d.Value", i)]
		if key == nil || value == nil {
			t.Fatalf("Tag.%d pair incomplete in CloneDBInstance request: key=%v value=%v", i, key, value)
		}
		paired[key.(string)] = value.(string)
	}
	for k, v := range tags {
		got, ok := paired[k]
		if !ok {
			t.Errorf("configured tag %q missing from CloneDBInstance request", k)
			continue
		}
		assert.Equal(t, v, got, "tag %q value mismatch in CloneDBInstance request", k)
	}
}
