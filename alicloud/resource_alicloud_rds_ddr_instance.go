package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDdrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDdrInstanceCreate,
		Read:   resourceAlicloudRdsDdrInstanceRead,
		Update: resourceAlicloudRdsDdrInstanceUpdate,
		Delete: resourceAlicloudRdsDdrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{string(MySQL), string(SQLServer)}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("storage_auto_scale"); ok && v.(string) == "Enable" {
						return true
					}
					return false
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"monitoring_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{5, 60, 300}),
				Optional:     true,
				Computed:     true,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 12),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidAndRenewDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If it is a new resource, do not suppress.
					if d.Id() == "" {
						return false
					}
					// If it is not a new resource and it is a multi-zone deployment, it needs to be suppressed.
					return len(strings.Split(new, ",")) > 1
				},
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"connection_string_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
				Computed:     true,
			},

			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"db_instance_ip_array_name": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"db_instance_ip_array_attribute": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"security_ip_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"whitelist_network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Classic", "VPC", "MIX"}, false),
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"modify_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Cover", "Append", "Delete"}, false),
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_ip_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{NormalMode, SafetyMode}, false),
				Optional:     true,
				Computed:     true,
			},

			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHashDdr,
				Optional: true,
				Computed: true,
			},
			"pg_hba_conf": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"option": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != string(PostgreSQL)
				},
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_upgrade_minor_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Auto", "Manual"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "MySQL"
				},
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"sql_collector_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
				Computed:     true,
			},
			"sql_collector_config_value": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{30, 180, 365, 1095, 1825}),
				Computed:     true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
			},
			"tde_status": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enabled"}, false),
				Optional:     true,
				ForceNew:     true,
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "MySQL" && d.Get("engine").(string) != "SQLServer"
				},
			},
			"ca_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_cert": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_ca_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_crl_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_cert_revocation_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"replication_acl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id_slave_a": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id_slave_b": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_db_instance_kernel_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"upgrade_time": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Immediate", "MaintainTime", "SpecifyTime"}, false),
				DiffSuppressFunc: kernelVersionDiffSuppressFunc,
			},
			"switch_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kernelVersionDiffSuppressFunc,
			},
			"target_minor_version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kernelVersionDiffSuppressFunc,
				Computed:         true,
			},
			"storage_auto_scale": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Optional:     true,
			},
			"storage_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{10, 20, 30, 40, 50}),
				DiffSuppressFunc: StorageAutoScaleDiffSuppressFunc,
				Optional:         true,
			},
			"storage_upper_bound": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntAtLeast(0),
				DiffSuppressFunc: StorageAutoScaleDiffSuppressFunc,
				Optional:         true,
			},
			"ha_config": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Auto", "Manual"}, false),
			},
			"manual_ha_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("ha_config"); ok && v.(string) == "Manual" {
						return false
					}
					return true
				},
			},
			"released_keep_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"None", "Lastest", "All"}, false),
			},
			"fresh_white_list_readins": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tcp_connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SHORT", "LONG"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Immediate", "MaintainTime"}, false),
			},
			"restore_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BackupSet", "BackupTime"}, false),
			},
			"backup_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"restore_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_db_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func parameterToHashDdr(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceAlicloudRdsDdrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var err error
	action := "CheckCreateDdrDBInstance"
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"Engine":            Trim(d.Get("engine").(string)),
		"EngineVersion":     Trim(d.Get("engine_version").(string)),
		"DBInstanceNetType": Intranet,
		"DBInstanceClass":   Trim(d.Get("instance_type").(string)),
		"DBInstanceStorage": d.Get("instance_storage"),
		"SourceIp":          client.SourceIp,
	}
	restoreType := d.Get("restore_type").(string)
	if restoreType == "BackupSet" {
		request["RestoreType"] = "0"
		if v, ok := d.GetOk("backup_set_id"); ok && v.(string) != "" {
			request["BackupSetId"] = v
		}
	}
	if restoreType == "BackupTime" {
		request["RestoreType"] = "1"
		if v, ok := d.GetOk("restore_time"); ok && v.(string) != "" {
			request["RestoreTime"] = v
		}
		if v, ok := d.GetOk("source_region"); ok && v.(string) != "" {
			request["SourceRegion"] = v
		}
		if v, ok := d.GetOk("source_db_instance_name"); ok && v.(string) != "" {
			request["SourceDBInstanceName"] = v
		}
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	action = "CreateDdrInstance"
	request, err = buildRdsDdrDBInstanceRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	d.SetId(response["DBInstanceId"].(string))
	d.SetPartial("restore_type")
	d.SetPartial("backup_set_id")
	d.SetPartial("restore_time")
	d.SetPartial("source_region")
	d.SetPartial("source_db_instance_name")

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRdsDdrInstanceUpdate(d, meta)
}

func resourceAlicloudRdsDdrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("pg_hba_conf") {
		err := rdsService.ModifyPgHbaConfig(d, "pg_hba_conf")
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("deletion_protection") && d.Get("payment_type") == "PayAsYouGo" {
		err := rdsService.ModifyDBInstanceDeletionProtection(d, "deletion_protection")
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("tcp_connection_type") {
		err := rdsService.ModifyHADiagnoseConfig(d, "tcp_connection_type")
		if err != nil {
			return WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}
	var err error

	if d.HasChanges("storage_auto_scale", "storage_threshold", "storage_upper_bound") {
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		action := "ModifyDasInstanceConfig"
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
			"RegionId":     client.RegionId,
			"SourceIp":     client.SourceIp,
		}

		if v, ok := d.GetOk("storage_auto_scale"); ok && v.(string) != "" {
			request["StorageAutoScale"] = v
		}
		if v, ok := d.GetOk("storage_threshold"); ok {
			request["StorageThreshold"] = v.(int)
		}
		if v, ok := d.GetOk("storage_upper_bound"); ok {
			request["StorageUpperBound"] = v.(int)
		}

		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		stateConf = BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		d.SetPartial("storage_auto_scale")
		d.SetPartial("storage_threshold")
		d.SetPartial("storage_upper_bound")
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	payType := PayType(convertRdsDdrPaymentTypeRequest(d.Get("payment_type")).(string))
	if !d.IsNewResource() && d.HasChange("payment_type") {
		action := "TransformDBInstancePayType"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"PayType":      payType,
			"SourceIp":     client.SourceIp,
		}
		if payType == Prepaid {
			period := d.Get("period").(int)
			request["UsedTime"] = d.Get("period")
			request["Period"] = Month
			if period > 9 {
				request["UsedTime"] = period / 12
				request["Period"] = Year
			}
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("payment_type")
		d.SetPartial("period")

	}

	if payType == Prepaid && (d.HasChange("auto_renew") || d.HasChange("auto_renew_period")) {
		action := "ModifyInstanceAutoRenewalAttribute"
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
			"RegionId":     client.RegionId,
			"SourceIp":     client.SourceIp,
		}
		auto_renew := d.Get("auto_renew").(bool)
		if auto_renew {
			request["AutoRenew"] = "True"
		} else {
			request["AutoRenew"] = "False"
		}
		request["Duration"] = strconv.Itoa(d.Get("auto_renew_period").(int))
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("security_group_ids") {
		groupIds := strings.Join(expandStringList(d.Get("security_group_ids").(*schema.Set).List())[:], COMMA_SEPARATED)
		err := rdsService.ModifySecurityGroupConfiguration(d.Id(), groupIds)
		if err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_ids")
	}

	if d.HasChange("monitoring_period") {
		period := d.Get("monitoring_period").(int)
		action := "ModifyDBInstanceMonitor"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"Period":       strconv.Itoa(period),
			"SourceIp":     client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("maintain_time") {
		action := "ModifyDBInstanceMaintainTime"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"MaintainTime": d.Get("maintain_time"),
			"ClientToken":  buildClientToken(action),
			"SourceIp":     client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("maintain_time")
	}
	if d.HasChange("auto_upgrade_minor_version") {
		action := "ModifyDBInstanceAutoUpgradeMinorVersion"
		request := map[string]interface{}{
			"RegionId":                client.SourceIp,
			"DBInstanceId":            d.Id(),
			"AutoUpgradeMinorVersion": d.Get("auto_upgrade_minor_version"),
			"ClientToken":             buildClientToken(action),
			"SourceIp":                client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_upgrade_minor_version")
	}

	if !d.IsNewResource() && d.HasChange("engine_version") && d.Get("engine").(string) == string(MySQL) {
		action := "UpgradeDBInstanceEngineVersion"
		request := map[string]interface{}{
			"RegionId":      client.SourceIp,
			"DBInstanceId":  d.Id(),
			"EngineVersion": d.Get("engine_version"),
			"EffectiveTime": d.Get("effective_time"),
			"ClientToken":   buildClientToken(action),
			"SourceIp":      client.SourceIp,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		d.SetPartial("engine_version")
		d.SetPartial("effective_time")
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("security_ip_mode") && d.Get("security_ip_mode").(string) == SafetyMode {
		action := "MigrateSecurityIPMode"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_ip_mode")
	}

	if d.HasChange("sql_collector_status") {
		action := "ModifySQLCollectorPolicy"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		if d.Get("sql_collector_status").(string) == "Enabled" {
			request["SQLCollectorStatus"] = "Enable"
		} else {
			request["SQLCollectorStatus"] = d.Get("sql_collector_status")
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("sql_collector_status")
	}

	if d.Get("sql_collector_status").(string) == "Enabled" && d.HasChange("sql_collector_config_value") && d.Get("engine").(string) == string(MySQL) {
		action := "ModifySQLCollectorRetention"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"ConfigValue":  strconv.Itoa(d.Get("sql_collector_config_value").(int)),
			"SourceIp":     client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("sql_collector_config_value")
	}

	if d.HasChange("ssl_action") {
		action := "ModifyDBInstanceSSL"
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
			"RegionId":     client.RegionId,
			"SourceIp":     client.SourceIp,
		}
		sslAction := d.Get("ssl_action").(string)
		if sslAction == "Close" {
			request["SSLEnabled"] = 0
		}
		if sslAction == "Open" {
			request["SSLEnabled"] = 1
		}
		if sslAction == "Update" {
			request["SSLEnabled"] = 2
		}

		if sslAction == "Update" && d.Get("engine").(string) == "PostgreSQL" {
			request["SSLEnabled"] = 1
		}

		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}

		if d.Get("engine").(string) == "PostgreSQL" {
			if v, ok := d.GetOk("ca_type"); ok && v.(string) != "" {
				request["CAType"] = v.(string)
			}
			if v, ok := d.GetOk("server_cert"); ok && v.(string) != "" {
				request["ServerCert"] = v.(string)
			}
			if v, ok := d.GetOk("server_key"); ok && v.(string) != "" {
				request["ServerKey"] = v.(string)
			}
			if v, ok := d.GetOk("client_ca_enabled"); ok {
				request["ClientCAEnabled"] = v.(int)
			}

			if v, ok := d.GetOk("client_ca_cert"); ok && v.(string) != "" {
				request["ClientCACert"] = v.(string)
			}

			if v, ok := d.GetOk("client_crl_enabled"); ok {
				request["ClientCrlEnabled"] = v.(int)
			}

			if v, ok := d.GetOk("client_cert_revocation_list"); ok && v.(string) != "" {
				request["ClientCertRevocationList"] = v.(string)
			}

			if v, ok := d.GetOk("acl"); ok && v.(string) != "" {
				request["ACL"] = v.(string)
			}

			if v, ok := d.GetOk("replication_acl"); ok && v.(string) != "" {
				request["ReplicationACL"] = v.(string)
			}
		}
		request["ConnectionString"] = instance["ConnectionString"]
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_action")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tde_status") {
		action := "ModifyDBInstanceTDE"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"TDEStatus":    d.Get("tde_status"),
			"SourceIp":     client.SourceIp,
		}

		if "MySQL" == d.Get("engine").(string) {
			if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
				request["EncryptionKey"] = v.(string)
				roleArn, err := rdsService.FindKmsRoleArnDdr(v.(string))
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				request["RoleARN"] = roleArn
			}
		}

		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("tde_status")

		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChanges("ha_config", "manual_ha_time") {
		action := "ModifyHASwitchConfig"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		if v, ok := d.GetOk("ha_config"); ok && v.(string) != "" {
			request["HAConfig"] = v
		}
		if v, ok := d.GetOk("manual_ha_time"); ok && v.(string) != "" {
			request["ManualHATime"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		d.SetPartial("ha_config")
		d.SetPartial("manual_ha_time")
		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	connectUpdate := false
	connectAction := "ModifyDBInstanceConnectionString"
	connectRequest := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     client.RegionId,
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("port") {
		connectUpdate = true
	}
	if d.HasChange("connection_string_prefix") {
		connectUpdate = true
	}
	if connectUpdate {
		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			return err
		}
		connectionStringPrefix := strings.Split(instance["ConnectionString"].(string), ".")[0]

		connectRequest["CurrentConnectionString"] = instance["ConnectionString"]
		connectRequest["Port"] = instance["Port"]
		connectRequest["ConnectionStringPrefix"] = connectionStringPrefix
		if v, ok := d.GetOk("port"); ok && v != instance["Port"] {
			connectRequest["Port"] = v
		}
		if v, ok := d.GetOk("connection_string_prefix"); ok && v != connectionStringPrefix {
			connectRequest["ConnectionStringPrefix"] = v
		}

		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", connectAction, nil, connectRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), connectAction, AlibabaCloudSdkGoERROR)
		}
		addDebug(connectAction, response, connectRequest)
		d.SetPartial("port")
		d.SetPartial("connection_string")
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudRdsDdrInstanceRead(d, meta)
	}

	if !d.IsNewResource() && d.HasChange("instance_name") {
		action := "ModifyDBInstanceDescription"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceId":          d.Id(),
			"DBInstanceDescription": d.Get("instance_name"),
			"SourceIp":              client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_name")
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}
		action := "ModifySecurityIps"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SecurityIps":  ipstr,
			"SourceIp":     client.SourceIp,
		}
		if v, ok := d.GetOk("db_instance_ip_array_name"); ok && v.(string) != "" {
			request["DBInstanceIPArrayName"] = v
		}
		if v, ok := d.GetOk("db_instance_ip_array_attribute"); ok && v.(string) != "" {
			request["DBInstanceIPArrayAttribute"] = v
		}
		if v, ok := d.GetOk("security_ip_type"); ok && v.(string) != "" {
			request["SecurityIPType"] = v
		}
		if v, ok := d.GetOk("whitelist_network_type"); ok && v.(string) != "" {
			request["WhitelistNetworkType"] = v
		}
		if v, ok := d.GetOk("modify_mode"); ok && v.(string) != "" {
			request["ModifyMode"] = v
		}
		if v, ok := d.GetOk("fresh_white_list_readins"); ok && v.(string) != "" {
			request["FreshWhiteListReadins"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_ips")
		d.SetPartial("db_instance_ip_array_name")
		d.SetPartial("db_instance_ip_array_attribute")
		d.SetPartial("security_ip_type")
		d.SetPartial("whitelist_network_type")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		action := "ModifyResourceGroup"
		request := map[string]interface{}{
			"DBInstanceId":    d.Id(),
			"ResourceGroupId": d.Get("resource_group_id"),
			"ClientToken":     buildClientToken(action),
			"SourceIp":        client.SourceIp,
		}
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update := false
	action := "ModifyDBInstanceSpec"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"PayType":      payType,
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("effective_time"); ok && v.(string) != "" {
		request["EffectiveTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("instance_type") {
		request["DBInstanceClass"] = d.Get("instance_type")
		update = true
	}

	if !d.IsNewResource() && d.HasChange("instance_storage") {
		request["DBInstanceStorage"] = d.Get("instance_storage")
		update = true
	}
	if !d.IsNewResource() && d.HasChange("db_instance_storage_type") {
		request["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
		update = true
	}
	if update {
		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			d.SetPartial("instance_type")
			d.SetPartial("instance_storage")
			d.SetPartial("db_instance_storage_type")
			d.SetPartial("effective_time")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		stateConf = BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	vpcService := VpcService{client}
	netUpdate := false
	netAction := "SwitchDBInstanceVpc"
	netRequest := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     client.RegionId,
		"SourceIp":     client.SourceIp,
	}
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		netUpdate = true
	}
	if d.HasChange("private_ip_address") {
		netUpdate = true
	}
	if netUpdate {
		v := d.Get("vswitch_id").(string)
		vsw, err := vpcService.DescribeVSwitch(v)
		if err != nil {
			return WrapError(err)
		}
		netRequest["VPCId"] = vsw.VpcId
		netRequest["VSwitchId"] = v
		if v, ok := d.GetOk("private_ip_address"); ok && v.(string) != "" {
			netRequest["PrivateIpAddress"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", netAction, nil, netRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), netAction, AlibabaCloudSdkGoERROR)
		}
		addDebug(netAction, response, netRequest)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("vswitch_id")
		d.SetPartial("private_ip_address")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("upgrade_db_instance_kernel_version") {
		action := "UpgradeDBInstanceKernelVersion"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		if v, ok := d.GetOk("upgrade_time"); ok && v.(string) != "" {
			request["UpgradeTime"] = v
		}
		if v, ok := d.GetOk("switch_time"); ok && v.(string) != "" {
			request["SwitchTime"] = v
		}
		if v, ok := d.GetOk("target_minor_version"); ok && v.(string) != "" {
			request["TargetMinorVersion"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("target_minor_version")
		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAlicloudRdsDdrInstanceRead(d, meta)
}

func resourceAlicloudRdsDdrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	dbInstanceIpArrayName := "default"
	if v, ok := d.GetOk("db_instance_ip_array_name"); ok {
		dbInstanceIpArrayName = v.(string)
	}

	ips, err := rdsService.GetSecurityIps(d.Id(), dbInstanceIpArrayName)
	if err != nil {
		return WrapError(err)
	}

	tags, err := rdsService.describeTags(d)
	if err != nil {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", rdsService.tagsToMap(tags))
	}

	monitoringPeriod, err := rdsService.DescribeDbInstanceMonitor(d.Id())
	if err != nil {
		return WrapError(err)
	}

	sqlCollectorPolicy, err := rdsService.DescribeSQLCollectorPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}

	sqlCollectorRetention, err := rdsService.DescribeSQLCollectorRetention(d.Id())
	if err != nil {
		return WrapError(err)
	}
	netInfoResponse, err := rdsService.DescribeDBInstanceNetInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var privateIpAddress string

	for _, item := range netInfoResponse {
		ipType := item.(map[string]interface{})["IPType"]
		if ipType == "Private" {
			privateIpAddress = item.(map[string]interface{})["IPAddress"].(string)
			break
		}
	}
	d.Set("private_ip_address", privateIpAddress)

	d.Set("storage_auto_scale", d.Get("storage_auto_scale"))
	d.Set("storage_threshold", d.Get("storage_threshold"))
	d.Set("storage_upper_bound", d.Get("storage_upper_bound"))

	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("monitoring_period", monitoringPeriod)

	d.Set("security_ips", ips)
	d.Set("db_instance_ip_array_name", d.Get("db_instance_ip_array_name"))
	d.Set("db_instance_ip_array_attribute", d.Get("db_instance_ip_array_attribute"))
	d.Set("security_ip_type", d.Get("security_ip_type"))
	d.Set("whitelist_network_type", d.Get("whitelist_network_type"))
	d.Set("security_ip_mode", instance["SecurityIPMode"])
	d.Set("engine", instance["Engine"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("instance_type", instance["DBInstanceClass"])
	d.Set("port", instance["Port"])
	d.Set("instance_storage", instance["DBInstanceStorage"])
	d.Set("db_instance_storage_type", instance["DBInstanceStorageType"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("payment_type", convertRdsDdrPaymentTypeResponse(instance["PayType"].(string)))
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance["VSwitchId"])
	// some instance class without connection string
	if instance["ConnectionString"] != nil {
		d.Set("connection_string", instance["ConnectionString"])
		d.Set("connection_string_prefix", strings.Split(fmt.Sprint(instance["ConnectionString"]), ".")[0])
	}
	d.Set("instance_name", instance["DBInstanceDescription"])
	d.Set("maintain_time", instance["MaintainTime"])
	d.Set("auto_upgrade_minor_version", instance["AutoUpgradeMinorVersion"])
	d.Set("target_minor_version", instance["CurrentKernelVersion"])
	d.Set("deletion_protection", instance["DeletionProtection"])
	d.Set("vpc_id", instance["VpcId"])
	d.Set("category", instance["Category"])
	slaveZones := instance["SlaveZones"].(map[string]interface{})["SlaveZone"].([]interface{})
	if len(slaveZones) == 2 {
		d.Set("zone_id_slave_a", slaveZones[0].(map[string]interface{})["ZoneId"])
		d.Set("zone_id_slave_b", slaveZones[1].(map[string]interface{})["ZoneId"])
	} else if len(slaveZones) == 1 {
		d.Set("zone_id_slave_a", slaveZones[0].(map[string]interface{})["ZoneId"])
	}
	if sqlCollectorPolicy["SQLCollectorStatus"] == "Enable" {
		d.Set("sql_collector_status", "Enabled")
	} else {
		d.Set("sql_collector_status", sqlCollectorPolicy["SQLCollectorStatus"])
	}
	configValue, err := strconv.Atoi(sqlCollectorRetention["ConfigValue"].(string))
	if err != nil {
		return WrapError(err)
	}
	d.Set("sql_collector_config_value", configValue)

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}
	if instance["Engine"].(string) == string(PostgreSQL) && instance["DBInstanceStorageType"].(string) != "local_ssd" {
		if err = rdsService.RefreshPgHbaConf(d, "pg_hba_conf"); err != nil {
			return WrapError(err)
		}
	}
	if err = rdsService.SetTimeZone(d); err != nil {
		return WrapError(err)
	}
	if instance["PayType"] == string(Prepaid) {
		action := "DescribeInstanceAutoRenewalAttribute"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		var err error
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		items := response["Items"].(map[string]interface{})["Item"].([]interface{})
		if response != nil && len(items) > 0 {
			renew := items[0].(map[string]interface{})
			d.Set("auto_renew", renew["AutoRenew"] == "True")
			d.Set("auto_renew_period", renew["Duration"])
		}
	}

	groups, err := rdsService.DescribeSecurityGroupConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(groups) > 0 {
		d.Set("security_group_ids", groups)
	}

	sslAction, err := rdsService.DescribeDBInstanceSSL(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction["RequireUpdate"])
	d.Set("ssl_action", d.Get("ssl_action"))
	d.Set("client_ca_enabled", d.Get("client_ca_enabled"))
	d.Set("client_crl_enabled", d.Get("client_crl_enabled"))
	d.Set("ca_type", sslAction["CAType"])
	d.Set("server_cert", sslAction["ServerCert"])
	d.Set("server_key", sslAction["ServerKey"])
	d.Set("client_ca_cert", sslAction["ClientCACert"])
	d.Set("client_cert_revocation_list", sslAction["ClientCertRevocationList"])
	d.Set("acl", sslAction["ACL"])
	d.Set("replication_acl", sslAction["ReplicationACL"])
	tdeInfo, err := rdsService.DescribeRdsTDEInfo(d.Id())
	if err != nil && !IsExpectedErrors(err, DBInstanceTDEErrors) {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo["TDEStatus"])

	res, err := rdsService.DescribeHASwitchConfig(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ha_config", res["HAConfig"])
	d.Set("manual_ha_time", res["ManualHATime"])

	res, err = rdsService.DescribeHADiagnoseConfig(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("tcp_connection_type", res["TcpConnectionType"])
	d.Set("restore_type", d.Get("restore_type"))
	d.Set("backup_set_id", d.Get("backup_set_id"))
	d.Set("restore_time", d.Get("restore_time"))
	d.Set("source_region", d.Get("source_region"))
	d.Set("source_db_instance_name", d.Get("source_db_instance_name"))
	return nil
}

func resourceAlicloudRdsDdrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(instance["PayType"].(string)) == Prepaid {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_rds_ddr_instance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	action := "DeleteDBInstance"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("released_keep_policy"); ok && v.(string) != "" {
		request["ReleasedKeepPolicy"] = v
	}
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Processing", "Pending", "NoStart", "Failed", "Default"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, rdsService.RdsTaskStateRefreshFunc(d.Id(), "DeleteDBInstance"))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildRdsDdrDBInstanceRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	vpcService := VpcService{client}
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"EngineVersion":     Trim(d.Get("engine_version").(string)),
		"Engine":            Trim(d.Get("engine").(string)),
		"DBInstanceStorage": d.Get("instance_storage"),
		"DBInstanceClass":   Trim(d.Get("instance_type").(string)),
		"DBInstanceNetType": Intranet,
		"SourceIp":          client.SourceIp,
		"RestoreType":       d.Get("restore_type").(string),
	}
	restoreType := d.Get("restore_type").(string)
	if restoreType == "BackupSet" {
		if v, ok := d.GetOk("backup_set_id"); ok && v.(string) != "" {
			request["BackupSetId"] = v
		}
	}
	if restoreType == "BackupTime" {
		if v, ok := d.GetOk("restore_time"); ok && v.(string) != "" {
			request["RestoreTime"] = v
		}
		if v, ok := d.GetOk("source_region"); ok && v.(string) != "" {
			request["SourceRegion"] = v
		}
		if v, ok := d.GetOk("source_db_instance_name"); ok && v.(string) != "" {
			request["SourceDBInstanceName"] = v
		}
	}
	if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
		request["DBInstanceDescription"] = v
	}
	if v, ok := d.GetOk("db_instance_storage_type"); ok && v.(string) != "" {
		request["DBInstanceStorageType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v
	}

	if request["Engine"] == "PostgreSQL" || request["Engine"] == "MySQL" || request["Engine"] == "SQLServer" {
		if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
			request["EncryptionKey"] = v.(string)

			roleArn, err := rdsService.FindKmsRoleArnDdr(v.(string))
			if err != nil {
				return nil, WrapError(err)
			}
			request["RoleARN"] = roleArn
		}
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	request["InstanceNetworkType"] = Classic
	if request["VSwitchId"] != nil {
		request["InstanceNetworkType"] = strings.ToUpper(string(Vpc))
		// check vswitchId in zone
		v := strings.Split(request["VSwitchId"].(string), COMMA_SEPARATED)[0]
		if request["ZoneId"] == nil || request["VPCId"] == nil {

			vsw, err := vpcService.DescribeVSwitch(v)
			if err != nil {
				return nil, WrapError(err)
			}

			if v, ok := request["VPCId"].(string); !ok || v == "" {
				request["VPCId"] = vsw.VpcId
			}
			if v, ok := request["ZoneId"].(string); !ok || v == "" {
				request["ZoneId"] = vsw.ZoneId
			}
		}
	}

	request["PayType"] = convertRdsDdrPaymentTypeRequest(d.Get("payment_type"))
	if PayType(request["PayType"].(string)) == Prepaid {
		period := d.Get("period").(int)
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = Month
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = Year
		}
	}

	request["SecurityIPList"] = LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		request["SecurityIPList"] = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request["ClientToken"] = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}

func convertRdsDdrPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}

func convertRdsDdrPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}
