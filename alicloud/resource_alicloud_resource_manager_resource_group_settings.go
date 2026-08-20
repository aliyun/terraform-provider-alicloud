// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudResourceManagerResourceGroupSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerResourceGroupSettingsCreate,
		Read:   resourceAliCloudResourceManagerResourceGroupSettingsRead,
		Update: resourceAliCloudResourceManagerResourceGroupSettingsUpdate,
		Delete: resourceAliCloudResourceManagerResourceGroupSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"resource_group_admin_setting_status": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"resource_group_notification_setting_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudResourceManagerResourceGroupSettingsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UpdateResourceGroupAdminSetting"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["CreatorAsAdmin"] = d.Get("resource_group_admin_setting_status")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_group_settings", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudResourceManagerResourceGroupSettingsUpdate(d, meta)
}

func resourceAliCloudResourceManagerResourceGroupSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerResourceGroupSettings(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_group_settings DescribeResourceManagerResourceGroupSettings Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_notification_setting_status", objectRaw["ResourceGroupNotificationEnableStatus"])

	objectRaw, err = resourceManagerServiceV2.DescribeResourceGroupSettingsGetResourceGroupAdminSetting(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("resource_group_admin_setting_status", objectRaw["CreatorAsAdmin"])

	return nil
}

func resourceAliCloudResourceManagerResourceGroupSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	objectRaw, _ := resourceManagerServiceV2.DescribeResourceManagerResourceGroupSettings(d.Id())

	initedResourceGroupNotificationSettingStatus := false
	if _, ok := d.GetOkExists("resource_group_notification_setting_status"); ok && d.IsNewResource() {
		initedResourceGroupNotificationSettingStatus = true
	}
	if initedResourceGroupNotificationSettingStatus || d.HasChange("resource_group_notification_setting_status") {
		var err error
		target := d.Get("resource_group_notification_setting_status").(bool)

		currentStatus := objectRaw["ResourceGroupNotificationEnableStatus"]
		if formatBool(currentStatus) != target {
			if target == false {
				action := "DisableResourceGroupNotification"
				request = make(map[string]interface{})
				query = make(map[string]interface{})

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == true {
				action := "EnableResourceGroupNotification"
				request = make(map[string]interface{})
				query = make(map[string]interface{})

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	var err error
	action := "UpdateResourceGroupAdminSetting"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("resource_group_admin_setting_status") {
		update = true
	}
	request["CreatorAsAdmin"] = d.Get("resource_group_admin_setting_status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudResourceManagerResourceGroupSettingsRead(d, meta)
}

func resourceAliCloudResourceManagerResourceGroupSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Resource Group Settings. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
