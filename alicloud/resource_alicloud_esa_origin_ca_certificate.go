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

func resourceAliCloudEsaOriginCaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaOriginCaCertificateCreate,
		Read:   resourceAliCloudEsaOriginCaCertificateRead,
		Delete: resourceAliCloudEsaOriginCaCertificateDelete,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"origin_ca_certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaOriginCaCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UploadOriginCaCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	request["RegionId"] = client.RegionId

	request["Certificate"] = d.Get("certificate")
	if v, ok := d.GetOk("name"); ok {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_origin_ca_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	return resourceAliCloudEsaOriginCaCertificateRead(d, meta)
}

func resourceAliCloudEsaOriginCaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaOriginCaCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_origin_ca_certificate DescribeEsaOriginCaCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("certificate", objectRaw["Certificate"])
	d.Set("site_id", objectRaw["SiteId"])

	resultRawObj, _ := jsonpath.Get("$.Result", objectRaw)
	resultRaw := make(map[string]interface{})
	if resultRawObj != nil {
		resultRaw = resultRawObj.(map[string]interface{})
	}
	d.Set("create_time", resultRaw["CreateTime"])
	d.Set("name", resultRaw["Name"])
	d.Set("status", resultRaw["Status"])
	d.Set("origin_ca_certificate_id", resultRaw["Id"])

	return nil
}

func resourceAliCloudEsaOriginCaCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteOriginCaCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SiteId"] = parts[0]
	query["Id"] = parts[1]
	query["RegionId"] = client.RegionId

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
