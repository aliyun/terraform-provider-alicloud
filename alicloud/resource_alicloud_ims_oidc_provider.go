// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImsOidcProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImsOidcProviderCreate,
		Read:   resourceAliCloudImsOidcProviderRead,
		Update: resourceAliCloudImsOidcProviderUpdate,
		Delete: resourceAliCloudImsOidcProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"client_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fingerprints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"issuance_limit_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 168),
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oidc_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudImsOidcProviderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOIDCProvider"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["OIDCProviderName"] = d.Get("oidc_provider_name")

	request["IssuerUrl"] = d.Get("issuer_url")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("issuance_limit_time"); ok {
		request["IssuanceLimitTime"] = v
	}
	if v, ok := d.GetOk("client_ids"); ok {
		request["ClientIds"] = convertArrayToString(v.(*schema.Set).List(), ",")
	}
	if v, ok := d.GetOk("fingerprints"); ok {
		request["Fingerprints"] = convertArrayToString(v.(*schema.Set).List(), ",")
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ims_oidc_provider", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.OIDCProvider.OIDCProviderName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudImsOidcProviderUpdate(d, meta)
}

func resourceAliCloudImsOidcProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	imsServiceV2 := ImsServiceV2{client}

	objectRaw, err := imsServiceV2.DescribeImsOidcProvider(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ims_oidc_provider DescribeImsOidcProvider Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("issuance_limit_time", objectRaw["IssuanceLimitTime"])
	d.Set("issuer_url", objectRaw["IssuerUrl"])
	d.Set("oidc_provider_name", objectRaw["OIDCProviderName"])
	d.Set("arn", objectRaw["Arn"])

	e := jsonata.MustCompile("$split($.ClientIds, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("client_ids", evaluation)
	e = jsonata.MustCompile("$split($.Fingerprints, \",\")")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("fingerprints", evaluation)

	return nil
}

func resourceAliCloudImsOidcProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateOIDCProvider"
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["OIDCProviderName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("issuance_limit_time") {
		update = true
		request["IssuanceLimitTime"] = d.Get("issuance_limit_time")
	}

	if !d.IsNewResource() && d.HasChange("client_ids") {
		update = true
		request["ClientIds"] = convertArrayToString(d.Get("client_ids").(*schema.Set).List(), ",")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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

	if !d.IsNewResource() && d.HasChange("fingerprints") {
		oldEntry, newEntry := d.GetChange("fingerprints")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if added.Len() > 0 {
			fingerprints := added.List()

			for _, item := range fingerprints {
				action := "AddFingerprintToOIDCProvider"
				conn, err := client.NewImsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["OIDCProviderName"] = d.Id()
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["Fingerprint"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
			d.SetPartial("fingerprints")
		}

		if removed.Len() > 0 {
			fingerprints := removed.List()

			for _, item := range fingerprints {
				action := "RemoveFingerprintFromOIDCProvider"
				conn, err := client.NewImsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["OIDCProviderName"] = d.Id()
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["Fingerprint"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
			d.SetPartial("fingerprints")
		}

	}
	return resourceAliCloudImsOidcProviderRead(d, meta)
}

func resourceAliCloudImsOidcProviderDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteOIDCProvider"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["OIDCProviderName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		if IsExpectedErrors(err, []string{"EntityNotExist.OIDCProvider"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
