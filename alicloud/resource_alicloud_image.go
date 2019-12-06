package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliCloudImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageCreate,
		Read:   resourceAliCloudImageRead,
		Update: resourceAliCloudImageUpdate,
		Delete: resourceAliCloudImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "x86_64",
				ValidateFunc: validation.StringInSlice([]string{
					"x86_64",
					"i386",
				}, false),
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"disk_device_mapping", "snapshot_id"},
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id", "disk_device_mapping"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Ubuntu",
				ValidateFunc: validation.StringInSlice([]string{
					"CentOS",
					"Ubuntu",
					"SUSE",
					"OpenSUSE",
					"RedHat",
					"Debian",
					"CoreOS",
					"Aliyun Linux",
					"Windows Server 2003",
					"Windows Server 2008",
					"Windows Server 2012",
					"Windows 7",
					"Customized Linux",
					"Others Linux",
				}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_device_mapping": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_id", "snapshot_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),
		},
	}
}
func resourceAliCloudImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	// Make sure the instance status is Running or Stopped
	if v, ok := d.GetOk("instance_id"); ok {
		instance, err := ecsService.DescribeInstance(v.(string))
		if err != nil {
			return WrapError(err)
		}
		status := Status(instance.Status)
		if status != Running && status != Stopped {
			return WrapError(Error("You must make sure that the status of the specified instance is Running or Stopped. "))
		}
	}

	// The snapshot cannot be a snapshot created before July 15, 2013 (inclusive)
	if snapshotId, ok := d.GetOk("snapshot_id"); ok {
		snapshot, err := ecsService.DescribeSnapshot(snapshotId.(string))
		if err != nil {
			return WrapError(err)
		}
		snapshotCreationTime, err := time.Parse("2006-01-02T15:04:05Z", snapshot.CreationTime)
		if err != nil {
			return WrapErrorf(err, IdMsg, snapshotId)
		}
		deadlineTime, _ := time.Parse("2006-01-02T15:04:05Z", "2013-07-16T00:00:00Z")
		if deadlineTime.After(snapshotCreationTime) {
			return WrapError(Error("the specified snapshot cannot be created on or before July 15, 2013."))
		}
	}
	request := ecs.CreateCreateImageRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Get("instance_id").(string)
	diskDeviceMappings := d.Get("disk_device_mapping").([]interface{})
	if diskDeviceMappings != nil && len(diskDeviceMappings) > 0 {
		mappings := make([]ecs.CreateImageDiskDeviceMapping, 0, len(diskDeviceMappings))
		for _, diskDeviceMapping := range diskDeviceMappings {
			mapping := diskDeviceMapping.(map[string]interface{})
			deviceMapping := ecs.CreateImageDiskDeviceMapping{
				SnapshotId: mapping["snapshot_id"].(string),
				Size:       mapping["size"].(string),
				DiskType:   mapping["disk_type"].(string),
				Device:     mapping["device"].(string),
			}
			mappings = append(mappings, deviceMapping)
		}
		request.DiskDeviceMapping = &mappings
	}
	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		imageTags := make([]ecs.CreateImageTag, 0, len(tags))
		for k, v := range tags {
			imageTag := ecs.CreateImageTag{
				Key:   k,
				Value: v.(string),
			}
			imageTags = append(imageTags, imageTag)
		}
		request.Tag = &imageTags
	}
	request.SnapshotId = d.Get("snapshot_id").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.Platform = d.Get("platform").(string)
	request.ImageName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.Architecture = d.Get("architecture").(string)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.CreateImageResponse)
	d.SetId(response.ImageId)
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	err := setTags(client, TagResourceImage, d)
	if err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	request := ecs.CreateModifyImageAttributeRequest()
	request.RegionId = client.RegionId
	request.ImageId = d.Id()

	if d.HasChange("description") || d.HasChange("name") {
		if description, ok := d.GetOk("description"); ok {
			request.Description = description.(string)
		}
		if imageName, ok := d.GetOk("name"); ok {
			request.ImageName = imageName.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyImageAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("description")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("platform", object.Platform)
	d.Set("name", object.ImageName)
	d.Set("description", object.Description)
	d.Set("architecture", object.Architecture)
	d.Set("disk_device_mapping", FlattenImageDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping))
	tags := object.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", tagsToMap(tags))
	}
	return WrapError(err)
}

func resourceAliCloudImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()

	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.RegionId = client.RegionId
	request.ImageId = object.ImageId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, ecsService.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func FlattenImageDiskDeviceMappings(list []ecs.DiskDeviceMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		size, _ := strconv.Atoi(i.Size)
		l := map[string]interface{}{
			"device":      i.Device,
			"size":        size,
			"snapshot_id": i.SnapshotId,
			"disk_type":   i.Type,
		}
		result = append(result, l)
	}

	return result
}
