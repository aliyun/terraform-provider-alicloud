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
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"resource_share_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_share_owner": {
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

func resourceAliCloudResourceManagerResourceShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateResourceShare"
	request := make(map[string]interface{})
	var err error

	request["ResourceShareName"] = d.Get("resource_share_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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

	if resp, err := jsonpath.Get("$.ResourceShare", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_resource_manager_resource_share")
	} else {
		resourceShareId := resp.(map[string]interface{})["ResourceShareId"]
		d.SetId(fmt.Sprint(resourceShareId))
	}

	return resourceAliCloudResourceManagerResourceShareRead(d, meta)
}
func resourceAliCloudResourceManagerResourceShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceSharingService := ResourcesharingService{client}

	object, err := resourceSharingService.DescribeResourceManagerResourceShare(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_share resourceSharingService.DescribeResourceManagerResourceShare Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_share_name", object["ResourceShareName"])
	d.Set("resource_share_owner", object["ResourceShareOwner"])
	d.Set("status", object["ResourceShareStatus"])

	return nil
}
func resourceAliCloudResourceManagerResourceShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"ResourceShareId": d.Id(),
	}

	if d.HasChange("resource_share_name") {
		update = true
	}
	request["ResourceShareName"] = d.Get("resource_share_name")

	if update {
		action := "UpdateResourceShare"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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
	}

	return resourceAliCloudResourceManagerResourceShareRead(d, meta)
}
func resourceAliCloudResourceManagerResourceShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceSharingService := ResourcesharingService{client}
	action := "DeleteResourceShare"
	var response map[string]interface{}
	var err error

	object, err := resourceSharingService.DescribeResourceManagerResourceShare(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if object["ResourceShareStatus"] == "Deleted" {
		return nil
	}

	request := map[string]interface{}{
		"ResourceShareId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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

	stateConf := BuildStateConf([]string{}, []string{"Deleted"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceSharingService.ResourceManagerResourceShareStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
