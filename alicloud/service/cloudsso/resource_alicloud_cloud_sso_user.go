// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudsso

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tags "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudSsoUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSsoUserCreate,
		Read:   resourceAliCloudCloudSsoUserRead,
		Update: resourceAliCloudCloudSsoUserUpdate,
		Delete: resourceAliCloudCloudSsoUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mfa_authentication_settings": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"tags": tags.TagsSchema(),
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudSsoUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateUser"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DirectoryId"] = d.Get("directory_id")

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMapWithTags(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("first_name"); ok {
		request["FirstName"] = v
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("last_name"); ok {
		request["LastName"] = v
	}
	request["UserName"] = d.Get("user_name")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_sso_user", action, errs.AlibabaCloudSdkGoERROR)
	}

	UserUserIdVar, _ := jsonpath.Get("$.User.UserId", response)
	d.SetId(fmt.Sprintf("%v:%v", request["DirectoryId"], UserUserIdVar))

	return resourceAliCloudCloudSsoUserUpdate(d, meta)
}

func resourceAliCloudCloudSsoUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSSOServiceV2 := CloudSSOServiceV2{client}

	objectRaw, err := cloudSSOServiceV2.DescribeCloudSsoUser(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_user DescribeCloudSsoUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("display_name", objectRaw["DisplayName"])
	d.Set("email", objectRaw["Email"])
	d.Set("first_name", objectRaw["FirstName"])
	d.Set("last_name", objectRaw["LastName"])
	d.Set("status", objectRaw["Status"])
	d.Set("user_name", objectRaw["UserName"])
	d.Set("user_id", objectRaw["UserId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	objectRaw, err = cloudSSOServiceV2.DescribeUserGetUserMFAAuthenticationSettings(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("mfa_authentication_settings", objectRaw["UserMFAAuthenticationSettings"])

	parts := strings.Split(d.Id(), ":")
	d.Set("directory_id", parts[0])

	return nil
}

func resourceAliCloudCloudSsoUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateUser"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["DirectoryId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("last_name") {
		update = true
		request["NewLastName"] = d.Get("last_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("first_name") {
		update = true
		request["NewFirstName"] = d.Get("first_name")
	}

	if !d.IsNewResource() && d.HasChange("email") {
		update = true
		request["NewEmail"] = d.Get("email")
	}

	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
		request["NewDisplayName"] = d.Get("display_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateUserStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["DirectoryId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("status") {
		update = true
		request["NewStatus"] = d.Get("status")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ResetUserPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["DirectoryId"] = parts[0]

	if d.HasChange("password") {
		update = true
		request["Password"] = d.Get("password")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateUserMFAAuthenticationSettings"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["DirectoryId"] = parts[0]

	if d.HasChange("mfa_authentication_settings") {
		update = true
	}
	request["UserMFAAuthenticationSettings"] = d.Get("mfa_authentication_settings")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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

	if !d.IsNewResource() && d.HasChange("tags") {
		cloudSSOServiceV2 := CloudSSOServiceV2{client}
		if err := cloudSSOServiceV2.SetResourceTags(d, "user"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudCloudSsoUserRead(d, meta)
}

func resourceAliCloudCloudSsoUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteUser"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["DirectoryId"] = parts[0]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeletionConflict.User.AccessAssigment"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.User"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
