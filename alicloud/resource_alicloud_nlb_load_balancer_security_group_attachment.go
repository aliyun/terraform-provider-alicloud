package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudNlbLoadBalancerSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentCreate,
		Read:   resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentRead,
		Update: resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentUpdate,
		Delete: resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("LoadBalancerJoinSecurityGroup")
	request["LoadBalancerId"] = d.Get("load_balancer_id")

	securityGroupId := d.Get("security_group_id").(string)
	request["SecurityGroupIds"] = []string{securityGroupId}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	var response map[string]interface{}
	action := "LoadBalancerJoinSecurityGroup"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict.Lock"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer_security_group_attachment", action, AlibabaCloudSdkGoERROR)
	}

	nlbService := NlbService{client}

	jobId := fmt.Sprint(response["JobId"])
	taskConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(jobId, []string{}))
	if _, err := taskConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"], ":", securityGroupId))

	return resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentRead(d, meta)
}

func resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}

	_, err := nlbService.DescribeNlbLoadBalancerSecurityGroupAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer_security_group_attachment nlbService.DescribeNlbLoadBalancerSecurityGroupAttachment Failed!!! %s", err)
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

func resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentRead(d, meta)
}

func resourceAlicloudNlbLoadBalancerSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ClientToken":        buildClientToken("LoadBalancerLeaveSecurityGroup"),
		"LoadBalancerId":     parts[0],
		"SecurityGroupIds.1": parts[1],
	}

	action := "LoadBalancerLeaveSecurityGroup"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict.Lock"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
