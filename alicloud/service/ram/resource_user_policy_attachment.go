package ram

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceUserPolicyAttachmentCreate,
		Read:   ResourceUserPolicyAttachmentRead,
		Delete: ResourceUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				ValidateFunc: validators.StringInSlice([]string{"System", "Custom"}, false),
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AttachPolicyToUser"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["PolicyName"] = d.Get("policy_name")
	request["UserName"] = d.Get("user_name")
	request["PolicyType"] = d.Get("policy_type")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_user_policy_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", "user", request["PolicyName"], request["PolicyType"], request["UserName"]))

	return ResourceUserPolicyAttachmentRead(d, meta)
}

func ResourceUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>
	if split := strings.Split(d.Id(), ":"); len(split) != 4 {
		id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, colonSeparated)
		d.SetId(id)
	}
	objectRaw, err := ramServiceV2.DescribeRamUserPolicyAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_user_policy_attachment DescribeRamUserPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("policy_name", objectRaw["PolicyName"])
	d.Set("policy_type", objectRaw["PolicyType"])
	d.Set("user_name", parts[3])

	return nil
}

func ResourceUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DetachPolicyFromUser"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>
	id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, colonSeparated)

	if d.Id() != id {
		d.SetId(id)
	}

	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request["PolicyType"] = parts[2]
	request["PolicyName"] = parts[1]
	request["UserName"] = parts[3]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.User", "EntityNotExist.User.Policy"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
