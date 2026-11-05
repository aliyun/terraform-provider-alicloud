package polardb

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		d := schema.TestResourceDataRaw(t, ResourcePolarFs().Schema, map[string]interface{}{"db_type": test.terraformValue})
		if got := buildPolarDBPolarFsCreateRequest(d, "cn-beijing")["DBType"]; got != test.apiEcho {
			t.Fatalf("CreatePolarFs DBType for %q: got %q, want %q", test.terraformValue, got, test.apiEcho)
		}
		if got := flattenPolarDBPolarFsDBType(test.apiEcho); got != test.terraformValue {
			t.Fatalf("flatten %q: got %q, want %q", test.apiEcho, got, test.terraformValue)
		}
	}
}
func TestUnitPolarDBPolarFsAttributeReadKeepsCreateOnlyArguments(t *testing.T) {
	resourceSchema := ResourcePolarFs()
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
	if ResourcePolarFs().Schema["accelerate_storage_size"].ValidateFunc != nil {
		t.Fatal("the API contract does not document a minimum or increment for accelerate_storage_size")
	}
}
func TestUnitPolarDBPolarFsDBTypeValidation(t *testing.T) {
	resourceSchema := ResourcePolarFs()
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
func TestUnitPolarDBPolarFsCustomBucketPathListRequest(t *testing.T) {
	resourceSchema := ResourcePolarFs()
	d := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"custom_bucket_path_list": []interface{}{
			map[string]interface{}{"bucket": "pfs-first.oss-cn-beijing-internal.aliyuncs.com", "path": "/data"},
			map[string]interface{}{"bucket": "second-bucket", "path": "/other"},
		},
	})
	request := buildPolarDBPolarFsCreateRequest(d, "cn-beijing")
	list, ok := request["CustomBucketPathList"].([]map[string]interface{})
	if !ok || len(list) != 2 {
		t.Fatalf("CustomBucketPathList was not expanded with both entries: %#v", request["CustomBucketPathList"])
	}
	pairs := make(map[string]string, len(list))
	for _, entry := range list {
		bucket, _ := entry["Bucket"].(string)
		path, _ := entry["Path"].(string)
		pairs[bucket] = path
	}
	if pairs["pfs-first.oss-cn-beijing-internal.aliyuncs.com"] != "/data" || pairs["second-bucket"] != "/other" {
		t.Fatalf("CustomBucketPathList lost bucket/path pairing: %#v", list)
	}

	empty := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"custom_bucket_path_list": []interface{}{
			map[string]interface{}{"bucket": "", "path": ""},
		},
	})
	if _, ok := buildPolarDBPolarFsCreateRequest(empty, "cn-beijing")["CustomBucketPathList"]; ok {
		t.Fatal("an empty custom_bucket_path_list entry must not be sent to CreatePolarFs")
	}

	unset := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{})
	if _, ok := buildPolarDBPolarFsCreateRequest(unset, "cn-beijing")["CustomBucketPathList"]; ok {
		t.Fatal("CustomBucketPathList must not be sent when custom_bucket_path_list is not configured")
	}
}
func TestUnitPolarDBPolarFsCustomBucketPathListRead(t *testing.T) {
	resourceSchema := ResourcePolarFs()
	d := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{})
	object := map[string]interface{}{
		"CustomBucketPathList": []interface{}{
			map[string]interface{}{"Bucket": "pfs-first.oss-cn-beijing-internal.aliyuncs.com", "Path": "/data"},
		},
	}
	if err := flattenPolarDBPolarFsAttribute(d, object); err != nil {
		t.Fatal(err)
	}
	list := d.Get("custom_bucket_path_list").(*schema.Set).List()
	if len(list) != 1 {
		t.Fatalf("custom_bucket_path_list was not read back: %#v", list)
	}
	entry := list[0].(map[string]interface{})
	if entry["bucket"] != "pfs-first.oss-cn-beijing-internal.aliyuncs.com" || entry["path"] != "/data" {
		t.Fatalf("custom_bucket_path_list lost bucket/path: %#v", entry)
	}

	empty := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"custom_bucket_path_list": []interface{}{
			map[string]interface{}{"bucket": "stale", "path": "/stale"},
		},
	})
	if err := flattenPolarDBPolarFsAttribute(empty, map[string]interface{}{"CustomBucketPathList": []interface{}{}}); err != nil {
		t.Fatal(err)
	}
	if got := empty.Get("custom_bucket_path_list").(*schema.Set).Len(); got != 0 {
		t.Fatalf("an empty API response must clear custom_bucket_path_list, got %d entries", got)
	}

	if _, err := flattenPolarDBPolarFsCustomBucketPathList("unexpected"); err == nil {
		t.Fatal("an unexpected CustomBucketPathList payload must return an error")
	}
}
func TestUnitPolarDBPolarFsCustomBucketPathListRequiredForAlluxio(t *testing.T) {
	resourceSchema := ResourcePolarFs()
	for _, test := range []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "alluxio without list",
			config:  map[string]interface{}{"accelerate_type": "alluxio"},
			wantErr: true,
		},
		{
			name: "alluxio with empty entry",
			config: map[string]interface{}{
				"accelerate_type":         "alluxio",
				"custom_bucket_path_list": []interface{}{map[string]interface{}{"bucket": "", "path": ""}},
			},
			wantErr: true,
		},
		{
			name: "alluxio with list",
			config: map[string]interface{}{
				"accelerate_type":         "alluxio",
				"custom_bucket_path_list": []interface{}{map[string]interface{}{"bucket": "tf-bucket.oss-cn-beijing-internal.aliyuncs.com", "path": "/data"}},
			},
		},
		{
			name:   "juice without list",
			config: map[string]interface{}{"accelerate_type": "juice"},
		},
		{
			name:   "accelerate_type unset",
			config: map[string]interface{}{},
		},
	} {
		d := schema.TestResourceDataRaw(t, resourceSchema.Schema, test.config)
		err := validatePolarDBPolarFsCreateRequest(d)
		if test.wantErr && err == nil {
			t.Fatalf("%s: expected an error for accelerate_type=alluxio without a bucket/path pair", test.name)
		}
		if !test.wantErr && err != nil {
			t.Fatalf("%s: unexpected error: %s", test.name, err)
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
