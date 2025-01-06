package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenInterRegionTrafficQosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInterRegionTrafficQosPolicyCreate,
		Read:   resourceAlicloudCenInterRegionTrafficQosPolicyRead,
		Update: resourceAlicloudCenInterRegionTrafficQosPolicyUpdate,
		Delete: resourceAlicloudCenInterRegionTrafficQosPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inter_region_traffic_qos_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inter_region_traffic_qos_policy_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenInterRegionTrafficQosPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateCenInterRegionTrafficQosPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken("CreateCenInterRegionTrafficQosPolicy")
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["TransitRouterAttachmentId"] = d.Get("transit_router_attachment_id")

	if v, ok := d.GetOk("inter_region_traffic_qos_policy_name"); ok {
		request["TrafficQosPolicyName"] = v
	}

	if v, ok := d.GetOk("inter_region_traffic_qos_policy_description"); ok {
		request["TrafficQosPolicyDescription"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenInterRegionTrafficQosPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenInterRegionTrafficQosPolicyRead(d, meta)
}

func resourceAlicloudCenInterRegionTrafficQosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenInterRegionTrafficQosPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("inter_region_traffic_qos_policy_name", object["TrafficQosPolicyName"])
	d.Set("inter_region_traffic_qos_policy_description", object["TrafficQosPolicyDescription"])
	d.Set("status", object["TrafficQosPolicyStatus"])

	return nil
}

func resourceAlicloudCenInterRegionTrafficQosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"ClientToken":        buildClientToken("UpdateCenInterRegionTrafficQosPolicyAttribute"),
		"TrafficQosPolicyId": d.Id(),
	}

	if d.HasChange("inter_region_traffic_qos_policy_name") {
		update = true
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_policy_name"); ok {
		request["TrafficQosPolicyName"] = v
	}

	if d.HasChange("inter_region_traffic_qos_policy_description") {
		update = true
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_policy_description"); ok {
		request["TrafficQosPolicyDescription"] = v
	}

	if update {
		action := "UpdateCenInterRegionTrafficQosPolicyAttribute"
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenInterRegionTrafficQosPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudCenInterRegionTrafficQosPolicyRead(d, meta)
}

func resourceAlicloudCenInterRegionTrafficQosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteCenInterRegionTrafficQosPolicy"
	var response map[string]interface{}

	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":        buildClientToken("DeleteCenInterRegionTrafficQosPolicy"),
		"TrafficQosPolicyId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenInterRegionTrafficQosPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
