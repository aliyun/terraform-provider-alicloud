package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backup_period": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},

			"backup_time": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(rds.BACKUP_TIME),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"retention_period": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(7, 730),
				Optional:     true,
				Default:      7,
			},

			"log_backup": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"log_retention_period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(7, 730),
				Optional:         true,
				Default:          7,
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

	resp, err := meta.(*AliyunClient).rdsconn.DescribeBackupPolicy(&rds.DescribeBackupPolicyArgs{
		DBInstanceId: d.Id(),
	})
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe DB backup policy: %#v", err)
	}

	d.Set("instance_id", d.Id())
	d.Set("backup_time", resp.PreferredBackupTime)
	d.Set("backup_period", strings.Split(resp.PreferredBackupPeriod, ","))
	d.Set("retention_period", resp.BackupRetentionPeriod)
	d.Set("log_backup", resp.BackupLog == "Enable")
	d.Set("log_retention_period", resp.LogBackupRetentionPeriod)

	return nil
}

func resourceAlicloudDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	update := false
	args := rds.ModifyBackupPolicyArgs{
		DBInstanceId: d.Id(),
	}
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	args.PreferredBackupPeriod = fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	args.PreferredBackupTime = d.Get("backup_time").(string)
	args.BackupRetentionPeriod = d.Get("retention_period").(int)
	args.BackupLog = "Enable"
	args.LogBackupRetentionPeriod = strconv.Itoa(d.Get("log_retention_period").(int))

	if d.HasChange("backup_period") {
		update = true
		d.SetPartial("backup_period")
	}

	if d.HasChange("backup_time") {
		update = true
		d.SetPartial("backup_time")
	}

	if d.HasChange("retention_period") {
		update = true
		d.SetPartial("retention_period")
	}

	if d.HasChange("log_backup") {
		if !d.Get("log_backup").(bool) {
			args.BackupLog = "Disabled"
		}
		update = true
		d.SetPartial("retention_period")
	}

	if d.HasChange("log_retention_period") {
		if d.Get("log_retention_period").(int) > args.BackupRetentionPeriod {
			args.LogBackupRetentionPeriod = strconv.Itoa(args.BackupRetentionPeriod)
		}
		update = true
		d.SetPartial("log_retention_period")
	}

	if update {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			ag := args
			if _, err := meta.(*AliyunClient).rdsconn.ModifyBackupPolicy(&ag); err != nil {
				if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
					return resource.RetryableError(fmt.Errorf("ModifyBackupPolicy got an error: %#v.", err))
				}
				return resource.NonRetryableError(fmt.Errorf("ModifyBackupPolicy got an error: %#v.", err))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceAlicloudDBBackupPolicyRead(d, meta)
}

func resourceAlicloudDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	args := &rds.ModifyBackupPolicyArgs{
		DBInstanceId: d.Id(),
	}
	args.PreferredBackupTime = "02:00Z-03:00Z"
	args.PreferredBackupPeriod = "Tuesday,Thursday,Saturday"
	args.BackupRetentionPeriod = 7
	args.BackupLog = "Enable"
	args.LogBackupRetentionPeriod = "7"

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := meta.(*AliyunClient).rdsconn.ModifyBackupPolicy(args); err != nil {
			return resource.RetryableError(fmt.Errorf("ModifyBackupPolicy got an error: %#v", err))
		}

		return nil
	})
}
