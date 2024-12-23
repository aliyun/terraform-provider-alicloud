package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaBasicAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaBasicAcceleratorCreate,
		Read:   resourceAliCloudGaBasicAcceleratorRead,
		Update: resourceAliCloudGaBasicAcceleratorUpdate,
		Delete: resourceAliCloudGaBasicAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"BandwidthPackage", "CDT", "CDT95"}, false),
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
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"auto_use_coupon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 12),
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
			"basic_accelerator_name": {
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

func resourceAliCloudGaBasicAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBasicAccelerator"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBasicAccelerator")

	if v, ok := d.GetOk("bandwidth_billing_type"); ok {
		request["BandwidthBillingType"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertGaBasicAcceleratorPaymentTypeRequest(v.(string))
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}

	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}

	if v, ok := d.GetOk("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOkExists("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["PromotionOptionNo"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_basic_accelerator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AcceleratorId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBasicAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaBasicAcceleratorUpdate(d, meta)
}

func resourceAliCloudGaBasicAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBasicAccelerator(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_billing_type", object["BandwidthBillingType"])
	d.Set("payment_type", convertGaBasicAcceleratorPaymentTypeResponse(object["InstanceChargeType"]))
	d.Set("cross_border_status", object["CrossBorderStatus"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("basic_accelerator_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "basicaccelerator")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaBasicAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": d.Id(),
		"ClientToken":   buildClientToken("UpdateBasicAccelerator"),
	}

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "basicaccelerator"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("basic_accelerator_name") {
		update = true

		if v, ok := d.GetOk("basic_accelerator_name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("description") {
		update = true

		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		action := "UpdateBasicAccelerator"
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaBasicAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("basic_accelerator_name")
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaBasicAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cross_border_status")
	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  buildClientToken("ChangeResourceGroup"),
		"ResourceId":   d.Id(),
		"ResourceType": "basicaccelerator",
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

	return resourceAliCloudGaBasicAcceleratorRead(d, meta)
}

func resourceAliCloudGaBasicAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteBasicAccelerator"
	var response map[string]interface{}

	var err error

	object, err := gaService.DescribeGaBasicAccelerator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if fmt.Sprint(object["InstanceChargeType"]) == "PREPAY" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudGaBasicAccelerator prepay type. Terraform will remove this resource from the state file, however resources may remain.")
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBasicAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertGaBasicAcceleratorPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}

	return source
}

func convertGaBasicAcceleratorPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}

	return source
}
