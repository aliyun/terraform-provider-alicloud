package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudFcv3CustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3CustomDomainCreate,
		Read:   resourceAliCloudFcv3CustomDomainRead,
		Update: resourceAliCloudFcv3CustomDomainUpdate,
		Delete: resourceAliCloudFcv3CustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_info": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"anonymous", "function", "jwt"}, false),
						},
					},
				},
			},
			"cert_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"cert_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cors_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_origins": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allow_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"allow_methods": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "HTTP,HTTPS"}, false),
			},
			"route_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"function_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"qualifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"rewrite_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"wildcard_rules": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"replacement": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"match": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"regex_rules": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"replacement": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"match": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"equal_rules": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"replacement": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"match": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"methods": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"subdomain_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"TLSv1.3", "TLSv1.2", "TLSv1.1", "TLSv1.0"}, false),
						},
						"max_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"TLSv1.3", "TLSv1.2", "TLSv1.1", "TLSv1.0"}, false),
						},
						"cipher_suites": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"waf_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_waf": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudFcv3CustomDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/2023-03-30/custom-domains")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("custom_domain_name"); ok {
		request["domainName"] = v
	}

	routeConfig := make(map[string]interface{})

	if v := d.Get("route_config"); !IsNil(v) {
		if v, ok := d.GetOk("route_config"); ok {
			localData, err := jsonpath.Get("$[0].routes", v)
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
				dataLoopMap["path"] = dataLoopTmp["path"]
				localData1 := make(map[string]interface{})
				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData2, err := jsonpath.Get("$[0].wildcard_rules", v)
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
						dataLoop2Map["replacement"] = dataLoop2Tmp["replacement"]
						dataLoop2Map["match"] = dataLoop2Tmp["match"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					localData1["wildcardRules"] = localMaps2

				}

				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData3, err := jsonpath.Get("$[0].regex_rules", v)
					if err != nil {
						localData3 = make([]interface{}, 0)
					}
					localMaps3 := make([]interface{}, 0)
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := make(map[string]interface{})
						if dataLoop3 != nil {
							dataLoop3Tmp = dataLoop3.(map[string]interface{})
						}
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["replacement"] = dataLoop3Tmp["replacement"]
						dataLoop3Map["match"] = dataLoop3Tmp["match"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					localData1["regexRules"] = localMaps3

				}

				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData4, err := jsonpath.Get("$[0].equal_rules", v)
					if err != nil {
						localData4 = make([]interface{}, 0)
					}
					localMaps4 := make([]interface{}, 0)
					for _, dataLoop4 := range convertToInterfaceArray(localData4) {
						dataLoop4Tmp := make(map[string]interface{})
						if dataLoop4 != nil {
							dataLoop4Tmp = dataLoop4.(map[string]interface{})
						}
						dataLoop4Map := make(map[string]interface{})
						dataLoop4Map["match"] = dataLoop4Tmp["match"]
						dataLoop4Map["replacement"] = dataLoop4Tmp["replacement"]
						localMaps4 = append(localMaps4, dataLoop4Map)
					}
					localData1["equalRules"] = localMaps4

				}

				if len(localData1) > 0 {
					dataLoopMap["rewriteConfig"] = localData1
				}
				dataLoopMap["qualifier"] = dataLoopTmp["qualifier"]
				dataLoopMap["functionName"] = dataLoopTmp["function_name"]
				dataLoopMap["methods"] = dataLoopTmp["methods"]
				localMaps = append(localMaps, dataLoopMap)
			}
			routeConfig["routes"] = localMaps
		}

		request["routeConfig"] = routeConfig
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["protocol"] = v
	}
	certConfig := make(map[string]interface{})

	if v := d.Get("cert_config"); !IsNil(v) {
		certName1, _ := jsonpath.Get("$[0].cert_name", v)
		if certName1 != nil && certName1 != "" {
			certConfig["certName"] = certName1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && certificate1 != "" {
			certConfig["certificate"] = certificate1
		}
		privateKey1, _ := jsonpath.Get("$[0].private_key", v)
		if privateKey1 != nil && privateKey1 != "" {
			certConfig["privateKey"] = privateKey1
		}

		request["certConfig"] = certConfig
	}

	authConfig := make(map[string]interface{})

	if v := d.Get("auth_config"); !IsNil(v) {
		authInfo1, _ := jsonpath.Get("$[0].auth_info", v)
		if authInfo1 != nil && authInfo1 != "" {
			authConfig["authInfo"] = authInfo1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", v)
		if authType1 != nil && authType1 != "" {
			authConfig["authType"] = authType1
		}

		request["authConfig"] = authConfig
	}

	corsConfig := make(map[string]interface{})

	if v := d.Get("cors_config"); !IsNil(v) {
		allowHeaders1, _ := jsonpath.Get("$[0].allow_headers", v)
		if allowHeaders1 != nil && allowHeaders1 != "" {
			corsConfig["allowHeaders"] = allowHeaders1
		}
		allowMethods1, _ := jsonpath.Get("$[0].allow_methods", v)
		if allowMethods1 != nil && allowMethods1 != "" {
			corsConfig["allowMethods"] = allowMethods1
		}
		exposeHeaders1, _ := jsonpath.Get("$[0].expose_headers", v)
		if exposeHeaders1 != nil && exposeHeaders1 != "" {
			corsConfig["exposeHeaders"] = exposeHeaders1
		}
		allowOrigins1, _ := jsonpath.Get("$[0].allow_origins", v)
		if allowOrigins1 != nil && allowOrigins1 != "" {
			corsConfig["allowOrigins"] = allowOrigins1
		}
		maxAge1, _ := jsonpath.Get("$[0].max_age", v)
		if maxAge1 != nil && maxAge1 != "" {
			corsConfig["maxAge"] = maxAge1
		}

		request["corsConfig"] = corsConfig
	}

	tlsConfig := make(map[string]interface{})

	if v := d.Get("tls_config"); !IsNil(v) {
		minVersion1, _ := jsonpath.Get("$[0].min_version", v)
		if minVersion1 != nil && minVersion1 != "" {
			tlsConfig["minVersion"] = minVersion1
		}
		maxVersion1, _ := jsonpath.Get("$[0].max_version", v)
		if maxVersion1 != nil && maxVersion1 != "" {
			tlsConfig["maxVersion"] = maxVersion1
		}
		cipherSuites1, _ := jsonpath.Get("$[0].cipher_suites", v)
		if cipherSuites1 != nil && cipherSuites1 != "" {
			tlsConfig["cipherSuites"] = cipherSuites1
		}

		request["tlsConfig"] = tlsConfig
	}

	wafConfig := make(map[string]interface{})

	if v := d.Get("waf_config"); !IsNil(v) {
		enableWaf, _ := jsonpath.Get("$[0].enable_waf", v)
		if enableWaf != nil && enableWaf != "" {
			wafConfig["enableWAF"] = enableWaf
		}

		request["wafConfig"] = wafConfig
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_custom_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["domainName"]))

	return resourceAliCloudFcv3CustomDomainRead(d, meta)
}

func resourceAliCloudFcv3CustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3CustomDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_custom_domain DescribeFcv3CustomDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_id", objectRaw["accountId"])
	d.Set("api_version", objectRaw["apiVersion"])
	d.Set("create_time", objectRaw["createdTime"])
	d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	d.Set("protocol", objectRaw["protocol"])
	d.Set("subdomain_count", objectRaw["subdomainCount"])
	d.Set("custom_domain_name", objectRaw["domainName"])

	authConfigMaps := make([]map[string]interface{}, 0)
	authConfigMap := make(map[string]interface{})
	authConfigRaw := make(map[string]interface{})
	if objectRaw["authConfig"] != nil {
		authConfigRaw = objectRaw["authConfig"].(map[string]interface{})
	}
	if len(authConfigRaw) > 0 {
		authConfigMap["auth_info"] = authConfigRaw["authInfo"]
		authConfigMap["auth_type"] = authConfigRaw["authType"]

		authConfigMaps = append(authConfigMaps, authConfigMap)
	}
	if err := d.Set("auth_config", authConfigMaps); err != nil {
		return err
	}
	certConfigMaps := make([]map[string]interface{}, 0)
	certConfigMap := make(map[string]interface{})
	certConfigRaw := make(map[string]interface{})
	if objectRaw["certConfig"] != nil {
		certConfigRaw = objectRaw["certConfig"].(map[string]interface{})
	}
	if len(certConfigRaw) > 0 {
		certConfigMap["cert_name"] = certConfigRaw["certName"]
		certConfigMap["certificate"] = certConfigRaw["certificate"]
		oldConfig := d.Get("cert_config").([]interface{})
		// The FC service will not return private key crendential for security reason.
		// Read it from the terraform file.
		if len(oldConfig) > 0 {
			certConfigMap["private_key"] = oldConfig[0].(map[string]interface{})["private_key"]
		}

		certConfigMaps = append(certConfigMaps, certConfigMap)
	}
	if err := d.Set("cert_config", certConfigMaps); err != nil {
		return err
	}
	corsConfigMaps := make([]map[string]interface{}, 0)
	corsConfigMap := make(map[string]interface{})
	corsConfigRaw := make(map[string]interface{})
	if objectRaw["corsConfig"] != nil {
		corsConfigRaw = objectRaw["corsConfig"].(map[string]interface{})
	}
	if len(corsConfigRaw) > 0 {
		corsConfigMap["allow_credentials"] = corsConfigRaw["allowCredentials"]
		corsConfigMap["max_age"] = corsConfigRaw["maxAge"]

		allowHeadersRaw := make([]interface{}, 0)
		if corsConfigRaw["allowHeaders"] != nil {
			allowHeadersRaw = convertToInterfaceArray(corsConfigRaw["allowHeaders"])
		}

		corsConfigMap["allow_headers"] = allowHeadersRaw
		allowMethodsRaw := make([]interface{}, 0)
		if corsConfigRaw["allowMethods"] != nil {
			allowMethodsRaw = convertToInterfaceArray(corsConfigRaw["allowMethods"])
		}

		corsConfigMap["allow_methods"] = allowMethodsRaw
		allowOriginsRaw := make([]interface{}, 0)
		if corsConfigRaw["allowOrigins"] != nil {
			allowOriginsRaw = convertToInterfaceArray(corsConfigRaw["allowOrigins"])
		}

		corsConfigMap["allow_origins"] = allowOriginsRaw
		exposeHeadersRaw := make([]interface{}, 0)
		if corsConfigRaw["exposeHeaders"] != nil {
			exposeHeadersRaw = convertToInterfaceArray(corsConfigRaw["exposeHeaders"])
		}

		corsConfigMap["expose_headers"] = exposeHeadersRaw
		corsConfigMaps = append(corsConfigMaps, corsConfigMap)
	}
	if err := d.Set("cors_config", corsConfigMaps); err != nil {
		return err
	}
	routeConfigMaps := make([]map[string]interface{}, 0)
	routeConfigMap := make(map[string]interface{})
	routesRaw, _ := jsonpath.Get("$.routeConfig.routes", objectRaw)

	routesMaps := make([]map[string]interface{}, 0)
	if routesRaw != nil {
		for _, routesChildRaw := range convertToInterfaceArray(routesRaw) {
			routesMap := make(map[string]interface{})
			routesChildRaw := routesChildRaw.(map[string]interface{})
			routesMap["function_name"] = routesChildRaw["functionName"]
			routesMap["path"] = routesChildRaw["path"]
			routesMap["qualifier"] = routesChildRaw["qualifier"]

			methodsRaw := make([]interface{}, 0)
			if routesChildRaw["methods"] != nil {
				methodsRaw = convertToInterfaceArray(routesChildRaw["methods"])
			}

			routesMap["methods"] = methodsRaw
			rewriteConfigMaps := make([]map[string]interface{}, 0)
			rewriteConfigMap := make(map[string]interface{})
			rewriteConfigRaw := make(map[string]interface{})
			if routesChildRaw["rewriteConfig"] != nil {
				rewriteConfigRaw = routesChildRaw["rewriteConfig"].(map[string]interface{})
			}
			if len(rewriteConfigRaw) > 0 {

				equalRulesRaw := rewriteConfigRaw["equalRules"]
				equalRulesMaps := make([]map[string]interface{}, 0)
				if equalRulesRaw != nil {
					for _, equalRulesChildRaw := range convertToInterfaceArray(equalRulesRaw) {
						equalRulesMap := make(map[string]interface{})
						equalRulesChildRaw := equalRulesChildRaw.(map[string]interface{})
						equalRulesMap["match"] = equalRulesChildRaw["match"]
						equalRulesMap["replacement"] = equalRulesChildRaw["replacement"]

						equalRulesMaps = append(equalRulesMaps, equalRulesMap)
					}
				}
				rewriteConfigMap["equal_rules"] = equalRulesMaps
				regexRulesRaw := rewriteConfigRaw["regexRules"]
				regexRulesMaps := make([]map[string]interface{}, 0)
				if regexRulesRaw != nil {
					for _, regexRulesChildRaw := range convertToInterfaceArray(regexRulesRaw) {
						regexRulesMap := make(map[string]interface{})
						regexRulesChildRaw := regexRulesChildRaw.(map[string]interface{})
						regexRulesMap["match"] = regexRulesChildRaw["match"]
						regexRulesMap["replacement"] = regexRulesChildRaw["replacement"]

						regexRulesMaps = append(regexRulesMaps, regexRulesMap)
					}
				}
				rewriteConfigMap["regex_rules"] = regexRulesMaps
				wildcardRulesRaw := rewriteConfigRaw["wildcardRules"]
				wildcardRulesMaps := make([]map[string]interface{}, 0)
				if wildcardRulesRaw != nil {
					for _, wildcardRulesChildRaw := range convertToInterfaceArray(wildcardRulesRaw) {
						wildcardRulesMap := make(map[string]interface{})
						wildcardRulesChildRaw := wildcardRulesChildRaw.(map[string]interface{})
						wildcardRulesMap["match"] = wildcardRulesChildRaw["match"]
						wildcardRulesMap["replacement"] = wildcardRulesChildRaw["replacement"]

						wildcardRulesMaps = append(wildcardRulesMaps, wildcardRulesMap)
					}
				}
				rewriteConfigMap["wildcard_rules"] = wildcardRulesMaps
				rewriteConfigMaps = append(rewriteConfigMaps, rewriteConfigMap)
			}
			routesMap["rewrite_config"] = rewriteConfigMaps
			routesMaps = append(routesMaps, routesMap)
		}
	}
	routeConfigMap["routes"] = routesMaps
	routeConfigMaps = append(routeConfigMaps, routeConfigMap)
	if err := d.Set("route_config", routeConfigMaps); err != nil {
		return err
	}
	tlsConfigMaps := make([]map[string]interface{}, 0)
	tlsConfigMap := make(map[string]interface{})
	tlsConfigRaw := make(map[string]interface{})
	if objectRaw["tlsConfig"] != nil {
		tlsConfigRaw = objectRaw["tlsConfig"].(map[string]interface{})
	}
	if len(tlsConfigRaw) > 0 {
		tlsConfigMap["max_version"] = tlsConfigRaw["maxVersion"]
		tlsConfigMap["min_version"] = tlsConfigRaw["minVersion"]

		cipherSuitesRaw := make([]interface{}, 0)
		if tlsConfigRaw["cipherSuites"] != nil {
			cipherSuitesRaw = convertToInterfaceArray(tlsConfigRaw["cipherSuites"])
		}

		tlsConfigMap["cipher_suites"] = cipherSuitesRaw
		tlsConfigMaps = append(tlsConfigMaps, tlsConfigMap)
	}
	if err := d.Set("tls_config", tlsConfigMaps); err != nil {
		return err
	}
	wafConfigMaps := make([]map[string]interface{}, 0)
	wafConfigMap := make(map[string]interface{})
	wafConfigRaw := make(map[string]interface{})
	if objectRaw["wafConfig"] != nil {
		wafConfigRaw = objectRaw["wafConfig"].(map[string]interface{})
	}
	if len(wafConfigRaw) > 0 {
		wafConfigMap["enable_waf"] = wafConfigRaw["enableWAF"]

		wafConfigMaps = append(wafConfigMaps, wafConfigMap)
	}
	if err := d.Set("waf_config", wafConfigMaps); err != nil {
		return err
	}

	d.Set("custom_domain_name", d.Id())

	return nil
}

func resourceAliCloudFcv3CustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	domainName := d.Id()
	action := fmt.Sprintf("/2023-03-30/custom-domains/%s", domainName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("cors_config") {
		update = true
		corsConfig := make(map[string]interface{})

		if v := d.Get("cors_config"); v != nil {
			allowCredentials1, _ := jsonpath.Get("$[0].allow_credentials", v)
			if allowCredentials1 != nil && allowCredentials1 != "" {
				corsConfig["allowCredentials"] = allowCredentials1
			}
			allowHeaders1, _ := jsonpath.Get("$[0].allow_headers", v)
			if allowHeaders1 != nil && allowHeaders1 != "" {
				corsConfig["allowHeaders"] = allowHeaders1
			}
			allowMethods1, _ := jsonpath.Get("$[0].allow_methods", v)
			if allowMethods1 != nil && allowMethods1 != "" {
				corsConfig["allowMethods"] = allowMethods1
			}
			exposeHeaders1, _ := jsonpath.Get("$[0].expose_headers", v)
			if exposeHeaders1 != nil && exposeHeaders1 != "" {
				corsConfig["exposeHeaders"] = exposeHeaders1
			}
			allowOrigins1, _ := jsonpath.Get("$[0].allow_origins", v)
			if allowOrigins1 != nil && allowOrigins1 != "" {
				corsConfig["allowOrigins"] = allowOrigins1
			}
			maxAge1, _ := jsonpath.Get("$[0].max_age", v)
			if maxAge1 != nil && maxAge1 != "" {
				corsConfig["maxAge"] = maxAge1
			}

			request["corsConfig"] = corsConfig
		}
	}

	if d.HasChange("route_config") {
		update = true
		routeConfig := make(map[string]interface{})

		if v := d.Get("route_config"); v != nil {
			if v, ok := d.GetOk("route_config"); ok {
				localData, err := jsonpath.Get("$[0].routes", v)
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
					dataLoopMap["path"] = dataLoopTmp["path"]
					if !IsNil(dataLoopTmp["rewrite_config"]) {
						localData1 := make(map[string]interface{})
						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData2, err := jsonpath.Get("$[0].wildcard_rules", v)
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
								dataLoop2Map["replacement"] = dataLoop2Tmp["replacement"]
								dataLoop2Map["match"] = dataLoop2Tmp["match"]
								localMaps2 = append(localMaps2, dataLoop2Map)
							}
							localData1["wildcardRules"] = localMaps2

						}

						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData3, err := jsonpath.Get("$[0].regex_rules", v)
							if err != nil {
								localData3 = make([]interface{}, 0)
							}
							localMaps3 := make([]interface{}, 0)
							for _, dataLoop3 := range convertToInterfaceArray(localData3) {
								dataLoop3Tmp := make(map[string]interface{})
								if dataLoop3 != nil {
									dataLoop3Tmp = dataLoop3.(map[string]interface{})
								}
								dataLoop3Map := make(map[string]interface{})
								dataLoop3Map["replacement"] = dataLoop3Tmp["replacement"]
								dataLoop3Map["match"] = dataLoop3Tmp["match"]
								localMaps3 = append(localMaps3, dataLoop3Map)
							}
							localData1["regexRules"] = localMaps3

						}

						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData4, err := jsonpath.Get("$[0].equal_rules", v)
							if err != nil {
								localData4 = make([]interface{}, 0)
							}
							localMaps4 := make([]interface{}, 0)
							for _, dataLoop4 := range convertToInterfaceArray(localData4) {
								dataLoop4Tmp := make(map[string]interface{})
								if dataLoop4 != nil {
									dataLoop4Tmp = dataLoop4.(map[string]interface{})
								}
								dataLoop4Map := make(map[string]interface{})
								dataLoop4Map["match"] = dataLoop4Tmp["match"]
								dataLoop4Map["replacement"] = dataLoop4Tmp["replacement"]
								localMaps4 = append(localMaps4, dataLoop4Map)
							}
							localData1["equalRules"] = localMaps4

						}

						if len(localData1) > 0 {
							dataLoopMap["rewriteConfig"] = localData1
						}
					}
					dataLoopMap["qualifier"] = dataLoopTmp["qualifier"]
					dataLoopMap["functionName"] = dataLoopTmp["function_name"]
					dataLoopMap["methods"] = dataLoopTmp["methods"]
					localMaps = append(localMaps, dataLoopMap)
				}
				routeConfig["routes"] = localMaps
			}

			request["routeConfig"] = routeConfig
		}
	}

	if d.HasChange("protocol") {
		update = true
		request["protocol"] = d.Get("protocol")
	}

	if d.HasChange("cert_config") {
		update = true
		certConfig := make(map[string]interface{})

		if v := d.Get("cert_config"); v != nil {
			certName1, _ := jsonpath.Get("$[0].cert_name", v)
			if certName1 != nil && certName1 != "" {
				certConfig["certName"] = certName1
			}
			certificate1, _ := jsonpath.Get("$[0].certificate", v)
			if certificate1 != nil && certificate1 != "" {
				certConfig["certificate"] = certificate1
			}
			privateKey1, _ := jsonpath.Get("$[0].private_key", v)
			if privateKey1 != nil && privateKey1 != "" {
				certConfig["privateKey"] = privateKey1
			}

			request["certConfig"] = certConfig
		}
	}

	if d.HasChange("auth_config") {
		update = true
		authConfig := make(map[string]interface{})

		if v := d.Get("auth_config"); v != nil {
			authInfo1, _ := jsonpath.Get("$[0].auth_info", v)
			if authInfo1 != nil && authInfo1 != "" {
				authConfig["authInfo"] = authInfo1
			}
			authType1, _ := jsonpath.Get("$[0].auth_type", v)
			if authType1 != nil && authType1 != "" {
				authConfig["authType"] = authType1
			}

			request["authConfig"] = authConfig
		}
	}

	if d.HasChange("tls_config") {
		update = true
		tlsConfig := make(map[string]interface{})

		if v := d.Get("tls_config"); v != nil {
			minVersion1, _ := jsonpath.Get("$[0].min_version", v)
			if minVersion1 != nil && minVersion1 != "" {
				tlsConfig["minVersion"] = minVersion1
			}
			maxVersion1, _ := jsonpath.Get("$[0].max_version", v)
			if maxVersion1 != nil && maxVersion1 != "" {
				tlsConfig["maxVersion"] = maxVersion1
			}
			cipherSuites1, _ := jsonpath.Get("$[0].cipher_suites", v)
			if cipherSuites1 != nil && cipherSuites1 != "" {
				tlsConfig["cipherSuites"] = cipherSuites1
			}

			request["tlsConfig"] = tlsConfig
		}
	}

	if d.HasChange("waf_config") {
		update = true
		wafConfig := make(map[string]interface{})

		if v := d.Get("waf_config"); v != nil {
			enableWaf, _ := jsonpath.Get("$[0].enable_waf", v)
			if enableWaf != nil && enableWaf != "" {
				wafConfig["enableWAF"] = enableWaf
			}

			request["wafConfig"] = wafConfig
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, header, body, true)
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

	return resourceAliCloudFcv3CustomDomainRead(d, meta)
}

func resourceAliCloudFcv3CustomDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	domainName := d.Id()
	action := fmt.Sprintf("/2023-03-30/custom-domains/%s", domainName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"429"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNameNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
