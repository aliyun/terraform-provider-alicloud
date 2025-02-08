// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNlbLoadBalancerSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentCreate,
		Read:   resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentRead,
		Update: resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentUpdate,
		Delete: resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentDelete,
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

func resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

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
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer_security_group_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["LoadBalancerId"], request["SecurityGroupIds.1"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentRead(d, meta)
}

func resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	_, err := nlbServiceV2.DescribeNlbLoadBalancerSecurityGroupAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer_security_group_attachment DescribeNlbLoadBalancerSecurityGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("load_balancer_id", parts[0])
	d.Set("security_group_id", parts[1])

	return nil
}

func resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Load Balancer Security Group Attachment.")
	return nil
}

func resourceAliCloudNlbLoadBalancerSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

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

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
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

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertNlbSecurityGroupIdsResponse(source interface{}) interface{} {
	switch source {
	}
	return source
}
