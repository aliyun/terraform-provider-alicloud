package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsInstanceCrossBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsInstanceCrossBackupPolicyCreate,
		Read:   resourceAlicloudRdsInstanceCrossBackupPolicyRead,
		Update: resourceAlicloudRdsInstanceCrossBackupPolicyUpdate,
		Delete: resourceAlicloudRdsInstanceCrossBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"backup_enabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_backup_enabled": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disabled"}, false),
			},
			"cross_backup_region": {
				Type:     schema.TypeString,
				Required: true,
			},

			"retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(7, 1825),
			},
			"backup_enabled_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_backup_enabled_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsInstanceCrossBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("instance_id").(string))

	return resourceAlicloudRdsInstanceCrossBackupPolicyUpdate(d, meta)
}

func resourceAlicloudRdsInstanceCrossBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeInstanceCrossBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", d.Id())
	d.Set("backup_enabled", object["BackupEnabled"])
	d.Set("log_backup_enabled", object["LogBackupEnabled"])
	d.Set("cross_backup_region", object["CrossBackupRegion"])
	d.Set("retention", formatInt(object["Retention"]))
	d.Set("backup_enabled_time", object["BackupEnabledTime"])
	d.Set("log_backup_enabled_time", object["LogBackupEnabledTime"])
	d.Set("db_instance_status", object["DBInstanceStatus"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("retent_type", object["RetentType"])
	d.Set("cross_backup_type", object["CrossBackupType"])
	return nil
}

func resourceAlicloudRdsInstanceCrossBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChanges("log_backup_enabled", "cross_backup_region", "retention") {
		action := "ModifyInstanceCrossBackupPolicy"
		request := map[string]interface{}{
			"RegionId":          client.RegionId,
			"DBInstanceId":      d.Id(),
			"DBProxyEndpointId": d.Get("cross_backup_region"),
			"CrossBackupType":   "1",
			"BackupEnabled":     "1",
			"RetentType":        1,
			"SourceIp":          client.SourceIp,
		}
		logBackupEnabled := d.Get("log_backup_enabled").(string)
		if logBackupEnabled == "Enable" {
			request["LogBackupEnabled"] = 1
		}
		if logBackupEnabled == "Disabled" {
			request["LogBackupEnabled"] = 0
		}
		if v, ok := d.GetOk("cross_backup_region"); ok {
			request["CrossBackupRegion"] = v
		}
		if v, ok := d.GetOk("retention"); ok {
			request["Retention"] = v
		}
		if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudRdsInstanceCrossBackupPolicyRead(d, meta)
}

func resourceAlicloudRdsInstanceCrossBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeInstanceCrossBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["LockMode"]; ok && v.(string) == "null" {
		d.SetId("")
		return nil
	}
	action := "ModifyInstanceCrossBackupPolicy"
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"DBInstanceId":  d.Id(),
		"BackupEnabled": "0",
		"SourceIp":      client.SourceIp,
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
