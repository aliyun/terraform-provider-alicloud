// Package alicloud. This file is hand-written equivalent to the generator output for the BizTrace resource.
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsBizTrace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsBizTraceCreate,
		Read:   resourceAliCloudCmsBizTraceRead,
		Update: resourceAliCloudCmsBizTraceUpdate,
		Delete: resourceAliCloudCmsBizTraceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"biz_trace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"biz_trace_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"biz_trace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_config": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bizTraceJsonSubsetDiffSuppress,
			},
			"advanced_config": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bizTraceJsonSubsetDiffSuppress,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsBizTraceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	action := "/bizTrace"
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	body["bizTraceCode"] = d.Get("biz_trace_code")
	if v, ok := d.GetOk("biz_trace_name"); ok {
		body["bizTraceName"] = v
	}
	if v, ok := d.GetOk("rule_config"); ok {
		body["ruleConfig"] = v
	}
	if v, ok := d.GetOk("advanced_config"); ok {
		body["advancedConfig"] = v
	}
	if v, ok := d.GetOk("workspace"); ok {
		body["workspace"] = v
	}

	// The EntityStore for the workspace must be explicitly initialised before
	// BizTrace creation; it is not auto-provisioned by workspace creation.
	if workspaceName, ok := d.GetOk("workspace"); ok {
		if err = cmsServiceV2.CreateCmsEntityStore(workspaceName.(string)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_biz_trace", "CreateEntityStore", AlibabaCloudSdkGoERROR)
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			// NotFound.EntityStore is returned right after workspace creation while the
			// backend entityStore resource is still being asynchronously provisioned;
			// retry until it becomes ready. Keep NeedRetry for throttling/5xx/network.
			if NeedRetry(err) || strings.Contains(err.Error(), "NotFound.EntityStore") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_biz_trace", action, AlibabaCloudSdkGoERROR)
	}

	bizTraceId := fmt.Sprint(response["bizTraceId"])
	if bizTraceId == "" {
		return WrapError(Error("failed to create BizTrace: empty bizTraceId in response"))
	}
	d.SetId(bizTraceId)

	return resourceAliCloudCmsBizTraceRead(d, meta)
}

func resourceAliCloudCmsBizTraceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	object, err := cmsServiceV2.DescribeCmsBizTrace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_biz_trace DescribeCmsBizTrace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("biz_trace_id", object["bizTraceId"])
	d.Set("biz_trace_code", object["bizTraceCode"])
	d.Set("biz_trace_name", object["bizTraceName"])
	d.Set("rule_config", object["ruleConfig"])
	d.Set("advanced_config", object["advancedConfig"])
	d.Set("workspace", object["workspace"])
	d.Set("create_time", object["createTime"])
	d.Set("region_id", object["regionId"])

	return nil
}

func resourceAliCloudCmsBizTraceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/bizTrace/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	update := false
	if d.HasChange("biz_trace_name") {
		update = true
	}
	if d.HasChange("rule_config") {
		update = true
	}
	if d.HasChange("advanced_config") {
		update = true
	}

	if v, ok := d.GetOk("biz_trace_name"); ok || d.HasChange("biz_trace_name") {
		body["bizTraceName"] = v
	}
	if v, ok := d.GetOk("rule_config"); ok || d.HasChange("rule_config") {
		body["ruleConfig"] = v
	}
	if v, ok := d.GetOk("advanced_config"); ok || d.HasChange("advanced_config") {
		body["advancedConfig"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsBizTraceRead(d, meta)
}

func resourceAliCloudCmsBizTraceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/bizTrace/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

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
	addDebug(action, response, query)

	if err != nil {
		if IsExpectedErrors(err, []string{"404", "BizTraceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// bizTraceJsonSubsetDiffSuppress suppresses diffs for rule_config / advanced_config.
// The BizTrace API enriches these JSON strings with server-resolved keys on read
// (e.g. callType, entranceServiceId, dyingScope), so the stored value is a superset
// of the configured value; a diff is real only when the configured JSON is not a
// recursive subset of the stored JSON.
func bizTraceJsonSubsetDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if old == "" || new == "" {
		return old == new
	}
	return bizTraceJsonIsSubset(new, old)
}

func bizTraceJsonIsSubset(sub, sup string) bool {
	var subValue, supValue interface{}
	if err := json.Unmarshal([]byte(sub), &subValue); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(sup), &supValue); err != nil {
		return false
	}
	return bizTraceJsonValueIsSubset(subValue, supValue)
}

func bizTraceJsonValueIsSubset(sub, sup interface{}) bool {
	switch subTyped := sub.(type) {
	case map[string]interface{}:
		supMap, ok := sup.(map[string]interface{})
		if !ok {
			return false
		}
		for key, value := range subTyped {
			supEntry, exists := supMap[key]
			if !exists || !bizTraceJsonValueIsSubset(value, supEntry) {
				return false
			}
		}
		return true
	case []interface{}:
		supArray, ok := sup.([]interface{})
		if !ok || len(subTyped) > len(supArray) {
			return false
		}
		for i, value := range subTyped {
			if !bizTraceJsonValueIsSubset(value, supArray[i]) {
				return false
			}
		}
		return true
	default:
		return fmt.Sprintf("%v", sub) == fmt.Sprintf("%v", sup)
	}
}
