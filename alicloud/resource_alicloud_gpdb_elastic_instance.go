package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 7),
				Default:      7,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_recovery_point": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"encryption_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("encryption_type").(string) == "Off"
				},
			},
			"encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "CloudDisk"}, false),
				Default:      "Off",
			},
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
			"force_restart_instance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
				Default:      "VPC",
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2C16G", "4C32G", "16C128G"}, false),
			},
			"master_node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 16),
				Default:      1,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					return hashcode.String(
						v.(map[string]interface{})["name"].(string) + "|" + v.(map[string]interface{})["value"].(string))
				},
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Default:      "PayAsYouGo",
			},
			"payment_duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"payment_duration_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"preferred_backup_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"recovery_point_period": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("enable_recovery_point").(bool)
				},
				ValidateFunc: validation.StringInSlice([]string{"1", "2", "4", "8"}, false),
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"seg_node_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 4 || v > 256 || v%4 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 256 inclusive, and multiple of 4, got: %d", key, v))
					}
					return
				},
			},
			"seg_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency"}, false),
			},
			"ssl_enabled": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"ssl_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 50 || v > 4000 || v%50 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 50 and 4000 inclusive, and multiple of 50, got: %d", key, v))
					}
					return
				},
			},
			"tags": tagsSchema(),
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
		},
	}
}

func resourceAlicloudGpdbElasticInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateECSDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version").(string)
	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}
	request["InstanceSpec"] = d.Get("instance_spec")
	if v, ok := d.GetOk("master_node_num"); ok {
		request["MasterNodeNum"] = v.(int)
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertGpdbElasticInstancePaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("payment_duration"); ok {
		request["UsedTime"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId

	request["SegNodeNum"] = d.Get("seg_node_num").(int)
	request["SegStorageType"] = d.Get("seg_storage_type")
	request["StorageSize"] = d.Get("storage_size").(int)

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
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
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_elastic_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	gpdbService := GpdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1200*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
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
	d.Set("db_instance_description", object["DBInstanceDescription"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("encryption_type", convertGpdbElasticInstanceEncryptionTypeResponse(object["EncryptionType"]))
	d.Set("encryption_key", object["EncryptionKey"])
	d.Set("instance_network_type", object["InstanceNetworkType"])
	d.Set("instance_spec", convertDBInstanceClassToInstanceSpec(object["DBInstanceClass"].(string)))
	d.Set("payment_type", convertGpdbElasticInstancePaymentTypeResponse(object["PayType"]))
	d.Set("seg_node_num", formatInt(object["SegNodeNum"]))
	d.Set("master_node_num", formatInt(object["MasterNodeNum"]))
	d.Set("seg_storage_type", object["StorageType"])
	d.Set("status", object["DBInstanceStatus"])
	d.Set("storage_size", formatInt(object["StorageSize"]))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("tags", gpdbService.tagsRespToMap(object["Tags"].(map[string]interface{})["Tag"].([]interface{})))
	describeBackupPolicyObject, err := gpdbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_retention_period", formatInt(describeBackupPolicyObject["BackupRetentionPeriod"]))
	d.Set("enable_recovery_point", describeBackupPolicyObject["EnableRecoveryPoint"])
	d.Set("preferred_backup_period", describeBackupPolicyObject["PreferredBackupPeriod"])
	d.Set("preferred_backup_time", describeBackupPolicyObject["PreferredBackupTime"])
	d.Set("recovery_point_period", describeBackupPolicyObject["RecoveryPointPeriod"])

	describeDBInstanceSSLObject, err := gpdbService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ssl_expired_time", describeDBInstanceSSLObject["SSLExpiredTime"])

	securityIps, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", securityIps)
	if err = gpdbService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudGpdbElasticInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
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
	if d.HasChange("db_instance_description") {
		update = true
	}
	request["DBInstanceDescription"] = d.Get("db_instance_description")
	if update {
		action := "ModifyDBInstanceDescription"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
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
	modifyDBInstanceSSLReq["SSLEnabled"] = d.Get("ssl_enabled")
	if update {
		action := "ModifyDBInstanceSSL"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 180*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_enabled")
	}
	update = false
	modifyParametersReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("parameters") {
		update = true
	}
	parametersMap := make(map[string]interface{}, 0)
	for _, parameters := range d.Get("parameters").(*schema.Set).List() {
		parametersArg := parameters.(map[string]interface{})
		parametersMap[parametersArg["name"].(string)] = parametersArg["value"]
	}
	if v, err := convertArrayObjectToJsonString(parametersMap); err == nil {
		modifyParametersReq["Parameters"] = v
	} else {
		return WrapError(err)
	}
	if update {
		if v, ok := d.GetOkExists("force_restart_instance"); ok {
			modifyParametersReq["ForceRestartInstance"] = v.(bool)
		}
		action := "ModifyParameters"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
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
	if d.HasChange("security_ip_list") {
		update = true
	}
	ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
	modifySecurityIpsReq["SecurityIPList"] = strings.Join(ipList[:], COMMA_SEPARATED)
	if update {
		action := "ModifySecurityIps"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
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
		d.SetPartial("security_ip_list")
	}
	update = false
	modifyBackupPolicyReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		modifyBackupPolicyReq["PreferredBackupPeriod"] = d.Get("preferred_backup_period")
		update = true
	}
	if d.HasChange("preferred_backup_time") {
		modifyBackupPolicyReq["PreferredBackupTime"] = d.Get("preferred_backup_time")
		update = true
	}
	if d.HasChange("backup_retention_period") {
		update = true
	}
	modifyBackupPolicyReq["BackupRetentionPeriod"] = d.Get("backup_retention_period").(int)
	if d.HasChange("enable_recovery_point") {
		update = true
	}
	modifyBackupPolicyReq["EnableRecoveryPoint"] = d.Get("enable_recovery_point").(bool)
	if d.HasChange("recovery_point_period") {
		modifyBackupPolicyReq["RecoveryPointPeriod"] = d.Get("recovery_point_period")
		update = true
	}
	if update {
		action := "ModifyBackupPolicy"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
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

func convertDBInstanceClassToInstanceSpec(instanceClass string) string {
	splitClass := strings.Split(instanceClass, ".")
	return strings.ToUpper(splitClass[len(splitClass)-1])
}

func convertGpdbElasticInstanceEncryptionTypeResponse(source interface{}) interface{} {
	switch source {
	case "":
		return "Off"
	}
	return source
}
