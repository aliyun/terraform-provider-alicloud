package hbr

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

func ResourceOssBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOssBackupPlanCreate,
		Read:   ResourceOssBackupPlanRead,
		Update: ResourceOssBackupPlanUpdate,
		Delete: ResourceOssBackupPlanDelete,
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
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
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

func ResourceOssBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	var err error
	request["BackupType"] = d.Get("backup_type")
	if v, ok := d.GetOk("bucket"); ok {
		request["Bucket"] = v
	}
	request["PlanName"] = d.Get("oss_backup_plan_name")
	if v, ok := d.GetOk("path"); ok {
		request["Path"] = v
	}
	if v, ok := d.GetOk("prefix"); ok {
		request["Prefix"] = v
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	request["Schedule"] = d.Get("schedule")
	request["SourceType"] = "OSS"
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hbr_oss_backup_plan", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return ResourceOssBackupPlanUpdate(d, meta)
}
func ResourceOssBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrOssBackupPlan(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_oss_backup_plan hbrService.DescribeHbrOssBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("backup_type", object["BackupType"])
	d.Set("bucket", object["Bucket"])
	d.Set("oss_backup_plan_name", object["PlanName"])
	d.Set("prefix", object["Prefix"])
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
func ResourceOssBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var err error
	var response map[string]interface{}
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
	if !d.IsNewResource() && d.HasChange("oss_backup_plan_name") {
		update = true
		request["PlanName"] = d.Get("oss_backup_plan_name")
	}
	if !d.IsNewResource() && d.HasChange("prefix") {
		update = true
		request["Prefix"] = d.Get("prefix")
	}
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	request["SourceType"] = "OSS"
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
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrOssBackupPlan(d.Id())
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
			request["SourceType"] = "OSS"

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
		}
	}
	d.Partial(false)
	return ResourceOssBackupPlanRead(d, meta)
}
func ResourceOssBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupPlan"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}

	request["SourceType"] = "OSS"
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
	return nil
}
