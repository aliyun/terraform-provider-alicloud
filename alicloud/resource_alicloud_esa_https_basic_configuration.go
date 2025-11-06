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

func resourceAliCloudEsaHttpsBasicConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaHttpsBasicConfigurationCreate,
		Read:   resourceAliCloudEsaHttpsBasicConfigurationRead,
		Update: resourceAliCloudEsaHttpsBasicConfigurationUpdate,
		Delete: resourceAliCloudEsaHttpsBasicConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ciphersuite": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ciphersuite_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"all", "strict", "custom"}, false),
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http2": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"http3": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"https": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"ocsp_stapling": {
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
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"tls10": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"tls11": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"tls12": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"tls13": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
		},
	}
}

func resourceAliCloudEsaHttpsBasicConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHttpsBasicConfiguration"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("tls11"); ok {
		request["Tls11"] = v
	}
	if v, ok := d.GetOk("tls10"); ok {
		request["Tls10"] = v
	}
	if v, ok := d.GetOkExists("sequence"); ok {
		request["Sequence"] = v
	}
	if v, ok := d.GetOk("tls13"); ok {
		request["Tls13"] = v
	}
	if v, ok := d.GetOk("tls12"); ok {
		request["Tls12"] = v
	}
	if v, ok := d.GetOk("ciphersuite"); ok {
		request["Ciphersuite"] = v
	}
	if v, ok := d.GetOk("ocsp_stapling"); ok {
		request["OcspStapling"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("ciphersuite_group"); ok {
		request["CiphersuiteGroup"] = v
	}
	if v, ok := d.GetOk("http2"); ok {
		request["Http2"] = v
	}
	if v, ok := d.GetOk("https"); ok {
		request["Https"] = v
	}
	if v, ok := d.GetOk("http3"); ok {
		request["Http3"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_https_basic_configuration", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaHttpsBasicConfigurationRead(d, meta)
}

func resourceAliCloudEsaHttpsBasicConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaHttpsBasicConfiguration(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_https_basic_configuration DescribeEsaHttpsBasicConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ciphersuite", objectRaw["Ciphersuite"])
	d.Set("ciphersuite_group", objectRaw["CiphersuiteGroup"])
	d.Set("http2", objectRaw["Http2"])
	d.Set("http3", objectRaw["Http3"])
	d.Set("https", objectRaw["Https"])
	d.Set("ocsp_stapling", objectRaw["OcspStapling"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("tls10", objectRaw["Tls10"])
	d.Set("tls11", objectRaw["Tls11"])
	d.Set("tls12", objectRaw["Tls12"])
	d.Set("tls13", objectRaw["Tls13"])
	d.Set("config_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaHttpsBasicConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateHttpsBasicConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("tls11") {
		update = true
		request["Tls11"] = d.Get("tls11")
	}

	if d.HasChange("tls10") {
		update = true
		request["Tls10"] = d.Get("tls10")
	}

	if d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if d.HasChange("tls13") {
		update = true
		request["Tls13"] = d.Get("tls13")
	}

	if d.HasChange("tls12") {
		update = true
		request["Tls12"] = d.Get("tls12")
	}

	if d.HasChange("ciphersuite") {
		update = true
		request["Ciphersuite"] = d.Get("ciphersuite")
	}

	if d.HasChange("ocsp_stapling") {
		update = true
		request["OcspStapling"] = d.Get("ocsp_stapling")
	}

	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("ciphersuite_group") {
		update = true
		request["CiphersuiteGroup"] = d.Get("ciphersuite_group")
	}

	if d.HasChange("http2") {
		update = true
		request["Http2"] = d.Get("http2")
	}

	if d.HasChange("https") {
		update = true
		request["Https"] = d.Get("https")
	}

	if d.HasChange("http3") {
		update = true
		request["Http3"] = d.Get("http3")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
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

	return resourceAliCloudEsaHttpsBasicConfigurationRead(d, meta)
}

func resourceAliCloudEsaHttpsBasicConfigurationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteHttpsBasicConfiguration"
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
