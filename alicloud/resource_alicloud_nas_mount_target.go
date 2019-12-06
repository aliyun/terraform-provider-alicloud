package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNasMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasMountTargetCreate,
		Read:   resourceAlicloudNasMountTargetRead,
		Update: resourceAlicloudNasMountTargetUpdate,
		Delete: resourceAlicloudNasMountTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
		},
	}
}

func resourceAlicloudNasMountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	vpcService := VpcService{client}
	request := nas.CreateCreateMountTargetRequest()
	request.RegionId = string(client.Region)
	request.FileSystemId = d.Get("file_system_id").(string)
	request.AccessGroupName = d.Get("access_group_name").(string)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	request.NetworkType = string(Classic)
	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.NetworkType = string(Vpc)
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request.VpcId = vsw.VpcId
	}
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateMountTarget(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_mount_target", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*nas.CreateMountTargetResponse)
	d.SetId(response.MountTargetDomain)
	err = nasService.WaitForNasMountTarget(d.Id(), Active, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlicloudNasMountTargetUpdate(d, meta)
}

func resourceAlicloudNasMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	split := strings.Split(d.Id(), "-")
	update := false
	request := nas.CreateModifyMountTargetRequest()
	request.RegionId = client.RegionId
	request.FileSystemId = split[0]
	request.MountTargetDomain = d.Id()
	if !d.IsNewResource() && d.HasChange("access_group_name") {
		request.AccessGroupName = d.Get("access_group_name").(string)
		update = true
	}
	if d.HasChange("status") {
		request.Status = d.Get("status").(string)
		update = true
	}
	if update {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyMountTarget(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudNasMountTargetRead(d, meta)
}

func resourceAlicloudNasMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasMountTarget(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object.Status)
	d.Set("access_group_name", object.AccessGroup)
	d.Set("vswitch_id", object.VswId)
	d.Set("file_system_id", strings.Split(object.MountTargetDomain, "-")[0])
	return nil
}

func resourceAlicloudNasMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), "-")
	request := nas.CreateDeleteMountTargetRequest()
	request.RegionId = client.RegionId
	request.FileSystemId = split[0]
	request.MountTargetDomain = d.Id()
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.DeleteMountTarget(request)
	})

	if err != nil {
		if IsExceptedError(err, ForbiddenNasNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(nasService.WaitForNasMountTarget(d.Id(), Deleted, DefaultTimeoutMedium))
}
