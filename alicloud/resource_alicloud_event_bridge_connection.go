package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Required: true,
				ForceNew: true,
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
						"network_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"PublicNetwork", "PrivateNetwork"}, false),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitche_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"auth_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorization_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"API_KEY_AUTH", "BASIC_AUTH", "OAUTH_AUTH"}, false),
						},
						"api_key_auth_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_key_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"api_key_value": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if old != "" && new != "" && old != new {
												return true
											}
											return false
										},
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
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if old != "" && new != "" && old != new {
												return true
											}
											return false
										},
									},
								},
							},
						},
						"oauth_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization_endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"http_method": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"GET", "POST", "HEAD", "DELETE", "PUT", "PATCH"}, false),
									},
									"client_parameters": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"client_secret": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
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
															"key": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"is_value_secret": {
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
															"key": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"is_value_secret": {
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
															"key": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"is_value_secret": {
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
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEventBridgeConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateConnection"
	request := make(map[string]interface{})
	var err error

	request["ConnectionName"] = d.Get("connection_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	networkParameters := d.Get("network_parameters")
	networkParametersMap := map[string]interface{}{}
	for _, networkParametersList := range networkParameters.([]interface{}) {
		networkParametersArg := networkParametersList.(map[string]interface{})

		networkParametersMap["NetworkType"] = networkParametersArg["network_type"]

		if vpcId, ok := networkParametersArg["vpc_id"]; ok {
			networkParametersMap["VpcId"] = vpcId
		}

		if vswitcheId, ok := networkParametersArg["vswitche_id"]; ok {
			networkParametersMap["VswitcheId"] = vswitcheId
		}

		if securityGroupId, ok := networkParametersArg["security_group_id"]; ok {
			networkParametersMap["SecurityGroupId"] = securityGroupId
		}
	}

	networkParametersJson, err := convertMaptoJsonString(networkParametersMap)
	if err != nil {
		return WrapError(err)
	}

	request["NetworkParameters"] = networkParametersJson

	if v, ok := d.GetOk("auth_parameters"); ok {
		authParametersMap := map[string]interface{}{}
		for _, authParametersList := range v.([]interface{}) {
			authParametersArg := authParametersList.(map[string]interface{})

			if authorizationType, ok := authParametersArg["authorization_type"]; ok {
				authParametersMap["AuthorizationType"] = authorizationType
			}

			if apiKeyAuthParameters, ok := authParametersArg["api_key_auth_parameters"]; ok {
				apiKeyAuthParametersMap := map[string]interface{}{}
				for _, apiKeyAuthParametersList := range apiKeyAuthParameters.([]interface{}) {
					apiKeyAuthParametersArg := apiKeyAuthParametersList.(map[string]interface{})

					if apiKeyName, ok := apiKeyAuthParametersArg["api_key_name"]; ok {
						apiKeyAuthParametersMap["ApiKeyName"] = apiKeyName
					}

					if apiKeyValue, ok := apiKeyAuthParametersArg["api_key_value"]; ok {
						apiKeyAuthParametersMap["ApiKeyValue"] = apiKeyValue
					}
				}

				authParametersMap["ApiKeyAuthParameters"] = apiKeyAuthParametersMap
			}

			if basicAuthParameters, ok := authParametersArg["basic_auth_parameters"]; ok {
				basicAuthParametersMap := map[string]interface{}{}
				for _, basicAuthParametersList := range basicAuthParameters.([]interface{}) {
					basicAuthParametersArg := basicAuthParametersList.(map[string]interface{})

					if username, ok := basicAuthParametersArg["username"]; ok {
						basicAuthParametersMap["Username"] = username
					}

					if password, ok := basicAuthParametersArg["password"]; ok {
						basicAuthParametersMap["Password"] = password
					}
				}

				authParametersMap["BasicAuthParameters"] = basicAuthParametersMap
			}

			if oAuthParameters, ok := authParametersArg["oauth_parameters"]; ok {
				oAuthParametersMap := map[string]interface{}{}
				for _, oAuthParametersList := range oAuthParameters.([]interface{}) {
					oAuthParametersArg := oAuthParametersList.(map[string]interface{})

					if authorizationEndpoint, ok := oAuthParametersArg["authorization_endpoint"]; ok {
						oAuthParametersMap["AuthorizationEndpoint"] = authorizationEndpoint
					}

					if httpMethod, ok := oAuthParametersArg["http_method"]; ok {
						oAuthParametersMap["HttpMethod"] = httpMethod
					}

					if clientParameters, ok := oAuthParametersArg["client_parameters"]; ok {
						clientParametersMap := map[string]interface{}{}
						for _, clientParametersList := range clientParameters.([]interface{}) {
							clientParametersArg := clientParametersList.(map[string]interface{})

							if clientID, ok := clientParametersArg["client_id"]; ok {
								clientParametersMap["ClientID"] = clientID
							}

							if clientSecret, ok := clientParametersArg["client_secret"]; ok {
								clientParametersMap["ClientSecret"] = clientSecret
							}
						}

						oAuthParametersMap["ClientParameters"] = clientParametersMap
					}

					if oAuthHttpParameters, ok := oAuthParametersArg["oauth_http_parameters"]; ok {
						oAuthHttpParametersMap := map[string]interface{}{}
						for _, oAuthHttpParametersList := range oAuthHttpParameters.([]interface{}) {
							oAuthHttpParametersArg := oAuthHttpParametersList.(map[string]interface{})

							if headerParameters, ok := oAuthHttpParametersArg["header_parameters"]; ok {
								headerParametersMaps := make([]map[string]interface{}, 0)
								for _, headerParametersList := range headerParameters.([]interface{}) {
									headerParametersMap := map[string]interface{}{}
									headerParametersArg := headerParametersList.(map[string]interface{})

									if key, ok := headerParametersArg["key"]; ok {
										headerParametersMap["Key"] = key
									}

									if value, ok := headerParametersArg["value"]; ok {
										headerParametersMap["Value"] = value
									}

									if isValueSecret, ok := headerParametersArg["is_value_secret"]; ok {
										headerParametersMap["IsValueSecret"] = isValueSecret
									}

									headerParametersMaps = append(headerParametersMaps, headerParametersMap)
								}

								oAuthHttpParametersMap["HeaderParameters"] = headerParametersMaps
							}

							if bodyParameters, ok := oAuthHttpParametersArg["body_parameters"]; ok {
								bodyParametersMaps := make([]map[string]interface{}, 0)
								for _, bodyParametersList := range bodyParameters.([]interface{}) {
									bodyParametersMap := map[string]interface{}{}
									bodyParametersArg := bodyParametersList.(map[string]interface{})

									if key, ok := bodyParametersArg["key"]; ok {
										bodyParametersMap["Key"] = key
									}

									if value, ok := bodyParametersArg["value"]; ok {
										bodyParametersMap["Value"] = value
									}

									if isValueSecret, ok := bodyParametersArg["is_value_secret"]; ok {
										bodyParametersMap["IsValueSecret"] = isValueSecret
									}

									bodyParametersMaps = append(bodyParametersMaps, bodyParametersMap)
								}

								oAuthHttpParametersMap["BodyParameters"] = bodyParametersMaps
							}

							if queryStringParameters, ok := oAuthHttpParametersArg["query_string_parameters"]; ok {
								queryStringParametersMaps := make([]map[string]interface{}, 0)
								for _, queryStringParametersList := range queryStringParameters.([]interface{}) {
									queryStringParametersMap := map[string]interface{}{}
									queryStringParametersArg := queryStringParametersList.(map[string]interface{})

									if key, ok := queryStringParametersArg["key"]; ok {
										queryStringParametersMap["Key"] = key
									}

									if value, ok := queryStringParametersArg["value"]; ok {
										queryStringParametersMap["Value"] = value
									}

									if isValueSecret, ok := queryStringParametersArg["is_value_secret"]; ok {
										queryStringParametersMap["IsValueSecret"] = isValueSecret
									}

									queryStringParametersMaps = append(queryStringParametersMaps, queryStringParametersMap)
								}

								oAuthHttpParametersMap["QueryStringParameters"] = queryStringParametersMaps
							}

						}

						oAuthParametersMap["OAuthHttpParameters"] = oAuthHttpParametersMap
					}
				}

				authParametersMap["OAuthParameters"] = oAuthParametersMap
			}
		}

		authParametersJson, err := convertMaptoJsonString(authParametersMap)
		if err != nil {
			return WrapError(err)
		}

		request["AuthParameters"] = authParametersJson
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_connection", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_event_bridge_connection")
	} else {
		connectionName := resp.(map[string]interface{})["ConnectionName"]
		d.SetId(fmt.Sprint(connectionName))
	}

	return resourceAliCloudEventBridgeConnectionRead(d, meta)
}

func resourceAliCloudEventBridgeConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeServiceV2 := EventBridgeServiceV2{client}

	object, err := eventBridgeServiceV2.DescribeEventBridgeConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("connection_name", object["ConnectionName"])
	d.Set("description", object["Description"])
	d.Set("create_time", object["GmtCreate"])

	if networkParameters, ok := object["NetworkParameters"]; ok {
		networkParametersMaps := make([]map[string]interface{}, 0)
		networkParametersArg := networkParameters.(map[string]interface{})
		networkParametersMap := make(map[string]interface{})

		if networkType, ok := networkParametersArg["NetworkType"]; ok {
			networkParametersMap["network_type"] = networkType
		}

		if vpcId, ok := networkParametersArg["VpcId"]; ok {
			networkParametersMap["vpc_id"] = vpcId
		}

		if vswitcheId, ok := networkParametersArg["VswitcheId"]; ok {
			networkParametersMap["vswitche_id"] = vswitcheId
		}

		if securityGroupId, ok := networkParametersArg["SecurityGroupId"]; ok {
			networkParametersMap["security_group_id"] = securityGroupId
		}

		networkParametersMaps = append(networkParametersMaps, networkParametersMap)

		d.Set("network_parameters", networkParametersMaps)
	}

	if authParameters, ok := object["AuthParameters"]; ok {
		authParametersMaps := make([]map[string]interface{}, 0)
		authParametersArg := authParameters.(map[string]interface{})
		authParametersMap := make(map[string]interface{})

		if authorizationType, ok := authParametersArg["AuthorizationType"]; ok {
			authParametersMap["authorization_type"] = authorizationType
		}

		if apiKeyAuthParameters, ok := authParametersArg["ApiKeyAuthParameters"]; ok {
			apiKeyAuthParametersMaps := make([]map[string]interface{}, 0)
			apiKeyAuthParametersArg := apiKeyAuthParameters.(map[string]interface{})
			apiKeyAuthParametersMap := make(map[string]interface{})

			if len(apiKeyAuthParametersArg) > 0 {
				if apiKeyName, ok := apiKeyAuthParametersArg["ApiKeyName"]; ok {
					apiKeyAuthParametersMap["api_key_name"] = apiKeyName
				}

				if apiKeyValue, ok := apiKeyAuthParametersArg["ApiKeyValue"]; ok {
					apiKeyAuthParametersMap["api_key_value"] = apiKeyValue
				}

				apiKeyAuthParametersMaps = append(apiKeyAuthParametersMaps, apiKeyAuthParametersMap)

				authParametersMap["api_key_auth_parameters"] = apiKeyAuthParametersMaps
			}
		}

		if basicAuthParameters, ok := authParametersArg["BasicAuthParameters"]; ok {
			basicAuthParametersMaps := make([]map[string]interface{}, 0)
			basicAuthParametersArg := basicAuthParameters.(map[string]interface{})
			basicAuthParametersMap := make(map[string]interface{})

			if len(basicAuthParametersArg) > 0 {
				if username, ok := basicAuthParametersArg["Username"]; ok {
					basicAuthParametersMap["username"] = username
				}

				if password, ok := basicAuthParametersArg["Password"]; ok {
					basicAuthParametersMap["password"] = password
				}

				basicAuthParametersMaps = append(basicAuthParametersMaps, basicAuthParametersMap)

				authParametersMap["basic_auth_parameters"] = basicAuthParametersMaps
			}
		}

		if oAuthParameters, ok := authParametersArg["OAuthParameters"]; ok {
			oAuthParametersMaps := make([]map[string]interface{}, 0)
			oAuthParametersArg := oAuthParameters.(map[string]interface{})
			oAuthParametersMap := make(map[string]interface{})

			if authorizationEndpoint, ok := oAuthParametersArg["AuthorizationEndpoint"]; ok {
				oAuthParametersMap["authorization_endpoint"] = authorizationEndpoint
			}

			if httpMethod, ok := oAuthParametersArg["HttpMethod"]; ok {
				oAuthParametersMap["http_method"] = httpMethod
			}

			if clientParameters, ok := oAuthParametersArg["ClientParameters"]; ok {
				clientParametersMaps := make([]map[string]interface{}, 0)
				clientParametersArg := clientParameters.(map[string]interface{})
				clientParametersMap := make(map[string]interface{})

				if len(clientParametersArg) > 0 {
					if clientID, ok := clientParametersArg["ClientID"]; ok {
						clientParametersMap["client_id"] = clientID
					}

					if clientSecret, ok := clientParametersArg["ClientSecret"]; ok {
						clientParametersMap["client_secret"] = clientSecret
					}

					clientParametersMaps = append(clientParametersMaps, clientParametersMap)

					oAuthParametersMap["client_parameters"] = clientParametersMaps
				}
			}

			if oAuthHttpParameters, ok := oAuthParametersArg["OAuthHttpParameters"]; ok {
				oAuthHttpParametersMaps := make([]map[string]interface{}, 0)
				oAuthHttpParametersArg := oAuthHttpParameters.(map[string]interface{})
				oAuthHttpParametersMap := make(map[string]interface{})

				if len(oAuthHttpParametersArg) > 0 {
					if headerParametersList, ok := oAuthHttpParametersArg["HeaderParameters"]; ok {
						headerParametersMaps := make([]map[string]interface{}, 0)
						for _, headerParameters := range headerParametersList.([]interface{}) {
							headerParametersArg := headerParameters.(map[string]interface{})
							headerParametersMap := map[string]interface{}{}

							if key, ok := headerParametersArg["Key"]; ok {
								headerParametersMap["key"] = key
							}

							if value, ok := headerParametersArg["Value"]; ok {
								headerParametersMap["value"] = value
							}

							if isValueSecret, ok := headerParametersArg["IsValueSecret"]; ok {
								headerParametersMap["is_value_secret"] = isValueSecret
							}

							headerParametersMaps = append(headerParametersMaps, headerParametersMap)
						}

						oAuthHttpParametersMap["header_parameters"] = headerParametersMaps
					}

					if bodyParametersList, ok := oAuthHttpParametersArg["BodyParameters"]; ok {
						bodyParametersMaps := make([]map[string]interface{}, 0)
						for _, bodyParameters := range bodyParametersList.([]interface{}) {
							bodyParametersArg := bodyParameters.(map[string]interface{})
							bodyParametersMap := map[string]interface{}{}

							if key, ok := bodyParametersArg["Key"]; ok {
								bodyParametersMap["key"] = key
							}

							if value, ok := bodyParametersArg["Value"]; ok {
								bodyParametersMap["value"] = value
							}

							if isValueSecret, ok := bodyParametersArg["IsValueSecret"]; ok {
								bodyParametersMap["is_value_secret"] = isValueSecret
							}

							bodyParametersMaps = append(bodyParametersMaps, bodyParametersMap)
						}

						oAuthHttpParametersMap["body_parameters"] = bodyParametersMaps
					}

					if queryStringParametersList, ok := oAuthHttpParametersArg["QueryStringParameters"]; ok {
						queryStringParametersMaps := make([]map[string]interface{}, 0)
						for _, queryStringParameters := range queryStringParametersList.([]interface{}) {
							queryStringParametersArg := queryStringParameters.(map[string]interface{})
							queryStringParametersMap := map[string]interface{}{}

							if key, ok := queryStringParametersArg["Key"]; ok {
								queryStringParametersMap["key"] = key
							}

							if value, ok := queryStringParametersArg["Value"]; ok {
								queryStringParametersMap["value"] = value
							}

							if isValueSecret, ok := queryStringParametersArg["IsValueSecret"]; ok {
								queryStringParametersMap["is_value_secret"] = isValueSecret
							}

							queryStringParametersMaps = append(queryStringParametersMaps, queryStringParametersMap)
						}

						oAuthHttpParametersMap["query_string_parameters"] = queryStringParametersMaps
					}

					oAuthHttpParametersMaps = append(oAuthHttpParametersMaps, oAuthHttpParametersMap)

					oAuthParametersMap["oauth_http_parameters"] = oAuthHttpParametersMaps
				}
			}

			if len(oAuthParametersMap) > 0 {
				oAuthParametersMaps = append(oAuthParametersMaps, oAuthParametersMap)

				authParametersMap["oauth_parameters"] = oAuthParametersMaps
			}
		}

		if len(authParametersMap) > 0 {
			authParametersMaps = append(authParametersMaps, authParametersMap)
		}

		d.Set("auth_parameters", authParametersMaps)
	}

	return nil
}

func resourceAliCloudEventBridgeConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"ConnectionName": d.Id(),
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("network_parameters") {
		update = true
	}
	if v, ok := d.GetOk("network_parameters"); ok {
		networkParametersMap := map[string]interface{}{}
		for _, networkParametersList := range v.([]interface{}) {
			networkParametersArg := networkParametersList.(map[string]interface{})

			networkParametersMap["NetworkType"] = networkParametersArg["network_type"]

			if vpcId, ok := networkParametersArg["vpc_id"]; ok {
				networkParametersMap["VpcId"] = vpcId
			}

			if vswitcheId, ok := networkParametersArg["vswitche_id"]; ok {
				networkParametersMap["VswitcheId"] = vswitcheId
			}

			if securityGroupId, ok := networkParametersArg["security_group_id"]; ok {
				networkParametersMap["SecurityGroupId"] = securityGroupId
			}
		}

		networkParametersJson, err := convertMaptoJsonString(networkParametersMap)
		if err != nil {
			return WrapError(err)
		}

		request["NetworkParameters"] = networkParametersJson
	}

	if d.HasChange("auth_parameters") {
		update = true
	}
	if v, ok := d.GetOk("auth_parameters"); ok {
		authParametersMap := map[string]interface{}{}
		for _, authParametersList := range v.([]interface{}) {
			authParametersArg := authParametersList.(map[string]interface{})

			if authorizationType, ok := authParametersArg["authorization_type"]; ok {
				authParametersMap["AuthorizationType"] = authorizationType
			}

			if apiKeyAuthParameters, ok := authParametersArg["api_key_auth_parameters"]; ok {
				apiKeyAuthParametersMap := map[string]interface{}{}
				for _, apiKeyAuthParametersList := range apiKeyAuthParameters.([]interface{}) {
					apiKeyAuthParametersArg := apiKeyAuthParametersList.(map[string]interface{})

					if apiKeyName, ok := apiKeyAuthParametersArg["api_key_name"]; ok {
						apiKeyAuthParametersMap["ApiKeyName"] = apiKeyName
					}

					if apiKeyValue, ok := apiKeyAuthParametersArg["api_key_value"]; ok {
						apiKeyAuthParametersMap["ApiKeyValue"] = apiKeyValue
					}
				}

				authParametersMap["ApiKeyAuthParameters"] = apiKeyAuthParametersMap
			}

			if basicAuthParameters, ok := authParametersArg["basic_auth_parameters"]; ok {
				basicAuthParametersMap := map[string]interface{}{}
				for _, basicAuthParametersList := range basicAuthParameters.([]interface{}) {
					basicAuthParametersArg := basicAuthParametersList.(map[string]interface{})

					if username, ok := basicAuthParametersArg["username"]; ok {
						basicAuthParametersMap["Username"] = username
					}

					if password, ok := basicAuthParametersArg["password"]; ok {
						basicAuthParametersMap["Password"] = password
					}
				}

				authParametersMap["BasicAuthParameters"] = basicAuthParametersMap
			}

			if oAuthParameters, ok := authParametersArg["oauth_parameters"]; ok {
				oAuthParametersMap := map[string]interface{}{}
				for _, oAuthParametersList := range oAuthParameters.([]interface{}) {
					oAuthParametersArg := oAuthParametersList.(map[string]interface{})

					if authorizationEndpoint, ok := oAuthParametersArg["authorization_endpoint"]; ok {
						oAuthParametersMap["AuthorizationEndpoint"] = authorizationEndpoint
					}

					if httpMethod, ok := oAuthParametersArg["http_method"]; ok {
						oAuthParametersMap["HttpMethod"] = httpMethod
					}

					if clientParameters, ok := oAuthParametersArg["client_parameters"]; ok {
						clientParametersMap := map[string]interface{}{}
						for _, clientParametersList := range clientParameters.([]interface{}) {
							clientParametersArg := clientParametersList.(map[string]interface{})

							if clientID, ok := clientParametersArg["client_id"]; ok {
								clientParametersMap["ClientID"] = clientID
							}

							if clientSecret, ok := clientParametersArg["client_secret"]; ok {
								clientParametersMap["ClientSecret"] = clientSecret
							}
						}

						oAuthParametersMap["ClientParameters"] = clientParametersMap
					}

					if oAuthHttpParameters, ok := oAuthParametersArg["oauth_http_parameters"]; ok {
						oAuthHttpParametersMap := map[string]interface{}{}
						for _, oAuthHttpParametersList := range oAuthHttpParameters.([]interface{}) {
							oAuthHttpParametersArg := oAuthHttpParametersList.(map[string]interface{})

							if headerParameters, ok := oAuthHttpParametersArg["header_parameters"]; ok {
								headerParametersMaps := make([]map[string]interface{}, 0)
								for _, headerParametersList := range headerParameters.([]interface{}) {
									headerParametersMap := map[string]interface{}{}
									headerParametersArg := headerParametersList.(map[string]interface{})

									if key, ok := headerParametersArg["key"]; ok {
										headerParametersMap["Key"] = key
									}

									if value, ok := headerParametersArg["value"]; ok {
										headerParametersMap["Value"] = value
									}

									if isValueSecret, ok := headerParametersArg["is_value_secret"]; ok {
										headerParametersMap["IsValueSecret"] = isValueSecret
									}

									headerParametersMaps = append(headerParametersMaps, headerParametersMap)
								}

								oAuthHttpParametersMap["HeaderParameters"] = headerParametersMaps
							}

							if bodyParameters, ok := oAuthHttpParametersArg["body_parameters"]; ok {
								bodyParametersMaps := make([]map[string]interface{}, 0)
								for _, bodyParametersList := range bodyParameters.([]interface{}) {
									bodyParametersMap := map[string]interface{}{}
									bodyParametersArg := bodyParametersList.(map[string]interface{})

									if key, ok := bodyParametersArg["key"]; ok {
										bodyParametersMap["Key"] = key
									}

									if value, ok := bodyParametersArg["value"]; ok {
										bodyParametersMap["Value"] = value
									}

									if isValueSecret, ok := bodyParametersArg["is_value_secret"]; ok {
										bodyParametersMap["IsValueSecret"] = isValueSecret
									}

									bodyParametersMaps = append(bodyParametersMaps, bodyParametersMap)
								}

								oAuthHttpParametersMap["BodyParameters"] = bodyParametersMaps
							}

							if queryStringParameters, ok := oAuthHttpParametersArg["query_string_parameters"]; ok {
								queryStringParametersMaps := make([]map[string]interface{}, 0)
								for _, queryStringParametersList := range queryStringParameters.([]interface{}) {
									queryStringParametersMap := map[string]interface{}{}
									queryStringParametersArg := queryStringParametersList.(map[string]interface{})

									if key, ok := queryStringParametersArg["key"]; ok {
										queryStringParametersMap["Key"] = key
									}

									if value, ok := queryStringParametersArg["value"]; ok {
										queryStringParametersMap["Value"] = value
									}

									if isValueSecret, ok := queryStringParametersArg["is_value_secret"]; ok {
										queryStringParametersMap["IsValueSecret"] = isValueSecret
									}

									queryStringParametersMaps = append(queryStringParametersMaps, queryStringParametersMap)
								}

								oAuthHttpParametersMap["QueryStringParameters"] = queryStringParametersMaps
							}

						}

						oAuthParametersMap["OAuthHttpParameters"] = oAuthHttpParametersMap
					}
				}

				authParametersMap["OAuthParameters"] = oAuthParametersMap
			}
		}

		authParametersJson, err := convertMaptoJsonString(authParametersMap)
		if err != nil {
			return WrapError(err)
		}

		request["AuthParameters"] = authParametersJson
	}

	if update {
		action := "UpdateConnection"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

	return resourceAliCloudEventBridgeConnectionRead(d, meta)
}

func resourceAliCloudEventBridgeConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConnection"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ConnectionName": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, false)
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
