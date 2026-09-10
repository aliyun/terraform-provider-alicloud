// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsSnapshotCreate,
		Read:   resourceAliCloudEnsSnapshotRead,
		Update: resourceAliCloudEnsSnapshotUpdate,
		Delete: resourceAliCloudEnsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEnsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["EnsRegionId"] = d.Get("ens_region_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DiskId"] = d.Get("disk_id")
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_snapshot", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.SnapShotId[0]", response)
	d.SetId(fmt.Sprint(id))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ensServiceV2.EnsSnapshotStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsSnapshotUpdate(d, meta)
}

func resourceAliCloudEnsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsSnapshot(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_snapshot DescribeEnsSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("disk_id", objectRaw["SourceDiskId"])
	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("snapshot_name", objectRaw["SnapshotName"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudEnsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifySnapshotAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SnapshotId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("snapshot_name") {
		update = true
		request["SnapshotName"] = d.Get("snapshot_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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
	}

	return resourceAliCloudEnsSnapshotRead(d, meta)
}

func resourceAliCloudEnsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSnapshot"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SnapshotId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsSnapshotStateRefreshFunc(d.Id(), "SnapshotId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
