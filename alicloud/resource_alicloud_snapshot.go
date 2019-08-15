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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultTimeout * time.Second),
			Delete: schema.DefaultTimeout(DefaultTimeout * time.Second),
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

	request := ecs.CreateCreateSnapshotRequest()
	request.RegionId = client.RegionId
	request.DiskId = d.Get("disk_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())
	if name, ok := d.GetOk("name"); ok {
		request.SnapshotName = name.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSnapshot(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultDebugMsg, "alicloud_snapshot", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.CreateSnapshotResponse)
	d.SetId(response.SnapshotId)

	ecsService := EcsService{client}

	stateConf := BuildStateConf([]string{}, []string{string(SnapshotCreatingAccomplished)}, d.Timeout(schema.TimeoutCreate), 0,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliyunSnapshotUpdate(d, meta)
}

func resourceAliyunSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	snapshot, err := ecsService.DescribeSnapshot(d.Id())
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
		return WrapError(err)
	}
	return resourceAliyunSnapshotRead(d, meta)
}

func resourceAliyunSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteSnapshotRequest()
	request.RegionId = client.RegionId
	request.SnapshotId = d.Id()

	var raw interface{}
	var err error
	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSnapshot(request)
		})
		if err != nil {
			if IsExceptedErrors(err, SnapshotInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if IsExceptedError(err, SnapshotNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil

}
