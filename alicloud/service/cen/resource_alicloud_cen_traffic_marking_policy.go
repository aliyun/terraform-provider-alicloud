// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cen

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCenTrafficMarkingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudCenTrafficMarkingPolicyCreate,
		Read:   ResourceAliCloudCenTrafficMarkingPolicyRead,
		Update: ResourceAliCloudCenTrafficMarkingPolicyUpdate,
		Delete: ResourceAliCloudCenTrafficMarkingPolicyDelete,
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
							ValidateFunc: validators.IntBetween(0, 63),
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

func ResourceAliCloudCenTrafficMarkingPolicyCreate(d *schema.ResourceData, meta interface{}) error {

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

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("traffic_match_rules"); ok {
		trafficMatchRulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_traffic_marking_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterId"], response["TrafficMarkingPolicyId"]))

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAliCloudCenTrafficMarkingPolicyRead(d, meta)
}

func ResourceAliCloudCenTrafficMarkingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := NewCenServiceV2(client)

	objectRaw, err := cenServiceV2.DescribeCenTrafficMarkingPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_traffic_marking_policy DescribeCenTrafficMarkingPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
		for _, trafficMatchRulesChildRaw := range convert.ConvertToInterfaceArray(trafficMatchRulesRaw) {
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
				dstPortRangeRaw = convert.ConvertToInterfaceArray(trafficMatchRulesChildRaw["DstPortRange"])
			}

			trafficMatchRulesMap["dst_port_range"] = dstPortRangeRaw
			srcPortRangeRaw := make([]interface{}, 0)
			if trafficMatchRulesChildRaw["SrcPortRange"] != nil {
				srcPortRangeRaw = convert.ConvertToInterfaceArray(trafficMatchRulesChildRaw["SrcPortRange"])
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

func ResourceAliCloudCenTrafficMarkingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
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

	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["TrafficMarkingPolicyDescription"] = d.Get("description")
	}

	if d.HasChange("traffic_marking_policy_name") {
		update = true
		request["TrafficMarkingPolicyName"] = d.Get("traffic_marking_policy_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		cenServiceV2 := NewCenServiceV2(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

			request["ClientToken"] = helper.BuildClientToken(action)
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

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficMarkingPolicy", "Throttling.User"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			cenServiceV2 := NewCenServiceV2(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UpdateTrafficMarkingPolicyAttribute"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["TrafficMarkingPolicyId"] = parts[1]

			request["ClientToken"] = helper.BuildClientToken(action)
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

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TrafficMarkingPolicy", "Throttling.User"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			cenServiceV2 := NewCenServiceV2(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}

	}
	return ResourceAliCloudCenTrafficMarkingPolicyRead(d, meta)
}

func ResourceAliCloudCenTrafficMarkingPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTrafficMarkingPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TrafficMarkingPolicyId"] = parts[1]

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTrafficMarkingPolicyStateRefreshFunc(d.Id(), "TrafficMarkingPolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
