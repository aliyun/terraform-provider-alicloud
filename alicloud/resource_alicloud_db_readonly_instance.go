package alicloud

import (
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBReadonlyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBReadonlyInstanceCreate,
		Read:   resourceAlicloudDBReadonlyInstanceRead,
		Update: resourceAlicloudDBReadonlyInstanceUpdate,
		Delete: resourceAlicloudDBReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"master_db_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Computed:     true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
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

			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_enabled": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Optional:     true,
				Computed:     true,
			},
			"ca_type": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"aliyun", "custom"}, false),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"server_cert": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"server_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"client_ca_enabled": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{0, 1}),
				Optional:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"client_ca_cert": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"client_crl_enabled": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{0, 1}),
				Optional:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"client_cert_revocation_list": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"acl": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"cert", "perfer", "verify-ca", "verify-full"}, false),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
			},
			"replication_acl": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"cert", "perfer", "verify-ca", "verify-full"}, false),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslEnabledDiffSuppressFunc,
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:     true,
				Default:      Postpaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidAndRenewDiffSuppressFunc,
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Immediate", "MaintainTime"}, false),
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Up", "Down", "TempUpgrade", "Serverless"}, false),
			},
		},
	}
}

func resourceAlicloudDBReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	request, err := buildDBReadonlyCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	var response map[string]interface{}
	action := "CreateReadOnlyDBInstance"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	wait := incrementalWait(2*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.PrimaryDBInstanceStatus"}) {
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
	d.SetId(response["DBInstanceId"].(string))

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 15*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDBReadonlyInstanceUpdate(d, meta)
}

func resourceAlicloudDBReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	var err error
	sslUpdate := false
	sslAction := "ModifyDBInstanceSSL"
	sslRequest := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     client.RegionId,
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("ssl_enabled") {
		sslRequest["SSLEnabled"] = d.Get("ssl_enabled").(int)
		sslUpdate = true
	}
	if d.HasChange("ca_type") {
		sslRequest["CAType"] = d.Get("ca_type")
		sslUpdate = true
	}
	if d.HasChange("server_cert") {
		sslRequest["ServerCert"] = d.Get("server_cert")
		sslUpdate = true
	}
	if d.HasChange("server_key") {
		sslRequest["ServerKey"] = d.Get("server_key")
		sslUpdate = true
	}
	if d.HasChange("client_ca_enabled") {
		sslRequest["ClientCAEnabled"] = d.Get("client_ca_enabled")
		sslUpdate = true
	}
	if d.HasChange("client_ca_cert") {
		sslRequest["ClientCACert"] = d.Get("client_ca_cert")
		sslUpdate = true
	}
	if d.HasChange("client_crl_enabled") {
		sslRequest["ClientCrlEnabled"] = d.Get("client_crl_enabled")
		sslUpdate = true
	}
	if d.HasChange("client_cert_revocation_list") {
		sslRequest["ClientCertRevocationList"] = d.Get("client_cert_revocation_list")
		sslUpdate = true
	}
	if d.HasChange("acl") {
		sslRequest["ACL"] = d.Get("acl")
		sslUpdate = true
	}
	if d.HasChange("replication_acl") {
		sslRequest["ReplicationACL"] = d.Get("replication_acl")
		sslUpdate = true
	}
	if sslUpdate {
		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}

		sslRequest["ConnectionString"] = instance["ConnectionString"]
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", sslAction, nil, sslRequest, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), sslAction, AlibabaCloudSdkGoERROR)
		}
		addDebug(sslAction, response, sslRequest)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_enabled")
		d.SetPartial("ca_type")
		d.SetPartial("server_cert")
		d.SetPartial("server_key")
		d.SetPartial("client_ca_enabled")
		d.SetPartial("client_ca_cert")
		d.SetPartial("client_crl_enabled")
		d.SetPartial("client_cert_revocation_list")
		d.SetPartial("acl")
		d.SetPartial("replication_acl")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("deletion_protection") {
		err := rdsService.ModifyDBInstanceDeletionProtection(d, "deletion_protection")
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChanges("security_ips", "db_instance_ip_array_name", "db_instance_ip_array_attribute", "whitelist_network_type") {
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

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDBReadonlyInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
		action := "ModifyDBInstanceDescription"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceId":          d.Id(),
			"DBInstanceDescription": d.Get("instance_name"),
			"SourceIp":              client.SourceIp,
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			d.SetPartial("instance_name")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

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
	if d.HasChanges("instance_type", "instance_storage", "db_instance_storage_type") {
		if v, ok := d.GetOk("instance_type"); ok && v.(string) != "" {
			request["DBInstanceClass"] = v
		}
		if v, ok := d.GetOk("direction"); ok && v.(string) != "" {
			request["Direction"] = v
		}
		if v, ok := d.GetOk("instance_storage"); ok {
			request["DBInstanceStorage"] = v
		}
		if v, ok := d.GetOk("db_instance_storage_type"); ok && v.(string) != "" {
			request["DBInstanceStorageType"] = v
		}
		if v, ok := d.GetOk("effective_time"); ok && v.(string) != "" {
			request["EffectiveTime"] = v
		}
		update = true
	}

	if update {
		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		_, err := stateConf.WaitForState()
		if err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			d.SetPartial("instance_type")
			d.SetPartial("instance_storage")
			d.SetPartial("db_instance_storage_type")
			d.SetPartial("db_instance_storage_type")
			d.SetPartial("effective_time")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		_, err = stateConf.WaitForState()
		if err != nil {
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
			request["UsedTime"] = period
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}

	d.Partial(false)
	return resourceAlicloudDBReadonlyInstanceRead(d, meta)
}

func resourceAlicloudDBReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("engine", instance["Engine"])
	d.Set("master_db_instance_id", instance["MasterInstanceId"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("instance_type", instance["DBInstanceClass"])
	d.Set("port", instance["Port"])
	d.Set("instance_storage", instance["DBInstanceStorage"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("vswitch_id", instance["VSwitchId"])
	d.Set("connection_string", instance["ConnectionString"])
	d.Set("instance_name", instance["DBInstanceDescription"])
	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("target_minor_version", instance["CurrentKernelVersion"])
	d.Set("deletion_protection", instance["DeletionProtection"])
	d.Set("security_ips", ips)
	d.Set("db_instance_ip_array_name", d.Get("db_instance_ip_array_name"))
	d.Set("db_instance_ip_array_attribute", d.Get("db_instance_ip_array_attribute"))
	d.Set("security_ip_type", d.Get("security_ip_type"))
	d.Set("whitelist_network_type", d.Get("whitelist_network_type"))
	d.Set("instance_charge_type", instance["PayType"])
	d.Set("db_instance_storage_type", instance["DBInstanceStorageType"])

	sslAction, err := rdsService.DescribeDBInstanceSSL(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction["RequireUpdate"])
	d.Set("ssl_enabled", d.Get("ssl_enabled"))
	d.Set("client_ca_enabled", d.Get("client_ca_enabled"))
	d.Set("client_crl_enabled", d.Get("client_crl_enabled"))
	d.Set("ca_type", sslAction["CAType"])
	d.Set("server_cert", sslAction["ServerCert"])
	d.Set("server_key", sslAction["ServerKey"])
	d.Set("client_ca_cert", sslAction["ClientCACert"])
	d.Set("client_cert_revocation_list", sslAction["ClientCertRevocationList"])
	d.Set("acl", sslAction["ACL"])
	d.Set("replication_acl", sslAction["ReplicationACL"])

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
	}

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return err
	}

	tags, err := rdsService.describeTags(d)
	if err != nil {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", rdsService.tagsToMap(tags))
	}

	return nil
}

func resourceAlicloudDBReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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
		return WrapError(Error("'Prepaid' instance cannot be deleted can wait it to be expired and release it automatically. Or change instance_charge_type to 'Postpaid' to deleted"))
	}
	action := "DeleteDBInstance"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"RwSplitNetType.Exist", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Creating", "Active", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildDBReadonlyCreateRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"DBInstanceId":          Trim(d.Get("master_db_instance_id").(string)),
		"EngineVersion":         Trim(d.Get("engine_version").(string)),
		"DBInstanceStorage":     d.Get("instance_storage"),
		"DBInstanceClass":       Trim(d.Get("instance_type").(string)),
		"DBInstanceDescription": d.Get("instance_name"),
		"SourceIp":              client.SourceIp,
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}
	if auto_renew, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = auto_renew
	}
	vswitchId := Trim(d.Get("vswitch_id").(string))

	request["InstanceNetworkType"] = Classic

	if vswitchId != "" {
		request["VSwitchId"] = vswitchId
		request["InstanceNetworkType"] = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request["ZoneId"] == nil || request["ZoneId"].(string) == "" {
			request["ZoneId"] = vsw.ZoneId
		} else if strings.Contains(request["ZoneId"].(string), MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request["ZoneId"].(string), "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %v.", vsw.VSwitchId, request["ZoneId"]))
			}
		} else if request["ZoneId"] != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %v.", vsw.VSwitchId, request["ZoneId"]))
		}

		request["VPCId"] = vsw.VpcId
	}

	request["SecurityIPList"] = LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		request["SecurityIPList"] = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	// SQLServer only supports incoming cloud disk storage types.
	if dbInstanceStorageType, ok := d.GetOkExists("db_instance_storage_type"); ok {
		request["DBInstanceStorageType"] = dbInstanceStorageType
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

	request["ClientToken"] = buildClientToken("CreateReadOnlyDBInstance")

	return request, nil
}
