package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunVpnRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnRouteEntryCreate,
		Read:   resourceAliyunVpnRouteEntryRead,
		Update: resourceAliyunVpnRouteEntryUpdate,
		Delete: resourceAliyunVpnRouteEntryDelete,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
			},

			"route_dest": {
				Type:     schema.TypeString,
				Required: true,
			},

			"weight": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateVpnBandwidth([]int{0, 100}),
			},

			"publish_vpc": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 1 || len(value) > 64 {
						errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters or shorter than 1", k))
					}
					return
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"create_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"old_weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnRouteEntryService := VpnRouteEntryService{client}
	request := vpc.CreateCreateVpnRouteEntryRequest()
	request.RegionId = client.RegionId
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.RouteDest = d.Get("route_dest").(string)
	request.NextHop = d.Get("next_hop").(string)
	request.Weight = requests.NewInteger(d.Get("weight").(int))
	request.PublishVpc = requests.NewBoolean(d.Get("publish_vpc").(bool))

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = d.Get("description").(string)
	}

	if v, ok := d.GetOk("client_token"); ok && v.(string) != "" {
		request.ClientToken = d.Get("client_token").(string)
	}

	time.Sleep(7 * time.Second)
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateVpnRouteEntry(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_route_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.CreateVpnRouteEntryResponse)

	id := response.VpnInstanceId + ":" + response.NextHop + response.RouteDest
	d.SetId(id)
	d.Set("old_weight", d.Get("weight").(int))

	time.Sleep(10 * time.Second)
	if err := vpnRouteEntryService.WaitForVpnRouteEntry(d.Id(), Active, 2*DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpnRouteEntryUpdate(d, meta)
}

func resourceAliyunVpnRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnRouteEntryService := VpnRouteEntryService{client}

	object, err := vpnRouteEntryService.DescribeVpnRouteEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", object.CreateTime)
	d.Set("state", object.State)
	return nil
}

func resourceAliyunVpnRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pubChange := d.HasChange("publish_vpc")
	weightChange := d.HasChange("weight")
	d.Partial(true)

	if pubChange {
		request := vpc.CreatePublishVpnRouteEntryRequest()
		request.RegionId = client.RegionId
		request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
		request.RouteDest = d.Get("route_dest").(string)
		request.NextHop = d.Get("next_hop").(string)
		request.RouteType = "dbr"
		request.PublishVpc = requests.NewBoolean(d.Get("publish_vpc").(bool))

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.PublishVpnRouteEntry(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("publish_vpc")
	}

	if weightChange {
		request := vpc.CreateModifyVpnRouteEntryWeightRequest()
		request.RegionId = client.RegionId
		request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
		request.RouteDest = d.Get("route_dest").(string)
		request.NextHop = d.Get("next_hop").(string)
		newWeight := d.Get("weight").(int)
		request.Weight = requests.NewInteger(d.Get("old_weight").(int))
		request.NewWeight = requests.NewInteger(newWeight)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnRouteEntryWeight(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)

		d.Set("weight", newWeight)
		d.Set("old_weight", newWeight)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunVpnRouteEntryRead(d, meta)
	}

	if d.HasChange("route_dest") {
		return fmt.Errorf("Now Cann't Support modify vpn RouteDest, try to modify on the web console")
	}
	if d.HasChange("next_hop") {
		return fmt.Errorf("Now Cann't Support modify vpn NextHop, try to modify on the web console")
	}

	d.Partial(false)
	return resourceAliyunVpnRouteEntryRead(d, meta)

}

func resourceAliyunVpnRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnRouteEntryService := VpnRouteEntryService{client}

	request := vpc.CreateDeleteVpnRouteEntryRequest()
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.RouteDest = d.Get("route_dest").(string)
	request.NextHop = d.Get("next_hop").(string)
	request.Weight = requests.NewInteger(d.Get("weight").(int))

	request.ClientToken = buildClientToken(request.GetActionName())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnRouteEntry(&args)
		})
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			/*Vpn known issue: while the vpn is configuring, it will return unknown error*/
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		if IsExceptedError(err, VpnNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(vpnRouteEntryService.WaitForVpnRouteEntry(d.Id(), Deleted, DefaultTimeoutMedium))
}
