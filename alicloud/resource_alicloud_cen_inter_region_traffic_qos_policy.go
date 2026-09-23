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

func resourceAliCloudCenInterRegionTrafficQosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenInterRegionTrafficQosPolicyCreate,
		Read:   resourceAliCloudCenInterRegionTrafficQosPolicyRead,
		Update: resourceAliCloudCenInterRegionTrafficQosPolicyUpdate,
		Delete: resourceAliCloudCenInterRegionTrafficQosPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_guarantee_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"inter_region_traffic_qos_policy_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inter_region_traffic_qos_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCenInterRegionTrafficQosPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCenInterRegionTrafficQosPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["TransitRouterId"] = d.Get("transit_router_id")
	if v, ok := d.GetOk("inter_region_traffic_qos_policy_name"); ok {
		request["TrafficQosPolicyName"] = v
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_policy_description"); ok {
		request["TrafficQosPolicyDescription"] = v
	}
	request["TransitRouterAttachmentId"] = d.Get("transit_router_attachment_id")
	if v, ok := d.GetOk("bandwidth_guarantee_mode"); ok {
		request["BandwidthGuaranteeMode"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TransitRouterAttachment", "IncorrectStatus.TransitRouterInstance", "Operation.Blocking"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_inter_region_traffic_qos_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficQosPolicyId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenInterRegionTrafficQosPolicyStateRefreshFunc(d.Id(), "TrafficQosPolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenInterRegionTrafficQosPolicyRead(d, meta)
}

func resourceAliCloudCenInterRegionTrafficQosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenInterRegionTrafficQosPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_inter_region_traffic_qos_policy DescribeCenInterRegionTrafficQosPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_guarantee_mode", objectRaw["BandwidthGuaranteeMode"])
	d.Set("inter_region_traffic_qos_policy_description", objectRaw["TrafficQosPolicyDescription"])
	d.Set("inter_region_traffic_qos_policy_name", objectRaw["TrafficQosPolicyName"])
	d.Set("status", objectRaw["TrafficQosPolicyStatus"])
	d.Set("transit_router_attachment_id", objectRaw["TransitRouterAttachmentId"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])

	return nil
}

func resourceAliCloudCenInterRegionTrafficQosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateCenInterRegionTrafficQosPolicyAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TrafficQosPolicyId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("inter_region_traffic_qos_policy_name") {
		update = true
		request["TrafficQosPolicyName"] = d.Get("inter_region_traffic_qos_policy_name")
	}

	if d.HasChange("inter_region_traffic_qos_policy_description") {
		update = true
		request["TrafficQosPolicyDescription"] = d.Get("inter_region_traffic_qos_policy_description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
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

	return resourceAliCloudCenInterRegionTrafficQosPolicyRead(d, meta)
}

func resourceAliCloudCenInterRegionTrafficQosPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCenInterRegionTrafficQosPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TrafficQosPolicyId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenInterRegionTrafficQosPolicyStateRefreshFunc(d.Id(), "TrafficQosPolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
