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
				Type:     schema.TypeString,
				Optional: true,
			},
			"idcard": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"webhooks": {
				Type:     schema.TypeString,
				Optional: true,
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

	if v, ok := d.GetOk("webhooks"); ok {
		request["Webhooks"] = v
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

	// email, mobile and webhooks are sensitive fields that GetContact/ListContact return
	// masked ("t****@example.com", "133******78"), empty, or as an empty JSON array ("[]")
	// when unset. These cannot be reconciled against the configured value, so only
	// overwrite state when the API returns a real (non-empty, non-masked, non-"[]") value.
	// This keeps the resource idempotent while still allowing genuine updates to surface.
	if v, ok := objectRaw["Email"].(string); ok && !sslContactSensitiveFieldMasked(v) {
		d.Set("email", v)
	}
	if v, ok := objectRaw["Mobile"].(string); ok && !sslContactSensitiveFieldMasked(v) {
		d.Set("mobile", v)
	}
	d.Set("name", objectRaw["Name"])
	if v, ok := objectRaw["Webhooks"].(string); ok && !sslContactSensitiveFieldMasked(v) {
		d.Set("webhooks", v)
	}

	return nil
}

// sslContactSensitiveFieldMasked reports whether the value returned by GetContact/ListContact
// for a sensitive field (email/mobile/webhooks) is not a real, reconcilable value — i.e. it is
// empty, an empty JSON array ("[]"), or server-masked (contains "*"). Such values must not be
// written into state, otherwise every plan would report a non-empty diff against the configured
// value.
func sslContactSensitiveFieldMasked(v string) bool {
	return v == "" || v == "[]" || strings.Contains(v, "*")
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
	if d.HasChange("webhooks") {
		update = true
		request["Webhooks"] = d.Get("webhooks")
	}

	if d.HasChange("mobile") {
		update = true
		request["Mobile"] = d.Get("mobile")
	}

	if d.HasChange("email") {
		update = true
		request["Email"] = d.Get("email")
	}

	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if v, ok := d.GetOk("idcard"); ok {
		request["Idcard"] = v
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
