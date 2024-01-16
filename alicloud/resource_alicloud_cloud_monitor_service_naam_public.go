// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceNaamPublic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceNaamPublicCreate,
		Read:   resourceAliCloudCloudMonitorServiceNaamPublicRead,
		Delete: resourceAliCloudCloudMonitorServiceNaamPublicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceNaamPublicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["ProductCode"] = "cms"
	request["ProductType"] = "cms_naam_public_cn"
	request["SubscriptionType"] = "PayAsYouGo"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cms_naam_public_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_monitor_service_naam_public", action, AlibabaCloudSdkGoERROR)
	}

	if response["Code"] == "Has.effect.suit" {
		parts := strings.Split(response["Message"].(string), ": ")
		if len(parts) < 2 {
			return WrapErrorf(err, ResponseCodeMsg, "alicloud_cloud_monitor_service_naam_public", action, response)
		}
		d.SetId(parts[1])
		return resourceAliCloudCloudMonitorServiceNaamPublicRead(d, meta)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudMonitorServiceNaamPublicRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceNaamPublicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceNaamPublic(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_naam_public DescribeCloudMonitorServiceNaamPublic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])

	return nil
}

func resourceAliCloudCloudMonitorServiceNaamPublicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "StopPostPayQuota"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewCloudmonitorserviceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["PostType"] = "postNaamV2"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), query, request, &runtime)

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

	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudMonitorServiceServiceV2.CloudMonitorServiceNaamPublicStateRefreshFunc(d.Id(), "InstanceID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertCloudMonitorServiceNaamPublicPostTypeRequest(source interface{}) interface{} {
	switch source {
	case "cms_naam_public_cn":
		return "postNaamV2"
	case "cms_naam_public_intl":
		return "postNaamV2"
	}
	return source
}
