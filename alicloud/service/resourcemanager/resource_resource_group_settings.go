// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceResourceGroupSettings() *schema.Resource {
	return &schema.Resource{
		Create: ResourceResourceGroupSettingsCreate,
		Read:   ResourceResourceGroupSettingsRead,
		Update: ResourceResourceGroupSettingsUpdate,
		Delete: ResourceResourceGroupSettingsDelete,
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

func ResourceResourceGroupSettingsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UpdateResourceGroupAdminSetting"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["CreatorAsAdmin"] = d.Get("resource_group_admin_setting_status")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_resource_group_settings", action, errs.AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return ResourceResourceGroupSettingsUpdate(d, meta)
}

func ResourceResourceGroupSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerResourceGroupSettings(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_group_settings DescribeResourceManagerResourceGroupSettings Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("resource_group_notification_setting_status", objectRaw["ResourceGroupNotificationEnableStatus"])

	objectRaw, err = resourceManagerServiceV2.DescribeResourceGroupSettingsGetResourceGroupAdminSetting(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("resource_group_admin_setting_status", objectRaw["CreatorAsAdmin"])

	return nil
}

func ResourceResourceGroupSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
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
		currentStatusValue, formatErr := convert.FormatBool(currentStatus)
		if formatErr != nil {
			return formatErr
		}
		if currentStatusValue != target {
			if target == false {
				action := "DisableResourceGroupNotification"
				request = make(map[string]interface{})
				query = make(map[string]interface{})

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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

			}
			if target == true {
				action := "EnableResourceGroupNotification"
				request = make(map[string]interface{})
				query = make(map[string]interface{})

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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
	}

	return ResourceResourceGroupSettingsRead(d, meta)
}

func ResourceResourceGroupSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Resource Group Settings. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
