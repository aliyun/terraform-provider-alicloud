package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_db_zones.default.ids.0
  cidr_block   = "172.16.0.0/24"
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}
`, name)
}

var cloneInstanceBasicMap = map[string]string{}

func TestAccAlicloudRdsCloneDBInstancePostgreSQLSSL(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_rds_clone_db_instance.default"
	ra := resourceAttrInit(resourceId, cloneInstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
					"db_instance_class": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
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
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "1",
						"ca_type":         "aliyun",
						"acl":             "perfer",
						"replication_acl": "perfer",
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
					"connection_string_prefix":    "rm-ccccccc",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "backup_id", "source_db_instance_id", "client_ca_enabled", "client_crl_enabled", "connection_string_prefix", "ssl_enabled"},
			},
		},
	})
}
