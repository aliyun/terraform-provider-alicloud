package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNetworkAclCreate,
		Read:   resourceAliyunNetworkAclRead,
		Update: resourceAliyunNetworkAclUpdate,
		Delete: resourceAliyunNetworkAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
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

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateNetworkAclRequest()
	request.RegionId = client.RegionId

	request.VpcId = d.Get("vpc_id").(string)
	if networkAclName, ok := d.GetOk("name"); ok {
		request.NetworkAclName = networkAclName.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateNetworkAcl(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_network_acl", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateNetworkAclResponse)
	d.SetId(response.NetworkAclId)

	if err := vpcService.WaitForNetworkAcl(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunNetworkAclRead(d, meta)
}

func resourceAliyunNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeNetworkAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("vpc_id", object.VpcId)
	d.Set("name", object.NetworkAclName)
	d.Set("description", object.Description)

	return nil
}

func resourceAliyunNetworkAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateModifyNetworkAclAttributesRequest()
	request.RegionId = client.RegionId
	request.NetworkAclId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		request.NetworkAclName = d.Get("name").(string)
	}
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyNetworkAclAttributes(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := vpcService.WaitForNetworkAcl(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunNetworkAclRead(d, meta)
}

func resourceAliyunNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if err := vpcService.WaitForNetworkAcl(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	request := vpc.CreateDeleteNetworkAclRequest()
	request.RegionId = client.RegionId
	request.NetworkAclId = d.Id()
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteNetworkAcl(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(vpcService.WaitForNetworkAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
