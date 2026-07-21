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

func resourceAliCloudMessageServiceEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceEventRuleCreate,
		Read:   resourceAliCloudMessageServiceEventRuleRead,
		Delete: resourceAliCloudMessageServiceEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"delivery_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"DIRECT", "BROADCAST"}, false),
			},
			"endpoint": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"endpoint_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"event_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_rules": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"suffix": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
							},
							"match_state": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
							},
							"prefix": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
							},
							"name": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
							},
						},
					},
				},
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMessageServiceEventRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEventRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["ProductName"] = "oss"
	if v, ok := d.GetOk("delivery_mode"); ok {
		request["DeliveryMode"] = v
	}
	if v, ok := d.GetOk("event_types"); ok {
		eventTypesMapsArray := v.([]interface{})
		eventTypesMapsJson, err := json.Marshal(eventTypesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["EventTypes"] = string(eventTypesMapsJson)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("endpoint"); !IsNil(v) {
		endpointType1, _ := jsonpath.Get("$[0].endpoint_type", v)
		if endpointType1 != nil && endpointType1 != "" {
			objectDataLocalMap["EndpointType"] = endpointType1
		}
		endpointValue1, _ := jsonpath.Get("$[0].endpoint_value", v)
		if endpointValue1 != nil && endpointValue1 != "" {
			objectDataLocalMap["EndpointValue"] = endpointValue1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Endpoint"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("match_rules"); ok {
		matchRulesMapsArray := make([][]map[string]interface{}, 0)
		for _, outerList := range v.([]interface{}) {
			innerList := make([]map[string]interface{}, 0)
			for _, innerItem := range outerList.([]interface{}) {
				dataMap := innerItem.(map[string]interface{})
				ruleMap := make(map[string]interface{})
				if val, exist := dataMap["name"]; exist {
					ruleMap["Name"] = val
				}
				if val, exist := dataMap["prefix"]; exist {
					ruleMap["Prefix"] = val
				}
				if val, exist := dataMap["suffix"]; exist {
					ruleMap["Suffix"] = val
				}
				if val, exist := dataMap["match_state"]; exist {
					ruleMap["MatchState"] = val
				}
				innerList = append(innerList, ruleMap)
			}
			matchRulesMapsArray = append(matchRulesMapsArray, innerList)
		}
		matchRulesJson, err := json.Marshal(matchRulesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["MatchRules"] = string(matchRulesJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_event_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleName"]))

	return resourceAliCloudMessageServiceEventRuleRead(d, meta)
}

func resourceAliCloudMessageServiceEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceEventRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_event_rule DescribeMessageServiceEventRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("delivery_mode", objectRaw["DeliveryMode"])
	d.Set("rule_name", objectRaw["RuleName"])

	endpointMaps := make([]map[string]interface{}, 0)
	endpointMap := make(map[string]interface{})
	endpointRaw := make(map[string]interface{})
	if objectRaw["Endpoint"] != nil {
		endpointRaw = objectRaw["Endpoint"].(map[string]interface{})
	}
	if len(endpointRaw) > 0 {
		endpointMap["endpoint_type"] = endpointRaw["EndpointType"]
		parts := strings.Split(fmt.Sprint(endpointRaw["EndpointValue"]), "/")
		if endpointMap["endpoint_type"] == "queue" {
			if len(parts) > 1 {
				endpointMap["endpoint_value"] = parts[len(parts)-1]
			} else {
				endpointMap["endpoint_value"] = endpointRaw["EndpointValue"]
			}
		} else {
			endpointMap["endpoint_value"] = endpointRaw["EndpointValue"]
		}
		endpointMaps = append(endpointMaps, endpointMap)
	}
	if err := d.Set("endpoint", endpointMaps); err != nil {
		return err
	}
	eventTypesRaw := make([]interface{}, 0)
	if objectRaw["EventTypes"] != nil {
		eventTypesRaw = objectRaw["EventTypes"].([]interface{})
	}

	d.Set("event_types", eventTypesRaw)

	matchRulesChildRaw := objectRaw["MatchRules"]
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
	if err := d.Set("match_rules", matchRulesMaps); err != nil {
		return err
	}

	d.Set("rule_name", d.Id())

	return nil
}

func resourceAliCloudMessageServiceEventRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RuleName"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ProductName"] = "oss"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

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
