package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bundle_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh-CN", "zh-HK", "en-US", "ja-JP"}, false),
			},
			"root_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"root_disk_size_gib": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"user_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"user_disk_size_gib": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcdBundleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBundle"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("bundle_name"); ok {
		request["BundleName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DesktopType"] = d.Get("desktop_type")
	request["ImageId"] = d.Get("image_id")
	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("root_disk_performance_level"); ok {
		request["RootDiskPerformanceLevel"] = v
	}
	request["RootDiskSizeGib"] = d.Get("root_disk_size_gib")
	if v, ok := d.GetOk("user_disk_performance_level"); ok {
		request["UserDiskPerformanceLevel"] = v
	}
	request["UserDiskSizeGib"] = d.Get("user_disk_size_gib")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_bundle ecdService.DescribeEcdBundle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bundle_name", object["BundleName"])
	d.Set("description", object["Description"])
	d.Set("desktop_type", object["DesktopType"])
	d.Set("image_id", object["ImageId"])
	diskArr := make([]int, 0)
	for _, item := range object["Disks"].([]interface{}) {
		disk := item.(map[string]interface{})
		if fmt.Sprint(disk["DiskType"]) == "DATA" {
			diskArr = append(diskArr, formatInt(disk["DiskSize"]))
			d.Set("user_disk_performance_level", disk["DiskPerformanceLevel"])
		}
		if fmt.Sprint(disk["DiskType"]) == "SYSTEM" {
			d.Set("root_disk_size_gib", formatInt(disk["DiskSize"]))
			d.Set("root_disk_performance_level", disk["DiskPerformanceLevel"])
		}
	}
	d.Set("user_disk_size_gib", diskArr)
	return nil
}
func resourceAlicloudEcdBundleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"BundleId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("bundle_name") {
		update = true
		if v, ok := d.GetOk("bundle_name"); ok {
			request["BundleName"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("image_id") {
		update = true
		request["ImageId"] = d.Get("image_id")
	}
	if update {
		if v, ok := d.GetOk("language"); ok {
			request["Language"] = v
		}
		action := "ModifyBundle"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BundleId": []string{d.Id()},
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
