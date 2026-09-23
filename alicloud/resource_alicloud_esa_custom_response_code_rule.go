// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaCustomResponseCodeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCustomResponseCodeRuleCreate,
		Read:   resourceAliCloudEsaCustomResponseCodeRuleRead,
		Update: resourceAliCloudEsaCustomResponseCodeRuleUpdate,
		Delete: resourceAliCloudEsaCustomResponseCodeRuleDelete,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"page_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"return_code": {
				Type:     schema.TypeString,
				Required: true,
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
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaCustomResponseCodeRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCustomResponseCodeRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOkExists("sequence"); ok {
		request["Sequence"] = v
	}
	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	request["ReturnCode"] = d.Get("return_code")
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	request["PageId"] = d.Get("page_id")
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_custom_response_code_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaCustomResponseCodeRuleRead(d, meta)
}

func resourceAliCloudEsaCustomResponseCodeRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCustomResponseCodeRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_custom_response_code_rule DescribeEsaCustomResponseCodeRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("page_id", objectRaw["PageId"])
	d.Set("return_code", objectRaw["ReturnCode"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("config_id", fmt.Sprint(objectRaw["ConfigId"]))

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))

	return nil
}

func resourceAliCloudEsaCustomResponseCodeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateCustomResponseCodeRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("sequence") {
		update = true

		if v, ok := d.GetOkExists("sequence"); ok {
			request["Sequence"] = v
		}
	}

	if d.HasChange("return_code") {
		update = true
	}
	request["ReturnCode"] = d.Get("return_code")
	if d.HasChange("rule_enable") {
		update = true
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}

	if d.HasChange("page_id") {
		update = true
	}
	request["PageId"] = d.Get("page_id")
	if d.HasChange("rule_name") {
		update = true
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if d.HasChange("rule") {
		update = true
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 0*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return resourceAliCloudEsaCustomResponseCodeRuleRead(d, meta)
}

func resourceAliCloudEsaCustomResponseCodeRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCustomResponseCodeRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
