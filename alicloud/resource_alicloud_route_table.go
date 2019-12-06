package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteTableCreate,
		Read:   resourceAliyunRouteTableRead,
		Update: resourceAliyunRouteTableUpdate,
		Delete: resourceAliyunRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateRouteTableRequest()
	request.RegionId = client.RegionId

	request.VpcId = d.Get("vpc_id").(string)
	request.RouteTableName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouteTable(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_table", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateRouteTableResponse)
	d.SetId(response.RouteTableId)

	if err := vpcService.WaitForRouteTable(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAliyunRouteTableUpdate(d, meta)
}

func resourceAliyunRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouteTable(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("vpc_id", object.VpcId)
	d.Set("name", object.RouteTableName)
	d.Set("description", object.Description)
	tags, err := vpcService.DescribeTags(d.Id(), nil, TagResourceRouteTable)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", vpcService.tagsToMap(tags))
	return nil
}

func resourceAliyunRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if err := vpcService.setInstanceTags(d, TagResourceRouteTable); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunRouteTableRead(d, meta)
	}
	request := vpc.CreateModifyRouteTableAttributesRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		request.RouteTableName = d.Get("name").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyRouteTableAttributes(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAliyunRouteTableRead(d, meta)
}

func resourceAliyunRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := VpcService{client}
	request := vpc.CreateDeleteRouteTableRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = d.Id()

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteRouteTable(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(routeTableService.WaitForRouteTable(d.Id(), Deleted, DefaultTimeoutMedium))
}
