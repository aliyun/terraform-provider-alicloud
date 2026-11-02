// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package expressconnect

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

func resourceAliCloudExpressConnectEcFailoverTestJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectEcFailoverTestJobCreate,
		Read:   resourceAliCloudExpressConnectEcFailoverTestJobRead,
		Update: resourceAliCloudExpressConnectEcFailoverTestJobUpdate,
		Delete: resourceAliCloudExpressConnectEcFailoverTestJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ec_failover_test_job_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 4320),
			},
			"job_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"StartNow", "StartLater"}, true),
			},
			"resource_id": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PHYSICALCONNECTION"}, true),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Init", "Starting", "Testing", "Stopping", "Stopped", "Deleted"}, true),
			},
		},
	}
}

func resourceAliCloudExpressConnectEcFailoverTestJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFailoverTestJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["ResourceType"] = d.Get("resource_type")
	if v, ok := d.GetOk("resource_id"); ok {
		resourceIdMaps := v.(*schema.Set).List()
		request["ResourceId"] = resourceIdMaps
	}

	request["JobDuration"] = d.Get("job_duration")
	request["JobType"] = d.Get("job_type")
	if v, ok := d.GetOk("ec_failover_test_job_name"); ok {
		request["Name"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_ec_failover_test_job", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["JobId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Init", "Testing", "Stopped"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectServiceV2.ExpressConnectEcFailoverTestJobStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectEcFailoverTestJobUpdate(d, meta)
}

func resourceAliCloudExpressConnectEcFailoverTestJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectEcFailoverTestJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_ec_failover_test_job DescribeExpressConnectEcFailoverTestJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("ec_failover_test_job_name", objectRaw["Name"])
	d.Set("job_duration", objectRaw["JobDuration"])
	d.Set("job_type", objectRaw["JobType"])
	d.Set("resource_type", objectRaw["ResourceType"])
	d.Set("status", objectRaw["Status"])

	resourceId1Raw := make([]interface{}, 0)
	if objectRaw["ResourceId"] != nil {
		resourceId1Raw = objectRaw["ResourceId"].([]interface{})
	}

	d.Set("resource_id", resourceId1Raw)

	return nil
}

func resourceAliCloudExpressConnectEcFailoverTestJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateFailoverTestJob"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["JobId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("resource_id") {
		update = true
		if v, ok := d.GetOk("resource_id"); ok {
			resourceIdMaps := v.(*schema.Set).List()
			request["ResourceId"] = resourceIdMaps
		}
	}

	if !d.IsNewResource() && d.HasChange("job_duration") {
		update = true
	}
	request["JobDuration"] = d.Get("job_duration")
	if !d.IsNewResource() && d.HasChange("ec_failover_test_job_name") {
		update = true
		request["Name"] = d.Get("ec_failover_test_job_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		object, err := expressConnectServiceV2.DescribeExpressConnectEcFailoverTestJob(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Testing" {
				action = "StartFailoverTestJob"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["JobId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
					request["ClientToken"] = helper.BuildClientToken(action)

					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(action, response, request)
					return nil
				})
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				expressConnectServiceV2 := ExpressConnectServiceV2{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Testing"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectEcFailoverTestJobStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				action = "StopFailoverTestJob"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["JobId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
					request["ClientToken"] = helper.BuildClientToken(action)

					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(action, response, request)
					return nil
				})
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				expressConnectServiceV2 := ExpressConnectServiceV2{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectEcFailoverTestJobStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	return resourceAliCloudExpressConnectEcFailoverTestJobRead(d, meta)
}

func resourceAliCloudExpressConnectEcFailoverTestJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFailoverTestJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["JobId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.OnlyForInitOrStopped"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"IllegalParam.JobId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
