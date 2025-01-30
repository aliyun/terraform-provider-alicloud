package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudExpressConnectRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectRouterInterfaceCreate,
		Read:   resourceAlicloudExpressConnectRouterInterfaceRead,
		Update: resourceAlicloudExpressConnectRouterInterfaceUpdate,
		Delete: resourceAlicloudExpressConnectRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"auto_pay": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"bandwidth": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"business_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"connected_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"cross_border": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"delete_health_check_ip": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"end_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"has_reservation_data": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"hc_rate": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"hc_threshold": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"health_check_source_ip": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"health_check_target_ip": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"opposite_access_point_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"opposite_bandwidth": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"opposite_interface_business_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"opposite_interface_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"opposite_interface_owner_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"opposite_interface_spec": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"opposite_interface_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"opposite_region_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"opposite_router_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"opposite_router_type": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"VRouter", "VBR"}, false),
			},
			"opposite_vpc_instance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"pricing_cycle": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"reservation_active_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"reservation_bandwidth": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"reservation_internet_charge_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"reservation_order_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"role": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"router_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"router_interface_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"router_interface_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"router_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"VRouter", "VBR"}, false),
			},
			"spec": {
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Large.1", "Large.2", "Large.5", "Middle.1", "Middle.2", "Middle.5", "Mini.2", "Mini.5", "Small.1", "Small.2", "Small.5", "Negative"}, false),
				Type:         schema.TypeString,
			},
			"status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Idle", "AcceptingConnecting", "Connecting", "Activating", "Active", "Modifying", "Deactivating", "Inactive", "Deleting"}, false),
			},
			"vpc_instance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudExpressConnectRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("access_point_id"); ok {
		request["AccessPointId"] = v
	}
	if v, ok := d.GetOk("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("health_check_source_ip"); ok {
		request["HealthCheckSourceIp"] = v
	}
	if v, ok := d.GetOk("health_check_target_ip"); ok {
		request["HealthCheckTargetIp"] = v
	}
	if v, ok := d.GetOk("opposite_access_point_id"); ok {
		request["OppositeAccessPointId"] = v
	}
	if v, ok := d.GetOk("opposite_interface_id"); ok {
		request["OppositeInterfaceId"] = v
	}
	if v, ok := d.GetOk("opposite_interface_owner_id"); ok {
		request["OppositeInterfaceOwnerId"] = v
	}
	if v, ok := d.GetOk("opposite_region_id"); ok {
		request["OppositeRegionId"] = v
	}
	if v, ok := d.GetOk("opposite_router_id"); ok {
		request["OppositeRouterId"] = v
	}
	if v, ok := d.GetOk("opposite_router_type"); ok {
		request["OppositeRouterType"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertRouterInterfaceRequest(v.(string))
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("role"); ok {
		request["Role"] = v
	}
	if v, ok := d.GetOk("router_id"); ok {
		request["RouterId"] = v
	}
	if v, ok := d.GetOk("router_interface_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("router_type"); ok {
		request["RouterType"] = v
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}

	var response map[string]interface{}
	action := "CreateRouterInterface"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateRouterInterface")
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_interface", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.RouterInterfaceId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_express_connect_router_interface")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudExpressConnectRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudExpressConnectRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeExpressConnectRouterInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_interface vpcService.DescribeExpressConnectRouterInterface Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("access_point_id", object["AccessPointId"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("business_status", object["BusinessStatus"])
	d.Set("connected_time", object["ConnectedTime"])
	d.Set("create_time", object["CreationTime"])
	d.Set("cross_border", object["CrossBorder"])
	d.Set("description", object["Description"])
	d.Set("end_time", object["EndTime"])
	d.Set("has_reservation_data", object["HasReservationData"])
	d.Set("hc_rate", object["HcRate"])
	d.Set("hc_threshold", object["HcThreshold"])
	d.Set("health_check_source_ip", object["HealthCheckSourceIp"])
	d.Set("health_check_target_ip", object["HealthCheckTargetIp"])
	d.Set("opposite_access_point_id", object["OppositeAccessPointId"])
	d.Set("opposite_bandwidth", object["OppositeBandwidth"])
	d.Set("opposite_interface_business_status", object["OppositeInterfaceBusinessStatus"])
	d.Set("opposite_interface_id", object["OppositeInterfaceId"])
	d.Set("opposite_interface_owner_id", object["OppositeInterfaceOwnerId"])
	d.Set("opposite_interface_spec", object["OppositeInterfaceSpec"])
	d.Set("opposite_interface_status", object["OppositeInterfaceStatus"])
	d.Set("opposite_region_id", object["OppositeRegionId"])
	d.Set("opposite_router_id", object["OppositeRouterId"])
	d.Set("opposite_router_type", object["OppositeRouterType"])
	d.Set("opposite_vpc_instance_id", object["OppositeVpcInstanceId"])
	d.Set("payment_type", convertRouterInterfaceResponse(object["ChargeType"]))
	d.Set("reservation_active_time", object["ReservationActiveTime"])
	d.Set("reservation_bandwidth", object["ReservationBandwidth"])
	d.Set("reservation_internet_charge_type", object["ReservationInternetChargeType"])
	d.Set("reservation_order_type", object["ReservationOrderType"])
	d.Set("role", object["Role"])
	d.Set("router_id", object["RouterId"])
	d.Set("router_interface_name", object["Name"])
	d.Set("router_type", object["RouterType"])
	d.Set("spec", object["Spec"])
	d.Set("status", object["Status"])
	d.Set("vpc_instance_id", object["VpcInstanceId"])

	return nil
}

func resourceAlicloudExpressConnectRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	vpcService := VpcService{client}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"RouterInterfaceId": d.Id(),
		"RegionId":          client.RegionId,
	}

	if v, ok := d.GetOk("delete_health_check_ip"); ok {
		request["DeleteHealthCheckIp"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("hc_rate") {
		update = true
		if v, ok := d.GetOk("hc_rate"); ok {
			request["HcRate"] = v
		}
	}
	if d.HasChange("hc_threshold") {
		update = true
		if v, ok := d.GetOk("hc_threshold"); ok {
			request["HcThreshold"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("health_check_source_ip") {
		update = true
		if v, ok := d.GetOk("health_check_source_ip"); ok {
			request["HealthCheckSourceIp"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("health_check_target_ip") {
		update = true
		if v, ok := d.GetOk("health_check_target_ip"); ok {
			request["HealthCheckTargetIp"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("opposite_interface_id") {
		update = true
		if v, ok := d.GetOk("opposite_interface_id"); ok {
			request["OppositeInterfaceId"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("opposite_interface_owner_id") {
		update = true
		if v, ok := d.GetOk("opposite_interface_owner_id"); ok {
			request["OppositeInterfaceOwnerId"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("opposite_router_id") {
		update = true
		if v, ok := d.GetOk("opposite_router_id"); ok {
			request["OppositeRouterId"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("opposite_router_type") {
		update = true
		if v, ok := d.GetOk("opposite_router_type"); ok {
			request["OppositeRouterType"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("router_interface_name") {
		update = true
		if v, ok := d.GetOk("router_interface_name"); ok {
			request["Name"] = v
		}
	}

	if update {
		action := "ModifyRouterInterfaceAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("delete_health_check_ip")
		d.SetPartial("description")
		d.SetPartial("hc_rate")
		d.SetPartial("hc_threshold")
		d.SetPartial("health_check_source_ip")
		d.SetPartial("health_check_target_ip")
		d.SetPartial("opposite_interface_id")
		d.SetPartial("opposite_interface_owner_id")
		d.SetPartial("opposite_router_id")
		d.SetPartial("opposite_router_type")
		d.SetPartial("router_interface_name")
	}

	update = false
	request = map[string]interface{}{
		"RouterInterfaceId": d.Id(),
		"RegionId":          client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
	}
	request["Spec"] = d.Get("spec")

	if update {
		action := "ModifyRouterInterfaceSpec"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyRouterInterfaceSpec")
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("spec")
	}

	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectRouterInterface(d.Id())
		if err != nil {
			WrapError(err)
		}
		target := fmt.Sprint(d.Get("status"))
		if fmt.Sprint(object["Status"]) != target {
			if target == "Activating" {
				request := map[string]interface{}{
					"RouterInterfaceId": d.Id(),
					"RegionId":          client.RegionId,
				}
				action := "ActivateRouterInterface"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "Active" {
				request := map[string]interface{}{
					"RouterInterfaceId": d.Id(),
					"RegionId":          client.RegionId,
				}

				action := "ConnectRouterInterface"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "Deactivating" {
				request := map[string]interface{}{
					"RouterInterfaceId": d.Id(),
					"RegionId":          client.RegionId,
				}

				action := "DeactivateRouterInterface"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudExpressConnectRouterInterfaceRead(d, meta)
}

func resourceAlicloudExpressConnectRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	vpcService := VpcService{client}

	if object, err := vpcService.DescribeExpressConnectRouterInterface(d.Id()); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	} else if object["Status"] == string(Active) {
		request := map[string]interface{}{
			"RouterInterfaceId": d.Id(),
			"RegionId":          client.RegionId,
		}

		action := "DeactivateRouterInterface"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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
			if IsExpectedErrors(err, []string{"InvalidRouterInterfaceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	request := map[string]interface{}{
		"RouterInterfaceId": d.Id(),
		"RegionId":          client.RegionId,
	}

	action := "DeleteRouterInterface"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteRouterInterface")
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertRouterInterfaceRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func convertRouterInterfaceResponse(source interface{}) interface{} {
	switch source {
	case "AfterPay":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
