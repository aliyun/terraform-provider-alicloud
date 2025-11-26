package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbRuleCreate,
		Read:   resourceAliCloudAlbRuleRead,
		Update: resourceAliCloudAlbRuleUpdate,
		Delete: resourceAliCloudAlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Request", "Response"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rule_actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntBetween(1, 50000),
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"ForwardGroup", "Redirect", "FixedResponse", "Rewrite", "InsertHeader", "TrafficLimit", "TrafficMirror", "Cors", "RemoveHeader"}, false),
						},
						"fixed_response_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"text/plain", "text/css", "text/html", "application/javascript", "application/json"}, false),
									},
									"http_code": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^[2-5][0-9]{2}$`), "The http code must be an HTTP_2xx,HTTP_4xx or HTTP_5xx.x is a digit."),
									},
								},
							},
						},
						"forward_group_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_id": {
													Type:     schema.TypeString,
													Computed: true,
													Optional: true,
												},
												"weight": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      100,
													ValidateFunc: IntBetween(0, 100),
												},
											},
										},
									},
									"server_group_sticky_session": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"timeout": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													ValidateFunc: IntBetween(1, 86400),
												},
											},
										},
									},
								},
							},
						},
						"insert_header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
									"value": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringLenBetween(1, 128),
									},
									"value_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"UserDefined", "ReferenceHeader", "SystemDefined"}, false),
									},
								},
							},
						},
						"remove_header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
								},
							},
						},
						"redirect_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^.{3,128}$`), "The host name must be 3 to 128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?)."),
									},
									"http_code": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"301", "302", "303", "307", "308"}, false),
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "${protocol}"}, false),
									},
									"query": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringLenBetween(1, 128),
									},
								},
							},
						},
						"rewrite_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^.{3,128}$`), "The host name must be 3 to 128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?)."),
									},
									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringLenBetween(1, 128),
									},
									"query": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringLenBetween(1, 128),
									},
								},
							},
						},
						"traffic_limit_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"qps": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(1, 100000),
									},
									"per_ip_qps": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(1, 100000),
									},
								},
							},
						},
						"traffic_mirror_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mirror_group_config": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_tuples": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"server_group_id": {
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
						"cors_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_origin": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allow_methods": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allow_headers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"expose_headers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allow_credentials": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"on", "off"}, false),
									},
									"max_age": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(-1, 172800),
									},
								},
							},
						},
					},
				},
			},
			"rule_conditions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Host", "Path", "Header", "QueryString", "Method", "Cookie", "SourceIp", "ResponseHeader", "ResponseStatusCode"}, false),
						},
						"cookie_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringLenBetween(1, 100),
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringLenBetween(1, 128),
												},
											},
										},
									},
								},
							},
						},
						"header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"response_header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"response_status_code_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"host_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"method_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"path_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"query_string_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringLenBetween(1, 100),
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringLenBetween(1, 128),
												},
											},
										},
									},
								},
							},
						},
						"source_ip_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 5,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlbRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	var direction string
	action := "CreateRule"
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = buildClientToken("CreateRule")
	request["ListenerId"] = d.Get("listener_id")
	request["RuleName"] = d.Get("rule_name")
	request["Priority"] = d.Get("priority")

	if v, ok := d.GetOk("direction"); ok {
		direction = v.(string)
		request["Direction"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	ruleActionsMaps := make([]map[string]interface{}, 0)
	for _, ruleActions := range d.Get("rule_actions").(*schema.Set).List() {
		ruleActionsArg := ruleActions.(map[string]interface{})
		ruleActionsMap := map[string]interface{}{}
		ruleActionsMap["Order"] = ruleActionsArg["order"]
		ruleActionsMap["Type"] = ruleActionsArg["type"]

		if ruleActionsMap["Type"] == "" {
			continue
		}

		fixedResponseConfigMap := map[string]interface{}{}
		for _, fixedResponseConfig := range ruleActionsArg["fixed_response_config"].(*schema.Set).List() {
			fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
			fixedResponseConfigMap["Content"] = fixedResponseConfigArg["content"]
			fixedResponseConfigMap["ContentType"] = fixedResponseConfigArg["content_type"]
			fixedResponseConfigMap["HttpCode"] = fixedResponseConfigArg["http_code"]
			ruleActionsMap["FixedResponseConfig"] = fixedResponseConfigMap
		}

		forwardGroupConfigMap := map[string]interface{}{}
		for _, forwardGroupConfig := range ruleActionsArg["forward_group_config"].([]interface{}) {
			forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
			serverGroupTuplesMaps := make([]map[string]interface{}, 0)
			for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
				serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
				serverGroupTuplesMap := map[string]interface{}{}
				serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
				if v, ok := serverGroupTuplesArg["weight"]; ok {
					serverGroupTuplesMap["Weight"] = v
				}

				serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
			}
			forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps

			serverGroupStickySessionMap := map[string]interface{}{}
			for _, serverGroupStickySession := range forwardGroupConfigArg["server_group_sticky_session"].(*schema.Set).List() {
				serverGroupStickySessionArg := serverGroupStickySession.(map[string]interface{})
				serverGroupStickySessionMap["Enabled"] = serverGroupStickySessionArg["enabled"]
				serverGroupStickySessionMap["Timeout"] = serverGroupStickySessionArg["timeout"]
			}
			forwardGroupConfigMap["ServerGroupStickySession"] = serverGroupStickySessionMap

			ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		insertHeaderConfigMap := map[string]interface{}{}
		for _, insertHeaderConfig := range ruleActionsArg["insert_header_config"].(*schema.Set).List() {
			insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
			insertHeaderConfigMap["Key"] = insertHeaderConfigArg["key"]
			insertHeaderConfigMap["Value"] = insertHeaderConfigArg["value"]
			insertHeaderConfigMap["ValueType"] = insertHeaderConfigArg["value_type"]
			ruleActionsMap["InsertHeaderConfig"] = insertHeaderConfigMap
		}

		removeHeaderConfigMap := map[string]interface{}{}
		for _, removeHeaderConfig := range ruleActionsArg["remove_header_config"].(*schema.Set).List() {
			removeHeaderConfigArg := removeHeaderConfig.(map[string]interface{})
			insertHeaderConfigMap["Key"] = removeHeaderConfigArg["key"]
			ruleActionsMap["RemoveHeaderConfig"] = removeHeaderConfigMap
		}

		redirectConfigMap := map[string]interface{}{}
		for _, redirectConfig := range ruleActionsArg["redirect_config"].(*schema.Set).List() {
			redirectConfigArg := redirectConfig.(map[string]interface{})
			redirectConfigMap["Host"] = redirectConfigArg["host"]
			redirectConfigMap["HttpCode"] = redirectConfigArg["http_code"]
			redirectConfigMap["Path"] = redirectConfigArg["path"]
			redirectConfigMap["Port"] = redirectConfigArg["port"]
			redirectConfigMap["Protocol"] = redirectConfigArg["protocol"]
			redirectConfigMap["Query"] = redirectConfigArg["query"]
			ruleActionsMap["RedirectConfig"] = redirectConfigMap
		}

		rewriteConfigMap := map[string]interface{}{}
		for _, rewriteConfig := range ruleActionsArg["rewrite_config"].(*schema.Set).List() {
			rewriteConfigArg := rewriteConfig.(map[string]interface{})
			rewriteConfigMap["Host"] = rewriteConfigArg["host"]
			rewriteConfigMap["Path"] = rewriteConfigArg["path"]
			rewriteConfigMap["Query"] = rewriteConfigArg["query"]
			ruleActionsMap["RewriteConfig"] = rewriteConfigMap
		}

		trafficLimitConfigList := ruleActionsArg["traffic_limit_config"].(*schema.Set).List()
		if len(trafficLimitConfigList) > 0 {
			trafficLimitConfigArg := trafficLimitConfigList[0].(map[string]interface{})
			ruleActionsMap["TrafficLimitConfig"] = map[string]interface{}{
				"QPS":      trafficLimitConfigArg["qps"],
				"PerIpQps": trafficLimitConfigArg["per_ip_qps"],
			}
		}

		trafficMirrorConfigList := ruleActionsArg["traffic_mirror_config"].(*schema.Set).List()
		if len(trafficMirrorConfigList) > 0 {
			trafficMirrorConfigArg := trafficMirrorConfigList[0].(map[string]interface{})
			mirrorGroupConfigMap := make(map[string]interface{}, 0)
			mirrorGroupConfigList := trafficMirrorConfigArg["mirror_group_config"].(*schema.Set).List()
			if len(mirrorGroupConfigList) > 0 {
				mirrorGroupConfigArg := mirrorGroupConfigList[0].(map[string]interface{})
				serverGroupTuplesMaps := make([]map[string]interface{}, 0)
				for _, serverGroupTuples := range mirrorGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
					serverGroupTuplesMap := map[string]interface{}{}
					serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
					serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
				mirrorGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
			}
			ruleActionsMap["TrafficMirrorConfig"] = map[string]interface{}{
				"MirrorGroupConfig": mirrorGroupConfigMap,
				"TargetType":        trafficMirrorConfigArg["target_type"],
			}
		}

		if v, ok := ruleActionsArg["cors_config"]; ok {
			corsConfigMap := map[string]interface{}{}
			for _, corsConfigList := range v.(*schema.Set).List() {
				corsConfigArg := corsConfigList.(map[string]interface{})

				if allowOrigin, ok := corsConfigArg["allow_origin"]; ok {
					corsConfigMap["AllowOrigin"] = allowOrigin
				}

				if allowMethods, ok := corsConfigArg["allow_methods"]; ok {
					corsConfigMap["AllowMethods"] = allowMethods
				}

				if allowHeaders, ok := corsConfigArg["allow_headers"]; ok {
					corsConfigMap["AllowHeaders"] = allowHeaders
				}

				if exposeHeaders, ok := corsConfigArg["expose_headers"]; ok {
					corsConfigMap["ExposeHeaders"] = exposeHeaders
				}

				if allowCredentials, ok := corsConfigArg["allow_credentials"]; ok {
					corsConfigMap["AllowCredentials"] = allowCredentials
				}

				if maxAge, ok := corsConfigArg["max_age"]; ok {
					corsConfigMap["MaxAge"] = maxAge
				}
			}

			ruleActionsMap["CorsConfig"] = corsConfigMap
		}

		ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
	}
	request["RuleActions"] = ruleActionsMaps

	ruleConditionsMaps := make([]map[string]interface{}, 0)
	for _, ruleConditions := range d.Get("rule_conditions").(*schema.Set).List() {
		ruleConditionsArg := ruleConditions.(map[string]interface{})
		ruleConditionsMap := map[string]interface{}{}
		ruleConditionsMap["Type"] = ruleConditionsArg["type"]

		cookieConfigMap := map[string]interface{}{}
		for _, cookieConfig := range ruleConditionsArg["cookie_config"].(*schema.Set).List() {
			cookieConfigArg := cookieConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range cookieConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			cookieConfigMap["Values"] = valuesMaps
			ruleConditionsMap["CookieConfig"] = cookieConfigMap
		}

		headerConfigMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["header_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			headerConfigMap["Key"] = headerConfigArg["key"]
			headerConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HeaderConfig"] = headerConfigMap
		}

		responseHeaderConfigMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["response_header_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			headerConfigMap["Key"] = headerConfigArg["key"]
			headerConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["ResponseHeaderConfig"] = responseHeaderConfigMap
		}

		responseStatusCodeMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["response_status_code_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			responseStatusCodeMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["ResponseStatusCodeConfig"] = responseStatusCodeMap
		}

		hostConfigMap := map[string]interface{}{}
		for _, hostConfig := range ruleConditionsArg["host_config"].(*schema.Set).List() {
			hostConfigArg := hostConfig.(map[string]interface{})
			hostConfigMap["Values"] = hostConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HostConfig"] = hostConfigMap
		}

		methodConfigMap := map[string]interface{}{}
		for _, methodConfig := range ruleConditionsArg["method_config"].(*schema.Set).List() {
			methodConfigArg := methodConfig.(map[string]interface{})
			methodConfigMap["Values"] = methodConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["MethodConfig"] = methodConfigMap
		}

		pathConfigMap := map[string]interface{}{}
		for _, pathConfig := range ruleConditionsArg["path_config"].(*schema.Set).List() {
			pathConfigArg := pathConfig.(map[string]interface{})
			pathConfigMap["Values"] = pathConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["PathConfig"] = pathConfigMap
		}

		sourceIpConfigMap := map[string]interface{}{}
		for _, sourceIpConfig := range ruleConditionsArg["source_ip_config"].(*schema.Set).List() {
			sourceIpConfigArg := sourceIpConfig.(map[string]interface{})
			sourceIpConfigMap["Values"] = sourceIpConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["SourceIpConfig"] = sourceIpConfigMap
		}

		queryStringConfigMap := map[string]interface{}{}
		for _, queryStringConfig := range ruleConditionsArg["query_string_config"].(*schema.Set).List() {
			queryStringConfigArg := queryStringConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range queryStringConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			queryStringConfigMap["Values"] = valuesMaps
			ruleConditionsMap["QueryStringConfig"] = queryStringConfigMap
		}

		ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
	}
	request["RuleConditions"] = ruleConditionsMaps

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing", "IncorrectStatus.Listener", "ResourceInConfiguring.Listener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RuleId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbRuleStateRefreshFunc(d.Id(), direction, []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbRuleRead(d, meta)
}

func resourceAliCloudAlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}

	var direction string
	if v, ok := d.GetOk("direction"); ok {
		direction = v.(string)
	}

	object, err := albService.DescribeAlbRule(d.Id(), direction)
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_rule albService.DescribeAlbRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("listener_id", object["ListenerId"])
	d.Set("rule_name", object["RuleName"])
	d.Set("direction", object["Direction"])
	d.Set("status", object["RuleStatus"])

	if v, ok := object["Priority"]; ok && fmt.Sprint(v) != "0" {
		d.Set("priority", formatInt(v))
	}

	if ruleActionsList, ok := object["RuleActions"]; ok {
		ruleActionsMaps := make([]map[string]interface{}, 0)
		for _, ruleActions := range ruleActionsList.([]interface{}) {
			ruleActionsArg := ruleActions.(map[string]interface{})
			ruleActionsMap := map[string]interface{}{}
			ruleActionsMap["type"] = ruleActionsArg["Type"]
			ruleActionsMap["order"] = formatInt(ruleActionsArg["Order"])

			if forwardGroupConfig, ok := ruleActionsArg["ForwardGroupConfig"]; ok {
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				if len(forwardGroupConfigArg) > 0 {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					if forwardGroupConfigArgs, ok := forwardGroupConfigArg["ServerGroupTuples"].([]interface{}); ok {
						for _, serverGroupTuples := range forwardGroupConfigArgs {
							serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
							serverGroupTuplesMap := map[string]interface{}{}
							serverGroupTuplesMap["server_group_id"] = serverGroupTuplesArg["ServerGroupId"]
							serverGroupTuplesMap["weight"] = formatInt(serverGroupTuplesArg["Weight"])
							serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
						}
					}

					serverGroupStickySessionMaps := make([]map[string]interface{}, 0)
					if serverGroupStickySessionArg, ok := forwardGroupConfigArg["ServerGroupStickySession"].(map[string]interface{}); ok && len(serverGroupStickySessionArg) > 0 {
						serverGroupStickySessionMap := map[string]interface{}{}
						if serverGroupStickySessionArgEnabled, ok := serverGroupStickySessionArg["Enabled"]; ok {
							serverGroupStickySessionMap["enabled"] = serverGroupStickySessionArgEnabled
						}
						if serverGroupStickySessionArgTimeout, ok := serverGroupStickySessionArg["Timeout"]; ok {
							serverGroupStickySessionMap["timeout"] = serverGroupStickySessionArgTimeout
						}
						serverGroupStickySessionMaps = append(serverGroupStickySessionMaps, serverGroupStickySessionMap)
					}

					if len(serverGroupTuplesMaps) > 0 {
						forwardGroupConfigMaps := make([]map[string]interface{}, 0)
						forwardGroupConfigMap := map[string]interface{}{}
						forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
						forwardGroupConfigMap["server_group_sticky_session"] = serverGroupStickySessionMaps
						forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
						ruleActionsMap["forward_group_config"] = forwardGroupConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}

				}
			}

			if fixedResponseConfig, ok := ruleActionsArg["FixedResponseConfig"]; ok {
				fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
				if len(fixedResponseConfigArg) > 0 {
					fixedResponseConfigMaps := make([]map[string]interface{}, 0)
					fixedResponseConfigMap := make(map[string]interface{}, 0)
					fixedResponseConfigMap["content"] = fixedResponseConfigArg["Content"]
					fixedResponseConfigMap["content_type"] = fixedResponseConfigArg["ContentType"]
					fixedResponseConfigMap["http_code"] = fixedResponseConfigArg["HttpCode"]
					fixedResponseConfigMaps = append(fixedResponseConfigMaps, fixedResponseConfigMap)
					ruleActionsMap["fixed_response_config"] = fixedResponseConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if insertHeaderConfig, ok := ruleActionsArg["InsertHeaderConfig"]; ok {
				insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
				if len(insertHeaderConfigArg) > 0 {
					insertHeaderConfigMaps := make([]map[string]interface{}, 0)
					insertHeaderConfigMap := make(map[string]interface{}, 0)
					insertHeaderConfigMap["key"] = insertHeaderConfigArg["Key"]
					insertHeaderConfigMap["value"] = insertHeaderConfigArg["Value"]
					insertHeaderConfigMap["value_type"] = insertHeaderConfigArg["ValueType"]
					insertHeaderConfigMaps = append(insertHeaderConfigMaps, insertHeaderConfigMap)
					ruleActionsMap["insert_header_config"] = insertHeaderConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if removeHeaderConfig, ok := ruleActionsArg["RemoveHeaderConfig"]; ok {
				removeHeaderConfigArg := removeHeaderConfig.(map[string]interface{})
				if len(removeHeaderConfigArg) > 0 {
					removeHeaderConfigMaps := make([]map[string]interface{}, 0)
					removeHeaderConfigMap := make(map[string]interface{}, 0)
					removeHeaderConfigMap["key"] = removeHeaderConfigArg["Key"]
					removeHeaderConfigMaps = append(removeHeaderConfigMaps, removeHeaderConfigMap)
					ruleActionsMap["remove_header_config"] = removeHeaderConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if redirectConfig, ok := ruleActionsArg["RedirectConfig"]; ok {
				redirectConfigArg := redirectConfig.(map[string]interface{})
				if len(redirectConfigArg) > 0 {
					redirectConfigMaps := make([]map[string]interface{}, 0)
					redirectConfigMap := make(map[string]interface{}, 0)
					redirectConfigMap["host"] = redirectConfigArg["Host"]
					redirectConfigMap["http_code"] = redirectConfigArg["HttpCode"]
					redirectConfigMap["path"] = redirectConfigArg["Path"]
					redirectConfigMap["port"] = redirectConfigArg["Port"]
					redirectConfigMap["protocol"] = redirectConfigArg["Protocol"]
					redirectConfigMap["query"] = redirectConfigArg["Query"]
					redirectConfigMaps = append(redirectConfigMaps, redirectConfigMap)
					ruleActionsMap["redirect_config"] = redirectConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if rewriteConfig, ok := ruleActionsArg["RewriteConfig"]; ok {
				rewriteConfigArg := rewriteConfig.(map[string]interface{})
				if len(rewriteConfigArg) > 0 {
					rewriteConfigMaps := make([]map[string]interface{}, 0)
					rewriteConfigMap := make(map[string]interface{}, 0)
					rewriteConfigMap["host"] = rewriteConfigArg["Host"]
					rewriteConfigMap["path"] = rewriteConfigArg["Path"]
					rewriteConfigMap["query"] = rewriteConfigArg["Query"]
					rewriteConfigMaps = append(rewriteConfigMaps, rewriteConfigMap)
					ruleActionsMap["rewrite_config"] = rewriteConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if trafficLimitConfig, ok := ruleActionsArg["TrafficLimitConfig"]; ok {
				trafficLimitConfigArg := trafficLimitConfig.(map[string]interface{})
				if len(trafficLimitConfigArg) > 0 {
					trafficLimitConfigMaps := make([]map[string]interface{}, 0)
					trafficLimitConfigMap := make(map[string]interface{}, 0)
					trafficLimitConfigMap["qps"] = trafficLimitConfigArg["QPS"]
					trafficLimitConfigMap["per_ip_qps"] = trafficLimitConfigArg["PerIpQps"]
					trafficLimitConfigMaps = append(trafficLimitConfigMaps, trafficLimitConfigMap)
					ruleActionsMap["traffic_limit_config"] = trafficLimitConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if trafficMirrorConfig, ok := ruleActionsArg["TrafficMirrorConfig"]; ok {
				if trafficMirrorConfigArg, ok := trafficMirrorConfig.(map[string]interface{}); ok && len(trafficMirrorConfigArg) > 0 {
					if mirrorGroupConfig, ok := trafficMirrorConfigArg["MirrorGroupConfig"]; ok {
						if mirrorGroupConfigArg, ok := mirrorGroupConfig.(map[string]interface{}); ok && len(mirrorGroupConfigArg) > 0 {
							if serverGroupTuples, ok := mirrorGroupConfigArg["ServerGroupTuples"].([]interface{}); ok {
								serverGroupTuplesMaps := make([]map[string]interface{}, 0)
								for _, item := range serverGroupTuples {
									serverGroupTuplesArg := item.(map[string]interface{})
									serverGroupTuplesMap := map[string]interface{}{}
									serverGroupTuplesMap["server_group_id"] = serverGroupTuplesArg["ServerGroupId"]
									serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
								}

								if len(serverGroupTuplesMaps) > 0 {
									mirrorGroupConfigMaps := make([]map[string]interface{}, 0)
									mirrorGroupConfigMaps = append(mirrorGroupConfigMaps, map[string]interface{}{
										"server_group_tuples": serverGroupTuplesMaps,
									})
									trafficMirrorConfigMaps := make([]map[string]interface{}, 0)
									trafficMirrorConfigMaps = append(trafficMirrorConfigMaps, map[string]interface{}{
										"mirror_group_config": mirrorGroupConfigMaps,
										"target_type":         trafficMirrorConfigArg["TargetType"],
									})
									ruleActionsMap["traffic_mirror_config"] = trafficMirrorConfigMaps
									ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
								}
							}
						}
					}
				}
			}

			if corsConfig, ok := ruleActionsArg["CorsConfig"]; ok {
				corsConfigMaps := make([]map[string]interface{}, 0)
				corsConfigArg := corsConfig.(map[string]interface{})

				if len(corsConfigArg) > 0 {
					corsConfigMap := map[string]interface{}{}

					if allowOrigin, ok := corsConfigArg["AllowOrigin"]; ok {
						corsConfigMap["allow_origin"] = allowOrigin
					}

					if allowMethods, ok := corsConfigArg["AllowMethods"]; ok {
						corsConfigMap["allow_methods"] = allowMethods
					}

					if allowHeaders, ok := corsConfigArg["AllowHeaders"]; ok {
						corsConfigMap["allow_headers"] = allowHeaders
					}

					if exposeHeaders, ok := corsConfigArg["ExposeHeaders"]; ok {
						corsConfigMap["expose_headers"] = exposeHeaders
					}

					if allowCredentials, ok := corsConfigArg["AllowCredentials"]; ok {
						corsConfigMap["allow_credentials"] = allowCredentials
					}

					if maxAge, ok := corsConfigArg["MaxAge"]; ok {
						corsConfigMap["max_age"] = maxAge
					}

					corsConfigMaps = append(corsConfigMaps, corsConfigMap)
					ruleActionsMap["cors_config"] = corsConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}
		}

		d.Set("rule_actions", ruleActionsMaps)
	}

	if ruleConditionsList, ok := object["RuleConditions"]; ok {
		ruleConditionsMaps := make([]map[string]interface{}, 0)
		for _, ruleConditions := range ruleConditionsList.([]interface{}) {
			ruleConditionsArg := ruleConditions.(map[string]interface{})
			ruleConditionsMap := map[string]interface{}{}
			ruleConditionsMap["type"] = ruleConditionsArg["Type"]

			if cookieConfig, ok := ruleConditionsArg["CookieConfig"]; ok {
				cookieConfigArg := cookieConfig.(map[string]interface{})
				if len(cookieConfigArg) > 0 {
					cookieConfigMaps := make([]map[string]interface{}, 0)
					valuesMaps := make([]map[string]interface{}, 0)
					for _, values := range cookieConfigArg["Values"].([]interface{}) {
						valuesArg := values.(map[string]interface{})
						valuesMap := map[string]interface{}{}
						valuesMap["key"] = valuesArg["Key"]
						valuesMap["value"] = valuesArg["Value"]
						valuesMaps = append(valuesMaps, valuesMap)
					}
					cookieConfigMap := map[string]interface{}{}
					cookieConfigMap["values"] = valuesMaps
					cookieConfigMaps = append(cookieConfigMaps, cookieConfigMap)
					ruleConditionsMap["cookie_config"] = cookieConfigMaps
				}
			}

			if headerConfig, ok := ruleConditionsArg["HeaderConfig"]; ok {
				headerConfigArg := headerConfig.(map[string]interface{})
				if len(headerConfigArg) > 0 {
					headerConfigMaps := make([]map[string]interface{}, 0)
					headerConfigMap := map[string]interface{}{}
					headerConfigMap["values"] = headerConfigArg["Values"].([]interface{})
					headerConfigMap["key"] = headerConfigArg["Key"]
					headerConfigMaps = append(headerConfigMaps, headerConfigMap)
					ruleConditionsMap["header_config"] = headerConfigMaps
				}
			}

			if headerConfig, ok := ruleConditionsArg["ResponseHeaderConfig"]; ok {
				headerConfigArg := headerConfig.(map[string]interface{})
				if len(headerConfigArg) > 0 {
					headerConfigMaps := make([]map[string]interface{}, 0)
					headerConfigMap := map[string]interface{}{}
					headerConfigMap["values"] = headerConfigArg["Values"].([]interface{})
					headerConfigMap["key"] = headerConfigArg["Key"]
					headerConfigMaps = append(headerConfigMaps, headerConfigMap)
					ruleConditionsMap["response_header_config"] = headerConfigMaps
				}
			}

			if headerConfig, ok := ruleConditionsArg["ResponseStatusCodeConfig"]; ok {
				headerConfigArg := headerConfig.(map[string]interface{})
				if len(headerConfigArg) > 0 {
					headerConfigMaps := make([]map[string]interface{}, 0)
					headerConfigMap := map[string]interface{}{}
					headerConfigMap["values"] = headerConfigArg["Values"].([]interface{})
					headerConfigMaps = append(headerConfigMaps, headerConfigMap)
					ruleConditionsMap["response_status_code_config"] = headerConfigMaps
				}
			}

			if queryStringConfig, ok := ruleConditionsArg["QueryStringConfig"]; ok {
				queryStringConfigArg := queryStringConfig.(map[string]interface{})
				if len(queryStringConfigArg) > 0 {
					queryStringConfigMaps := make([]map[string]interface{}, 0)
					queryStringValuesMaps := make([]map[string]interface{}, 0)
					for _, values := range queryStringConfigArg["Values"].([]interface{}) {
						valuesArg := values.(map[string]interface{})
						valuesMap := map[string]interface{}{}
						valuesMap["key"] = valuesArg["Key"]
						valuesMap["value"] = valuesArg["Value"]
						queryStringValuesMaps = append(queryStringValuesMaps, valuesMap)
					}
					queryStringConfigMap := map[string]interface{}{}
					queryStringConfigMap["values"] = queryStringValuesMaps
					queryStringConfigMaps = append(queryStringConfigMaps, queryStringConfigMap)
					ruleConditionsMap["query_string_config"] = queryStringConfigMaps
				}
			}

			if hostConfig, ok := ruleConditionsArg["HostConfig"]; ok {
				hostConfigArg := hostConfig.(map[string]interface{})
				if len(hostConfigArg) > 0 {
					hostConfigMaps := make([]map[string]interface{}, 0)
					hostConfigMap := map[string]interface{}{}
					hostConfigMap["values"] = hostConfigArg["Values"].([]interface{})
					hostConfigMaps = append(hostConfigMaps, hostConfigMap)
					ruleConditionsMap["host_config"] = hostConfigMaps
				}
			}

			if methodConfig, ok := ruleConditionsArg["MethodConfig"]; ok {
				methodConfigArg := methodConfig.(map[string]interface{})
				if len(methodConfigArg) > 0 {
					methodConfigMaps := make([]map[string]interface{}, 0)
					methodConfigMap := map[string]interface{}{}
					methodConfigMap["values"] = methodConfigArg["Values"].([]interface{})
					methodConfigMaps = append(methodConfigMaps, methodConfigMap)
					ruleConditionsMap["method_config"] = methodConfigMaps
				}
			}

			if pathConfig, ok := ruleConditionsArg["PathConfig"]; ok {
				pathConfigArg := pathConfig.(map[string]interface{})
				if len(pathConfigArg) > 0 {
					pathConfigMaps := make([]map[string]interface{}, 0)
					pathConfigMap := map[string]interface{}{}
					pathConfigMap["values"] = pathConfigArg["Values"].([]interface{})
					pathConfigMaps = append(pathConfigMaps, pathConfigMap)
					ruleConditionsMap["path_config"] = pathConfigMaps
				}
			}

			if sourceIpConfig, ok := ruleConditionsArg["SourceIpConfig"]; ok {
				sourceIpConfigArg := sourceIpConfig.(map[string]interface{})
				if len(sourceIpConfigArg) > 0 {
					sourceIpConfigMaps := make([]map[string]interface{}, 0)
					sourceIpConfigMap := map[string]interface{}{}
					sourceIpConfigMap["values"] = sourceIpConfigArg["Values"].([]interface{})
					sourceIpConfigMaps = append(sourceIpConfigMaps, sourceIpConfigMap)
					ruleConditionsMap["source_ip_config"] = sourceIpConfigMaps
				}
			}

			ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
		}

		d.Set("rule_conditions", ruleConditionsMaps)
	}

	return nil
}

func resourceAliCloudAlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"ClientToken": buildClientToken("UpdateRuleAttribute"),
		"RuleId":      d.Id(),
	}

	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")

	if d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")

	var direction string
	if v, ok := d.GetOk("direction"); ok {
		direction = v.(string)
	}

	if d.HasChange("rule_actions") {
		update = true
	}
	ruleActionsMaps := make([]map[string]interface{}, 0)
	for _, ruleActions := range d.Get("rule_actions").(*schema.Set).List() {
		ruleActionsArg := ruleActions.(map[string]interface{})
		ruleActionsMap := map[string]interface{}{}
		ruleActionsMap["Order"] = ruleActionsArg["order"]
		ruleActionsMap["Type"] = ruleActionsArg["type"]

		if ruleActionsMap["Type"] == "" {
			continue
		}

		fixedResponseConfigMap := map[string]interface{}{}
		for _, fixedResponseConfig := range ruleActionsArg["fixed_response_config"].(*schema.Set).List() {
			fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
			fixedResponseConfigMap["Content"] = fixedResponseConfigArg["content"]
			fixedResponseConfigMap["ContentType"] = fixedResponseConfigArg["content_type"]
			fixedResponseConfigMap["HttpCode"] = fixedResponseConfigArg["http_code"]
			ruleActionsMap["FixedResponseConfig"] = fixedResponseConfigMap
		}

		forwardGroupConfigMap := map[string]interface{}{}
		for _, forwardGroupConfig := range ruleActionsArg["forward_group_config"].([]interface{}) {
			forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
			serverGroupTuplesMaps := make([]map[string]interface{}, 0)
			for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
				serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
				serverGroupTuplesMap := map[string]interface{}{}
				serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
				if v, ok := serverGroupTuplesArg["weight"]; ok {
					serverGroupTuplesMap["Weight"] = v
				}

				serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
			}
			forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps

			serverGroupStickySessionMap := map[string]interface{}{}
			for _, serverGroupStickySession := range forwardGroupConfigArg["server_group_sticky_session"].(*schema.Set).List() {
				serverGroupStickySessionArg := serverGroupStickySession.(map[string]interface{})
				serverGroupStickySessionMap["Enabled"] = serverGroupStickySessionArg["enabled"]
				serverGroupStickySessionMap["Timeout"] = serverGroupStickySessionArg["timeout"]
			}
			forwardGroupConfigMap["ServerGroupStickySession"] = serverGroupStickySessionMap

			ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		insertHeaderConfigMap := map[string]interface{}{}
		for _, insertHeaderConfig := range ruleActionsArg["insert_header_config"].(*schema.Set).List() {
			insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
			insertHeaderConfigMap["Key"] = insertHeaderConfigArg["key"]
			insertHeaderConfigMap["Value"] = insertHeaderConfigArg["value"]
			insertHeaderConfigMap["ValueType"] = insertHeaderConfigArg["value_type"]
			ruleActionsMap["InsertHeaderConfig"] = insertHeaderConfigMap
		}

		removeHeaderConfigMap := map[string]interface{}{}
		for _, removeHeaderConfig := range ruleActionsArg["remove_header_config"].(*schema.Set).List() {
			removeHeaderConfigArg := removeHeaderConfig.(map[string]interface{})
			removeHeaderConfigMap["Key"] = removeHeaderConfigArg["key"]
			ruleActionsMap["RemoveHeaderConfig"] = removeHeaderConfigMap
		}

		redirectConfigMap := map[string]interface{}{}
		for _, redirectConfig := range ruleActionsArg["redirect_config"].(*schema.Set).List() {
			redirectConfigArg := redirectConfig.(map[string]interface{})
			redirectConfigMap["Host"] = redirectConfigArg["host"]
			redirectConfigMap["HttpCode"] = redirectConfigArg["http_code"]
			redirectConfigMap["Path"] = redirectConfigArg["path"]
			redirectConfigMap["Port"] = redirectConfigArg["port"]
			redirectConfigMap["Protocol"] = redirectConfigArg["protocol"]
			redirectConfigMap["Query"] = redirectConfigArg["query"]
			ruleActionsMap["RedirectConfig"] = redirectConfigMap
		}

		rewriteConfigMap := map[string]interface{}{}
		for _, rewriteConfig := range ruleActionsArg["rewrite_config"].(*schema.Set).List() {
			rewriteConfigArg := rewriteConfig.(map[string]interface{})
			rewriteConfigMap["Host"] = rewriteConfigArg["host"]
			rewriteConfigMap["Path"] = rewriteConfigArg["path"]
			rewriteConfigMap["Query"] = rewriteConfigArg["query"]
			ruleActionsMap["RewriteConfig"] = rewriteConfigMap
		}

		trafficLimitConfigList := ruleActionsArg["traffic_limit_config"].(*schema.Set).List()
		if len(trafficLimitConfigList) > 0 {
			trafficLimitConfigArg := trafficLimitConfigList[0].(map[string]interface{})
			ruleActionsMap["TrafficLimitConfig"] = map[string]interface{}{
				"QPS":      trafficLimitConfigArg["qps"],
				"PerIpQps": trafficLimitConfigArg["per_ip_qps"],
			}
		}

		trafficMirrorConfigList := ruleActionsArg["traffic_mirror_config"].(*schema.Set).List()
		if len(trafficMirrorConfigList) > 0 {
			trafficMirrorConfigArg := trafficMirrorConfigList[0].(map[string]interface{})
			mirrorGroupConfigMap := make(map[string]interface{}, 0)
			mirrorGroupConfigList := trafficMirrorConfigArg["mirror_group_config"].(*schema.Set).List()
			if len(mirrorGroupConfigList) > 0 {
				mirrorGroupConfigArg := mirrorGroupConfigList[0].(map[string]interface{})
				serverGroupTuplesMaps := make([]map[string]interface{}, 0)
				for _, serverGroupTuples := range mirrorGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
					serverGroupTuplesMap := map[string]interface{}{}
					serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
					serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
				mirrorGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
			}
			ruleActionsMap["TrafficMirrorConfig"] = map[string]interface{}{
				"MirrorGroupConfig": mirrorGroupConfigMap,
				"TargetType":        trafficMirrorConfigArg["target_type"],
			}
		}

		if v, ok := ruleActionsArg["cors_config"]; ok {
			corsConfigMap := map[string]interface{}{}
			for _, corsConfigList := range v.(*schema.Set).List() {
				corsConfigArg := corsConfigList.(map[string]interface{})

				if allowOrigin, ok := corsConfigArg["allow_origin"]; ok {
					corsConfigMap["AllowOrigin"] = allowOrigin
				}

				if allowMethods, ok := corsConfigArg["allow_methods"]; ok {
					corsConfigMap["AllowMethods"] = allowMethods
				}

				if allowHeaders, ok := corsConfigArg["allow_headers"]; ok {
					corsConfigMap["AllowHeaders"] = allowHeaders
				}

				if exposeHeaders, ok := corsConfigArg["expose_headers"]; ok {
					corsConfigMap["ExposeHeaders"] = exposeHeaders
				}

				if allowCredentials, ok := corsConfigArg["allow_credentials"]; ok {
					corsConfigMap["AllowCredentials"] = allowCredentials
				}

				if maxAge, ok := corsConfigArg["max_age"]; ok {
					corsConfigMap["MaxAge"] = maxAge
				}
			}

			ruleActionsMap["CorsConfig"] = corsConfigMap
		}

		ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
	}
	request["RuleActions"] = ruleActionsMaps

	if d.HasChange("rule_conditions") {
		update = true
	}
	ruleConditionsMaps := make([]map[string]interface{}, 0)
	for _, ruleConditions := range d.Get("rule_conditions").(*schema.Set).List() {
		ruleConditionsArg := ruleConditions.(map[string]interface{})
		ruleConditionsMap := map[string]interface{}{}
		ruleConditionsMap["Type"] = ruleConditionsArg["type"]

		cookieConfigMap := map[string]interface{}{}
		for _, cookieConfig := range ruleConditionsArg["cookie_config"].(*schema.Set).List() {
			cookieConfigArg := cookieConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range cookieConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			cookieConfigMap["Values"] = valuesMaps
			ruleConditionsMap["CookieConfig"] = cookieConfigMap
		}

		headerConfigMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["header_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			headerConfigMap["Key"] = headerConfigArg["key"]
			headerConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HeaderConfig"] = headerConfigMap
		}

		responseHeaderConfigMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["response_header_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			responseHeaderConfigMap["Key"] = headerConfigArg["key"]
			responseHeaderConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["ResponseHeaderConfig"] = responseHeaderConfigMap
		}

		responseStatusCodeMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["response_status_code_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			responseStatusCodeMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["ResponseStatusCodeConfig"] = responseStatusCodeMap
		}

		hostConfigMap := map[string]interface{}{}
		for _, hostConfig := range ruleConditionsArg["host_config"].(*schema.Set).List() {
			hostConfigArg := hostConfig.(map[string]interface{})
			hostConfigMap["Values"] = hostConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HostConfig"] = hostConfigMap
		}

		methodConfigMap := map[string]interface{}{}
		for _, methodConfig := range ruleConditionsArg["method_config"].(*schema.Set).List() {
			methodConfigArg := methodConfig.(map[string]interface{})
			methodConfigMap["Values"] = methodConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["MethodConfig"] = methodConfigMap
		}

		pathConfigMap := map[string]interface{}{}
		for _, pathConfig := range ruleConditionsArg["path_config"].(*schema.Set).List() {
			pathConfigArg := pathConfig.(map[string]interface{})
			pathConfigMap["Values"] = pathConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["PathConfig"] = pathConfigMap
		}

		sourceIpConfigMap := map[string]interface{}{}
		for _, sourceIpConfig := range ruleConditionsArg["source_ip_config"].(*schema.Set).List() {
			sourceIpConfigArg := sourceIpConfig.(map[string]interface{})
			sourceIpConfigMap["Values"] = sourceIpConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["SourceIpConfig"] = sourceIpConfigMap
		}

		queryStringConfigMap := map[string]interface{}{}
		for _, queryStringConfig := range ruleConditionsArg["query_string_config"].(*schema.Set).List() {
			queryStringConfigArg := queryStringConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range queryStringConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			queryStringConfigMap["Values"] = valuesMaps
			ruleConditionsMap["QueryStringConfig"] = queryStringConfigMap
		}

		ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
	}
	request["RuleConditions"] = ruleConditionsMaps

	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}

		action := "UpdateRuleAttribute"

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ResourceInConfiguring.Listener"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbRuleStateRefreshFunc(d.Id(), direction, []string{"CreateFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudAlbRuleRead(d, meta)
}

func resourceAliCloudAlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRule"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"ClientToken": buildClientToken("DeleteRule"),
		"RuleId":      d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.Rule", "SystemBusy", "-21013", "ResourceInConfiguring.Listener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Rule"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
