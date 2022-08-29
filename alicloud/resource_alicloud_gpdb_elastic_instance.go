package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceAlicloudGpdbElasticInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbElasticInstanceCreate,
		Read:   resourceAlicloudGpdbElasticInstanceRead,
		Update: resourceAlicloudGpdbElasticInstanceUpdate,
		Delete: resourceAlicloudGpdbElasticInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"gpdb"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"6.0"}, false),
			},
			"seg_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency"}, false),
			},
			"seg_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 4 || v > 256 || v%4 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 256 inclusive, and multiple of 4, got: %d", key, v))
					}
					return
				},
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 50 || v > 4000 || v%50 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 50 and 4000 inclusive, and multiple of 50, got: %d", key, v))
					}
					return
				},
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2C16G", "4C32G", "16C128G", "2C8G", "4C16G", "8C32G", "16C64G"}, false),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "HighAvailability"}, false),
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
				Default:      "VPC",
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"payment_duration_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"payment_duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"force_restart_instance": {
				Type:         schema.TypeBool,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"ip_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_group_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"master_node_num": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preferred_backup_period": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Friday", "Monday", "Saturday", "Sunday", "Thursday", "Tuesday", "Wednesday"}, false),
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recovery_point_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8}),
			},
			"sql_collector_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enable"}, false),
			},
			"src_db_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 7),
			},
			"enable_recovery_point": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGpdbElasticInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	var response map[string]interface{}
	action := "CreateECSDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["SegStorageType"] = d.Get("seg_storage_type")
	request["SegNodeNum"] = d.Get("seg_node_num")
	request["StorageSize"] = d.Get("storage_size")
	request["InstanceSpec"] = d.Get("instance_spec")
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertGpdbElasticInstancePaymentTypeRequest(v)
	}
	if request["PayType"] == "Prepaid" {
		request["Period"] = d.Get("payment_duration_unit")
		paymentDuration := d.Get("payment_duration").(int)
		request["UsedTime"] = strconv.Itoa(paymentDuration)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("db_instance_category"); ok {
		request["DBInstanceCategory"] = v
	}
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}

	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}
	if v, ok := d.GetOk("master_node_num"); ok {
		request["MasterNodeNum"] = v
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	request["RegionId"] = client.RegionId

	request["SecurityIPList"] = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request["SecurityIPList"] = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}
	if v, ok := d.GetOk("src_db_instance_name"); ok {
		request["SrcDbInstanceName"] = v
	}
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["VPCId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	request["ClientToken"] = buildClientToken("CreateECSDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_elastic_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGpdbElasticInstanceUpdate(d, meta)
}
func resourceAlicloudGpdbElasticInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbElasticInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_elastic_instance gpdbService.DescribeGpdbElasticInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("seg_storage_type", object["StorageType"])
	if v, ok := object["SegNodeNum"]; ok && fmt.Sprint(v) != "0" {
		d.Set("seg_node_num", formatInt(v))
	}
	if v, ok := object["StorageSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("storage_size", formatInt(v))
	}
	d.Set("status", object["DBInstanceStatus"])
	d.Set("db_instance_description", object["DBInstanceDescription"])
	d.Set("instance_network_type", object["InstanceNetworkType"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("db_instance_category", object["DBInstanceCategory"])
	d.Set("encryption_key", object["EncryptionKey"])
	d.Set("encryption_type", object["EncryptionType"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("master_node_num", fmt.Sprint(formatInt(object["MasterNodeNum"])))
	d.Set("payment_type", convertGpdbElasticInstancePaymentTypeResponse(object["PayType"]))
	securityIps, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", securityIps)

	describeBackupPolicyObject, err := gpdbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := describeBackupPolicyObject["BackupRetentionPeriod"]; ok && fmt.Sprint(v) != "0" {
		d.Set("backup_retention_period", formatInt(v))
	}
	d.Set("enable_recovery_point", describeBackupPolicyObject["EnableRecoveryPoint"])
	d.Set("preferred_backup_period", describeBackupPolicyObject["PreferredBackupPeriod"])
	d.Set("preferred_backup_time", describeBackupPolicyObject["PreferredBackupTime"])
	if v, ok := describeBackupPolicyObject["RecoveryPointPeriod"]; ok && v.(string) != "" {
		d.Set("recovery_point_period", formatInt(v))
	}
	describeDBInstanceAttributeObject, err := gpdbService.DescribeDBInstanceAttributes(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("maintain_end_time", describeDBInstanceAttributeObject["MaintainEndTime"])
	d.Set("maintain_start_time", describeDBInstanceAttributeObject["MaintainStartTime"])
	describeDBInstanceSSLObject, err := gpdbService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ssl_enabled", describeDBInstanceSSLObject["SSLEnabled"])
	describeSQLCollectorPolicyObject, err := gpdbService.DescribeSQLCollectorPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("sql_collector_status", describeSQLCollectorPolicyObject["SQLCollectorStatus"])
	return nil
}
func resourceAlicloudGpdbElasticInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gpdbService.SetResourceTags(d, "ALIYUN::GPDB::INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		update = true
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("db_instance_description")
	}
	update = false
	modifyDBInstanceSSLReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("ssl_enabled") {
		update = true
	}
	if v, ok := d.GetOk("ssl_enabled"); ok {
		modifyDBInstanceSSLReq["SSLEnabled"] = v
	}
	if update {
		action := "ModifyDBInstanceSSL"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifyDBInstanceSSLReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceSSLReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_enabled")
	}
	//这个接口无法使用
	update = false
	modifySQLCollectorPolicyReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("sql_collector_status") {
		update = true
	}
	if v, ok := d.GetOk("sql_collector_status"); ok {
		modifySQLCollectorPolicyReq["SQLCollectorStatus"] = v
	}
	if update {
		action := "ModifySQLCollectorPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifySQLCollectorPolicyReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySQLCollectorPolicyReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("sql_collector_status")
	}
	update = false
	modifyDBInstanceMaintainTimeReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("maintain_end_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_end_time"); ok {
		modifyDBInstanceMaintainTimeReq["EndTime"] = v
	}
	if d.HasChange("maintain_start_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_start_time"); ok {
		modifyDBInstanceMaintainTimeReq["StartTime"] = v
	}
	if update {
		action := "ModifyDBInstanceMaintainTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifyDBInstanceMaintainTimeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceMaintainTimeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	update = false
	modifyParametersReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("parameters") {
		update = true
	}
	if v, ok := d.GetOk("parameters"); ok {
		parametersMap := map[string]interface{}{}
		for _, parameters := range v.(*schema.Set).List() {
			parametersArg := parameters.(map[string]interface{})
			parametersMap[parametersArg["parameter_name"].(string)] = parametersArg["current_value"]
		}
		if v, err := convertArrayObjectToJsonString(parametersMap); err == nil {
			modifyParametersReq["Parameters"] = v
		} else {
			return WrapError(err)
		}
	}
	if update {
		if v, ok := d.GetOkExists("force_restart_instance"); ok {
			modifyParametersReq["ForceRestartInstance"] = v
		}
		action := "ModifyParameters"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifyParametersReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyParametersReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("parameters")
	}
	update = false
	modifySecurityIpsReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("ip_whitelist") {
		update = true
		for _, ipWhitelist := range d.Get("ip_whitelist").(*schema.Set).List() {
			ipWhitelistArg := ipWhitelist.(map[string]interface{})
			modifySecurityIpsReq["DBInstanceIPArrayAttribute"] = ipWhitelistArg["ip_group_attribute"]
			modifySecurityIpsReq["DBInstanceIPArrayName"] = ipWhitelistArg["ip_group_name"]
			modifySecurityIpsReq["SecurityIpList"] = ipWhitelistArg["security_ip_list"]
		}
	}
	if update {
		action := "ModifySecurityIps"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifySecurityIpsReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySecurityIpsReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ip_whitelist")
		d.SetPartial("ip_group_attribute")
		d.SetPartial("ip_group_name")
	}
	update = false
	upgradeDBInstanceReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("instance_spec") {
		update = true
		upgradeDBInstanceReq["InstanceSpec"] = d.Get("instance_spec")
	}
	if !d.IsNewResource() && d.HasChange("master_node_num") {
		update = true
		upgradeDBInstanceReq["MasterNodeNum"] = d.Get("master_node_num")
	}
	if !d.IsNewResource() && d.HasChange("seg_node_num") {
		update = true
		upgradeDBInstanceReq["SegNodeNum"] = d.Get("seg_node_num")
	}
	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true
		upgradeDBInstanceReq["StorageSize"] = d.Get("storage_size")
	}
	upgradeDBInstanceReq["RegionId"] = client.RegionId
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, upgradeDBInstanceReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, upgradeDBInstanceReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_spec")
	}
	update = false
	modifyBackupPolicyReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		update = true
	}
	if v, ok := d.GetOk("preferred_backup_period"); ok {
		modifyBackupPolicyReq["PreferredBackupPeriod"] = v
	}
	if d.HasChange("preferred_backup_time") {
		update = true
	}
	if v, ok := d.GetOk("preferred_backup_time"); ok {
		modifyBackupPolicyReq["PreferredBackupTime"] = v
	}
	if d.HasChange("backup_retention_period") {
		update = true
		if v, ok := d.GetOk("backup_retention_period"); ok {
			modifyBackupPolicyReq["BackupRetentionPeriod"] = v
		}
	}
	if d.HasChange("enable_recovery_point") {
		update = true
		if v, ok := d.GetOkExists("enable_recovery_point"); ok {
			modifyBackupPolicyReq["EnableRecoveryPoint"] = v
		}
	}
	if d.HasChange("recovery_point_period") {
		update = true
		if v, ok := d.GetOk("recovery_point_period"); ok {
			modifyBackupPolicyReq["RecoveryPointPeriod"] = v
		}
	}
	if update {
		action := "ModifyBackupPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifyBackupPolicyReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyBackupPolicyReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("preferred_backup_period")
		d.SetPartial("preferred_backup_time")
		d.SetPartial("backup_retention_period")
		d.SetPartial("enable_recovery_point")
		d.SetPartial("recovery_point_period")
	}
	d.Partial(false)
	return resourceAlicloudGpdbElasticInstanceRead(d, meta)
}
func resourceAlicloudGpdbElasticInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertGpdbElasticInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertGpdbElasticInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}
