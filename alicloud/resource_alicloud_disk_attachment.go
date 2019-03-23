package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskAttachmentCreate,
		Read:   resourceAliyunDiskAttachmentRead,
		Delete: resourceAliyunDiskAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Attribute device_name is deprecated on disk attachment resource. Suggest to remove it from your template.",
			},
		},
	}
}

func resourceAliyunDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	diskID := d.Get("disk_id").(string)
	instanceID := d.Get("instance_id").(string)

	request := ecs.CreateAttachDiskRequest()
	request.InstanceId = instanceID
	request.DiskId = diskID

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachDisk(request)
		})

		if err != nil {
			if IsExceptedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_disk_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(request.DiskId + ":" + request.InstanceId)

	if err := ecsService.WaitForDiskAttachment(d.Id(), DiskInUse, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunDiskAttachmentRead(d, meta)
}

func resourceAliyunDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	disk, err := ecsService.DescribeDiskAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", disk.InstanceId)
	d.Set("disk_id", disk.DiskId)
	d.Set("device_name", disk.Device)

	return nil
}

func resourceAliyunDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := ecs.CreateDetachDiskRequest()
	request.InstanceId = parts[1]
	request.DiskId = parts[0]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachDisk(request)
		})
		if err != nil {
			if IsExceptedErrors(err, DiskInvalidOperation) {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForDiskAttachment(d.Id(), Deleted, DefaultTimeout))
}
