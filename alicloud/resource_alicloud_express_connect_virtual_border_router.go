// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudExpressConnectVirtualBorderRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectVirtualBorderRouterCreate,
		Read:   resourceAliCloudExpressConnectVirtualBorderRouterRead,
		Update: resourceAliCloudExpressConnectVirtualBorderRouterUpdate,
		Delete: resourceAliCloudExpressConnectVirtualBorderRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detect_multiplier": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"local_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"min_rx_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_tx_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{0, 1500, 8500}),
			},
			"peer_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_ipv6_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_subnet_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sitelink_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vbr_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_border_router_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"associated_physical_connections": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `associated_physical_connections` has been deprecated from provider version 1.263.0. Please use the resource `alicloud_express_connect_vbr_pconn_association` instead.",
			},
			"include_cross_account_vbr": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field `include_cross_account_vbr` has been removed from provider version 1.263.0.",
			},
		},
	}
}

func resourceAliCloudExpressConnectVirtualBorderRouterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVirtualBorderRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("vbr_owner_id"); ok {
		request["VbrOwnerId"] = v
	}
	request["VlanId"] = d.Get("vlan_id")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("circuit_code"); ok {
		request["CircuitCode"] = v
	}
	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
		request["PeeringIpv6SubnetMask"] = v
	}
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
		request["LocalIpv6GatewayIp"] = v
	}
	request["LocalGatewayIp"] = d.Get("local_gateway_ip")
	if v, ok := d.GetOk("virtual_border_router_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
		request["PeerIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}
	request["PeeringSubnetMask"] = d.Get("peering_subnet_mask")
	request["PhysicalConnectionId"] = d.Get("physical_connection_id")
	request["PeerGatewayIp"] = d.Get("peer_gateway_ip")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_virtual_border_router", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VbrId"]))

	return resourceAliCloudExpressConnectVirtualBorderRouterUpdate(d, meta)
}

func resourceAliCloudExpressConnectVirtualBorderRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectVirtualBorderRouter(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_virtual_border_router DescribeExpressConnectVirtualBorderRouter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth", formatInt(objectRaw["Bandwidth"]))
	d.Set("circuit_code", objectRaw["CircuitCode"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	if v, ok := objectRaw["DetectMultiplier"]; ok && fmt.Sprint(v) != "0" {
		d.Set("detect_multiplier", formatInt(v))
	}
	d.Set("enable_ipv6", objectRaw["EnableIpv6"])
	d.Set("local_gateway_ip", objectRaw["LocalGatewayIp"])
	d.Set("local_ipv6_gateway_ip", objectRaw["LocalIpv6GatewayIp"])
	if v, ok := objectRaw["MinRxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_rx_interval", formatInt(v))
	}
	if v, ok := objectRaw["MinTxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_tx_interval", formatInt(v))
	}
	d.Set("mtu", objectRaw["Mtu"])
	d.Set("peer_gateway_ip", objectRaw["PeerGatewayIp"])
	d.Set("peer_ipv6_gateway_ip", objectRaw["PeerIpv6GatewayIp"])
	d.Set("peering_ipv6_subnet_mask", objectRaw["PeeringIpv6SubnetMask"])
	d.Set("peering_subnet_mask", objectRaw["PeeringSubnetMask"])
	d.Set("physical_connection_id", objectRaw["PhysicalConnectionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("route_table_id", objectRaw["RouteTableId"])
	d.Set("sitelink_enable", objectRaw["SitelinkEnable"])
	d.Set("status", objectRaw["Status"])
	d.Set("virtual_border_router_name", objectRaw["Name"])
	d.Set("vlan_id", formatInt(objectRaw["VlanId"]))

	tagsMaps, _ := jsonpath.Get("$.Tags.Tags", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudExpressConnectVirtualBorderRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	objectRaw, _ := expressConnectServiceV2.DescribeExpressConnectVirtualBorderRouter(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "terminated" {
				action := "TerminateVirtualBorderRouter"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VbrId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
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
				stateConf := BuildStateConf([]string{}, []string{"terminated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "active" {
				action := "RecoverVirtualBorderRouter"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VbrId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
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
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ModifyVirtualBorderRouterAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VbrId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("vlan_id") {
		update = true
	}
	request["VlanId"] = d.Get("vlan_id")
	if !d.IsNewResource() && d.HasChange("circuit_code") {
		update = true
	}
	if v, ok := d.GetOk("circuit_code"); ok || d.HasChange("circuit_code") {
		request["CircuitCode"] = v
	}
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true

		if v, ok := d.GetOkExists("bandwidth"); ok || d.HasChange("bandwidth") {
			request["Bandwidth"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if !d.IsNewResource() && d.HasChange("peering_ipv6_subnet_mask") {
		update = true
	}
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok || d.HasChange("peering_ipv6_subnet_mask") {
		request["PeeringIpv6SubnetMask"] = v
	}
	if !d.IsNewResource() && d.HasChange("local_ipv6_gateway_ip") {
		update = true
	}
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok || d.HasChange("local_ipv6_gateway_ip") {
		request["LocalIpv6GatewayIp"] = v
	}
	if d.HasChange("min_rx_interval") {
		update = true
	}
	if v, ok := d.GetOkExists("min_rx_interval"); ok || d.HasChange("min_rx_interval") {
		request["MinRxInterval"] = v
	}
	if !d.IsNewResource() && d.HasChange("local_gateway_ip") {
		update = true
	}
	request["LocalGatewayIp"] = d.Get("local_gateway_ip")
	if !d.IsNewResource() && d.HasChange("virtual_border_router_name") {
		update = true
	}
	if v, ok := d.GetOk("virtual_border_router_name"); ok || d.HasChange("virtual_border_router_name") {
		request["Name"] = v
	}
	if !d.IsNewResource() && d.HasChange("enable_ipv6") {
		update = true

		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request["EnableIpv6"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("peer_ipv6_gateway_ip") {
		update = true
	}
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok || d.HasChange("peer_ipv6_gateway_ip") {
		request["PeerIpv6GatewayIp"] = v
	}
	if !d.IsNewResource() && d.HasChange("peering_subnet_mask") {
		update = true
	}
	request["PeeringSubnetMask"] = d.Get("peering_subnet_mask")
	if d.HasChange("mtu") {
		update = true
	}
	if v, ok := d.GetOkExists("mtu"); ok || d.HasChange("mtu") {
		request["Mtu"] = v
	}
	if d.HasChange("min_tx_interval") {
		update = true
	}
	if v, ok := d.GetOkExists("min_tx_interval"); ok || d.HasChange("min_tx_interval") {
		request["MinTxInterval"] = v
	}
	if d.HasChange("detect_multiplier") {
		update = true
	}
	if v, ok := d.GetOkExists("detect_multiplier"); ok || d.HasChange("detect_multiplier") {
		request["DetectMultiplier"] = v
	}

	if d.HasChange("sitelink_enable") {
		update = true

		if v, ok := d.GetOkExists("sitelink_enable"); ok || d.HasChange("sitelink_enable") {
			request["SitelinkEnable"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("peer_gateway_ip") {
		update = true
	}
	request["PeerGatewayIp"] = d.Get("peer_gateway_ip")
	if update {
		if v, ok := d.GetOk("associated_physical_connections"); ok {
			request["AssociatedPhysicalConnections"] = v
		}
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
	request["ResourceType"] = "VIRTUALBORDERROUTER"
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
		if err := expressConnectServiceV2.SetResourceTags(d, "VIRTUALBORDERROUTER"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudExpressConnectVirtualBorderRouterRead(d, meta)
}

func resourceAliCloudExpressConnectVirtualBorderRouterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVirtualBorderRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VbrId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.BgpGroup"}) || NeedRetry(err) {
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

	return nil
}
