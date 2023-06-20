// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeConnectionCreate,
		Read:   resourceAliCloudEventBridgeConnectionRead,
		Update: resourceAliCloudEventBridgeConnectionUpdate,
		Delete: resourceAliCloudEventBridgeConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_parameters": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_secret": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"client_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"oauth_http_parameters": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"header_parameters": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_value_secret": {
																Type:     schema.TypeString,
																Optional: true,
															},
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
												"query_string_parameters": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_value_secret": {
																Type:     schema.TypeString,
																Optional: true,
															},
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
												"body_parameters": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_value_secret": {
																Type:     schema.TypeString,
																Optional: true,
															},
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
											},
										},
									},
									"authorization_endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"http_method": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"basic_auth_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
						"api_key_auth_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_key_value": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
									"api_key_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"authorization_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_parameters": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitche_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudEventBridgeConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEventBridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ConnectionName"] = d.Get("connection_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("network_parameters"); ok {
		nodeNative, _ := jsonpath.Get("$[0].network_type", v)
		if nodeNative != "" {
			objectDataLocalMap["NetworkType"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].vpc_id", v)
		if nodeNative1 != "" {
			objectDataLocalMap["VpcId"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].vswitche_id", v)
		if nodeNative2 != "" {
			objectDataLocalMap["VswitcheId"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].security_group_id", v)
		if nodeNative3 != "" {
			objectDataLocalMap["SecurityGroupId"] = nodeNative3
		}
	}
	request["NetworkParameters"] = objectDataLocalMap

	objectDataLocalMap1 := make(map[string]interface{})
	if v, ok := d.GetOk("auth_parameters"); ok {
		nodeNative4, _ := jsonpath.Get("$[0].authorization_type", v)
		if nodeNative4 != "" {
			objectDataLocalMap1["AuthorizationType"] = nodeNative4
		}
		apiKeyAuthParameters := make(map[string]interface{})
		nodeNative5, _ := jsonpath.Get("$[0].api_key_auth_parameters[0].api_key_name", v)
		if nodeNative5 != "" {
			apiKeyAuthParameters["ApiKeyName"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].api_key_auth_parameters[0].api_key_value", v)
		if nodeNative6 != "" {
			apiKeyAuthParameters["ApiKeyValue"] = nodeNative6
		}
		objectDataLocalMap1["ApiKeyAuthParameters"] = apiKeyAuthParameters
		basicAuthParameters := make(map[string]interface{})
		nodeNative7, _ := jsonpath.Get("$[0].basic_auth_parameters[0].password", v)
		if nodeNative7 != "" {
			basicAuthParameters["Password"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].basic_auth_parameters[0].username", v)
		if nodeNative8 != "" {
			basicAuthParameters["Username"] = nodeNative8
		}
		objectDataLocalMap1["BasicAuthParameters"] = basicAuthParameters
		oAuthParameters := make(map[string]interface{})
		clientParameters := make(map[string]interface{})
		nodeNative9, _ := jsonpath.Get("$[0].oauth_parameters[0].client_parameters[0].client_id", v)
		if nodeNative9 != "" {
			clientParameters["ClientID"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].oauth_parameters[0].client_parameters[0].client_secret", v)
		if nodeNative10 != "" {
			clientParameters["ClientSecret"] = nodeNative10
		}
		oAuthParameters["ClientParameters"] = clientParameters
		nodeNative11, _ := jsonpath.Get("$[0].oauth_parameters[0].http_method", v)
		if nodeNative11 != "" {
			oAuthParameters["HttpMethod"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].oauth_parameters[0].authorization_endpoint", v)
		if nodeNative12 != "" {
			oAuthParameters["AuthorizationEndpoint"] = nodeNative12
		}
		oAuthHttpParameters := make(map[string]interface{})
		if v, ok := d.GetOk("auth_parameters"); ok {
			localData, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].body_parameters", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["IsValueSecret"] = dataLoopTmp["is_value_secret"]
				dataLoopMap["Key"] = dataLoopTmp["key"]
				dataLoopMap["Value"] = dataLoopTmp["value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			oAuthHttpParameters["BodyParameters"] = localMaps
		}
		if v, ok := d.GetOk("auth_parameters"); ok {
			localData1, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].header_parameters", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps1 := make([]map[string]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["IsValueSecret"] = dataLoop1Tmp["is_value_secret"]
				dataLoop1Map["Key"] = dataLoop1Tmp["key"]
				dataLoop1Map["Value"] = dataLoop1Tmp["value"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			oAuthHttpParameters["HeaderParameters"] = localMaps1
		}
		if v, ok := d.GetOk("auth_parameters"); ok {
			localData2, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].query_string_parameters", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps2 := make([]map[string]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["IsValueSecret"] = dataLoop2Tmp["is_value_secret"]
				dataLoop2Map["Key"] = dataLoop2Tmp["key"]
				dataLoop2Map["Value"] = dataLoop2Tmp["value"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			oAuthHttpParameters["QueryStringParameters"] = localMaps2
		}
		oAuthParameters["OAuthHttpParameters"] = oAuthHttpParameters
		objectDataLocalMap1["OAuthParameters"] = oAuthParameters
	}
	request["AuthParameters"] = objectDataLocalMap1

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_connection", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.Data.ConnectionName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudEventBridgeConnectionUpdate(d, meta)
}

func resourceAliCloudEventBridgeConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeServiceV2 := EventBridgeServiceV2{client}

	objectRaw, err := eventBridgeServiceV2.DescribeEventBridgeConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_connection DescribeEventBridgeConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	authParametersMaps := make([]map[string]interface{}, 0)
	authParametersMap := make(map[string]interface{})
	authParameters1Raw := make(map[string]interface{})
	if objectRaw["AuthParameters"] != nil {
		authParameters1Raw = objectRaw["AuthParameters"].(map[string]interface{})
	}
	if len(authParameters1Raw) > 0 {
		authParametersMap["authorization_type"] = authParameters1Raw["AuthorizationType"]
		apiKeyAuthParametersMaps := make([]map[string]interface{}, 0)
		apiKeyAuthParametersMap := make(map[string]interface{})
		apiKeyAuthParameters1Raw := make(map[string]interface{})
		if authParameters1Raw["ApiKeyAuthParameters"] != nil {
			apiKeyAuthParameters1Raw = authParameters1Raw["ApiKeyAuthParameters"].(map[string]interface{})
		}
		if len(apiKeyAuthParameters1Raw) > 0 {
			apiKeyAuthParametersMap["api_key_name"] = apiKeyAuthParameters1Raw["ApiKeyName"]
			apiKeyAuthParametersMap["api_key_value"] = apiKeyAuthParameters1Raw["ApiKeyValue"]
			apiKeyAuthParametersMaps = append(apiKeyAuthParametersMaps, apiKeyAuthParametersMap)
		}
		authParametersMap["api_key_auth_parameters"] = apiKeyAuthParametersMaps
		basicAuthParametersMaps := make([]map[string]interface{}, 0)
		basicAuthParametersMap := make(map[string]interface{})
		basicAuthParameters1Raw := make(map[string]interface{})
		if authParameters1Raw["BasicAuthParameters"] != nil {
			basicAuthParameters1Raw = authParameters1Raw["BasicAuthParameters"].(map[string]interface{})
		}
		if len(basicAuthParameters1Raw) > 0 {
			basicAuthParametersMap["password"] = basicAuthParameters1Raw["Password"]
			basicAuthParametersMap["username"] = basicAuthParameters1Raw["Username"]
			basicAuthParametersMaps = append(basicAuthParametersMaps, basicAuthParametersMap)
		}
		authParametersMap["basic_auth_parameters"] = basicAuthParametersMaps
		oauthParametersMaps := make([]map[string]interface{}, 0)
		oauthParametersMap := make(map[string]interface{})
		oAuthParameters1Raw := make(map[string]interface{})
		if authParameters1Raw["OAuthParameters"] != nil {
			oAuthParameters1Raw = authParameters1Raw["OAuthParameters"].(map[string]interface{})
		}
		if len(oAuthParameters1Raw) > 0 {
			oauthParametersMap["authorization_endpoint"] = oAuthParameters1Raw["AuthorizationEndpoint"]
			oauthParametersMap["http_method"] = oAuthParameters1Raw["HttpMethod"]
			clientParametersMaps := make([]map[string]interface{}, 0)
			clientParametersMap := make(map[string]interface{})
			clientParameters1Raw := make(map[string]interface{})
			if oAuthParameters1Raw["ClientParameters"] != nil {
				clientParameters1Raw = oAuthParameters1Raw["ClientParameters"].(map[string]interface{})
			}
			if len(clientParameters1Raw) > 0 {
				clientParametersMap["client_id"] = clientParameters1Raw["ClientID"]
				clientParametersMap["client_secret"] = clientParameters1Raw["ClientSecret"]
				clientParametersMaps = append(clientParametersMaps, clientParametersMap)
			}
			oauthParametersMap["client_parameters"] = clientParametersMaps
			oauthHttpParametersMaps := make([]map[string]interface{}, 0)
			oauthHttpParametersMap := make(map[string]interface{})
			oAuthHttpParameters1Raw := make(map[string]interface{})
			if oAuthParameters1Raw["OAuthHttpParameters"] != nil {
				oAuthHttpParameters1Raw = oAuthParameters1Raw["OAuthHttpParameters"].(map[string]interface{})
			}
			if len(oAuthHttpParameters1Raw) > 0 {
				bodyParameters1Raw := oAuthHttpParameters1Raw["BodyParameters"]
				bodyParametersMaps := make([]map[string]interface{}, 0)
				if bodyParameters1Raw != nil {
					for _, bodyParametersChild1Raw := range bodyParameters1Raw.([]interface{}) {
						bodyParametersMap := make(map[string]interface{})
						bodyParametersChild1Raw := bodyParametersChild1Raw.(map[string]interface{})
						bodyParametersMap["is_value_secret"] = bodyParametersChild1Raw["IsValueSecret"]
						bodyParametersMap["key"] = bodyParametersChild1Raw["Key"]
						bodyParametersMap["value"] = bodyParametersChild1Raw["Value"]
						bodyParametersMaps = append(bodyParametersMaps, bodyParametersMap)
					}
				}
				oauthHttpParametersMap["body_parameters"] = bodyParametersMaps
				headerParameters1Raw := oAuthHttpParameters1Raw["HeaderParameters"]
				headerParametersMaps := make([]map[string]interface{}, 0)
				if headerParameters1Raw != nil {
					for _, headerParametersChild1Raw := range headerParameters1Raw.([]interface{}) {
						headerParametersMap := make(map[string]interface{})
						headerParametersChild1Raw := headerParametersChild1Raw.(map[string]interface{})
						headerParametersMap["is_value_secret"] = headerParametersChild1Raw["IsValueSecret"]
						headerParametersMap["key"] = headerParametersChild1Raw["Key"]
						headerParametersMap["value"] = headerParametersChild1Raw["Value"]
						headerParametersMaps = append(headerParametersMaps, headerParametersMap)
					}
				}
				oauthHttpParametersMap["header_parameters"] = headerParametersMaps
				queryStringParameters1Raw := oAuthHttpParameters1Raw["QueryStringParameters"]
				queryStringParametersMaps := make([]map[string]interface{}, 0)
				if queryStringParameters1Raw != nil {
					for _, queryStringParametersChild1Raw := range queryStringParameters1Raw.([]interface{}) {
						queryStringParametersMap := make(map[string]interface{})
						queryStringParametersChild1Raw := queryStringParametersChild1Raw.(map[string]interface{})
						queryStringParametersMap["is_value_secret"] = queryStringParametersChild1Raw["IsValueSecret"]
						queryStringParametersMap["key"] = queryStringParametersChild1Raw["Key"]
						queryStringParametersMap["value"] = queryStringParametersChild1Raw["Value"]
						queryStringParametersMaps = append(queryStringParametersMaps, queryStringParametersMap)
					}
				}
				oauthHttpParametersMap["query_string_parameters"] = queryStringParametersMaps
				oauthHttpParametersMaps = append(oauthHttpParametersMaps, oauthHttpParametersMap)
			}
			oauthParametersMap["oauth_http_parameters"] = oauthHttpParametersMaps
			oauthParametersMaps = append(oauthParametersMaps, oauthParametersMap)
		}
		authParametersMap["oauth_parameters"] = oauthParametersMaps
		authParametersMaps = append(authParametersMaps, authParametersMap)
	}
	d.Set("auth_parameters", authParametersMaps)
	networkParametersMaps := make([]map[string]interface{}, 0)
	networkParametersMap := make(map[string]interface{})
	networkParameters1Raw := make(map[string]interface{})
	if objectRaw["NetworkParameters"] != nil {
		networkParameters1Raw = objectRaw["NetworkParameters"].(map[string]interface{})
	}
	if len(networkParameters1Raw) > 0 {
		networkParametersMap["network_type"] = networkParameters1Raw["NetworkType"]
		networkParametersMap["security_group_id"] = networkParameters1Raw["SecurityGroupId"]
		networkParametersMap["vpc_id"] = networkParameters1Raw["VpcId"]
		networkParametersMap["vswitche_id"] = networkParameters1Raw["VswitcheId"]
		networkParametersMaps = append(networkParametersMaps, networkParametersMap)
	}
	d.Set("network_parameters", networkParametersMaps)

	return nil
}

func resourceAliCloudEventBridgeConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateConnection"
	conn, err := client.NewEventBridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ConnectionName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("network_parameters") {
		update = true
		objectDataLocalMap := make(map[string]interface{})
		if v, ok := d.GetOk("network_parameters"); ok {
			nodeNative, _ := jsonpath.Get("$[0].network_type", v)
			if nodeNative != "" {
				objectDataLocalMap["NetworkType"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].vpc_id", v)
			if nodeNative1 != "" {
				objectDataLocalMap["VpcId"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].vswitche_id", v)
			if nodeNative2 != "" {
				objectDataLocalMap["VswitcheId"] = nodeNative2
			}
			nodeNative3, _ := jsonpath.Get("$[0].security_group_id", v)
			if nodeNative3 != "" {
				objectDataLocalMap["SecurityGroupId"] = nodeNative3
			}
		}
		request["NetworkParameters"] = objectDataLocalMap
	}

	if d.HasChange("auth_parameters") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})
		if v, ok := d.GetOk("auth_parameters"); ok {
			nodeNative4, _ := jsonpath.Get("$[0].authorization_type", v)
			if nodeNative4 != "" {
				objectDataLocalMap1["AuthorizationType"] = nodeNative4
			}
			apiKeyAuthParameters := make(map[string]interface{})
			nodeNative5, _ := jsonpath.Get("$[0].api_key_auth_parameters[0].api_key_name", v)
			if nodeNative5 != "" {
				apiKeyAuthParameters["ApiKeyName"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].api_key_auth_parameters[0].api_key_value", v)
			if nodeNative6 != "" {
				apiKeyAuthParameters["ApiKeyValue"] = nodeNative6
			}
			objectDataLocalMap1["ApiKeyAuthParameters"] = apiKeyAuthParameters
			basicAuthParameters := make(map[string]interface{})
			nodeNative7, _ := jsonpath.Get("$[0].basic_auth_parameters[0].password", v)
			if nodeNative7 != "" {
				basicAuthParameters["Password"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].basic_auth_parameters[0].username", v)
			if nodeNative8 != "" {
				basicAuthParameters["Username"] = nodeNative8
			}
			objectDataLocalMap1["BasicAuthParameters"] = basicAuthParameters
			oAuthParameters := make(map[string]interface{})
			nodeNative9, _ := jsonpath.Get("$[0].oauth_parameters[0].authorization_endpoint", v)
			if nodeNative9 != "" {
				oAuthParameters["AuthorizationEndpoint"] = nodeNative9
			}
			clientParameters := make(map[string]interface{})
			nodeNative10, _ := jsonpath.Get("$[0].oauth_parameters[0].client_parameters[0].client_id", v)
			if nodeNative10 != "" {
				clientParameters["ClientID"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].oauth_parameters[0].client_parameters[0].client_secret", v)
			if nodeNative11 != "" {
				clientParameters["ClientSecret"] = nodeNative11
			}
			oAuthParameters["ClientParameters"] = clientParameters
			nodeNative12, _ := jsonpath.Get("$[0].oauth_parameters[0].http_method", v)
			if nodeNative12 != "" {
				oAuthParameters["HttpMethod"] = nodeNative12
			}
			oAuthHttpParameters := make(map[string]interface{})
			if v, ok := d.GetOk("auth_parameters"); ok {
				localData, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].body_parameters", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := dataLoop.(map[string]interface{})
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["IsValueSecret"] = dataLoopTmp["is_value_secret"]
					dataLoopMap["Key"] = dataLoopTmp["key"]
					dataLoopMap["Value"] = dataLoopTmp["value"]
					localMaps = append(localMaps, dataLoopMap)
				}
				oAuthHttpParameters["BodyParameters"] = localMaps
			}
			if v, ok := d.GetOk("auth_parameters"); ok {
				localData1, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].header_parameters", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps1 := make([]map[string]interface{}, 0)
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := dataLoop1.(map[string]interface{})
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["IsValueSecret"] = dataLoop1Tmp["is_value_secret"]
					dataLoop1Map["Key"] = dataLoop1Tmp["key"]
					dataLoop1Map["Value"] = dataLoop1Tmp["value"]
					localMaps1 = append(localMaps1, dataLoop1Map)
				}
				oAuthHttpParameters["HeaderParameters"] = localMaps1
			}
			if v, ok := d.GetOk("auth_parameters"); ok {
				localData2, err := jsonpath.Get("$[0].oauth_parameters[0].oauth_http_parameters[0].query_string_parameters", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps2 := make([]map[string]interface{}, 0)
				for _, dataLoop2 := range localData2.([]interface{}) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["IsValueSecret"] = dataLoop2Tmp["is_value_secret"]
					dataLoop2Map["Key"] = dataLoop2Tmp["key"]
					dataLoop2Map["Value"] = dataLoop2Tmp["value"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				oAuthHttpParameters["QueryStringParameters"] = localMaps2
			}
			oAuthParameters["OAuthHttpParameters"] = oAuthHttpParameters
			objectDataLocalMap1["OAuthParameters"] = oAuthParameters
		}
		request["AuthParameters"] = objectDataLocalMap1
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
	return resourceAliCloudEventBridgeConnectionRead(d, meta)
}

func resourceAliCloudEventBridgeConnectionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEventBridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ConnectionName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
