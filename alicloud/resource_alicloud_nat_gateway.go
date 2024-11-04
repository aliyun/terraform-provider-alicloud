package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNatGatewayCreate,
		Read:   resourceAliCloudNatGatewayRead,
		Update: resourceAliCloudNatGatewayUpdate,
		Delete: resourceAliCloudNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"forward_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayByLcu", "PayBySpec"}, false),
			},
			"nat_gateway_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nat_gateway_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.121.0. New field `nat_gateway_name` instead.",
			},
			"nat_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enhanced", "Normal"}, false),
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ConflictsWith: []string{"instance_charge_type"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"payment_type"},
				Deprecated:    "Field `instance_charge_type` has been deprecated from provider version 1.121.0. New field `payment_type` instead.",
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},
			"bandwidth_packages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Removed:  "Field `ip_count` has been removed from provider version 1.121.0.",
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
							Removed:  "Field `bandwidth` has been removed from provider version 1.121.0.",
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Removed:  "Field `zone` has been removed from provider version 1.121.0.",
						},
						"public_ip_addresses": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field `public_ip_addresses` has been removed from provider version 1.121.0.",
						},
					},
				},
				MaxItems: 4,
				Optional: true,
				Removed:  "Field `bandwidth_packages` has been removed from provider version 1.121.0.",
			},
			"bandwidth_package_ids": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field `bandwidth_package_ids` has been removed from provider version 1.121.0.",
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `spec` has been removed from provider version 1.121.0. New field `specification` instead.",
			},
			"snat_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Large", "Middle", "Small", "XLarge.1"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("internet_charge_type").(string) == "PayByLcu"
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("nat_type").(string) != "Enhanced"
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"internet", "intranet"}, false),
			},
			"eip_bind_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"MULTI_BINDED", "NAT"}, false),
			},
			"icmp_reply_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"private_link_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"access_mode": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"tunnel_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateNatGateway"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := d.GetOk("nat_gateway_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	request["NatType"] = d.Get("nat_type")
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertNatGatewayPaymentTypeRequest(v.(string))
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}

	if v, ok := request["InstanceChargeType"]; ok && v.(string) == "PrePaid" {
		period := d.Get("period").(int)
		request["Duration"] = strconv.Itoa(period)
		request["PricingCycle"] = "Month"
		if period > 9 {
			request["Duration"] = strconv.Itoa(period / 12)
			request["PricingCycle"] = string(Year)
		}
		request["AutoPay"] = true
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("specification"); ok {
		request["Spec"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("eip_bind_mode"); ok {
		request["EipBindMode"] = v
	}

	if v, ok := d.GetOkExists("icmp_reply_enabled"); ok {
		request["IcmpReplyEnabled"] = v
	}

	if v, ok := d.GetOkExists("private_link_enabled"); ok {
		request["PrivateLinkEnabled"] = v
	}

	if v, ok := d.GetOk("access_mode"); ok {
		accessModeMap := map[string]interface{}{}
		for _, accessModeList := range v.([]interface{}) {
			accessModeArg := accessModeList.(map[string]interface{})

			if modeValue, ok := accessModeArg["mode_value"]; ok && modeValue.(string) != "" {
				accessModeMap["ModeValue"] = modeValue
			}

			if tunnelType, ok := accessModeArg["tunnel_type"]; ok && tunnelType.(string) != "" {
				accessModeMap["TunnelType"] = tunnelType
			}
		}

		accessModeJson, err := convertMaptoJsonString(accessModeMap)
		if err != nil {
			return WrapError(err)
		}

		request["AccessMode"] = accessModeJson
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateNatGateway")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "VswitchStatusError", "IncorrectStatus.VSWITCH", "OperationConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nat_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["NatGatewayId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNatGatewayUpdate(d, meta)
}

func resourceAliCloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nat_gateway vpcService.DescribeNatGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", object["Description"])
	if v, ok := object["ForwardTableIds"].(map[string]interface{})["ForwardTableId"].([]interface{}); ok {
		ids := []string{}
		for _, id := range v {
			ids = append(ids, id.(string))
		}
		d.Set("forward_table_ids", strings.Join(ids, ","))
	}
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("nat_gateway_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("nat_type", object["NatType"])
	d.Set("payment_type", convertNatGatewayPaymentTypeResponse(object["InstanceChargeType"].(string)))
	d.Set("instance_charge_type", object["InstanceChargeType"])
	d.Set("network_type", object["NetworkType"])
	d.Set("eip_bind_mode", object["EipBindMode"])
	//if object["InstanceChargeType"] == "PrePaid" {
	//	period, err := computePeriodByUnit(object["CreationTime"], object["ExpiredTime"], d.Get("period").(int), "Month")
	//	if err != nil {
	//		return WrapError(err)
	//	}
	//	d.Set("period", period)
	//}
	if v, ok := object["SnatTableIds"].(map[string]interface{})["SnatTableId"].([]interface{}); ok {
		ids := []string{}
		for _, id := range v {
			ids = append(ids, id.(string))
		}
		d.Set("snat_table_ids", strings.Join(ids, ","))
	}
	d.Set("specification", object["Spec"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["NatGatewayPrivateInfo"].(map[string]interface{})["VswitchId"])
	d.Set("vpc_id", object["VpcId"])

	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "NATGATEWAY")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))
	d.Set("deletion_protection", object["DeletionProtection"])
	d.Set("icmp_reply_enabled", object["IcmpReplyEnabled"])
	d.Set("private_link_enabled", object["PrivateLinkEnabled"])

	if accessMode, ok := object["AccessMode"]; ok {
		accessModeMaps := make([]map[string]interface{}, 0)
		accessModeArg := accessMode.(map[string]interface{})
		accessModeMap := make(map[string]interface{})

		if modeValue, ok := accessModeArg["ModeValue"]; ok {
			accessModeMap["mode_value"] = modeValue
		}

		if tunnelType, ok := accessModeArg["TunnelType"]; ok {
			accessModeMap["tunnel_type"] = tunnelType
		}

		accessModeMaps = append(accessModeMaps, accessModeMap)

		d.Set("access_mode", accessModeMaps)
	}

	return nil
}

func resourceAliCloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "NATGATEWAY"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("deletion_protection") {
		var response map[string]interface{}
		action := "DeletionProtection"
		request := map[string]interface{}{
			"RegionId":         client.RegionId,
			"InstanceId":       d.Id(),
			"ProtectionEnable": d.Get("deletion_protection"),
			"Type":             "NATGW",
		}
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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

		d.SetPartial("deletion_protection")
	}

	update := false
	request := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("nat_gateway_name") {
		update = true
		request["Name"] = d.Get("nat_gateway_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if !d.IsNewResource() && d.HasChange("eip_bind_mode") {
		update = true
		request["EipBindMode"] = d.Get("eip_bind_mode")
	}

	if !d.IsNewResource() && d.HasChange("icmp_reply_enabled") {
		update = true

		if v, ok := d.GetOkExists("icmp_reply_enabled"); ok {
			request["IcmpReplyEnabled"] = v
		}
	}

	if update {
		action := "ModifyNatGatewayAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.NATGW"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("nat_gateway_name")
	}

	if !d.IsNewResource() && d.HasChange("specification") {
		request := map[string]interface{}{
			"NatGatewayId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["Spec"] = d.Get("specification")
		action := "ModifyNatGatewaySpec"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.NatGateway"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("specification")
	}
	update = false
	updateNatGatewayNatTypeReq := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("nat_type") {
		update = true
	}
	updateNatGatewayNatTypeReq["NatType"] = d.Get("nat_type")
	updateNatGatewayNatTypeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
	}
	updateNatGatewayNatTypeReq["VSwitchId"] = d.Get("vswitch_id")
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			updateNatGatewayNatTypeReq["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateNatGatewayNatType"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		updateNatGatewayNatTypeReq["ClientToken"] = buildClientToken("UpdateNatGatewayNatType")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, updateNatGatewayNatTypeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationFailed.NatGwRouteInMiddleStatus", "TaskConflict", "UnknownError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateNatGatewayNatTypeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("nat_type")
		d.SetPartial("vswitch_id")
		d.SetPartial("dry_run")
	}

	d.Partial(false)

	return resourceAliCloudNatGatewayRead(d, meta)
}

func resourceAliCloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" || d.Get("instance_charge_type").(string) == "Prepaid" {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_nat_gateway. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteNatGateway"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages", "DependencyViolation.EIPS", "OperationConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXISTS", "IncorrectStatus.NatGateway", "InvalidNatGatewayId.NotFound", "InvalidRegionId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertNatGatewayPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func convertNatGatewayPaymentTypeResponse(source string) string {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
