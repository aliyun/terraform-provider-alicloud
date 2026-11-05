package ga

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaAcceleratorCreate,
		Read:   resourceAliCloudGaAcceleratorRead,
		Update: resourceAliCloudGaAcceleratorUpdate,
		Delete: resourceAliCloudGaAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"1", "2", "3", "5", "8", "10"}, false),
			},
			"bandwidth_billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"BandwidthPackage", "CDT"}, false),
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(200, 5000),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"cross_border_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cross_border_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"bgpPro", "private"}, false),
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(1, 9),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Month", "Year"}, false),
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"auto_renew": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"renewal_status"},
			},
			"renewal_status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"auto_renew"},
				ValidateFunc: validators.StringInSlice([]string{
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
			"ip_set_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validators.StringInSlice([]string{"UserDefine", "Anycast"}, false),
						},
					},
				},
			},
			"ddos_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `ddos_id` has been deprecated from provider version 2.3.0. Use `ddos_config_list` instead.",
			},
			"ddos_region_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `ddos_region_id` has been deprecated from provider version 2.3.0. Use `ddos_config_list` instead.",
			},
			"ddos_config_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddos_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ddos_region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": tags.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateAccelerator"
	request := make(map[string]interface{})
	var err error
	// there is an api bug that the name can not effect
	//if v, ok := d.GetOk("accelerator_name"); ok {
	//	request["Name"] = v
	//}

	request["RegionId"] = client.RegionId
	request["AutoPay"] = d.Get("auto_pay")
	request["ClientToken"] = helper.BuildClientToken("CreateAccelerator")

	request["AutoRenew"] = d.Get("auto_renew")
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

	if v, ok := d.GetOk("ip_set_config"); ok {
		ipSetConfigList := v.([]interface{})
		if len(ipSetConfigList) > 0 {
			ipSetConfigMap := ipSetConfigList[0].(map[string]interface{})
			if accessMode, ok := ipSetConfigMap["access_mode"]; ok && accessMode.(string) != "" {
				request["IpSetConfig.AccessMode"] = accessMode
			}
		}
	}

	if v, ok := d.GetOk("ip_set_config"); ok {
		ipSetConfigList := v.([]interface{})
		if len(ipSetConfigList) > 0 {
			ipSetConfigMap := ipSetConfigList[0].(map[string]interface{})
			if mode, ok := ipSetConfigMap["access_mode"]; ok && mode.(string) == "Anycast" {
				if v, ok := d.GetOk("bandwidth"); ok {
					request["Bandwidth"] = v
				} else {
					return errs.WrapErrorf(nil, errs.DefaultErrorMsg, "alicloud_ga_accelerator", "bandwidth is required when access_mode is set to Anycast", errs.AlibabaCloudSdkGoERROR)
				}
			} else if _, ok := d.GetOk("bandwidth"); ok {
				return errs.WrapErrorf(nil, errs.DefaultErrorMsg, "alicloud_ga_accelerator", "bandwidth is valid only when access_mode is set to Anycast", errs.AlibabaCloudSdkGoERROR)
			}
		}
	} else if _, ok := d.GetOk("bandwidth"); ok {
		return errs.WrapErrorf(nil, errs.DefaultErrorMsg, "alicloud_ga_accelerator", "bandwidth is valid only when ip_set_config.access_mode is set to Anycast", errs.AlibabaCloudSdkGoERROR)
	}

	response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_accelerator", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	d.SetId(fmt.Sprint(response["AcceleratorId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaAcceleratorUpdate(d, meta)
}

func resourceAliCloudGaAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaAccelerator(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_accelerator gaService.DescribeGaAccelerator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("spec", object["Spec"])
	d.Set("bandwidth_billing_type", object["BandwidthBillingType"])
	if v, ok := object["Bandwidth"]; ok {
		bandwidthVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("bandwidth", bandwidthVal)
	}
	d.Set("payment_type", convertGaAcceleratorPaymentTypeResponse(object["InstanceChargeType"]))
	d.Set("cross_border_status", object["CrossBorderStatus"])
	d.Set("cross_border_mode", object["CrossBorderMode"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("accelerator_name", object["Name"])
	d.Set("description", object["Description"])
	if ipSetConfig, ok := object["IpSetConfig"]; ok {
		ipSetConfigMaps := make([]map[string]interface{}, 0)
		ipSetConfigMap := make(map[string]interface{})
		if ipSetConfigMapData, ok := ipSetConfig.(map[string]interface{}); ok {
			if accessMode, ok := ipSetConfigMapData["AccessMode"]; ok {
				ipSetConfigMap["access_mode"] = accessMode
			}
		}
		ipSetConfigMaps = append(ipSetConfigMaps, ipSetConfigMap)
		d.Set("ip_set_config", ipSetConfigMaps)
	}
	d.Set("status", object["State"])

	ddosConfigs := flattenGaAcceleratorDdosConfigList(object["DdosConfigList"])
	d.Set("ddos_config_list", ddosConfigs)

	if v, ok := object["DdosId"]; ok && fmt.Sprint(v) != "" {
		d.Set("ddos_id", v)
	} else if len(ddosConfigs) > 0 {
		d.Set("ddos_id", ddosConfigs[0]["ddos_id"])
	} else {
		d.Set("ddos_id", "")
	}

	if len(ddosConfigs) > 0 {
		d.Set("ddos_region_id", ddosConfigs[0]["ddos_region_id"])
	} else {
		d.Set("ddos_region_id", "")
	}

	describeAcceleratorAutoRenewAttributeObject, err := gaService.DescribeAcceleratorAutoRenewAttribute(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("renewal_status", describeAcceleratorAutoRenewAttributeObject["RenewalStatus"])
	d.Set("auto_renew", fmt.Sprint(describeAcceleratorAutoRenewAttributeObject["RenewalStatus"]) == string(RenewAutoRenewal))

	if v, ok := describeAcceleratorAutoRenewAttributeObject["AutoRenewDuration"]; ok && fmt.Sprint(v) != "0" {
		auto_renew_durationVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("auto_renew_duration", auto_renew_durationVal)
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "accelerator")
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "accelerator"); err != nil {
			return errs.WrapError(err)
		}
	}

	update := false
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   helper.BuildClientToken("UpdateAcceleratorAutoRenewAttribute"),
		"AcceleratorId": d.Id(),
	}

	autoRenewChanged := !d.IsNewResource() && d.HasChange("auto_renew")
	if !autoRenewChanged {
		if v, ok := d.GetOk("renewal_status"); ok {
			request["RenewalStatus"] = v
		}
	}
	if d.HasChange("renewal_status") {
		update = true
	}

	if d.HasChange("auto_renew_duration") {
		update = true
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}

	if autoRenewChanged {
		update = true
		request["AutoRenew"] = d.Get("auto_renew")
	}

	if update {
		if _, ok := request["RenewalStatus"]; !ok {
			if _, ok := request["AutoRenew"]; !ok {
				request["AutoRenew"] = d.Get("auto_renew")
			}
		}
	}

	if update {
		action := "UpdateAcceleratorAutoRenewAttribute"
		var err error

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
	updateAcceleratorReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   helper.BuildClientToken("UpdateAccelerator"),
		"AcceleratorId": d.Id(),
		"AutoPay":       d.Get("auto_pay"),
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

	if d.HasChange("bandwidth") {
		update = true
		if v, ok := d.GetOk("bandwidth"); ok {
			if accessMode, ok := d.GetOk("ip_set_config"); ok {
				ipSetConfigList := accessMode.([]interface{})
				if len(ipSetConfigList) > 0 {
					ipSetConfigMap := ipSetConfigList[0].(map[string]interface{})
					if mode, ok := ipSetConfigMap["access_mode"]; ok && mode.(string) == "Anycast" {
						updateAcceleratorReq["Bandwidth"] = v
					} else {
						return errs.WrapErrorf(nil, errs.DefaultErrorMsg, "alicloud_ga_accelerator", "bandwidth is valid only when access_mode is set to Anycast", errs.AlibabaCloudSdkGoERROR)
					}
				}
			} else {
				return errs.WrapErrorf(nil, errs.DefaultErrorMsg, "alicloud_ga_accelerator", "bandwidth is valid only when ip_set_config.access_mode is set to Anycast", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	if update {
		if v, ok := d.GetOkExists("auto_use_coupon"); ok {
			updateAcceleratorReq["AutoUseCoupon"] = v
		}

		action := "UpdateAccelerator"
		var err error

		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorReq, true)
		logs.AddDebug(action, response, updateAcceleratorReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	updateAcceleratorCrossBorderStatusReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   helper.BuildClientToken("UpdateAcceleratorCrossBorderStatus"),
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

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorCrossBorderStatusReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateAcceleratorCrossBorderStatusReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	updateAcceleratorCrossBorderModeReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   helper.BuildClientToken("UpdateAcceleratorCrossBorderMode"),
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

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateAcceleratorCrossBorderModeReq, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateAcceleratorCrossBorderModeReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  helper.BuildClientToken("ChangeResourceGroup"),
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

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	if d.HasChange("ddos_config_list") {
		oldDdosConfigList, newDdosConfigList := d.GetChange("ddos_config_list")

		oldDdosConfigs := gaAcceleratorDdosConfigsFromSet(oldDdosConfigList)
		newDdosConfigs := gaAcceleratorDdosConfigsFromSet(newDdosConfigList)

		addDdosConfigs, removeDdosConfigs := gaAcceleratorDdosConfigDiff(oldDdosConfigs, newDdosConfigs)

		if len(removeDdosConfigs) > 0 {
			detachDdosReq := map[string]interface{}{
				"RegionId":      client.RegionId,
				"ClientToken":   helper.BuildClientToken("DetachDdosFromAccelerator"),
				"AcceleratorId": d.Id(),
			}
			for k, ddosConfig := range removeDdosConfigs {
				detachDdosReq[fmt.Sprintf("DdosConfigList.%d.DdosId", k+1)] = ddosConfig["ddos_id"]
				detachDdosReq[fmt.Sprintf("DdosConfigList.%d.DdosRegionId", k+1)] = ddosConfig["ddos_region_id"]
			}

			action := "DetachDdosFromAccelerator"
			var err error
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, detachDdosReq, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"Accelerator.NotBindingDdos"}) {
						return nil
					}
					if errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, detachDdosReq)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		if len(addDdosConfigs) > 0 {
			attachDdosReq := map[string]interface{}{
				"RegionId":      client.RegionId,
				"ClientToken":   helper.BuildClientToken("AttachDdosToAccelerator"),
				"AcceleratorId": d.Id(),
			}
			for k, ddosConfig := range addDdosConfigs {
				attachDdosReq[fmt.Sprintf("DdosConfigList.%d.DdosId", k+1)] = ddosConfig["ddos_id"]
				attachDdosReq[fmt.Sprintf("DdosConfigList.%d.DdosRegionId", k+1)] = ddosConfig["ddos_region_id"]
			}

			action := "AttachDdosToAccelerator"
			var err error
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, attachDdosReq, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, attachDdosReq)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}

	d.Partial(false)

	return resourceAliCloudGaAcceleratorRead(d, meta)
}

func resourceAliCloudGaAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteAccelerator"
	var response map[string]interface{}

	var err error

	object, err := gaService.DescribeGaAccelerator(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if fmt.Sprint(object["InstanceChargeType"]) == "PREPAY" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudGaAccelerator. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) || helper.NeedRetry(err) {
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

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

func gaAcceleratorDdosConfigsFromSet(v interface{}) []map[string]string {
	set, ok := v.(*schema.Set)
	if !ok || set.Len() == 0 {
		return nil
	}
	configs := make([]map[string]string, 0, set.Len())
	for _, item := range set.List() {
		config := item.(map[string]interface{})
		configs = append(configs, map[string]string{
			"ddos_id":        fmt.Sprint(config["ddos_id"]),
			"ddos_region_id": fmt.Sprint(config["ddos_region_id"]),
		})
	}
	return configs
}

func flattenGaAcceleratorDdosConfigList(v interface{}) []map[string]interface{} {
	items, ok := v.([]interface{})
	if !ok {
		return nil
	}
	configs := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		config, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		configs = append(configs, map[string]interface{}{
			"ddos_id":        fmt.Sprint(config["DdosId"]),
			"ddos_region_id": fmt.Sprint(config["DdosRegionId"]),
		})
	}
	return configs
}

func gaAcceleratorDdosConfigDiff(oldConfigs, newConfigs []map[string]string) (add, remove []map[string]string) {
	oldKeys := make(map[string]struct{}, len(oldConfigs))
	for _, config := range oldConfigs {
		oldKeys[config["ddos_id"]+"|"+config["ddos_region_id"]] = struct{}{}
	}
	newKeys := make(map[string]struct{}, len(newConfigs))
	for _, config := range newConfigs {
		newKeys[config["ddos_id"]+"|"+config["ddos_region_id"]] = struct{}{}
	}
	for _, config := range newConfigs {
		if _, ok := oldKeys[config["ddos_id"]+"|"+config["ddos_region_id"]]; !ok {
			add = append(add, config)
		}
	}
	for _, config := range oldConfigs {
		if _, ok := newKeys[config["ddos_id"]+"|"+config["ddos_region_id"]]; !ok {
			remove = append(remove, config)
		}
	}
	return add, remove
}
