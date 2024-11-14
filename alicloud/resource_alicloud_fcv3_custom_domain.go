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
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("custom_domain_name"); ok {
		request["domainName"] = v
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("auth_config"); !IsNil(v) {
		authInfo1, _ := jsonpath.Get("$[0].auth_info", d.Get("auth_config"))
		if authInfo1 != nil && authInfo1 != "" {
			objectDataLocalMap["authInfo"] = authInfo1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", d.Get("auth_config"))
		if authType1 != nil && authType1 != "" {
			objectDataLocalMap["authType"] = authType1
		}

		request["authConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("cert_config"); !IsNil(v) {
		certName1, _ := jsonpath.Get("$[0].cert_name", d.Get("cert_config"))
		if certName1 != nil && certName1 != "" {
			objectDataLocalMap1["certName"] = certName1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", d.Get("cert_config"))
		if certificate1 != nil && certificate1 != "" {
			objectDataLocalMap1["certificate"] = certificate1
		}
		privateKey1, _ := jsonpath.Get("$[0].private_key", d.Get("cert_config"))
		if privateKey1 != nil && privateKey1 != "" {
			objectDataLocalMap1["privateKey"] = privateKey1
		}

		request["certConfig"] = objectDataLocalMap1
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["protocol"] = v
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("tls_config"); !IsNil(v) {
		cipherSuites1, _ := jsonpath.Get("$[0].cipher_suites", v)
		if cipherSuites1 != nil && cipherSuites1 != "" {
			objectDataLocalMap2["cipherSuites"] = cipherSuites1
		}
		maxVersion1, _ := jsonpath.Get("$[0].max_version", d.Get("tls_config"))
		if maxVersion1 != nil && maxVersion1 != "" {
			objectDataLocalMap2["maxVersion"] = maxVersion1
		}
		minVersion1, _ := jsonpath.Get("$[0].min_version", d.Get("tls_config"))
		if minVersion1 != nil && minVersion1 != "" {
			objectDataLocalMap2["minVersion"] = minVersion1
		}

		request["tlsConfig"] = objectDataLocalMap2
	}

	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("route_config"); !IsNil(v) {
		if v, ok := d.GetOk("route_config"); ok {
			localData, err := jsonpath.Get("$[0].routes", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["methods"] = dataLoopTmp["methods"]
				dataLoopMap["functionName"] = dataLoopTmp["function_name"]
				dataLoopMap["path"] = dataLoopTmp["path"]
				dataLoopMap["qualifier"] = dataLoopTmp["qualifier"]
				localData1 := make(map[string]interface{})
				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData2, err := jsonpath.Get("$[0].equal_rules", v)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps2 := make([]interface{}, 0)
					for _, dataLoop2 := range localData2.([]interface{}) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["match"] = dataLoop2Tmp["match"]
						dataLoop2Map["replacement"] = dataLoop2Tmp["replacement"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					localData1["equalRules"] = localMaps2
				}

				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData3, err := jsonpath.Get("$[0].regex_rules", v)
					if err != nil {
						localData3 = make([]interface{}, 0)
					}
					localMaps3 := make([]interface{}, 0)
					for _, dataLoop3 := range localData3.([]interface{}) {
						dataLoop3Tmp := make(map[string]interface{})
						if dataLoop3 != nil {
							dataLoop3Tmp = dataLoop3.(map[string]interface{})
						}
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["match"] = dataLoop3Tmp["match"]
						dataLoop3Map["replacement"] = dataLoop3Tmp["replacement"]
						localMaps3 = append(localMaps3, dataLoop3Map)
					}
					localData1["regexRules"] = localMaps3
				}

				if v, ok := dataLoopTmp["rewrite_config"]; ok {
					localData4, err := jsonpath.Get("$[0].wildcard_rules", v)
					if err != nil {
						localData4 = make([]interface{}, 0)
					}
					localMaps4 := make([]interface{}, 0)
					for _, dataLoop4 := range localData4.([]interface{}) {
						dataLoop4Tmp := make(map[string]interface{})
						if dataLoop4 != nil {
							dataLoop4Tmp = dataLoop4.(map[string]interface{})
						}
						dataLoop4Map := make(map[string]interface{})
						dataLoop4Map["match"] = dataLoop4Tmp["match"]
						dataLoop4Map["replacement"] = dataLoop4Tmp["replacement"]
						localMaps4 = append(localMaps4, dataLoop4Map)
					}
					localData1["wildcardRules"] = localMaps4
				}

				dataLoopMap["rewriteConfig"] = localData1
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap3["routes"] = localMaps
		}

		request["routeConfig"] = objectDataLocalMap3
	}

	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("waf_config"); !IsNil(v) {
		enableWaf, _ := jsonpath.Get("$[0].enable_waf", d.Get("waf_config"))
		if enableWaf != nil && enableWaf != "" {
			objectDataLocalMap4["enableWAF"] = enableWaf
		}

		request["wafConfig"] = objectDataLocalMap4
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	id, _ := jsonpath.Get("$.body.domainName", response)
	d.SetId(fmt.Sprint(id))

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

	if objectRaw["accountId"] != nil {
		d.Set("account_id", objectRaw["accountId"])
	}
	if objectRaw["apiVersion"] != nil {
		d.Set("api_version", objectRaw["apiVersion"])
	}
	if objectRaw["createdTime"] != nil {
		d.Set("create_time", objectRaw["createdTime"])
	}
	if objectRaw["lastModifiedTime"] != nil {
		d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	}
	if objectRaw["protocol"] != nil {
		d.Set("protocol", objectRaw["protocol"])
	}
	if objectRaw["subdomainCount"] != nil {
		d.Set("subdomain_count", objectRaw["subdomainCount"])
	}
	if objectRaw["domainName"] != nil {
		d.Set("custom_domain_name", objectRaw["domainName"])
	}

	authConfigMaps := make([]map[string]interface{}, 0)
	authConfigMap := make(map[string]interface{})
	authConfig1Raw := make(map[string]interface{})
	if objectRaw["authConfig"] != nil {
		authConfig1Raw = objectRaw["authConfig"].(map[string]interface{})
	}
	if len(authConfig1Raw) > 0 {
		authConfigMap["auth_info"] = authConfig1Raw["authInfo"]
		authConfigMap["auth_type"] = authConfig1Raw["authType"]

		authConfigMaps = append(authConfigMaps, authConfigMap)
	}
	if objectRaw["authConfig"] != nil {
		if err := d.Set("auth_config", authConfigMaps); err != nil {
			return err
		}
	}
	certConfigMaps := make([]map[string]interface{}, 0)
	certConfigMap := make(map[string]interface{})
	certConfig1Raw := make(map[string]interface{})
	if objectRaw["certConfig"] != nil {
		certConfig1Raw = objectRaw["certConfig"].(map[string]interface{})
	}
	if len(certConfig1Raw) > 0 {
		certConfigMap["cert_name"] = certConfig1Raw["certName"]
		certConfigMap["certificate"] = certConfig1Raw["certificate"]
		oldConfig := d.Get("cert_config").([]interface{})
		// The FC service will not return private key crendential for security reason.
		// Read it from the terraform file.
		if len(oldConfig) > 0 {
			certConfigMap["private_key"] = oldConfig[0].(map[string]interface{})["private_key"]
		}

		certConfigMaps = append(certConfigMaps, certConfigMap)
	}
	if objectRaw["certConfig"] != nil {
		if err := d.Set("cert_config", certConfigMaps); err != nil {
			return err
		}
	}
	routeConfigMaps := make([]map[string]interface{}, 0)
	routeConfigMap := make(map[string]interface{})
	routes1Raw, _ := jsonpath.Get("$.routeConfig.routes", objectRaw)

	routesMaps := make([]map[string]interface{}, 0)
	if routes1Raw != nil {
		for _, routesChild1Raw := range routes1Raw.([]interface{}) {
			routesMap := make(map[string]interface{})
			routesChild1Raw := routesChild1Raw.(map[string]interface{})
			routesMap["function_name"] = routesChild1Raw["functionName"]
			routesMap["path"] = routesChild1Raw["path"]
			routesMap["qualifier"] = routesChild1Raw["qualifier"]

			methods1Raw := make([]interface{}, 0)
			if routesChild1Raw["methods"] != nil {
				methods1Raw = routesChild1Raw["methods"].([]interface{})
			}

			routesMap["methods"] = methods1Raw
			rewriteConfigMaps := make([]map[string]interface{}, 0)
			rewriteConfigMap := make(map[string]interface{})
			rewriteConfig1Raw := make(map[string]interface{})
			if routesChild1Raw["rewriteConfig"] != nil {
				rewriteConfig1Raw = routesChild1Raw["rewriteConfig"].(map[string]interface{})
			}

			equalRules1Raw := rewriteConfig1Raw["equalRules"]
			equalRulesMaps := make([]map[string]interface{}, 0)
			if equalRules1Raw != nil {
				for _, equalRulesChild1Raw := range equalRules1Raw.([]interface{}) {
					equalRulesMap := make(map[string]interface{})
					equalRulesChild1Raw := equalRulesChild1Raw.(map[string]interface{})
					equalRulesMap["match"] = equalRulesChild1Raw["match"]
					equalRulesMap["replacement"] = equalRulesChild1Raw["replacement"]

					equalRulesMaps = append(equalRulesMaps, equalRulesMap)
				}
			}
			rewriteConfigMap["equal_rules"] = equalRulesMaps
			regexRules1Raw := rewriteConfig1Raw["regexRules"]
			regexRulesMaps := make([]map[string]interface{}, 0)
			if regexRules1Raw != nil {
				for _, regexRulesChild1Raw := range regexRules1Raw.([]interface{}) {
					regexRulesMap := make(map[string]interface{})
					regexRulesChild1Raw := regexRulesChild1Raw.(map[string]interface{})
					regexRulesMap["match"] = regexRulesChild1Raw["match"]
					regexRulesMap["replacement"] = regexRulesChild1Raw["replacement"]

					regexRulesMaps = append(regexRulesMaps, regexRulesMap)
				}
			}
			rewriteConfigMap["regex_rules"] = regexRulesMaps
			wildcardRules1Raw := rewriteConfig1Raw["wildcardRules"]
			wildcardRulesMaps := make([]map[string]interface{}, 0)
			if wildcardRules1Raw != nil {
				for _, wildcardRulesChild1Raw := range wildcardRules1Raw.([]interface{}) {
					wildcardRulesMap := make(map[string]interface{})
					wildcardRulesChild1Raw := wildcardRulesChild1Raw.(map[string]interface{})
					wildcardRulesMap["match"] = wildcardRulesChild1Raw["match"]
					wildcardRulesMap["replacement"] = wildcardRulesChild1Raw["replacement"]

					wildcardRulesMaps = append(wildcardRulesMaps, wildcardRulesMap)
				}
			}
			rewriteConfigMap["wildcard_rules"] = wildcardRulesMaps
			rewriteConfigMaps = append(rewriteConfigMaps, rewriteConfigMap)
			routesMap["rewrite_config"] = rewriteConfigMaps
			routesMaps = append(routesMaps, routesMap)
		}
	}
	routeConfigMap["routes"] = routesMaps
	routeConfigMaps = append(routeConfigMaps, routeConfigMap)
	if routes1Raw != nil {
		if err := d.Set("route_config", routeConfigMaps); err != nil {
			return err
		}
	}
	tlsConfigMaps := make([]map[string]interface{}, 0)
	tlsConfigMap := make(map[string]interface{})
	tlsConfig1Raw := make(map[string]interface{})
	if objectRaw["tlsConfig"] != nil {
		tlsConfig1Raw = objectRaw["tlsConfig"].(map[string]interface{})
	}
	if len(tlsConfig1Raw) > 0 {
		tlsConfigMap["max_version"] = tlsConfig1Raw["maxVersion"]
		tlsConfigMap["min_version"] = tlsConfig1Raw["minVersion"]

		cipherSuites1Raw := make([]interface{}, 0)
		if tlsConfig1Raw["cipherSuites"] != nil {
			cipherSuites1Raw = tlsConfig1Raw["cipherSuites"].([]interface{})
		}

		tlsConfigMap["cipher_suites"] = cipherSuites1Raw
		tlsConfigMaps = append(tlsConfigMaps, tlsConfigMap)
	}
	if objectRaw["tlsConfig"] != nil {
		if err := d.Set("tls_config", tlsConfigMaps); err != nil {
			return err
		}
	}
	wafConfigMaps := make([]map[string]interface{}, 0)
	wafConfigMap := make(map[string]interface{})
	wafConfig1Raw := make(map[string]interface{})
	if objectRaw["wafConfig"] != nil {
		wafConfig1Raw = objectRaw["wafConfig"].(map[string]interface{})
	}
	if len(wafConfig1Raw) > 0 {
		wafConfigMap["enable_waf"] = wafConfig1Raw["enableWAF"]

		wafConfigMaps = append(wafConfigMaps, wafConfigMap)
	}
	if objectRaw["wafConfig"] != nil {
		if err := d.Set("waf_config", wafConfigMaps); err != nil {
			return err
		}
	}

	d.Set("custom_domain_name", d.Id())

	return nil
}

func resourceAliCloudFcv3CustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	domainName := d.Id()
	action := fmt.Sprintf("/2023-03-30/custom-domains/%s", domainName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["domainName"] = d.Id()

	if d.HasChange("auth_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("auth_config"); !IsNil(v) {
			authInfo1, _ := jsonpath.Get("$[0].auth_info", v)
			if authInfo1 != nil && (d.HasChange("auth_config.0.auth_info") || authInfo1 != "") {
				objectDataLocalMap["authInfo"] = authInfo1
			}
			authType1, _ := jsonpath.Get("$[0].auth_type", v)
			if authType1 != nil && (d.HasChange("auth_config.0.auth_type") || authType1 != "") {
				objectDataLocalMap["authType"] = authType1
			}

			request["authConfig"] = objectDataLocalMap
		}
	}

	if d.HasChange("cert_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("cert_config"); !IsNil(v) {
			certName1, _ := jsonpath.Get("$[0].cert_name", v)
			if certName1 != nil && (d.HasChange("cert_config.0.cert_name") || certName1 != "") {
				objectDataLocalMap1["certName"] = certName1
			}
			certificate1, _ := jsonpath.Get("$[0].certificate", v)
			if certificate1 != nil && (d.HasChange("cert_config.0.certificate") || certificate1 != "") {
				objectDataLocalMap1["certificate"] = certificate1
			}
			privateKey1, _ := jsonpath.Get("$[0].private_key", v)
			if privateKey1 != nil && (d.HasChange("cert_config.0.private_key") || privateKey1 != "") {
				objectDataLocalMap1["privateKey"] = privateKey1
			}

			request["certConfig"] = objectDataLocalMap1
		}
	}

	if d.HasChange("protocol") {
		update = true
		request["protocol"] = d.Get("protocol")
	}

	if d.HasChange("tls_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})

		if v := d.Get("tls_config"); !IsNil(v) {
			cipherSuites1, _ := jsonpath.Get("$[0].cipher_suites", d.Get("tls_config"))
			if cipherSuites1 != nil && (d.HasChange("tls_config.0.cipher_suites") || cipherSuites1 != "") {
				objectDataLocalMap2["cipherSuites"] = cipherSuites1
			}
			maxVersion1, _ := jsonpath.Get("$[0].max_version", v)
			if maxVersion1 != nil && (d.HasChange("tls_config.0.max_version") || maxVersion1 != "") {
				objectDataLocalMap2["maxVersion"] = maxVersion1
			}
			minVersion1, _ := jsonpath.Get("$[0].min_version", v)
			if minVersion1 != nil && (d.HasChange("tls_config.0.min_version") || minVersion1 != "") {
				objectDataLocalMap2["minVersion"] = minVersion1
			}

			request["tlsConfig"] = objectDataLocalMap2
		}
	}

	if d.HasChange("route_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

		if v := d.Get("route_config"); !IsNil(v) {
			if v, ok := d.GetOk("route_config"); ok {
				localData, err := jsonpath.Get("$[0].routes", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["methods"] = dataLoopTmp["methods"]
					dataLoopMap["functionName"] = dataLoopTmp["function_name"]
					dataLoopMap["path"] = dataLoopTmp["path"]
					dataLoopMap["qualifier"] = dataLoopTmp["qualifier"]
					if !IsNil(dataLoopTmp["rewrite_config"]) {
						localData1 := make(map[string]interface{})
						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData2, err := jsonpath.Get("$[0].equal_rules", v)
							if err != nil {
								localData2 = make([]interface{}, 0)
							}
							localMaps2 := make([]interface{}, 0)
							for _, dataLoop2 := range localData2.([]interface{}) {
								dataLoop2Tmp := make(map[string]interface{})
								if dataLoop2 != nil {
									dataLoop2Tmp = dataLoop2.(map[string]interface{})
								}
								dataLoop2Map := make(map[string]interface{})
								dataLoop2Map["match"] = dataLoop2Tmp["match"]
								dataLoop2Map["replacement"] = dataLoop2Tmp["replacement"]
								localMaps2 = append(localMaps2, dataLoop2Map)
							}
							localData1["equalRules"] = localMaps2
						}

						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData3, err := jsonpath.Get("$[0].regex_rules", v)
							if err != nil {
								localData3 = make([]interface{}, 0)
							}
							localMaps3 := make([]interface{}, 0)
							for _, dataLoop3 := range localData3.([]interface{}) {
								dataLoop3Tmp := make(map[string]interface{})
								if dataLoop3 != nil {
									dataLoop3Tmp = dataLoop3.(map[string]interface{})
								}
								dataLoop3Map := make(map[string]interface{})
								dataLoop3Map["match"] = dataLoop3Tmp["match"]
								dataLoop3Map["replacement"] = dataLoop3Tmp["replacement"]
								localMaps3 = append(localMaps3, dataLoop3Map)
							}
							localData1["regexRules"] = localMaps3
						}

						if v, ok := dataLoopTmp["rewrite_config"]; ok {
							localData4, err := jsonpath.Get("$[0].wildcard_rules", v)
							if err != nil {
								localData4 = make([]interface{}, 0)
							}
							localMaps4 := make([]interface{}, 0)
							for _, dataLoop4 := range localData4.([]interface{}) {
								dataLoop4Tmp := make(map[string]interface{})
								if dataLoop4 != nil {
									dataLoop4Tmp = dataLoop4.(map[string]interface{})
								}
								dataLoop4Map := make(map[string]interface{})
								dataLoop4Map["match"] = dataLoop4Tmp["match"]
								dataLoop4Map["replacement"] = dataLoop4Tmp["replacement"]
								localMaps4 = append(localMaps4, dataLoop4Map)
							}
							localData1["wildcardRules"] = localMaps4
						}

						dataLoopMap["rewriteConfig"] = localData1
					}
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap3["routes"] = localMaps
			}

			request["routeConfig"] = objectDataLocalMap3
		}
	}

	if d.HasChange("waf_config") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})

		if v := d.Get("waf_config"); !IsNil(v) {
			enableWaf, _ := jsonpath.Get("$[0].enable_waf", v)
			if enableWaf != nil && (d.HasChange("waf_config.0.enable_waf") || enableWaf != "") {
				objectDataLocalMap4["enableWAF"] = enableWaf
			}

			request["wafConfig"] = objectDataLocalMap4
		}
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["domainName"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
