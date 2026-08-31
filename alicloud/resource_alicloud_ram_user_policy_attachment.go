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

func resourceAliCloudRamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamUserPolicyAttachmentCreate,
		Read:   resourceAliCloudRamUserPolicyAttachmentRead,
		Delete: resourceAliCloudRamUserPolicyAttachmentDelete,
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
				ValidateFunc: StringInSlice([]string{"System", "Custom"}, false),
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_user_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", "user", request["PolicyName"], request["PolicyType"], request["UserName"]))

	return resourceAliCloudRamUserPolicyAttachmentRead(d, meta)
}

func resourceAliCloudRamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>
	if split := strings.Split(d.Id(), ":"); len(split) != 4 {
		id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, COLON_SEPARATED)
		d.SetId(id)
	}
	objectRaw, err := ramServiceV2.DescribeRamUserPolicyAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_user_policy_attachment DescribeRamUserPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	d.Set("policy_name", objectRaw["PolicyName"])
	d.Set("policy_type", objectRaw["PolicyType"])
	d.Set("user_name", parts[3])

	return nil
}

func resourceAliCloudRamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DetachPolicyFromUser"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>
	id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, COLON_SEPARATED)

	if d.Id() != id {
		d.SetId(id)
	}

	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}

	request["PolicyType"] = parts[2]
	request["PolicyName"] = parts[1]
	request["UserName"] = parts[3]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.User", "EntityNotExist.User.Policy"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
