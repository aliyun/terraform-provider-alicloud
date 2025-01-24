package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(1, 100000000),
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
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Cloud platform configuration check scan times, interval type, value range:[15000,9999999999].> You must have sas_cspm_switch = 1 to purchase this module. The step size is 55000, that is, only multiples of 55000 can be filled in."),
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
			"threat_analysis_flow": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "Threat analysis and response log access traffic. After ThreatAnalysisSwitch1 is selected, it must be selected. Interval type, value interval:[0,9999999999].> Step size is 1."),
			},
			"threat_analysis_sls_storage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "Threat analysis and response log storage capacity. Interval type, value interval:[0,9999999999].> The step size is 10, that is, only multiples of 10 can be filled in."),
			},
			"threat_analysis_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
			"threat_analysis_switch1": {
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
	query := make(map[string]interface{})
	var err error
	var endpoint string
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
	if v, ok := d.GetOk("threat_analysis_switch1"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch_1",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_flow"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_flow",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_sls_storage",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductType"] = "sas"
	request["ProductCode"] = "sas"
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
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
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_instance", action, AlibabaCloudSdkGoERROR)
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
	if instanceComponentValue1Raw["BuyNumber"] != nil {
		d.Set("buy_number", instanceComponentValue1Raw["BuyNumber"])
	}
	if instanceComponentValue1Raw["ContainerImageScan"] != nil {
		d.Set("container_image_scan", instanceComponentValue1Raw["ContainerImageScan"])
	}
	if instanceComponentValue1Raw["ContainerImageScanNew"] != nil {
		d.Set("container_image_scan_new", instanceComponentValue1Raw["ContainerImageScanNew"])
	}
	if instanceComponentValue1Raw["Honeypot"] != nil {
		d.Set("honeypot", instanceComponentValue1Raw["Honeypot"])
	}
	if instanceComponentValue1Raw["HoneypotSwitch"] != nil {
		d.Set("honeypot_switch", instanceComponentValue1Raw["HoneypotSwitch"])
	}
	if instanceComponentValue1Raw["RaspCount"] != nil {
		d.Set("rasp_count", instanceComponentValue1Raw["RaspCount"])
	}
	if instanceComponentValue1Raw["SasAntiRansomware"] != nil {
		d.Set("sas_anti_ransomware", instanceComponentValue1Raw["SasAntiRansomware"])
	}
	if instanceComponentValue1Raw["SasCspm"] != nil {
		d.Set("sas_cspm", instanceComponentValue1Raw["SasCspm"])
	}
	if instanceComponentValue1Raw["SasCspmSwitch"] != nil {
		d.Set("sas_cspm_switch", instanceComponentValue1Raw["SasCspmSwitch"])
	}
	if instanceComponentValue1Raw["SasSc"] != nil && instanceComponentValue1Raw["SasSc"] != "" {
		d.Set("sas_sc", convertStringToBool(instanceComponentValue1Raw["SasSc"].(string)))
	}
	if instanceComponentValue1Raw["SasSdk"] != nil {
		d.Set("sas_sdk", instanceComponentValue1Raw["SasSdk"])
	}
	if instanceComponentValue1Raw["SasSdkSwitch"] != nil {
		d.Set("sas_sdk_switch", instanceComponentValue1Raw["SasSdkSwitch"])
	}
	if instanceComponentValue1Raw["SasSlsStorage"] != nil {
		d.Set("sas_sls_storage", instanceComponentValue1Raw["SasSlsStorage"])
	}
	if instanceComponentValue1Raw["SasWebguardBoolean"] != nil {
		d.Set("sas_webguard_boolean", instanceComponentValue1Raw["SasWebguardBoolean"])
	}
	if instanceComponentValue1Raw["SasWebguardOrderNum"] != nil {
		d.Set("sas_webguard_order_num", instanceComponentValue1Raw["SasWebguardOrderNum"])
	}
	if instanceComponentValue1Raw["ThreatAnalysis"] != nil {
		d.Set("threat_analysis", instanceComponentValue1Raw["ThreatAnalysis"])
	}
	if instanceComponentValue1Raw["ThreatAnalysisFlow"] != nil {
		d.Set("threat_analysis_flow", instanceComponentValue1Raw["ThreatAnalysisFlow"])
	}
	if instanceComponentValue1Raw["ThreatAnalysisSlsStorage"] != nil {
		d.Set("threat_analysis_sls_storage", instanceComponentValue1Raw["ThreatAnalysisSlsStorage"])
	}
	if instanceComponentValue1Raw["ThreatAnalysisSwitch"] != nil {
		d.Set("threat_analysis_switch", instanceComponentValue1Raw["ThreatAnalysisSwitch"])
	}
	if instanceComponentValue1Raw["ThreatAnalysisSwitch1"] != nil {
		d.Set("threat_analysis_switch1", instanceComponentValue1Raw["ThreatAnalysisSwitch1"])
	}
	if instanceComponentValue1Raw["VCore"] != nil {
		d.Set("v_core", instanceComponentValue1Raw["VCore"])
	}
	if instanceComponentValue1Raw["VersionCode"] != nil {
		d.Set("version_code", instanceComponentValue1Raw["VersionCode"])
	}
	if instanceComponentValue1Raw["VulCount"] != nil {
		d.Set("vul_count", instanceComponentValue1Raw["VulCount"])
	}
	if instanceComponentValue1Raw["VulSwitch"] != nil {
		d.Set("vul_switch", instanceComponentValue1Raw["VulSwitch"])
	}

	objectRaw, err = threatDetectionServiceV2.DescribeQueryAvailableInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["SubscriptionType"] != nil {
		d.Set("payment_type", objectRaw["SubscriptionType"])
	}
	if v, ok := objectRaw["RenewalDuration"]; ok && formatInt(v) != 0 {
		d.Set("renew_period", formatInt(objectRaw["RenewalDuration"]))
	}
	if objectRaw["RenewalDurationUnit"] != nil {
		d.Set("renewal_period_unit", objectRaw["RenewalDurationUnit"])
	}
	if objectRaw["RenewStatus"] != nil {
		d.Set("renewal_status", objectRaw["RenewStatus"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	return nil
}

func resourceAliCloudThreatDetectionInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	request["ModifyType"] = d.Get("modify_type")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if !d.IsNewResource() && d.HasChanges("version_code", "buy_number", "v_core", "threat_analysis", "threat_analysis_switch", "sas_sdk", "sas_sdk_switch", "honeypot", "honeypot_switch", "sas_anti_ransomware", "sas_sc", "sas_sls_storage", "sas_webguard_order_num", "sas_webguard_boolean", "container_image_scan", "vul_switch", "vul_count", "sas_cspm_switch", "sas_cspm", "rasp_count", "container_image_scan_new", "threat_analysis_switch1", "threat_analysis_flow", "threat_analysis_sls_storage") {
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
	if v, ok := d.GetOk("threat_analysis_switch1"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch_1",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_flow"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_flow",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_sls_storage",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductType"] = "sas"
	request["ProductCode"] = "sas"
	if update {
		var endpoint string
		if client.IsInternationalAccount() {
			request["ProductType"] = ""
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
	update = false
	action = "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = d.Id()

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
	if v, ok := d.GetOk("subscription_type"); ok {
		request["SubscriptionType"] = v
	}
	if update {
		var endpoint string
		if client.IsInternationalAccount() {
			request["ProductType"] = ""
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
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

	d.Partial(false)
	return resourceAliCloudThreatDetectionInstanceRead(d, meta)
}

func resourceAliCloudThreatDetectionInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
