package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDBInstanceCreate,
		Read:   resourceAliCloudDBInstanceRead,
		Update: resourceAliCloudDBInstanceUpdate,
		Delete: resourceAliCloudDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"engine_version": {
				Type: schema.TypeString,
				// Remove this limitation and refer to https://www.alibabacloud.com/help/doc-detail/26228.htm each time
				//ValidateFunc: validateAllowedStringValue([]string{"5.5", "5.6", "5.7", "2008r2", "2012", "9.4", "9.3", "10.0"}),
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
					if v, ok := d.GetOk("storage_auto_scale"); ok && v.(string) == "Enable" && old != "" && new != "" && old != new {
						return true
					}
					if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "Serverless" && old != "" && new != "" && old != new {
						return true
					}
					return false
				},
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(Postpaid), string(Prepaid), string(Serverless)}, false),
				Optional:     true,
				Default:      Postpaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"monitoring_period": {
				Type:         schema.TypeInt,
				ValidateFunc: IntInSlice([]int{5, 10, 60, 300}),
				Optional:     true,
				Computed:     true,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntBetween(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidAndRenewDiffSuppressFunc,
			},
			"force": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Yes", "No"}, false),
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"db_time_zone": {
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
				ValidateFunc: StringLenBetween(2, 256),
			},

			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"connection_string_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: StringLenBetween(8, 64),
				Computed: true,
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
			"db_instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"whitelist_network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"Classic", "VPC", "MIX"}, false),
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"modify_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"Cover", "Append", "Delete"}, false),
				DiffSuppressFunc: securityIpsDiffSuppressFunc,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_ids"},
				Deprecated:    "Attribute `security_group_id` has been deprecated from 1.69.0 and use `security_group_ids` instead.",
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_ip_mode": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{NormalMode, SafetyMode}, false),
				Optional:     true,
				Default:      NormalMode,
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
				Set:      parameterToHash,
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
							// if attribute contains Optional feature, need to add Default: "", otherwise when terraform plan is executed, unmodified items wil detect differences.
							Default: "",
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
							Default:  "",
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
				Default:  false,
			},
			"tags": tagsSchema(),
			"babelfish_config": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"babelfish_enabled": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"migration_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"single-db", "multi-db"}, false),
						},
						"master_username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"master_user_password": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"babelfish_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Default to Manual
			"auto_upgrade_minor_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Auto", "Manual"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "MySQL" && d.Get("engine").(string) != "PostgreSQL"
				},
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3", "general_essd"}, false),
			},
			"sql_collector_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
				Computed:     true,
			},
			"sql_collector_config_value": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{30, 180, 365, 1095, 1825}),
				Default:      30,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("sql_collector_status"); ok && strings.ToLower(v.(string)) == "enabled" {
						return false
					}
					return true
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// currently, only mysql serverless support setting ssl_action
					return d.Get("instance_charge_type").(string) == "Serverless" && d.Get("engine").(string) != "MySQL"
				},
			},
			"ssl_connection_string": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslActionDiffSuppressFunc,
			},
			"tde_status": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
				Optional:     true,
				Computed:     true,
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					engine := d.Get("engine").(string)
					encryptionKey := d.Get("encryption_key").(string)
					if engine != "PostgreSQL" && engine != "MySQL" && engine != "SQLServer" {
						return true
					}
					if engine == "PostgreSQL" {
						if encryptionKey == "ServiceKey" && old != "" {
							return true
						}
						if encryptionKey == "disabled" && old == "" {
							return true
						}
					}
					return false
				},
			},
			"tde_encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id_slave_a": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id_slave_b": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ca_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_cert": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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

			"upgrade_db_instance_kernel_version": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Attribute `upgrade_db_instance_kernel_version` has been deprecated from 1.198.0 and use `target_minor_version` instead.",
			},
			"upgrade_time": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"Immediate", "MaintainTime", "SpecifyTime"}, false),
				DiffSuppressFunc: kernelSmallVersionDiffSuppressFunc,
			},
			"switch_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kernelSmallVersionDiffSuppressFunc,
			},
			"target_minor_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_auto_scale": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
				Optional:     true,
			},
			"storage_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{0, 10, 20, 30, 40, 50}),
				DiffSuppressFunc: StorageAutoScaleDiffSuppressFunc,
				Optional:         true,
			},
			"storage_upper_bound": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntAtLeast(0),
				DiffSuppressFunc: StorageAutoScaleDiffSuppressFunc,
				Optional:         true,
			},
			"ha_config": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Auto", "Manual"}, false),
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
				ValidateFunc: StringInSlice([]string{"None", "Lastest", "All"}, false),
			},
			"fresh_white_list_readins": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"db_is_ignore_case": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tcp_connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"SHORT", "LONG"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "HighAvailability", "AlwaysOn", "Finance", "cluster", "serverless_basic", "serverless_standard", "serverless_ha"}, false),
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Immediate", "MaintainTime"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serverless_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_capacity": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"min_capacity": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"auto_pause": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"switch_force": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "Serverless" {
						return true
					}
					return false
				},
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Up", "Down", "TempUpgrade", "Serverless"}, false),
			},
			"pg_bouncer_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "PostgreSQL"
				},
			},
			"recovery_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "SQLServer"
				},
			},
			"bursting_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"optimized_writes": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"optimized", "none"}, false),
			},
			"template_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"cold_data_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceAliCloudDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var err error
	action := "CreateDBInstance"
	request, err := buildDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	d.SetId(response["DBInstanceId"].(string))

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDBInstanceUpdate(d, meta)
}

func resourceAliCloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "CONFIG_ENCRYPTING", "SSL_MODIFYING", "TDE_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("deletion_protection") && (d.Get("instance_charge_type") == string(Postpaid) || d.Get("instance_charge_type") == string(Serverless)) {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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

		stateConf = BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		d.SetPartial("storage_auto_scale")
		d.SetPartial("storage_threshold")
		d.SetPartial("storage_upper_bound")
		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	payType := PayType(d.Get("instance_charge_type").(string))
	if !d.IsNewResource() && d.HasChange("instance_charge_type") {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_charge_type")
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("security_group_ids") || d.HasChange("security_group_id") {
		groupIds := d.Get("security_group_id").(string)
		if d.HasChange("security_group_ids") {
			groupIds = strings.Join(expandStringList(d.Get("security_group_ids").(*schema.Set).List())[:], COMMA_SEPARATED)
		}
		err := rdsService.ModifySecurityGroupConfiguration(d.Id(), groupIds)
		if err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_ids")
		d.SetPartial("security_group_id")
	}

	if d.HasChange("template_id_list") {
		oldIDs, newIDs := d.GetChange("template_id_list")
		oldIDSet := convertInterfaceListToSet(oldIDs.([]interface{}))
		newIDSet := convertInterfaceListToSet(newIDs.([]interface{}))

		addedIDs := getDifference(newIDSet, oldIDSet)
		removedIDs := getDifference(oldIDSet, newIDSet)

		for _, id := range addedIDs {
			action := "AttachWhitelistTemplateToInstance"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InsName":    d.Id(),
				"TemplateId": id,
			}
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, request)

			stateConf := BuildStateConf(
				[]string{},
				[]string{"Running"},
				d.Timeout(schema.TimeoutUpdate),
				5*time.Second,
				rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}),
			)
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		for _, id := range removedIDs {
			action := "DetachWhitelistTemplateToInstance"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InsName":    d.Id(),
				"TemplateId": id,
			}
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, request)

			stateConf := BuildStateConf(
				[]string{},
				[]string{"Running"},
				d.Timeout(schema.TimeoutUpdate),
				5*time.Second,
				rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}),
			)
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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

	if d.HasChange("tde_status") {
		action := "ModifyDBInstanceTDE"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"TDEStatus":    d.Get("tde_status"),
			"SourceIp":     client.SourceIp,
		}

		if "MySQL" == d.Get("engine").(string) || string(PostgreSQL) == d.Get("engine").(string) {
			if v, ok := d.GetOk("role_arn"); ok && v.(string) != "" {
				request["RoleARN"] = v.(string)
			}
			if v, ok := d.GetOk("tde_encryption_key"); ok && v.(string) != "" {
				request["EncryptionKey"] = v.(string)
				if ro, ok := request["RoleARN"].(string); !ok || ro == "" {
					roleArn, err := findKmsRoleArn(client, v.(string))
					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}
					request["RoleARN"] = roleArn
				}
			}
		}

		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("tde_status")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if d.HasChanges("node_id") {
		action := "SwitchDBInstanceHA"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"NodeId":       d.Get("node_id"),
			"SourceIp":     client.SourceIp,
		}
		if v, ok := d.GetOk("force"); ok && v.(string) != "" {
			request["Force"] = v
		}
		if v, ok := d.GetOk("effective_time"); ok && v.(string) != "" {
			request["EffectiveTime"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		nodeId := d.Get("node_id").(string)
		stateConfNodeId := BuildStateConf([]string{}, []string{nodeId}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceNodeIdRefreshFunc(d.Id()))
		if _, err := stateConfNodeId.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_id")
		d.SetPartial("force")
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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

	// port default to 3306 and if setting port to 3306, there will have a change
	if d.HasChanges("port", "connection_string_prefix", "babelfish_port") {
		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			return err
		}
		connectionStringPrefix := strings.Split(instance["ConnectionString"].(string), ".")[0]

		connectRequest["CurrentConnectionString"] = instance["ConnectionString"]
		connectRequest["Port"] = instance["Port"]
		connectRequest["ConnectionStringPrefix"] = connectionStringPrefix
		if v, ok := d.GetOk("port"); ok && v != instance["Port"] {
			connectUpdate = true
			connectRequest["Port"] = v
		}
		if v, ok := d.GetOk("connection_string_prefix"); ok && v != connectionStringPrefix {
			connectUpdate = true
			connectRequest["ConnectionStringPrefix"] = v
		}
		if d.HasChange("babelfish_port") {
			connectUpdate = true
		}
		if v, ok := d.GetOk("babelfish_port"); ok {
			connectRequest["BabelfishPort"] = v
		}
		if connectUpdate {
			var response map[string]interface{}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Rds", "2014-08-15", connectAction, nil, connectRequest, true)
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
			d.SetPartial("babelfish_port")
			// wait instance status is running after modifying
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	if d.HasChanges("ssl_action", "ssl_connection_string") {
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

		if sslAction == "Update" && (d.Get("engine").(string) == "PostgreSQL" || d.Get("engine").(string) == "MySQL") {
			request["SSLEnabled"] = 1
		}

		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if d.Get("engine").(string) == "PostgreSQL" {
			if d.HasChange("ca_type") {
				if v, ok := d.GetOk("ca_type"); ok && v.(string) != "" {
					request["CAType"] = v.(string)
				}
			}
			if d.HasChange("server_cert") {
				if v, ok := d.GetOk("server_cert"); ok && v.(string) != "" {
					request["ServerCert"] = v.(string)
				}
			}
			if d.HasChange("server_key") {
				if v, ok := d.GetOk("server_key"); ok && v.(string) != "" {
					request["ServerKey"] = v.(string)
				}
			}
			if d.HasChange("client_ca_enabled") {
				if v, ok := d.GetOk("client_ca_enabled"); ok {
					request["ClientCAEnabled"] = v.(int)
				}
			}
			if d.HasChange("client_ca_cert") {
				if v, ok := d.GetOk("client_ca_cert"); ok && v.(string) != "" {
					request["ClientCACert"] = v.(string)
				}
			}
			if d.HasChange("client_crl_enabled") {
				if v, ok := d.GetOk("client_crl_enabled"); ok {
					request["ClientCrlEnabled"] = v.(int)
				}
			}
			if d.HasChange("client_cert_revocation_list") {
				if v, ok := d.GetOk("client_cert_revocation_list"); ok && v.(string) != "" {
					request["ClientCertRevocationList"] = v.(string)
				}
			}
			if d.HasChange("acl") {
				if v, ok := d.GetOk("acl"); ok && v.(string) != "" {
					request["ACL"] = v.(string)
				}
			}
			if d.HasChange("replication_acl") {
				if v, ok := d.GetOk("replication_acl"); ok && v.(string) != "" {
					request["ReplicationACL"] = v.(string)
				}
			}
		}

		if d.Get("engine").(string) == "MySQL" {
			if d.HasChange("ca_type") {
				if v, ok := d.GetOk("ca_type"); ok && v.(string) != "" {
					request["CAType"] = v.(string)
				}
			}
			if d.HasChange("server_cert") {
				if v, ok := d.GetOk("server_cert"); ok && v.(string) != "" {
					request["ServerCert"] = v.(string)
				}
			}
			if d.HasChange("server_key") {
				if v, ok := d.GetOk("server_key"); ok && v.(string) != "" {
					request["ServerKey"] = v.(string)
				}
			}
		}

		request["ConnectionString"] = instance["ConnectionString"]
		if d.HasChange("ssl_connection_string") {
			if v, ok := d.GetOk("ssl_connection_string"); ok && v.(string) != "" {
				request["ConnectionString"] = v.(string)
			}
		}

		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_action")
		d.SetPartial("ssl_connection_string")
		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("pg_hba_conf") {
		err := rdsService.ModifyPgHbaConfig(d, "pg_hba_conf")
		if err != nil {
			return WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliCloudDBInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
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
		// default disable connect from outside
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		"PayType":      d.Get("instance_charge_type"),
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("effective_time"); ok && v.(string) != "" {
		request["EffectiveTime"] = v
	}
	if d.HasChange("instance_type") {
		update = true
	}
	if v, ok := d.GetOk("direction"); ok && v.(string) != "" {
		request["Direction"] = v
	}
	request["DBInstanceClass"] = d.Get("instance_type")

	if d.HasChange("instance_storage") {
		update = true
	}
	request["DBInstanceStorage"] = d.Get("instance_storage")

	if d.HasChange("bursting_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("bursting_enabled"); ok {
		request["BurstingEnabled"] = v
	}

	if d.HasChange("serverless_config") {
		update = true
		if v, ok := d.GetOk("serverless_config"); ok {
			v := v.([]interface{})[0].(map[string]interface{})
			if string(MySQL) == d.Get("engine").(string) || string(PostgreSQL) == d.Get("engine") {
				serverlessConfig, err := json.Marshal(struct {
					MaxCapacity float64 `json:"MaxCapacity"`
					MinCapacity float64 `json:"MinCapacity"`
					AutoPause   bool    `json:"AutoPause"`
					SwitchForce bool    `json:"SwitchForce"`
				}{
					v["max_capacity"].(float64),
					v["min_capacity"].(float64),
					v["auto_pause"].(bool),
					v["switch_force"].(bool),
				})
				if err != nil {
					return WrapError(err)
				}
				request["ServerlessConfiguration"] = string(serverlessConfig)
				if category, ok := d.GetOk("category"); ok {
					request["Category"] = category
				}
				request["Direction"] = "Serverless"
			} else if string(SQLServer) == d.Get("engine") {
				serverlessConfig, err := json.Marshal(struct {
					MaxCapacity float64 `json:"MaxCapacity"`
					MinCapacity float64 `json:"MinCapacity"`
				}{
					v["max_capacity"].(float64),
					v["min_capacity"].(float64),
				})
				if err != nil {
					return WrapError(err)
				}
				request["ServerlessConfiguration"] = string(serverlessConfig)
				if category, ok := d.GetOk("category"); ok {
					request["Category"] = category
				}
				request["Direction"] = "Serverless"
			}
		}
	}
	if d.HasChange("optimized_writes") && d.Get("engine").(string) == "MySQL" {
		update = true
		if optimizedWrites, ok := d.GetOk("optimized_writes"); ok && optimizedWrites.(string) != "" {
			request["OptimizedWrites"] = optimizedWrites.(string)
		}
	}

	if d.HasChange("db_instance_storage_type") {
		update = true
	}
	request["DBInstanceStorageType"] = d.Get("db_instance_storage_type")

	if d.HasChange("cold_data_enabled") {
		if v, ok := d.GetOkExists("cold_data_enabled"); ok {
			request["ColdDataEnabled"] = v
			update = true
		}
	}

	if update {
		// wait instance status is running before modifying
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
			d.SetPartial("serverless_config")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
	if d.HasChanges("vswitch_id") {
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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

	handleConfigChange := func(configName string, configValue interface{}) error {
		action := "ModifyDBInstanceConfig"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"ConfigName":   configName,
			"ConfigValue":  configValue,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		return nil
	}
	if "PostgreSQL" == d.Get("engine").(string) {
		if d.HasChange("pg_bouncer_enabled") {
			if err := handleConfigChange("pgbouncer", d.Get("pg_bouncer_enabled")); err != nil {
				return err
			}
		}

		if d.HasChange("encryption_key") {
			if v, ok := d.GetOk("encryption_key"); ok {
				var configValue string
				if v.(string) == "disabled" {
					configValue = "disabled"
				} else {
					configValue = v.(string)
				}
				if err := handleConfigChange("encryptionKey", configValue); err != nil {
					return err
				}
			}
		}

	}

	if "SQLServer" == d.Get("engine").(string) {
		if d.HasChange("recovery_model") {
			if err := handleConfigChange("backup_recovery_model", d.Get("recovery_model")); err != nil {
				return err
			}
		}
	}

	if !d.IsNewResource() && (d.HasChange("target_minor_version") || d.HasChange("upgrade_db_instance_kernel_version")) {
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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
	return resourceAliCloudDBInstanceRead(d, meta)
}

func resourceAliCloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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
	d.Set("tags", rdsService.tagsToMap(tags))

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
	if "PostgreSQL" == instance["Engine"] {
		DBInstanceEncryptionKey, err := rdsService.DescribeDBInstanceEncryptionKey(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("encryption_key", "")
		if DBInstanceEncryptionKey["EncryptionKey"] != "" {
			d.Set("encryption_key", DBInstanceEncryptionKey["EncryptionKey"])
		}
	}

	describeDBInstanceHAConfigObject, err := rdsService.DescribeDBInstanceHAConfig(d.Id())
	hostInstanceInfos := describeDBInstanceHAConfigObject["HostInstanceInfos"].(map[string]interface{})["NodeInfo"].([]interface{})
	var nodeId string
	for _, val := range hostInstanceInfos {
		item := val.(map[string]interface{})
		nodeType := item["NodeType"].(string)
		if nodeType == "Master" {
			nodeId = item["NodeId"].(string)
			break // 停止遍历
		}
	}
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
	d.Set("node_id", nodeId)
	d.Set("storage_auto_scale", d.Get("storage_auto_scale"))
	d.Set("storage_threshold", d.Get("storage_threshold"))
	d.Set("storage_upper_bound", d.Get("storage_upper_bound"))

	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("monitoring_period", monitoringPeriod)

	d.Set("security_ips", ips)
	d.Set("db_instance_ip_array_name", d.Get("db_instance_ip_array_name"))
	d.Set("db_instance_ip_array_attribute", d.Get("db_instance_ip_array_attribute"))
	d.Set("db_instance_type", instance["DBInstanceType"])
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
	d.Set("status", instance["DBInstanceStatus"])
	d.Set("create_time", instance["CreationTime"])
	d.Set("bursting_enabled", instance["BurstingEnabled"])
	d.Set("pg_bouncer_enabled", instance["PGBouncerEnabled"])
	if v, ok := instance["OptimizedWritesInfo"]; ok {
		value := ConvertMySQLInstanceOptimizedWritesResponse(fmt.Sprint(v))
		if value != "" {
			d.Set("optimized_writes", value)
		}
	}
	if val, exists := instance["ColdDataEnabled"]; exists && val != nil {
		d.Set("cold_data_enabled", val)
	}

	// MySQL Serverless instance query PayType return SERVERLESS, need to be consistent with the participant.
	payType := instance["PayType"]
	if instance["PayType"] == "SERVERLESS" {
		payType = "Serverless"
	}

	serverlessConfig := make([]map[string]interface{}, 0)
	slc := instance["ServerlessConfig"].(map[string]interface{})
	if payType == "Serverless" && (string(MySQL) == instance["Engine"] || string(PostgreSQL) == instance["Engine"]) {
		slcMaps := map[string]interface{}{
			"max_capacity": slc["ScaleMax"],
			"min_capacity": slc["ScaleMin"],
			"auto_pause":   slc["AutoPause"],
			"switch_force": slc["SwitchForce"],
		}
		serverlessConfig = append(serverlessConfig, slcMaps)
		d.Set("serverless_config", serverlessConfig)
	} else if payType == "Serverless" && string(SQLServer) == instance["Engine"] {
		slcMaps := map[string]interface{}{
			"max_capacity": slc["ScaleMax"],
			"min_capacity": slc["ScaleMin"],
		}
		serverlessConfig = append(serverlessConfig, slcMaps)
		d.Set("serverless_config", serverlessConfig)
	}
	d.Set("instance_charge_type", payType)
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance["VSwitchId"])
	// some instance class without connection string
	if instance["ConnectionString"] != nil {
		d.Set("connection_string", instance["ConnectionString"])
		d.Set("connection_string_prefix", strings.Split(fmt.Sprint(instance["ConnectionString"]), ".")[0])
		connection, err := rdsService.DescribeDBConnection(d.Id() + ":" + strings.Split(fmt.Sprint(instance["ConnectionString"]), ".")[0])
		if err != nil {
			return WrapError(err)
		}
		if connection["BabelfishPort"] != nil {
			d.Set("babelfish_port", connection["BabelfishPort"])
		}
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
	recoveryModel := instance["Extra"].(map[string]interface{})["RecoveryModel"]
	d.Set("recovery_model", recoveryModel)
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
		//period, err := computePeriodByUnit(instance["CreationTime"], instance["ExpireTime"], d.Get("period").(int), "Month")
		//if err != nil {
		//	return WrapError(err)
		//}
		//d.Set("period", period)
	}

	groups, err := rdsService.DescribeSecurityGroupConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(groups) > 0 {
		d.Set("security_group_id", strings.Join(groups, COMMA_SEPARATED))
		d.Set("security_group_ids", groups)
	}

	sslAction, err := rdsService.DescribeDBInstanceSSL(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction["RequireUpdate"])
	d.Set("ssl_action", convertRdsInstanceSslActionResponse(sslAction["SSLEnabled"], d.Get("ssl_action")))
	d.Set("client_ca_enabled", d.Get("client_ca_enabled"))
	d.Set("client_crl_enabled", d.Get("client_crl_enabled"))
	d.Set("ca_type", sslAction["CAType"])
	d.Set("server_cert", sslAction["ServerCert"])
	d.Set("server_key", sslAction["ServerKey"])
	d.Set("client_ca_cert", sslAction["ClientCACert"])
	d.Set("client_cert_revocation_list", sslAction["ClientCertRevocationList"])
	d.Set("acl", sslAction["ACL"])
	d.Set("replication_acl", sslAction["ReplicationACL"])
	d.Set("ssl_connection_string", sslAction["ConnectionString"])

	//When the instance schema is docker on ECS, TDE encryption is not supported, so the query is not executed.
	if kindCode, ok := instance["kindCode"]; ok && kindCode != "3" {
		tdeInfo, err := rdsService.DescribeRdsTDEInfo(d.Id())
		if err != nil && !IsExpectedErrors(err, DBInstanceTDEErrors) {
			return WrapError(err)
		}
		d.Set("tde_status", tdeInfo["TDEStatus"])
	}
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

	response, err := rdsService.DescribeParameters(d.Id())
	dbParamGroupInfo := response["ParamGroupInfo"].(map[string]interface{})
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_param_group_id", dbParamGroupInfo["ParamGroupId"].(string))

	WhitelistTemplate, err := rdsService.DescribeInstanceLinkedWhitelistTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if templateIds, ok := WhitelistTemplate["TemplateIds"].([]int); ok {
		d.Set("template_id_list", templateIds)
	}

	if templates, ok := WhitelistTemplate["Templates"].([]interface{}); ok {
		var terraformTemplates []interface{}
		for _, item := range templates {
			if templateMap, ok := item.(map[string]interface{}); ok {
				terraformTemplates = append(terraformTemplates, templateMap)
			}
		}
		d.Set("templates", terraformTemplates)
	} else {
		d.Set("templates", []interface{}{})
	}
	return nil
}

func resourceAliCloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_db_instance. Terraform will remove this resource from the state file, however resources may remain.")
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
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
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

func buildDBCreateRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"EngineVersion":         Trim(d.Get("engine_version").(string)),
		"Engine":                Trim(d.Get("engine").(string)),
		"DBInstanceStorage":     d.Get("instance_storage"),
		"DBInstanceClass":       Trim(d.Get("instance_type").(string)),
		"DBInstanceNetType":     Intranet,
		"DBInstanceDescription": d.Get("instance_name"),
		"SourceIp":              client.SourceIp,
	}
	if v, ok := d.GetOk("db_instance_storage_type"); ok && v.(string) != "" {
		request["DBInstanceStorageType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("db_param_group_id"); ok && v.(string) != "" {
		request["DBParamGroupId"] = v
	}
	if v, ok := d.GetOk("target_minor_version"); ok && v.(string) != "" {
		request["TargetMinorVersion"] = v
	}
	if v, ok := d.GetOk("port"); ok && v.(string) != "" {
		request["Port"] = v
	}
	if v, ok := d.GetOkExists("bursting_enabled"); ok {
		request["BurstingEnabled"] = v
	}
	if v, ok := d.GetOk("optimized_writes"); ok && v.(string) != "" {
		request["OptimizedWrites"] = v
	}

	if request["Engine"] == "MySQL" || request["Engine"] == "PostgreSQL" || request["Engine"] == "SQLServer" {
		if v, ok := d.GetOk("role_arn"); ok && v.(string) != "" {
			request["RoleARN"] = v.(string)
		}
		if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
			request["EncryptionKey"] = v.(string)
		}
	}

	if request["Engine"] == "MySQL" {
		if v, ok := d.GetOkExists("cold_data_enabled"); ok {
			request["ColdDataEnabled"] = v
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
			//else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			//	zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			//	if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
			//		return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			//	}
			//} else if request.ZoneId != vsw.ZoneId {
			//	return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
			//}
		}
	}

	request["PayType"] = Trim(d.Get("instance_charge_type").(string))

	// if charge type is postpaid, the commodity code must set to bards
	//args.CommodityCode = rds.Bards
	// At present, API supports two charge options about 'Prepaid'.
	// 'Month': valid period ranges [1-9]; 'Year': valid period range [1-3]
	// This resource only supports to input Month period [1-9, 12, 24, 36] and the values need to be converted before using them.
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

	if v, ok := d.GetOk("zone_id_slave_a"); ok {
		request["ZoneIdSlave1"] = v
	}

	if v, ok := d.GetOk("zone_id_slave_b"); ok {
		request["ZoneIdSlave2"] = v
	}

	if v, ok := d.GetOk("db_time_zone"); ok {
		request["DBTimeZone"] = v
	}
	if v, ok := d.GetOkExists("db_is_ignore_case"); ok {
		request["DBIsIgnoreCase"] = v
	}

	if request["Engine"] == string(PostgreSQL) {
		if v, ok := d.GetOk("babelfish_config"); ok {
			v := v.(*schema.Set).List()[0].(map[string]interface{})
			babelfishConfig, err := json.Marshal(struct {
				BabelfishEnabled   string `json:"babelfishEnabled"`
				MigrationMode      string `json:"migrationMode"`
				MasterUsername     string `json:"masterUsername"`
				MasterUserPassword string `json:"masterUserPassword"`
			}{v["babelfish_enabled"].(string),
				v["migration_mode"].(string),
				v["master_username"].(string),
				v["master_user_password"].(string),
			})
			if err != nil {
				return nil, err
			}
			request["BabelfishConfig"] = string(babelfishConfig)
		}
	}

	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}

	if v, ok := d.GetOk("private_ip_address"); ok && v.(string) != "" {
		request["PrivateIpAddress"] = v
	}

	if request["PayType"] == string(Serverless) {
		if v, ok := d.GetOk("serverless_config"); ok {
			v := v.([]interface{})[0].(map[string]interface{})
			if string(MySQL) == request["Engine"] || string(PostgreSQL) == request["Engine"] {
				serverlessConfig, err := json.Marshal(struct {
					MaxCapacity float64 `json:"MaxCapacity"`
					MinCapacity float64 `json:"MinCapacity"`
					AutoPause   bool    `json:"AutoPause"`
					SwitchForce bool    `json:"SwitchForce"`
				}{
					v["max_capacity"].(float64),
					v["min_capacity"].(float64),
					v["auto_pause"].(bool),
					v["switch_force"].(bool),
				})
				if err != nil {
					return nil, WrapError(err)
				}
				request["ServerlessConfig"] = string(serverlessConfig)
			} else if string(SQLServer) == request["Engine"] {
				serverlessConfig, err := json.Marshal(struct {
					MaxCapacity float64 `json:"MaxCapacity"`
					MinCapacity float64 `json:"MinCapacity"`
				}{
					v["max_capacity"].(float64),
					v["min_capacity"].(float64),
				})
				if err != nil {
					return nil, WrapError(err)
				}
				request["ServerlessConfig"] = string(serverlessConfig)
			}
		}
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request["ClientToken"] = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}

func findKmsRoleArn(client *connectivity.AliyunClient, k string) (string, error) {
	action := "DescribeKey"
	var response map[string]interface{}
	var err error
	request := make(map[string]interface{})
	request["KeyId"] = k
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, true)
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
		return "", WrapErrorf(err, DataDefaultErrorMsg, k, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.KeyMetadata.Creator", response)
	if err != nil {
		return "", WrapErrorf(err, FailedGetAttributeMsg, action, "$.VersionIds.VersionId", response)
	}
	return strings.Join([]string{"acs:ram::", fmt.Sprint(resp), ":role/aliyunrdsinstanceencryptiondefaultrole"}, ""), nil
}

func convertRdsInstanceSslActionResponse(sslEnable, sslActionRead interface{}) string {
	if fmt.Sprint(sslActionRead) == "Update" {
		return "Update"
	}
	switch fmt.Sprint(sslEnable) {
	case "Yes", "on":
		return "Open"
	case "No", "off":
		return "Close"
	}
	return fmt.Sprint(sslActionRead)
}
func convertInterfaceListToSet(list []interface{}) map[int]struct{} {
	set := make(map[int]struct{})
	for _, v := range list {
		if val, ok := v.(int); ok {
			set[val] = struct{}{}
		}
	}
	return set
}
func getDifference(a, b map[int]struct{}) []int {
	var result []int
	for key := range a {
		if _, exists := b[key]; !exists {
			result = append(result, key)
		}
	}
	return result
}

func ConvertMySQLInstanceOptimizedWritesResponse(source string) string {
	if source == "nil" {
		return ""
	}

	switch source {
	case "{\"optimized_writes\":true,\"init_optimized_writes\":true}":
		return "optimized"
	case "{\"optimized_writes\":false,\"init_optimized_writes\":true}":
		return "none"
	case "{\"optimized_writes\":false,\"init_optimized_writes\":false}":
		return ""
	default:
		return source
	}
}
