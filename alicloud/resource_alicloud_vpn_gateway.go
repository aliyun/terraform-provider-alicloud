// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVPNGatewayVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVPNGatewayVPNGatewayCreate,
		Read:   resourceAliCloudVPNGatewayVPNGatewayRead,
		Update: resourceAliCloudVPNGatewayVPNGatewayUpdate,
		Delete: resourceAliCloudVPNGatewayVPNGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_propagate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disaster_recovery_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ipsec": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"public", "private"}, true),
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_charge_type"},
				ForceNew:      true,
				ValidateFunc:  StringInSlice([]string{"Subscription", "PayAsYouGo", "PostPaid", "PrePaid"}, true),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.Any(IntBetween(1, 9), IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_gateway_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"vpn_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "NationalStandard"}, true),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'instance_charge_type' has been deprecated since provider version 1.215.0. New field 'payment_type' instead.",
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo", "PostPaid", "PrePaid"}, true),
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_vpn_internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disaster_recovery_internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.215.0. New field 'vpn_gateway_name' instead.",
			},
		},
	}
}

func resourceAliCloudVPNGatewayVPNGatewayCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpnGateway"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("vpn_gateway_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	} else {
		request["AutoPay"] = true
	}
	if v, ok := d.GetOkExists("enable_ipsec"); ok {
		request["EnableIpsec"] = v
	}
	if v, ok := d.GetOkExists("enable_ssl"); ok {
		request["EnableSsl"] = v
	}
	if v, ok := d.GetOk("ssl_connections"); ok && d.Get("enable_ssl").(bool) {
		request["SslConnections"] = v
	}
	if v, ok := d.GetOk("vpn_type"); ok {
		request["VpnType"] = v
	}
	if v, ok := d.GetOk("disaster_recovery_vswitch_id"); ok {
		request["DisasterRecoveryVSwitchId"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = convertVPNGatewayInstanceChargeTypeRequest(convertChargeTypeToPaymentType(v.(string)))
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertVPNGatewayInstanceChargeTypeRequest(v.(string))
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("period"); ok && v.(int) != 0 && request["InstanceChargeType"] == "PREPAY" {
		request["Period"] = requests.NewInteger(v.(int))
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.SslNotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpnGatewayId"]))

	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, vPNGatewayServiceV2.VPNGatewayVPNGatewayStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVPNGatewayVPNGatewayUpdate(d, meta)
}

func resourceAliCloudVPNGatewayVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}

	objectRaw, err := vPNGatewayServiceV2.DescribeVPNGatewayVPNGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway DescribeVPNGatewayVPNGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_propagate", objectRaw["AutoPropagate"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("disaster_recovery_vswitch_id", objectRaw["DisasterRecoveryVSwitchId"])
	d.Set("network_type", objectRaw["NetworkType"])
	d.Set("payment_type", convertVPNGatewayChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("ssl_connections", objectRaw["SslMaxConnections"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vpn_gateway_name", objectRaw["Name"])
	d.Set("vpn_type", objectRaw["VpnType"])
	d.Set("internet_ip", objectRaw["InternetIp"])
	d.Set("business_status", objectRaw["BusinessStatus"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	e := jsonata.MustCompile("$substringBefore($.Spec, \"M\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("bandwidth", formatInt(evaluation))

	d.Set("instance_charge_type", convertPaymentTypeToChargeType(d.Get("payment_type")))
	d.Set("name", d.Get("vpn_gateway_name"))

	d.Set("enable_ipsec", "enable" == objectRaw["IpsecVpn"])
	d.Set("enable_ssl", "enable" == objectRaw["SslVpn"])
	d.Set("ssl_connections", objectRaw["SslMaxConnections"])
	d.Set("ssl_vpn_internet_ip", objectRaw["SslVpnInternetIp"])
	d.Set("disaster_recovery_internet_ip", objectRaw["DisasterRecoveryInternetIp"])
	return nil
}

func resourceAliCloudVPNGatewayVPNGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyVpnGatewayAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["VpnGatewayId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("auto_propagate") {
		update = true
		request["AutoPropagate"] = d.Get("auto_propagate")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("vpn_gateway_name") {
		update = true
		request["Name"] = d.Get("vpn_gateway_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vPNGatewayServiceV2.VPNGatewayVPNGatewayStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("auto_propagate")
		d.SetPartial("vpn_gateway_name")
	}
	update = false
	action = "MoveVpnResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "VpnGateWay"
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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		if err := vPNGatewayServiceV2.SetResourceTags(d, "VpnGateWay"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVPNGatewayVPNGatewayRead(d, meta)
}

func resourceAliCloudVPNGatewayVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_vpn_gateway which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpnGateway"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["VpnGatewayId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXISTS", "IncorrectStatus.VpnGateway", "InvalidVpnGatewayId.NotFound", "InvalidRegionId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vPNGatewayServiceV2.VPNGatewayVPNGatewayStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertVPNGatewayChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "Prepay":
		return "Subscription"
	case "PostpayByFlow":
		return "PayAsYouGo"
	}
	return source
}
func convertVPNGatewayInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PREPAY"
	case "PrePaid":
		return "PREPAY"
	case "PayAsYouGo":
		return "POSTPAY"
	case "PostPaid":
		return "POSTPAY"
	}
	return source
}
