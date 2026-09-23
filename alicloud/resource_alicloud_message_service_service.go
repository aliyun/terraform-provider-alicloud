// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMessageServiceService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceServiceCreate,
		Read:   resourceAliCloudMessageServiceServiceRead,
		Delete: resourceAliCloudMessageServiceServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudMessageServiceServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["SubscriptionType"] = "PayAsYouGo"
	dataList := make(map[string]interface{})

	if v, ok := d.GetOk("code"); ok {
		dataList["Code"] = v
	}

	if v, ok := d.GetOk("value"); ok {
		dataList["Value"] = v
	}

	dataList["Code"] = "commodity_type"
	if client.IsInternationalAccount() {
		dataList["Value"] = "mns"
	} else {
		dataList["Value"] = "commodity_type:mns"
	}
	ParameterMap := make([]interface{}, 0)
	ParameterMap = append(ParameterMap, dataList)
	request["Parameter"] = ParameterMap

	var endpoint string
	request["ProductCode"] = "mns"
	request["ProductType"] = ""
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if IsExpectedErrors(err, []string{"INSTANCE_ID_IS_NOT_UNIQUE", "ORDER.OPEND"}) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "mns"
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_service", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudMessageServiceServiceRead(d, meta)
}

func resourceAliCloudMessageServiceServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_service DescribeMessageServiceService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	instanceListRawArrayObj, _ := jsonpath.Get("$.Data.InstanceList[*]", objectRaw)
	instanceListRawArray := make([]interface{}, 0)
	if instanceListRawArrayObj != nil {
		instanceListRawArray = convertToInterfaceArray(instanceListRawArrayObj)
	}
	instanceListRaw := make(map[string]interface{})
	if len(instanceListRawArray) > 0 {
		instanceListRaw = instanceListRawArray[0].(map[string]interface{})
	}

	d.Set("status", instanceListRaw["Status"])

	return nil
}

func resourceAliCloudMessageServiceServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Service. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
