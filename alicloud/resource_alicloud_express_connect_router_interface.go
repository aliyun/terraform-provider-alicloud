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

func resourceAliCloudExpressConnectRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectRouterInterfaceCreate,
		Read:   resourceAliCloudExpressConnectRouterInterfaceRead,
		Update: resourceAliCloudExpressConnectRouterInterfaceUpdate,
		Delete: resourceAliCloudExpressConnectRouterInterfaceDelete,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connected_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_border": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"delete_health_check_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fast_link_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"has_reservation_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hc_rate": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hc_threshold": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_target_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_access_point_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"opposite_interface_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opposite_interface_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_interface_spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opposite_interface_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opposite_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_router_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VBR", "VRouter"}, false),
			},
			"opposite_vpc_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reservation_active_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation_bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation_internet_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation_order_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"InitiatingSide", "AcceptingSide"}, false),
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_interface_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VBR", "VRouter"}, false),
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Mini.2", "Mini.5", "Small.1", "Small.2", "Small.5", "Middle.1", "Middle.2", "Middle.5", "Large.1", "Large.2", "Large.5", "Xlarge.1", "Negative"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Idle", "Connecting", "AcceptingConnecting", "Activating", "Active", "Modifying", "Deactivating", "Inactive", "Deleting", "Deleted"}, false),
			},
			"tags": tagsSchema(),
			"vpc_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_pay": {
				Optional:   true,
				Type:       schema.TypeBool,
				Deprecated: "Field 'auto_pay' has been deprecated since provider version 1.263.0.",
			},
			"opposite_interface_id": {
				Optional:   true,
				Computed:   true,
				Type:       schema.TypeString,
				Deprecated: "Field 'opposite_interface_id' has been deprecated since provider version 1.263.0.",
			},
		},
	}
}

func resourceAliCloudExpressConnectRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRouterInterface"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertExpressConnectRouterInterfaceInstanceChargeTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("health_check_source_ip"); ok {
		request["HealthCheckSourceIp"] = v
	}
	request["AutoPay"] = "true"
	if v, ok := d.GetOk("opposite_router_id"); ok {
		request["OppositeRouterId"] = v
	}
	request["RouterType"] = d.Get("router_type")
	if v, ok := d.GetOk("access_point_id"); ok {
		request["AccessPointId"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("opposite_access_point_id"); ok {
		request["OppositeAccessPointId"] = v
	}
	if v, ok := d.GetOkExists("fast_link_mode"); ok {
		request["FastLinkMode"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	request["OppositeRegionId"] = d.Get("opposite_region_id")
	if v, ok := d.GetOk("opposite_interface_owner_id"); ok {
		request["OppositeInterfaceOwnerId"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["RouterId"] = d.Get("router_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("router_interface_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("health_check_target_ip"); ok {
		request["HealthCheckTargetIp"] = v
	}
	if v, ok := d.GetOk("opposite_router_type"); ok {
		request["OppositeRouterType"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	request["Spec"] = d.Get("spec")
	request["Role"] = d.Get("role")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_interface", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RouterInterfaceId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Idle", "Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectRouterInterfaceUpdate(d, meta)
}

func resourceAliCloudExpressConnectRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectRouterInterface(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_interface DescribeExpressConnectRouterInterface Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_point_id", objectRaw["AccessPointId"])
	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("description", objectRaw["Description"])
	d.Set("fast_link_mode", formatBool(objectRaw["FastLinkMode"]))
	d.Set("hc_rate", objectRaw["HcRate"])
	if v, ok := objectRaw["HcThreshold"]; ok {
		d.Set("hc_threshold", v)
	}

	d.Set("health_check_source_ip", objectRaw["HealthCheckSourceIp"])
	d.Set("health_check_target_ip", objectRaw["HealthCheckTargetIp"])
	d.Set("opposite_access_point_id", objectRaw["OppositeAccessPointId"])
	d.Set("opposite_bandwidth", objectRaw["OppositeBandwidth"])
	d.Set("opposite_interface_business_status", objectRaw["OppositeInterfaceBusinessStatus"])
	if v, ok := objectRaw["OppositeInterfaceOwnerId"]; ok {
		d.Set("opposite_interface_owner_id", v)
	}

	d.Set("opposite_interface_spec", objectRaw["OppositeInterfaceSpec"])
	d.Set("opposite_interface_status", objectRaw["OppositeInterfaceStatus"])
	d.Set("opposite_region_id", objectRaw["OppositeRegionId"])
	d.Set("opposite_router_id", objectRaw["OppositeRouterId"])
	d.Set("opposite_router_type", objectRaw["OppositeRouterType"])
	d.Set("opposite_vpc_instance_id", objectRaw["OppositeVpcInstanceId"])
	d.Set("payment_type", convertExpressConnectRouterInterfaceChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("role", objectRaw["Role"])
	d.Set("router_id", objectRaw["RouterId"])
	d.Set("router_interface_name", objectRaw["Name"])
	d.Set("router_type", objectRaw["RouterType"])
	d.Set("spec", objectRaw["Spec"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_instance_id", objectRaw["VpcInstanceId"])
	d.Set("business_status", objectRaw["BusinessStatus"])
	d.Set("connected_time", objectRaw["ConnectedTime"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("cross_border", objectRaw["CrossBorder"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("has_reservation_data", objectRaw["HasReservationData"])
	d.Set("reservation_active_time", objectRaw["ReservationActiveTime"])
	d.Set("reservation_bandwidth", objectRaw["ReservationBandwidth"])
	d.Set("reservation_internet_charge_type", objectRaw["ReservationInternetChargeType"])
	d.Set("reservation_order_type", objectRaw["ReservationOrderType"])
	d.Set("router_interface_id", objectRaw["RouterInterfaceId"])
	d.Set("opposite_interface_id", objectRaw["OppositeInterfaceId"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tags", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudExpressConnectRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	objectRaw, _ := expressConnectServiceV2.DescribeExpressConnectRouterInterface(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)
		if objectRaw["Status"].(string) != target {
			enableActivateRouterInterfaceActive := false
			checkValue00 := objectRaw["Status"]
			if checkValue00 == "Inactive" {
				enableActivateRouterInterfaceActive = true
			}
			if enableActivateRouterInterfaceActive && target == "Active" {
				action := "ActivateRouterInterface"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RouterInterfaceId"] = d.Id()
				request["RegionId"] = client.RegionId
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
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				expressConnectServiceV2 := ExpressConnectServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			enableDeactivateRouterInterfaceInactive := false
			checkValue00 = objectRaw["Status"]
			if checkValue00 == "Active" {
				enableDeactivateRouterInterfaceInactive = true
			}
			if enableDeactivateRouterInterfaceInactive && target == "Inactive" {
				action := "DeactivateRouterInterface"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RouterInterfaceId"] = d.Id()
				request["RegionId"] = client.RegionId
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
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				expressConnectServiceV2 := ExpressConnectServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Inactive"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ModifyRouterInterfaceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RouterInterfaceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("hc_threshold") {
		update = true
		request["HcThreshold"] = d.Get("hc_threshold")
	}

	if !d.IsNewResource() && d.HasChange("health_check_source_ip") {
		update = true
		request["HealthCheckSourceIp"] = d.Get("health_check_source_ip")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("router_interface_name") {
		update = true
		request["Name"] = d.Get("router_interface_name")
	}

	if d.HasChange("hc_rate") {
		update = true
		request["HcRate"] = d.Get("hc_rate")
	}

	if !d.IsNewResource() && d.HasChange("health_check_target_ip") {
		update = true
		request["HealthCheckTargetIp"] = d.Get("health_check_target_ip")
	}

	if v, ok := d.GetOk("delete_health_check_ip"); ok {
		request["DeleteHealthCheckIp"] = v
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Idle", "Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyRouterInterfaceSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RouterInterfaceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
	}
	request["Spec"] = d.Get("spec")
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "ROUTERINTERFACE"
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		if err := expressConnectServiceV2.SetResourceTags(d, "ROUTERINTERFACE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudExpressConnectRouterInterfaceRead(d, meta)
}

func resourceAliCloudExpressConnectRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := false
	if v, ok := d.GetOkExists("fast_link_mode"); ok {
		if InArray(fmt.Sprint(v), []string{"true"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := "DeleteExpressConnect"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RouterInterfaceId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		request["Force"] = "true"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}

	enableDelete = false
	if v, ok := d.GetOkExists("fast_link_mode"); ok {
		if InArray(fmt.Sprint(v), []string{"false"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		enableDelete = false
		if v, ok := d.GetOkExists("payment_type"); ok {
			if InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
				enableDelete = true
			}
		}
		if enableDelete {
			client := meta.(*connectivity.AliyunClient)
			action := "DeleteRouterInterface"
			var request map[string]interface{}
			var response map[string]interface{}
			query := make(map[string]interface{})
			var err error
			request = make(map[string]interface{})
			request["RouterInterfaceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				if NotFoundError(err) {
					return nil
				}
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	enableDelete = false
	if v, ok := d.GetOkExists("fast_link_mode"); ok {
		if InArray(fmt.Sprint(v), []string{"false"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		enableDelete = false
		if v, ok := d.GetOkExists("payment_type"); ok {
			if InArray(fmt.Sprint(v), []string{"Subscription"}) {
				enableDelete = true
			}
		}
		if enableDelete {
			client := meta.(*connectivity.AliyunClient)
			action := "RefundInstance"
			var request map[string]interface{}
			var response map[string]interface{}
			query := make(map[string]interface{})
			var err error
			request = make(map[string]interface{})
			request["InstanceId"] = d.Id()

			request["ClientToken"] = buildClientToken(action)

			request["ImmediatelyRelease"] = "1"
			var endpoint string
			request["ProductCode"] = "ri"
			request["ProductType"] = ""
			if client.IsInternationalAccount() {
				request["ProductType"] = ""
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
						request["ProductCode"] = "ri_pre_intl"
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

			expressConnectServiceV2 := ExpressConnectServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectServiceV2.ExpressConnectRouterInterfaceStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}
	return nil
}

func convertExpressConnectRouterInterfaceChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Prepaid":
		return "Subscription"
	case "AfterPay":
		return "PayAsYouGo"
	}
	return source
}
func convertExpressConnectRouterInterfaceInstanceChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
