package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

func resourceAlicloudNAS_MountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNAS_MountTargetCreate,
		Read:   resourceAlicloudNAS_MountTargetRead,
		Update: resourceAlicloudNAS_MountTargetUpdate,
		Delete: resourceAlicloudNAS_MountTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"filesystem_id": &schema.Schema{
				Type:         schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"networktype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Vpc",
					"Classic",
				}),
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,

			},
			"vsw_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional: true,

			},
			"accessgroup_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,

			},
		},
	}
}

func resourceAlicloudNAS_MountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateCreateMountTargetRequest()
	request.RegionId = string(client.Region)
	request.FileSystemId = d.Get("filesystem_id").(string)
	request.AccessGroupName = d.Get("accessgroup_name").(string)
	request.NetworkType = d.Get("networktype").(string)
	if d.Get("networktype").(string) == "Vpc" {
		request.VpcId = d.Get("vpc_id").(string)
		request.VSwitchId = d.Get("vsw_id").(string)
	}
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateMountTarget(request)
	})
	fs, _ := raw.(*nas.CreateMountTargetResponse)
	if err != nil {
		return fmt.Errorf("Error Waitting for NAS available: %#v", err)
	}

	d.SetId(d.Get("filesystem_id").(string) + ":" + fs.MountTargetDomain)
	return resourceAlicloudNAS_MountTargetUpdate(d, meta)
}

func resourceAlicloudNAS_MountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	split := strings.Split(d.Id(), ":")
	attributeUpdate := false
	request := nas.CreateModifyMountTargetRequest()
	request.FileSystemId = split[0]
	request.MountTargetDomain = split[1]
	if d.HasChange("AccessGroupName") {
		attributeUpdate = true
		d.SetPartial("AccessGroupName")
		request.AccessGroupName = d.Get("accessgroup_name").(string)
	}
	if d.HasChange("status") {
		attributeUpdate = true
		d.SetPartial("status")
		request.AccessGroupName = d.Get("status").(string)
	}
	if attributeUpdate {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyMountTarget(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudNAS_MountTargetRead(d, meta)
}

func resourceAlicloudNAS_MountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), ":")

	resp, err := nasService.DescribeMountTargets(split[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("mounttarget_domain", resp.MountTargetDomain)
	d.Set("network_type", resp.NetworkType)
	d.Set("access_group", resp.AccessGroup)
	if resp.VpcId != "" {
		d.Set("vpc_id", resp.VpcId)
	}
	if resp.VswId != "" {
		d.Set("vsw_id", resp.VswId)
	}

	return nil
}

func resourceAlicloudNAS_MountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), ":")
	request := nas.CreateDeleteMountTargetRequest()
	request.FileSystemId = split[0]
	request.MountTargetDomain = split[1]
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteMountTarget(request)
		})

		if err != nil {
			if IsExceptedError(err, ForbiddenNasNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete NAS timeout and got an error: %#v.", err))
		}

		if _, err := nasService.DescribeMountTargets(split[0]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

