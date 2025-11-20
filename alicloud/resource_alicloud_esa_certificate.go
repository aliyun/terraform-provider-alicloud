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

func resourceAliCloudEsaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCertificateCreate,
		Read:   resourceAliCloudEsaCertificateRead,
		Update: resourceAliCloudEsaCertificateUpdate,
		Delete: resourceAliCloudEsaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cas_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"cas", "upload", "free"}, false),
			},
			"domains": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"lets_encrypt"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("created_type"); ok && v.(string) == "free" {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceAliCloudEsaCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("created_type"); ok && InArray(fmt.Sprint(v), []string{"cas", "upload"}) {
		action := "SetCertificate"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("site_id"); ok {
			request["SiteId"] = v
		}
		if v, ok := d.GetOk("cert_id"); ok {
			request["Id"] = v
		}

		if v, ok := d.GetOk("private_key"); ok {
			request["PrivateKey"] = v
		}
		request["Type"] = d.Get("created_type")
		if v, ok := d.GetOk("region"); ok {
			request["Region"] = v
		}
		if v, ok := d.GetOk("cas_id"); ok {
			request["CasId"] = v
		}
		if v, ok := d.GetOk("certificate"); ok {
			request["Certificate"] = v
		}
		if v, ok := d.GetOk("cert_name"); ok {
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_certificate", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	}

	if v, ok := d.GetOk("created_type"); ok && InArray(fmt.Sprint(v), []string{"free"}) {
		action := "ApplyCertificate"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("site_id"); ok {
			query["SiteId"] = v
		}

		if v, ok := d.GetOk("domains"); ok {
			query["Domains"] = v.(string)
		}

		if v, ok := d.GetOk("type"); ok {
			query["Type"] = v.(string)
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_certificate", action, AlibabaCloudSdkGoERROR)
		}

		ResultIdVar, _ := jsonpath.Get("$.Result[0].Id", response)
		d.SetId(fmt.Sprintf("%v:%v", query["SiteId"], ResultIdVar))

	}

	return resourceAliCloudEsaCertificateRead(d, meta)
}

func resourceAliCloudEsaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_certificate DescribeEsaCertificate Failed!!! %s", err)
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
	d.Set("cas_id", resultRaw["CasId"])
	d.Set("cert_name", resultRaw["Name"])
	d.Set("create_time", resultRaw["CreateTime"])
	d.Set("domains", resultRaw["SAN"])
	d.Set("region", resultRaw["Region"])
	d.Set("status", resultRaw["Status"])
	d.Set("type", resultRaw["Type"])
	d.Set("cert_id", resultRaw["Id"])

	return nil
}

func resourceAliCloudEsaCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "SetCertificate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	if v, ok := d.GetOk("private_key"); ok {
		request["PrivateKey"] = v
	}
	request["Type"] = d.Get("created_type")
	if _, ok := d.GetOk("region"); ok || d.HasChange("region") {
		update = true
		request["Region"] = d.Get("region")
	}

	if _, ok := d.GetOk("cas_id"); ok || d.HasChange("cas_id") {
		update = true
		request["CasId"] = d.Get("cas_id")
	}

	if _, ok := d.GetOk("certificate"); ok || d.HasChange("certificate") {
		update = true
		request["Certificate"] = d.Get("certificate")
	}

	if _, ok := d.GetOk("cert_name"); ok || d.HasChange("cert_name") {
		update = true
		request["Name"] = d.Get("cert_name")
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

	return resourceAliCloudEsaCertificateRead(d, meta)
}

func resourceAliCloudEsaCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCertificate"
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
