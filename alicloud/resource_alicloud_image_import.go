package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImageImport() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageImportCreate,
		Read:   resourceAliCloudImageImportRead,
		Update: resourceAliCloudImageImportUpdate,
		Delete: resourceAliCloudImageImportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "x86_64",
				ValidateFunc: StringInSlice([]string{"x86_64", "i386"}, false),
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "linux",
				ValidateFunc: StringInSlice([]string{"windows", "linux"}, false),
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"boot_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"BIOS", "UEFI"}, false),
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Auto",
				ValidateFunc: StringInSlice([]string{"Auto", "Aliyun", "BYOL"}, false),
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_device_mapping": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"RAW", "VHD", "qcow2"}, false),
						},
						"oss_bucket": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"oss_object": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"disk_image_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  5,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudImageImportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client: client}
	request := ecs.CreateImportImageRequest()
	request.RegionId = client.RegionId

	if v, ok := d.GetOk("architecture"); ok {
		request.Architecture = v.(string)
	}

	if v, ok := d.GetOk("os_type"); ok {
		request.OSType = v.(string)
	}

	if v, ok := d.GetOk("platform"); ok {
		request.Platform = v.(string)
	}

	if v, ok := d.GetOk("boot_mode"); ok {
		request.BootMode = v.(string)
	}

	if v, ok := d.GetOk("license_type"); ok {
		request.LicenseType = v.(string)
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	diskDeviceMappings := d.Get("disk_device_mapping")
	diskDeviceMappingsMaps := make([]ecs.ImportImageDiskDeviceMapping, 0)
	for _, diskDeviceMappingsList := range diskDeviceMappings.([]interface{}) {
		var diskDeviceMappingsMap ecs.ImportImageDiskDeviceMapping
		diskDeviceMappingsArg := diskDeviceMappingsList.(map[string]interface{})

		if format, ok := diskDeviceMappingsArg["format"]; ok {
			diskDeviceMappingsMap.Format = format.(string)
		}

		if ossBucket, ok := diskDeviceMappingsArg["oss_bucket"]; ok {
			diskDeviceMappingsMap.OSSBucket = ossBucket.(string)
		}

		if ossObject, ok := diskDeviceMappingsArg["oss_object"]; ok {
			diskDeviceMappingsMap.OSSObject = ossObject.(string)
		}

		if device, ok := diskDeviceMappingsArg["device"]; ok {
			diskDeviceMappingsMap.Device = device.(string)
		}

		if diskImageSize, ok := diskDeviceMappingsArg["disk_image_size"]; ok {
			diskImageSizeStr := strconv.Itoa(diskImageSize.(int))
			diskDeviceMappingsMap.DiskImageSize = diskImageSizeStr
		}

		diskDeviceMappingsMaps = append(diskDeviceMappingsMaps, diskDeviceMappingsMap)
	}

	request.DiskDeviceMapping = &diskDeviceMappingsMaps

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ImportImage(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_import", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	resp, _ := raw.(*ecs.ImportImageResponse)
	d.SetId(resp.ImageId)

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudImageImportRead(d, meta)
}

func resourceAliCloudImageImportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("architecture", object.Architecture)
	d.Set("os_type", object.OSType)
	d.Set("platform", object.Platform)
	d.Set("boot_mode", object.BootMode)
	d.Set("image_name", object.ImageName)
	d.Set("description", object.Description)

	diskDeviceMappings, err := FlattenImageImportDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping)
	if err != nil {
		return WrapError(err)
	}

	d.Set("disk_device_mapping", diskDeviceMappings)

	return nil
}

func resourceAliCloudImageImportUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}

	return resourceAliCloudImageRead(d, meta)
}

func resourceAliCloudImageImportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	return ecsService.deleteImage(d)
}

func FlattenImageImportDiskDeviceMappings(list []ecs.DiskDeviceMapping) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		diskImageSize, err := strconv.Atoi(v.Size)
		if err != nil {
			return nil, WrapError(err)
		}

		diskDeviceMappings := map[string]interface{}{
			"format":          v.Format,
			"oss_bucket":      v.ImportOSSBucket,
			"oss_object":      v.ImportOSSObject,
			"device":          v.Device,
			"disk_image_size": diskImageSize,
		}

		result = append(result, diskDeviceMappings)
	}

	return result, nil
}
