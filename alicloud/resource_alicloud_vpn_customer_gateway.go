package alicloud

import (
	"fmt"
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

		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIpAddress,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": &schema.Schema{
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
	var cgw *vpc.CreateCustomerGatewayResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := buildAliyunCustomerGatewayArgs(d, meta)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateCustomerGateway(args)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		cgw, _ = raw.(*vpc.CreateCustomerGatewayResponse)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create vpn customer gateway got an error :%#v", err)
	}

	d.SetId(cgw.CustomerGatewayId)

	err = vpnGatewayService.WaitForCustomerGateway(d.Id(), Available, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitforCustomerGateway: %#v", err)
	}

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	resp, err := vpnGatewayService.DescribeCustomerGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("ip_address", resp.IpAddress)
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)

	return nil
}

func resourceAliyunVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	attributeUpdate := false
	request := vpc.CreateModifyCustomerGatewayAttributeRequest()
	request.CustomerGatewayId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCustomerGatewayAttribute(request)
		})
		if err != nil {
			return err
		}
	}

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteCustomerGatewayRequest()
	request.CustomerGatewayId = d.Id()
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCustomerGateway(request)
		})

		if err != nil {
			if IsExceptedError(err, CgwNotFound) {
				return nil
			}
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Error deleting vpn customer gateway failed: %#v", err))
			}

			return resource.NonRetryableError(fmt.Errorf("Delete CustomerGateway timeout and got an error: %#v.", err))
		}

		if _, err := vpnGatewayService.DescribeCustomerGateway(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(err)
		}

		return nil
	})
}

func buildAliyunCustomerGatewayArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateCustomerGatewayRequest {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateCreateCustomerGatewayRequest()
	request.RegionId = client.RegionId
	request.IpAddress = d.Get("ip_address").(string)

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	return request
}
