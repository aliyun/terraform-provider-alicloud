package alicloud

import (
	"fmt"
	"log"

	util "github.com/alibabacloud-go/tea-utils/service"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclEntryAttachmentCreate,
		Read:   resourceAlicloudSlbAclEntryAttachmentRead,
		Delete: resourceAlicloudSlbAclEntryAttachmentDelete,
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
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSlbAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddAccessControlListEntry"
	request := make(map[string]interface{})
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = map[string]interface{}{
		"AclId":    d.Id(),
		"RegionId": client.RegionId,
	}
	request["AclId"] = d.Get("acl_id")
	aclMaps := make([]map[string]interface{}, 1)
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
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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

func resourceAlicloudSlbAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
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

func resourceAlicloudSlbAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveAccessControlListEntry"
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId":    parts[0],
		"RegionId": client.RegionId,
	}

	aclMaps := make([]map[string]interface{}, 1)
	aclMaps = append(aclMaps, map[string]interface{}{
		"entry": parts[1],
	})
	aclEntries, err := convertListMapToJsonString(aclMaps)
	if err != nil {
		return WrapError(err)
	}
	request["AclEntrys"] = aclEntries
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"AclEntryProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
