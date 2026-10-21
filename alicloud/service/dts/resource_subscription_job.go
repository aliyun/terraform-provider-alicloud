package dts

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSubscriptionJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionJobCreate,
		Read:   resourceSubscriptionJobRead,
		Update: ResourceSubscriptionJobUpdate,
		Delete: resourceSubscriptionJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"checkpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"compute_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"db_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delay_notice": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delay_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delay_rule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"ADS", "DB2", "DRDS", "DataHub", "Greenplum", "MSSQL", "MySQL", "PolarDB", "PostgreSQL", "Redis", "Tablestore", "as400", "clickhouse", "kafka", "mongodb", "odps", "oracle", "polardb_o", "polardb_pg", "tidb"}, false),
			},
			"destination_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"error_notice": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"error_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"large", "medium", "micro", "small", "xlarge", "xxlarge"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"payment_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Month", "Year"}, false),
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"reserve": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"MySQL", "Oracle"}, false),
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"CEN", "DRDS", "ECS", "Express", "LocalInstance", "PolarDB", "RDS", "dg"}, false),
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_endpoint_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_ssl": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"0", "1"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Abnormal", "Downgrade", "Locked", "Normal", "NotStarted", "NotStarted", "PreCheckPass", "PrecheckFailed", "Prechecking", "Retrying", "Starting", "Upgrade"}, false),
			},
			"subscription_data_type_ddl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"subscription_data_type_dml": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"subscription_instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"classic", "vpc"}, false),
			},
			"subscription_instance_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscription_instance_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"bidirectional", "oneway"}, false),
			},
			"synchronization_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.TagsSchemaWithElements(),
		},
	}
}

func resourceSubscriptionJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDtsInstance"
	request := make(map[string]interface{})
	request["AutoPay"] = false
	request["AutoStart"] = true
	var err error
	if v, ok := d.GetOk("compute_unit"); ok {
		request["ComputeUnit"] = v
	}
	if v, ok := d.GetOk("database_count"); ok {
		request["DatabaseCount"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request["DestinationEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("destination_region"); ok {
		request["DestinationRegion"] = v
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request["SourceEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceRegion"] = v
	}
	if v, ok := d.GetOk("sync_architecture"); ok {
		request["SyncArchitecture"] = v
	}
	request["Type"] = "SUBSCRIBE"
	if v, ok := d.GetOk("payment_duration"); ok {
		request["UsedTime"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dts_subscription_job", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["JobId"]))
	d.Set("dts_instance_id", response["InstanceId"])

	return ResourceSubscriptionJobUpdate(d, meta)
}
func resourceSubscriptionJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := Service{client}
	object, err := dtsService.DescribeSubscriptionJob(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_subscription_job dtsService.DescribeSubscriptionJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	actualCheckpoint, err := convert.FormatInt(object["Checkpoint"])
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("checkpoint", fmt.Sprint(actualCheckpoint))
	d.Set("db_list", object["DbObject"])
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	d.Set("payment_type", convertPaymentTypeResponse(object["PayType"]))
	d.Set("source_endpoint_database_name", object["SourceEndpoint"].(map[string]interface{})["DatabaseName"])
	d.Set("source_endpoint_engine_name", object["SourceEndpoint"].(map[string]interface{})["EngineName"])
	d.Set("source_endpoint_ip", object["SourceEndpoint"].(map[string]interface{})["Ip"])
	d.Set("source_endpoint_instance_id", object["SourceEndpoint"].(map[string]interface{})["InstanceID"])
	d.Set("source_endpoint_instance_type", object["SourceEndpoint"].(map[string]interface{})["InstanceType"])
	d.Set("source_endpoint_oracle_sid", object["SourceEndpoint"].(map[string]interface{})["OracleSID"])
	d.Set("source_endpoint_owner_id", object["SourceEndpoint"].(map[string]interface{})["AliyunUid"])
	d.Set("source_endpoint_port", object["SourceEndpoint"].(map[string]interface{})["Port"])
	d.Set("source_endpoint_region", object["SourceEndpoint"].(map[string]interface{})["Region"])
	d.Set("source_endpoint_role", object["SourceEndpoint"].(map[string]interface{})["RoleName"])
	d.Set("source_endpoint_user_name", object["SourceEndpoint"].(map[string]interface{})["UserName"])
	if ssl := convertEndpointSslResponse(object["SourceEndpoint"].(map[string]interface{})["SslSolutionEnum"]); ssl != nil {
		d.Set("source_endpoint_ssl", ssl)
	}
	d.Set("status", object["Status"])
	d.Set("subscription_data_type_ddl", object["SubscriptionDataType"].(map[string]interface{})["Ddl"])
	d.Set("subscription_data_type_dml", object["SubscriptionDataType"].(map[string]interface{})["Dml"])

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(object["Reserved"].(string)), &jsonData)
	if jsonData["netType"] != nil {
		d.Set("subscription_instance_network_type", strings.ToLower(jsonData["netType"].(string)))
	}
	d.Set("subscription_instance_vpc_id", jsonData["vpcId"])
	d.Set("subscription_instance_vswitch_id", jsonData["vswitchId"])
	listTagResourcesObject, err := dtsService.ListTagResources(object["DtsInstanceID"].(string), "ALIYUN::DTS::INSTANCE")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}
func ResourceSubscriptionJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := Service{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := dtsService.SetResourceTags(d, "ALIYUN::DTS::INSTANCE:JOB"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if d.HasChange("dts_job_name") {
		update = true
		if v, ok := d.GetOk("dts_job_name"); ok {
			request["DtsJobName"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyDtsJobName"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if isJobNotFoundSignal(err, response) {
			// An instance that carries historical deletion records can lose the server-side
			// mapping between the job ID stored in state and the instance, so ModifyDtsJobName
			// reports that the instance does not exist. Recover the currently valid job ID from
			// the instance and retry the rename once before failing.
			resolvedId, resolveErr := subscriptionJobResolveJobId(d, &dtsService)
			if resolveErr != nil {
				log.Printf("[WARN] Failed to resolve the valid DTS job ID under instance %s: %s", d.Get("dts_instance_id"), resolveErr)
			}
			if resolveErr == nil && resolvedId != "" && resolvedId != d.Id() {
				log.Printf("[INFO] The DTS job ID %s of alicloud_dts_subscription_job is no longer valid, retrying %s with the resolved job ID %s.", d.Id(), action, resolvedId)
				request["DtsJobId"] = resolvedId
				var retryErr error
				retryWait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				retryErr = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, retryErr = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
					if retryErr != nil {
						if helper.NeedRetry(retryErr) {
							retryWait()
							return retry.RetryableError(retryErr)
						}
						return retry.NonRetryableError(retryErr)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if retryErr == nil && fmt.Sprint(response["Success"]) != "false" {
					d.SetId(resolvedId)
					err = nil
				} else if retryErr != nil {
					err = retryErr
				} else {
					err = fmt.Errorf("%s failed, response: %v", action, response)
				}
			}
		}
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	update = false
	modifyDtsJobPasswordReq := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	modifyDtsJobPasswordReq["Endpoint"] = "src"
	modifyDtsJobPasswordReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("source_endpoint_password") {
		update = true
		if v, ok := d.GetOk("source_endpoint_password"); ok {
			modifyDtsJobPasswordReq["Password"] = v
		}
		if v, ok := d.GetOk("source_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq["UserName"] = v
		}
	}
	if update {
		action := "ModifyDtsJobPassword"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, modifyDtsJobPasswordReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDtsJobPasswordReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		target := d.Get("status").(string)
		err = resourceSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, d.Get("status")))
		}
	}

	if !d.IsNewResource() && d.HasChange("status") {
		target := d.Get("status").(string)
		err := resourceSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, d.Get("status")))
		}
	}

	update = false
	if d.IsNewResource() {
		update = true
	}
	configureSubscriptionReq := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if d.HasChange("db_list") {
		update = true
	}
	if v, ok := d.GetOk("db_list"); ok {
		configureSubscriptionReq["DbList"] = v
	}
	if v, ok := d.GetOk("dts_job_name"); ok {
		configureSubscriptionReq["DtsJobName"] = v
	}
	configureSubscriptionReq["RegionId"] = client.RegionId
	if d.HasChange("subscription_instance_network_type") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_network_type"); ok {
		configureSubscriptionReq["SubscriptionInstanceNetworkType"] = v
	}
	if d.HasChange("checkpoint") {
		update = true
	}
	if v, ok := d.GetOk("checkpoint"); ok {
		configureSubscriptionReq["Checkpoint"] = v
	}
	if d.HasChange("source_endpoint_database_name") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		configureSubscriptionReq["SourceEndpointDatabaseName"] = v
	}
	if !d.IsNewResource() && d.HasChange("source_endpoint_engine_name") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		configureSubscriptionReq["SourceEndpointEngineName"] = v
	}
	if d.HasChange("source_endpoint_ip") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		configureSubscriptionReq["SourceEndpointIP"] = v
	}
	if d.HasChange("source_endpoint_instance_id") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		configureSubscriptionReq["SourceEndpointInstanceID"] = v
	}
	if d.HasChange("source_endpoint_instance_type") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_instance_type"); ok {
		configureSubscriptionReq["SourceEndpointInstanceType"] = v
	}
	if d.HasChange("source_endpoint_oracle_sid") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		configureSubscriptionReq["SourceEndpointOracleSID"] = v
	}
	if d.HasChange("source_endpoint_owner_id") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		configureSubscriptionReq["SourceEndpointOwnerID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		configureSubscriptionReq["SourceEndpointPassword"] = v
	}
	if d.HasChange("source_endpoint_port") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		configureSubscriptionReq["SourceEndpointPort"] = v
	}
	if d.HasChange("source_endpoint_region") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		configureSubscriptionReq["SourceEndpointRegion"] = v
	}
	if d.HasChange("source_endpoint_role") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		configureSubscriptionReq["SourceEndpointRole"] = v
	}

	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		configureSubscriptionReq["SourceEndpointUserName"] = v
	}
	if d.HasChange("source_endpoint_ssl") {
		update = true
	}
	if d.HasChange("subscription_data_type_ddl") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("subscription_data_type_ddl"); ok {
		configureSubscriptionReq["SubscriptionDataTypeDDL"] = v
	}
	if d.HasChange("subscription_data_type_dml") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("subscription_data_type_dml"); ok {
		configureSubscriptionReq["SubscriptionDataTypeDML"] = v
	}
	if d.HasChange("subscription_instance_vpc_id") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_vpc_id"); ok {
		configureSubscriptionReq["SubscriptionInstanceVPCId"] = v
	}
	if d.HasChange("subscription_instance_vswitch_id") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_vswitch_id"); ok {
		configureSubscriptionReq["SubscriptionInstanceVSwitchId"] = v
	}
	if update {

		if v, ok := d.GetOkExists("delay_notice"); ok {
			configureSubscriptionReq["DelayNotice"] = v
		}
		if v, ok := d.GetOk("delay_phone"); ok {
			configureSubscriptionReq["DelayPhone"] = v
		}
		if v, ok := d.GetOk("delay_rule_time"); ok {
			configureSubscriptionReq["DelayRuleTime"] = v
		}
		if v, ok := d.GetOkExists("error_notice"); ok {
			configureSubscriptionReq["ErrorNotice"] = v
		}
		if v, ok := d.GetOk("error_phone"); ok {
			configureSubscriptionReq["ErrorPhone"] = v
		}
		if v, ok := d.GetOk("reserve"); ok {
			configureSubscriptionReq["Reserve"] = v
		}
		if d.HasChange("source_endpoint_ssl") {
			if err := setEndpointSSL(d, configureSubscriptionReq); err != nil {
				return errs.WrapError(err)
			}
		}
		action := "ConfigureSubscription"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, configureSubscriptionReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, configureSubscriptionReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		target := d.Get("status").(string)
		err = resourceSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, d.Get("status")))
		}
	}
	d.Partial(false)
	return resourceSubscriptionJobRead(d, meta)
}
func resourceSubscriptionJobDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v.(string) == "Subscription" {
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDtsJob"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}

	if v, ok := d.GetOk("dts_instance_id"); ok {
		request["DtsInstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("synchronization_direction"); ok {
		request["SynchronizationDirection"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
func resourceSubscriptionJobStatusFlow(d *schema.ResourceData, meta interface{}, target string) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := Service{client}
	var response map[string]interface{}
	var err error
	object, err := dtsService.DescribeSubscriptionJob(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	if object["Status"].(string) != target {
		if target == "NotConfigured" {
			request := map[string]interface{}{
				"DtsJobId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request["SynchronizationDirection"] = v
			}
			action := "ResetDtsJob"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"NotConfigured"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.SubscriptionJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
		if target == "Normal" || (target == "Abnormal" && object["Status"].(string) == "NotStarted") {
			request := map[string]interface{}{
				"DtsJobId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request["SynchronizationDirection"] = v
			}
			action := "StartDtsJob"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, dtsService.SubscriptionJobStateRefreshFunc(d.Id(), []string{"PrecheckFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
		if target == "Abnormal" {
			request := map[string]interface{}{
				"DtsJobId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request["SynchronizationDirection"] = v
			}
			action := "SuspendDtsJob"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Abnormal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, dtsService.SubscriptionJobStateRefreshFunc(d.Id(), []string{"PrecheckFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}

	return nil
}

// subscriptionJobResolveJobId recovers the currently valid DTS job ID from the DTS instance
// backing this resource. The instance ID remains stable while the job attached to it can be
// deleted and re-created, leaving the job ID stored in state stale.
func subscriptionJobResolveJobId(d *schema.ResourceData, dtsService *Service) (string, error) {
	instanceId := ""
	if v, ok := d.GetOk("dts_instance_id"); ok {
		instanceId = fmt.Sprint(v)
	}
	if instanceId == "" {
		return "", fmt.Errorf("dts_instance_id is empty, cannot resolve the valid job ID of alicloud_dts_subscription_job %s", d.Id())
	}
	jobs, err := dtsService.DescribeJobsByInstanceId(instanceId, "SUBSCRIBE")
	if err != nil {
		return "", errs.WrapError(err)
	}
	if len(jobs) == 0 {
		return "", fmt.Errorf("no DTS subscription job was found under instance %s", instanceId)
	}
	if len(jobs) == 1 {
		return fmt.Sprint(jobs[0]["DtsJobId"]), nil
	}
	jobIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		jobIds = append(jobIds, fmt.Sprint(job["DtsJobId"]))
	}
	return "", fmt.Errorf("found multiple DTS subscription jobs (%s) under instance %s, cannot determine the valid job ID automatically", strings.Join(jobIds, ","), instanceId)
}

func convertPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "POSTPAY":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
func convertPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
