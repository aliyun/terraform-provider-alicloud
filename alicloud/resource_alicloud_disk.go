package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskCreate,
		Read:   resourceAliyunDiskRead,
		Update: resourceAliyunDiskUpdate,
		Delete: resourceAliyunDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskDescription,
			},

			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DiskCloud), string(DiskCloudESSD), string(DiskCloudSSD), string(DiskCloudEfficiency)}),
				Default:      DiskCloudEfficiency,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"encrypted"},
			},

			"encrypted": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id"},
			},

			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	availabilityZone, err := ecsService.DescribeZone(d.Get("availability_zone").(string))
	if err != nil {
		return WrapError(err)
	}

	request := ecs.CreateCreateDiskRequest()
	request.RegionId = client.RegionId
	request.ZoneId = availabilityZone.ZoneId

	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		category := DiskCategory(v.(string))
		if err := ecsService.DiskAvailable(availabilityZone, category); err != nil {
			return WrapError(err)
		}
		request.DiskCategory = v.(string)
	}

	request.Size = requests.NewInteger(d.Get("size").(int))

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		request.SnapshotId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.DiskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("encrypted"); ok {
		request.Encrypted = requests.NewBoolean(v.(bool))
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateDisk(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_disk", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateDiskResponse)
	d.SetId(response.DiskId)
	if err := ecsService.WaitForDisk(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunDiskUpdate(d, meta)
}

func resourceAliyunDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeDisk(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("availability_zone", object.ZoneId)
	d.Set("category", object.Category)
	d.Set("size", object.Size)
	d.Set("status", object.Status)
	d.Set("name", object.DiskName)
	d.Set("description", object.Description)
	d.Set("snapshot_id", object.SourceSnapshotId)
	d.Set("encrypted", object.Encrypted)
	d.Set("delete_auto_snapshot", object.DeleteAutoSnapshot)
	d.Set("delete_with_instance", object.DeleteWithInstance)
	d.Set("enable_auto_snapshot", object.EnableAutoSnapshot)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceDisk)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if err := setTags(client, TagResourceDisk, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}
	update := false
	request := ecs.CreateModifyDiskAttributeRequest()
	request.DiskId = d.Id()

	if !d.IsNewResource() && d.HasChange("name") {
		request.DiskName = d.Get("name").(string)
		update = true
		d.SetPartial("name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
		d.SetPartial("description")
	}

	if d.IsNewResource() || d.HasChange("delete_auto_snapshot") {
		v := d.Get("delete_auto_snapshot")
		request.DeleteAutoSnapshot = requests.NewBoolean(v.(bool))
		update = true
		d.SetPartial("delete_auto_snapshot")
	}

	if d.IsNewResource() || d.HasChange("delete_with_instance") {
		v := d.Get("delete_with_instance")
		request.DeleteWithInstance = requests.NewBoolean(v.(bool))
		update = true
		d.SetPartial("delete_with_instance")
	}

	if d.IsNewResource() || d.HasChange("enable_auto_snapshot") {
		v := d.Get("enable_auto_snapshot")
		request.EnableAutoSnapshot = requests.NewBoolean(v.(bool))
		update = true
		d.SetPartial("enable_auto_snapshot")
	}

	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDiskAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunDiskRead(d, meta)
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		request := ecs.CreateResizeDiskRequest()
		request.DiskId = d.Id()
		request.NewSize = requests.NewInteger(size)
		request.Type = string(DiskResizeTypeOnline)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ResizeDisk(request)
		})
		if IsExceptedErrors(err, DiskNotSupportOnlineChangeErrors) {
			request.Type = string(DiskResizeTypeOffline)
			raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ResizeDisk(request)
			})
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("size")
	}

	d.Partial(false)
	return resourceAliyunDiskRead(d, meta)
}

func resourceAliyunDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteDiskRequest()
	request.DiskId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(request)
		})
		if err != nil {
			if IsExceptedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForDisk(d.Id(), Deleted, DefaultTimeout))
}
