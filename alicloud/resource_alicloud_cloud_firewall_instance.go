package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallInstanceCreate,
		Read:   resourceAliCloudCloudFirewallInstanceRead,
		Update: resourceAliCloudCloudFirewallInstanceUpdate,
		Delete: resourceAliCloudCloudFirewallInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 12, 2, 3, 6}),
				Computed:     true,
				Deprecated:   "Attribute 'renew_period' has been deprecated since 1.209.1. Using 'renewal_duration' instead.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 12, 2, 3, 6}),
				Computed:     true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
				ConflictsWith: []string{"renew_period"},
			},
			"renewal_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal", "NotRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"logistics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cfw_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Attribute 'cfw_service' does not support longer, and it has been removed since v1.209.1",
			},
			"cfw_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
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
			"fw_vpc_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(2, 500),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("spec"); ok && v.(string) == "premium_version" {
						return true
					}
					return false
				},
			},
			"ip_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(20, 4000),
			},
			"cfw_log_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1000, 500000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cfw_log"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"cfw_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(10, 15000),
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(5, 5000),
			},
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"premium_version", "enterprise_version", "ultimate_version"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Downgrade", "Upgrade"}, false),
			},
		},
	}
}
func resourceAliCloudCloudFirewallInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken(action)
	request["ProductCode"] = "vipcloudfw"
	request["ProductType"] = "vipcloudfw"
	request["SubscriptionType"] = d.Get("payment_type")

	if fmt.Sprint(request["SubscriptionType"]) == "PayAsYouGo" {
		request["ProductCode"] = "cfw"
		request["ProductType"] = "cfw_elasticity_public_cn"
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	} else if d.Get("payment_type").(string) == "Subscription" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v", "period", "payment_type", "Subscription"))
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration", "renewal_status", d.Get("renewal_status")))
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOk("logistics"); ok {
		request["Logistics"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Spec",
		"Value": convertCloudFirewallInstanceVersion(d.Get("spec").(string)),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "IpNumber",
		"Value": d.Get("ip_number"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "BandWidth",
		"Value": d.Get("band_width"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "CfwLog",
		"Value": d.Get("cfw_log"),
	})
	if v, ok := d.GetOk("cfw_log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLogStorage",
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
	request["Parameter"] = parameterMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "cfw"
				request["ProductType"] = "cfw_pre_intl"

				if fmt.Sprint(request["SubscriptionType"]) == "PayAsYouGo" {
					request["ProductType"] = "cfw_elasticity_public_intl"
				}

				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))
	return resourceAliCloudCloudFirewallInstanceRead(d, meta)
}
func resourceAliCloudCloudFirewallInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bssOpenApiService := BssOpenApiService{client}
	getQueryInstanceObject, err := bssOpenApiService.QueryAvailableInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance bssOpenApiService.QueryAvailableInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("create_time", getQueryInstanceObject["CreateTime"])
	d.Set("renewal_status", getQueryInstanceObject["RenewStatus"])
	d.Set("renewal_duration_unit", convertCloudFirewallInstanceRenewalDurationUnitResponse(getQueryInstanceObject["RenewalDurationUnit"]))
	d.Set("renewal_duration", getQueryInstanceObject["RenewalDuration"])
	d.Set("renew_period", getQueryInstanceObject["RenewalDuration"])
	d.Set("status", getQueryInstanceObject["Status"])
	d.Set("payment_type", getQueryInstanceObject["SubscriptionType"])
	d.Set("end_time", getQueryInstanceObject["EndTime"])
	return nil
}
func resourceAliCloudCloudFirewallInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	setRenewalReq := map[string]interface{}{
		"InstanceIDs": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		setRenewalReq["RenewalStatus"] = v
	}

	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true

		if v, ok := d.GetOk("renew_period"); ok {
			setRenewalReq["RenewalPeriod"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true

		if v, ok := d.GetOk("renewal_duration"); ok {
			setRenewalReq["RenewalPeriod"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("renewal_duration_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		setRenewalReq["RenewalPeriodUnit"] = convertCloudFirewallInstanceRenewalDurationUnitRequest(v.(string))
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration_unit", "renewal_status", d.Get("renewal_status")))
	}

	setRenewalReq["SubscriptionType"] = d.Get("payment_type")
	setRenewalReq["ProductCode"] = "vipcloudfw"
	setRenewalReq["ProductType"] = "vipcloudfw"

	if update {
		action := "SetRenewal"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, setRenewalReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					setRenewalReq["ProductCode"] = "cfw"
					setRenewalReq["ProductType"] = "cfw_pre_intl"

					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setRenewalReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("renewal_status")
		d.SetPartial("payment_type")
		d.SetPartial("renewal_duration")
		d.SetPartial("renewal_duration_unit")
	}

	update = false
	modifyInstanceRequest := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	modifyInstanceRequest["ProductType"] = "vipcloudfw"
	modifyInstanceRequest["ProductCode"] = "vipcloudfw"
	modifyInstanceRequest["SubscriptionType"] = d.Get("payment_type")

	if fmt.Sprint(modifyInstanceRequest["SubscriptionType"]) == "PayAsYouGo" {
		modifyInstanceRequest["ProductCode"] = "cfw"
		modifyInstanceRequest["ProductType"] = "cfw_elasticity_public_cn"
	}

	parameterMapList := make([]map[string]interface{}, 0)

	if d.HasChange("cfw_account") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "CfwAccount",
		"Value": d.Get("cfw_account"),
	})

	if d.HasChange("account_number") {
		update = true
	}
	if v, ok := d.GetOk("account_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "AccountNum",
			"Value": v,
		})
	}

	if d.HasChange("fw_vpc_number") {
		update = true
		if v, ok := d.GetOk("fw_vpc_number"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "FwVpcNumber",
				"Value": v,
			})
		}
	}
	if d.HasChange("ip_number") {
		update = true
	}
	if v, ok := d.GetOk("ip_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpNumber",
			"Value": v,
		})
	}
	if d.HasChange("cfw_log_storage") {
		update = true
	}
	if v, ok := d.GetOk("cfw_log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLogStorage",
			"Value": v,
		})
	}
	if d.HasChange("cfw_log") {
		update = true
	}
	if v, ok := d.GetOk("cfw_log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLog",
			"Value": v,
		})
	}
	if d.HasChange("band_width") {
		update = true
	}
	if v, ok := d.GetOk("band_width"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BandWidth",
			"Value": v,
		})
	}
	if d.HasChange("spec") {
		update = true
	}
	if v, ok := d.GetOk("spec"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Spec",
			"Value": convertCloudFirewallInstanceVersion(v.(string)),
		})
	}
	if d.HasChange("instance_count") {
		update = true
	}
	if v, ok := d.GetOk("instance_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceCount",
			"Value": v,
		})
	}
	modifyInstanceRequest["Parameter"] = parameterMapList
	if update {
		if v, ok := d.GetOk("modify_type"); ok {
			modifyInstanceRequest["ModifyType"] = v
		}
		action := "ModifyInstance"
		modifyInstanceRequest["ClientToken"] = buildClientToken(action)
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, modifyInstanceRequest, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					modifyInstanceRequest["ProductCode"] = "cfw"
					modifyInstanceRequest["ProductType"] = "cfw_pre_intl"

					if fmt.Sprint(modifyInstanceRequest["SubscriptionType"]) == "PayAsYouGo" {
						modifyInstanceRequest["ProductType"] = "cfw_elasticity_public_intl"
					}

					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
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
		addDebug(action, response, modifyInstanceRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("payment_type")
		//d.SetPartial("cfw_service")
		d.SetPartial("fw_vpc_number")
		d.SetPartial("ip_number")
		d.SetPartial("cfw_log_storage")
		d.SetPartial("cfw_log")
		d.SetPartial("band_width")
		d.SetPartial("spec")
		d.SetPartial("instance_count")
	}
	d.Partial(false)
	return resourceAliCloudCloudFirewallInstanceRead(d, meta)
}
func resourceAliCloudCloudFirewallInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudCloudFirewallInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	action := "ReleasePostInstance"
	var response map[string]interface{}

	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if fmt.Sprint(response["Message"]) == "not buy user" {
			conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
			return resource.RetryableError(fmt.Errorf("%s", response))
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

	if fmt.Sprint(response["Success"]) == "false" || fmt.Sprint(response["ReleaseStatus"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}

func convertCloudFirewallInstanceVersion(source string) interface{} {
	switch source {
	case "premium_version":
		return 2
	case "enterprise_version":
		return 3
	case "ultimate_version":
		return 4
	}

	return source
}

func convertCloudFirewallInstanceRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}

	return source
}

func convertCloudFirewallInstanceRenewalDurationUnitRequest(source interface{}) interface{} {
	switch source {
	case "Month":
		return "M"
	case "Year":
		return "Y"
	}

	return source
}
