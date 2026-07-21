// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSslCertificatesServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceCertificateCreate,
		Read:   resourceAliCloudSslCertificatesServiceCertificateRead,
		Update: resourceAliCloudSslCertificatesServiceCertificateUpdate,
		Delete: resourceAliCloudSslCertificatesServiceCertificateDelete,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"certificate_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"certificate_name", "name"},
			},
			"encrypt_cert": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encrypt_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sign_cert": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sign_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"tags": tagsSchema(),
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated from provider version 1.129.0 and it will be removed in the future version. Please use the new attribute 'certificate_name' instead.",
			},
			"lang": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'lang' has been deprecated from provider version 1.260.1 and it will be removed in the future version.",
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UploadUserCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMapWithTags(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("key"); ok {
		request["Key"] = v
	}
	if v, ok := d.GetOk("encrypt_cert"); ok {
		request["EncryptCert"] = v
	}
	if v, ok := d.GetOk("sign_cert"); ok {
		request["SignCert"] = v
	}
	if v, ok := d.GetOk("certificate_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("encrypt_private_key"); ok {
		request["EncryptPrivateKey"] = v
	}
	if v, ok := d.GetOk("sign_private_key"); ok {
		request["SignPrivateKey"] = v
	}
	if v, ok := d.GetOk("cert"); ok {
		request["Cert"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CertId"]))

	return resourceAliCloudSslCertificatesServiceCertificateRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_certificate DescribeSslCertificatesServiceCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cert", objectRaw["Cert"])
	d.Set("certificate_name", objectRaw["Name"])
	d.Set("encrypt_cert", objectRaw["EncryptCert"])
	d.Set("encrypt_private_key", objectRaw["EncryptPrivateKey"])
	d.Set("key", objectRaw["Key"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("sign_cert", objectRaw["SignCert"])
	d.Set("sign_private_key", objectRaw["SignPrivateKey"])
	d.Set("name", objectRaw["Name"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudSslCertificatesServiceCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "CreateCertificateName"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CertificateId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			request = expandTagsToMapWithTags(request, tagsMap)
		} else {
			request["Tags.1.Key"] = ""
			request["Tags.1.Value"] = ""
		}
	}

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("certificate_name") {
		update = true
		request["Name"] = d.Get("certificate_name")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2018-08-13", action, query, request, true)
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

	return resourceAliCloudSslCertificatesServiceCertificateRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteUserCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["CertId"] = d.Id()

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
