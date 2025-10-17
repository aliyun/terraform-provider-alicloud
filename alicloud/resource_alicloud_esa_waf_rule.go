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
	"github.com/tidwall/sjson"
)

func resourceAliCloudEsaWafRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaWafRuleCreate,
		Read:   resourceAliCloudEsaWafRuleRead,
		Update: resourceAliCloudEsaWafRuleUpdate,
		Delete: resourceAliCloudEsaWafRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"response": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"code": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"bypass": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"skip": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"regular_rules": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeInt},
												},
												"custom_rules": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeInt},
												},
												"regular_types": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"tags": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"managed_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"managed_rulesets": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protection_level": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"managed_rules": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"action": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"attack_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"number_total": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"number_enabled": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"sigchl": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"app_sdk": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_sign": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"custom_sign_status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"feature_abnormal": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"rate_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"characteristics": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"criteria": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"criteria": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"match_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"criteria": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"match_type": {
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																				},
																			},
																		},
																		"logic": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"logic": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"logic": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"on_hit": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"managed_rules_blocked": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"distinct_managed_rules": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"response_status": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ratio": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"count": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"code": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"traffic": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"request": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"app_package": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package_signs": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sign": {
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
						"managed_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timer": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"periods": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"end": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"scopes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"zone": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weekly_periods": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"daily_periods": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"end": {
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
							},
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_level": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"phase": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ruleset_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"shared": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"response": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"code": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cross_site_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"match": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"criteria": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"criteria": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"criteria": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"match_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"logic": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"logic": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"logic": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"waf_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaWafRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("phase"); ok && InArray(fmt.Sprint(v), []string{"http_anti_scan", "http_bot"}) {
		action := "BatchCreateWafRules"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("site_id"); ok {
			request["SiteId"] = v
		}

		dataList := make(map[string]interface{})

		if v := d.Get("config"); !IsNil(v) {
			if v, ok := d.GetOk("config"); ok {
				localData, err := jsonpath.Get("$[0].managed_rulesets", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(localData) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["ProtectionLevel"] = dataLoopTmp["protection_level"]
					dataLoopMap["AttackType"] = dataLoopTmp["attack_type"]
					dataLoopMap["Action"] = dataLoopTmp["action"]
					localMaps1 := make([]interface{}, 0)
					localData1 := dataLoopTmp["managed_rules"]
					for _, dataLoop1 := range convertToInterfaceArray(localData1) {
						dataLoop1Tmp := dataLoop1.(map[string]interface{})
						dataLoop1Map := make(map[string]interface{})
						dataLoop1Map["Action"] = dataLoop1Tmp["action"]
						dataLoop1Map["Status"] = dataLoop1Tmp["status"]
						dataLoop1Map["Id"] = dataLoop1Tmp["id"]
						localMaps1 = append(localMaps1, dataLoop1Map)
					}
					dataLoopMap["ManagedRules"] = localMaps1
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList["ManagedRulesets"] = localMaps
			}

		}

		if v, ok := d.GetOk("config"); ok {
			sigchl1, _ := jsonpath.Get("$[0].sigchl", v)
			if sigchl1 != nil && sigchl1 != "" {
				dataList["Sigchl"] = convertToInterfaceArray(sigchl1)
			}
		}

		if v, ok := d.GetOk("config"); ok {
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && type1 != "" {
				dataList["Type"] = type1
			}
		}

		if v := d.Get("config"); !IsNil(v) {
			rateLimit := make(map[string]interface{})
			threshold := make(map[string]interface{})
			managedRulesBlocked1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].managed_rules_blocked", d.Get("config"))
			if managedRulesBlocked1 != nil && managedRulesBlocked1 != "" {
				threshold["ManagedRulesBlocked"] = managedRulesBlocked1
			}
			traffic1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].traffic", d.Get("config"))
			if traffic1 != nil && traffic1 != "" {
				threshold["Traffic"] = traffic1
			}
			responseStatus := make(map[string]interface{})
			ratio1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].ratio", d.Get("config"))
			if ratio1 != nil && ratio1 != "" {
				responseStatus["Ratio"] = ratio1
			}
			count1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].count", d.Get("config"))
			if count1 != nil && count1 != "" {
				responseStatus["Count"] = count1
			}
			code1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].code", d.Get("config"))
			if code1 != nil && code1 != "" {
				responseStatus["Code"] = code1
			}

			threshold["ResponseStatus"] = responseStatus
			request1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].request", d.Get("config"))
			if request1 != nil && request1 != "" {
				threshold["Request"] = request1
			}
			distinctManagedRules1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].distinct_managed_rules", d.Get("config"))
			if distinctManagedRules1 != nil && distinctManagedRules1 != "" {
				threshold["DistinctManagedRules"] = distinctManagedRules1
			}

			rateLimit["Threshold"] = threshold
			characteristics := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData2, err := jsonpath.Get("$[0].rate_limit[0].characteristics[0].criteria", v)
				if err != nil {
					localData2 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := make(map[string]interface{})
					if dataLoop2 != nil {
						dataLoop2Tmp = dataLoop2.(map[string]interface{})
					}
					dataLoop2Map := make(map[string]interface{})
					localMaps3 := make([]interface{}, 0)
					localData3 := dataLoop2Tmp["criteria"]
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := dataLoop3.(map[string]interface{})
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["Logic"] = dataLoop3Tmp["logic"]
						localMaps4 := make([]interface{}, 0)
						localData4 := dataLoop3Tmp["criteria"]
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := dataLoop4.(map[string]interface{})
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["MatchType"] = dataLoop4Tmp["match_type"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						dataLoop3Map["Criteria"] = localMaps4
						dataLoop3Map["MatchType"] = dataLoop3Tmp["match_type"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					dataLoop2Map["Criteria"] = localMaps3
					dataLoop2Map["Logic"] = dataLoop2Tmp["logic"]
					dataLoop2Map["MatchType"] = dataLoop2Tmp["match_type"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				characteristics["Criteria"] = localMaps2
			}

			logic5, _ := jsonpath.Get("$[0].rate_limit[0].characteristics[0].logic", d.Get("config"))
			if logic5 != nil && logic5 != "" {
				characteristics["Logic"] = logic5
			}

			rateLimit["Characteristics"] = characteristics
			interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
			if interval1 != nil && interval1 != "" {
				rateLimit["Interval"] = interval1
			}
			onHit1, _ := jsonpath.Get("$[0].rate_limit[0].on_hit", d.Get("config"))
			if onHit1 != nil && onHit1 != "" {
				rateLimit["OnHit"] = onHit1
			}
			ttl, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
			if ttl != nil && ttl != "" {
				rateLimit["TTL"] = ttl
			}

			dataList["RateLimit"] = rateLimit
		}

		if v := d.Get("config"); !IsNil(v) {
			timer := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData5, err := jsonpath.Get("$[0].timer[0].periods", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps5 := make([]interface{}, 0)
				for _, dataLoop5 := range convertToInterfaceArray(localData5) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["End"] = dataLoop5Tmp["end"]
					dataLoop5Map["Start"] = dataLoop5Tmp["start"]
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				timer["Periods"] = localMaps5
			}

			if v, ok := d.GetOk("config"); ok {
				localData6, err := jsonpath.Get("$[0].timer[0].weekly_periods", v)
				if err != nil {
					localData6 = make([]interface{}, 0)
				}
				localMaps6 := make([]interface{}, 0)
				for _, dataLoop6 := range convertToInterfaceArray(localData6) {
					dataLoop6Tmp := make(map[string]interface{})
					if dataLoop6 != nil {
						dataLoop6Tmp = dataLoop6.(map[string]interface{})
					}
					dataLoop6Map := make(map[string]interface{})
					localMaps7 := make([]interface{}, 0)
					localData7 := dataLoop6Tmp["daily_periods"]
					for _, dataLoop7 := range convertToInterfaceArray(localData7) {
						dataLoop7Tmp := dataLoop7.(map[string]interface{})
						dataLoop7Map := make(map[string]interface{})
						dataLoop7Map["End"] = dataLoop7Tmp["end"]
						dataLoop7Map["Start"] = dataLoop7Tmp["start"]
						localMaps7 = append(localMaps7, dataLoop7Map)
					}
					dataLoop6Map["DailyPeriods"] = localMaps7
					dataLoop6Map["Days"] = dataLoop6Tmp["days"]
					localMaps6 = append(localMaps6, dataLoop6Map)
				}
				timer["WeeklyPeriods"] = localMaps6
			}

			zone1, _ := jsonpath.Get("$[0].timer[0].zone", d.Get("config"))
			if zone1 != nil && zone1 != "" {
				timer["Zone"] = zone1
			}
			scopes1, _ := jsonpath.Get("$[0].timer[0].scopes", d.Get("config"))
			if scopes1 != nil && scopes1 != "" {
				timer["Scopes"] = scopes1
			}

			dataList["Timer"] = timer
		}

		if v, ok := d.GetOk("config"); ok {
			status3, _ := jsonpath.Get("$[0].status", v)
			if status3 != nil && status3 != "" {
				dataList["Status"] = status3
			}
		}

		if v, ok := d.GetOk("config"); ok {
			id3, _ := jsonpath.Get("$[0].id", v)
			if id3 != nil && id3 != "" {
				dataList["Id"] = id3
			}
		}

		if v, ok := d.GetOk("config"); ok {
			notes1, _ := jsonpath.Get("$[0].notes", v)
			if notes1 != nil && notes1 != "" {
				dataList["Notes"] = notes1
			}
		}

		if v, ok := d.GetOk("config"); ok {
			action5, _ := jsonpath.Get("$[0].action", v)
			if action5 != nil && action5 != "" {
				dataList["Action"] = action5
			}
		}

		if v := d.Get("config"); !IsNil(v) {
			appSdk := make(map[string]interface{})
			customSign := make(map[string]interface{})
			value1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].value", d.Get("config"))
			if value1 != nil && value1 != "" {
				customSign["Value"] = value1
			}
			key1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].key", d.Get("config"))
			if key1 != nil && key1 != "" {
				customSign["Key"] = key1
			}

			appSdk["CustomSign"] = customSign
			customSignStatus1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign_status", d.Get("config"))
			if customSignStatus1 != nil && customSignStatus1 != "" {
				appSdk["CustomSignStatus"] = customSignStatus1
			}
			featureAbnormal1, _ := jsonpath.Get("$[0].app_sdk[0].feature_abnormal", d.Get("config"))
			if featureAbnormal1 != nil && featureAbnormal1 != "" {
				appSdk["FeatureAbnormal"] = convertToInterfaceArray(featureAbnormal1)
			}

			dataList["AppSdk"] = appSdk
		}

		if v := d.Get("config"); !IsNil(v) {
			securityLevel := make(map[string]interface{})
			value3, _ := jsonpath.Get("$[0].security_level[0].value", d.Get("config"))
			if value3 != nil && value3 != "" {
				securityLevel["Value"] = value3
			}

			dataList["SecurityLevel"] = securityLevel
		}

		if v, ok := d.GetOk("config"); ok {
			value5, _ := jsonpath.Get("$[0].value", v)
			if value5 != nil && value5 != "" {
				dataList["Value"] = value5
			}
		}

		if v, ok := d.GetOk("config"); ok {
			expression1, _ := jsonpath.Get("$[0].expression", v)
			if expression1 != nil && expression1 != "" {
				dataList["Expression"] = expression1
			}
		}

		if v, ok := d.GetOk("config"); ok {
			managedGroupId1, _ := jsonpath.Get("$[0].managed_group_id", v)
			if managedGroupId1 != nil && managedGroupId1 != "" {
				dataList["ManagedGroupId"] = managedGroupId1
			}
		}

		if v, ok := d.GetOk("config"); ok {
			managedList1, _ := jsonpath.Get("$[0].managed_list", v)
			if managedList1 != nil && managedList1 != "" {
				dataList["ManagedList"] = managedList1
			}
		}

		if v := d.Get("config"); !IsNil(v) {
			appPackage := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData8, err := jsonpath.Get("$[0].app_package[0].package_signs", v)
				if err != nil {
					localData8 = make([]interface{}, 0)
				}
				localMaps8 := make([]interface{}, 0)
				for _, dataLoop8 := range convertToInterfaceArray(localData8) {
					dataLoop8Tmp := make(map[string]interface{})
					if dataLoop8 != nil {
						dataLoop8Tmp = dataLoop8.(map[string]interface{})
					}
					dataLoop8Map := make(map[string]interface{})
					dataLoop8Map["Sign"] = dataLoop8Tmp["sign"]
					dataLoop8Map["Name"] = dataLoop8Tmp["name"]
					localMaps8 = append(localMaps8, dataLoop8Map)
				}
				appPackage["PackageSigns"] = localMaps8
			}

			dataList["AppPackage"] = appPackage
		}

		if v, ok := d.GetOk("config"); ok {
			name3, _ := jsonpath.Get("$[0].name", v)
			if name3 != nil && name3 != "" {
				dataList["Name"] = name3
			}
		}

		if v := d.Get("config"); !IsNil(v) {
			actions := make(map[string]interface{})
			bypass := make(map[string]interface{})
			skip1, _ := jsonpath.Get("$[0].actions[0].bypass[0].skip", d.Get("config"))
			if skip1 != nil && skip1 != "" {
				bypass["Skip"] = skip1
			}
			customRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].custom_rules", d.Get("config"))
			if customRules1 != nil && customRules1 != "" {
				bypass["CustomRules"] = convertToInterfaceArray(customRules1)
			}
			regularRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_rules", d.Get("config"))
			if regularRules1 != nil && regularRules1 != "" {
				bypass["RegularRules"] = convertToInterfaceArray(regularRules1)
			}
			regularTypes1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_types", d.Get("config"))
			if regularTypes1 != nil && regularTypes1 != "" {
				bypass["RegularTypes"] = convertToInterfaceArray(regularTypes1)
			}
			tags1, _ := jsonpath.Get("$[0].actions[0].bypass[0].tags", d.Get("config"))
			if tags1 != nil && tags1 != "" {
				bypass["Tags"] = convertToInterfaceArray(tags1)
			}

			actions["Bypass"] = bypass
			response := make(map[string]interface{})
			code3, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("config"))
			if code3 != nil && code3 != "" {
				response["Code"] = code3
			}
			id5, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("config"))
			if id5 != nil && id5 != "" {
				response["Id"] = id5
			}

			actions["Response"] = response

			dataList["Actions"] = actions
		}

		ConfigsMap := make([]interface{}, 0)
		ConfigsMap = append(ConfigsMap, dataList)
		dataListJson, err := json.Marshal(ConfigsMap)
		if err != nil {
			return WrapError(err)
		}
		request["Configs"] = string(dataListJson)

		if v, ok := d.GetOkExists("site_version"); ok {
			request["SiteVersion"] = v
		}
		dataList1 := make(map[string]interface{})

		if v := d.Get("shared"); !IsNil(v) {
			name5, _ := jsonpath.Get("$[0].name", v)
			if name5 != nil && name5 != "" {
				dataList1["Name"] = name5
			}
			action7, _ := jsonpath.Get("$[0].action", v)
			if action7 != nil && action7 != "" {
				dataList1["Action"] = action7
			}
			crossSiteId1, _ := jsonpath.Get("$[0].cross_site_id", v)
			if crossSiteId1 != nil && crossSiteId1 != "" {
				dataList1["CrossSiteId"] = crossSiteId1
			}
			expression3, _ := jsonpath.Get("$[0].expression", v)
			if expression3 != nil && expression3 != "" {
				dataList1["Expression"] = expression3
			}
			target1, _ := jsonpath.Get("$[0].target", v)
			if target1 != nil && target1 != "" {
				dataList1["Target"] = target1
			}
			match := make(map[string]interface{})
			if v, ok := d.GetOk("shared"); ok {
				localData9, err := jsonpath.Get("$[0].match[0].criteria", v)
				if err != nil {
					localData9 = make([]interface{}, 0)
				}
				localMaps9 := make([]interface{}, 0)
				for _, dataLoop9 := range convertToInterfaceArray(localData9) {
					dataLoop9Tmp := make(map[string]interface{})
					if dataLoop9 != nil {
						dataLoop9Tmp = dataLoop9.(map[string]interface{})
					}
					dataLoop9Map := make(map[string]interface{})
					dataLoop9Map["Logic"] = dataLoop9Tmp["logic"]
					dataLoop9Map["MatchType"] = dataLoop9Tmp["match_type"]
					localMaps10 := make([]interface{}, 0)
					localData10 := dataLoop9Tmp["criteria"]
					for _, dataLoop10 := range convertToInterfaceArray(localData10) {
						dataLoop10Tmp := dataLoop10.(map[string]interface{})
						dataLoop10Map := make(map[string]interface{})
						localMaps11 := make([]interface{}, 0)
						localData11 := dataLoop10Tmp["criteria"]
						for _, dataLoop11 := range convertToInterfaceArray(localData11) {
							dataLoop11Tmp := dataLoop11.(map[string]interface{})
							dataLoop11Map := make(map[string]interface{})
							dataLoop11Map["MatchType"] = dataLoop11Tmp["match_type"]
							localMaps11 = append(localMaps11, dataLoop11Map)
						}
						dataLoop10Map["Criteria"] = localMaps11
						dataLoop10Map["Logic"] = dataLoop10Tmp["logic"]
						dataLoop10Map["MatchType"] = dataLoop10Tmp["match_type"]
						localMaps10 = append(localMaps10, dataLoop10Map)
					}
					dataLoop9Map["Criteria"] = localMaps10
					localMaps9 = append(localMaps9, dataLoop9Map)
				}
				match["Criteria"] = localMaps9
			}

			matchType13, _ := jsonpath.Get("$[0].match[0].match_type", d.Get("shared"))
			if matchType13 != nil && matchType13 != "" {
				match["MatchType"] = matchType13
			}
			logic11, _ := jsonpath.Get("$[0].match[0].logic", d.Get("shared"))
			if logic11 != nil && logic11 != "" {
				match["Logic"] = logic11
			}

			dataList1["Match"] = match
			mode1, _ := jsonpath.Get("$[0].mode", v)
			if mode1 != nil && mode1 != "" {
				dataList1["Mode"] = mode1
			}
			actions1 := make(map[string]interface{})
			response1 := make(map[string]interface{})
			code5, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("shared"))
			if code5 != nil && code5 != "" {
				response1["Code"] = code5
			}
			id7, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("shared"))
			if id7 != nil && id7 != "" {
				response1["Id"] = id7
			}

			actions1["Response"] = response1

			dataList1["Actions"] = actions1

			dataList1Json, err := json.Marshal(dataList1)
			if err != nil {
				return WrapError(err)
			}
			request["Shared"] = string(dataList1Json)
		}

		if v, ok := d.GetOkExists("ruleset_id"); ok {
			request["RulesetId"] = v
		}
		request["Phase"] = d.Get("phase")
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_waf_rule", action, AlibabaCloudSdkGoERROR)
		}

		IdsVar, _ := jsonpath.Get("$.Ids[0]", response)
		d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], IdsVar))

	}

	invalidCreate := false
	if v, ok := d.GetOk("phase"); ok {
		if InArray(fmt.Sprint(v), []string{"http_anti_scan", "http_bot"}) {
			invalidCreate = true
		}
	}
	if !invalidCreate {

		action := "CreateWafRule"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("site_id"); ok {
			request["SiteId"] = v
		}

		dataList := make(map[string]interface{})

		if v := d.Get("config"); !IsNil(v) {
			if v, ok := d.GetOk("config"); ok {
				localData, err := jsonpath.Get("$[0].managed_rulesets", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(localData) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["ProtectionLevel"] = dataLoopTmp["protection_level"]
					dataLoopMap["AttackType"] = dataLoopTmp["attack_type"]
					dataLoopMap["Action"] = dataLoopTmp["action"]
					localMaps1 := make([]interface{}, 0)
					localData1 := dataLoopTmp["managed_rules"]
					for _, dataLoop1 := range convertToInterfaceArray(localData1) {
						dataLoop1Tmp := dataLoop1.(map[string]interface{})
						dataLoop1Map := make(map[string]interface{})
						dataLoop1Map["Action"] = dataLoop1Tmp["action"]
						dataLoop1Map["Status"] = dataLoop1Tmp["status"]
						dataLoop1Map["Id"] = dataLoop1Tmp["id"]
						localMaps1 = append(localMaps1, dataLoop1Map)
					}
					dataLoopMap["ManagedRules"] = localMaps1
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList["ManagedRulesets"] = localMaps
			}

			sigchl1, _ := jsonpath.Get("$[0].sigchl", v)
			if sigchl1 != nil && sigchl1 != "" {
				dataList["Sigchl"] = convertToInterfaceArray(sigchl1)
			}
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && type1 != "" {
				dataList["Type"] = type1
			}
			rateLimit := make(map[string]interface{})
			threshold := make(map[string]interface{})
			managedRulesBlocked1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].managed_rules_blocked", d.Get("config"))
			if managedRulesBlocked1 != nil && managedRulesBlocked1 != "" {
				threshold["ManagedRulesBlocked"] = managedRulesBlocked1
			}
			traffic1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].traffic", d.Get("config"))
			if traffic1 != nil && traffic1 != "" {
				threshold["Traffic"] = traffic1
			}
			responseStatus := make(map[string]interface{})
			ratio1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].ratio", d.Get("config"))
			if ratio1 != nil && ratio1 != "" {
				responseStatus["Ratio"] = ratio1
			}
			count1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].count", d.Get("config"))
			if count1 != nil && count1 != "" {
				responseStatus["Count"] = count1
			}
			code1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].code", d.Get("config"))
			if code1 != nil && code1 != "" {
				responseStatus["Code"] = code1
			}

			threshold["ResponseStatus"] = responseStatus
			request1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].request", d.Get("config"))
			if request1 != nil && request1 != "" {
				threshold["Request"] = request1
			}
			distinctManagedRules1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].distinct_managed_rules", d.Get("config"))
			if distinctManagedRules1 != nil && distinctManagedRules1 != "" {
				threshold["DistinctManagedRules"] = distinctManagedRules1
			}

			rateLimit["Threshold"] = threshold
			characteristics := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData2, err := jsonpath.Get("$[0].rate_limit[0].characteristics[0].criteria", v)
				if err != nil {
					localData2 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := make(map[string]interface{})
					if dataLoop2 != nil {
						dataLoop2Tmp = dataLoop2.(map[string]interface{})
					}
					dataLoop2Map := make(map[string]interface{})
					localMaps3 := make([]interface{}, 0)
					localData3 := dataLoop2Tmp["criteria"]
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := dataLoop3.(map[string]interface{})
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["Logic"] = dataLoop3Tmp["logic"]
						localMaps4 := make([]interface{}, 0)
						localData4 := dataLoop3Tmp["criteria"]
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := dataLoop4.(map[string]interface{})
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["MatchType"] = dataLoop4Tmp["match_type"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						dataLoop3Map["Criteria"] = localMaps4
						dataLoop3Map["MatchType"] = dataLoop3Tmp["match_type"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					dataLoop2Map["Criteria"] = localMaps3
					dataLoop2Map["Logic"] = dataLoop2Tmp["logic"]
					dataLoop2Map["MatchType"] = dataLoop2Tmp["match_type"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				characteristics["Criteria"] = localMaps2
			}

			logic5, _ := jsonpath.Get("$[0].rate_limit[0].characteristics[0].logic", d.Get("config"))
			if logic5 != nil && logic5 != "" {
				characteristics["Logic"] = logic5
			}

			rateLimit["Characteristics"] = characteristics
			interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
			if interval1 != nil && interval1 != "" {
				rateLimit["Interval"] = interval1
			}
			onHit1, _ := jsonpath.Get("$[0].rate_limit[0].on_hit", d.Get("config"))
			if onHit1 != nil && onHit1 != "" {
				rateLimit["OnHit"] = onHit1
			}
			ttl, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
			if ttl != nil && ttl != "" {
				rateLimit["TTL"] = ttl
			}

			dataList["RateLimit"] = rateLimit
			timer := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData5, err := jsonpath.Get("$[0].timer[0].periods", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps5 := make([]interface{}, 0)
				for _, dataLoop5 := range convertToInterfaceArray(localData5) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["End"] = dataLoop5Tmp["end"]
					dataLoop5Map["Start"] = dataLoop5Tmp["start"]
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				timer["Periods"] = localMaps5
			}

			if v, ok := d.GetOk("config"); ok {
				localData6, err := jsonpath.Get("$[0].timer[0].weekly_periods", v)
				if err != nil {
					localData6 = make([]interface{}, 0)
				}
				localMaps6 := make([]interface{}, 0)
				for _, dataLoop6 := range convertToInterfaceArray(localData6) {
					dataLoop6Tmp := make(map[string]interface{})
					if dataLoop6 != nil {
						dataLoop6Tmp = dataLoop6.(map[string]interface{})
					}
					dataLoop6Map := make(map[string]interface{})
					localMaps7 := make([]interface{}, 0)
					localData7 := dataLoop6Tmp["daily_periods"]
					for _, dataLoop7 := range convertToInterfaceArray(localData7) {
						dataLoop7Tmp := dataLoop7.(map[string]interface{})
						dataLoop7Map := make(map[string]interface{})
						dataLoop7Map["End"] = dataLoop7Tmp["end"]
						dataLoop7Map["Start"] = dataLoop7Tmp["start"]
						localMaps7 = append(localMaps7, dataLoop7Map)
					}
					dataLoop6Map["DailyPeriods"] = localMaps7
					dataLoop6Map["Days"] = dataLoop6Tmp["days"]
					localMaps6 = append(localMaps6, dataLoop6Map)
				}
				timer["WeeklyPeriods"] = localMaps6
			}

			zone1, _ := jsonpath.Get("$[0].timer[0].zone", d.Get("config"))
			if zone1 != nil && zone1 != "" {
				timer["Zone"] = zone1
			}
			scopes1, _ := jsonpath.Get("$[0].timer[0].scopes", d.Get("config"))
			if scopes1 != nil && scopes1 != "" {
				timer["Scopes"] = scopes1
			}

			dataList["Timer"] = timer
			status3, _ := jsonpath.Get("$[0].status", v)
			if status3 != nil && status3 != "" {
				dataList["Status"] = status3
			}
			id3, _ := jsonpath.Get("$[0].id", v)
			if id3 != nil && id3 != "" {
				dataList["Id"] = id3
			}
			notes1, _ := jsonpath.Get("$[0].notes", v)
			if notes1 != nil && notes1 != "" {
				dataList["Notes"] = notes1
			}
			action5, _ := jsonpath.Get("$[0].action", v)
			if action5 != nil && action5 != "" {
				dataList["Action"] = action5
			}
			appSdk := make(map[string]interface{})
			customSign := make(map[string]interface{})
			value1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].value", d.Get("config"))
			if value1 != nil && value1 != "" {
				customSign["Value"] = value1
			}
			key1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].key", d.Get("config"))
			if key1 != nil && key1 != "" {
				customSign["Key"] = key1
			}

			appSdk["CustomSign"] = customSign
			customSignStatus1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign_status", d.Get("config"))
			if customSignStatus1 != nil && customSignStatus1 != "" {
				appSdk["CustomSignStatus"] = customSignStatus1
			}
			featureAbnormal1, _ := jsonpath.Get("$[0].app_sdk[0].feature_abnormal", d.Get("config"))
			if featureAbnormal1 != nil && featureAbnormal1 != "" {
				appSdk["FeatureAbnormal"] = convertToInterfaceArray(featureAbnormal1)
			}

			dataList["AppSdk"] = appSdk
			securityLevel := make(map[string]interface{})
			value3, _ := jsonpath.Get("$[0].security_level[0].value", d.Get("config"))
			if value3 != nil && value3 != "" {
				securityLevel["Value"] = value3
			}

			dataList["SecurityLevel"] = securityLevel
			value5, _ := jsonpath.Get("$[0].value", v)
			if value5 != nil && value5 != "" {
				dataList["Value"] = value5
			}
			expression1, _ := jsonpath.Get("$[0].expression", v)
			if expression1 != nil && expression1 != "" {
				dataList["Expression"] = expression1
			}
			managedGroupId1, _ := jsonpath.Get("$[0].managed_group_id", v)
			if managedGroupId1 != nil && managedGroupId1 != "" {
				dataList["ManagedGroupId"] = managedGroupId1
			}
			managedList1, _ := jsonpath.Get("$[0].managed_list", v)
			if managedList1 != nil && managedList1 != "" {
				dataList["ManagedList"] = managedList1
			}
			appPackage := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData8, err := jsonpath.Get("$[0].app_package[0].package_signs", v)
				if err != nil {
					localData8 = make([]interface{}, 0)
				}
				localMaps8 := make([]interface{}, 0)
				for _, dataLoop8 := range convertToInterfaceArray(localData8) {
					dataLoop8Tmp := make(map[string]interface{})
					if dataLoop8 != nil {
						dataLoop8Tmp = dataLoop8.(map[string]interface{})
					}
					dataLoop8Map := make(map[string]interface{})
					dataLoop8Map["Sign"] = dataLoop8Tmp["sign"]
					dataLoop8Map["Name"] = dataLoop8Tmp["name"]
					localMaps8 = append(localMaps8, dataLoop8Map)
				}
				appPackage["PackageSigns"] = localMaps8
			}

			dataList["AppPackage"] = appPackage
			name3, _ := jsonpath.Get("$[0].name", v)
			if name3 != nil && name3 != "" {
				dataList["Name"] = name3
			}
			actions := make(map[string]interface{})
			bypass := make(map[string]interface{})
			skip1, _ := jsonpath.Get("$[0].actions[0].bypass[0].skip", d.Get("config"))
			if skip1 != nil && skip1 != "" {
				bypass["Skip"] = skip1
			}
			customRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].custom_rules", d.Get("config"))
			if customRules1 != nil && customRules1 != "" {
				bypass["CustomRules"] = convertToInterfaceArray(customRules1)
			}
			regularRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_rules", d.Get("config"))
			if regularRules1 != nil && regularRules1 != "" {
				bypass["RegularRules"] = convertToInterfaceArray(regularRules1)
			}
			regularTypes1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_types", d.Get("config"))
			if regularTypes1 != nil && regularTypes1 != "" {
				bypass["RegularTypes"] = convertToInterfaceArray(regularTypes1)
			}
			tags1, _ := jsonpath.Get("$[0].actions[0].bypass[0].tags", d.Get("config"))
			if tags1 != nil && tags1 != "" {
				bypass["Tags"] = convertToInterfaceArray(tags1)
			}

			actions["Bypass"] = bypass
			response := make(map[string]interface{})
			code3, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("config"))
			if code3 != nil && code3 != "" {
				response["Code"] = code3
			}
			id5, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("config"))
			if id5 != nil && id5 != "" {
				response["Id"] = id5
			}

			actions["Response"] = response

			dataList["Actions"] = actions

			dataListJson, err := json.Marshal(dataList)
			if err != nil {
				return WrapError(err)
			}
			request["Config"] = string(dataListJson)
		}

		if v, ok := d.GetOkExists("site_version"); ok {
			request["SiteVersion"] = v
		}
		if v, ok := d.GetOkExists("ruleset_id"); ok {
			request["RulesetId"] = v
		}
		request["Phase"] = d.Get("phase")
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_waf_rule", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	}

	return resourceAliCloudEsaWafRuleUpdate(d, meta)
}

func resourceAliCloudEsaWafRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaWafRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_waf_rule DescribeEsaWafRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("phase", objectRaw["Phase"])
	d.Set("ruleset_id", objectRaw["RulesetId"])
	d.Set("waf_rule_id", objectRaw["Id"])

	configMaps := make([]map[string]interface{}, 0)
	configMap := make(map[string]interface{})
	configRaw := make(map[string]interface{})
	if objectRaw["Config"] != nil {
		configRaw = objectRaw["Config"].(map[string]interface{})
	}
	if len(configRaw) > 0 {
		configMap["action"] = configRaw["Action"]
		configMap["expression"] = configRaw["Expression"]
		configMap["id"] = configRaw["Id"]
		configMap["managed_group_id"] = configRaw["ManagedGroupId"]
		configMap["managed_list"] = configRaw["ManagedList"]
		configMap["name"] = configRaw["Name"]
		configMap["notes"] = configRaw["Notes"]
		configMap["status"] = configRaw["Status"]
		configMap["type"] = configRaw["Type"]
		configMap["value"] = configRaw["Value"]

		actionsMaps := make([]map[string]interface{}, 0)
		actionsMap := make(map[string]interface{})
		actionsRaw := make(map[string]interface{})
		if configRaw["Actions"] != nil {
			actionsRaw = configRaw["Actions"].(map[string]interface{})
		}
		if len(actionsRaw) > 0 {

			bypassMaps := make([]map[string]interface{}, 0)
			bypassMap := make(map[string]interface{})
			bypassRaw := make(map[string]interface{})
			if actionsRaw["Bypass"] != nil {
				bypassRaw = actionsRaw["Bypass"].(map[string]interface{})
			}
			if len(bypassRaw) > 0 {
				bypassMap["skip"] = bypassRaw["Skip"]

				customRulesRaw := make([]interface{}, 0)
				if bypassRaw["CustomRules"] != nil {
					customRulesRaw = convertToInterfaceArray(bypassRaw["CustomRules"])
				}

				bypassMap["custom_rules"] = customRulesRaw
				regularRulesRaw := make([]interface{}, 0)
				if bypassRaw["RegularRules"] != nil {
					regularRulesRaw = convertToInterfaceArray(bypassRaw["RegularRules"])
				}

				bypassMap["regular_rules"] = regularRulesRaw
				regularTypesRaw := make([]interface{}, 0)
				if bypassRaw["RegularTypes"] != nil {
					regularTypesRaw = convertToInterfaceArray(bypassRaw["RegularTypes"])
				}

				bypassMap["regular_types"] = regularTypesRaw
				tagsRaw := make([]interface{}, 0)
				if bypassRaw["Tags"] != nil {
					tagsRaw = convertToInterfaceArray(bypassRaw["Tags"])
				}

				bypassMap["tags"] = tagsRaw
				bypassMaps = append(bypassMaps, bypassMap)
			}
			actionsMap["bypass"] = bypassMaps
			responseMaps := make([]map[string]interface{}, 0)
			responseMap := make(map[string]interface{})
			responseRaw := make(map[string]interface{})
			if actionsRaw["Response"] != nil {
				responseRaw = actionsRaw["Response"].(map[string]interface{})
			}
			if len(responseRaw) > 0 {
				responseMap["code"] = responseRaw["Code"]
				responseMap["id"] = responseRaw["Id"]

				responseMaps = append(responseMaps, responseMap)
			}
			actionsMap["response"] = responseMaps
			actionsMaps = append(actionsMaps, actionsMap)
		}
		configMap["actions"] = actionsMaps
		appPackageMaps := make([]map[string]interface{}, 0)
		appPackageMap := make(map[string]interface{})
		packageSignsRaw, _ := jsonpath.Get("$.Config.AppPackage.PackageSigns", objectRaw)

		packageSignsMaps := make([]map[string]interface{}, 0)
		if packageSignsRaw != nil {
			for _, packageSignsChildRaw := range convertToInterfaceArray(packageSignsRaw) {
				packageSignsMap := make(map[string]interface{})
				packageSignsChildRaw := packageSignsChildRaw.(map[string]interface{})
				packageSignsMap["name"] = packageSignsChildRaw["Name"]
				packageSignsMap["sign"] = packageSignsChildRaw["Sign"]

				packageSignsMaps = append(packageSignsMaps, packageSignsMap)
			}
		}
		appPackageMap["package_signs"] = packageSignsMaps
		appPackageMaps = append(appPackageMaps, appPackageMap)
		configMap["app_package"] = appPackageMaps
		appSdkMaps := make([]map[string]interface{}, 0)
		appSdkMap := make(map[string]interface{})
		appSdkRaw := make(map[string]interface{})
		if configRaw["AppSdk"] != nil {
			appSdkRaw = configRaw["AppSdk"].(map[string]interface{})
		}
		if len(appSdkRaw) > 0 {
			appSdkMap["custom_sign_status"] = appSdkRaw["CustomSignStatus"]

			customSignMaps := make([]map[string]interface{}, 0)
			customSignMap := make(map[string]interface{})
			customSignRaw := make(map[string]interface{})
			if appSdkRaw["CustomSign"] != nil {
				customSignRaw = appSdkRaw["CustomSign"].(map[string]interface{})
			}
			if len(customSignRaw) > 0 {
				customSignMap["key"] = customSignRaw["Key"]
				customSignMap["value"] = customSignRaw["Value"]

				customSignMaps = append(customSignMaps, customSignMap)
			}
			appSdkMap["custom_sign"] = customSignMaps
			featureAbnormalRaw := make([]interface{}, 0)
			if appSdkRaw["FeatureAbnormal"] != nil {
				featureAbnormalRaw = convertToInterfaceArray(appSdkRaw["FeatureAbnormal"])
			}

			appSdkMap["feature_abnormal"] = featureAbnormalRaw
			appSdkMaps = append(appSdkMaps, appSdkMap)
		}
		configMap["app_sdk"] = appSdkMaps
		managedRulesetsRaw := configRaw["ManagedRulesets"]
		managedRulesetsMaps := make([]map[string]interface{}, 0)
		if managedRulesetsRaw != nil {
			for _, managedRulesetsChildRaw := range convertToInterfaceArray(managedRulesetsRaw) {
				managedRulesetsMap := make(map[string]interface{})
				managedRulesetsChildRaw := managedRulesetsChildRaw.(map[string]interface{})
				managedRulesetsMap["action"] = managedRulesetsChildRaw["Action"]
				managedRulesetsMap["attack_type"] = managedRulesetsChildRaw["AttackType"]
				managedRulesetsMap["number_enabled"] = managedRulesetsChildRaw["NumberEnabled"]
				managedRulesetsMap["number_total"] = managedRulesetsChildRaw["NumberTotal"]
				managedRulesetsMap["protection_level"] = managedRulesetsChildRaw["ProtectionLevel"]

				managedRulesRaw := managedRulesetsChildRaw["ManagedRules"]
				managedRulesMaps := make([]map[string]interface{}, 0)
				if managedRulesRaw != nil {
					for _, managedRulesChildRaw := range convertToInterfaceArray(managedRulesRaw) {
						managedRulesMap := make(map[string]interface{})
						managedRulesChildRaw := managedRulesChildRaw.(map[string]interface{})
						managedRulesMap["action"] = managedRulesChildRaw["Action"]
						managedRulesMap["id"] = managedRulesChildRaw["Id"]
						managedRulesMap["status"] = managedRulesChildRaw["Status"]

						managedRulesMaps = append(managedRulesMaps, managedRulesMap)
					}
				}
				managedRulesetsMap["managed_rules"] = managedRulesMaps
				managedRulesetsMaps = append(managedRulesetsMaps, managedRulesetsMap)
			}
		}
		configMap["managed_rulesets"] = managedRulesetsMaps
		rateLimitMaps := make([]map[string]interface{}, 0)
		rateLimitMap := make(map[string]interface{})
		rateLimitRaw := make(map[string]interface{})
		if configRaw["RateLimit"] != nil {
			rateLimitRaw = configRaw["RateLimit"].(map[string]interface{})
		}
		if len(rateLimitRaw) > 0 {
			rateLimitMap["interval"] = rateLimitRaw["Interval"]
			rateLimitMap["on_hit"] = rateLimitRaw["OnHit"]
			rateLimitMap["ttl"] = rateLimitRaw["TTL"]

			characteristicsMaps := make([]map[string]interface{}, 0)
			characteristicsMap := make(map[string]interface{})
			characteristicsRaw := make(map[string]interface{})
			if rateLimitRaw["Characteristics"] != nil {
				characteristicsRaw = rateLimitRaw["Characteristics"].(map[string]interface{})
			}
			if len(characteristicsRaw) > 0 {
				characteristicsMap["logic"] = characteristicsRaw["Logic"]

				criteriaRaw := characteristicsRaw["Criteria"]
				criteriaMaps := make([]map[string]interface{}, 0)
				if criteriaRaw != nil {
					for _, criteriaChildRaw := range convertToInterfaceArray(criteriaRaw) {
						criteriaMap := make(map[string]interface{})
						criteriaChildRaw := criteriaChildRaw.(map[string]interface{})
						criteriaMap["logic"] = criteriaChildRaw["Logic"]
						criteriaMap["match_type"] = criteriaChildRaw["MatchType"]

						criteriaRawLevel2 := criteriaChildRaw["Criteria"]
						criteriaMapsLevel2 := make([]map[string]interface{}, 0)
						if criteriaRawLevel2 != nil {
							for _, criteriaChildRawLevel2 := range convertToInterfaceArray(criteriaRawLevel2) {
								criteriaMapLevel2 := make(map[string]interface{})
								criteriaChildRawLevel2 := criteriaChildRawLevel2.(map[string]interface{})
								criteriaMapLevel2["logic"] = criteriaChildRawLevel2["Logic"]
								criteriaMapLevel2["match_type"] = criteriaChildRawLevel2["MatchType"]

								criteriaRawLevel3 := criteriaChildRawLevel2["Criteria"]
								criteriaMapsLevel3 := make([]map[string]interface{}, 0)
								if criteriaRawLevel3 != nil {
									for _, criteriaChildRawLevel3 := range convertToInterfaceArray(criteriaRawLevel3) {
										criteriaMapLevel3 := make(map[string]interface{})
										criteriaChildRawLevel3 := criteriaChildRawLevel3.(map[string]interface{})
										criteriaMapLevel3["match_type"] = criteriaChildRawLevel3["MatchType"]

										criteriaMapsLevel3 = append(criteriaMapsLevel3, criteriaMapLevel3)
									}
								}
								criteriaMapLevel2["criteria"] = criteriaMapsLevel3
								criteriaMapsLevel2 = append(criteriaMapsLevel2, criteriaMapLevel2)
							}
						}
						criteriaMap["criteria"] = criteriaMapsLevel2
						criteriaMaps = append(criteriaMaps, criteriaMap)
					}
				}
				characteristicsMap["criteria"] = criteriaMaps
				characteristicsMaps = append(characteristicsMaps, characteristicsMap)
			}
			rateLimitMap["characteristics"] = characteristicsMaps
			thresholdMaps := make([]map[string]interface{}, 0)
			thresholdMap := make(map[string]interface{})
			thresholdRaw := make(map[string]interface{})
			if rateLimitRaw["Threshold"] != nil {
				thresholdRaw = rateLimitRaw["Threshold"].(map[string]interface{})
			}
			if len(thresholdRaw) > 0 {
				thresholdMap["distinct_managed_rules"] = thresholdRaw["DistinctManagedRules"]
				thresholdMap["managed_rules_blocked"] = thresholdRaw["ManagedRulesBlocked"]
				thresholdMap["request"] = thresholdRaw["Request"]
				thresholdMap["traffic"] = thresholdRaw["Traffic"]

				responseStatusMaps := make([]map[string]interface{}, 0)
				responseStatusMap := make(map[string]interface{})
				responseStatusRaw := make(map[string]interface{})
				if thresholdRaw["ResponseStatus"] != nil {
					responseStatusRaw = thresholdRaw["ResponseStatus"].(map[string]interface{})
				}
				if len(responseStatusRaw) > 0 {
					responseStatusMap["code"] = responseStatusRaw["Code"]
					responseStatusMap["count"] = responseStatusRaw["Count"]
					responseStatusMap["ratio"] = responseStatusRaw["Ratio"]

					responseStatusMaps = append(responseStatusMaps, responseStatusMap)
				}
				thresholdMap["response_status"] = responseStatusMaps
				thresholdMaps = append(thresholdMaps, thresholdMap)
			}
			rateLimitMap["threshold"] = thresholdMaps
			rateLimitMaps = append(rateLimitMaps, rateLimitMap)
		}
		configMap["rate_limit"] = rateLimitMaps
		securityLevelMaps := make([]map[string]interface{}, 0)
		securityLevelMap := make(map[string]interface{})
		securityLevelRaw := make(map[string]interface{})
		if configRaw["SecurityLevel"] != nil {
			securityLevelRaw = configRaw["SecurityLevel"].(map[string]interface{})
		}
		if len(securityLevelRaw) > 0 {
			securityLevelMap["value"] = securityLevelRaw["Value"]

			securityLevelMaps = append(securityLevelMaps, securityLevelMap)
		}
		configMap["security_level"] = securityLevelMaps
		sigchlRaw := make([]interface{}, 0)
		if configRaw["Sigchl"] != nil {
			sigchlRaw = convertToInterfaceArray(configRaw["Sigchl"])
		}

		configMap["sigchl"] = sigchlRaw
		timerMaps := make([]map[string]interface{}, 0)
		timerMap := make(map[string]interface{})
		timerRaw := make(map[string]interface{})
		if configRaw["Timer"] != nil {
			timerRaw = configRaw["Timer"].(map[string]interface{})
		}
		if len(timerRaw) > 0 {
			timerMap["scopes"] = timerRaw["Scopes"]
			timerMap["zone"] = timerRaw["Zone"]

			periodsRaw := timerRaw["Periods"]
			periodsMaps := make([]map[string]interface{}, 0)
			if periodsRaw != nil {
				for _, periodsChildRaw := range convertToInterfaceArray(periodsRaw) {
					periodsMap := make(map[string]interface{})
					periodsChildRaw := periodsChildRaw.(map[string]interface{})
					periodsMap["end"] = periodsChildRaw["End"]
					periodsMap["start"] = periodsChildRaw["Start"]

					periodsMaps = append(periodsMaps, periodsMap)
				}
			}
			timerMap["periods"] = periodsMaps
			weeklyPeriodsRaw := timerRaw["WeeklyPeriods"]
			weeklyPeriodsMaps := make([]map[string]interface{}, 0)
			if weeklyPeriodsRaw != nil {
				for _, weeklyPeriodsChildRaw := range convertToInterfaceArray(weeklyPeriodsRaw) {
					weeklyPeriodsMap := make(map[string]interface{})
					weeklyPeriodsChildRaw := weeklyPeriodsChildRaw.(map[string]interface{})
					weeklyPeriodsMap["days"] = weeklyPeriodsChildRaw["Days"]

					dailyPeriodsRaw := weeklyPeriodsChildRaw["DailyPeriods"]
					dailyPeriodsMaps := make([]map[string]interface{}, 0)
					if dailyPeriodsRaw != nil {
						for _, dailyPeriodsChildRaw := range convertToInterfaceArray(dailyPeriodsRaw) {
							dailyPeriodsMap := make(map[string]interface{})
							dailyPeriodsChildRaw := dailyPeriodsChildRaw.(map[string]interface{})
							dailyPeriodsMap["end"] = dailyPeriodsChildRaw["End"]
							dailyPeriodsMap["start"] = dailyPeriodsChildRaw["Start"]

							dailyPeriodsMaps = append(dailyPeriodsMaps, dailyPeriodsMap)
						}
					}
					weeklyPeriodsMap["daily_periods"] = dailyPeriodsMaps
					weeklyPeriodsMaps = append(weeklyPeriodsMaps, weeklyPeriodsMap)
				}
			}
			timerMap["weekly_periods"] = weeklyPeriodsMaps
			timerMaps = append(timerMaps, timerMap)
		}
		configMap["timer"] = timerMaps
		configMaps = append(configMaps, configMap)
	}
	if err := d.Set("config", configMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaWafRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	esaServiceV2 := EsaServiceV2{client}
	objectRaw, _ := esaServiceV2.DescribeEsaWafRule(d.Id())

	var err error
	enableBatchUpdateWafRules := false
	checkValue00 := objectRaw["Phase"]
	if InArray(fmt.Sprint(checkValue00), []string{"http_anti_scan", "http_bot"}) {
		enableBatchUpdateWafRules = true
	}
	parts := strings.Split(d.Id(), ":")
	action := "BatchUpdateWafRules"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]

	dataList := make(map[string]interface{})

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			if v, ok := d.GetOk("config"); ok {
				localData, err := jsonpath.Get("$[0].managed_rulesets", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(localData) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["ProtectionLevel"] = dataLoopTmp["protection_level"]
					dataLoopMap["AttackType"] = dataLoopTmp["attack_type"]
					dataLoopMap["Action"] = dataLoopTmp["action"]
					localMaps1 := make([]interface{}, 0)
					localData1 := dataLoopTmp["managed_rules"]
					for _, dataLoop1 := range convertToInterfaceArray(localData1) {
						dataLoop1Tmp := dataLoop1.(map[string]interface{})
						dataLoop1Map := make(map[string]interface{})
						dataLoop1Map["Action"] = dataLoop1Tmp["action"]
						dataLoop1Map["Status"] = dataLoop1Tmp["status"]
						dataLoop1Map["Id"] = dataLoop1Tmp["id"]
						localMaps1 = append(localMaps1, dataLoop1Map)
					}
					dataLoopMap["ManagedRules"] = localMaps1
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList["ManagedRulesets"] = localMaps
			}

		}
	}

	if d.HasChange("config") {
		update = true
		sigchl1, _ := jsonpath.Get("$[0].sigchl", d.Get("config"))
		if sigchl1 != nil && (d.HasChange("config.0.sigchl") || sigchl1 != "") {
			dataList["Sigchl"] = convertToInterfaceArray(sigchl1)
		}
	}

	if d.HasChange("config") {
		update = true
		type1, _ := jsonpath.Get("$[0].type", d.Get("config"))
		if type1 != nil && (d.HasChange("config.0.type") || type1 != "") {
			dataList["Type"] = type1
		}
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			rateLimit := make(map[string]interface{})
			threshold := make(map[string]interface{})
			managedRulesBlocked1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].managed_rules_blocked", d.Get("config"))
			if managedRulesBlocked1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.managed_rules_blocked") || managedRulesBlocked1 != "") {
				threshold["ManagedRulesBlocked"] = managedRulesBlocked1
			}
			traffic1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].traffic", d.Get("config"))
			if traffic1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.traffic") || traffic1 != "") {
				threshold["Traffic"] = traffic1
			}
			responseStatus := make(map[string]interface{})
			ratio1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].ratio", d.Get("config"))
			if ratio1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.ratio") || ratio1 != "") {
				responseStatus["Ratio"] = ratio1
			}
			count1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].count", d.Get("config"))
			if count1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.count") || count1 != "") {
				responseStatus["Count"] = count1
			}
			code1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].code", d.Get("config"))
			if code1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.code") || code1 != "") {
				responseStatus["Code"] = code1
			}

			threshold["ResponseStatus"] = responseStatus
			request1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].request", d.Get("config"))
			if request1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.request") || request1 != "") {
				threshold["Request"] = request1
			}
			distinctManagedRules1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].distinct_managed_rules", d.Get("config"))
			if distinctManagedRules1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.distinct_managed_rules") || distinctManagedRules1 != "") {
				threshold["DistinctManagedRules"] = distinctManagedRules1
			}

			rateLimit["Threshold"] = threshold
			characteristics := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData2, err := jsonpath.Get("$[0].rate_limit[0].characteristics[0].criteria", v)
				if err != nil {
					localData2 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := make(map[string]interface{})
					if dataLoop2 != nil {
						dataLoop2Tmp = dataLoop2.(map[string]interface{})
					}
					dataLoop2Map := make(map[string]interface{})
					localMaps3 := make([]interface{}, 0)
					localData3 := dataLoop2Tmp["criteria"]
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := dataLoop3.(map[string]interface{})
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["Logic"] = dataLoop3Tmp["logic"]
						localMaps4 := make([]interface{}, 0)
						localData4 := dataLoop3Tmp["criteria"]
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := dataLoop4.(map[string]interface{})
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["MatchType"] = dataLoop4Tmp["match_type"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						dataLoop3Map["Criteria"] = localMaps4
						dataLoop3Map["MatchType"] = dataLoop3Tmp["match_type"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					dataLoop2Map["Criteria"] = localMaps3
					dataLoop2Map["Logic"] = dataLoop2Tmp["logic"]
					dataLoop2Map["MatchType"] = dataLoop2Tmp["match_type"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				characteristics["Criteria"] = localMaps2
			}

			logic5, _ := jsonpath.Get("$[0].rate_limit[0].characteristics[0].logic", d.Get("config"))
			if logic5 != nil && (d.HasChange("config.0.rate_limit.0.characteristics.0.logic") || logic5 != "") {
				characteristics["Logic"] = logic5
			}

			rateLimit["Characteristics"] = characteristics
			interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
			if interval1 != nil && (d.HasChange("config.0.rate_limit.0.interval") || interval1 != "") {
				rateLimit["Interval"] = interval1
			}
			onHit1, _ := jsonpath.Get("$[0].rate_limit[0].on_hit", d.Get("config"))
			if onHit1 != nil && (d.HasChange("config.0.rate_limit.0.on_hit") || onHit1 != "") {
				rateLimit["OnHit"] = onHit1
			}
			ttl, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
			if ttl != nil && (d.HasChange("config.0.rate_limit.0.ttl") || ttl != "") {
				rateLimit["TTL"] = ttl
			}

			dataList["RateLimit"] = rateLimit
		}
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			timer := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData5, err := jsonpath.Get("$[0].timer[0].periods", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps5 := make([]interface{}, 0)
				for _, dataLoop5 := range convertToInterfaceArray(localData5) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["End"] = dataLoop5Tmp["end"]
					dataLoop5Map["Start"] = dataLoop5Tmp["start"]
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				timer["Periods"] = localMaps5
			}

			if v, ok := d.GetOk("config"); ok {
				localData6, err := jsonpath.Get("$[0].timer[0].weekly_periods", v)
				if err != nil {
					localData6 = make([]interface{}, 0)
				}
				localMaps6 := make([]interface{}, 0)
				for _, dataLoop6 := range convertToInterfaceArray(localData6) {
					dataLoop6Tmp := make(map[string]interface{})
					if dataLoop6 != nil {
						dataLoop6Tmp = dataLoop6.(map[string]interface{})
					}
					dataLoop6Map := make(map[string]interface{})
					localMaps7 := make([]interface{}, 0)
					localData7 := dataLoop6Tmp["daily_periods"]
					for _, dataLoop7 := range convertToInterfaceArray(localData7) {
						dataLoop7Tmp := dataLoop7.(map[string]interface{})
						dataLoop7Map := make(map[string]interface{})
						dataLoop7Map["End"] = dataLoop7Tmp["end"]
						dataLoop7Map["Start"] = dataLoop7Tmp["start"]
						localMaps7 = append(localMaps7, dataLoop7Map)
					}
					dataLoop6Map["DailyPeriods"] = localMaps7
					dataLoop6Map["Days"] = dataLoop6Tmp["days"]
					localMaps6 = append(localMaps6, dataLoop6Map)
				}
				timer["WeeklyPeriods"] = localMaps6
			}

			zone1, _ := jsonpath.Get("$[0].timer[0].zone", d.Get("config"))
			if zone1 != nil && (d.HasChange("config.0.timer.0.zone") || zone1 != "") {
				timer["Zone"] = zone1
			}
			scopes1, _ := jsonpath.Get("$[0].timer[0].scopes", d.Get("config"))
			if scopes1 != nil && (d.HasChange("config.0.timer.0.scopes") || scopes1 != "") {
				timer["Scopes"] = scopes1
			}

			dataList["Timer"] = timer
		}
	}

	if d.HasChange("config") {
		update = true
		status3, _ := jsonpath.Get("$[0].status", d.Get("config"))
		if status3 != nil && (d.HasChange("config.0.status") || status3 != "") {
			dataList["Status"] = status3
		}
	}

	if d.HasChange("config") {
		update = true
		notes1, _ := jsonpath.Get("$[0].notes", d.Get("config"))
		if notes1 != nil && (d.HasChange("config.0.notes") || notes1 != "") {
			dataList["Notes"] = notes1
		}
	}

	if d.HasChange("config") {
		update = true
		action5, _ := jsonpath.Get("$[0].action", d.Get("config"))
		if action5 != nil && (d.HasChange("config.0.action") || action5 != "") {
			dataList["Action"] = action5
		}
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			appSdk := make(map[string]interface{})
			customSign := make(map[string]interface{})
			value1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].value", d.Get("config"))
			if value1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign.0.value") || value1 != "") {
				customSign["Value"] = value1
			}
			key1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].key", d.Get("config"))
			if key1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign.0.key") || key1 != "") {
				customSign["Key"] = key1
			}

			appSdk["CustomSign"] = customSign
			customSignStatus1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign_status", d.Get("config"))
			if customSignStatus1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign_status") || customSignStatus1 != "") {
				appSdk["CustomSignStatus"] = customSignStatus1
			}
			featureAbnormal1, _ := jsonpath.Get("$[0].app_sdk[0].feature_abnormal", d.Get("config"))
			if featureAbnormal1 != nil && (d.HasChange("config.0.app_sdk.0.feature_abnormal") || featureAbnormal1 != "") {
				appSdk["FeatureAbnormal"] = convertToInterfaceArray(featureAbnormal1)
			}

			dataList["AppSdk"] = appSdk
		}
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			securityLevel := make(map[string]interface{})
			value3, _ := jsonpath.Get("$[0].security_level[0].value", d.Get("config"))
			if value3 != nil && (d.HasChange("config.0.security_level.0.value") || value3 != "") {
				securityLevel["Value"] = value3
			}

			dataList["SecurityLevel"] = securityLevel
		}
	}

	if d.HasChange("config") {
		update = true
		value5, _ := jsonpath.Get("$[0].value", d.Get("config"))
		if value5 != nil && (d.HasChange("config.0.value") || value5 != "") {
			dataList["Value"] = value5
		}
	}

	if d.HasChange("config") {
		update = true
		expression1, _ := jsonpath.Get("$[0].expression", d.Get("config"))
		if expression1 != nil && (d.HasChange("config.0.expression") || expression1 != "") {
			dataList["Expression"] = expression1
		}
	}

	if d.HasChange("config") {
		update = true
		managedGroupId1, _ := jsonpath.Get("$[0].managed_group_id", d.Get("config"))
		if managedGroupId1 != nil && (d.HasChange("config.0.managed_group_id") || managedGroupId1 != "") {
			dataList["ManagedGroupId"] = managedGroupId1
		}
	}

	if d.HasChange("config") {
		update = true
		managedList1, _ := jsonpath.Get("$[0].managed_list", d.Get("config"))
		if managedList1 != nil && (d.HasChange("config.0.managed_list") || managedList1 != "") {
			dataList["ManagedList"] = managedList1
		}
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			appPackage := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData8, err := jsonpath.Get("$[0].app_package[0].package_signs", v)
				if err != nil {
					localData8 = make([]interface{}, 0)
				}
				localMaps8 := make([]interface{}, 0)
				for _, dataLoop8 := range convertToInterfaceArray(localData8) {
					dataLoop8Tmp := make(map[string]interface{})
					if dataLoop8 != nil {
						dataLoop8Tmp = dataLoop8.(map[string]interface{})
					}
					dataLoop8Map := make(map[string]interface{})
					dataLoop8Map["Sign"] = dataLoop8Tmp["sign"]
					dataLoop8Map["Name"] = dataLoop8Tmp["name"]
					localMaps8 = append(localMaps8, dataLoop8Map)
				}
				appPackage["PackageSigns"] = localMaps8
			}

			dataList["AppPackage"] = appPackage
		}
	}

	if d.HasChange("config") {
		update = true
		name3, _ := jsonpath.Get("$[0].name", d.Get("config"))
		if name3 != nil && (d.HasChange("config.0.name") || name3 != "") {
			dataList["Name"] = name3
		}
	}

	dataList["Id"] = parts[1]
	if d.HasChange("waf_rule_id") {
		update = true
	}
	if v, ok := d.GetOk("waf_rule_id"); ok {
		dataList["Id"] = v
	}

	if d.HasChange("config") {
		update = true
		if v := d.Get("config"); v != nil {
			actions := make(map[string]interface{})
			bypass := make(map[string]interface{})
			skip1, _ := jsonpath.Get("$[0].actions[0].bypass[0].skip", d.Get("config"))
			if skip1 != nil && (d.HasChange("config.0.actions.0.bypass.0.skip") || skip1 != "") {
				bypass["Skip"] = skip1
			}
			customRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].custom_rules", d.Get("config"))
			if customRules1 != nil && (d.HasChange("config.0.actions.0.bypass.0.custom_rules") || customRules1 != "") {
				bypass["CustomRules"] = convertToInterfaceArray(customRules1)
			}
			regularRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_rules", d.Get("config"))
			if regularRules1 != nil && (d.HasChange("config.0.actions.0.bypass.0.regular_rules") || regularRules1 != "") {
				bypass["RegularRules"] = convertToInterfaceArray(regularRules1)
			}
			regularTypes1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_types", d.Get("config"))
			if regularTypes1 != nil && (d.HasChange("config.0.actions.0.bypass.0.regular_types") || regularTypes1 != "") {
				bypass["RegularTypes"] = convertToInterfaceArray(regularTypes1)
			}
			tags1, _ := jsonpath.Get("$[0].actions[0].bypass[0].tags", d.Get("config"))
			if tags1 != nil && (d.HasChange("config.0.actions.0.bypass.0.tags") || tags1 != "") {
				bypass["Tags"] = convertToInterfaceArray(tags1)
			}

			actions["Bypass"] = bypass
			response := make(map[string]interface{})
			code3, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("config"))
			if code3 != nil && (d.HasChange("config.0.actions.0.response.0.code") || code3 != "") {
				response["Code"] = code3
			}
			id4, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("config"))
			if id4 != nil && (d.HasChange("config.0.actions.0.response.0.id") || id4 != "") {
				response["Id"] = id4
			}

			actions["Response"] = response

			dataList["Actions"] = actions
		}
	}

	ConfigsMap := make([]interface{}, 0)
	ConfigsMap = append(ConfigsMap, dataList)
	dataListJson, err := json.Marshal(ConfigsMap)
	if err != nil {
		return WrapError(err)
	}
	request["Configs"] = string(dataListJson)

	if v, ok := d.GetOk("site_version"); ok {
		request["SiteVersion"] = v
	}
	dataList1 := make(map[string]interface{})

	if v := d.Get("shared"); v != nil {
		name5, _ := jsonpath.Get("$[0].name", v)
		if name5 != nil && (d.HasChange("shared.0.name") || name5 != "") {
			dataList1["Name"] = name5
		}
		action7, _ := jsonpath.Get("$[0].action", v)
		if action7 != nil && (d.HasChange("shared.0.action") || action7 != "") {
			dataList1["Action"] = action7
		}
		crossSiteId1, _ := jsonpath.Get("$[0].cross_site_id", v)
		if crossSiteId1 != nil && (d.HasChange("shared.0.cross_site_id") || crossSiteId1 != "") {
			dataList1["CrossSiteId"] = crossSiteId1
		}
		expression3, _ := jsonpath.Get("$[0].expression", v)
		if expression3 != nil && (d.HasChange("shared.0.expression") || expression3 != "") {
			dataList1["Expression"] = expression3
		}
		target1, _ := jsonpath.Get("$[0].target", v)
		if target1 != nil && (d.HasChange("shared.0.target") || target1 != "") {
			dataList1["Target"] = target1
		}
		match := make(map[string]interface{})
		if v, ok := d.GetOk("shared"); ok {
			localData9, err := jsonpath.Get("$[0].match[0].criteria", v)
			if err != nil {
				localData9 = make([]interface{}, 0)
			}
			localMaps9 := make([]interface{}, 0)
			for _, dataLoop9 := range convertToInterfaceArray(localData9) {
				dataLoop9Tmp := make(map[string]interface{})
				if dataLoop9 != nil {
					dataLoop9Tmp = dataLoop9.(map[string]interface{})
				}
				dataLoop9Map := make(map[string]interface{})
				dataLoop9Map["Logic"] = dataLoop9Tmp["logic"]
				dataLoop9Map["MatchType"] = dataLoop9Tmp["match_type"]
				localMaps10 := make([]interface{}, 0)
				localData10 := dataLoop9Tmp["criteria"]
				for _, dataLoop10 := range convertToInterfaceArray(localData10) {
					dataLoop10Tmp := dataLoop10.(map[string]interface{})
					dataLoop10Map := make(map[string]interface{})
					localMaps11 := make([]interface{}, 0)
					localData11 := dataLoop10Tmp["criteria"]
					for _, dataLoop11 := range convertToInterfaceArray(localData11) {
						dataLoop11Tmp := dataLoop11.(map[string]interface{})
						dataLoop11Map := make(map[string]interface{})
						dataLoop11Map["MatchType"] = dataLoop11Tmp["match_type"]
						localMaps11 = append(localMaps11, dataLoop11Map)
					}
					dataLoop10Map["Criteria"] = localMaps11
					dataLoop10Map["Logic"] = dataLoop10Tmp["logic"]
					dataLoop10Map["MatchType"] = dataLoop10Tmp["match_type"]
					localMaps10 = append(localMaps10, dataLoop10Map)
				}
				dataLoop9Map["Criteria"] = localMaps10
				localMaps9 = append(localMaps9, dataLoop9Map)
			}
			match["Criteria"] = localMaps9
		}

		matchType13, _ := jsonpath.Get("$[0].match[0].match_type", d.Get("shared"))
		if matchType13 != nil && (d.HasChange("shared.0.match.0.match_type") || matchType13 != "") {
			match["MatchType"] = matchType13
		}
		logic11, _ := jsonpath.Get("$[0].match[0].logic", d.Get("shared"))
		if logic11 != nil && (d.HasChange("shared.0.match.0.logic") || logic11 != "") {
			match["Logic"] = logic11
		}

		dataList1["Match"] = match
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && (d.HasChange("shared.0.mode") || mode1 != "") {
			dataList1["Mode"] = mode1
		}
		actions1 := make(map[string]interface{})
		response1 := make(map[string]interface{})
		code5, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("shared"))
		if code5 != nil && (d.HasChange("shared.0.actions.0.response.0.code") || code5 != "") {
			response1["Code"] = code5
		}
		id6, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("shared"))
		if id6 != nil && (d.HasChange("shared.0.actions.0.response.0.id") || id6 != "") {
			response1["Id"] = id6
		}

		actions1["Response"] = response1

		dataList1["Actions"] = actions1

		dataList1Json, err := json.Marshal(dataList1)
		if err != nil {
			return WrapError(err)
		}
		request["Shared"] = string(dataList1Json)
	}

	if !d.IsNewResource() && d.HasChange("ruleset_id") {
		update = true
	}
	request["RulesetId"] = d.Get("ruleset_id")

	if !d.IsNewResource() && d.HasChange("phase") {
		update = true
	}
	request["Phase"] = d.Get("phase")
	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Configs.0.Id", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

	if update && enableBatchUpdateWafRules {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	enableUpdateWafRule := false
	checkValue00 = objectRaw["Phase"]
	if !(InArray(fmt.Sprint(checkValue00), []string{"http_anti_scan", "http_bot"})) {
		enableUpdateWafRule = true
	}
	parts = strings.Split(d.Id(), ":")
	action = "UpdateWafRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	if !d.IsNewResource() && d.HasChange("config") {
		update = true
		dataList := make(map[string]interface{})

		if v := d.Get("config"); v != nil {
			if v, ok := d.GetOk("config"); ok {
				localData, err := jsonpath.Get("$[0].managed_rulesets", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(localData) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["ProtectionLevel"] = dataLoopTmp["protection_level"]
					dataLoopMap["AttackType"] = dataLoopTmp["attack_type"]
					dataLoopMap["Action"] = dataLoopTmp["action"]
					localMaps1 := make([]interface{}, 0)
					localData1 := dataLoopTmp["managed_rules"]
					for _, dataLoop1 := range convertToInterfaceArray(localData1) {
						dataLoop1Tmp := dataLoop1.(map[string]interface{})
						dataLoop1Map := make(map[string]interface{})
						dataLoop1Map["Action"] = dataLoop1Tmp["action"]
						dataLoop1Map["Status"] = dataLoop1Tmp["status"]
						dataLoop1Map["Id"] = dataLoop1Tmp["id"]
						localMaps1 = append(localMaps1, dataLoop1Map)
					}
					dataLoopMap["ManagedRules"] = localMaps1
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList["ManagedRulesets"] = localMaps
			}

			sigchl1, _ := jsonpath.Get("$[0].sigchl", v)
			if sigchl1 != nil && (d.HasChange("config.0.sigchl") || sigchl1 != "") {
				dataList["Sigchl"] = convertToInterfaceArray(sigchl1)
			}
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && (d.HasChange("config.0.type") || type1 != "") {
				dataList["Type"] = type1
			}
			rateLimit := make(map[string]interface{})
			threshold := make(map[string]interface{})
			managedRulesBlocked1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].managed_rules_blocked", d.Get("config"))
			if managedRulesBlocked1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.managed_rules_blocked") || managedRulesBlocked1 != "") {
				threshold["ManagedRulesBlocked"] = managedRulesBlocked1
			}
			traffic1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].traffic", d.Get("config"))
			if traffic1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.traffic") || traffic1 != "") {
				threshold["Traffic"] = traffic1
			}
			responseStatus := make(map[string]interface{})
			ratio1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].ratio", d.Get("config"))
			if ratio1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.ratio") || ratio1 != "") {
				responseStatus["Ratio"] = ratio1
			}
			count1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].count", d.Get("config"))
			if count1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.count") || count1 != "") {
				responseStatus["Count"] = count1
			}
			code1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].response_status[0].code", d.Get("config"))
			if code1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.response_status.0.code") || code1 != "") {
				responseStatus["Code"] = code1
			}

			threshold["ResponseStatus"] = responseStatus
			request1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].request", d.Get("config"))
			if request1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.request") || request1 != "") {
				threshold["Request"] = request1
			}
			distinctManagedRules1, _ := jsonpath.Get("$[0].rate_limit[0].threshold[0].distinct_managed_rules", d.Get("config"))
			if distinctManagedRules1 != nil && (d.HasChange("config.0.rate_limit.0.threshold.0.distinct_managed_rules") || distinctManagedRules1 != "") {
				threshold["DistinctManagedRules"] = distinctManagedRules1
			}

			rateLimit["Threshold"] = threshold
			characteristics := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData2, err := jsonpath.Get("$[0].rate_limit[0].characteristics[0].criteria", v)
				if err != nil {
					localData2 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := make(map[string]interface{})
					if dataLoop2 != nil {
						dataLoop2Tmp = dataLoop2.(map[string]interface{})
					}
					dataLoop2Map := make(map[string]interface{})
					localMaps3 := make([]interface{}, 0)
					localData3 := dataLoop2Tmp["criteria"]
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := dataLoop3.(map[string]interface{})
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["Logic"] = dataLoop3Tmp["logic"]
						localMaps4 := make([]interface{}, 0)
						localData4 := dataLoop3Tmp["criteria"]
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := dataLoop4.(map[string]interface{})
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["MatchType"] = dataLoop4Tmp["match_type"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						dataLoop3Map["Criteria"] = localMaps4
						dataLoop3Map["MatchType"] = dataLoop3Tmp["match_type"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					dataLoop2Map["Criteria"] = localMaps3
					dataLoop2Map["Logic"] = dataLoop2Tmp["logic"]
					dataLoop2Map["MatchType"] = dataLoop2Tmp["match_type"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				characteristics["Criteria"] = localMaps2
			}

			logic5, _ := jsonpath.Get("$[0].rate_limit[0].characteristics[0].logic", d.Get("config"))
			if logic5 != nil && (d.HasChange("config.0.rate_limit.0.characteristics.0.logic") || logic5 != "") {
				characteristics["Logic"] = logic5
			}

			rateLimit["Characteristics"] = characteristics
			interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
			if interval1 != nil && (d.HasChange("config.0.rate_limit.0.interval") || interval1 != "") {
				rateLimit["Interval"] = interval1
			}
			onHit1, _ := jsonpath.Get("$[0].rate_limit[0].on_hit", d.Get("config"))
			if onHit1 != nil && (d.HasChange("config.0.rate_limit.0.on_hit") || onHit1 != "") {
				rateLimit["OnHit"] = onHit1
			}
			ttl, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
			if ttl != nil && (d.HasChange("config.0.rate_limit.0.ttl") || ttl != "") {
				rateLimit["TTL"] = ttl
			}

			dataList["RateLimit"] = rateLimit
			timer := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData5, err := jsonpath.Get("$[0].timer[0].periods", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps5 := make([]interface{}, 0)
				for _, dataLoop5 := range convertToInterfaceArray(localData5) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["End"] = dataLoop5Tmp["end"]
					dataLoop5Map["Start"] = dataLoop5Tmp["start"]
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				timer["Periods"] = localMaps5
			}

			if v, ok := d.GetOk("config"); ok {
				localData6, err := jsonpath.Get("$[0].timer[0].weekly_periods", v)
				if err != nil {
					localData6 = make([]interface{}, 0)
				}
				localMaps6 := make([]interface{}, 0)
				for _, dataLoop6 := range convertToInterfaceArray(localData6) {
					dataLoop6Tmp := make(map[string]interface{})
					if dataLoop6 != nil {
						dataLoop6Tmp = dataLoop6.(map[string]interface{})
					}
					dataLoop6Map := make(map[string]interface{})
					localMaps7 := make([]interface{}, 0)
					localData7 := dataLoop6Tmp["daily_periods"]
					for _, dataLoop7 := range convertToInterfaceArray(localData7) {
						dataLoop7Tmp := dataLoop7.(map[string]interface{})
						dataLoop7Map := make(map[string]interface{})
						dataLoop7Map["End"] = dataLoop7Tmp["end"]
						dataLoop7Map["Start"] = dataLoop7Tmp["start"]
						localMaps7 = append(localMaps7, dataLoop7Map)
					}
					dataLoop6Map["DailyPeriods"] = localMaps7
					dataLoop6Map["Days"] = dataLoop6Tmp["days"]
					localMaps6 = append(localMaps6, dataLoop6Map)
				}
				timer["WeeklyPeriods"] = localMaps6
			}

			zone1, _ := jsonpath.Get("$[0].timer[0].zone", d.Get("config"))
			if zone1 != nil && (d.HasChange("config.0.timer.0.zone") || zone1 != "") {
				timer["Zone"] = zone1
			}
			scopes1, _ := jsonpath.Get("$[0].timer[0].scopes", d.Get("config"))
			if scopes1 != nil && (d.HasChange("config.0.timer.0.scopes") || scopes1 != "") {
				timer["Scopes"] = scopes1
			}

			dataList["Timer"] = timer
			notes1, _ := jsonpath.Get("$[0].notes", v)
			if notes1 != nil && (d.HasChange("config.0.notes") || notes1 != "") {
				dataList["Notes"] = notes1
			}
			action5, _ := jsonpath.Get("$[0].action", v)
			if action5 != nil && (d.HasChange("config.0.action") || action5 != "") {
				dataList["Action"] = action5
			}
			appSdk := make(map[string]interface{})
			customSign := make(map[string]interface{})
			value1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].value", d.Get("config"))
			if value1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign.0.value") || value1 != "") {
				customSign["Value"] = value1
			}
			key1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign[0].key", d.Get("config"))
			if key1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign.0.key") || key1 != "") {
				customSign["Key"] = key1
			}

			appSdk["CustomSign"] = customSign
			customSignStatus1, _ := jsonpath.Get("$[0].app_sdk[0].custom_sign_status", d.Get("config"))
			if customSignStatus1 != nil && (d.HasChange("config.0.app_sdk.0.custom_sign_status") || customSignStatus1 != "") {
				appSdk["CustomSignStatus"] = customSignStatus1
			}
			featureAbnormal1, _ := jsonpath.Get("$[0].app_sdk[0].feature_abnormal", d.Get("config"))
			if featureAbnormal1 != nil && (d.HasChange("config.0.app_sdk.0.feature_abnormal") || featureAbnormal1 != "") {
				appSdk["FeatureAbnormal"] = convertToInterfaceArray(featureAbnormal1)
			}

			dataList["AppSdk"] = appSdk
			securityLevel := make(map[string]interface{})
			value3, _ := jsonpath.Get("$[0].security_level[0].value", d.Get("config"))
			if value3 != nil && (d.HasChange("config.0.security_level.0.value") || value3 != "") {
				securityLevel["Value"] = value3
			}

			dataList["SecurityLevel"] = securityLevel
			value5, _ := jsonpath.Get("$[0].value", v)
			if value5 != nil && (d.HasChange("config.0.value") || value5 != "") {
				dataList["Value"] = value5
			}
			expression1, _ := jsonpath.Get("$[0].expression", v)
			if expression1 != nil && (d.HasChange("config.0.expression") || expression1 != "") {
				dataList["Expression"] = expression1
			}
			managedGroupId1, _ := jsonpath.Get("$[0].managed_group_id", v)
			if managedGroupId1 != nil && (d.HasChange("config.0.managed_group_id") || managedGroupId1 != "") {
				dataList["ManagedGroupId"] = managedGroupId1
			}
			managedList1, _ := jsonpath.Get("$[0].managed_list", v)
			if managedList1 != nil && (d.HasChange("config.0.managed_list") || managedList1 != "") {
				dataList["ManagedList"] = managedList1
			}
			appPackage := make(map[string]interface{})
			if v, ok := d.GetOk("config"); ok {
				localData8, err := jsonpath.Get("$[0].app_package[0].package_signs", v)
				if err != nil {
					localData8 = make([]interface{}, 0)
				}
				localMaps8 := make([]interface{}, 0)
				for _, dataLoop8 := range convertToInterfaceArray(localData8) {
					dataLoop8Tmp := make(map[string]interface{})
					if dataLoop8 != nil {
						dataLoop8Tmp = dataLoop8.(map[string]interface{})
					}
					dataLoop8Map := make(map[string]interface{})
					dataLoop8Map["Sign"] = dataLoop8Tmp["sign"]
					dataLoop8Map["Name"] = dataLoop8Tmp["name"]
					localMaps8 = append(localMaps8, dataLoop8Map)
				}
				appPackage["PackageSigns"] = localMaps8
			}

			dataList["AppPackage"] = appPackage
			name3, _ := jsonpath.Get("$[0].name", v)
			if name3 != nil && (d.HasChange("config.0.name") || name3 != "") {
				dataList["Name"] = name3
			}
			actions := make(map[string]interface{})
			bypass := make(map[string]interface{})
			skip1, _ := jsonpath.Get("$[0].actions[0].bypass[0].skip", d.Get("config"))
			if skip1 != nil && (d.HasChange("config.0.actions.0.bypass.0.skip") || skip1 != "") {
				bypass["Skip"] = skip1
			}
			customRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].custom_rules", d.Get("config"))
			if customRules1 != nil && (d.HasChange("config.0.actions.0.bypass.0.custom_rules") || customRules1 != "") {
				bypass["CustomRules"] = convertToInterfaceArray(customRules1)
			}
			regularRules1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_rules", d.Get("config"))
			if regularRules1 != nil && (d.HasChange("config.0.actions.0.bypass.0.regular_rules") || regularRules1 != "") {
				bypass["RegularRules"] = convertToInterfaceArray(regularRules1)
			}
			regularTypes1, _ := jsonpath.Get("$[0].actions[0].bypass[0].regular_types", d.Get("config"))
			if regularTypes1 != nil && (d.HasChange("config.0.actions.0.bypass.0.regular_types") || regularTypes1 != "") {
				bypass["RegularTypes"] = convertToInterfaceArray(regularTypes1)
			}
			tags1, _ := jsonpath.Get("$[0].actions[0].bypass[0].tags", d.Get("config"))
			if tags1 != nil && (d.HasChange("config.0.actions.0.bypass.0.tags") || tags1 != "") {
				bypass["Tags"] = convertToInterfaceArray(tags1)
			}

			actions["Bypass"] = bypass
			response := make(map[string]interface{})
			code3, _ := jsonpath.Get("$[0].actions[0].response[0].code", d.Get("config"))
			if code3 != nil && (d.HasChange("config.0.actions.0.response.0.code") || code3 != "") {
				response["Code"] = code3
			}
			id3, _ := jsonpath.Get("$[0].actions[0].response[0].id", d.Get("config"))
			if id3 != nil && (d.HasChange("config.0.actions.0.response.0.id") || id3 != "") {
				response["Id"] = id3
			}

			actions["Response"] = response

			dataList["Actions"] = actions

			dataListJson, err := json.Marshal(dataList)
			if err != nil {
				return WrapError(err)
			}
			request["Config"] = string(dataListJson)
		}
	}

	if v, ok := d.GetOk("site_version"); ok {
		request["SiteVersion"] = v
	}
	if !d.IsNewResource() && d.HasChange("config.0.status") {
		update = true
		configStatusJsonPath, err := jsonpath.Get("$[0].status", d.Get("config"))
		if err == nil {
			request["Status"] = configStatusJsonPath
		}
	}

	if update && enableUpdateWafRule {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudEsaWafRuleRead(d, meta)
}

func resourceAliCloudEsaWafRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteWafRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = parts[1]
	request["SiteId"] = parts[0]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
