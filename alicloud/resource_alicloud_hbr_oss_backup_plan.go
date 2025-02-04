package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbrOssBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrOssBackupPlanCreate,
		Read:   resourceAlicloudHbrOssBackupPlanRead,
		Update: resourceAlicloudHbrOssBackupPlanUpdate,
		Delete: resourceAlicloudHbrOssBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceAlicloudHbrOssBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_oss_backup_plan", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return resourceAlicloudHbrOssBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrOssBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrOssBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_oss_backup_plan hbrService.DescribeHbrOssBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_type", object["BackupType"])
	d.Set("bucket", object["Bucket"])
	d.Set("oss_backup_plan_name", object["PlanName"])
	d.Set("prefix", object["Prefix"])
	d.Set("retention", fmt.Sprint(formatInt(object["Retention"])))
	d.Set("schedule", object["Schedule"])
	d.Set("vault_id", object["VaultId"])
	d.Set("disabled", object["Disabled"])
	d.Set("cross_account_type", object["CrossAccountType"])
	d.Set("cross_account_user_id", formatInt(object["CrossAccountUserId"]))
	d.Set("cross_account_role_name", object["CrossAccountRoleName"])

	return nil
}
func resourceAlicloudHbrOssBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		d.SetPartial("retention")
		d.SetPartial("vault_id")
		d.SetPartial("oss_backup_plan_name")
		d.SetPartial("prefix")
		d.SetPartial("schedule")
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrOssBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
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

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		}
	}
	d.SetPartial("disabled")
	d.Partial(false)
	return resourceAlicloudHbrOssBackupPlanRead(d, meta)
}
func resourceAlicloudHbrOssBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
