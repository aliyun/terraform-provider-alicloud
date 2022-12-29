package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionInstanceCreate,
		Read:   resourceAlicloudThreatDetectionInstanceRead,
		Update: resourceAlicloudThreatDetectionInstanceUpdate,
		Delete: resourceAlicloudThreatDetectionInstanceDelete,
		Schema: map[string]*schema.Schema{
			"buy_number": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"container_image_scan": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"honeypot": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"honeypot_switch": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
				Type:         schema.TypeString,
			},
			"instance_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"modify_type": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Downgrade", "Upgrade"}, false),
				Type:         schema.TypeString,
			},
			"payment_type": {
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
				Type:         schema.TypeString,
			},
			"period": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"renew_period": {
				Optional: true,
				Type:     schema.TypeInt,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"M", "Y"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_status": {
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				Type:         schema.TypeString,
			},
			"sas_anti_ransomware": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"sas_sc": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"sas_sdk": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"sas_sdk_switch": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
				Type:         schema.TypeString,
			},
			"sas_sls_storage": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"sas_webguard_boolean": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
				Type:         schema.TypeString,
			},
			"sas_webguard_order_num": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"threat_analysis": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"threat_analysis_switch": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
				Type:         schema.TypeString,
			},
			"v_core": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"version_code": {
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"level10", "level2", "level3", "level7", "level8"}, false),
				Type:         schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOk("sas_anti_ransomware"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_anti_ransomware",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sc",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renew_period", "renewal_status", d.Get("renewal_status")))
	}
	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"

	request["ClientToken"] = buildClientToken("CreateInstance")
	var response map[string]interface{}
	action := "CreateInstance"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				request["ProductType"] = ""
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	if v, err := jsonpath.Get("$.Data.InstanceId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionInstanceRead(d, meta)
}

func resourceAlicloudThreatDetectionInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_instance threatDetectionService.DescribeThreatDetectionInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("instance_id", object["InstanceID"])
	d.Set("payment_type", object["SubscriptionType"])
	d.Set("status", object["Status"])
	d.Set("renewal_status", object["RenewStatus"])
	if v, ok := object["RenewalDuration"]; ok && formatInt(v) != 0 {
		d.Set("renew_period", formatInt(v))
	}

	d.Set("renewal_period_unit", object["RenewalDurationUnit"])

	return nil
}

func resourceAlicloudThreatDetectionInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	update := false
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	request["ModifyType"] = d.Get("modify_type")
	if d.HasChanges("buy_number", "container_image_scan", "honeypot", "honeypot_switch", "sas_anti_ransomware", "sas_sc", "sas_sdk", "sas_sdk_switch", "sas_sls_storage", "sas_webguard_boolean", "threat_analysis_switch", "threat_analysis", "sas_webguard_order_num", "version_code", "v_core") {
		update = true
	}

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_anti_ransomware"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_anti_ransomware",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sc",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = d.Get("payment_type")
	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"

	if update {
		action := "ModifyInstance"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					request["ProductType"] = ""
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("parameter")
	}

	var setRenewalResponse map[string]interface{}
	update = false
	setRenewalReq := map[string]interface{}{
		"InstanceIDs":      d.Id(),
		"ProductCode":      "sas",
		"ProductType":      "sas",
		"SubscriptionType": d.Get("payment_type"),
	}

	if d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		setRenewalReq["RenewalStatus"] = v
	}

	if d.HasChange("renew_period") {
		update = true
		if v, ok := d.GetOk("renew_period"); ok {
			setRenewalReq["RenewalPeriod"] = v
		}
	}

	if d.HasChange("renewal_period_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_period_unit"); ok {
		setRenewalReq["RenewalPeriodUnit"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_period_unit", "renewal_status", d.Get("renewal_status")))
	}

	if update {
		action := "SetRenewal"
		conn, err := client.NewBssopenapiClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			setRenewalResponse, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, setRenewalReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					setRenewalReq["ProductType"] = ""
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, setRenewalResponse, setRenewalReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(setRenewalResponse["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, setRenewalResponse))
		}

		d.SetPartial("renewal_status")
		d.SetPartial("renew_period")
		d.SetPartial("renewal_period_unit")
	}

	d.Partial(false)

	return resourceAlicloudThreatDetectionInstanceRead(d, meta)
}

func resourceAlicloudThreatDetectionInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
