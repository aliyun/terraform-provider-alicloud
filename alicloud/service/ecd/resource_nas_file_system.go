package ecd

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNasFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceNasFileSystemCreate,
		Read:   resourceNasFileSystemRead,
		Update: resourceNasFileSystemUpdate,
		Delete: resourceNasFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mount_target_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nas_file_system_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reset": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNasFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNASFileSystem"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("nas_file_system_name"); ok {
		request["Name"] = v
	}
	request["OfficeSiteId"] = d.Get("office_site_id")
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_nas_file_system", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))
	ecdService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdNasFileSystemStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceNasFileSystemUpdate(d, meta)
}
func resourceNasFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdNasFileSystem(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_nas_file_system ecdService.DescribeEcdNasFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("mount_target_domain", object["MountTargetDomain"])
	d.Set("nas_file_system_name", object["FileSystemName"])
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("status", object["FileSystemStatus"])
	d.Set("file_system_id", object["FileSystemId"])
	return nil
}
func resourceNasFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
		"RegionId":     client.RegionId,
	}
	if d.HasChange("reset") {
		if v, ok := d.GetOk("reset"); ok && v.(bool) {
			update = true
		}
	}
	if update {
		action := "ResetNASDefaultMountTarget"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecdService.EcdNasFileSystemStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceNasFileSystemRead(d, meta)
}
func resourceNasFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	action := "DeleteNASFileSystems"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"FileSystemId": []string{d.Id()},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNasFileSystemStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
