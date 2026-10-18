// Hand-written resource for Cms AlertAction (CloudSpec resource published,
// @terraform mapping not yet generated to ACube; written from cspec definition).
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsAlertAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertActionCreate,
		Read:   resourceAliCloudCmsAlertActionRead,
		Update: resourceAliCloudCmsAlertActionUpdate,
		Delete: resourceAliCloudCmsAlertActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert_action_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"webhook_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"JSON", "FORM"}, false),
						},
						"headers": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"mns_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mns_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"sls_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logstore": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ess_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ess_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ess_rule_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"fc_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"pager_duty_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"fc3_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"eb_param": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"event_bus_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eb_source": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertActionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertAction"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	request["type"] = d.Get("type").(string)
	request["alertActionName"] = d.Get("alert_action_name").(string)

	if v := d.Get("webhook_param"); !IsNil(v) {
		request["webhookParam"] = expandCmsAlertActionWebhookParam(v)
	}
	if v := d.Get("mns_param"); !IsNil(v) {
		request["mnsParam"] = expandCmsAlertActionMnsParam(v)
	}
	if v := d.Get("sls_param"); !IsNil(v) {
		request["slsParam"] = expandCmsAlertActionSlsParam(v)
	}
	if v := d.Get("ess_param"); !IsNil(v) {
		request["essParam"] = expandCmsAlertActionEssParam(v)
	}
	if v := d.Get("fc_param"); !IsNil(v) {
		request["fcParam"] = expandCmsAlertActionFcParam(v)
	}
	if v := d.Get("pager_duty_param"); !IsNil(v) {
		request["pagerDutyParam"] = expandCmsAlertActionPagerDutyParam(v)
	}
	if v := d.Get("fc3_param"); !IsNil(v) {
		request["fc3Param"] = expandCmsAlertActionFc3Param(v)
	}
	if v := d.Get("eb_param"); !IsNil(v) {
		request["ebParam"] = expandCmsAlertActionEbParam(v)
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_action", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.alertActionId", response)
	if id == nil || fmt.Sprint(id) == "" {
		id, _ = jsonpath.Get("$.data", response)
	}
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsAlertActionRead(d, meta)
}

func resourceAliCloudCmsAlertActionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertAction(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_action DescribeCmsAlertAction Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alert_action_id", objectRaw["alertActionId"])
	d.Set("alert_action_name", objectRaw["alertActionName"])
	d.Set("type", objectRaw["type"])
	d.Set("region_id", objectRaw["regionId"])

	if v, ok := objectRaw["webhookParam"].(map[string]interface{}); ok && len(v) > 0 {
		webhookMaps := make([]map[string]interface{}, 0)
		webhookMap := make(map[string]interface{})
		webhookMap["method"] = v["method"]
		webhookMap["url"] = v["url"]
		webhookMap["content_type"] = v["contentType"]
		webhookMap["headers"] = v["headers"]
		webhookMaps = append(webhookMaps, webhookMap)
		if err := d.Set("webhook_param", webhookMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["mnsParam"].(map[string]interface{}); ok && len(v) > 0 {
		mnsMaps := make([]map[string]interface{}, 0)
		mnsMap := make(map[string]interface{})
		mnsMap["mns_type"] = v["mnsType"]
		mnsMap["name"] = v["name"]
		mnsMap["region_id"] = v["regionId"]
		mnsMaps = append(mnsMaps, mnsMap)
		if err := d.Set("mns_param", mnsMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["slsParam"].(map[string]interface{}); ok && len(v) > 0 {
		slsMaps := make([]map[string]interface{}, 0)
		slsMap := make(map[string]interface{})
		slsMap["logstore"] = v["logstore"]
		slsMap["project"] = v["project"]
		slsMap["region_id"] = v["regionId"]
		slsMaps = append(slsMaps, slsMap)
		if err := d.Set("sls_param", slsMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["essParam"].(map[string]interface{}); ok && len(v) > 0 {
		essMaps := make([]map[string]interface{}, 0)
		essMap := make(map[string]interface{})
		essMap["ess_group_id"] = v["essGroupId"]
		essMap["ess_rule_id"] = v["essRuleId"]
		essMap["region_id"] = v["regionId"]
		essMaps = append(essMaps, essMap)
		if err := d.Set("ess_param", essMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["fcParam"].(map[string]interface{}); ok && len(v) > 0 {
		fcMaps := make([]map[string]interface{}, 0)
		fcMap := make(map[string]interface{})
		fcMap["function"] = v["function"]
		fcMap["region_id"] = v["regionId"]
		fcMap["service"] = v["service"]
		fcMaps = append(fcMaps, fcMap)
		if err := d.Set("fc_param", fcMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["pagerDutyParam"].(map[string]interface{}); ok && len(v) > 0 {
		pdMaps := make([]map[string]interface{}, 0)
		pdMap := make(map[string]interface{})
		pdMap["key"] = v["key"]
		pdMap["url"] = v["url"]
		pdMaps = append(pdMaps, pdMap)
		if err := d.Set("pager_duty_param", pdMaps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["fc3Param"].(map[string]interface{}); ok && len(v) > 0 {
		fc3Maps := make([]map[string]interface{}, 0)
		fc3Map := make(map[string]interface{})
		fc3Map["region_id"] = v["regionId"]
		fc3Map["function"] = v["function"]
		fc3Map["qualifier"] = v["qualifier"]
		fc3Maps = append(fc3Maps, fc3Map)
		if err := d.Set("fc3_param", fc3Maps); err != nil {
			return err
		}
	}
	if v, ok := objectRaw["ebParam"].(map[string]interface{}); ok && len(v) > 0 {
		ebMaps := make([]map[string]interface{}, 0)
		ebMap := make(map[string]interface{})
		ebMap["region_id"] = v["regionId"]
		ebMap["event_bus_name"] = v["eventBusName"]
		ebMap["subject"] = v["subject"]
		ebMap["eb_source"] = v["ebSource"]
		ebMaps = append(ebMaps, ebMap)
		if err := d.Set("eb_param", ebMaps); err != nil {
			return err
		}
	}

	d.Set("alert_action_id", d.Id())
	return nil
}

func resourceAliCloudCmsAlertActionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/alertAction/%s", d.Id())
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	update := false
	if d.HasChange("alert_action_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_action_name"); ok || d.HasChange("alert_action_name") {
		request["alertActionName"] = v
	}
	if d.HasChange("webhook_param") {
		update = true
	}
	if v := d.Get("webhook_param"); !IsNil(v) || d.HasChange("webhook_param") {
		request["webhookParam"] = expandCmsAlertActionWebhookParam(v)
	}
	if d.HasChange("mns_param") {
		update = true
	}
	if v := d.Get("mns_param"); !IsNil(v) || d.HasChange("mns_param") {
		request["mnsParam"] = expandCmsAlertActionMnsParam(v)
	}
	if d.HasChange("sls_param") {
		update = true
	}
	if v := d.Get("sls_param"); !IsNil(v) || d.HasChange("sls_param") {
		request["slsParam"] = expandCmsAlertActionSlsParam(v)
	}
	if d.HasChange("ess_param") {
		update = true
	}
	if v := d.Get("ess_param"); !IsNil(v) || d.HasChange("ess_param") {
		request["essParam"] = expandCmsAlertActionEssParam(v)
	}
	if d.HasChange("fc_param") {
		update = true
	}
	if v := d.Get("fc_param"); !IsNil(v) || d.HasChange("fc_param") {
		request["fcParam"] = expandCmsAlertActionFcParam(v)
	}
	if d.HasChange("pager_duty_param") {
		update = true
	}
	if v := d.Get("pager_duty_param"); !IsNil(v) || d.HasChange("pager_duty_param") {
		request["pagerDutyParam"] = expandCmsAlertActionPagerDutyParam(v)
	}
	if d.HasChange("fc3_param") {
		update = true
	}
	if v := d.Get("fc3_param"); !IsNil(v) || d.HasChange("fc3_param") {
		request["fc3Param"] = expandCmsAlertActionFc3Param(v)
	}
	if d.HasChange("eb_param") {
		update = true
	}
	if v := d.Get("eb_param"); !IsNil(v) || d.HasChange("eb_param") {
		request["ebParam"] = expandCmsAlertActionEbParam(v)
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
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
	}

	return resourceAliCloudCmsAlertActionRead(d, meta)
}

func resourceAliCloudCmsAlertActionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertActions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	idsJSON, _ := json.Marshal([]string{d.Id()})
	query["alertActionIds"] = StringPointer(string(idsJSON))

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func expandCmsAlertActionWebhookParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["method"].(string); ok && s != "" {
		result["method"] = s
	}
	if s, ok := m["url"].(string); ok && s != "" {
		result["url"] = s
	}
	if s, ok := m["content_type"].(string); ok && s != "" {
		result["contentType"] = s
	}
	if hm, ok := m["headers"].(map[string]interface{}); ok && len(hm) > 0 {
		result["headers"] = hm
	}
	return result
}

func expandCmsAlertActionMnsParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["mns_type"].(string); ok && s != "" {
		result["mnsType"] = s
	}
	if s, ok := m["name"].(string); ok && s != "" {
		result["name"] = s
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	return result
}

func expandCmsAlertActionSlsParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["logstore"].(string); ok && s != "" {
		result["logstore"] = s
	}
	if s, ok := m["project"].(string); ok && s != "" {
		result["project"] = s
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	return result
}

func expandCmsAlertActionEssParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["ess_group_id"].(string); ok && s != "" {
		result["essGroupId"] = s
	}
	if s, ok := m["ess_rule_id"].(string); ok && s != "" {
		result["essRuleId"] = s
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	return result
}

func expandCmsAlertActionFcParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["function"].(string); ok && s != "" {
		result["function"] = s
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	if s, ok := m["service"].(string); ok && s != "" {
		result["service"] = s
	}
	return result
}

func expandCmsAlertActionPagerDutyParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["key"].(string); ok && s != "" {
		result["key"] = s
	}
	if s, ok := m["url"].(string); ok && s != "" {
		result["url"] = s
	}
	return result
}

func expandCmsAlertActionFc3Param(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	if s, ok := m["function"].(string); ok && s != "" {
		result["function"] = s
	}
	if s, ok := m["qualifier"].(string); ok && s != "" {
		result["qualifier"] = s
	}
	return result
}

func expandCmsAlertActionEbParam(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return result
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return result
	}
	if s, ok := m["region_id"].(string); ok && s != "" {
		result["regionId"] = s
	}
	if s, ok := m["event_bus_name"].(string); ok && s != "" {
		result["eventBusName"] = s
	}
	if s, ok := m["subject"].(string); ok && s != "" {
		result["subject"] = s
	}
	if s, ok := m["eb_source"].(string); ok && s != "" {
		result["ebSource"] = s
	}
	return result
}
