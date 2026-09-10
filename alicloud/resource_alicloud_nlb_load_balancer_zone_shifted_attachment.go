// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudNlbLoadBalancerZoneShiftedAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentCreate,
		Read:   resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentRead,
		Delete: resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "StartShiftLoadBalancerZones"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.VSwitchId", d.Get("vswitch_id"))
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.ZoneId", d.Get("zone_id"))
	_ = json.Unmarshal([]byte(jsonString), &request)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer_zone_shifted_attachment", action, AlibabaCloudSdkGoERROR)
	}

	ZoneMappingsZoneIdVar, _ := jsonpath.Get("ZoneMappings[0].ZoneId", request)
	ZoneMappingsVSwitchIdVar, _ := jsonpath.Get("ZoneMappings[0].VSwitchId", request)
	d.SetId(fmt.Sprintf("%v:%v:%v", request["LoadBalancerId"], ZoneMappingsZoneIdVar, ZoneMappingsVSwitchIdVar))

	return resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentRead(d, meta)
}

func resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbLoadBalancerZoneShiftedAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer_zone_shifted_attachment DescribeNlbLoadBalancerZoneShiftedAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("load_balancer_id", parts[0])

	return nil
}

func resourceAliCloudNlbLoadBalancerZoneShiftedAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "CancelShiftLoadBalancerZones"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.VSwitchId", parts[2])
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.ZoneId", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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

	return nil
}
