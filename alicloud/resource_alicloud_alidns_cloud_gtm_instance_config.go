// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlidnsCloudGtmInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsCloudGtmInstanceConfigCreate,
		Read:   resourceAliCloudAlidnsCloudGtmInstanceConfigRead,
		Update: resourceAliCloudAlidnsCloudGtmInstanceConfigUpdate,
		Delete: resourceAliCloudAlidnsCloudGtmInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_pool_lb_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"round_robin", "sequence", "weight", "source_nearest"}, false),
			},
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disable"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_rr_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"A", "AAAA", "CNAME"}, false),
			},
			"schedule_zone_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"custom"}, false),
			},
			"schedule_zone_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence_lb_strategy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"preemptive", "nonPreemptive"}, false),
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAliCloudAlidnsCloudGtmInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCloudGtmInstanceConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["EnableStatus"] = d.Get("enable_status")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}
	request["Ttl"] = d.Get("ttl")
	request["ChargeType"] = "postpay"
	request["ScheduleZoneMode"] = d.Get("schedule_zone_mode")
	if v, ok := d.GetOk("schedule_host_name"); ok {
		request["ScheduleHostname"] = v
	}
	request["ScheduleRrType"] = d.Get("schedule_rr_type")
	if v, ok := d.GetOk("schedule_zone_name"); ok {
		request["ScheduleZoneName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_cloud_gtm_instance_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["ConfigId"], response["InstanceId"]))

	return resourceAliCloudAlidnsCloudGtmInstanceConfigUpdate(d, meta)
}

func resourceAliCloudAlidnsCloudGtmInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsServiceV2 := AlidnsServiceV2{client}

	objectRaw, err := alidnsServiceV2.DescribeAlidnsCloudGtmInstanceConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_cloud_gtm_instance_config DescribeAlidnsCloudGtmInstanceConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_pool_lb_strategy", objectRaw["AddressPoolLbStrategy"])
	d.Set("enable_status", objectRaw["EnableStatus"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("schedule_host_name", objectRaw["ScheduleHostname"])
	d.Set("schedule_rr_type", objectRaw["ScheduleRrType"])
	d.Set("schedule_zone_mode", objectRaw["ScheduleZoneMode"])
	d.Set("schedule_zone_name", objectRaw["ScheduleZoneName"])
	d.Set("sequence_lb_strategy_mode", objectRaw["SequenceLbStrategyMode"])
	d.Set("ttl", objectRaw["Ttl"])
	d.Set("config_id", objectRaw["ConfigId"])
	d.Set("instance_id", objectRaw["InstanceId"])

	return nil
}

func resourceAliCloudAlidnsCloudGtmInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateCloudGtmInstanceConfigBasic"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("ttl") {
		update = true
	}
	request["Ttl"] = d.Get("ttl")
	if !d.IsNewResource() && d.HasChange("schedule_host_name") {
		update = true
	}
	if v, ok := d.GetOk("schedule_host_name"); ok || d.HasChange("schedule_host_name") {
		request["ScheduleHostname"] = v
	}
	if !d.IsNewResource() && d.HasChange("schedule_zone_name") {
		update = true
	}
	if v, ok := d.GetOk("schedule_zone_name"); ok || d.HasChange("schedule_zone_name") {
		request["ScheduleZoneName"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateCloudGtmInstanceConfigEnableStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("enable_status") {
		update = true
	}
	request["EnableStatus"] = d.Get("enable_status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateCloudGtmInstanceConfigLbStrategy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("sequence_lb_strategy_mode") {
		update = true
	}
	if v, ok := d.GetOk("sequence_lb_strategy_mode"); ok || d.HasChange("sequence_lb_strategy_mode") {
		request["SequenceLbStrategyMode"] = v
	}
	if d.HasChange("address_pool_lb_strategy") {
		update = true
	}
	request["AddressPoolLbStrategy"] = d.Get("address_pool_lb_strategy")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateCloudGtmInstanceConfigRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
	}
	if v, ok := d.GetOk("remark"); ok || d.HasChange("remark") {
		request["Remark"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudAlidnsCloudGtmInstanceConfigRead(d, meta)
}

func resourceAliCloudAlidnsCloudGtmInstanceConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCloudGtmInstanceConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
