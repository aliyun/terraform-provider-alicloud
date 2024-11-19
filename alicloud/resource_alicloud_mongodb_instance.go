package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMongoDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongoDBInstanceCreate,
		Read:   resourceAliCloudMongoDBInstanceRead,
		Update: resourceAliCloudMongoDBInstanceUpdate,
		Delete: resourceAliCloudMongoDBInstanceDelete,
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
				Required: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": {
				Type:     schema.TypeInt,
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
				ForceNew: true,
			},
			"hidden_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"cloud_disk_encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{1, 3, 5, 7}),
			},
			"readonly_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 5),
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
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable_backup_log": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"log_backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
			"ssl_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Open", "Close", "Update"}, false),
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Immediately", "MaintainTime"}, false),
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"tde_status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"enabled"}, false),
				ConflictsWith: []string{"encrypted", "cloud_disk_encryption_key"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine_version").(string) < "4.0"
				},
			},
			"encryptor_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"encrypted", "cloud_disk_encryption_key"},
			},
			"encryption_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"encrypted", "cloud_disk_encryption_key"},
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      parameterToHash,
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
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_sets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica_set_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudMongoDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = string(client.Region)
	request["ClientToken"] = buildClientToken(action)
	request["Engine"] = "MongoDB"
	request["EngineVersion"] = Trim(d.Get("engine_version").(string))
	request["DBInstanceClass"] = Trim(d.Get("db_instance_class").(string))
	request["DBInstanceStorage"] = d.Get("db_instance_storage")

	if v, ok := d.GetOk("storage_engine"); ok {
		request["StorageEngine"] = v
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v
	}

	if v, ok := d.GetOkExists("provisioned_iops"); ok {
		request["ProvisionedIops"] = v
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
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if request["ZoneId"] == nil {
			request["ZoneId"] = vsw["ZoneId"]
		} else if strings.Contains(request["ZoneId"].(string), MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(request["ZoneId"].(string), "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(vsw["ZoneId"].(string))[len(vsw["ZoneId"].(string))-1])) {
				return WrapError(Error("The specified vswitch " + request["VSwitchId"].(string) + " isn't in multi the zone " + request["ZoneId"].(string)))
			}
		} else if request["ZoneId"].(string) != vsw["ZoneId"] {
			return WrapError(Error("The specified vswitch " + request["VSwitchId"].(string) + " isn't in the zone " + request["ZoneId"].(string)))
		}
		if request["VpcId"] == nil {
			request["VpcId"] = vsw["VpcId"]
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

	if v, ok := d.GetOk("encrypted"); ok {
		request["Encrypted"] = v
	}

	if v, ok := d.GetOk("cloud_disk_encryption_key"); ok {
		request["EncryptionKey"] = v
	}

	if v, ok := d.GetOkExists("replication_factor"); ok {
		request["ReplicationFactor"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOkExists("readonly_replicas"); ok {
		request["ReadonlyReplicas"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = strconv.FormatBool(v.(bool))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudMongoDBInstanceUpdate(d, meta)
}

func resourceAliCloudMongoDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	object, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine_version", object["EngineVersion"])
	d.Set("db_instance_class", object["DBInstanceClass"])
	d.Set("db_instance_storage", object["DBInstanceStorage"])
	d.Set("storage_engine", object["StorageEngine"])
	d.Set("storage_type", convertMongoDBInstanceStorageTypeResponse(fmt.Sprint(object["StorageType"])))
	d.Set("provisioned_iops", formatInt(object["ProvisionedIops"]))
	d.Set("vpc_id", object["VPCId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("secondary_zone_id", object["SecondaryZoneId"])
	d.Set("hidden_zone_id", object["HiddenZoneId"])
	d.Set("network_type", object["NetworkType"])
	d.Set("name", object["DBInstanceDescription"])
	d.Set("instance_charge_type", object["ChargeType"])
	d.Set("encrypted", object["Encrypted"])
	d.Set("cloud_disk_encryption_key", object["EncryptionKey"])
	d.Set("replication_factor", formatInt(object["ReplicationFactor"]))
	d.Set("readonly_replicas", formatInt(object["ReadonlyReplicas"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("maintain_start_time", object["MaintainStartTime"])
	d.Set("maintain_end_time", object["MaintainEndTime"])
	d.Set("replica_set_name", object["ReplicaSetName"])

	if fmt.Sprint(object["ChargeType"]) == "PrePaid" {
		period, err := computePeriodByUnit(object["CreationTime"], object["ExpireTime"], d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}

		d.Set("period", period)
	}

	groupIp, err := ddsService.DescribeMongoDBSecurityGroupId(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if len(groupIp) > 0 {
		if groupIpItem, ok := groupIp[0].(map[string]interface{}); ok {
			d.Set("security_group_id", groupIpItem["SecurityGroupId"])
		}
	}

	securityIpList, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ip_list", securityIpList)

	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("backup_time", backupPolicy["PreferredBackupTime"])

	if backupPeriod, ok := backupPolicy["PreferredBackupPeriod"]; ok && fmt.Sprint(backupPeriod) != "" {
		d.Set("backup_period", strings.Split(backupPeriod.(string), ","))
	}

	d.Set("backup_retention_period", formatInt(backupPolicy["BackupRetentionPeriod"]))
	d.Set("backup_retention_policy_on_cluster_deletion", formatInt(backupPolicy["BackupRetentionPolicyOnClusterDeletion"]))
	d.Set("enable_backup_log", formatInt(backupPolicy["EnableBackupLog"]))
	d.Set("log_backup_retention_period", formatInt(backupPolicy["LogBackupRetentionPeriod"]))
	d.Set("snapshot_backup_type", backupPolicy["SnapshotBackupType"])
	d.Set("backup_interval", backupPolicy["BackupInterval"])
	d.Set("retention_period", formatInt(backupPolicy["BackupRetentionPeriod"]))

	if object["ReplicationFactor"] != "" && object["ReplicationFactor"] != "1" {
		tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
		if err != nil {
			return WrapError(err)
		}

		d.Set("tde_status", tdeInfo["TDEStatus"])
		d.Set("encryptor_name", tdeInfo["EncryptorName"])
		d.Set("encryption_key", tdeInfo["EncryptionKey"])
		d.Set("role_arn", tdeInfo["RoleARN"])
	}

	sslAction, err := ddsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		if !IsExpectedErrors(err, []string{"StorageTypeOrInstanceTypeNotSupported", "SingleNodeNotSupport"}) {
			return WrapError(err)
		}
	} else {
		d.Set("ssl_status", sslAction["SSLStatus"])
	}

	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	if err = ddsService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}

	if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
			replicaSetsMaps := make([]map[string]interface{}, 0)
			for _, replicaSets := range replicaSetsList.([]interface{}) {
				replicaSetsArg := replicaSets.(map[string]interface{})
				replicaSetsItemMap := make(map[string]interface{})

				if vpcId, ok := replicaSetsArg["VPCId"]; ok {
					replicaSetsItemMap["vpc_id"] = vpcId
				}

				if vswitchId, ok := replicaSetsArg["VSwitchId"]; ok {
					replicaSetsItemMap["vswitch_id"] = vswitchId
				}

				if networkType, ok := replicaSetsArg["NetworkType"]; ok {
					replicaSetsItemMap["network_type"] = networkType
				}

				if vpcCloudInstanceId, ok := replicaSetsArg["VPCCloudInstanceId"]; ok {
					replicaSetsItemMap["vpc_cloud_instance_id"] = vpcCloudInstanceId
				}

				if replicaSetRole, ok := replicaSetsArg["ReplicaSetRole"]; ok {
					replicaSetsItemMap["replica_set_role"] = replicaSetRole
				}

				if connectionDomain, ok := replicaSetsArg["ConnectionDomain"]; ok {
					replicaSetsItemMap["connection_domain"] = connectionDomain
				}

				if connectionPort, ok := replicaSetsArg["ConnectionPort"]; ok {
					replicaSetsItemMap["connection_port"] = connectionPort
				}

				replicaSetsMaps = append(replicaSetsMaps, replicaSetsItemMap)
			}

			d.Set("replica_sets", replicaSetsMaps)
		}
	}

	return nil
}

func resourceAliCloudMongoDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	upgradeDBInstanceEngineVersionReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
	}
	upgradeDBInstanceEngineVersionReq["EngineVersion"] = d.Get("engine_version").(string)

	if update {
		action := "UpgradeDBInstanceEngineVersion"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, upgradeDBInstanceEngineVersionReq, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("engine_version")
	}

	update = false
	modifyDBInstanceSpecReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("db_instance_class") {
		update = true
	}
	modifyDBInstanceSpecReq["DBInstanceClass"] = d.Get("db_instance_class").(string)

	if !d.IsNewResource() && d.HasChange("db_instance_storage") {
		update = true
	}
	modifyDBInstanceSpecReq["DBInstanceStorage"] = strconv.Itoa(d.Get("db_instance_storage").(int))

	if !d.IsNewResource() && d.HasChange("replication_factor") {
		update = true
	}
	if v, ok := d.GetOkExists("replication_factor"); ok {
		modifyDBInstanceSpecReq["ReplicationFactor"] = strconv.Itoa(v.(int))
	}

	if !d.IsNewResource() && d.HasChange("readonly_replicas") {
		update = true
	}
	if v, ok := d.GetOkExists("readonly_replicas"); ok {
		modifyDBInstanceSpecReq["ReadonlyReplicas"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("effective_time"); ok {
		modifyDBInstanceSpecReq["EffectiveTime"] = v
	}

	if v, ok := d.GetOk("order_type"); ok {
		modifyDBInstanceSpecReq["OrderType"] = v.(string)
	}

	if update {
		action := "ModifyDBInstanceSpec"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceSpecReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Task.Conflict", "OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceSpecReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"order_wait_for_produce"}, []string{"all_completed"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceOrderStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		stateConf = BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "NodeCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("db_instance_class")
		d.SetPartial("db_instance_storage")
		d.SetPartial("replication_factor")
		d.SetPartial("readonly_replicas")
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
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceDiskTypeReq, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("storage_type")
		d.SetPartial("provisioned_iops")
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
		action := "ModifySecurityGroupConfiguration"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifySecurityGroupConfigurationReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"InstanceStatusInvalid", "OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceDescriptionReq, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		action := "ModifyResourceGroup"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyResourceGroupReq, &runtime)
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

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		action := "TransformToPrePaid"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		transformToPrePaidReq := map[string]interface{}{
			"InstanceId": d.Id(),
		}

		transformToPrePaidReq["AutoPay"] = requests.NewBoolean(true)
		transformToPrePaidReq["Period"] = requests.NewInteger(d.Get("period").(int))

		if v, ok := d.GetOkExists("auto_renew"); ok {
			transformToPrePaidReq["AutoRenew"] = strconv.FormatBool(v.(bool))
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, transformToPrePaidReq, &runtime)
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
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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

		err := ddsService.ResetAccountPassword(d, accountPassword, "instance")
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("backup_time") || d.HasChange("backup_period") || d.HasChange("backup_retention_period") || d.HasChange("backup_retention_policy_on_cluster_deletion") || d.HasChange("enable_backup_log") || d.HasChange("log_backup_retention_period") || d.HasChange("snapshot_backup_type") || d.HasChange("backup_interval") {
		if err := ddsService.ModifyMongoDBBackupPolicy(d); err != nil {
			return WrapError(err)
		}

		d.SetPartial("backup_time")
		d.SetPartial("backup_period")
		d.SetPartial("backup_retention_period")
		d.SetPartial("backup_retention_policy_on_cluster_deletion")
		d.SetPartial("enable_backup_log")
		d.SetPartial("log_backup_retention_period")
		d.SetPartial("snapshot_backup_type")
		d.SetPartial("backup_interval")
	}

	if d.HasChange("ssl_action") {
		action := "ModifyDBInstanceSSL"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		modifyDBInstanceSSLReq := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}

		if v, ok := d.GetOk("ssl_action"); ok {
			modifyDBInstanceSSLReq["SSLAction"] = v
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceSSLReq, &runtime)
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

		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("ssl_action")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		action := "ModifyDBInstanceMaintainTime"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		modifyDBInstanceMaintainTimeReq := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}

		if v, ok := d.GetOk("maintain_start_time"); ok {
			modifyDBInstanceMaintainTimeReq["MaintainStartTime"] = v
		}

		if v, ok := d.GetOk("maintain_end_time"); ok {
			modifyDBInstanceMaintainTimeReq["MaintainEndTime"] = v
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceMaintainTimeReq, &runtime)
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

		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	if d.HasChange("tde_status") || d.HasChange("encryptor_name") || d.HasChange("encryption_key") {
		action := "ModifyDBInstanceTDE"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}

		modifyDBInstanceTDEReq := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}

		if d.HasChange("tde_status") {

			if v, ok := d.GetOk("tde_status"); ok {
				modifyDBInstanceTDEReq["TDEStatus"] = v
			}
		}

		if d.HasChange("encryptor_name") {

			if v, ok := d.GetOk("encryptor_name"); ok {
				modifyDBInstanceTDEReq["EncryptorName"] = v
			}
		}

		if d.HasChange("encryption_key") {

			if v, ok := d.GetOk("encryption_key"); ok {
				modifyDBInstanceTDEReq["EncryptionKey"] = v
			}
		}

		if d.HasChange("role_arn") {

			if v, ok := d.GetOk("role_arn"); ok {
				modifyDBInstanceTDEReq["RoleARN"] = v
			}
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, modifyDBInstanceTDEReq, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tde_status")
		d.SetPartial("encryptor_name")
		d.SetPartial("encryption_key")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("parameters") {
		if err := ddsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)

	return resourceAliCloudMongoDBInstanceRead(d, meta)
}

func resourceAliCloudMongoDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	action := "DeleteDBInstance"
	var response map[string]interface{}

	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}

	object, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(object["ChargeType"]) == "PrePaid" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudMongoDBInstance prepay type. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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

	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertMongoDBInstanceStorageTypeResponse(source string) string {
	switch source {
	case "cloud_essd":
		return "cloud_essd1"
	}

	return source
}
