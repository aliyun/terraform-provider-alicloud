package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaCacheReserveInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCacheReserveInstanceCreate,
		Read:   resourceAliCloudEsaCacheReserveInstanceRead,
		Update: resourceAliCloudEsaCacheReserveInstanceUpdate,
		Delete: resourceAliCloudEsaCacheReserveInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cr_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"quota_gb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaCacheReserveInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PurchaseCacheReserve"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["ChargeType"] = convertEsaCacheReserveInstanceChargeTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOkExists("quota_gb"); ok {
		request["QuotaGb"] = v
	}
	if v, ok := d.GetOk("cr_region"); ok {
		request["CrRegion"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_cache_reserve_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, esaServiceV2.EsaCacheReserveInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaCacheReserveInstanceRead(d, meta)
}

func resourceAliCloudEsaCacheReserveInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCacheReserveInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_cache_reserve_instance DescribeEsaCacheReserveInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cr_region", objectRaw["CacheReserveRegion"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", convertEsaCacheReserveInstanceInstanceInfoChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("period", objectRaw["Duration"])
	d.Set("quota_gb", objectRaw["CacheReserveCapacity"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudEsaCacheReserveInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateCacheReserveSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	if _, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEsaCacheReserveInstanceChargeTypeRequest(d.Get("payment_type").(string))
	}

	if v, ok := d.GetOk("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if d.HasChange("quota_gb") {
		update = true
		request["TargetQuotaGb"] = d.Get("quota_gb")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		esaServiceV2 := EsaServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, esaServiceV2.EsaCacheReserveInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEsaCacheReserveInstanceRead(d, meta)
}

func resourceAliCloudEsaCacheReserveInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ProductCode"] = "dcdn"
	request["ProductType"] = "dcdn_dpscachereserve_public_cn"
	request["ImmediatelyRelease"] = "0"
	var endpoint string
	if client.IsInternationalAccount() {
		request["ProductType"] = "dcdn_dpscachereserve_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
				request["ProductCode"] = ""
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Second, esaServiceV2.EsaCacheReserveInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertEsaCacheReserveInstanceInstanceInfoChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PREPAY":
		return "Subscription"
	}
	return source
}
func convertEsaCacheReserveInstanceChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PREPAY"
	}
	return source
}
