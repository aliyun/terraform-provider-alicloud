// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerSharedTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerSharedTargetCreate,
		Read:   resourceAliCloudResourceManagerSharedTargetRead,
		Delete: resourceAliCloudResourceManagerSharedTargetDelete,
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

func resourceAliCloudResourceManagerSharedTargetCreate(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_shared_target", action, AlibabaCloudSdkGoERROR)
	}

	ResourceShareAssociationsResourceShareIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].ResourceShareId", response)
	ResourceShareAssociationsEntityIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityId", response)
	d.SetId(fmt.Sprintf("%v:%v", ResourceShareAssociationsResourceShareIdVar, ResourceShareAssociationsEntityIdVar))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudResourceManagerSharedTargetRead(d, meta)
}

func resourceAliCloudResourceManagerSharedTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerSharedTarget(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_shared_target DescribeResourceManagerSharedTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("status", objectRaw["AssociationStatus"])
	d.Set("resource_share_id", objectRaw["ResourceShareId"])
	d.Set("target_id", objectRaw["EntityId"])

	return nil
}

func resourceAliCloudResourceManagerSharedTargetDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)

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

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Disassociated"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
