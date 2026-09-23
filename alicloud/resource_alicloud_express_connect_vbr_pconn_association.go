// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectVbrPconnAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectVbrPconnAssociationCreate,
		Read:   resourceAliCloudExpressConnectVbrPconnAssociationRead,
		Delete: resourceAliCloudExpressConnectVbrPconnAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"circuit_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"local_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"local_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peer_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peer_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peering_ipv6_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peering_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"physical_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectVbrPconnAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociatePhysicalConnectionToVirtualBorderRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("physical_connection_id"); ok {
		request["PhysicalConnectionId"] = v
	}
	if v, ok := d.GetOk("vbr_id"); ok {
		request["VbrId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["VlanId"] = d.Get("vlan_id")
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
		request["PeeringIpv6SubnetMask"] = v
	}
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
		request["LocalIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOk("local_gateway_ip"); ok {
		request["LocalGatewayIp"] = v
	}
	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
		request["PeerIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOk("peering_subnet_mask"); ok {
		request["PeeringSubnetMask"] = v
	}
	if v, ok := d.GetOk("peer_gateway_ip"); ok {
		request["PeerGatewayIp"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_vbr_pconn_association", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VbrId"], request["PhysicalConnectionId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectServiceV2.ExpressConnectVbrPconnAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectVbrPconnAssociationRead(d, meta)
}

func resourceAliCloudExpressConnectVbrPconnAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectVbrPconnAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_vbr_pconn_association DescribeExpressConnectVbrPconnAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("circuit_code", objectRaw["CircuitCode"])
	d.Set("enable_ipv6", objectRaw["EnableIpv6"])
	d.Set("local_gateway_ip", objectRaw["LocalGatewayIp"])
	d.Set("local_ipv6_gateway_ip", objectRaw["LocalIpv6GatewayIp"])
	d.Set("peer_gateway_ip", objectRaw["PeerGatewayIp"])
	d.Set("peer_ipv6_gateway_ip", objectRaw["PeerIpv6GatewayIp"])
	d.Set("peering_ipv6_subnet_mask", objectRaw["PeeringIpv6SubnetMask"])
	d.Set("peering_subnet_mask", objectRaw["PeeringSubnetMask"])
	d.Set("status", objectRaw["Status"])
	d.Set("vlan_id", formatInt(objectRaw["VlanId"]))
	d.Set("physical_connection_id", objectRaw["PhysicalConnectionId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("vbr_id", parts[0])

	return nil
}

func resourceAliCloudExpressConnectVbrPconnAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnassociatePhysicalConnectionFromVirtualBorderRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PhysicalConnectionId"] = parts[1]
	request["VbrId"] = parts[0]
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

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectServiceV2.ExpressConnectVbrPconnAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
