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

func resourceAliCloudEsaCompressionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCompressionRuleCreate,
		Read:   resourceAliCloudEsaCompressionRuleRead,
		Update: resourceAliCloudEsaCompressionRuleUpdate,
		Delete: resourceAliCloudEsaCompressionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"brotli": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gzip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
				ForceNew: true,
			},
			"zstd": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
		},
	}
}

func resourceAliCloudEsaCompressionRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCompressionRule"
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
	if v, ok := d.GetOk("zstd"); ok {
		request["Zstd"] = v
	}
	if v, ok := d.GetOkExists("site_version"); ok {
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
	if v, ok := d.GetOk("gzip"); ok {
		request["Gzip"] = v
	}
	if v, ok := d.GetOk("brotli"); ok {
		request["Brotli"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_compression_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaCompressionRuleUpdate(d, meta)
}

func resourceAliCloudEsaCompressionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCompressionRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_compression_rule DescribeEsaCompressionRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("brotli", objectRaw["Brotli"])
	d.Set("gzip", objectRaw["Gzip"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("zstd", objectRaw["Zstd"])
	d.Set("config_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprintf(parts[0]))

	return nil
}

func resourceAliCloudEsaCompressionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateCompressionRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if !d.IsNewResource() && d.HasChange("zstd") {
		update = true
		request["Zstd"] = d.Get("zstd")
	}

	if !d.IsNewResource() && d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if !d.IsNewResource() && d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if !d.IsNewResource() && d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if !d.IsNewResource() && d.HasChange("gzip") {
		update = true
		request["Gzip"] = d.Get("gzip")
	}

	if !d.IsNewResource() && d.HasChange("brotli") {
		update = true
		request["Brotli"] = d.Get("brotli")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
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

	return resourceAliCloudEsaCompressionRuleRead(d, meta)
}

func resourceAliCloudEsaCompressionRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCompressionRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
