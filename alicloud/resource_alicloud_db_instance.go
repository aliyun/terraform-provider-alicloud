package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

	request, err := buildDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.CreateDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*rds.CreateDBInstanceResponse)
	d.SetId(response.DBInstanceId)

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
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "CONFIG_ENCRYPTING", "SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	payType := PayType(d.Get("instance_charge_type").(string))
	if !d.IsNewResource() && d.HasChange("instance_charge_type") && payType == Prepaid {
		prePaidRequest := rds.CreateModifyDBInstancePayTypeRequest()
		prePaidRequest.RegionId = client.RegionId
		prePaidRequest.DBInstanceId = d.Id()
		prePaidRequest.PayType = string(payType)
		prePaidRequest.AutoPay = "true"
		period := d.Get("period").(int)
		prePaidRequest.UsedTime = requests.Integer(strconv.Itoa(period))
		prePaidRequest.Period = string(Month)
		if period > 9 {
			prePaidRequest.UsedTime = requests.Integer(strconv.Itoa(period / 12))
			prePaidRequest.Period = string(Year)
		}
		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstancePayType(prePaidRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("period")

	}

	if payType == Prepaid && (d.HasChange("auto_renew") || d.HasChange("auto_renew_period")) {
		request := rds.CreateModifyInstanceAutoRenewalAttributeRequest()
		request.DBInstanceId = d.Id()
		request.RegionId = client.RegionId
		auto_renew := d.Get("auto_renew").(bool)
		if auto_renew {
			request.AutoRenew = "True"
		} else {
			request.AutoRenew = "False"
		}
		request.Duration = strconv.Itoa(d.Get("auto_renew_period").(int))

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

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
		request := rds.CreateModifyDBInstanceMonitorRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.Period = strconv.Itoa(period)

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMonitor(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("maintain_time") {
		request := rds.CreateModifyDBInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_time")
	}
	if d.HasChange("auto_upgrade_minor_version") {
		request := rds.CreateModifyDBInstanceAutoUpgradeMinorVersionRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.AutoUpgradeMinorVersion = d.Get("auto_upgrade_minor_version").(string)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceAutoUpgradeMinorVersion(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("auto_upgrade_minor_version")
	}

	if d.HasChange("security_ip_mode") && d.Get("security_ip_mode").(string) == SafetyMode {
		request := rds.CreateMigrateSecurityIPModeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.MigrateSecurityIPMode(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_ip_mode")
	}

	if d.HasChange("sql_collector_status") {
		request := rds.CreateModifySQLCollectorPolicyRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		if d.Get("sql_collector_status").(string) == "Enabled" {
			request.SQLCollectorStatus = "Enable"
		} else {
			request.SQLCollectorStatus = d.Get("sql_collector_status").(string)
		}
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifySQLCollectorPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("sql_collector_status")
	}

	if d.Get("sql_collector_status").(string) == "Enabled" && d.HasChange("sql_collector_config_value") {
		request := rds.CreateModifySQLCollectorRetentionRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.ConfigValue = strconv.Itoa(d.Get("sql_collector_config_value").(int))
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifySQLCollectorRetention(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("sql_collector_config_value")
	}

	if d.HasChange("ssl_action") {
		request := rds.CreateModifyDBInstanceSSLRequest()
		request.DBInstanceId = d.Id()
		request.RegionId = client.RegionId
		sslAction := d.Get("ssl_action").(string)
		if sslAction == "Close" {
			request.SSLEnabled = requests.NewInteger(0)
		}
		if sslAction == "Open" {
			request.SSLEnabled = requests.NewInteger(1)
		}
		if sslAction == "Update" {
			request.SSLEnabled = requests.NewInteger(2)
		}

		instance, err := rdsService.DescribeDBInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		request.ConnectionString = instance.ConnectionString

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBInstanceSSL(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("ssl_action")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tde_status") {
		request := rds.CreateModifyDBInstanceTDERequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.TDEStatus = d.Get("tde_status").(string)
		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("instance_name").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	request := rds.CreateModifyDBInstanceSpecRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()
	request.PayType = d.Get("instance_charge_type").(string)

	if d.HasChange("instance_type") {
		request.DBInstanceClass = d.Get("instance_type").(string)
		update = true
	}

	if d.HasChange("instance_storage") {
		request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
		update = true
	}
	if d.HasChange("db_instance_storage_type") {
		request.DBInstanceStorageType = d.Get("db_instance_storage_type").(string)
		update = true
	}
	if update {
		// wait instance status is running before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			d.SetPartial("instance_type")
			d.SetPartial("instance_storage")
			d.SetPartial("db_instance_storage_type")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
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

	d.Set("resource_group_id", instance.ResourceGroupId)
	d.Set("monitoring_period", monitoringPeriod)

	d.Set("security_ips", ips)
	d.Set("security_ip_mode", instance.SecurityIPMode)

	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("instance_type", instance.DBInstanceClass)
	d.Set("port", instance.Port)
	d.Set("instance_storage", instance.DBInstanceStorage)
	d.Set("db_instance_storage_type", instance.DBInstanceStorageType)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.PayType)
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	d.Set("instance_name", instance.DBInstanceDescription)
	d.Set("maintain_time", instance.MaintainTime)
	d.Set("auto_upgrade_minor_version", instance.AutoUpgradeMinorVersion)
	if len(instance.SlaveZones.SlaveZone) == 2 {
		d.Set("zone_id_slave_a", instance.SlaveZones.SlaveZone[0].ZoneId)
		d.Set("zone_id_slave_b", instance.SlaveZones.SlaveZone[1].ZoneId)
	} else if len(instance.SlaveZones.SlaveZone) == 1 {
		d.Set("zone_id_slave_a", instance.SlaveZones.SlaveZone[0].ZoneId)
	}
	if sqlCollectorPolicy.SQLCollectorStatus == "Enable" {
		d.Set("sql_collector_status", "Enabled")
	} else {
		d.Set("sql_collector_status", sqlCollectorPolicy.SQLCollectorStatus)
	}
	configValue, err := strconv.Atoi(sqlCollectorRetention.ConfigValue)
	if err != nil {
		return WrapError(err)
	}
	d.Set("sql_collector_config_value", configValue)

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}

	if instance.PayType == string(Prepaid) {
		request := rds.CreateDescribeInstanceAutoRenewalAttributeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*rds.DescribeInstanceAutoRenewalAttributeResponse)
		if response != nil && len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			d.Set("auto_renew", renew.AutoRenew == "True")
			d.Set("auto_renew_period", renew.Duration)
		}
		period, err := computePeriodByUnit(instance.CreationTime, instance.ExpireTime, d.Get("period").(int), "Month")
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
	d.Set("ssl_status", sslAction.RequireUpdate)
	d.Set("ssl_action", d.Get("ssl_action"))

	tdeInfo, err := rdsService.DescribeRdsTDEInfo(d.Id())
	if err != nil && !IsExpectedErrors(err, []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType"}) {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo.TDEStatus)

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
	if PayType(instance.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(request)
		})

		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Processing", "Pending", "NoStart", "Failed", "Default"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, rdsService.RdsTaskStateRefreshFunc(d.Id(), "DeleteDBInstance"))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildDBCreateRequest(d *schema.ResourceData, meta interface{}) (*rds.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := rds.CreateCreateDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("instance_type").(string))
	request.DBInstanceNetType = string(Intranet)
	request.DBInstanceDescription = d.Get("instance_name").(string)
	request.DBInstanceStorageType = d.Get("db_instance_storage_type").(string)

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request.InstanceNetworkType = string(Classic)

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		v := strings.Split(vswitchId, COMMA_SEPARATED)[0]

		vsw, err := vpcService.DescribeVSwitch(v)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		}

		if request.VPCId == "" {
			request.VPCId = vsw.VpcId
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

	request.PayType = Trim(d.Get("instance_charge_type").(string))

	// if charge type is postpaid, the commodity code must set to bards
	//args.CommodityCode = rds.Bards
	// At present, API supports two charge options about 'Prepaid'.
	// 'Month': valid period ranges [1-9]; 'Year': valid period range [1-3]
	// This resource only supports to input Month period [1-9, 12, 24, 36] and the values need to be converted before using them.
	if PayType(request.PayType) == Prepaid {

		period := d.Get("period").(int)
		request.UsedTime = strconv.Itoa(period)
		request.Period = string(Month)
		if period > 9 {
			request.UsedTime = strconv.Itoa(period / 12)
			request.Period = string(Year)
		}
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	if v, ok := d.GetOk("zone_id_slave_a"); ok {
		request.ZoneIdSlave1 = v.(string)
	}

	if v, ok := d.GetOk("zone_id_slave_b"); ok {
		request.ZoneIdSlave2 = v.(string)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request.ClientToken = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}
