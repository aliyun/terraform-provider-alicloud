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

func resourceAliCloudCenInterRegionTrafficQosQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenInterRegionTrafficQosQueueCreate,
		Read:   resourceAliCloudCenInterRegionTrafficQosQueueRead,
		Update: resourceAliCloudCenInterRegionTrafficQosQueueUpdate,
		Delete: resourceAliCloudCenInterRegionTrafficQosQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dscps": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"inter_region_traffic_qos_queue_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inter_region_traffic_qos_queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remain_bandwidth_percent": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_qos_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCenInterRegionTrafficQosQueueCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCenInterRegionTrafficQosQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("inter_region_traffic_qos_queue_name"); ok {
		request["QosQueueName"] = v
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_queue_description"); ok {
		request["QosQueueDescription"] = v
	}
	if v, ok := d.GetOkExists("remain_bandwidth_percent"); ok {
		request["RemainBandwidthPercent"] = v
	}
	if v, ok := d.GetOk("dscps"); ok {
		dscpsMapsArray := v.([]interface{})
		request["Dscps"] = dscpsMapsArray
	}

	request["TrafficQosPolicyId"] = d.Get("traffic_qos_policy_id")
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficQosPolicy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_inter_region_traffic_qos_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["QosQueueId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenInterRegionTrafficQosQueueStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenInterRegionTrafficQosQueueRead(d, meta)
}

func resourceAliCloudCenInterRegionTrafficQosQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenInterRegionTrafficQosQueue(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_inter_region_traffic_qos_queue DescribeCenInterRegionTrafficQosQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("inter_region_traffic_qos_queue_description", objectRaw["TrafficQosQueueDescription"])
	d.Set("inter_region_traffic_qos_queue_name", objectRaw["TrafficQosQueueName"])
	d.Set("remain_bandwidth_percent", objectRaw["RemainBandwidthPercent"])
	d.Set("status", objectRaw["Status"])
	d.Set("traffic_qos_policy_id", objectRaw["TrafficQosPolicyId"])

	dscpsRaw := make([]interface{}, 0)
	if objectRaw["Dscps"] != nil {
		dscpsRaw = objectRaw["Dscps"].([]interface{})
	}

	d.Set("dscps", dscpsRaw)

	return nil
}

func resourceAliCloudCenInterRegionTrafficQosQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateCenInterRegionTrafficQosQueueAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["QosQueueId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("remain_bandwidth_percent") {
		update = true
		request["RemainBandwidthPercent"] = d.Get("remain_bandwidth_percent")
	}

	if d.HasChange("dscps") {
		update = true
	}
	if v, ok := d.GetOk("dscps"); ok || d.HasChange("dscps") {
		dscpsMapsArray := v.([]interface{})
		request["Dscps"] = dscpsMapsArray
	}

	if d.HasChange("inter_region_traffic_qos_queue_name") {
		update = true
		request["QosQueueName"] = d.Get("inter_region_traffic_qos_queue_name")
	}

	if d.HasChange("inter_region_traffic_qos_queue_description") {
		update = true
		request["QosQueueDescription"] = d.Get("inter_region_traffic_qos_queue_description")
	}

	if d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficQosPolicy"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenInterRegionTrafficQosQueueStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCenInterRegionTrafficQosQueueRead(d, meta)
}

func resourceAliCloudCenInterRegionTrafficQosQueueDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCenInterRegionTrafficQosQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["QosQueueId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficQosPolicy"}) || NeedRetry(err) {
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

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenInterRegionTrafficQosQueueStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
