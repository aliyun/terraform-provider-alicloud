// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsValueAddedService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsValueAddedServiceCreate,
		Read:   resourceAliCloudKmsValueAddedServiceRead,
		Update: resourceAliCloudKmsValueAddedServiceUpdate,
		Delete: resourceAliCloudKmsValueAddedServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3),
			},
			"renew_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_added_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudKmsValueAddedServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("renew_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok && v.(int) > 0 {
		request["RenewPeriod"] = convertKmsValueAddedServiceRenewPeriodRequest(v.(int))
	}
	if v, ok := d.GetOkExists("period"); ok && v.(int) > 0 {
		request["Period"] = convertKmsValueAddedServicePeriodRequest(v.(int))
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("value_added_service"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ValueAddedService",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "ProductVersion",
		"Value": "4",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList

	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_value_added_service", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudKmsValueAddedServiceRead(d, meta)
}

func resourceAliCloudKmsValueAddedServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsValueAddedService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_value_added_service DescribeKmsValueAddedService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("region_id", objectRaw["Region"])
	d.Set("renew_period", formatInt(convertKmsValueAddedServiceDataInstanceListRenewalDurationResponse(objectRaw["RenewalDuration"])))
	d.Set("renew_status", objectRaw["RenewStatus"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudKmsValueAddedServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if d.HasChange("renew_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renew_status")
	if d.HasChange("renew_period") {
		update = true
		request["RenewalPeriod"] = d.Get("renew_period")
	}

	request["RenewalPeriodUnit"] = "Y"
	request["SubscriptionType"] = "Subscription"
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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

	return resourceAliCloudKmsValueAddedServiceRead(d, meta)
}

func resourceAliCloudKmsValueAddedServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 60*time.Second, kmsServiceV2.KmsValueAddedServiceStateRefreshFunc(d.Id(), "$.InstanceID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertKmsValueAddedServiceDataInstanceListRenewalDurationResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "12":
		return "1"
	case "24":
		return "2"
	case "36":
		return "3"
	}
	return source
}
func convertKmsValueAddedServiceDataInstanceListRenewalDurationUnitResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "M":
		return "Y"
	}
	return source
}
func convertKmsValueAddedServiceRenewPeriodRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "12"
	case "2":
		return "24"
	case "3":
		return "36"
	}
	return source
}
func convertKmsValueAddedServicePeriodRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "12"
	case "2":
		return "24"
	case "3":
		return "36"
	}
	return source
}
