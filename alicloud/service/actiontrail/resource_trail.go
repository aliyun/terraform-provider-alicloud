// Package actiontrail. This file is generated automatically. Please do not modify it manually, thank you!
package actiontrail

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTrail() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTrailCreate,
		Read:   ResourceTrailRead,
		Update: ResourceTrailUpdate,
		Delete: ResourceTrailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_event_trail_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_rw": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Write", "Read", "All"}, false),
			},
			"event_selectors": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_organization_trail": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"max_compute_project_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_compute_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_project_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trail_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"trail_name", "name"},
				Computed:     true,
				ForceNew:     true,
			},
			"trail_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.95.0. New field 'trail_name' instead.",
				ForceNew:   true,
			},
			"mns_topic_arn": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `mns_topic_arn` has been deprecated from version 1.118.0",
			},
			"role_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `role_name` has been deprecated from version 1.118.0",
			},
		},
	}
}

func ResourceTrailCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrail"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("trail_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "trail_name" or "name" must be set one!`))
	}

	if v, ok := d.GetOk("sls_project_arn"); ok {
		request["SlsProjectArn"] = v
	}
	if v, ok := d.GetOkExists("is_organization_trail"); ok {
		request["IsOrganizationTrail"] = v
	}
	if v, ok := d.GetOk("max_compute_project_arn"); ok {
		request["MaxComputeProjectArn"] = v
	}
	if v, ok := d.GetOk("event_rw"); ok {
		request["EventRW"] = v
	}
	if v, ok := d.GetOk("oss_write_role_arn"); ok {
		request["OssWriteRoleArn"] = v
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok {
		request["OssBucketName"] = v
	}
	if v, ok := d.GetOk("oss_key_prefix"); ok {
		request["OssKeyPrefix"] = v
	}
	if v, ok := d.GetOk("sls_write_role_arn"); ok {
		request["SlsWriteRoleArn"] = v
	}
	if v, ok := d.GetOk("trail_region"); ok {
		request["TrailRegion"] = v
	}
	if v, ok := d.GetOk("max_compute_write_role_arn"); ok {
		request["MaxComputeWriteRoleArn"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_actiontrail_trail", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	actionTrailServiceV2 := NewActionTrailServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Fresh"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, actionTrailServiceV2.ActionTrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceTrailUpdate(d, meta)
}

func ResourceTrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actionTrailServiceV2 := NewActionTrailServiceV2(client)

	objectRaw, err := actionTrailServiceV2.DescribeActionTrailTrail(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_trail DescribeActionTrailTrail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("event_rw", objectRaw["EventRW"])
	d.Set("is_organization_trail", objectRaw["IsOrganizationTrail"])
	d.Set("max_compute_project_arn", objectRaw["MaxComputeProjectArn"])
	d.Set("max_compute_write_role_arn", objectRaw["MaxComputeWriteRoleArn"])
	d.Set("oss_bucket_name", objectRaw["OssBucketName"])
	d.Set("oss_key_prefix", objectRaw["OssKeyPrefix"])
	d.Set("oss_write_role_arn", objectRaw["OssWriteRoleArn"])
	d.Set("region_id", objectRaw["HomeRegion"])
	d.Set("sls_project_arn", objectRaw["SlsProjectArn"])
	d.Set("sls_write_role_arn", objectRaw["SlsWriteRoleArn"])
	d.Set("status", objectRaw["Status"])
	d.Set("trail_region", objectRaw["TrailRegion"])
	d.Set("trail_name", objectRaw["Name"])

	objectRaw, err = actionTrailServiceV2.DescribeTrailGetDataEventSelector(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("event_selectors", objectRaw["DataEventSelectors"])

	slsDeliveryConfigsRawArrayObj, _ := jsonpath.Get("$.SlsDeliveryConfigs[*]", objectRaw)
	slsDeliveryConfigsRawArray := make([]interface{}, 0)
	if slsDeliveryConfigsRawArrayObj != nil {
		slsDeliveryConfigsRawArray = convert.ConvertToInterfaceArray(slsDeliveryConfigsRawArrayObj)
	}
	slsDeliveryConfigsRaw := make(map[string]interface{})
	if len(slsDeliveryConfigsRawArray) > 0 {
		slsDeliveryConfigsRaw = slsDeliveryConfigsRawArray[0].(map[string]interface{})
	}

	d.Set("data_event_trail_region", slsDeliveryConfigsRaw["TrailRegion"])

	d.Set("trail_name", d.Id())

	d.Set("name", d.Get("trail_name"))
	return nil
}

func ResourceTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	actionTrailServiceV2 := NewActionTrailServiceV2(client)
	objectRaw, _ := actionTrailServiceV2.DescribeActionTrailTrail(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "Enable" {
				action := "StartLogging"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Name"] = d.Id()

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || helper.NeedRetry(err) {
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
				actionTrailServiceV2 := NewActionTrailServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Enable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actionTrailServiceV2.ActionTrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "Disable" {
				action := "StopLogging"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Name"] = d.Id()

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcGet("Actiontrail", "2020-07-06", action, query, request)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || helper.NeedRetry(err) {
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
				actionTrailServiceV2 := NewActionTrailServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Disable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actionTrailServiceV2.ActionTrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "UpdateTrail"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Name"] = d.Id()

	if !d.IsNewResource() && d.HasChange("sls_project_arn") {
		update = true
	}
	if v, ok := d.GetOk("sls_project_arn"); ok || d.HasChange("sls_project_arn") {
		request["SlsProjectArn"] = v
	}
	if !d.IsNewResource() && d.HasChange("max_compute_project_arn") {
		update = true
	}
	if v, ok := d.GetOk("max_compute_project_arn"); ok || d.HasChange("max_compute_project_arn") {
		request["MaxComputeProjectArn"] = v
	}
	if !d.IsNewResource() && d.HasChange("event_rw") {
		update = true
	}
	if v, ok := d.GetOk("event_rw"); ok || d.HasChange("event_rw") {
		request["EventRW"] = v
	}
	if !d.IsNewResource() && d.HasChange("oss_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("oss_write_role_arn"); ok || d.HasChange("oss_write_role_arn") {
		request["OssWriteRoleArn"] = v
	}
	if !d.IsNewResource() && d.HasChange("oss_bucket_name") {
		update = true
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok || d.HasChange("oss_bucket_name") {
		request["OssBucketName"] = v
	}
	if !d.IsNewResource() && d.HasChange("oss_key_prefix") {
		update = true
	}
	if v, ok := d.GetOk("oss_key_prefix"); ok || d.HasChange("oss_key_prefix") {
		request["OssKeyPrefix"] = v
	}
	if !d.IsNewResource() && d.HasChange("sls_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("sls_write_role_arn"); ok || d.HasChange("sls_write_role_arn") {
		request["SlsWriteRoleArn"] = v
	}
	if !d.IsNewResource() && d.HasChange("trail_region") {
		update = true
	}
	if v, ok := d.GetOk("trail_region"); ok || d.HasChange("trail_region") {
		request["TrailRegion"] = v
	}
	if !d.IsNewResource() && d.HasChange("max_compute_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("max_compute_write_role_arn"); ok || d.HasChange("max_compute_write_role_arn") {
		request["MaxComputeWriteRoleArn"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || helper.NeedRetry(err) {
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
	update = false
	action = "PutDataEventSelector"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TrailName"] = d.Id()

	if d.HasChange("data_event_trail_region") {
		update = true
	}
	if v, ok := d.GetOk("data_event_trail_region"); ok || d.HasChange("data_event_trail_region") {
		request["TrailRegionIds"] = v
	}
	if d.HasChange("event_selectors") {
		update = true
	}
	request["EventSelectors"] = d.Get("event_selectors")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
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

	d.Partial(false)
	return ResourceTrailRead(d, meta)
}

func ResourceTrailDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrail"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"TrailNotFoundException"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
