package alicloud

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAlikafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlikafkaInstanceCreate,
		Read:   resourceAliCloudAlikafkaInstanceRead,
		Update: resourceAliCloudAlikafkaInstanceUpdate,
		Delete: resourceAliCloudAlikafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"deploy_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{4, 5}),
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"alikafka", "alikafka_serverless", "alikafka_confluent"}, false),
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"topic_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, _ := strconv.Atoi(old)
					partitionNum := d.Get("partition_num").(int)
					if o > 0 {
						return o-1000 == partitionNum
					}
					return false
				},
				Deprecated: "Attribute `topic_quota` has been deprecated since 1.194.0 and it will be removed in the next future. Using new attribute `partition_num` instead.",
			},
			"io_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"io_max_spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringLenBetween(3, 64),
			},
			"paid_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "normal",
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(int) == 5
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: alikafkaInstanceConfigDiffSuppressFunc,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"enable_auto_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_topic": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disable"}, false),
			},
			"default_topic_partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"selected_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"serverless_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserved_publish_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"reserved_subscribe_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"confluent_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kafka_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"kafka_storage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"kafka_replica": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"kafka_rest_proxy_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"kafka_rest_proxy_replica": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"zookeeper_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"zookeeper_storage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"zookeeper_replica": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"control_center_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"control_center_storage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"control_center_replica": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"schema_registry_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"schema_registry_replica": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"connect_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"connect_replica": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ksql_cu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ksql_storage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ksql_replica": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"end_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sasl_domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_num_of_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_partition_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}
	var err error

	// 1. Create order
	var createOrderAction string
	var instanceType string
	createOrderResponse := make(map[string]interface{})
	createOrderReq := make(map[string]interface{})
	createOrderReq["RegionId"] = client.RegionId
	createOrderReq["DeployType"] = d.Get("deploy_type")

	createOrderReq["DiskType"] = d.Get("disk_type")

	createOrderReq["DiskSize"] = d.Get("disk_size")

	if v, ok := d.GetOk("partition_num"); ok {
		createOrderReq["PartitionNum"] = v
	} else if v, ok := d.GetOk("topic_quota"); ok {
		createOrderReq["TopicQuota"] = v
	}

	if v, ok := d.GetOk("instance_type"); ok {
		instanceType = fmt.Sprint(v)
	}

	if v, ok := d.GetOk("io_max"); ok {
		createOrderReq["IoMax"] = v
	}

	if v, ok := d.GetOk("io_max_spec"); ok {
		createOrderReq["IoMaxSpec"] = v
	}

	if v, ok := d.GetOk("spec_type"); ok {
		createOrderReq["SpecType"] = v
	}

	if v, ok := d.GetOkExists("eip_max"); ok {
		createOrderReq["EipMax"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		createOrderReq["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("serverless_config"); ok {
		serverlessConfigMap := map[string]interface{}{}
		for _, serverlessConfigList := range v.([]interface{}) {
			serverlessConfigArg := serverlessConfigList.(map[string]interface{})

			if reservedPublishCapacity, ok := serverlessConfigArg["reserved_publish_capacity"]; ok {
				serverlessConfigMap["ReservedPublishCapacity"] = reservedPublishCapacity
			}

			if reservedSubscribeCapacity, ok := serverlessConfigArg["reserved_subscribe_capacity"]; ok {
				serverlessConfigMap["ReservedSubscribeCapacity"] = reservedSubscribeCapacity
			}
		}

		serverlessConfigJson, err := convertMaptoJsonString(serverlessConfigMap)
		if err != nil {
			return WrapError(err)
		}

		createOrderReq["ServerlessConfig"] = serverlessConfigJson
	}

	if v, ok := d.GetOk("confluent_config"); ok {
		confluentConfigMap := map[string]interface{}{}
		for _, confluentConfigList := range v.([]interface{}) {
			confluentConfigArg := confluentConfigList.(map[string]interface{})

			if kafkaCU, ok := confluentConfigArg["kafka_cu"]; ok {
				confluentConfigMap["KafkaCU"] = kafkaCU
			}

			if kafkaStorage, ok := confluentConfigArg["kafka_storage"]; ok {
				confluentConfigMap["KafkaStorage"] = kafkaStorage
			}

			if kafkaReplica, ok := confluentConfigArg["kafka_replica"]; ok {
				confluentConfigMap["KafkaReplica"] = kafkaReplica
			}

			if kafkaRestProxyCU, ok := confluentConfigArg["kafka_rest_proxy_cu"]; ok {
				confluentConfigMap["KafkaRestProxyCU"] = kafkaRestProxyCU
			}

			if kafkaRestProxyReplica, ok := confluentConfigArg["kafka_rest_proxy_replica"]; ok {
				confluentConfigMap["KafkaRestProxyReplica"] = kafkaRestProxyReplica
			}

			if zookeeperCU, ok := confluentConfigArg["zookeeper_cu"]; ok {
				confluentConfigMap["ZooKeeperCU"] = zookeeperCU
			}

			if zookeeperStorage, ok := confluentConfigArg["zookeeper_storage"]; ok {
				confluentConfigMap["ZooKeeperStorage"] = zookeeperStorage
			}

			if zookeeperReplica, ok := confluentConfigArg["zookeeper_replica"]; ok {
				confluentConfigMap["ZooKeeperReplica"] = zookeeperReplica
			}

			if controlCenterCU, ok := confluentConfigArg["control_center_cu"]; ok {
				confluentConfigMap["ControlCenterCU"] = controlCenterCU
			}

			if controlCenterStorage, ok := confluentConfigArg["control_center_storage"]; ok {
				confluentConfigMap["ControlCenterStorage"] = controlCenterStorage
			}

			if controlCenterReplica, ok := confluentConfigArg["control_center_replica"]; ok {
				confluentConfigMap["ControlCenterReplica"] = controlCenterReplica
			}

			if schemaRegistryCU, ok := confluentConfigArg["schema_registry_cu"]; ok {
				confluentConfigMap["SchemaRegistryCU"] = schemaRegistryCU
			}

			if schemaRegistryReplica, ok := confluentConfigArg["schema_registry_replica"]; ok {
				confluentConfigMap["SchemaRegistryReplica"] = schemaRegistryReplica
			}

			if connectCU, ok := confluentConfigArg["connect_cu"]; ok {
				confluentConfigMap["ConnectCU"] = connectCU
			}

			if connectReplica, ok := confluentConfigArg["connect_replica"]; ok {
				confluentConfigMap["ConnectReplica"] = connectReplica
			}

			if ksqlCU, ok := confluentConfigArg["ksql_cu"]; ok {
				confluentConfigMap["KsqlCU"] = ksqlCU
			}

			if ksqlStorage, ok := confluentConfigArg["ksql_storage"]; ok {
				confluentConfigMap["KsqlStorage"] = ksqlStorage
			}

			if ksqlReplica, ok := confluentConfigArg["ksql_replica"]; ok {
				confluentConfigMap["KsqlReplica"] = ksqlReplica
			}
		}

		confluentConfigJson, err := convertMaptoJsonString(confluentConfigMap)
		if err != nil {
			return WrapError(err)
		}

		createOrderReq["ConfluentConfig"] = confluentConfigJson
	}

	if v, ok := d.GetOk("paid_type"); ok {
		switch v.(string) {
		case "PostPaid":
			createOrderAction = "CreatePostPayOrder"

			createOrderReq["PaidType"] = convertAliKafkaPostPaidTypeRequest(instanceType)
		case "PrePaid":
			createOrderAction = "CreatePrePayOrder"

			createOrderReq["PaidType"] = convertAliKafkaPrePaidTypeRequest(instanceType)
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		createOrderResponse, err = client.RpcPost("alikafka", "2019-09-16", createOrderAction, nil, createOrderReq, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL", "ONS_SYSTEM_ERROR"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(createOrderAction, createOrderResponse, createOrderReq)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", createOrderAction, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(createOrderResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", createOrderAction, createOrderResponse))
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAliKafkaInstanceByOrderId(fmt.Sprint(createOrderResponse["OrderId"]), 60)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(fmt.Sprint(alikafkaInstanceVO["InstanceId"]))

	// 2. Start instance
	startInstanceAction := "StartInstance"
	startInstanceResponse := make(map[string]interface{})
	startInstanceReq := make(map[string]interface{})
	startInstanceReq["RegionId"] = client.RegionId
	startInstanceReq["InstanceId"] = alikafkaInstanceVO["InstanceId"]
	startInstanceReq["VSwitchId"] = d.Get("vswitch_id")

	if v, ok := d.GetOk("vpc_id"); ok {
		startInstanceReq["VpcId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		startInstanceReq["ZoneId"] = v
	}

	if startInstanceReq["VpcId"] == nil {
		vsw, err := vpcService.DescribeVswitch(startInstanceReq["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}

		if v, ok := startInstanceReq["VpcId"].(string); !ok || v == "" {
			startInstanceReq["VpcId"] = vsw["VpcId"]
		}
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		startInstanceReq["VSwitchIds"] = v
	}

	if _, ok := d.GetOkExists("eip_max"); ok {
		startInstanceReq["DeployModule"] = "eip"
		startInstanceReq["IsEipInner"] = true
	}

	if v, ok := d.GetOk("name"); ok {
		startInstanceReq["Name"] = v
	}

	if v, ok := d.GetOk("security_group"); ok {
		startInstanceReq["SecurityGroup"] = v
	}

	if v, ok := d.GetOk("service_version"); ok {
		startInstanceReq["ServiceVersion"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		startInstanceReq["Config"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		startInstanceReq["KMSKeyId"] = v
	}

	if v, ok := d.GetOk("password"); ok {
		startInstanceReq["Password"] = v
	}

	if v, ok := d.GetOk("selected_zones"); ok {
		startInstanceReq["SelectedZones"] = formatSelectedZonesReq(v.([]interface{}))
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		startInstanceResponse, err = client.RpcPost("alikafka", "2019-09-16", startInstanceAction, nil, startInstanceReq, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(startInstanceAction, startInstanceResponse, startInstanceReq)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceAction, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(startInstanceResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", startInstanceAction, startInstanceResponse))
	}

	// 3. wait until running
	stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlikafkaInstanceUpdate(d, meta)
}

func resourceAliCloudAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAliKafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikakfa_instance alikafkaService.DescribeAliKafkaInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("disk_type", object["DiskType"])
	d.Set("disk_size", object["DiskSize"])
	d.Set("deploy_type", object["DeployType"])
	d.Set("instance_type", convertAliKafkaInstanceTypeResponse(fmt.Sprint(object["PaidType"])))
	d.Set("io_max", object["IoMax"])
	d.Set("io_max_spec", object["IoMaxSpec"])
	d.Set("eip_max", object["EipMax"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("paid_type", convertAliKafkaPaidTypeResponse(fmt.Sprint(object["PaidType"])))
	d.Set("spec_type", object["SpecType"])
	d.Set("security_group", object["SecurityGroup"])
	d.Set("end_point", object["EndPoint"])
	d.Set("ssl_endpoint", object["SslEndPoint"])
	d.Set("domain_endpoint", object["DomainEndpoint"])
	d.Set("ssl_domain_endpoint", object["SslDomainEndpoint"])
	d.Set("sasl_domain_endpoint", object["SaslDomainEndpoint"])
	d.Set("status", object["ServiceStatus"])
	// object.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO[0].Current2OpenSourceVersion can guaranteed not to be null
	d.Set("service_version", object["UpgradeServiceDetailInfo"].(map[string]interface{})["Current2OpenSourceVersion"])
	d.Set("kms_key_id", object["KmsKeyId"])
	d.Set("enable_auto_group", object["AutoCreateGroupEnable"])
	d.Set("enable_auto_topic", convertAliKafkaAutoCreateTopicEnableResponse(object["AutoCreateTopicEnable"]))
	d.Set("default_topic_partition_num", formatInt(object["DefaultPartitionNum"]))

	if allConfig, ok := object["AllConfig"]; ok && fmt.Sprint(allConfig) != "" {
		allConfigMap, err := convertJsonStringToMap(fmt.Sprint(allConfig))
		if err != nil {
			return WrapError(err)
		}

		configMap := make(map[string]interface{})

		for k, v := range allConfigMap {
			if k == "enable.tiered" || k == "cloud.maxTieredStoreSpace" || k == "enable.compact" {
				continue
			}

			configMap[k] = v
		}

		configMapJson, err := convertMaptoJsonString(configMap)
		if err != nil {
			return WrapError(err)
		}

		d.Set("config", configMapJson)
	}

	if vSwitchIds, ok := object["VSwitchIds"]; ok {
		vSwitchIdsArg := vSwitchIds.(map[string]interface{})

		if vSwitchIdsList, ok := vSwitchIdsArg["VSwitchIds"]; ok {
			d.Set("vswitch_ids", vSwitchIdsList)
		}
	}

	serverlessConfigMaps := make([]map[string]interface{}, 0)
	serverlessConfigMap := map[string]interface{}{}
	serverlessConfigMap["reserved_publish_capacity"] = object["ReservedPublishCapacity"]
	serverlessConfigMap["reserved_subscribe_capacity"] = object["ReservedSubscribeCapacity"]
	serverlessConfigMaps = append(serverlessConfigMaps, serverlessConfigMap)
	d.Set("serverless_config", serverlessConfigMaps)

	if confluentConfig, ok := object["ConfluentConfig"]; ok {
		confluentConfigMaps := make([]map[string]interface{}, 0)
		confluentConfigMap := map[string]interface{}{}
		confluentConfigArg := confluentConfig.(map[string]interface{})

		if kafkaCU, ok := confluentConfigArg["KafkaCU"]; ok {
			confluentConfigMap["kafka_cu"] = kafkaCU
		}

		if kafkaStorage, ok := confluentConfigArg["KafkaStorage"]; ok {
			confluentConfigMap["kafka_storage"] = kafkaStorage
		}

		if kafkaReplica, ok := confluentConfigArg["KafkaReplica"]; ok {
			confluentConfigMap["kafka_replica"] = kafkaReplica
		}

		if kafkaRestProxyCU, ok := confluentConfigArg["KafkaRestProxyCU"]; ok {
			confluentConfigMap["kafka_rest_proxy_cu"] = kafkaRestProxyCU
		}

		if kafkaRestProxyReplica, ok := confluentConfigArg["KafkaRestProxyReplica"]; ok {
			confluentConfigMap["kafka_rest_proxy_replica"] = kafkaRestProxyReplica
		}

		if zookeeperCU, ok := confluentConfigArg["ZooKeeperCU"]; ok {
			confluentConfigMap["zookeeper_cu"] = zookeeperCU
		}

		if zookeeperStorage, ok := confluentConfigArg["ZooKeeperStorage"]; ok {
			confluentConfigMap["zookeeper_storage"] = zookeeperStorage
		}

		if zookeeperReplica, ok := confluentConfigArg["ZooKeeperReplica"]; ok {
			confluentConfigMap["zookeeper_replica"] = zookeeperReplica
		}

		if controlCenterCU, ok := confluentConfigArg["ControlCenterCU"]; ok {
			confluentConfigMap["control_center_cu"] = controlCenterCU
		}

		if controlCenterStorage, ok := confluentConfigArg["ControlCenterStorage"]; ok {
			confluentConfigMap["control_center_storage"] = controlCenterStorage
		}

		if controlCenterReplica, ok := confluentConfigArg["ControlCenterReplica"]; ok {
			confluentConfigMap["control_center_replica"] = controlCenterReplica
		}

		if schemaRegistryCU, ok := confluentConfigArg["SchemaRegistryCU"]; ok {
			confluentConfigMap["schema_registry_cu"] = schemaRegistryCU
		}

		if schemaRegistryReplica, ok := confluentConfigArg["SchemaRegistryReplica"]; ok {
			confluentConfigMap["schema_registry_replica"] = schemaRegistryReplica
		}

		if connectCU, ok := confluentConfigArg["ConnectCU"]; ok {
			confluentConfigMap["connect_cu"] = connectCU
		}

		if connectReplica, ok := confluentConfigArg["ConnectReplica"]; ok {
			confluentConfigMap["connect_replica"] = connectReplica
		}

		if ksqlCU, ok := confluentConfigArg["KsqlCU"]; ok {
			confluentConfigMap["ksql_cu"] = ksqlCU
		}

		if ksqlStorage, ok := confluentConfigArg["KsqlStorage"]; ok {
			confluentConfigMap["ksql_storage"] = ksqlStorage
		}

		if ksqlReplica, ok := confluentConfigArg["KsqlReplica"]; ok {
			confluentConfigMap["ksql_replica"] = ksqlReplica
		}

		confluentConfigMaps = append(confluentConfigMaps, confluentConfigMap)

		d.Set("confluent_config", confluentConfigMaps)
	}

	if fmt.Sprint(convertAliKafkaInstanceTypeResponse(fmt.Sprint(object["PaidType"]))) == "alikafka" {
		quota, err := alikafkaService.GetQuotaTip(d.Id())
		if err != nil {
			return WrapError(err)
		}

		d.Set("topic_quota", quota["TopicQuota"])
		d.Set("partition_num", quota["PartitionNumOfBuy"])
		d.Set("topic_num_of_buy", quota["TopicNumOfBuy"])
		d.Set("topic_used", quota["TopicUsed"])
		d.Set("topic_left", quota["TopicLeft"])
		d.Set("partition_used", quota["PartitionUsed"])
		d.Set("partition_left", quota["PartitionLeft"])
		d.Set("group_used", quota["GroupUsed"])
		d.Set("group_left", quota["GroupLeft"])
		d.Set("is_partition_buy", quota["IsPartitionBuy"])
	}

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAliCloudAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if err := alikafkaService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}

	// Process change instance name.
	if !d.IsNewResource() && d.HasChange("name") {
		action := "ModifyInstanceName"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("name"); ok {
			request["InstanceName"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("name")
	}

	// Process paid type change, note only support change from post to pre pay.
	if !d.IsNewResource() && d.HasChange("paid_type") {
		o, n := d.GetChange("paid_type")
		oldPaidType := o.(string)
		newPaidType := n.(string)
		oldPaidTypeInt := 1
		newPaidTypeInt := 1
		if oldPaidType == string(PrePaid) {
			oldPaidTypeInt = 0
		}
		if newPaidType == string(PrePaid) {
			newPaidTypeInt = 0
		}
		if oldPaidTypeInt == 1 && newPaidTypeInt == 0 {
			action := "ConvertPostPayOrder"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": d.Id(),
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

			stateConf := BuildStateConf([]string{}, []string{strconv.Itoa(newPaidTypeInt)}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "PaidType", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		} else {
			return WrapError(errors.New("paid type only support change from post pay to pre pay"))
		}

		d.SetPartial("paid_type")
	}

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}
	// updating topic_quota only by updating partition_num
	if !d.IsNewResource() && d.HasChange("partition_num") {
		update = true
	}
	request["PartitionNum"] = d.Get("partition_num")

	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")

	if !d.IsNewResource() && d.HasChange("io_max") {
		update = true

		if v, ok := d.GetOk("io_max"); ok {
			request["IoMax"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("io_max_spec") {
		update = true

		if v, ok := d.GetOk("io_max_spec"); ok {
			request["IoMaxSpec"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("spec_type") {
		update = true
	}
	request["SpecType"] = d.Get("spec_type")

	if !d.IsNewResource() && d.HasChange("deploy_type") {
		update = true
	}
	if d.Get("deploy_type").(int) == 4 {
		request["EipModel"] = true
	} else {
		request["EipModel"] = false
	}
	if !d.IsNewResource() && d.HasChange("eip_max") {
		update = true
	}
	request["EipMax"] = d.Get("eip_max").(int)

	if !d.IsNewResource() && d.HasChange("serverless_config") {
		update = true

		if v, ok := d.GetOk("serverless_config"); ok {
			serverlessConfigMap := map[string]interface{}{}
			for _, serverlessConfigList := range v.([]interface{}) {
				serverlessConfigArg := serverlessConfigList.(map[string]interface{})

				if reservedPublishCapacity, ok := serverlessConfigArg["reserved_publish_capacity"]; ok {
					serverlessConfigMap["ReservedPublishCapacity"] = reservedPublishCapacity
				}

				if reservedSubscribeCapacity, ok := serverlessConfigArg["reserved_subscribe_capacity"]; ok {
					serverlessConfigMap["ReservedSubscribeCapacity"] = reservedSubscribeCapacity
				}
			}

			serverlessConfigJson, err := convertMaptoJsonString(serverlessConfigMap)
			if err != nil {
				return WrapError(err)
			}

			request["ServerlessConfig"] = serverlessConfigJson
		}
	}

	if !d.IsNewResource() && d.HasChange("confluent_config") {
		update = true

		if v, ok := d.GetOk("confluent_config"); ok {
			confluentConfigMap := map[string]interface{}{}
			for _, confluentConfigList := range v.([]interface{}) {
				confluentConfigArg := confluentConfigList.(map[string]interface{})

				if kafkaCU, ok := confluentConfigArg["kafka_cu"]; ok {
					confluentConfigMap["KafkaCU"] = kafkaCU
				}

				if kafkaStorage, ok := confluentConfigArg["kafka_storage"]; ok {
					confluentConfigMap["KafkaStorage"] = kafkaStorage
				}

				if kafkaReplica, ok := confluentConfigArg["kafka_replica"]; ok {
					confluentConfigMap["KafkaReplica"] = kafkaReplica
				}

				if kafkaRestProxyCU, ok := confluentConfigArg["kafka_rest_proxy_cu"]; ok {
					confluentConfigMap["KafkaRestProxyCU"] = kafkaRestProxyCU
				}

				if kafkaRestProxyReplica, ok := confluentConfigArg["kafka_rest_proxy_replica"]; ok {
					confluentConfigMap["KafkaRestProxyReplica"] = kafkaRestProxyReplica
				}

				if zookeeperCU, ok := confluentConfigArg["zookeeper_cu"]; ok {
					confluentConfigMap["ZooKeeperCU"] = zookeeperCU
				}

				if zookeeperStorage, ok := confluentConfigArg["zookeeper_storage"]; ok {
					confluentConfigMap["ZooKeeperStorage"] = zookeeperStorage
				}

				if zookeeperReplica, ok := confluentConfigArg["zookeeper_replica"]; ok {
					confluentConfigMap["ZooKeeperReplica"] = zookeeperReplica
				}

				if controlCenterCU, ok := confluentConfigArg["control_center_cu"]; ok {
					confluentConfigMap["ControlCenterCU"] = controlCenterCU
				}

				if controlCenterStorage, ok := confluentConfigArg["control_center_storage"]; ok {
					confluentConfigMap["ControlCenterStorage"] = controlCenterStorage
				}

				if controlCenterReplica, ok := confluentConfigArg["control_center_replica"]; ok {
					confluentConfigMap["ControlCenterReplica"] = controlCenterReplica
				}

				if schemaRegistryCU, ok := confluentConfigArg["schema_registry_cu"]; ok {
					confluentConfigMap["SchemaRegistryCU"] = schemaRegistryCU
				}

				if schemaRegistryReplica, ok := confluentConfigArg["schema_registry_replica"]; ok {
					confluentConfigMap["SchemaRegistryReplica"] = schemaRegistryReplica
				}

				if connectCU, ok := confluentConfigArg["connect_cu"]; ok {
					confluentConfigMap["ConnectCU"] = connectCU
				}

				if connectReplica, ok := confluentConfigArg["connect_replica"]; ok {
					confluentConfigMap["ConnectReplica"] = connectReplica
				}

				if ksqlCU, ok := confluentConfigArg["ksql_cu"]; ok {
					confluentConfigMap["KsqlCU"] = ksqlCU
				}

				if ksqlStorage, ok := confluentConfigArg["ksql_storage"]; ok {
					confluentConfigMap["KsqlStorage"] = ksqlStorage
				}

				if ksqlReplica, ok := confluentConfigArg["ksql_replica"]; ok {
					confluentConfigMap["KsqlReplica"] = ksqlReplica
				}
			}

			confluentConfigJson, err := convertMaptoJsonString(confluentConfigMap)
			if err != nil {
				return WrapError(err)
			}

			request["ConfluentConfig"] = confluentConfigJson
		}
	}

	if update {
		action := "UpgradePostPayOrder"

		if d.Get("paid_type").(string) == string(PrePaid) {
			request["PaidType"] = convertAliKafkaPrePaidTypeRequest(d.Get("instance_type").(string))
			action = "UpgradePrePayOrder"
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL", "ScheduledTask.AlreadyHasSameTaskType"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if d.HasChange("disk_size") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("disk_size"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "DiskSize", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if d.HasChange("io_max") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("io_max"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "IoMax", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if d.HasChange("io_max_spec") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("io_max_spec"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "IoMaxSpec", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if d.HasChange("spec_type") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("spec_type"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "SpecType", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if d.HasChange("eip_max") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("eip_max"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "EipMax", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("partition_num")
		d.SetPartial("disk_size")
		d.SetPartial("io_max")
		d.SetPartial("io_max_spec")
		d.SetPartial("spec_type")
		d.SetPartial("eip_max")
		d.SetPartial("serverless_config")
		d.SetPartial("confluent_config")
	}

	if !d.IsNewResource() && d.HasChange("service_version") {
		action := "UpgradeInstanceVersion"
		request := map[string]interface{}{
			"InstanceId": d.Id(),
			"RegionId":   client.RegionId,
		}

		if v, ok := d.GetOk("service_version"); ok {
			request["TargetVersion"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				// means no need to update version
				if IsExpectedErrors(err, []string{"ONS_INIT_ENV_ERROR"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		// upgrade service may be last a long time
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("service_version")
	}

	if !d.IsNewResource() && d.HasChange("config") {
		action := "UpdateInstanceConfig"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("config"); ok {
			request["Config"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("config")
	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ResourceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, changeResourceGroupReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	object, err := alikafkaService.DescribeAliKafkaInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(convertAliKafkaInstanceTypeResponse(fmt.Sprint(object["PaidType"]))) == "alikafka" {
		update = false
		enableAutoGroupCreationReq := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if d.HasChange("enable_auto_group") {
			update = true

			if v, ok := d.GetOkExists("enable_auto_group"); ok {
				enableAutoGroupCreationReq["Enable"] = v
			}
		}

		if update {
			action := "EnableAutoGroupCreation"

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, enableAutoGroupCreationReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, enableAutoGroupCreationReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}

			d.SetPartial("enable_auto_group")
		}
	}

	if fmt.Sprint(convertAliKafkaInstanceTypeResponse(fmt.Sprint(object["PaidType"]))) != "alikafka_confluent" {
		update = false
		enableAutoTopicCreationReq := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if d.HasChange("enable_auto_topic") {
			update = true
		}
		if v, ok := d.GetOk("enable_auto_topic"); ok {
			enableAutoTopicCreationReq["Operate"] = v
		}

		if update {
			action := "EnableAutoTopicCreation"

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, enableAutoTopicCreationReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, enableAutoTopicCreationReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}

			d.SetPartial("enable_auto_topic")
		}

		update = false
		updateTopicPartitionNumReq := map[string]interface{}{
			"RegionId":        client.RegionId,
			"Operate":         "updatePartition",
			"UpdatePartition": true,
			"InstanceId":      d.Id(),
		}

		defaultTopicPartitionNum, ok := d.GetOkExists("default_topic_partition_num")
		if ok && fmt.Sprint(object["DefaultPartitionNum"]) != fmt.Sprint(defaultTopicPartitionNum) {
			update = true
			updateTopicPartitionNumReq["PartitionNum"] = defaultTopicPartitionNum
		}

		if update {
			action := "EnableAutoTopicCreation"

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, updateTopicPartitionNumReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action+" updateTopicPartitionNum", response, updateTopicPartitionNumReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}

			d.SetPartial("default_topic_partition_num")
		}
	}

	d.Partial(false)

	return resourceAliCloudAlikafkaInstanceRead(d, meta)
}

func resourceAliCloudAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	action := "ReleaseInstance"
	var err error
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId":          d.Id(),
		"RegionId":            client.RegionId,
		"ForceDeleteInstance": true,
	}

	// Pre paid instance can not be release.
	if d.Get("paid_type").(string) == string(PrePaid) {
		return nil
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"15"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	action = "DeleteInstance"
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
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

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf = BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func formatSelectedZonesReq(configured []interface{}) string {
	doubleList := make([][]interface{}, len(configured))
	for i, v := range configured {
		doubleList[i] = []interface{}{v}
	}

	if len(doubleList) < 1 {
		return ""
	}

	if len(doubleList) == 1 {
		return "[[\"" + doubleList[0][0].(string) + "\"],[]]"
	}

	result := "[["

	for i := 0; i < len(doubleList); i++ {
		switch i {
		case len(doubleList) - 2:
			result += "\"" + doubleList[i][0].(string) + "\""
		case len(doubleList) - 1:
			result += "],[\"" + doubleList[i][0].(string) + "\"]"
		default:
			result += "\"" + doubleList[i][0].(string) + "\","
		}
	}

	result += "]"

	return result
}

func convertAliKafkaAutoCreateTopicEnableResponse(source interface{}) interface{} {
	switch source {
	case true:
		return "enable"
	case false:
		return "disable"
	}

	return source
}

func convertAliKafkaPaidTypeResponse(source interface{}) interface{} {
	switch source {
	case "0", "4":
		return "PrePaid"
	case "1", "3":
		return "PostPaid"
	}

	return "PostPaid"
}

func convertAliKafkaPostPaidTypeRequest(source string) int {
	switch source {
	case "alikafka":
		return 1
	case "alikafka_serverless":
		return 3
	}

	return 1
}

func convertAliKafkaPrePaidTypeRequest(source string) int {
	switch source {
	case "alikafka":
		return 0
	case "alikafka_confluent":
		return 4
	}

	return 0
}

func convertAliKafkaInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "0", "1":
		return "alikafka"
	case "3":
		return "alikafka_serverless"
	case "4":
		return "alikafka_confluent"
	}

	return "alikafka"
}
