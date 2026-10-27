// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package nlb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

func ResourceLoadBalancerZoneShiftedAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerZoneShiftedAttachmentCreate,
		Read:   resourceLoadBalancerZoneShiftedAttachmentRead,
		Delete: resourceLoadBalancerZoneShiftedAttachmentDelete,
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

func resourceLoadBalancerZoneShiftedAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

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
	request["ClientToken"] = helper.BuildClientToken(action)

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.VSwitchId", d.Get("vswitch_id"))
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.ZoneId", d.Get("zone_id"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_nlb_load_balancer_zone_shifted_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	ZoneMappingsZoneIdVar, _ := jsonpath.Get("ZoneMappings[0].ZoneId", request)
	ZoneMappingsVSwitchIdVar, _ := jsonpath.Get("ZoneMappings[0].VSwitchId", request)
	d.SetId(fmt.Sprintf("%v:%v:%v", request["LoadBalancerId"], ZoneMappingsZoneIdVar, ZoneMappingsVSwitchIdVar))

	return resourceLoadBalancerZoneShiftedAttachmentRead(d, meta)
}

func resourceLoadBalancerZoneShiftedAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NewService(client)

	objectRaw, err := nlbServiceV2.DescribeNlbLoadBalancerZoneShiftedAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer_zone_shifted_attachment DescribeNlbLoadBalancerZoneShiftedAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("load_balancer_id", parts[0])

	return nil
}

func resourceLoadBalancerZoneShiftedAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

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
	request["ClientToken"] = helper.BuildClientToken(action)

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.VSwitchId", parts[2])
	jsonString, _ = sjson.Set(jsonString, "ZoneMappings.0.ZoneId", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
