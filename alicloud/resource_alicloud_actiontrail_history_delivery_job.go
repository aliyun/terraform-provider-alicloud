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

func resourceAliCloudActionTrailHistoryDeliveryJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudActionTrailHistoryDeliveryJobCreate,
		Read:   resourceAliCloudActionTrailHistoryDeliveryJobRead,
		Delete: resourceAliCloudActionTrailHistoryDeliveryJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"trail_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudActionTrailHistoryDeliveryJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDeliveryHistoryJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["TrailName"] = d.Get("trail_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_history_delivery_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["JobId"]))

	actionTrailServiceV2 := ActionTrailServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"1", "2", "3"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, actionTrailServiceV2.ActionTrailHistoryDeliveryJobStateRefreshFunc(d.Id(), "JobStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudActionTrailHistoryDeliveryJobRead(d, meta)
}

func resourceAliCloudActionTrailHistoryDeliveryJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actionTrailServiceV2 := ActionTrailServiceV2{client}

	objectRaw, err := actionTrailServiceV2.DescribeActionTrailHistoryDeliveryJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_history_delivery_job DescribeActionTrailHistoryDeliveryJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("status", objectRaw["JobStatus"])
	d.Set("trail_name", objectRaw["TrailName"])

	return nil
}

func resourceAliCloudActionTrailHistoryDeliveryJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeliveryHistoryJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["JobId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
