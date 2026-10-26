package alb

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

func ResourceAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAclEntryAttachmentCreate,
		Read:   resourceAclEntryAttachmentRead,
		Update: resourceAclEntryAttachmentUpdate,
		Delete: resourceAclEntryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Deprecated:   "Field 'entry' has been deprecated from provider version 1.292.0 and it will be removed in the future version. Please use the new field 'entries'.",
				ExactlyOneOf: []string{"entry", "entries"},
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(1, 256),
				ConflictsWith: []string{"entries"},
			},
			"entries": {
				Type:         schema.TypeSet,
				Optional:     true,
				ExactlyOneOf: []string{"entry", "entries"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// batch mode: the attachment owns all entries of the acl and its ID is the acl ID
	if _, ok := d.GetOk("entries"); ok {
		return resourceAclEntryAttachmentCreateBatch(d, meta)
	}
	albService := NewService(client)
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
	request["ClientToken"] = helper.BuildClientToken("AddEntriesToAcl")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "IncorrectStatus.Acl", "ResourceInConfiguring", "AclStatusNotStable"}) || helper.NeedRetry(err) {
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

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntriesMap["Entry"]))
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AclEntryAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	// The acl entry turns Available before the acl itself leaves Configuring status.
	// Wait for the acl to become Available, so that dependent resources (e.g. alicloud_alb_listener_acl_attachment)
	// will not be created while the server side still considers the acl Configuring.
	stateConf = helper.BuildStateConf([]string{"Creating", "Configuring"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AclStateRefreshFunc(fmt.Sprint(request["AclId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAclEntryAttachmentRead(d, meta)
}

func resourceAclEntryAttachmentCreateBatch(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aclId := d.Get("acl_id").(string)
	entries := d.Get("entries").(*schema.Set).List()
	if len(entries) == 0 {
		// the exactly-one-of constraint rejects zero entry blocks at plan time;
		// keep an explicit error here in case an empty set reaches the apply path anyway
		return errs.WrapError(errs.Error("at least one entry block is required, to remove all entries please remove the resource"))
	}
	if err := addAclEntries(client, aclId, entries, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}
	d.SetId(aclId)
	return resourceAclEntryAttachmentRead(d, meta)
}

func resourceAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)
	// batch mode: the ID is the acl ID and the attachment owns all entries of the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := albService.ListAclEntries(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				log.Printf("[DEBUG] Resource alicloud_alb_acl_entry_attachment albService.ListAclEntries Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return errs.WrapError(err)
		}
		if len(entries) == 0 {
			log.Printf("[DEBUG] Resource alicloud_alb_acl_entry_attachment batch attachment %s has no entries", d.Id())
			d.SetId("")
			return nil
		}
		entriesMaps := make([]interface{}, 0, len(entries))
		for _, entry := range entries {
			entriesMaps = append(entriesMaps, map[string]interface{}{
				"description": entry["Description"],
				"entry":       entry["Entry"],
				"status":      entry["Status"],
			})
		}
		d.Set("acl_id", d.Id())
		d.Set("entries", entriesMaps)
		return nil
	}
	object, err := albService.DescribeAclEntryAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_acl_entry_attachment albService.DescribeAclEntryAttachment Failed!!! %s", err)
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
	d.Set("entry", parts[1])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])
	return nil
}

func resourceAclEntryAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// `acl_id`, `entry` and `description` are ForceNew: in-place update only applies to `entries` (batch mode)
	if !d.HasChange("entries") {
		return nil
	}
	oraw, nraw := d.GetChange("entries")
	remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
	create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

	if len(remove) > 0 {
		entries := make([]interface{}, 0, len(remove))
		for _, item := range remove {
			entries = append(entries, item.(map[string]interface{})["entry"])
		}
		if err := removeAclEntries(client, d.Id(), entries, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	if len(create) > 0 {
		if err := addAclEntries(client, d.Id(), create, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	return resourceAclEntryAttachmentRead(d, meta)
}

func resourceAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)
	// batch mode: remove all entries owned by this attachment from the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := albService.ListAclEntries(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapError(err)
		}
		if len(entries) == 0 {
			return nil
		}
		entryList := make([]interface{}, 0, len(entries))
		for _, entry := range entries {
			entryList = append(entryList, entry["Entry"])
		}
		return removeAclEntries(client, d.Id(), entryList, d.Timeout(schema.TimeoutDelete))
	}
	action := "RemoveEntriesFromAcl"
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"AclId":    parts[0],
		"RegionId": client.RegionId,
	}

	aclEntriesMaps := make([]string, 0)
	aclEntriesMaps = append(aclEntriesMaps, parts[1])
	request["Entries"] = aclEntriesMaps
	request["ClientToken"] = helper.BuildClientToken("RemoveEntriesFromAcl")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Acl", "OperationFailed.ResourceGroupStatusCheckFail", "ResourceInConfiguring", "AclStatusNotStable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AclEntry"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AclEntryAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

// addAclEntries adds the given acl entries in batches of at most 20 entries
// per call, which is the limit of the AddEntriesToAcl API.
func addAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	albService := NewService(client)
	action := "AddEntriesToAcl"
	for _, item := range helper.SplitSlice(entries, 20) {
		aclEntriesMaps := make([]map[string]interface{}, 0)
		for _, aclEntry := range item {
			aclEntryArg := aclEntry.(map[string]interface{})
			aclEntriesMap := map[string]interface{}{
				"Entry": aclEntryArg["entry"],
			}
			if v, ok := aclEntryArg["description"].(string); ok && v != "" {
				aclEntriesMap["Description"] = v
			}
			aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
		}
		request := map[string]interface{}{
			"AclId":       aclId,
			"AclEntries":  aclEntriesMaps,
			"ClientToken": helper.BuildClientToken(action),
		}
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "IncorrectStatus.Acl", "ResourceInConfiguring", "AclStatusNotStable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			// an entry already on the acl is kept as it is: it either belongs to this
			// full-set attachment or will converge on the next apply
			if errs.IsExpectedErrors(err, []string{"ResourceAlreadyExist.AclEntry"}) {
				continue
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, aclId, action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, timeout, 5*time.Second, albService.AclStateRefreshFunc(aclId, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, aclId)
		}
	}
	return nil
}

// removeAclEntries removes the given acl entries in batches of at most 20
// entries per call, which is the limit of the RemoveEntriesFromAcl API.
func removeAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	albService := NewService(client)
	action := "RemoveEntriesFromAcl"
	for _, item := range helper.SplitSlice(entries, 20) {
		aclEntriesMaps := make([]string, 0)
		for _, aclEntry := range item {
			aclEntriesMaps = append(aclEntriesMaps, fmt.Sprint(aclEntry))
		}
		request := map[string]interface{}{
			"AclId":       aclId,
			"RegionId":    client.RegionId,
			"Entries":     aclEntriesMaps,
			"ClientToken": helper.BuildClientToken(action),
		}
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Acl", "OperationFailed.ResourceGroupStatusCheckFail", "ResourceInConfiguring", "AclStatusNotStable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AclEntry"}) {
				continue
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, aclId, action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, timeout, 5*time.Second, albService.AclStateRefreshFunc(aclId, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, aclId)
		}
	}
	return nil
}
