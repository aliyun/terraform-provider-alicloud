package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRosStackInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosStackInstancesCreate,
		Read:   resourceAlicloudRosStackInstancesRead,
		Update: resourceAlicloudRosStackInstancesUpdate,
		Delete: resourceAlicloudRosStackInstancesDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"stack_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deployment_targets": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"account_ids"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:             schema.TypeList,
							Optional:         true,
							MaxItems:         50,
							Elem:             &schema.Schema{Type: schema.TypeString},
							Description:      "List of target account IDs for service-managed permissions.",
							DiffSuppressFunc: suppressDeploymentTargetAccountIdsDiff,
						},
						"rd_folder_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    20,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of Resource Directory folder IDs.",
						},
					},
				},
				Description: "Configuration block defining deployment targets. Conflicts with account_ids.",
			},
			"region_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				MaxItems:    20,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of target region IDs. Maximum 20 regions.",
			},
			"account_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"deployment_targets"},
				MinItems:      1,
				MaxItems:      50,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of target account IDs for self-managed permissions. Maximum 50 accounts.",
			},
			"parameter_overrides": {
				Type:      schema.TypeSet,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key":   {Type: schema.TypeString, Required: true},
						"parameter_value": {Type: schema.TypeString, Optional: true, Sensitive: true},
					},
				},
				Description: "Parameters to override in the stack instances.",
			},
			"operation_preferences": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_concurrent_count": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 20),
							ConflictsWith: []string{"operation_preferences.0.max_concurrent_percentage"},
							Description:   "Maximum number of concurrent operations per region. Range: 1-20.",
						},
						"max_concurrent_percentage": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 100),
							ConflictsWith: []string{"operation_preferences.0.max_concurrent_count"},
							Description:   "Maximum percentage of concurrent targets per region. Range: 1-100.",
						},
						"failure_tolerance_count": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(0, 20),
							ConflictsWith: []string{"operation_preferences.0.failure_tolerance_percentage"},
							Description:   "Number of failures tolerated per region. Range: 0-20.",
						},
						"failure_tolerance_percentage": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(0, 100),
							ConflictsWith: []string{"operation_preferences.0.failure_tolerance_count"},
							Description:   "Percentage of failures tolerated per region. Range: 0-100.",
						},
						"region_concurrency_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"SEQUENTIAL", "PARALLEL"}, false),
							Description:  "Concurrency type for regions. Valid values: SEQUENTIAL, PARALLEL.",
						},
					},
				},
				Description: "Preferences for how the operation is performed across multiple accounts and regions.",
			},
			"timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(1, 1440),
				Description:  "The amount of time that can elapse before the stack operation status is set to TIMED_OUT.",
			},
			"operation_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				Description:  "Description of the stack instances operation. Length: 1-256 characters.",
			},
			"disable_rollback": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "When creating stack instances fails, whether to disable the rollback policy.",
			},
			"deployment_options": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"IgnoreExisting"}, false)},
				Description: "Deployment options for service-managed permissions. Only supports 'IgnoreExisting'.",
			},
			"stack_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_operation_id":    {Type: schema.TypeString, Computed: true},
						"status":               {Type: schema.TypeString, Computed: true},
						"stack_group_id":       {Type: schema.TypeString, Computed: true},
						"stack_id":             {Type: schema.TypeString, Computed: true},
						"drift_detection_time": {Type: schema.TypeString, Computed: true},
						"stack_drift_status":   {Type: schema.TypeString, Computed: true},
						"status_reason":        {Type: schema.TypeString, Computed: true},
						"stack_group_name":     {Type: schema.TypeString, Computed: true},
						"account_id":           {Type: schema.TypeString, Computed: true},
						"region_id":            {Type: schema.TypeString, Computed: true},
						"rd_folder_id":         {Type: schema.TypeString, Computed: true},
					},
				},
				Description: "List of stack instances with their latest operation tracking ID.",
			},
		},
	}
}

func resourceAlicloudRosStackInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateStackInstances"
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": d.Get("stack_group_name").(string),
	}

	if v, ok := d.GetOk("region_ids"); ok {
		request["RegionIds"] = convertListToJsonString(v.([]interface{}))
	}

	if dtJSON := expandDeploymentTargets(d); dtJSON != "" {
		request["DeploymentTargets"] = dtJSON
	} else if v, ok := d.GetOk("account_ids"); ok {
		request["AccountIds"] = convertListToJsonString(v.([]interface{}))
	} else {
		return fmt.Errorf("either deployment_targets or account_ids must be specified")
	}

	if v, ok := d.GetOk("parameter_overrides"); ok {
		overrides := v.(*schema.Set).List()
		sort.SliceStable(overrides, func(i, j int) bool {
			return overrides[i].(map[string]interface{})["parameter_key"].(string) < overrides[j].(map[string]interface{})["parameter_key"].(string)
		})

		for i, param := range overrides {
			p := param.(map[string]interface{})
			idx := i + 1
			request[fmt.Sprintf("ParameterOverrides.%d.ParameterKey", idx)] = p["parameter_key"]
			if val, ok := p["parameter_value"]; ok && val.(string) != "" {
				request[fmt.Sprintf("ParameterOverrides.%d.ParameterValue", idx)] = val
			}
		}
	}

	if opPrefJSON := expandOperationPreferences(d); opPrefJSON != "" {
		request["OperationPreferences"] = opPrefJSON
	}

	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}

	if v, ok := d.GetOk("operation_description"); ok {
		request["OperationDescription"] = v
	}

	request["DisableRollback"] = d.Get("disable_rollback").(bool)

	if v, ok := d.GetOk("deployment_options"); ok {
		for i, option := range v.([]interface{}) {
			request[fmt.Sprintf("DeploymentOptions.%d", i+1)] = option
		}
	}

	request["ClientToken"] = buildClientToken(action)

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StackGroupOperationInProgress"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_stack_instances", action, AlibabaCloudSdkGoERROR)
	}

	operationId := fmt.Sprint(response["OperationId"])
	if operationId == "" || operationId == "<nil>" {
		return WrapError(fmt.Errorf("ROS CreateStackInstances did not return a valid OperationId"))
	}

	if err := waitForRosStackGroupOperationAndCheckResults(client, d.Get("stack_group_name").(string), operationId, d.Timeout(schema.TimeoutCreate)); err != nil {
		log.Printf("[WARN] Create operation %s ended with issues.", operationId)
		return WrapErrorf(err, "ROS CreateStackInstances operation %s failed", operationId)
	}

	d.SetId(d.Get("stack_group_name").(string))
	return resourceAlicloudRosStackInstancesRead(d, meta)
}

func resourceAlicloudRosStackInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	stackGroupName := d.Id()
	if stackGroupName == "" {
		if v, ok := d.GetOk("stack_group_name"); ok {
			stackGroupName = v.(string)
		} else {
			d.SetId("")
			return nil
		}
	}
	d.Set("stack_group_name", stackGroupName)

	// Collect targets from state for filtering
	var targetRegions, targetAccounts []string
	if v, ok := d.GetOk("region_ids"); ok {
		for _, r := range v.([]interface{}) {
			targetRegions = append(targetRegions, r.(string))
		}
	}
	if v, ok := d.GetOk("account_ids"); ok {
		for _, a := range v.([]interface{}) {
			targetAccounts = append(targetAccounts, a.(string))
		}
	}

	action := "ListStackInstances"
	var matchedInstances []map[string]interface{}
	pageNum, pageSize := 1, 50

	getStr := func(m map[string]interface{}, key string) string {
		if val, ok := m[key]; ok && val != nil {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	for {
		req := map[string]interface{}{
			"RegionId":       client.RegionId,
			"StackGroupName": stackGroupName,
			"PageNumber":     pageNum,
			"PageSize":       pageSize,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = client.RpcPost("ROS", "2019-09-10", action, nil, req, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotFound", "StackGroupNotFound"}) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, req)

		result, ok := response["StackInstances"].([]interface{})
		if !ok || len(result) == 0 {
			break
		}

		for _, item := range result {
			inst := item.(map[string]interface{})
			regId := getStr(inst, "RegionId")
			accId := getStr(inst, "AccountId")
			if regId == "" || accId == "" {
				continue
			}
			matchReg := len(targetRegions) == 0 || stringSliceContains(targetRegions, regId)
			matchAcc := len(targetAccounts) == 0 || stringSliceContains(targetAccounts, accId)
			if matchReg && matchAcc {
				matchedInstances = append(matchedInstances, inst)
			}
		}

		totalRaw, hasTotal := response["TotalCount"]
		if hasTotal {
			var total int
			switch t := totalRaw.(type) {
			case float64:
				total = int(t)
			case int:
				total = t
			case string:
				if parsed, err := strconv.Atoi(t); err == nil {
					total = parsed
				}
			}
			if pageNum*pageSize >= total {
				break
			}
		} else {
			break
		}
		pageNum++
	}

	var foundRegionIds, foundAccountIds []string
	var stateInstances []map[string]interface{}

	for _, inst := range matchedInstances {
		accId := getStr(inst, "AccountId")
		regId := getStr(inst, "RegionId")
		foundRegionIds = append(foundRegionIds, regId)
		foundAccountIds = append(foundAccountIds, accId)

		stateInstances = append(stateInstances, map[string]interface{}{
			"last_operation_id":    getStr(inst, "LastOperationId"),
			"status":               getStr(inst, "Status"),
			"stack_group_id":       getStr(inst, "StackGroupId"),
			"stack_id":             getStr(inst, "StackId"),
			"drift_detection_time": getStr(inst, "DriftDetectionTime"),
			"stack_drift_status":   getStr(inst, "StackDriftStatus"),
			"status_reason":        getStr(inst, "StatusReason"),
			"stack_group_name":     getStr(inst, "StackGroupName"),
			"account_id":           accId,
			"region_id":            regId,
			"rd_folder_id":         getStr(inst, "RdFolderId"),
		})
	}

	// Sort to prevent state drift diffs
	sort.Strings(foundRegionIds)
	sort.Strings(foundAccountIds)
	sort.SliceStable(stateInstances, func(i, j int) bool {
		if stateInstances[i]["account_id"].(string) == stateInstances[j]["account_id"].(string) {
			return stateInstances[i]["region_id"].(string) < stateInstances[j]["region_id"].(string)
		}
		return stateInstances[i]["account_id"].(string) < stateInstances[j]["account_id"].(string)
	})

	if len(foundRegionIds) > 0 {
		d.Set("region_ids", convertStringListToInterfaceList(foundRegionIds))
	}
	if len(foundAccountIds) > 0 {
		if _, isServiceManaged := d.GetOk("deployment_targets"); !isServiceManaged {
			d.Set("account_ids", convertStringListToInterfaceList(foundAccountIds))
		}
	} else if _, isServiceManaged := d.GetOk("deployment_targets"); !isServiceManaged {
		d.Set("account_ids", nil)
	}

	d.Set("stack_instances", stateInstances)
	// Persist non-returned attributes from state/config
	d.Set("timeout_in_minutes", d.Get("timeout_in_minutes"))
	d.Set("disable_rollback", d.Get("disable_rollback"))

	// Safely handle deployment_targets preservation
	if dtList, ok := d.Get("deployment_targets").([]interface{}); ok && len(dtList) > 0 {
		if dtMap, ok := dtList[0].(map[string]interface{}); ok {
			// Ensure empty lists are preserved correctly for Terraform diff suppression
			if accIds, ok := dtMap["account_ids"]; ok && accIds == nil {
				dtMap["account_ids"] = []interface{}{}
			}
			if folderIds, ok := dtMap["rd_folder_ids"]; ok && folderIds == nil {
				dtMap["rd_folder_ids"] = []interface{}{}
			}
			d.Set("deployment_targets", []map[string]interface{}{dtMap})
		}
	}

	if v := d.Get("operation_preferences"); v != nil {
		d.Set("operation_preferences", v)
	}
	if v := d.Get("operation_description"); v != nil {
		d.Set("operation_description", v)
	}
	if v := d.Get("parameter_overrides"); v != nil {
		d.Set("parameter_overrides", v)
	}
	if v := d.Get("deployment_options"); v != nil {
		d.Set("deployment_options", v)
	}

	return nil
}

func resourceAlicloudRosStackInstancesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateStackInstances"

	if !d.HasChanges("parameter_overrides", "operation_preferences", "timeout_in_minutes", "operation_description") {
		return nil
	}

	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": d.Get("stack_group_name").(string),
		"ClientToken":    buildClientToken(action),
	}

	if v, ok := d.GetOk("region_ids"); ok {
		request["RegionIds"] = convertListToJsonString(v.([]interface{}))
	}
	if dtJSON := expandDeploymentTargets(d); dtJSON != "" {
		request["DeploymentTargets"] = dtJSON
	} else if v, ok := d.GetOk("account_ids"); ok {
		request["AccountIds"] = convertListToJsonString(v.([]interface{}))
	}

	if d.HasChange("parameter_overrides") {
		v := d.Get("parameter_overrides").(*schema.Set).List()
		sort.SliceStable(v, func(i, j int) bool {
			return v[i].(map[string]interface{})["parameter_key"].(string) < v[j].(map[string]interface{})["parameter_key"].(string)
		})
		for i, param := range v {
			p := param.(map[string]interface{})
			idx := i + 1
			request[fmt.Sprintf("ParameterOverrides.%d.ParameterKey", idx)] = p["parameter_key"]
			if val, ok := p["parameter_value"]; ok && val.(string) != "" {
				request[fmt.Sprintf("ParameterOverrides.%d.ParameterValue", idx)] = val
			}
		}
	}

	if opPrefJSON := expandOperationPreferences(d); opPrefJSON != "" {
		request["OperationPreferences"] = opPrefJSON
	}
	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}
	if v, ok := d.GetOk("operation_description"); ok {
		request["OperationDescription"] = v
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StackGroupOperationInProgress"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	operationId := fmt.Sprint(response["OperationId"])
	if err := waitForRosStackGroupOperationAndCheckResults(client, d.Get("stack_group_name").(string), operationId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, "ROS UpdateStackInstances operation %s failed", operationId)
	}

	return resourceAlicloudRosStackInstancesRead(d, meta)
}

func resourceAlicloudRosStackInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStackInstances"
	stackGroupName := d.Get("stack_group_name").(string)

	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": stackGroupName,
		"RetainStacks":   false,
		"ClientToken":    buildClientToken(action),
	}

	hasTargets := false

	if v, ok := d.GetOk("region_ids"); ok && len(v.([]interface{})) > 0 {
		request["RegionIds"] = convertListToJsonString(uniqueInterfaceSlice(v.([]interface{})))
		hasTargets = true
	}

	if dtJSON := expandDeploymentTargets(d); dtJSON != "" {
		request["DeploymentTargets"] = dtJSON
		hasTargets = true
	} else if v, ok := d.GetOk("account_ids"); ok && len(v.([]interface{})) > 0 {
		request["AccountIds"] = convertListToJsonString(uniqueInterfaceSlice(v.([]interface{})))
		hasTargets = true
	}

	if !hasTargets {
		if instances, ok := d.GetOk("stack_instances"); ok {
			accSet := make(map[string]struct{})
			regSet := make(map[string]struct{})
			var accIds, regIds []string
			for _, inst := range instances.([]interface{}) {
				m := inst.(map[string]interface{})
				if aid, ok := m["account_id"].(string); ok && aid != "" {
					if _, exists := accSet[aid]; !exists {
						accSet[aid] = struct{}{}
						accIds = append(accIds, aid)
					}
				}
				if rid, ok := m["region_id"].(string); ok && rid != "" {
					if _, exists := regSet[rid]; !exists {
						regSet[rid] = struct{}{}
						regIds = append(regIds, rid)
					}
				}
			}
			if len(accIds) > 0 {
				request["AccountIds"] = convertListToJsonString(convertStringListToInterfaceList(accIds))
				hasTargets = true
			}
			if len(regIds) > 0 && request["RegionIds"] == nil {
				request["RegionIds"] = convertListToJsonString(convertStringListToInterfaceList(regIds))
				hasTargets = true
			}
		}
	}

	if !hasTargets {
		d.SetId("")
		return nil
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotFound", "InstanceNotFound", "StackGroupNotFound", "MissingRegionIds"}) {
				return nil
			}
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StackGroupOperationInProgress"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	operationId := fmt.Sprint(response["OperationId"])
	if operationId != "" && operationId != "<nil>" {
		if err := waitForRosStackGroupOperationAndCheckResults(client, stackGroupName, operationId, d.Timeout(schema.TimeoutDelete)); err != nil {
			if !IsExpectedErrors(err, []string{"ResourceNotFound", "InstanceNotFound"}) {
				log.Printf("[WARN] ROS DeleteStackInstances operation warning: %v", err)
			}
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		req := map[string]interface{}{
			"RegionId":       client.RegionId,
			"StackGroupName": stackGroupName,
			"PageNumber":     1,
			"PageSize":       50,
		}
		resp, err := client.RpcPost("ROS", "2019-09-10", "ListStackInstances", nil, req, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StackGroupNotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if resp == nil {
			return resource.RetryableError(fmt.Errorf("waiting for stack instances deletion verification"))
		}

		instances, _ := resp["StackInstances"].([]interface{})
		for _, item := range instances {
			inst := item.(map[string]interface{})
			accId := fmt.Sprintf("%v", inst["AccountId"])
			regId := fmt.Sprintf("%v", inst["RegionId"])

			targetAccounts := request["AccountIds"]
			targetRegions := request["RegionIds"]

			matchReg := targetRegions == nil || strings.Contains(fmt.Sprintf("%v", targetRegions), regId)
			matchAcc := targetAccounts == nil || strings.Contains(fmt.Sprintf("%v", targetAccounts), accId)

			if matchAcc && matchReg {
				return resource.RetryableError(fmt.Errorf("stack instances still exist in account %s region %s", accId, regId))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[WARN] Final cleanup wait warning: %v", err)
	}

	d.SetId("")
	return nil
}

// ================= Helper Functions =================

func expandDeploymentTargets(d *schema.ResourceData) string {
	v := d.Get("deployment_targets").([]interface{})
	if len(v) == 0 {
		return ""
	}
	cfg, ok := v[0].(map[string]interface{})
	if !ok {
		return ""
	}

	targets := map[string]interface{}{}
	if accs, ok := cfg["account_ids"]; ok {
		if arr, ok := accs.([]interface{}); ok && len(arr) > 0 {
			targets["AccountIds"] = arr
		}
	}
	if folders, ok := cfg["rd_folder_ids"]; ok {
		if arr, ok := folders.([]interface{}); ok && len(arr) > 0 {
			targets["RdFolderIds"] = arr
		}
	}
	if len(targets) == 0 {
		return ""
	}
	b, _ := json.Marshal(targets)
	return string(b)
}

func expandOperationPreferences(d *schema.ResourceData) string {
	v := d.Get("operation_preferences").([]interface{})
	if len(v) == 0 {
		return ""
	}
	cfg, ok := v[0].(map[string]interface{})
	if !ok {
		return ""
	}

	prefs := map[string]interface{}{}
	if maxCount, ok := cfg["max_concurrent_count"].(int); ok && maxCount > 0 {
		prefs["MaxConcurrentCount"] = maxCount
	} else if maxPct, ok := cfg["max_concurrent_percentage"].(int); ok && maxPct > 0 {
		prefs["MaxConcurrentPercentage"] = maxPct
	}

	if failCount, ok := cfg["failure_tolerance_count"].(int); ok && failCount > 0 {
		prefs["FailureToleranceCount"] = failCount
	} else if failPct, ok := cfg["failure_tolerance_percentage"].(int); ok && failPct > 0 {
		prefs["FailureTolerancePercentage"] = failPct
	}

	if val, ok := cfg["region_concurrency_type"].(string); ok && val != "" {
		prefs["RegionConcurrencyType"] = val
	}

	if len(prefs) == 0 {
		return ""
	}
	b, _ := json.Marshal(prefs)
	return string(b)
}

func waitForRosStackGroupOperationAndCheckResults(client *connectivity.AliyunClient, stackGroupName, operationId string, timeout time.Duration) error {
	if err := waitForRosStackGroupOperation(client, operationId, timeout); err != nil {
		return err
	}
	action := "ListStackGroupOperationResults"
	req := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": stackGroupName,
		"OperationId":    operationId,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, req, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, operationId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, req)

	if results, ok := response["Results"].([]interface{}); ok {
		for _, resultItem := range results {
			if result, ok := resultItem.(map[string]interface{}); ok {
				status := fmt.Sprint(result["Status"])
				if status == "FAILED" || status == "CANCELLED" {
					return fmt.Errorf("stack instance operation failed for account %v in region %v: %v",
						result["StackInstanceAccountId"], result["StackInstanceRegionId"], result["Reason"])
				}
			}
		}
	}
	return nil
}

func waitForRosStackGroupOperation(client *connectivity.AliyunClient, operationId string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		req := map[string]interface{}{"RegionId": client.RegionId, "OperationId": operationId}
		var response map[string]interface{}
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", "GetStackGroupOperation", nil, req, true)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("GetStackGroupOperation", response, req)

		opRaw, ok := response["StackGroupOperation"]
		if !ok || opRaw == nil {
			time.Sleep(5 * time.Second)
			return resource.RetryableError(fmt.Errorf("operation not available"))
		}
		op, _ := opRaw.(map[string]interface{})
		status := fmt.Sprint(op["Status"])
		switch status {
		case "SUCCEEDED":
			return nil
		case "FAILED", "CANCELLED":
			return resource.NonRetryableError(fmt.Errorf("operation ended: %s, Reason: %v", status, op["Reason"]))
		default:
			time.Sleep(5 * time.Second)
			return resource.RetryableError(fmt.Errorf("operation status: %s", status))
		}
	})
}

func suppressJsonStringDiff(k, old, new string, d *schema.ResourceData) bool {
	if old == new {
		return true
	}
	if old == "" || new == "" {
		return false
	}
	var o, n interface{}
	if err := json.Unmarshal([]byte(old), &o); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &n); err != nil {
		return false
	}
	return reflect.DeepEqual(o, n)
}

func stringSliceContains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func convertStringListToInterfaceList(list []string) []interface{} {
	result := make([]interface{}, len(list))
	for i, s := range list {
		result[i] = s
	}
	return result
}

func uniqueInterfaceSlice(list []interface{}) []interface{} {
	seen := make(map[string]struct{})
	var unique []interface{}
	for _, val := range list {
		if str, ok := val.(string); ok {
			if _, exists := seen[str]; !exists {
				seen[str] = struct{}{}
				unique = append(unique, str)
			}
		}
	}
	return unique
}

func suppressDeploymentTargetAccountIdsDiff(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "0" {
		return true
	}
	if old == "0" && new == "" {
		return true
	}
	return false
}
