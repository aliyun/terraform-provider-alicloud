package swasopen

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotCreate,
		Read:   resourceSnapshotRead,
		Delete: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z一-龥][一-龥A-Za-z0-9:._-]{2,50}$"), "The name must be 2 to 50 characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-)."),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSnapshot"
	request := make(map[string]interface{})
	var err error
	request["DiskId"] = d.Get("disk_id")
	request["RegionId"] = client.RegionId
	request["SnapshotName"] = d.Get("snapshot_name")
	request["ClientToken"] = helper.BuildClientToken("CreateSnapshot")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_simple_application_server_snapshot", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotId"]))
	swasOpenService := NewSwasOpenService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, swasOpenService.SimpleApplicationServerSnapshotStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceSnapshotRead(d, meta)
}
func resourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := NewSwasOpenService(client)
	object, err := swasOpenService.DescribeSimpleApplicationServerSnapshot(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_snapshot swasOpenService.DescribeSimpleApplicationServerSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("disk_id", object["SourceDiskId"])
	d.Set("snapshot_name", object["SnapshotName"])
	d.Set("status", object["Status"])
	return nil
}
func resourceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSnapshot"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"SnapshotId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("DeleteSnapshot")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
	return nil
}
