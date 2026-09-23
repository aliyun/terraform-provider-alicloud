package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDbsBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbsBackupPlanCreate,
		Read:   resourceAlicloudDbsBackupPlanRead,
		Update: resourceAlicloudDbsBackupPlanUpdate,
		Delete: resourceAlicloudDbsBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_log_interval_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_method": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"duplication", "logical", "physical"}, false),
			},
			"backup_objects": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_rate_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 1825),
			},
			"backup_speed_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system"}, false),
			},
			"backup_strategy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"simple", "manual"}, false),
			},
			"cross_aliyun_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cross_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"database_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DRDS", "FIle", "MSSQL", "MariaDB", "MongoDB", "MySQL", "Oracle", "PPAS", "PostgreSQL", "Redis"}, false),
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duplication_archive_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"duplication_infrequent_access_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_backup_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"large", "medium", "micro", "small", "xlarge"}, false),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDS", "PolarDB", "DDS", "Kvstore", "Other"}, false),
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Year", "Month"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"used_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 11),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDS", "ECS", "Express", "Agent", "DDS", "Other"}, false),
			},
			"source_endpoint_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_endpoint_sid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"pause", "running"}, false),
			},
			"storage_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDbsBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAndStartBackupPlan"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("backup_gateway_id"); ok {
		request["BackupGatewayId"] = v
	}
	if v, ok := d.GetOk("backup_log_interval_seconds"); ok {
		request["BackupLogIntervalSeconds"] = v
	}
	request["BackupMethod"] = d.Get("backup_method")
	if v, ok := d.GetOk("backup_objects"); ok {
		request["BackupObjects"] = v
	}
	if v, ok := d.GetOk("backup_period"); ok {
		request["BackupPeriod"] = v
	}
	request["BackupPlanName"] = d.Get("backup_plan_name")
	if v, ok := d.GetOk("backup_rate_limit"); ok {
		request["BackupRateLimit"] = v
	}
	if v, ok := d.GetOk("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	if v, ok := d.GetOk("backup_speed_limit"); ok {
		request["BackupSpeedLimit"] = v
	}
	if v, ok := d.GetOk("backup_start_time"); ok {
		request["BackupStartTime"] = v
	}
	if v, ok := d.GetOk("backup_storage_type"); ok {
		request["BackupStorageType"] = v
	}
	if v, ok := d.GetOk("backup_strategy_type"); ok {
		request["BackupStrategyType"] = v
	}
	if v, ok := d.GetOk("cross_aliyun_id"); ok {
		request["CrossAliyunId"] = v
	}
	if v, ok := d.GetOk("cross_role_name"); ok {
		request["CrossRoleName"] = v
	}
	if v, ok := d.GetOk("database_region"); ok {
		request["DatabaseRegion"] = v
	}
	request["DatabaseType"] = d.Get("database_type")
	if v, ok := d.GetOk("duplication_archive_period"); ok {
		request["DuplicationArchivePeriod"] = v
	}
	if v, ok := d.GetOk("duplication_infrequent_access_period"); ok {
		request["DuplicationInfrequentAccessPeriod"] = v
	}
	if v, ok := d.GetOkExists("enable_backup_log"); ok {
		request["EnableBackupLog"] = v
	}
	request["InstanceClass"] = d.Get("instance_class")
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok {
		request["OSSBucketName"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertDbsBackupPlanPaymentTypeRequest(v)
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		request["SourceEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["Region"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		request["SourceEndpointDatabaseName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		request["SourceEndpointIP"] = v
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		request["SourceEndpointInstanceID"] = v
	}
	request["SourceEndpointInstanceType"] = d.Get("source_endpoint_instance_type")
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		request["SourceEndpointPassword"] = v
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		request["SourceEndpointPort"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceEndpointRegion"] = v
	}
	if v, ok := d.GetOk("source_endpoint_sid"); ok {
		request["SourceEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		request["SourceEndpointUserName"] = v
	}
	if v, ok := d.GetOk("storage_region"); ok {
		request["StorageRegion"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	request["ClientToken"] = buildClientToken("CreateAndStartBackupPlan")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dbs_backup_plan", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BackupPlanId"]))
	dbsService := DbsService{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbsService.DbsBackupPlanStateRefreshFunc(d.Id(), []string{"stop"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDbsBackupPlanUpdate(d, meta)
}
func resourceAlicloudDbsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbsService := DbsService{client}
	object, err := dbsService.DescribeDbsBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbs_backup_plan dbsService.DescribeDbsBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["BackupGatewayId"]; ok && fmt.Sprint(object["BackupGatewayId"]) != "0" {
		d.Set("backup_gateway_id", fmt.Sprint(v))
	}
	d.Set("backup_method", object["BackupMethod"])
	d.Set("backup_start_time", object["BackupStartTime"])
	d.Set("backup_storage_type", object["BackupStorageType"])
	d.Set("cross_aliyun_id", object["CrossAliyunId"])
	d.Set("cross_role_name", object["CrossRoleName"])
	d.Set("database_type", object["DatabaseType"])
	d.Set("duplication_archive_period", object["DuplicationArchivePeriod"])
	d.Set("duplication_infrequent_access_period", object["DuplicationInfrequentAccessPeriod"])
	d.Set("enable_backup_log", object["EnableBackupLog"])
	d.Set("instance_class", object["InstanceClass"])
	d.Set("oss_bucket_name", object["OSSBucketName"])
	d.Set("source_endpoint_database_name", object["SourceEndpointDatabaseName"])
	d.Set("source_endpoint_instance_id", object["SourceEndpointInstanceID"])
	d.Set("source_endpoint_instance_type", convertDbsSourceEndpointInstanceTypeResponse(object["SourceEndpointInstanceType"]))
	d.Set("source_endpoint_region", object["SourceEndpointRegion"])
	d.Set("source_endpoint_sid", object["SourceEndpointOracleSID"])
	d.Set("source_endpoint_user_name", object["SourceEndpointUserName"])
	d.Set("status", object["BackupPlanStatus"])
	d.Set("backup_objects", object["BackupObjects"])
	d.Set("backup_period", convertDbsBackupPeriodResponse(object["BackupPeriod"]))
	d.Set("backup_plan_name", object["BackupPlanName"])
	d.Set("backup_retention_period", object["BackupRetentionPeriod"])
	d.Set("source_endpoint_oracle_sid", object["SourceEndpointOracleSID"])

	describeBackupPlanBillingObject, err := dbsService.DescribeBackupPlanBilling(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("payment_type", convertDbsBackupPlanPaymentTypeResponse(describeBackupPlanBillingObject["BuyChargeType"]))
	return nil
}
func resourceAlicloudDbsBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbsService := DbsService{client}
	var response map[string]interface{}
	if d.HasChange("status") {
		object, err := dbsService.DescribeDbsBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["BackupPlanStatus"].(string) != target {
			if target == "pause" {
				request := map[string]interface{}{
					"BackupPlanId": d.Id(),
				}
				request["StopMethod"] = "ALL"
				action := "StopBackupPlan"
				request["ClientToken"] = buildClientToken("StopBackupPlan")
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
				stateConf := BuildStateConf([]string{}, []string{"pause"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbsService.DbsBackupPlanStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "running" {
				request := map[string]interface{}{
					"BackupPlanId": d.Id(),
				}
				action := "StartBackupPlan"
				request["ClientToken"] = buildClientToken("StartBackupPlan")
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
				stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbsService.DbsBackupPlanStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}
	return resourceAlicloudDbsBackupPlanRead(d, meta)
}
func resourceAlicloudDbsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbsService := DbsService{client}
	action := "ReleaseBackupPlan"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"BackupPlanId": d.Id(),
	}

	request["ClientToken"] = buildClientToken("ReleaseBackupPlan")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidJobId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dbsService.DbsBackupPlanStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertDbsBackupPlanPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "postpay"
	case "Subscription":
		return "prepay"
	}
	return source
}
func convertDbsBackupPlanPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertDbsSourceEndpointInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "rds":
		return "RDS"
	case "ecs":
		return "ECS"
	case "express":
		return "Express"
	case "agent":
		return "Agent"
	case "dds":
		return "DDS"
	case "other":
		return "Other"
	}
	return source
}

func convertDbsBackupPeriodResponse(source interface{}) interface{} {
	switch source {
	case "MONDAY":
		return "Monday"
	case "TUESDAY":
		return "Tuesday"
	case "WEDNESDAY":
		return "Wednesday"
	case "THURSDAY":
		return "Thursday"
	case "FRIDAY":
		return "Friday"
	case "SATURDAY":
		return "Saturday"
	case "SUNDAY":
		return "Sunday"

	}
	return source
}
