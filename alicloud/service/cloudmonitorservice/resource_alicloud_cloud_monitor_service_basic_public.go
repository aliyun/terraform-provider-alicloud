// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudMonitorServiceBasicPublic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceBasicPublicCreate,
		Read:   resourceAliCloudCloudMonitorServiceBasicPublicRead,
		Delete: resourceAliCloudCloudMonitorServiceBasicPublicDelete,
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

func resourceAliCloudCloudMonitorServiceBasicPublicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	query := make(map[string]interface{})
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)
	request["ProductCode"] = "cms"
	request["ProductType"] = "cms_basic_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "cms_basic_public_intl"
	}
	request["SubscriptionType"] = "PayAsYouGo"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = helper.BuildClientToken(action)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cms_basic_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Has.effect.suit"}) {
			parts := strings.Split(response["Message"].(string), ": ")
			if len(parts) < 2 {
				return errs.WrapErrorf(err, errs.ResponseCodeMsg, "alicloud_cloud_monitor_service_basic_public", action, response)
			}
			d.SetId(parts[1])
			return resourceAliCloudCloudMonitorServiceBasicPublicRead(d, meta)
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_monitor_service_basic_public", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudMonitorServiceBasicPublicRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceBasicPublicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceBasicPublic(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_basic_public DescribeCloudMonitorServiceBasicPublic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])

	return nil
}

func resourceAliCloudCloudMonitorServiceBasicPublicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "StopPostPayQuota"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	request["PostType"] = "postPayV2"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudMonitorServiceServiceV2.CloudMonitorServiceBasicPublicStateRefreshFunc(d.Id(), "InstanceID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
