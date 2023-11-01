// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionInstanceCreate,
		Read:   resourceAliCloudThreatDetectionInstanceRead,
		Update: resourceAliCloudThreatDetectionInstanceUpdate,
		Delete: resourceAliCloudThreatDetectionInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"buy_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_image_scan": {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Field 'container_image_scan' has been deprecated from provider version 1.212.0. Container Image security scan. Interval type, value interval:[0,200000].> The step size is 20, that is, only multiples of 20 can be filled in.",
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Container Image security scan. Interval type, value interval:[0,200000].> The step size is 20, that is, only multiples of 20 can be filled in."),
			},
			"container_image_scan_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Container Image security scan. Interval type, value interval:[0,200000].> The step size is 20, that is, only multiples of 20 can be filled in."),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"honeypot": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Number of cloud honeypot licenses. Interval type, value interval:[20,500].> This module can only be purchased when honeypot_switch = 1, starting with 20."),
			},
			"honeypot_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"1", "2"}, false),
				Computed:     true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Upgrade", "Downgrade"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rasp_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
				Computed: true,
			},
			"renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"M", "Y"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"sas_anti_ransomware": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Anti-ransomware capacity. Unit: GB. Interval type, value interval:[0,9999999999].> The step size is 10, that is, only multiples of 10 can be filled in."),
			},
			"sas_cspm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Cloud platform configuration check scan times, interval type, value range:[1000,9999999999].> You must have sas_cspm_switch = 1 to purchase this module. The step size is 100, that is, only multiples of 10 can be filled in."),
			},
			"sas_cspm_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
			"sas_sc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sas_sdk": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Number of malicious file detections. Unit: 10,000 times. Interval type, value interval:[10,9999999999].> This module can only be purchased when sas_sdk_switch = 1. The step size is 10, that is, only multiples of 10 can be filled in."),
			},
			"sas_sdk_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
				Computed:     true,
			},
			"sas_sls_storage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Log analysis storage capacity. Unit: GB. Interval type, value interval:[0,600000].> The step size is 10, that is, only multiples of 10 can be filled in."),
			},
			"sas_webguard_boolean": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
			"sas_webguard_order_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Tamper-proof authorization number. Value:-0: No1: Yes."),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threat_analysis": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Threat Analysis log storage capacity. Interval type, value interval:[0,9999999999].> This module can only be purchased when Threat_analysis_switch = 1. The step size is 10, that is, only multiples of 10 can be filled in."),
			},
			"threat_analysis_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
			"v_core": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"level2", "level8", "level7", "level3", "level10"}, false),
			},
			"vul_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vul_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	request["SubscriptionType"] = d.Get("payment_type")
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
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
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("rasp_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "rasp_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan_new"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan_new",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductType"] = "sas"
	request["ProductCode"] = "sas"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		addDebug(action, response, request)

		if fmt.Sprint(response["Code"]) != "Success" {
			return resource.NonRetryableError(WrapError(fmt.Errorf("%s failed, response: %v", action, response)))
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionInstanceUpdate(d, meta)
}

func resourceAliCloudThreatDetectionInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_instance DescribeThreatDetectionInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	instanceComponentValue1RawObj, _ := jsonpath.Get("$.InstanceComponentValue", objectRaw)
	instanceComponentValue1Raw := make(map[string]interface{})
	if instanceComponentValue1RawObj != nil {
		instanceComponentValue1Raw = instanceComponentValue1RawObj.(map[string]interface{})
	}
	d.Set("buy_number", instanceComponentValue1Raw["BuyNumber"])
	d.Set("container_image_scan", instanceComponentValue1Raw["ContainerImageScan"])
	d.Set("container_image_scan_new", instanceComponentValue1Raw["ContainerImageScanNew"])
	d.Set("honeypot", instanceComponentValue1Raw["Honeypot"])
	d.Set("honeypot_switch", instanceComponentValue1Raw["HoneypotSwitch"])
	d.Set("rasp_count", instanceComponentValue1Raw["RaspCount"])
	d.Set("sas_anti_ransomware", instanceComponentValue1Raw["SasAntiRansomware"])
	d.Set("sas_cspm", instanceComponentValue1Raw["SasCspm"])
	d.Set("sas_cspm_switch", instanceComponentValue1Raw["SasCspmSwitch"])
	if instanceComponentValue1Raw["SasSc"] != nil && instanceComponentValue1Raw["SasSc"] != "" {
		d.Set("sas_sc", convertStringToBool(instanceComponentValue1Raw["SasSc"].(string)))
	}
	d.Set("sas_sdk", instanceComponentValue1Raw["SasSdk"])
	d.Set("sas_sdk_switch", instanceComponentValue1Raw["SasSdkSwitch"])
	d.Set("sas_sls_storage", instanceComponentValue1Raw["SasSlsStorage"])
	d.Set("sas_webguard_boolean", instanceComponentValue1Raw["SasWebguardBoolean"])
	d.Set("sas_webguard_order_num", instanceComponentValue1Raw["SasWebguardOrderNum"])
	d.Set("threat_analysis", instanceComponentValue1Raw["ThreatAnalysis"])
	d.Set("threat_analysis_switch", instanceComponentValue1Raw["ThreatAnalysisSwitch"])
	d.Set("v_core", instanceComponentValue1Raw["VCore"])
	d.Set("version_code", instanceComponentValue1Raw["VersionCode"])
	d.Set("vul_count", instanceComponentValue1Raw["VulCount"])
	d.Set("vul_switch", instanceComponentValue1Raw["VulSwitch"])

	objectRaw, err = threatDetectionServiceV2.DescribeQueryAvailableInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	if v, ok := objectRaw["RenewalDuration"]; ok && formatInt(v) != 0 {
		d.Set("renew_period", formatInt(v))
	}
	d.Set("renewal_period_unit", objectRaw["RenewalDurationUnit"])
	d.Set("renewal_status", objectRaw["RenewStatus"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudThreatDetectionInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstance"
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	request["ModifyType"] = d.Get("modify_type")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if !d.IsNewResource() && d.HasChanges("version_code", "buy_number", "v_core", "threat_analysis", "threat_analysis_switch", "sas_sdk", "sas_sdk_switch", "honeypot", "honeypot_switch", "sas_anti_ransomware", "sas_sc", "sas_sls_storage", "sas_webguard_order_num", "sas_webguard_boolean", "container_image_scan", "vul_switch", "vul_count", "sas_cspm_switch", "sas_cspm", "rasp_count", "container_image_scan_new") {
		update = true
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
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
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("rasp_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "rasp_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan_new"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan_new",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductType"] = "sas"
	request["ProductCode"] = "sas"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
			addDebug(action, response, request)

			if fmt.Sprint(response["Code"]) != "Success" {
				return resource.NonRetryableError(WrapError(fmt.Errorf("%s failed, response: %v", action, response)))
			}

			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "SetRenewal"
	conn, err = client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()
	request["ProductType"] = "sas"
	if d.HasChange("renewal_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renewal_status")

	if d.HasChange("renew_period") {
		update = true
	}
	request["RenewalPeriod"] = d.Get("renew_period")

	if d.HasChange("renewal_period_unit") {
		update = true
	}
	request["RenewalPeriodUnit"] = d.Get("renewal_period_unit")

	request["ProductCode"] = "sas"
	request["SubscriptionType"] = d.Get("payment_type")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
			addDebug(action, response, request)

			if fmt.Sprint(response["Code"]) != "Success" {
				return resource.NonRetryableError(WrapError(fmt.Errorf("%s failed, response: %v", action, response)))
			}

			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("renewal_status")
		d.SetPartial("renew_period")
		d.SetPartial("renewal_period_unit")
	}

	d.Partial(false)
	return resourceAliCloudThreatDetectionInstanceRead(d, meta)
}

func resourceAliCloudThreatDetectionInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
