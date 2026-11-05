// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager

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

func ResourceSharedTarget() *schema.Resource {
	return &schema.Resource{
		Create: ResourceSharedTargetCreate,
		Read:   ResourceSharedTargetRead,
		Delete: ResourceSharedTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceSharedTargetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ResourceShareId"] = d.Get("resource_share_id")
	request["Targets"] = []string{d.Get("target_id").(string)}
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_shared_target", action, errs.AlibabaCloudSdkGoERROR)
	}

	ResourceShareAssociationsResourceShareIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].ResourceShareId", response)
	ResourceShareAssociationsEntityIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityId", response)
	d.SetId(fmt.Sprintf("%v:%v", ResourceShareAssociationsResourceShareIdVar, ResourceShareAssociationsEntityIdVar))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceSharedTargetRead(d, meta)
}

func ResourceSharedTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerSharedTarget(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_shared_target DescribeResourceManagerSharedTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("status", objectRaw["AssociationStatus"])
	d.Set("resource_share_id", objectRaw["ResourceShareId"])
	d.Set("target_id", objectRaw["EntityId"])

	return nil
}

func ResourceSharedTargetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DisassociateResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ResourceShareId"] = parts[0]
	request["Targets"] = []string{parts[1]}
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)

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

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Disassociated"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
