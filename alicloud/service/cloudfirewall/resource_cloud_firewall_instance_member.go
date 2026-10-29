package cloudfirewall

import (
	"fmt"
	"log"
	"time"

	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCloudFirewallInstanceMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceMemberCreate,
		Read:   resourceInstanceMemberRead,
		Update: resourceInstanceMemberUpdate,
		Delete: resourceInstanceMemberDelete,
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

func resourceInstanceMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := NewService(client)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_instance_member", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("member_uid")))

	stateConf := helper.BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudfwService.CloudFirewallInstanceMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceInstanceMemberRead(d, meta)
}

func resourceInstanceMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := NewService(client)

	object, err := cloudfwService.DescribeCloudFirewallInstanceMember(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance_member cloudfwService.DescribeCloudFirewallInstanceMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("member_uid", object["MemberUid"])
	d.Set("create_time", object["CreateTime"])
	d.Set("member_desc", object["MemberDesc"])
	d.Set("member_display_name", object["MemberDisplayName"])
	d.Set("modify_time", object["ModifyTime"])
	d.Set("status", object["MemberStatus"])

	return nil
}

func resourceInstanceMemberUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err := client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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
	}

	return resourceInstanceMemberRead(d, meta)
}

func resourceInstanceMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := NewService(client)
	var err error
	var endpoint string

	request := map[string]interface{}{
		"MemberUids": []string{d.Id()},
	}

	action := "DeleteInstanceMembers"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err := client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, cloudfwService.CloudFirewallInstanceMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
