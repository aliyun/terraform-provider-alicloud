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

func resourceAliCloudExpressConnectTrafficQosQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectTrafficQosQueueCreate,
		Read:   resourceAliCloudExpressConnectTrafficQosQueueRead,
		Update: resourceAliCloudExpressConnectTrafficQosQueueUpdate,
		Delete: resourceAliCloudExpressConnectTrafficQosQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_percent": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"qos_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"High", "Medium"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectTrafficQosQueueCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExpressConnectTrafficQosQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = d.Get("qos_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["QueueType"] = d.Get("queue_type")
	if v, ok := d.GetOk("bandwidth_percent"); ok {
		request["BandwidthPercent"] = v
	}
	if v, ok := d.GetOk("queue_description"); ok {
		request["QueueDescription"] = v
	}
	if v, ok := d.GetOk("queue_name"); ok {
		request["QueueName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_traffic_qos_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["QosId"], response["QueueId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosQueueStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudExpressConnectTrafficQosQueueRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectTrafficQosQueue(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_traffic_qos_queue DescribeExpressConnectTrafficQosQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_percent", objectRaw["BandwidthPercent"])
	d.Set("queue_description", objectRaw["QueueDescription"])
	d.Set("queue_name", objectRaw["QueueName"])
	d.Set("queue_type", objectRaw["QueueType"])
	d.Set("status", objectRaw["Status"])
	d.Set("qos_id", objectRaw["QosId"])
	d.Set("queue_id", objectRaw["QueueId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("qos_id", parts[0])
	d.Set("queue_id", parts[1])

	return nil
}

func resourceAliCloudExpressConnectTrafficQosQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectTrafficQosQueue"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["QueueId"] = parts[1]
	query["QosId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("bandwidth_percent") {
		update = true
		request["BandwidthPercent"] = d.Get("bandwidth_percent")
	}

	if d.HasChange("queue_description") {
		update = true
		request["QueueDescription"] = d.Get("queue_description")
	}

	if d.HasChange("queue_name") {
		update = true
		request["QueueName"] = d.Get("queue_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosQueueStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	return resourceAliCloudExpressConnectTrafficQosQueueRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosQueueDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExpressConnectTrafficQosQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = parts[0]
	query["QueueId"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"EcQoSConflict", "IncorrectStatus.Qos"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"IllegalParam.%s"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosQueueStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}
	return nil
}
