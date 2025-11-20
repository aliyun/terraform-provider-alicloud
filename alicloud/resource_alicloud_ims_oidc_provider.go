package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
				ValidateFunc: IntBetween(0, 168),
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
		},
	}
}

func resourceAliCloudImsOidcProviderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOIDCProvider"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("oidc_provider_name"); ok {
		request["OIDCProviderName"] = v
	}

	request["IssuerUrl"] = d.Get("issuer_url")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("issuance_limit_time"); ok && v.(int) > 0 {
		request["IssuanceLimitTime"] = v
	}
	if v, ok := d.GetOk("client_ids"); ok {
		jsonPathResult3, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult3 != "" {
			request["ClientIds"] = convertListToCommaSeparate(jsonPathResult3.(*schema.Set).List())
		}
	}
	if v, ok := d.GetOk("fingerprints"); ok {
		jsonPathResult4, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult4 != "" {
			request["Fingerprints"] = convertListToCommaSeparate(jsonPathResult4.(*schema.Set).List())
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ims_oidc_provider", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.OIDCProvider.OIDCProviderName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudImsOidcProviderRead(d, meta)
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

	d.Set("arn", objectRaw["Arn"])
	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("issuance_limit_time", objectRaw["IssuanceLimitTime"])
	d.Set("issuer_url", objectRaw["IssuerUrl"])
	d.Set("oidc_provider_name", objectRaw["OIDCProviderName"])

	e := jsonata.MustCompile("$split($.ClientIds, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("client_ids", evaluation)
	e = jsonata.MustCompile("$split($.Fingerprints, \",\")")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("fingerprints", evaluation)

	d.Set("oidc_provider_name", d.Id())

	return nil
}

func resourceAliCloudImsOidcProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateOIDCProvider"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["OIDCProviderName"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if d.HasChange("issuance_limit_time") {
		update = true
		request["IssuanceLimitTime"] = d.Get("issuance_limit_time")
	}

	if d.HasChange("client_ids") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$", d.Get("client_ids"))
		if err == nil {
			request["ClientIds"] = convertListToCommaSeparate(jsonPathResult2.(*schema.Set).List())
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

	if d.HasChange("fingerprints") {
		oldEntry, newEntry := d.GetChange("fingerprints")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		// API constraints: fingerprints must have at least 1 and at most 5 entries
		// Calculate remaining count after removal
		currentCount := oldEntrySet.Len()
		afterRemovalCount := currentCount - removed.Len()

		// Determine the operation order based on constraints
		if afterRemovalCount < 1 && added.Len() > 0 {
			// Need to replace all fingerprints, must maintain at least 1 at all times
			removedList := removed.List()
			addedList := added.List()

			// Check if adding would exceed the limit of 5
			if currentCount+added.Len() > 5 {
				// Interleaved operation: remove some, add all, remove rest
				// Step 1: Remove all but one to make room for new ones
				for i := 0; i < removed.Len()-1; i++ {
					item := removedList[i]
					action := "RemoveFingerprintFromOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

				// Step 2: Add all new fingerprints
				for _, item := range addedList {
					action := "AddFingerprintToOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

				// Step 3: Remove the last old fingerprint
				lastItem := removedList[removed.Len()-1]
				action := "RemoveFingerprintFromOIDCProvider"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["OIDCProviderName"] = d.Id()

				if v, ok := lastItem.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["Fingerprint"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
			} else {
				// Can add all new ones first without exceeding limit, then remove old ones
				// Step 1: Add all new fingerprints
				for _, item := range addedList {
					action := "AddFingerprintToOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

				// Step 2: Remove all old fingerprints
				for _, item := range removedList {
					action := "RemoveFingerprintFromOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
			}
		} else {
			// Remove first to avoid exceeding the limit of 5 when adding
			if removed.Len() > 0 {
				fingerprints := removed.List()
				for _, item := range fingerprints {
					action := "RemoveFingerprintFromOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
			}

			// Then add new fingerprints
			if added.Len() > 0 {
				fingerprints := added.List()
				for _, item := range fingerprints {
					action := "AddFingerprintToOIDCProvider"
					request = make(map[string]interface{})
					query = make(map[string]interface{})
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
						response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
			}
		}

	}
	return resourceAliCloudImsOidcProviderRead(d, meta)
}

func resourceAliCloudImsOidcProviderDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteOIDCProvider"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["OIDCProviderName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.OIDCProvider"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
