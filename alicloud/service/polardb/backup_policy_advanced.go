package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var polarDBAdvancedPolicyRequestFields = map[string]string{
	"retention_type":             "RetentionType",
	"retention_value":            "RetentionValue",
	"filter_type":                "FilterType",
	"filter_key":                 "FilterKey",
	"filter_value":               "FilterValue",
	"src_type":                   "SrcType",
	"src_region":                 "SrcRegion",
	"dest_type":                  "DestType",
	"dest_region":                "DestRegion",
	"dump_action":                "DumpAction",
	"only_preserve_one_each_day": "OnlyPreserveOneEachDay",
	"storage_class":              "StorageClass",
}

func polarDBAdvancedPolicyBoolConfigured(value cty.Value) bool {
	return value != cty.NilVal && value.IsKnown() && !value.IsNull()
}

func polarDBAdvancedPolicySchema() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"policy_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"retention_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validators.StringInSlice([]string{"never", "delay"}, false),
		},
		"retention_value": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validators.StringInSlice([]string{"crontab", "event"}, false),
		},
		"filter_key": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validators.StringInSlice([]string{"dayOfWeek", "dayOfMonth", "dayOfYear", "backupInterval"}, false),
		},
		"filter_value": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"src_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validators.StringInSlice([]string{"db", "level1", "level2", "level2Cross"}, false),
		},
		"src_region": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"dest_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validators.StringInSlice([]string{"level1", "level2", "level2Cross"}, false),
		},
		"dest_region": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"dump_action": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validators.StringInSlice([]string{"copy"}, false),
		},
		"only_preserve_one_each_day": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"storage_class": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"auto_created": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"only_preserve_one_each_hour": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"bak_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}}
}

func polarDBObservedAdvancedPolicySchema() *schema.Resource {
	managed := polarDBAdvancedPolicySchema()
	fields := make(map[string]*schema.Schema, len(managed.Schema))
	for key, field := range managed.Schema {
		fields[key] = &schema.Schema{Type: field.Type, Computed: true}
	}
	return &schema.Resource{Schema: fields}
}

func validatePolarDBAdvancedBackupPolicyDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	oldLevel, newLevel := d.GetChange("backup_policy_level")
	if oldLevel == "Advanced" && newLevel == "Normal" {
		return fmt.Errorf("backup_policy_level cannot be changed from Advanced back to Normal")
	}
	policies := polarDBAdvancedPolicyList(d.Get("advanced_data_policies"))
	if len(policies) > 0 && newLevel != "Advanced" {
		return fmt.Errorf("advanced_data_policies requires backup_policy_level = Advanced")
	}
	if newLevel == "Advanced" {
		rawConfig := d.GetRawConfig()
		for _, key := range []string{
			"preferred_backup_period",
			"preferred_backup_time",
			"backup_retention_period",
			"data_level1_backup_retention_period",
			"data_level2_backup_retention_period",
			"backup_frequency",
			"data_level1_backup_frequency",
			"data_level1_backup_time",
			"data_level1_backup_period",
			"data_level2_backup_period",
			"data_level2_backup_another_region_region",
			"data_level2_backup_another_region_retention_period",
			"log_backup_retention_period",
			"log_backup_another_region_region",
			"log_backup_another_region_retention_period",
		} {
			if (rawConfig.IsNull() || !rawConfig.GetAttr(key).IsNull()) && d.HasChange(key) {
				return fmt.Errorf("%s cannot be changed when backup_policy_level is Advanced; use advanced_data_policies", key)
			}
		}
	}
	for index, policy := range policies {
		if policy["retention_type"] == "delay" && strings.TrimSpace(fmt.Sprint(policy["retention_value"])) == "" {
			return fmt.Errorf("advanced_data_policies[%d].retention_value is required for delay retention", index)
		}
		if policy["filter_type"] == "crontab" && strings.TrimSpace(fmt.Sprint(policy["filter_key"])) == "" {
			return fmt.Errorf("advanced_data_policies[%d].filter_key is required for crontab filters", index)
		}
		if polarDBAdvancedPolicyIsCrossRegion(policy) && (polarDBAdvancedPolicyString(policy["src_region"]) == "" || polarDBAdvancedPolicyString(policy["dest_region"]) == "") {
			return fmt.Errorf("advanced_data_policies[%d] requires explicit src_region and dest_region for cross-region rules", index)
		}
	}
	oldValue, newValue := d.GetChange("advanced_data_policies")
	previous := polarDBAdvancedPolicyList(oldValue)
	for _, desired := range polarDBAdvancedPolicyList(newValue) {
		for _, old := range previous {
			if polarDBAdvancedPolicyString(desired["policy_id"]) != "" && desired["policy_id"] != old["policy_id"] {
				continue
			}
			if polarDBAdvancedPolicyString(desired["policy_id"]) == "" && polarDBAdvancedPolicyIdentity(desired) != polarDBAdvancedPolicyIdentity(old) {
				continue
			}
			for _, key := range []string{
				"retention_value",
				"filter_key",
				"filter_value",
				"src_region",
				"dest_region",
				"dump_action",
				"storage_class",
			} {
				if polarDBAdvancedPolicyString(old[key]) != "" && polarDBAdvancedPolicyString(desired[key]) == "" {
					return fmt.Errorf("advanced_data_policies.%s cannot be cleared by setting an empty value; remove the policy block to delete the rule", key)
				}
			}
			break
		}
	}
	return nil
}

func polarDBAdvancedPolicyList(value interface{}) []map[string]interface{} {
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if policy, ok := item.(map[string]interface{}); ok {
			result = append(result, policy)
		}
	}
	return result
}

func polarDBAdvancedPolicyIsCrossRegion(policy map[string]interface{}) bool {
	return policy["src_type"] == "level2Cross" || policy["dest_type"] == "level2Cross"
}

func preparePolarDBAdvancedPolicyRegions(policies []map[string]interface{}, region string, requireCrossRegions bool) error {
	for index, policy := range policies {
		if polarDBAdvancedPolicyIsCrossRegion(policy) {
			if requireCrossRegions && (polarDBAdvancedPolicyString(policy["src_region"]) == "" || polarDBAdvancedPolicyString(policy["dest_region"]) == "") {
				return fmt.Errorf("advanced_data_policies[%d] requires explicit src_region and dest_region for cross-region rules", index)
			}
			continue
		}
		if polarDBAdvancedPolicyString(policy["src_region"]) == "" {
			policy["src_region"] = region
		}
		if polarDBAdvancedPolicyString(policy["dest_region"]) == "" {
			policy["dest_region"] = region
		}
	}
	return nil
}

func polarDBAdvancedPolicyIdentity(policy map[string]interface{}) string {
	return strings.Join([]string{
		polarDBAdvancedPolicyString(policy["src_type"]),
		polarDBAdvancedPolicyString(policy["src_region"]),
		polarDBAdvancedPolicyString(policy["dest_type"]),
		polarDBAdvancedPolicyString(policy["dest_region"]),
		polarDBAdvancedPolicyString(policy["filter_type"]),
		polarDBAdvancedPolicyString(policy["filter_key"]),
	}, "|")
}

type polarDBAdvancedPolicyAction struct {
	Action string
	Policy map[string]interface{}
}

func polarDBAdvancedPolicyActions(previous, desired []map[string]interface{}) []polarDBAdvancedPolicyAction {
	used := make([]bool, len(previous))
	actions := make([]polarDBAdvancedPolicyAction, 0)
	for _, policy := range desired {
		match := matchPolarDBAdvancedPolicy(previous, used, policy)
		policyID := strings.TrimSpace(fmt.Sprint(policy["policy_id"]))
		if match < 0 {
			if policyID != "" && policyID != "<nil>" {
				actions = append(actions, polarDBAdvancedPolicyAction{Action: "UPDATE", Policy: policy})
			} else {
				actions = append(actions, polarDBAdvancedPolicyAction{Action: "CREATE", Policy: policy})
			}
			continue
		}
		used[match] = true
		if !polarDBAdvancedPolicyEqual(previous[match], policy) {
			updated := make(map[string]interface{}, len(policy)+1)
			for key, value := range policy {
				updated[key] = value
			}
			updated["policy_id"] = previous[match]["policy_id"]
			actions = append(actions, polarDBAdvancedPolicyAction{Action: "UPDATE", Policy: updated})
		}
	}
	for i, old := range previous {
		if !used[i] {
			actions = append(actions, polarDBAdvancedPolicyAction{Action: "DELETE", Policy: old})
		}
	}
	return actions
}

func matchPolarDBAdvancedPolicy(previous []map[string]interface{}, used []bool, desired map[string]interface{}) int {
	policyID := strings.TrimSpace(fmt.Sprint(desired["policy_id"]))
	if policyID != "" && policyID != "<nil>" {
		for i, old := range previous {
			if !used[i] &&
				strings.TrimSpace(fmt.Sprint(old["policy_id"])) == policyID &&
				polarDBAdvancedPolicyIdentity(old) == polarDBAdvancedPolicyIdentity(desired) {
				return i
			}
		}
	}
	for i, old := range previous {
		if !used[i] && polarDBAdvancedPolicyIdentity(old) == polarDBAdvancedPolicyIdentity(desired) {
			return i
		}
	}
	if policyID != "" && policyID != "<nil>" {
		for i, old := range previous {
			if !used[i] && strings.TrimSpace(fmt.Sprint(old["policy_id"])) == policyID {
				return i
			}
		}
	}
	return -1
}

func polarDBAdvancedPolicyIDsInDesiredOrder(previous, desired []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(desired))
	used := make([]bool, len(previous))
	for _, policy := range desired {
		updated := make(map[string]interface{}, len(policy))
		for key, value := range policy {
			updated[key] = value
		}
		delete(updated, "only_preserve_one_each_day_set")
		if match := matchPolarDBAdvancedPolicy(previous, used, policy); match >= 0 {
			used[match] = true
			updated["policy_id"] = previous[match]["policy_id"]
		}
		result = append(result, updated)
	}
	return result
}

func polarDBAdvancedPolicyEqual(a, b map[string]interface{}) bool {
	for key := range polarDBAdvancedPolicyRequestFields {
		if key == "only_preserve_one_each_day" && b["only_preserve_one_each_day_set"] != true {
			continue
		}
		if polarDBAdvancedPolicyString(a[key]) != polarDBAdvancedPolicyString(b[key]) {
			return false
		}
	}
	return true
}

func polarDBAdvancedPolicyApplied(observed, requested map[string]interface{}) bool {
	for key := range polarDBAdvancedPolicyRequestFields {
		if key == "only_preserve_one_each_day" && requested["only_preserve_one_each_day_set"] != true {
			continue
		}
		want := polarDBAdvancedPolicyString(requested[key])
		if want == "" && key != "only_preserve_one_each_day" {
			continue
		}
		if polarDBAdvancedPolicyString(observed[key]) != want {
			return false
		}
	}
	return true
}

func expandPolarDBAdvancedPolicyActions(request map[string]interface{}, actions []polarDBAdvancedPolicyAction) error {
	policies := make([]map[string]interface{}, 0, len(actions))
	for _, action := range actions {
		policy := map[string]interface{}{
			"ActionType": action.Action,
		}
		if action.Action != "CREATE" {
			policy["PolicyId"] = action.Policy["policy_id"]
		}
		if action.Action != "DELETE" {
			for key, apiName := range polarDBAdvancedPolicyRequestFields {
				value := action.Policy[key]
				if key == "only_preserve_one_each_day" {
					if action.Policy["only_preserve_one_each_day_set"] == true {
						policy[apiName] = value
					}
					continue
				}
				if value != nil && fmt.Sprint(value) != "" {
					policy[apiName] = value
				}
			}
		}
		policies = append(policies, policy)
	}
	encoded, err := json.Marshal(policies)
	if err != nil {
		return err
	}
	request["AdvancedDataPolicies"] = string(encoded)
	return nil
}

func callPolarDBModifyBackupPolicy(client *connectivity.AliyunClient, request map[string]interface{}, timeout time.Duration) error {
	var response map[string]interface{}
	var err error
	action := "ModifyBackupPolicy"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return err
}

func waitPolarDBAdvancedPolicyActions(client *connectivity.AliyunClient, clusterID string, actions []polarDBAdvancedPolicyAction, timeout time.Duration) error {
	service := Service{client}
	return retry.Retry(timeout, func() *retry.RetryError {
		object, err := service.DescribeBackupPolicy(clusterID)
		if err != nil {
			if helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if object["BackupPolicyLevel"] != "Advanced" {
			return retry.RetryableError(fmt.Errorf("backup policy level has not become Advanced"))
		}
		observed := polarDBDescribeAdvancedPolicies(object)
		for _, action := range actions {
			found := false
			for _, policy := range observed {
				if policy["auto_created"] == true {
					continue
				}
				if action.Action == "CREATE" && polarDBAdvancedPolicyIdentity(policy) == polarDBAdvancedPolicyIdentity(action.Policy) && policy["policy_id"] != "" && polarDBAdvancedPolicyApplied(policy, action.Policy) {
					found = true
					break
				}
				if action.Action == "UPDATE" && policy["policy_id"] == action.Policy["policy_id"] && polarDBAdvancedPolicyApplied(policy, action.Policy) {
					found = true
					break
				}
				if action.Action == "DELETE" && policy["policy_id"] == action.Policy["policy_id"] {
					found = true
					break
				}
			}
			if action.Action == "DELETE" && found {
				return retry.RetryableError(fmt.Errorf("advanced backup policy %v has not been deleted", action.Policy["policy_id"]))
			}
			if action.Action != "DELETE" && !found {
				return retry.RetryableError(fmt.Errorf("advanced backup policy %v is not yet visible", action.Policy["policy_id"]))
			}
		}
		return nil
	})
}

func polarDBDescribeAdvancedPolicies(object map[string]interface{}) []map[string]interface{} {
	wrapper, ok := object["AdvancedDataPolicies"].(map[string]interface{})
	if !ok {
		return nil
	}
	items, ok := wrapper["AdvancedDataPolicy"].([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		policy := map[string]interface{}{
			"policy_id":                   strings.TrimSpace(polarDBAdvancedPolicyString(entry["PolicyId"])),
			"retention_type":              entry["RetentionType"],
			"retention_value":             entry["RetentionValue"],
			"filter_type":                 entry["FilterType"],
			"filter_key":                  entry["FilterKey"],
			"filter_value":                entry["FilterValue"],
			"src_type":                    entry["SrcType"],
			"src_region":                  entry["SrcRegion"],
			"dest_type":                   entry["DestType"],
			"dest_region":                 entry["DestRegion"],
			"dump_action":                 entry["DumpAction"],
			"only_preserve_one_each_day":  entry["OnlyPreserveOneEachDay"],
			"storage_class":               entry["StorageClass"],
			"auto_created":                entry["AutoCreated"],
			"only_preserve_one_each_hour": entry["OnlyPreserveOneEachHour"],
			"bak_type":                    entry["BakType"],
		}
		result = append(result, policy)
	}
	sort.Slice(result, func(i, j int) bool {
		return fmt.Sprint(result[i]["policy_id"]) < fmt.Sprint(result[j]["policy_id"])
	})
	return result
}

func polarDBAdvancedPolicyString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func reconcilePolarDBAdvancedPolicies(previous []map[string]interface{}, observed []map[string]interface{}) ([]map[string]interface{}, []map[string]interface{}) {
	managed := make([]map[string]interface{}, 0, len(previous))
	autoCreated := make([]map[string]interface{}, 0)
	used := make([]bool, len(observed))
	for _, policy := range observed {
		if policy["auto_created"] == true {
			autoCreated = append(autoCreated, policy)
		}
	}
	for _, old := range previous {
		match := -1
		oldID := strings.TrimSpace(fmt.Sprint(old["policy_id"]))
		for index, policy := range observed {
			if used[index] || policy["auto_created"] == true {
				continue
			}
			if oldID != "" && oldID != "<nil>" && oldID == policy["policy_id"] {
				match = index
				break
			}
		}
		if match < 0 && (oldID == "" || oldID == "<nil>") {
			for index, policy := range observed {
				if used[index] || policy["auto_created"] == true {
					continue
				}
				if polarDBAdvancedPolicyIdentity(old) == polarDBAdvancedPolicyIdentity(policy) && polarDBAdvancedPolicyApplied(policy, old) {
					match = index
					break
				}
			}
		}
		if match < 0 && (oldID == "" || oldID == "<nil>") {
			for index, policy := range observed {
				if used[index] || policy["auto_created"] == true {
					continue
				}
				if polarDBAdvancedPolicyIdentity(old) != polarDBAdvancedPolicyIdentity(policy) {
					continue
				}
				match = index
				break
			}
		}
		if match >= 0 {
			used[match] = true
			managed = append(managed, observed[match])
		}
	}
	return managed, autoCreated
}
