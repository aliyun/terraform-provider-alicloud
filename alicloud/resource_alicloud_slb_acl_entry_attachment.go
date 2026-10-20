package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudSlbAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclEntryAttachmentCreate,
		Read:   resourceAlicloudSlbAclEntryAttachmentRead,
		Update: resourceAlicloudSlbAclEntryAttachmentUpdate,
		Delete: resourceAlicloudSlbAclEntryAttachmentDelete,
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

func resourceAlicloudSlbAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	// batch mode: the attachment owns all entries of the acl and its ID is the acl ID
	if _, ok := d.GetOk("entries"); ok {
		return resourceAlicloudSlbAclEntryAttachmentCreateBatch(d, meta)
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
	aclEntries, err := convertListMapToJsonString(aclMaps)
	if err != nil {
		return WrapError(err)
	}
	request["AclEntrys"] = aclEntries
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_acl_entry_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntry["entry"]))
	return resourceAlicloudSlbAclEntryAttachmentRead(d, meta)
}

func resourceAlicloudSlbAclEntryAttachmentCreateBatch(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aclId := d.Get("acl_id").(string)
	entries := d.Get("entries").(*schema.Set).List()
	if len(entries) == 0 {
		// the exactly-one-of constraint rejects zero entry blocks at plan time;
		// keep an explicit error here in case an empty set reaches the apply path anyway
		return WrapError(Error("at least one entry block is required, to remove all entries please remove the resource"))
	}
	if err := addSlbAclEntries(client, aclId, entries, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}
	d.SetId(aclId)
	return resourceAlicloudSlbAclEntryAttachmentRead(d, meta)
}

func resourceAlicloudSlbAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	// batch mode: the ID is the acl ID and the attachment owns all entries of the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := slbService.ListSlbAclEntries(d.Id())
		if err != nil {
			if NotFoundError(err) {
				log.Printf("[DEBUG] Resource alicloud_slb_acl_entry_attachment slbService.ListSlbAclEntries Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return WrapError(err)
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
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_acl_entry_attachment slbService.DescribeSlbAclEntryAttachment Failed!!! %s", err)
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
	d.Set("comment", object["AclEntryComment"])
	return nil
}

func resourceAlicloudSlbAclEntryAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
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
	return resourceAlicloudSlbAclEntryAttachmentRead(d, meta)
}

func resourceAlicloudSlbAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	// batch mode: remove all entries owned by this attachment from the acl
	if !strings.Contains(d.Id(), ":") {
		entries, err := slbService.ListSlbAclEntries(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
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

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId":    parts[0],
		"RegionId": client.RegionId,
	}

	aclMaps := make([]map[string]interface{}, 0, 1)
	aclMaps = append(aclMaps, map[string]interface{}{
		"entry": parts[1],
	})
	aclEntries, err := convertListMapToJsonString(aclMaps)
	if err != nil {
		return WrapError(err)
	}
	request["AclEntrys"] = aclEntries
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// addSlbAclEntries adds the given acl entries in batches of at most 50 entries
// per call, which is the limit of the AddAccessControlListEntry API.
func addSlbAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	action := "AddAccessControlListEntry"
	for _, item := range SplitSlice(entries, max_num_per_time) {
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
		aclEntries, err := convertListMapToJsonString(aclMaps)
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"AclId":     aclId,
			"RegionId":  client.RegionId,
			"AclEntrys": aclEntries,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, aclId, action, AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

// removeSlbAclEntries removes the given acl entries in batches of at most 50
// entries per call, which is the limit of the RemoveAccessControlListEntry API.
func removeSlbAclEntries(client *connectivity.AliyunClient, aclId string, entries []interface{}, timeout time.Duration) error {
	action := "RemoveAccessControlListEntry"
	for _, item := range SplitSlice(entries, max_num_per_time) {
		aclMaps := make([]map[string]interface{}, 0, len(item))
		for _, aclEntry := range item {
			aclMaps = append(aclMaps, map[string]interface{}{
				"entry": fmt.Sprint(aclEntry),
			})
		}
		aclEntries, err := convertListMapToJsonString(aclMaps)
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"AclId":     aclId,
			"RegionId":  client.RegionId,
			"AclEntrys": aclEntries,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			// the acl entry is already gone: keep the attachment in sync
			if IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
				continue
			}
			return WrapErrorf(err, DefaultErrorMsg, aclId, action, AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
