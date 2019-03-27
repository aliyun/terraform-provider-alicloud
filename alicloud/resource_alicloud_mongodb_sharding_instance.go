package alicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(WiredTiger), string(RocksDB)}),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(PrePaid), string(PostPaid)}),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: mongoDBPeriodDiffSuppressFunc,
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
				ValidateFunc: validateDBInstanceName,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(10, 1000),
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
		},
	}
}

func buildMongoDBShardingCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateShardingDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := dds.CreateCreateShardingDBInstanceRequest()

	request.RegionId = string(client.Region)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
	request.DBInstanceDescription = d.Get("name").(string)
	request.AccountPassword = d.Get("account_password").(string)
	request.ZoneId = d.Get("zone_id").(string)

	shardList, ok := d.GetOk("shard_list")
	if ok {
		replicaSets := []dds.CreateShardingDBInstanceReplicaSet{}
		for _, rew := range shardList.([]interface{}) {
			item := rew.(map[string]interface{})
			class := item["node_class"].(string)
			nodeStorage := item["node_storage"].(int)
			replicaSets = append(replicaSets, dds.CreateShardingDBInstanceReplicaSet{strconv.Itoa(nodeStorage), class})
		}
		request.ReplicaSet = &replicaSets
	}

	mongoList, ok := d.GetOk("mongo_list")
	if ok {
		mongos := []dds.CreateShardingDBInstanceMongos{}
		for _, rew := range mongoList.([]interface{}) {
			item := rew.(map[string]interface{})
			class := item["node_class"].(string)
			mongos = append(mongos, dds.CreateShardingDBInstanceMongos{class})
		}
		request.Mongos = &mongos
	}

	request.ConfigServer = &[]dds.CreateShardingDBInstanceConfigServer{{"20", "dds.cs.mid"}}

	request.NetworkType = string(Classic)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVswitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in multi the zone " + request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in the zone " + request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if ok && PayType(request.ChargeType) == PrePaid {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List()), COMMA_SEPARATED)
	}

	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceAlicloudMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBShardingCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateShardingDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*dds.CreateShardingDBInstanceResponse)
	addDebug(request.GetActionName(), response)

	d.SetId(response.DBInstanceId)

	if err := ddsService.WaitForMongoDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("name", instance.DBInstanceDescription)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("storage_engine", instance.StorageEngine)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.ChargeType)
	d.Set("vswitch_id", instance.VSwitchId)

	mongosList := []map[string]interface{}{}
	for _, item := range instance.MongosList.MongosAttribute {
		mongo := map[string]interface{}{
			"node_class":     item.NodeClass,
			"node_id":        item.NodeId,
			"port":           item.Port,
			"connect_string": item.ConnectSting,
		}
		mongosList = append(mongosList, mongo)
	}
	err = d.Set("mongo_list", mongosList)
	if err != nil {
		return WrapError(err)
	}

	shardList := []map[string]interface{}{}
	for _, item := range instance.ShardList.ShardAttribute {
		shard := map[string]interface{}{
			"node_id":      item.NodeId,
			"node_storage": item.NodeStorage,
			"node_class":   item.NodeClass,
		}
		shardList = append(shardList, shard)
	}
	err = d.Set("shard_list", shardList)
	if err != nil {
		return WrapError(err)
	}

	ips, err := ddsService.GetSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ip_list", ips)

	return nil
}

func resourceAlicloudMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	d.Partial(true)

	if d.HasChange("shard_list") {
		state, diff := d.GetChange("shard_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d.Id(), MongoDBShardingNodeShard, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("shard_list")
	}

	if d.HasChange("mongo_list") {
		state, diff := d.GetChange("mongo_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d.Id(), MongoDBShardingNodeMongos, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("mongo_list")
	}

	if d.HasChange("name") {
		request := dds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("name").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceDescription(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("name")
	}

	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := ddsService.ModifyMongoDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}
	d.Partial(false)
	return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			if ddsService.NotFoundMongoDBInstance(err) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ddsService.WaitForMongoDBInstance(d.Id(), Deleted, DefaultTimeout))
}
