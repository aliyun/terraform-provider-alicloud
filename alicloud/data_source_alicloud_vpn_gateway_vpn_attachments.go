// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"hash/crc32"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudVpnGatewayVpnAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpnGatewayVpnAttachmentRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"init", "active", "attaching", "attached", "detaching", "financialLocked", "provisioning", "updating", "upgrading", "deleted"}, false),
			},
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
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attach_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_asn": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"tunnel_cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_bgp_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effect_immediately": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_dpd": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_nat_traversal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_tunnels_bgp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"dip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"retry": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ike_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ike_auth_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_enc_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_lifetime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"psk": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_pfs": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipsec_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_pfs": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_enc_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_auth_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_lifetime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"local_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"tunnel_options_specification": {
							Set: func(v interface{}) int {
								return int(crc32.ChecksumIEEE([]byte(fmt.Sprint(v.(map[string]interface{})["tunnel_index"]))))
							},
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"customer_gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_no": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_ike_config": {
										Type:     schema.TypeSet,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ike_auth_alg": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"local_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ike_enc_alg": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ike_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ike_mode": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ike_lifetime": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"psk": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"remote_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ike_pfs": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"internet_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_bgp_config": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_asn": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"tunnel_cidr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"bgp_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"peer_bgp_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"peer_asn": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"local_bgp_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_index": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable_nat_traversal": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tunnel_ipsec_config": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipsec_pfs": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ipsec_enc_alg": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ipsec_auth_alg": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ipsec_lifetime": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"enable_dpd": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"vpn_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_connection_id": {
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudVpnGatewayVpnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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

	status, statusOk := d.GetOk("status")
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeVpnConnections"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.VpnConnections.VpnConnection[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VpnConnectionId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != fmt.Sprint(item["State"]) {
				continue
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["VpnConnectionId"]

		mapping["attach_type"] = objectRaw["AttachType"]
		mapping["connection_status"] = objectRaw["State"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["customer_gateway_id"] = objectRaw["CustomerGatewayId"]
		mapping["effect_immediately"] = objectRaw["EffectImmediately"]
		mapping["enable_dpd"] = objectRaw["EnableDpd"]
		mapping["enable_nat_traversal"] = objectRaw["EnableNatTraversal"]
		mapping["enable_tunnels_bgp"] = objectRaw["EnableTunnelsBgp"]
		mapping["internet_ip"] = objectRaw["InternetIp"]
		mapping["local_subnet"] = objectRaw["LocalSubnet"]
		mapping["network_type"] = objectRaw["NetworkType"]
		mapping["remote_subnet"] = objectRaw["RemoteSubnet"]
		mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
		mapping["status"] = objectRaw["State"]
		mapping["vpn_attachment_name"] = objectRaw["Name"]
		mapping["vpn_connection_id"] = objectRaw["VpnConnectionId"]

		bgpConfigMaps := make([]map[string]interface{}, 0)
		bgpConfigMap := make(map[string]interface{})
		vpnBgpConfigRaw := make(map[string]interface{})
		if objectRaw["VpnBgpConfig"] != nil {
			vpnBgpConfigRaw = objectRaw["VpnBgpConfig"].(map[string]interface{})
		}
		if len(vpnBgpConfigRaw) > 0 {
			bgpConfigMap["local_asn"] = vpnBgpConfigRaw["LocalAsn"]
			bgpConfigMap["local_bgp_ip"] = vpnBgpConfigRaw["LocalBgpIp"]
			bgpConfigMap["status"] = vpnBgpConfigRaw["Status"]
			bgpConfigMap["tunnel_cidr"] = vpnBgpConfigRaw["TunnelCidr"]

			bgpConfigMaps = append(bgpConfigMaps, bgpConfigMap)
		}
		mapping["bgp_config"] = bgpConfigMaps
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		vcoHealthCheckRaw := make(map[string]interface{})
		if objectRaw["VcoHealthCheck"] != nil {
			vcoHealthCheckRaw = objectRaw["VcoHealthCheck"].(map[string]interface{})
		}
		if len(vcoHealthCheckRaw) > 0 {
			healthCheckConfigMap["dip"] = vcoHealthCheckRaw["Dip"]
			healthCheckConfigMap["enable"] = formatBool(vcoHealthCheckRaw["Enable"])
			healthCheckConfigMap["interval"] = vcoHealthCheckRaw["Interval"]
			healthCheckConfigMap["policy"] = vcoHealthCheckRaw["Policy"]
			healthCheckConfigMap["retry"] = vcoHealthCheckRaw["Retry"]
			healthCheckConfigMap["sip"] = vcoHealthCheckRaw["Sip"]
			healthCheckConfigMap["status"] = vcoHealthCheckRaw["Status"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		mapping["health_check_config"] = healthCheckConfigMaps
		ikeConfigMaps := make([]map[string]interface{}, 0)
		ikeConfigMap := make(map[string]interface{})
		ikeConfigRaw := make(map[string]interface{})
		if objectRaw["IkeConfig"] != nil {
			ikeConfigRaw = objectRaw["IkeConfig"].(map[string]interface{})
		}
		if len(ikeConfigRaw) > 0 {
			ikeConfigMap["ike_auth_alg"] = ikeConfigRaw["IkeAuthAlg"]
			ikeConfigMap["ike_enc_alg"] = ikeConfigRaw["IkeEncAlg"]
			ikeConfigMap["ike_lifetime"] = ikeConfigRaw["IkeLifetime"]
			ikeConfigMap["ike_mode"] = ikeConfigRaw["IkeMode"]
			ikeConfigMap["ike_pfs"] = ikeConfigRaw["IkePfs"]
			ikeConfigMap["ike_version"] = ikeConfigRaw["IkeVersion"]
			ikeConfigMap["local_id"] = ikeConfigRaw["LocalId"]
			ikeConfigMap["psk"] = ikeConfigRaw["Psk"]
			ikeConfigMap["remote_id"] = ikeConfigRaw["RemoteId"]

			ikeConfigMaps = append(ikeConfigMaps, ikeConfigMap)
		}
		mapping["ike_config"] = ikeConfigMaps
		ipsecConfigMaps := make([]map[string]interface{}, 0)
		ipsecConfigMap := make(map[string]interface{})
		ipsecConfigRaw := make(map[string]interface{})
		if objectRaw["IpsecConfig"] != nil {
			ipsecConfigRaw = objectRaw["IpsecConfig"].(map[string]interface{})
		}
		if len(ipsecConfigRaw) > 0 {
			ipsecConfigMap["ipsec_auth_alg"] = ipsecConfigRaw["IpsecAuthAlg"]
			ipsecConfigMap["ipsec_enc_alg"] = ipsecConfigRaw["IpsecEncAlg"]
			ipsecConfigMap["ipsec_lifetime"] = ipsecConfigRaw["IpsecLifetime"]
			ipsecConfigMap["ipsec_pfs"] = ipsecConfigRaw["IpsecPfs"]

			ipsecConfigMaps = append(ipsecConfigMaps, ipsecConfigMap)
		}
		mapping["ipsec_config"] = ipsecConfigMaps
		tagsMaps, _ := jsonpath.Get("$.Tag.Tag", objectRaw)
		mapping["tags"] = tagsToMap(tagsMaps)
		tunnelOptionsRaw, _ := jsonpath.Get("$.TunnelOptionsSpecification.TunnelOptions", objectRaw)
		tunnelOptionsSpecificationMaps := make([]map[string]interface{}, 0)
		if tunnelOptionsRaw != nil {
			for _, tunnelOptionsChildRaw := range tunnelOptionsRaw.([]interface{}) {
				tunnelOptionsSpecificationMap := make(map[string]interface{})
				tunnelOptionsChildRaw := tunnelOptionsChildRaw.(map[string]interface{})
				tunnelOptionsSpecificationMap["customer_gateway_id"] = tunnelOptionsChildRaw["CustomerGatewayId"]
				tunnelOptionsSpecificationMap["enable_dpd"] = formatBool(tunnelOptionsChildRaw["EnableDpd"])
				tunnelOptionsSpecificationMap["enable_nat_traversal"] = formatBool(tunnelOptionsChildRaw["EnableNatTraversal"])
				tunnelOptionsSpecificationMap["internet_ip"] = tunnelOptionsChildRaw["InternetIp"]
				tunnelOptionsSpecificationMap["role"] = tunnelOptionsChildRaw["Role"]
				tunnelOptionsSpecificationMap["state"] = tunnelOptionsChildRaw["State"]
				tunnelOptionsSpecificationMap["status"] = tunnelOptionsChildRaw["Status"]
				tunnelOptionsSpecificationMap["tunnel_id"] = tunnelOptionsChildRaw["TunnelId"]
				tunnelOptionsSpecificationMap["tunnel_index"] = tunnelOptionsChildRaw["TunnelIndex"]
				tunnelOptionsSpecificationMap["zone_no"] = tunnelOptionsChildRaw["ZoneNo"]

				tunnelBgpConfigMaps := make([]map[string]interface{}, 0)
				tunnelBgpConfigMap := make(map[string]interface{})
				tunnelBgpConfigRawObj, _ := jsonpath.Get("$.TunnelBgpConfig", tunnelOptionsChildRaw)
				tunnelBgpConfigRaw := make(map[string]interface{})
				if tunnelBgpConfigRawObj != nil {
					tunnelBgpConfigRaw = tunnelBgpConfigRawObj.(map[string]interface{})
				}
				if len(tunnelBgpConfigRaw) > 0 {
					tunnelBgpConfigMap["bgp_status"] = tunnelBgpConfigRaw["BgpStatus"]
					tunnelBgpConfigMap["local_asn"] = formatInt(tunnelBgpConfigRaw["LocalAsn"])
					tunnelBgpConfigMap["local_bgp_ip"] = tunnelBgpConfigRaw["LocalBgpIp"]
					tunnelBgpConfigMap["peer_asn"] = tunnelBgpConfigRaw["PeerAsn"]
					tunnelBgpConfigMap["peer_bgp_ip"] = tunnelBgpConfigRaw["PeerBgpIp"]
					tunnelBgpConfigMap["tunnel_cidr"] = tunnelBgpConfigRaw["TunnelCidr"]

					tunnelBgpConfigMaps = append(tunnelBgpConfigMaps, tunnelBgpConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_bgp_config"] = tunnelBgpConfigMaps
				tunnelIkeConfigMaps := make([]map[string]interface{}, 0)
				tunnelIkeConfigMap := make(map[string]interface{})
				tunnelIkeConfigRawObj, _ := jsonpath.Get("$.TunnelIkeConfig", tunnelOptionsChildRaw)
				tunnelIkeConfigRaw := make(map[string]interface{})
				if tunnelIkeConfigRawObj != nil {
					tunnelIkeConfigRaw = tunnelIkeConfigRawObj.(map[string]interface{})
				}
				if len(tunnelIkeConfigRaw) > 0 {
					tunnelIkeConfigMap["ike_auth_alg"] = tunnelIkeConfigRaw["IkeAuthAlg"]
					tunnelIkeConfigMap["ike_enc_alg"] = tunnelIkeConfigRaw["IkeEncAlg"]
					tunnelIkeConfigMap["ike_lifetime"] = formatInt(tunnelIkeConfigRaw["IkeLifetime"])
					tunnelIkeConfigMap["ike_mode"] = tunnelIkeConfigRaw["IkeMode"]
					tunnelIkeConfigMap["ike_pfs"] = tunnelIkeConfigRaw["IkePfs"]
					tunnelIkeConfigMap["ike_version"] = tunnelIkeConfigRaw["IkeVersion"]
					tunnelIkeConfigMap["local_id"] = tunnelIkeConfigRaw["LocalId"]
					tunnelIkeConfigMap["psk"] = tunnelIkeConfigRaw["Psk"]
					tunnelIkeConfigMap["remote_id"] = tunnelIkeConfigRaw["RemoteId"]

					tunnelIkeConfigMaps = append(tunnelIkeConfigMaps, tunnelIkeConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_ike_config"] = tunnelIkeConfigMaps
				tunnelIpsecConfigMaps := make([]map[string]interface{}, 0)
				tunnelIpsecConfigMap := make(map[string]interface{})
				tunnelIpsecConfigRawObj, _ := jsonpath.Get("$.TunnelIpsecConfig", tunnelOptionsChildRaw)
				tunnelIpsecConfigRaw := make(map[string]interface{})
				if tunnelIpsecConfigRawObj != nil {
					tunnelIpsecConfigRaw = tunnelIpsecConfigRawObj.(map[string]interface{})
				}
				if len(tunnelIpsecConfigRaw) > 0 {
					tunnelIpsecConfigMap["ipsec_auth_alg"] = tunnelIpsecConfigRaw["IpsecAuthAlg"]
					tunnelIpsecConfigMap["ipsec_enc_alg"] = tunnelIpsecConfigRaw["IpsecEncAlg"]
					tunnelIpsecConfigMap["ipsec_lifetime"] = formatInt(tunnelIpsecConfigRaw["IpsecLifetime"])
					tunnelIpsecConfigMap["ipsec_pfs"] = tunnelIpsecConfigRaw["IpsecPfs"]

					tunnelIpsecConfigMaps = append(tunnelIpsecConfigMaps, tunnelIpsecConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_ipsec_config"] = tunnelIpsecConfigMaps
				tunnelOptionsSpecificationMaps = append(tunnelOptionsSpecificationMaps, tunnelOptionsSpecificationMap)
			}
		}
		mapping["tunnel_options_specification"] = tunnelOptionsSpecificationMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, mapping["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
