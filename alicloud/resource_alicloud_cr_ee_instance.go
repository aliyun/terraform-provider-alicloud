// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrInstanceCreate,
		Read:   resourceAliCloudCrInstanceRead,
		Update: resourceAliCloudCrInstanceUpdate,
		Delete: resourceAliCloudCrInstanceDelete,
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
			"custom_oss_bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_oss_bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"true", "false"}, false),
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_scanner": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ACR", "SAS"}, false),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "Standard", "Advanced"}, false),
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'created_time' has been deprecated since provider version 1.235.0. New field 'create_time' instead.",
			},
		},
	}
}

func resourceAliCloudCrInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("instance_type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceType",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("instance_name"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceName",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	if v, ok := d.GetOk("default_oss_bucket"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "DefaultOssBucket",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("custom_oss_bucket"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceStorageName",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("image_scanner"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "image_scanner",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = d.Get("payment_type")
	request["ProductCode"] = "acr"
	request["ProductType"] = "acr_ee_public_cn"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "acr"
				request["ProductType"] = "acr_ee_public_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	crServiceV2 := CrServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, crServiceV2.CrInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCrInstanceUpdate(d, meta)
}

func resourceAliCloudCrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_ee_instance DescribeCrInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["InstanceName"] != nil {
		d.Set("instance_name", objectRaw["InstanceName"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["InstanceSpecification"] != nil {
		d.Set("instance_type", strings.TrimPrefix(objectRaw["InstanceSpecification"].(string), "Enterprise_"))
	}
	if objectRaw["InstanceStatus"] != nil {
		d.Set("status", objectRaw["InstanceStatus"])
	}

	objectRaw, err = crServiceV2.DescribeQueryAvailableInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["EndTime"] != nil {
		d.Set("end_time", objectRaw["EndTime"])
	}
	if objectRaw["SubscriptionType"] != nil {
		d.Set("payment_type", objectRaw["SubscriptionType"])
	}
	if objectRaw["Region"] != nil {
		d.Set("region_id", objectRaw["Region"])
	}
	if objectRaw["RenewalDuration"] == nil {
		d.Set("renew_period", objectRaw["RenewalDuration"])
	}
	if objectRaw["RenewalDuration"] != nil {
		d.Set("renew_period", objectRaw["RenewalDuration"])
	}
	if objectRaw["RenewStatus"] != nil {
		d.Set("renewal_status", objectRaw["RenewStatus"])
	}

	d.Set("created_time", d.Get("create_time"))
	return nil
}

func resourceAliCloudCrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ChangeResourceGroup"
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), query, request, &runtime)
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
		crServiceV2 := CrServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, crServiceV2.CrInstanceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ResetLoginPassword"
	conn, err = client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChanges("password", "kms_encrypted_password") {
		update = true
	}
	request["Password"] = d.Get("password")
	if update && request["Password"] == "" {
		kmsPassword, ok := d.Get("kms_encrypted_password").(string)
		if ok {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Password"] = decryptResp
		}
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), query, request, &runtime)
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

	d.Partial(false)
	return resourceAliCloudCrInstanceRead(d, meta)
}

func resourceAliCloudCrInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	request["ProductCode"] = "acr"
	request["ProductType"] = "acr_ee_public_cn"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "acr"
				request["ProductType"] = "acr_ee_public_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotExists"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	crServiceV2 := CrServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"INSTANCE_NOT_EXIST"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, crServiceV2.DescribeAsyncCrInstanceStateRefreshFunc(d, response, "$.Code", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
