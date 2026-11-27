package ots

import (
	"github.com/alibabacloud-go/tea/dara"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOtsInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceAttachmentCreate,
		Read:   resourceAliyunOtsInstanceAttachmentRead,
		Delete: resourceAliyunOtsInstanceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},

			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunOtsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)

	request := ots.CreateBindInstance2VpcRequest()
	request.RegionId = client.RegionId
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceVpcName = d.Get("vpc_name").(string)
	request.VirtualSwitchId = d.Get("vswitch_id").(string)

	if vsw, err := otsService.describeVSwitch(d.Get("vswitch_id").(string)); err != nil {
		return errs.WrapError(err)
	} else {
		request.VpcId = vsw.VpcId
	}

	var raw interface{}
	var err error
	if client.UseOtsV4() {
		// sign_version.default = v4 signs the control plane through the
		// common V4 RPC path; the legacy POP client cannot sign V4.
		raw, err = client.OtsPopV4("POST", "BindInstance2Vpc", map[string]interface{}{
			"InstanceName":    request.InstanceName,
			"InstanceVpcName": request.InstanceVpcName,
			"VpcId":           request.VpcId,
			"VirtualSwitchId": request.VirtualSwitchId,
		})
	} else {
		raw, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.BindInstance2Vpc(request)
		})
	}
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ots_instance_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(request.InstanceName)
	return resourceAliyunOtsInstanceAttachmentRead(d, meta)
}

func resourceAliyunOtsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	// There is a bug that inst does not contain instance name and vswitch ID, so this resource does not support import function.
	//d.Set("instance_name", inst.InstanceName)
	d.Set("vpc_name", object.InstanceVpcName)
	d.Set("vpc_id", object.VpcId)
	return nil
}

func resourceAliyunOtsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	request := ots.CreateUnbindInstance2VpcRequest()
	request.RegionId = client.RegionId
	request.InstanceName = d.Id()
	request.InstanceVpcName = object.InstanceVpcName

	var raw interface{}
	if client.UseOtsV4() {
		// sign_version.default = v4 signs the control plane through the
		// common V4 RPC path; the legacy POP client cannot sign V4.
		raw, err = client.OtsPopV4("POST", "UnbindInstance2Vpc", map[string]interface{}{
			"InstanceName":    request.InstanceName,
			"InstanceVpcName": request.InstanceVpcName,
		})
	} else {
		raw, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UnbindInstance2Vpc(request)
		})
	}
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(otsService.WaitForOtsInstanceVpc(d.Id(), Deleted, DefaultTimeout))
}

// describeVSwitch mirrors alicloud.VpcService.DescribeVSwitch so the OTS
// instance attachment can resolve a VPC id from a vswitch without importing the
// (not yet split) root alicloud package. Behavior is preserved verbatim.
func (s *Service) describeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	if s.client.UseSignV4("vpc") {
		// The vswitch lookup is a VPC request, so it follows the VPC effective
		// sign version (the default), not the "ots" override: with
		// sign_version{default = v4, ots = v1} the lookup still signs V4 while
		// the bind itself goes through the legacy POP client. The shared
		// legacy VPC client cannot sign V4, so the lookup goes through the
		// common layer and the raw response is converted into the same struct,
		// mirroring listVpcInfoByInstance.
		data, err := s.client.RpcPost("Vpc", "2016-04-28", "DescribeVSwitchAttributes", map[string]interface{}{
			"RegionId":  s.client.RegionId,
			"VSwitchId": id,
		}, nil, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return v, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return v, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeVSwitchAttributes", errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug("DescribeVSwitchAttributes", data)
		resp := vpc.CreateDescribeVSwitchAttributesResponse()
		if err := dara.Convert(data, resp); err != nil {
			return v, err
		}
		if resp.VSwitchId != id {
			return v, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return *resp, nil
	}
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}
