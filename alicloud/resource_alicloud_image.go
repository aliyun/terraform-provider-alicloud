// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsImageCreate,
		Read:   resourceAliCloudEcsImageRead,
		Update: resourceAliCloudEcsImageUpdate,
		Delete: resourceAliCloudEcsImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "x86_64",
				ValidateFunc: StringInSlice([]string{"i386", "x86_64", "arm64"}, false),
			},
			"boot_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detection_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_device_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"system", "data"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: IntAtMost(32768),
						},
						"import_oss_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remain_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"import_oss_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"features": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nvme_support": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"image_family": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"image_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Aliyun", "Anolis", "CentOS", "Ubuntu", "CoreOS", "SUSE", "Debian", "OpenSUSE", "FreeBSD", "RedHat", "Kylin", "UOS", "Fedora", "Fedora CoreOS", "CentOS Stream", "AlmaLinux", "Rocky Linux", "Gentoo", "Customized Linux", "Others Linux", "Windows Server 2022", "Windows Server 2019", "Windows Server 2016", "Windows Server 2012", "Windows Server 2008", "Windows Server 2003"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			// Not the public attribute and it used to automatically delete dependence snapshots while deleting the image.
			// Available in 1.136.0
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.227.0. New field 'image_name' instead.",
			},
		},
	}
}

func resourceAliCloudEcsImageCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateImage"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("architecture"); ok {
		request["Architecture"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["ImageName"] = v
	}

	if v, ok := d.GetOk("image_name"); ok {
		request["ImageName"] = v
	}
	if v, ok := d.GetOk("image_version"); ok {
		request["ImageVersion"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOk("image_family"); ok {
		request["ImageFamily"] = v
	}
	if v, ok := d.GetOk("boot_mode"); ok {
		request["BootMode"] = v
	}
	if v, ok := d.GetOk("detection_strategy"); ok {
		request["DetectionStrategy"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("disk_device_mapping"); ok {
		diskDeviceMappingMaps := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["SnapshotId"] = dataLoop1Tmp["snapshot_id"]
			dataLoop1Map["Size"] = dataLoop1Tmp["size"]
			dataLoop1Map["Device"] = dataLoop1Tmp["device"]
			dataLoop1Map["DiskType"] = dataLoop1Tmp["disk_type"]
			diskDeviceMappingMaps = append(diskDeviceMappingMaps, dataLoop1Map)
		}
		request["DiskDeviceMapping"] = diskDeviceMappingMaps
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ImageId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsServiceV2.EcsImageStateRefreshFunc(d.Id(), "Status", []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsImageRead(d, meta)
}

func resourceAliCloudEcsImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsImage(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_image DescribeEcsImage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Architecture"] != nil {
		d.Set("architecture", objectRaw["Architecture"])
	}
	if objectRaw["BootMode"] != nil {
		d.Set("boot_mode", objectRaw["BootMode"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["ImageFamily"] != nil {
		d.Set("image_family", objectRaw["ImageFamily"])
	}
	if objectRaw["ImageName"] != nil {
		d.Set("image_name", objectRaw["ImageName"])
	}
	if objectRaw["ImageVersion"] != nil {
		d.Set("image_version", objectRaw["ImageVersion"])
	}
	if objectRaw["Platform"] != nil {
		d.Set("platform", objectRaw["Platform"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	diskDeviceMapping1Raw, _ := jsonpath.Get("$.DiskDeviceMappings.DiskDeviceMapping", objectRaw)
	diskDeviceMappingMaps := make([]map[string]interface{}, 0)
	if diskDeviceMapping1Raw != nil {
		for _, diskDeviceMappingChild1Raw := range diskDeviceMapping1Raw.([]interface{}) {
			diskDeviceMappingMap := make(map[string]interface{})
			diskDeviceMappingChild1Raw := diskDeviceMappingChild1Raw.(map[string]interface{})
			diskDeviceMappingMap["device"] = diskDeviceMappingChild1Raw["Device"]
			diskDeviceMappingMap["disk_type"] = diskDeviceMappingChild1Raw["Type"]
			diskDeviceMappingMap["format"] = diskDeviceMappingChild1Raw["Format"]
			diskDeviceMappingMap["import_oss_object"] = diskDeviceMappingChild1Raw["ImportOSSObject"]
			diskDeviceMappingMap["import_oss_bucket"] = diskDeviceMappingChild1Raw["ImportOSSBucket"]
			diskDeviceMappingMap["progress"] = diskDeviceMappingChild1Raw["Progress"]
			diskDeviceMappingMap["remain_time"] = diskDeviceMappingChild1Raw["RemainTime"]
			diskDeviceMappingMap["size"] = diskDeviceMappingChild1Raw["Size"]
			diskDeviceMappingMap["snapshot_id"] = diskDeviceMappingChild1Raw["SnapshotId"]

			diskDeviceMappingMaps = append(diskDeviceMappingMaps, diskDeviceMappingMap)
		}
	}
	d.Set("disk_device_mapping", diskDeviceMappingMaps)
	featuresMaps := make([]map[string]interface{}, 0)
	featuresMap := make(map[string]interface{})
	features1Raw := make(map[string]interface{})
	if objectRaw["Features"] != nil {
		features1Raw = objectRaw["Features"].(map[string]interface{})
	}
	if len(features1Raw) > 0 {
		featuresMap["nvme_support"] = features1Raw["NvmeSupport"]

		featuresMaps = append(featuresMaps, featuresMap)
	}
	if objectRaw["Features"] != nil {
		d.Set("features", featuresMaps)
	}
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("image_name"))
	return nil
}

func resourceAliCloudEcsImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyImageAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ImageId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("name") {
		update = true
		request["ImageName"] = d.Get("name")
	}

	if d.HasChange("image_name") {
		update = true
		request["ImageName"] = d.Get("image_name")
	}

	if d.HasChange("image_family") {
		update = true
		request["ImageFamily"] = d.Get("image_family")
	}

	if d.HasChange("boot_mode") {
		update = true
		request["BootMode"] = d.Get("boot_mode")
	}

	if v, ok := d.GetOk("license_type"); ok {
		request["LicenseType"] = v
	}
	if v, ok := d.GetOk("features"); ok {
		jsonPathResult5, err := jsonpath.Get("$[0].nvme_support", v)
		if err == nil && jsonPathResult5 != "" {
			request["Features.NvmeSupport"] = jsonPathResult5
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "JoinResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "image"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ecsServiceV2 := EcsServiceV2{client}
		if err := ecsServiceV2.SetResourceTags(d, "image"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{d.Get("description").(string)}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsServiceV2.EcsImageStateRefreshFunc(d.Id(), "Description", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudEcsImageRead(d, meta)
}

func resourceAliCloudEcsImageDelete(d *schema.ResourceData, meta interface{}) error {

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
	action := "DeleteImage"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	query["ImageId"] = d.Id()
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsServiceV2.EcsImageStateRefreshFunc(d.Id(), "ImageId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if v, ok := d.GetOk("delete_auto_snapshot"); ok && v.(bool) {
		errs := map[string]error{}

		for _, item := range object.DiskDeviceMappings.DiskDeviceMapping {
			if item.SnapshotId == "" {
				continue
			}
			request := ecs.CreateDeleteSnapshotRequest()
			request.RegionId = ecsService.client.RegionId
			request.SnapshotId = item.SnapshotId
			request.Force = requests.NewBoolean(true)
			raw, err := ecsService.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DeleteSnapshot(request)
			})
			if err != nil {
				errs[item.SnapshotId] = err
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
		if len(errs) > 0 {
			errParts := []string{"Errors while deleting associated snapshots:"}
			for snapshotId, err := range errs {
				errParts = append(errParts, fmt.Sprintf("%s: %s", snapshotId, err))
			}
			errParts = append(errParts, "These are no longer managed by Terraform and must be deleted manually.")
			return WrapError(fmt.Errorf(strings.Join(errParts, "\n")))
		}
	}
	return nil
}

func resourceAliCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_image ecsService.DescribeImageById Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("platform", object.Platform)
	d.Set("image_name", object.ImageName)
	d.Set("name", object.ImageName)
	d.Set("description", object.Description)
	d.Set("architecture", object.Architecture)
	d.Set("disk_device_mapping", FlattenImageDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping))
	tags, err := ecsService.ListTagResources(d.Id(), "image")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}
	return WrapError(err)
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
