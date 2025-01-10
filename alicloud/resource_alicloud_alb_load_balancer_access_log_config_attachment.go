// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlbLoadBalancerAccessLogConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentCreate,
		Read:   resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentRead,
		Delete: resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnableLoadBalancerAccessLog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["LogProject"] = d.Get("log_project")
	request["LogStore"] = d.Get("log_store")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer_access_log_config_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutCreate), 0, albServiceV2.DescribeAsyncAlbLoadBalancerAccessLogConfigAttachmentStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentRead(d, meta)
}

func resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbLoadBalancerAccessLogConfigAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer_access_log_config_attachment DescribeAlbLoadBalancerAccessLogConfigAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}

	accessLogConfig1RawObj, _ := jsonpath.Get("$.AccessLogConfig", objectRaw)
	accessLogConfig1Raw := make(map[string]interface{})
	if accessLogConfig1RawObj != nil {
		accessLogConfig1Raw = accessLogConfig1RawObj.(map[string]interface{})
	}
	if accessLogConfig1Raw["LogProject"] != nil {
		d.Set("log_project", accessLogConfig1Raw["LogProject"])
	}
	if accessLogConfig1Raw["LogStore"] != nil {
		d.Set("log_store", accessLogConfig1Raw["LogStore"])
	}

	d.Set("load_balancer_id", d.Id())

	return nil
}

func resourceAliCloudAlbLoadBalancerAccessLogConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DisableLoadBalancerAccessLog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
