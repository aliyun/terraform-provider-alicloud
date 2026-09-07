package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlbAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbAclEntryAttachmentCreate,
		Read:   resourceAlicloudAlbAclEntryAttachmentRead,
		Delete: resourceAlicloudAlbAclEntryAttachmentDelete,
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
			"description": {
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

func resourceAlicloudAlbAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	var err error
	action := "AddEntriesToAcl"
	request := map[string]interface{}{
		"AclId": d.Get("acl_id"),
	}
	aclEntriesMaps := make([]map[string]interface{}, 0)
	aclEntriesMap := map[string]interface{}{}
	if v, ok := d.GetOk("description"); ok {
		aclEntriesMap["Description"] = v

	}
	aclEntriesMap["Entry"] = d.Get("entry")
	aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
	request["AclEntries"] = aclEntriesMaps
	request["ClientToken"] = buildClientToken("AddEntriesToAcl")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "IncorrectStatus.Acl", "ResourceInConfiguring"}) || NeedRetry(err) {
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

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntriesMap["Entry"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbAclEntryAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// Wait for the ACL itself to settle to Available before returning. The
	// entry-level status may flip to Available while the ACL is still
	// Configuring, which breaks downstream operations such as associating the
	// ACL with a listener (IncorrectStatus.Acl). Aligns with the deprecated
	// acl_entries path in resourceAlicloudAlbAclUpdate.
	aclStateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbAclStateRefreshFunc(d.Get("acl_id").(string), []string{}))
	if _, err := aclStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudAlbAclEntryAttachmentRead(d, meta)
}

func resourceAlicloudAlbAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbAclEntryAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_acl_entry_attachment AlbService.DescribeAlbAclEntryAttachment Failed!!! %s", err)
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
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])
	return nil
}

func resourceAlicloudAlbAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	action := "RemoveEntriesFromAcl"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId":    parts[0],
		"RegionId": client.RegionId,
	}

	aclEntriesMaps := make([]string, 0)
	aclEntriesMaps = append(aclEntriesMaps, parts[1])
	request["Entries"] = aclEntriesMaps
	request["ClientToken"] = buildClientToken("RemoveEntriesFromAcl")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Acl", "OperationFailed.ResourceGroupStatusCheckFail", "ResourceInConfiguring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.AclEntry"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbAclEntryAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// Wait for the ACL itself to settle to Available before returning. The
	// entry is gone once the remove call succeeds, but the ACL may still be
	// Configuring, which breaks downstream operations on the same ACL.
	aclStateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbAclStateRefreshFunc(parts[0], []string{}))
	if _, err := aclStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
