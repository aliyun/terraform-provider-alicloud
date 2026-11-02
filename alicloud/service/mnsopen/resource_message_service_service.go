// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package mnsopen

import (
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	request["ClientToken"] = tfhelper.BuildClientToken(action)

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
	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"INSTANCE_ID_IS_NOT_UNIQUE", "ORDER.OPEND"}) {
				return nil
			}
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && tferrs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "mns"
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_message_service_service", action, tferrs.AlibabaCloudSdkGoERROR)
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
		if !d.IsNewResource() && tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_service DescribeMessageServiceService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}

	instanceListRawArrayObj, _ := jsonpath.Get("$.Data.InstanceList[*]", objectRaw)
	instanceListRawArray := make([]interface{}, 0)
	if instanceListRawArrayObj != nil {
		instanceListRawArray = tfconvert.ConvertToInterfaceArray(instanceListRawArrayObj)
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
