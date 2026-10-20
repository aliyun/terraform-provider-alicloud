package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudEcsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSnapshotCreate,
		Read:   resourceAliCloudEcsSnapshotRead,
		Update: resourceAliCloudEcsSnapshotUpdate,
		Delete: resourceAliCloudEcsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"standard", "flash"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instant_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "Field `instant_access` has been deprecated from provider version 1.231.0. ESSD normal snapshots have been upgraded to instant-access by default; no additional configuration is required.",
				Description: "Whether to enable the instant access feature. Valid values: `true` (enabled, ESSD only) and `false` (disabled, creates a normal snapshot). Default value: `false`. Deprecated since v1.231.0: ESSD normal snapshots have been upgraded to instant-access by default, so this field no longer needs to be configured and no extra cost is incurred.",
			},
			"instant_access_retention_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "Field `instant_access_retention_days` has been deprecated from provider version 1.231.0.",
				Description: "The retention period of the instant access feature. When the retention period ends, the feature is disabled and the IA snapshot is automatically released. This parameter takes effect only when `instant_access` is set to `true`. Unit: days. Valid values: 1 to 65535. By default, the value of this parameter is the same as that of `retention_days`.",
			},
			"lock_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The lock duration of the snapshot in compliance mode. Unit: days. Valid values: 1 to 36500 (0 is accepted as an explicit request to unlock; negative values are rejected). When set to a positive value, the snapshot is locked in compliance mode via LockSnapshot with `CoolOffPeriod=0` (immediate compliance, no cooling-off period) and cannot be deleted by any user until the lock duration expires. Once a snapshot enters compliance mode the lock duration can only be extended, not shortened, and must not exceed the snapshot retention period (`retention_days`). A compliance-locked snapshot (CoolOffPeriod=0) CANNOT be unlocked early: explicitly setting `lock_duration` to 0 calls UnlockSnapshot, which returns `InvalidOperation.CoolOffPeriodExpired` for a compliance-locked snapshot; the snapshot only becomes deletable after `lock_duration` elapses (auto-expire). Because `lock_duration` is also Computed, removing the field from the configuration keeps the previously Read-backfilled value (no UnlockSnapshot is issued); to unlock, set `lock_duration` to 0 explicitly.",
				ValidateFunc: validation.IntBetween(0, 36500),
			},
			"lock_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lock status of the snapshot. Valid values: `compliance-cooloff` (locked in compliance mode, still in cooling-off period; can be unlocked, and cooling-off period and lock duration can be extended or shortened), `compliance` (locked in compliance mode, cooling-off period ended; cannot be unlocked or deleted, lock duration can only be extended), and `expired` (lock expired, snapshot is not locked and can be deleted).",
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snapshot_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"wait_until": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"accomplished", "available"}, false),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"snapshot_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.120.0. New field `snapshot_name` instead.",
			},
		},
	}
}

func resourceAliCloudEcsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["SnapshotName"] = v
	}
	if v, ok := d.GetOkExists("retention_days"); ok {
		request["RetentionDays"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["DiskId"] = d.Get("disk_id")
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_snapshot", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateField, target := ecsSnapshotWaitCondition(d)
	stateConf := BuildStateConf([]string{}, []string{target}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsServiceV2.EcsSnapshotStateRefreshFunc(d.Id(), stateField, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if lockDuration := d.Get("lock_duration").(int); lockDuration > 0 {
		// LockSnapshot requires Status=accomplished; if wait_until=available
		// returned while the snapshot was still progressing, wait for
		// accomplished before locking.
		lockStateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsServiceV2.EcsSnapshotStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := lockStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		lockAction := "LockSnapshot"
		lockRequest := make(map[string]interface{})
		lockQuery := make(map[string]interface{})
		lockRequest["SnapshotId"] = d.Id()
		lockRequest["RegionId"] = client.RegionId
		lockRequest["LockMode"] = "compliance"
		lockRequest["LockDuration"] = lockDuration
		lockRequest["CoolOffPeriod"] = 0
		lockRequest["ClientToken"] = buildClientToken(lockAction)

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", lockAction, lockQuery, lockRequest, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.SnapshotNotAvailable"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(lockAction, response, lockRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), lockAction, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEcsSnapshotRead(d, meta)
}

func ecsSnapshotWaitCondition(d *schema.ResourceData) (string, string) {
	if v, ok := d.GetOk("wait_until"); ok && v == "available" {
		return "Available", "true"
	}
	return "Status", "accomplished"
}

func resourceAliCloudEcsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsSnapshot(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_snapshot DescribeEcsSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	available := objectRaw["Available"]
	if available == nil {
		available = false
	}
	if err := d.Set("available", available); err != nil {
		return WrapError(err)
	}
	if objectRaw["Category"] != nil {
		d.Set("category", objectRaw["Category"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["SourceDiskId"] != nil {
		d.Set("disk_id", objectRaw["SourceDiskId"])
	}
	if objectRaw["Encrypted"] != nil {
		d.Set("encrypted", objectRaw["Encrypted"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["RetentionDays"] != nil {
		d.Set("retention_days", objectRaw["RetentionDays"])
	}
	if objectRaw["SnapshotName"] != nil {
		d.Set("snapshot_name", objectRaw["SnapshotName"])
		d.Set("name", objectRaw["SnapshotName"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	lockInfo, lockErr := ecsServiceV2.DescribeEcsSnapshotLock(d.Id())
	if lockErr != nil {
		// A permission, network or service error from DescribeLockedSnapshots
		// must not be ignored. Swallowing it and clearing lock_status/lock_duration
		// to zero here would conflate a real failure with a genuine no-lock snapshot
		// and could mask an actually-locked snapshot. Propagate the error so the
		// caller knows the lock state could not be determined.
		return WrapError(lockErr)
	}
	if lockInfo != nil {
		if lockInfo["LockStatus"] != nil {
			d.Set("lock_status", lockInfo["LockStatus"])
		}
		if lockInfo["LockDuration"] != nil {
			d.Set("lock_duration", lockInfo["LockDuration"])
		}
	} else {
		// DescribeEcsSnapshotLock returns (nil, nil) only for a genuine no-lock
		// snapshot (InvalidSnapshotLock.NotFound or an empty result list). Clear
		// the fields so an unlocked snapshot reads lock_status="" / lock_duration=0.
		d.Set("lock_status", "")
		d.Set("lock_duration", 0)
	}

	return nil
}

func resourceAliCloudEcsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	retentionDaysChanged := false
	d.Partial(true)

	action := "ModifySnapshotAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SnapshotId"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("snapshot_name") {
		update = true
		request["SnapshotName"] = d.Get("snapshot_name")
	}

	if d.HasChange("retention_days") {
		update = true
		retentionDaysChanged = true

		if v, ok := d.GetOkExists("retention_days"); ok {
			request["RetentionDays"] = v
		}
	}

	if d.HasChange("name") {
		update = true

		if v, ok := d.GetOk("name"); ok {
			request["SnapshotName"] = v
		}
	}

	if update {
		if retentionDaysChanged {
			// ECS only accepts retention-day changes after the snapshot upload is complete.
			stateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsServiceV2.EcsSnapshotStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "JoinResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	request["ResourceType"] = "snapshot"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ecsServiceV2 := EcsServiceV2{client}
		if err := ecsServiceV2.SetResourceTags(d, "snapshot"); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("lock_duration") {
		lockDuration := d.Get("lock_duration").(int)
		if lockDuration > 0 {
			// LockSnapshot requires Status=accomplished.
			lockStateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsServiceV2.EcsSnapshotStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := lockStateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

			lockAction := "LockSnapshot"
			lockRequest := make(map[string]interface{})
			lockQuery := make(map[string]interface{})
			lockRequest["SnapshotId"] = d.Id()
			lockRequest["RegionId"] = client.RegionId
			lockRequest["LockMode"] = "compliance"
			lockRequest["LockDuration"] = lockDuration
			lockRequest["CoolOffPeriod"] = 0
			lockRequest["ClientToken"] = buildClientToken(lockAction)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", lockAction, lockQuery, lockRequest, true)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.SnapshotNotAvailable"}) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(lockAction, response, lockRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), lockAction, AlibabaCloudSdkGoERROR)
			}
		} else {
			unlockAction := "UnlockSnapshot"
			unlockRequest := make(map[string]interface{})
			unlockQuery := make(map[string]interface{})
			unlockRequest["SnapshotId"] = d.Id()
			unlockRequest["RegionId"] = client.RegionId
			unlockRequest["ClientToken"] = buildClientToken(unlockAction)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", unlockAction, unlockQuery, unlockRequest, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound"}) {
						return nil
					}
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(unlockAction, response, unlockRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), unlockAction, AlibabaCloudSdkGoERROR)
			}
		}
	}

	d.Partial(false)
	return resourceAliCloudEcsSnapshotRead(d, meta)
}

func resourceAliCloudEcsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SnapshotId"] = d.Id()

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		// InvalidOperation.SnapshotIsLocked means the snapshot is locked in
		// compliance mode and cannot be deleted by any user until the lock
		// duration expires (auto-expire). The error must NOT be swallowed: if
		// Delete returned nil here Terraform would remove the resource from
		// state while the snapshot still exists on the cloud, causing a
		// silent state/cloud divergence. Propagating the real error keeps the
		// resource in state so it is not destroyed prematurely; the snapshot
		// auto-expires after lock_duration. Only snapshots that are already
		// gone (NotFound/InvalidParameter) are treated as successfully deleted.
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound", "InvalidParameter"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
