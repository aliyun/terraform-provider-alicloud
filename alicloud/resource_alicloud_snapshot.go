package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSnapshotCreate,
		Read:   resourceAliyunSnapshotRead,
		Update: resourceAliyunSnapshotUpdate,
		Delete: resourceAliyunSnapshotDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ecs.CreateCreateSnapshotRequest()
	args.DiskId = d.Get("disk_id").(string)
	args.ClientToken = buildClientToken(args.GetActionName())
	if name, ok := d.GetOk("name"); ok {
		args.SnapshotName = name.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		args.Description = description.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSnapshot(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultDebugMsg, "alicloud_snapshot", args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(args.GetActionName(), raw)
	resp := raw.(*ecs.CreateSnapshotResponse)
	d.SetId(resp.SnapshotId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshot(d.Id(), SnapshotCreatingAccomplished, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSnapshotUpdate(d, meta)
}

func resourceAliyunSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	snapshot, err := ecsService.DescribeSnapshotById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", snapshot.SnapshotName)
	d.Set("disk_id", snapshot.SourceDiskId)
	d.Set("description", snapshot.Description)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceSnapshot)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if err := setTags(client, TagResourceSnapshot, d); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "setTags", ProviderERROR)
	}

	return resourceAliyunSnapshotRead(d, meta)
}

func resourceAliyunSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	args := ecs.CreateDeleteSnapshotRequest()
	args.SnapshotId = d.Id()

	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSnapshot(args)
		})
		if err != nil {
			if IsExceptedError(err, SnapshotNotFound) {
				return nil
			}
			if IsExceptedErrors(err, SnapshotInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForSnapshot(d.Id(), Deleted, DefaultLongTimeout))
}
