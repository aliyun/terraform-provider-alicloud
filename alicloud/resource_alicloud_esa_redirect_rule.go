// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaRedirectRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRedirectRuleCreate,
		Read:   resourceAliCloudEsaRedirectRuleRead,
		Update: resourceAliCloudEsaRedirectRuleUpdate,
		Delete: resourceAliCloudEsaRedirectRuleDelete,
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
			"reserve_query_string": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
			"status_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"301", "302"}, false),
			},
			"target_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"static"}, false),
			},
		},
	}
}

func resourceAliCloudEsaRedirectRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRedirectRule"
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

	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	request["ReserveQueryString"] = d.Get("reserve_query_string")
	request["TargetUrl"] = d.Get("target_url")
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	request["Type"] = d.Get("type")
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	request["StatusCode"] = d.Get("status_code")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_redirect_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaRedirectRuleRead(d, meta)
}

func resourceAliCloudEsaRedirectRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRedirectRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_redirect_rule DescribeEsaRedirectRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["ReserveQueryString"] != nil {
		d.Set("reserve_query_string", objectRaw["ReserveQueryString"])
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
	if objectRaw["StatusCode"] != nil {
		d.Set("status_code", objectRaw["StatusCode"])
	}
	if objectRaw["TargetUrl"] != nil {
		d.Set("target_url", objectRaw["TargetUrl"])
	}
	if objectRaw["Type"] != nil {
		d.Set("type", objectRaw["Type"])
	}
	if objectRaw["ConfigId"] != nil {
		d.Set("config_id", objectRaw["ConfigId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaRedirectRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "UpdateRedirectRule"
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("reserve_query_string") {
		update = true
	}
	request["ReserveQueryString"] = d.Get("reserve_query_string")
	if d.HasChange("target_url") {
		update = true
	}
	request["TargetUrl"] = d.Get("target_url")
	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("type") {
		update = true
	}
	request["Type"] = d.Get("type")
	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if d.HasChange("status_code") {
		update = true
	}
	request["StatusCode"] = d.Get("status_code")
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

	return resourceAliCloudEsaRedirectRuleRead(d, meta)
}

func resourceAliCloudEsaRedirectRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteRedirectRule"
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
