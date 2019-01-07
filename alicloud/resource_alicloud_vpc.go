package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpcCreate,
		Read:   resourceAliyunVpcRead,
		Update: resourceAliyunVpcUpdate,
		Delete: resourceAliyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(2, 256),
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	var vpcResponse *vpc.CreateVpcResponse
	request := buildAliyunVpcArgs(d, meta)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpc(&args)
		})
		if err != nil {
			if IsExceptedError(err, VpcQuotaExceeded) {
				return resource.NonRetryableError(fmt.Errorf("The number of VPC has quota has reached the quota limit in your account, and please use existing VPCs or remove some of them."))
			}
			if IsExceptedErrors(err, []string{TaskConflict, UnknownError, Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create vpc timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}
		vpcResponse, _ = raw.(*vpc.CreateVpcResponse)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create vpc got an error :%#v", err)
	}

	d.SetId(vpcResponse.VpcId)

	err = vpcService.WaitForVpc(d.Id(), Available, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForVpcAvailable")
	}

	return resourceAliyunVpcUpdate(d, meta)
}

func resourceAliyunVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	resp, err := vpcService.DescribeVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cidr_block", resp.CidrBlock)
	d.Set("name", resp.VpcName)
	d.Set("description", resp.Description)
	d.Set("router_id", resp.VRouterId)
	request := vpc.CreateDescribeVRoutersRequest()
	request.RegionId = client.RegionId
	request.VRouterId = resp.VRouterId
	var response vpc.DescribeVRoutersResponse
	if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, e := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVRouters(request)
		})
		if e != nil && IsExceptedErrors(err, []string{Throttling}) {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(e)
		}
		r, _ := raw.(*vpc.DescribeVRoutersResponse)
		response = *r
		return resource.NonRetryableError(e)
	}); err != nil {
		return fmt.Errorf("DescribeVRouters got an error: %#v.", err)
	}
	if len(response.VRouters.VRouter) > 0 && len(response.VRouters.VRouter[0].RouteTableIds.RouteTableId) > 0 {
		d.Set("router_table_id", response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
		d.Set("route_table_id", response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
	} else {
		d.Set("router_table_id", "")
		d.Set("route_table_id", "")
	}

	return nil
}

func resourceAliyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	request.VpcId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		request.VpcName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpcAttribute(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpc(request)
		})

		if err != nil {
			if IsExceptedError(err, InvalidVpcIDNotFound) || IsExceptedError(err, ForbiddenVpcNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete VPC timeout and got an error: %#v.", err))
		}

		if _, err := vpcService.DescribeVpc(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunVpcArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateVpcRequest {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateCreateVpcRequest()
	request.RegionId = string(client.Region)
	request.CidrBlock = d.Get("cidr_block").(string)

	if v := d.Get("name").(string); v != "" {
		request.VpcName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	request.ClientToken = buildClientToken("TF-CreateVpc")

	return request
}
