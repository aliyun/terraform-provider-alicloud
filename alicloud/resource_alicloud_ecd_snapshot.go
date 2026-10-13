package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudEcdSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdSnapshotCreate,
		Read:   resourceAlicloudEcdSnapshotRead,
		Delete: resourceAlicloudEcdSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"desktop_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"snapshot_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_disk_type": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DATA", "SYSTEM"}, false),
				Type:         schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudEcdSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("desktop_id"); ok {
		request["DesktopId"] = v
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	}
	if v, ok := d.GetOk("source_disk_type"); ok {
		request["SourceDiskType"] = v
	}

	var response map[string]interface{}
	action := "CreateSnapshot"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		resp, err := client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_snapshot", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotId"]))
	return resourceAlicloudEcdSnapshotRead(d, meta)
}
func resourceAlicloudEcdSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}

	object, err := ecdService.DescribeEcdSnapshot(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_snapshot ecdService.DescribeEcdSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("desktop_id", object["DesktopId"])
	d.Set("snapshot_name", object["SnapshotName"])
	d.Set("source_disk_type", object["SourceDiskType"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudEcdSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	request := map[string]interface{}{
		"SnapshotId.1": d.Id(),
		"RegionId":     client.RegionId,
	}

	action := "DeleteSnapshot"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
