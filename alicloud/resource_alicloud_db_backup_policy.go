package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBBackupPolicyCreate,
		Read:   resourceAlicloudDBBackupPolicyRead,
		Update: resourceAlicloudDBBackupPolicyUpdate,
		Delete: resourceAlicloudDBBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},

			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"retention_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(7, 730),
				Optional:     true,
				Default:      7,
			},

			"log_backup": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "Attribute 'log_backup' has been deprecated from version 1.67.0. Use `enable_backup_log`",
			},

			"enable_backup_log": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"log_retention_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(7, 730),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},

			"local_log_retention_hours": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(0, 7*24),
				Optional:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},

			"local_log_retention_space": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(5, 50),
				Optional:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},

			"high_space_usage_protection": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Default:          "Enable",
				Optional:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},
		},
	}
}

func resourceAlicloudDBBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("instance_id").(string))

	return resourceAlicloudDBBackupPolicyUpdate(d, meta)
}

func resourceAlicloudDBBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("backup_time", object.PreferredBackupTime)
	d.Set("backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("retention_period", object.BackupRetentionPeriod)
	d.Set("log_backup", object.BackupLog == "Enable")
	d.Set("enable_backup_log", object.EnableBackupLog == "1")
	d.Set("log_retention_period", object.LogBackupRetentionPeriod)
	d.Set("local_log_retention_hours", object.LocalLogRetentionHours)
	d.Set("local_log_retention_space", object.LocalLogRetentionSpace)
	d.Set("high_space_usage_protection", object.HighSpaceUsageProtection)

	return nil
}

func resourceAlicloudDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	updateForData := false
	updateForLog := false
	if d.HasChange("backup_period") || d.HasChange("backup_time") || d.HasChange("retention_period") {
		updateForData = true
	}

	if d.HasChange("log_backup") || d.HasChange("enable_backup_log") || d.HasChange("log_retention_period") ||
		d.HasChange("local_log_retention_hours") || d.HasChange("local_log_retention_space") || d.HasChange("high_space_usage_protection") {
		updateForLog = true
	}

	if updateForData || updateForLog {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := rdsService.ModifyDBBackupPolicy(d, updateForData, updateForLog); err != nil {
				if IsExceptedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudDBBackupPolicyRead(d, meta)
}

func resourceAlicloudDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	request := rds.CreateModifyBackupPolicyRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = "Tuesday,Thursday,Saturday"
	request.BackupRetentionPeriod = "7"
	request.PreferredBackupTime = "02:00Z-03:00Z"
	request.EnableBackupLog = "1"
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if instance.Engine != "SQLServer" {
		request.LogBackupRetentionPeriod = "7"
	}

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium)
}
