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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"1", "2"}, false),
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
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"post_paid_flag": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"post_pay_module_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"rasp_count": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(1, 100000000),
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				Computed:     true,
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
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
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
			"subscription_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
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
				ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "Threat analysis and response log storage capacity. Interval type, value interval:[100,9999999999].> The step size is 100, that is, only multiples of 100 can be filled in."),
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
				Optional:     true,
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
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
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
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan_new"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan_new",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("post_paid_flag"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "PostPaidFlag",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch1"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch_1",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_anti_ransomware"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_anti_ransomware",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_flow"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_flow",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("rasp_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "rasp_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sc",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = d.Get("payment_type")
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	var endpoint string
	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "sas"
		request["ProductType"] = "sas_postpaid_public_cn"
		if client.IsInternationalAccount() {
			request["ProductType"] = "sas_postpaid_public_intl"
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
				request["ProductCode"] = "sas"
				request["ProductType"] = ""
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "sas"
					request["ProductType"] = "sas_postpaid_public_intl"
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

	instanceComponentValueRawObj, _ := jsonpath.Get("$.InstanceComponentValue", objectRaw)
	instanceComponentValueRaw := make(map[string]interface{})
	if instanceComponentValueRawObj != nil {
		instanceComponentValueRaw = instanceComponentValueRawObj.(map[string]interface{})
	}
	d.Set("buy_number", instanceComponentValueRaw["BuyNumber"])
	d.Set("container_image_scan", instanceComponentValueRaw["ContainerImageScan"])
	d.Set("container_image_scan_new", instanceComponentValueRaw["ContainerImageScanNew"])
	d.Set("honeypot", instanceComponentValueRaw["Honeypot"])
	d.Set("honeypot_switch", instanceComponentValueRaw["HoneypotSwitch"])
	d.Set("rasp_count", instanceComponentValueRaw["RaspCount"])
	d.Set("sas_anti_ransomware", instanceComponentValueRaw["SasAntiRansomware"])
	d.Set("sas_cspm", instanceComponentValueRaw["SasCspm"])
	d.Set("sas_cspm_switch", instanceComponentValueRaw["SasCspmSwitch"])
	if instanceComponentValueRaw["SasSc"] != nil && instanceComponentValueRaw["SasSc"] != "" {
		d.Set("sas_sc", convertStringToBool(instanceComponentValueRaw["SasSc"].(string)))
	}
	d.Set("sas_sdk", instanceComponentValueRaw["SasSdk"])
	d.Set("sas_sdk_switch", instanceComponentValueRaw["SasSdkSwitch"])
	d.Set("sas_sls_storage", instanceComponentValueRaw["SasSlsStorage"])
	d.Set("sas_webguard_boolean", instanceComponentValueRaw["SasWebguardBoolean"])
	d.Set("sas_webguard_order_num", instanceComponentValueRaw["SasWebguardOrderNum"])
	d.Set("threat_analysis", instanceComponentValueRaw["ThreatAnalysis"])
	d.Set("threat_analysis_flow", instanceComponentValueRaw["ThreatAnalysisFlow"])
	d.Set("threat_analysis_sls_storage", instanceComponentValueRaw["ThreatAnalysisSlsStorage"])
	d.Set("threat_analysis_switch", instanceComponentValueRaw["ThreatAnalysisSwitch"])
	d.Set("threat_analysis_switch1", instanceComponentValueRaw["ThreatAnalysisSwitch1"])
	d.Set("v_core", instanceComponentValueRaw["VCore"])
	d.Set("version_code", instanceComponentValueRaw["VersionCode"])
	d.Set("vul_count", instanceComponentValueRaw["VulCount"])
	d.Set("vul_switch", instanceComponentValueRaw["VulSwitch"])

	objectRaw, err = threatDetectionServiceV2.DescribeInstanceQueryAvailableInstances(d)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	if v, ok := objectRaw["RenewalDuration"]; ok && formatInt(v) != 0 {
		d.Set("renew_period", formatInt(objectRaw["RenewalDuration"]))
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
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChanges("vul_count", "sas_cspm_switch", "sas_webguard_order_num", "sas_sls_storage", "sas_sdk", "container_image_scan_new", "container_image_scan", "honeypot_switch", "threat_analysis_switch", "honeypot", "threat_analysis", "sas_sdk_switch", "threat_analysis_switch1", "sas_anti_ransomware", "threat_analysis_flow", "v_core", "rasp_count", "vul_switch", "sas_sc", "version_code", "threat_analysis_sls_storage", "sas_webguard_boolean", "buy_number", "sas_cspm") {
		update = true
	}
	parameterMapList := make([]map[string]interface{}, 0)
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
	if v, ok := d.GetOk("sas_webguard_order_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_order_num",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan_new"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan_new",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("container_image_scan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Container_Image_Scan",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("honeypot"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "honeypot",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sdk_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_SDK_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_switch1"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_switch_1",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_anti_ransomware"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_anti_ransomware",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_flow"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_flow",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("v_core"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "vCore",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("rasp_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "rasp_count",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("vul_switch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Vul_switch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_sc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_sc",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("version_code"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VersionCode",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("threat_analysis_sls_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Threat_analysis_sls_storage",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_webguard_boolean"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_webguard_boolean",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("buy_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BuyNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("sas_cspm"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "sas_cspm",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	request["ModifyType"] = d.Get("modify_type")
	var endpoint string
	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "sas"
		request["ProductType"] = "sas_postpaid_public_cn"
		if client.IsInternationalAccount() {
			request["ProductType"] = "sas_postpaid_public_intl"
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
					request["ProductCode"] = "sas"
					request["ProductType"] = ""
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "sas"
						request["ProductType"] = "sas_postpaid_public_intl"
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
	update = false
	action = "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if v, ok := d.GetOk("subscription_type"); ok {
		request["SubscriptionType"] = v
	}
	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renewal_status")
	if d.HasChange("renewal_period_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_period_unit"); ok {
		request["RenewalPeriodUnit"] = v
	}

	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true
	}
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewalPeriod"] = v
	}

	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "sas"
		request["ProductType"] = "sas_postpaid_public_cn"
		if client.IsInternationalAccount() {
			request["ProductType"] = "sas_postpaid_public_intl"
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
					request["ProductCode"] = "sas"
					request["ProductType"] = ""
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "sas"
						request["ProductType"] = "sas_postpaid_public_intl"
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
	update = false
	action = "ModifyPostPayModuleSwitch"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PostPayInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("post_pay_module_switch") {
		update = true
	}
	if v, ok := d.GetOk("post_pay_module_switch"); ok {
		request["PostPayModuleSwitch"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
	return resourceAliCloudThreatDetectionInstanceRead(d, meta)
}

func resourceAliCloudThreatDetectionInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
