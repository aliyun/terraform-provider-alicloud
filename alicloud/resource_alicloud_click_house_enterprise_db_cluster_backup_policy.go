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

func resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyRead,
		Update: resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyDelete,
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
				Required: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateBackupPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["BackupRetentionPeriod"] = d.Get("backup_retention_period")
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	request["PreferredBackupPeriod"] = d.Get("preferred_backup_period")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster_backup_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"]))

	return resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDbClusterBackupPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster_backup_policy DescribeClickHouseEnterpriseDbClusterBackupPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_retention_period", objectRaw["BackupRetentionPeriod"])
	d.Set("preferred_backup_period", objectRaw["PreferredBackupPeriod"])
	d.Set("preferred_backup_time", objectRaw["PreferredBackupTime"])

	d.Set("db_instance_id", d.Id())

	return nil
}

func resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyBackupPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("backup_retention_period") {
		update = true
	}
	request["BackupRetentionPeriod"] = d.Get("backup_retention_period")
	if d.HasChange("preferred_backup_time") {
		update = true
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	if d.HasChange("preferred_backup_period") {
		update = true
	}
	request["PreferredBackupPeriod"] = d.Get("preferred_backup_period")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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

	return resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
