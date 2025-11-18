// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaSiteOriginClientCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaSiteOriginClientCertificateCreate,
		Read:   resourceAliCloudEsaSiteOriginClientCertificateRead,
		Delete: resourceAliCloudEsaSiteOriginClientCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_origin_client_certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_origin_client_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaSiteOriginClientCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UploadSiteOriginClientCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	request["Certificate"] = d.Get("certificate")
	request["PrivateKey"] = d.Get("private_key")
	if v, ok := d.GetOk("site_origin_client_certificate_name"); ok {
		request["Name"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_site_origin_client_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	return resourceAliCloudEsaSiteOriginClientCertificateRead(d, meta)
}

func resourceAliCloudEsaSiteOriginClientCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaSiteOriginClientCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_site_origin_client_certificate DescribeEsaSiteOriginClientCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("certificate", objectRaw["Certificate"])
	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", v)
	}

	resultRawObj, _ := jsonpath.Get("$.Result", objectRaw)
	resultRaw := make(map[string]interface{})
	if resultRawObj != nil {
		resultRaw = resultRawObj.(map[string]interface{})
	}
	d.Set("create_time", resultRaw["CreateTime"])
	d.Set("site_origin_client_certificate_name", resultRaw["Name"])
	d.Set("status", resultRaw["Status"])
	d.Set("site_origin_client_certificate_id", resultRaw["Id"])

	return nil
}

func resourceAliCloudEsaSiteOriginClientCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteSiteOriginClientCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SiteId"] = parts[0]
	query["Id"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
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
