package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnCustomerGatewayCreate,
		Read:   resourceAliyunVpnCustomerGatewayRead,
		Update: resourceAliyunVpnCustomerGatewayUpdate,
		Delete: resourceAliyunVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIpAddress,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},
		},
	}
}

func resourceAliyunVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateCustomerGatewayRequest()
	request.RegionId = client.RegionId
	request.IpAddress = d.Get("ip_address").(string)

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	var response *vpc.CreateCustomerGatewayResponse
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateCustomerGateway(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_customer_gateway", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*vpc.CreateCustomerGatewayResponse)

	d.SetId(response.CustomerGatewayId)

	err = vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Null, 60)
	if err != nil {
		return WrapError(err)
	}
	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeVpnCustomerGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ip_address", object.IpAddress)
	d.Set("name", object.Name)
	d.Set("description", object.Description)

	return nil
}

func resourceAliyunVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateModifyCustomerGatewayAttributeRequest()
	request.CustomerGatewayId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyCustomerGatewayAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteCustomerGatewayRequest()
	request.CustomerGatewayId = d.Id()
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCustomerGateway(request)
		})

		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		if IsExceptedError(err, CgwNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Deleted, DefaultTimeout))
}
