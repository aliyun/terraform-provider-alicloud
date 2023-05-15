package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
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
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
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
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "LATEST", "NONE"}, false),
			},
			"data_level1_backup_retention_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_retention_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "2/24H", "3/24H", "4/24H"}, false),
			},
			"data_level1_backup_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "2/24H", "3/24H", "4/24H"}, false),
			},
			"data_level1_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"data_level1_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_another_region_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_another_region_retention_period": {
				Type:     schema.TypeString,
				Optional: true,
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
	return nil
}

func resourceAlicloudPolarDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		update = true
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	if d.HasChange("preferred_backup_time") {
		update = true
		if v, ok := d.GetOk("preferred_backup_time"); ok {
			request["PreferredBackupTime"] = v
		}
	}
	if d.HasChange("backup_retention_policy_on_cluster_deletion") {
		update = true
		if v, ok := d.GetOk("backup_retention_policy_on_cluster_deletion"); ok {
			request["BackupRetentionPolicyOnClusterDeletion"] = v
		}
	}
	if d.HasChange("data_level1_backup_retention_period") {
		update = true
		if v, ok := d.GetOk("data_level1_backup_retention_period"); ok {
			request["DataLevel1BackupRetentionPeriod"] = v
		}
	}
	if d.HasChange("data_level2_backup_retention_period") {
		update = true
		if v, ok := d.GetOk("data_level2_backup_retention_period"); ok {
			request["DataLevel2BackupRetentionPeriod"] = v
		}
	}
	if d.HasChange("backup_frequency") {
		update = true
		if v, ok := d.GetOk("backup_frequency"); ok {
			request["BackupFrequency"] = v
		}
	}
	if d.HasChange("data_level1_backup_frequency") {
		update = true
		if v, ok := d.GetOk("data_level1_backup_frequency"); ok {
			request["DataLevel1BackupFrequency"] = v
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
	if d.HasChange("data_level2_backup_another_region_region") {
		update = true
		if v, ok := d.GetOk("data_level2_backup_another_region_region"); ok {
			request["DataLevel2BackupAnotherRegionRegion"] = v
		}
	}
	if d.HasChange("data_level2_backup_another_region_retention_period") {
		update = true
		if v, ok := d.GetOk("data_level2_backup_another_region_retention_period"); ok {
			request["DataLevel2BackupAnotherRegionRetentionPeriod"] = v
		}
	}
	if d.HasChange("backup_retention_period") {
		update = true
		if v, ok := d.GetOk("backup_retention_period"); ok {
			request["DataLevel1BackupRetentionPeriod"] = v
		}
	}

	if update {
		action := "ModifyBackupPolicy"
		conn, err := client.NewPolarDBClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
