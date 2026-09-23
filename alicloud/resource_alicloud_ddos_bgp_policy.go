// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudDdosBgpPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosBgpPolicyCreate,
		Read:   resourceAliCloudDdosBgpPolicyRead,
		Update: resourceAliCloudDdosBgpPolicyUpdate,
		Delete: resourceAliCloudDdosBgpPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"whiten_gfbr_nets": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"port_rule_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"src_port_end": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"port_rule_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"src_port_start": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"seq_no": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"dst_port_start": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"dst_port_end": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"match_action": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"region_block_country_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"region_block_province_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"finger_print_rule_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"src_port_end": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"finger_print_rule_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"src_port_start": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"seq_no": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"payload_bytes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"rate_value": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"dst_port_start": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"match_action": {
										Type:     schema.TypeString,
										Required: true,
									},
									"offset": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"max_pkt_len": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"dst_port_end": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"min_pkt_len": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"source_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pps": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"syn_bps": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"bps": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"syn_pps": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"enable_intelligence": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_defense": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_block_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: IntBetween(3, 6),
									},
									"exceed_limit_times": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"every_seconds": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"block_expire_seconds": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"enable_drop_icmp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"intelligence_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"reflect_block_udp_port_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"black_ip_list_expire_at": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"layer4_rule_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},
									"condition_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"position": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"arg": {
													Type:     schema.TypeString,
													Required: true,
												},
												"depth": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"priority": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"method": {
										Type:     schema.TypeString,
										Required: true,
									},
									"limited": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"match": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudDdosBgpPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Name"] = d.Get("policy_name")
	request["Type"] = d.Get("type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddos_bgp_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudDdosBgpPolicyUpdate(d, meta)
}

func resourceAliCloudDdosBgpPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosBgpServiceV2 := DdosBgpServiceV2{client}

	objectRaw, err := ddosBgpServiceV2.DescribeDdosBgpPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddos_bgp_policy DescribeDdosBgpPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Name"] != nil {
		d.Set("policy_name", objectRaw["Name"])
	}
	if objectRaw["Type"] != nil {
		d.Set("type", objectRaw["Type"])
	}

	contentMaps := make([]map[string]interface{}, 0)
	contentMap := make(map[string]interface{})
	content1Raw := make(map[string]interface{})
	if objectRaw["Content"] != nil {
		content1Raw = objectRaw["Content"].(map[string]interface{})
	}
	if len(content1Raw) > 0 {
		contentMap["black_ip_list_expire_at"] = content1Raw["BlackIpListExpireAt"]
		contentMap["enable_defense"] = content1Raw["EnableL4Defense"]
		contentMap["enable_drop_icmp"] = content1Raw["EnableDropIcmp"]
		contentMap["enable_intelligence"] = content1Raw["EnableIntelligence"]
		contentMap["intelligence_level"] = content1Raw["IntelligenceLevel"]
		contentMap["whiten_gfbr_nets"] = content1Raw["WhitenGfbrNets"]

		fingerPrintRuleList1Raw := content1Raw["FingerPrintRuleList"]
		fingerPrintRuleListMaps := make([]map[string]interface{}, 0)
		if fingerPrintRuleList1Raw != nil {
			for _, fingerPrintRuleListChild1Raw := range fingerPrintRuleList1Raw.([]interface{}) {
				fingerPrintRuleListMap := make(map[string]interface{})
				fingerPrintRuleListChild1Raw := fingerPrintRuleListChild1Raw.(map[string]interface{})
				fingerPrintRuleListMap["dst_port_end"] = fingerPrintRuleListChild1Raw["DstPortEnd"]
				fingerPrintRuleListMap["dst_port_start"] = fingerPrintRuleListChild1Raw["DstPortStart"]
				fingerPrintRuleListMap["finger_print_rule_id"] = fingerPrintRuleListChild1Raw["Id"]
				fingerPrintRuleListMap["match_action"] = fingerPrintRuleListChild1Raw["MatchAction"]
				fingerPrintRuleListMap["max_pkt_len"] = fingerPrintRuleListChild1Raw["MaxPktLen"]
				fingerPrintRuleListMap["min_pkt_len"] = fingerPrintRuleListChild1Raw["MinPktLen"]
				fingerPrintRuleListMap["offset"] = fingerPrintRuleListChild1Raw["Offset"]
				fingerPrintRuleListMap["payload_bytes"] = fingerPrintRuleListChild1Raw["PayloadBytes"]
				fingerPrintRuleListMap["protocol"] = fingerPrintRuleListChild1Raw["Protocol"]
				fingerPrintRuleListMap["rate_value"] = fingerPrintRuleListChild1Raw["RateValue"]
				fingerPrintRuleListMap["seq_no"] = fingerPrintRuleListChild1Raw["SeqNo"]
				fingerPrintRuleListMap["src_port_end"] = fingerPrintRuleListChild1Raw["SrcPortEnd"]
				fingerPrintRuleListMap["src_port_start"] = fingerPrintRuleListChild1Raw["SrcPortStart"]

				fingerPrintRuleListMaps = append(fingerPrintRuleListMaps, fingerPrintRuleListMap)
			}
		}
		contentMap["finger_print_rule_list"] = fingerPrintRuleListMaps
		l4RuleList1Raw := content1Raw["L4RuleList"]
		layer4RuleListMaps := make([]map[string]interface{}, 0)
		if l4RuleList1Raw != nil {
			for _, l4RuleListChild1Raw := range l4RuleList1Raw.([]interface{}) {
				layer4RuleListMap := make(map[string]interface{})
				l4RuleListChild1Raw := l4RuleListChild1Raw.(map[string]interface{})
				layer4RuleListMap["action"] = l4RuleListChild1Raw["Action"]
				layer4RuleListMap["limited"] = l4RuleListChild1Raw["Limited"]
				layer4RuleListMap["match"] = l4RuleListChild1Raw["Match"]
				layer4RuleListMap["method"] = l4RuleListChild1Raw["Method"]
				layer4RuleListMap["name"] = l4RuleListChild1Raw["Name"]
				layer4RuleListMap["priority"] = l4RuleListChild1Raw["Priority"]

				conditionList1Raw := l4RuleListChild1Raw["ConditionList"]
				conditionListMaps := make([]map[string]interface{}, 0)
				if conditionList1Raw != nil {
					for _, conditionListChild1Raw := range conditionList1Raw.([]interface{}) {
						conditionListMap := make(map[string]interface{})
						conditionListChild1Raw := conditionListChild1Raw.(map[string]interface{})
						conditionListMap["arg"] = conditionListChild1Raw["Arg"]
						conditionListMap["depth"] = conditionListChild1Raw["Depth"]
						conditionListMap["position"] = conditionListChild1Raw["Position"]

						conditionListMaps = append(conditionListMaps, conditionListMap)
					}
				}
				layer4RuleListMap["condition_list"] = conditionListMaps
				layer4RuleListMaps = append(layer4RuleListMaps, layer4RuleListMap)
			}
		}
		contentMap["layer4_rule_list"] = layer4RuleListMaps
		portRuleList1Raw := content1Raw["PortRuleList"]
		portRuleListMaps := make([]map[string]interface{}, 0)
		if portRuleList1Raw != nil {
			for _, portRuleListChild1Raw := range portRuleList1Raw.([]interface{}) {
				portRuleListMap := make(map[string]interface{})
				portRuleListChild1Raw := portRuleListChild1Raw.(map[string]interface{})
				portRuleListMap["dst_port_end"] = portRuleListChild1Raw["DstPortEnd"]
				portRuleListMap["dst_port_start"] = portRuleListChild1Raw["DstPortStart"]
				portRuleListMap["match_action"] = portRuleListChild1Raw["MatchAction"]
				portRuleListMap["port_rule_id"] = portRuleListChild1Raw["Id"]
				portRuleListMap["protocol"] = portRuleListChild1Raw["Protocol"]
				portRuleListMap["seq_no"] = portRuleListChild1Raw["SeqNo"]
				portRuleListMap["src_port_end"] = portRuleListChild1Raw["SrcPortEnd"]
				portRuleListMap["src_port_start"] = portRuleListChild1Raw["SrcPortStart"]

				portRuleListMaps = append(portRuleListMaps, portRuleListMap)
			}
		}
		contentMap["port_rule_list"] = portRuleListMaps
		reflectBlockUdpPortList1Raw := make([]interface{}, 0)
		if content1Raw["ReflectBlockUdpPortList"] != nil {
			reflectBlockUdpPortList1Raw = content1Raw["ReflectBlockUdpPortList"].([]interface{})
		}

		contentMap["reflect_block_udp_port_list"] = reflectBlockUdpPortList1Raw
		regionBlockCountryList1Raw := make([]interface{}, 0)
		if content1Raw["RegionBlockCountryList"] != nil {
			regionBlockCountryList1Raw = content1Raw["RegionBlockCountryList"].([]interface{})
		}

		contentMap["region_block_country_list"] = regionBlockCountryList1Raw
		regionBlockProvinceList1Raw := make([]interface{}, 0)
		if content1Raw["RegionBlockProvinceList"] != nil {
			regionBlockProvinceList1Raw = content1Raw["RegionBlockProvinceList"].([]interface{})
		}

		contentMap["region_block_province_list"] = regionBlockProvinceList1Raw
		sourceBlockList1Raw := content1Raw["SourceBlockList"]
		sourceBlockListMaps := make([]map[string]interface{}, 0)
		if sourceBlockList1Raw != nil {
			for _, sourceBlockListChild1Raw := range sourceBlockList1Raw.([]interface{}) {
				sourceBlockListMap := make(map[string]interface{})
				sourceBlockListChild1Raw := sourceBlockListChild1Raw.(map[string]interface{})
				sourceBlockListMap["block_expire_seconds"] = sourceBlockListChild1Raw["BlockExpireSeconds"]
				sourceBlockListMap["every_seconds"] = sourceBlockListChild1Raw["EverySeconds"]
				sourceBlockListMap["exceed_limit_times"] = sourceBlockListChild1Raw["ExceedLimitTimes"]
				sourceBlockListMap["type"] = sourceBlockListChild1Raw["Type"]

				sourceBlockListMaps = append(sourceBlockListMaps, sourceBlockListMap)
			}
		}
		contentMap["source_block_list"] = sourceBlockListMaps
		sourceLimitMaps := make([]map[string]interface{}, 0)
		sourceLimitMap := make(map[string]interface{})
		sourceLimit1Raw := make(map[string]interface{})
		if content1Raw["SourceLimit"] != nil {
			sourceLimit1Raw = content1Raw["SourceLimit"].(map[string]interface{})
		}
		if len(sourceLimit1Raw) > 0 {
			sourceLimitMap["bps"] = sourceLimit1Raw["Bps"]
			sourceLimitMap["pps"] = sourceLimit1Raw["Pps"]
			sourceLimitMap["syn_bps"] = sourceLimit1Raw["SynBps"]
			sourceLimitMap["syn_pps"] = sourceLimit1Raw["SynPps"]

			sourceLimitMaps = append(sourceLimitMaps, sourceLimitMap)
		}
		contentMap["source_limit"] = sourceLimitMaps
		contentMaps = append(contentMaps, contentMap)
	}
	if objectRaw["Content"] != nil {
		d.Set("content", contentMaps)
	}

	return nil
}

func resourceAliCloudDdosBgpPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyPolicyContent"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()
	if d.HasChange("content") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("content"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].black_ip_list_expire_at", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["BlackIpListExpireAt"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].enable_intelligence", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["EnableIntelligence"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].intelligence_level", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["IntelligenceLevel"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].whiten_gfbr_nets", v)
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["WhitenGfbrNets"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].enable_drop_icmp", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			objectDataLocalMap["EnableDropIcmp"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].region_block_country_list", d.Get("content"))
		if nodeNative5 != nil && nodeNative5 != "" {
			objectDataLocalMap["RegionBlockCountryList"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].region_block_province_list", d.Get("content"))
		if nodeNative6 != nil && nodeNative6 != "" {
			objectDataLocalMap["RegionBlockProvinceList"] = nodeNative6
		}
		sourceLimit := make(map[string]interface{})
		nodeNative7, _ := jsonpath.Get("$[0].source_limit[0].pps", v)
		if nodeNative7 != nil && nodeNative7 != "" {
			sourceLimit["Pps"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].source_limit[0].bps", v)
		if nodeNative8 != nil && nodeNative8 != "" {
			sourceLimit["Bps"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].source_limit[0].syn_pps", v)
		if nodeNative9 != nil && nodeNative9 != "" {
			sourceLimit["SynPps"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].source_limit[0].syn_bps", v)
		if nodeNative10 != nil && nodeNative10 != "" {
			sourceLimit["SynBps"] = nodeNative10
		}

		objectDataLocalMap["SourceLimit"] = sourceLimit
		if v, ok := d.GetOk("content"); ok {
			localData, err := jsonpath.Get("$[0].source_block_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				if dataLoopTmp["type"].(int) > 0 {
					dataLoopMap["Type"] = dataLoopTmp["type"]
				}
				dataLoopMap["BlockExpireSeconds"] = dataLoopTmp["block_expire_seconds"]
				dataLoopMap["EverySeconds"] = dataLoopTmp["every_seconds"]
				dataLoopMap["ExceedLimitTimes"] = dataLoopTmp["exceed_limit_times"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["SourceBlockList"] = localMaps
		}
		nodeNative15, _ := jsonpath.Get("$[0].reflect_block_udp_port_list", d.Get("content"))
		if nodeNative15 != nil && nodeNative15 != "" {
			objectDataLocalMap["ReflectBlockUdpPortList"] = nodeNative15
		}
		if v, ok := d.GetOk("content"); ok {
			localData1, err := jsonpath.Get("$[0].port_rule_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Id"] = dataLoop1Tmp["port_rule_id"]
				dataLoop1Map["Protocol"] = dataLoop1Tmp["protocol"]
				dataLoop1Map["SrcPortStart"] = dataLoop1Tmp["src_port_start"]
				dataLoop1Map["SrcPortEnd"] = dataLoop1Tmp["src_port_end"]
				dataLoop1Map["DstPortStart"] = dataLoop1Tmp["dst_port_start"]
				dataLoop1Map["DstPortEnd"] = dataLoop1Tmp["dst_port_end"]
				dataLoop1Map["MatchAction"] = dataLoop1Tmp["match_action"]
				dataLoop1Map["SeqNo"] = dataLoop1Tmp["seq_no"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			objectDataLocalMap["PortRuleList"] = localMaps1
		}
		if v, ok := d.GetOk("content"); ok {
			localData2, err := jsonpath.Get("$[0].finger_print_rule_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps2 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["Id"] = dataLoop2Tmp["finger_print_rule_id"]
				dataLoop2Map["Protocol"] = dataLoop2Tmp["protocol"]
				dataLoop2Map["SrcPortStart"] = dataLoop2Tmp["src_port_start"]
				dataLoop2Map["SrcPortEnd"] = dataLoop2Tmp["src_port_end"]
				dataLoop2Map["DstPortStart"] = dataLoop2Tmp["dst_port_start"]
				dataLoop2Map["DstPortEnd"] = dataLoop2Tmp["dst_port_end"]
				dataLoop2Map["MinPktLen"] = dataLoop2Tmp["min_pkt_len"]
				dataLoop2Map["MaxPktLen"] = dataLoop2Tmp["max_pkt_len"]
				dataLoop2Map["Offset"] = dataLoop2Tmp["offset"]
				dataLoop2Map["PayloadBytes"] = dataLoop2Tmp["payload_bytes"]
				dataLoop2Map["MatchAction"] = dataLoop2Tmp["match_action"]
				dataLoop2Map["RateValue"] = dataLoop2Tmp["rate_value"]
				dataLoop2Map["SeqNo"] = dataLoop2Tmp["seq_no"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			objectDataLocalMap["FingerPrintRuleList"] = localMaps2
		}
		nodeNative37, _ := jsonpath.Get("$[0].enable_defense", v)
		if nodeNative37 != nil && nodeNative37 != "" {
			objectDataLocalMap["EnableL4Defense"] = nodeNative37
		}
		if v, ok := d.GetOk("content"); ok {
			localData3, err := jsonpath.Get("$[0].layer4_rule_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := dataLoop3.(map[string]interface{})
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["Name"] = dataLoop3Tmp["name"]
				dataLoop3Map["Priority"] = dataLoop3Tmp["priority"]
				dataLoop3Map["Method"] = dataLoop3Tmp["method"]
				dataLoop3Map["Match"] = dataLoop3Tmp["match"]
				dataLoop3Map["Action"] = dataLoop3Tmp["action"]
				dataLoop3Map["Limited"] = dataLoop3Tmp["limited"]
				localMaps4 := make([]interface{}, 0)
				localData4 := dataLoop3Tmp["condition_list"]
				for _, dataLoop4 := range localData4.([]interface{}) {
					dataLoop4Tmp := dataLoop4.(map[string]interface{})
					dataLoop4Map := make(map[string]interface{})
					dataLoop4Map["Arg"] = dataLoop4Tmp["arg"]
					dataLoop4Map["Position"] = dataLoop4Tmp["position"]
					dataLoop4Map["Depth"] = dataLoop4Tmp["depth"]
					localMaps4 = append(localMaps4, dataLoop4Map)
				}
				dataLoop3Map["ConditionList"] = localMaps4
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			objectDataLocalMap["L4RuleList"] = localMaps3
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Content"] = string(objectDataLocalMapJson)
	}

	if !d.IsNewResource() && d.HasChange("policy_name") {
		update = true
	}
	request["Name"] = d.Get("policy_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddosbgp", "2018-07-20", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
	}

	return resourceAliCloudDdosBgpPolicyRead(d, meta)
}

func resourceAliCloudDdosBgpPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
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
