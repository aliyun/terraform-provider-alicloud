// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
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
			"auto_use_coupon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_billing_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cross_border_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cross_border_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ddos_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ddos_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_set_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Month",
			},
			"promotion_option_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "Normal", "NotRenewal"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudGaAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccelerator"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}
	if v, ok := d.GetOk("accelerator_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("ip_set_config"); ok {
		jsonPathResult4, err := jsonpath.Get("$[0].access_mode", v)
		if err != nil {
			return WrapError(err)
		}
		request["IpSetConfig.AccessMode"] = jsonPathResult4
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}
	if v, ok := d.GetOk("bandwidth_billing_type"); ok {
		request["BandwidthBillingType"] = v
	}
	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["PromotionOptionNo"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Value"] = dataLoopTmp["tag_value"]
			dataLoopMap["Key"] = dataLoopTmp["tag_key"]
			tagMaps = append(tagMaps, dataLoopMap)
		}
		request["Tag"] = tagMaps
	}

	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOk("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_accelerator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AcceleratorId"]))

	gaServiceV2 := GaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaServiceV2.GaAcceleratorStateRefreshFunc(d.Id(), "State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaAcceleratorUpdate(d, meta)
}

func resourceAliCloudGaAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaServiceV2 := GaServiceV2{client}

	objectRaw, err := gaServiceV2.DescribeGaAccelerator(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_accelerator DescribeGaAccelerator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_name", objectRaw["Name"])
	d.Set("bandwidth_billing_type", objectRaw["BandwidthBillingType"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("cross_border_mode", objectRaw["CrossBorderMode"])
	d.Set("cross_border_status", objectRaw["CrossBorderStatus"])
	d.Set("ddos_id", objectRaw["DdosId"])
	d.Set("description", objectRaw["Description"])
	d.Set("payment_type", convertGaInstanceChargeTypeResponse(objectRaw["InstanceChargeType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("spec", objectRaw["Spec"])
	d.Set("status", objectRaw["State"])

	ipSetConfigMaps := make([]map[string]interface{}, 0)
	ipSetConfigMap := make(map[string]interface{})
	ipSetConfig1Raw := make(map[string]interface{})
	if objectRaw["IpSetConfig"] != nil {
		ipSetConfig1Raw = objectRaw["IpSetConfig"].(map[string]interface{})
	}
	if len(ipSetConfig1Raw) > 0 {
		ipSetConfigMap["access_mode"] = ipSetConfig1Raw["AccessMode"]
		ipSetConfigMaps = append(ipSetConfigMaps, ipSetConfigMap)
	}
	d.Set("ip_set_config", ipSetConfigMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = gaServiceV2.DescribeDescribeAcceleratorAutoRenewAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("renewal_status", objectRaw["RenewalStatus"])
	d.Set("auto_renew", objectRaw["AutoRenew"])
	d.Set("auto_renew_duration", objectRaw["AutoRenewDuration"])

	return nil
}

func resourceAliCloudGaAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateAccelerator"
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("accelerator_name") {
		update = true
		request["Name"] = d.Get("accelerator_name")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
		request["Spec"] = d.Get("spec")
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		gaServiceV2 := GaServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaServiceV2.GaAcceleratorStateRefreshFunc(d.Id(), "State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("accelerator_name")
		d.SetPartial("description")
		d.SetPartial("spec")
	}
	update = false
	action = "UpdateAcceleratorAutoRenewAttribute"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("renewal_status") {
		update = true
		request["RenewalStatus"] = d.Get("renewal_status")
	}

	if !d.IsNewResource() && d.HasChange("accelerator_name") {
		update = true
		request["Name"] = d.Get("accelerator_name")
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		d.SetPartial("renewal_status")
		d.SetPartial("accelerator_name")
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_duration")
	}
	update = false
	action = "ChangeResourceGroup"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "accelerator"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		d.SetPartial("resource_group_id")
	}
	update = false
	action = "AttachDdosToAccelerator"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("ddos_id") {
		update = true
		request["DdosId"] = d.Get("ddos_id")
	}

	if v, ok := d.GetOk("ddos_region_id"); ok {
		request["DdosRegionId"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		d.SetPartial("ddos_id")
	}
	update = false
	action = "DetachDdosFromAccelerator"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
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
	action = "DeployAccelerator"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
	action = "UpdateAcceleratorCrossBorderMode"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("cross_border_mode") {
		update = true
		request["CrossBorderMode"] = d.Get("cross_border_mode")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		d.SetPartial("cross_border_mode")
	}
	update = false
	action = "UpdateAcceleratorCrossBorderStatus"
	conn, err = client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("cross_border_status") {
		update = true
		request["CrossBorderStatus"] = d.Get("cross_border_status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
		d.SetPartial("cross_border_status")
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		gaServiceV2 := GaServiceV2{client}
		if err := gaServiceV2.SetResourceTags(d, "accelerator"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudGaAcceleratorRead(d, meta)
}

func resourceAliCloudGaAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccelerator"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AcceleratorId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.Accelerator"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	gaServiceV2 := GaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaServiceV2.GaAcceleratorStateRefreshFunc(d.Id(), "AcceleratorId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertGaInstanceChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PREPAY":
		return "Subscription"
	}
	return source
}
