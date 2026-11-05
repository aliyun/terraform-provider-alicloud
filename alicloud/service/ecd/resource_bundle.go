package ecd

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBundle() *schema.Resource {
	return &schema.Resource{
		Create: resourceBundleCreate,
		Read:   resourceBundleRead,
		Update: resourceBundleUpdate,
		Delete: resourceBundleDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"user_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"zh-CN", "zh-HK", "en-US", "ja-JP"}, false),
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

func resourceBundleCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_bundle", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BundleId"]))

	return resourceBundleRead(d, meta)
}

func resourceBundleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)

	object, err := ecdService.DescribeEcdBundle(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_bundle ecdService.DescribeEcdBundle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
			root_disk_size_gibVal, err := convert.FormatInt(disk["DiskSize"])
			if err != nil {
				return err
			}
			d.Set("root_disk_size_gib", root_disk_size_gibVal)
			d.Set("root_disk_performance_level", disk["DiskPerformanceLevel"])
		}

		if fmt.Sprint(disk["DiskType"]) == "DATA" {
			diskArrItem, err := convert.FormatInt(disk["DiskSize"])
			if err != nil {
				return err
			}
			diskArr = append(diskArr, diskArrItem)
			d.Set("user_disk_performance_level", disk["DiskPerformanceLevel"])
		}
	}

	d.Set("user_disk_size_gib", diskArr)

	return nil
}

func resourceBundleUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceBundleRead(d, meta)
}

func resourceBundleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBundles"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"BundleId": []string{d.Id()},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
