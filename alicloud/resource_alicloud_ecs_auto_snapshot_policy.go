// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudEcsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsAutoSnapshotPolicyCreate,
		Read:   resourceAliCloudEcsAutoSnapshotPolicyRead,
		Update: resourceAliCloudEcsAutoSnapshotPolicyUpdate,
		Delete: resourceAliCloudEcsAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"copied_snapshots_retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_cross_region_copy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"repeat_weekdays": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"target_copy_regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"time_points": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudEcsAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateAutoSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["retentionDays"] = d.Get("retention_days")
	jsonPathResult1, err := jsonpath.Get("$", d.Get("repeat_weekdays"))
	if err != nil {
		return WrapError(err)
	}
	request["repeatWeekdays"] = jsonPathResult1

	jsonPathResult2, err := jsonpath.Get("$", d.Get("time_points"))
	if err != nil {
		return WrapError(err)
	}
	request["timePoints"] = jsonPathResult2

	if v, ok := d.GetOk("copied_snapshots_retention_days"); ok {
		request["CopiedSnapshotsRetentionDays"] = v
	}
	if v, ok := d.GetOkExists("enable_cross_region_copy"); ok {
		request["EnableCrossRegionCopy"] = v
	}
	if v, ok := d.GetOk("target_copy_regions"); ok {
		jsonPathResult5, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		request["TargetCopyRegions"] = jsonPathResult5
	}
	if v, ok := d.GetOk("auto_snapshot_policy_name"); ok {
		request["autoSnapshotPolicyName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AutoSnapshotPolicyId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsServiceV2.EcsAutoSnapshotPolicyStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsAutoSnapshotPolicyUpdate(d, meta)
}

func resourceAliCloudEcsAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsAutoSnapshotPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_auto_snapshot_policy DescribeEcsAutoSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_snapshot_policy_name", objectRaw["AutoSnapshotPolicyName"])
	d.Set("copied_snapshots_retention_days", objectRaw["CopiedSnapshotsRetentionDays"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("enable_cross_region_copy", objectRaw["EnableCrossRegionCopy"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("retention_days", objectRaw["RetentionDays"])
	d.Set("status", objectRaw["Status"])
	repeatWeekdays1Raw := make(map[string]interface{})
	if objectRaw["RepeatWeekdays"] != nil {
		repeatWeekdays1Raw = objectRaw["RepeatWeekdays"].(map[string]interface{})
	}
	if len(repeatWeekdays1Raw) > 0 {
		d.Set("repeat_weekdays", repeatWeekdays1Raw)
		tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
		d.Set("tags", tagsToMap(tagsMaps))
		targetCopyRegions1Raw := make(map[string]interface{})
		if objectRaw["TargetCopyRegions"] != nil {
			targetCopyRegions1Raw = objectRaw["TargetCopyRegions"].(map[string]interface{})
		}
		if len(targetCopyRegions1Raw) > 0 {
			d.Set("target_copy_regions", targetCopyRegions1Raw)
			timePoints1Raw := make(map[string]interface{})
			if objectRaw["TimePoints"] != nil {
				timePoints1Raw = objectRaw["TimePoints"].(map[string]interface{})
			}
			if len(timePoints1Raw) > 0 {
				d.Set("time_points", timePoints1Raw)

				return nil
			}
		}
	}
	return nil
}
func resourceAliCloudEcsAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyAutoSnapshotPolicyEx"
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["autoSnapshotPolicyId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("retention_days") {
		update = true
	}
	request["retentionDays"] = d.Get("retention_days")
	if !d.IsNewResource() && d.HasChange("repeat_weekdays") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("repeat_weekdays"))
	if err != nil {
		return WrapError(err)
	}
	request["repeatWeekdays"] = jsonPathResult1

	if !d.IsNewResource() && d.HasChange("time_points") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$", d.Get("time_points"))
	if err != nil {
		return WrapError(err)
	}
	request["timePoints"] = jsonPathResult2

	if !d.IsNewResource() && d.HasChange("copied_snapshots_retention_days") {
		update = true
		request["CopiedSnapshotsRetentionDays"] = d.Get("copied_snapshots_retention_days")
	}

	if !d.IsNewResource() && d.HasChange("enable_cross_region_copy") {
		update = true
		request["EnableCrossRegionCopy"] = d.Get("enable_cross_region_copy")
	}

	if !d.IsNewResource() && d.HasChange("target_copy_regions") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$", d.Get("target_copy_regions"))
		if err != nil {
			return WrapError(err)
		}
		request["TargetCopyRegions"] = jsonPathResult5
	}

	if !d.IsNewResource() && d.HasChange("auto_snapshot_policy_name") {
		update = true
		request["autoSnapshotPolicyName"] = d.Get("auto_snapshot_policy_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		if err := ecsServiceV2.SetResourceTags(d, "snapshotpolicy"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	return resourceAliCloudEcsAutoSnapshotPolicyRead(d, meta)
}

func resourceAliCloudEcsAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAutoSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["autoSnapshotPolicyId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterInvalid"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
