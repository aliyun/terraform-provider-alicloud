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

func resourceAliCloudActiontrailTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudActiontrailTrailCreate,
		Read:   resourceAliCloudActiontrailTrailRead,
		Update: resourceAliCloudActiontrailTrailUpdate,
		Delete: resourceAliCloudActiontrailTrailDelete,
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
			"event_rw": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Write", "Read", "All"}, false),
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Enable",
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
			},
			"trail_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"trail_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trail_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.95.0. New field `trail_name` instead.",
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

func resourceAliCloudActiontrailTrailCreate(d *schema.ResourceData, meta interface{}) error {

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
		return WrapError(Error(`[ERROR] Argument "trail_name" or "name" must be set one!`))
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_trail", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	actiontrailServiceV2 := ActiontrailServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Fresh"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, actiontrailServiceV2.ActiontrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudActiontrailTrailUpdate(d, meta)
}

func resourceAliCloudActiontrailTrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailServiceV2 := ActiontrailServiceV2{client}

	objectRaw, err := actiontrailServiceV2.DescribeActiontrailTrail(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_trail DescribeActiontrailTrail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	d.Set("name", objectRaw["Name"])

	return nil
}

func resourceAliCloudActiontrailTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		var err error
		actiontrailServiceV2 := ActiontrailServiceV2{client}
		object, err := actiontrailServiceV2.DescribeActiontrailTrail(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Enable" {
				action := "StartLogging"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Name"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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

				stateConf := BuildStateConf([]string{}, []string{"Enable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailServiceV2.ActiontrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if target == "Disable" {
				action := "StopLogging"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Name"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcGet("Actiontrail", "2020-07-06", action, query, request)
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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

				stateConf := BuildStateConf([]string{}, []string{"Disable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailServiceV2.ActiontrailTrailStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}

	var err error
	action := "UpdateTrail"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Name"] = d.Id()

	if !d.IsNewResource() && d.HasChange("max_compute_project_arn") {
		update = true
	}
	if v, ok := d.GetOk("max_compute_project_arn"); ok {
		request["MaxComputeProjectArn"] = v
	}

	if !d.IsNewResource() && d.HasChange("event_rw") {
		update = true
	}
	if v, ok := d.GetOk("event_rw"); ok {
		request["EventRW"] = v
	}

	if !d.IsNewResource() && d.HasChange("oss_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("oss_write_role_arn"); ok {
		request["OssWriteRoleArn"] = v
	}

	if !d.IsNewResource() && d.HasChange("oss_bucket_name") {
		update = true
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok {
		request["OssBucketName"] = v
	}

	if !d.IsNewResource() && d.HasChange("oss_key_prefix") {
		update = true
	}
	if v, ok := d.GetOk("oss_key_prefix"); ok {
		request["OssKeyPrefix"] = v
	}

	if !d.IsNewResource() && d.HasChange("sls_project_arn") {
		update = true
	}
	if v, ok := d.GetOk("sls_project_arn"); ok {
		request["SlsProjectArn"] = v
	}

	if !d.IsNewResource() && d.HasChange("sls_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("sls_write_role_arn"); ok {
		request["SlsWriteRoleArn"] = v
	}

	if !d.IsNewResource() && d.HasChange("trail_region") {
		update = true
	}
	if v, ok := d.GetOk("trail_region"); ok {
		request["TrailRegion"] = v
	}

	if !d.IsNewResource() && d.HasChange("max_compute_write_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("max_compute_write_role_arn"); ok {
		request["MaxComputeWriteRoleArn"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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

	d.Partial(false)
	return resourceAliCloudActiontrailTrailRead(d, meta)
}

func resourceAliCloudActiontrailTrailDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrail"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"TrailNotFoundException"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
