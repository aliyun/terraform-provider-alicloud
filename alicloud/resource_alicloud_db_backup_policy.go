package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:       schema.TypeSet,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'backup_period' has been deprecated from version 1.69.0. Use `preferred_backup_period` instead",
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
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
				Deprecated:   "Attribute 'backup_time' has been deprecated from version 1.69.0. Use `preferred_backup_time` instead",
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"retention_period": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'retention_period' has been deprecated from version 1.69.0. Use `backup_retention_period` instead",
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
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
				Deprecated:       "Attribute 'log_retention_period' has been deprecated from version 1.69.0. Use `log_backup_retention_period` instead",
			},

			"log_backup_retention_period": {
				Type:             schema.TypeInt,
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
				ValidateFunc:     validation.IntBetween(0, 50),
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
			"released_keep_policy": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"None", "Lastest", "All"}, false),
				Optional:     true,
				Computed:     true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Flash", "Standard"}, false),
			},
			"backup_interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"-1", "15", "30", "60", "120", "180", "240", "360", "480", "720"}, false),
			},
			"backup_priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 2),
			},

			"enable_increment_data_backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_backup_local_retention_number": {
				Type:         schema.TypeInt,
				Computed:     true,
				ValidateFunc: validation.Any(IntBetween(6, 100), IntInSlice([]int{-1})),
				Optional:     true,
			},
			"backup_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Physical", "Snapshot"}, false),
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", d.Id())
	d.Set("backup_time", object["PreferredBackupTime"])
	d.Set("backup_period", strings.Split(object["PreferredBackupPeriod"].(string), ","))
	d.Set("retention_period", formatInt(object["BackupRetentionPeriod"]))
	d.Set("preferred_backup_time", object["PreferredBackupTime"])
	d.Set("preferred_backup_period", strings.Split(object["PreferredBackupPeriod"].(string), ","))
	d.Set("backup_retention_period", formatInt(object["BackupRetentionPeriod"]))
	d.Set("log_backup", object["BackupLog"] == "Enable")
	d.Set("enable_backup_log", object["EnableBackupLog"] == "1")
	d.Set("log_retention_period", formatInt(object["LogBackupRetentionPeriod"]))
	d.Set("log_backup_retention_period", formatInt(object["LogBackupRetentionPeriod"]))
	d.Set("local_log_retention_hours", formatInt(object["LocalLogRetentionHours"]))
	d.Set("local_log_retention_space", formatInt(object["LocalLogRetentionSpace"]))
	d.Set("released_keep_policy", object["ReleasedKeepPolicy"])
	d.Set("category", object["Category"])
	d.Set("enable_increment_data_backup", object["EnableIncrementDataBackup"].(bool))
	d.Set("log_backup_local_retention_number", object["LogBackupLocalRetentionNumber"])
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	// At present, the sql server database does not support setting high_space_usage_protection and it`s has default value
	if instance["Engine"] == "SQLServer" {
		d.Set("high_space_usage_protection", "Enable")
		d.Set("backup_method", object["BackupMethod"])
	} else {
		d.Set("high_space_usage_protection", object["HighSpaceUsageProtection"])
		d.Set("backup_interval", object["BackupInterval"])
	}
	if instance["Engine"] == "SQLServer" && instance["Category"] == "AlwaysOn" {
		d.Set("backup_priority", formatInt(object["BackupPriority"]))
	}
	d.Set("log_backup_frequency", object["LogBackupFrequency"])
	d.Set("compress_type", object["CompressType"])
	d.Set("archive_backup_retention_period", formatInt(object["ArchiveBackupRetentionPeriod"]))
	d.Set("archive_backup_keep_count", formatInt(object["ArchiveBackupKeepCount"]))
	d.Set("archive_backup_keep_policy", object["ArchiveBackupKeepPolicy"])
	return nil
}

func resourceAlicloudDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	updateForData := false
	updateForLog := false

	if d.HasChanges("backup_period", "backup_time", "retention_period", "preferred_backup_period", "preferred_backup_time",
		"backup_retention_period", "compress_type", "log_backup_frequency", "archive_backup_retention_period", "archive_backup_keep_count",
		"archive_backup_keep_policy", "released_keep_policy", "category", "backup_interval", "backup_priority", "backup_method",
		"enable_increment_data_backup") {
		updateForData = true
	}

	if d.HasChanges("log_backup", "enable_backup_log", "log_backup_retention_period", "log_retention_period", "local_log_retention_hours",
		"local_log_retention_space", "high_space_usage_protection", "log_backup_local_retention_number") {
		updateForLog = true
	}

	if _, ok := d.GetOkExists("enable_backup_log"); ok {
		updateForLog = true
	}
	if updateForData || updateForLog {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := rdsService.ModifyDBBackupPolicy(d, updateForData, updateForLog); err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	var err error
	action := "ModifyBackupPolicy"
	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"DBInstanceId":          d.Id(),
		"PreferredBackupPeriod": "Tuesday,Thursday,Saturday",
		"BackupRetentionPeriod": "7",
		"PreferredBackupTime":   "02:00Z-03:00Z",
		"EnableBackupLog":       "1",
		"SourceIp":              client.SourceIp,
	}
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if instance["Engine"] != "SQLServer" {
		request["LogBackupRetentionPeriod"] = "7"
	}
	if instance["Engine"] == "MySQL" && instance["DBInstanceStorageType"] == "local_ssd" {
		request["ArchiveBackupRetentionPeriod"] = "0"
		request["ArchiveBackupKeepCount"] = "1"
		request["ArchiveBackupKeepPolicy"] = "ByMonth"
	}
	response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium)
}
