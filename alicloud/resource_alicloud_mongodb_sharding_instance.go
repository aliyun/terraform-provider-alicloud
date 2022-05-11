package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMongoDBShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBShardingInstanceCreate,
		Read:   resourceAlicloudMongoDBShardingInstanceRead,
		Update: resourceAlicloudMongoDBShardingInstanceUpdate,
		Delete: resourceAlicloudMongoDBShardingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger", "RocksDB"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Optional:     true,
				Computed:     true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" || d.Get("engine_version").(string) < "4.0"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled"}, false),
				Optional:     true,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"readonly_replicas": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 5),
							Computed:     true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},
			"mongo_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"tags": tagsSchema(),
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"config_server_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "VPC"}, false),
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"mongodb", "dynamodb"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateShardingDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["Engine"] = "MongoDB"
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["DBInstanceDescription"] = v
	}

	request["AccountPassword"] = d.Get("account_password").(string)
	if request["AccountPassword"] == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["AccountPassword"] = decryptResp
		}
	}

	if v, ok := d.GetOk("shard_list"); ok {
		shardListMaps := make([]map[string]interface{}, 0)
		shardList := v.([]interface{})
		for _, rew := range shardList {
			shardListMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			shardListMap["Class"] = item["node_class"]
			shardListMap["Storage"] = item["node_storage"]
			shardListMap["ReadonlyReplicas"] = item["readonly_replicas"]
			shardListMaps = append(shardListMaps, shardListMap)
		}
		request["ReplicaSet"] = shardListMaps
	}

	if v, ok := d.GetOk("mongo_list"); ok {
		mongoListMaps := make([]map[string]interface{}, 0)
		mongoList := v.([]interface{})
		for _, rew := range mongoList {
			mongoListMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			mongoListMap["Class"] = item["node_class"]
			mongoListMaps = append(mongoListMaps, mongoListMap)
		}
		request["Mongos"] = mongoListMaps
	}
	request["ConfigServer"] = []map[string]interface{}{{"Class": "dds.cs.mid", "Storage": 20}}

	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if (request["ZoneId"] == nil || request["VpcId"] == nil) && request["VSwitchId"] != nil {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if request["ZoneId"] == nil {
			request["ZoneId"] = vsw.ZoneId
		} else if strings.Contains(request["ZoneId"].(string), MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(request["ZoneId"].(string), "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in multi the zone " + request["ZoneId"].(string)))
			}
		} else if request["ZoneId"].(string) != vsw.ZoneId {
			return WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in the zone " + request["ZoneId"].(string)))
		}
		if request["VpcId"] == nil {
			request["VpcId"] = vsw.VpcId
		}
	}
	if request["NetworkType"] == nil && request["VSwitchId"] != nil {
		request["NetworkType"] = "VPC"
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["ChargeType"] = v
		if vv, ok := d.GetOk("period"); ok && PayType(v.(string)) == PrePaid {
			request["Period"] = vv
		}
	}

	request["SecurityIPList"] = LOCAL_HOST_IP
	if v, ok := d.GetOk("security_ip_list"); ok {
		request["SecurityIPList"] = strings.Join(expandStringList(v.(*schema.Set).List()), COMMA_SEPARATED)
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = strconv.FormatBool(v.(bool))
	}
	if v, ok := d.GetOk("storage_engine"); ok {
		request["StorageEngine"] = v
	}
	if v, ok := d.GetOk("protocol_type"); ok {
		request["ProtocolType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_instance", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["DBInstanceId"]))

	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBShardingInstanceUpdate(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	instance, err := ddsService.DescribeMongoDBShardingInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	backupPolicy, err := ddsService.DescribeMongoDBShardingBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_time", backupPolicy["PreferredBackupTime"])
	d.Set("backup_period", strings.Split(backupPolicy["PreferredBackupPeriod"].(string), ","))
	d.Set("retention_period", formatInt(backupPolicy["BackupRetentionPeriod"]))

	d.Set("name", instance["DBInstanceDescription"])
	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("network_type", instance["NetworkType"])
	d.Set("storage_engine", instance["StorageEngine"])
	d.Set("protocol_type", instance["ProtocolType"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("vpc_id", instance["VPCId"])
	d.Set("instance_charge_type", instance["ChargeType"])
	if instance["ChargeType"] == "PrePaid" {
		period, err := computePeriodByUnit(instance["CreationTime"], instance["ExpireTime"], d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}
	d.Set("vswitch_id", instance["VSwitchId"])

	if v, ok := instance["MongosList"]; ok {
		mongosMaps := make([]map[string]interface{}, 0)
		mongosList := v.(map[string]interface{})
		if v, ok := mongosList["MongosAttribute"]; ok {
			if mongosAttributeList, ok := v.([]interface{}); ok && len(mongosAttributeList) > 0 {
				for _, mongosMap := range mongosAttributeList {
					item := mongosMap.(map[string]interface{})
					mongo := map[string]interface{}{
						"node_class":     item["NodeClass"],
						"node_id":        item["NodeId"],
						"port":           formatInt(item["Port"]),
						"connect_string": item["ConnectSting"],
					}
					mongosMaps = append(mongosMaps, mongo)
				}
				d.Set("mongo_list", mongosMaps)
			}
		}
	}

	if v, ok := instance["ShardList"]; ok {
		shardMaps := make([]map[string]interface{}, 0)
		shardList := v.(map[string]interface{})
		if v, ok := shardList["ShardAttribute"]; ok {
			if shardAttributeList, ok := v.([]interface{}); ok && len(shardAttributeList) > 0 {
				for _, shardMap := range shardAttributeList {
					item := shardMap.(map[string]interface{})
					shard := map[string]interface{}{
						"node_id":           item["NodeId"],
						"node_storage":      formatInt(item["NodeStorage"]),
						"node_class":        item["NodeClass"],
						"readonly_replicas": formatInt(item["ReadonlyReplicas"]),
					}
					shardMaps = append(shardMaps, shard)
				}
				d.Set("shard_list", shardMaps)
			}
		}
	}

	tdeInfo, err := ddsService.DescribeMongoDBShardingTDEInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo["TDEStatus"])

	ips, err := ddsService.DescribeMongoDBShardingSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ip_list", ips)
	groupIp, err := ddsService.DescribeMongoDBShardingSecurityGroupId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(groupIp) > 0 {
		d.Set("security_group_id", groupIp[0].(map[string]interface{})["SecurityGroupId"])
	}

	d.Set("tags", tagsToMap(instance["Tags"].(map[string]interface{})["Tag"]))

	if v, ok := instance["ConfigserverList"]; ok {
		configserverMaps := make([]map[string]interface{}, 0)
		configserverList := v.(map[string]interface{})
		if v, ok := configserverList["ConfigserverAttribute"]; ok {
			if configserverAttributeList, ok := v.([]interface{}); ok && len(configserverAttributeList) > 0 {
				for _, configserverMap := range configserverAttributeList {
					item := configserverMap.(map[string]interface{})
					configserver := map[string]interface{}{
						"max_iops":         item["MaxIOPS"],
						"connect_string":   item["ConnectString"],
						"node_class":       item["NodeClass"],
						"max_connections":  item["MaxConnections"],
						"port":             formatInt(item["Port"]),
						"node_description": item["NodeDescription"],
						"node_id":          item["NodeId"],
						"node_storage":     item["NodeStorage"],
					}
					configserverMaps = append(configserverMaps, configserver)
				}
				d.Set("config_server_list", configserverMaps)
			}
		}
	}

	return nil
}

func resourceAlicloudMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("backup_time") || d.HasChange("backup_period") {
		if err := ddsService.MotifyMongoDBBackupPolicy(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("backup_time")
		d.SetPartial("backup_period")
	}
	if d.HasChange("tde_status") {
		var response map[string]interface{}
		action := "ModifyDBInstanceTDE"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		request["TDEStatus"] = d.Get("tde_status").(string)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("tde_status")
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("security_group_id") {
		var response map[string]interface{}
		action := "ModifySecurityGroupConfiguration"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		request["SecurityGroupId"] = d.Get("security_group_id").(string)

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"InstanceStatusInvalid"}) || NeedRetry(err) {
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
		d.SetPartial("security_group_id")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
	}

	if d.HasChange("resource_group_id") {
		if v, ok := d.GetOk("resource_group_id"); ok {
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
			action := "ModifyResourceGroup"
			request := map[string]interface{}{
				"DBInstanceId":    d.Id(),
				"ResourceGroupId": v,
				"RegionId":        client.RegionId,
			}
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, request)
			d.SetPartial("resource_group_id")
		}
	}

	if d.HasChange("shard_list") {
		state, diff := d.GetChange("shard_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeShard, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("shard_list")
	}

	if d.HasChange("mongo_list") {
		state, diff := d.GetChange("mongo_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeMongos, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("mongo_list")
	}

	if d.HasChange("name") {
		var response map[string]interface{}
		action := "ModifyDBInstanceDescription"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		if v, ok := d.GetOk("name"); ok {
			request["DBInstanceDescription"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("name")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		var accountPassword string
		if accountPassword = d.Get("account_password").(string); accountPassword != "" {
			d.SetPartial("account_password")
		} else if kmsPassword := d.Get("kms_encrypted_password").(string); kmsPassword != "" {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			accountPassword = decryptResp
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
		}

		err := ddsService.ResetAccountPassword(d, accountPassword)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("account_password")
	}

	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := ddsService.ModifyMongoDBSecurityIps(d, ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}
	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		var response map[string]interface{}
		action := "TransformToPrePaid"
		prePaidRequest := make(map[string]interface{})
		prePaidRequest["InstanceId"] = d.Id()
		prePaidRequest["AutoPay"] = requests.NewBoolean(true)
		prePaidRequest["Period"] = requests.NewInteger(d.Get("period").(int))
		if v, ok := d.GetOk("auto_renew"); ok {
			prePaidRequest["AutoRenew"] = strconv.FormatBool(v.(bool))
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, prePaidRequest, &util.RuntimeOptions{})
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
		addDebug(action, response, prePaidRequest)
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("period")
	}
	d.Partial(false)
	return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		log.Printf("[WARN] Cannot destroy resourceAlicloudMongoDBShardingInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return nil
}
