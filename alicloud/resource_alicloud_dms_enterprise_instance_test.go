package alicloud

import (
	"fmt"
	"os"
	"testing"

	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDmsEnterprise(t *testing.T) {
	resourceId := "alicloud_dms_enterprise_instance.default"
	var v dms_enterprise.Instance
	ra := resourceAttrInit(resourceId, testAccCheckKeyValueInMapsForDMS)

	serviceFunc := func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDmsConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"dba_uid":           os.Getenv("DBA_UID"),
					"host":              "${alicloud_db_connection.foo.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "13429",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "dev",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "600",
					"ecs_region":        "cn-shanghai",
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"dba_uid":           os.Getenv("DBA_UID"),
						"host":              "${alicloud_db_connection.foo.connection_string}",
						"port":              "3306",
						"network_type":      "VPC",
						"safe_rule":         "自由操作",
						"tid":               "13429",
						"instance_type":     "mysql",
						"instance_source":   "RDS",
						"env_type":          "dev",
						"database_user":     "${alicloud_db_account.account.name}",
						"database_password": CHECKSET,
						"instance_alias":    name,
						"query_timeout":     "70",
						"export_timeout":    "600",
						"ecs_region":        "cn-shanghai",
						"ddl_online":        "0",
						"use_dsql":          "0",
						"data_link_name":    "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password", "dba_uid", "network_type", "port", "safe_rule", "tid"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"safe_rule":      "安全协同",
					"use_dsql":       "1",
					"data_link_name": "testname",
					"ddl_online":     "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"safe_rule":      "安全协同",
						"use_dsql":       "1",
						"data_link_name": "testname",
						"ddl_online":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_type": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": "other_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": "other_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_dsql": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_dsql": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_timeout": "77",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_timeout": "77",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"dba_uid":           os.Getenv("DBA_UID"),
					"host":              "${alicloud_db_connection.foo.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "13429",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "dev",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "600",
					"ecs_region":        "cn-shanghai",
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"dba_uid":           os.Getenv("DBA_UID"),
						"host":              "${alicloud_db_connection.foo.connection_string}",
						"port":              "3306",
						"network_type":      "VPC",
						"safe_rule":         "自由操作",
						"tid":               "13429",
						"instance_type":     "mysql",
						"instance_source":   "RDS",
						"env_type":          "dev",
						"database_user":     "${alicloud_db_account.account.name}",
						"database_password": CHECKSET,
						"instance_alias":    name,
						"query_timeout":     "70",
						"export_timeout":    "600",
						"ecs_region":        "cn-shanghai",
						"ddl_online":        "0",
						"use_dsql":          "0",
						"data_link_name":    "",
					}),
				),
			},
		},
	})
}

var testAccCheckKeyValueInMapsForDMS = map[string]string{}

func resourceDmsConfigDependence(name string) string {
	return fmt.Sprintf(`    
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbconnectionbasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.7"
  instance_type    = "rds.mysql.t1.small"
  instance_storage = "10"
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  security_ips     = ["100.104.5.0/24","192.168.0.6"]
}

resource "alicloud_db_connection" "foo" {
  instance_id       = "${alicloud_db_instance.instance.id}"
  port              = 3306
}
resource "alicloud_db_account" "account" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "tftestnormal"
  password    = "Test12345"
}
`)
}
