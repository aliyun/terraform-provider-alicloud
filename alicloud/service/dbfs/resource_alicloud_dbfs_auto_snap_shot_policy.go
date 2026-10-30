package dbfs

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAutoSnapShotPolicy() *schema.Resource {
	return &schema.Resource{
		Create:             resourceAlicloudDbfsAutoSnapShotPolicyCreate,
		Read:               resourceAlicloudDbfsAutoSnapShotPolicyRead,
		Update:             resourceAlicloudDbfsAutoSnapShotPolicyUpdate,
		Delete:             resourceAlicloudDbfsAutoSnapShotPolicyDelete,
		DeprecationMessage: "This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"applied_dbfs_number": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"last_modified": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"policy_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"policy_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"repeat_weekdays": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"1", "2", "3", "4", "5", "6", "7"}, false),
				},
			},
			"retention_days": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status_detail": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"time_points": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}, false),
				},
			},
		},
	}
}

func resourceAlicloudDbfsAutoSnapShotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["PolicyName"] = d.Get("policy_name")
	request["RetentionDays"] = d.Get("retention_days")
	request["RepeatWeekdays"] = convert.ConvertListToJsonString(d.Get("repeat_weekdays").([]interface{}))
	request["TimePoints"] = convert.ConvertListToJsonString(d.Get("time_points").([]interface{}))

	var response map[string]interface{}
	action := "CreateAutoSnapshotPolicy"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("DBFS", "2020-04-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dbfs_auto_snap_shot_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PolicyId", response); err != nil || v == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_dbfs_auto_snap_shot_policy")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudDbfsAutoSnapShotPolicyUpdate(d, meta)
}

func resourceAlicloudDbfsAutoSnapShotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := Service{client}

	object, err := dbfsService.DescribeDbfsAutoSnapShotPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_auto_snap_shot_policy dbfsService.DescribeDbfsAutoSnapShotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("applied_dbfs_number", object["AppliedDbfsNumber"])
	d.Set("create_time", object["CreatedTime"])
	d.Set("last_modified", object["LastModified"])
	d.Set("policy_name", object["PolicyName"])
	repeatWeekdays, _ := jsonpath.Get("$.RepeatWeekdays", object)
	d.Set("repeat_weekdays", repeatWeekdays)
	d.Set("retention_days", object["RetentionDays"])
	d.Set("status", object["Status"])
	d.Set("status_detail", object["StatusDetail"])
	timePoints, _ := jsonpath.Get("$.TimePoints", object)
	d.Set("time_points", timePoints)

	return nil
}

func resourceAlicloudDbfsAutoSnapShotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"PolicyId": d.Id(),
		"RegionId": client.RegionId,
	}

	if d.HasChange("policy_name") {
		update = true
		if v, ok := d.GetOk("policy_name"); ok {
			request["PolicyName"] = v
		}
	}
	if d.HasChange("repeat_weekdays") {
		update = true
		if v, ok := d.GetOk("repeat_weekdays"); ok {
			request["RepeatWeekdays"] = convert.ConvertListToJsonString(v.([]interface{}))
		}
	}
	if d.HasChange("retention_days") {
		update = true
		if v, ok := d.GetOk("retention_days"); ok {
			request["RetentionDays"] = v
		}
	}
	if d.HasChange("time_points") {
		update = true
		if v, ok := d.GetOk("time_points"); ok {
			request["TimePoints"] = convert.ConvertListToJsonString(v.([]interface{}))
		}
	}

	if update {
		action := "ModifyAutoSnapshotPolicy"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcPost("DBFS", "2020-04-18", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudDbfsAutoSnapShotPolicyRead(d, meta)
}

func resourceAlicloudDbfsAutoSnapShotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"PolicyId": d.Id(),
	}

	action := "DeleteAutoSnapshotPolicy"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		resp, err := client.RpcPost("DBFS", "2020-04-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AutoSnapshotPolicyNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
