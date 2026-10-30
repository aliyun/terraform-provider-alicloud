package hbr

import (
	"fmt"
	"log"
	"strconv"
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

func ResourceNasBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: ResourceNasBackupPlanCreate,
		Read:   ResourceNasBackupPlanRead,
		Update: ResourceNasBackupPlanUpdate,
		Delete: ResourceNasBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.249.0 and will be removed in the future. Please use 'alicloud_hbr_policy' and 'alicloud_hbr_policy_binding' instead.",
		Schema: map[string]*schema.Schema{
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE"}, false),
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'create_time' has been deprecated from provider version 1.153.0.",
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nas_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cross_account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SELF_ACCOUNT", "CROSS_ACCOUNT"}, false),
			},
			"cross_account_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"cross_account_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceNasBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	var err error
	request["BackupType"] = d.Get("backup_type")
	if v, ok := d.GetOk("create_time"); ok {
		request["CreateTime"] = ConvertNasFileSystemStringToUnix(v.(string))
	}
	if v, ok := d.GetOk("file_system_id"); ok {
		request["FileSystemId"] = v
	}
	request["PlanName"] = d.Get("nas_backup_plan_name")
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	request["Schedule"] = d.Get("schedule")
	request["SourceType"] = "NAS"
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	if v, ok := d.GetOk("cross_account_type"); ok {
		request["CrossAccountType"] = v
	}
	if v, ok := d.GetOk("cross_account_user_id"); ok {
		request["CrossAccountUserId"] = v
	}
	if v, ok := d.GetOk("cross_account_role_name"); ok {
		request["CrossAccountRoleName"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hbr_nas_backup_plan", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return ResourceNasBackupPlanUpdate(d, meta)
}

func ResourceNasBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrNasBackupPlan(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_nas_backup_plan hbrService.DescribeHbrNasBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("backup_type", object["BackupType"])
	createTime, err := convert.FormatInt(object["CreateTime"])
	if err != nil {
		return err
	}
	t := int64(createTime)
	d.Set("create_time", ConvertNasFileSystemUnixToString(d.Get("create_time").(string), t))
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("nas_backup_plan_name", object["PlanName"])
	d.Set("options", object["Options"])
	if object["Paths"] != nil {
		d.Set("path", object["Paths"].(map[string]interface{})["Path"])
	}
	actualRetention, err := convert.FormatInt(object["Retention"])
	if err != nil {
		return err
	}
	d.Set("retention", fmt.Sprint(actualRetention))
	d.Set("schedule", object["Schedule"])
	d.Set("vault_id", object["VaultId"])
	d.Set("disabled", object["Disabled"])
	d.Set("cross_account_type", object["CrossAccountType"])
	cross_account_user_idVal, err := convert.FormatInt(object["CrossAccountUserId"])
	if err != nil {
		return err
	}
	d.Set("cross_account_user_id", cross_account_user_idVal)
	d.Set("cross_account_role_name", object["CrossAccountRoleName"])

	return nil
}

func ResourceNasBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("retention") {
		update = true
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}

	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	if !d.IsNewResource() && d.HasChange("nas_backup_plan_name") {
		update = true
		request["PlanName"] = d.Get("nas_backup_plan_name")
	}
	if !d.IsNewResource() && d.HasChange("options") {
		update = true
		if v, ok := d.GetOk("options"); ok {
			request["Options"] = v
		}
	}
	if d.HasChange("path") {
		update = true
		request["UpdatePaths"] = true
		if v, ok := d.GetOk("path"); ok {
			request["Path"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	request["SourceType"] = "NAS"
	if update {
		action := "UpdateBackupPlan"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrNasBackupPlan(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := strconv.FormatBool(d.Get("disabled").(bool))
		if strconv.FormatBool(object["Disabled"].(bool)) != target {
			action := "EnableBackupPlan"
			if target == "false" {
				action = "EnableBackupPlan"
			} else {
				action = "DisableBackupPlan"
			}
			request := map[string]interface{}{
				"PlanId": d.Id(),
			}
			request["VaultId"] = d.Get("vault_id")
			request["SourceType"] = "NAS"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
			if fmt.Sprint(response["Success"]) == "false" {
				return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
		}
	}
	d.Partial(false)
	return ResourceNasBackupPlanRead(d, meta)
}

func ResourceNasBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupPlan"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}

	request["SourceType"] = "NAS"
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	{
		action := "DetachNasFileSystem"
		request := map[string]interface{}{
			"PlanId": d.Id(),
		}
		request["VaultId"] = d.Get("vault_id")
		if v, ok := d.GetOk("file_system_id"); ok {
			request["FileSystemId"] = v
		}
		if v, ok := d.GetOk("create_time"); ok {
			request["CreateTime"] = ConvertNasFileSystemStringToUnix(v.(string))
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
			if errs.IsExpectedErrors(err, []string{"NasFileSystemNotExist"}) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		time.Sleep(1 * time.Minute)
	}
	return nil
}

func ConvertNasFileSystemUnixToString(v string, t int64) string {
	var s string
	c := time.Unix(t, 0)
	if strings.HasSuffix(v, "CST") {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05CST"))
	} else if strings.HasSuffix(v, "Z") {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05Z"))
	} else {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05CST"))
	}
	return s
}

func ConvertNasFileSystemStringToUnix(v string) int64 {
	var t time.Time
	if strings.HasSuffix(v, "CST") {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05CST", v, time.Local)
	} else if strings.HasSuffix(v, "Z") {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05Z", v, time.FixedZone("CST", 8*3600))
	} else {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05CST", v, time.Local)
	}
	return t.Unix()
}
