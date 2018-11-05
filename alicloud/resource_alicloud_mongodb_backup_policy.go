package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudMongoDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBBackupPolicyCreate,
		Read:   resourceAlicloudMongoDBBackupPolicyRead,
		Update: resourceAlicloudMongoDBBackupPolicyUpdate,
		Delete: resourceAlicloudMongoDBBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_time": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(BACKUP_TIME),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},
			"preferred_backup_period": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMongoDBBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	aliyunClient := meta.(*connectivity.AliyunClient)
	client := MongoDBService{aliyunClient}
	request := CommonRequestInit(aliyunClient.RegionId, MONGODBCode, MongoDBDomain)
	request.RegionId = aliyunClient.RegionId
	request.QueryParams["DBInstanceId"] = d.Get("instance_id").(string)
	request.QueryParams["PreferredBackupTime"] = d.Get("preferred_backup_time").(string)
	periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	request.QueryParams["PreferredBackupPeriod"] = backupPeriod

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := client.ModifyMongoDBBackupPolicy(request, aliyunClient); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Create security whitelist ips got an error: %#v", err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.QueryParams["DBInstanceId"], COLON_SEPARATED, resource.UniqueId()))
	return resourceAlicloudMongoDBBackupPolicyRead(d, meta)
}

func resourceAlicloudMongoDBBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	aliyunClient := meta.(*connectivity.AliyunClient)
	client := MongoDBService{aliyunClient}
	instanceID := strings.Split(d.Id(), COLON_SEPARATED)[0]

	request := CommonRequestInit(aliyunClient.RegionId, MONGODBCode, MongoDBDomain)
	request.RegionId = aliyunClient.RegionId
	request.QueryParams["DBInstanceId"] = instanceID
	policy, err := client.DescribeMongoDBBackupPolicy(request, aliyunClient)
	if err != nil {
		if client.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe MongoDB Backup Policy: %#v", err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", instanceID)
	d.Set("preferred_backup_time", policy.PreferredBackupTime)
	d.Set("preferred_backup_period", strings.Split(policy.PreferredBackupPeriod, ","))

	return nil
}

func resourceAlicloudMongoDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	aliyunClient := meta.(*connectivity.AliyunClient)
	client := MongoDBService{aliyunClient}
	request := CommonRequestInit(aliyunClient.RegionId, MONGODBCode, MongoDBDomain)
	request.RegionId = aliyunClient.RegionId
	request.QueryParams["DBInstanceId"] = strings.Split(d.Id(), COLON_SEPARATED)[0]
	update := false

	if d.HasChange("preferred_backup_time") {
		update = true
		request.QueryParams["PreferredBackupTime"] = d.Get("preferred_backup_time").(string)

	}

	if d.HasChange("preferred_backup_period") {
		update = true
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		request.QueryParams["PreferredBackupPeriod"] = backupPeriod
	}

	if update {
		if err := client.ModifyMongoDBBackupPolicy(request, aliyunClient); err != nil {
			return err
		}
	}

	return resourceAlicloudMongoDBBackupPolicyRead(d, meta)
}

func resourceAlicloudMongoDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no explicit delete, only update with modified security ips
	return nil
}
