package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRedisTairInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRedisTairInstanceCreate,
		Read:   resourceAliCloudRedisTairInstanceRead,
		Update: resourceAliCloudRedisTairInstanceUpdate,
		Delete: resourceAliCloudRedisTairInstanceDelete,
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
			"cluster_backup_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: StringInSlice([]string{"Immediately", "MaintainTime"}, false),
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
				ValidateFunc: StringInSlice([]string{"tair_rdb", "tair_scm", "tair_essd"}, false),
			},
			"intranet_bandwidth": {
				Type:     schema.TypeInt,
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
				ValidateFunc: StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE", "double", "single"}, false),
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
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60}),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(1024, 65535),
			},
			"read_only_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 5),
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
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: StringInSlice([]string{"Disable", "Enable", "Update"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"storage_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
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

func resourceAliCloudRedisTairInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTairInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

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
		request["ChargeType"] = convertRedisTairInstanceChargeTypeRequest(v.(string))
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
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOkExists("storage_size_gb"); ok && v.(int) > 0 {
		request["Storage"] = v
	}
	if v, ok := d.GetOk("storage_performance_level"); ok {
		request["StorageType"] = convertRedisTairInstanceStorageTypeRequest(v.(string))
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
		request["ShardType"] = convertRedisTairInstanceShardTypeRequest(v.(string))
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_redis_tair_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	redisServiceV2 := RedisServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRedisTairInstanceUpdate(d, meta)
}

func resourceAliCloudRedisTairInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisServiceV2 := RedisServiceV2{client}

	objectRaw, err := redisServiceV2.DescribeRedisTairInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_tair_instance DescribeRedisTairInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	if objectRaw["Connections"] != nil {
		d.Set("max_connections", objectRaw["Connections"])
	}
	if objectRaw["NetworkType"] != nil {
		d.Set("network_type", objectRaw["NetworkType"])
	}
	if objectRaw["NodeType"] != nil {
		d.Set("node_type", convertRedisTairInstanceInstancesDBInstanceAttributeNodeTypeResponse(objectRaw["NodeType"]))
	}
	if objectRaw["ChargeType"] != nil {
		d.Set("payment_type", convertRedisTairInstanceInstancesDBInstanceAttributeChargeTypeResponse(objectRaw["ChargeType"]))
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
		d.Set("storage_performance_level", convertRedisTairInstanceInstancesDBInstanceAttributeStorageTypeResponse(objectRaw["StorageType"]))
	}
	if objectRaw["Storage"] != nil {
		d.Set("storage_size_gb", formatInt(objectRaw["Storage"]))
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
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeInstanceConfig(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
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

	var securityIpGroupName string
	if v, ok := d.GetOk("security_ip_group_name"); ok {
		securityIpGroupName = v.(string)
	}
	objectRaw, err = redisServiceV2.DescribeDescribeSecurityIps(d.Id(), securityIpGroupName)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if objectRaw["SecurityIpGroupName"] != nil {
		d.Set("security_ip_group_name", objectRaw["SecurityIpGroupName"])
	}
	if objectRaw["SecurityIpList"] != nil {
		d.Set("security_ips", objectRaw["SecurityIpList"])
	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeSecurityGroupConfiguration(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["SecurityGroupId"] != nil {
		d.Set("security_group_id", objectRaw["SecurityGroupId"])
	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeInstanceSSL(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if objectRaw["SSLEnabled"] != nil {
		d.Set("ssl_enabled", objectRaw["SSLEnabled"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("tair_instance_id", objectRaw["InstanceId"])
	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeIntranetAttribute(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if objectRaw["IntranetBandwidth"] != nil {
		d.Set("intranet_bandwidth", objectRaw["IntranetBandwidth"])
	}

	return nil
}

func resourceAliCloudRedisTairInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 6*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	request["ChargeType"] = convertRedisTairInstanceChargeTypeRequest(d.Get("payment_type").(string))
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	redisServiceV2 := RedisServiceV2{client}
	objectRaw, _ := redisServiceV2.DescribeTairInstanceDescribeInstanceSSL(d.Id())
	if objectRaw["SSLEnabled"] != nil && objectRaw["SSLEnabled"] == d.Get("ssl_enabled") {
		update = false
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added > 0 {
			action := "AddShardingNode"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["ShardCount"] = added

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}
	if d.HasChange("tags") {
		redisServiceV2 := RedisServiceV2{client}
		if err := redisServiceV2.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudRedisTairInstanceRead(d, meta)
}

func resourceAliCloudRedisTairInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_redis_tair_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	redisServiceV2 := RedisServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertRedisTairInstanceInstancesDBInstanceAttributeNodeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "double":
		return "MASTER_SLAVE"
	case "single":
		return "STAND_ALONE"
	}
	return source
}
func convertRedisTairInstanceInstancesDBInstanceAttributeChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertRedisTairInstanceInstancesDBInstanceAttributeStorageTypeResponse(source interface{}) interface{} {
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
func convertRedisTairInstanceChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
func convertRedisTairInstanceStorageTypeRequest(source interface{}) interface{} {
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
func convertRedisTairInstanceShardTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
