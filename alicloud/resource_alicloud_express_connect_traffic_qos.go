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

func resourceAliCloudExpressConnectTrafficQos() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectTrafficQosCreate,
		Read:   resourceAliCloudExpressConnectTrafficQosRead,
		Update: resourceAliCloudExpressConnectTrafficQosUpdate,
		Delete: resourceAliCloudExpressConnectTrafficQosDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"qos_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"qos_name": {
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

func resourceAliCloudExpressConnectTrafficQosCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "CreateExpressConnectTrafficQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("qos_name"); ok {
		request["QosName"] = v
	}
	if v, ok := d.GetOk("qos_description"); ok {
		request["QosDescription"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_traffic_qos", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["QosId"]))

	return resourceAliCloudExpressConnectTrafficQosRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectTrafficQos(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_traffic_qos DescribeExpressConnectTrafficQos Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("qos_description", objectRaw["QosDescription"])
	d.Set("qos_name", objectRaw["QosName"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudExpressConnectTrafficQosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyExpressConnectTrafficQos"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["QosId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("qos_description") {
		update = true
		request["QosDescription"] = d.Get("qos_description")
	}

	if d.HasChange("qos_name") {
		update = true
		request["QosName"] = d.Get("qos_name")
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
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 0, expressConnectServiceV2.ExpressConnectTrafficQosStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudExpressConnectTrafficQosRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExpressConnectTrafficQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"IllegalParam.%s"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
