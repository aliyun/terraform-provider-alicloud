package alicloud

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestUnitPolarDBPolarFsBucketPaths(t *testing.T) {
	input := []interface{}{map[string]interface{}{"bucket": "example.oss-cn-beijing-internal.aliyuncs.com", "path": "/data"}}
	expanded := expandPolarDBPolarFsBucketPaths(input)
	wantExpanded := []map[string]interface{}{{"Bucket": "example.oss-cn-beijing-internal.aliyuncs.com", "Path": "/data"}}
	if !reflect.DeepEqual(expanded, wantExpanded) {
		t.Fatalf("unexpected expanded bucket paths: %#v", expanded)
	}
	raw := []interface{}{map[string]interface{}{"Bucket": "example.oss-cn-beijing-internal.aliyuncs.com", "Path": "/data"}}
	wantFlattened := []map[string]interface{}{{"bucket": "example.oss-cn-beijing-internal.aliyuncs.com", "path": "/data"}}
	if got := flattenPolarDBPolarFsBucketPaths(raw); !reflect.DeepEqual(got, wantFlattened) {
		t.Fatalf("unexpected flattened bucket paths: %#v", got)
	}
}

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

func TestUnitPolarDBPolarFsSensitiveFields(t *testing.T) {
	resourceSchema := resourceAlicloudPolarDBPolarFs()
	for _, field := range []string{"custom_oss_ak", "custom_oss_sk"} {
		if !resourceSchema.Schema[field].Sensitive {
			t.Fatalf("%s must be sensitive", field)
		}
	}
	if resourceSchema.Schema["custom_bucket_path_list"].Type != schema.TypeSet {
		t.Fatal("custom_bucket_path_list must be an unordered set")
	}
	validateDBType := resourceSchema.Schema["db_type"].ValidateFunc
	for _, value := range []string{"polardb_mysql", "polardb_pg", "polar_agentic_db"} {
		if _, errors := validateDBType(value, "db_type"); len(errors) != 0 {
			t.Fatalf("db_type %q must be accepted: %v", value, errors)
		}
	}
	for _, value := range []string{"MySQL", "PostgreSQL"} {
		if _, errors := validateDBType(value, "db_type"); len(errors) == 0 {
			t.Fatalf("display value %q must not be sent as the OpenAPI DBType", value)
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

func TestAccAliCloudPolarDBPolarFs_planOnly(t *testing.T) {
	resourceId := "alicloud_polardb_polar_fs.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, "polarfs", func(string) string { return "" })
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: testAccConfig(map[string]interface{}{
					"storage_type":            "essdpl1",
					"authorized_user_ids":     []string{"1234567890"},
					"db_type":                 "polardb_pg",
					"vpc_id":                  "vpc-example",
					"vswitch_id":              "vsw-example",
					"zone_id":                 "cn-beijing-a",
					"storage_space":           100,
					"db_cluster_id":           "pc-example",
					"pay_type":                "Prepaid",
					"period":                  "Month",
					"used_time":               "1",
					"auto_renew":              false,
					"accelerate_switch":       "ON",
					"accelerate_storage_size": 500,
					"creation_category":       "high_performance",
					"custom_bucket_count":     1,
					"custom_bucket_path":      "/data",
					"custom_oss_ak":           "test-ak",
					"custom_oss_sk":           "test-sk",
					"auto_use_coupon":         false,
					"promotion_code":          "test-code",
					"accelerate_type":         "alluxio",
					"custom_bucket_path_list": []map[string]interface{}{{
						"bucket": "example.oss-cn-beijing-internal.aliyuncs.com",
						"path":   "/data",
					}},
				}),
				ExpectNonEmptyPlan: true,
			},
			{
				PlanOnly: true,
				Config: testAccConfig(map[string]interface{}{
					"custom_bucket_path_list": []map[string]interface{}{{
						"bucket": "example-2.oss-cn-beijing-internal.aliyuncs.com",
						"path":   "/data-2",
					}},
				}),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
