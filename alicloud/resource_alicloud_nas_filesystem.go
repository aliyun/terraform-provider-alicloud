package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNASFilesystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNASFilesystemCreate,
		Read:   resourceAlicloudNASFilesystemRead,
		Update: resourceAlicloudNASFilesystemUpdate,
		Delete: resourceAlicloudNASFilesystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Capacity",
					"Performance",
				}),
			},
			"protocol_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"NFS",
					"SMB",
				}),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNASDescription,
			},
		},
	}
}

func resourceAlicloudNASFilesystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateFileSystemRequest()
	request.RegionId = string(client.Region)
	request.ProtocolType = d.Get("protocol_type").(string)
	request.StorageType = d.Get("storage_type").(string)

	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateFileSystem(request)
	})
	fs, _ := raw.(*nas.CreateFileSystemResponse)
	if err != nil {
		return fmt.Errorf("Error Waitting for NAT available: %#v", err)
	}

	d.SetId(fs.FileSystemId)
	return resourceAlicloudNASFilesystemUpdate(d, meta)
}

func resourceAlicloudNASFilesystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	//split := strings.Split(d.Id(), COLON_SEPARATED)
	attributeUpdate := false
	request := nas.CreateModifyFileSystemRequest()
	request.FileSystemId = d.Id()

	if d.HasChange("description") {
		attributeUpdate = true
		d.SetPartial("description")
		request.Description = d.Get("description").(string)
	}
	if attributeUpdate {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyFileSystem(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudNASFilesystemRead(d, meta)
}

func resourceAlicloudNASFilesystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	//split := strings.Split(d.Id(), COLON_SEPARATED)
	resp, err := nasService.DescribeFileSystems(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("create_time", resp.CreateTime)
	d.Set("description", resp.Destription)
	d.Set("metered_size", resp.MeteredSize)
	d.Set("protocol_type", resp.ProtocolType)
	d.Set("storage_type", resp.StorageType)
	d.Set("FileSystemId", resp.FileSystemId)

	return nil
}

func resourceAlicloudNASFilesystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	//split := strings.Split(d.Id(), COLON_SEPARATED)
	request := nas.CreateDeleteFileSystemRequest()
	request.FileSystemId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteFileSystem(request)
		})

		if err != nil {
			if IsExceptedError(err, ForbiddenNasNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete NAS timeout and got an error: %#v.", err))
		}

		if _, err := nasService.DescribeFileSystems(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}
