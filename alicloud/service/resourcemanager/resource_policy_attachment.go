package resourcemanager

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePolicyAttachmentCreate,
		Read:   ResourcePolicyAttachmentRead,
		Delete: ResourcePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"principal_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMSGroup", "IMSUser", "ServiceRole"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func ResourcePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPolicy"
	request := make(map[string]interface{})
	var err error
	request["PolicyName"] = d.Get("policy_name")
	request["PolicyType"] = d.Get("policy_type")
	request["PrincipalName"] = d.Get("principal_name")
	request["PrincipalType"] = d.Get("principal_type")
	request["ResourceGroupId"] = d.Get("resource_group_id")
	wait := helper.IncrementalWait(30*time.Second, 30*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_policy_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["PolicyName"], ":", request["PolicyType"], ":", request["PrincipalName"], ":", request["PrincipalType"], ":", request["ResourceGroupId"]))

	return ResourcePolicyAttachmentRead(d, meta)
}
func ResourcePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	_, err := resourcemanagerService.DescribeResourceManagerPolicyAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_attachment resourcemanagerService.DescribeResourceManagerPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 5)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("policy_type", parts[1])
	d.Set("principal_name", parts[2])
	d.Set("principal_type", parts[3])
	d.Set("resource_group_id", parts[4])
	return nil
}
func ResourcePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	parts, err := helper.ParseResourceId(d.Id(), 5)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DetachPolicy"
	var response map[string]interface{}
	request := map[string]interface{}{
		"PolicyName":      parts[0],
		"PolicyType":      parts[1],
		"PrincipalName":   parts[2],
		"PrincipalType":   parts[3],
		"ResourceGroupId": parts[4],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
