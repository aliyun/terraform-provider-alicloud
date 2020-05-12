package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudVirtualBorderRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVirtualBorderRouterCreate,
		Read:   resourceAliyunVirtualBorderRouterRead,
		Update: resourceAliyunVirtualBorderRouterUpdate,
		Delete: resourceAliyunVirtualBorderRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"physical_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2999),
			},
			"local_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peering_subnet_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"vlan_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVirtualBorderRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateCreateVirtualBorderRouterRequest()
	request.RegionId = client.RegionId
	request.PhysicalConnectionId = d.Get("physical_connection_id").(string)
	request.VlanId = requests.NewInteger(d.Get("vlan_id").(int))
	request.LocalGatewayIp = d.Get("local_gateway_ip").(string)
	request.PeerGatewayIp = d.Get("peer_gateway_ip").(string)
	request.PeeringSubnetMask = d.Get("peering_subnet_mask").(string)

	if v, ok := d.GetOk("name"); ok && v != "" {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v != "" {
		request.Description = v.(string)
	}

	request.ClientToken = buildClientToken(request.GetActionName())

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVirtualBorderRouter(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateVirtualBorderRouterResponse)
		d.SetId(response.VbrId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_virtual_border_router", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunVirtualBorderRouterRead(d, meta)
}

func resourceAliyunVirtualBorderRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}
	object, err := vpcService.DescribeVirtualBorderRouter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("physical_connection_id", object.PhysicalConnectionId)
	d.Set("vlan_id", object.VlanId)
	d.Set("local_gateway_ip", object.LocalGatewayIp)
	d.Set("peer_gateway_ip", object.PeerGatewayIp)
	d.Set("peering_subnet_mask", object.PeeringSubnetMask)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("vlan_interface_id", object.VlanInterfaceId)
	d.Set("route_table_id", object.RouteTableId)
	return nil
}

func resourceAliyunVirtualBorderRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	update := false
	request := vpc.CreateModifyVirtualBorderRouterAttributeRequest()
	request.RegionId = client.RegionId
	request.VbrId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	if d.HasChange("vlan_id") {
		request.VlanId = requests.NewInteger(d.Get("vlan_id").(int))
		update = true
		d.SetPartial("vlan_id")
	}

	if d.HasChange("local_gateway_ip") {
		request.LocalGatewayIp = d.Get("local_gateway_ip").(string)
		update = true
		d.SetPartial("local_gateway_ip")
	}

	if d.HasChange("peer_gateway_ip") {
		request.PeerGatewayIp = d.Get("peer_gateway_ip").(string)
		update = true
		d.SetPartial("peer_gateway_ip")
	}

	if d.HasChange("peering_subnet_mask") {
		request.PeeringSubnetMask = d.Get("peering_subnet_mask").(string)
		update = true
		d.SetPartial("peering_subnet_mask")
	}

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
		d.SetPartial("name")
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
		d.SetPartial("description")
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVirtualBorderRouterAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAliyunVirtualBorderRouterRead(d, meta)
}

func resourceAliyunVirtualBorderRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVirtualBorderRouterRequest()
	request.RegionId = client.RegionId
	request.VbrId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVirtualBorderRouter(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVbrId.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForVirtualBorderRouter(d.Id(), Deleted, DefaultTimeout))

}
