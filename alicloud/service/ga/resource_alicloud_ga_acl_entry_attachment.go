package ga

import (
	"fmt"
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

func ResourceAliCloudGaAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaAclEntryAttachmentCreate,
		Read:   resourceAliCloudGaAclEntryAttachmentRead,
		Delete: resourceAliCloudGaAclEntryAttachmentDelete,
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
				ValidateFunc: validators.StringLenBetween(1, 256),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "AddEntriesToAcl"
	request := make(map[string]interface{})
	var entry string
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("AddEntriesToAcl")
	request["AclId"] = d.Get("acl_id")
	entry = d.Get("entry").(string)

	aclEntriesMaps := make([]map[string]interface{}, 0)
	aclEntriesMap := map[string]interface{}{}
	aclEntriesMap["Entry"] = entry

	if strings.Contains(entry, ":") {
		entry = strings.Replace(entry, ":", "_", -1)
	}

	if v, ok := d.GetOk("entry_description"); ok {
		aclEntriesMap["EntryDescription"] = v

	}

	aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
	request["AclEntries"] = aclEntriesMaps

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Acl", "NotExist.Acl", "ACL_NOT_STEADY"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_acl_entry_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AclId"], entry))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaAclStateRefreshFunc(fmt.Sprint(request["AclId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaAclEntryAttachmentRead(d, meta)
}

func resourceAliCloudGaAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaAclEntryAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("acl_id", parts[0])
	d.Set("entry", object["Entry"])
	d.Set("entry_description", object["EntryDescription"])

	describeGaAclObject, err := gaService.DescribeGaAcl(parts[0])
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("status", describeGaAclObject["AclStatus"])

	return nil
}

func resourceAliCloudGaAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveEntriesFromAcl"
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("RemoveEntriesFromAcl"),
		"AclId":       parts[0],
	}

	aclEntriesMaps := make([]map[string]interface{}, 0)
	aclEntriesMap := map[string]interface{}{}

	if strings.Contains(fmt.Sprint(parts[1]), "_") {
		parts[1] = strings.Replace(fmt.Sprint(parts[1]), "_", ":", -1)
	}

	aclEntriesMap["Entry"] = parts[1]
	aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
	request["AclEntries"] = aclEntriesMaps

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Acl", "ACL_NOT_STEADY"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
