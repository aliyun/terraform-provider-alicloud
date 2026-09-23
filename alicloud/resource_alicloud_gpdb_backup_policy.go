// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbBackupPolicyCreate,
		Read:   resourceAliCloudGpdbBackupPolicyRead,
		Update: resourceAliCloudGpdbBackupPolicyUpdate,
		Delete: resourceAliCloudGpdbBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable_recovery_point": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recovery_point_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyBackupPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Get("db_instance_id")

	if v, ok := d.GetOkExists("enable_recovery_point"); ok {
		request["EnableRecoveryPoint"] = v
	}
	if v, ok := d.GetOk("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	if v, ok := d.GetOk("recovery_point_period"); ok {
		request["RecoveryPointPeriod"] = v
	}
	request["PreferredBackupPeriod"] = d.Get("preferred_backup_period")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_backup_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"]))

	return resourceAliCloudGpdbBackupPolicyRead(d, meta)
}

func resourceAliCloudGpdbBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbBackupPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_backup_policy DescribeGpdbBackupPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_retention_period", objectRaw["BackupRetentionPeriod"])
	d.Set("enable_recovery_point", objectRaw["EnableRecoveryPoint"])
	d.Set("preferred_backup_period", objectRaw["PreferredBackupPeriod"])
	d.Set("preferred_backup_time", objectRaw["PreferredBackupTime"])
	d.Set("recovery_point_period", objectRaw["RecoveryPointPeriod"])

	d.Set("db_instance_id", d.Id())

	return nil
}

func resourceAliCloudGpdbBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyBackupPolicy"
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	if d.HasChange("enable_recovery_point") {
		update = true
		request["EnableRecoveryPoint"] = d.Get("enable_recovery_point")
	}

	if d.HasChange("backup_retention_period") {
		update = true
		request["BackupRetentionPeriod"] = d.Get("backup_retention_period")
	}

	if d.HasChange("preferred_backup_time") {
		update = true
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	if d.HasChange("recovery_point_period") {
		update = true
		request["RecoveryPointPeriod"] = d.Get("recovery_point_period")
	}

	if d.HasChange("preferred_backup_period") {
		update = true
	}
	request["PreferredBackupPeriod"] = d.Get("preferred_backup_period")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudGpdbBackupPolicyRead(d, meta)
}

func resourceAliCloudGpdbBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Backup Policy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
