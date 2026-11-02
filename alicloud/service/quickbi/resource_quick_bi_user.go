// Package quickbi. This file is generated automatically. Please do not modify it manually, thank you!
package quickbi

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceQuickBiUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuickBiUserCreate,
		Read:   resourceQuickBiUserRead,
		Update: resourceQuickBiUserUpdate,
		Delete: resourceQuickBiUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					parts := strings.Split(new, ":")
					if len(parts) < 2 {
						return false
					}
					return parts[1] == old
				},
			},
			"admin_user": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"auth_admin_user": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"nick_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Analyst", "Developer", "Visitor"}, false),
			},
		},
	}
}

func resourceQuickBiUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddUser"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("account_id"); ok {
		request["AccountId"] = v
	}
	request["AccountName"] = d.Get("account_name")
	request["AdminUser"] = d.Get("admin_user")
	request["AuthAdminUser"] = d.Get("auth_admin_user")
	request["NickName"] = d.Get("nick_name")
	request["UserType"] = convertQuickBiUserUserTypeRequest(d.Get("user_type").(string))
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("quickbi-public", "2020-08-01", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_quick_bi_user", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseResult := response["Result"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseResult["UserId"]))

	return resourceQuickBiUserRead(d, meta)
}
func resourceQuickBiUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quickbiPublicService := NewService(client)
	object, err := quickbiPublicService.DescribeQuickBiUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quick_bi_user quickbiPublicService.DescribeQuickBiUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("account_id", object["AccountId"])
	d.Set("account_name", object["AccountName"])
	d.Set("admin_user", object["AdminUser"])
	d.Set("auth_admin_user", object["AuthAdminUser"])
	d.Set("nick_name", object["NickName"])
	userType, err := convert.FormatInt(object["UserType"])
	if err != nil {
		return err
	}
	d.Set("user_type", convertQuickBiUserUserTypeResponse(userType))
	return nil
}
func resourceQuickBiUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"UserId": d.Id(),
	}
	if d.HasChange("admin_user") || d.IsNewResource() {
		update = true
	}
	request["AdminUser"] = d.Get("admin_user")
	if d.HasChange("auth_admin_user") || d.IsNewResource() {
		update = true
	}
	request["AuthAdminUser"] = d.Get("auth_admin_user")

	request["NickName"] = d.Get("nick_name")
	if d.HasChange("user_type") {
		update = true
	}
	request["UserType"] = convertQuickBiUserUserTypeRequest(d.Get("user_type").(string))
	if update {
		action := "UpdateUser"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("quickbi-public", "2020-08-01", action, nil, request, false)
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
	return resourceQuickBiUserRead(d, meta)
}
func resourceQuickBiUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteUser"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"UserId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("quickbi-public", "2020-08-01", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertQuickBiUserUserTypeRequest(source interface{}) interface{} {
	switch source {
	case "Analyst":
		return 3
	case "Developer":
		return 1
	case "Visitor":
		return 2
	}
	return 0
}
func convertQuickBiUserUserTypeResponse(source interface{}) interface{} {
	switch source {
	case 3:
		return "Analyst"
	case 1:
		return "Developer"
	case 2:
		return "Visitor"
	}
	return ""
}
