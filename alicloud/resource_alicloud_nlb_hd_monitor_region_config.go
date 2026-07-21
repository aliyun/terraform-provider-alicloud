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

func resourceAliCloudNlbHdMonitorRegionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbHdMonitorRegionConfigCreate,
		Read:   resourceAliCloudNlbHdMonitorRegionConfigRead,
		Update: resourceAliCloudNlbHdMonitorRegionConfigUpdate,
		Delete: resourceAliCloudNlbHdMonitorRegionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"log_project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_store": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudNlbHdMonitorRegionConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SetHdMonitorRegionConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
	}
	request["RegionId"] = client.RegionId

	request["LogProject"] = d.Get("log_project")
	request["MetricStore"] = d.Get("metric_store")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_hd_monitor_region_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RegionId"]))

	return resourceAliCloudNlbHdMonitorRegionConfigRead(d, meta)
}

func resourceAliCloudNlbHdMonitorRegionConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbHdMonitorRegionConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_hd_monitor_region_config DescribeNlbHdMonitorRegionConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("log_project", objectRaw["LogProject"])
	d.Set("metric_store", objectRaw["MetricStore"])
	d.Set("region_id", objectRaw["RegionId"])

	d.Set("region_id", d.Id())

	return nil
}

func resourceAliCloudNlbHdMonitorRegionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetHdMonitorRegionConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("log_project") {
		update = true
	}
	request["LogProject"] = d.Get("log_project")
	if d.HasChange("metric_store") {
		update = true
	}
	request["MetricStore"] = d.Get("metric_store")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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

	return resourceAliCloudNlbHdMonitorRegionConfigRead(d, meta)
}

func resourceAliCloudNlbHdMonitorRegionConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHdMonitorRegionConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.RegionHdMonitorConfig"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
