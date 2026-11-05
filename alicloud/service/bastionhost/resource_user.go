package bastionhost

import (
	"fmt"
	"log"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AE", "AU", "CH", "CN", "DE", "GB", "HK", "ID", "IN", "JP", "KR", "MO", "MY", "PH", "RU", "SE", "SG", "TW", "US"}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Ram", "AD", "LDAP"}, false),
			},
			"source_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Frozen", "Normal", "RemoteDNChanged"}, false),
			},
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

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUser"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("mobile"); ok {
		request["Mobile"] = v
	} else if v, ok := d.GetOk("mobile_country_code"); ok && v.(string) == "" {
		return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "mobile", "mobile_country_code", d.Get("mobile_country_code")))
	}
	if v, ok := d.GetOk("mobile_country_code"); ok {
		request["MobileCountryCode"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	} else if v, ok := d.GetOk("source"); ok && v.(string) == "Local" {
		return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "password", "source", d.Get("source")))
	}
	request["RegionId"] = client.RegionId
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_user_id"); ok {
		request["SourceUserId"] = v
	} else if v, ok := d.GetOk("source"); ok && v.(string) == "Ram" {
		return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "source_user_id", "source", d.Get("source")))
	}
	request["UserName"] = d.Get("user_name")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_bastionhost_user", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["UserId"]))

	return resourceUserUpdate(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	object, err := yundunBastionhostService.DescribeBastionhostUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_user yundunBastionhostService.DescribeBastionhostUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("user_id", parts[1])
	d.Set("comment", object["Comment"])
	d.Set("display_name", object["DisplayName"])
	d.Set("email", object["Email"])
	d.Set("mobile", object["Mobile"])
	d.Set("mobile_country_code", object["MobileCountryCode"])
	d.Set("source", object["Source"])
	d.Set("source_user_id", object["SourceUserId"])
	d.Set("status", convert.ConvertArrayToString(object["UserState"], ""))
	d.Set("user_name", object["UserName"])
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	if !d.IsNewResource() && d.HasChange("comment") {
		update = true
		if v, ok := d.GetOk("comment"); ok {
			request["Comment"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
		if v, ok := d.GetOk("display_name"); ok {
			request["DisplayName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("email") {
		update = true
		if v, ok := d.GetOk("email"); ok {
			request["Email"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("mobile") {
		update = true
		if v, ok := d.GetOk("mobile"); ok {
			request["Mobile"] = v
		} else if v, ok := d.GetOk("mobile_country_code"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "mobile", "mobile_country_code", d.Get("mobile_country_code")))
		}
		request["MobileCountryCode"] = d.Get("mobile_country_code")
	}
	if !d.IsNewResource() && d.HasChange("mobile_country_code") {
		update = true
		if v, ok := d.GetOk("mobile_country_code"); ok {
			request["MobileCountryCode"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			request["Password"] = v
		} else if v, ok := d.GetOk("source"); ok && v.(string) == "Local" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "password", "source", d.Get("source")))
		}
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyUser"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
	if d.HasChange("status") {
		object, err := yundunBastionhostService.DescribeBastionhostUser(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if convert.ConvertArrayToString(object["UserState"], "") != target {
			if target == "Frozen" {
				request := map[string]interface{}{
					"InstanceId": parts[0],
					"UserIds":    convert.ConvertListToJsonString([]interface{}{parts[1]}),
				}
				request["RegionId"] = client.RegionId
				action := "LockUsers"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
					response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
			if target == "Normal" {
				request := map[string]interface{}{
					"InstanceId": parts[0],
					"UserIds":    convert.ConvertListToJsonString([]interface{}{parts[1]}),
				}
				request["RegionId"] = client.RegionId
				action := "UnlockUsers"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
					response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
	d.Partial(false)
	return resourceUserRead(d, meta)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteUser"
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"OBJECT_NOT_FOUND"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
