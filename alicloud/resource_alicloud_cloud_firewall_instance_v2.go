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

func resourceAliCloudCloudFirewallInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallInstanceV2Create,
		Read:   resourceAliCloudCloudFirewallInstanceV2Read,
		Update: resourceAliCloudCloudFirewallInstanceV2Update,
		Delete: resourceAliCloudCloudFirewallInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 1000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cfw_account"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"auto_asset_protection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(10, 15000),
			},
			"cfw_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cfw_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cfw_log_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1000, 500000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cfw_log"); ok && v.(bool) && d.Get("payment_type").(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fw_vpc_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 500),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("spec"); ok && v.(string) == "premium_version" {
						return true
					}
					return false
				},
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(5, 5000),
			},
			"ip_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(20, 4000),
			},
			"logistics": {
				Type:     schema.TypeString,
				Optional: true,
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
			"product_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cfw"}, false),
			},
			"product_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cfw_elasticity_public_cn", "cfw_elasticity_public_intl", "cfw_sub_public_cn", "cfw_sub_public_intl"}, false),
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renewal_duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sdl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"payg_version", "premium_version", "enterprise_version", "ultimate_version"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallInstanceV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOkExists("sdl"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "cfw_ndlp_enable",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOkExists("cfw_log"); ok {
		if d.Get("payment_type").(string) == "Subscription" {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "cfw_log",
				"Value": fmt.Sprint(v),
			})
		} else {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "CfwLog",
				"Value": fmt.Sprint(v),
			})
		}
	}
	if v, ok := d.GetOk("spec"); ok {
		if d.Get("payment_type").(string) == "PayAsYouGo" {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		} else {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "cfw_spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		}
	}
	if v, ok := d.GetOk("ip_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("band_width"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BandWidth",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("cfw_log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLogStorage",
			"Value": v,
		})
	}
	if v, ok := d.GetOkExists("cfw_account"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwAccount",
			"Value": v,
		})
	}
	if v, ok := d.GetOkExists("account_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwAccountNum",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("fw_vpc_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "FwVpcNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("instance_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceCount",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("auto_asset_protection"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "AutoAssetProtection",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")
	request["ProductType"] = d.Get("product_type")

	if v, ok := d.GetOk("logistics"); ok {
		request["Logistics"] = v
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renewal_duration"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	var endpoint string
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable", NotFoundArticle}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				if _, ok := d.GetOkExists("account_number"); ok {
					for _, v := range parameterMapList {
						if fmt.Sprint(v["Code"]) == "CfwAccountNum" {
							v["Code"] = "CfwAccountIntlNum"
						}
					}
					request["Parameter"] = parameterMapList
				}
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_instance_v2", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFuncWithApi(d.Id(), "$.InstanceStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallInstanceV2Update(d, meta)
}

func resourceAliCloudCloudFirewallInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance_v2 DescribeCloudFirewallInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("product_type", objectRaw["ProductType"])
	d.Set("release_time", objectRaw["ReleaseTime"])
	d.Set("renewal_duration", objectRaw["RenewalDuration"])
	d.Set("renewal_duration_unit", objectRaw["RenewalDurationUnit"])
	d.Set("renewal_status", objectRaw["RenewStatus"])

	objectRaw, err = cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cfw_log", objectRaw["LogStatus"])
	d.Set("sdl", convertCloudFirewallInstanceSdlResponse(objectRaw["Sdl"]))
	d.Set("spec", convertCloudFirewallInstanceVersionResponse(objectRaw["Version"]))
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("user_status", fmt.Sprint(objectRaw["UserStatus"]))
	// DescribeUserBuyVersion reports the purchased spec counts (IpNumber /
	// VpcNumber / LogStorage) as 0 for Subscription instances even when
	// InstanceStatus is normal, while PayAsYouGo returns the real values.
	// Overwriting the user-configured counts with 0 produces a spurious diff
	// on every plan, so these fields are only set when the API reports a
	// non-zero value. Drift detection for these counts is therefore limited to
	// instances where the API actually surfaces the value.
	if v := objectRaw["LogStorage"]; v != nil && fmt.Sprint(v) != "0" {
		d.Set("cfw_log_storage", v)
	}
	if v := objectRaw["VpcNumber"]; v != nil && fmt.Sprint(v) != "0" {
		d.Set("fw_vpc_number", v)
	}
	if v := objectRaw["IpNumber"]; v != nil && fmt.Sprint(v) != "0" {
		d.Set("ip_number", v)
	}

	assetStatistic, err := cloudFirewallServiceV2.DescribeInstanceDescribeAssetStatistic(d.Id())
	if err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
	} else {
		d.Set("auto_asset_protection", fmt.Sprint(assetStatistic["AutoResourceEnable"]))
	}

	return nil
}

func resourceAliCloudCloudFirewallInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyCfwInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	updateList := make([]map[string]interface{}, 0)

	objectRaw, err := cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(convertCloudFirewallInstanceSdlResponse(objectRaw["Sdl"])) != fmt.Sprint(d.Get("sdl")) {
		update = true
	}
	if v, ok := d.GetOkExists("sdl"); ok {
		updateList = append(updateList, map[string]interface{}{
			"Code":  "Sdl",
			"Value": convertCloudFirewallInstanceSdlRequest(v),
		})
	}

	request["UpdateList"] = updateList

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
	update = false
	action = "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	var endpoint string

	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	request["ProductType"] = d.Get("product_type")
	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")

	if v, ok := d.GetOk("modify_type"); ok {
		request["ModifyType"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)

	// Append every order Parameter that has a value so ModifyInstance receives
	// the complete component list; HasChange only gates the update trigger.
	// Sending the full set (including unchanged specs) is required for BssOpenApi
	// ModifyInstance to accept the Upgrade/Downgrade without
	// UPDOWNGRADE_CONFIG_NO_CHANGE.
	if !d.IsNewResource() && d.HasChange("cfw_account") {
		update = true
	}
	if v, ok := d.GetOkExists("cfw_account"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwAccount",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("account_number") {
		update = true
	}
	if v, ok := d.GetOkExists("account_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwAccountNum",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("fw_vpc_number") {
		update = true
	}
	if v, ok := d.GetOk("fw_vpc_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "FwVpcNumber",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("ip_number") {
		update = true
	}
	if v, ok := d.GetOk("ip_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpNumber",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("cfw_log") {
		update = true
	}
	if v, ok := d.GetOkExists("cfw_log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLog",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("cfw_log_storage") {
		update = true
	}
	if v, ok := d.GetOk("cfw_log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLogStorage",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("band_width") {
		update = true
	}
	if v, ok := d.GetOk("band_width"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BandWidth",
			"Value": v,
		})
	}

	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
	}
	if v, ok := d.GetOk("spec"); ok {
		if d.Get("payment_type").(string) == "PayAsYouGo" {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		} else {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "cfw_spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_count") {
		update = true
	}
	if v, ok := d.GetOk("instance_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceCount",
			"Value": v,
		})
	}

	request["Parameter"] = parameterMapList

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
					endpoint = connectivity.BssOpenAPIEndpointInternational
					if _, ok := d.GetOkExists("account_number"); ok {
						for _, v := range parameterMapList {
							if fmt.Sprint(v["Code"]) == "CfwAccountNum" {
								v["Code"] = "CfwAccountIntlNum"
							}
						}
						request["Parameter"] = parameterMapList
					}
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			if fmt.Sprint(response["Code"]) == "SYSTEM.CONCURRENT_OPERATE" {
				wait()
				return resource.RetryableError(fmt.Errorf("%s", response))
			}

			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		if d.HasChange("cfw_log") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("cfw_log"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFuncWithApi(d.Id(), "LogStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		bssOpenApiService := BssOpenApiService{client}
		stateConf := BuildStateConf([]string{}, []string{"Paid"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, bssOpenApiService.CloudFirewallInstanceOrderDetailStateRefreshFunc(fmt.Sprint(response["Data"].(map[string]interface{})["OrderId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	request["ProductCode"] = d.Get("product_code")
	request["ProductType"] = d.Get("product_type")
	request["SubscriptionType"] = d.Get("payment_type")

	if d.HasChange("renewal_duration_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		request["RenewalPeriodUnit"] = v
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewalPeriod"] = v
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
	setAutoProtectNewAssetsRequest := make(map[string]interface{})
	setAutoProtectNewAssetsQuery := make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("auto_asset_protection") {
		update = true
	}
	setAutoProtectNewAssetsRequest["AutoProtect"] = d.Get("auto_asset_protection")

	if update {
		action := "SetAutoProtectNewAssets"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, setAutoProtectNewAssetsQuery, setAutoProtectNewAssetsRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setAutoProtectNewAssetsRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("auto_asset_protection")
	}

	d.Partial(false)
	return resourceAliCloudCloudFirewallInstanceV2Read(d, meta)
}

func resourceAliCloudCloudFirewallInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.Get("payment_type").(string) == "Subscription" {
		// Prepaid (Subscription) Cloud Firewall instances cannot be released
		// through Cloudfw ReleasePostInstance; unsubscribe them via BssOpenApi
		// RefundInstance so that no cloud resource is left behind on destroy.
		action := "RefundInstance"
		request := map[string]interface{}{
			"InstanceId":         d.Id(),
			"ClientToken":        buildClientToken(action),
			"ImmediatelyRelease": "1",
			"ProductCode":        "vipcloudfw",
			"ProductType":        "vipcloudfw",
		}
		if client.IsInternationalAccount() {
			request["ProductCode"] = "cfw"
			request["ProductType"] = "cfw_pre_intl"
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var response map[string]interface{}
		var err error
		var endpoint string
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "cfw"
					request["ProductType"] = "cfw_pre_intl"
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
		bssOpenApiService := BssOpenApiService{client}
		bssWait := incrementalWait(10*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			_, err = bssOpenApiService.QueryAvailableInstance(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				if NeedRetry(err) {
					bssWait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			bssWait()
			return resource.RetryableError(fmt.Errorf("Cloud Firewall instance %s still exists", d.Id()))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "WaitForRefundInstance", AlibabaCloudSdkGoERROR)
		}
		return nil
	}

	enableDelete := false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "ReleasePostInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		var endpoint string
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFunc(d.Id(), "#$.Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return nil
}

func convertCloudFirewallInstanceVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "2":
		return "premium_version"
	case "3":
		return "enterprise_version"
	case "4":
		return "ultimate_version"
	case "10":
		return "payg_version"
	}

	return source
}

func convertCloudFirewallInstanceSdlRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return 1
	case "false":
		return 0
	}

	return source
}

func convertCloudFirewallInstanceSdlResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "0":
		return false
	case "1":
		return true
	}

	return source
}
