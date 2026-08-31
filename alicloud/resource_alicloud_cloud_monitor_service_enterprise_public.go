// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceEnterprisePublic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceEnterprisePublicCreate,
		Read:   resourceAliCloudCloudMonitorServiceEnterprisePublicRead,
		Delete: resourceAliCloudCloudMonitorServiceEnterprisePublicDelete,
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

func resourceAliCloudCloudMonitorServiceEnterprisePublicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	query := make(map[string]interface{})
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)
	request["ProductCode"] = "cms"
	request["ProductType"] = "cms_enterprise_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "cms_enterprise_public_intl"
	}
	request["SubscriptionType"] = "PayAsYouGo"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cms_enterprise_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_monitor_service_enterprise_public", action, AlibabaCloudSdkGoERROR)
	}

	if response["Code"] == "Has.effect.suit" {
		parts := strings.Split(response["Message"].(string), ": ")
		if len(parts) < 2 {
			return WrapErrorf(err, ResponseCodeMsg, "alicloud_cloud_monitor_service_enterprise_public", action, response)
		}
		d.SetId(parts[1])
		return resourceAliCloudCloudMonitorServiceEnterprisePublicRead(d, meta)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudMonitorServiceEnterprisePublicRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceEnterprisePublicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceEnterprisePublic(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_enterprise_public DescribeCloudMonitorServiceEnterprisePublic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])

	return nil
}

func resourceAliCloudCloudMonitorServiceEnterprisePublicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "StopPostPayQuota"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["PostType"] = "postEnterprise"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, false)

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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudMonitorServiceServiceV2.CloudMonitorServiceEnterprisePublicStateRefreshFunc(d.Id(), "InstanceID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertCloudMonitorServiceEnterprisePublicPostTypeRequest(source interface{}) interface{} {
	switch source {
	case "cms_enterprise_public_cn":
		return "postEnterprise"
	case "cms_enterprise_public_intl":
		return "postEnterprise"
	}
	return source
}
