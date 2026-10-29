package rkvstore

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAliCloudRedisTairInstance() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudRedisTairInstanceCreate,
		Read:   ResourceAliCloudRedisTairInstanceRead,
		Update: ResourceAliCloudRedisTairInstanceUpdate,
		Delete: ResourceAliCloudRedisTairInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("password"), cty.GetAttrPath("password_wo")),
		},
		Schema: map[string]*schema.Schema{
			"architecture_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_bandwidth": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\S+$`), "must be a data shard name without surrounding whitespace"),
						},
						"bandwidth": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
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
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				// The coupon is only used by CreateTairInstance, so changes after creation are
				// suppressed instead of forcing a new instance. The check is based on the state
				// ID because a plan-time diff is not marked as a new resource yet: keying on
				// IsNewResource() would also suppress the create-time value.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"bandwidth_burst": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cluster_backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"coupon_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRedisConfig,
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
			"instance_endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"AzIndependentEndpoint", "UnifiedEndpoint"}, false),
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
			"password_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				WriteOnly:     true,
				RequiredWith:  []string{"password_wo_version"},
				ConflictsWith: []string{"password"},
			},
			"password_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				RequiredWith: []string{"password_wo"},
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
			"replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"slave_replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"storage_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

func ResourceAliCloudRedisTairInstanceCreate(d *schema.ResourceData, meta interface{}) error {

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
		request["ChargeType"] = convertRedisTairInstanceChargeTypeRequest(v.(string))
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if woValue, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("password_wo"), cty.String); err != nil {
		return errs.WrapError(err)
	} else if !woValue.IsNull() {
		request["Password"] = woValue.AsString()
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
	applyTairInstanceReplicaCountRequest(d, request)
	if v, ok := d.GetOk("maintain_start_time"); ok {
		request["MaintainStartTime"] = v
	}
	if v, ok := d.GetOk("maintain_end_time"); ok {
		request["MaintainEndTime"] = v
	}
	if v, ok := d.GetOk("coupon_no"); ok {
		request["CouponNo"] = v
	}
	if v, ok := tairInstanceConfiguredValue(d, "auto_use_coupon"); ok {
		request["AutoUseCoupon"] = convert.ConvertBoolToString(v.(bool))
	}
	if v, ok := d.GetOk("instance_endpoint_type"); ok {
		request["InstanceEndpointType"] = v
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

	redisServiceV2 := RedisServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAliCloudRedisTairInstanceUpdate(d, meta)
}

func ResourceAliCloudRedisTairInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisServiceV2 := RedisServiceV2{client}

	objectRaw, err := redisServiceV2.DescribeRedisTairInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_tair_instance DescribeRedisTairInstance Failed!!! %s", err)
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
		storage_size_gbVal, err := convert.FormatInt(objectRaw["Storage"])
		if err != nil {
			return err
		}
		d.Set("storage_size_gb", storage_size_gbVal)
	}
	if objectRaw["ReplicaCount"] != nil {
		replicaCountVal, err := convert.FormatInt(objectRaw["ReplicaCount"])
		if err != nil {
			return err
		}
		d.Set("replica_count", replicaCountVal)
	}
	if objectRaw["SlaveReplicaCount"] != nil {
		slaveReplicaCountVal, err := convert.FormatInt(objectRaw["SlaveReplicaCount"])
		if err != nil {
			return err
		}
		d.Set("slave_replica_count", slaveReplicaCountVal)
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

	// additional_bandwidth is not part of DescribeInstanceAttribute, but DescribeRoleZoneInfo
	// reports each data node's current and default bandwidth, so the purchased additional
	// bandwidth is observable instead of being write-only.
	dataNodes, err := redisServiceV2.DescribeTairInstanceDataNodes(d.Id())
	if err != nil {
		// The instance can be removed between the attribute read and this call, and an
		// incomplete page listing must not overwrite the state with a partial value.
		if !errs.NotFoundError(err) && !errors.Is(err, errIncompleteAdditionalBandwidthNodes) {
			return errs.WrapError(err)
		}
		log.Printf("[WARN] Keeping the previous additional_bandwidth of %s: %s", d.Id(), err)
	} else if !refreshAdditionalBandwidthState(d, dataNodes) {
		log.Printf("[WARN] Keeping the previous additional_bandwidth of %s: the DescribeRoleZoneInfo listing has no row for a configured data shard", d.Id())
	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeInstanceConfig(d.Id())
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

	// OpenAPI returns every instance config, not only the ones the user configured. When config
	// is not managed in the configuration the full template is kept, so an import read and later
	// refreshes return the same map; filtering on the previous state instead would freeze the
	// template seen at the first read and make parameters changed later appear as import-only
	// additions.
	configuredConfigKeys := configuredConfigKeysFromRawConfig(d)
	if objectRaw["Config"] != nil && fmt.Sprint(objectRaw["Config"]) != "" {
		configMap := make(map[string]string)
		config, err := convert.ConvertJsonStringToMap(fmt.Sprint(objectRaw["Config"]))
		if err != nil {
			return errs.WrapError(err)
		}
		for k, v := range config {
			if len(configuredConfigKeys) > 0 {
				if _, ok := configuredConfigKeys[k]; !ok {
					continue
				}
			}
			configMap[k] = fmt.Sprint(v)
		}
		d.Set("config", configMap)
	}

	var securityIpGroupName string
	if v, ok := d.GetOk("security_ip_group_name"); ok {
		securityIpGroupName = v.(string)
	}
	objectRaw, err = redisServiceV2.DescribeDescribeSecurityIps(d.Id(), securityIpGroupName)
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["SecurityIpGroupName"] != nil {
		d.Set("security_ip_group_name", objectRaw["SecurityIpGroupName"])
	}
	if objectRaw["SecurityIpList"] != nil {
		d.Set("security_ips", objectRaw["SecurityIpList"])
	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeSecurityGroupConfiguration(d.Id())
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
		objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeInstanceTDEStatus(d.Id())
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		d.Set("tde_status", objectRaw["TDEStatus"])

	}

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeInstanceSSL(d.Id())
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

	objectRaw, err = redisServiceV2.DescribeTairInstanceDescribeIntranetAttribute(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	if objectRaw["IntranetBandwidth"] != nil {
		d.Set("intranet_bandwidth", objectRaw["IntranetBandwidth"])
	}

	return nil
}

func ResourceAliCloudRedisTairInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	// Resolve the order type before the first write of this update: a prepaid request whose
	// direction cannot be derived deterministically is refused here, so the earlier write
	// steps of this function cannot be left half-applied. Only a request that this update
	// really sends is classified: a new resource creates the instance through
	// CreateTairInstance and this update carries no specification change for it.
	var (
		orderType      string
		slaveOrderType string
		err            error
	)
	if !d.IsNewResource() {
		if tairInstanceMainSpecRequested(d) {
			orderType, err = tairInstanceSpecOrderType(d, tairInstanceMainSpecOrderTypeFields...)
			if err != nil {
				return errs.WrapError(err)
			}
		}
		if d.HasChange("slave_replica_count") {
			slaveOrderType, err = tairInstanceSpecOrderType(d, "slave_replica_count")
			if err != nil {
				return errs.WrapError(err)
			}
		}
	}
	// ModifyInstanceSpec carries the two replica counts in separate requests and the API keeps
	// ReplicaCount+SlaveReplicaCount within 4. When the primary count grows while the secondary
	// count shrinks, the secondary decrease is sent first so the intermediate sum stays within
	// the limit; every other combination keeps the primary request first.
	slaveSpecFirst := false
	if !d.IsNewResource() && d.HasChange("replica_count") && d.HasChange("slave_replica_count") {
		oldReplicaCount, newReplicaCount := d.GetChange("replica_count")
		oldSlaveReplicaCount, newSlaveReplicaCount := d.GetChange("slave_replica_count")
		slaveSpecFirst = newReplicaCount.(int) > oldReplicaCount.(int) && newSlaveReplicaCount.(int) < oldSlaveReplicaCount.(int)
	}
	d.Partial(true)

	action := "ModifyInstanceAttribute"
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
	if !d.IsNewResource() && d.HasChange("password_wo_version") {
		update = true

		if woValue, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("password_wo"), cty.String); err != nil {
			return errs.WrapError(err)
		} else if !woValue.IsNull() {
			request["NewPassword"] = woValue.AsString()
		}
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	// The secondary zone replica count is applied in its own ModifyInstanceSpec request: the two
	// counts can move in opposite directions, and a prepaid order can only declare one OrderType,
	// so each request carries the direction of the fields it actually changes.
	sendSlaveReplicaCountSpec := func() error {
		slaveRequest := buildRedisTairInstanceSlaveReplicaCountRequest(
			d.Id(),
			client.RegionId,
			helper.BuildClientToken("ModifyInstanceSpec"),
			d.Get("instance_class").(string),
			d.Get("slave_replica_count").(int),
			optionalTairInstanceUpdateOption(d, "force_upgrade"),
			optionalTairInstanceUpdateOption(d, "effective_time"),
			slaveOrderType,
		)
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		var slaveResponse map[string]interface{}
		var slaveErr error
		slaveErr = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			slaveResponse, slaveErr = client.RpcPost("R-kvstore", "2015-01-01", "ModifyInstanceSpec", nil, slaveRequest, true)
			if slaveErr != nil {
				// The first order may still be processing when this request arrives, so the
				// in-progress order code has to stay retryable here as well.
				if redisTairInstanceOrderRetryable(slaveErr) {
					wait()
					return retry.RetryableError(slaveErr)
				}
				return retry.NonRetryableError(slaveErr)
			}
			return nil
		})
		logs.AddDebug("ModifyInstanceSpec", slaveResponse, slaveRequest)
		if slaveErr != nil {
			return errs.WrapErrorf(slaveErr, errs.DefaultErrorMsg, d.Id(), "ModifyInstanceSpec", errs.AlibabaCloudSdkGoERROR)
		}
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
		return nil
	}
	if slaveSpecFirst {
		if err := sendSlaveReplicaCountSpec(); err != nil {
			return err
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
	if v := optionalTairInstanceUpdateOption(d, "force_upgrade"); v != nil {
		request["ForceUpgrade"] = v
	}
	if v := optionalTairInstanceUpdateOption(d, "effective_time"); v != nil {
		request["EffectiveTime"] = v
	}
	request["AutoPay"] = "true"
	if orderType != "" {
		request["OrderType"] = orderType
	}
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
	if !d.IsNewResource() && d.HasChange("replica_count") {
		update = true
		request["ReplicaCount"] = d.Get("replica_count")
	}
	if !d.IsNewResource() && d.HasChange("storage_size_gb") {
		update = true
		request["Storage"] = d.Get("storage_size_gb")
		// ModifyInstanceSpec requires StorageType whenever the storage size of a tair_essd
		// instance changes, even when the performance level itself is unchanged.
		if v, ok := d.Get("storage_performance_level").(string); ok && v != "" {
			request["StorageType"] = convertRedisTairInstanceStorageTypeRequest(v)
		}
	}
	if !d.IsNewResource() && d.HasChange("storage_performance_level") {
		// Removing the attribute must not serialize an empty StorageType: the API only
		// accepts essd_pl1, essd_pl2 and essd_pl3.
		if v, ok := d.Get("storage_performance_level").(string); ok && v != "" {
			update = true
			request["StorageType"] = convertRedisTairInstanceStorageTypeRequest(v)
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
			if err != nil {
				if redisTairInstanceOrderRetryable(err) {
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	if !slaveSpecFirst && !d.IsNewResource() && d.HasChange("slave_replica_count") {
		if err := sendSlaveReplicaCountSpec(); err != nil {
			return err
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 6*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
	if v, ok := d.GetOk("coupon_no"); ok {
		request["CouponNo"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
			if err != nil {
				if redisTairInstanceOrderRetryable(err) {
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
	redisServiceV2 := RedisServiceV2{client}
	objectRaw, _ := redisServiceV2.DescribeTairInstanceDescribeInstanceSSL(d.Id())
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
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
	if d.HasChanges("bandwidth_burst", "additional_bandwidth") {
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["InstanceId"] = d.Id()
		request["ChargeType"] = "PostPaid"
		// GetOkExists is intentional in this update path: a removed bandwidth_burst reads as
		// (false, true), so removal sends BandWidthBurst=false and disables burst. That is the
		// documented removal semantics and is covered by the acceptance test's removal step;
		// the GetRawConfig() alternative would change it to "leave burst as-is".
		if v, ok := d.GetOkExists("bandwidth_burst"); ok {
			request["BandWidthBurst"] = v
		}

		// additional_bandwidth is Optional+Computed, so d.Get would also return the value read
		// back from DescribeRoleZoneInfo. Only the raw configuration says what the user actually
		// configured; otherwise a value removed from the configuration would be re-applied.
		if targets, err := additionalBandwidthEntriesFromRawConfig(d); err != nil {
			return errs.WrapError(err)
		} else if len(targets) > 0 {
			tairService := RedisServiceV2{client}
			instance, err := tairService.DescribeRedisTairInstance(d.Id())
			if err != nil {
				return errs.WrapError(err)
			}
			architectureType := fmt.Sprint(instance["ArchitectureType"])
			shardCount, _ := strconv.Atoi(fmt.Sprint(instance["ShardCount"]))
			var shardNames []string
			if additionalBandwidthNeedsShardNames(targets) {
				dataNodes, err := tairService.DescribeTairInstanceDataNodes(d.Id())
				if err != nil {
					return errs.WrapError(err)
				}
				shardNames = additionalBandwidthShardNames(additionalBandwidthNodeExtras(dataNodes))
			}
			nodeIds, bandwidths, err := buildAdditionalBandwidthRequest(architectureType, shardCount, shardNames, targets)
			if err != nil {
				return errs.WrapError(err)
			}
			if nodeIds != "" {
				request["NodeId"] = nodeIds
			}
			request["Bandwidth"] = bandwidths
		}

		if len(request) > 2 {
			action = "EnableAdditionalBandwidth"
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				// autoRetry=false matches the kvstore copy of this mechanism and the adjacent
				// ModifyInstanceBandwidth flow; the outer retry.Retry still absorbs lock/task conflicts.
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, false)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"CanNotAcquireLock", "Task.Conflict"}) || helper.NeedRetry(err) {
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
			redisServiceV2 := RedisServiceV2{client}
			instanceStatusConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
			if _, err := instanceStatusConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "$.IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}
	if d.HasChange("tags") {
		redisServiceV2 := RedisServiceV2{client}
		if err := redisServiceV2.SetResourceTags(d, "INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return ResourceAliCloudRedisTairInstanceRead(d, meta)
}

func ResourceAliCloudRedisTairInstanceDelete(d *schema.ResourceData, meta interface{}) error {

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

	redisServiceV2 := RedisServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

// configuredConfigKeysFromRawConfig returns the config map keys the user actually configured.
// config is Optional+Computed, so d.Get would also return keys that were only read back from
// the API.
func configuredConfigKeysFromRawConfig(d *schema.ResourceData) map[string]struct{} {
	keys := make(map[string]struct{})
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsNull() && rawConfig.IsKnown() {
		configValue := rawConfig.GetAttr("config")
		if configValue.IsNull() || !configValue.IsKnown() {
			// config is not managed by the configuration: keep the full API template so the
			// state does not freeze a subset that later diverges from an import read.
			return keys
		}
		for k := range configValue.AsValueMap() {
			keys[k] = struct{}{}
		}
		return keys
	}
	// A plain refresh carries no raw configuration; fall back to the keys already tracked in
	// state so a configured subset is not replaced by the full API template.
	if v, ok := d.GetOk("config"); ok && v != nil {
		for k := range v.(map[string]interface{}) {
			keys[k] = struct{}{}
		}
	}
	return keys
}

// tairInstanceConfiguredValue returns the value the user configured for an attribute and
// whether that attribute is present and non-null in the configuration. Values that come only
// from API read-back (Optional+Computed) or from state are not treated as configured, so an
// explicit zero value stays distinguishable from an unset or removed attribute.
func tairInstanceConfiguredValue(d *schema.ResourceData, key string) (interface{}, bool) {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil, false
	}
	rawType := rawConfig.Type()
	if !rawType.IsObjectType() || !rawType.HasAttribute(key) {
		return nil, false
	}
	value := rawConfig.GetAttr(key)
	if value.IsNull() || !value.IsKnown() {
		return nil, false
	}
	return d.Get(key), true
}

// optionalTairInstanceUpdateOption returns the configured value of an execution option and
// nil when the user did not set it or removed it, so the option is omitted from the request
// instead of being sent with its zero value or with a stale state value.
func optionalTairInstanceUpdateOption(d *schema.ResourceData, key string) interface{} {
	if v, ok := tairInstanceConfiguredValue(d, key); ok {
		return v
	}
	return nil
}

// applyTairInstanceReplicaCountRequest adds the replica counts the user configured to a
// CreateTairInstance request. Both attributes are Optional+Computed, so the configuration is
// inspected instead of GetOkExists to keep an explicit 0 distinguishable from an unset value.
func applyTairInstanceReplicaCountRequest(d *schema.ResourceData, request map[string]interface{}) {
	if v, ok := tairInstanceConfiguredValue(d, "replica_count"); ok {
		request["ReplicaCount"] = v
	}
	if v, ok := tairInstanceConfiguredValue(d, "slave_replica_count"); ok {
		request["SlaveReplicaCount"] = v
	}
}

const (
	tairInstanceSpecDirectionUnchanged = iota
	tairInstanceSpecDirectionUpgrade
	tairInstanceSpecDirectionDowngrade
	tairInstanceSpecDirectionUnknown
)

// tairInstanceMainSpecOrderTypeFields are every price-bearing field the main
// ModifyInstanceSpec request can carry. read_only_count and slave_read_only_count are carried
// by that request as well, so the classification covers the whole serialized request.
var tairInstanceMainSpecOrderTypeFields = []string{
	"instance_class",
	"node_type",
	"replica_count",
	"storage_size_gb",
	"storage_performance_level",
	"read_only_count",
	"slave_read_only_count",
}

// tairInstanceMainSpecRequested reports whether this update sends the main ModifyInstanceSpec
// request. It mirrors the branch conditions of the request builder, so a create-time update
// that only writes the newly created instance is not classified as a specification change.
func tairInstanceMainSpecRequested(d *schema.ResourceData) bool {
	if d.HasChange("instance_class") || d.HasChange("node_type") || d.HasChange("replica_count") || d.HasChange("storage_size_gb") {
		return true
	}
	if d.HasChange("storage_performance_level") {
		level, _ := d.Get("storage_performance_level").(string)
		return level != ""
	}
	return false
}

// tairInstanceSpecOrderType returns the ModifyInstanceSpec OrderType for one prepaid
// request. The API compares the price of the target specification carried by that request
// with the current one, so the direction is derived from every price-bearing field of the
// request: a request mixing increases and decreases, or combining a field whose direction
// cannot be compared locally with another changed field, is refused instead of being sent
// with a contradictory OrderType. A single field that cannot be compared locally keeps the
// API default (UPGRADE), and an empty result keeps the default for every other request.
func tairInstanceSpecOrderType(d *schema.ResourceData, fields ...string) (string, error) {
	// ModifyInstanceSpec runs before TransformInstanceChargeType in this update, so the order
	// is priced with the billing type the instance still has in state. Reading the planned
	// payment_type would classify a conversion that has not happened yet.
	previousPaymentType, _ := d.GetChange("payment_type")
	if paymentType, _ := previousPaymentType.(string); paymentType != "Subscription" {
		return "", nil
	}
	var upgrades, downgrades, unknown []string
	for _, field := range fields {
		if !d.HasChange(field) {
			continue
		}
		switch tairInstanceSpecFieldDirection(d, field) {
		case tairInstanceSpecDirectionUpgrade:
			upgrades = append(upgrades, field)
		case tairInstanceSpecDirectionDowngrade:
			downgrades = append(downgrades, field)
		case tairInstanceSpecDirectionUnknown:
			unknown = append(unknown, field)
		}
	}
	if len(unknown) > 0 {
		if len(upgrades) > 0 || len(downgrades) > 0 {
			changed := append(append([]string{}, upgrades...), downgrades...)
			return "", fmt.Errorf("cannot determine the ModifyInstanceSpec OrderType for the prepaid instance: the price direction of %s cannot be compared with %s in the same request, so the update is refused instead of risking a wrong order", strings.Join(unknown, ", "), strings.Join(changed, ", "))
		}
		return "", nil
	}
	if len(upgrades) > 0 && len(downgrades) > 0 {
		return "", fmt.Errorf("cannot determine the ModifyInstanceSpec OrderType for the prepaid instance: %s increase while %s decrease in the same request, so the update is refused instead of risking a wrong order", strings.Join(upgrades, ", "), strings.Join(downgrades, ", "))
	}
	if len(downgrades) > 0 {
		return "DOWNGRADE", nil
	}
	return "", nil
}

// tairInstanceSpecFieldDirection classifies one price-bearing field of a ModifyInstanceSpec
// request. Numeric sizes, read-only counts and ESSD performance levels are comparable locally;
// the instance class and the node architecture are not, because a smaller capacity or a
// different node layout can cost more than the current one.
func tairInstanceSpecFieldDirection(d *schema.ResourceData, field string) int {
	switch field {
	case "replica_count", "slave_replica_count", "storage_size_gb", "read_only_count", "slave_read_only_count":
		oldValue, newValue := d.GetChange(field)
		oldCount, _ := oldValue.(int)
		newCount, _ := newValue.(int)
		switch {
		case newCount > oldCount:
			return tairInstanceSpecDirectionUpgrade
		case newCount < oldCount:
			return tairInstanceSpecDirectionDowngrade
		default:
			return tairInstanceSpecDirectionUnchanged
		}
	case "storage_performance_level":
		oldValue, newValue := d.GetChange(field)
		newLevel := fmt.Sprint(newValue)
		if newLevel == "" {
			// Removing the attribute does not serialize StorageType.
			return tairInstanceSpecDirectionUnchanged
		}
		oldRank, oldKnown := redisTairInstanceStoragePerformanceLevelRank(fmt.Sprint(oldValue))
		newRank, newKnown := redisTairInstanceStoragePerformanceLevelRank(newLevel)
		if !oldKnown || !newKnown {
			return tairInstanceSpecDirectionUnknown
		}
		switch {
		case newRank > oldRank:
			return tairInstanceSpecDirectionUpgrade
		case newRank < oldRank:
			return tairInstanceSpecDirectionDowngrade
		default:
			return tairInstanceSpecDirectionUnchanged
		}
	default:
		return tairInstanceSpecDirectionUnknown
	}
}

// redisTairInstanceStoragePerformanceLevelRank orders the ESSD performance levels from the
// lowest to the highest so a level change can be classified as an upgrade or a downgrade.
func redisTairInstanceStoragePerformanceLevelRank(level string) (int, bool) {
	switch level {
	case "PL1":
		return 1, true
	case "PL2":
		return 2, true
	case "PL3":
		return 3, true
	default:
		return 0, false
	}
}

// redisTairInstanceOrderRetryable reports whether an order-style RPC should be retried:
// transient errors plus the order code returned while the previous order is still being
// processed. Every ModifyInstanceSpec and TransformInstanceChargeType retry uses it so the
// same in-progress condition is handled consistently.
func redisTairInstanceOrderRetryable(err error) bool {
	return helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Order.LatestOrderIsHanding"})
}

// buildRedisTairInstanceSlaveReplicaCountRequest builds the standalone ModifyInstanceSpec
// request that changes only SlaveReplicaCount after the request carrying ReplicaCount has
// converged, so each prepaid order can declare the direction of the count it changes. User
// configured execution options are carried over; omitting them would make the API fall back
// to its defaults (force upgrade enabled, run immediately).
func buildRedisTairInstanceSlaveReplicaCountRequest(instanceId, regionId, clientToken, instanceClass string, slaveReplicaCount int, forceUpgrade, effectiveTime interface{}, orderType string) map[string]interface{} {
	request := map[string]interface{}{
		"InstanceId":        instanceId,
		"RegionId":          regionId,
		"ClientToken":       clientToken,
		"InstanceClass":     instanceClass,
		"SlaveReplicaCount": slaveReplicaCount,
	}
	if forceUpgrade != nil {
		request["ForceUpgrade"] = forceUpgrade
	}
	if effectiveTime != nil {
		request["EffectiveTime"] = effectiveTime
	}
	if orderType != "" {
		request["OrderType"] = orderType
	}
	return request
}
