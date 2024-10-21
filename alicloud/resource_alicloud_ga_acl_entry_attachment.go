package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudGaAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaAclEntryAttachmentCreate,
		Read:   resourceAlicloudGaAclEntryAttachmentRead,
		Delete: resourceAlicloudGaAclEntryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "AddEntriesToAcl"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("AddEntriesToAcl")
	request["AclId"] = d.Get("acl_id")

	aclEntriesMaps := make([]map[string]interface{}, 0)
	aclEntriesMap := map[string]interface{}{}
	aclEntriesMap["Entry"] = d.Get("entry")
	if v, ok := d.GetOk("entry_description"); ok {
		aclEntriesMap["EntryDescription"] = v

	}
	aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
	request["AclEntries"] = aclEntriesMaps

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Acl", "NotExist.Acl", "ACL_NOT_STEADY"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_acl_entry_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AclId"], aclEntriesMap["Entry"]))

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaAclStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaAclEntryAttachmentRead(d, meta)
}

func resourceAlicloudGaAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaAclEntryAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("acl_id", parts[0])
	d.Set("entry", parts[1])
	d.Set("entry_description", object["EntryDescription"])

	describeGaAclObject, err := gaService.DescribeGaAcl(parts[0])
	if err != nil {
		return WrapError(err)
	}

	d.Set("status", describeGaAclObject["AclStatus"])

	return nil
}

func resourceAlicloudGaAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveEntriesFromAcl"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("RemoveEntriesFromAcl"),
		"AclId":       parts[0],
	}

	aclEntriesMaps := make([]map[string]interface{}, 0)
	aclEntriesMap := map[string]interface{}{}
	aclEntriesMap["Entry"] = parts[1]
	aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
	request["AclEntries"] = aclEntriesMaps

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Acl", "ACL_NOT_STEADY"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
