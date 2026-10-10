// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsSnapshotLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSnapshotLockCreate,
		Read:   resourceAliCloudEcsSnapshotLockRead,
		Update: resourceAliCloudEcsSnapshotLockUpdate,
		Delete: resourceAliCloudEcsSnapshotLockDelete,
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
				ValidateFunc: IntBetween(0, 72),
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
				ValidateFunc: IntBetween(1, 36500),
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
				ValidateFunc: StringInSlice([]string{"compliance"}, false),
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

func resourceAliCloudEcsSnapshotLockCreate(d *schema.ResourceData, meta interface{}) error {

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
	request["ClientToken"] = buildClientToken(action)

	request["CoolOffPeriod"] = d.Get("cool_off_period")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["LockDuration"] = d.Get("lock_duration")
	request["LockMode"] = d.Get("lock_mode")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_snapshot_lock", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["SnapshotId"]))

	return resourceAliCloudEcsSnapshotLockRead(d, meta)
}

func resourceAliCloudEcsSnapshotLockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsSnapshotLock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_snapshot_lock DescribeEcsSnapshotLock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

	d.Set("snapshot_lock_id", d.Id())

	return nil
}

func resourceAliCloudEcsSnapshotLockUpdate(d *schema.ResourceData, meta interface{}) error {
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
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("cool_off_period") {
		update = true
	}
	request["CoolOffPeriod"] = d.Get("cool_off_period")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("lock_duration") {
		update = true
	}
	request["LockDuration"] = d.Get("lock_duration")
	if d.HasChange("lock_mode") {
		update = true
	}
	request["LockMode"] = d.Get("lock_mode")
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEcsSnapshotLockRead(d, meta)
}

func resourceAliCloudEcsSnapshotLockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "UnlockSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SnapshotId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound", "InvalidSnapshotId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
