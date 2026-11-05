// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package esa

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEsaClientCaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaClientCaCertificateCreate,
		Read:   resourceAliCloudEsaClientCaCertificateRead,
		Update: resourceAliCloudEsaClientCaCertificateUpdate,
		Delete: resourceAliCloudEsaClientCaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_ca_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_ca_cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_ca_certificate_hostnames": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
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

func resourceAliCloudEsaClientCaCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UploadClientCaCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("client_ca_cert_name"); ok {
		request["Name"] = v
	}
	request["Certificate"] = d.Get("certificate")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_esa_client_ca_certificate", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	return resourceAliCloudEsaClientCaCertificateUpdate(d, meta)
}

func resourceAliCloudEsaClientCaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaClientCaCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_client_ca_certificate DescribeEsaClientCaCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	d.Set("client_ca_cert_name", resultRaw["Name"])
	d.Set("create_time", resultRaw["CreateTime"])
	d.Set("status", resultRaw["Status"])
	d.Set("client_ca_cert_id", resultRaw["Id"])

	objectRaw, err = esaServiceV2.DescribeClientCaCertificateGetClientCaCertificateHostnames(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	hostnamesRaw := make([]interface{}, 0)
	if objectRaw["Hostnames"] != nil {
		hostnamesRaw = convert.ConvertToInterfaceArray(objectRaw["Hostnames"])
	}

	d.Set("client_ca_certificate_hostnames", hostnamesRaw)

	return nil
}

func resourceAliCloudEsaClientCaCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "SetClientCaCertificateHostnames"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	if d.HasChange("client_ca_certificate_hostnames") {
		update = true
	}
	if v, ok := d.GetOk("client_ca_certificate_hostnames"); ok || d.HasChange("client_ca_certificate_hostnames") {
		hostnamesMapsArray := convert.ConvertToInterfaceArray(v)

		request["Hostnames"] = convert.ConvertListToJsonString(hostnamesMapsArray)
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEsaClientCaCertificateRead(d, meta)
}

func resourceAliCloudEsaClientCaCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteClientCaCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SiteId"] = parts[0]
	query["Id"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
