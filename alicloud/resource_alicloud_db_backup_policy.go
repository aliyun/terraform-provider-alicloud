package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateAllowedStringValue(BACKUP_TIME),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"retention_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(7, 730),
				Optional:     true,
				Default:      7,
			},

			"log_backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"log_retention_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(7, 730),
				Optional:         true,
				Computed:         true,
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
	d.Set("log_retention_period", object.LogBackupRetentionPeriod)

	return nil
}

func resourceAlicloudDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	update := false

	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	backupTime := d.Get("backup_time").(string)
	backupLog := "Enable"

	retentionPeriod := ""
	if temp, ok := d.GetOk("retention_period"); ok {
		retentionPeriod = strconv.Itoa(temp.(int))
	}

	logBackupRetentionPeriod := ""
	if temp, ok := d.GetOk("log_retention_period"); ok {
		logBackupRetentionPeriod = strconv.Itoa(temp.(int))
	}

	if d.HasChange("backup_period") {
		update = true
	}

	if d.HasChange("backup_time") {
		update = true
	}

	if d.HasChange("retention_period") {
		update = true
	}

	if d.HasChange("log_backup") {
		if !d.Get("log_backup").(bool) {
			backupLog = "Disabled"
		}
		update = true
	}

	if d.HasChange("log_retention_period") {
		if d.Get("log_retention_period").(int) > d.Get("retention_period").(int) {
			logBackupRetentionPeriod = retentionPeriod
		}
		update = true
	}

	if update {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := rdsService.ModifyDBBackupPolicy(d.Id(), backupTime, backupPeriod, retentionPeriod, backupLog, logBackupRetentionPeriod); err != nil {
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
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = "Tuesday,Thursday,Saturday"
	request.BackupRetentionPeriod = "7"
	request.PreferredBackupTime = "02:00Z-03:00Z"
	request.BackupLog = "Enable"
	request.LogBackupRetentionPeriod = "7"

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium)
}
