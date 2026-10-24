package alicloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestUnitDBConnectionOptionalParams pins the connection-string parameters that
// are optional on both AllocateInstancePublicConnection and
// ModifyDBInstanceConnectionString: pgbouncer_port and general_group_name are
// only sent when they are configured, so an omitted value never changes the
// remote defaults.
func TestUnitDBConnectionOptionalParams(t *testing.T) {
	testCases := []struct {
		name    string
		raw     map[string]interface{}
		want    map[string]interface{}
		notWant []string
	}{
		{
			name: "both configured",
			raw: map[string]interface{}{
				"instance_id":        "rm-test123",
				"pgbouncer_port":     "6432",
				"general_group_name": "rgc-test123",
			},
			want: map[string]interface{}{
				"PGBouncerPort":    "6432",
				"GeneralGroupName": "rgc-test123",
			},
		},
		{
			name: "both omitted",
			raw: map[string]interface{}{
				"instance_id": "rm-test123",
			},
			notWant: []string{"PGBouncerPort", "GeneralGroupName"},
		},
		{
			name: "only pgbouncer port",
			raw: map[string]interface{}{
				"instance_id":    "rm-test123",
				"pgbouncer_port": "6433",
			},
			want: map[string]interface{}{
				"PGBouncerPort": "6433",
			},
			notWant: []string{"GeneralGroupName"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAlicloudDBConnection().Schema, tc.raw)
			request := map[string]interface{}{}
			setDBConnectionOptionalParams(request, d)
			for k, v := range tc.want {
				if got := request[k]; got != v {
					t.Fatalf("request[%s] = %v, want %v", k, got, v)
				}
			}
			for _, k := range tc.notWant {
				if got, ok := request[k]; ok {
					t.Fatalf("request[%s] = %v, want the key to be absent", k, got)
				}
			}
		})
	}
}

func TestAccAliCloudRdsDBConnectionConfigMySQL(t *testing.T) {
	var v map[string]interface{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)
	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("^tf-testacc%s.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "alicloud_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_db_instance.default.id}",
					"connection_prefix":  fmt.Sprintf("tf-testacc%s", rand),
					"general_group_name": "rgc-tftest123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					resource.TestCheckResourceAttr(resourceId, "general_group_name", "rgc-tftest123"),
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
					"port":               "3333",
					"general_group_name": "rgc-tftest123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":               "3333",
						"general_group_name": "rgc-tftest123",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBConnectionConfigPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)
	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("^tf-testacc%s.pgsql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "alicloud_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("tf-testacc%s", suffix)
	name := fmt.Sprintf("tf-testAccDBconnection%s", suffix)
	resourceId := "alicloud_db_connection.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"instance_id":    CHECKSET,
		"port":           "3306",
		"pgbouncer_port": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigPgbouncerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		object, err := (&RdsService{client}).DescribeDBConnection(instanceId + ":" + prefix)
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

func dbConnectionMockResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceAlicloudDBConnection().Schema, map[string]interface{}{
		"instance_id":       "rm-test",
		"connection_prefix": "abc",
		"pgbouncer_port":    "6432",
	})
	d.SetId("rm-test:abc")
	return d
}

func dbConnectionMockHandler(t *testing.T, instance map[string]interface{}, netInfo map[string]interface{}) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{netInfo},
				},
			})
		case "DescribeDBInstanceAttribute":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{instance},
				},
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	}
}

func dbConnectionMockNetInfo(pgbouncerPort interface{}) map[string]interface{} {
	netInfo := map[string]interface{}{
		"ConnectionString": "abc.mysql.rds.aliyuncs.com",
		"Port":             "3306",
		"IPAddress":        "10.0.0.1",
	}
	if pgbouncerPort != nil {
		netInfo["PGBouncerPort"] = pgbouncerPort
	}
	return netInfo
}

// TestUnitDBConnectionPgBouncerPortReadOmission covers QA-NEW-03: on PostgreSQL
// an omitted PGBouncerPort means the configured value is unconfirmed and must
// not stay in state, whether the instance reports PgBouncer as enabled or its
// flag is unknown. Only an instance that explicitly reports PgBouncer disabled
// (or a non-PostgreSQL engine, where the parameter does not apply) keeps the
// configured value instead of producing a permanent diff.
func TestUnitDBConnectionPgBouncerPortReadOmission(t *testing.T) {
	cases := []struct {
		name          string
		engine        string
		pgbouncerFlag interface{}
		netInfoPort   interface{}
		wantPort      string
	}{
		{name: "enabled and present", engine: "PostgreSQL", pgbouncerFlag: "true", netInfoPort: "6432", wantPort: "6432"},
		{name: "enabled and omitted", engine: "PostgreSQL", pgbouncerFlag: "true", netInfoPort: nil, wantPort: ""},
		{name: "disabled and omitted", engine: "PostgreSQL", pgbouncerFlag: "false", netInfoPort: nil, wantPort: "6432"},
		{name: "unknown and omitted", engine: "PostgreSQL", pgbouncerFlag: nil, netInfoPort: nil, wantPort: ""},
		{name: "mysql and omitted", engine: "MySQL", pgbouncerFlag: nil, netInfoPort: nil, wantPort: "6432"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			instance := map[string]interface{}{"DBInstanceStatus": "Running", "Engine": tc.engine}
			if tc.pgbouncerFlag != nil {
				instance["PGBouncerEnabled"] = tc.pgbouncerFlag
			}
			client := dbInstanceMockClient(t, dbConnectionMockHandler(t, instance, dbConnectionMockNetInfo(tc.netInfoPort)))
			d := dbConnectionMockResourceData(t)
			if err := resourceAlicloudDBConnectionRead(d, client); err != nil {
				t.Fatalf("read failed: %v", err)
			}
			if got := d.Get("pgbouncer_port").(string); got != tc.wantPort {
				t.Fatalf("pgbouncer_port = %q, want %q", got, tc.wantPort)
			}
		})
	}
}

// TestUnitDBConnectionParentInstanceNotFound covers QA-NEW-09: a confirmed
// missing parent instance clears the connection from state instead of keeping a
// dangling record.
func TestUnitDBConnectionParentInstanceNotFound(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
				},
			})
		case "DescribeDBInstanceAttribute":
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code": "InvalidDBInstanceId.NotFound", "Message": "instance not found",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	d := dbConnectionMockResourceData(t)
	if err := resourceAlicloudDBConnectionRead(d, client); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("connection must be removed from state, id = %q", d.Id())
	}
}

// TestUnitDBConnectionPgBouncerPortConvergence covers the shared Create/Update
// waiter: an omitted port is not convergence, and the waiter succeeds once the
// API reports the requested value.
func TestUnitDBConnectionPgBouncerPortConvergence(t *testing.T) {
	t.Run("eventually converged", func(t *testing.T) {
		var calls int
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceNetInfo" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			calls++
			port := interface{}(nil)
			if calls > 1 {
				port = "6433"
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo(port)},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortConverged(RdsService{client}, "rm-test:abc", "6433", 5*time.Second); err != nil {
			t.Fatalf("waiter failed: %v", err)
		}
	})

	t.Run("never confirmed", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo(nil)},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortConverged(RdsService{client}, "rm-test:abc", "6432", 2*time.Second); err == nil {
			t.Fatal("an omitted port must not be reported as converged")
		}
	})
}

// TestUnitDBConnectionPgBouncerPortSkipWhenDisabled covers the gate in front of
// the port waiter: when the instance reports PgBouncer disabled the parameter is
// not applicable, so the waiter must not block on a value the API never returns.
func TestUnitDBConnectionPgBouncerPortSkipWhenDisabled(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "PostgreSQL", "PGBouncerEnabled": "false"}},
				},
			})
		case "DescribeDBInstanceNetInfo":
			t.Error("the port waiter must not query the connection when PgBouncer is disabled")
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err != nil {
		t.Fatalf("waiter must skip instances with PgBouncer disabled, got %v", err)
	}
}

// TestUnitDBConnectionPgBouncerPortUnknownFlag covers QA-NEW-03's remaining
// hole: an omitted PGBouncerEnabled is unknown, not "disabled". The waiter must
// resolve the flag before deciding, and must fail instead of accepting an
// unconfirmed port when the flag never becomes known.
func TestUnitDBConnectionPgBouncerPortUnknownFlag(t *testing.T) {
	t.Run("unknown then enabled", func(t *testing.T) {
		var attributeCalls int
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeDBInstanceAttribute":
				attributeCalls++
				attribute := map[string]interface{}{"Engine": "PostgreSQL", "DBInstanceStatus": "CONFIG_SWITCHING"}
				if attributeCalls > 1 {
					attribute["PGBouncerEnabled"] = "true"
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Items": map[string]interface{}{"DBInstanceAttribute": []interface{}{attribute}},
				})
			case "DescribeDBInstanceNetInfo":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"DBInstanceNetInfos": map[string]interface{}{
						"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
					},
				})
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 5*time.Second); err != nil {
			t.Fatalf("waiter must converge once the flag becomes known, got %v", err)
		}
		if attributeCalls < 2 {
			t.Fatalf("expected the flag to be re-read, got %d calls", attributeCalls)
		}
	})

	t.Run("never known", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceAttribute" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "PostgreSQL", "DBInstanceStatus": "CONFIG_SWITCHING"}},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err == nil {
			t.Fatal("an unconfirmed PgBouncer state must not be accepted")
		}
	})

	t.Run("non postgresql", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceAttribute" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "MySQL"}},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err != nil {
			t.Fatalf("the parameter does not apply to other engines, got %v", err)
		}
	})
}

// TestUnitDBConnectionParentInstanceLookupError covers N-06: a parent lookup
// that fails for a reason other than NotFound must propagate through the public
// Read and must not drop the connection from state.
func TestUnitDBConnectionParentInstanceLookupError(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
				},
			})
		case "DescribeDBInstanceAttribute":
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code": "Forbidden.RAM", "Message": "not authorized",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	d := dbConnectionMockResourceData(t)
	if err := resourceAlicloudDBConnectionRead(d, client); err == nil {
		t.Fatal("a failed parent lookup must be returned as an error")
	}
	if d.Id() != "rm-test:abc" {
		t.Fatalf("connection id must be preserved on a failed lookup, got %q", d.Id())
	}
}

// TestUnitDBConnectionPgBouncerPortSharedBudget covers N-04: the flag wait and
// the port wait share one budget, so consuming most of it in the first stage
// must not give the second stage a full budget again.
func TestUnitDBConnectionPgBouncerPortSharedBudget(t *testing.T) {
	// The port stage only receives the budget the flag stage left behind: an
	// expired deadline yields an error instead of a fresh budget.
	limit := 15 * time.Minute
	remaining, err := pgBouncerPortRemainingBudget(time.Now().Add(2*time.Minute), limit)
	if err != nil {
		t.Fatalf("a live deadline must yield a budget, got %v", err)
	}
	if remaining <= 0 || remaining > 2*time.Minute {
		t.Fatalf("remaining budget must derive from the deadline, got %s", remaining)
	}

	_, err = pgBouncerPortRemainingBudget(time.Now().Add(-time.Second), limit)
	if err == nil {
		t.Fatal("an exhausted budget must be reported instead of starting a new wait")
	}
	if !strings.Contains(err.Error(), "budget") {
		t.Fatalf("error must explain the exhausted budget, got %v", err)
	}
}
