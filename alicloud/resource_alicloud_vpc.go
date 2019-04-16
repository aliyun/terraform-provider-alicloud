package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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
				return resource.NonRetryableError(WrapErrorf(err, "The number of VPC has quota has reached the quota limit in your account, and please use existing VPCs or remove some of them."))
			}
			if IsExceptedErrors(err, []string{TaskConflict, UnknownError, Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		vpcResponse, _ = raw.(*vpc.CreateVpcResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "new", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(vpcResponse.VpcId)

	err = vpcService.WaitForVpc(d.Id(), Available, 60)
	if err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpcRead(d, meta)
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
		return WrapError(err)
	}

	d.Set("cidr_block", resp.CidrBlock)
	d.Set("name", resp.VpcName)
	d.Set("description", resp.Description)
	d.Set("router_id", resp.VRouterId)
	// Retrieve all route tables and filter to get system
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = client.RegionId
	request.VRouterId = resp.VRouterId
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	var routeTabls []vpc.RouteTable
	for {
		total := 0
		if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			raw, e := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			if e != nil && IsExceptedErrors(e, []string{Throttling}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(e)
			}
			r, _ := raw.(*vpc.DescribeRouteTablesResponse)
			routeTabls = append(routeTabls, r.RouteTables.RouteTable...)
			total = len(r.RouteTables.RouteTable)
			return resource.NonRetryableError(e)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if total < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	// Generally, the system route table is the last one
	for i := len(routeTabls) - 1; i >= 0; i-- {
		if routeTabls[i].RouteTableType == "System" {
			d.Set("router_table_id", routeTabls[i].RouteTableId)
			d.Set("route_table_id", routeTabls[i].RouteTableId)
			break
		}
	}

	return nil
}

func resourceAliyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	request.VpcId = d.Id()

	if d.HasChange("name") {
		request.VpcName = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpcAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

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
			if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound}) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		if _, err := vpcService.DescribeVpc(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
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
	request.ClientToken = buildClientToken(request.GetActionName())

	return request
}
