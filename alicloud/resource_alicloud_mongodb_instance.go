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

func resourceAlicloudMongoDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBInstanceCreate,
		Read:   resourceAlicloudMongoDBInstanceRead,
		Update: resourceAlicloudMongoDBInstanceUpdate,
		Delete: resourceAlicloudMongoDBInstanceDelete,
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
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 5, 7}),
				Optional:     true,
				Computed:     true,
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
				Default:      PostPaid,
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
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
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
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
			},
			//Computed
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" || d.Get("engine_version").(string) < "4.0"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled"}, false),
				Optional:     true,
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
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"replica_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_port": {
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
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "VPC"}, false),
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

func resourceAlicloudMongoDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = string(client.Region)
	request["EngineVersion"] = Trim(d.Get("engine_version").(string))
	request["Engine"] = "MongoDB"
	request["DBInstanceStorage"] = d.Get("db_instance_storage")
	request["DBInstanceClass"] = Trim(d.Get("db_instance_class").(string))
	if v, ok := d.GetOk("name"); ok {
		request["DBInstanceDescription"] = v
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
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("storage_engine"); ok {
		request["StorageEngine"] = v
	}
	if v, ok := d.GetOk("replication_factor"); ok {
		request["ReplicationFactor"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
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
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ClientToken"] = buildClientToken(action)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBInstanceUpdate(d, meta)
}

func resourceAlicloudMongoDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	instance, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := instance["ReplicaSets"]; ok {
		replicaSets := v.(map[string]interface{})
		if replicaSetList, ok := replicaSets["ReplicaSet"]; ok {
			if replicaSet, ok := replicaSetList.([]interface{}); ok && len(replicaSet) > 0 {
				replicaSets := make([]map[string]interface{}, 0)
				for _, v := range replicaSet {
					item := v.(map[string]interface{})
					replicaSets = append(replicaSets, map[string]interface{}{
						"vswitch_id":            item["VSwitchId"],
						"connection_port":       item["ConnectionPort"],
						"replica_set_role":      item["ReplicaSetRole"],
						"connection_domain":     item["ConnectionDomain"],
						"vpc_cloud_instance_id": item["VPCCloudInstanceId"],
						"network_type":          item["NetworkType"],
						"vpc_id":                item["VPCId"],
					})
				}
				d.Set("replica_sets", replicaSets)
			}
		}
	}

	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_time", backupPolicy["PreferredBackupTime"])
	d.Set("backup_period", strings.Split(backupPolicy["PreferredBackupPeriod"].(string), ","))
	d.Set("retention_period", backupPolicy["BackupRetentionPeriod"])

	ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", ips)

	groupIp, err := ddsService.DescribeMongoDBSecurityGroupId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(groupIp) > 0 {
		if groupIpItem, ok := groupIp[0].(map[string]interface{}); ok {
			d.Set("security_group_id", groupIpItem["SecurityGroupId"])
		}
	}

	d.Set("name", instance["DBInstanceDescription"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("db_instance_class", instance["DBInstanceClass"])
	d.Set("db_instance_storage", instance["DBInstanceStorage"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("network_type", instance["NetworkType"])
	d.Set("instance_charge_type", instance["ChargeType"])
	if fmt.Sprint(instance["ChargeType"]) == "PrePaid" {
		period, err := computePeriodByUnit(instance["CreationTime"], instance["ExpireTime"], d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}
	d.Set("vswitch_id", instance["VSwitchId"])
	d.Set("storage_engine", instance["StorageEngine"])
	d.Set("maintain_start_time", instance["MaintainStartTime"])
	d.Set("maintain_end_time", instance["MaintainEndTime"])
	d.Set("replica_set_name", instance["ReplicaSetName"])
	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("vpc_id", instance["VPCId"])

	sslAction, err := ddsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction["SSLStatus"])

	d.Set("replication_factor", formatInt(instance["ReplicationFactor"]))
	if instance["ReplicationFactor"] != "" && instance["ReplicationFactor"] != "1" {
		tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("tde_Status", tdeInfo["TDEStatus"])
	}

	if v, ok := instance["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	return nil
}

func resourceAlicloudMongoDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

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
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("period")
	}

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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		var response map[string]interface{}
		action := "ModifyDBInstanceMaintainTime"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		request["MaintainStartTime"] = d.Get("maintain_start_time").(string)
		request["MaintainEndTime"] = d.Get("maintain_end_time").(string)

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
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
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
				if IsExpectedErrors(err, []string{"InstanceStatusInvalid", "OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
		return resourceAlicloudMongoDBInstanceRead(d, meta)
	}

	if d.HasChange("resource_group_id") {
		if v, ok := d.GetOk("resource_group_id"); ok {
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
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
	}

	if d.HasChange("ssl_action") {
		var response map[string]interface{}
		action := "ModifyDBInstanceSSL"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		request["SSLAction"] = d.Get("ssl_action")

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
		d.SetPartial("ssl_action")
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("db_instance_storage") ||
		d.HasChange("db_instance_class") ||
		d.HasChange("replication_factor") {

		var response map[string]interface{}
		action := "ModifyDBInstanceSpec"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["DBInstanceId"] = d.Id()
		if v, ok := d.GetOk("name"); ok {
			request["DBInstanceDescription"] = v
		}

		if d.Get("instance_charge_type").(string) == "PrePaid" {
			if v, ok := d.GetOk("order_type"); ok {
				request["OrderType"] = v.(string)
			}
		}
		request["DBInstanceClass"] = d.Get("db_instance_class").(string)
		request["DBInstanceStorage"] = strconv.Itoa(d.Get("db_instance_storage").(int))
		request["ReplicationFactor"] = strconv.Itoa(d.Get("replication_factor").(int))

		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
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

		d.SetPartial("db_instance_class")
		d.SetPartial("db_instance_storage")
		d.SetPartial("replication_factor")

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAlicloudMongoDBInstanceRead(d, meta)
}

func resourceAlicloudMongoDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// Pre paid instance can not be release.
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		log.Printf("[WARN] Cannot destroy resourceAlicloudMongoDBInstance. Terraform will remove this resource from the state file, however resources may remain.")
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
	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return nil
}
