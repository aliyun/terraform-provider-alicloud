// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenTrafficMarkingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTrafficMarkingPolicyCreate,
		Read:   resourceAliCloudCenTrafficMarkingPolicyRead,
		Update: resourceAliCloudCenTrafficMarkingPolicyUpdate,
		Delete: resourceAliCloudCenTrafficMarkingPolicyDelete,
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"marking_dscp": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_marking_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_marking_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"traffic_match_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"traffic_match_rule_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"match_dscp": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 63),
						},
						"dst_port_range": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"src_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"traffic_match_rule_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"address_family": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"src_port_range": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"dst_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCenTrafficMarkingPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrafficMarkingPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["TransitRouterId"] = d.Get("transit_router_id")

	request["ClientToken"] = buildClientToken(action)

	request["Priority"] = d.Get("priority")
	request["MarkingDscp"] = d.Get("marking_dscp")
	if v, ok := d.GetOk("traffic_marking_policy_name"); ok {
		request["TrafficMarkingPolicyName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["TrafficMarkingPolicyDescription"] = v
	}
	if v, ok := d.GetOk("traffic_match_rules"); ok {
		trafficMatchRulesMaps := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
			dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
			dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
			dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
			dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
			dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
			dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
			dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
			dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
			trafficMatchRulesMaps = append(trafficMatchRulesMaps, dataLoopMap)
		}
		request["TrafficMatchRules"] = trafficMatchRulesMaps
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_traffic_marking_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["TransitRouterId"], response["TrafficMarkingPolicyId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTrafficMarkingPolicyRead(d, meta)
}

func resourceAliCloudCenTrafficMarkingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTrafficMarkingPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_traffic_marking_policy DescribeCenTrafficMarkingPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["TrafficMarkingPolicyDescription"] != nil {
		d.Set("description", objectRaw["TrafficMarkingPolicyDescription"])
	}
	if objectRaw["MarkingDscp"] != nil {
		d.Set("marking_dscp", objectRaw["MarkingDscp"])
	}
	if objectRaw["Priority"] != nil {
		d.Set("priority", objectRaw["Priority"])
	}
	if objectRaw["TrafficMarkingPolicyStatus"] != nil {
		d.Set("status", objectRaw["TrafficMarkingPolicyStatus"])
	}
	if objectRaw["TrafficMarkingPolicyName"] != nil {
		d.Set("traffic_marking_policy_name", objectRaw["TrafficMarkingPolicyName"])
	}
	if objectRaw["TrafficMarkingPolicyId"] != nil {
		d.Set("traffic_marking_policy_id", objectRaw["TrafficMarkingPolicyId"])
	}
	if objectRaw["TransitRouterId"] != nil {
		d.Set("transit_router_id", objectRaw["TransitRouterId"])
	}

	trafficMatchRules1Raw := objectRaw["TrafficMatchRules"]
	trafficMatchRulesMaps := make([]map[string]interface{}, 0)
	if trafficMatchRules1Raw != nil {
		for _, trafficMatchRulesChild1Raw := range trafficMatchRules1Raw.([]interface{}) {
			trafficMatchRulesMap := make(map[string]interface{})
			trafficMatchRulesChild1Raw := trafficMatchRulesChild1Raw.(map[string]interface{})
			trafficMatchRulesMap["address_family"] = trafficMatchRulesChild1Raw["AddressFamily"]
			trafficMatchRulesMap["dst_cidr"] = trafficMatchRulesChild1Raw["DstCidr"]
			trafficMatchRulesMap["match_dscp"] = trafficMatchRulesChild1Raw["MatchDscp"]
			trafficMatchRulesMap["protocol"] = trafficMatchRulesChild1Raw["Protocol"]
			trafficMatchRulesMap["src_cidr"] = trafficMatchRulesChild1Raw["SrcCidr"]
			trafficMatchRulesMap["traffic_match_rule_description"] = trafficMatchRulesChild1Raw["TrafficMatchRuleDescription"]
			trafficMatchRulesMap["traffic_match_rule_name"] = trafficMatchRulesChild1Raw["TrafficMatchRuleName"]

			dstPortRange1Raw := make([]interface{}, 0)
			if trafficMatchRulesChild1Raw["DstPortRange"] != nil {
				dstPortRange1Raw = trafficMatchRulesChild1Raw["DstPortRange"].([]interface{})
			}

			trafficMatchRulesMap["dst_port_range"] = dstPortRange1Raw
			srcPortRange1Raw := make([]interface{}, 0)
			if trafficMatchRulesChild1Raw["SrcPortRange"] != nil {
				srcPortRange1Raw = trafficMatchRulesChild1Raw["SrcPortRange"].([]interface{})
			}

			trafficMatchRulesMap["src_port_range"] = srcPortRange1Raw
			trafficMatchRulesMaps = append(trafficMatchRulesMaps, trafficMatchRulesMap)
		}
	}
	if objectRaw["TrafficMatchRules"] != nil {
		if err := d.Set("traffic_match_rules", trafficMatchRulesMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("transit_router_id", parts[0])
	d.Set("traffic_marking_policy_id", parts[1])

	return nil
}

func resourceAliCloudCenTrafficMarkingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateTrafficMarkingPolicyAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TrafficMarkingPolicyId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("traffic_marking_policy_name") {
		update = true
		request["TrafficMarkingPolicyName"] = d.Get("traffic_marking_policy_name")
	}

	if d.HasChange("description") {
		update = true
		request["TrafficMarkingPolicyDescription"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cenServiceV2 := CenServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("traffic_match_rules") {
		oldEntry, newEntry := d.GetChange("traffic_match_rules")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UpdateTrafficMarkingPolicyAttribute"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["TrafficMarkingPolicyId"] = parts[1]

			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			deleteTrafficMatchRulesMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
				dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
				dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
				dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
				dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
				dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
				dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
				dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
				dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
				deleteTrafficMatchRulesMaps = append(deleteTrafficMatchRulesMaps, dataLoopMap)
			}
			request["DeleteTrafficMatchRules"] = deleteTrafficMatchRulesMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.TrafficMarkingPolicy"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			cenServiceV2 := CenServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UpdateTrafficMarkingPolicyAttribute"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["TrafficMarkingPolicyId"] = parts[1]

			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			addTrafficMatchRulesMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
				dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
				dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
				dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
				dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
				dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
				dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
				dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
				dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
				addTrafficMatchRulesMaps = append(addTrafficMatchRulesMaps, dataLoopMap)
			}
			request["AddTrafficMatchRules"] = addTrafficMatchRulesMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.TrafficMarkingPolicy"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			cenServiceV2 := CenServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	return resourceAliCloudCenTrafficMarkingPolicyRead(d, meta)
}

func resourceAliCloudCenTrafficMarkingPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTrafficMarkingPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["TrafficMarkingPolicyId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
