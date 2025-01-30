// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEipAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEipAddressCreate,
		Read:   resourceAliCloudEipAddressRead,
		Update: resourceAliCloudEipAddressUpdate,
		Delete: resourceAliCloudEipAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(9 * time.Minute),
			Update: schema.DefaultTimeout(9 * time.Minute),
			Delete: schema.DefaultTimeout(9 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"activity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"allocation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"high_definition_monitor_log_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"OFF", "ON"}, false),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"BGP", "BGP_PRO", "ChinaTelecom", "ChinaUnicom", "ChinaMobile", "ChinaTelecom_L2", "ChinaUnicom_L2", "ChinaMobile_L2", "BGP_FinanceCloud", "BGP_International"}, false),
			},
			"log_project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_store": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"NAT", "BINDED", "MULTI_BINDED"}, false),
			},
			"netmode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"public"}, false),
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_charge_type"},
				ForceNew:      true,
				ValidateFunc:  StringInSlice([]string{"Subscription", "PayAsYouGo", "PrePaid", "PostPaid"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"public_ip_address_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_protection_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.126.0. New field 'address_name' instead.",
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'instance_charge_type' has been deprecated since provider version 1.126.0. New field 'payment_type' instead.",
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo", "PrePaid", "PostPaid"}, false),
			},
		},
	}
}

func resourceAliCloudEipAddressCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AllocateEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = d.Get("allocation_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("ip_address"); ok {
		request["IpAddress"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("address_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("netmode"); ok {
		request["Netmode"] = v
	}
	if v, ok := d.GetOk("security_protection_types"); ok {
		securityProtectionTypesMaps := v.([]interface{})
		request["SecurityProtectionTypes"] = securityProtectionTypesMaps
	}

	if v, ok := d.GetOk("isp"); ok {
		request["ISP"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("activity_id"); ok {
		request["ActivityId"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("public_ip_address_pool_id"); ok {
		request["PublicIpAddressPoolId"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = convertEipAddressInstanceChargeTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("payment_type"); ok || d.Get("instance_charge_type").(string) == "Prepaid" {
		request["InstanceChargeType"] = convertEipAddressInstanceChargeTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("zone"); ok {
		request["Zone"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := request["InstanceChargeType"]; ok && v.(string) == "PrePaid" {
		period := 1
		if v, ok := d.GetOk("period"); ok {
			period = v.(int)
		}
		request["Period"] = period
		request["PricingCycle"] = string(Month)
		if period > 9 {
			request["Period"] = period / 12
			request["PricingCycle"] = string(Year)
		}
		if v, ok := d.GetOk("pricing_cycle"); ok {
			request["PricingCycle"] = v
		}
		autoPay := true
		if v, ok := d.GetOkExists("auto_pay"); ok {
			autoPay = v.(bool)
		}
		request["AutoPay"] = autoPay
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "FrequentPurchase.EIP"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AllocationId"]))

	eipServiceV2 := EipServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eipServiceV2.EipAddressStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEipAddressUpdate(d, meta)
}

func resourceAliCloudEipAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipServiceV2 := EipServiceV2{client}

	objectRaw, err := eipServiceV2.DescribeEipAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eip_address DescribeEipAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Name"] != nil {
		d.Set("address_name", objectRaw["Name"])
	}
	if objectRaw["Bandwidth"] != nil || objectRaw["EipBandwidth"] != nil {
		if eipBandwidth, ok := objectRaw["EipBandwidth"]; ok {
			d.Set("bandwidth", eipBandwidth)
		} else {
			d.Set("bandwidth", objectRaw["Bandwidth"])
		}
	}
	if objectRaw["AllocationTime"] != nil {
		d.Set("create_time", objectRaw["AllocationTime"])
	}
	if objectRaw["DeletionProtection"] != nil {
		d.Set("deletion_protection", objectRaw["DeletionProtection"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["HDMonitorStatus"] != nil {
		d.Set("high_definition_monitor_log_status", objectRaw["HDMonitorStatus"])
	}
	if objectRaw["InternetChargeType"] != nil {
		d.Set("internet_charge_type", objectRaw["InternetChargeType"])
	}
	if objectRaw["IpAddress"] != nil {
		d.Set("ip_address", objectRaw["IpAddress"])
	}
	if objectRaw["ISP"] != nil {
		d.Set("isp", objectRaw["ISP"])
	}
	if objectRaw["Mode"] != nil {
		d.Set("mode", objectRaw["Mode"])
	}
	if objectRaw["Netmode"] != nil {
		d.Set("netmode", objectRaw["Netmode"])
	}
	if convertEipAddressEipAddressesEipAddressChargeTypeResponse(objectRaw["ChargeType"]) != nil {
		d.Set("payment_type", convertEipAddressEipAddressesEipAddressChargeTypeResponse(objectRaw["ChargeType"]))
	}
	if objectRaw["PublicIpAddressPoolId"] != nil {
		d.Set("public_ip_address_pool_id", objectRaw["PublicIpAddressPoolId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["Zone"] != nil {
		d.Set("zone", objectRaw["Zone"])
	}
	if objectRaw["AllocationId"] != nil {
		d.Set("allocation_id", objectRaw["AllocationId"])
	}

	securityProtectionType1Raw, _ := jsonpath.Get("$.SecurityProtectionTypes.SecurityProtectionType", objectRaw)
	d.Set("security_protection_types", securityProtectionType1Raw)
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	checkValue00 := d.Get("high_definition_monitor_log_status")
	if checkValue00 == "ON" {

		objectRaw, err = eipServiceV2.DescribeDescribeHighDefinitionMonitorLogAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if objectRaw["LogProject"] != nil {
			d.Set("log_project", objectRaw["LogProject"])
		}
		if objectRaw["LogStore"] != nil {
			d.Set("log_store", objectRaw["LogStore"])
		}
		if objectRaw["InstanceId"] != nil {
			d.Set("allocation_id", objectRaw["InstanceId"])
		}

	}

	d.Set("allocation_id", d.Id())

	d.Set("name", d.Get("address_name"))
	d.Set("instance_charge_type", objectRaw["ChargeType"])
	return nil
}

func resourceAliCloudEipAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyEipAddressAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AllocationId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("address_name") {
		update = true
		request["Name"] = d.Get("address_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "IncorrectEipStatus", "IncorrectStatus.ResourceStatus", "VPC_TASK_CONFLICT"}) || NeedRetry(err) {
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
		eipServiceV2 := EipServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, eipServiceV2.EipAddressStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "SetHighDefinitionMonitorLogStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("log_project") {
		update = true
		request["LogProject"] = d.Get("log_project")
	}

	if d.HasChange("log_store") {
		update = true
		request["LogStore"] = d.Get("log_store")
	}

	if d.HasChange("high_definition_monitor_log_status") {
		update = true
		request["Status"] = d.Get("high_definition_monitor_log_status")
	}

	request["InstanceType"] = "EIP"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "IncorrectEipStatus", "IncorrectInstanceStatus", "InvalidBindingStatus", "IncorrectStatus.NatGateway", "InvalidStatus.EcsStatusNotSupport", "InvalidStatus.InstanceHasBandWidth", "InvalidStatus.EniStatusNotSupport", "TaskConflict"}) || NeedRetry(err) {
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
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "EIP"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
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
	action = "DeletionProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("deletion_protection") {
		update = true
		request["ProtectionEnable"] = d.Get("deletion_protection")
	}

	request["Type"] = "EIP"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
	action = "ModifyEipForwardMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("mode") {
		update = true
		request["Mode"] = d.Get("mode")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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

	if d.HasChange("tags") {
		eipServiceV2 := EipServiceV2{client}
		if err := eipServiceV2.SetResourceTags(d, "EIP"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEipAddressRead(d, meta)
}

func resourceAliCloudEipAddressDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok || d.Get("instance_charge_type").(string) == "Prepaid" {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_eip_address which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AllocationId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "IncorrectEipStatus", "TaskConflict.AssociateGlobalAccelerationInstance"}) || NeedRetry(err) {
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

	eipServiceV2 := EipServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eipServiceV2.EipAddressStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertEipAddressPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func convertEipAddressPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertEipAddressEipAddressesEipAddressChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertEipAddressInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
