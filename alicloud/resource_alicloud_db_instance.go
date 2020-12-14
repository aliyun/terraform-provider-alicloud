package alicloud

import (
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBInstanceCreate,
		Read:   resourceAlicloudDBInstanceRead,
		Update: resourceAlicloudDBInstanceUpdate,
		Delete: resourceAlicloudDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
				ForceNew: true,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
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
				Default:          1,
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
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
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
				ValidateFunc: validation.StringInSlice([]string{NormalMode, SafetyMode}, false),
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
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),

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
				Default:      30,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				},
			"zone_id_slave_a": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id_slave_b": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceAlicloudDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	action := "CreateDBInstance"
	request, err := buildDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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

	return resourceAlicloudDBInstanceUpdate(d, meta)
}

func resourceAlicloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "CONFIG_ENCRYPTING", "SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	payType := PayType(d.Get("instance_charge_type").(string))
	if !d.IsNewResource() && d.HasChange("instance_charge_type") && payType == Prepaid {
		action := "ModifyDBInstancePayType"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"PayType":      payType,
			"AutoPay":      "true",
			"UsedTime":     d.Get("period"),
			"Period":       Month,
			"SourceIp":     client.SourceIp,
		}
		period := d.Get("period").(int)
		if period > 9 {
			request["UsedTime"] = period / 12
			request["Period"] = Year
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		// wait instance status is Normal after modifying
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
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		d.SetPartial("security_group_ids")
		d.SetPartial("security_group_id")
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
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("auto_upgrade_minor_version")
	}

	if d.HasChange("security_ip_mode") && d.Get("security_ip_mode").(string) == SafetyMode {
		action := "MigrateSecurityIPMode"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("sql_collector_status")
	}

	if d.Get("sql_collector_status").(string) == "Enabled" && d.HasChange("sql_collector_config_value") {
		action := "ModifySQLCollectorRetention"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"ConfigValue":  strconv.Itoa(d.Get("sql_collector_config_value").(int)),
			"SourceIp":     client.SourceIp,
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		request["ConnectionString"] = instance["ConnectionString"]
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
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
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDBInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
		action := "ModifyDBInstanceDescription"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceId":          d.Id(),
			"DBInstanceDescription": d.Get("instance_name"),
			"SourceIp":              client.SourceIp,
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("instance_name")
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := rdsService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ips")
	}

	update := false
	action := "ModifyDBInstanceSpec"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"PayType":      d.Get("instance_charge_type"),
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("instance_type") {
		request["DBInstanceClass"] = d.Get("instance_type")
		update = true
	}

	if d.HasChange("instance_storage") {
		request["DBInstanceStorage"] = d.Get("instance_storage")
		update = true
	}
	if d.HasChange("db_instance_storage_type") {
		request["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
		update = true
	}
	if update {
		// wait instance status is running before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		runtime := util.RuntimeOptions{}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAlicloudDBInstanceRead(d, meta)
}

func resourceAlicloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	ips, err := rdsService.GetSecurityIps(d.Id())
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

	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("monitoring_period", monitoringPeriod)

	d.Set("security_ips", ips)
	d.Set("security_ip_mode", instance["SecurityIPMode"])

	d.Set("engine", instance["Engine"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("instance_type", instance["DBInstanceClass"])
	d.Set("port", instance["Port"])
	d.Set("instance_storage", instance["DBInstanceStorage"])
	d.Set("db_instance_storage_type", instance["DBInstanceStorageType"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("instance_charge_type", instance["PayType"])
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance["VSwitchId"])
	d.Set("connection_string", instance["ConnectionString"])
	d.Set("instance_name", instance["DBInstanceDescription"])
	d.Set("maintain_time", instance["MaintainTime"])
	d.Set("auto_upgrade_minor_version", instance["AutoUpgradeMinorVersion"])
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

	if instance["PayType"] == string(Prepaid) {
		action := "DescribeInstanceAutoRenewalAttribute"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"SourceIp":     client.SourceIp,
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		items := response["Items"].(map[string]interface{})["Item"].([]interface{})
		if response != nil && len(items) > 0 {
			renew := items[0].(map[string]interface{})
			d.Set("auto_renew", renew["AutoRenew"] == "True")
			d.Set("auto_renew_period", renew["Duration"])
		}
		period, err := computePeriodByUnit(instance["CreationTime"], instance["ExpireTime"], d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}

	groups, err := rdsService.DescribeSecurityGroupConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_group_id", strings.Join(groups, COMMA_SEPARATED))
	d.Set("security_group_ids", groups)

	sslAction, err := rdsService.DescribeDBInstanceSSL(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction["RequireUpdate"])
	d.Set("ssl_action", d.Get("ssl_action"))

	tdeInfo, err := rdsService.DescribeRdsTDEInfo(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo["TDEStatus"])

	return nil
}

func resourceAlicloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}
	action := "DeleteDBInstance"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
		"DBInstanceStorageType": d.Get("db_instance_storage_type"),
		"SourceIp":              client.SourceIp,
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
		request.EncryptionKey = v.(string)

		roleArn, err := findKmsRoleArn(client, v.(string))
		if err != nil {
			return nil, err
		}
		request.RoleARN = roleArn
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request["InstanceNetworkType"] = Classic

	if vswitchId != "" {
		request["VSwitchId"] = vswitchId
		request["InstanceNetworkType"] = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		v := strings.Split(vswitchId, COMMA_SEPARATED)[0]

		vsw, err := vpcService.DescribeVSwitch(v)
		if err != nil {
			return nil, WrapError(err)
		}

		if request["ZoneId"] == nil {
			request["ZoneId"] = vsw.ZoneId
		}

		if request["VPCId"] == nil {
			request["VPCId"] = vsw.VpcId
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

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request["ClientToken"] = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}

func findKmsRoleArn(client *connectivity.AliyunClient, k string) (string, error) {
	request := kms.CreateDescribeKeyRequest()
	request.KeyId = k

	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	if err != nil {
		return "", WrapErrorf(err, DataDefaultErrorMsg, k, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	key, _ := raw.(*kms.DescribeKeyResponse)
	return strings.Join([]string{"acs:ram::", key.KeyMetadata.Creator, ":role/aliyunrdsinstanceencryptiondefaultrole"}, ""), nil
}
