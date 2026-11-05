// Package alb. This file is generated automatically. Please do not modify it manually, thank you!
package alb

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLoadBalancerAccessLogConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerAccessLogConfigAttachmentCreate,
		Read:   resourceLoadBalancerAccessLogConfigAttachmentRead,
		Delete: resourceLoadBalancerAccessLogConfigAttachmentDelete,
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

func resourceLoadBalancerAccessLogConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnableLoadBalancerAccessLog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["LogProject"] = d.Get("log_project")
	request["LogStore"] = d.Get("log_store")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_load_balancer_access_log_config_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutCreate), 0, albServiceV2.DescribeAsyncAlbLoadBalancerAccessLogConfigAttachmentStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return resourceLoadBalancerAccessLogConfigAttachmentRead(d, meta)
}

func resourceLoadBalancerAccessLogConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbLoadBalancerAccessLogConfigAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer_access_log_config_attachment DescribeAlbLoadBalancerAccessLogConfigAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

func resourceLoadBalancerAccessLogConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DisableLoadBalancerAccessLog"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
