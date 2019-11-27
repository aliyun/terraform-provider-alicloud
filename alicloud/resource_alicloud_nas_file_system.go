package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNasFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasFileSystemCreate,
		Read:   resourceAlicloudNasFileSystemRead,
		Update: resourceAlicloudNasFileSystemUpdate,
		Delete: resourceAlicloudNasFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Capacity",
					"Performance",
				}, false),
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS",
					"SMB",
				}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
		},
	}
}

func resourceAlicloudNasFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateFileSystemRequest()
	request.RegionId = string(client.RegionId)
	request.ProtocolType = d.Get("protocol_type").(string)
	request.StorageType = d.Get("storage_type").(string)
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateFileSystem(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_file_system", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*nas.CreateFileSystemResponse)
	d.SetId(response.FileSystemId)
	return resourceAlicloudNasFileSystemUpdate(d, meta)
}

func resourceAlicloudNasFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateModifyFileSystemRequest()
	request.RegionId = client.RegionId
	request.FileSystemId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyFileSystem(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudNasFileSystemRead(d, meta)
}

func resourceAlicloudNasFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasFileSystem(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", object.Description)
	d.Set("protocol_type", object.ProtocolType)
	d.Set("storage_type", object.StorageType)
	return nil
}

func resourceAlicloudNasFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	request := nas.CreateDeleteFileSystemRequest()
	request.RegionId = client.RegionId
	request.FileSystemId = d.Id()

	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.DeleteFileSystem(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidFileSystemIDNotFound, ForbiddenNasNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(nasService.WaitForNasFileSystem(d.Id(), Deleted, DefaultTimeoutMedium))
}
