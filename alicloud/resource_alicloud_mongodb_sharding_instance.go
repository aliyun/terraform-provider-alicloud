package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMongoDBShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongoDBShardingInstanceCreate,
		Read:   resourceAliCloudMongoDBShardingInstanceRead,
		Update: resourceAliCloudMongoDBShardingInstanceUpdate,
		Delete: resourceAliCloudMongoDBShardingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"WiredTiger", "RocksDB"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"cloud_essd1", "cloud_essd2", "cloud_essd3", "cloud_auto", "local_ssd"}, false),
			},
			"provisioned_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"mongodb", "dynamodb"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hidden_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Classic", "VPC"}, false),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Elem:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice(BACKUP_TIME, false),
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snapshot_backup_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "Flash"}, false),
			},
			"backup_interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"-1", "15", "30", "60", "120", "180", "240", "360", "480", "720"}, false),
			},
			"tde_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine_version").(string) < "4.0"
				},
			},
			"db_instance_release_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"global_security_group_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mongo_list": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 2,
				MaxItems: 32,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
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
			},
			"shard_list": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 2,
				MaxItems: 32,
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
							Computed:     true,
							ValidateFunc: IntBetween(0, 5),
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"config_server_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
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
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
		},
	}
}

func resourceAliCloudMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	var response map[string]interface{}
	action := "CreateShardingDBInstance"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["Engine"] = "MongoDB"
	request["EngineVersion"] = d.Get("engine_version")

	if v, ok := d.GetOk("storage_engine"); ok {
		request["StorageEngine"] = v
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v
	}

	if v, ok := d.GetOkExists("provisioned_iops"); ok {
		request["ProvisionedIops"] = v
	}

	if v, ok := d.GetOk("protocol_type"); ok {
		request["ProtocolType"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
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
				return WrapError(Error("The specified vswitch %s isn't in multi the zone %s", vsw.VSwitchId, request["ZoneId"].(string)))
			}
		} else if request["ZoneId"].(string) != vsw.ZoneId {
			return WrapError(Error("The specified vswitch %s isn't in the zone %s", vsw.VSwitchId, request["ZoneId"].(string)))
		}
		if request["VpcId"] == nil {
			request["VpcId"] = vsw.VpcId
		}
	}

	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}

	if v, ok := d.GetOk("hidden_zone_id"); ok {
		request["HiddenZoneId"] = v
	}

	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}

	if request["NetworkType"] == nil && request["VSwitchId"] != nil {
		request["NetworkType"] = "VPC"
	}

	if v, ok := d.GetOk("name"); ok {
		request["DBInstanceDescription"] = v
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["ChargeType"] = v

		if period, ok := d.GetOk("period"); ok && PayType(v.(string)) == PrePaid {
			request["Period"] = period
		}
	}

	request["SecurityIPList"] = LOCAL_HOST_IP
	if v, ok := d.GetOk("security_ip_list"); ok {
		request["SecurityIPList"] = strings.Join(expandStringList(v.(*schema.Set).List()), COMMA_SEPARATED)
	}

	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	} else if v, ok := d.GetOk("kms_encrypted_password"); ok {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}

		request["AccountPassword"] = decryptResp
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = strconv.FormatBool(v.(bool))
	}

	mongoList := d.Get("mongo_list")
	mongoListMaps := make([]map[string]interface{}, 0)
	for _, mongoLists := range mongoList.([]interface{}) {
		mongoListMap := map[string]interface{}{}
		mongoListArg := mongoLists.(map[string]interface{})

		mongoListMap["Class"] = mongoListArg["node_class"]

		mongoListMaps = append(mongoListMaps, mongoListMap)
	}

	request["Mongos"] = mongoListMaps

	shardList := d.Get("shard_list")
	shardListMaps := make([]map[string]interface{}, 0)
	for _, shardLists := range shardList.([]interface{}) {
		shardListMap := map[string]interface{}{}
		shardListArg := shardLists.(map[string]interface{})

		shardListMap["Class"] = shardListArg["node_class"]
		shardListMap["Storage"] = shardListArg["node_storage"]

		if readonlyReplicas, ok := shardListArg["readonly_replicas"]; ok {
			shardListMap["ReadonlyReplicas"] = readonlyReplicas
		}

		shardListMaps = append(shardListMaps, shardListMap)
	}

	request["ReplicaSet"] = shardListMaps

	if v, ok := d.GetOk("config_server_list"); ok {
		configServerListMaps := make([]map[string]interface{}, 0)
		for _, configServerList := range v.([]interface{}) {
			configServerListMap := map[string]interface{}{}
			configServerListArg := configServerList.(map[string]interface{})

			if class, ok := configServerListArg["node_class"]; ok {
				configServerListMap["Class"] = class
			}

			if storage, ok := configServerListArg["node_storage"]; ok {
				configServerListMap["Storage"] = storage
			}

			configServerListMaps = append(configServerListMaps, configServerListMap)
		}

		request["ConfigServer"] = configServerListMaps
	} else {
		request["ConfigServer"] = []map[string]interface{}{{"Class": "dds.cs.mid", "Storage": 20}}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudMongoDBShardingInstanceUpdate(d, meta)
}

func resourceAliCloudMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	object, err := ddsService.DescribeMongoDBShardingInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine_version", object["EngineVersion"])
	d.Set("storage_engine", object["StorageEngine"])
	d.Set("storage_type", object["StorageType"])
	d.Set("provisioned_iops", formatInt(object["ProvisionedIops"]))
	d.Set("protocol_type", object["ProtocolType"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("secondary_zone_id", object["SecondaryZoneId"])
	d.Set("hidden_zone_id", object["HiddenZoneId"])
	d.Set("network_type", object["NetworkType"])
	d.Set("name", object["DBInstanceDescription"])
	d.Set("instance_charge_type", object["ChargeType"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("db_instance_release_protection", object["DBInstanceReleaseProtection"])

	groupIp, err := ddsService.DescribeMongoDBShardingSecurityGroupId(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if len(groupIp) > 0 {
		if groupIpItem, ok := groupIp[0].(map[string]interface{}); ok {
			d.Set("security_group_id", groupIpItem["SecurityGroupId"])
		}
	}

	if fmt.Sprint(object["ChargeType"]) == "PrePaid" {
		period, err := computePeriodByUnit(object["CreationTime"], object["ExpireTime"], d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}

		d.Set("period", period)
	}

	securityIpList, err := ddsService.DescribeMongoDBShardingSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ip_list", securityIpList)

	backupPolicy, err := ddsService.DescribeMongoDBShardingBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("backup_time", backupPolicy["PreferredBackupTime"])
	d.Set("backup_period", strings.Split(backupPolicy["PreferredBackupPeriod"].(string), ","))
	d.Set("retention_period", formatInt(backupPolicy["BackupRetentionPeriod"]))
	d.Set("backup_retention_policy_on_cluster_deletion", formatInt(backupPolicy["BackupRetentionPolicyOnClusterDeletion"]))
	d.Set("snapshot_backup_type", backupPolicy["SnapshotBackupType"])
	d.Set("backup_interval", backupPolicy["BackupInterval"])

	tdeInfo, err := ddsService.DescribeMongoDBShardingTDEInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("tde_status", tdeInfo["TDEStatus"])

	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	globalSecurityGroupIds, err := ddsService.DescribeMongoDBGlobalSecurityGroupIds(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("global_security_group_list", globalSecurityGroupIds)

	if MongosListMap, ok := object["MongosList"].(map[string]interface{}); ok && MongosListMap != nil {
		if MongosList, ok := MongosListMap["MongosAttribute"]; ok && MongosList != nil {
			MongosListMaps := make([]map[string]interface{}, 0)
			for _, MongosLists := range MongosList.([]interface{}) {
				MongosListItemMap := make(map[string]interface{})
				MongosListArg := MongosLists.(map[string]interface{})

				if nodeClass, ok := MongosListArg["NodeClass"]; ok {
					MongosListItemMap["node_class"] = nodeClass
				}

				if nodeId, ok := MongosListArg["NodeId"]; ok {
					MongosListItemMap["node_id"] = nodeId
				}

				if connectSting, ok := MongosListArg["ConnectSting"]; ok {
					MongosListItemMap["connect_string"] = connectSting
				}

				if port, ok := MongosListArg["Port"]; ok {
					MongosListItemMap["port"] = formatInt(port)
				}

				MongosListMaps = append(MongosListMaps, MongosListItemMap)
			}

			d.Set("mongo_list", MongosListMaps)
		}
	}

	if shardListMap, ok := object["ShardList"].(map[string]interface{}); ok && shardListMap != nil {
		if shardList, ok := shardListMap["ShardAttribute"]; ok && shardList != nil {
			shardListMaps := make([]map[string]interface{}, 0)
			for _, shardLists := range shardList.([]interface{}) {
				shardListItemMap := make(map[string]interface{})
				shardListArg := shardLists.(map[string]interface{})

				if nodeClass, ok := shardListArg["NodeClass"]; ok {
					shardListItemMap["node_class"] = nodeClass
				}

				if nodeStorage, ok := shardListArg["NodeStorage"]; ok {
					shardListItemMap["node_storage"] = formatInt(nodeStorage)
				}

				if readonlyReplicas, ok := shardListArg["ReadonlyReplicas"]; ok {
					shardListItemMap["readonly_replicas"] = formatInt(readonlyReplicas)
				}

				if nodeId, ok := shardListArg["NodeId"]; ok {
					shardListItemMap["node_id"] = nodeId
				}

				shardListMaps = append(shardListMaps, shardListItemMap)
			}

			d.Set("shard_list", shardListMaps)
		}
	}

	if configServerListMap, ok := object["ConfigserverList"].(map[string]interface{}); ok && configServerListMap != nil {
		if configServerList, ok := configServerListMap["ConfigserverAttribute"]; ok && configServerList != nil {
			configServerListMaps := make([]map[string]interface{}, 0)
			for _, configServerLists := range configServerList.([]interface{}) {
				configServerListItemMap := make(map[string]interface{})
				configServerListArg := configServerLists.(map[string]interface{})

				if nodeClass, ok := configServerListArg["NodeClass"]; ok {
					configServerListItemMap["node_class"] = nodeClass
				}

				if nodeStorage, ok := configServerListArg["NodeStorage"]; ok {
					configServerListItemMap["node_storage"] = nodeStorage
				}

				if nodeId, ok := configServerListArg["NodeId"]; ok {
					configServerListItemMap["node_id"] = nodeId
				}

				if connectString, ok := configServerListArg["ConnectString"]; ok {
					configServerListItemMap["connect_string"] = connectString
				}

				if port, ok := configServerListArg["Port"]; ok {
					configServerListItemMap["port"] = formatInt(port)
				}

				if maxConnections, ok := configServerListArg["MaxConnections"]; ok {
					configServerListItemMap["max_connections"] = maxConnections
				}

				if maxIOPS, ok := configServerListArg["MaxIOPS"]; ok {
					configServerListItemMap["max_iops"] = maxIOPS
				}

				if nodeDescription, ok := configServerListArg["NodeDescription"]; ok {
					configServerListItemMap["node_description"] = nodeDescription
				}

				configServerListMaps = append(configServerListMaps, configServerListItemMap)
			}

			d.Set("config_server_list", configServerListMaps)
		}
	}

	return nil
}

func resourceAliCloudMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	upgradeDBInstanceEngineVersionReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
	}
	upgradeDBInstanceEngineVersionReq["EngineVersion"] = d.Get("engine_version")

	if update {
		action := "UpgradeDBInstanceEngineVersion"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, upgradeDBInstanceEngineVersionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, upgradeDBInstanceEngineVersionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("engine_version")
	}

	update = false
	modifyDBInstanceDiskTypeReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("storage_type") {
		update = true
	}
	if v, ok := d.GetOk("storage_type"); ok {
		modifyDBInstanceDiskTypeReq["DbInstanceStorageType"] = v
	}

	if !d.IsNewResource() && d.HasChange("provisioned_iops") {
		update = true

		if v, ok := d.GetOkExists("provisioned_iops"); ok {
			modifyDBInstanceDiskTypeReq["ProvisionedIops"] = v
		}
	}

	if update {
		action := "ModifyDBInstanceDiskType"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyDBInstanceDiskTypeReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceDiskTypeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("storage_type")
		d.SetPartial("provisioned_iops")
	}

	update = false
	migrateAvailableZoneReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"Vswitch":      d.Get("vswitch_id").(string),
		"ZoneId":       d.Get("zone_id").(string),
	}

	if !d.IsNewResource() && d.HasChange("secondary_zone_id") {
		update = true
	}
	if v, ok := d.GetOk("secondary_zone_id"); ok {
		migrateAvailableZoneReq["SecondaryZoneId"] = v
	}

	if !d.IsNewResource() && d.HasChange("hidden_zone_id") {
		update = true
	}
	if v, ok := d.GetOk("hidden_zone_id"); ok {
		migrateAvailableZoneReq["HiddenZoneId"] = v
	}

	if update {
		action := "MigrateAvailableZone"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, migrateAvailableZoneReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, migrateAvailableZoneReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("secondary_zone_id")
		d.SetPartial("hidden_zone_id")
	}

	update = false
	modifySecurityGroupConfigurationReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if d.HasChange("security_group_id") {
		update = true
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		modifySecurityGroupConfigurationReq["SecurityGroupId"] = v.(string)
	}

	if update {
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		action := "ModifySecurityGroupConfiguration"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifySecurityGroupConfigurationReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InstanceStatusInvalid"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySecurityGroupConfigurationReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("security_group_id")
	}

	update = false
	modifyDBInstanceDescriptionReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
	}
	if v, ok := d.GetOk("name"); ok {
		modifyDBInstanceDescriptionReq["DBInstanceDescription"] = v
	}

	if update {
		action := "ModifyDBInstanceDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyDBInstanceDescriptionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceDescriptionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("name")
	}

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		action := "TransformToPrePaid"
		transformToPrePaidReq := map[string]interface{}{
			"InstanceId": d.Id(),
		}

		transformToPrePaidReq["AutoPay"] = requests.NewBoolean(true)
		transformToPrePaidReq["Period"] = requests.NewInteger(d.Get("period").(int))

		if v, ok := d.GetOkExists("auto_renew"); ok {
			transformToPrePaidReq["AutoRenew"] = strconv.FormatBool(v.(bool))
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, transformToPrePaidReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, transformToPrePaidReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("instance_charge_type")
		d.SetPartial("period")
	}

	if !d.IsNewResource() && d.HasChange("security_ip_list") {
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

	if !d.IsNewResource() && (d.HasChange("account_password") || d.HasChange("kms_encrypted_password")) {
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

		err := ddsService.ResetAccountPassword(d, accountPassword, "shardingInstance")
		if err != nil {
			return WrapError(err)
		}
	}

	update = false
	modifyResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		modifyResourceGroupReq["ResourceGroupId"] = v
	}

	if update {
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		action := "ModifyResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	if d.HasChange("backup_time") || d.HasChange("backup_period") || d.HasChange("backup_retention_policy_on_cluster_deletion") || d.HasChange("snapshot_backup_type") || d.HasChange("backup_interval") {
		if err := ddsService.ModifyMongoDBBackupPolicy(d); err != nil {
			return WrapError(err)
		}

		d.SetPartial("backup_time")
		d.SetPartial("backup_period")
		d.SetPartial("backup_retention_policy_on_cluster_deletion")
		d.SetPartial("snapshot_backup_type")
		d.SetPartial("backup_interval")
	}

	if d.HasChange("tde_status") {
		action := "ModifyDBInstanceTDE"
		modifyDBInstanceTDEReq := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}

		if v, ok := d.GetOk("tde_status"); ok {
			modifyDBInstanceTDEReq["TDEStatus"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyDBInstanceTDEReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceTDEReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tde_status")
	}

	if d.HasChange("db_instance_release_protection") {
		action := "ModifyDBInstanceAttribute"

		modifyDBInstanceAttributeReq := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}

		if v, ok := d.GetOkExists("db_instance_release_protection"); ok {
			modifyDBInstanceAttributeReq["DBInstanceReleaseProtection"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyDBInstanceAttributeReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceAttributeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("db_instance_release_protection")
	}

	update = false
	modifyGlobalSecurityIPGroupRelationReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DBClusterId": d.Id(),
	}

	if d.HasChange("global_security_group_list") {
		update = true
	}
	if v, ok := d.GetOk("global_security_group_list"); ok {
		globalSecurityGroupList := v.(*schema.Set).List()

		modifyGlobalSecurityIPGroupRelationReq["GlobalSecurityGroupId"] = convertListToCommaSeparate(globalSecurityGroupList)
	}

	if update {
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		action := "ModifyGlobalSecurityIPGroupRelation"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, modifyGlobalSecurityIPGroupRelationReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyGlobalSecurityIPGroupRelationReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("global_security_group_list")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if !d.IsNewResource() && d.HasChange("mongo_list") {
		state, diff := d.GetChange("mongo_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeMongos, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}

		d.SetPartial("mongo_list")
	}

	if !d.IsNewResource() && d.HasChange("shard_list") {
		state, diff := d.GetChange("shard_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeShard, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}

		d.SetPartial("shard_list")
	}

	d.Partial(false)

	return resourceAliCloudMongoDBShardingInstanceRead(d, meta)
}

func resourceAliCloudMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	action := "DeleteDBInstance"
	var response map[string]interface{}

	var err error

	object, err := ddsService.DescribeMongoDBShardingInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(object["ChargeType"]) == "PrePaid" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudMongoDBShardingInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
