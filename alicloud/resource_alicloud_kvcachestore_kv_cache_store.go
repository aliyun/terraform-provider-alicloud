// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKvcachestoreKvCacheStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKvcachestoreKvCacheStoreCreate,
		Read:   resourceAliCloudKvcachestoreKvCacheStoreRead,
		Update: resourceAliCloudKvcachestoreKvCacheStoreUpdate,
		Delete: resourceAliCloudKvcachestoreKvCacheStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hpn_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudKvcachestoreKvCacheStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tag"); ok {
		tagMaps := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			tagMap := make(map[string]interface{})
			if tagKey, ok := itemMap["tag_key"].(string); ok && tagKey != "" {
				tagMap["TagKey"] = tagKey
			}
			if tagValue, ok := itemMap["tag_value"].(string); ok && tagValue != "" {
				tagMap["TagValue"] = tagValue
			}
			tagMaps = append(tagMaps, tagMap)
		}
		if len(tagMaps) > 0 {
			request["Tag"] = tagMaps
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PaymentType"] = v
	}
	request["HpnZone"] = d.Get("hpn_zone")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Capacity"] = d.Get("capacity")
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvcachestore_kv_cache_store", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["KvcsId"]))

	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, kvcachestoreServiceV2.KvcachestoreKvCacheStoreStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudKvcachestoreKvCacheStoreRead(d, meta)
}

func resourceAliCloudKvcachestoreKvCacheStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}

	objectRaw, err := kvcachestoreServiceV2.DescribeKvcachestoreKvCacheStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvcachestore_kv_cache_store DescribeKvcachestoreKvCacheStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("capacity", objectRaw["Capacity"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("name", objectRaw["Name"])
	d.Set("payment_type", objectRaw["PaymentType"])
	// NOTE: the read API never returns ResourceGroupId, so the configured value
	// is retained in state; it is listed in ImportStateVerifyIgnore accordingly.
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	if tagsRaw, ok := objectRaw["Tags"]; ok && tagsRaw != nil {
		tagMaps := make([]map[string]interface{}, 0)
		for _, tagsChildRaw := range convertToInterfaceArray(tagsRaw) {
			if tagsChild, ok := tagsChildRaw.(map[string]interface{}); ok {
				tagMap := make(map[string]interface{})
				tagMap["tag_key"] = tagsChild["TagKey"]
				tagMap["tag_value"] = tagsChild["TagValue"]
				tagMaps = append(tagMaps, tagMap)
			}
		}
		d.Set("tag", tagMaps)
	}

	return nil
}

func resourceAliCloudKvcachestoreKvCacheStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error

	// Capacity changes are processed by the backend through an asynchronous order
	// workflow; submit them in a dedicated request so that metadata updates are
	// not coupled to (or dropped by) that workflow.
	if d.HasChange("capacity") {
		action := "UpdateKVCacheStore"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["KvcsId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)
		request["Capacity"] = d.Get("capacity")

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	action := "UpdateKVCacheStore"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KvcsId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("tag") {
		update = true
		tagMaps := make([]interface{}, 0)
		if v, ok := d.GetOk("tag"); ok {
			for _, item := range v.([]interface{}) {
				itemMap := item.(map[string]interface{})
				tagMap := make(map[string]interface{})
				if tagKey, ok := itemMap["tag_key"].(string); ok && tagKey != "" {
					tagMap["TagKey"] = tagKey
				}
				if tagValue, ok := itemMap["tag_value"].(string); ok && tagValue != "" {
					tagMap["TagValue"] = tagValue
				}
				tagMaps = append(tagMaps, tagMap)
			}
		}
		request["Tag"] = tagMaps
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	// NOTE: KVCacheStore has no public ChangeResourceGroup API, so
	// resource_group_id is ForceNew and no in-place resource group change is
	// attempted here.

	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}

	// UpdateKVCacheStore applies Name and Description asynchronously; wait until
	// GetKVCacheStore reflects the target values before finishing the update.
	if d.HasChange("name") || d.HasChange("description") {
		if rerr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			objectRaw, derr := kvcachestoreServiceV2.DescribeKvcachestoreKvCacheStore(d.Id())
			if derr != nil {
				return resource.NonRetryableError(derr)
			}
			if fmt.Sprint(objectRaw["Name"]) == d.Get("name") && fmt.Sprint(objectRaw["Description"]) == d.Get("description") {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("waiting for KvCacheStore %s name/description update to propagate", d.Id()))
		}); rerr != nil {
			return WrapError(rerr)
		}
	}

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, kvcachestoreServiceV2.KvcachestoreKvCacheStoreStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudKvcachestoreKvCacheStoreRead(d, meta)
}

func resourceAliCloudKvcachestoreKvCacheStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["KvcsId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// Deletion is asynchronous: GetKVCacheStore keeps returning the instance in
	// "Deleting" status after DeleteKVCacheStore succeeds. Wait until the read
	// reports the instance as gone. The refresh func must return a literal nil
	// result on NotFound so WaitForState treats it as absence.
	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
	deleteRefresh := func() (interface{}, string, error) {
		objectRaw, derr := kvcachestoreServiceV2.DescribeKvcachestoreKvCacheStore(d.Id())
		if derr != nil {
			if NotFoundError(derr) {
				return nil, "", nil
			}
			return nil, "", derr
		}
		return objectRaw, fmt.Sprint(objectRaw["Status"]), nil
	}
	stateConf := BuildStateConf([]string{"Creating", "Available", "InUse", "Stopping", "Stopped", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, deleteRefresh)
	if _, err := stateConf.WaitForState(); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
