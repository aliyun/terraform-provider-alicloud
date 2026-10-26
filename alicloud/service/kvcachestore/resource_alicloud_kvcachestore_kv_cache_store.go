package kvcachestore

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceKvCacheStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceKvCacheStoreCreate,
		Read:   resourceKvCacheStoreRead,
		Update: resourceKvCacheStoreUpdate,
		Delete: resourceKvCacheStoreDelete,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKvCacheStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kvcachestore_kv_cache_store", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["KvcsId"]))

	kvCacheStoreService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, kvCacheStoreService.KvCacheStoreStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceKvCacheStoreRead(d, meta)
}

func resourceKvCacheStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvCacheStoreService := NewService(client)

	objectRaw, err := kvCacheStoreService.DescribeKvCacheStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvcachestore_kv_cache_store DescribeKvCacheStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("capacity", objectRaw["Capacity"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("name", objectRaw["Name"])
	d.Set("payment_type", objectRaw["PaymentType"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	return nil
}

func resourceKvCacheStoreUpdate(d *schema.ResourceData, meta interface{}) error {
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
		request["ClientToken"] = helper.BuildClientToken(action)
		request["Capacity"] = d.Get("capacity")

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
	}

	action := "UpdateKVCacheStore"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KvcsId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
	}
	kvCacheStoreService := NewService(client)

	// UpdateKVCacheStore applies Name and Description asynchronously; wait until
	// GetKVCacheStore reflects the target values before finishing the update.
	if d.HasChange("name") || d.HasChange("description") {
		if rerr := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			objectRaw, derr := kvCacheStoreService.DescribeKvCacheStore(d.Id())
			if derr != nil {
				return retry.NonRetryableError(derr)
			}
			if fmt.Sprint(objectRaw["Name"]) == d.Get("name") && fmt.Sprint(objectRaw["Description"]) == d.Get("description") {
				return nil
			}
			return retry.RetryableError(fmt.Errorf("waiting for KvCacheStore %s name/description update to propagate", d.Id()))
		}); rerr != nil {
			return errs.WrapError(rerr)
		}
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, kvCacheStoreService.KvCacheStoreStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceKvCacheStoreRead(d, meta)
}

func resourceKvCacheStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKVCacheStore"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["KvcsId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	// Deletion is asynchronous: GetKVCacheStore keeps returning the instance in
	// "Deleting" status after DeleteKVCacheStore succeeds. Wait until the read
	// reports the instance as gone. The refresh func must return a literal nil
	// result on NotFound so WaitForState treats it as absence.
	kvCacheStoreService := NewService(client)
	deleteRefresh := func() (interface{}, string, error) {
		objectRaw, derr := kvCacheStoreService.DescribeKvCacheStore(d.Id())
		if derr != nil {
			if errs.NotFoundError(derr) {
				return nil, "", nil
			}
			return nil, "", derr
		}
		return objectRaw, fmt.Sprint(objectRaw["Status"]), nil
	}
	stateConf := helper.BuildStateConf([]string{"Creating", "Available", "InUse", "Stopping", "Stopped", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, deleteRefresh)
	if _, err := stateConf.WaitForState(); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
