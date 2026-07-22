// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServiceCompany() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceCompanyCreate,
		Read:   resourceAliCloudSslCertificatesServiceCompanyRead,
		Update: resourceAliCloudSslCertificatesServiceCompanyUpdate,
		Delete: resourceAliCloudSslCertificatesServiceCompanyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"city": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_phone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"department": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Required: true,
			},
			"post_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"province": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceCompanyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCompany"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["CountryCode"] = d.Get("country_code")
	request["Lang"] = d.Get("lang")
	request["City"] = d.Get("city")
	request["Province"] = d.Get("province")
	request["PostCode"] = d.Get("post_code")
	request["CompanyCode"] = d.Get("company_code")
	request["CompanyName"] = d.Get("company_name")
	if v, ok := d.GetOk("department"); ok {
		request["Department"] = v
	}
	request["CompanyPhone"] = d.Get("company_phone")
	request["CompanyAddress"] = d.Get("company_address")
	request["CompanyType"] = d.Get("company_type")
	if v, ok := d.GetOk("company_email"); ok {
		request["CompanyEmail"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_company", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CompanyId"]))

	return resourceAliCloudSslCertificatesServiceCompanyRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCompanyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceCompany(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_company DescribeSslCertificatesServiceCompany Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("city", objectRaw["City"])
	d.Set("company_address", objectRaw["CompanyAddress"])
	d.Set("company_code", objectRaw["CompanyCode"])
	d.Set("company_email", objectRaw["CompanyEmail"])
	d.Set("company_name", objectRaw["CompanyName"])
	d.Set("company_phone", objectRaw["CompanyPhone"])
	d.Set("company_type", objectRaw["CompanyType"])
	d.Set("country_code", objectRaw["CountryCode"])
	d.Set("department", objectRaw["Department"])
	d.Set("lang", objectRaw["Lang"])
	d.Set("post_code", objectRaw["PostCode"])
	d.Set("province", objectRaw["Province"])

	return nil
}

func resourceAliCloudSslCertificatesServiceCompanyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateCompany"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CompanyId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("country_code") {
		update = true
	}
	request["CountryCode"] = d.Get("country_code")
	if d.HasChange("lang") {
		update = true
	}
	request["Lang"] = d.Get("lang")

	if d.HasChange("city") {
		update = true
	}
	request["City"] = d.Get("city")
	if d.HasChange("province") {
		update = true
	}
	request["Province"] = d.Get("province")
	if d.HasChange("post_code") {
		update = true
	}
	request["PostCode"] = d.Get("post_code")
	if d.HasChange("company_code") {
		update = true
	}
	request["CompanyCode"] = d.Get("company_code")
	if d.HasChange("company_name") {
		update = true
	}
	request["CompanyName"] = d.Get("company_name")
	if d.HasChange("department") {
		update = true
		request["Department"] = d.Get("department")
	}

	if d.HasChange("company_phone") {
		update = true
	}
	request["CompanyPhone"] = d.Get("company_phone")
	if d.HasChange("company_address") {
		update = true
	}
	request["CompanyAddress"] = d.Get("company_address")
	if d.HasChange("company_type") {
		update = true
	}
	request["CompanyType"] = d.Get("company_type")
	if d.HasChange("company_email") {
		update = true
		request["CompanyEmail"] = d.Get("company_email")
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

	return resourceAliCloudSslCertificatesServiceCompanyRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCompanyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCompany"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["CompanyId"] = d.Id()

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
