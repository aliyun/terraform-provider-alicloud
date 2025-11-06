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

func resourceAliCloudEsaOriginRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaOriginRuleCreate,
		Read:   resourceAliCloudEsaOriginRuleRead,
		Update: resourceAliCloudEsaOriginRuleUpdate,
		Delete: resourceAliCloudEsaOriginRuleDelete,
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
			"dns_record": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"follow302_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"follow302_max_tries": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"follow302_retain_args": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"follow302_retain_header": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"follow302_target_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_http_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_https_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_mtls": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"origin_read_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_scheme": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"http", "https", "follow"}, false),
			},
			"origin_sni": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin_verify": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"range": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off", "force"}, false),
			},
			"range_chunk_size": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceAliCloudEsaOriginRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOriginRule"
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
	if v, ok := d.GetOk("range"); ok {
		request["Range"] = v
	}
	if v, ok := d.GetOk("origin_https_port"); ok {
		request["OriginHttpsPort"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("follow302_max_tries"); ok {
		request["Follow302MaxTries"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("origin_host"); ok {
		request["OriginHost"] = v
	}
	if v, ok := d.GetOk("origin_scheme"); ok {
		request["OriginScheme"] = v
	}
	if v, ok := d.GetOk("follow302_enable"); ok {
		request["Follow302Enable"] = v
	}
	if v, ok := d.GetOk("follow302_retain_args"); ok {
		request["Follow302RetainArgs"] = v
	}
	if v, ok := d.GetOk("follow302_target_host"); ok {
		request["Follow302TargetHost"] = v
	}
	if v, ok := d.GetOk("origin_mtls"); ok {
		request["OriginMtls"] = v
	}
	if v, ok := d.GetOk("origin_sni"); ok {
		request["OriginSni"] = v
	}
	if v, ok := d.GetOk("dns_record"); ok {
		request["DnsRecord"] = v
	}
	if v, ok := d.GetOk("origin_verify"); ok {
		request["OriginVerify"] = v
	}
	if v, ok := d.GetOk("follow302_retain_header"); ok {
		request["Follow302RetainHeader"] = v
	}
	if v, ok := d.GetOk("range_chunk_size"); ok {
		request["RangeChunkSize"] = v
	}
	if v, ok := d.GetOk("origin_http_port"); ok {
		request["OriginHttpPort"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	if v, ok := d.GetOk("origin_read_timeout"); ok {
		request["OriginReadTimeout"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_origin_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaOriginRuleRead(d, meta)
}

func resourceAliCloudEsaOriginRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaOriginRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_origin_rule DescribeEsaOriginRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dns_record", objectRaw["DnsRecord"])
	d.Set("follow302_enable", objectRaw["Follow302Enable"])
	d.Set("follow302_max_tries", objectRaw["Follow302MaxTries"])
	d.Set("follow302_retain_args", objectRaw["Follow302RetainArgs"])
	d.Set("follow302_retain_header", objectRaw["Follow302RetainHeader"])
	d.Set("follow302_target_host", objectRaw["Follow302TargetHost"])
	d.Set("origin_host", objectRaw["OriginHost"])
	d.Set("origin_http_port", objectRaw["OriginHttpPort"])
	d.Set("origin_https_port", objectRaw["OriginHttpsPort"])
	d.Set("origin_mtls", objectRaw["OriginMtls"])
	d.Set("origin_read_timeout", objectRaw["OriginReadTimeout"])
	d.Set("origin_scheme", objectRaw["OriginScheme"])
	d.Set("origin_sni", objectRaw["OriginSni"])
	d.Set("origin_verify", objectRaw["OriginVerify"])
	d.Set("range", objectRaw["Range"])
	d.Set("range_chunk_size", objectRaw["RangeChunkSize"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("config_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaOriginRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateOriginRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if d.HasChange("range") {
		update = true
		request["Range"] = d.Get("range")
	}

	if d.HasChange("origin_https_port") {
		update = true
		request["OriginHttpsPort"] = d.Get("origin_https_port")
	}

	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("follow302_max_tries") {
		update = true
		request["Follow302MaxTries"] = d.Get("follow302_max_tries")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("origin_host") {
		update = true
		request["OriginHost"] = d.Get("origin_host")
	}

	if d.HasChange("origin_scheme") {
		update = true
		request["OriginScheme"] = d.Get("origin_scheme")
	}

	if d.HasChange("follow302_enable") {
		update = true
		request["Follow302Enable"] = d.Get("follow302_enable")
	}

	if d.HasChange("follow302_retain_args") {
		update = true
		request["Follow302RetainArgs"] = d.Get("follow302_retain_args")
	}

	if d.HasChange("follow302_target_host") {
		update = true
		request["Follow302TargetHost"] = d.Get("follow302_target_host")
	}

	if d.HasChange("origin_mtls") {
		update = true
		request["OriginMtls"] = d.Get("origin_mtls")
	}

	if d.HasChange("origin_sni") {
		update = true
		request["OriginSni"] = d.Get("origin_sni")
	}

	if d.HasChange("dns_record") {
		update = true
		request["DnsRecord"] = d.Get("dns_record")
	}

	if d.HasChange("origin_verify") {
		update = true
		request["OriginVerify"] = d.Get("origin_verify")
	}

	if d.HasChange("follow302_retain_header") {
		update = true
		request["Follow302RetainHeader"] = d.Get("follow302_retain_header")
	}

	if d.HasChange("range_chunk_size") {
		update = true
		request["RangeChunkSize"] = d.Get("range_chunk_size")
	}

	if d.HasChange("origin_http_port") {
		update = true
		request["OriginHttpPort"] = d.Get("origin_http_port")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if d.HasChange("origin_read_timeout") {
		update = true
		request["OriginReadTimeout"] = d.Get("origin_read_timeout")
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

	return resourceAliCloudEsaOriginRuleRead(d, meta)
}

func resourceAliCloudEsaOriginRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteOriginRule"
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
