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

func dataSourceAliCloudApigGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigGatewayRead,
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
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"environments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"alias": {
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
						"expire_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway_edition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv4_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"address_ip_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_default": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"load_balancer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"address_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
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
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_domain_infos": {
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"target_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vswitch": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vswitch_id": {
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
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
						"zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudApigGatewayRead(d *schema.ResourceData, meta interface{}) error {
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
	// ListGateways
	action := fmt.Sprintf("/v1/gateways")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("gateway_name"); ok {
		query["name"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	tagsFilter := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			tagsFilter[key] = value.(string)
		}
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
				if _, ok := idsMap[fmt.Sprint(item["gatewayId"])]; !ok {
					continue
				}
			}
			if len(tagsFilter) > 0 {
				itemTags := tagsToMap(item["tags"])
				matched := true
				for key, value := range tagsFilter {
					if v, ok := itemTags[key]; !ok || v != value {
						matched = false
						break
					}
				}
				if !matched {
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

		mapping["id"] = fmt.Sprint(objectRaw["gatewayId"])

		mapping["create_from"] = objectRaw["createFrom"]
		mapping["create_time"] = objectRaw["createTimestamp"]
		mapping["expire_time"] = objectRaw["expireTimestamp"]
		mapping["gateway_edition"] = objectRaw["gatewayEdition"]
		mapping["gateway_name"] = objectRaw["name"]
		mapping["gateway_type"] = objectRaw["gatewayType"]
		mapping["payment_type"] = convertApigGatewaydatachargeTypeResponse(objectRaw["chargeType"])
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["spec"] = objectRaw["spec"]
		mapping["status"] = objectRaw["status"]
		mapping["target_version"] = objectRaw["targetVersion"]
		mapping["update_time"] = objectRaw["updateTimestamp"]
		mapping["version"] = objectRaw["version"]
		mapping["gateway_id"] = objectRaw["gatewayId"]

		loadBalancersRaw := objectRaw["loadBalancers"]
		loadBalancersMaps := make([]map[string]interface{}, 0)
		if loadBalancersRaw != nil {
			for _, loadBalancersChildRaw := range convertToInterfaceArray(loadBalancersRaw) {
				loadBalancersMap := make(map[string]interface{})
				loadBalancersChildRaw := loadBalancersChildRaw.(map[string]interface{})
				loadBalancersMap["address"] = loadBalancersChildRaw["address"]
				loadBalancersMap["address_ip_version"] = loadBalancersChildRaw["addressIpVersion"]
				loadBalancersMap["address_type"] = loadBalancersChildRaw["addressType"]
				loadBalancersMap["gateway_default"] = loadBalancersChildRaw["gatewayDefault"]
				loadBalancersMap["load_balancer_id"] = loadBalancersChildRaw["loadBalancerId"]
				loadBalancersMap["mode"] = loadBalancersChildRaw["mode"]
				loadBalancersMap["status"] = loadBalancersChildRaw["status"]
				loadBalancersMap["type"] = loadBalancersChildRaw["type"]

				ipv4AddressesRaw := make([]interface{}, 0)
				if loadBalancersChildRaw["ipv4Addresses"] != nil {
					ipv4AddressesRaw = convertToInterfaceArray(loadBalancersChildRaw["ipv4Addresses"])
				}

				loadBalancersMap["ipv4_addresses"] = ipv4AddressesRaw
				ipv6AddressesRaw := make([]interface{}, 0)
				if loadBalancersChildRaw["ipv6Addresses"] != nil {
					ipv6AddressesRaw = convertToInterfaceArray(loadBalancersChildRaw["ipv6Addresses"])
				}

				loadBalancersMap["ipv6_addresses"] = ipv6AddressesRaw
				portsRaw := loadBalancersChildRaw["ports"]
				portsMaps := make([]map[string]interface{}, 0)
				if portsRaw != nil {
					for _, portsChildRaw := range convertToInterfaceArray(portsRaw) {
						portsMap := make(map[string]interface{})
						portsChildRaw := portsChildRaw.(map[string]interface{})
						portsMap["port"] = portsChildRaw["port"]
						portsMap["protocol"] = portsChildRaw["protocol"]

						portsMaps = append(portsMaps, portsMap)
					}
				}
				loadBalancersMap["ports"] = portsMaps
				loadBalancersMaps = append(loadBalancersMaps, loadBalancersMap)
			}
		}
		mapping["load_balancers"] = loadBalancersMaps
		securityGroupMaps := make([]map[string]interface{}, 0)
		securityGroupMap := make(map[string]interface{})
		securityGroupRaw := make(map[string]interface{})
		if objectRaw["securityGroup"] != nil {
			securityGroupRaw = objectRaw["securityGroup"].(map[string]interface{})
		}
		if len(securityGroupRaw) > 0 {
			securityGroupMap["name"] = securityGroupRaw["name"]
			securityGroupMap["security_group_id"] = securityGroupRaw["securityGroupId"]

			securityGroupMaps = append(securityGroupMaps, securityGroupMap)
		}
		mapping["security_group"] = securityGroupMaps
		subDomainInfosRaw := objectRaw["subDomainInfos"]
		subDomainInfosMaps := make([]map[string]interface{}, 0)
		if subDomainInfosRaw != nil {
			for _, subDomainInfosChildRaw := range convertToInterfaceArray(subDomainInfosRaw) {
				subDomainInfosMap := make(map[string]interface{})
				subDomainInfosChildRaw := subDomainInfosChildRaw.(map[string]interface{})
				subDomainInfosMap["domain_id"] = subDomainInfosChildRaw["domainId"]
				subDomainInfosMap["name"] = subDomainInfosChildRaw["name"]
				subDomainInfosMap["network_type"] = subDomainInfosChildRaw["networkType"]
				subDomainInfosMap["protocol"] = subDomainInfosChildRaw["protocol"]

				subDomainInfosMaps = append(subDomainInfosMaps, subDomainInfosMap)
			}
		}
		mapping["sub_domain_infos"] = subDomainInfosMaps
		tagsMaps := objectRaw["tags"]
		mapping["tags"] = tagsToMap(tagsMaps)
		vpcMaps := make([]map[string]interface{}, 0)
		vpcMap := make(map[string]interface{})
		vpcRaw := make(map[string]interface{})
		if objectRaw["vpc"] != nil {
			vpcRaw = objectRaw["vpc"].(map[string]interface{})
		}
		if len(vpcRaw) > 0 {
			vpcMap["name"] = vpcRaw["name"]
			vpcMap["vpc_id"] = vpcRaw["vpcId"]

			vpcMaps = append(vpcMaps, vpcMap)
		}
		mapping["vpc"] = vpcMaps
		zonesRaw := objectRaw["zones"]
		zonesMaps := make([]map[string]interface{}, 0)
		if zonesRaw != nil {
			for _, zonesChildRaw := range convertToInterfaceArray(zonesRaw) {
				zonesMap := make(map[string]interface{})
				zonesChildRaw := zonesChildRaw.(map[string]interface{})
				zonesMap["name"] = zonesChildRaw["name"]
				zonesMap["zone_id"] = zonesChildRaw["zoneId"]

				vSwitchRawObj, _ := jsonpath.Get("$.vSwitch", zonesChildRaw)
				vSwitchRaw := make(map[string]interface{})
				if vSwitchRawObj != nil {
					vSwitchRaw = vSwitchRawObj.(map[string]interface{})
				}
				zonesMap["vswitch_id"] = vSwitchRaw["vSwitchId"]

				zonesMaps = append(zonesMaps, zonesMap)
			}
		}
		mapping["zones"] = zonesMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["gatewayId"])
		mapping, err = dataSourceAliCloudApigGatewayReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

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
	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudApigGatewayReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	apigServiceV2 := ApigServiceV2{client}
	getResp, err := apigServiceV2.DescribeApigGateway(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_from"] = objectRaw["createFrom"]
	mapping["create_time"] = objectRaw["createTimestamp"]
	mapping["expire_time"] = objectRaw["expireTimestamp"]
	mapping["gateway_edition"] = objectRaw["gatewayEdition"]
	mapping["gateway_name"] = objectRaw["name"]
	mapping["gateway_type"] = objectRaw["gatewayType"]
	mapping["payment_type"] = convertApigGatewaydatachargeTypeResponse(objectRaw["chargeType"])
	mapping["resource_group_id"] = objectRaw["resourceGroupId"]
	mapping["spec"] = objectRaw["spec"]
	mapping["status"] = objectRaw["status"]
	mapping["target_version"] = objectRaw["targetVersion"]
	mapping["update_time"] = objectRaw["updateTimestamp"]
	mapping["version"] = objectRaw["version"]
	mapping["gateway_id"] = objectRaw["gatewayId"]

	environmentsRaw := objectRaw["environments"]
	environmentsMaps := make([]map[string]interface{}, 0)
	if environmentsRaw != nil {
		for _, environmentsChildRaw := range convertToInterfaceArray(environmentsRaw) {
			environmentsMap := make(map[string]interface{})
			environmentsChildRaw := environmentsChildRaw.(map[string]interface{})
			environmentsMap["alias"] = environmentsChildRaw["alias"]
			environmentsMap["environment_id"] = environmentsChildRaw["environmentId"]
			environmentsMap["name"] = environmentsChildRaw["name"]

			environmentsMaps = append(environmentsMaps, environmentsMap)
		}
	}
	mapping["environments"] = environmentsMaps
	loadBalancersRaw := objectRaw["loadBalancers"]
	loadBalancersMaps := make([]map[string]interface{}, 0)
	if loadBalancersRaw != nil {
		for _, loadBalancersChildRaw := range convertToInterfaceArray(loadBalancersRaw) {
			loadBalancersMap := make(map[string]interface{})
			loadBalancersChildRaw := loadBalancersChildRaw.(map[string]interface{})
			loadBalancersMap["address"] = loadBalancersChildRaw["address"]
			loadBalancersMap["address_ip_version"] = loadBalancersChildRaw["addressIpVersion"]
			loadBalancersMap["address_type"] = loadBalancersChildRaw["addressType"]
			loadBalancersMap["gateway_default"] = loadBalancersChildRaw["gatewayDefault"]
			loadBalancersMap["load_balancer_id"] = loadBalancersChildRaw["loadBalancerId"]
			loadBalancersMap["mode"] = loadBalancersChildRaw["mode"]
			loadBalancersMap["status"] = loadBalancersChildRaw["status"]
			loadBalancersMap["type"] = loadBalancersChildRaw["type"]

			ipv4AddressesRaw := make([]interface{}, 0)
			if loadBalancersChildRaw["ipv4Addresses"] != nil {
				ipv4AddressesRaw = convertToInterfaceArray(loadBalancersChildRaw["ipv4Addresses"])
			}

			loadBalancersMap["ipv4_addresses"] = ipv4AddressesRaw
			ipv6AddressesRaw := make([]interface{}, 0)
			if loadBalancersChildRaw["ipv6Addresses"] != nil {
				ipv6AddressesRaw = convertToInterfaceArray(loadBalancersChildRaw["ipv6Addresses"])
			}

			loadBalancersMap["ipv6_addresses"] = ipv6AddressesRaw
			portsRaw := loadBalancersChildRaw["ports"]
			portsMaps := make([]map[string]interface{}, 0)
			if portsRaw != nil {
				for _, portsChildRaw := range convertToInterfaceArray(portsRaw) {
					portsMap := make(map[string]interface{})
					portsChildRaw := portsChildRaw.(map[string]interface{})
					portsMap["port"] = portsChildRaw["port"]
					portsMap["protocol"] = portsChildRaw["protocol"]

					portsMaps = append(portsMaps, portsMap)
				}
			}
			loadBalancersMap["ports"] = portsMaps
			loadBalancersMaps = append(loadBalancersMaps, loadBalancersMap)
		}
	}
	mapping["load_balancers"] = loadBalancersMaps
	securityGroupMaps := make([]map[string]interface{}, 0)
	securityGroupMap := make(map[string]interface{})
	securityGroupRaw := make(map[string]interface{})
	if objectRaw["securityGroup"] != nil {
		securityGroupRaw = objectRaw["securityGroup"].(map[string]interface{})
	}
	if len(securityGroupRaw) > 0 {
		securityGroupMap["name"] = securityGroupRaw["name"]
		securityGroupMap["security_group_id"] = securityGroupRaw["securityGroupId"]

		securityGroupMaps = append(securityGroupMaps, securityGroupMap)
	}
	mapping["security_group"] = securityGroupMaps
	tagsMaps := objectRaw["tags"]
	mapping["tags"] = tagsToMap(tagsMaps)
	vSwitchMaps := make([]map[string]interface{}, 0)
	vSwitchMap := make(map[string]interface{})
	vSwitchRaw := make(map[string]interface{})
	if objectRaw["vSwitch"] != nil {
		vSwitchRaw = objectRaw["vSwitch"].(map[string]interface{})
	}
	if len(vSwitchRaw) > 0 {
		vSwitchMap["name"] = vSwitchRaw["name"]
		vSwitchMap["vswitch_id"] = vSwitchRaw["vSwitchId"]

		vSwitchMaps = append(vSwitchMaps, vSwitchMap)
	}
	mapping["vswitch"] = vSwitchMaps
	vpcMaps := make([]map[string]interface{}, 0)
	vpcMap := make(map[string]interface{})
	vpcRaw := make(map[string]interface{})
	if objectRaw["vpc"] != nil {
		vpcRaw = objectRaw["vpc"].(map[string]interface{})
	}
	if len(vpcRaw) > 0 {
		vpcMap["name"] = vpcRaw["name"]
		vpcMap["vpc_id"] = vpcRaw["vpcId"]

		vpcMaps = append(vpcMaps, vpcMap)
	}
	mapping["vpc"] = vpcMaps
	zonesRaw := objectRaw["zones"]
	zonesMaps := make([]map[string]interface{}, 0)
	if zonesRaw != nil {
		for _, zonesChildRaw := range convertToInterfaceArray(zonesRaw) {
			zonesMap := make(map[string]interface{})
			zonesChildRaw := zonesChildRaw.(map[string]interface{})
			zonesMap["name"] = zonesChildRaw["name"]
			zonesMap["zone_id"] = zonesChildRaw["zoneId"]

			vSwitchRawObj, _ := jsonpath.Get("$.vSwitch", zonesChildRaw)
			vSwitchRaw := make(map[string]interface{})
			if vSwitchRawObj != nil {
				vSwitchRaw = vSwitchRawObj.(map[string]interface{})
			}
			zonesMap["vswitch_id"] = vSwitchRaw["vSwitchId"]

			zonesMaps = append(zonesMaps, zonesMap)
		}
	}
	mapping["zones"] = zonesMaps

	return mapping, nil
}
