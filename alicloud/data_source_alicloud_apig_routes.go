// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigRouteRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"environment_info": {
				Type:     schema.TypeList,
				Optional: true,
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
			"http_api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"services": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"weight": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"service_id": {
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
									"scene": {
										Type:     schema.TypeString,
										Computed: true,
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
							Computed: true,
						},
						"domain_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
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
						"environment_info": {
							Type:     schema.TypeList,
							Computed: true,
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
										Computed: true,
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
						"match": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
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
									"query_params": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
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
									"methods": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ignore_uri_case": {
										Type:     schema.TypeBool,
										Computed: true,
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
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudApigRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListHttpApiRoutes
	httpApiId := d.Get("http_api_id")
	action := fmt.Sprintf("/v1/http-apis/%s/routes", httpApiId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("environment_info.0.environment_id"); ok {
		query["environmentId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("environment_info.0.gateway_info.0.gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("http_api_id"); ok {
		request["httpApiId"] = v
	}
	if v, ok := d.GetOk("route_name"); ok {
		query["name"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		query["deployStatuses"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(httpApiId, ":", item["routeId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(httpApiId, ":", objectRaw["routeId"])

		mapping["builtin"] = objectRaw["builtin"]
		mapping["create_time"] = objectRaw["createTimestamp"]
		mapping["description"] = objectRaw["description"]
		mapping["gateway_status"] = objectRaw["gatewayStatus"]
		mapping["status"] = objectRaw["deployStatus"]
		mapping["update_time"] = objectRaw["updateTimestamp"]
		mapping["route_id"] = objectRaw["routeId"]
		mapping["route_name"] = objectRaw["name"]

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
		mapping["backend"] = backendMaps
		domainInfosRaw := objectRaw["domainInfos"]
		domainInfosMaps := make([]map[string]interface{}, 0)
		if domainInfosRaw != nil {
			for _, domainInfosChildRaw := range convertToInterfaceArray(domainInfosRaw) {
				domainInfosMap := make(map[string]interface{})
				domainInfosChildRaw := domainInfosChildRaw.(map[string]interface{})
				domainInfosMap["domain_id"] = domainInfosChildRaw["domainId"]
				domainInfosMap["name"] = domainInfosChildRaw["name"]
				domainInfosMap["protocol"] = domainInfosChildRaw["protocol"]

				domainInfosMaps = append(domainInfosMaps, domainInfosMap)
			}
		}
		mapping["domain_infos"] = domainInfosMaps
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
		mapping["environment_info"] = environmentInfoMaps
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
		mapping["match"] = matchMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("routes", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
