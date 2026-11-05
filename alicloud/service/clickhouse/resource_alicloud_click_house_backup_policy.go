package clickhouse

import (
	"fmt"
	"log"
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

// backupTimeWindows lists the allowed preferred_backup_time values for
// ClickHouse backup policies. It mirrors the root alicloud package's
// backupTimeWindows (extension_rds.go) so the leaf package does not import root.
var backupTimeWindows = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

func ResourceClickHouseBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudClickHouseBackupPolicyCreate,
		Read:   resourceAlicloudClickHouseBackupPolicyRead,
		Update: resourceAlicloudClickHouseBackupPolicyUpdate,
		Delete: resourceAlicloudClickHouseBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(7, 730),
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"preferred_backup_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(backupTimeWindows, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudClickHouseBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPolicy"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	request["DBClusterId"] = d.Get("db_cluster_id")
	if v, ok := d.Get("preferred_backup_period").(*schema.Set); ok {
		periodList := convert.ExpandStringList(v.List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], helper.CommaSeparated)
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_click_house_backup_policy", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DBClusterId"]))
	clickhouseService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseBackupPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudClickHouseBackupPolicyRead(d, meta)
}
func resourceAlicloudClickHouseBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := NewService(client)
	object, err := clickhouseService.DescribeClickHouseBackupPolicy(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_backup_policy clickhouseService.DescribeClickHouseBackupPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	if v, ok := object["BackupRetentionPeriod"]; ok && fmt.Sprint(v) != "0" {
		backup_retention_periodVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("backup_retention_period", backup_retention_periodVal)
	}
	d.Set("preferred_backup_period", strings.Split(object["PreferredBackupPeriod"].(string), ","))
	d.Set("preferred_backup_time", object["PreferredBackupTime"])
	d.Set("status", fmt.Sprint(object["Switch"]))
	return nil
}
func resourceAlicloudClickHouseBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		update = true
	}
	if v, ok := d.Get("preferred_backup_period").(*schema.Set); ok {
		periodList := convert.ExpandStringList(v.List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], helper.CommaSeparated)
	}
	if d.HasChange("preferred_backup_time") {
		update = true
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	if d.HasChange("backup_retention_period") {
		update = true
		if v, ok := d.GetOk("backup_retention_period"); ok {
			request["BackupRetentionPeriod"] = v
		}
	}
	if update {
		action := "ModifyBackupPolicy"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
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
	return resourceAlicloudClickHouseBackupPolicyRead(d, meta)
}
func resourceAlicloudClickHouseBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy ResourceClickHouseBackupPolicy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
