package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcdBundle() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdBundleCreate,
		Read:   resourceAlicloudEcdBundleRead,
		Update: resourceAlicloudEcdBundleUpdate,
		Delete: resourceAlicloudEcdBundleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desktop_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_disk_size_gib": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"user_disk_size_gib": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"root_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"user_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh-CN", "zh-HK", "en-US", "ja-JP"}, false),
			},
			"bundle_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcdBundleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBundle"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ImageId"] = d.Get("image_id")
	request["DesktopType"] = d.Get("desktop_type")
	request["RootDiskSizeGib"] = d.Get("root_disk_size_gib")
	request["UserDiskSizeGib"] = d.Get("user_disk_size_gib")

	if v, ok := d.GetOk("root_disk_performance_level"); ok {
		request["RootDiskPerformanceLevel"] = v
	}

	if v, ok := d.GetOk("user_disk_performance_level"); ok {
		request["UserDiskPerformanceLevel"] = v
	}

	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}

	if v, ok := d.GetOk("bundle_name"); ok {
		request["BundleName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_bundle", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BundleId"]))

	return resourceAlicloudEcdBundleRead(d, meta)
}

func resourceAlicloudEcdBundleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}

	object, err := ecdService.DescribeEcdBundle(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_bundle ecdService.DescribeEcdBundle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("image_id", object["ImageId"])
	d.Set("desktop_type", object["DesktopType"])
	d.Set("language", object["Language"])
	d.Set("bundle_name", object["BundleName"])
	d.Set("description", object["Description"])

	diskArr := make([]int, 0)
	for _, item := range object["Disks"].([]interface{}) {
		disk := item.(map[string]interface{})
		if fmt.Sprint(disk["DiskType"]) == "SYSTEM" {
			d.Set("root_disk_size_gib", formatInt(disk["DiskSize"]))
			d.Set("root_disk_performance_level", disk["DiskPerformanceLevel"])
		}

		if fmt.Sprint(disk["DiskType"]) == "DATA" {
			diskArr = append(diskArr, formatInt(disk["DiskSize"]))
			d.Set("user_disk_performance_level", disk["DiskPerformanceLevel"])
		}
	}

	d.Set("user_disk_size_gib", diskArr)

	return nil
}

func resourceAlicloudEcdBundleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"BundleId": d.Id(),
	}

	if d.HasChange("image_id") {
		update = true
	}
	request["ImageId"] = d.Get("image_id")

	if d.HasChange("language") {
		update = true
		if v, ok := d.GetOk("language"); ok {
			request["Language"] = v
		}
	}

	if d.HasChange("bundle_name") {
		update = true
		if v, ok := d.GetOk("bundle_name"); ok {
			request["BundleName"] = v
		}
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyBundle"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudEcdBundleRead(d, meta)
}

func resourceAlicloudEcdBundleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBundles"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"BundleId": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
