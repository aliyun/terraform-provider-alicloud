// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
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
			"endpoints": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"endpoint_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"http", "queue", "topic"}, false),
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
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeList},
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

	request["ProductName"] = "oss"
	if v, ok := d.GetOk("event_types"); ok {
		eventTypesMapsArray := v.([]interface{})
		eventTypesMapsJson, err := json.Marshal(eventTypesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["EventTypes"] = string(eventTypesMapsJson)
	}

	if v, ok := d.GetOk("endpoints"); ok {
		endpointsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["EndpointType"] = dataLoop1Tmp["endpoint_type"]
			dataLoop1Map["EndpointValue"] = dataLoop1Tmp["endpoint_value"]
			endpointsMapsArray = append(endpointsMapsArray, dataLoop1Map)
		}
		endpointsMapsJson, err := json.Marshal(endpointsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Endpoints"] = string(endpointsMapsJson)
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

	d.Set("rule_name", objectRaw["RuleName"])

	endpointsRaw := objectRaw["Endpoints"]
	endpointsMaps := make([]map[string]interface{}, 0)
	if endpointsRaw != nil {
		for _, endpointsChildRaw := range endpointsRaw.([]interface{}) {
			endpointsMap := make(map[string]interface{})
			endpointsChild := endpointsChildRaw.(map[string]interface{})
			endpointsMap["endpoint_value"] = endpointsChild["EndpointValue"]
			endpointsMap["endpoint_type"] = endpointsChild["EndpointType"]
			endpointsMaps = append(endpointsMaps, endpointsMap)
		}
	}
	if err := d.Set("endpoints", endpointsRaw); err != nil {
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
				matchRulesMap["match_state"] = matchRulesChildChildChild["MatchState"]
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

	e := jsonata.MustCompile("$contains($.Endpoints.EndpointType, \"queue\") and $contains($.TopicName, \"mns-en-topics-oss\") \n    ? \"queue\" :  $contains($.Endpoints.EndpointType, \"queue\") ? \"topic\" : \"http\"\n")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("endpoints", evaluation)
	e = jsonata.MustCompile("$split($.Endpoints.EndpointValue, \"/queues/\")[-1]")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("endpoints", evaluation)

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
