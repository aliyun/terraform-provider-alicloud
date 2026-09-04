package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// alicloud_threat_detection_rd_default_sync_list is a singleton account-level
// resource that manages the default synchronization list of resource directory
// folders for Threat Detection (云安全中心). There is no dedicated Get/Update/Delete
// API: CreateRdDefaultSyncList applies the whole folder list (set/replace semantics,
// an empty value clears the existing list), and ListRdDefaultSyncList reads it back.
// The resource id is the Alibaba Cloud account id.
func resourceAliCloudThreatDetectionRdDefaultSyncList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionRdDefaultSyncListCreate,
		Read:   resourceAliCloudThreatDetectionRdDefaultSyncListRead,
		Update: resourceAliCloudThreatDetectionRdDefaultSyncListUpdate,
		Delete: resourceAliCloudThreatDetectionRdDefaultSyncListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"folder_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudThreatDetectionRdDefaultSyncListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateRdDefaultSyncList"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	// FolderIds is a query parameter that accepts a comma-separated list of
	// resource directory folder ids. CreateRdDefaultSyncList is set/replace:
	// the whole folder list is replaced on every call, and an empty value
	// clears the existing synchronized folders. Only send FolderIds when the
	// user declares it in the configuration, so that adopting an account that
	// already holds a synchronization list (for example via import without
	// folder_ids) does not wipe it out.
	if d.HasChange("folder_ids") {
		request["FolderIds"] = strings.Join(expandStringList(d.Get("folder_ids").([]interface{})), COMMA_SEPARATED)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_rd_default_sync_list", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(accountId))

	return resourceAliCloudThreatDetectionRdDefaultSyncListRead(d, meta)
}

func resourceAliCloudThreatDetectionRdDefaultSyncListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	object, err := threatDetectionServiceV2.DescribeThreatDetectionRdDefaultSyncList(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_rd_default_sync_list DescribeThreatDetectionRdDefaultSyncList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	folderIds := make([]interface{}, 0)
	if v, ok := object["FolderIds"]; ok && v != nil {
		if str, ok := v.(string); ok {
			str = strings.TrimSpace(str)
			if str != "" {
				for _, id := range strings.Split(str, COMMA_SEPARATED) {
					id = strings.TrimSpace(id)
					if id != "" {
						folderIds = append(folderIds, id)
					}
				}
			}
		}
	}
	if err := d.Set("folder_ids", folderIds); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudThreatDetectionRdDefaultSyncListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateRdDefaultSyncList"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	// CreateRdDefaultSyncList is set/replace: re-applying the whole folder list
	// with the new value replaces the previously synchronized folders. Skip the
	// API call when the folder list is unchanged.
	if !d.HasChange("folder_ids") {
		return resourceAliCloudThreatDetectionRdDefaultSyncListRead(d, meta)
	}
	request["FolderIds"] = strings.Join(expandStringList(d.Get("folder_ids").([]interface{})), COMMA_SEPARATED)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudThreatDetectionRdDefaultSyncListRead(d, meta)
}

func resourceAliCloudThreatDetectionRdDefaultSyncListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateRdDefaultSyncList"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	// No dedicated Delete API: clearing the folder list (empty FolderIds) is
	// equivalent to disabling the default synchronization.
	request["FolderIds"] = ""

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
