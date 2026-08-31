// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenFlowLogCreate,
		Read:   resourceAliCloudCenFlowLogRead,
		Update: resourceAliCloudCenFlowLogUpdate,
		Delete: resourceAliCloudCenFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 60, 600}),
			},
			"log_format_string": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCenFlowLogCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFlowlog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["ProjectName"] = d.Get("project_name")
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["LogStoreName"] = d.Get("log_store_name")
	if v, ok := d.GetOk("flow_log_name"); ok {
		request["FlowLogName"] = v
	}
	if v, ok := d.GetOkExists("interval"); ok {
		request["Interval"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	if v, ok := d.GetOk("log_format_string"); ok {
		request["LogFormatString"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_flowlog", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FlowLogId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenFlowLogUpdate(d, meta)
}

func resourceAliCloudCenFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenFlowLog(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_flowlog DescribeCenFlowLog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CenId"] != nil {
		d.Set("cen_id", objectRaw["CenId"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["FlowLogName"] != nil {
		d.Set("flow_log_name", objectRaw["FlowLogName"])
	}
	if objectRaw["Interval"] != nil {
		d.Set("interval", objectRaw["Interval"])
	}
	if objectRaw["LogFormatString"] != nil {
		d.Set("log_format_string", objectRaw["LogFormatString"])
	}
	if objectRaw["LogStoreName"] != nil {
		d.Set("log_store_name", objectRaw["LogStoreName"])
	}
	if objectRaw["ProjectName"] != nil {
		d.Set("project_name", objectRaw["ProjectName"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["TransitRouterAttachmentId"] != nil {
		d.Set("transit_router_attachment_id", objectRaw["TransitRouterAttachmentId"])
	}
	if objectRaw["TransitRouterId"] != nil {
		d.Set("transit_router_id", objectRaw["TransitRouterId"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudCenFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if d.HasChange("status") {
		cenServiceV2 := CenServiceV2{client}
		object, err := cenServiceV2.DescribeCenFlowLog(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Active" {
				action := "ActiveFlowLog"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FlowLogId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR", "Operation.Blocking", "IncorrectStatus.TrFlowlog"}) || NeedRetry(err) {
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
				cenServiceV2 := CenServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Inactive" {
				action := "DeactiveFlowLog"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FlowLogId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR", "Operation.Blocking", "IncorrectStatus.TrFlowlog"}) || NeedRetry(err) {
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
				cenServiceV2 := CenServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Inactive"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	action := "ModifyFlowLogAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FlowLogId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("flow_log_name") {
		update = true
		request["FlowLogName"] = d.Get("flow_log_name")
	}

	if !d.IsNewResource() && d.HasChange("interval") {
		update = true
		request["Interval"] = d.Get("interval")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"LOCK_ERROR", "Operation.Blocking", "IncorrectStatus.flowlog"}) || NeedRetry(err) {
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
		cenServiceV2 := CenServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		cenServiceV2 := CenServiceV2{client}
		if err := cenServiceV2.SetResourceTags(d, "flowlog"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenFlowLogRead(d, meta)
}

func resourceAliCloudCenFlowLogDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFlowlog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["FlowLogId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrFlowlog"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"SourceProjectNotExist", "InvalidFlowlogId.NotFound", "ProjectOrLogstoreNotExist", "SourceProjectNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
