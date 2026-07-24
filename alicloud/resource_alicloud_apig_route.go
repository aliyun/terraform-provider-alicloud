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

func resourceAliCloudApigRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigRouteCreate,
		Read:   resourceAliCloudApigRouteRead,
		Update: resourceAliCloudApigRouteUpdate,
		Delete: resourceAliCloudApigRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backend": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"services": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"service_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scene": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"builtin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"environment_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_edition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gateway_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_api_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"match": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
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
						"query_params": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
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
						"methods": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ignore_uri_case": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"route_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigRouteCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	httpApiId := d.Get("http_api_id")
	action := fmt.Sprintf("/v1/http-apis/%s/routes", httpApiId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	match := make(map[string]interface{})

	if v := d.Get("match"); !IsNil(v) {
		if v, ok := d.GetOk("match"); ok {
			localData, err := jsonpath.Get("$[0].headers", v)
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
				dataLoopMap["type"] = dataLoopTmp["type"]
				dataLoopMap["name"] = dataLoopTmp["name"]
				dataLoopMap["value"] = dataLoopTmp["value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			match["headers"] = localMaps
		}

		path := make(map[string]interface{})
		value3, _ := jsonpath.Get("$[0].path[0].value", d.Get("match"))
		if value3 != nil && value3 != "" {
			path["value"] = value3
		}
		type3, _ := jsonpath.Get("$[0].path[0].type", d.Get("match"))
		if type3 != nil && type3 != "" {
			path["type"] = type3
		}

		if len(path) > 0 {
			match["path"] = path
		}
		if v, ok := d.GetOk("match"); ok {
			localData1, err := jsonpath.Get("$[0].query_params", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["type"] = dataLoop1Tmp["type"]
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				dataLoop1Map["name"] = dataLoop1Tmp["name"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			match["queryParams"] = localMaps1
		}

		methods1, _ := jsonpath.Get("$[0].methods", v)
		if methods1 != nil && methods1 != "" {
			match["methods"] = methods1
		}
		ignoreUriCase1, _ := jsonpath.Get("$[0].ignore_uri_case", v)
		if ignoreUriCase1 != nil && ignoreUriCase1 != "" {
			match["ignoreUriCase"] = ignoreUriCase1
		}

		request["match"] = match
	}

	if v, ok := d.GetOk("route_name"); ok {
		request["name"] = v
	}
	backendConfig := make(map[string]interface{})

	if v := d.Get("backend"); !IsNil(v) {
		if v, ok := d.GetOk("backend"); ok {
			localData2, err := jsonpath.Get("$[0].services", v)
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
				dataLoop2Map["version"] = dataLoop2Tmp["version"]
				dataLoop2Map["port"] = dataLoop2Tmp["port"]
				dataLoop2Map["protocol"] = dataLoop2Tmp["protocol"]
				dataLoop2Map["serviceId"] = dataLoop2Tmp["service_id"]
				dataLoop2Map["weight"] = dataLoop2Tmp["weight"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			backendConfig["services"] = localMaps2
		}

		scene1, _ := jsonpath.Get("$[0].scene", v)
		if scene1 != nil && scene1 != "" {
			backendConfig["scene"] = scene1
		}

		request["backendConfig"] = backendConfig
	}

	if v, ok := d.GetOk("environment_info"); ok {
		environmentInfoEnvironmentIdJsonPath, err := jsonpath.Get("$[0].environment_id", v)
		if err == nil && environmentInfoEnvironmentIdJsonPath != "" {
			request["environmentId"] = environmentInfoEnvironmentIdJsonPath
		}
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("domain_ids"); ok {
		domainIdsMapsArray := convertToInterfaceArray(v)

		request["domainIds"] = domainIdsMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_route", action, AlibabaCloudSdkGoERROR)
	}

	datarouteIdVar, _ := jsonpath.Get("$.data.routeId", response)
	d.SetId(fmt.Sprintf("%v:%v", httpApiId, datarouteIdVar))

	return resourceAliCloudApigRouteUpdate(d, meta)
}

func resourceAliCloudApigRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigRoute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_route DescribeApigRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("builtin", objectRaw["builtin"])
	d.Set("create_time", objectRaw["createTimestamp"])
	d.Set("description", objectRaw["description"])
	d.Set("gateway_status", objectRaw["gatewayStatus"])
	d.Set("status", objectRaw["deployStatus"])
	d.Set("update_time", objectRaw["updateTimestamp"])
	d.Set("route_id", objectRaw["routeId"])
	d.Set("route_name", objectRaw["name"])

	backendMaps := make([]map[string]interface{}, 0)
	backendMap := make(map[string]interface{})
	backendRaw := make(map[string]interface{})
	if objectRaw["backend"] != nil {
		backendRaw = objectRaw["backend"].(map[string]interface{})
	}
	if len(backendRaw) > 0 {
		backendMap["scene"] = backendRaw["scene"]

		servicesRaw := backendRaw["services"]
		servicesMaps := make([]map[string]interface{}, 0)
		if servicesRaw != nil {
			for _, servicesChildRaw := range convertToInterfaceArray(servicesRaw) {
				servicesMap := make(map[string]interface{})
				servicesChildRaw := servicesChildRaw.(map[string]interface{})
				servicesMap["name"] = servicesChildRaw["name"]
				servicesMap["port"] = servicesChildRaw["port"]
				servicesMap["protocol"] = servicesChildRaw["protocol"]
				servicesMap["service_id"] = servicesChildRaw["serviceId"]
				servicesMap["version"] = servicesChildRaw["version"]
				servicesMap["weight"] = servicesChildRaw["weight"]

				servicesMaps = append(servicesMaps, servicesMap)
			}
		}
		backendMap["services"] = servicesMaps
		backendMaps = append(backendMaps, backendMap)
	}
	if err := d.Set("backend", backendMaps); err != nil {
		return err
	}
	environmentInfoMaps := make([]map[string]interface{}, 0)
	environmentInfoMap := make(map[string]interface{})
	environmentInfoRaw := make(map[string]interface{})
	if objectRaw["environmentInfo"] != nil {
		environmentInfoRaw = objectRaw["environmentInfo"].(map[string]interface{})
	}
	if len(environmentInfoRaw) > 0 {
		environmentInfoMap["alias"] = environmentInfoRaw["alias"]
		environmentInfoMap["environment_id"] = environmentInfoRaw["environmentId"]
		environmentInfoMap["name"] = environmentInfoRaw["name"]

		gatewayInfoMaps := make([]map[string]interface{}, 0)
		gatewayInfoMap := make(map[string]interface{})
		gatewayInfoRaw := make(map[string]interface{})
		if environmentInfoRaw["gatewayInfo"] != nil {
			gatewayInfoRaw = environmentInfoRaw["gatewayInfo"].(map[string]interface{})
		}
		if len(gatewayInfoRaw) > 0 {
			gatewayInfoMap["gateway_edition"] = gatewayInfoRaw["gatewayEdition"]
			gatewayInfoMap["gateway_id"] = gatewayInfoRaw["gatewayId"]
			gatewayInfoMap["name"] = gatewayInfoRaw["name"]

			gatewayInfoMaps = append(gatewayInfoMaps, gatewayInfoMap)
		}
		environmentInfoMap["gateway_info"] = gatewayInfoMaps
		subDomainsRaw := environmentInfoRaw["subDomains"]
		subDomainsMaps := make([]map[string]interface{}, 0)
		if subDomainsRaw != nil {
			for _, subDomainsChildRaw := range convertToInterfaceArray(subDomainsRaw) {
				subDomainsMap := make(map[string]interface{})
				subDomainsChildRaw := subDomainsChildRaw.(map[string]interface{})
				subDomainsMap["domain_id"] = subDomainsChildRaw["domainId"]
				subDomainsMap["name"] = subDomainsChildRaw["name"]
				subDomainsMap["network_type"] = subDomainsChildRaw["networkType"]
				subDomainsMap["protocol"] = subDomainsChildRaw["protocol"]

				subDomainsMaps = append(subDomainsMaps, subDomainsMap)
			}
		}
		environmentInfoMap["sub_domains"] = subDomainsMaps
		environmentInfoMaps = append(environmentInfoMaps, environmentInfoMap)
	}
	if err := d.Set("environment_info", environmentInfoMaps); err != nil {
		return err
	}
	matchMaps := make([]map[string]interface{}, 0)
	matchMap := make(map[string]interface{})
	matchRaw := make(map[string]interface{})
	if objectRaw["match"] != nil {
		matchRaw = objectRaw["match"].(map[string]interface{})
	}
	if len(matchRaw) > 0 {
		matchMap["ignore_uri_case"] = matchRaw["ignoreUriCase"]

		headersRaw := matchRaw["headers"]
		headersMaps := make([]map[string]interface{}, 0)
		if headersRaw != nil {
			for _, headersChildRaw := range convertToInterfaceArray(headersRaw) {
				headersMap := make(map[string]interface{})
				headersChildRaw := headersChildRaw.(map[string]interface{})
				headersMap["name"] = headersChildRaw["name"]
				headersMap["type"] = headersChildRaw["type"]
				headersMap["value"] = headersChildRaw["value"]

				headersMaps = append(headersMaps, headersMap)
			}
		}
		matchMap["headers"] = headersMaps
		methodsRaw := make([]interface{}, 0)
		if matchRaw["methods"] != nil {
			methodsRaw = convertToInterfaceArray(matchRaw["methods"])
		}

		matchMap["methods"] = methodsRaw
		pathMaps := make([]map[string]interface{}, 0)
		pathMap := make(map[string]interface{})
		pathRaw := make(map[string]interface{})
		if matchRaw["path"] != nil {
			pathRaw = matchRaw["path"].(map[string]interface{})
		}
		if len(pathRaw) > 0 {
			pathMap["type"] = pathRaw["type"]
			pathMap["value"] = pathRaw["value"]

			pathMaps = append(pathMaps, pathMap)
		}
		matchMap["path"] = pathMaps
		queryParamsRaw := matchRaw["queryParams"]
		queryParamsMaps := make([]map[string]interface{}, 0)
		if queryParamsRaw != nil {
			for _, queryParamsChildRaw := range convertToInterfaceArray(queryParamsRaw) {
				queryParamsMap := make(map[string]interface{})
				queryParamsChildRaw := queryParamsChildRaw.(map[string]interface{})
				queryParamsMap["name"] = queryParamsChildRaw["name"]
				queryParamsMap["type"] = queryParamsChildRaw["type"]
				queryParamsMap["value"] = queryParamsChildRaw["value"]

				queryParamsMaps = append(queryParamsMaps, queryParamsMap)
			}
		}
		matchMap["query_params"] = queryParamsMaps
		matchMaps = append(matchMaps, matchMap)
	}
	if err := d.Set("match", matchMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("http_api_id", parts[0])

	return nil
}

func resourceAliCloudApigRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	httpApiId := parts[0]
	routeId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpApiId, routeId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("match") {
		update = true
	}
	match := make(map[string]interface{})

	if v := d.Get("match"); !IsNil(v) || d.HasChange("match") {
		if v, ok := d.GetOk("match"); ok {
			localData, err := jsonpath.Get("$[0].headers", v)
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
				dataLoopMap["value"] = dataLoopTmp["value"]
				dataLoopMap["type"] = dataLoopTmp["type"]
				dataLoopMap["name"] = dataLoopTmp["name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			match["headers"] = localMaps
		}

		if v, ok := d.GetOk("match"); ok {
			localData1, err := jsonpath.Get("$[0].query_params", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				dataLoop1Map["name"] = dataLoop1Tmp["name"]
				dataLoop1Map["type"] = dataLoop1Tmp["type"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			match["queryParams"] = localMaps1
		}

		methods1, _ := jsonpath.Get("$[0].methods", v)
		if methods1 != nil && methods1 != "" {
			match["methods"] = methods1
		}
		path := make(map[string]interface{})
		type5, _ := jsonpath.Get("$[0].path[0].type", d.Get("match"))
		if type5 != nil && type5 != "" {
			path["type"] = type5
		}
		value5, _ := jsonpath.Get("$[0].path[0].value", d.Get("match"))
		if value5 != nil && value5 != "" {
			path["value"] = value5
		}

		if len(path) > 0 {
			match["path"] = path
		}
		ignoreUriCase1, _ := jsonpath.Get("$[0].ignore_uri_case", v)
		if ignoreUriCase1 != nil && ignoreUriCase1 != "" {
			match["ignoreUriCase"] = ignoreUriCase1
		}

		request["match"] = match
	}

	if !d.IsNewResource() && d.HasChange("backend") {
		update = true
	}
	backendConfig := make(map[string]interface{})

	if v := d.Get("backend"); !IsNil(v) || d.HasChange("backend") {
		if v, ok := d.GetOk("backend"); ok {
			localData2, err := jsonpath.Get("$[0].services", v)
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
				dataLoop2Map["port"] = dataLoop2Tmp["port"]
				dataLoop2Map["protocol"] = dataLoop2Tmp["protocol"]
				dataLoop2Map["version"] = dataLoop2Tmp["version"]
				dataLoop2Map["serviceId"] = dataLoop2Tmp["service_id"]
				dataLoop2Map["weight"] = dataLoop2Tmp["weight"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			backendConfig["services"] = localMaps2
		}

		scene1, _ := jsonpath.Get("$[0].scene", v)
		if scene1 != nil && scene1 != "" {
			backendConfig["scene"] = scene1
		}

		request["backendConfig"] = backendConfig
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if !d.IsNewResource() && d.HasChange("domain_ids") {
		update = true
		if v, ok := d.GetOk("domain_ids"); ok || d.HasChange("domain_ids") {
			domainIdsMapsArray := convertToInterfaceArray(v)

			request["domainIds"] = domainIdsMapsArray
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigRouteRead(d, meta)
}

func resourceAliCloudApigRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	httpApiId := parts[0]
	routeId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpApiId, routeId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
