package alicloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUnitNormalizePolarDBPolarFsInteger(t *testing.T) {
	for _, value := range []interface{}{0, float64(100), "0.0", json.Number("100.0")} {
		got, err := normalizePolarDBPolarFsInteger(value)
		if err != nil {
			t.Fatalf("unexpected normalization error for %v: %s", value, err)
		}
		if value == float64(100) || value == json.Number("100.0") {
			if got != 100 {
				t.Fatalf("expected 100 for %v, got %d", value, got)
			}
		} else if got != 0 {
			t.Fatalf("expected 0 for %v, got %d", value, got)
		}
	}
	if _, err := normalizePolarDBPolarFsInteger("1.5"); err == nil {
		t.Fatal("fractional values must not be silently truncated")
	}
}

func TestUnitPolarDBPolarFsDBTypeWireValue(t *testing.T) {
	tests := []struct {
		terraformValue string
		apiEcho        string
	}{
		{terraformValue: "MySQL", apiEcho: "polardb_mysql"},
		{terraformValue: "PostgreSQL", apiEcho: "polardb_pg"},
	}
	for _, test := range tests {
		d := schema.TestResourceDataRaw(t, resourceAlicloudPolarDBPolarFs().Schema, map[string]interface{}{"db_type": test.terraformValue})
		if got := buildPolarDBPolarFsCreateRequest(d, "cn-beijing")["DBType"]; got != test.apiEcho {
			t.Fatalf("CreatePolarFs DBType for %q: got %q, want %q", test.terraformValue, got, test.apiEcho)
		}
		if got := flattenPolarDBPolarFsDBType(test.apiEcho); got != test.terraformValue {
			t.Fatalf("flatten %q: got %q, want %q", test.apiEcho, got, test.terraformValue)
		}
	}
}

func TestUnitPolarDBPolarFsAttributeReadKeepsCreateOnlyArguments(t *testing.T) {
	resourceSchema := resourceAlicloudPolarDBPolarFs()
	object := map[string]interface{}{
		"AcceleratingEnable": "ON", "AcceleratedStorageSpace": float64(500),
	}
	for _, test := range []struct {
		config     map[string]interface{}
		wantSwitch string
		wantSize   int
	}{
		{config: map[string]interface{}{"accelerate_switch": "ONLY", "accelerate_storage_size": 500}, wantSwitch: "ONLY", wantSize: 500},
		{config: map[string]interface{}{}},
	} {
		d := schema.TestResourceDataRaw(t, resourceSchema.Schema, test.config)
		if err := flattenPolarDBPolarFsAttribute(d, object); err != nil {
			t.Fatal(err)
		}
		if got := d.Get("accelerate_switch"); got != test.wantSwitch {
			t.Fatalf("accelerate_switch changed during Read: got %q, config %#v", got, test.config)
		}
		if got := d.Get("accelerate_storage_size"); got != test.wantSize {
			t.Fatalf("accelerate_storage_size changed during Read: got %v, config %#v", got, test.config)
		}
		if got := d.Get("accelerating_enable"); got != "ON" {
			t.Fatalf("accelerating_enable was not read: %v", got)
		}
		if got := d.Get("accelerated_storage_space"); got != float64(500) {
			t.Fatalf("accelerated_storage_space was not read: %v", got)
		}
	}
}

func TestUnitPolarDBRpcPostRetryableLockTimeout(t *testing.T) {
	if !polarDBRpcPostRetryable(errors.New("LockTimeout")) {
		t.Fatal("LockTimeout must be retried")
	}
	if polarDBRpcPostRetryable(errors.New("InvalidParameter")) {
		t.Fatal("InvalidParameter must not be retried")
	}
}

func TestUnitPolarDBPolarFsAccelerateStorageSizeValidation(t *testing.T) {
	if resourceAlicloudPolarDBPolarFs().Schema["accelerate_storage_size"].ValidateFunc != nil {
		t.Fatal("the API contract does not document a minimum or increment for accelerate_storage_size")
	}
}

func TestUnitPolarDBPolarFsDBTypeValidation(t *testing.T) {
	resourceSchema := resourceAlicloudPolarDBPolarFs()
	validateDBType := resourceSchema.Schema["db_type"].ValidateFunc
	for _, value := range []string{"MySQL", "PostgreSQL"} {
		if _, errors := validateDBType(value, "db_type"); len(errors) != 0 {
			t.Fatalf("db_type %q must be accepted: %v", value, errors)
		}
	}
	for _, value := range []string{"polardb_mysql", "polardb_pg", "polar_agentic_db"} {
		if _, errors := validateDBType(value, "db_type"); len(errors) == 0 {
			t.Fatalf("internal enum value %q must not be sent as the OpenAPI DBType", value)
		}
	}
}

func TestUnitRedactPolarDBSensitiveValue(t *testing.T) {
	input := map[string]interface{}{
		"CustomOssAk": "ak", "CustomOssSk": "sk",
		"MountInfo": map[string]interface{}{"Token": "token", "PolarFsCluster": "pfs-1"},
	}
	got := redactPolarDBSensitiveValue(input).(map[string]interface{})
	if got["CustomOssAk"] != "***" || got["CustomOssSk"] != "***" {
		t.Fatalf("OSS credentials were not redacted: %#v", got)
	}
	mountInfo := got["MountInfo"].(map[string]interface{})
	if mountInfo["Token"] != "***" || mountInfo["PolarFsCluster"] != "pfs-1" {
		t.Fatalf("nested token redaction failed: %#v", mountInfo)
	}
}

func TestAccAliCloudPolarDBPolarFs_complete(t *testing.T) {
	resourceId := "alicloud_polardb_polar_fs.default"
	endpointId := "alicloud_polardb_polar_fs_endpoint.default"
	instancesId := "data.alicloud_polardb_polar_fs_instances.default"
	tokenId := "data.alicloud_polardb_lakebase_tenant_token.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, "default", testAccPolarDBPolarFsDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBPolarFsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                 "PostgreSQL",
					"creation_category":       "basic",
					"storage_type":            "city_redundancy",
					"storage_space":           50,
					"pay_type":                "Postpaid",
					"accelerate_switch":       "ON",
					"accelerate_storage_size": 5000,
					"accelerate_type":         "juice",
					"db_cluster_id":           "${alicloud_polardb_cluster.polarfs.id}",
					"vpc_id":                  "${data.alicloud_vswitches.polarfs.vswitches[0].vpc_id}",
					"vswitch_id":              "${data.alicloud_vswitches.polarfs.vswitches[0].id}",
					"zone_id":                 "${data.alicloud_vswitches.polarfs.vswitches[0].zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "db_type", "PostgreSQL"),
					resource.TestCheckResourceAttr(resourceId, "storage_type", "city_redundancy"),
					resource.TestCheckResourceAttr(resourceId, "storage_space", "50"),
					resource.TestCheckResourceAttr(resourceId, "pay_type", "Postpaid"),
					resource.TestCheckResourceAttr(resourceId, "creation_category", "basic"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_switch", "ON"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_storage_size", "5000"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_type", "juice"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "meta_instance_name"),
					resource.TestCheckResourceAttr(endpointId, "endpoint_type", "S3Gateway"),
					resource.TestCheckResourceAttrSet(endpointId, "db_endpoint_id"),
					resource.TestCheckResourceAttr(instancesId, "polar_fs_instances.#", "1"),
					resource.TestCheckResourceAttrPair(instancesId, "polar_fs_instances.0.id", resourceId, "id"),
					resource.TestCheckResourceAttrSet(tokenId, "token"),
					resource.TestCheckResourceAttr(tokenId, "subdir", "/terraform-acc"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"polar_fs_path", "storage_used", "accelerate_switch", "accelerate_storage_size"},
			},
		},
	})
}

func testAccPolarDBPolarFsDependence(string) string {
	return `
data "alicloud_vswitches" "polarfs" {
  status  = "Available"
  zone_id = "cn-beijing-l"
}

resource "alicloud_polardb_cluster" "polarfs" {
  db_type           = "PostgreSQL"
  db_version        = "14"
  pay_type          = "PostPaid"
  db_node_class     = "polar.pg.g1.tiny.c"
  db_node_num       = 1
  vswitch_id        = data.alicloud_vswitches.polarfs.vswitches[0].id
  zone_id           = data.alicloud_vswitches.polarfs.vswitches[0].zone_id
  creation_category = "SENormal"
  storage_type      = "ESSDAUTOPL"
  storage_space     = 40
  storage_pay_type  = "PostPaid"
  hot_standby_cluster = "OFF"
  encrypt_new_tables = "OFF"
  description       = "tf-testAccPolarDBPolarFs"
}

resource "alicloud_polardb_polar_fs_endpoint" "default" {
  db_cluster_id           = alicloud_polardb_cluster.polarfs.id
  polar_fs_instance_id    = alicloud_polardb_polar_fs.default.id
  endpoint_type           = "S3Gateway"
  vpc_id                  = data.alicloud_vswitches.polarfs.vswitches[0].vpc_id
  vswitch_id              = data.alicloud_vswitches.polarfs.vswitches[0].id
  db_endpoint_description = "terraform-acc"
}

data "alicloud_polardb_polar_fs_instances" "default" {
  ids = [alicloud_polardb_polar_fs.default.id]
}

data "alicloud_polardb_lakebase_tenant_token" "default" {
  polar_fs_instance_id = alicloud_polardb_polar_fs.default.id
  db_cluster_id        = alicloud_polardb_cluster.polarfs.id
  subdir               = "/terraform-acc"
}
`
}

func testAccCheckPolarDBPolarFsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_polar_fs" {
			continue
		}
		_, err := describePolarDBPolarFsAttribute(client, rs.Primary.ID, rs.Primary.Attributes["db_cluster_id"])
		if err == nil {
			return fmt.Errorf("PolarFS instance %s still exists", rs.Primary.ID)
		}
		if !NotFoundError(err) {
			return err
		}
	}
	return nil
}
