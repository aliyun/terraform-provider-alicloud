package ecd

import (
	"fmt"
	"log"
	"strconv"
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

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Unlocked", "Locked"}, false),
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUsers"
	request := make(map[string]interface{})
	var err error

	requestUsers := make(map[string]interface{})
	requestUsersMap := make([]interface{}, 0)

	requestUsers["Email"] = d.Get("email")
	requestUsers["EndUserId"] = d.Get("end_user_id")
	if v, ok := d.GetOk("phone"); ok {
		requestUsers["Phone"] = v
	}

	requestUsersMap = append(requestUsersMap, requestUsers)
	request["Users"] = requestUsersMap
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_user", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(requestUsers["EndUserId"]))
	edsUserService := EdsUserService{client}

	stateConf := helper.BuildStateConf([]string{}, []string{"Unlocked", "Locked"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edsUserService.EcdUserStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceUserUpdate(d, meta)
}
func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	object, err := edsUserService.DescribeEcdUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_user edsUserService.DescribeEcdUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("end_user_id", d.Id())
	d.Set("email", object["Email"])
	d.Set("phone", object["Phone"])
	actualStatus, err := convert.FormatInt(object["Status"])
	if err != nil {
		return err
	}
	d.Set("status", convertEcdUserStatusResponse(fmt.Sprint(actualStatus)))
	return nil
}
func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := edsUserService.DescribeEcdUser(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		if target, ok := d.GetOk("status"); ok {
			actualStatus, formatErr := convert.FormatInt(object["Status"])
			if formatErr != nil {
				return formatErr
			}
			if convertEcdUserStatusRequest(strconv.Itoa(actualStatus)) != target {
				if target == "Unlocked" {
					request := map[string]interface{}{
						"Users": []string{d.Id()},
					}
					action := "UnlockUsers"
					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
						response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
				if target == "Locked" {
					request := map[string]interface{}{
						"Users": []string{d.Id()},
					}
					action := "LockUsers"
					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
						response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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

	}
	d.Partial(false)
	return resourceUserRead(d, meta)
}
func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveUsers"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Users": []string{d.Id()},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
	return nil
}
func convertEcdUserStatusResponse(source string) string {
	switch source {
	case "0":
		return "Unlocked"
	case "9":
		return "Locked"
	}
	return source
}
func convertEcdUserStatusRequest(source string) string {
	switch source {
	case "Unlocked":
		return "0"
	case "Locked":
		return "9"
	}
	return source
}
