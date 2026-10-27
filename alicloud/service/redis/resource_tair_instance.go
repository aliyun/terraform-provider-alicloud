package redis

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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

func ResourceTairInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTairInstanceCreate,
		Read:   resourceTairInstanceRead,
		Update: resourceTairInstanceUpdate,
		Delete: resourceTairInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_period": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"backup_time"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"backup_period"},
			},
			"cluster_backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateConfig,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Immediately", "MaintainTime"}, false),
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"global_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"tair_rdb", "tair_scm", "tair_essd"}, false),
			},
			"intranet_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE", "double", "single"}, false),
			},
			"param_no_loose_sentinel_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_no_loose_sentinel_password_free_access": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_no_loose_sentinel_password_free_commands": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_repl_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_semisync_repl_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_sentinel_compat_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60}),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntBetween(1024, 65535),
			},
			"read_only_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 5),
			},
			"recover_config_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new != "" && old != new {
						oldParts := strings.Split(old, ",")
						sort.Strings(oldParts)
						newParts := strings.Split(new, ",")
						sort.Strings(newParts)
						return reflect.DeepEqual(newParts, oldParts)
					}
					return false
				},
			},
			"security_ip_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new != "" && old != new {
						oldParts := strings.Split(old, ",")
						sort.Strings(oldParts)
						newParts := strings.Split(new, ",")
						sort.Strings(newParts)
						return reflect.DeepEqual(newParts, oldParts)
					}
					return false
				},
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"slave_read_only_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"src_db_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Disable", "Enable", "Update"}, false),
			},
			"ca_cert_download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert_common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"storage_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tde_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"tags": tags.TagsSchema(),
			"tair_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tair_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_auth_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTairInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTairInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["InstanceClass"] = d.Get("instance_class")
	request["VpcId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}
	if v, ok := d.GetOkExists("port"); ok && v.(int) > 0 {
		request["Port"] = v
	}
	if v, ok := d.GetOk("tair_instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertTairInstanceChargeTypeRequest(v.(string))
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	request["AutoPay"] = "true"
	if v, ok := d.GetOkExists("shard_count"); ok {
		request["ShardCount"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOkExists("storage_size_gb"); ok && v.(int) > 0 {
		request["Storage"] = v
	}
	if v, ok := d.GetOk("storage_performance_level"); ok {
		request["StorageType"] = convertTairInstanceStorageTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("cluster_backup_id"); ok {
		request["ClusterBackupId"] = v
	}
	if v, ok := d.GetOkExists("slave_read_only_count"); ok {
		request["SlaveReadOnlyCount"] = v
	}
	if v, ok := d.GetOkExists("read_only_count"); ok {
		request["ReadOnlyCount"] = v
	}
	if v, ok := d.GetOk("node_type"); ok {
		request["ShardType"] = convertTairInstanceShardTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("global_instance_id"); ok {
		request["GlobalInstanceId"] = v
	}
	if v, ok := d.GetOk("src_db_instance_id"); ok {
		request["SrcDBInstanceId"] = v
	}
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("recover_config_mode"); ok {
		request["RecoverConfigMode"] = v
	}
	if v, ok := d.GetOk("connection_string_prefix"); ok {
		request["ConnectionStringPrefix"] = v
	}
	wait := helper.IncrementalWait(180*time.Second, 180*time.Second)
	internalErrRetryCount := 0
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
		if err != nil {
			if (helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"SocketTimeout.Exception", "InternalError"})) && internalErrRetryCount < 2 {
				wait()
				internalErrRetryCount++
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_redis_tair_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	redisService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceTairInstanceUpdate(d, meta)
}

func resourceTairInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisService := NewService(client)

	objectRaw, err := redisService.DescribeTairInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_tair_instance DescribeTairInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	// IsSupportTDE reflects whether the instance actually supports Transparent Data Encryption.
	// It must be captured from the DescribeInstanceAttribute response before objectRaw gets
	// reassigned by the subsequent describe calls below.
	isSupportTDE := false
	if v, ok := objectRaw["IsSupportTDE"].(bool); ok {
		isSupportTDE = v
	}

	if objectRaw["ArchitectureType"] != nil {
		d.Set("architecture_type", objectRaw["ArchitectureType"])
	}
	if objectRaw["ConnectionDomain"] != nil {
		d.Set("connection_domain", objectRaw["ConnectionDomain"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["EngineVersion"] != nil {
		d.Set("engine_version", objectRaw["EngineVersion"])
	}
	if objectRaw["RealInstanceClass"] != nil {
		d.Set("instance_class", objectRaw["RealInstanceClass"])
	}
	if objectRaw["InstanceType"] != nil {
		d.Set("instance_type", objectRaw["InstanceType"])
	}
	if objectRaw["MaintainEndTime"] != nil {
		d.Set("maintain_end_time", objectRaw["MaintainEndTime"])
	}
	if objectRaw["MaintainStartTime"] != nil {
		d.Set("maintain_start_time", objectRaw["MaintainStartTime"])
	}
	if objectRaw["Connections"] != nil {
		d.Set("max_connections", objectRaw["Connections"])
	}
	if objectRaw["NetworkType"] != nil {
		d.Set("network_type", objectRaw["NetworkType"])
	}
	if objectRaw["NodeType"] != nil {
		d.Set("node_type", convertTairInstanceInstancesDBInstanceAttributeNodeTypeResponse(objectRaw["NodeType"]))
	}
	if objectRaw["ChargeType"] != nil {
		d.Set("payment_type", convertTairInstanceInstancesDBInstanceAttributeChargeTypeResponse(objectRaw["ChargeType"]))
	}
	if objectRaw["Port"] != nil {
		d.Set("port", objectRaw["Port"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["SecondaryZoneId"] != nil {
		d.Set("secondary_zone_id", objectRaw["SecondaryZoneId"])
	}
	if objectRaw["ShardCount"] != nil && objectRaw["ShardCount"] != 0 {
		d.Set("shard_count", objectRaw["ShardCount"])
	}
	if objectRaw["InstanceStatus"] != nil {
		d.Set("status", objectRaw["InstanceStatus"])
	}
	if objectRaw["StorageType"] != nil {
		d.Set("storage_performance_level", convertTairInstanceInstancesDBInstanceAttributeStorageTypeResponse(objectRaw["StorageType"]))
	}
	if objectRaw["Storage"] != nil {
		storageSizeGbVal, err := convert.FormatInt(objectRaw["Storage"])
		if err != nil {
			return err
		}
		d.Set("storage_size_gb", storageSizeGbVal)
	}
	if objectRaw["InstanceName"] != nil {
		d.Set("tair_instance_name", objectRaw["InstanceName"])
	}
	if objectRaw["VSwitchId"] != nil {
		d.Set("vswitch_id", objectRaw["VSwitchId"])
	}
	if objectRaw["VpcAuthMode"] != nil {
		d.Set("vpc_auth_mode", objectRaw["VpcAuthMode"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}
	if objectRaw["ZoneId"] != nil {
		d.Set("zone_id", objectRaw["ZoneId"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("tair_instance_id", objectRaw["InstanceId"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	objectRaw, err = redisService.DescribeBackupPolicy(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}
	if v, ok := objectRaw["PreferredBackupPeriod"]; ok && fmt.Sprint(v) != "" {
		d.Set("backup_period", strings.Split(fmt.Sprint(v), ","))
	}
	if v, ok := objectRaw["PreferredBackupTime"]; ok && fmt.Sprint(v) != "" {
		d.Set("backup_time", fmt.Sprint(v))
	}

	objectRaw, err = redisService.DescribeTairInstanceDescribeInstanceConfig(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["ParamNoLooseSentinelEnabled"] != nil {
		d.Set("param_no_loose_sentinel_enabled", objectRaw["ParamNoLooseSentinelEnabled"])
	}
	if objectRaw["ParamNoLooseSentinelPasswordFreeAccess"] != nil {
		d.Set("param_no_loose_sentinel_password_free_access", objectRaw["ParamNoLooseSentinelPasswordFreeAccess"])
	}
	if objectRaw["ParamNoLooseSentinelPasswordFreeCommands"] != nil {
		d.Set("param_no_loose_sentinel_password_free_commands", objectRaw["ParamNoLooseSentinelPasswordFreeCommands"])
	}
	if objectRaw["ParamReplMode"] != nil {
		d.Set("param_repl_mode", objectRaw["ParamReplMode"])
	}
	if objectRaw["ParamReplTimeout"] != nil {
		d.Set("param_semisync_repl_timeout", objectRaw["ParamReplTimeout"])
	}
	if objectRaw["ParamSentinelCompatEnable"] != nil {
		d.Set("param_sentinel_compat_enable", objectRaw["ParamSentinelCompatEnable"])
	}

	gotConfigs := make(map[string]interface{})
	if v, ok := d.GetOk("config"); ok && v != nil {
		gotConfigs = v.(map[string]interface{})
	}
	if objectRaw["Config"] != nil && fmt.Sprint(objectRaw["Config"]) != "" {
		configMap := make(map[string]string)
		config, err := convert.ConvertJsonStringToMap(fmt.Sprint(objectRaw["Config"]))
		if err != nil {
			return errs.WrapError(err)
		}
		for k, v := range config {
			// OpenAPI returns all configs even if the user did not specify them. Filter via gotConfigs.
			if _, ok := gotConfigs[k]; !ok && len(gotConfigs) > 0 {
				continue
			}
			configMap[k] = fmt.Sprint(v)
		}
		d.Set("config", configMap)
	}

	var securityIpGroupName string
	if v, ok := d.GetOk("security_ip_group_name"); ok {
		securityIpGroupName = v.(string)
	}
	objectRaw, err = redisService.DescribeDescribeSecurityIps(d.Id(), securityIpGroupName)
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["SecurityIpGroupName"] != nil {
		d.Set("security_ip_group_name", objectRaw["SecurityIpGroupName"])
	}
	if objectRaw["SecurityIpList"] != nil {
		d.Set("security_ips", objectRaw["SecurityIpList"])
	}

	objectRaw, err = redisService.DescribeTairInstanceDescribeSecurityGroupConfiguration(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["SecurityGroupId"] != nil {
		d.Set("security_group_id", objectRaw["SecurityGroupId"])
	}

	// Only query the TDE status when the instance actually supports it. Gating on the
	// runtime IsSupportTDE flag (instead of a static instance_type/engine_version match)
	// avoids the InstanceType.NotSupport 400 error that DescribeInstanceTDEStatus raises
	// for instances that do not support TDE, which otherwise bubbles up as a fatal read error.
	if isSupportTDE {
		objectRaw, err = redisService.DescribeTairInstanceDescribeInstanceTDEStatus(d.Id())
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		d.Set("tde_status", objectRaw["TDEStatus"])

	}

	objectRaw, err = redisService.DescribeTairInstanceDescribeInstanceSSL(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["SSLEnabled"] != nil {
		d.Set("ssl_enabled", objectRaw["SSLEnabled"])
	}
	d.Set("ca_cert_download_url", objectRaw["CertDownloadURL"])
	d.Set("ca_cert_common_name", objectRaw["CertCommonName"])
	d.Set("ssl_expired_time", objectRaw["SSLExpiredTime"])
	if objectRaw["InstanceId"] != nil {
		d.Set("tair_instance_id", objectRaw["InstanceId"])
	}

	objectRaw, err = redisService.DescribeTairInstanceDescribeIntranetAttribute(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["IntranetBandwidth"] != nil {
		d.Set("intranet_bandwidth", objectRaw["IntranetBandwidth"])
	}

	return nil
}

func resourceTairInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ModifyInstanceAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("tair_instance_name") {
		update = true
		request["InstanceName"] = d.Get("tair_instance_name")
	}

	if !d.IsNewResource() && d.HasChange("password") {
		update = true
		request["NewPassword"] = d.Get("password")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
	}
	request["InstanceClass"] = d.Get("instance_class")
	if v, ok := d.GetOkExists("force_upgrade"); ok {
		request["ForceUpgrade"] = v
	}
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	request["AutoPay"] = "true"
	if v, ok := d.GetOkExists("slave_read_only_count"); ok {
		request["SlaveReadOnlyCount"] = v
	}
	if v, ok := d.GetOkExists("read_only_count"); ok {
		request["ReadOnlyCount"] = v
	}
	if !d.IsNewResource() && d.HasChange("node_type") {
		update = true
		request["NodeType"] = d.Get("node_type")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	}
	update = false
	action = "ModifyInstanceMajorVersion"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
	}
	request["MajorVersion"] = d.Get("engine_version")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 6*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifySecurityGroupConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("security_group_id") {
		update = true
	}
	request["SecurityGroupId"] = d.Get("security_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	}
	update = false
	action = "TransformInstanceChargeType"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["ChargeType"] = convertTairInstanceChargeTypeRequest(d.Get("payment_type").(string))
	request["AutoPay"] = "true"
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	}
	update = false
	action = "ModifyInstanceSSL"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("ssl_enabled") {
		update = true
	}
	request["SSLEnabled"] = d.Get("ssl_enabled")
	redisService := NewService(client)
	objectRaw, _ := redisService.DescribeTairInstanceDescribeInstanceSSL(d.Id())
	if objectRaw["SSLEnabled"] != nil && objectRaw["SSLEnabled"] == d.Get("ssl_enabled") {
		update = false
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceBandwidth"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("intranet_bandwidth") {
		update = true
	}
	request["TargetIntranetBandwidth"] = d.Get("intranet_bandwidth")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceVpcAuthMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("vpc_auth_mode") {
		update = true
	}
	request["VpcAuthMode"] = d.Get("vpc_auth_mode")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifySecurityIps"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("security_ip_group_name") {
		update = true
	}
	if v, ok := d.GetOk("security_ip_group_name"); ok {
		request["SecurityIpGroupName"] = v
	}

	if v, ok := d.GetOk("modify_mode"); ok {
		request["ModifyMode"] = v
	}
	if d.HasChange("security_ips") {
		update = true
	}
	request["SecurityIps"] = d.Get("security_ips")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("param_repl_mode") {
		update = true
		request["ParamReplMode"] = d.Get("param_repl_mode")
	}

	if d.HasChange("param_semisync_repl_timeout") {
		update = true
		request["ParamSemisyncReplTimeout"] = d.Get("param_semisync_repl_timeout")
	}

	if d.HasChange("param_no_loose_sentinel_enabled") {
		update = true
		request["ParamNoLooseSentinelEnabled"] = d.Get("param_no_loose_sentinel_enabled")
	}

	if d.HasChange("param_sentinel_compat_enable") {
		update = true
		request["ParamSentinelCompatEnable"] = d.Get("param_sentinel_compat_enable")
	}

	if d.HasChange("param_no_loose_sentinel_password_free_access") {
		update = true
		request["ParamNoLooseSentinelPasswordFreeAccess"] = d.Get("param_no_loose_sentinel_password_free_access")
	}

	if d.HasChange("param_no_loose_sentinel_password_free_commands") {
		update = true
		request["ParamNoLooseSentinelPasswordFreeCommands"] = d.Get("param_no_loose_sentinel_password_free_commands")
	}

	if d.HasChange("config") {
		update = true
		respJson, err := convert.ConvertMaptoJsonString(d.Get("config").(map[string]interface{}))
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_redis_tair_instance", action, errs.AlibabaCloudSdkGoERROR)
		}
		request["Config"] = respJson
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceTDE"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("encryption_name"); ok {
		request["EncryptionName"] = v
	}
	if v, ok := d.GetOk("role_arn"); ok {
		request["RoleArn"] = v
	}
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if d.HasChange("tde_status") {
		update = true
	}
	request["TDEStatus"] = d.Get("tde_status")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("shard_count") {
		oldEntry, newEntry := d.GetChange("shard_count")
		oldEntryValue := oldEntry.(int)
		newEntryValue := newEntry.(int)
		removed := oldEntryValue - newEntryValue
		added := newEntryValue - oldEntryValue

		if removed > 0 {
			action := "DeleteShardingNode"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ShardCount"] = removed

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
			redisService := NewService(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}

		if added > 0 {
			action := "AddShardingNode"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = helper.BuildClientToken(action)
			request["ShardCount"] = added

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
			redisService := NewService(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisService.TairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}
	}
	update = false
	action = "ModifyInstanceMaintainTime"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("maintain_end_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_end_time"); ok {
		request["MaintainEndTime"] = v
	}
	if d.HasChange("maintain_start_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_start_time"); ok {
		request["MaintainStartTime"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	}
	if d.HasChange("backup_period") || d.HasChange("backup_time") {
		action = "ModifyBackupPolicy"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["InstanceId"] = d.Id()
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("backup_time"); ok {
			request["PreferredBackupTime"] = v
		}
		if v, ok := d.GetOk("backup_period"); ok && v.(*schema.Set).Len() > 0 {
			request["PreferredBackupPeriod"] = convert.ConvertListToCommaSeparate(v.(*schema.Set).List())
		}
		// ModifyBackupPolicy requires the full policy; wait until both parts are known
		// so a partially configured pair never reaches the API as a 400.
		if _, timeOk := request["PreferredBackupTime"]; timeOk {
			if _, periodOk := request["PreferredBackupPeriod"]; periodOk {
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
				redisService := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}
	}
	if d.HasChange("tags") {
		redisService := NewService(client)
		if err := redisService.SetResourceTags(d, "INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceTairInstanceRead(d, meta)
}

func resourceTairInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please convert it to 'PayAsYouGo' (set payment_type = \"PayAsYouGo\" and re-apply) or unsubscribe it in the console. To remove it from the state only, run 'terraform state rm'."))
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	redisService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, redisService.TairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertTairInstanceInstancesDBInstanceAttributeNodeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "double":
		return "MASTER_SLAVE"
	case "single":
		return "STAND_ALONE"
	}
	return source
}
func convertTairInstanceInstancesDBInstanceAttributeChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertTairInstanceInstancesDBInstanceAttributeStorageTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "essd_pl1":
		return "PL1"
	case "essd_pl2":
		return "PL2"
	case "essd_pl3":
		return "PL3"
	}
	return source
}
func convertTairInstanceChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
func convertTairInstanceStorageTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PL1":
		return "essd_pl1"
	case "PL2":
		return "essd_pl2"
	case "PL3":
		return "essd_pl3"
	}
	return source
}
func convertTairInstanceShardTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}

func validateConfig(v interface{}, k string) (ws []string, errors []error) {
	value, _ := v.(map[string]interface{})

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("invalid value for %s (%s)", k, value))
	}

	return
}
