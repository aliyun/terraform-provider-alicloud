package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskAttachmentCreate,
		Read:   resourceAliyunDiskAttachmentRead,
		Delete: resourceAliyunDiskAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"device_name": &schema.Schema{
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
	client := meta.(*AliyunClient)

	diskID := d.Get("disk_id").(string)
	instanceID := d.Get("instance_id").(string)

	args := ecs.CreateAttachDiskRequest()
	args.InstanceId = instanceID
	args.DiskId = diskID

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.ecsconn.AttachDisk(args)

		if err != nil {
			if IsExceptedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Attach Disk %s timeout and got an error: %#v", diskID, err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Attaching disk %s to instance %s got an error: %#v.", diskID, instanceID, err)
	}

	if err := client.WaitForEcsDisk(diskID, DiskInUse, DefaultTimeout); err != nil {
		return fmt.Errorf("Waitting for disk %s %s got an error: %#v.", diskID, DiskInUse, err)
	}

	d.SetId(d.Get("disk_id").(string) + ":" + d.Get("instance_id").(string))

	return resourceAliyunDiskAttachmentRead(d, meta)
}

func resourceAliyunDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	diskId, instanceId, err := getDiskIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	disk, err := meta.(*AliyunClient).DescribeDiskById(instanceId, diskId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeDiskAttribute: %#v", err)
	}

	d.Set("instance_id", disk.InstanceId)
	d.Set("disk_id", disk.DiskId)
	d.Set("device_name", disk.Device)

	return nil
}

func resourceAliyunDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	diskID, instanceID, err := getDiskIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	req := ecs.CreateDetachDiskRequest()
	req.InstanceId = instanceID
	req.DiskId = diskID

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		disk, err := client.DescribeDiskById(instanceID, diskID)

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While detach disk %s, describing disk got an error: %#v.", diskID, err))
		}

		if disk.Status == string(Available) {
			return nil
		}

		_, err = client.ecsconn.DetachDisk(req)
		if err != nil {
			if IsExceptedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Detach Disk %s timeout and got an error: %#v", diskID, err))
			}
			if IsExceptedErrors(err, []string{DependencyViolation}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Detaching disk %s got an error: %#v.", diskID, err))
		}
		return resource.RetryableError(fmt.Errorf("Detach Disk timeout and got an error: %#v", err))
	})
}

func getDiskIDAndInstanceID(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
