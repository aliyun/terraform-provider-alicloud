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

func resourceAliCloudResourceManagerResourceShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerResourceShareCreate,
		Read:   resourceAliCloudResourceManagerResourceShareRead,
		Update: resourceAliCloudResourceManagerResourceShareUpdate,
		Delete: resourceAliCloudResourceManagerResourceShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_external_targets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_share_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_share_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudResourceManagerResourceShareCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("targets"); ok {
		targetsMapsArray := convertToInterfaceArray(v)

		request["Targets"] = targetsMapsArray
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("resources"); ok {
		resourcesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ResourceType"] = dataLoop1Tmp["resource_type"]
			dataLoop1Map["ResourceId"] = dataLoop1Tmp["resource_id"]
			resourcesMapsArray = append(resourcesMapsArray, dataLoop1Map)
		}
		request["Resources"] = resourcesMapsArray
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOkExists("allow_external_targets"); ok {
		request["AllowExternalTargets"] = v
	}
	if v, ok := d.GetOk("permission_names"); ok {
		permissionNamesMapsArray := convertToInterfaceArray(v)

		request["PermissionNames"] = permissionNamesMapsArray
	}

	request["ResourceShareName"] = d.Get("resource_share_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_share", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.ResourceShare.ResourceShareId", response)
	d.SetId(fmt.Sprint(id))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerResourceShareStateRefreshFunc(d.Id(), "ResourceShareStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudResourceManagerResourceShareUpdate(d, meta)
}

func resourceAliCloudResourceManagerResourceShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerResourceShare(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_share DescribeResourceManagerResourceShare Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("allow_external_targets", objectRaw["AllowExternalTargets"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("resource_share_name", objectRaw["ResourceShareName"])
	d.Set("resource_share_owner", objectRaw["ResourceShareOwner"])
	d.Set("status", objectRaw["ResourceShareStatus"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudResourceManagerResourceShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateResourceShare"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceShareId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("allow_external_targets") {
		update = true
		request["AllowExternalTargets"] = d.Get("allow_external_targets")
	}

	if !d.IsNewResource() && d.HasChange("resource_share_name") {
		update = true
	}
	request["ResourceShareName"] = d.Get("resource_share_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourceManagerServiceV2.ResourceManagerResourceShareStateRefreshFunc(d.Id(), "ResourceShareStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourceManagerServiceV2.ResourceManagerResourceShareStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		if err := resourceManagerServiceV2.SetResourceTagsForResourceSharing(d, "ResourceShare"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudResourceManagerResourceShareRead(d, meta)
}

func resourceAliCloudResourceManagerResourceShareDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ResourceShareId"] = d.Id()
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.ResourceManagerResourceShareStateRefreshFunc(d.Id(), "ResourceShareStatus", []string{"Active"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
