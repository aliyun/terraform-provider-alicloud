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
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"preferred_backup_period"},
				Deprecated:    "Attribute 'backup_period' has been deprecated from version 1.69.0. Use `preferred_backup_period` instead",
			},

			"preferred_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},

			"backup_time": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice(BACKUP_TIME, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"preferred_backup_time"},
				Deprecated:    "Attribute 'backup_time' has been deprecated from version 1.69.0. Use `preferred_backup_time` instead",
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"retention_period": {
				Type:          schema.TypeInt,
				ValidateFunc:  validation.IntBetween(7, 730),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"backup_retention_period"},
				Deprecated:    "Attribute 'retention_period' has been deprecated from version 1.69.0. Use `backup_retention_period` instead",
			},

			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  7,
			},

			"log_backup": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'log_backup' has been deprecated from version 1.68.0. Use `enable_backup_log` instead",
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
				ConflictsWith:    []string{"log_backup_retention_period"},
				Deprecated:       "Attribute 'log_retention_period' has been deprecated from version 1.69.0. Use `log_backup_retention_period` instead",
			},

			"log_backup_retention_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(7, 730),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},

			"local_log_retention_hours": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(0, 7*24),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"local_log_retention_space": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(5, 50),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"high_space_usage_protection": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Default:          "Enable",
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"log_backup_frequency": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"compress_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"1", "4", "8"}, false),
				Computed:     true,
				Optional:     true,
			},

			"archive_backup_retention_period": {
				Type:             schema.TypeInt,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: archiveBackupPeriodDiffSuppressFunc,
			},

			"archive_backup_keep_count": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 31),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"archive_backup_keep_policy": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"ByMonth", "ByWeek", "KeepAll"}, false),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
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
	d.Set("preferred_backup_time", object.PreferredBackupTime)
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("backup_retention_period", object.BackupRetentionPeriod)
	d.Set("log_backup", object.BackupLog == "Enable")
	d.Set("enable_backup_log", object.EnableBackupLog == "1")
	d.Set("log_retention_period", object.LogBackupRetentionPeriod)
	d.Set("log_backup_retention_period", object.LogBackupRetentionPeriod)
	d.Set("local_log_retention_hours", object.LocalLogRetentionHours)
	d.Set("local_log_retention_space", object.LocalLogRetentionSpace)
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	// At present, the sql server database does not support setting high_space_usage_protection and it`s has default value
	if instance.Engine == "SQLServer" {
		d.Set("high_space_usage_protection", "Enable")
	} else {
		d.Set("high_space_usage_protection", object.HighSpaceUsageProtection)
	}
	d.Set("log_backup_frequency", object.LogBackupFrequency)
	d.Set("compress_type", object.CompressType)
	d.Set("archive_backup_retention_period", object.ArchiveBackupRetentionPeriod)
	d.Set("archive_backup_keep_count", object.ArchiveBackupKeepCount)
	d.Set("archive_backup_keep_policy", object.ArchiveBackupKeepPolicy)
	return nil
}

func resourceAlicloudDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	updateForData := false
	updateForLog := false
	if d.HasChange("backup_period") || d.HasChange("backup_time") || d.HasChange("retention_period") ||
		d.HasChange("preferred_backup_period") || d.HasChange("preferred_backup_time") || d.HasChange("backup_retention_period") ||
		d.HasChange("compress_type") || d.HasChange("log_backup_frequency") || d.HasChange("archive_backup_retention_period") ||
		d.HasChange("archive_backup_keep_count") || d.HasChange("archive_backup_keep_policy") {
		updateForData = true
	}

	if d.HasChange("log_backup") || d.HasChange("enable_backup_log") || d.HasChange("log_backup_retention_period") || d.HasChange("log_retention_period") ||
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
	if instance.Engine == "MySQL" && instance.DBInstanceStorageType == "local_ssd" {
		request.ArchiveBackupRetentionPeriod = "0"
		request.ArchiveBackupKeepCount = "1"
		request.ArchiveBackupKeepPolicy = "ByMonth"
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
