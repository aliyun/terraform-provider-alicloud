package dbfs

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

func ResourceInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create:             resourceAliCloudDbfsInstanceAttachmentCreate,
		Read:               resourceAliCloudDbfsInstanceAttachmentRead,
		Delete:             resourceAliCloudDbfsInstanceAttachmentDelete,
		DeprecationMessage: "This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ecs_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDbfsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := Service{client}
	var response map[string]interface{}
	action := "AttachDbfs"
	request := make(map[string]interface{})
	var err error

	request["FsId"] = d.Get("instance_id")
	request["ECSInstanceId"] = d.Get("ecs_id")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dbfs_instance_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["FsId"], request["ECSInstanceId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(fmt.Sprint(request["FsId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudDbfsInstanceAttachmentRead(d, meta)
}

func resourceAliCloudDbfsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := Service{client}

	object, err := dbfsService.DescribeDbfsInstanceAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_instance_attachment dbfsService.DescribeDbfsInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("instance_id", object["FsId"])
	d.Set("ecs_id", parts[1])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudDbfsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := Service{client}
	action := "DetachDbfs"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"FsId":          parts[0],
		"ECSInstanceId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, false)
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

	stateConf := helper.BuildStateConf([]string{}, []string{"unattached"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
