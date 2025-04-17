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
			Create: schema.DefaultTimeout(8 * time.Minute),
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
			"force": {
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
	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("traffic_match_rules"); ok {
		trafficMatchRulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
			dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
			dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
			dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
			dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
			dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
			dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
			dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
			dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
			trafficMatchRulesMapsArray = append(trafficMatchRulesMapsArray, dataLoopMap)
		}
		request["TrafficMatchRules"] = trafficMatchRulesMapsArray
	}

	if v, ok := d.GetOk("description"); ok {
		request["TrafficMarkingPolicyDescription"] = v
	}
	request["Priority"] = d.Get("priority")
	request["MarkingDscp"] = d.Get("marking_dscp")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("traffic_marking_policy_name"); ok {
		request["TrafficMarkingPolicyName"] = v
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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_traffic_marking_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterId"], response["TrafficMarkingPolicyId"]))

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

	d.Set("description", objectRaw["TrafficMarkingPolicyDescription"])
	d.Set("marking_dscp", objectRaw["MarkingDscp"])
	d.Set("priority", objectRaw["Priority"])
	d.Set("status", objectRaw["TrafficMarkingPolicyStatus"])
	d.Set("traffic_marking_policy_name", objectRaw["TrafficMarkingPolicyName"])
	d.Set("traffic_marking_policy_id", objectRaw["TrafficMarkingPolicyId"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])

	trafficMatchRulesRaw := objectRaw["TrafficMatchRules"]
	trafficMatchRulesMaps := make([]map[string]interface{}, 0)
	if trafficMatchRulesRaw != nil {
		for _, trafficMatchRulesChildRaw := range trafficMatchRulesRaw.([]interface{}) {
			trafficMatchRulesMap := make(map[string]interface{})
			trafficMatchRulesChildRaw := trafficMatchRulesChildRaw.(map[string]interface{})
			trafficMatchRulesMap["address_family"] = trafficMatchRulesChildRaw["AddressFamily"]
			trafficMatchRulesMap["dst_cidr"] = trafficMatchRulesChildRaw["DstCidr"]
			trafficMatchRulesMap["match_dscp"] = trafficMatchRulesChildRaw["MatchDscp"]
			trafficMatchRulesMap["protocol"] = trafficMatchRulesChildRaw["Protocol"]
			trafficMatchRulesMap["src_cidr"] = trafficMatchRulesChildRaw["SrcCidr"]
			trafficMatchRulesMap["traffic_match_rule_description"] = trafficMatchRulesChildRaw["TrafficMatchRuleDescription"]
			trafficMatchRulesMap["traffic_match_rule_name"] = trafficMatchRulesChildRaw["TrafficMatchRuleName"]

			dstPortRangeRaw := make([]interface{}, 0)
			if trafficMatchRulesChildRaw["DstPortRange"] != nil {
				dstPortRangeRaw = trafficMatchRulesChildRaw["DstPortRange"].([]interface{})
			}

			trafficMatchRulesMap["dst_port_range"] = dstPortRangeRaw
			srcPortRangeRaw := make([]interface{}, 0)
			if trafficMatchRulesChildRaw["SrcPortRange"] != nil {
				srcPortRangeRaw = trafficMatchRulesChildRaw["SrcPortRange"].([]interface{})
			}

			trafficMatchRulesMap["src_port_range"] = srcPortRangeRaw
			trafficMatchRulesMaps = append(trafficMatchRulesMaps, trafficMatchRulesMap)
		}
	}
	if err := d.Set("traffic_match_rules", trafficMatchRulesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCenTrafficMarkingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateTrafficMarkingPolicyAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TrafficMarkingPolicyId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["TrafficMarkingPolicyDescription"] = d.Get("description")
	}

	if d.HasChange("traffic_marking_policy_name") {
		update = true
		request["TrafficMarkingPolicyName"] = d.Get("traffic_marking_policy_name")
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
			return nil
		})
		addDebug(action, response, request)
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
			request["TrafficMarkingPolicyId"] = parts[1]

			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			deleteTrafficMatchRulesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
				dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
				dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
				dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
				dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
				dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
				dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
				dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
				dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
				deleteTrafficMatchRulesMapsArray = append(deleteTrafficMatchRulesMapsArray, dataLoopMap)
			}
			request["DeleteTrafficMatchRules"] = deleteTrafficMatchRulesMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficMarkingPolicy", "Throttling.User"}) || NeedRetry(err) {
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
			request["TrafficMarkingPolicyId"] = parts[1]

			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			addTrafficMatchRulesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["SrcPortRange"] = dataLoopTmp["src_port_range"]
				dataLoopMap["TrafficMatchRuleDescription"] = dataLoopTmp["traffic_match_rule_description"]
				dataLoopMap["MatchDscp"] = dataLoopTmp["match_dscp"]
				dataLoopMap["AddressFamily"] = dataLoopTmp["address_family"]
				dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
				dataLoopMap["TrafficMatchRuleName"] = dataLoopTmp["traffic_match_rule_name"]
				dataLoopMap["SrcCidr"] = dataLoopTmp["src_cidr"]
				dataLoopMap["DstCidr"] = dataLoopTmp["dst_cidr"]
				dataLoopMap["DstPortRange"] = dataLoopTmp["dst_port_range"]
				addTrafficMatchRulesMapsArray = append(addTrafficMatchRulesMapsArray, dataLoopMap)
			}
			request["AddTrafficMatchRules"] = addTrafficMatchRulesMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficMarkingPolicy", "Throttling.User"}) || NeedRetry(err) {
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
	request["TrafficMarkingPolicyId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
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
