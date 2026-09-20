// Package alicloud. Provides the alicloud_kv_cache_store resource.
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

func resourceAliCloudKvCacheStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKvCacheStoreCreate,
		Read:   resourceAliCloudKvCacheStoreRead,
		Update: resourceAliCloudKvCacheStoreUpdate,
		Delete: resourceAliCloudKvCacheStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Storage capacity in GiB. Minimum 307200 GiB (300 TiB), expanding in steps of 300 TiB (307200 GiB).",
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HPN cluster code used for affinity scheduling between the KVCacheStore and a specified HPN cluster. This is a cluster identifier (e.g. \"default\", \"B6\"), not a zone ID.",
			},
			"kvcs_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"POSTPAY", "PREPAY"}, false),
				Description:  "Payment type. Valid values: `POSTPAY` (pay-as-you-go), `PREPAY` (subscription). Create currently only supports `POSTPAY`; defaults to `POSTPAY` when omitted.",
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID. Query available zones for a region via the DescribeZones API.",
			},
		},
	}
}

func resourceAliCloudKvCacheStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["Capacity"] = d.Get("capacity")
	request["ZoneId"] = d.Get("zone_id")
	request["HpnZone"] = d.Get("hpn_zone")

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PaymentType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = ConvertTagsForKms(v.(map[string]interface{}))
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kv_cache_store", action, AlibabaCloudSdkGoERROR)
	}

	id, ok := response["KvcsId"]
	if !ok || id == nil || fmt.Sprint(id) == "" {
		return WrapError(fmt.Errorf("failed to create KVCacheStore: KvcsId not found in response"))
	}
	d.SetId(fmt.Sprint(id))

	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, kvcachestoreServiceV2.KvcachestoreKVCacheStoreStateRefreshFunc(d.Id(), "Status", []string{"Creating", "Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudKvCacheStoreUpdate(d, meta)
}

func resourceAliCloudKvCacheStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}

	objectRaw, err := kvcachestoreServiceV2.DescribeKvcachestoreKVCacheStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kv_cache_store DescribeKvcachestoreKVCacheStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("capacity", objectRaw["Capacity"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("kvcs_id", objectRaw["KvcsId"])
	d.Set("name", objectRaw["Name"])
	d.Set("payment_type", objectRaw["PaymentType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	d.Set("tags", tagsToMap(objectRaw["Tags"]))
	return nil
}

func resourceAliCloudKvCacheStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	// 1. ChangeResourceGroup for resource_group_id changes
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		action := "ChangeResourceGroup"
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["ResourceId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)
		request["NewResourceGroupId"] = d.Get("resource_group_id")
		request["ResourceType"] = "KVCacheStore"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var response map[string]interface{}
		var err error
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
		d.SetPartial("resource_group_id")
	}

	// 2. UpdateKVCacheStore for updatable fields (capacity, name, description, tags)
	if !d.IsNewResource() && (d.HasChange("capacity") || d.HasChange("name") || d.HasChange("description") || d.HasChange("tags")) {
		action := "UpdateKVCacheStore"
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["KvcsId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)
		if d.HasChange("capacity") {
			request["Capacity"] = d.Get("capacity")
		}
		if d.HasChange("name") {
			request["Name"] = d.Get("name")
		}
		if d.HasChange("description") {
			request["Description"] = d.Get("description")
		}
		if d.HasChange("tags") {
			if v, ok := d.GetOk("tags"); ok {
				request["Tag"] = ConvertTagsForKms(v.(map[string]interface{}))
			} else {
				request["Tag"] = []map[string]interface{}{}
			}
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var response map[string]interface{}
		var err error
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
		kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, kvcachestoreServiceV2.KvcachestoreKVCacheStoreStateRefreshFunc(d.Id(), "Status", []string{"Creating", "Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("capacity")
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("tags")
	}

	d.Partial(false)
	return resourceAliCloudKvCacheStoreRead(d, meta)
}

func resourceAliCloudKvCacheStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KvcsId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

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

	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, kvcachestoreServiceV2.KvcachestoreKVCacheStoreStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
