package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
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
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    50,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of target account IDs for service-managed permissions.",
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
	var response map[string]interface{}
	action := "CreateStackInstances"
	request := make(map[string]interface{})

	request["RegionId"] = client.RegionId
	request["StackGroupName"] = d.Get("stack_group_name")

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
	if err := resourceAlicloudRosStackInstancesRead(d, meta); err != nil {
		return err
	}
	injectLastOperationId(d, operationId)
	return nil
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

	var targetRegions []string
	if v, ok := d.GetOk("region_ids"); ok {
		for _, r := range v.([]interface{}) {
			targetRegions = append(targetRegions, r.(string))
		}
	}
	var targetAccounts []string
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
		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"StackGroupName": stackGroupName,
			"PageNumber":     pageNum,
			"PageSize":       pageSize,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
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
				log.Printf("[DEBUG] Stack Group %s not found, marking instances resource as deleted", stackGroupName)
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

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
			if (len(targetRegions) == 0 || stringSliceContains(targetRegions, regId)) &&
				(len(targetAccounts) == 0 || stringSliceContains(targetAccounts, accId)) {
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
			if float64(pageNum*pageSize) >= float64(total) {
				break
			}
		} else {
			break
		}
		pageNum++
	}

	var stateInstances []map[string]interface{}
	var foundRegionIds []interface{}
	var foundAccountIds []interface{}

	for _, inst := range matchedInstances {
		accId := getStr(inst, "AccountId")
		regId := getStr(inst, "RegionId")
		foundRegionIds = append(foundRegionIds, regId)
		foundAccountIds = append(foundAccountIds, accId)

		stateInst := map[string]interface{}{
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
		}
		stateInstances = append(stateInstances, stateInst)
	}

	if len(foundRegionIds) > 0 {
		d.Set("region_ids", foundRegionIds)
	} else {
		d.Set("region_ids", d.Get("region_ids"))
	}
	if len(foundAccountIds) > 0 {
		if _, isServiceManaged := d.GetOk("deployment_targets"); !isServiceManaged {
			d.Set("account_ids", foundAccountIds)
		}
	} else {
		d.Set("account_ids", d.Get("account_ids"))
	}

	if len(stateInstances) > 0 {
		d.Set("stack_instances", stateInstances)
	} else {
		d.Set("stack_instances", d.Get("stack_instances"))
	}

	d.Set("timeout_in_minutes", d.Get("timeout_in_minutes"))
	d.Set("disable_rollback", d.Get("disable_rollback"))

	if dtRaw, ok := d.GetOk("deployment_targets"); ok {
		d.Set("deployment_targets", dtRaw)
	}
	if opPrefRaw, ok := d.GetOk("operation_preferences"); ok {
		d.Set("operation_preferences", opPrefRaw)
	}
	if v, ok := d.GetOk("operation_description"); ok {
		d.Set("operation_description", v)
	}
	if v, ok := d.GetOk("parameter_overrides"); ok {
		d.Set("parameter_overrides", v)
	}
	if v, ok := d.GetOk("deployment_options"); ok {
		d.Set("deployment_options", v)
	}

	return nil
}

func resourceAlicloudRosStackInstancesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateStackInstances"

	if !d.HasChanges("parameter_overrides", "operation_preferences", "timeout_in_minutes", "operation_description") {
		return resourceAlicloudRosStackInstancesRead(d, meta)
	}

	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": d.Get("stack_group_name"),
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

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)

	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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

	if err := resourceAlicloudRosStackInstancesRead(d, meta); err != nil {
		return err
	}
	injectLastOperationId(d, operationId)
	return nil
}

func resourceAlicloudRosStackInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStackInstances"
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"StackGroupName": d.Get("stack_group_name"),
	}

	hasRegionIds := false
	if rawList, ok := d.Get("region_ids").([]interface{}); ok && len(rawList) > 0 {
		seen := make(map[string]struct{})
		var uniqueRegions []interface{}
		for _, val := range rawList {
			if str, ok := val.(string); ok {
				if _, exists := seen[str]; !exists {
					seen[str] = struct{}{}
					uniqueRegions = append(uniqueRegions, str)
				}
			}
		}
		if len(uniqueRegions) > 0 {
			request["RegionIds"] = convertListToJsonString(uniqueRegions)
			hasRegionIds = true
		}
	}

	hasDeploymentTargets := false
	if dtJSON := expandDeploymentTargets(d); dtJSON != "" {
		request["DeploymentTargets"] = dtJSON
		hasDeploymentTargets = true
	}

	hasAccountIds := false
	if !hasDeploymentTargets {
		if rawList, ok := d.Get("account_ids").([]interface{}); ok && len(rawList) > 0 {
			request["AccountIds"] = convertListToJsonString(rawList)
			hasAccountIds = true
		}
	}

	if !hasAccountIds && !hasDeploymentTargets {
		if v, ok := d.GetOk("stack_instances"); ok {
			instances := v.([]interface{})
			if len(instances) > 0 {
				var accIds []string
				for _, inst := range instances {
					m := inst.(map[string]interface{})
					if aid, ok := m["account_id"].(string); ok && aid != "" {
						accIds = append(accIds, aid)
					}
				}
				if len(accIds) > 0 {
					request["AccountIds"] = convertListToJsonString(convertStringListToInterfaceList(accIds))
					hasAccountIds = true
				}
				if !hasRegionIds {
					var regIds []string
					for _, inst := range instances {
						m := inst.(map[string]interface{})
						if rid, ok := m["region_id"].(string); ok && rid != "" {
							regIds = append(regIds, rid)
						}
					}
					if len(regIds) > 0 {
						request["RegionIds"] = convertListToJsonString(convertStringListToInterfaceList(regIds))
						hasRegionIds = true
					}
				}
			}
		}
	}

	if !hasAccountIds && !hasDeploymentTargets {
		d.SetId("")
		return nil
	}

	request["RetainStacks"] = false
	request["ClientToken"] = buildClientToken(action)

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
	if response == nil {
		d.SetId("")
		return nil
	}

	operationId := fmt.Sprint(response["OperationId"])
	if err := waitForRosStackGroupOperationAndCheckResults(client, d.Get("stack_group_name").(string), operationId, d.Timeout(schema.TimeoutDelete)); err != nil {
		if !IsExpectedErrors(err, []string{"ResourceNotFound", "InstanceNotFound"}) {
			log.Printf("[WARN] ROS DeleteStackInstances operation warning: %v", err)
		}
	}

	time.Sleep(15 * time.Second)

	stackGroupName := d.Get("stack_group_name").(string)
	var targetAccounts []string
	if hasAccountIds {
		if v, ok := d.GetOk("account_ids"); ok {
			for _, a := range v.([]interface{}) {
				targetAccounts = append(targetAccounts, a.(string))
			}
		}
	}
	var targetRegions []string
	if v, ok := d.GetOk("region_ids"); ok {
		for _, r := range v.([]interface{}) {
			targetRegions = append(targetRegions, r.(string))
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
		instances, _ := resp["StackInstances"].([]interface{})
		for _, item := range instances {
			inst := item.(map[string]interface{})
			accId := fmt.Sprintf("%v", inst["AccountId"])
			regId := fmt.Sprintf("%v", inst["RegionId"])

			matchAcc := len(targetAccounts) == 0 || stringSliceContains(targetAccounts, accId)
			matchReg := len(targetRegions) == 0 || stringSliceContains(targetRegions, regId)

			if matchAcc && matchReg {
				return resource.RetryableError(fmt.Errorf("waiting for stack instances to be fully deleted from cloud"))
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
	if v, ok := d.GetOk("deployment_targets"); ok {
		cfg := v.([]interface{})[0].(map[string]interface{})
		targets := map[string]interface{}{}
		if accs, ok := cfg["account_ids"]; ok && len(accs.([]interface{})) > 0 {
			targets["AccountIds"] = accs
		}
		if folders, ok := cfg["rd_folder_ids"]; ok && len(folders.([]interface{})) > 0 {
			targets["RdFolderIds"] = folders
		}
		if len(targets) > 0 {
			b, _ := json.Marshal(targets)
			return string(b)
		}
	}
	return ""
}

func expandOperationPreferences(d *schema.ResourceData) string {
	if v, ok := d.GetOk("operation_preferences"); ok {
		cfg := v.([]interface{})[0].(map[string]interface{})
		prefs := map[string]interface{}{}

		maxCount := cfg["max_concurrent_count"].(int)
		maxPct := cfg["max_concurrent_percentage"].(int)
		if maxCount > 0 {
			prefs["MaxConcurrentCount"] = maxCount
		} else if maxPct > 0 {
			prefs["MaxConcurrentPercentage"] = maxPct
		}

		failCount := cfg["failure_tolerance_count"].(int)
		failPct := cfg["failure_tolerance_percentage"].(int)
		if failCount > 0 {
			prefs["FailureToleranceCount"] = failCount
		} else if failPct > 0 {
			prefs["FailureTolerancePercentage"] = failPct
		}

		if val := cfg["region_concurrency_type"].(string); val != "" {
			prefs["RegionConcurrencyType"] = val
		}

		if len(prefs) > 0 {
			b, _ := json.Marshal(prefs)
			return string(b)
		}
	}
	return ""
}

func injectLastOperationId(d *schema.ResourceData, operationId string) {
	if instances, ok := d.GetOk("stack_instances"); ok {
		list := instances.([]interface{})
		for i, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				m["last_operation_id"] = operationId
				list[i] = m
			}
		}
		d.Set("stack_instances", list)
	}
}

func waitForRosStackGroupOperationAndCheckResults(client *connectivity.AliyunClient, stackGroupName, operationId string, timeout time.Duration) error {
	if err := waitForRosStackGroupOperation(client, operationId, timeout); err != nil {
		return err
	}
	action := "ListStackGroupOperationResults"
	request := map[string]interface{}{
		"RegionId": client.RegionId, "StackGroupName": stackGroupName, "OperationId": operationId,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
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
	addDebug(action, response, request)

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
		request := map[string]interface{}{"RegionId": client.RegionId, "OperationId": operationId}
		var response map[string]interface{}
		var err error
		response, err = client.RpcPost("ROS", "2019-09-10", "GetStackGroupOperation", nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("GetStackGroupOperation", response, request)

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
	var result []interface{}
	for _, s := range list {
		result = append(result, s)
	}
	return result
}
