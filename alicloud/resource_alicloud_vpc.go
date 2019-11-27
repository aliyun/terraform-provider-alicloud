package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"tags": tagsSchema(),
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

	var response *vpc.CreateVpcResponse
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
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateVpcResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.VpcId)

	err = vpcService.WaitForVpc(d.Id(), Available, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpcUpdate(d, meta)
}

func resourceAliyunVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", object.CidrBlock)
	d.Set("name", object.VpcName)
	d.Set("description", object.Description)
	d.Set("router_id", object.VRouterId)
	d.Set("resource_group_id", object.ResourceGroupId)
	tags, err := vpcService.DescribeTags(d.Id(), nil, TagResourceVpc)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", vpcService.tagsToMap(tags))
	// Retrieve all route tables and filter to get system
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = client.RegionId
	request.VRouterId = object.VRouterId
	request.ResourceGroupId = object.ResourceGroupId
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	var routeTabls []vpc.RouteTable
	for {
		total := 0
		err = resource.Retry(6*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			if err != nil && IsExceptedErrors(err, []string{Throttling}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*vpc.DescribeRouteTablesResponse)
			routeTabls = append(routeTabls, response.RouteTables.RouteTable...)
			total = len(response.RouteTables.RouteTable)
			return resource.NonRetryableError(err)
		})
		if err != nil {
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
	vpcService := VpcService{client}
	if err := vpcService.setInstanceTags(d, TagResourceVpc); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunVpcRead(d, meta)
	}
	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	request.RegionId = client.RegionId
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
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpcAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVpcRequest()
	request.RegionId = client.RegionId
	request.VpcId = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpc(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound}) {
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
	return WrapError(vpcService.WaitForVpc(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildAliyunVpcArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateVpcRequest {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateCreateVpcRequest()
	request.RegionId = client.RegionId
	request.CidrBlock = d.Get("cidr_block").(string)

	if v := d.Get("name").(string); v != "" {
		request.VpcName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}

	request.ClientToken = buildClientToken(request.GetActionName())

	return request
}
