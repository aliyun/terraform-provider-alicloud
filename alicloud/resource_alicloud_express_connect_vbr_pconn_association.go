package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudExpressConnectVbrPconnAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectVbrPconnAssociationCreate,
		Read:   resourceAlicloudExpressConnectVbrPconnAssociationRead,
		Delete: resourceAlicloudExpressConnectVbrPconnAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"circuit_code": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"enable_ipv6": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"local_gateway_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"local_ipv6_gateway_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peer_gateway_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peer_ipv6_gateway_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peering_ipv6_subnet_mask": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peering_subnet_mask": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"physical_connection_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vbr_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vlan_id": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 2999),
			},
		},
	}
}

func resourceAlicloudExpressConnectVbrPconnAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectService := ExpressConnectService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error
	if v, ok := d.GetOk("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}
	if v, ok := d.GetOk("local_gateway_ip"); ok {
		request["LocalGatewayIp"] = v
	}
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
		request["LocalIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOk("peer_gateway_ip"); ok {
		request["PeerGatewayIp"] = v
	}
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
		request["PeerIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
		request["PeeringIpv6SubnetMask"] = v
	}
	if v, ok := d.GetOk("peering_subnet_mask"); ok {
		request["PeeringSubnetMask"] = v
	}
	if v, ok := d.GetOk("physical_connection_id"); ok {
		request["PhysicalConnectionId"] = v
	}
	if v, ok := d.GetOk("vbr_id"); ok {
		request["VbrId"] = v
	}
	if v, ok := d.GetOk("vlan_id"); ok {
		request["VlanId"] = v
	}

	var response map[string]interface{}
	action := "AssociatePhysicalConnectionToVirtualBorderRouter"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("AssociatePhysicalConnectionToVirtualBorderRouter")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_vbr_pconn_association", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VbrId"], ":", request["PhysicalConnectionId"]))
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectService.ExpressConnectVbrPconnAssociationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudExpressConnectVbrPconnAssociationRead(d, meta)
}

func resourceAlicloudExpressConnectVbrPconnAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectService := ExpressConnectService{client}

	object, err := expressConnectService.DescribeExpressConnectVbrPconnAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_vbr_pconn_association expressConnectService.DescribeExpressConnectVbrPconnAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vbr_id", parts[0])
	d.Set("physical_connection_id", parts[1])
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("enable_ipv6", object["EnableIpv6"])
	d.Set("local_gateway_ip", object["LocalGatewayIp"])
	d.Set("local_ipv6_gateway_ip", object["LocalIpv6GatewayIp"])
	d.Set("peer_gateway_ip", object["PeerGatewayIp"])
	d.Set("peer_ipv6_gateway_ip", object["PeerIpv6GatewayIp"])
	d.Set("peering_ipv6_subnet_mask", object["PeeringIpv6SubnetMask"])
	d.Set("peering_subnet_mask", object["PeeringSubnetMask"])
	d.Set("status", object["Status"])
	d.Set("vlan_id", formatInt(object["VlanId"]))

	return nil
}

func resourceAlicloudExpressConnectVbrPconnAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectService := ExpressConnectService{client}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VbrId":                parts[0],
		"PhysicalConnectionId": parts[1],
		"RegionId":             client.RegionId,
	}

	action := "UnassociatePhysicalConnectionFromVirtualBorderRouter"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("UnassociatePhysicalConnectionFromVirtualBorderRouter")
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectService.ExpressConnectVbrPconnAssociationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
