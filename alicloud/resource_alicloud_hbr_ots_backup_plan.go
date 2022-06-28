package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbrOtsBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrOtsBackupPlanCreate,
		Read:   resourceAlicloudHbrOtsBackupPlanRead,
		Update: resourceAlicloudHbrOtsBackupPlanUpdate,
		Delete: resourceAlicloudHbrOtsBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE"}, false),
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"ots_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retention": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'schedule' has been deprecated from version 1.163.0. Use 'rules' instead.",
			},
			"vault_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ots_detail": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_names": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudHbrOtsBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}

	request["SourceType"] = "OTS"
	request["PlanName"] = d.Get("ots_backup_plan_name")
	request["BackupType"] = d.Get("backup_type")
	request["VaultId"] = d.Get("vault_id")
	request["Schedule"] = d.Get("schedule")
	request["Retention"] = d.Get("retention")

	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("ots_detail"); ok {
		otsDetails := make(map[string]interface{}, 0)
		for _, otsdetail := range v.(*schema.Set).List() {
			otsDetails["TableNames"] = otsdetail.(map[string]interface{})["table_names"]
		}
		ots, _ := json.Marshal(otsDetails)
		request["OtsDetail"] = string(ots)
	}

	if v, ok := d.GetOk("rules"); ok {
		for index, raw := range v.(*schema.Set).List() {
			ruleObj := raw.(map[string]interface{})
			if v, ok := ruleObj["schedule"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.Schedule", index+1)] = ruleObj["schedule"]
			}
			if v, ok := ruleObj["retention"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.Retention", index+1)] = ruleObj["retention"]
			}
			if _, ok := ruleObj["disabled"]; ok {
				request[fmt.Sprintf("Rule.%d.Disabled", index+1)] = ruleObj["disabled"]
			}
			if v, ok := ruleObj["rule_name"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.RuleName", index+1)] = ruleObj["rule_name"]
			}
			if v, ok := ruleObj["backup_type"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.BackupType", index+1)] = ruleObj["backup_type"]
			}
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_ots_backup_plan", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return resourceAlicloudHbrOtsBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrOtsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrOtsBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_ots_backup_plan hbrService.DescribeHbrOtsBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backup_type", object["BackupType"])
	d.Set("ots_backup_plan_name", object["PlanName"])
	d.Set("instance_name", object["InstanceName"])
	d.Set("retention", fmt.Sprint(formatInt(object["Retention"])))
	d.Set("schedule", object["Schedule"])
	d.Set("vault_id", object["VaultId"])
	d.Set("disabled", object["Disabled"])

	otsDetails := make([]map[string]interface{}, 0)
	if v, ok := object["OtsDetail"].(map[string]interface{}); ok {
		otsDetail := make(map[string]interface{}, 0)
		otsDetail["table_names"] = v["TableNames"].(map[string]interface{})["TableName"]

		otsDetails = append(otsDetails, otsDetail)
	}
	d.Set("ots_detail", otsDetails)

	rules := make([]map[string]interface{}, 0)
	if v, ok := object["Rules"].(map[string]interface{}); ok {
		for _, val := range v["Rule"].([]interface{}) {
			raw := val.(map[string]interface{})
			obj := map[string]interface{}{
				"schedule":    raw["Schedule"],
				"retention":   raw["Retention"],
				"disabled":    raw["Disabled"],
				"rule_name":   raw["RuleName"],
				"backup_type": raw["BackupType"],
			}
			rules = append(rules, obj)
		}
	}
	d.Set("rules", rules)

	return nil
}
func resourceAlicloudHbrOtsBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PlanId":     d.Id(),
		"SourceType": "OTS",
	}
	if !d.IsNewResource() && d.HasChange("ots_backup_plan_name") {
		update = true
	}
	request["PlanName"] = d.Get("ots_backup_plan_name")

	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
	}
	request["Schedule"] = d.Get("schedule")

	if !d.IsNewResource() && d.HasChange("retention") {
		update = true
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}

	if !d.IsNewResource() && d.HasChange("vault_id") {
		update = true
	}
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}

	if !d.IsNewResource() && d.HasChange("ots_detail") {
		update = true
	}
	if v, ok := d.GetOk("ots_detail"); ok {
		otsDetails := make(map[string]interface{}, 0)
		for _, otsdetail := range v.(*schema.Set).List() {
			otsDetails["TableNames"] = otsdetail.(map[string]interface{})["table_names"]
		}
		ots, _ := json.Marshal(otsDetails)
		request["OtsDetail"] = string(ots)
	}

	if !d.IsNewResource() && d.HasChange("rules") {
		update = true
	}
	if v, ok := d.GetOk("rules"); ok {
		for index, raw := range v.(*schema.Set).List() {
			ruleObj := raw.(map[string]interface{})
			if v, ok := ruleObj["schedule"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.Schedule", index+1)] = ruleObj["schedule"]
			}
			if v, ok := ruleObj["retention"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.Retention", index+1)] = ruleObj["retention"]
			}
			if _, ok := ruleObj["disabled"]; ok {
				request[fmt.Sprintf("Rule.%d.Disabled", index+1)] = ruleObj["disabled"]
			}
			if v, ok := ruleObj["rule_name"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.RuleName", index+1)] = ruleObj["rule_name"]
			}
			if v, ok := ruleObj["backup_type"]; ok && len(v.(string)) != 0 {
				request[fmt.Sprintf("Rule.%d.BackupType", index+1)] = ruleObj["backup_type"]
			}
		}
	}

	if update {
		action := "UpdateBackupPlan"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("retention")
		d.SetPartial("vault_id")
		d.SetPartial("ots_backup_plan_name")
		d.SetPartial("schedule")
		d.SetPartial("update_paths")
		d.SetPartial("path")
		d.SetPartial("ots_detail")
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrOtsBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("disabled").(bool))
		if strconv.FormatBool(object["Disabled"].(bool)) != target {
			action := "EnableBackupPlan"
			if target == "false" {
				action = "EnableBackupPlan"
			} else {
				action = "DisableBackupPlan"
			}
			request := map[string]interface{}{
				"PlanId": d.Id(),
			}
			request["VaultId"] = d.Get("vault_id")
			request["SourceType"] = "OTS"

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
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
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
		}
	}
	d.SetPartial("disabled")
	d.Partial(false)
	return resourceAlicloudHbrOtsBackupPlanRead(d, meta)
}
func resourceAlicloudHbrOtsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupPlan"
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}

	request["SourceType"] = "OTS"
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
