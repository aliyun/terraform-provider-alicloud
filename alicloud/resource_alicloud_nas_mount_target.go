package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"file_system_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Active", "Inactive"}),
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
		vsw, err := vpcService.DescribeVswitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request.VpcId = vsw.VpcId
	}
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateMountTarget(request)
	})
	fs, _ := raw.(*nas.CreateMountTargetResponse)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "nas_mount_target", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fs.MountTargetDomain)
	err = nasService.WaitForMountTarget(d.Id(), Active, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudNasMountTargetUpdate(d, meta)
}

func resourceAlicloudNasMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	split := strings.Split(d.Id(), "-")
	update := false
	request := nas.CreateModifyMountTargetRequest()
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
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyMountTarget(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudNasMountTargetRead(d, meta)
}

func resourceAlicloudNasMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	resp, err := nasService.DescribeNasMountTarget(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", resp.Status)
	d.Set("access_group_name", resp.AccessGroup)
	d.Set("vswitch_id", resp.VswId)
	d.Set("file_system_id", strings.Split(resp.MountTargetDomain, "-")[0])
	return nil
}

func resourceAlicloudNasMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), "-")
	request := nas.CreateDeleteMountTargetRequest()
	request.FileSystemId = split[0]
	request.MountTargetDomain = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteMountTarget(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidFileSystemIDNotFound, ForbiddenNasNotFound}) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		if _, err := nasService.DescribeNasMountTarget(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
