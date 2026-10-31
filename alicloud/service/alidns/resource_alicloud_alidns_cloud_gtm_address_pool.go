// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alidns

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudAlidnsCloudGtmAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsCloudGtmAddressPoolCreate,
		Read:   resourceAliCloudAlidnsCloudGtmAddressPoolRead,
		Update: resourceAliCloudAlidnsCloudGtmAddressPoolUpdate,
		Delete: resourceAliCloudAlidnsCloudGtmAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_lb_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"round_robin", "sequence", "weight", "source_nearest"}, false),
			},
			"address_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_pool_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"IPv4", "IPv6", "domain"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"enable", "disable"}, false),
			},
			"health_judgement": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"any_ok", "p30_ok", "p50_ok", "p70_ok", "all_ok"}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence_lb_strategy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"preemptive", "nonPreemptive"}, false),
			},
		},
	}
}

func resourceAliCloudAlidnsCloudGtmAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCloudGtmAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

	request["EnableStatus"] = d.Get("enable_status")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}
	request["AddressPoolName"] = d.Get("address_pool_name")
	request["AddressPoolType"] = d.Get("address_pool_type")
	request["HealthJudgement"] = d.Get("health_judgement")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alidns_cloud_gtm_address_pool", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AddressPoolId"]))

	return resourceAliCloudAlidnsCloudGtmAddressPoolUpdate(d, meta)
}

func resourceAliCloudAlidnsCloudGtmAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsServiceV2 := NewAlidnsServiceV2(client)

	objectRaw, err := alidnsServiceV2.DescribeAlidnsCloudGtmAddressPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_cloud_gtm_address_pool DescribeAlidnsCloudGtmAddressPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("address_lb_strategy", objectRaw["AddressLbStrategy"])
	d.Set("address_pool_name", objectRaw["AddressPoolName"])
	d.Set("address_pool_type", objectRaw["AddressPoolType"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_status", objectRaw["EnableStatus"])
	d.Set("health_judgement", objectRaw["HealthJudgement"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("sequence_lb_strategy_mode", objectRaw["SequenceLbStrategyMode"])

	return nil
}

func resourceAliCloudAlidnsCloudGtmAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateCloudGtmAddressPoolBasicConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressPoolId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("address_pool_name") {
		update = true
	}
	request["AddressPoolName"] = d.Get("address_pool_name")
	if !d.IsNewResource() && d.HasChange("health_judgement") {
		update = true
	}
	request["HealthJudgement"] = d.Get("health_judgement")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateCloudGtmAddressPoolEnableStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressPoolId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("enable_status") {
		update = true
	}
	request["EnableStatus"] = d.Get("enable_status")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateCloudGtmAddressPoolLbStrategy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressPoolId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("sequence_lb_strategy_mode") {
		update = true
		request["SequenceLbStrategyMode"] = d.Get("sequence_lb_strategy_mode")
	}

	if d.HasChange("address_lb_strategy") {
		update = true
	}
	request["AddressLbStrategy"] = d.Get("address_lb_strategy")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateCloudGtmAddressPoolRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressPoolId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudAlidnsCloudGtmAddressPoolRead(d, meta)
}

func resourceAliCloudAlidnsCloudGtmAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCloudGtmAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AddressPoolId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
