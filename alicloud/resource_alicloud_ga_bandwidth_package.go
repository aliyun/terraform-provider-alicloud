package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaBandwidthPackageCreate,
		Read:   resourceAlicloudGaBandwidthPackageRead,
		Update: resourceAlicloudGaBandwidthPackageUpdate,
		Delete: resourceAlicloudGaBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"bandwidth": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"bandwidth_package_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"bandwidth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Advanced", "Basic", "Enhanced"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "Basic"
				},
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayBy95", "PayByTraffic"}, false),
			},
			"cbn_geographic_region_ida": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "China-mainland",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"cbn_geographic_region_idb": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Global",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"duration": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Default:      "Subscription",
			},
			"ratio": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "CrossDomain"}, false),
			},
		},
	}
}

func resourceAlicloudGaBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("bandwidth_type"); ok {
		request["BandwidthType"] = v
	}
	if v, ok := d.GetOk("billing_type"); ok {
		request["BillingType"] = v
	}
	if v, ok := d.GetOk("cbn_geographic_region_ida"); ok {
		request["CbnGeographicRegionIdA"] = v
	}
	if v, ok := d.GetOk("cbn_geographic_region_idb"); ok {
		request["CbnGeographicRegionIdB"] = v
	}

	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertGaBandwidthPackagePaymentTypeRequest(v.(string))
		if request["ChargeType"].(string) == "PREPAY" {
			request["PricingCycle"] = "Month"
		}
	}
	if v, ok := d.GetOk("ratio"); ok {
		request["Ratio"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	var response map[string]interface{}
	request["ClientToken"] = buildClientToken("CreateBandwidthPackage")
	action := "CreateBandwidthPackage"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_bandwidth_package", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BandwidthPackageId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaBandwidthPackageUpdate(d, meta)
}
func resourceAlicloudGaBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_bandwidth_package gaService.DescribeGaBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("bandwidth_package_name", object["Name"])
	d.Set("bandwidth_type", object["BandwidthType"])
	d.Set("cbn_geographic_region_ida", object["CbnGeographicRegionIdA"])
	d.Set("cbn_geographic_region_idb", object["CbnGeographicRegionIdB"])
	d.Set("description", object["Description"])
	d.Set("payment_type", convertGaBandwidthPackagePaymentTypeResponse(object["ChargeType"].(string)))
	d.Set("status", object["State"])
	d.Set("type", object["Type"])
	if val, ok := d.GetOkExists("auto_use_coupon"); ok {
		d.Set("auto_use_coupon", val)
	}
	return nil
}
func resourceAlicloudGaBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
		"RegionId":           client.RegionId,
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		if v, ok := d.GetOk("bandwidth"); ok {
			request["Bandwidth"] = v
		}
	}
	if d.HasChange("bandwidth_package_name") {
		update = true
		if v, ok := d.GetOk("bandwidth_package_name"); ok {
			request["Name"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("bandwidth_type") {
		update = true
		if v, ok := d.GetOk("bandwidth_type"); ok {
			request["BandwidthType"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		action := "UpdateBandwidthPackage"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"NotExist.BandwidthPackage", "StateError.BandwidthPackage", "UpgradeError.BandwidthPackage"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"active", "binded"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudGaBandwidthPackageRead(d, meta)
}

func resourceAlicloudGaBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type") == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceAlicloudGaBandwidthPackage prepay type. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
		"RegionId":           client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("DeleteBandwidthPackage")
	action := "DeleteBandwidthPackage"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.BandwidthPackage"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
