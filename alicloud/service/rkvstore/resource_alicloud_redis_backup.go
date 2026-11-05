package rkvstore

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRedisBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRedisBackupCreate,
		Read:   resourceAliCloudRedisBackupRead,
		Update: resourceAliCloudRedisBackupUpdate,
		Delete: resourceAliCloudRedisBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRedisBackupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateBackup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"BackupJobExists"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_redis_backup", action, errs.AlibabaCloudSdkGoERROR)
	}

	// Get BackupJobID from CreateBackup response
	backupJobID, ok := response["BackupJobID"]
	if !ok {
		return errs.WrapErrorf(errs.Error("BackupJobID not found in response"), errs.DefaultErrorMsg, "alicloud_redis_backup", action, errs.AlibabaCloudSdkGoERROR)
	}

	// Set temporary ID with instance_id:job_id format for waiting
	instanceID := request["InstanceId"].(string)
	d.SetId(fmt.Sprintf("%v:%v", instanceID, backupJobID))

	// Wait for backup job to complete and get BackupId
	redisServiceV2 := RedisServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, redisServiceV2.DescribeAsyncRedisBackupStateRefreshFunc(d, response, "#$.Backups.Backup[*].BackupId", []string{}))
	jobDetail, err := stateConf.WaitForState()
	if err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	// Get the actual BackupId from the job result
	if backupDetail, ok := jobDetail.(map[string]interface{}); ok {
		if backups, err := jsonpath.Get("$.Backups.Backup", backupDetail); err == nil {
			if backupList, ok := backups.([]interface{}); ok && len(backupList) > 0 {
				if backup, ok := backupList[0].(map[string]interface{}); ok {
					if backupID, ok := backup["BackupId"]; ok {
						// Update ID with the actual BackupId
						d.SetId(fmt.Sprintf("%v:%v", instanceID, backupID))
					}
				}
			}
		}
	}

	return resourceAliCloudRedisBackupRead(d, meta)
}

func resourceAliCloudRedisBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisServiceV2 := RedisServiceV2{client}

	objectRaw, err := redisServiceV2.DescribeRedisBackup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_backup DescribeRedisBackup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("status", objectRaw["BackupStatus"])
	d.Set("backup_id", objectRaw["BackupId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudRedisBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Backup.")
	return nil
}

func resourceAliCloudRedisBackupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteBackup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BackupId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
