// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrChain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrChainCreate,
		Read:   resourceAliCloudCrChainRead,
		Update: resourceAliCloudCrChainUpdate,
		Delete: resourceAliCloudCrChainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"chain_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chain_config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"to": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"is_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"node_config": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"deny_policy": {
													Type:     schema.TypeList,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"issue_count": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"logic": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"issue_level": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"retry": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"scan_engine": {
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
			"chain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"chain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"is_success": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	}
}

func resourceAliCloudCrChainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateChain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	chainConfig := make(map[string]interface{})

	if v := d.Get("chain_config"); !IsNil(v) {
		localData, err := jsonpath.Get("$[0].nodes", v)
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
			dataLoopMap["NodeName"] = dataLoopTmp["node_name"]
			dataLoopMap["Enable"] = dataLoopTmp["enable"]
			localData1 := make(map[string]interface{})
			scanEngine1, _ := jsonpath.Get("$[0].scan_engine", dataLoopTmp["node_config"])
			if scanEngine1 != nil && scanEngine1 != "" {
				localData1["ScanEngine"] = scanEngine1
			}
			if len(localData1) > 0 {
				dataLoopMap["NodeConfig"] = localData1
			}
			localMaps = append(localMaps, dataLoopMap)
		}
		chainConfig["Nodes"] = localMaps

		localData2, err := jsonpath.Get("$[0].routers", v)
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
			localData3 := make(map[string]interface{})
			nodeName3, _ := jsonpath.Get("$[0].node_name", dataLoop2Tmp["to"])
			if nodeName3 != nil && nodeName3 != "" {
				localData3["NodeName"] = nodeName3
			}
			if len(localData3) > 0 {
				dataLoop2Map["To"] = localData3
			}
			localData4 := make(map[string]interface{})
			nodeName5, _ := jsonpath.Get("$[0].node_name", dataLoop2Tmp["from"])
			if nodeName5 != nil && nodeName5 != "" {
				localData4["NodeName"] = nodeName5
			}
			if len(localData4) > 0 {
				dataLoop2Map["From"] = localData4
			}
			localMaps2 = append(localMaps2, dataLoop2Map)
		}
		chainConfig["Routers"] = localMaps2

		request["ChainConfig"] = chainConfig
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["RepoNamespaceName"] = d.Get("repo_namespace_name")
	request["RepoName"] = d.Get("repo_name")
	request["Name"] = d.Get("chain_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["ChainId"]))

	return resourceAliCloudCrChainRead(d, meta)
}

func resourceAliCloudCrChainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrChain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_chain DescribeCrChain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("chain_name", objectRaw["Name"])
	d.Set("code", objectRaw["Code"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("is_success", objectRaw["IsSuccess"])
	d.Set("modified_time", objectRaw["ModifiedTime"])
	d.Set("scope_id", objectRaw["ScopeId"])
	d.Set("scope_type", objectRaw["ScopeType"])
	d.Set("chain_id", objectRaw["ChainId"])
	d.Set("instance_id", objectRaw["InstanceId"])

	chainConfigMaps := make([]map[string]interface{}, 0)
	chainConfigMap := make(map[string]interface{})
	chainConfigRaw := make(map[string]interface{})
	if objectRaw["ChainConfig"] != nil {
		chainConfigRaw = objectRaw["ChainConfig"].(map[string]interface{})
	}
	if len(chainConfigRaw) > 0 {
		chainConfigMap["chain_config_id"] = chainConfigRaw["ChainConfigId"]
		chainConfigMap["is_active"] = chainConfigRaw["IsActive"]
		chainConfigMap["version"] = chainConfigRaw["Version"]

		nodesRaw := chainConfigRaw["Nodes"]
		nodesMaps := make([]map[string]interface{}, 0)
		if nodesRaw != nil {
			for _, nodesChildRaw := range convertToInterfaceArray(nodesRaw) {
				nodesMap := make(map[string]interface{})
				nodesChildRaw := nodesChildRaw.(map[string]interface{})
				nodesMap["enable"] = nodesChildRaw["Enable"]
				nodesMap["node_name"] = nodesChildRaw["NodeName"]

				nodeConfigMaps := make([]map[string]interface{}, 0)
				nodeConfigMap := make(map[string]interface{})
				nodeConfigRaw := make(map[string]interface{})
				if nodesChildRaw["NodeConfig"] != nil {
					nodeConfigRaw = nodesChildRaw["NodeConfig"].(map[string]interface{})
				}
				if len(nodeConfigRaw) > 0 {
					nodeConfigMap["retry"] = nodeConfigRaw["Retry"]
					nodeConfigMap["scan_engine"] = nodeConfigRaw["ScanEngine"]
					nodeConfigMap["timeout"] = nodeConfigRaw["Timeout"]

					denyPolicyMaps := make([]map[string]interface{}, 0)
					denyPolicyMap := make(map[string]interface{})
					denyPolicyRaw := make(map[string]interface{})
					if nodeConfigRaw["DenyPolicy"] != nil {
						denyPolicyRaw = nodeConfigRaw["DenyPolicy"].(map[string]interface{})
					}
					if len(denyPolicyRaw) > 0 {
						denyPolicyMap["action"] = denyPolicyRaw["Action"]
						denyPolicyMap["issue_count"] = denyPolicyRaw["IssueCount"]
						denyPolicyMap["issue_level"] = denyPolicyRaw["IssueLevel"]
						denyPolicyMap["logic"] = denyPolicyRaw["Logic"]

						denyPolicyMaps = append(denyPolicyMaps, denyPolicyMap)
					}
					nodeConfigMap["deny_policy"] = denyPolicyMaps
					nodeConfigMaps = append(nodeConfigMaps, nodeConfigMap)
				}
				nodesMap["node_config"] = nodeConfigMaps
				nodesMaps = append(nodesMaps, nodesMap)
			}
		}
		chainConfigMap["nodes"] = nodesMaps
		routersRaw := chainConfigRaw["Routers"]
		routersMaps := make([]map[string]interface{}, 0)
		if routersRaw != nil {
			for _, routersChildRaw := range convertToInterfaceArray(routersRaw) {
				routersMap := make(map[string]interface{})
				routersChildRaw := routersChildRaw.(map[string]interface{})

				fromMaps := make([]map[string]interface{}, 0)
				fromMap := make(map[string]interface{})
				fromRaw := make(map[string]interface{})
				if routersChildRaw["From"] != nil {
					fromRaw = routersChildRaw["From"].(map[string]interface{})
				}
				if len(fromRaw) > 0 {
					fromMap["node_name"] = fromRaw["NodeName"]

					fromMaps = append(fromMaps, fromMap)
				}
				routersMap["from"] = fromMaps
				toMaps := make([]map[string]interface{}, 0)
				toMap := make(map[string]interface{})
				toRaw := make(map[string]interface{})
				if routersChildRaw["To"] != nil {
					toRaw = routersChildRaw["To"].(map[string]interface{})
				}
				if len(toRaw) > 0 {
					toMap["node_name"] = toRaw["NodeName"]

					toMaps = append(toMaps, toMap)
				}
				routersMap["to"] = toMaps
				routersMaps = append(routersMaps, routersMap)
			}
		}
		chainConfigMap["routers"] = routersMaps
		chainConfigMaps = append(chainConfigMaps, chainConfigMap)
	}
	if err := d.Set("chain_config", chainConfigMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCrChainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateChain"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ChainId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("chain_config") {
		update = true
	}
	chainConfig := make(map[string]interface{})

	if v := d.Get("chain_config"); v != nil {
		localData, err := jsonpath.Get("$[0].nodes", v)
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
			dataLoopMap["NodeName"] = dataLoopTmp["node_name"]
			dataLoopMap["Enable"] = dataLoopTmp["enable"]
			if !IsNil(dataLoopTmp["node_config"]) {
				localData1 := make(map[string]interface{})
				scanEngine1, _ := jsonpath.Get("$[0].scan_engine", dataLoopTmp["node_config"])
				if scanEngine1 != nil && scanEngine1 != "" {
					localData1["ScanEngine"] = scanEngine1
				}
				if len(localData1) > 0 {
					dataLoopMap["NodeConfig"] = localData1
				}
			}
			localMaps = append(localMaps, dataLoopMap)
		}
		chainConfig["Nodes"] = localMaps

		localData2, err := jsonpath.Get("$[0].routers", v)
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
			if !IsNil(dataLoop2Tmp["to"]) {
				localData3 := make(map[string]interface{})
				nodeName3, _ := jsonpath.Get("$[0].node_name", dataLoop2Tmp["to"])
				if nodeName3 != nil && nodeName3 != "" {
					localData3["NodeName"] = nodeName3
				}
				if len(localData3) > 0 {
					dataLoop2Map["To"] = localData3
				}
			}
			if !IsNil(dataLoop2Tmp["from"]) {
				localData4 := make(map[string]interface{})
				nodeName5, _ := jsonpath.Get("$[0].node_name", dataLoop2Tmp["from"])
				if nodeName5 != nil && nodeName5 != "" {
					localData4["NodeName"] = nodeName5
				}
				if len(localData4) > 0 {
					dataLoop2Map["From"] = localData4
				}
			}
			localMaps2 = append(localMaps2, dataLoop2Map)
		}
		chainConfig["Routers"] = localMaps2

		request["ChainConfig"] = chainConfig
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("chain_name") {
		update = true
	}
	request["Name"] = d.Get("chain_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	return resourceAliCloudCrChainRead(d, meta)
}

func resourceAliCloudCrChainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteChain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ChainId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
