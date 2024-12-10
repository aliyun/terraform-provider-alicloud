package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudFirewallInstanceMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallInstanceMemberCreate,
		Read:   resourceAlicloudCloudFirewallInstanceMemberRead,
		Update: resourceAlicloudCloudFirewallInstanceMemberUpdate,
		Delete: resourceAlicloudCloudFirewallInstanceMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"member_desc": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"member_display_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"member_uid": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"modify_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCloudFirewallInstanceMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	request := make(map[string]interface{})
	var err error
	var endpoint string

	if v, ok := d.GetOk("member_desc"); ok {
		request["Members.1.MemberDesc"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["Members.1.MemberUid"] = v
	}

	var response map[string]interface{}
	action := "AddInstanceMembers"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_instance_member", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("member_uid")))

	stateConf := BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudfwService.CloudFirewallInstanceMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCloudFirewallInstanceMemberRead(d, meta)
}

func resourceAlicloudCloudFirewallInstanceMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallInstanceMember(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance_member cloudfwService.DescribeCloudFirewallInstanceMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("member_uid", object["MemberUid"])
	d.Set("create_time", object["CreateTime"])
	d.Set("member_desc", object["MemberDesc"])
	d.Set("member_display_name", object["MemberDisplayName"])
	d.Set("modify_time", object["ModifyTime"])
	d.Set("status", object["MemberStatus"])

	return nil
}

func resourceAlicloudCloudFirewallInstanceMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	var endpoint string
	update := false
	request := map[string]interface{}{
		"Members.1.MemberUid": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("member_desc") {
		update = true
	}
	request["Members.1.MemberDesc"] = d.Get("member_desc")

	if update {
		action := "ModifyInstanceMemberAttributes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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
	}

	return resourceAlicloudCloudFirewallInstanceMemberRead(d, meta)
}

func resourceAlicloudCloudFirewallInstanceMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	var err error
	var endpoint string

	request := map[string]interface{}{
		"MemberUids": []string{d.Id()},
	}

	action := "DeleteInstanceMembers"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, cloudfwService.CloudFirewallInstanceMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
