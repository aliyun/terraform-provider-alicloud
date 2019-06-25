package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/hashcode"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(MySQL), string(SQLServer), string(PostgreSQL), string(PPAS)}),
				ForceNew:     true,
				Required:     true,
			},
			"engine_version": {
				Type: schema.TypeString,
				// Remove this limitation and refer to https://www.alibabacloud.com/help/doc-detail/26228.htm each time
				//ValidateFunc: validateAllowedStringValue([]string{"5.5", "5.6", "5.7", "2008r2", "2012", "9.4", "9.3", "10.0"}),
				ForceNew: true,
				Required: true,
			},
			"db_instance_class": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'db_instance_class' has been deprecated from provider version 1.5.0. New field 'instance_type' replaces it.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'db_instance_storage' has been deprecated from provider version 1.5.0. New field 'instance_storage' replaces it.",
			},

			"instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(Postpaid), string(Prepaid)}),
				Optional:     true,
				ForceNew:     true,
				Default:      Postpaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rdsPostPaidDiffSuppressFunc,
			},
			"monitoring_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue([]int{5, 60, 300}),
				Optional:     true,
				Computed:     true,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: rdsPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rdsPostPaidDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"multi_az": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'multi_az' has been deprecated from provider version 1.8.1. Please use field 'zone_id' to specify multiple availability zone.",
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBInstanceName,
			},

			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"db_instance_net_type": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'db_instance_net_type' has been deprecated from provider version 1.5.0.",
			},
			"allocate_public_connection": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'allocate_public_connection' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_connection' replaces it.",
			},

			"instance_network_type": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_network_type' has been deprecated from provider version 1.5.0.",
			},

			"master_user_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_user_name' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_account' field 'name' replaces it.",
			},

			"master_user_password": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_user_password' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_account' field 'password' replaces it.",
			},

			"preferred_backup_period": {
				Type:       schema.TypeList,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Optional:   true,
				Deprecated: "Field 'preferred_backup_period' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'backup_period' replaces it.",
			},

			"preferred_backup_time": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'preferred_backup_time' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'backup_time' replaces it.",
			},

			"backup_retention_period": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'backup_retention_period' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'retention_period' replaces it.",
			},

			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'connections' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_connection' replaces it.",
			},

			"db_mappings": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"character_set_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'db_mappings' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_database' replaces it.",
			},

			"parameters": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
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
		return err
	}

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.CreateDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*rds.CreateDBInstanceResponse)
	d.SetId(response.DBInstanceId)

	// wait instance status change from Creating to running
	if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBInstanceUpdate(d, meta)
}

func resourceAlicloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("auto_renew") || d.HasChange("auto_renew_period") {
		request := rds.CreateModifyInstanceAutoRenewalAttributeRequest()
		request.DBInstanceId = d.Id()

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
		addDebug(request.GetActionName(), raw)

		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("security_group_id") {
		err := rdsService.ModifySecurityGroupConfiguration(d.Id(), d.Get("security_group_id").(string))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_group_id")
	}

	if d.HasChange("monitoring_period") {
		period := d.Get("monitoring_period").(int)
		request := rds.CreateModifyDBInstanceMonitorRequest()
		request.DBInstanceId = d.Id()
		request.Period = strconv.Itoa(period)

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMonitor(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDBInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("instance_name").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
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
	request.DBInstanceId = d.Id()
	request.PayType = string(Postpaid)

	if d.HasChange("instance_type") {
		request.DBInstanceClass = d.Get("instance_type").(string)
		update = true
	}

	if d.HasChange("instance_storage") {
		request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
		update = true
	}

	if update {
		// wait instance status is running before modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{"InvalidOrderTask.NotSupport"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			d.SetPartial("instance_type")
			d.SetPartial("instance_storage")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, 2*DefaultLongTimeout); err != nil {
			return WrapError(err)
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
		if rdsService.NotFoundDBInstance(err) {
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

	d.Set("monitoring_period", monitoringPeriod)

	d.Set("security_ips", ips)

	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("instance_type", instance.DBInstanceClass)
	d.Set("port", instance.Port)
	storage := instance.DBInstanceStorage
	if instance.Engine == string(SQLServer) {
		// It is product design and gives more 5GB for SqlServer Instance
		storage = instance.DBInstanceStorage - 5
	}
	d.Set("instance_storage", storage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.PayType)
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	d.Set("instance_name", instance.DBInstanceDescription)

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}

	if instance.PayType == string(Prepaid) {
		request := rds.CreateDescribeInstanceAutoRenewalAttributeRequest()
		request.DBInstanceId = d.Id()

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*rds.DescribeInstanceAutoRenewalAttributeResponse)
		if response != nil && len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			auto_renew := bool(renew.AutoRenew == "True")

			d.Set("auto_renew", auto_renew)
			d.Set("auto_renew_period", renew.Duration)
		}
	}

	object, err := rdsService.DescribeSecurityGroupConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_group_id", object)

	return nil
}

func resourceAlicloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(instance.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(request)
		})

		if err != nil && !rdsService.NotFoundDBInstance(err) {
			if IsExceptedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(rdsService.WaitForDBInstance(d.Id(), Deleted, DefaultTimeoutMedium))
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

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request.InstanceNetworkType = string(Classic)

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VPCId = vsw.VpcId
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

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request.ClientToken = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}
