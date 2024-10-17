// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudDdosCooDomainResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosCooDomainResourceCreate,
		Read:   resourceAliCloudDdosCooDomainResourceRead,
		Update: resourceAliCloudDdosCooDomainResourceUpdate,
		Delete: resourceAliCloudDdosCooDomainResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cert": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cert_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_region": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"https_ext": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ocsp_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"proxy_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"real_servers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rs_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAliCloudDdosCooDomainResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDomainResource"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		query["Domain"] = v
	}

	if v, ok := d.GetOk("proxy_types"); ok {
		proxyTypesMaps := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ProxyPorts"] = dataLoopTmp["proxy_ports"].(*schema.Set).List()
			dataLoopMap["ProxyType"] = dataLoopTmp["proxy_type"]
			proxyTypesMaps = append(proxyTypesMaps, dataLoopMap)
		}
		request["ProxyTypes"] = proxyTypesMaps
	}

	request["RsType"] = d.Get("rs_type")
	if v, ok := d.GetOk("real_servers"); ok {
		realServersMaps := v.(*schema.Set).List()
		request["RealServers"] = realServersMaps
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsMaps := v.(*schema.Set).List()
		request["InstanceIds"] = instanceIdsMaps
	}

	if v, ok := d.GetOk("https_ext"); ok {
		request["HttpsExt"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_domain_resource", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(query["Domain"]))

	return resourceAliCloudDdosCooDomainResourceUpdate(d, meta)
}

func resourceAliCloudDdosCooDomainResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosCooServiceV2 := DdosCooServiceV2{client}

	objectRaw, err := ddosCooServiceV2.DescribeDdosCooDomainResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_domain_resource DescribeDdosCooDomainResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Cname"] != nil {
		d.Set("cname", objectRaw["Cname"])
	}
	if objectRaw["HttpsExt"] != nil {
		d.Set("https_ext", objectRaw["HttpsExt"])
	}
	if objectRaw["OcspEnabled"] != nil {
		d.Set("ocsp_enabled", formatBool(convertDdosCooDomainResourceWebRulesOcspEnabledResponse(objectRaw["OcspEnabled"])))
	}
	if objectRaw["RsType"] != nil {
		d.Set("rs_type", objectRaw["RsType"])
	}
	if objectRaw["Domain"] != nil {
		d.Set("domain", objectRaw["Domain"])
	}

	instanceIds1Raw := make([]interface{}, 0)
	if objectRaw["InstanceIds"] != nil {
		instanceIds1Raw = objectRaw["InstanceIds"].([]interface{})
	}

	d.Set("instance_ids", instanceIds1Raw)
	proxyTypes1Raw := objectRaw["ProxyTypes"]
	proxyTypesMaps := make([]map[string]interface{}, 0)
	if proxyTypes1Raw != nil {
		for _, proxyTypesChild1Raw := range proxyTypes1Raw.([]interface{}) {
			proxyTypesMap := make(map[string]interface{})
			proxyTypesChild1Raw := proxyTypesChild1Raw.(map[string]interface{})
			proxyTypesMap["proxy_type"] = proxyTypesChild1Raw["ProxyType"]

			proxyPorts1Raw := make([]interface{}, 0)
			if proxyTypesChild1Raw["ProxyPorts"] != nil {
				proxyPorts1Raw = proxyTypesChild1Raw["ProxyPorts"].([]interface{})
			}

			proxyTypesMap["proxy_ports"] = proxyPorts1Raw
			proxyTypesMaps = append(proxyTypesMaps, proxyTypesMap)
		}
	}
	if objectRaw["ProxyTypes"] != nil {
		if err := d.Set("proxy_types", proxyTypesMaps); err != nil {
			return err
		}
	}
	realServers1Raw := make([]interface{}, 0)
	if objectRaw["RealServers"] != nil {
		realServers1Raw = objectRaw["RealServers"].([]interface{})
	}

	d.Set("real_servers", realServers1Raw)

	objectRaw, err = ddosCooServiceV2.DescribeDescribeWebRules(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["UserCertName"] != nil {
		d.Set("cert_name", objectRaw["UserCertName"])
	}
	if objectRaw["Domain"] != nil {
		d.Set("domain", objectRaw["Domain"])
	}

	d.Set("domain", d.Id())

	return nil
}

func resourceAliCloudDdosCooDomainResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyDomainResource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Domain"] = d.Id()

	if !d.IsNewResource() && d.HasChange("real_servers") {
		update = true
	}
	if v, ok := d.GetOk("real_servers"); ok || d.HasChange("real_servers") {
		realServersMaps := v.(*schema.Set).List()
		request["RealServers"] = realServersMaps
	}

	if !d.IsNewResource() && d.HasChange("rs_type") {
		update = true
	}
	request["RsType"] = d.Get("rs_type")
	if !d.IsNewResource() && d.HasChange("proxy_types") {
		update = true
	}
	if v, ok := d.GetOk("proxy_types"); ok || d.HasChange("proxy_types") {
		proxyTypesMaps := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ProxyType"] = dataLoop1Tmp["proxy_type"]
			dataLoop1Map["ProxyPorts"] = dataLoop1Tmp["proxy_ports"].(*schema.Set).List()
			proxyTypesMaps = append(proxyTypesMaps, dataLoop1Map)
		}
		request["ProxyTypes"] = proxyTypesMaps
	}

	if !d.IsNewResource() && d.HasChange("instance_ids") {
		update = true
	}
	if v, ok := d.GetOk("instance_ids"); ok || d.HasChange("instance_ids") {
		instanceIdsMaps := v.(*schema.Set).List()
		request["InstanceIds"] = instanceIdsMaps
	}

	if !d.IsNewResource() && d.HasChange("https_ext") {
		update = true
		request["HttpsExt"] = d.Get("https_ext")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyOcspStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Domain"] = d.Id()

	if d.HasChange("ocsp_enabled") {
		update = true
	}
	request["Enable"] = convertDdosCooDomainResourceEnableRequest(d.Get("ocsp_enabled").(bool))
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "AssociateWebCert"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()

	if d.HasChange("cert_name") {
		update = true
		request["CertName"] = d.Get("cert_name")
	}

	if d.HasChange("cert") {
		update = true
		request["Cert"] = d.Get("cert")
	}

	if d.HasChange("key") {
		update = true
		request["Key"] = d.Get("key")
	}

	if d.HasChange("cert_region") {
		update = true
		request["CertRegion"] = d.Get("cert_region")
	}

	if v, ok := d.GetOk("cert_identifier"); ok {
		request["CertId"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudDdosCooDomainResourceRead(d, meta)
}

func resourceAliCloudDdosCooDomainResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDomainResource"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	query["Domain"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertDdosCooDomainResourceWebRulesOcspEnabledResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
func convertDdosCooDomainResourceEnableRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return "1"
	case "false":
		return "0"
	}
	return source
}
