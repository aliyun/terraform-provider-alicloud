package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudPolarDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create:        resourceAlicloudPolarDBBackupPolicyCreate,
		Read:          resourceAlicloudPolarDBBackupPolicyRead,
		Update:        resourceAlicloudPolarDBBackupPolicyUpdate,
		Delete:        resourceAlicloudPolarDBBackupPolicyDelete,
		CustomizeDiff: validatePolarDBAdvancedBackupPolicyDiff,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backup_policy_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "Advanced"}, false),
			},
			"advanced_policy_option": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advanced_data_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     polarDBAdvancedPolicySchema(),
			},
			"auto_created_advanced_data_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     polarDBObservedAdvancedPolicySchema(),
			},
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
	if level, ok := object["BackupPolicyLevel"]; ok {
		if err = d.Set("backup_policy_level", level); err != nil {
			return WrapError(err)
		}
	}
	if option, ok := object["AdvancedPolicyOption"]; ok {
		if err = d.Set("advanced_policy_option", option); err != nil {
			return WrapError(err)
		}
	}
	observedPolicies := polarDBDescribeAdvancedPolicies(object)
	managedPolicies, autoCreatedPolicies := reconcilePolarDBAdvancedPolicies(
		polarDBAdvancedPolicyList(d.Get("advanced_data_policies")), observedPolicies,
	)
	if err = d.Set("advanced_data_policies", managedPolicies); err != nil {
		return WrapError(err)
	}
	if err = d.Set("auto_created_advanced_data_policies", autoCreatedPolicies); err != nil {
		return WrapError(err)
	}
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
	if d.Get("backup_policy_level") == "Advanced" {
		if d.HasChange("backup_policy_level") {
			service := PolarDBService{client}
			current, err := service.DescribeBackupPolicy(d.Id())
			if err != nil {
				return WrapError(err)
			}
			if current["AdvancedPolicyOption"] == "notSupport" {
				return fmt.Errorf("cluster %s does not support advanced backup policies", d.Id())
			}
		}
		request := map[string]interface{}{
			"DBClusterId":       d.Id(),
			"BackupPolicyLevel": "Advanced",
		}
		previousValue, desiredValue := d.GetChange("advanced_data_policies")
		previousPolicies := polarDBAdvancedPolicyList(previousValue)
		desiredPolicies := polarDBAdvancedPolicyList(desiredValue)
		if err := preparePolarDBAdvancedPolicyRegions(previousPolicies, client.RegionId, false); err != nil {
			return WrapError(err)
		}
		if err := preparePolarDBAdvancedPolicyRegions(desiredPolicies, client.RegionId, true); err != nil {
			return WrapError(err)
		}
		for index, policy := range desiredPolicies {
			rawValue, diagnostics := d.GetRawConfigAt(cty.Path{
				cty.GetAttrStep{Name: "advanced_data_policies"},
				cty.IndexStep{Key: cty.NumberIntVal(int64(index))},
				cty.GetAttrStep{Name: "only_preserve_one_each_day"},
			})
			policy["only_preserve_one_each_day_set"] = !diagnostics.HasError() && polarDBAdvancedPolicyBoolConfigured(rawValue)
		}
		actions := polarDBAdvancedPolicyActions(
			previousPolicies,
			desiredPolicies,
		)
		if d.HasChange("backup_retention_policy_on_cluster_deletion") {
			if value, ok := d.GetOk("backup_retention_policy_on_cluster_deletion"); ok {
				request["BackupRetentionPolicyOnClusterDeletion"] = value
			}
		}
		if d.HasChange("backup_policy_level") || len(actions) > 0 || len(request) > 2 {
			if len(actions) > 0 {
				if err := expandPolarDBAdvancedPolicyActions(request, actions); err != nil {
					return WrapError(err)
				}
			}
			if err := callPolarDBModifyBackupPolicy(client, request, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyBackupPolicy", AlibabaCloudSdkGoERROR)
			}
			if err := waitPolarDBAdvancedPolicyActions(client, d.Id(), actions, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeBackupPolicy", AlibabaCloudSdkGoERROR)
			}
		}
		if err := d.Set("advanced_data_policies", polarDBAdvancedPolicyIDsInDesiredOrder(previousPolicies, desiredPolicies)); err != nil {
			return WrapError(err)
		}
		return resourceAlicloudPolarDBBackupPolicyRead(d, meta)
	}
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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	// A cluster's backup policy cannot be removed. Delete only the custom advanced
	// rules managed by this resource; the cluster and its backup level remain.
	if d.Get("backup_policy_level") == "Advanced" {
		policies := polarDBAdvancedPolicyList(d.Get("advanced_data_policies"))
		if len(policies) > 0 {
			actions := make([]polarDBAdvancedPolicyAction, 0, len(policies))
			for _, policy := range policies {
				if id := strings.TrimSpace(fmt.Sprint(policy["policy_id"])); id != "" && id != "<nil>" {
					actions = append(actions, polarDBAdvancedPolicyAction{Action: "DELETE", Policy: policy})
				}
			}
			if len(actions) > 0 {
				request := map[string]interface{}{
					"DBClusterId":       d.Id(),
					"BackupPolicyLevel": "Advanced",
				}
				if err := expandPolarDBAdvancedPolicyActions(request, actions); err != nil {
					return WrapError(err)
				}
				client := meta.(*connectivity.AliyunClient)
				if err := callPolarDBModifyBackupPolicy(client, request, d.Timeout(schema.TimeoutDelete)); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyBackupPolicy", AlibabaCloudSdkGoERROR)
				}
				if err := waitPolarDBAdvancedPolicyActions(client, d.Id(), actions, d.Timeout(schema.TimeoutDelete)); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeBackupPolicy", AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	return nil
}
