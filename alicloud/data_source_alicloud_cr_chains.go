package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCrChains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrChainsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"repo_namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"chains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chain_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nodes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"node_config": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"deny_policy": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"issue_count": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"issue_level": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"logic": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"action": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"node_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"routers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"from": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"node_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"to": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"node_name": {
																Type:     schema.TypeString,
																Computed: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudCrChainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListChain"
	request := make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("repo_name"); ok {
		request["RepoName"] = v
	}
	if v, ok := d.GetOk("repo_namespace_name"); ok {
		request["RepoNamespaceName"] = v
	}
	var objects []map[string]interface{}
	var chainNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		chainNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}

	pageNo, pageSize := 1, PageSizeLarge

	for {

		request["PageNo"] = pageNo
		request["PageSize"] = pageSize
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cr_chains", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Chains", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Chains", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if chainNameRegex != nil && !chainNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ChainId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < pageSize {
			break
		}
		pageNo++
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := make(map[string]interface{})
		mapping["id"] = fmt.Sprint(object["InstanceId"], ":", object["ChainId"])
		mapping["chain_id"] = object["ChainId"]
		mapping["chain_name"] = object["Name"]
		mapping["create_time"] = object["CreateTime"]
		mapping["description"] = object["Description"]
		mapping["instance_id"] = object["InstanceId"]
		mapping["modified_time"] = object["ModifiedTime"]
		mapping["scope_id"] = object["ScopeId"]
		mapping["scope_type"] = object["ScopeType"]

		id := mapping["id"].(string)
		ids = append(ids, id)
		names = append(names, fmt.Sprint(mapping["chain_name"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		crService := CrService{client}
		getResp, err := crService.DescribeCrChain(id)
		if err != nil {
			return WrapError(err)
		}

		if chainConfig, ok := getResp["ChainConfig"].(map[string]interface{}); ok {

			chainConfigParams := make([]map[string]interface{}, 0)
			chainConfigParam := make(map[string]interface{})

			if routers, ok := chainConfig["Routers"].([]interface{}); ok {
				routersConfigs := make([]map[string]interface{}, 0)
				for _, router := range routers {
					if v, ok := router.(map[string]interface{}); ok {
						routerConfig := make(map[string]interface{})
						if fromNode, ok := v["From"].(map[string]interface{}); ok {
							fromNodeParams := make([]map[string]interface{}, 0)
							fromNodeParam := map[string]interface{}{
								"node_name": fromNode["NodeName"],
							}
							fromNodeParams = append(fromNodeParams, fromNodeParam)
							routerConfig["from"] = fromNodeParams
						}
						if toNode, ok := v["To"].(map[string]interface{}); ok {
							toNodeParams := make([]map[string]interface{}, 0)
							toNodeParam := map[string]interface{}{
								"node_name": toNode["NodeName"],
							}
							toNodeParams = append(toNodeParams, toNodeParam)
							routerConfig["to"] = toNodeParams
						}
						routersConfigs = append(routersConfigs, routerConfig)
					}
				}
				chainConfigParam["routers"] = routersConfigs
			}

			if nodes, ok := chainConfig["Nodes"].([]interface{}); ok {
				nodesParams := make([]map[string]interface{}, 0)
				for _, node := range nodes {
					nodeParam := make(map[string]interface{})
					if v, ok := node.(map[string]interface{}); ok {
						if enable, ok := v["Enable"]; ok {
							nodeParam["enable"] = enable
						}
						if nodeName, ok := v["NodeName"]; ok {
							nodeParam["node_name"] = nodeName
						}
						nodeConfigParams := make([]map[string]interface{}, 0)
						nodeConfigParam := make(map[string]interface{})
						if nodeConfig, ok := v["NodeConfig"].(map[string]interface{}); ok {
							denyPolicyParams := make([]map[string]interface{}, 0)
							denyPolicyParam := make(map[string]interface{})
							if denyPolicy, ok := nodeConfig["DenyPolicy"].(map[string]interface{}); ok {
								denyPolicyParam["issue_count"] = denyPolicy["IssueCount"]
								denyPolicyParam["issue_level"] = denyPolicy["IssueLevel"]
								denyPolicyParam["logic"] = denyPolicy["Logic"]
								denyPolicyParam["action"] = denyPolicy["Action"]
							}
							denyPolicyParams = append(denyPolicyParams, denyPolicyParam)
							nodeConfigParam["deny_policy"] = denyPolicyParams
						}
						nodeConfigParams = append(nodeConfigParams, nodeConfigParam)
						nodeParam["node_config"] = nodeConfigParams
					}
					nodesParams = append(nodesParams, nodeParam)
				}
				chainConfigParam["nodes"] = nodesParams
			}

			chainConfigParams = append(chainConfigParams, chainConfigParam)
			mapping["chain_config"] = chainConfigParams
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("chains", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
