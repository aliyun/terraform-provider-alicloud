package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpcCidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcCidrBlockCreate,
		Read:   resourceAlicloudVpcCidrBlockRead,
		Delete: resourceAlicloudVpcCidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secondary_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudVpcCidrBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	vpcId := d.Get("vpc_id").(string)
	secondaryCidrBlock := d.Get("secondary_cidr_block").(string)
	_, err := vpcService.DescribeVpc(vpcId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	request := vpc.CreateAssociateVpcCidrBlockRequest()
	request.RegionId = client.RegionId
	request.VpcId = vpcId
	request.SecondaryCidrBlock = secondaryCidrBlock

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AssociateVpcCidrBlock(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_cidr_block", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(vpcId + ":" + secondaryCidrBlock)
	return resourceAlicloudVpcCidrBlockRead(d, meta)

}

func resourceAlicloudVpcCidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := vpcService.DescribeVpc(parts[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	for _, cidr := range object.SecondaryCidrBlocks.SecondaryCidrBlock {
		if cidr == parts[1] {
			d.Set("secondary_cidr_block", cidr)
		}
	}

	return nil
}

func resourceAlicloudVpcCidrBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateUnassociateVpcCidrBlockRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.VpcId = parts[0]
	request.SecondaryCidrBlock = parts[1]
	request.RegionId = client.RegionId

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.UnassociateVpcCidrBlock(request)
	})

	if err != nil {
		//if IsExceptedError(err, "InvalidInstanceId") {
		//	return nil
		//}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
