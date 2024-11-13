// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudEsaRatePlanInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRatePlanInstanceCreate,
		Read:   resourceAliCloudEsaRatePlanInstanceRead,
		Update: resourceAliCloudEsaRatePlanInstanceUpdate,
		Delete: resourceAliCloudEsaRatePlanInstanceDelete,
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
			"coverage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEsaRatePlanInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PurchaseRatePlan"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	if v, ok := d.GetOk("plan_name"); ok {
		request["PlanName"] = v
	}
	if v, ok := d.GetOk("coverage"); ok {
		request["Coverage"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEsaRatePlanInstanceChargeTypeRequest(v.(string))
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-09-10"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_rate_plan_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, esaServiceV2.EsaRatePlanInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaRatePlanInstanceRead(d, meta)
}

func resourceAliCloudEsaRatePlanInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRatePlanInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_rate_plan_instance DescribeEsaRatePlanInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["BillingMode"] != nil {
		d.Set("payment_type", convertEsaRatePlanInstanceInstanceInfoBillingModeResponse(objectRaw["BillingMode"]))
	}
	if objectRaw["PlanName"] != nil {
		d.Set("plan_name", objectRaw["PlanName"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	objectRaw, err = esaServiceV2.DescribeDescribeRatePlanInstanceStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["InstanceStatus"] != nil {
		d.Set("instance_status", objectRaw["InstanceStatus"])
	}

	return nil
}

func resourceAliCloudEsaRatePlanInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateRatePlanSpec"
	conn, err := client.NewEsaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["OrderType"] = "UPGRADE"
	if d.HasChange("plan_name") {
		update = true
		request["TargetPlanName"] = d.Get("plan_name")
	}

	if d.HasChange("payment_type") {
		update = true
	}
	request["ChargeType"] = convertEsaRatePlanInstanceChargeTypeRequest(d.Get("payment_type").(string))

	request["AutoPay"] = true
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-09-10"), StringPointer("AK"), query, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, esaServiceV2.EsaRatePlanInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEsaRatePlanInstanceRead(d, meta)
}

func resourceAliCloudEsaRatePlanInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Rate Plan Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertEsaRatePlanInstanceInstanceInfoBillingModeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PREPAY":
		return "Subscription"
	case "POSTPAY":
		return "PayAsYouGo"
	}
	return source
}
func convertEsaRatePlanInstanceChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}
