// Package ebs. This file is generated automatically. Please do not modify it manually, thank you!
package ebs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceReplicaGroupDrill() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicaGroupDrillCreate,
		Read:   resourceReplicaGroupDrillRead,
		Delete: resourceReplicaGroupDrillDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replica_group_drill_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceReplicaGroupDrillCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "StartReplicaGroupDrill"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["GroupId"] = d.Get("group_id")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_replica_group_drill", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["GroupId"], response["DrillId"]))

	return resourceReplicaGroupDrillRead(d, meta)
}

func resourceReplicaGroupDrillRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := NewEbsServiceV2(client)

	objectRaw, err := ebsServiceV2.DescribeEbsReplicaGroupDrill(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_replica_group_drill DescribeEbsReplicaGroupDrill Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("group_id", objectRaw["GroupId"])
	d.Set("replica_group_drill_id", objectRaw["DrillId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("group_id", parts[0])
	d.Set("replica_group_drill_id", parts[1])

	return nil
}

func resourceReplicaGroupDrillDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "ClearReplicaGroupDrill"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["GroupId"] = parts[0]
	query["DrillId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.InvalidStatus"}) || helper.NeedRetry(err) {
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

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ebsServiceV2.EbsReplicaGroupDrillStateRefreshFunc(d.Id(), "DrillId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
