package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDBReadonlyInstance_update(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "5",
		"engine_version":        "5.6",
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
					"instance_type":         "${alicloud_db_instance.default.instance_type}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": "10"}),
				),
			},
			// upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
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
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${alicloud_db_instance.default.instance_type}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage + 2*data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":    name,
						"instance_storage": "15",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstanceSSL_update(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_readonly_instance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "20",
		"engine_version":        "11.0",
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

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigSSLDependence)
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
					"instance_type":         "${alicloud_db_instance.default.instance_type}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
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
			//upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": "25"}),
				),
			},
			//upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
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
					"instance_type":               "${alicloud_db_instance.default.instance_type}",
					"instance_storage":            "${alicloud_db_instance.default.instance_storage + 2*data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "${alicloud_vswitch.default.id}",
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
						"instance_name":               name,
						"instance_storage":            "30",
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
		},
	})
}

func TestAccAlicloudDBReadonlyInstance_multi(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_readonly_instance.default.1"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "5",
		"engine_version":        "5.6",
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
					"count":                 "2",
					"master_db_instance_id": "${alicloud_db_instance.default.id}",
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${alicloud_db_instance.default.instance_type}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceDBReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "%s"
	}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

	resource "alicloud_db_instance" "default" {
		engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
`, RdsCommonTestCase, name)
}

func resourceDBReadonlyInstanceConfigSSLDependence(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "%s"
	}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PostgreSQL"
  engine_version       = "11.0"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PostgreSQL"
  engine_version       = "11.0"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

	resource "alicloud_db_instance" "default" {
		engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
`, RdsCommonTestCase, name)
}
