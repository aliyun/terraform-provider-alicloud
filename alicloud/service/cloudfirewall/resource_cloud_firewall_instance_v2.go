package cloudfirewall

import (
	"context"
	"fmt"
	"log"
	"time"

	jsonpath "github.com/PaesslerAG/jsonpath"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	validators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cloudFirewallInstancePaymentTypeChangeError = "cannot change payment_type of Cloud Firewall instance %s from %s to %s: converting the billing method of an existing instance is not supported by the API. Convert it in the Cloud Firewall console (Overview -> More -> Convert Pay-As-You-Go to Subscription, or Convert Subscription to Pay-As-You-Go) by following https://help.aliyun.com/cloud-firewall/cloudfirewall/product-overview/smoothly-change-the-billing-method-between-pay-as-you-go-and-subscription, then align payment_type, product_type and spec in the configuration or import the converted instance. The existing instance is kept unchanged (no destroy or recreate)."

func ResourceCloudFirewallInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceV2Create,
		Read:   resourceInstanceV2Read,
		Update: resourceInstanceV2Update,
		Delete: resourceInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceInstanceV2CustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cfw_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Upgrade", "Downgrade"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"product_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"cfw"}, false),
			},
			"product_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"cfw_elasticity_public_cn", "cfw_elasticity_public_intl", "cfw_sub_public_cn", "cfw_sub_public_intl"}, false),
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
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"payg_version", "premium_version", "enterprise_version", "ultimate_version"}, false),
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

// resourceInstanceV2CustomizeDiff rejects a billing change at plan time when
// product_type or spec (both ForceNew) change in the same plan. Such a
// combination requires a replacement, which would release the existing
// instance and bypass the Update guard instead of failing with the console
// conversion guidance.
func resourceInstanceV2CustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Id() == "" || !d.HasChange("payment_type") {
		return nil
	}
	if !d.HasChange("product_type") && !d.HasChange("spec") {
		return nil
	}
	oldPaymentType, newPaymentType := d.GetChange("payment_type")
	return errs.Error(cloudFirewallInstancePaymentTypeChangeError, d.Id(), oldPaymentType, newPaymentType)
}

func resourceInstanceV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

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
	request["Parameter"] = parameterMapList

	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")
	request["ProductType"] = d.Get("product_type")

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable", errs.NotFoundArticle}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_instance_v2", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	cloudFirewallServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudFirewallServiceV2.InstanceStateRefreshFuncWithApi(d.Id(), "$.InstanceStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceInstanceV2Update(d, meta)
}

func resourceInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance_v2 DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("cfw_log", objectRaw["LogStatus"])
	d.Set("sdl", convertCloudFirewallInstanceSdlResponse(objectRaw["Sdl"]))
	d.Set("spec", convertCloudFirewallInstanceVersionResponse(objectRaw["Version"]))
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("user_status", fmt.Sprint(objectRaw["UserStatus"]))

	return nil
}

func resourceInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	// The billing method of an existing instance cannot be changed through the
	// API: the console conversion is the only supported path. Reconcile against
	// the actual billing status so a stale state converges without any mutation,
	// and reject the change deterministically otherwise. Creating an instance
	// passes through this function too, hence the IsNewResource guard.
	if !d.IsNewResource() && d.HasChange("payment_type") {
		objectRaw, err := cloudFirewallServiceV2.DescribeInstance(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		actualPaymentType := fmt.Sprint(objectRaw["SubscriptionType"])
		expectedPaymentType := d.Get("payment_type").(string)
		if actualPaymentType != expectedPaymentType {
			return errs.Error(cloudFirewallInstancePaymentTypeChangeError, d.Id(), actualPaymentType, expectedPaymentType)
		}
	}

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
		return errs.WrapError(err)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	var endpoint string

	request["InstanceId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	request["ProductType"] = d.Get("product_type")
	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")

	if v, ok := d.GetOk("modify_type"); ok {
		request["ModifyType"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)

	if !d.IsNewResource() && d.HasChange("cfw_log") {
		update = true
	}
	if v, ok := d.GetOkExists("cfw_log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLog",
			"Value": v,
		})
	}

	request["Parameter"] = parameterMapList

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}

			if fmt.Sprint(response["Code"]) == "SYSTEM.CONCURRENT_OPERATE" {
				wait()
				return retry.RetryableError(fmt.Errorf("%s", response))
			}

			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if d.HasChange("cfw_log") {
			stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("cfw_log"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudFirewallServiceV2.InstanceStateRefreshFuncWithApi(d.Id(), "LogStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		bssOpenApiService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Paid"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, bssOpenApiService.InstanceOrderDetailStateRefreshFunc(fmt.Sprint(response["Data"].(map[string]interface{})["OrderId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceInstanceV2Read(d, meta)
}

func resourceInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var action string
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	if d.Get("payment_type").(string) == "Subscription" {
		// A prepaid instance is released by refunding the subscription through
		// BSS; the Cloud Firewall release API only handles pay-as-you-go ones.
		action = "RefundInstance"
		request["ClientToken"] = helper.BuildClientToken(action)
		request["ImmediatelyRelease"] = "1"
		request["ProductCode"] = d.Get("product_code")
		request["ProductType"] = d.Get("product_type")
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
	} else {
		action = "ReleasePostInstance"
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
	}
	logs.AddDebug(action, response, request)

	if err != nil {
		// ResourceNotExists and ExistRefundingOrderError both mean the instance is
		// already released or being refunded; the wait below still has to confirm
		// that the instance is gone instead of dropping it from the state blindly.
		if !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"ResourceNotExists", "ExistRefundingOrderError"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	cloudFirewallServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudFirewallServiceV2.InstanceStateRefreshFunc(d.Id(), "#$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

func convertCloudFirewallInstanceSpecRequest(source interface{}) interface{} {
	switch source {
	case "premium_version":
		return 2
	case "enterprise_version":
		return 3
	case "ultimate_version":
		return 4
	case "payg_version":
		return 10
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
