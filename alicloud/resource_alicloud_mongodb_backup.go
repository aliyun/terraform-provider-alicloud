// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudMongodbBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbBackupCreate,
		Read:   resourceAliCloudMongodbBackupRead,
		Delete: resourceAliCloudMongodbBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(23 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Snapshot", "Physical", "Logical"}, false),
			},
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_db_names": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_intranet_download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudMongodbBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateBackup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("backup_method"); ok {
		request["BackupMethod"] = v
	}
	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_backup", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", request["DBInstanceId"]))

	mongodbServiceV2 := MongodbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, mongodbServiceV2.DescribeAsyncMongodbBackupStateRefreshFunc(d, response, "#$.Backups.Backup[*].BackupId", []string{}))
	jobDetail, err := stateConf.WaitForState()
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	// Extract BackupId from the async polling result. The BackupId in the
	// CreateBackup response is deprecated and may be empty.
	backupId, _ := jsonpath.Get("$.Backups.Backup[0].BackupId", jobDetail)
	if backupId == nil {
		return fmt.Errorf("failed to get BackupId after create")
	}
	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], backupId))

	return resourceAliCloudMongodbBackupRead(d, meta)
}

func resourceAliCloudMongodbBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}

	objectRaw, err := mongodbServiceV2.DescribeMongodbBackup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_backup DescribeMongodbBackup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_db_names", objectRaw["BackupDBNames"])
	d.Set("backup_download_url", objectRaw["BackupDownloadURL"])
	d.Set("backup_intranet_download_url", objectRaw["BackupIntranetDownloadURL"])
	d.Set("backup_job_id", objectRaw["BackupJobId"])
	d.Set("backup_method", objectRaw["BackupMethod"])
	d.Set("backup_mode", objectRaw["BackupMode"])
	d.Set("backup_size", objectRaw["BackupSize"])
	d.Set("backup_type", objectRaw["BackupType"])
	d.Set("status", objectRaw["BackupStatus"])
	d.Set("backup_id", objectRaw["BackupId"])
	d.Set("backup_start_time", objectRaw["BackupStartTime"])
	d.Set("backup_end_time", objectRaw["BackupEndTime"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudMongodbBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteBackup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BackupId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
