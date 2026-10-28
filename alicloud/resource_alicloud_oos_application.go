package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudOosApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosApplicationCreate,
		Read:   resourceAlicloudOosApplicationRead,
		Update: resourceAlicloudOosApplicationUpdate,
		Delete: resourceAlicloudOosApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"application_source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alarm_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"health_check_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"contact_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"delete_alarm_rules_before_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

// expandOosApplicationAlarmConfig builds the AlarmConfig JSON string from the
// schema list. It is shared by Create and Update.
func expandOosApplicationAlarmConfig(d *schema.ResourceData) (string, bool, error) {
	v, ok := d.GetOk("alarm_config")
	if !ok {
		return "", false, nil
	}
	alarmConfig := make(map[string]interface{})
	for _, raw := range v.([]interface{}) {
		args, ok := raw.(map[string]interface{})
		if !ok || args == nil {
			continue
		}
		if templateIds := flattenOosApplicationSet(args["template_ids"]); len(templateIds) > 0 {
			alarmConfig["TemplateIds"] = templateIds
		}
		if healthCheckUrl, ok := args["health_check_url"].(string); ok && healthCheckUrl != "" {
			alarmConfig["HealthCheckUrl"] = healthCheckUrl
		}
		if contactGroups := flattenOosApplicationSet(args["contact_groups"]); len(contactGroups) > 0 {
			alarmConfig["ContactGroups"] = contactGroups
		}
	}
	if len(alarmConfig) == 0 {
		return "", false, nil
	}
	respJson, err := json.Marshal(alarmConfig)
	if err != nil {
		return "", false, WrapError(err)
	}
	return string(respJson), true, nil
}

// flattenOosApplicationSet normalizes a nested alarm_config collection field
// (the SDK returns *schema.Set for TypeSet, []interface{} for other forms)
// into a stable slice of strings for the AlarmConfig JSON payload.
func flattenOosApplicationSet(v interface{}) []string {
	var list []interface{}
	switch val := v.(type) {
	case *schema.Set:
		list = val.List()
	case []interface{}:
		list = val
	default:
		return nil
	}
	out := make([]string, 0, len(list))
	for _, item := range list {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func resourceAlicloudOosApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateApplication"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("application_source"); ok {
		request["ApplicationSource"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("deployment_type"); ok {
		request["DeploymentType"] = v
	}
	request["Name"] = d.Get("application_name")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("service_id"); ok {
		request["ServiceId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	if alarmConfigJson, hasAlarmConfig, err := expandOosApplicationAlarmConfig(d); err != nil {
		return err
	} else if hasAlarmConfig {
		request["AlarmConfig"] = alarmConfigJson
	}
	request["ClientToken"] = buildClientToken("CreateApplication")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_application", action, AlibabaCloudSdkGoERROR)
	}
	responseApplication := response["Application"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseApplication["Name"]))

	return resourceAlicloudOosApplicationRead(d, meta)
}
func resourceAlicloudOosApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_application oosService.DescribeOosApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("application_source", object["ApplicationSource"])
	d.Set("description", object["Description"])
	d.Set("application_name", object["Name"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("service_id", object["ServiceId"])
	d.Set("deployment_type", object["DeploymentType"])
	d.Set("tags", tagsToMap(object["Tags"]))
	alarmConfigSli := make([]map[string]interface{}, 0)
	if alarmConfig, ok := object["AlarmConfig"].(map[string]interface{}); ok && len(alarmConfig) > 0 {
		alarmConfigMap := make(map[string]interface{})
		if templateIds, ok := alarmConfig["TemplateIds"].([]interface{}); ok {
			alarmConfigMap["template_ids"] = templateIds
		}
		alarmConfigMap["health_check_url"] = alarmConfig["HealthCheckUrl"]
		if contactGroups, ok := alarmConfig["ContactGroups"].([]interface{}); ok {
			alarmConfigMap["contact_groups"] = contactGroups
		}
		alarmConfigSli = append(alarmConfigSli, alarmConfigMap)
	}
	d.Set("alarm_config", alarmConfigSli)
	return nil
}
func resourceAlicloudOosApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if d.HasChange("application_source") {
		update = true
		if v, ok := d.GetOk("application_source"); ok {
			request["ApplicationSource"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("alarm_config") {
		update = true
		if alarmConfigJson, hasAlarmConfig, err := expandOosApplicationAlarmConfig(d); err != nil {
			return err
		} else if hasAlarmConfig {
			request["AlarmConfig"] = alarmConfigJson
		} else {
			request["AlarmConfig"] = "{}"
		}
	}
	if d.HasChange("delete_alarm_rules_before_update") {
		update = true
		request["DeleteAlarmRulesBeforeUpdate"] = d.Get("delete_alarm_rules_before_update")
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Tags"] = respJson
		}
	}
	if update {
		action := "UpdateApplication"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudOosApplicationRead(d, meta)
}
func resourceAlicloudOosApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApplication"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"DeleteConflict.Application"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Application"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
