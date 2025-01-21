// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaHttpRequestHeaderModificationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaHttpRequestHeaderModificationRuleCreate,
		Read:   resourceAliCloudEsaHttpRequestHeaderModificationRuleRead,
		Update: resourceAliCloudEsaHttpRequestHeaderModificationRuleUpdate,
		Delete: resourceAliCloudEsaHttpRequestHeaderModificationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_header_modification": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaHttpRequestHeaderModificationRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHttpRequestHeaderModificationRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("request_header_modification"); ok {
		requestHeaderModificationMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Operation"] = dataLoopTmp["operation"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			requestHeaderModificationMapsArray = append(requestHeaderModificationMapsArray, dataLoopMap)
		}
		requestHeaderModificationMapsJson, err := json.Marshal(requestHeaderModificationMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RequestHeaderModification"] = string(requestHeaderModificationMapsJson)
	}

	if v, ok := d.GetOk("site_version"); ok {
		request["SiteVersion"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-09-10"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_http_request_header_modification_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaHttpRequestHeaderModificationRuleRead(d, meta)
}

func resourceAliCloudEsaHttpRequestHeaderModificationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaHttpRequestHeaderModificationRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_http_request_header_modification_rule DescribeEsaHttpRequestHeaderModificationRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Rule"] != nil {
		d.Set("rule", objectRaw["Rule"])
	}
	if objectRaw["RuleEnable"] != nil {
		d.Set("rule_enable", objectRaw["RuleEnable"])
	}
	if objectRaw["RuleName"] != nil {
		d.Set("rule_name", objectRaw["RuleName"])
	}
	if objectRaw["SiteVersion"] != nil {
		d.Set("site_version", objectRaw["SiteVersion"])
	}
	if objectRaw["ConfigId"] != nil {
		d.Set("config_id", objectRaw["ConfigId"])
	}

	requestHeaderModification1Raw := objectRaw["RequestHeaderModification"]
	requestHeaderModificationMaps := make([]map[string]interface{}, 0)
	if requestHeaderModification1Raw != nil {
		for _, requestHeaderModificationChild1Raw := range requestHeaderModification1Raw.([]interface{}) {
			requestHeaderModificationMap := make(map[string]interface{})
			requestHeaderModificationChild1Raw := requestHeaderModificationChild1Raw.(map[string]interface{})
			requestHeaderModificationMap["name"] = requestHeaderModificationChild1Raw["Name"]
			requestHeaderModificationMap["operation"] = requestHeaderModificationChild1Raw["Operation"]
			requestHeaderModificationMap["value"] = requestHeaderModificationChild1Raw["Value"]

			requestHeaderModificationMaps = append(requestHeaderModificationMaps, requestHeaderModificationMap)
		}
	}
	if objectRaw["RequestHeaderModification"] != nil {
		if err := d.Set("request_header_modification", requestHeaderModificationMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaHttpRequestHeaderModificationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "UpdateHttpRequestHeaderModificationRule"
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("request_header_modification") {
		update = true
	}
	if v, ok := d.GetOk("request_header_modification"); ok || d.HasChange("request_header_modification") {
		requestHeaderModificationMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Operation"] = dataLoopTmp["operation"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			requestHeaderModificationMapsArray = append(requestHeaderModificationMapsArray, dataLoopMap)
		}
		requestHeaderModificationMapsJson, err := json.Marshal(requestHeaderModificationMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RequestHeaderModification"] = string(requestHeaderModificationMapsJson)
	}

	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-09-10"), StringPointer("AK"), query, request, &runtime)
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

	return resourceAliCloudEsaHttpRequestHeaderModificationRuleRead(d, meta)
}

func resourceAliCloudEsaHttpRequestHeaderModificationRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteHttpRequestHeaderModificationRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-09-10"), StringPointer("AK"), query, request, &runtime)

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
