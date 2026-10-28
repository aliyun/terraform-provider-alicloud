// Package ebs. This file is generated automatically. Please do not modify it manually, thank you!
package ebs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnterpriseSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnterpriseSnapshotPolicyAttachmentCreate,
		Read:   resourceEnterpriseSnapshotPolicyAttachmentRead,
		Delete: resourceEnterpriseSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEnterpriseSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "BindEnterpriseSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PolicyId"] = d.Get("policy_id")
	query["DiskTargets.1"] = d.Get("disk_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_enterprise_snapshot_policy_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["PolicyId"], query["DiskTargets.1"]))

	return resourceEnterpriseSnapshotPolicyAttachmentRead(d, meta)
}

func resourceEnterpriseSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := NewEbsServiceV2(client)

	objectRaw, err := ebsServiceV2.DescribeEbsEnterpriseSnapshotPolicyAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_enterprise_snapshot_policy_attachment DescribeEbsEnterpriseSnapshotPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("policy_id", objectRaw["PolicyId"])

	diskIds1Raw, _ := jsonpath.Get("$.DiskIds", objectRaw)
	d.Set("disk_id", diskIds1Raw.([]interface{})[0])

	parts := strings.Split(d.Id(), ":")
	d.Set("policy_id", parts[0])
	d.Set("disk_id", parts[1])

	return nil
}

func resourceEnterpriseSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnbindEnterpriseSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PolicyId"] = parts[0]
	query["DiskTargets.1"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
