package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAmqpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAmqpInstanceCreate,
		Read:   resourceAliCloudAmqpInstanceRead,
		Update: resourceAliCloudAmqpInstanceUpdate,
		Delete: resourceAliCloudAmqpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"professional", "enterprise", "vip", "serverless"}, true),
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_eip_tps": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("support_eip"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"max_tps": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Upgrade", "Downgrade"}, true),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, true),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3, 6, 12, 24}),
			},
			"period_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, true),
			},
			"queue_capacity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"renewal_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3, 6, 12}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, true),
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
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal", "NotRenewal"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"serverless_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("instance_type"); ok && v.(string) == "vip" {
						return false
					}
					return true
				},
			},
			"support_eip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"support_tracing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tracing_storage_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{-1, 0, 3, 7, 15}),
			},
		},
	}
}

func resourceAliCloudAmqpInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["PaymentType"] = d.Get("payment_type")
	if v, ok := d.GetOkExists("support_eip"); ok {
		request["SupportEip"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("queue_capacity"); ok {
		request["QueueCapacity"] = v
	}
	if v, ok := d.GetOk("max_eip_tps"); ok {
		request["MaxEipTps"] = v
	}
	if v, ok := d.GetOk("max_connections"); ok {
		request["MaxConnections"] = v
	}
	if v, ok := d.GetOk("storage_size"); ok {
		request["StorageSize"] = v
	}
	if v, ok := d.GetOkExists("support_tracing"); ok {
		request["SupportTracing"] = v
	}
	if v, ok := d.GetOk("tracing_storage_time"); ok {
		request["TracingStorageTime"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewStatus"] = v
	}
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		request["RenewalDurationUnit"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_cycle"); ok {
		request["PeriodCycle"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("serverless_charge_type"); ok {
		request["ServerlessChargeType"] = v
	}
	if v, ok := d.GetOk("max_tps"); ok {
		request["MaxPrivateTps"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	amqpServiceV2 := AmqpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"SERVING"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, amqpServiceV2.AmqpInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAmqpInstanceRead(d, meta)
}

func resourceAliCloudAmqpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpServiceV2 := AmqpServiceV2{client}

	objectRaw, err := amqpServiceV2.DescribeAmqpInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_instance DescribeAmqpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["OrderCreateTime"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", convertAmqpInstanceDataInstanceTypeResponse(objectRaw["InstanceType"]))
	d.Set("max_connections", objectRaw["MaxConnections"])
	d.Set("max_eip_tps", objectRaw["MaxEipTps"])
	d.Set("max_tps", objectRaw["MaxTps"])
	d.Set("queue_capacity", objectRaw["MaxQueue"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_size", objectRaw["StorageSize"])
	d.Set("support_eip", objectRaw["SupportEIP"])
	d.Set("support_tracing", objectRaw["SupportTracing"])
	d.Set("tracing_storage_time", objectRaw["TracingStorageTime"])

	if convertAmqpInstanceDataInstanceTypeResponse(objectRaw["InstanceType"]) == "SERVERLESS" {
		d.Set("payment_type", "PayAsYouGo")
		return nil
	}

	objectRaw, err = amqpServiceV2.DescribeQueryAvailableInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("renewal_duration", objectRaw["RenewalDuration"])
	d.Set("renewal_duration_unit", convertAmqpInstanceDataInstanceListRenewalDurationUnitResponse(objectRaw["RenewalDurationUnit"]))
	d.Set("renewal_status", objectRaw["RenewStatus"])

	return nil
}

func resourceAliCloudAmqpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = d.Id()
	if d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if d.HasChange("renewal_duration_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		request["RenewalPeriodUnit"] = convertAmqpInstanceRenewalPeriodUnitRequest(fmt.Sprint(v))
	}

	if d.HasChange("renewal_duration") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewalPeriod"] = v
	}

	if d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	request["ProductCode"] = "ons"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "ons_onsproxy_public_intl"
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
		d.SetPartial("renewal_duration_unit")
		d.SetPartial("renewal_duration")
		d.SetPartial("renewal_status")
	}
	update = false
	action = "UpdateInstanceName"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, false)

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
		d.SetPartial("instance_name")
	}
	update = false
	action = "UpdateInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOk("modify_type"); ok {
		request["ModifyType"] = convertAmqpInstanceModifyTypeRequest(v.(string))
	}
	if d.HasChange("instance_type") {
		update = true
	}
	if v, ok := d.GetOk("instance_type"); ok && fmt.Sprint(v) != "serverless" {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("serverless_charge_type"); ok {
		request["ServerlessChargeType"] = v
	}
	if d.HasChange("support_eip") {
		update = true
	}
	if v, ok := d.GetOkExists("support_eip"); ok {
		request["SupportEip"] = v
	}

	if d.HasChange("max_eip_tps") {
		update = true
	}
	if v, ok := d.GetOk("max_eip_tps"); ok && fmt.Sprint(v) != "-1" && d.Get("support_eip").(bool) {
		request["MaxEipTps"] = v
	}

	if d.HasChange("queue_capacity") {
		update = true
	}
	if v, ok := d.GetOk("queue_capacity"); ok && fmt.Sprint(v) != "-1" {
		request["QueueCapacity"] = v
	}

	if d.HasChange("max_connections") {
		update = true
	}
	if v, ok := d.GetOk("max_connections"); ok && fmt.Sprint(v) != "-1" {
		request["MaxConnections"] = v
	}

	if d.HasChange("storage_size") {
		update = true
	}
	if v, ok := d.GetOk("storage_size"); ok && fmt.Sprint(v) != "-1" {
		request["StorageSize"] = v
	}

	if d.HasChange("support_tracing") {
		update = true
	}
	if v, ok := d.GetOk("support_tracing"); ok && fmt.Sprint(v) != "-1" {
		request["SupportTracing"] = v
	}

	if d.HasChange("tracing_storage_time") {
		update = true
	}
	if v, ok := d.GetOk("tracing_storage_time"); ok && fmt.Sprint(v) != "-1" {
		request["TracingStorageTime"] = v
	}

	if d.HasChange("max_tps") {
		update = true
	}
	if v, ok := d.GetOk("max_tps"); ok && fmt.Sprint(v) != "-1" {
		request["MaxPrivateTps"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)
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
		code, _ := jsonpath.Get("$.Code", response)
		if fmt.Sprint(code) != "200" {
			log.Printf("[DEBUG] Resource alicloud_amqp_instance UpdateInstance Failed!!! %s", response)
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_instance", action, AlibabaCloudSdkGoERROR, response)
		}
		amqpServiceV2 := AmqpServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"SERVING"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, amqpServiceV2.AmqpInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_type")
		d.SetPartial("support_eip")
		d.SetPartial("max_eip_tps")
		d.SetPartial("queue_capacity")
		d.SetPartial("max_connections")
		d.SetPartial("storage_size")
		d.SetPartial("support_tracing")
		d.SetPartial("tracing_storage_time")
		d.SetPartial("max_tps")
	}

	d.Partial(false)
	return resourceAliCloudAmqpInstanceRead(d, meta)
}

func resourceAliCloudAmqpInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "PayAsYouGo" {
			log.Printf("[WARN] Cannot destroy resource alicloud_amqp_instance which payment_type valued PayAsYouGo. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	query["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	request["ProductCode"] = "ons"
	request["ProductType"] = "ons_onsproxy_pre"
	if client.IsInternationalAccount() {
		request["ProductType"] = "ons_onsproxy_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "ons_onsproxy_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	amqpServiceV2 := AmqpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, amqpServiceV2.AmqpInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertAmqpInstanceSupportEipRequest(source interface{}) interface{} {
	switch source {
	case false:
		return "eip_false"
	case true:
		return "eip_true"
	}
	return ""
}
func convertAmqpInstanceInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "PROFESSIONAL":
		return "professional"
	case "ENTERPRISE":
		return "enterprise"
	case "VIP":
		return "vip"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitRequest(source interface{}) interface{} {
	switch source {
	case "Month":
		return "M"
	case "Year":
		return "Y"
	}
	return source
}

func convertAmqpInstanceDataInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "PROFESSIONAL":
		return "professional"
	case "VIP":
		return "vip"
	case "ENTERPRISE":
		return "enterprise"
	case "SERVERLESS":
		return "serverless"
	}
	return source
}
func convertAmqpInstanceDataInstanceListRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}
	return source
}
func convertAmqpInstanceRenewalPeriodUnitRequest(source interface{}) interface{} {
	switch source {
	case "Month":
		return "M"
	case "Year":
		return "Y"
	}
	return source
}
func convertAmqpInstanceModifyTypeRequest(source interface{}) interface{} {
	switch source {
	case "Downgrade":
		return "DOWNGRADE"
	case "Upgrade":
		return "UPGRADE"
	}
	return source
}
