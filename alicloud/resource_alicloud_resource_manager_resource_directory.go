package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerResourceDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerResourceDirectoryCreate,
		Read:   resourceAliCloudResourceManagerResourceDirectoryRead,
		Update: resourceAliCloudResourceManagerResourceDirectoryUpdate,
		Delete: resourceAliCloudResourceManagerResourceDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"member_deletion_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"root_folder_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudResourceManagerResourceDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "InitResourceDirectory"
	request := make(map[string]interface{})
	var err error

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_directory", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.ResourceDirectory", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_resource_manager_resource_directory")
	} else {
		resourceDirectoryId := resp.(map[string]interface{})["ResourceDirectoryId"]
		d.SetId(fmt.Sprint(resourceDirectoryId))
	}

	return resourceAliCloudResourceManagerResourceDirectoryUpdate(d, meta)
}

func resourceAliCloudResourceManagerResourceDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}

	object, err := resourceManagerService.DescribeResourceManagerResourceDirectory(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_directory resourceManagerService.DescribeResourceManagerResourceDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", object["ScpStatus"])
	d.Set("member_deletion_status", object["MemberDeletionStatus"])
	d.Set("root_folder_id", object["RootFolderId"])
	d.Set("master_account_id", object["MasterAccountId"])
	d.Set("master_account_name", object["MasterAccountName"])

	return nil
}

func resourceAliCloudResourceManagerResourceDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := resourceManagerService.DescribeResourceManagerResourceDirectory(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["ScpStatus"].(string) != target {
			if target == "Disabled" {
				request := map[string]interface{}{}

				action := "DisableControlPolicy"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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

				stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourceManagerService.ResourceManagerResourceDirectoryStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if target == "Enabled" {
				request := map[string]interface{}{}

				action := "EnableControlPolicy"

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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

				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourceManagerService.ResourceManagerResourceDirectoryStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			d.SetPartial("status")
		}
	}

	if d.HasChange("member_deletion_status") {
		object, err := resourceManagerService.DescribeResourceManagerResourceDirectory(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("member_deletion_status").(string)
		if object["MemberDeletionStatus"].(string) != target {
			request := map[string]interface{}{
				"Status": target,
			}

			action := "SetMemberDeletionPermission"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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

			d.SetPartial("member_deletion_status")
		}
	}

	d.Partial(false)

	return resourceAliCloudResourceManagerResourceDirectoryRead(d, meta)
}

func resourceAliCloudResourceManagerResourceDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	action := "DestroyResourceDirectory"
	request := map[string]interface{}{}
	var response map[string]interface{}
	var err error
	object, err := resourceManagerService.DescribeResourceManagerResourceDirectory(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(object["ScpStatus"]) == "Enabled" {
		disableAction := "DisableControlPolicy"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", disableAction, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(disableAction, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), disableAction, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourceManagerService.ResourceManagerResourceDirectoryStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcGet("ResourceManager", "2020-03-31", action, request, nil)
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

	return nil
}
