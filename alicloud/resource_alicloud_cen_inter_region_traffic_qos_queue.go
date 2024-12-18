package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenInterRegionTrafficQosQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInterRegionTrafficQosQueueCreate,
		Read:   resourceAlicloudCenInterRegionTrafficQosQueueRead,
		Update: resourceAlicloudCenInterRegionTrafficQosQueueUpdate,
		Delete: resourceAlicloudCenInterRegionTrafficQosQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"traffic_qos_policy_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"remain_bandwidth_percent": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"dscps": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"inter_region_traffic_qos_queue_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"inter_region_traffic_qos_queue_description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCenInterRegionTrafficQosQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("dscps"); ok {
		dscps := v.([]interface{})
		for index, dscp := range dscps {
			request["Dscps."+fmt.Sprintf(strconv.Itoa(index+1))] = dscp
		}
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_queue_description"); ok {
		request["QosQueueDescription"] = v
	}
	if v, ok := d.GetOk("inter_region_traffic_qos_queue_name"); ok {
		request["QosQueueName"] = v
	}
	if v, ok := d.GetOk("remain_bandwidth_percent"); ok {
		request["RemainBandwidthPercent"] = v
	}
	if v, ok := d.GetOk("traffic_qos_policy_id"); ok {
		request["TrafficQosPolicyId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateCenInterRegionTrafficQosQueue")
	var response map[string]interface{}
	action := "CreateCenInterRegionTrafficQosQueue"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficQosPolicy"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_inter_region_traffic_qos_queue", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.QosQueueId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cen_inter_region_traffic_qos_queue")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenInterRegionTrafficQosQueueStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCenInterRegionTrafficQosQueueRead(d, meta)
}

func resourceAlicloudCenInterRegionTrafficQosQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenInterRegionTrafficQosQueue(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_inter_region_traffic_qos_queue cbnService.DescribeCenInterRegionTrafficQosQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("dscps", convertJsonStringToStringList(object["Dscps"]))
	d.Set("inter_region_traffic_qos_queue_description", object["TrafficQosQueueDescription"])
	d.Set("inter_region_traffic_qos_queue_name", object["TrafficQosQueueName"])
	d.Set("remain_bandwidth_percent", object["RemainBandwidthPercent"])
	d.Set("status", object["Status"])
	d.Set("traffic_qos_policy_id", object["TrafficQosPolicyId"])

	return nil
}

func resourceAlicloudCenInterRegionTrafficQosQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cbnService := CbnService{client}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"QosQueueId": d.Id(),
	}

	if d.HasChange("dscps") {
		update = true
		if v, ok := d.GetOk("dscps"); ok {
			dscps := v.([]interface{})
			for index, dscp := range dscps {
				request["Dscps."+fmt.Sprintf(strconv.Itoa(index+1))] = dscp
			}
		}
	}
	if d.HasChange("inter_region_traffic_qos_queue_description") {
		update = true
		if v, ok := d.GetOk("inter_region_traffic_qos_queue_description"); ok {
			request["QosQueueDescription"] = v
		}
	}
	if d.HasChange("inter_region_traffic_qos_queue_name") {
		update = true
		if v, ok := d.GetOk("inter_region_traffic_qos_queue_name"); ok {
			request["QosQueueName"] = v
		}
	}
	if d.HasChange("remain_bandwidth_percent") {
		update = true
		if v, ok := d.GetOk("remain_bandwidth_percent"); ok {
			request["RemainBandwidthPercent"] = v
		}
	}

	if update {
		action := "UpdateCenInterRegionTrafficQosQueueAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficQosPolicy"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenInterRegionTrafficQosQueueStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudCenInterRegionTrafficQosQueueRead(d, meta)
}

func resourceAlicloudCenInterRegionTrafficQosQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"QosQueueId": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteCenInterRegionTrafficQosQueue")
	action := "DeleteCenInterRegionTrafficQosQueue"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking","IncorrectStatus.TrafficQosPolicy"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenInterRegionTrafficQosQueueStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
