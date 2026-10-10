// Package ecs. This file is generated automatically. Please do not modify it manually, thank you!
package ecs

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSnapshotLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotLockCreate,
		Read:   resourceSnapshotLockRead,
		Update: resourceSnapshotLockUpdate,
		Delete: resourceSnapshotLockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cool_off_period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: helper.IntBetween(0, 72),
			},
			"cool_off_period_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"lock_creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: helper.IntBetween(1, 36500),
			},
			"lock_duration_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: helper.StringInSlice([]string{"compliance"}, false),
			},
			"lock_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_lock_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSnapshotLockCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "LockSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("snapshot_lock_id"); ok {
		request["SnapshotId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["CoolOffPeriod"] = d.Get("cool_off_period")
	request["DryRun"] = d.Get("dry_run")
	request["LockDuration"] = d.Get("lock_duration")
	request["LockMode"] = d.Get("lock_mode")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecs_snapshot_lock", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["SnapshotId"]))

	return resourceSnapshotLockRead(d, meta)
}

func resourceSnapshotLockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	objectRaw, err := NewService(client).DescribeSnapshotLock(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_snapshot_lock DescribeSnapshotLock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cool_off_period", objectRaw["CoolOffPeriod"])
	d.Set("cool_off_period_expired_time", objectRaw["CoolOffPeriodExpiredTime"])
	d.Set("lock_creation_time", objectRaw["LockCreationTime"])
	d.Set("lock_duration", objectRaw["LockDuration"])
	d.Set("lock_duration_start_time", objectRaw["LockDurationStartTime"])
	d.Set("lock_expired_time", objectRaw["LockExpiredTime"])
	d.Set("lock_mode", objectRaw["LockMode"])
	d.Set("lock_status", objectRaw["LockStatus"])
	d.Set("snapshot_lock_id", objectRaw["SnapshotId"])

	return nil
}

func resourceSnapshotLockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "LockSnapshot"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SnapshotId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("cool_off_period") {
		update = true
	}
	request["CoolOffPeriod"] = d.Get("cool_off_period")
	request["DryRun"] = d.Get("dry_run")
	if d.HasChange("lock_duration") {
		update = true
	}
	request["LockDuration"] = d.Get("lock_duration")
	if d.HasChange("lock_mode") {
		update = true
	}
	request["LockMode"] = d.Get("lock_mode")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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

	return resourceSnapshotLockRead(d, meta)
}

func resourceSnapshotLockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "UnlockSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SnapshotId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["DryRun"] = d.Get("dry_run")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		if helper.IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound", "InvalidSnapshotId.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
