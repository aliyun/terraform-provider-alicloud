// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package nlb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLoadBalancerSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerSecurityGroupAttachmentCreate,
		Read:   resourceLoadBalancerSecurityGroupAttachmentRead,
		Update: resourceLoadBalancerSecurityGroupAttachmentUpdate,
		Delete: resourceLoadBalancerSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoadBalancerSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "LoadBalancerJoinSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["SecurityGroupIds.1"] = d.Get("security_group_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock", "OperationFailed.ResourceIsConfiguring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_nlb_load_balancer_security_group_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["LoadBalancerId"], request["SecurityGroupIds.1"]))

	nlbServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceLoadBalancerSecurityGroupAttachmentRead(d, meta)
}

func resourceLoadBalancerSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NewService(client)

	_, err := nlbServiceV2.DescribeNlbLoadBalancerSecurityGroupAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer_security_group_attachment DescribeNlbLoadBalancerSecurityGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("load_balancer_id", parts[0])
	d.Set("security_group_id", parts[1])

	return nil
}

func resourceLoadBalancerSecurityGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Load Balancer Security Group Attachment.")
	return nil
}

func resourceLoadBalancerSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "LoadBalancerLeaveSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = parts[0]
	request["SecurityGroupIds.1"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock"}) || helper.NeedRetry(err) {
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

	nlbServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

func convertNlbSecurityGroupIdsResponse(source interface{}) interface{} {
	switch source {
	}
	return source
}
