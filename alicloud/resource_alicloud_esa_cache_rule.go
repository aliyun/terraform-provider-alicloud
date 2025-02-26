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

func resourceAliCloudEsaCacheRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCacheRuleCreate,
		Read:   resourceAliCloudEsaCacheRuleRead,
		Update: resourceAliCloudEsaCacheRuleUpdate,
		Delete: resourceAliCloudEsaCacheRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"additional_cacheable_ports": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"8880", "2052", "2082", "2086", "2095", "2053", "2083", "2087", "2096"}, false),
			},
			"browser_cache_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"no_cache", "follow_origin", "override_origin"}, false),
			},
			"browser_cache_ttl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bypass_cache": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"cache_all", "bypass_all"}, false),
			},
			"cache_deception_armor": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"cache_reserve_eligibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"bypass_cache_reserve", "eligible_for_cache_reserve"}, false),
			},
			"cache_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_presence_cookie": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_presence_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edge_cache_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"follow_origin", "no_cache", "override_origin", "follow_origin_bypass"}, false),
			},
			"edge_cache_ttl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edge_status_code_cache_ttl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_cookie": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_string_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ignore_all", "exclude_query_string", "reserve_all", "include_query_string"}, false),
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
			"serve_stale": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
			"sort_query_string_for_cache": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"user_device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"user_geo": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"user_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
		},
	}
}

func resourceAliCloudEsaCacheRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCacheRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("check_presence_cookie"); ok {
		request["CheckPresenceCookie"] = v
	}
	if v, ok := d.GetOk("additional_cacheable_ports"); ok {
		request["AdditionalCacheablePorts"] = v
	}
	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	if v, ok := d.GetOk("query_string_mode"); ok {
		request["QueryStringMode"] = v
	}
	if v, ok := d.GetOk("sort_query_string_for_cache"); ok {
		request["SortQueryStringForCache"] = v
	}
	if v, ok := d.GetOk("include_cookie"); ok {
		request["IncludeCookie"] = v
	}
	if v, ok := d.GetOk("edge_cache_mode"); ok {
		request["EdgeCacheMode"] = v
	}
	if v, ok := d.GetOk("check_presence_header"); ok {
		request["CheckPresenceHeader"] = v
	}
	if v, ok := d.GetOk("edge_cache_ttl"); ok {
		request["EdgeCacheTtl"] = v
	}
	if v, ok := d.GetOk("include_header"); ok {
		request["IncludeHeader"] = v
	}
	if v, ok := d.GetOk("user_geo"); ok {
		request["UserGeo"] = v
	}
	if v, ok := d.GetOk("query_string"); ok {
		request["QueryString"] = v
	}
	if v, ok := d.GetOk("serve_stale"); ok {
		request["ServeStale"] = v
	}
	if v, ok := d.GetOk("bypass_cache"); ok {
		request["BypassCache"] = v
	}
	if v, ok := d.GetOk("cache_deception_armor"); ok {
		request["CacheDeceptionArmor"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("browser_cache_ttl"); ok {
		request["BrowserCacheTtl"] = v
	}
	if v, ok := d.GetOk("cache_reserve_eligibility"); ok {
		request["CacheReserveEligibility"] = v
	}
	if v, ok := d.GetOk("user_device_type"); ok {
		request["UserDeviceType"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	if v, ok := d.GetOk("edge_status_code_cache_ttl"); ok {
		request["EdgeStatusCodeCacheTtl"] = v
	}
	if v, ok := d.GetOk("user_language"); ok {
		request["UserLanguage"] = v
	}
	if v, ok := d.GetOk("browser_cache_mode"); ok {
		request["BrowserCacheMode"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_cache_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaCacheRuleRead(d, meta)
}

func resourceAliCloudEsaCacheRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCacheRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_cache_rule DescribeEsaCacheRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("additional_cacheable_ports", objectRaw["AdditionalCacheablePorts"])
	d.Set("browser_cache_mode", objectRaw["BrowserCacheMode"])
	d.Set("browser_cache_ttl", objectRaw["BrowserCacheTtl"])
	d.Set("bypass_cache", objectRaw["BypassCache"])
	d.Set("cache_deception_armor", objectRaw["CacheDeceptionArmor"])
	d.Set("cache_reserve_eligibility", objectRaw["CacheReserveEligibility"])
	d.Set("check_presence_cookie", objectRaw["CheckPresenceCookie"])
	d.Set("check_presence_header", objectRaw["CheckPresenceHeader"])
	d.Set("edge_cache_mode", objectRaw["EdgeCacheMode"])
	d.Set("edge_cache_ttl", objectRaw["EdgeCacheTtl"])
	d.Set("edge_status_code_cache_ttl", objectRaw["EdgeStatusCodeCacheTtl"])
	d.Set("include_cookie", objectRaw["IncludeCookie"])
	d.Set("include_header", objectRaw["IncludeHeader"])
	d.Set("query_string", objectRaw["QueryString"])
	d.Set("query_string_mode", objectRaw["QueryStringMode"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("serve_stale", objectRaw["ServeStale"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("sort_query_string_for_cache", objectRaw["SortQueryStringForCache"])
	d.Set("user_device_type", objectRaw["UserDeviceType"])
	d.Set("user_geo", objectRaw["UserGeo"])
	d.Set("user_language", objectRaw["UserLanguage"])
	d.Set("cache_rule_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaCacheRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateCacheRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["ConfigId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("check_presence_cookie") {
		update = true
		request["CheckPresenceCookie"] = d.Get("check_presence_cookie")
	}

	if d.HasChange("additional_cacheable_ports") {
		update = true
		request["AdditionalCacheablePorts"] = d.Get("additional_cacheable_ports")
	}

	if d.HasChange("query_string_mode") {
		update = true
		request["QueryStringMode"] = d.Get("query_string_mode")
	}

	if d.HasChange("sort_query_string_for_cache") {
		update = true
		request["SortQueryStringForCache"] = d.Get("sort_query_string_for_cache")
	}

	if d.HasChange("include_cookie") {
		update = true
		request["IncludeCookie"] = d.Get("include_cookie")
	}

	if d.HasChange("edge_cache_mode") {
		update = true
		request["EdgeCacheMode"] = d.Get("edge_cache_mode")
	}

	if d.HasChange("check_presence_header") {
		update = true
		request["CheckPresenceHeader"] = d.Get("check_presence_header")
	}

	if d.HasChange("edge_cache_ttl") {
		update = true
		request["EdgeCacheTtl"] = d.Get("edge_cache_ttl")
	}

	if d.HasChange("include_header") {
		update = true
		request["IncludeHeader"] = d.Get("include_header")
	}

	if d.HasChange("user_geo") {
		update = true
		request["UserGeo"] = d.Get("user_geo")
	}

	if d.HasChange("query_string") {
		update = true
		request["QueryString"] = d.Get("query_string")
	}

	if d.HasChange("serve_stale") {
		update = true
		request["ServeStale"] = d.Get("serve_stale")
	}

	if d.HasChange("bypass_cache") {
		update = true
		request["BypassCache"] = d.Get("bypass_cache")
	}

	if d.HasChange("cache_deception_armor") {
		update = true
		request["CacheDeceptionArmor"] = d.Get("cache_deception_armor")
	}

	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("browser_cache_ttl") {
		update = true
		request["BrowserCacheTtl"] = d.Get("browser_cache_ttl")
	}

	if d.HasChange("cache_reserve_eligibility") {
		update = true
		request["CacheReserveEligibility"] = d.Get("cache_reserve_eligibility")
	}

	if d.HasChange("user_device_type") {
		update = true
		request["UserDeviceType"] = d.Get("user_device_type")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if d.HasChange("edge_status_code_cache_ttl") {
		update = true
		request["EdgeStatusCodeCacheTtl"] = d.Get("edge_status_code_cache_ttl")
	}

	if d.HasChange("user_language") {
		update = true
		request["UserLanguage"] = d.Get("user_language")
	}

	if d.HasChange("browser_cache_mode") {
		update = true
		request["BrowserCacheMode"] = d.Get("browser_cache_mode")
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

	return resourceAliCloudEsaCacheRuleRead(d, meta)
}

func resourceAliCloudEsaCacheRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCacheRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]
	request["RegionId"] = client.RegionId

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
