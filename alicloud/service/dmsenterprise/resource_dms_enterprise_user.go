package dmsenterprise

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudDmsEnterpriseUser() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudDmsEnterpriseUserCreate,
		Read:   ResourceAlicloudDmsEnterpriseUserRead,
		Update: ResourceAlicloudDmsEnterpriseUserUpdate,
		Delete: ResourceAlicloudDmsEnterpriseUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"max_execute_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_result_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "NORMAL"}, false),
				Default:      "NORMAL",
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nick_name"},
			},
			"nick_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'nick_name' has been deprecated from version 1.100.0. Use 'user_name' instead.",
				ConflictsWith: []string{"user_name"},
			},
		},
	}
}

func ResourceAlicloudDmsEnterpriseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RegisterUser"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("mobile"); ok {
		request["Mobile"] = v
	}

	if v, ok := d.GetOk("role_names"); ok && v != nil {
		request["RoleNames"] = tfconvert.ConvertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}

	request["Uid"] = d.Get("uid")
	if v, ok := d.GetOk("user_name"); ok {
		request["UserNick"] = v
	} else if v, ok := d.GetOk("nick_name"); ok {
		request["UserNick"] = v
	}

	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_dms_enterprise_user", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Uid"]))

	return ResourceAlicloudDmsEnterpriseUserUpdate(d, meta)
}
func ResourceAlicloudDmsEnterpriseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}
	object, err := dmsEnterpriseService.DescribeDmsEnterpriseUser(d.Id())
	if err != nil {
		if tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_user dmsEnterpriseService.DescribeDmsEnterpriseUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}

	d.Set("uid", d.Id())
	d.Set("mobile", object["Mobile"])
	d.Set("role_names", object["RoleNameList"].(map[string]interface{})["RoleNames"])
	d.Set("status", object["State"])
	d.Set("user_name", object["NickName"])
	d.Set("nick_name", object["NickName"])
	return nil
}
func ResourceAlicloudDmsEnterpriseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"Uid": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("mobile") {
		update = true
		request["Mobile"] = d.Get("mobile")
	}
	if !d.IsNewResource() && d.HasChange("role_names") {
		update = true
		request["RoleNames"] = tfconvert.ConvertListToCommaSeparate(d.Get("role_names").(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("user_name") {
		update = true
		request["UserNick"] = d.Get("user_name")
	}
	if !d.IsNewResource() && d.HasChange("nick_name") {
		update = true
		request["UserNick"] = d.Get("nick_name")
	}
	if update {
		if _, ok := d.GetOk("max_execute_count"); ok {
			request["MaxExecuteCount"] = d.Get("max_execute_count")
		}
		if _, ok := d.GetOk("max_result_count"); ok {
			request["MaxResultCount"] = d.Get("max_result_count")
		}
		if _, ok := d.GetOk("tid"); ok {
			request["Tid"] = d.Get("tid")
		}
		action := "UpdateUser"
		wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
			if err != nil {
				if tfhelper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			tflogs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("status") {
		object, err := dmsEnterpriseService.DescribeDmsEnterpriseUser(d.Id())
		if err != nil {
			return tferrs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "DISABLE" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				if v, ok := d.GetOk("tid"); ok {
					request["Tid"] = v
				}
				action := "DisableUser"
				wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
					if err != nil {
						if tfhelper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					tflogs.AddDebug(action, response, request)
					return nil
				})
				if err != nil {
					return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
				}
			}
			if target == "NORMAL" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				if v, ok := d.GetOk("tid"); ok {
					request["Tid"] = v
				}
				action := "EnableUser"
				wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
					if err != nil {
						if tfhelper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					tflogs.AddDebug(action, response, request)
					return nil
				})
				if err != nil {
					return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	d.Partial(false)
	return ResourceAlicloudDmsEnterpriseUserRead(d, meta)
}
func ResourceAlicloudDmsEnterpriseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteUser"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Uid": d.Id(),
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
