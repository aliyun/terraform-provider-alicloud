// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaLoadBalancerCreate,
		Read:   resourceAliCloudEsaLoadBalancerRead,
		Update: resourceAliCloudEsaLoadBalancerUpdate,
		Delete: resourceAliCloudEsaLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"adaptive_routing": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failover_across_pools": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"default_pools": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"fallback_pool": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"monitor": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expected_codes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"follow_redirects": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 10),
						},
						"consecutive_up": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"consecutive_down": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"monitoring_region": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Global", "ChineseMainland", "OutsideChineseMainland"}, false),
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"random_steering": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_weights": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"default_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 100),
						},
					},
				},
			},
			"region_pools": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_enable": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"overrides": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sequence": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"terminates": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"fixed_response": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"status_code": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"location": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"message_body": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"session_affinity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"steering_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sub_region_pools": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(10, 600),
			},
		},
	}
}

func resourceAliCloudEsaLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("region_pools"); ok {
		request["RegionPools"] = v
	}
	if v, ok := d.GetOk("rules"); ok {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["RuleEnable"] = dataLoopTmp["rule_enable"]
			dataLoopMap["Rule"] = dataLoopTmp["rule"]
			if overrides, ok := dataLoopTmp["overrides"]; ok && fmt.Sprint(overrides) != "" {
				dataLoopMap["Overrides"] = overrides
			}

			localData1 := make(map[string]interface{})
			messageBody1, _ := jsonpath.Get("$[0].message_body", dataLoopTmp["fixed_response"])
			if messageBody1 != nil && messageBody1 != "" {
				localData1["MessageBody"] = messageBody1
			}
			location1, _ := jsonpath.Get("$[0].location", dataLoopTmp["fixed_response"])
			if location1 != nil && location1 != "" {
				localData1["Location"] = location1
			}
			statusCode1, _ := jsonpath.Get("$[0].status_code", dataLoopTmp["fixed_response"])
			if statusCode1 != nil && statusCode1 != "" {
				localData1["StatusCode"] = statusCode1
			}
			contentType1, _ := jsonpath.Get("$[0].content_type", dataLoopTmp["fixed_response"])
			if contentType1 != nil && contentType1 != "" {
				localData1["ContentType"] = contentType1
			}
			dataLoopMap["FixedResponse"] = localData1
			dataLoopMap["Sequence"] = dataLoopTmp["sequence"]
			dataLoopMap["Terminates"] = dataLoopTmp["terminates"]
			dataLoopMap["RuleName"] = dataLoopTmp["rule_name"]
			rulesMapsArray = append(rulesMapsArray, dataLoopMap)
		}
		rulesMapsJson, err := json.Marshal(rulesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Rules"] = string(rulesMapsJson)
	}

	dataList := make(map[string]interface{})

	if v := d.Get("monitor"); v != nil {
		path1, _ := jsonpath.Get("$[0].path", v)
		if path1 != nil && path1 != "" {
			dataList["Path"] = path1
		}
		expectedCodes1, _ := jsonpath.Get("$[0].expected_codes", v)
		if expectedCodes1 != nil && expectedCodes1 != "" {
			dataList["ExpectedCodes"] = expectedCodes1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			dataList["Type"] = type1
		}
		consecutiveUp1, _ := jsonpath.Get("$[0].consecutive_up", v)
		if consecutiveUp1 != nil && consecutiveUp1 != "" {
			dataList["ConsecutiveUp"] = consecutiveUp1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			dataList["Port"] = port1
		}
		monitoringRegion1, _ := jsonpath.Get("$[0].monitoring_region", v)
		if monitoringRegion1 != nil && monitoringRegion1 != "" {
			dataList["MonitoringRegion"] = monitoringRegion1
		}
		header1, _ := jsonpath.Get("$[0].header", v)
		if header1 != nil && header1 != "" {
			dataList["Header"] = header1
		}
		method1, _ := jsonpath.Get("$[0].method", v)
		if method1 != nil && method1 != "" {
			dataList["Method"] = method1
		}
		followRedirects1, _ := jsonpath.Get("$[0].follow_redirects", v)
		if followRedirects1 != nil && followRedirects1 != "" {
			dataList["FollowRedirects"] = followRedirects1
		}
		consecutiveDown1, _ := jsonpath.Get("$[0].consecutive_down", v)
		if consecutiveDown1 != nil && consecutiveDown1 != "" {
			dataList["ConsecutiveDown"] = consecutiveDown1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && interval1 != "" {
			dataList["Interval"] = interval1
		}
		timeout1, _ := jsonpath.Get("$[0].timeout", v)
		if timeout1 != nil && timeout1 != "" && timeout1.(int) > 0 {
			dataList["Timeout"] = timeout1
		}

		dataListJson, err := json.Marshal(dataList)
		if err != nil {
			return WrapError(err)
		}
		request["Monitor"] = string(dataListJson)
	}

	if v, ok := d.GetOkExists("enabled"); ok {
		request["Enabled"] = v
	}
	if v, ok := d.GetOk("session_affinity"); ok {
		request["SessionAffinity"] = v
	}
	if v, ok := d.GetOkExists("ttl"); ok && v.(int) > 0 {
		request["Ttl"] = v
	}
	dataList1 := make(map[string]interface{})

	if v := d.Get("random_steering"); !IsNil(v) {
		poolWeights1, _ := jsonpath.Get("$[0].pool_weights", v)
		if poolWeights1 != nil && poolWeights1 != "" {
			dataList1["PoolWeights"] = poolWeights1
		}
		defaultWeight1, _ := jsonpath.Get("$[0].default_weight", v)
		if defaultWeight1 != nil && defaultWeight1 != "" {
			dataList1["DefaultWeight"] = defaultWeight1
		}

		dataList1Json, err := json.Marshal(dataList1)
		if err != nil {
			return WrapError(err)
		}
		request["RandomSteering"] = string(dataList1Json)
	}

	if v, ok := d.GetOk("default_pools"); ok {
		defaultPoolsMapsArray := convertToInterfaceArray(v)

		defaultPoolsMapsJson, err := json.Marshal(defaultPoolsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DefaultPools"] = string(defaultPoolsMapsJson)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	dataList2 := make(map[string]interface{})

	if v := d.Get("adaptive_routing"); !IsNil(v) {
		failoverAcrossPools1, _ := jsonpath.Get("$[0].failover_across_pools", v)
		if failoverAcrossPools1 != nil && failoverAcrossPools1 != "" {
			dataList2["FailoverAcrossPools"] = failoverAcrossPools1
		}

		dataList2Json, err := json.Marshal(dataList2)
		if err != nil {
			return WrapError(err)
		}
		request["AdaptiveRouting"] = string(dataList2Json)
	}

	request["Name"] = d.Get("load_balancer_name")
	request["SteeringPolicy"] = d.Get("steering_policy")
	request["FallbackPool"] = d.Get("fallback_pool")
	if v, ok := d.GetOk("sub_region_pools"); ok {
		request["SubRegionPools"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	return resourceAliCloudEsaLoadBalancerRead(d, meta)
}

func resourceAliCloudEsaLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_load_balancer DescribeEsaLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("enabled", objectRaw["Enabled"])
	d.Set("fallback_pool", objectRaw["FallbackPool"])
	d.Set("load_balancer_name", objectRaw["Name"])
	if objectRaw["RegionPools"] != nil {
		d.Set("region_pools", convertObjectToJsonString(objectRaw["RegionPools"]))
	}
	d.Set("session_affinity", objectRaw["SessionAffinity"])
	d.Set("status", objectRaw["Status"])
	d.Set("steering_policy", objectRaw["SteeringPolicy"])
	if objectRaw["SubRegionPools"] != nil {
		d.Set("sub_region_pools", convertObjectToJsonString(objectRaw["SubRegionPools"]))
	}
	d.Set("ttl", objectRaw["Ttl"])
	d.Set("load_balancer_id", objectRaw["Id"])
	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", v)
	}

	adaptiveRoutingMaps := make([]map[string]interface{}, 0)
	adaptiveRoutingMap := make(map[string]interface{})
	adaptiveRoutingRaw := make(map[string]interface{})
	if objectRaw["AdaptiveRouting"] != nil {
		adaptiveRoutingRaw = objectRaw["AdaptiveRouting"].(map[string]interface{})
	}
	if len(adaptiveRoutingRaw) > 0 {
		adaptiveRoutingMap["failover_across_pools"] = adaptiveRoutingRaw["FailoverAcrossPools"]

		adaptiveRoutingMaps = append(adaptiveRoutingMaps, adaptiveRoutingMap)
	}
	if err := d.Set("adaptive_routing", adaptiveRoutingMaps); err != nil {
		return err
	}
	defaultPoolsRaw := make([]interface{}, 0)
	if objectRaw["DefaultPools"] != nil {
		defaultPoolsRaw = convertToInterfaceArray(objectRaw["DefaultPools"])
	}

	d.Set("default_pools", defaultPoolsRaw)
	monitorMaps := make([]map[string]interface{}, 0)
	monitorMap := make(map[string]interface{})
	monitorRaw := make(map[string]interface{})
	if objectRaw["Monitor"] != nil {
		monitorRaw = objectRaw["Monitor"].(map[string]interface{})
	}
	if len(monitorRaw) > 0 {
		monitorMap["consecutive_down"] = monitorRaw["ConsecutiveDown"]
		monitorMap["consecutive_up"] = monitorRaw["ConsecutiveUp"]
		monitorMap["expected_codes"] = monitorRaw["ExpectedCodes"]
		monitorMap["follow_redirects"] = monitorRaw["FollowRedirects"]
		if monitorRaw["Header"] != nil {
			monitorMap["header"] = convertObjectToJsonString(monitorRaw["Header"])
		}
		monitorMap["interval"] = monitorRaw["Interval"]
		monitorMap["method"] = monitorRaw["Method"]
		monitorMap["monitoring_region"] = monitorRaw["MonitoringRegion"]
		monitorMap["path"] = monitorRaw["Path"]
		monitorMap["port"] = monitorRaw["Port"]
		monitorMap["timeout"] = monitorRaw["Timeout"]
		monitorMap["type"] = monitorRaw["Type"]

		monitorMaps = append(monitorMaps, monitorMap)
	}
	if err := d.Set("monitor", monitorMaps); err != nil {
		return err
	}
	randomSteeringMaps := make([]map[string]interface{}, 0)
	randomSteeringMap := make(map[string]interface{})
	randomSteeringRaw := make(map[string]interface{})
	if objectRaw["RandomSteering"] != nil {
		randomSteeringRaw = objectRaw["RandomSteering"].(map[string]interface{})
	}
	if len(randomSteeringRaw) > 0 {
		randomSteeringMap["default_weight"] = randomSteeringRaw["DefaultWeight"]
		randomSteeringMap["pool_weights"] = randomSteeringRaw["PoolWeights"]

		randomSteeringMaps = append(randomSteeringMaps, randomSteeringMap)
	}
	if err := d.Set("random_steering", randomSteeringMaps); err != nil {
		return err
	}
	rulesRaw := objectRaw["Rules"]
	rulesMaps := make([]map[string]interface{}, 0)
	if rulesRaw != nil {
		for _, rulesChildRaw := range convertToInterfaceArray(rulesRaw) {
			rulesMap := make(map[string]interface{})
			rulesChildRaw := rulesChildRaw.(map[string]interface{})
			if rulesChildRaw["Overrides"] != nil {
				rulesMap["overrides"] = convertObjectToJsonString(rulesChildRaw["Overrides"])
			}
			rulesMap["rule"] = rulesChildRaw["Rule"]
			rulesMap["rule_enable"] = rulesChildRaw["RuleEnable"]
			rulesMap["rule_name"] = rulesChildRaw["RuleName"]
			rulesMap["sequence"] = rulesChildRaw["Sequence"]
			rulesMap["terminates"] = rulesChildRaw["Terminates"]

			fixedResponseMaps := make([]map[string]interface{}, 0)
			fixedResponseMap := make(map[string]interface{})
			fixedResponseRaw := make(map[string]interface{})
			if rulesChildRaw["FixedResponse"] != nil {
				fixedResponseRaw = rulesChildRaw["FixedResponse"].(map[string]interface{})
			}
			if len(fixedResponseRaw) > 0 {
				fixedResponseMap["content_type"] = fixedResponseRaw["ContentType"]
				fixedResponseMap["location"] = fixedResponseRaw["Location"]
				fixedResponseMap["message_body"] = fixedResponseRaw["MessageBody"]
				fixedResponseMap["status_code"] = fixedResponseRaw["StatusCode"]

				fixedResponseMaps = append(fixedResponseMaps, fixedResponseMap)
			}
			rulesMap["fixed_response"] = fixedResponseMaps
			rulesMaps = append(rulesMaps, rulesMap)
		}
	}
	if err := d.Set("rules", rulesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEsaLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateLoadBalancer"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	if d.HasChange("region_pools") {
		update = true
		request["RegionPools"] = d.Get("region_pools")
	}

	if d.HasChange("rules") {
		update = true
		if v, ok := d.GetOk("rules"); ok || d.HasChange("rules") {
			rulesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["RuleEnable"] = dataLoopTmp["rule_enable"]
				dataLoopMap["Rule"] = dataLoopTmp["rule"]

				if overrides, ok := dataLoopTmp["overrides"]; ok && fmt.Sprint(overrides) != "" {
					dataLoopMap["Overrides"] = overrides
				}
				if !IsNil(dataLoopTmp["fixed_response"]) {
					localData1 := make(map[string]interface{})
					messageBody1, _ := jsonpath.Get("$[0].message_body", dataLoopTmp["fixed_response"])
					if messageBody1 != nil && messageBody1 != "" {
						localData1["MessageBody"] = messageBody1
					}
					location1, _ := jsonpath.Get("$[0].location", dataLoopTmp["fixed_response"])
					if location1 != nil && location1 != "" {
						localData1["Location"] = location1
					}
					statusCode1, _ := jsonpath.Get("$[0].status_code", dataLoopTmp["fixed_response"])
					if statusCode1 != nil && statusCode1 != "" {
						localData1["StatusCode"] = statusCode1
					}
					contentType1, _ := jsonpath.Get("$[0].content_type", dataLoopTmp["fixed_response"])
					if contentType1 != nil && contentType1 != "" {
						localData1["ContentType"] = contentType1
					}
					dataLoopMap["FixedResponse"] = localData1
				}
				dataLoopMap["Sequence"] = dataLoopTmp["sequence"]
				dataLoopMap["Terminates"] = dataLoopTmp["terminates"]
				dataLoopMap["RuleName"] = dataLoopTmp["rule_name"]
				rulesMapsArray = append(rulesMapsArray, dataLoopMap)
			}
			rulesMapsJson, err := json.Marshal(rulesMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["Rules"] = string(rulesMapsJson)
		}
	}

	if d.HasChange("monitor") {
		update = true
	}
	dataList := make(map[string]interface{})

	if v := d.Get("monitor"); v != nil {
		path1, _ := jsonpath.Get("$[0].path", v)
		if path1 != nil && (d.HasChange("monitor.0.path") || path1 != "") {
			dataList["Path"] = path1
		}
		expectedCodes1, _ := jsonpath.Get("$[0].expected_codes", v)
		if expectedCodes1 != nil && (d.HasChange("monitor.0.expected_codes") || expectedCodes1 != "") {
			dataList["ExpectedCodes"] = expectedCodes1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && (d.HasChange("monitor.0.type") || type1 != "") {
			dataList["Type"] = type1
		}
		consecutiveUp1, _ := jsonpath.Get("$[0].consecutive_up", v)
		if consecutiveUp1 != nil && (d.HasChange("monitor.0.consecutive_up") || consecutiveUp1 != "") {
			dataList["ConsecutiveUp"] = consecutiveUp1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && (d.HasChange("monitor.0.port") || port1 != "") {
			dataList["Port"] = port1
		}
		monitoringRegion1, _ := jsonpath.Get("$[0].monitoring_region", v)
		if monitoringRegion1 != nil && (d.HasChange("monitor.0.monitoring_region") || monitoringRegion1 != "") {
			dataList["MonitoringRegion"] = monitoringRegion1
		}
		header1, _ := jsonpath.Get("$[0].header", v)
		if header1 != nil && (d.HasChange("monitor.0.header") || header1 != "") {
			dataList["Header"] = header1
		}
		method1, _ := jsonpath.Get("$[0].method", v)
		if method1 != nil && (d.HasChange("monitor.0.method") || method1 != "") {
			dataList["Method"] = method1
		}
		followRedirects1, _ := jsonpath.Get("$[0].follow_redirects", v)
		if followRedirects1 != nil && (d.HasChange("monitor.0.follow_redirects") || followRedirects1 != "") {
			dataList["FollowRedirects"] = followRedirects1
		}
		consecutiveDown1, _ := jsonpath.Get("$[0].consecutive_down", v)
		if consecutiveDown1 != nil && (d.HasChange("monitor.0.consecutive_down") || consecutiveDown1 != "") {
			dataList["ConsecutiveDown"] = consecutiveDown1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && (d.HasChange("monitor.0.interval") || interval1 != "") {
			dataList["Interval"] = interval1
		}
		timeout1, _ := jsonpath.Get("$[0].timeout", v)
		if timeout1 != nil && (d.HasChange("monitor.0.timeout") || timeout1 != "") && timeout1.(int) > 0 {
			dataList["Timeout"] = timeout1
		}

		dataListJson, err := json.Marshal(dataList)
		if err != nil {
			return WrapError(err)
		}
		request["Monitor"] = string(dataListJson)
	}

	if d.HasChange("enabled") {
		update = true
		request["Enabled"] = d.Get("enabled")
	}

	if d.HasChange("session_affinity") {
		update = true
		request["SessionAffinity"] = d.Get("session_affinity")
	}

	if d.HasChange("ttl") {
		update = true
		request["Ttl"] = d.Get("ttl")
	}

	if d.HasChange("random_steering") {
		update = true
		dataList1 := make(map[string]interface{})

		if v := d.Get("random_steering"); v != nil {
			poolWeights1, _ := jsonpath.Get("$[0].pool_weights", v)
			if poolWeights1 != nil && (d.HasChange("random_steering.0.pool_weights") || poolWeights1 != "") {
				dataList1["PoolWeights"] = poolWeights1
			}
			defaultWeight1, _ := jsonpath.Get("$[0].default_weight", v)
			if defaultWeight1 != nil && (d.HasChange("random_steering.0.default_weight") || defaultWeight1 != "") {
				dataList1["DefaultWeight"] = defaultWeight1
			}

			dataList1Json, err := json.Marshal(dataList1)
			if err != nil {
				return WrapError(err)
			}
			request["RandomSteering"] = string(dataList1Json)
		}
	}

	if d.HasChange("default_pools") {
		update = true
	}
	if v, ok := d.GetOk("default_pools"); ok || d.HasChange("default_pools") {
		defaultPoolsMapsArray := convertToInterfaceArray(v)

		defaultPoolsMapsJson, err := json.Marshal(defaultPoolsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DefaultPools"] = string(defaultPoolsMapsJson)
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("adaptive_routing") {
		update = true
		dataList2 := make(map[string]interface{})

		if v := d.Get("adaptive_routing"); v != nil {
			failoverAcrossPools1, _ := jsonpath.Get("$[0].failover_across_pools", v)
			if failoverAcrossPools1 != nil && (d.HasChange("adaptive_routing.0.failover_across_pools") || failoverAcrossPools1 != "") {
				dataList2["FailoverAcrossPools"] = failoverAcrossPools1
			}

			dataList2Json, err := json.Marshal(dataList2)
			if err != nil {
				return WrapError(err)
			}
			request["AdaptiveRouting"] = string(dataList2Json)
		}
	}

	if d.HasChange("steering_policy") {
		update = true
	}
	request["SteeringPolicy"] = d.Get("steering_policy")
	if d.HasChange("fallback_pool") {
		update = true
	}
	request["FallbackPool"] = d.Get("fallback_pool")
	if d.HasChange("sub_region_pools") {
		update = true
		request["SubRegionPools"] = d.Get("sub_region_pools")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"LockFailed"}) || NeedRetry(err) {
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

	return resourceAliCloudEsaLoadBalancerRead(d, meta)
}

func resourceAliCloudEsaLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"LockFailed"}) || NeedRetry(err) {
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
