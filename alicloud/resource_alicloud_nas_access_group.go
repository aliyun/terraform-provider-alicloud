package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNasAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAccessGroupCreate,
		Read:   resourceAlicloudNasAccessGroupRead,
		Update: resourceAlicloudNasAccessGroupUpdate,
		Delete: resourceAlicloudNasAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Vpc", "Classic"}, false),
				ForceNew:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudNasAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessGroupRequest()
	request.RegionId = string(client.Region)
	request.AccessGroupName = d.Get("name").(string)
	request.AccessGroupType = d.Get("type").(string)
	request.Description = d.Get("description").(string)
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*nas.CreateAccessGroupResponse)
	d.SetId(response.AccessGroupName)
	return resourceAlicloudNasAccessGroupRead(d, meta)
}

func resourceAlicloudNasAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateModifyAccessGroupRequest()
	request.RegionId = client.RegionId
	request.AccessGroupName = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyAccessGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudNasAccessGroupRead(d, meta)
}

func resourceAlicloudNasAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasAccessGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.AccessGroupName)
	d.Set("type", object.AccessGroupType)
	d.Set("description", object.Description)

	return nil
}

func resourceAlicloudNasAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	request := nas.CreateDeleteAccessGroupRequest()
	request.RegionId = client.RegionId
	request.AccessGroupName = d.Id()

	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.DeleteAccessGroup(request)
	})

	if err != nil {
		if IsExceptedError(err, ForbiddenNasNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(nasService.WaitForNasAccessGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
