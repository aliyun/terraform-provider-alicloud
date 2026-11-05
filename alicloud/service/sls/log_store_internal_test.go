package sls

import (
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitBuildMetricStore(t *testing.T) {
	tests := []struct {
		name     string
		raw      map[string]interface{}
		validate func(t *testing.T, ms *sls.MetricStore)
	}{
		{
			name: "all_fields_set",
			raw: map[string]interface{}{
				"name":                  "test-metric-store",
				"retention_period":      30,
				"shard_count":           2,
				"enable_web_tracking":   true,
				"auto_split":            true,
				"max_split_shard_count": 64,
				"append_meta":           true,
				"hot_ttl":               10,
				"infrequent_access_ttl": 20,
				"mode":                  "standard",
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.Name != "test-metric-store" {
					t.Errorf("expected Name %q, got %q", "test-metric-store", ms.Name)
				}
				if ms.TTL != 30 {
					t.Errorf("expected TTL %d, got %d", 30, ms.TTL)
				}
				if ms.ShardCount != 2 {
					t.Errorf("expected ShardCount %d, got %d", 2, ms.ShardCount)
				}
				if !ms.WebTracking {
					t.Errorf("expected WebTracking true, got false")
				}
				if !ms.AutoSplit {
					t.Errorf("expected AutoSplit true, got false")
				}
				if ms.MaxSplitShard != 64 {
					t.Errorf("expected MaxSplitShard %d, got %d", 64, ms.MaxSplitShard)
				}
				if !ms.AppendMeta {
					t.Errorf("expected AppendMeta true, got false")
				}
				if ms.HotTTL != 10 {
					t.Errorf("expected HotTTL %d, got %d", 10, ms.HotTTL)
				}
				if ms.InfrequentAccessTTL == nil || *ms.InfrequentAccessTTL != 20 {
					t.Errorf("expected InfrequentAccessTTL %d, got %v", 20, ms.InfrequentAccessTTL)
				}
				if ms.Mode != "standard" {
					t.Errorf("expected Mode %q, got %q", "standard", ms.Mode)
				}
			},
		},
		{
			name: "name_field_overrides_logstore_name",
			raw: map[string]interface{}{
				"name":             "metric-store-name",
				"logstore_name":    "logstore-name",
				"retention_period": 30,
				"shard_count":      2,
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.Name != "metric-store-name" {
					t.Errorf("expected Name %q, got %q", "metric-store-name", ms.Name)
				}
			},
		},
		{
			name: "no_optional_hot_ttl_and_infrequent_access_ttl",
			raw: map[string]interface{}{
				"name":             "test-metric-store",
				"retention_period": 30,
				"shard_count":      2,
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.HotTTL != 0 {
					t.Errorf("expected HotTTL %d, got %d", 0, ms.HotTTL)
				}
				if ms.InfrequentAccessTTL != nil {
					t.Errorf("expected InfrequentAccessTTL nil, got %v", *ms.InfrequentAccessTTL)
				}
			},
		},
		{
			name: "hot_ttl_set_without_infrequent_access_ttl",
			raw: map[string]interface{}{
				"name":             "test-metric-store",
				"retention_period": 30,
				"shard_count":      2,
				"hot_ttl":          10,
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.HotTTL != 10 {
					t.Errorf("expected HotTTL %d, got %d", 10, ms.HotTTL)
				}
				if ms.InfrequentAccessTTL != nil {
					t.Errorf("expected InfrequentAccessTTL nil, got %v", *ms.InfrequentAccessTTL)
				}
			},
		},
		{
			name: "infrequent_access_ttl_set_without_hot_ttl",
			raw: map[string]interface{}{
				"name":                  "test-metric-store",
				"retention_period":      30,
				"shard_count":           2,
				"infrequent_access_ttl": 20,
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.HotTTL != 0 {
					t.Errorf("expected HotTTL %d, got %d", 0, ms.HotTTL)
				}
				if ms.InfrequentAccessTTL == nil || *ms.InfrequentAccessTTL != 20 {
					t.Errorf("expected InfrequentAccessTTL %d, got %v", 20, ms.InfrequentAccessTTL)
				}
			},
		},
		{
			name: "encrypt_conf_with_user_cmk_info",
			raw: map[string]interface{}{
				"name":             "test-metric-store",
				"retention_period": 30,
				"shard_count":      2,
				"encrypt_conf": []interface{}{map[string]interface{}{
					"enable":       true,
					"encrypt_type": "default",
					"user_cmk_info": []interface{}{map[string]interface{}{
						"cmk_key_id": "key-abc",
						"arn":        "acs:kms:cn-shanghai:1234567890:key/key-abc",
						"region_id":  "cn-shanghai",
					}},
				}},
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.EncryptConf == nil {
					t.Fatalf("expected EncryptConf to be set, got nil")
				}
				if !ms.EncryptConf.Enable {
					t.Errorf("expected EncryptConf.Enable true, got false")
				}
				if ms.EncryptConf.EncryptType != "default" {
					t.Errorf("expected EncryptType %q, got %q", "default", ms.EncryptConf.EncryptType)
				}
				cmk := ms.EncryptConf.UserCmkInfo
				if cmk == nil {
					t.Fatalf("expected UserCmkInfo to be set, got nil")
				}
				if cmk.CmkKeyId != "key-abc" {
					t.Errorf("expected CmkKeyId %q, got %q", "key-abc", cmk.CmkKeyId)
				}
				if cmk.Arn != "acs:kms:cn-shanghai:1234567890:key/key-abc" {
					t.Errorf("expected Arn %q, got %q", "acs:kms:cn-shanghai:1234567890:key/key-abc", cmk.Arn)
				}
				if cmk.RegionId != "cn-shanghai" {
					t.Errorf("expected RegionId %q, got %q", "cn-shanghai", cmk.RegionId)
				}
			},
		},
		{
			name: "encrypt_conf_without_user_cmk_info",
			raw: map[string]interface{}{
				"name":             "test-metric-store",
				"retention_period": 30,
				"shard_count":      2,
				"encrypt_conf": []interface{}{map[string]interface{}{
					"enable":        true,
					"encrypt_type":  "default",
					"user_cmk_info": []interface{}{},
				}},
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.EncryptConf == nil {
					t.Fatalf("expected EncryptConf to be set, got nil")
				}
				if !ms.EncryptConf.Enable {
					t.Errorf("expected EncryptConf.Enable true, got false")
				}
				if ms.EncryptConf.UserCmkInfo != nil {
					t.Errorf("expected UserCmkInfo nil, got %v", ms.EncryptConf.UserCmkInfo)
				}
			},
		},
		{
			name: "no_encrypt_conf",
			raw: map[string]interface{}{
				"name":             "test-metric-store",
				"retention_period": 30,
				"shard_count":      2,
			},
			validate: func(t *testing.T, ms *sls.MetricStore) {
				if ms.EncryptConf != nil {
					t.Errorf("expected EncryptConf nil, got %v", ms.EncryptConf)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAliCloudSlsLogStore().Schema, tt.raw)
			tt.validate(t, buildMetricStore(d))
		})
	}
}
