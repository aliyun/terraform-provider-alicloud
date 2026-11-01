package slb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Deprecated:   "Field 'entry' has been deprecated from provider version 2.1.0 and it will be removed in the future version. Please use the new field 'entries'.",
				ExactlyOneOf: []string{"entry", "entries"},
			},
			"comment": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
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
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	// batch mode: the attachment owns all entries of the acl and its ID is the acl ID
	if _, ok := d.GetOk("entries"); ok {
		return resourceAclEntryAttachmentCreateBatch(d, meta)
	}
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddAccessControlListEntry"
	request := make(map[string]interface{})
	var err error
	request = map[string]interface{}{
		"AclId":    d.Get("acl_id"),
		"RegionId": client.RegionId,
	}
	aclMaps := make([]map[string]interface{}, 0, 1)
	aclEntry := make(map[string]interface{}, 2)
	aclEntry["entry"] = d.Get("entry")
	if v, ok := d.GetOk("comment"); ok {
		aclEntry["comment"] = v
	}
	aclMaps = append(aclMaps, aclEntry)
	aclEntries, err := convert.ConvertListMapToJsonString(aclMaps)
	if err != nil {
		return errs.WrapError(err)
	}
	request["AclEntrys"] = aclEntries
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"AclEntryProcessing"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_acl_entry_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntry["entry"]))
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
	if err := addSlbAclEntries(client, aclId, entries, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}
	d.SetId(aclId)
	return resourceAclEntryAttachmentRead(d, meta)
}

func resourceAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	// batch mode: the ID is the acl ID and the attachment owns all entries of the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := slbService.ListSlbAclEntries(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				log.Printf("[DEBUG] Resource alicloud_slb_acl_entry_attachment slbService.ListSlbAclEntries Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return errs.WrapError(err)
		}
		if len(entries) == 0 {
			log.Printf("[DEBUG] Resource alicloud_slb_acl_entry_attachment batch attachment %s has no entries", d.Id())
			d.SetId("")
			return nil
		}
		entriesMaps := make([]interface{}, 0, len(entries))
		for _, entry := range entries {
			entriesMaps = append(entriesMaps, map[string]interface{}{
				"comment": entry["comment"],
				"entry":   entry["entry"],
			})
		}
		d.Set("acl_id", d.Id())
		d.Set("entries", entriesMaps)
		return nil
	}
	object, err := slbService.DescribeSlbAclEntryAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_acl_entry_attachment slbService.DescribeSlbAclEntryAttachment Failed!!! %s", err)
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
	d.Set("comment", object["AclEntryComment"])
	return nil
}

func resourceAclEntryAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// `acl_id`, `entry` and `comment` are ForceNew: in-place update only applies to `entries` (batch mode)
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
		if err := removeSlbAclEntries(client, d.Id(), entries, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	if len(create) > 0 {
		if err := addSlbAclEntries(client, d.Id(), create, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	return resourceAclEntryAttachmentRead(d, meta)
}

func resourceAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	// batch mode: remove all entries owned by this attachment from the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := slbService.ListSlbAclEntries(d.Id())
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
			entryList = append(entryList, entry["entry"])
		}
		return removeSlbAclEntries(client, d.Id(), entryList, d.Timeout(schema.TimeoutDelete))
	}
	action := "RemoveAccessControlListEntry"
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

	aclMaps := make([]map[string]interface{}, 0, 1)
	aclMaps = append(aclMaps, map[string]interface{}{
		"entry": parts[1],
	})
	aclEntries, err := convert.ConvertListMapToJsonString(aclMaps)
	if err != nil {
		return errs.WrapError(err)
	}
	request["AclEntrys"] = aclEntries
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"AclEntryProcessing"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

// addSlbAclEntries adds the given acl entries in batches of at most 50 entries
// per call, which is the limit of the AddAccessControlListEntry API.
func addSlbAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	action := "AddAccessControlListEntry"
	for _, item := range helper.SplitSlice(entries, max_num_per_time) {
		aclMaps := make([]map[string]interface{}, 0, len(item))
		for _, aclEntry := range item {
			aclEntryArg := aclEntry.(map[string]interface{})
			m := map[string]interface{}{
				"entry": aclEntryArg["entry"],
			}
			if v, ok := aclEntryArg["comment"].(string); ok && v != "" {
				m["comment"] = v
			}
			aclMaps = append(aclMaps, m)
		}
		aclEntries, err := convert.ConvertListMapToJsonString(aclMaps)
		if err != nil {
			return errs.WrapError(err)
		}
		request := map[string]interface{}{
			"AclId":     aclId,
			"RegionId":  client.RegionId,
			"AclEntrys": aclEntries,
		}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"AclEntryProcessing"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, aclId, action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

// removeSlbAclEntries removes the given acl entries in batches of at most 50
// entries per call, which is the limit of the RemoveAccessControlListEntry API.
func removeSlbAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	action := "RemoveAccessControlListEntry"
	for _, item := range helper.SplitSlice(entries, max_num_per_time) {
		aclMaps := make([]map[string]interface{}, 0, len(item))
		for _, aclEntry := range item {
			aclMaps = append(aclMaps, map[string]interface{}{
				"entry": fmt.Sprint(aclEntry),
			})
		}
		aclEntries, err := convert.ConvertListMapToJsonString(aclMaps)
		if err != nil {
			return errs.WrapError(err)
		}
		request := map[string]interface{}{
			"AclId":     aclId,
			"RegionId":  client.RegionId,
			"AclEntrys": aclEntries,
		}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"AclEntryProcessing"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			// the acl entry is already gone: keep the attachment in sync
			if errs.IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
				continue
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, aclId, action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
