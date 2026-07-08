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
			"cluster_backup_id": {
				Type:     schema.TypeString,
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"BackupJobExists"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_redis_backup", action, AlibabaCloudSdkGoERROR)
	}

	// Get BackupJobID from CreateBackup response
	backupJobID, ok := response["BackupJobID"]
	if !ok {
		return WrapErrorf(Error("BackupJobID not found in response"), DefaultErrorMsg, "alicloud_redis_backup", action, AlibabaCloudSdkGoERROR)
	}

	// Set temporary ID with instance_id:job_id format for waiting
	instanceID := request["InstanceId"].(string)
	d.SetId(fmt.Sprintf("%v:%v", instanceID, backupJobID))

	// Wait for backup job to complete and get BackupId
	redisServiceV2 := RedisServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, redisServiceV2.DescribeAsyncRedisBackupStateRefreshFunc(d, response, "#$.Backups.Backup[*].BackupId", []string{}))
	jobDetail, err := stateConf.WaitForState()
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// Get the actual BackupId from the job result
	if backupDetail, ok := jobDetail.(map[string]interface{}); ok {
		if backups, err := jsonpath.Get("$.Backups.Backup", backupDetail); err == nil {
			if backupList, ok := backups.([]interface{}); ok && len(backupList) > 0 {
				if backup, ok := backupList[0].(map[string]interface{}); ok {
					if backupID, ok := backup["BackupId"]; ok {
						// Update ID with the actual BackupId
						d.SetId(fmt.Sprintf("%v:%v", instanceID, backupID))

						// Resolve the owning cluster backup set id (cb-*) for cluster-architecture
						// instances. Standalone instances legitimately return an empty result — that
						// is not an error. A real API failure (network / auth / bad parameter) must
						// surface: a silent empty cluster_backup_id would break downstream clone
						// resources with an unactionable InvalidParameter.Missing.
						clusterBackupID, cbErr := redisServiceV2.DescribeRedisClusterBackupId(instanceID, fmt.Sprint(backupID))
						if cbErr != nil {
							return WrapErrorf(cbErr, DefaultErrorMsg, d.Id(), "DescribeClusterBackupList", AlibabaCloudSdkGoERROR)
						}
						if clusterBackupID != "" {
							d.Set("cluster_backup_id", clusterBackupID)
						} else {
							log.Printf("[DEBUG] alicloud_redis_backup: no cluster backup set for %s (standalone instance)", d.Id())
						}
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
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_backup DescribeRedisBackup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	redisServiceV2 := RedisServiceV2{client}
	parts := strings.Split(d.Id(), ":")
	instanceID := parts[0]

	// Cluster-architecture instances hold one BackupId per shard under a single cluster
	// backup set id (cb-*). R-kvstore has no cluster-level DeleteClusterBackup API, so we
	// have to fan out shard-level DeleteBackup calls; otherwise only the state-recorded
	// shard is removed and the rest linger. Standalone instances (cluster_backup_id empty
	// or unset — including pre-fix state that never populated the field) take the single
	// shard path unchanged.
	backupIDs := []string{parts[1]}
	if clusterBackupID, ok := d.Get("cluster_backup_id").(string); ok && clusterBackupID != "" {
		shardIDs, listErr := redisServiceV2.ListClusterBackupShardIds(instanceID, clusterBackupID)
		if listErr != nil {
			return WrapErrorf(listErr, DefaultErrorMsg, d.Id(), "DescribeClusterBackupList", AlibabaCloudSdkGoERROR)
		}
		if len(shardIDs) > 0 {
			backupIDs = shardIDs
		}
	}

	action := "DeleteBackup"
	query := make(map[string]interface{})
	for _, backupID := range backupIDs {
		request := map[string]interface{}{
			"BackupId":   backupID,
			"InstanceId": instanceID,
			"RegionId":   client.RegionId,
		}

		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
				continue
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}
