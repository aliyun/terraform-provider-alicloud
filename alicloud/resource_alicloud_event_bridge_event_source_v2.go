package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeEventSourceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeEventSourceV2Create,
		Read:   resourceAliCloudEventBridgeEventSourceV2Read,
		Update: resourceAliCloudEventBridgeEventSourceV2Update,
		Delete: resourceAliCloudEventBridgeEventSourceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"linked_external_source": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_http_event_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"referer": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"public_web_hook_url": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"method": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_web_hook_url": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"source_kafka_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"offset_reset": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"consumer_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source_mns_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_base64_decode": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source_oss_event_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sts_role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"event_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"match_state": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"suffix": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"prefix": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"name": {
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
			"source_rabbit_mq_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source_rocketmq_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timestamp": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"offset": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_network": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_vswitch_ids": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auth_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source_sls_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"log_store": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"consume_position": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"source_scheduled_event_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudEventBridgeEventSourceV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEventSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("event_source_name"); ok {
		request["EventSourceName"] = v
	}

	sourceHttpEventParameters := make(map[string]interface{})

	if v := d.Get("source_http_event_parameters"); !IsNil(v) {
		ip1, _ := jsonpath.Get("$[0].ip", v)
		if ip1 != nil && ip1 != "" {
			sourceHttpEventParameters["Ip"] = ip1
		}
		method1, _ := jsonpath.Get("$[0].method", v)
		if method1 != nil && method1 != "" {
			sourceHttpEventParameters["Method"] = method1
		}
		securityConfig1, _ := jsonpath.Get("$[0].security_config", v)
		if securityConfig1 != nil && securityConfig1 != "" {
			sourceHttpEventParameters["SecurityConfig"] = securityConfig1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			sourceHttpEventParameters["Type"] = type1
		}
		referer1, _ := jsonpath.Get("$[0].referer", v)
		if referer1 != nil && referer1 != "" {
			sourceHttpEventParameters["Referer"] = referer1
		}

		sourceHttpEventParametersJson, err := json.Marshal(sourceHttpEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceHttpEventParameters"] = string(sourceHttpEventParametersJson)
	}

	sourceKafkaParameters := make(map[string]interface{})

	if v := d.Get("source_kafka_parameters"); !IsNil(v) {
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			sourceKafkaParameters["VpcId"] = vpcId1
		}
		securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
		if securityGroupId1 != nil && securityGroupId1 != "" {
			sourceKafkaParameters["SecurityGroupId"] = securityGroupId1
		}
		topic1, _ := jsonpath.Get("$[0].topic", v)
		if topic1 != nil && topic1 != "" {
			sourceKafkaParameters["Topic"] = topic1
		}
		vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", v)
		if vSwitchIds1 != nil && vSwitchIds1 != "" {
			sourceKafkaParameters["VSwitchIds"] = vSwitchIds1
		}
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			sourceKafkaParameters["InstanceId"] = instanceId1
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && regionId1 != "" {
			sourceKafkaParameters["RegionId"] = regionId1
		}
		network1, _ := jsonpath.Get("$[0].network", v)
		if network1 != nil && network1 != "" {
			sourceKafkaParameters["Network"] = network1
		}
		consumerGroup1, _ := jsonpath.Get("$[0].consumer_group", v)
		if consumerGroup1 != nil && consumerGroup1 != "" {
			sourceKafkaParameters["ConsumerGroup"] = consumerGroup1
		}
		offsetReset1, _ := jsonpath.Get("$[0].offset_reset", v)
		if offsetReset1 != nil && offsetReset1 != "" {
			sourceKafkaParameters["OffsetReset"] = offsetReset1
		}

		sourceKafkaParametersJson, err := json.Marshal(sourceKafkaParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceKafkaParameters"] = string(sourceKafkaParametersJson)
	}

	sourceRabbitMQParameters := make(map[string]interface{})

	if v := d.Get("source_rabbit_mq_parameters"); !IsNil(v) {
		virtualHostName1, _ := jsonpath.Get("$[0].virtual_host_name", v)
		if virtualHostName1 != nil && virtualHostName1 != "" {
			sourceRabbitMQParameters["VirtualHostName"] = virtualHostName1
		}
		instanceId3, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId3 != nil && instanceId3 != "" {
			sourceRabbitMQParameters["InstanceId"] = instanceId3
		}
		queueName1, _ := jsonpath.Get("$[0].queue_name", v)
		if queueName1 != nil && queueName1 != "" {
			sourceRabbitMQParameters["QueueName"] = queueName1
		}
		regionId3, _ := jsonpath.Get("$[0].region_id", v)
		if regionId3 != nil && regionId3 != "" {
			sourceRabbitMQParameters["RegionId"] = regionId3
		}

		sourceRabbitMQParametersJson, err := json.Marshal(sourceRabbitMQParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceRabbitMQParameters"] = string(sourceRabbitMQParametersJson)
	}

	sourceMNSParameters := make(map[string]interface{})

	if v := d.Get("source_mns_parameters"); !IsNil(v) {
		isBase64Decode1, _ := jsonpath.Get("$[0].is_base64_decode", v)
		if isBase64Decode1 != nil && isBase64Decode1 != "" {
			sourceMNSParameters["IsBase64Decode"] = isBase64Decode1
		}
		regionId5, _ := jsonpath.Get("$[0].region_id", v)
		if regionId5 != nil && regionId5 != "" {
			sourceMNSParameters["RegionId"] = regionId5
		}
		queueName3, _ := jsonpath.Get("$[0].queue_name", v)
		if queueName3 != nil && queueName3 != "" {
			sourceMNSParameters["QueueName"] = queueName3
		}

		sourceMNSParametersJson, err := json.Marshal(sourceMNSParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceMNSParameters"] = string(sourceMNSParametersJson)
	}

	sourceRocketMQParameters := make(map[string]interface{})

	if v := d.Get("source_rocketmq_parameters"); !IsNil(v) {
		instanceNetwork1, _ := jsonpath.Get("$[0].instance_network", v)
		if instanceNetwork1 != nil && instanceNetwork1 != "" {
			sourceRocketMQParameters["InstanceNetwork"] = instanceNetwork1
		}
		topic3, _ := jsonpath.Get("$[0].topic", v)
		if topic3 != nil && topic3 != "" {
			sourceRocketMQParameters["Topic"] = topic3
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && tag1 != "" {
			sourceRocketMQParameters["Tag"] = tag1
		}
		instanceUsername1, _ := jsonpath.Get("$[0].instance_username", v)
		if instanceUsername1 != nil && instanceUsername1 != "" {
			sourceRocketMQParameters["InstanceUsername"] = instanceUsername1
		}
		instanceVSwitchIds1, _ := jsonpath.Get("$[0].instance_vswitch_ids", v)
		if instanceVSwitchIds1 != nil && instanceVSwitchIds1 != "" {
			sourceRocketMQParameters["InstanceVSwitchIds"] = instanceVSwitchIds1
		}
		instanceSecurityGroupId1, _ := jsonpath.Get("$[0].instance_security_group_id", v)
		if instanceSecurityGroupId1 != nil && instanceSecurityGroupId1 != "" {
			sourceRocketMQParameters["InstanceSecurityGroupId"] = instanceSecurityGroupId1
		}
		timestamp1, _ := jsonpath.Get("$[0].timestamp", v)
		if timestamp1 != nil && timestamp1 != "" {
			sourceRocketMQParameters["Timestamp"] = timestamp1
		}
		groupId, _ := jsonpath.Get("$[0].group_id", v)
		if groupId != nil && groupId != "" {
			sourceRocketMQParameters["GroupID"] = groupId
		}
		instancePassword1, _ := jsonpath.Get("$[0].instance_password", v)
		if instancePassword1 != nil && instancePassword1 != "" {
			sourceRocketMQParameters["InstancePassword"] = instancePassword1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", v)
		if authType1 != nil && authType1 != "" {
			sourceRocketMQParameters["AuthType"] = authType1
		}
		offset1, _ := jsonpath.Get("$[0].offset", v)
		if offset1 != nil && offset1 != "" {
			sourceRocketMQParameters["Offset"] = offset1
		}
		instanceId5, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId5 != nil && instanceId5 != "" {
			sourceRocketMQParameters["InstanceId"] = instanceId5
		}
		instanceType1, _ := jsonpath.Get("$[0].instance_type", v)
		if instanceType1 != nil && instanceType1 != "" {
			sourceRocketMQParameters["InstanceType"] = instanceType1
		}
		instanceVpcId1, _ := jsonpath.Get("$[0].instance_vpc_id", v)
		if instanceVpcId1 != nil && instanceVpcId1 != "" {
			sourceRocketMQParameters["InstanceVpcId"] = instanceVpcId1
		}
		regionId7, _ := jsonpath.Get("$[0].region_id", v)
		if regionId7 != nil && regionId7 != "" {
			sourceRocketMQParameters["RegionId"] = regionId7
		}
		instanceEndpoint1, _ := jsonpath.Get("$[0].instance_endpoint", v)
		if instanceEndpoint1 != nil && instanceEndpoint1 != "" {
			sourceRocketMQParameters["InstanceEndpoint"] = instanceEndpoint1
		}

		sourceRocketMQParametersJson, err := json.Marshal(sourceRocketMQParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceRocketMQParameters"] = string(sourceRocketMQParametersJson)
	}

	sourceOSSEventParameters := make(map[string]interface{})

	if v := d.Get("source_oss_event_parameters"); !IsNil(v) {
		if v, ok := d.GetOk("source_oss_event_parameters"); ok {
			localData, err := jsonpath.Get("$[0].match_rules", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([][]map[string]interface{}, 0)
			for _, outerList := range convertToInterfaceArray(localData) {
				innerList := make([]map[string]interface{}, 0)
				for _, innerItem := range outerList.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if innerItem != nil {
						dataLoopTmp = innerItem.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["MatchState"] = dataLoopTmp["match_state"]
					dataLoopMap["Suffix"] = dataLoopTmp["suffix"]
					dataLoopMap["Prefix"] = dataLoopTmp["prefix"]
					dataLoopMap["Name"] = dataLoopTmp["name"]

					innerList = append(innerList, dataLoopMap)
				}

				localMaps = append(localMaps, innerList)
			}
			sourceOSSEventParameters["MatchRules"] = localMaps
		}

		stsRoleArn1, _ := jsonpath.Get("$[0].sts_role_arn", v)
		if stsRoleArn1 != nil && stsRoleArn1 != "" {
			sourceOSSEventParameters["StsRoleArn"] = stsRoleArn1
		}
		eventTypes1, _ := jsonpath.Get("$[0].event_types", v)
		if eventTypes1 != nil && eventTypes1 != "" {
			sourceOSSEventParameters["EventTypes"] = eventTypes1
		}

		sourceOSSEventParametersJson, err := json.Marshal(sourceOSSEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceOSSEventParameters"] = string(sourceOSSEventParametersJson)
	}

	sourceScheduledEventParameters := make(map[string]interface{})

	if v := d.Get("source_scheduled_event_parameters"); !IsNil(v) {
		schedule1, _ := jsonpath.Get("$[0].schedule", v)
		if schedule1 != nil && schedule1 != "" {
			sourceScheduledEventParameters["Schedule"] = schedule1
		}
		userData1, _ := jsonpath.Get("$[0].user_data", v)
		if userData1 != nil && userData1 != "" {
			sourceScheduledEventParameters["UserData"] = userData1
		}
		timeZone1, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone1 != nil && timeZone1 != "" {
			sourceScheduledEventParameters["TimeZone"] = timeZone1
		}

		sourceScheduledEventParametersJson, err := json.Marshal(sourceScheduledEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceScheduledEventParameters"] = string(sourceScheduledEventParametersJson)
	}

	request["EventBusName"] = d.Get("event_bus_name")
	sourceSLSParameters := make(map[string]interface{})

	if v := d.Get("source_sls_parameters"); !IsNil(v) {
		logStore1, _ := jsonpath.Get("$[0].log_store", v)
		if logStore1 != nil && logStore1 != "" {
			sourceSLSParameters["LogStore"] = logStore1
		}
		consumePosition1, _ := jsonpath.Get("$[0].consume_position", v)
		if consumePosition1 != nil && consumePosition1 != "" {
			sourceSLSParameters["ConsumePosition"] = consumePosition1
		}
		roleName1, _ := jsonpath.Get("$[0].role_name", v)
		if roleName1 != nil && roleName1 != "" {
			sourceSLSParameters["RoleName"] = roleName1
		}
		project1, _ := jsonpath.Get("$[0].project", v)
		if project1 != nil && project1 != "" {
			sourceSLSParameters["Project"] = project1
		}

		sourceSLSParametersJson, err := json.Marshal(sourceSLSParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceSLSParameters"] = string(sourceSLSParametersJson)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("linked_external_source"); ok {
		request["LinkedExternalSource"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_event_source_v2", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["EventSourceName"]))

	return resourceAliCloudEventBridgeEventSourceV2Read(d, meta)
}

func resourceAliCloudEventBridgeEventSourceV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeServiceV2 := EventBridgeServiceV2{client}

	objectRaw, err := eventBridgeServiceV2.DescribeEventBridgeEventSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_source_v2 DescribeEventBridgeEventSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("event_bus_name", objectRaw["EventBusName"])
	d.Set("event_source_name", objectRaw["Name"])

	sourceHttpEventParametersMaps := make([]map[string]interface{}, 0)
	sourceHttpEventParametersMap := make(map[string]interface{})
	sourceHttpEventParametersRaw := make(map[string]interface{})
	if objectRaw["SourceHttpEventParameters"] != nil {
		sourceHttpEventParametersRaw = objectRaw["SourceHttpEventParameters"].(map[string]interface{})
	}
	if len(sourceHttpEventParametersRaw) > 0 {
		sourceHttpEventParametersMap["security_config"] = sourceHttpEventParametersRaw["SecurityConfig"]
		sourceHttpEventParametersMap["type"] = sourceHttpEventParametersRaw["Type"]

		ipRaw := make([]interface{}, 0)
		if sourceHttpEventParametersRaw["Ip"] != nil {
			ipRaw = convertToInterfaceArray(sourceHttpEventParametersRaw["Ip"])
		}

		sourceHttpEventParametersMap["ip"] = ipRaw
		methodRaw := make([]interface{}, 0)
		if sourceHttpEventParametersRaw["Method"] != nil {
			methodRaw = convertToInterfaceArray(sourceHttpEventParametersRaw["Method"])
		}

		sourceHttpEventParametersMap["method"] = methodRaw
		publicWebHookUrlRaw := make([]interface{}, 0)
		if sourceHttpEventParametersRaw["PublicWebHookUrl"] != nil {
			publicWebHookUrlRaw = convertToInterfaceArray(sourceHttpEventParametersRaw["PublicWebHookUrl"])
		}

		sourceHttpEventParametersMap["public_web_hook_url"] = publicWebHookUrlRaw
		refererRaw := make([]interface{}, 0)
		if sourceHttpEventParametersRaw["Referer"] != nil {
			refererRaw = convertToInterfaceArray(sourceHttpEventParametersRaw["Referer"])
		}

		sourceHttpEventParametersMap["referer"] = refererRaw
		vpcWebHookUrlRaw := make([]interface{}, 0)
		if sourceHttpEventParametersRaw["VpcWebHookUrl"] != nil {
			vpcWebHookUrlRaw = convertToInterfaceArray(sourceHttpEventParametersRaw["VpcWebHookUrl"])
		}

		sourceHttpEventParametersMap["vpc_web_hook_url"] = vpcWebHookUrlRaw
		sourceHttpEventParametersMaps = append(sourceHttpEventParametersMaps, sourceHttpEventParametersMap)
	}
	if err := d.Set("source_http_event_parameters", sourceHttpEventParametersMaps); err != nil {
		return err
	}
	sourceKafkaParametersMaps := make([]map[string]interface{}, 0)
	sourceKafkaParametersMap := make(map[string]interface{})
	sourceKafkaParametersRaw := make(map[string]interface{})
	if objectRaw["SourceKafkaParameters"] != nil {
		sourceKafkaParametersRaw = objectRaw["SourceKafkaParameters"].(map[string]interface{})
	}
	if len(sourceKafkaParametersRaw) > 0 {
		sourceKafkaParametersMap["consumer_group"] = sourceKafkaParametersRaw["ConsumerGroup"]
		sourceKafkaParametersMap["instance_id"] = sourceKafkaParametersRaw["InstanceId"]
		sourceKafkaParametersMap["network"] = sourceKafkaParametersRaw["Network"]
		sourceKafkaParametersMap["offset_reset"] = sourceKafkaParametersRaw["OffsetReset"]
		sourceKafkaParametersMap["region_id"] = sourceKafkaParametersRaw["RegionId"]
		sourceKafkaParametersMap["security_group_id"] = sourceKafkaParametersRaw["SecurityGroupId"]
		sourceKafkaParametersMap["topic"] = sourceKafkaParametersRaw["Topic"]
		sourceKafkaParametersMap["vswitch_ids"] = sourceKafkaParametersRaw["VSwitchIds"]
		sourceKafkaParametersMap["vpc_id"] = sourceKafkaParametersRaw["VpcId"]

		sourceKafkaParametersMaps = append(sourceKafkaParametersMaps, sourceKafkaParametersMap)
	}
	if err := d.Set("source_kafka_parameters", sourceKafkaParametersMaps); err != nil {
		return err
	}
	sourceMNSParametersMaps := make([]map[string]interface{}, 0)
	sourceMNSParametersMap := make(map[string]interface{})
	sourceMNSParametersRaw := make(map[string]interface{})
	if objectRaw["SourceMNSParameters"] != nil {
		sourceMNSParametersRaw = objectRaw["SourceMNSParameters"].(map[string]interface{})
	}
	if len(sourceMNSParametersRaw) > 0 {
		sourceMNSParametersMap["is_base64_decode"] = sourceMNSParametersRaw["IsBase64Decode"]
		sourceMNSParametersMap["queue_name"] = sourceMNSParametersRaw["QueueName"]
		sourceMNSParametersMap["region_id"] = sourceMNSParametersRaw["RegionId"]

		sourceMNSParametersMaps = append(sourceMNSParametersMaps, sourceMNSParametersMap)
	}
	if err := d.Set("source_mns_parameters", sourceMNSParametersMaps); err != nil {
		return err
	}
	sourceOSSEventParametersMaps := make([]map[string]interface{}, 0)
	sourceOSSEventParametersMap := make(map[string]interface{})
	sourceOSSEventParametersRaw := make(map[string]interface{})
	if objectRaw["SourceOSSEventParameters"] != nil {
		sourceOSSEventParametersRaw = objectRaw["SourceOSSEventParameters"].(map[string]interface{})
	}
	if len(sourceOSSEventParametersRaw) > 0 {
		sourceOSSEventParametersMap["sts_role_arn"] = sourceOSSEventParametersRaw["StsRoleArn"]

		eventTypesRaw := make([]interface{}, 0)
		if sourceOSSEventParametersRaw["EventTypes"] != nil {
			eventTypesRaw = convertToInterfaceArray(sourceOSSEventParametersRaw["EventTypes"])
		}

		sourceOSSEventParametersMap["event_types"] = eventTypesRaw

		matchRulesChildRaw := sourceOSSEventParametersRaw["MatchRules"]
		matchRulesMaps := make([][]map[string]interface{}, 0)
		if matchRulesChildRaw != nil {
			for _, matchRulesChildChildRaw := range matchRulesChildRaw.([]interface{}) {
				matchRulesChildMaps := make([]map[string]interface{}, 0)
				for _, matchRulesChildChildChildRaw := range matchRulesChildChildRaw.([]interface{}) {
					matchRulesMap := make(map[string]interface{})
					matchRulesChildChildChild := matchRulesChildChildChildRaw.(map[string]interface{})
					matchRulesMap["match_state"] = fmt.Sprint(matchRulesChildChildChild["MatchState"])
					matchRulesMap["name"] = matchRulesChildChildChild["Name"]
					matchRulesMap["prefix"] = matchRulesChildChildChild["Prefix"]
					matchRulesMap["suffix"] = matchRulesChildChildChild["Suffix"]
					matchRulesChildMaps = append(matchRulesChildMaps, matchRulesMap)
				}

				matchRulesMaps = append(matchRulesMaps, matchRulesChildMaps)
			}
		}
		sourceOSSEventParametersMap["match_rules"] = matchRulesMaps
		sourceOSSEventParametersMaps = append(sourceOSSEventParametersMaps, sourceOSSEventParametersMap)
	}
	if err := d.Set("source_oss_event_parameters", sourceOSSEventParametersMaps); err != nil {
		return err
	}
	sourceRabbitMQParametersMaps := make([]map[string]interface{}, 0)
	sourceRabbitMQParametersMap := make(map[string]interface{})
	sourceRabbitMQParametersRaw := make(map[string]interface{})
	if objectRaw["SourceRabbitMQParameters"] != nil {
		sourceRabbitMQParametersRaw = objectRaw["SourceRabbitMQParameters"].(map[string]interface{})
	}
	if len(sourceRabbitMQParametersRaw) > 0 {
		sourceRabbitMQParametersMap["instance_id"] = sourceRabbitMQParametersRaw["InstanceId"]
		sourceRabbitMQParametersMap["queue_name"] = sourceRabbitMQParametersRaw["QueueName"]
		sourceRabbitMQParametersMap["region_id"] = sourceRabbitMQParametersRaw["RegionId"]
		sourceRabbitMQParametersMap["virtual_host_name"] = sourceRabbitMQParametersRaw["VirtualHostName"]

		sourceRabbitMQParametersMaps = append(sourceRabbitMQParametersMaps, sourceRabbitMQParametersMap)
	}
	if err := d.Set("source_rabbit_mq_parameters", sourceRabbitMQParametersMaps); err != nil {
		return err
	}
	sourceRocketMQParametersMaps := make([]map[string]interface{}, 0)
	sourceRocketMQParametersMap := make(map[string]interface{})
	sourceRocketMQParametersRaw := make(map[string]interface{})
	if objectRaw["SourceRocketMQParameters"] != nil {
		sourceRocketMQParametersRaw = objectRaw["SourceRocketMQParameters"].(map[string]interface{})
	}
	if len(sourceRocketMQParametersRaw) > 0 {
		sourceRocketMQParametersMap["auth_type"] = sourceRocketMQParametersRaw["AuthType"]
		sourceRocketMQParametersMap["group_id"] = sourceRocketMQParametersRaw["GroupId"]
		sourceRocketMQParametersMap["instance_endpoint"] = sourceRocketMQParametersRaw["InstanceEndpoint"]
		sourceRocketMQParametersMap["instance_id"] = sourceRocketMQParametersRaw["InstanceId"]
		sourceRocketMQParametersMap["instance_network"] = sourceRocketMQParametersRaw["InstanceNetwork"]
		sourceRocketMQParametersMap["instance_password"] = sourceRocketMQParametersRaw["InstancePassword"]
		sourceRocketMQParametersMap["instance_security_group_id"] = sourceRocketMQParametersRaw["InstanceSecurityGroupId"]
		sourceRocketMQParametersMap["instance_type"] = sourceRocketMQParametersRaw["InstanceType"]
		sourceRocketMQParametersMap["instance_username"] = sourceRocketMQParametersRaw["InstanceUsername"]
		sourceRocketMQParametersMap["instance_vswitch_ids"] = sourceRocketMQParametersRaw["InstanceVSwitchIds"]
		sourceRocketMQParametersMap["instance_vpc_id"] = sourceRocketMQParametersRaw["InstanceVpcId"]
		sourceRocketMQParametersMap["offset"] = sourceRocketMQParametersRaw["Offset"]
		sourceRocketMQParametersMap["region_id"] = sourceRocketMQParametersRaw["RegionId"]
		sourceRocketMQParametersMap["tag"] = sourceRocketMQParametersRaw["Tag"]
		sourceRocketMQParametersMap["timestamp"] = sourceRocketMQParametersRaw["Timestamp"]
		sourceRocketMQParametersMap["topic"] = sourceRocketMQParametersRaw["Topic"]

		sourceRocketMQParametersMaps = append(sourceRocketMQParametersMaps, sourceRocketMQParametersMap)
	}
	if err := d.Set("source_rocketmq_parameters", sourceRocketMQParametersMaps); err != nil {
		return err
	}
	sourceSLSParametersMaps := make([]map[string]interface{}, 0)
	sourceSLSParametersMap := make(map[string]interface{})
	sourceSLSParametersRaw := make(map[string]interface{})
	if objectRaw["SourceSLSParameters"] != nil {
		sourceSLSParametersRaw = objectRaw["SourceSLSParameters"].(map[string]interface{})
	}
	if len(sourceSLSParametersRaw) > 0 {
		sourceSLSParametersMap["consume_position"] = sourceSLSParametersRaw["ConsumePosition"]
		sourceSLSParametersMap["log_store"] = sourceSLSParametersRaw["LogStore"]
		sourceSLSParametersMap["project"] = sourceSLSParametersRaw["Project"]
		sourceSLSParametersMap["role_name"] = sourceSLSParametersRaw["RoleName"]

		sourceSLSParametersMaps = append(sourceSLSParametersMaps, sourceSLSParametersMap)
	}
	if err := d.Set("source_sls_parameters", sourceSLSParametersMaps); err != nil {
		return err
	}
	sourceScheduledEventParametersMaps := make([]map[string]interface{}, 0)
	sourceScheduledEventParametersMap := make(map[string]interface{})
	sourceScheduledEventParametersRaw := make(map[string]interface{})
	if objectRaw["SourceScheduledEventParameters"] != nil {
		sourceScheduledEventParametersRaw = objectRaw["SourceScheduledEventParameters"].(map[string]interface{})
	}
	if len(sourceScheduledEventParametersRaw) > 0 {
		sourceScheduledEventParametersMap["schedule"] = sourceScheduledEventParametersRaw["Schedule"]
		sourceScheduledEventParametersMap["time_zone"] = sourceScheduledEventParametersRaw["TimeZone"]
		sourceScheduledEventParametersMap["user_data"] = sourceScheduledEventParametersRaw["UserData"]

		sourceScheduledEventParametersMaps = append(sourceScheduledEventParametersMaps, sourceScheduledEventParametersMap)
	}
	if err := d.Set("source_scheduled_event_parameters", sourceScheduledEventParametersMaps); err != nil {
		return err
	}

	d.Set("event_source_name", d.Id())

	return nil
}

func resourceAliCloudEventBridgeEventSourceV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateEventSource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EventSourceName"] = d.Id()
	request["EventBusName"] = d.Get("event_bus_name")

	if d.HasChange("source_http_event_parameters") {
		update = true
	}
	sourceHttpEventParameters := make(map[string]interface{})

	if v := d.Get("source_http_event_parameters"); v != nil {
		ip1, _ := jsonpath.Get("$[0].ip", v)
		if ip1 != nil && (d.HasChange("source_http_event_parameters.0.ip") || ip1 != "") {
			sourceHttpEventParameters["Ip"] = ip1
		}
		method1, _ := jsonpath.Get("$[0].method", v)
		if method1 != nil && (d.HasChange("source_http_event_parameters.0.method") || method1 != "") {
			sourceHttpEventParameters["Method"] = method1
		}
		securityConfig1, _ := jsonpath.Get("$[0].security_config", v)
		if securityConfig1 != nil && (d.HasChange("source_http_event_parameters.0.security_config") || securityConfig1 != "") {
			sourceHttpEventParameters["SecurityConfig"] = securityConfig1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && (d.HasChange("source_http_event_parameters.0.type") || type1 != "") {
			sourceHttpEventParameters["Type"] = type1
		}
		referer1, _ := jsonpath.Get("$[0].referer", v)
		if referer1 != nil && (d.HasChange("source_http_event_parameters.0.referer") || referer1 != "") {
			sourceHttpEventParameters["Referer"] = referer1
		}

		sourceHttpEventParametersJson, err := json.Marshal(sourceHttpEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceHttpEventParameters"] = string(sourceHttpEventParametersJson)
	}

	if d.HasChange("source_kafka_parameters") {
		update = true
	}
	sourceKafkaParameters := make(map[string]interface{})

	if v := d.Get("source_kafka_parameters"); v != nil {
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && (d.HasChange("source_kafka_parameters.0.vpc_id") || vpcId1 != "") {
			sourceKafkaParameters["VpcId"] = vpcId1
		}
		securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
		if securityGroupId1 != nil && (d.HasChange("source_kafka_parameters.0.security_group_id") || securityGroupId1 != "") {
			sourceKafkaParameters["SecurityGroupId"] = securityGroupId1
		}
		topic1, _ := jsonpath.Get("$[0].topic", v)
		if topic1 != nil && (d.HasChange("source_kafka_parameters.0.topic") || topic1 != "") {
			sourceKafkaParameters["Topic"] = topic1
		}
		vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", v)
		if vSwitchIds1 != nil && (d.HasChange("source_kafka_parameters.0.vswitch_ids") || vSwitchIds1 != "") {
			sourceKafkaParameters["VSwitchIds"] = vSwitchIds1
		}
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && (d.HasChange("source_kafka_parameters.0.instance_id") || instanceId1 != "") {
			sourceKafkaParameters["InstanceId"] = instanceId1
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && (d.HasChange("source_kafka_parameters.0.region_id") || regionId1 != "") {
			sourceKafkaParameters["RegionId"] = regionId1
		}
		network1, _ := jsonpath.Get("$[0].network", v)
		if network1 != nil && (d.HasChange("source_kafka_parameters.0.network") || network1 != "") {
			sourceKafkaParameters["Network"] = network1
		}
		consumerGroup1, _ := jsonpath.Get("$[0].consumer_group", v)
		if consumerGroup1 != nil && (d.HasChange("source_kafka_parameters.0.consumer_group") || consumerGroup1 != "") {
			sourceKafkaParameters["ConsumerGroup"] = consumerGroup1
		}
		offsetReset1, _ := jsonpath.Get("$[0].offset_reset", v)
		if offsetReset1 != nil && (d.HasChange("source_kafka_parameters.0.offset_reset") || offsetReset1 != "") {
			sourceKafkaParameters["OffsetReset"] = offsetReset1
		}

		sourceKafkaParametersJson, err := json.Marshal(sourceKafkaParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceKafkaParameters"] = string(sourceKafkaParametersJson)
	}

	if d.HasChange("source_rabbit_mq_parameters") {
		update = true
	}
	sourceRabbitMQParameters := make(map[string]interface{})

	if v := d.Get("source_rabbit_mq_parameters"); v != nil {
		virtualHostName1, _ := jsonpath.Get("$[0].virtual_host_name", v)
		if virtualHostName1 != nil && (d.HasChange("source_rabbit_mq_parameters.0.virtual_host_name") || virtualHostName1 != "") {
			sourceRabbitMQParameters["VirtualHostName"] = virtualHostName1
		}
		instanceId3, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId3 != nil && (d.HasChange("source_rabbit_mq_parameters.0.instance_id") || instanceId3 != "") {
			sourceRabbitMQParameters["InstanceId"] = instanceId3
		}
		queueName1, _ := jsonpath.Get("$[0].queue_name", v)
		if queueName1 != nil && (d.HasChange("source_rabbit_mq_parameters.0.queue_name") || queueName1 != "") {
			sourceRabbitMQParameters["QueueName"] = queueName1
		}
		regionId3, _ := jsonpath.Get("$[0].region_id", v)
		if regionId3 != nil && (d.HasChange("source_rabbit_mq_parameters.0.region_id") || regionId3 != "") {
			sourceRabbitMQParameters["RegionId"] = regionId3
		}

		sourceRabbitMQParametersJson, err := json.Marshal(sourceRabbitMQParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceRabbitMQParameters"] = string(sourceRabbitMQParametersJson)
	}

	if d.HasChange("source_mns_parameters") {
		update = true
	}
	sourceMNSParameters := make(map[string]interface{})

	if v := d.Get("source_mns_parameters"); v != nil {
		isBase64Decode1, _ := jsonpath.Get("$[0].is_base64_decode", v)
		if isBase64Decode1 != nil && (d.HasChange("source_mns_parameters.0.is_base64_decode") || isBase64Decode1 != "") {
			sourceMNSParameters["IsBase64Decode"] = isBase64Decode1
		}
		regionId5, _ := jsonpath.Get("$[0].region_id", v)
		if regionId5 != nil && (d.HasChange("source_mns_parameters.0.region_id") || regionId5 != "") {
			sourceMNSParameters["RegionId"] = regionId5
		}
		queueName3, _ := jsonpath.Get("$[0].queue_name", v)
		if queueName3 != nil && (d.HasChange("source_mns_parameters.0.queue_name") || queueName3 != "") {
			sourceMNSParameters["QueueName"] = queueName3
		}

		sourceMNSParametersJson, err := json.Marshal(sourceMNSParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceMNSParameters"] = string(sourceMNSParametersJson)
	}

	if d.HasChange("source_rocketmq_parameters") {
		update = true
	}
	sourceRocketMQParameters := make(map[string]interface{})

	if v := d.Get("source_rocketmq_parameters"); v != nil {
		instanceNetwork1, _ := jsonpath.Get("$[0].instance_network", v)
		if instanceNetwork1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_network") || instanceNetwork1 != "") {
			sourceRocketMQParameters["InstanceNetwork"] = instanceNetwork1
		}
		topic3, _ := jsonpath.Get("$[0].topic", v)
		if topic3 != nil && (d.HasChange("source_rocketmq_parameters.0.topic") || topic3 != "") {
			sourceRocketMQParameters["Topic"] = topic3
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && (d.HasChange("source_rocketmq_parameters.0.tag") || tag1 != "") {
			sourceRocketMQParameters["Tag"] = tag1
		}
		instanceUsername1, _ := jsonpath.Get("$[0].instance_username", v)
		if instanceUsername1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_username") || instanceUsername1 != "") {
			sourceRocketMQParameters["InstanceUsername"] = instanceUsername1
		}
		instanceVSwitchIds1, _ := jsonpath.Get("$[0].instance_vswitch_ids", v)
		if instanceVSwitchIds1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_vswitch_ids") || instanceVSwitchIds1 != "") {
			sourceRocketMQParameters["InstanceVSwitchIds"] = instanceVSwitchIds1
		}
		instanceSecurityGroupId1, _ := jsonpath.Get("$[0].instance_security_group_id", v)
		if instanceSecurityGroupId1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_security_group_id") || instanceSecurityGroupId1 != "") {
			sourceRocketMQParameters["InstanceSecurityGroupId"] = instanceSecurityGroupId1
		}
		timestamp1, _ := jsonpath.Get("$[0].timestamp", v)
		if timestamp1 != nil && (d.HasChange("source_rocketmq_parameters.0.timestamp") || timestamp1 != "") {
			sourceRocketMQParameters["Timestamp"] = timestamp1
		}
		groupId, _ := jsonpath.Get("$[0].group_id", v)
		if groupId != nil && (d.HasChange("source_rocketmq_parameters.0.group_id") || groupId != "") {
			sourceRocketMQParameters["GroupID"] = groupId
		}
		instancePassword1, _ := jsonpath.Get("$[0].instance_password", v)
		if instancePassword1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_password") || instancePassword1 != "") {
			sourceRocketMQParameters["InstancePassword"] = instancePassword1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", v)
		if authType1 != nil && (d.HasChange("source_rocketmq_parameters.0.auth_type") || authType1 != "") {
			sourceRocketMQParameters["AuthType"] = authType1
		}
		offset1, _ := jsonpath.Get("$[0].offset", v)
		if offset1 != nil && (d.HasChange("source_rocketmq_parameters.0.offset") || offset1 != "") {
			sourceRocketMQParameters["Offset"] = offset1
		}
		instanceId5, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId5 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_id") || instanceId5 != "") {
			sourceRocketMQParameters["InstanceId"] = instanceId5
		}
		instanceType1, _ := jsonpath.Get("$[0].instance_type", v)
		if instanceType1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_type") || instanceType1 != "") {
			sourceRocketMQParameters["InstanceType"] = instanceType1
		}
		instanceVpcId1, _ := jsonpath.Get("$[0].instance_vpc_id", v)
		if instanceVpcId1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_vpc_id") || instanceVpcId1 != "") {
			sourceRocketMQParameters["InstanceVpcId"] = instanceVpcId1
		}
		regionId7, _ := jsonpath.Get("$[0].region_id", v)
		if regionId7 != nil && (d.HasChange("source_rocketmq_parameters.0.region_id") || regionId7 != "") {
			sourceRocketMQParameters["RegionId"] = regionId7
		}
		instanceEndpoint1, _ := jsonpath.Get("$[0].instance_endpoint", v)
		if instanceEndpoint1 != nil && (d.HasChange("source_rocketmq_parameters.0.instance_endpoint") || instanceEndpoint1 != "") {
			sourceRocketMQParameters["InstanceEndpoint"] = instanceEndpoint1
		}

		sourceRocketMQParametersJson, err := json.Marshal(sourceRocketMQParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceRocketMQParameters"] = string(sourceRocketMQParametersJson)
	}

	if d.HasChange("source_oss_event_parameters") {
		update = true
	}
	sourceOSSEventParameters := make(map[string]interface{})

	if v := d.Get("source_oss_event_parameters"); !IsNil(v) {
		if v, ok := d.GetOk("source_oss_event_parameters"); ok {
			localData, err := jsonpath.Get("$[0].match_rules", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([][]map[string]interface{}, 0)
			for _, outerList := range convertToInterfaceArray(localData) {
				innerList := make([]map[string]interface{}, 0)
				for _, innerItem := range outerList.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if innerItem != nil {
						dataLoopTmp = innerItem.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["MatchState"] = dataLoopTmp["match_state"]
					dataLoopMap["Suffix"] = dataLoopTmp["suffix"]
					dataLoopMap["Prefix"] = dataLoopTmp["prefix"]
					dataLoopMap["Name"] = dataLoopTmp["name"]

					innerList = append(innerList, dataLoopMap)
				}

				localMaps = append(localMaps, innerList)
			}
			sourceOSSEventParameters["MatchRules"] = localMaps
		}

		stsRoleArn1, _ := jsonpath.Get("$[0].sts_role_arn", v)
		if stsRoleArn1 != nil && stsRoleArn1 != "" {
			sourceOSSEventParameters["StsRoleArn"] = stsRoleArn1
		}
		eventTypes1, _ := jsonpath.Get("$[0].event_types", v)
		if eventTypes1 != nil && eventTypes1 != "" {
			sourceOSSEventParameters["EventTypes"] = eventTypes1
		}

		sourceOSSEventParametersJson, err := json.Marshal(sourceOSSEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceOSSEventParameters"] = string(sourceOSSEventParametersJson)
	}

	if d.HasChange("source_scheduled_event_parameters") {
		update = true
	}
	sourceScheduledEventParameters := make(map[string]interface{})

	if v := d.Get("source_scheduled_event_parameters"); v != nil {
		schedule1, _ := jsonpath.Get("$[0].schedule", v)
		if schedule1 != nil && (d.HasChange("source_scheduled_event_parameters.0.schedule") || schedule1 != "") {
			sourceScheduledEventParameters["Schedule"] = schedule1
		}
		userData1, _ := jsonpath.Get("$[0].user_data", v)
		if userData1 != nil && (d.HasChange("source_scheduled_event_parameters.0.user_data") || userData1 != "") {
			sourceScheduledEventParameters["UserData"] = userData1
		}
		timeZone1, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone1 != nil && (d.HasChange("source_scheduled_event_parameters.0.time_zone") || timeZone1 != "") {
			sourceScheduledEventParameters["TimeZone"] = timeZone1
		}

		sourceScheduledEventParametersJson, err := json.Marshal(sourceScheduledEventParameters)
		if err != nil {
			return WrapError(err)
		}
		request["SourceScheduledEventParameters"] = string(sourceScheduledEventParametersJson)
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("linked_external_source"); ok {
		request["LinkedExternalSource"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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

	return resourceAliCloudEventBridgeEventSourceV2Read(d, meta)
}

func resourceAliCloudEventBridgeEventSourceV2Delete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EventSourceName"] = d.Id()

	request["EventBusName"] = d.Get("event_bus_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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

	return nil
}
