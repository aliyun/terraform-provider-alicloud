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

func ResourceAliCloudGaBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaBandwidthPackageCreate,
		Read:   resourceAliCloudGaBandwidthPackageRead,
		Update: resourceAliCloudGaBandwidthPackageUpdate,
		Delete: resourceAliCloudGaBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Basic", "CrossDomain"}, false),
			},
			"bandwidth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Advanced", "Basic", "Enhanced"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "Basic"
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Subscription",
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayBy95", "PayByTraffic"}, false),
			},
			"ratio": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"cbn_geographic_region_ida": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"cbn_geographic_region_idb": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"bandwidth_package_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateBandwidthPackage"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateBandwidthPackage")
	request["Bandwidth"] = d.Get("bandwidth")
	request["Type"] = d.Get("type")

	if v, ok := d.GetOk("bandwidth_type"); ok {
		request["BandwidthType"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertGaBandwidthPackagePaymentTypeRequest(v.(string))

		if request["ChargeType"].(string) == "PREPAY" {
			request["PricingCycle"] = "Month"
		}
	}

	if v, ok := d.GetOk("billing_type"); ok {
		request["BillingType"] = v
	}

	if v, ok := d.GetOk("ratio"); ok {
		request["Ratio"] = v
	}

	if v, ok := d.GetOk("cbn_geographic_region_ida"); ok {
		request["CbnGeographicRegionIdA"] = v
	}

	if v, ok := d.GetOk("cbn_geographic_region_idb"); ok {
		request["CbnGeographicRegionIdB"] = v
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}

	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_bandwidth_package", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BandwidthPackageId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaBandwidthPackageUpdate(d, meta)
}

func resourceAliCloudGaBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaBandwidthPackage(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_bandwidth_package gaService.DescribeGaBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	bandwidthVal, err := convert.FormatInt(object["Bandwidth"])
	if err != nil {
		return err
	}
	d.Set("bandwidth", bandwidthVal)
	d.Set("type", object["Type"])
	d.Set("bandwidth_type", object["BandwidthType"])
	d.Set("payment_type", convertGaBandwidthPackagePaymentTypeResponse(object["ChargeType"].(string)))
	d.Set("billing_type", object["BillingType"])
	ratioVal, err := convert.FormatInt(object["Ratio"])
	if err != nil {
		return err
	}
	d.Set("ratio", ratioVal)
	d.Set("cbn_geographic_region_ida", object["CbnGeographicRegionIdA"])
	d.Set("cbn_geographic_region_idb", object["CbnGeographicRegionIdB"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("bandwidth_package_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	describeBandwidthPackageAutoRenewAttributeObject, err := gaService.DescribeBandwidthPackageAutoRenewAttribute(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("renewal_status", describeBandwidthPackageAutoRenewAttributeObject["RenewalStatus"])

	if v, ok := describeBandwidthPackageAutoRenewAttributeObject["AutoRenewDuration"]; ok && fmt.Sprint(v) != "0" {
		auto_renew_durationVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("auto_renew_duration", auto_renew_durationVal)
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "bandwidthpackage")
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "bandwidthpackage"); err != nil {
			return errs.WrapError(err)
		}
	}

	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("UpdateBandwidthPackagaAutoRenewAttribute"),
		"InstanceId":  d.Id(),
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
		action := "UpdateBandwidthPackagaAutoRenewAttribute"
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
	updateBandwidthPackageReq := map[string]interface{}{
		"RegionId":           client.RegionId,
		"BandwidthPackageId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		updateBandwidthPackageReq["Bandwidth"] = d.Get("bandwidth")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth_type") {
		update = true
		updateBandwidthPackageReq["BandwidthType"] = d.Get("bandwidth_type")
	}

	if d.HasChange("bandwidth_package_name") {
		update = true
		updateBandwidthPackageReq["Name"] = d.Get("bandwidth_package_name")
	}

	if d.HasChange("description") {
		update = true
		updateBandwidthPackageReq["Description"] = d.Get("description")
	}

	if update {
		if _, ok := d.GetOkExists("auto_pay"); ok {
			updateBandwidthPackageReq["AutoPay"] = d.Get("auto_pay")
		}

		if _, ok := d.GetOkExists("auto_use_coupon"); ok {
			updateBandwidthPackageReq["AutoUseCoupon"] = d.Get("auto_use_coupon")
		}

		action := "UpdateBandwidthPackage"
		var err error

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, updateBandwidthPackageReq, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NotExist.BandwidthPackage", "StateError.Accelerator", "StateError.BandwidthPackage", "UpgradeError.BandwidthPackage", "GreaterThanGa.IpSetBandwidth", "BandwidthIllegal.BandwidthPackage"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateBandwidthPackageReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active", "binded"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  helper.BuildClientToken("ChangeResourceGroup"),
		"ResourceId":   d.Id(),
		"ResourceType": "bandwidthpackage",
	}

	if d.HasChange("resource_group_id") {
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

	d.Partial(false)

	return resourceAliCloudGaBandwidthPackageRead(d, meta)
}

func resourceAliCloudGaBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type") == "Subscription" {
		return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please unsubscribe it in the console or wait for it to expire. To remove it from the state only, run 'terraform state rm'."))
	}

	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteBandwidthPackage"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ClientToken":        helper.BuildClientToken("DeleteBandwidthPackage"),
		"BandwidthPackageId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "BindExist.BandwidthPackage"}) || helper.NeedRetry(err) {
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

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertGaBandwidthPackagePaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}

	return source
}

func convertGaBandwidthPackagePaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}

	return source
}
