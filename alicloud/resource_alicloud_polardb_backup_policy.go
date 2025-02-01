package alicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBBackupPolicyCreate,
		Read:   resourceAlicloudPolarDBBackupPolicyRead,
		Update: resourceAlicloudPolarDBBackupPolicyUpdate,
		Delete: resourceAlicloudPolarDBBackupPolicyDelete,
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
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
				},
				Optional: true,
				Computed: true,
			},
			"backup_retention_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ALL", "LATEST", "NONE"}, false),
			},
			"data_level1_backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(3, 30),
				Computed:     true,
			},
			"data_level2_backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					IntInSlice([]int{-1, 0}),
					IntBetween(30, 7300)),
				Computed: true,
			},
			"backup_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "2/24H", "3/24H", "4/24H"}, false),
				Computed:     true,
			},
			"data_level1_backup_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "2/24H", "3/24H", "4/24H"}, false),
				Computed:     true,
			},
			"data_level1_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"data_level1_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
				},
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
				},
				Optional: true,
			},
			"data_level2_backup_another_region_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_level2_backup_another_region_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					IntInSlice([]int{-1, 0}),
					IntBetween(30, 7300)),
				Computed: true,
			},
			"log_backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					IntInSlice([]int{-1}),
					IntBetween(3, 7300)),
				Computed: true,
			},
			"log_backup_another_region_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_backup_another_region_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					IntInSlice([]int{-1, 0}),
					IntBetween(30, 7300)),
			},
			"enable_backup_log": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("db_cluster_id").(string))

	return resourceAlicloudPolarDBBackupPolicyUpdate(d, meta)
}

func resourceAlicloudPolarDBBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbService := PolarDBService{client}
	object, err := polardbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	if "" != object["PreferredBackupPeriod"].(string) {
		preferredBackupPeriods := strings.Split(object["PreferredBackupPeriod"].(string), ",")
		d.Set("preferred_backup_period", preferredBackupPeriods)
	}
	d.Set("backup_retention_period", object["DataLevel1BackupRetentionPeriod"])
	d.Set("preferred_backup_time", object["PreferredBackupTime"])
	d.Set("backup_retention_policy_on_cluster_deletion", object["BackupRetentionPolicyOnClusterDeletion"])
	d.Set("data_level1_backup_retention_period", object["DataLevel1BackupRetentionPeriod"])
	d.Set("data_level2_backup_retention_period", object["DataLevel2BackupRetentionPeriod"])
	d.Set("backup_frequency", object["BackupFrequency"])
	d.Set("data_level1_backup_frequency", object["DataLevel1BackupFrequency"])
	d.Set("data_level1_backup_time", object["DataLevel1BackupTime"])
	if "" != object["DataLevel1BackupPeriod"].(string) {
		dataLevel1BackupPeriods := strings.Split(object["DataLevel1BackupPeriod"].(string), ",")
		d.Set("data_level1_backup_period", dataLevel1BackupPeriods)
	}
	if "" != object["DataLevel2BackupPeriod"].(string) {
		dataLevel2BackupPeriods := strings.Split(object["DataLevel2BackupPeriod"].(string), ",")
		d.Set("data_level2_backup_period", dataLevel2BackupPeriods)
	}
	d.Set("data_level2_backup_another_region_region", object["DataLevel2BackupAnotherRegionRegion"])
	d.Set("data_level2_backup_another_region_retention_period", object["DataLevel2BackupAnotherRegionRetentionPeriod"])
	logBackupPolicy, err := polardbService.DescribeLogBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("log_backup_retention_period", logBackupPolicy["LogBackupRetentionPeriod"])
	d.Set("enable_backup_log", logBackupPolicy["EnableBackupLog"])

	return nil
}

func resourceAlicloudPolarDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		update = true
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	if d.HasChanges("backup_retention_policy_on_cluster_deletion", "data_level2_backup_retention_period", "backup_frequency",
		"data_level1_backup_frequency", "data_level2_backup_another_region_region", "data_level2_backup_another_region_retention_period") {
		update = true

		if v, ok := d.GetOk("backup_retention_policy_on_cluster_deletion"); ok {
			request["BackupRetentionPolicyOnClusterDeletion"] = v
		}
		if v, ok := d.GetOk("data_level2_backup_retention_period"); ok {
			dataLevel2BackupRetentionPeriod := v.(int)
			request["DataLevel2BackupRetentionPeriod"] = strconv.Itoa(dataLevel2BackupRetentionPeriod)
		}
		if v, ok := d.GetOk("backup_frequency"); ok {
			request["BackupFrequency"] = v
		}
		if v, ok := d.GetOk("data_level1_backup_frequency"); ok {
			request["DataLevel1BackupFrequency"] = v
		}
		if v, ok := d.GetOk("data_level2_backup_another_region_region"); ok {
			request["DataLevel2BackupAnotherRegionRegion"] = v
		}
		if v, ok := d.GetOk("data_level2_backup_another_region_retention_period"); ok {
			dataLevel2BackupAnotherRegionRetentionPeriod := v.(int)
			request["DataLevel2BackupAnotherRegionRetentionPeriod"] = strconv.Itoa(dataLevel2BackupAnotherRegionRetentionPeriod)
		}
	}
	if d.HasChange("backup_retention_period") {
		update = true
		if v, ok := d.GetOk("backup_retention_period"); ok {
			request["DataLevel1BackupRetentionPeriod"] = v
		}
	}
	if d.HasChange("data_level1_backup_retention_period") {
		update = true
		if v, ok := d.GetOk("data_level1_backup_retention_period"); ok {
			dataLevel1BackupRetentionPeriod := v.(int)
			request["DataLevel1BackupRetentionPeriod"] = strconv.Itoa(dataLevel1BackupRetentionPeriod)
		}
	}
	if d.HasChange("preferred_backup_time") {
		update = true
		if v, ok := d.GetOk("preferred_backup_time"); ok {
			request["PreferredBackupTime"] = v
		}
	}
	if d.HasChange("data_level1_backup_time") {
		update = true
		if v, ok := d.GetOk("data_level1_backup_time"); ok {
			request["DataLevel1BackupTime"] = v
		}
	}
	if d.HasChange("data_level1_backup_period") {
		update = true
		periodList := expandStringList(d.Get("data_level1_backup_period").(*schema.Set).List())
		request["DataLevel1BackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	if d.HasChange("data_level2_backup_period") {
		update = true
		periodList := expandStringList(d.Get("data_level2_backup_period").(*schema.Set).List())
		request["DataLevel2BackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	if update {
		action := "ModifyBackupPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
	if d.HasChanges("log_backup_retention_period", "log_backup_another_region_region", "log_backup_another_region_retention_period") {
		if v, ok := d.GetOk("log_backup_retention_period"); ok {
			logBackupRetentionPeriod := v.(int)
			request["LogBackupRetentionPeriod"] = strconv.Itoa(logBackupRetentionPeriod)
		}
		if v, ok := d.GetOk("log_backup_another_region_region"); ok {
			request["LogBackupAnotherRegionRegion"] = v
		}
		if v, ok := d.GetOk("log_backup_another_region_retention_period"); ok {
			logBackupAnotherRegionRetentionPeriod := v.(int)
			request["LogBackupAnotherRegionRetentionPeriod"] = strconv.Itoa(logBackupAnotherRegionRetentionPeriod)
		}

		action := "ModifyLogBackupPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
	return resourceAlicloudPolarDBBackupPolicyRead(d, meta)
}

func resourceAlicloudPolarDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	//  Terraform can not destroy it..
	return nil
}
