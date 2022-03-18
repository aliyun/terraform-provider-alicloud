package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCrChain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrChainCreate,
		Read:   resourceAlicloudCrChainRead,
		Update: resourceAlicloudCrChainUpdate,
		Delete: resourceAlicloudCrChainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"chain_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"node_config": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deny_policy": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"issue_count": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"issue_level": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"LOW", "MEDIUM", "HIGH", "UNKNOWN"}, false),
															},
															"logic": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"AND", "OR"}, false),
															},
															"action": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"BLOCK", "BLOCK_RETAG", "BLOCK_DELETE_TAG"}, false),
															},
														},
													},
												},
											},
										},
									},
									"node_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"DOCKER_IMAGE_BUILD", "DOCKER_IMAGE_PUSH", "VULNERABILITY_SCANNING", "ACTIVATE_REPLICATION", "TRIGGER", "SNAPSHOT", "TRIGGER_SNAPSHOT"}, false),
									},
								},
							},
						},
						"routers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_name": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"DOCKER_IMAGE_BUILD", "DOCKER_IMAGE_PUSH", "VULNERABILITY_SCANNING", "ACTIVATE_REPLICATION", "TRIGGER", "SNAPSHOT", "TRIGGER_SNAPSHOT"}, false),
												},
											},
										},
									},
									"to": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_name": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"DOCKER_IMAGE_BUILD", "DOCKER_IMAGE_PUSH", "VULNERABILITY_SCANNING", "ACTIVATE_REPLICATION", "TRIGGER", "SNAPSHOT", "TRIGGER_SNAPSHOT"}, false),
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
			"chain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"chain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceAlicloudCrChainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateChain"
	request := make(map[string]interface{})
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("chain_config"); ok {
		request["ChainConfig"], _ = convertCrChainConfigToJsonString(v.(*schema.Set).List())
	}
	request["Name"] = d.Get("chain_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("repo_name"); ok {
		request["RepoName"] = v
	}
	if v, ok := d.GetOk("repo_namespace_name"); ok {
		request["RepoNamespaceName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_chain", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["ChainId"]))

	return resourceAlicloudCrChainRead(d, meta)
}
func resourceAlicloudCrChainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	object, err := crService.DescribeCrChain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_chain crService.DescribeCrChain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("chain_id", parts[1])
	d.Set("instance_id", parts[0])
	if chainConfig, ok := object["ChainConfig"].(map[string]interface{}); ok {

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
		if err := d.Set("chain_config", chainConfigParams); err != nil {
			return WrapError(err)
		}
	}
	d.Set("chain_name", object["Name"])
	d.Set("description", object["Description"])
	return nil
}
func resourceAlicloudCrChainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ChainId":    parts[1],
		"InstanceId": parts[0],
	}
	if d.HasChange("chain_config") {
		update = true
	}
	if v, ok := d.GetOk("chain_config"); ok {
		request["ChainConfig"], _ = convertCrChainConfigToJsonString(v.(*schema.Set).List())
	}
	if d.HasChange("chain_name") {
		update = true
	}
	request["Name"] = d.Get("chain_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if update {
		action := "UpdateChain"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Code"]) != "success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
	}
	return resourceAlicloudCrChainRead(d, meta)
}
func resourceAlicloudCrChainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteChain"
	var response map[string]interface{}
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ChainId":    parts[1],
		"InstanceId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

func convertCrChainConfigToJsonString(configs []interface{}) (string, error) {
	chainConfig := make(map[string]interface{})

	for _, config := range configs {
		if v, ok := config.(map[string]interface{}); ok {

			if v["routers"] != nil {
				routersConfigs := make([]map[string]interface{}, 0)

				for _, router := range v["routers"].(*schema.Set).List() {
					routersConfig := make(map[string]interface{})
					if routerParameters, ok := router.(map[string]interface{}); ok {
						if _, ok := routerParameters["from"].(*schema.Set); ok {
							for _, fromNode := range routerParameters["from"].(*schema.Set).List() {
								routersConfigFromNode := make(map[string]interface{})
								routersConfigFromNode["NodeName"] = fromNode.(map[string]interface{})["node_name"]
								routersConfig["From"] = routersConfigFromNode
							}
						}
						if _, ok := routerParameters["to"].(*schema.Set); ok {
							for _, toNode := range routerParameters["to"].(*schema.Set).List() {
								routersConfigToNode := make(map[string]interface{})
								routersConfigToNode["NodeName"] = toNode.(map[string]interface{})["node_name"]
								routersConfig["To"] = routersConfigToNode
							}
						}
						routersConfigs = append(routersConfigs, routersConfig)
					}
				}
				chainConfig["Routers"] = routersConfigs
			}
			if v["nodes"] != nil {
				nodesConfigs := make([]map[string]interface{}, 0)

				for _, node := range v["nodes"].(*schema.Set).List() {

					if nodeParameters, ok := node.(map[string]interface{}); ok {
						nodesConfig := make(map[string]interface{})
						nodesConfig["Enable"] = nodeParameters["enable"]
						nodesConfig["NodeName"] = nodeParameters["node_name"]

						nodeConfigParameter := make(map[string]interface{})
						if _, ok := nodeParameters["node_config"].(*schema.Set); ok {
							for _, nodeConfig := range nodeParameters["node_config"].(*schema.Set).List() {
								nodeConfigArg := nodeConfig.(map[string]interface{})
								denyPolicyParameter := make(map[string]interface{})
								if _, ok := nodeConfigArg["deny_policy"].(*schema.Set); ok {
									for _, denyPolicy := range nodeConfigArg["deny_policy"].(*schema.Set).List() {
										if vv, ok := denyPolicy.(map[string]interface{})["issue_count"]; ok {
											denyPolicyParameter["IssueCount"] = vv
										}
										if vv, ok := denyPolicy.(map[string]interface{})["issue_level"]; ok {
											denyPolicyParameter["IssueLevel"] = vv
										}
										if vv, ok := denyPolicy.(map[string]interface{})["logic"]; ok {
											denyPolicyParameter["Logic"] = vv
										}
										if vv, ok := denyPolicy.(map[string]interface{})["action"]; ok {
											denyPolicyParameter["Action"] = vv
										}
									}
									nodeConfigParameter["DenyPolicy"] = denyPolicyParameter
								}
							}
							nodesConfig["node_config"] = nodeConfigParameter
						}
						nodesConfigs = append(nodesConfigs, nodesConfig)
					}
				}
				chainConfig["Nodes"] = nodesConfigs
			}
		}
	}
	if v, err := convertArrayObjectToJsonString(chainConfig); err != nil {
		return "", WrapError(err)
	} else {
		return v, nil
	}
}
