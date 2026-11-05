package dts

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMigrationJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceMigrationJobCreate,
		Read:   resourceMigrationJobRead,
		Update: resourceMigrationJobUpdate,
		Delete: resourceMigrationJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"data_synchronization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"db_list": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_check_configure": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"destination_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"RDS", "PolarDB", "POLARDBX20", "ADS", "MONGODB", "GREENPLUM", "DATAHUB", "OTHER", "ECS", "EXPRESS", "CEN", "DG"}, false),
			},

			"destination_endpoint_engine_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validators.StringInSlice([]string{"ADS", "ADB30", "AS400", "DATAHUB", "DB2", "GREENPLUM", "KAFKA", "MARIADB", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "PostgreSQL"}, true),
				DiffSuppressFunc: validators.UpperLowerCaseDiffSuppressFunc,
			},
			"destination_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dest_primary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dest_secondary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"CEN", "DG", "DISTRIBUTED_DMSLOGICDB", "ECS", "EXPRESS", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS"}, false),
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_engine_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validators.StringInSlice([]string{"AS400", "DB2", "DMSPOLARDB", "HBASE", "MARIADB", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "POSTGRESQL", "TERADATA"}, true),
				DiffSuppressFunc: validators.UpperLowerCaseDiffSuppressFunc,
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"src_primary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"src_secondary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_ssl": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"0", "1"}, false),
			},
			"destination_endpoint_ssl": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"0", "1", "3"}, false),
			},
			"structure_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"checkpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"xxlarge", "xlarge", "large", "medium", "small"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Suspending", "Migrating"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// A migration without incremental synchronization runs once and stays
					// Finished; the terminal state already satisfies a requested running
					// state and cannot be started or suspended again.
					return old == "Finished" && (new == "Migrating" || new == "Suspending")
				},
			},
		},
	}
}

func resourceMigrationJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ConfigureDtsJob"
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("dts_job_name"); ok {
		request["DtsJobName"] = v
	}
	if v, ok := d.GetOk("checkpoint"); ok {
		request["Checkpoint"] = v
	}

	request["DtsInstanceId"] = d.Get("dts_instance_id")
	request["DataInitialization"] = d.Get("data_initialization")
	request["DataSynchronization"] = d.Get("data_synchronization")
	request["StructureInitialization"] = d.Get("structure_initialization")
	request["DbList"] = d.Get("db_list")
	if err := setEndpointSSL(d, request); err != nil {
		return errs.WrapError(err)
	}

	request["DestinationEndpointInstanceType"] = d.Get("destination_endpoint_instance_type")
	request["DestinationEndpointEngineName"] = d.Get("destination_endpoint_engine_name")
	if v, ok := d.GetOk("destination_endpoint_instance_id"); ok {
		request["DestinationEndpointInstanceID"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request["DestinationEndpointRegion"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_ip"); ok {
		request["DestinationEndpointIP"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_port"); ok {
		request["DestinationEndpointPort"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_database_name"); ok {
		request["DestinationEndpointDataBaseName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
		request["DestinationEndpointUserName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_password"); ok {
		request["DestinationEndpointPassword"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_oracle_sid"); ok {
		request["DestinationEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("dest_primary_vswitch_id"); ok {
		request["DestPrimaryVswId"] = v
	}
	if v, ok := d.GetOk("dest_secondary_vswitch_id"); ok {
		request["DestSecondaryVswId"] = v
	}

	request["SourceEndpointInstanceType"] = d.Get("source_endpoint_instance_type")
	request["SourceEndpointEngineName"] = d.Get("source_endpoint_engine_name")
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		request["SourceEndpointInstanceID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceEndpointRegion"] = v
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		request["SourceEndpointIP"] = v
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		request["SourceEndpointPort"] = v
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		request["SourceEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		request["SourceEndpointDatabaseName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		request["SourceEndpointPassword"] = v
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		request["SourceEndpointUserName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		request["SourceEndpointOwnerID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		request["SourceEndpointRole"] = v
	}
	if v, ok := d.GetOk("src_primary_vswitch_id"); ok {
		request["SrcPrimaryVswId"] = v
	}
	if v, ok := d.GetOk("src_secondary_vswitch_id"); ok {
		request["SrcSecondaryVswId"] = v
	}

	request["JobType"] = "MIGRATION"
	if v, ok := d.GetOk("data_check_configure"); ok {
		request["DataCheckConfigure"] = v
	}
	request["RegionId"] = client.RegionId
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dts_migration_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DtsJobId"]))
	dtsService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"PreCheckPass", "Migrating", "PrecheckFailed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dtsService.MigrationJobStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceMigrationJobUpdate(d, meta)
}
func resourceMigrationJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := NewService(client)
	object, err := dtsService.DescribeMigrationJob(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_migration_job dtsService.DescribeMigrationJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("checkpoint", fmt.Sprint(object["Checkpoint"]))
	if v, ok := object["MigrationMode"].(map[string]interface{}); ok {
		d.Set("data_initialization", v["DataInitialization"])
		d.Set("data_synchronization", v["DataSynchronization"])
		d.Set("structure_initialization", v["StructureInitialization"])

	}
	d.Set("db_list", object["DbObject"])
	if v, ok := object["DestinationEndpoint"].(map[string]interface{}); ok {
		d.Set("destination_endpoint_database_name", v["DatabaseName"])
		d.Set("destination_endpoint_engine_name", v["EngineName"])
		d.Set("destination_endpoint_ip", v["Ip"])
		d.Set("destination_endpoint_instance_id", v["InstanceID"])
		d.Set("destination_endpoint_instance_type", v["InstanceType"])
		d.Set("destination_endpoint_oracle_sid", v["OracleSID"])
		d.Set("destination_endpoint_port", v["Port"])
		d.Set("destination_endpoint_region", v["Region"])
		d.Set("destination_endpoint_user_name", v["UserName"])
		if ssl := convertEndpointSslResponse(v["SslSolutionEnum"]); ssl != nil {
			d.Set("destination_endpoint_ssl", ssl)
		}
	}
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	if v, ok := object["SourceEndpoint"].(map[string]interface{}); ok {
		d.Set("source_endpoint_database_name", v["DatabaseName"])
		d.Set("source_endpoint_engine_name", v["EngineName"])
		d.Set("source_endpoint_ip", v["Ip"])
		d.Set("source_endpoint_instance_id", v["InstanceID"])
		d.Set("source_endpoint_instance_type", v["InstanceType"])
		d.Set("source_endpoint_oracle_sid", v["OracleSID"])
		d.Set("source_endpoint_owner_id", v["AliyunUid"])
		d.Set("source_endpoint_port", v["Port"])
		d.Set("source_endpoint_region", v["Region"])
		d.Set("source_endpoint_role", v["RoleName"])
		d.Set("source_endpoint_user_name", v["UserName"])
		if ssl := convertEndpointSslResponse(v["SslSolutionEnum"]); ssl != nil {
			d.Set("source_endpoint_ssl", ssl)
		}
	}
	d.Set("status", convertMigrationJobStatusResponse(object["Status"]))
	return nil
}
func resourceMigrationJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := NewService(client)
	var response map[string]interface{}

	if d.HasChange("status") {
		object, err := dtsService.DescribeMigrationJob(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Migrating" {
				request := map[string]interface{}{
					"DtsJobId": d.Id(),
				}
				request["RegionId"] = client.RegionId
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
				// StartDtsJob only submits the request; a job with data_check_configure
				// stays Prechecking for a while. Waiting here keeps the post-apply plan
				// empty instead of leaving the server-driven status mid-transition.
				// A job without incremental synchronization legitimately completes
				// without ever reporting Migrating, so Finished also satisfies the wait.
				stateConf := helper.BuildStateConf([]string{}, []string{"Migrating", "Finished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.MigrationJobStateRefreshFunc(d.Id(), []string{"PrecheckFailed", "InitializeFailed", "Failed", "MigrationFailed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
			if target == "Suspending" {
				request := map[string]interface{}{
					"DtsJobId": d.Id(),
				}
				request["RegionId"] = client.RegionId
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
				// SuspendDtsJob is asynchronous as well; without the wait the refreshed
				// status can still be Migrating when the post-apply plan is computed.
				stateConf := helper.BuildStateConf([]string{}, []string{"Suspending"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.MigrationJobStateRefreshFunc(d.Id(), []string{"PrecheckFailed", "InitializeFailed", "Failed", "MigrationFailed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}
	}

	if d.HasChange("source_endpoint_ssl") || d.HasChange("destination_endpoint_ssl") {
		reserved, err := endpointSSLReserved(d)
		if err != nil {
			return errs.WrapError(err)
		}
		modifyReservedReq := map[string]interface{}{
			"DtsInstanceId":  d.Get("dts_instance_id"),
			"RegionId":       client.RegionId,
			"ModifyTypeEnum": "UPDATE_RESERVED",
			"Reserved":       reserved,
		}

		action := "ModifyDtsJob"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, modifyReservedReq, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InvalidJobStatus", "InvalidTaskStatus", "DTS.Msg.OperationDenied.JobStatusModifying", "DTS.Msg.ModifyDenied.JobStatusNotRunning"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyReservedReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceMigrationJobRead(d, meta)
}
func resourceMigrationJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDtsJob"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}

	request["RegionId"] = client.RegionId
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
		if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
