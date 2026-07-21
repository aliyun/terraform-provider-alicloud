// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVPNGatewayVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVPNGatewayVpnConnectionCreate,
		Read:   resourceAliCloudVPNGatewayVpnConnectionRead,
		Update: resourceAliCloudVPNGatewayVpnConnectionUpdate,
		Delete: resourceAliCloudVPNGatewayVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_config_route": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bgp_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_asn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tunnel_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"local_bgp_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_dpd": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_nat_traversal": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_tunnels_bgp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"dip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"retry": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ike_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_auth_alg": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_enc_alg": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_local_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_remote_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_lifetime": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"psk": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_pfs": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ipsec_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_pfs": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipsec_enc_alg": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipsec_auth_alg": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipsec_lifetime": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"local_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"tunnel_options_specification": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_ike_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ike_auth_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"local_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ike_enc_alg": {
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
									"ike_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"psk": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"remote_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ike_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_bgp_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_asn": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tunnel_cidr": {
										Type:     schema.TypeString,
										Optional: true,
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
										Optional: true,
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
							Optional: true,
						},
						"tunnel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_ipsec_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_enc_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_auth_alg": {
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
						"enable_dpd": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"zone_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpn_connection_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.216.0. New field 'vpn_connection_name' instead.",
			},
		},
	}
}

func resourceAliCloudVPNGatewayVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpnConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("local_subnet"); ok {
		request["LocalSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("remote_subnet"); ok {
		request["RemoteSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}
	if v, ok := d.GetOkExists("effect_immediately"); ok {
		request["EffectImmediately"] = v
	}
	if v, ok := d.GetOkExists("enable_dpd"); ok {
		request["EnableDpd"] = v
	}
	if v, ok := d.GetOkExists("enable_nat_traversal"); ok {
		request["EnableNatTraversal"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("ike_config"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].psk", d.Get("ike_config"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["Psk"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].ike_version", d.Get("ike_config"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["IkeVersion"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].ike_mode", d.Get("ike_config"))
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["IkeMode"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].ike_enc_alg", d.Get("ike_config"))
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["IkeEncAlg"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].ike_auth_alg", d.Get("ike_config"))
		if nodeNative4 != nil && nodeNative4 != "" {
			objectDataLocalMap["IkeAuthAlg"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].ike_pfs", d.Get("ike_config"))
		if nodeNative5 != nil && nodeNative5 != "" {
			objectDataLocalMap["IkePfs"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].ike_lifetime", d.Get("ike_config"))
		if nodeNative6 != nil && nodeNative6 != "" {
			objectDataLocalMap["IkeLifetime"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].ike_local_id", d.Get("ike_config"))
		if nodeNative7 != nil && nodeNative7 != "" {
			objectDataLocalMap["LocalId"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].ike_remote_id", d.Get("ike_config"))
		if nodeNative8 != nil && nodeNative8 != "" {
			objectDataLocalMap["RemoteId"] = nodeNative8
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["IkeConfig"] = string(objectDataLocalMapJson)
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("ipsec_config"); !IsNil(v) {
		nodeNative9, _ := jsonpath.Get("$[0].ipsec_enc_alg", d.Get("ipsec_config"))
		if nodeNative9 != nil && nodeNative9 != "" {
			objectDataLocalMap1["IpsecEncAlg"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].ipsec_auth_alg", d.Get("ipsec_config"))
		if nodeNative10 != nil && nodeNative10 != "" {
			objectDataLocalMap1["IpsecAuthAlg"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].ipsec_pfs", d.Get("ipsec_config"))
		if nodeNative11 != nil && nodeNative11 != "" {
			objectDataLocalMap1["IpsecPfs"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].ipsec_lifetime", d.Get("ipsec_config"))
		if nodeNative12 != nil && nodeNative12 != "" {
			objectDataLocalMap1["IpsecLifetime"] = nodeNative12
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		request["IpsecConfig"] = string(objectDataLocalMap1Json)
	}

	objectDataLocalMap2 := make(map[string]interface{})
	if v := d.Get("bgp_config"); !IsNil(v) {
		nodeNative13, _ := jsonpath.Get("$[0].local_asn", d.Get("bgp_config"))
		if nodeNative13 != nil && nodeNative13 != "" {
			objectDataLocalMap2["LocalAsn"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].tunnel_cidr", d.Get("bgp_config"))
		if nodeNative14 != nil && nodeNative14 != "" {
			objectDataLocalMap2["TunnelCidr"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].local_bgp_ip", d.Get("bgp_config"))
		if nodeNative15 != nil && nodeNative15 != "" {
			objectDataLocalMap2["LocalBgpIp"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].enable", d.Get("bgp_config"))
		if nodeNative16 != nil && nodeNative16 != "" {
			objectDataLocalMap2["EnableBgp"] = nodeNative16
		}

		objectDataLocalMap2Json, err := json.Marshal(objectDataLocalMap2)
		if err != nil {
			return WrapError(err)
		}
		request["BgpConfig"] = string(objectDataLocalMap2Json)
	}

	if v, ok := d.GetOk("customer_gateway_id"); ok {
		request["CustomerGatewayId"] = v
	}
	if v, ok := d.GetOkExists("auto_config_route"); ok {
		request["AutoConfigRoute"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("vpn_connection_name"); ok {
		request["Name"] = v
	}
	request["VpnGatewayId"] = d.Get("vpn_gateway_id")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("tunnel_options_specification"); ok {
		tunnelOptionsSpecificationMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CustomerGatewayId"] = dataLoop1Tmp["customer_gateway_id"]
			dataLoop1Map["EnableDpd"] = dataLoop1Tmp["enable_dpd"]
			dataLoop1Map["EnableNatTraversal"] = dataLoop1Tmp["enable_nat_traversal"]
			dataLoop1Map["Role"] = dataLoop1Tmp["role"]
			if !IsNil(dataLoop1Tmp["tunnel_bgp_config"]) {
				localData2 := make(map[string]interface{})
				nodeNative23, _ := jsonpath.Get("$.tunnel_bgp_config[0].local_asn", dataLoop1Tmp)
				if nodeNative23 != nil && nodeNative23 != "" {
					localData2["LocalAsn"] = nodeNative23
				}
				nodeNative24, _ := jsonpath.Get("$.tunnel_bgp_config[0].local_bgp_ip", dataLoop1Tmp)
				if nodeNative24 != nil && nodeNative24 != "" {
					localData2["LocalBgpIp"] = nodeNative24
				}
				nodeNative25, _ := jsonpath.Get("$.tunnel_bgp_config[0].tunnel_cidr", dataLoop1Tmp)
				if nodeNative25 != nil && nodeNative25 != "" {
					localData2["TunnelCidr"] = nodeNative25
				}
				dataLoop1Map["TunnelBgpConfig"] = localData2
			}
			if !IsNil(dataLoop1Tmp["tunnel_ike_config"]) {
				localData3 := make(map[string]interface{})
				nodeNative26, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_auth_alg", dataLoop1Tmp)
				if nodeNative26 != nil && nodeNative26 != "" {
					localData3["IkeAuthAlg"] = nodeNative26
				}
				nodeNative27, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_enc_alg", dataLoop1Tmp)
				if nodeNative27 != nil && nodeNative27 != "" {
					localData3["IkeEncAlg"] = nodeNative27
				}
				nodeNative28, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_lifetime", dataLoop1Tmp)
				if nodeNative28 != nil && nodeNative28 != "" {
					localData3["IkeLifetime"] = nodeNative28
				}
				nodeNative29, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_mode", dataLoop1Tmp)
				if nodeNative29 != nil && nodeNative29 != "" {
					localData3["IkeMode"] = nodeNative29
				}
				nodeNative30, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_pfs", dataLoop1Tmp)
				if nodeNative30 != nil && nodeNative30 != "" {
					localData3["IkePfs"] = nodeNative30
				}
				nodeNative31, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_version", dataLoop1Tmp)
				if nodeNative31 != nil && nodeNative31 != "" {
					localData3["IkeVersion"] = nodeNative31
				}
				nodeNative32, _ := jsonpath.Get("$.tunnel_ike_config[0].local_id", dataLoop1Tmp)
				if nodeNative32 != nil && nodeNative32 != "" {
					localData3["LocalId"] = nodeNative32
				}
				nodeNative33, _ := jsonpath.Get("$.tunnel_ike_config[0].psk", dataLoop1Tmp)
				if nodeNative33 != nil && nodeNative33 != "" {
					localData3["Psk"] = nodeNative33
				}
				nodeNative34, _ := jsonpath.Get("$.tunnel_ike_config[0].remote_id", dataLoop1Tmp)
				if nodeNative34 != nil && nodeNative34 != "" {
					localData3["RemoteId"] = nodeNative34
				}
				dataLoop1Map["TunnelIkeConfig"] = localData3
			}
			if !IsNil(dataLoop1Tmp["tunnel_ipsec_config"]) {
				localData4 := make(map[string]interface{})
				nodeNative35, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_auth_alg", dataLoop1Tmp)
				if nodeNative35 != nil && nodeNative35 != "" {
					localData4["IpsecAuthAlg"] = nodeNative35
				}
				nodeNative36, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_enc_alg", dataLoop1Tmp)
				if nodeNative36 != nil && nodeNative36 != "" {
					localData4["IpsecEncAlg"] = nodeNative36
				}
				nodeNative37, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_lifetime", dataLoop1Tmp)
				if nodeNative37 != nil && nodeNative37 != "" {
					localData4["IpsecLifetime"] = nodeNative37
				}
				nodeNative38, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_pfs", dataLoop1Tmp)
				if nodeNative38 != nil && nodeNative38 != "" {
					localData4["IpsecPfs"] = nodeNative38
				}
				dataLoop1Map["TunnelIpsecConfig"] = localData4
			}
			tunnelOptionsSpecificationMaps = append(tunnelOptionsSpecificationMaps, dataLoop1Map)
		}
		request["TunnelOptionsSpecification"] = tunnelOptionsSpecificationMaps
	}

	if v, ok := d.GetOkExists("enable_tunnels_bgp"); ok {
		request["EnableTunnelsBgp"] = v
	}
	objectDataLocalMap3 := make(map[string]interface{})
	if v := d.Get("health_check_config"); !IsNil(v) {
		nodeNative39, _ := jsonpath.Get("$[0].enable", d.Get("health_check_config"))
		if nodeNative39 != nil && nodeNative39 != "" {
			objectDataLocalMap3["enable"] = nodeNative39
		}
		nodeNative40, _ := jsonpath.Get("$[0].dip", d.Get("health_check_config"))
		if nodeNative40 != nil && nodeNative40 != "" {
			objectDataLocalMap3["dip"] = nodeNative40
		}
		nodeNative41, _ := jsonpath.Get("$[0].sip", d.Get("health_check_config"))
		if nodeNative41 != nil && nodeNative41 != "" {
			objectDataLocalMap3["sip"] = nodeNative41
		}
		nodeNative42, _ := jsonpath.Get("$[0].interval", d.Get("health_check_config"))
		if nodeNative42 != nil && nodeNative42 != "" {
			objectDataLocalMap3["interval"] = nodeNative42
		}
		nodeNative43, _ := jsonpath.Get("$[0].retry", d.Get("health_check_config"))
		if nodeNative43 != nil && nodeNative43 != "" {
			objectDataLocalMap3["retry"] = nodeNative43
		}

		objectDataLocalMap3Json, err := json.Marshal(objectDataLocalMap3)
		if err != nil {
			return WrapError(err)
		}
		request["HealthCheckConfig"] = string(objectDataLocalMap3Json)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Appliance.Configuring", "VpnGateway.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpnConnectionId"]))

	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, vPNGatewayServiceV2.VPNGatewayVpnConnectionStateRefreshFunc(d.Id(), "State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVPNGatewayVpnConnectionUpdate(d, meta)
}

func resourceAliCloudVPNGatewayVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}

	objectRaw, err := vPNGatewayServiceV2.DescribeVPNGatewayVpnConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_connection DescribeVPNGatewayVpnConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("customer_gateway_id", objectRaw["CustomerGatewayId"])
	d.Set("effect_immediately", objectRaw["EffectImmediately"])
	d.Set("enable_dpd", objectRaw["EnableDpd"])
	d.Set("enable_nat_traversal", objectRaw["EnableNatTraversal"])
	d.Set("enable_tunnels_bgp", objectRaw["EnableTunnelsBgp"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpn_connection_name", objectRaw["Name"])
	d.Set("vpn_gateway_id", objectRaw["VpnGatewayId"])

	bgpConfigMaps := make([]map[string]interface{}, 0)
	bgpConfigMap := make(map[string]interface{})
	vpnBgpConfig1Raw := make(map[string]interface{})
	if objectRaw["VpnBgpConfig"] != nil {
		vpnBgpConfig1Raw = objectRaw["VpnBgpConfig"].(map[string]interface{})
	}
	if len(vpnBgpConfig1Raw) > 0 {
		bgpConfigMap["enable"] = formatBool(vpnBgpConfig1Raw["EnableBgp"])
		bgpConfigMap["local_asn"] = vpnBgpConfig1Raw["LocalAsn"]
		bgpConfigMap["local_bgp_ip"] = vpnBgpConfig1Raw["LocalBgpIp"]
		bgpConfigMap["status"] = vpnBgpConfig1Raw["Status"]
		bgpConfigMap["tunnel_cidr"] = vpnBgpConfig1Raw["TunnelCidr"]

		bgpConfigMaps = append(bgpConfigMaps, bgpConfigMap)
	}
	d.Set("bgp_config", bgpConfigMaps)
	healthCheckConfigMaps := make([]map[string]interface{}, 0)
	healthCheckConfigMap := make(map[string]interface{})
	vcoHealthCheck1Raw := make(map[string]interface{})
	if objectRaw["VcoHealthCheck"] != nil {
		vcoHealthCheck1Raw = objectRaw["VcoHealthCheck"].(map[string]interface{})
	}
	if len(vcoHealthCheck1Raw) > 0 {
		healthCheckConfigMap["dip"] = vcoHealthCheck1Raw["Dip"]
		healthCheckConfigMap["enable"] = formatBool(vcoHealthCheck1Raw["Enable"])
		healthCheckConfigMap["interval"] = vcoHealthCheck1Raw["Interval"]
		healthCheckConfigMap["retry"] = vcoHealthCheck1Raw["Retry"]
		healthCheckConfigMap["sip"] = vcoHealthCheck1Raw["Sip"]

		healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
	}
	d.Set("health_check_config", healthCheckConfigMaps)
	ikeConfigMaps := make([]map[string]interface{}, 0)
	ikeConfigMap := make(map[string]interface{})
	ikeConfig1Raw := make(map[string]interface{})
	if objectRaw["IkeConfig"] != nil {
		ikeConfig1Raw = objectRaw["IkeConfig"].(map[string]interface{})
	}
	if len(ikeConfig1Raw) > 0 {
		ikeConfigMap["ike_auth_alg"] = ikeConfig1Raw["IkeAuthAlg"]
		ikeConfigMap["ike_enc_alg"] = ikeConfig1Raw["IkeEncAlg"]
		ikeConfigMap["ike_lifetime"] = ikeConfig1Raw["IkeLifetime"]
		ikeConfigMap["ike_local_id"] = ikeConfig1Raw["LocalId"]
		ikeConfigMap["ike_mode"] = ikeConfig1Raw["IkeMode"]
		ikeConfigMap["ike_pfs"] = ikeConfig1Raw["IkePfs"]
		ikeConfigMap["ike_remote_id"] = ikeConfig1Raw["RemoteId"]
		ikeConfigMap["ike_version"] = ikeConfig1Raw["IkeVersion"]
		ikeConfigMap["psk"] = ikeConfig1Raw["Psk"]

		ikeConfigMaps = append(ikeConfigMaps, ikeConfigMap)
	}
	d.Set("ike_config", ikeConfigMaps)
	ipsecConfigMaps := make([]map[string]interface{}, 0)
	ipsecConfigMap := make(map[string]interface{})
	ipsecConfig1Raw := make(map[string]interface{})
	if objectRaw["IpsecConfig"] != nil {
		ipsecConfig1Raw = objectRaw["IpsecConfig"].(map[string]interface{})
	}
	if len(ipsecConfig1Raw) > 0 {
		ipsecConfigMap["ipsec_auth_alg"] = ipsecConfig1Raw["IpsecAuthAlg"]
		ipsecConfigMap["ipsec_enc_alg"] = ipsecConfig1Raw["IpsecEncAlg"]
		ipsecConfigMap["ipsec_lifetime"] = ipsecConfig1Raw["IpsecLifetime"]
		ipsecConfigMap["ipsec_pfs"] = ipsecConfig1Raw["IpsecPfs"]

		ipsecConfigMaps = append(ipsecConfigMaps, ipsecConfigMap)
	}
	d.Set("ipsec_config", ipsecConfigMaps)
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	tunnelOptions1Raw, _ := jsonpath.Get("$.TunnelOptionsSpecification.TunnelOptions", objectRaw)
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
	d.Set("tunnel_options_specification", tunnelOptionsSpecificationMaps)

	e := jsonata.MustCompile("$split($.LocalSubnet, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("local_subnet", evaluation)
	e = jsonata.MustCompile("$split($.RemoteSubnet, \",\")")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("remote_subnet", evaluation)

	d.Set("name", d.Get("vpn_connection_name"))
	return nil
}

func resourceAliCloudVPNGatewayVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyVpnConnectionAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["VpnConnectionId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("local_subnet") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("remote_subnet") {
		update = true
	}
	if v, ok := d.GetOk("local_subnet"); ok {
		request["LocalSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("remote_subnet"); ok {
		request["RemoteSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("effect_immediately") {
		update = true
		request["EffectImmediately"] = d.Get("effect_immediately")
	}

	if v, ok := d.GetOkExists("auto_config_route"); ok {
		request["AutoConfigRoute"] = v
	}
	if !d.IsNewResource() && d.HasChange("enable_dpd") {
		update = true
		request["EnableDpd"] = d.Get("enable_dpd")
	}

	if !d.IsNewResource() && d.HasChange("enable_nat_traversal") {
		update = true
		request["EnableNatTraversal"] = d.Get("enable_nat_traversal")
	}

	if !d.IsNewResource() && d.HasChange("enable_tunnels_bgp") {
		request["EnableTunnelsBgp"] = d.Get("enable_tunnels_bgp")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("vpn_connection_name") {
		update = true
		request["Name"] = d.Get("vpn_connection_name")
	}

	if !d.IsNewResource() && d.HasChange("ike_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("ike_config"); v != nil {
			nodeNative, _ := jsonpath.Get("$[0].ike_lifetime", v)
			if nodeNative != nil && nodeNative != "" {
				objectDataLocalMap["IkeLifetime"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].ike_enc_alg", v)
			if nodeNative1 != nil && nodeNative1 != "" {
				objectDataLocalMap["IkeEncAlg"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].ike_mode", v)
			if nodeNative2 != nil && nodeNative2 != "" {
				objectDataLocalMap["IkeMode"] = nodeNative2
			}
			nodeNative3, _ := jsonpath.Get("$[0].ike_version", v)
			if nodeNative3 != nil && nodeNative3 != "" {
				objectDataLocalMap["IkeVersion"] = nodeNative3
			}
			nodeNative4, _ := jsonpath.Get("$[0].ike_pfs", v)
			if nodeNative4 != nil && nodeNative4 != "" {
				objectDataLocalMap["IkePfs"] = nodeNative4
			}
			nodeNative5, _ := jsonpath.Get("$[0].psk", v)
			if nodeNative5 != nil && nodeNative5 != "" {
				objectDataLocalMap["Psk"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].ike_auth_alg", v)
			if nodeNative6 != nil && nodeNative6 != "" {
				objectDataLocalMap["IkeAuthAlg"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].ike_local_id", v)
			if nodeNative7 != nil && nodeNative7 != "" {
				objectDataLocalMap["LocalId"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].ike_remote_id", v)
			if nodeNative8 != nil && nodeNative8 != "" {
				objectDataLocalMap["RemoteId"] = nodeNative8
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["IkeConfig"] = string(objectDataLocalMapJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("ipsec_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})
		if v := d.Get("ipsec_config"); v != nil {
			nodeNative9, _ := jsonpath.Get("$[0].ipsec_auth_alg", v)
			if nodeNative9 != nil && nodeNative9 != "" {
				objectDataLocalMap1["IpsecAuthAlg"] = nodeNative9
			}
			nodeNative10, _ := jsonpath.Get("$[0].ipsec_lifetime", v)
			if nodeNative10 != nil && nodeNative10 != "" {
				objectDataLocalMap1["IpsecLifetime"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].ipsec_enc_alg", v)
			if nodeNative11 != nil && nodeNative11 != "" {
				objectDataLocalMap1["IpsecEncAlg"] = nodeNative11
			}
			nodeNative12, _ := jsonpath.Get("$[0].ipsec_pfs", v)
			if nodeNative12 != nil && nodeNative12 != "" {
				objectDataLocalMap1["IpsecPfs"] = nodeNative12
			}

			objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
			if err != nil {
				return WrapError(err)
			}
			request["IpsecConfig"] = string(objectDataLocalMap1Json)
		}
	}

	if !d.IsNewResource() && d.HasChange("bgp_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})
		if v := d.Get("bgp_config"); v != nil {
			nodeNative13, _ := jsonpath.Get("$[0].tunnel_cidr", v)
			if nodeNative13 != nil && nodeNative13 != "" {
				objectDataLocalMap2["TunnelCidr"] = nodeNative13
			}
			nodeNative14, _ := jsonpath.Get("$[0].enable", v)
			if nodeNative14 != nil && nodeNative14 != "" {
				objectDataLocalMap2["EnableBgp"] = nodeNative14
			}
			nodeNative15, _ := jsonpath.Get("$[0].local_bgp_ip", v)
			if nodeNative15 != nil && nodeNative15 != "" {
				objectDataLocalMap2["LocalBgpIp"] = nodeNative15
			}
			nodeNative16, _ := jsonpath.Get("$[0].local_asn", v)
			if nodeNative16 != nil && nodeNative16 != "" {
				objectDataLocalMap2["LocalAsn"] = nodeNative16
			}

			objectDataLocalMap2Json, err := json.Marshal(objectDataLocalMap2)
			if err != nil {
				return WrapError(err)
			}
			request["BgpConfig"] = string(objectDataLocalMap2Json)
		}
	}

	if !d.IsNewResource() && d.HasChange("tunnel_options_specification") {
		update = true
		if v, ok := d.GetOk("tunnel_options_specification"); ok {
			tunnelOptionsSpecificationMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["EnableDpd"] = dataLoopTmp["enable_dpd"]
				dataLoopMap["EnableNatTraversal"] = dataLoopTmp["enable_nat_traversal"]
				if !IsNil(dataLoopTmp["tunnel_bgp_config"]) {
					localData1 := make(map[string]interface{})
					nodeNative19, _ := jsonpath.Get("$.tunnel_bgp_config[0].local_asn", dataLoopTmp)
					if nodeNative19 != nil && nodeNative19 != "" {
						localData1["LocalAsn"] = nodeNative19
					}
					nodeNative20, _ := jsonpath.Get("$.tunnel_bgp_config[0].local_bgp_ip", dataLoopTmp)
					if nodeNative20 != nil && nodeNative20 != "" {
						localData1["LocalBgpIp"] = nodeNative20
					}
					nodeNative21, _ := jsonpath.Get("$.tunnel_bgp_config[0].tunnel_cidr", dataLoopTmp)
					if nodeNative21 != nil && nodeNative21 != "" {
						localData1["TunnelCidr"] = nodeNative21
					}
					dataLoopMap["TunnelBgpConfig"] = localData1
				}
				if !IsNil(dataLoopTmp["tunnel_ike_config"]) {
					localData2 := make(map[string]interface{})
					nodeNative22, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_auth_alg", dataLoopTmp)
					if nodeNative22 != nil && nodeNative22 != "" {
						localData2["IkeAuthAlg"] = nodeNative22
					}
					nodeNative23, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_enc_alg", dataLoopTmp)
					if nodeNative23 != nil && nodeNative23 != "" {
						localData2["IkeEncAlg"] = nodeNative23
					}
					nodeNative24, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_lifetime", dataLoopTmp)
					if nodeNative24 != nil && nodeNative24 != "" {
						localData2["IkeLifetime"] = nodeNative24
					}
					nodeNative25, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_mode", dataLoopTmp)
					if nodeNative25 != nil && nodeNative25 != "" {
						localData2["IkeMode"] = nodeNative25
					}
					nodeNative26, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_pfs", dataLoopTmp)
					if nodeNative26 != nil && nodeNative26 != "" {
						localData2["IkePfs"] = nodeNative26
					}
					nodeNative27, _ := jsonpath.Get("$.tunnel_ike_config[0].ike_version", dataLoopTmp)
					if nodeNative27 != nil && nodeNative27 != "" {
						localData2["IkeVersion"] = nodeNative27
					}
					nodeNative28, _ := jsonpath.Get("$.tunnel_ike_config[0].local_id", dataLoopTmp)
					if nodeNative28 != nil && nodeNative28 != "" {
						localData2["LocalId"] = nodeNative28
					}
					nodeNative29, _ := jsonpath.Get("$.tunnel_ike_config[0].psk", dataLoopTmp)
					if nodeNative29 != nil && nodeNative29 != "" {
						localData2["Psk"] = nodeNative29
					}
					nodeNative30, _ := jsonpath.Get("$.tunnel_ike_config[0].remote_id", dataLoopTmp)
					if nodeNative30 != nil && nodeNative30 != "" {
						localData2["RemoteId"] = nodeNative30
					}
					dataLoopMap["TunnelIkeConfig"] = localData2
				}
				if !IsNil(dataLoopTmp["tunnel_ipsec_config"]) {
					localData3 := make(map[string]interface{})
					nodeNative31, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_auth_alg", dataLoopTmp)
					if nodeNative31 != nil && nodeNative31 != "" {
						localData3["IpsecAuthAlg"] = nodeNative31
					}
					nodeNative32, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_enc_alg", dataLoopTmp)
					if nodeNative32 != nil && nodeNative32 != "" {
						localData3["IpsecEncAlg"] = nodeNative32
					}
					nodeNative33, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_lifetime", dataLoopTmp)
					if nodeNative33 != nil && nodeNative33 != "" {
						localData3["IpsecLifetime"] = nodeNative33
					}
					nodeNative34, _ := jsonpath.Get("$.tunnel_ipsec_config[0].ipsec_pfs", dataLoopTmp)
					if nodeNative34 != nil && nodeNative34 != "" {
						localData3["IpsecPfs"] = nodeNative34
					}
					dataLoopMap["TunnelIpsecConfig"] = localData3
				}
				dataLoopMap["Role"] = dataLoopTmp["role"]
				tunnelOptionsSpecificationMaps = append(tunnelOptionsSpecificationMaps, dataLoopMap)
			}
			request["TunnelOptionsSpecification"] = tunnelOptionsSpecificationMaps
		}
	}

	if d.HasChange("health_check_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})
		if v := d.Get("health_check_config"); v != nil {
			nodeNative36, _ := jsonpath.Get("$[0].enable", v)
			if nodeNative36 != nil && nodeNative36 != "" {
				objectDataLocalMap3["Enable"] = nodeNative36
			}
			nodeNative37, _ := jsonpath.Get("$[0].dip", v)
			if nodeNative37 != nil && nodeNative37 != "" {
				objectDataLocalMap3["Dip"] = nodeNative37
			}
			nodeNative38, _ := jsonpath.Get("$[0].sip", v)
			if nodeNative38 != nil && nodeNative38 != "" {
				objectDataLocalMap3["Sip"] = nodeNative38
			}
			nodeNative39, _ := jsonpath.Get("$[0].interval", v)
			if nodeNative39 != nil && nodeNative39 != "" {
				objectDataLocalMap3["Interval"] = nodeNative39
			}
			nodeNative40, _ := jsonpath.Get("$[0].retry", v)
			if nodeNative40 != nil && nodeNative40 != "" {
				objectDataLocalMap3["Retry"] = nodeNative40
			}

			objectDataLocalMap3Json, err := json.Marshal(objectDataLocalMap3)
			if err != nil {
				return WrapError(err)
			}
			request["HealthCheckConfig"] = string(objectDataLocalMap3Json)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"VpnGateway.Configuring", "VpnTask.CONFLICT"}) || NeedRetry(err) {
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
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vPNGatewayServiceV2.VPNGatewayVpnConnectionStateRefreshFunc(d.Id(), "State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		if err := vPNGatewayServiceV2.SetResourceTags(d, "VPNCONNECTION"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudVPNGatewayVpnConnectionRead(d, meta)
}

func resourceAliCloudVPNGatewayVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpnConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["VpnConnectionId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || NeedRetry(err) {
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

	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vPNGatewayServiceV2.VPNGatewayVpnConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
