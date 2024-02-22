package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effect_immediately": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
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
						"ipsec_config": {
							Type:     schema.TypeList,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_enc_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_auth_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"ike_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"psk": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_enc_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_auth_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ike_local_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_remote_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"vco_health_check": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"retry": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"sip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"vpn_bgp_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"peer_bgp_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"peer_asn": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"local_asn": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"auth_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tunnel_cidr": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"local_bgp_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"enable_tunnels_bgp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tunnel_options_specification": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_ike_config": {
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
									"customer_gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
										ForceNew: true,
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
													Type:     schema.TypeString,
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
									"enable_nat_traversal": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tunnel_id": {
										Type:     schema.TypeString,
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
									"zone_no": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeVpnConnectionsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allVpnConns []vpc.VpnConnection

	if v, ok := d.GetOk("vpn_gateway_id"); ok && v.(string) != "" {
		request.VpnGatewayId = v.(string)
	}

	if v, ok := d.GetOk("customer_gateway_id"); ok && v.(string) != "" {
		request.CustomerGatewayId = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnConnections(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpn_connections", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpnConnectionsResponse)

		if len(response.VpnConnections.VpnConnection) < 1 {
			break
		}

		allVpnConns = append(allVpnConns, response.VpnConnections.VpnConnection...)

		if len(response.VpnConnections.VpnConnection) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredVpnConns []vpc.VpnConnection
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return WrapError(err)
		}
	}

	for _, vpnConn := range allVpnConns {
		if reg != nil {
			if !reg.MatchString(vpnConn.Name) {
				continue
			}
		}
		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if vpnConn.VpnConnectionId == id {
					filteredVpnConns = append(filteredVpnConns, vpnConn)
				}
			}
		} else {
			filteredVpnConns = append(filteredVpnConns, vpnConn)
		}
	}

	return vpnConnectionsDecriptionAttributes(d, filteredVpnConns, meta)
}

func vpnConnectionsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnConnection, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, conn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"customer_gateway_id":  conn.CustomerGatewayId,
			"vpn_gateway_id":       conn.VpnGatewayId,
			"id":                   conn.VpnConnectionId,
			"name":                 conn.Name,
			"local_subnet":         conn.LocalSubnet,
			"remote_subnet":        conn.RemoteSubnet,
			"create_time":          TimestampToStr(conn.CreateTime),
			"effect_immediately":   conn.EffectImmediately,
			"status":               conn.Status,
			"enable_dpd":           conn.EnableDpd,
			"enable_nat_traversal": conn.EnableNatTraversal,
			"enable_tunnels_bgp":   conn.EnableTunnelsBgp,
			"ike_config":           vpnGatewayService.ParseIkeConfig(conn.IkeConfig),
			"ipsec_config":         vpnGatewayService.ParseIpsecConfig(conn.IpsecConfig),
			"vco_health_check":     vpnGatewayService.VcoHealthCheck(conn.VcoHealthCheck),
			"vpn_bgp_config":       vpnGatewayService.VpnBgpConfig(conn.VpnBgpConfig),
		}

		tunnelOptions1Raw, _ := jsonpath.Get("$.TunnelOptionsSpecification.TunnelOptions", conn)
		tunnelOptionsSpecificationMaps := make([]map[string]interface{}, 0)
		if tunnelOptions1Raw != nil {
			for _, tunnelOptionsChild1Raw := range tunnelOptions1Raw.([]interface{}) {
				tunnelOptionsSpecificationMap := make(map[string]interface{})
				tunnelOptionsChild1Raw := tunnelOptionsChild1Raw.(map[string]interface{})
				tunnelOptionsSpecificationMap["customer_gateway_id"] = tunnelOptionsChild1Raw["CustomerGatewayId"]
				tunnelOptionsSpecificationMap["enable_dpd"] = tunnelOptionsChild1Raw["EnableDpd"]
				tunnelOptionsSpecificationMap["enable_nat_traversal"] = tunnelOptionsChild1Raw["EnableNatTraversal"]
				tunnelOptionsSpecificationMap["internet_ip"] = tunnelOptionsChild1Raw["InternetIp"]
				tunnelOptionsSpecificationMap["role"] = tunnelOptionsChild1Raw["Role"]
				tunnelOptionsSpecificationMap["state"] = tunnelOptionsChild1Raw["State"]
				tunnelOptionsSpecificationMap["status"] = tunnelOptionsChild1Raw["Status"]
				tunnelOptionsSpecificationMap["tunnel_id"] = tunnelOptionsChild1Raw["TunnelId"]
				tunnelOptionsSpecificationMap["zone_no"] = tunnelOptionsChild1Raw["ZoneNo"]

				tunnelBgpConfigMaps := make([]map[string]interface{}, 0)
				tunnelBgpConfigMap := make(map[string]interface{})
				tunnelBgpConfig1RawObj, _ := jsonpath.Get("$.TunnelBgpConfig", tunnelOptionsChild1Raw)
				tunnelBgpConfig1Raw := make(map[string]interface{})
				if tunnelBgpConfig1RawObj != nil {
					tunnelBgpConfig1Raw = tunnelBgpConfig1RawObj.(map[string]interface{})
				}
				if len(tunnelBgpConfig1Raw) > 0 {
					tunnelBgpConfigMap["bgp_status"] = tunnelBgpConfig1Raw["BgpStatus"]
					tunnelBgpConfigMap["local_asn"] = tunnelBgpConfig1Raw["LocalAsn"]
					tunnelBgpConfigMap["local_bgp_ip"] = tunnelBgpConfig1Raw["LocalBgpIp"]
					tunnelBgpConfigMap["peer_asn"] = tunnelBgpConfig1Raw["PeerAsn"]
					tunnelBgpConfigMap["peer_bgp_ip"] = tunnelBgpConfig1Raw["PeerBgpIp"]
					tunnelBgpConfigMap["tunnel_cidr"] = tunnelBgpConfig1Raw["TunnelCidr"]

					tunnelBgpConfigMaps = append(tunnelBgpConfigMaps, tunnelBgpConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_bgp_config"] = tunnelBgpConfigMaps
				tunnelIkeConfigMaps := make([]map[string]interface{}, 0)
				tunnelIkeConfigMap := make(map[string]interface{})
				tunnelIkeConfig1RawObj, _ := jsonpath.Get("$.TunnelIkeConfig", tunnelOptionsChild1Raw)
				tunnelIkeConfig1Raw := make(map[string]interface{})
				if tunnelIkeConfig1RawObj != nil {
					tunnelIkeConfig1Raw = tunnelIkeConfig1RawObj.(map[string]interface{})
				}
				if len(tunnelIkeConfig1Raw) > 0 {
					tunnelIkeConfigMap["ike_auth_alg"] = tunnelIkeConfig1Raw["IkeAuthAlg"]
					tunnelIkeConfigMap["ike_enc_alg"] = tunnelIkeConfig1Raw["IkeEncAlg"]
					tunnelIkeConfigMap["ike_lifetime"] = tunnelIkeConfig1Raw["IkeLifetime"]
					tunnelIkeConfigMap["ike_mode"] = tunnelIkeConfig1Raw["IkeMode"]
					tunnelIkeConfigMap["ike_pfs"] = tunnelIkeConfig1Raw["IkePfs"]
					tunnelIkeConfigMap["ike_version"] = tunnelIkeConfig1Raw["IkeVersion"]
					tunnelIkeConfigMap["local_id"] = tunnelIkeConfig1Raw["LocalId"]
					tunnelIkeConfigMap["psk"] = tunnelIkeConfig1Raw["Psk"]
					tunnelIkeConfigMap["remote_id"] = tunnelIkeConfig1Raw["RemoteId"]

					tunnelIkeConfigMaps = append(tunnelIkeConfigMaps, tunnelIkeConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_ike_config"] = tunnelIkeConfigMaps
				tunnelIpsecConfigMaps := make([]map[string]interface{}, 0)
				tunnelIpsecConfigMap := make(map[string]interface{})
				tunnelIpsecConfig1RawObj, _ := jsonpath.Get("$.TunnelIpsecConfig", tunnelOptionsChild1Raw)
				tunnelIpsecConfig1Raw := make(map[string]interface{})
				if tunnelIpsecConfig1RawObj != nil {
					tunnelIpsecConfig1Raw = tunnelIpsecConfig1RawObj.(map[string]interface{})
				}
				if len(tunnelIpsecConfig1Raw) > 0 {
					tunnelIpsecConfigMap["ipsec_auth_alg"] = tunnelIpsecConfig1Raw["IpsecAuthAlg"]
					tunnelIpsecConfigMap["ipsec_enc_alg"] = tunnelIpsecConfig1Raw["IpsecEncAlg"]
					tunnelIpsecConfigMap["ipsec_lifetime"] = tunnelIpsecConfig1Raw["IpsecLifetime"]
					tunnelIpsecConfigMap["ipsec_pfs"] = tunnelIpsecConfig1Raw["IpsecPfs"]

					tunnelIpsecConfigMaps = append(tunnelIpsecConfigMaps, tunnelIpsecConfigMap)
				}
				tunnelOptionsSpecificationMap["tunnel_ipsec_config"] = tunnelIpsecConfigMaps
				tunnelOptionsSpecificationMaps = append(tunnelOptionsSpecificationMaps, tunnelOptionsSpecificationMap)
			}
		}
		mapping["tunnel_options_specification"] = tunnelOptionsSpecificationMaps

		ids = append(ids, conn.VpnConnectionId)
		names = append(names, conn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("connections", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
