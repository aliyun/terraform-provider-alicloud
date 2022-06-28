package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsDBReadonlyInstance_update(t *testing.T) {
	var instance map[string]interface{}
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
					"vswitch_id":            "${local.vswitch_id}",
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
			// upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alicloud_db_instance.default.instance_storage + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": CHECKSET}),
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
					"vswitch_id":            "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":    name,
						"instance_storage": CHECKSET,
					}),
				),
			},
		},
	})

}

func TestAccAlicloudRdsDBReadonlyInstancePostgreSQL_update(t *testing.T) {
	var instance map[string]interface{}
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
					"zone_id":               "${alicloud_db_instance.default.zone_id}",
					"engine_version":        "${alicloud_db_instance.default.engine_version}",
					"instance_type":         "${data.alicloud_db_instance_classes.ro.instance_classes.0.instance_class}",
					"instance_storage":      "${alicloud_db_instance.default.instance_storage}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${local.vswitch_id}",
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
					"vswitch_id":                  "${local.vswitch_id}",
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
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_ips = ["10.168.1.12", "100.69.7.112"]
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
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
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
	vswitch_id = local.vswitch_id
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
