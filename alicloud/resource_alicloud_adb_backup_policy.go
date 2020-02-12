package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAdbBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbBackupPolicyCreate,
		Read:   resourceAlicloudAdbBackupPolicyRead,
		Update: resourceAlicloudAdbBackupPolicyUpdate,
		Delete: resourceAlicloudAdbBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"preferred_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},
			"backup_retention_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAdbBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("db_cluster_id").(string))

	return resourceAlicloudAdbBackupPolicyUpdate(d, meta)
}

func resourceAlicloudAdbBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	d.Set("backup_retention_period", object.BackupRetentionPeriod)
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("preferred_backup_time", object.PreferredBackupTime)

	return nil
}

func resourceAlicloudAdbBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	if d.HasChange("preferred_backup_period") || d.HasChange("preferred_backup_time") {
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		preferredBackupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		preferredBackupTime := d.Get("preferred_backup_time").(string)

		// wait instance running before modifying
		if err := adbService.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := adbService.ModifyDBBackupPolicy(d.Id(), preferredBackupTime, preferredBackupPeriod); err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudAdbBackupPolicyRead(d, meta)
}

func resourceAlicloudAdbBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	//  Terraform can not destroy it..
	return nil
}
