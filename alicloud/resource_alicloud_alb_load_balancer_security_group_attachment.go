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

func resourceAliCloudAlbLoadBalancerSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentCreate,
		Read:   resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentRead,
		Update: resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentUpdate,
		Delete: resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentDelete,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "LoadBalancerJoinSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupIds.1"] = v
	}
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer_security_group_attachment", action, AlibabaCloudSdkGoERROR)
	}

	SecurityGroupIdsVar, _ := request["SecurityGroupIds.1"]
	d.SetId(fmt.Sprintf("%v:%v", request["LoadBalancerId"], SecurityGroupIdsVar))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.DescribeAsyncAlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentRead(d, meta)
}

func resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbLoadBalancerSecurityGroupAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer_security_group_attachment DescribeAlbLoadBalancerSecurityGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("security_group_id", objectRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("load_balancer_id", parts[0])

	return nil
}

func resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Load Balancer Security Group Attachment.")
	return nil
}

func resourceAliCloudAlbLoadBalancerSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "LoadBalancerLeaveSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SecurityGroupIds.1"] = parts[1]
	request["LoadBalancerId"] = parts[0]

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.DescribeAsyncAlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
