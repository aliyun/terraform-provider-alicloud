package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaAcceleratorCreate,
		Read:   resourceAliCloudGaAcceleratorRead,
		Update: resourceAliCloudGaAcceleratorUpdate,
		Delete: resourceAliCloudGaAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"1", "2", "3", "5", "8", "10"}, false),
			},
			"bandwidth_billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"BandwidthPackage", "CDT"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"cross_border_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cross_border_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"bgpPro", "private"}, false),
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 9),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
			},
			"auto_renew_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"promotion_option_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"accelerator_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateAccelerator"
	request := make(map[string]interface{})
	var err error
	// there is an api bug that the name can not effect
	//if v, ok := d.GetOk("accelerator_name"); ok {
	//	request["Name"] = v
	//}

	request["RegionId"] = client.RegionId
	request["AutoPay"] = true
	request["ClientToken"] = buildClientToken("CreateAccelerator")

	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}

	if v, ok := d.GetOk("bandwidth_billing_type"); ok {
		request["BandwidthBillingType"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertGaAcceleratorPaymentTypeRequest(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}

	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	} else {
		request["PricingCycle"] = "Month"
	}

	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["PromotionOptionNo"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_accelerator", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["AcceleratorId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaAcceleratorUpdate(d, meta)
}

func resourceAliCloudGaAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaAccelerator(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_accelerator gaService.DescribeGaAccelerator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("spec", object["Spec"])
	d.Set("bandwidth_billing_type", object["BandwidthBillingType"])
	d.Set("payment_type", convertGaAcceleratorPaymentTypeResponse(object["InstanceChargeType"]))
	d.Set("cross_border_status", object["CrossBorderStatus"])
	d.Set("cross_border_mode", object["CrossBorderMode"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("accelerator_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	describeAcceleratorAutoRenewAttributeObject, err := gaService.DescribeAcceleratorAutoRenewAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("renewal_status", describeAcceleratorAutoRenewAttributeObject["RenewalStatus"])

	if v, ok := describeAcceleratorAutoRenewAttributeObject["AutoRenewDuration"]; ok && fmt.Sprint(v) != "0" {
		d.Set("auto_renew_duration", formatInt(v))
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "accelerator")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "accelerator"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateAcceleratorAutoRenewAttribute"),
		"AcceleratorId": d.Id(),
	}

	if d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if d.HasChange("auto_renew_duration") {
		update = true
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}

	if update {
		action := "UpdateAcceleratorAutoRenewAttribute"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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

		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_duration")
	}

	update = false
	updateAcceleratorReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateAccelerator"),
		"AcceleratorId": d.Id(),
		"AutoPay":       true,
	}

	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
		updateAcceleratorReq["Spec"] = d.Get("spec")
	}

	if d.HasChange("accelerator_name") {
		update = true
		if v, ok := d.GetOk("accelerator_name"); ok {
			updateAcceleratorReq["Name"] = v
		}
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			updateAcceleratorReq["Description"] = v
		}
	}

	if update {
		if v, ok := d.GetOkExists("auto_use_coupon"); ok {
			updateAcceleratorReq["AutoUseCoupon"] = v
		}

		action := "UpdateAccelerator"
		var err error

		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorReq, true)
		addDebug(action, response, updateAcceleratorReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("spec")
		d.SetPartial("accelerator_name")
		d.SetPartial("description")
	}

	update = false
	updateAcceleratorCrossBorderStatusReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateAcceleratorCrossBorderStatus"),
		"AcceleratorId": d.Id(),
	}

	if d.HasChange("cross_border_status") {
		update = true
	}
	if v, ok := d.GetOkExists("cross_border_status"); ok {
		updateAcceleratorCrossBorderStatusReq["CrossBorderStatus"] = v
	}

	if update {
		action := "UpdateAcceleratorCrossBorderStatus"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorCrossBorderStatusReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAcceleratorCrossBorderStatusReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cross_border_status")
	}

	update = false
	updateAcceleratorCrossBorderModeReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateAcceleratorCrossBorderMode"),
		"AcceleratorId": d.Id(),
	}

	if d.HasChange("cross_border_mode") {
		update = true
	}
	if v, ok := d.GetOk("cross_border_mode"); ok {
		updateAcceleratorCrossBorderModeReq["CrossBorderMode"] = v
	}

	if update {
		action := "UpdateAcceleratorCrossBorderMode"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorCrossBorderModeReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAcceleratorCrossBorderModeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cross_border_mode")
	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  buildClientToken("ChangeResourceGroup"),
		"ResourceId":   d.Id(),
		"ResourceType": "accelerator",
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAliCloudGaAcceleratorRead(d, meta)
}

func resourceAliCloudGaAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteAccelerator"
	var response map[string]interface{}

	var err error

	object, err := gaService.DescribeGaAccelerator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if fmt.Sprint(object["InstanceChargeType"]) == "PREPAY" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudGaAccelerator. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertGaAcceleratorPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}

	return source
}

func convertGaAcceleratorPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}

	return source
}
