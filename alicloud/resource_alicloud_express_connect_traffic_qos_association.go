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

func resourceAliCloudExpressConnectTrafficQosAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectTrafficQosAssociationCreate,
		Read:   resourceAliCloudExpressConnectTrafficQosAssociationRead,
		Delete: resourceAliCloudExpressConnectTrafficQosAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PHYSICALCONNECTION"}, false),
			},
			"qos_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectTrafficQosAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyExpressConnectTrafficQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = d.Get("qos_id")
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "AddInstanceList.0.InstanceId", d.Get("instance_id"))
	jsonString, _ = sjson.Set(jsonString, "AddInstanceList.0.InstanceType", d.Get("instance_type"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_traffic_qos_association", action, AlibabaCloudSdkGoERROR)
	}

	AddInstanceListInstanceId, _ := jsonpath.Get("AddInstanceList[0].InstanceId", request)
	AddInstanceListInstanceType, _ := jsonpath.Get("AddInstanceList[0].InstanceType", request)
	d.SetId(fmt.Sprintf("%v:%v:%v", query["QosId"], AddInstanceListInstanceId, AddInstanceListInstanceType))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosAssociationStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudExpressConnectTrafficQosAssociationRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectTrafficQosAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_traffic_qos_association DescribeExpressConnectTrafficQosAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("instance_type", objectRaw["InstanceType"])

	parts := strings.Split(d.Id(), ":")
	d.Set("qos_id", parts[0])
	d.Set("instance_id", parts[1])
	d.Set("instance_type", parts[2])

	return nil
}

func resourceAliCloudExpressConnectTrafficQosAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectTrafficQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = parts[0]
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "RemoveInstanceList.0.InstanceId", parts[1])
	jsonString, _ = sjson.Set(jsonString, "RemoveInstanceList.0.InstanceType", parts[2])
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosAssociationStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}
	return nil
}
