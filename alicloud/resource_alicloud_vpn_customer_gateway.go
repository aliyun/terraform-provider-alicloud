package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

	client := meta.(*AliyunClient)

	//var vpc *vpc.CreateVpcResponse
	var cgw *vpc.CreateCustomerGatewayResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunCustomerGatewayArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateCustomerGateway got an error: %#v", err))
		}
		resp, err := client.vpcconn.CreateCustomerGateway(args)
		if err != nil {
			/*if IsExceptedError(err, ResQuotaFull) {
				return resource.NonRetryableError(fmt.Errorf("The quota of resource is full"))
			}
			if IsExceptedError(err, VpnForbidden) {
				return resource.RetryableError(fmt.Errorf("User not authorized to operate on the specified resource. %#v.", err))
			}
			if IsExceptedError(err, InvalidIpAddress) {
				return resource.RetryableError(fmt.Errorf("Specified IpAddress is already exist. %#v.", err))
			}*/
			return resource.NonRetryableError(err)
		}
		cgw = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create vpn customer gateway got an error :%#v", err)
	}

	d.SetId(cgw.CustomerGatewayId)

	err = client.WaitForCustomerGateway(d.Id(), Available, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitforCustomerGateway")
	}

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	resp, err := client.DescribeCustomerGateway(d.Id())
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

	//d.Partial(true)
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
		if _, err := meta.(*AliyunClient).vpcconn.ModifyCustomerGatewayAttribute(request); err != nil {
			return err
		}
	}

	//d.Partial(false)

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := vpc.CreateDeleteCustomerGatewayRequest()
	request.CustomerGatewayId = d.Id()
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteCustomerGateway(request)

		if err != nil {
			if IsExceptedError(err, CgwNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete CustomerGateway timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeCustomerGateway(d.Id()); err != nil {
			if IsExceptedError(err, CgwNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunCustomerGatewayArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateCustomerGatewayRequest, error) {
	request := vpc.CreateCreateCustomerGatewayRequest()
	request.RegionId = string(getRegion(d, meta))
	request.IpAddress = d.Get("ip_address").(string)

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	return request, nil
}
