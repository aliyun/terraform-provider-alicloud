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

func resourceAliCloudEsaHttpsApplicationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaHttpsApplicationConfigurationCreate,
		Read:   resourceAliCloudEsaHttpsApplicationConfigurationRead,
		Update: resourceAliCloudEsaHttpsApplicationConfigurationUpdate,
		Delete: resourceAliCloudEsaHttpsApplicationConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alt_svc": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"alt_svc_clear": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"alt_svc_ma": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alt_svc_persist": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hsts": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"hsts_include_subdomains": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"hsts_max_age": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hsts_preload": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"https_force": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"https_force_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"301", "302", "307", "308"}, false),
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

func resourceAliCloudEsaHttpsApplicationConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHttpsApplicationConfiguration"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("https_force"); ok {
		request["HttpsForce"] = v
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
	if v, ok := d.GetOk("hsts"); ok {
		request["Hsts"] = v
	}
	if v, ok := d.GetOk("https_force_code"); ok {
		request["HttpsForceCode"] = v
	}
	if v, ok := d.GetOk("hsts_max_age"); ok {
		request["HstsMaxAge"] = v
	}
	if v, ok := d.GetOk("hsts_include_subdomains"); ok {
		request["HstsIncludeSubdomains"] = v
	}
	if v, ok := d.GetOk("alt_svc_clear"); ok {
		request["AltSvcClear"] = v
	}
	if v, ok := d.GetOk("alt_svc_persist"); ok {
		request["AltSvcPersist"] = v
	}
	if v, ok := d.GetOk("alt_svc"); ok {
		request["AltSvc"] = v
	}
	if v, ok := d.GetOk("hsts_preload"); ok {
		request["HstsPreload"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	if v, ok := d.GetOk("alt_svc_ma"); ok {
		request["AltSvcMa"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_https_application_configuration", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaHttpsApplicationConfigurationUpdate(d, meta)
}

func resourceAliCloudEsaHttpsApplicationConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaHttpsApplicationConfiguration(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_https_application_configuration DescribeEsaHttpsApplicationConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alt_svc", objectRaw["AltSvc"])
	d.Set("alt_svc_clear", objectRaw["AltSvcClear"])
	d.Set("alt_svc_ma", objectRaw["AltSvcMa"])
	d.Set("alt_svc_persist", objectRaw["AltSvcPersist"])
	d.Set("hsts", objectRaw["Hsts"])
	d.Set("hsts_include_subdomains", objectRaw["HstsIncludeSubdomains"])
	d.Set("hsts_max_age", objectRaw["HstsMaxAge"])
	d.Set("hsts_preload", objectRaw["HstsPreload"])
	d.Set("https_force", objectRaw["HttpsForce"])
	d.Set("https_force_code", objectRaw["HttpsForceCode"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("config_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaHttpsApplicationConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateHttpsApplicationConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("https_force") {
		update = true
		request["HttpsForce"] = d.Get("https_force")
	}

	if !d.IsNewResource() && d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if !d.IsNewResource() && d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if !d.IsNewResource() && d.HasChange("hsts") {
		update = true
		request["Hsts"] = d.Get("hsts")
	}

	if !d.IsNewResource() && d.HasChange("https_force_code") {
		update = true
		request["HttpsForceCode"] = d.Get("https_force_code")
	}

	if !d.IsNewResource() && d.HasChange("hsts_max_age") {
		update = true
		request["HstsMaxAge"] = d.Get("hsts_max_age")
	}

	if !d.IsNewResource() && d.HasChange("hsts_include_subdomains") {
		update = true
		request["HstsIncludeSubdomains"] = d.Get("hsts_include_subdomains")
	}

	if !d.IsNewResource() && d.HasChange("alt_svc_clear") {
		update = true
		request["AltSvcClear"] = d.Get("alt_svc_clear")
	}

	if !d.IsNewResource() && d.HasChange("alt_svc_persist") {
		update = true
		request["AltSvcPersist"] = d.Get("alt_svc_persist")
	}

	if !d.IsNewResource() && d.HasChange("alt_svc") {
		update = true
		request["AltSvc"] = d.Get("alt_svc")
	}

	if !d.IsNewResource() && d.HasChange("hsts_preload") {
		update = true
		request["HstsPreload"] = d.Get("hsts_preload")
	}

	if !d.IsNewResource() && d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if !d.IsNewResource() && d.HasChange("alt_svc_ma") {
		update = true
		request["AltSvcMa"] = d.Get("alt_svc_ma")
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

	return resourceAliCloudEsaHttpsApplicationConfigurationRead(d, meta)
}

func resourceAliCloudEsaHttpsApplicationConfigurationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteHttpsApplicationConfiguration"
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
