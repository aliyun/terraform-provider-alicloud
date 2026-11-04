package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudRdsDBConnectionConfigMySQL covers the MySQL connection
// lifecycle. general_group_name is set on create so the Allocate request path
// is exercised, but it only applies to MyBase MySQL Standard Edition instances;
// a regular MySQL instance does not return it, so this case does not prove its
// read-back or import round-trip.
func TestAccAliCloudRdsDBConnectionConfigMySQL(t *testing.T) {
	var v map[string]interface{}
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)
	var basicMap = map[string]string{
		"instance_id":       acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + fmt.Sprintf("^tf-testacc%s.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
		"ip_address":        acctest.CHECKSET,
	}
	resourceId := "alicloud_db_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDBConnectionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_db_instance.default.id}",
					"connection_prefix":  fmt.Sprintf("tf-testacc%s", rand),
					"general_group_name": "rgc-tftest123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// GeneralGroupName is an instance attribute that
				// DescribeDBInstanceAttribute returns only for MyBase MySQL Standard
				// Edition instances (see the live AllocateInstancePublicConnection
				// contract); on a regular MySQL instance the configured value is
				// kept and cannot be verified on import.
				ImportStateVerifyIgnore: []string{"general_group_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBConnectionConfigPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)
	var basicMap = map[string]string{
		"instance_id":       acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + fmt.Sprintf("^tf-testacc%s.pgsql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
		"ip_address":        acctest.CHECKSET,
	}
	resourceId := "alicloud_db_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDBConnectionConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
					"babelfish_port":    "1533",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"babelfish_port": "1533",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"babelfish_port": "1433",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"babelfish_port": "1433",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudRdsDBConnection_pgbouncerPort verifies that pgbouncer_port is
// sent through AllocateInstancePublicConnection on create and through
// ModifyDBInstanceConnectionString on update, and that both values are read back
// from DescribeDBInstanceNetInfo (asserted against the API response, not only
// against the configured state); the import step then verifies the read-back.
// PgBouncer is enabled by the instance fixture, which the instance resource
// applies during creation, so the port is accepted when the connection is
// created.
func TestAccAliCloudRdsDBConnection_pgbouncerPort(t *testing.T) {
	var v map[string]interface{}
	suffix := tfacctest.RandStringFromCharSet(6, tfacctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("tf-testacc%s", suffix)
	name := fmt.Sprintf("tf-testAccDBconnection%s", suffix)
	resourceId := "alicloud_db_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"instance_id":    acctest.CHECKSET,
		"port":           "3306",
		"pgbouncer_port": acctest.CHECKSET,
	})
	serviceFunc := func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDBConnectionConfigPgbouncerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_instance.default.id}",
					"connection_prefix": prefix,
					"pgbouncer_port":    "6432",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pgbouncer_port": "6432",
					}),
					resource.TestCheckResourceAttr(resourceId, "connection_prefix", prefix),
					testAccCheckDBConnectionPgBouncerPort(prefix, "6432"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pgbouncer_port": "6433",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pgbouncer_port": "6433",
					}),
					testAccCheckDBConnectionPgBouncerPort(prefix, "6433"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckDBConnectionPgBouncerPort reads the connection back from
// DescribeDBInstanceNetInfo and compares the PGBouncerPort returned by the API.
// It distinguishes a real API read-back from a value that was merely kept in
// the local state.
func testAccCheckDBConnectionPgBouncerPort(prefix, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_db_connection.default"]
		if !ok {
			return fmt.Errorf("alicloud_db_connection.default is not in the state")
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := (rds.NewRdsService(client)).DescribeDBConnection(instanceId + ":" + prefix)
		if err != nil {
			return err
		}
		if got := fmt.Sprint(object["PGBouncerPort"]); got != expected {
			return fmt.Errorf("DescribeDBInstanceNetInfo returned PGBouncerPort %q, want %q", got, expected)
		}
		return nil
	}
}

func resourceDBConnectionConfigDependence(name string) string {
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

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
}
`, name)
}

func resourceDBConnectionConfigPostgreSQLDependence(name string) string {
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
    zone_id = data.alicloud_db_zones.default.zones[0].id
	engine = "PostgreSQL"
	engine_version = "13.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	commodity_code = "bards"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids[0]
 zone_id = data.alicloud_db_zones.default.ids[0]
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
  # Babelfish requires a dedicated (x4) PostgreSQL family; both the class and its
  # storage range come from the data source so a sold-out class is not selected.
  pg_babelfish_classes = [for c in data.alicloud_db_instance_classes.default.instance_classes : c if can(regex("^pg\\.x4\\.", c.instance_class))]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "PostgreSQL"
	engine_version = "13.0"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_type = local.pg_babelfish_classes[0].instance_class
	instance_storage = local.pg_babelfish_classes[0].storage_range[0].min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_ips = ["10.168.1.12", "100.69.7.112"]
	babelfish_config {
      babelfish_enabled = "true"
	  migration_mode = "single-db"
	  master_username = "tftest"
	  master_user_password = "!Q2w3e4r"
    }
}
`, name)
}

func resourceDBConnectionConfigPgbouncerDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "17.0"
}

data "alicloud_db_instance_classes" "default" {
	zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "17.0"
	category = "HighAvailability"
	db_instance_storage_type = "general_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
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
}

resource "alicloud_db_instance" "default" {
	engine = "PostgreSQL"
	engine_version = "17.0"
	category = "HighAvailability"
	db_instance_storage_type = "general_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	instance_charge_type = "Postpaid"
	vswitch_id = local.vswitch_id
	instance_name = var.name
	pg_bouncer_enabled = "true"
	security_ips = ["10.168.1.12", "100.69.7.112"]
}
`, name)
}
