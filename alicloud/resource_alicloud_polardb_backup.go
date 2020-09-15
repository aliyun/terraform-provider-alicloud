package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBBackupCreate,
		Read:   resourceAlicloudPolarDBBackupRead,
		Delete: resourceAlicloudPolarDBBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_status": {
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
		},
	}
}

func resourceAlicloudPolarDBBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request := polardb.CreateCreateBackupRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	currentTime := time.Now().UTC().Format(time.RFC3339)
	startTime := strings.ReplaceAll(timeFormatWithoutSecond(currentTime), ":", "-")
	_, err := client.WithPolarDBClient(func(polarClient *polardb.Client) (interface{}, error) {
		return polarClient.CreateBackup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_backup", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var backupId string
	if backupId, err = polarDBService.WaitForPolarDBBackupFinished(request.DBClusterId, currentTime, 3600); err != nil {
		return WrapError(err)
	}
	currentTime = time.Now().UTC().Format(TimeFormat)
	endTime := strings.ReplaceAll(currentTime, ":", "-")
	d.SetId(fmt.Sprintf("%s%s%s%s%s%s%s", request.DBClusterId, COLON_SEPARATED, startTime, COLON_SEPARATED, backupId, COLON_SEPARATED, endTime))

	return resourceAlicloudPolarDBBackupRead(d, meta)
}

func resourceAlicloudPolarDBBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]

	backup, err := polardbService.DescribePolarDBBackup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", dbClusterId)
	d.Set("backup_id", backup.BackupId)
	d.Set("backup_status", backup.BackupStatus)
	d.Set("backup_start_time", backup.BackupStartTime)
	d.Set("backup_end_time", backup.BackupEndTime)

	return nil
}

func resourceAlicloudPolarDBBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	backupId := parts[2]
	request := polardb.CreateDeleteBackupRequest()
	request.DBClusterId = dbClusterId
	request.BackupId = backupId

	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteBackup(request)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil
}
