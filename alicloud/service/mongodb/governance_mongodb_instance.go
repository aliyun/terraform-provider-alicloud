package mongodb

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mongoDBInstanceGovernanceSchema covers the Dds APIs that are exposed on the
// existing alicloud_mongodb_instance type instead of a dedicated Terraform
// resource. All capabilities are optional so existing configurations keep
// working unchanged.
func mongoDBInstanceGovernanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"databases": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"sync_source": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role_id":        {Type: schema.TypeString, Computed: true},
					"role":           {Type: schema.TypeString, Computed: true},
					"target_role_id": {Type: schema.TypeString, Computed: true},
					"target_role":    {Type: schema.TypeString, Computed: true},
					"ip":             {Type: schema.TypeString, Computed: true},
					"port":           {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"ai_assistant": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"db_instance_id":  {Type: schema.TypeString, Computed: true},
					"charge_type":     {Type: schema.TypeString, Computed: true},
					"status":          {Type: schema.TypeString, Computed: true},
					"base_quota":      {Type: schema.TypeInt, Computed: true},
					"extra_quota":     {Type: schema.TypeInt, Computed: true},
					"total_quota":     {Type: schema.TypeInt, Computed: true},
					"start_time":      {Type: schema.TypeString, Computed: true},
					"expire_time":     {Type: schema.TypeString, Computed: true},
					"product_edition": {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"cold_data_tables": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"db_name":    {Type: schema.TypeString, Computed: true},
					"table_name": {Type: schema.TypeString, Computed: true},
					"table_size": {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"hot_data_tables": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"db_name":    {Type: schema.TypeString, Computed: true},
					"table_name": {Type: schema.TypeString, Computed: true},
					"table_size": {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"table_recovery_support": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"database_lazy_load_support": {Type: schema.TypeBool, Computed: true},
					"db_tables_support":          {Type: schema.TypeBool, Computed: true},
					"instant_db_tables_support":  {Type: schema.TypeBool, Computed: true},
				},
			},
		},
		"advanced_backup_policy": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"preferred_backup_time": {Type: schema.TypeString, Optional: true, Computed: true},
					"snapshot_backup_type": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validators.StringInSlice([]string{"Standard", "Flash"}, false),
					},
					"backup_retention_policy_on_cluster_deletion": {Type: schema.TypeInt, Optional: true, Computed: true},
					"backup_tamper_proof_enabled":                 {Type: schema.TypeBool, Optional: true, Computed: true},
					"advanced_backup_data_policies":               {Type: schema.TypeString, Optional: true},
					"advanced_backup_log_policies":                {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"advanced_backup_status": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"archive_operation": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validators.StringInSlice([]string{"enable", "disable"}, false),
		},
		"archived_tables": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"restored_tables": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"restore_category": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"table_restore": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"table_meta":   {Type: schema.TypeString, Required: true},
					"restore_type": {Type: schema.TypeString, Optional: true, Default: "RestoreTime"},
					"restore_time": {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"instance_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config_name":  {Type: schema.TypeString, Required: true},
					"config_value": {Type: schema.TypeString, Required: true},
				},
			},
		},
		"transform_to_ecs": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"target_class":             {Type: schema.TypeString, Required: true},
					"coupon_no":                {Type: schema.TypeString, Optional: true},
					"auto_pay":                 {Type: schema.TypeBool, Optional: true},
					"storage":                  {Type: schema.TypeInt, Optional: true},
					"order_type":               {Type: schema.TypeString, Optional: true},
					"db_instance_storage_type": {Type: schema.TypeString, Optional: true},
					"effective_time":           {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"transform_order_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ai_gateway": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"db_instance_class": {Type: schema.TypeString, Optional: true},
					"effective_time":    {Type: schema.TypeInt, Optional: true},
				},
			},
		},
		"ai_gateway_instance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ai_gateway_task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ai_gateway_api_key_reset": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ai_gateway_api_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"replication_link": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"source_instance_id":        {Type: schema.TypeString, Required: true},
					"source_instance_region_id": {Type: schema.TypeString, Required: true},
				},
			},
		},
		"ai_performance_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ai_performance_start_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ai_performance_end_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ai_performance_interval": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ai_performance_values": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {Type: schema.TypeString, Computed: true},
					"date":  {Type: schema.TypeString, Computed: true},
				},
			},
		},
	}
}

func mongoDBGovernanceRpcPost(client *connectivity.AliyunClient, timeout time.Duration, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
	return response, err
}

// mongoDBGovernancePick looks a response key up at the top level and in the
// gateway's "data" envelope, because Dds APIs mix both shapes.
func mongoDBGovernancePick(response map[string]interface{}, keys ...string) interface{} {
	for _, key := range keys {
		if v, ok := response[key]; ok && v != nil {
			return v
		}
	}
	if data, ok := response["data"].(map[string]interface{}); ok {
		for _, key := range keys {
			if v, ok := data[key]; ok && v != nil {
				return v
			}
		}
	}
	return nil
}

func mongoDBGovernanceStringList(v interface{}) []string {
	out := make([]string, 0)
	list, ok := v.([]interface{})
	if !ok {
		return out
	}
	for _, item := range list {
		if s := fmt.Sprint(item); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func mongoDBGovernanceMapList(v interface{}) []map[string]interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return make([]map[string]interface{}, 0)
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func mongoDBGovernanceInt(v interface{}) int {
	n, _ := strconv.Atoi(fmt.Sprint(v))
	return n
}

func mongoDBGovernanceReadInstance(d *schema.ResourceData, client *connectivity.AliyunClient) {
	id := d.Id()

	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeDatabases", map[string]interface{}{"DBInstanceId": id}); err != nil {
		log.Printf("[DEBUG] DescribeDatabases for %s: %s", id, err)
	} else {
		d.Set("databases", mongoDBGovernanceStringList(mongoDBGovernancePick(response, "Databases", "databases")))
	}

	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeDBInstanceSyncSource", map[string]interface{}{"DBInstanceId": id}); err != nil {
		log.Printf("[DEBUG] DescribeDBInstanceSyncSource for %s: %s", id, err)
	} else {
		syncSources := make([]map[string]interface{}, 0)
		for _, item := range mongoDBGovernanceMapList(mongoDBGovernancePick(response, "ReplicaSetSyncSource", "replicaSetSyncSource")) {
			syncSources = append(syncSources, map[string]interface{}{
				"role_id":        item["RoleId"],
				"role":           item["Role"],
				"target_role_id": item["TargetRoleId"],
				"target_role":    item["TargetRole"],
				"ip":             item["Ip"],
				"port":           item["Port"],
			})
		}
		d.Set("sync_source", syncSources)
	}

	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeAIAssistant", map[string]interface{}{}); err != nil {
		log.Printf("[DEBUG] DescribeAIAssistant for %s: %s", id, err)
	} else if assistantId := fmt.Sprint(mongoDBGovernancePick(response, "DBInstanceId", "dbInstanceId")); assistantId == id {
		d.Set("ai_assistant", []map[string]interface{}{{
			"db_instance_id":  assistantId,
			"charge_type":     mongoDBGovernancePick(response, "ChargeType", "chargeType"),
			"status":          mongoDBGovernancePick(response, "Status", "status"),
			"base_quota":      mongoDBGovernanceInt(mongoDBGovernancePick(response, "BaseQuota", "baseQuota")),
			"extra_quota":     mongoDBGovernanceInt(mongoDBGovernancePick(response, "ExtraQuota", "extraQuota")),
			"total_quota":     mongoDBGovernanceInt(mongoDBGovernancePick(response, "TotalQuota", "totalQuota")),
			"start_time":      mongoDBGovernancePick(response, "StartTime", "startTime"),
			"expire_time":     mongoDBGovernancePick(response, "ExpireTime", "expireTime"),
			"product_edition": mongoDBGovernancePick(response, "ProductEdition", "productEdition"),
		}})
	}

	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "SupportDBTableRecovery", map[string]interface{}{"DBInstanceId": id}); err != nil {
		log.Printf("[DEBUG] SupportDBTableRecovery for %s: %s", id, err)
	} else {
		d.Set("table_recovery_support", []map[string]interface{}{{
			"database_lazy_load_support": mongoDBGovernancePick(response, "DatabaseLazyLoadSupport", "databaseLazyLoadSupport"),
			"db_tables_support":          mongoDBGovernancePick(response, "DbTablesSupport", "dbTablesSupport"),
			"instant_db_tables_support":  mongoDBGovernancePick(response, "InstantDBTablesSupport", "instantDBTablesSupport"),
		}})
	}

	coldRequest := map[string]interface{}{"DBInstanceId": id, "PageSize": 100, "PageNumber": 1}
	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeDBInstanceColdData", coldRequest); err != nil {
		log.Printf("[DEBUG] DescribeDBInstanceColdData for %s: %s", id, err)
	} else {
		tables := make([]map[string]interface{}, 0)
		for _, item := range mongoDBGovernanceMapList(mongoDBGovernancePick(response, "ArchivedTables", "archivedTables")) {
			tables = append(tables, map[string]interface{}{
				"db_name":    item["DbName"],
				"table_name": item["TableName"],
				"table_size": fmt.Sprint(item["TableSize"]),
			})
		}
		d.Set("cold_data_tables", tables)
	}

	hotRequest := map[string]interface{}{"DBInstanceId": id, "PageSize": 100, "PageNumber": 1}
	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeDBInstanceHotData", hotRequest); err != nil {
		log.Printf("[DEBUG] DescribeDBInstanceHotData for %s: %s", id, err)
	} else {
		tables := make([]map[string]interface{}, 0)
		for _, item := range mongoDBGovernanceMapList(mongoDBGovernancePick(response, "DbTables", "dbTables")) {
			tables = append(tables, map[string]interface{}{
				"db_name":    item["DbName"],
				"table_name": item["TableName"],
				"table_size": fmt.Sprint(item["TableSize"]),
			})
		}
		d.Set("hot_data_tables", tables)
	}

	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeAdvancedBackupPolicy", map[string]interface{}{"DBInstanceId": id}); err != nil {
		log.Printf("[DEBUG] DescribeAdvancedBackupPolicy for %s: %s", id, err)
	} else {
		d.Set("advanced_backup_policy", []map[string]interface{}{{
			"preferred_backup_time":                       mongoDBGovernancePick(response, "PreferredBackupTime", "preferredBackupTime"),
			"snapshot_backup_type":                        mongoDBGovernancePick(response, "SnapshotBackupType", "snapshotBackupType"),
			"backup_retention_policy_on_cluster_deletion": mongoDBGovernanceInt(mongoDBGovernancePick(response, "BackupRetentionPolicyOnClusterDeletion", "backupRetentionPolicyOnClusterDeletion")),
			"backup_tamper_proof_enabled":                 mongoDBGovernancePick(response, "BackupTamperProofEnabled", "backupTamperProofEnabled"),
		}})
	}

	key := d.Get("ai_performance_key").(string)
	startTime := d.Get("ai_performance_start_time").(string)
	endTime := d.Get("ai_performance_end_time").(string)
	if key != "" && startTime != "" && endTime != "" {
		request := map[string]interface{}{
			"DBInstanceId": id,
			"Key":          key,
			"StartTime":    startTime,
			"EndTime":      endTime,
		}
		if v, ok := d.GetOk("ai_performance_interval"); ok {
			request["Interval"] = v
		}
		if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeDBInstanceAiPerformance", request); err != nil {
			log.Printf("[DEBUG] DescribeDBInstanceAiPerformance for %s: %s", id, err)
		} else {
			values := make([]map[string]interface{}, 0)
			for _, keyItem := range mongoDBGovernanceMapList(mongoDBGovernancePick(response, "PerformanceKeys", "performanceKeys")) {
				for _, valueItem := range mongoDBGovernanceMapList(mongoDBGovernancePick(keyItem, "PerformanceValues", "Values", "values")) {
					values = append(values, map[string]interface{}{
						"value": fmt.Sprint(valueItem["Value"]),
						"date":  fmt.Sprint(valueItem["Date"]),
					})
				}
			}
			d.Set("ai_performance_values", values)
		}
	}
}

func mongoDBGovernanceUpdateInstance(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	id := d.Id()

	if d.HasChange("ai_gateway") {
		oldRaw, newRaw := d.GetChange("ai_gateway")
		oldList, _ := oldRaw.([]interface{})
		newList, _ := newRaw.([]interface{})
		oldClass, newClass := "", ""
		if len(oldList) > 0 {
			if oldMap, ok := oldList[0].(map[string]interface{}); ok {
				oldClass, _ = oldMap["db_instance_class"].(string)
			}
		}
		if len(newList) > 0 {
			if newMap, ok := newList[0].(map[string]interface{}); ok {
				newClass, _ = newMap["db_instance_class"].(string)
			}
		}
		// effective_time is an action parameter, not gateway state, so only a
		// class change or block removal/add re-creates the gateway.
		if oldClass != newClass {
			if len(oldList) > 0 {
				request := map[string]interface{}{"DBInstanceId": id}
				if oldMap, ok := oldList[0].(map[string]interface{}); ok {
					if v, ok := oldMap["effective_time"].(int); ok && v > 0 {
						request["EffectiveTime"] = v
					}
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "DeleteAIGatewayInstance", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteAIGatewayInstance", errs.AlibabaCloudSdkGoERROR)
				}
			}
			if len(newList) > 0 {
				request := map[string]interface{}{"DBInstanceId": id}
				if newMap, ok := newList[0].(map[string]interface{}); ok {
					if v, ok := newMap["db_instance_class"].(string); ok && v != "" {
						request["DbInstanceClass"] = v
					}
					if v, ok := newMap["effective_time"].(int); ok && v > 0 {
						request["EffectiveTime"] = v
					}
				}
				response, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "CreateAIGatewayInstance", request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "CreateAIGatewayInstance", errs.AlibabaCloudSdkGoERROR)
				}
				d.Set("ai_gateway_instance_name", mongoDBGovernancePick(response, "AIGatewayInstanceName", "DBInstanceName", "dbInstanceName"))
				d.Set("ai_gateway_task_id", fmt.Sprint(mongoDBGovernancePick(response, "TaskId", "taskId")))
			}
		}
	}

	if d.HasChange("ai_gateway_api_key_reset") && d.Get("ai_gateway_api_key_reset").(bool) {
		response, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ResetAIGatewayApiKey", map[string]interface{}{"DBInstanceId": id})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ResetAIGatewayApiKey", errs.AlibabaCloudSdkGoERROR)
		}
		d.Set("ai_gateway_api_key", mongoDBGovernancePick(response, "ApiKey", "apiKey"))
	}

	if d.HasChange("replication_link") {
		oldRaw, newRaw := d.GetChange("replication_link")
		oldList, _ := oldRaw.([]interface{})
		newList, _ := newRaw.([]interface{})
		if len(oldList) > 0 {
			if oldMap, ok := oldList[0].(map[string]interface{}); ok {
				request := map[string]interface{}{
					"DBInstanceId":           id,
					"SourceInstanceId":       oldMap["source_instance_id"],
					"SourceInstanceRegionId": oldMap["source_instance_region_id"],
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "DeleteReplicationLink", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteReplicationLink", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
		if len(newList) > 0 {
			if newMap, ok := newList[0].(map[string]interface{}); ok {
				request := map[string]interface{}{
					"DBInstanceId":           id,
					"SourceInstanceId":       newMap["source_instance_id"],
					"SourceInstanceRegionId": newMap["source_instance_region_id"],
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "CreateReplicationLink", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "CreateReplicationLink", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	if d.HasChange("advanced_backup_policy") {
		policies, _ := d.Get("advanced_backup_policy").([]interface{})
		if len(policies) > 0 {
			if policy, ok := policies[0].(map[string]interface{}); ok {
				request := map[string]interface{}{"DBInstanceId": id}
				if v, ok := policy["backup_retention_policy_on_cluster_deletion"].(int); ok {
					request["BackupRetentionPolicyOnClusterDeletion"] = v
				}
				if v, ok := policy["preferred_backup_time"].(string); ok && v != "" {
					request["PreferredBackupTime"] = v
				}
				if v, ok := policy["snapshot_backup_type"].(string); ok && v != "" {
					request["SnapshotBackupType"] = v
				}
				if v, ok := policy["backup_tamper_proof_enabled"].(bool); ok && v {
					request["BackupTamperProofEnabled"] = v
				}
				if v, ok := policy["advanced_backup_data_policies"].(string); ok && v != "" {
					request["AdvancedBackupPolicies"] = v
				}
				if v, ok := policy["advanced_backup_log_policies"].(string); ok && v != "" {
					request["AdvancedBackupLogPolicies"] = v
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyAdvancedBackupPolicy", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyAdvancedBackupPolicy", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	if d.HasChange("advanced_backup_status") && d.Get("advanced_backup_status").(bool) {
		if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyAdvancedBackupStatus", map[string]interface{}{"DBInstanceId": id}); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyAdvancedBackupStatus", errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("archive_operation") {
		if v, ok := d.GetOk("archive_operation"); ok {
			request := map[string]interface{}{"DBInstanceId": id, "Operation": v}
			if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyDBInstanceArchive", request); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyDBInstanceArchive", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	if d.HasChange("archived_tables") || d.HasChange("restored_tables") {
		request := map[string]interface{}{"DBInstanceId": id}
		if v, ok := d.GetOk("archived_tables"); ok {
			request["ArchivedTables"] = v
		}
		if v, ok := d.GetOk("restored_tables"); ok {
			request["RestoredTables"] = v
		}
		if len(request) > 1 {
			if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyDBInstanceArchivedTables", request); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyDBInstanceArchivedTables", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	if d.HasChange("restore_category") {
		if v, ok := d.GetOk("restore_category"); ok {
			request := map[string]interface{}{"DBInstanceId": id, "RestoreCategory": v}
			if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyDBTablesRecoveryState", request); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyDBTablesRecoveryState", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	if d.HasChange("table_restore") {
		restores, _ := d.Get("table_restore").([]interface{})
		if len(restores) > 0 {
			if restore, ok := restores[0].(map[string]interface{}); ok {
				request := map[string]interface{}{
					"DBInstanceId": id,
					"TableMeta":    restore["table_meta"],
				}
				if v, ok := restore["restore_type"].(string); ok && v != "" {
					request["RestoreType"] = v
				}
				if v, ok := restore["restore_time"].(string); ok && v != "" {
					request["RestoreTime"] = v
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "RestoreTable", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "RestoreTable", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	if d.HasChange("instance_config") {
		configs, _ := d.Get("instance_config").([]interface{})
		if len(configs) > 0 {
			if config, ok := configs[0].(map[string]interface{}); ok {
				request := map[string]interface{}{
					"DBInstanceId": id,
					"ConfigName":   config["config_name"],
					"ConfigValue":  config["config_value"],
				}
				if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyDBInstanceConfig", request); err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyDBInstanceConfig", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	if d.HasChange("transform_to_ecs") {
		transforms, _ := d.Get("transform_to_ecs").([]interface{})
		if len(transforms) > 0 {
			if transform, ok := transforms[0].(map[string]interface{}); ok {
				request := map[string]interface{}{
					"DBInstanceId": id,
					"TargetClass":  transform["target_class"],
				}
				if v, ok := transform["coupon_no"].(string); ok && v != "" {
					request["CouponNo"] = v
				}
				if v, ok := transform["auto_pay"].(bool); ok && v {
					request["AutoPay"] = v
				}
				if v, ok := transform["storage"].(int); ok && v > 0 {
					request["Storage"] = v
				}
				if v, ok := transform["order_type"].(string); ok && v != "" {
					request["OrderType"] = v
				}
				if v, ok := transform["db_instance_storage_type"].(string); ok && v != "" {
					request["DBInstanceStorageType"] = v
				}
				if v, ok := transform["effective_time"].(string); ok && v != "" {
					request["EffectiveTime"] = v
				}
				response, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "TransformDBInstanceToEcs", request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "TransformDBInstanceToEcs", errs.AlibabaCloudSdkGoERROR)
				}
				d.Set("transform_order_id", fmt.Sprint(mongoDBGovernancePick(response, "OrderId", "orderId")))
			}
		}
	}

	return nil
}

func mongoDBGovernanceDeleteInstance(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	id := d.Id()
	if links, ok := d.Get("replication_link").([]interface{}); ok && len(links) > 0 {
		if link, ok := links[0].(map[string]interface{}); ok {
			request := map[string]interface{}{
				"DBInstanceId":           id,
				"SourceInstanceId":       link["source_instance_id"],
				"SourceInstanceRegionId": link["source_instance_region_id"],
			}
			if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutDelete), "DeleteReplicationLink", request); err != nil && !errs.NotFoundError(err) {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteReplicationLink", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	if gateways, ok := d.Get("ai_gateway").([]interface{}); ok && len(gateways) > 0 {
		request := map[string]interface{}{"DBInstanceId": id}
		if gateway, ok := gateways[0].(map[string]interface{}); ok {
			if v, ok := gateway["effective_time"].(int); ok && v > 0 {
				request["EffectiveTime"] = v
			}
		}
		if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutDelete), "DeleteAIGatewayInstance", request); err != nil && !errs.NotFoundError(err) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteAIGatewayInstance", errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
