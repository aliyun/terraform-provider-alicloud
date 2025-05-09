// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

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
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	if v, ok := d.GetOk("proxy_types"); ok {
		proxyTypesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ProxyPorts"] = dataLoopTmp["proxy_ports"].(*schema.Set).List()
			dataLoopMap["ProxyType"] = dataLoopTmp["proxy_type"]
			proxyTypesMapsArray = append(proxyTypesMapsArray, dataLoopMap)
		}
		request["ProxyTypes"] = proxyTypesMapsArray
	}

	request["RsType"] = d.Get("rs_type")
	if v, ok := d.GetOk("real_servers"); ok {
		realServersMapsArray := v.(*schema.Set).List()
		request["RealServers"] = realServersMapsArray
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsMapsArray := v.(*schema.Set).List()
		request["InstanceIds"] = instanceIdsMapsArray
	}

	if v, ok := d.GetOk("https_ext"); ok {
		request["HttpsExt"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_domain_resource", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Domain"]))

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

	d.Set("cname", objectRaw["Cname"])
	d.Set("https_ext", objectRaw["HttpsExt"])
	d.Set("ocsp_enabled", formatBool(convertDdosCooDomainResourceWebRulesOcspEnabledResponse(objectRaw["OcspEnabled"])))
	d.Set("rs_type", objectRaw["RsType"])
	d.Set("domain", objectRaw["Domain"])

	instanceIdsRaw := make([]interface{}, 0)
	if objectRaw["InstanceIds"] != nil {
		instanceIdsRaw = objectRaw["InstanceIds"].([]interface{})
	}

	d.Set("instance_ids", instanceIdsRaw)
	proxyTypesRaw := objectRaw["ProxyTypes"]
	proxyTypesMaps := make([]map[string]interface{}, 0)
	if proxyTypesRaw != nil {
		for _, proxyTypesChildRaw := range proxyTypesRaw.([]interface{}) {
			proxyTypesMap := make(map[string]interface{})
			proxyTypesChildRaw := proxyTypesChildRaw.(map[string]interface{})
			proxyTypesMap["proxy_type"] = proxyTypesChildRaw["ProxyType"]

			proxyPortsRaw := make([]interface{}, 0)
			if proxyTypesChildRaw["ProxyPorts"] != nil {
				proxyPortsRaw = proxyTypesChildRaw["ProxyPorts"].([]interface{})
			}

			proxyTypesMap["proxy_ports"] = proxyPortsRaw
			proxyTypesMaps = append(proxyTypesMaps, proxyTypesMap)
		}
	}
	if err := d.Set("proxy_types", proxyTypesMaps); err != nil {
		return err
	}
	realServersRaw := make([]interface{}, 0)
	if objectRaw["RealServers"] != nil {
		realServersRaw = objectRaw["RealServers"].([]interface{})
	}

	d.Set("real_servers", realServersRaw)

	objectRaw, err = ddosCooServiceV2.DescribeDomainResourceDescribeWebRules(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("cert_name", objectRaw["UserCertName"])
	d.Set("domain", objectRaw["Domain"])

	return nil
}

func resourceAliCloudDdosCooDomainResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyDomainResource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()

	if !d.IsNewResource() && d.HasChange("real_servers") {
		update = true
	}
	if v, ok := d.GetOk("real_servers"); ok {
		realServersMapsArray := v.(*schema.Set).List()
		request["RealServers"] = realServersMapsArray
	}

	if !d.IsNewResource() && d.HasChange("rs_type") {
		update = true
	}
	request["RsType"] = d.Get("rs_type")
	if !d.IsNewResource() && d.HasChange("proxy_types") {
		update = true
	}
	if v, ok := d.GetOk("proxy_types"); ok {
		proxyTypesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ProxyType"] = dataLoop1Tmp["proxy_type"]
			dataLoop1Map["ProxyPorts"] = dataLoop1Tmp["proxy_ports"].(*schema.Set).List()
			proxyTypesMapsArray = append(proxyTypesMapsArray, dataLoop1Map)
		}
		request["ProxyTypes"] = proxyTypesMapsArray
	}

	if !d.IsNewResource() && d.HasChange("instance_ids") {
		update = true
	}
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsMapsArray := v.(*schema.Set).List()
		request["InstanceIds"] = instanceIdsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("https_ext") {
		update = true
		request["HttpsExt"] = d.Get("https_ext")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
	update = false
	action = "ModifyOcspStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()

	if d.HasChange("ocsp_enabled") {
		update = true
	}
	request["Enable"] = convertDdosCooDomainResourceEnableRequest(d.Get("ocsp_enabled").(bool))
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
		request["CertIdentifier"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudDdosCooDomainResourceRead(d, meta)
}

func resourceAliCloudDdosCooDomainResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDomainResource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Domain"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
