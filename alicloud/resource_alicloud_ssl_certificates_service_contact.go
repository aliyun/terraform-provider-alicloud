// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServiceContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceContactCreate,
		Read:   resourceAliCloudSslCertificatesServiceContactRead,
		Update: resourceAliCloudSslCertificatesServiceContactUpdate,
		Delete: resourceAliCloudSslCertificatesServiceContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"idcard": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mobile": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"webhook_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceContactCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("webhook_list"); ok {
		if rawUrls, err := json.Marshal(v.([]interface{})); err == nil {
			request["Webhooks"] = string(rawUrls)
		}
	}
	request["Mobile"] = d.Get("mobile")
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	}
	request["Name"] = d.Get("name")
	if v, ok := d.GetOk("idcard"); ok {
		request["Idcard"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_contact", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactId"]))

	return resourceAliCloudSslCertificatesServiceContactRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceContact(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_contact DescribeSslCertificatesServiceContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// email/mobile/idcard are PII that GetContact returns server-masked
	// ("tes****1@example.com", "133******78", "1****************5"). Writing the
	// masked value back to state would diverge from the configured plaintext and
	// cause a perpetual plan diff, so these fields are write-only from the API's
	// perspective: the configured value is authoritative, out-of-band changes are
	// not detectable, and import will not populate them. Do NOT add d.Set for them.
	d.Set("name", objectRaw["Name"])
	// webhook_list: GetContact returns WebhookList as a plaintext string array,
	// directly reconcilable with the configured value.
	if v, ok := objectRaw["WebhookList"].([]interface{}); ok {
		d.Set("webhook_list", v)
	}

	return nil
}

func resourceAliCloudSslCertificatesServiceContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateContact"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ContactId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	// The CAS product requires UpdateContact to carry all mutable parameters together,
	// not only the changed one, so send every writable field on each update.
	if d.HasChange("webhook_list") || d.HasChange("mobile") || d.HasChange("email") || d.HasChange("name") || d.HasChange("idcard") {
		update = true
	}
	if update {
		request["Mobile"] = d.Get("mobile")
		request["Name"] = d.Get("name")
		if v, ok := d.GetOk("email"); ok {
			request["Email"] = v
		}
		if v, ok := d.GetOk("idcard"); ok {
			request["Idcard"] = v
		}
		if v, ok := d.GetOk("webhook_list"); ok {
			if rawUrls, err := json.Marshal(v.([]interface{})); err == nil {
				request["Webhooks"] = string(rawUrls)
			}
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
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

	return resourceAliCloudSslCertificatesServiceContactRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceContactDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ContactId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
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
