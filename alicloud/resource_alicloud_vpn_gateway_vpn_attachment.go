package alicloud

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

func resourceAliCloudVpnGatewayVpnAttachment() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourceAliCloudVpnGatewayVpnAttachmentV0(),
			Upgrade: resourceAliCloudVpnGatewayVpnAttachmentStateUpgradeV0,
		}},
		CustomizeDiff: vpnAttachmentValidateTunnelChanges,
		Create:        resourceAliCloudVpnGatewayVpnAttachmentCreate,
		Read:          resourceAliCloudVpnGatewayVpnAttachmentRead,
		Update:        resourceAliCloudVpnGatewayVpnAttachmentUpdate,
		Delete:        resourceAliCloudVpnGatewayVpnAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
							Type:     schema.TypeInt,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			"tunnel_bandwidth": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "Large"}, false),
			},
			"health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"reserve_route", "revoke_route"}, false),
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"md5", "sha1", "sha256", "sha384", "sha512"}, false),
						},
						"local_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des"}, false),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"ikev1", "ikev2"}, false),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"main", "aggressive"}, false),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 86400),
						},
						"psk": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"group2", "group1", "group5", "group14"}, false),
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
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"group2", "disabled", "group1", "group5", "group14"}, false),
						},
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des"}, false),
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"md5", "sha1", "sha256", "sha384", "sha512"}, false),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 86400),
						},
					},
				},
			},
			"local_subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"public", "private"}, false),
			},
			"remote_subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Required: true,
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
							Optional: true,
							Computed: true,
						},
						"tunnel_ike_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ike_auth_alg": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"sha1", "md5", "sha256", "sha384", "sha512"}, false),
									},
									"local_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ike_enc_alg": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des"}, false),
									},
									"ike_version": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"ikev1", "ikev2"}, false),
									},
									"ike_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"main", "aggressive"}, false),
									},
									"ike_lifetime": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 86400),
									},
									"psk": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"remote_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ike_pfs": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"group1", "group2", "group5", "group14"}, false),
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
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_asn": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"tunnel_cidr": {
										Type:     schema.TypeString,
										Optional: true,
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
										Optional: true,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntInSlice([]int{0, 1, 2}),
						},
						"enable_nat_traversal": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"tunnel_ipsec_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_pfs": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"disabled", "group1", "group2", "group5", "group14"}, false),
									},
									"ipsec_enc_alg": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des"}, false),
									},
									"ipsec_auth_alg": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"sha1", "md5", "sha256", "sha384", "sha512"}, false),
									},
									"ipsec_lifetime": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 86400),
									},
								},
							},
						},
						"enable_dpd": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"vpn_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudVpnGatewayVpnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpnAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	request["LocalSubnet"] = d.Get("local_subnet")
	request["RemoteSubnet"] = d.Get("remote_subnet")
	if v, ok := d.GetOkExists("effect_immediately"); ok {
		request["EffectImmediately"] = v
	}
	if v, ok := d.GetOkExists("enable_dpd"); ok {
		request["EnableDpd"] = v
	}
	if v, ok := d.GetOkExists("enable_nat_traversal"); ok {
		request["EnableNatTraversal"] = v
	}
	if v, ok := d.GetOk("vpn_attachment_name"); ok {
		request["Name"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("ike_config"); !IsNil(v) {
		psk1, _ := jsonpath.Get("$[0].psk", v)
		if psk1 != nil && psk1 != "" {
			objectDataLocalMap["Psk"] = psk1
		}
		ikeVersion1, _ := jsonpath.Get("$[0].ike_version", v)
		if ikeVersion1 != nil && ikeVersion1 != "" {
			objectDataLocalMap["IkeVersion"] = ikeVersion1
		}
		ikeMode1, _ := jsonpath.Get("$[0].ike_mode", v)
		if ikeMode1 != nil && ikeMode1 != "" {
			objectDataLocalMap["IkeMode"] = ikeMode1
		}
		ikeEncAlg1, _ := jsonpath.Get("$[0].ike_enc_alg", v)
		if ikeEncAlg1 != nil && ikeEncAlg1 != "" {
			objectDataLocalMap["IkeEncAlg"] = ikeEncAlg1
		}
		ikeAuthAlg1, _ := jsonpath.Get("$[0].ike_auth_alg", v)
		if ikeAuthAlg1 != nil && ikeAuthAlg1 != "" {
			objectDataLocalMap["IkeAuthAlg"] = ikeAuthAlg1
		}
		ikePfs1, _ := jsonpath.Get("$[0].ike_pfs", v)
		if ikePfs1 != nil && ikePfs1 != "" {
			objectDataLocalMap["IkePfs"] = ikePfs1
		}
		ikeLifetime1, _ := jsonpath.Get("$[0].ike_lifetime", v)
		if ikeLifetime1 != nil && ikeLifetime1 != "" {
			objectDataLocalMap["IkeLifetime"] = ikeLifetime1
		}
		remoteId1, _ := jsonpath.Get("$[0].remote_id", v)
		if remoteId1 != nil && remoteId1 != "" {
			objectDataLocalMap["RemoteId"] = remoteId1
		}
		localId1, _ := jsonpath.Get("$[0].local_id", v)
		if localId1 != nil && localId1 != "" {
			objectDataLocalMap["LocalId"] = localId1
		}

		request["IkeConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("ipsec_config"); !IsNil(v) {
		ipsecEncAlg1, _ := jsonpath.Get("$[0].ipsec_enc_alg", v)
		if ipsecEncAlg1 != nil && ipsecEncAlg1 != "" {
			objectDataLocalMap1["IpsecEncAlg"] = ipsecEncAlg1
		}
		ipsecAuthAlg1, _ := jsonpath.Get("$[0].ipsec_auth_alg", v)
		if ipsecAuthAlg1 != nil && ipsecAuthAlg1 != "" {
			objectDataLocalMap1["IpsecAuthAlg"] = ipsecAuthAlg1
		}
		ipsecPfs1, _ := jsonpath.Get("$[0].ipsec_pfs", v)
		if ipsecPfs1 != nil && ipsecPfs1 != "" {
			objectDataLocalMap1["IpsecPfs"] = ipsecPfs1
		}
		ipsecLifetime1, _ := jsonpath.Get("$[0].ipsec_lifetime", v)
		if ipsecLifetime1 != nil && ipsecLifetime1 != "" {
			objectDataLocalMap1["IpsecLifetime"] = ipsecLifetime1
		}

		request["IpsecConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap1)
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("bgp_config"); !IsNil(v) {
		localAsn1, _ := jsonpath.Get("$[0].local_asn", v)
		if localAsn1 != nil && localAsn1 != "" {
			objectDataLocalMap2["LocalAsn"] = localAsn1
		}
		tunnelCidr1, _ := jsonpath.Get("$[0].tunnel_cidr", v)
		if tunnelCidr1 != nil && tunnelCidr1 != "" {
			objectDataLocalMap2["TunnelCidr"] = tunnelCidr1
		}
		localBgpIp1, _ := jsonpath.Get("$[0].local_bgp_ip", v)
		if localBgpIp1 != nil && localBgpIp1 != "" {
			objectDataLocalMap2["LocalBgpIp"] = localBgpIp1
		}
		enable, _ := jsonpath.Get("$[0].enable", v)
		if enable != nil && enable != "" {
			objectDataLocalMap2["EnableBgp"] = enable
		}

		request["BgpConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap2)
	}

	if v, ok := d.GetOk("customer_gateway_id"); ok {
		request["CustomerGatewayId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("enable_tunnels_bgp"); ok {
		request["EnableTunnelsBgp"] = v
	}
	if v, ok := d.GetOk("tunnel_bandwidth"); ok {
		request["TunnelBandwidth"] = v
	}
	if v, ok := d.GetOk("tunnel_options_specification"); ok {
		tunnelOptionsSpecificationMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CustomerGatewayId"] = dataLoop1Tmp["customer_gateway_id"]
			dataLoop1Map["EnableDpd"] = dataLoop1Tmp["enable_dpd"]
			dataLoop1Map["EnableNatTraversal"] = dataLoop1Tmp["enable_nat_traversal"]
			dataLoop1Map["Role"] = dataLoop1Tmp["role"]
			dataLoop1Map["TunnelIndex"] = dataLoop1Tmp["tunnel_index"]
			localData2 := make(map[string]interface{})
			localAsn3, _ := jsonpath.Get("$[0].local_asn", dataLoop1Tmp["tunnel_bgp_config"])
			if localAsn3 != nil && localAsn3 != "" {
				localData2["LocalAsn"] = localAsn3
			}
			localBgpIp3, _ := jsonpath.Get("$[0].local_bgp_ip", dataLoop1Tmp["tunnel_bgp_config"])
			if localBgpIp3 != nil && localBgpIp3 != "" {
				localData2["LocalBgpIp"] = localBgpIp3
			}
			tunnelCidr3, _ := jsonpath.Get("$[0].tunnel_cidr", dataLoop1Tmp["tunnel_bgp_config"])
			if tunnelCidr3 != nil && tunnelCidr3 != "" {
				localData2["TunnelCidr"] = tunnelCidr3
			}
			dataLoop1Map["TunnelBgpConfig"] = localData2
			localData3 := make(map[string]interface{})
			ikeAuthAlg3, _ := jsonpath.Get("$[0].ike_auth_alg", dataLoop1Tmp["tunnel_ike_config"])
			if ikeAuthAlg3 != nil && ikeAuthAlg3 != "" {
				localData3["IkeAuthAlg"] = ikeAuthAlg3
			}
			ikeEncAlg3, _ := jsonpath.Get("$[0].ike_enc_alg", dataLoop1Tmp["tunnel_ike_config"])
			if ikeEncAlg3 != nil && ikeEncAlg3 != "" {
				localData3["IkeEncAlg"] = ikeEncAlg3
			}
			ikeLifetime3, _ := jsonpath.Get("$[0].ike_lifetime", dataLoop1Tmp["tunnel_ike_config"])
			if ikeLifetime3 != nil && ikeLifetime3 != "" {
				localData3["IkeLifetime"] = ikeLifetime3
			}
			ikeMode3, _ := jsonpath.Get("$[0].ike_mode", dataLoop1Tmp["tunnel_ike_config"])
			if ikeMode3 != nil && ikeMode3 != "" {
				localData3["IkeMode"] = ikeMode3
			}
			ikePfs3, _ := jsonpath.Get("$[0].ike_pfs", dataLoop1Tmp["tunnel_ike_config"])
			if ikePfs3 != nil && ikePfs3 != "" {
				localData3["IkePfs"] = ikePfs3
			}
			ikeVersion3, _ := jsonpath.Get("$[0].ike_version", dataLoop1Tmp["tunnel_ike_config"])
			if ikeVersion3 != nil && ikeVersion3 != "" {
				localData3["IkeVersion"] = ikeVersion3
			}
			localId3, _ := jsonpath.Get("$[0].local_id", dataLoop1Tmp["tunnel_ike_config"])
			if localId3 != nil && localId3 != "" {
				localData3["LocalId"] = localId3
			}
			psk3, _ := jsonpath.Get("$[0].psk", dataLoop1Tmp["tunnel_ike_config"])
			if psk3 != nil && psk3 != "" {
				localData3["Psk"] = psk3
			}
			remoteId3, _ := jsonpath.Get("$[0].remote_id", dataLoop1Tmp["tunnel_ike_config"])
			if remoteId3 != nil && remoteId3 != "" {
				localData3["RemoteId"] = remoteId3
			}
			dataLoop1Map["TunnelIkeConfig"] = localData3
			localData4 := make(map[string]interface{})
			ipsecAuthAlg3, _ := jsonpath.Get("$[0].ipsec_auth_alg", dataLoop1Tmp["tunnel_ipsec_config"])
			if ipsecAuthAlg3 != nil && ipsecAuthAlg3 != "" {
				localData4["IpsecAuthAlg"] = ipsecAuthAlg3
			}
			ipsecEncAlg3, _ := jsonpath.Get("$[0].ipsec_enc_alg", dataLoop1Tmp["tunnel_ipsec_config"])
			if ipsecEncAlg3 != nil && ipsecEncAlg3 != "" {
				localData4["IpsecEncAlg"] = ipsecEncAlg3
			}
			ipsecLifetime3, _ := jsonpath.Get("$[0].ipsec_lifetime", dataLoop1Tmp["tunnel_ipsec_config"])
			if ipsecLifetime3 != nil && ipsecLifetime3 != "" {
				localData4["IpsecLifetime"] = ipsecLifetime3
			}
			ipsecPfs3, _ := jsonpath.Get("$[0].ipsec_pfs", dataLoop1Tmp["tunnel_ipsec_config"])
			if ipsecPfs3 != nil && ipsecPfs3 != "" {
				localData4["IpsecPfs"] = ipsecPfs3
			}
			dataLoop1Map["TunnelIpsecConfig"] = localData4
			tunnelOptionsSpecificationMapsArray = append(tunnelOptionsSpecificationMapsArray, dataLoop1Map)
		}
		request["TunnelOptionsSpecification"] = tunnelOptionsSpecificationMapsArray
	}

	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("health_check_config"); !IsNil(v) {
		enable2, _ := jsonpath.Get("$[0].enable", v)
		if enable2 != nil && enable2 != "" {
			objectDataLocalMap3["enable"] = enable2
		}
		dip1, _ := jsonpath.Get("$[0].dip", v)
		if dip1 != nil && dip1 != "" {
			objectDataLocalMap3["dip"] = dip1
		}
		sip1, _ := jsonpath.Get("$[0].sip", v)
		if sip1 != nil && sip1 != "" {
			objectDataLocalMap3["sip"] = sip1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && interval1 != "" {
			objectDataLocalMap3["interval"] = interval1
		}
		retry1, _ := jsonpath.Get("$[0].retry", v)
		if retry1 != nil && retry1 != "" {
			objectDataLocalMap3["retry"] = retry1
		}
		policy1, _ := jsonpath.Get("$[0].policy", v)
		if policy1 != nil && policy1 != "" {
			objectDataLocalMap3["Policy"] = policy1
		}

		request["HealthCheckConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap3)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Appliance.Configuring", "VpnGateway.Configuring", "VpnTask.CONFLICT", "TaskConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway_vpn_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpnConnectionId"]))

	return resourceAliCloudVpnGatewayVpnAttachmentUpdate(d, meta)
}

func resourceAliCloudVpnGatewayVpnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}

	objectRaw, err := vPNGatewayServiceV2.DescribeVpnGatewayVpnAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway_vpn_attachment DescribeVpnGatewayVpnAttachment Failed!!! %s", err)
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
	if tunnelBandwidth, ok := objectRaw["TunnelBandwidth"].(string); ok {
		switch strings.ToLower(tunnelBandwidth) {
		case "large":
			d.Set("tunnel_bandwidth", "Large")
		case "standard":
			d.Set("tunnel_bandwidth", "Standard")
		default:
			d.Set("tunnel_bandwidth", tunnelBandwidth)
		}
	} else {
		d.Set("tunnel_bandwidth", objectRaw["TunnelBandwidth"])
	}
	d.Set("local_subnet", objectRaw["LocalSubnet"])
	d.Set("network_type", objectRaw["NetworkType"])
	d.Set("remote_subnet", objectRaw["RemoteSubnet"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["State"])
	d.Set("vpn_attachment_name", objectRaw["Name"])

	bgpConfigMaps := make([]map[string]interface{}, 0)
	bgpConfigMap := make(map[string]interface{})
	vpnBgpConfigRaw := make(map[string]interface{})
	if objectRaw["VpnBgpConfig"] != nil {
		vpnBgpConfigRaw = objectRaw["VpnBgpConfig"].(map[string]interface{})
	}
	if len(vpnBgpConfigRaw) > 0 {
		bgpConfigMap["enable"] = formatBool(vpnBgpConfigRaw["EnableBgp"])
		bgpConfigMap["local_asn"] = vpnBgpConfigRaw["LocalAsn"]
		bgpConfigMap["local_bgp_ip"] = vpnBgpConfigRaw["LocalBgpIp"]
		bgpConfigMap["status"] = vpnBgpConfigRaw["Status"]
		bgpConfigMap["tunnel_cidr"] = vpnBgpConfigRaw["TunnelCidr"]

		bgpConfigMaps = append(bgpConfigMaps, bgpConfigMap)
	}
	if err := d.Set("bgp_config", bgpConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("health_check_config", healthCheckConfigMaps); err != nil {
		return err
	}
	priorIkePsk, _ := d.Get("ike_config.0.psk").(string)
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
		ikeConfigMap["psk"] = vpnAttachmentReadPsk(ikeConfigRaw["Psk"], priorIkePsk)
		ikeConfigMap["remote_id"] = ikeConfigRaw["RemoteId"]

		ikeConfigMaps = append(ikeConfigMaps, ikeConfigMap)
	}
	if err := d.Set("ike_config", ikeConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("ipsec_config", ipsecConfigMaps); err != nil {
		return err
	}
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	tunnelOptionsRaw, _ := jsonpath.Get("$.TunnelOptionsSpecification.TunnelOptions", objectRaw)
	priorTunnelPsk := make(map[int]string)
	for _, raw := range d.Get("tunnel_options_specification").([]interface{}) {
		tunnel, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		psk := ""
		if list, ok := tunnel["tunnel_ike_config"].([]interface{}); ok && len(list) > 0 {
			if ike, ok := list[0].(map[string]interface{}); ok {
				psk, _ = ike["psk"].(string)
			}
		}
		priorTunnelPsk[formatInt(tunnel["tunnel_index"])] = psk
	}
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
				tunnelIkeConfigMap["psk"] = vpnAttachmentReadPsk(tunnelIkeConfigRaw["Psk"], priorTunnelPsk[formatInt(tunnelOptionsChildRaw["TunnelIndex"])])
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
	vpnAttachmentOrderTunnels(tunnelOptionsSpecificationMaps, d.Get("tunnel_options_specification").([]interface{}))
	if err := d.Set("tunnel_options_specification", tunnelOptionsSpecificationMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudVpnGatewayVpnAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyVpnAttachmentAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpnConnectionId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("effect_immediately") {
		update = true
		request["EffectImmediately"] = d.Get("effect_immediately")
	}

	if !d.IsNewResource() && d.HasChange("local_subnet") {
		update = true
	}
	request["LocalSubnet"] = d.Get("local_subnet")
	if !d.IsNewResource() && d.HasChange("remote_subnet") {
		update = true
	}
	request["RemoteSubnet"] = d.Get("remote_subnet")
	if !d.IsNewResource() && d.HasChange("vpn_attachment_name") {
		update = true
		request["Name"] = d.Get("vpn_attachment_name")
	}

	if !d.IsNewResource() && d.HasChange("enable_dpd") {
		update = true
		request["EnableDpd"] = d.Get("enable_dpd")
	}

	if !d.IsNewResource() && d.HasChange("enable_nat_traversal") {
		update = true
		request["EnableNatTraversal"] = d.Get("enable_nat_traversal")
	}

	if !d.IsNewResource() && d.HasChange("bgp_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("bgp_config"); v != nil {
			enable, _ := jsonpath.Get("$[0].enable", v)
			if enable != nil && (d.HasChange("bgp_config.0.enable") || enable != "") {
				objectDataLocalMap["EnableBgp"] = enable
			}
			tunnelCidr1, _ := jsonpath.Get("$[0].tunnel_cidr", v)
			if tunnelCidr1 != nil && (d.HasChange("bgp_config.0.tunnel_cidr") || tunnelCidr1 != "") {
				objectDataLocalMap["TunnelCidr"] = tunnelCidr1
			}
			localBgpIp1, _ := jsonpath.Get("$[0].local_bgp_ip", v)
			if localBgpIp1 != nil && (d.HasChange("bgp_config.0.local_bgp_ip") || localBgpIp1 != "") {
				objectDataLocalMap["LocalBgpIp"] = localBgpIp1
			}
			localAsn1, _ := jsonpath.Get("$[0].local_asn", v)
			if localAsn1 != nil && (d.HasChange("bgp_config.0.local_asn") || localAsn1 != "") {
				objectDataLocalMap["LocalAsn"] = localAsn1
			}

			request["BgpConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
		}
	}

	if !d.IsNewResource() && d.HasChange("customer_gateway_id") {
		update = true
		request["CustomerGatewayId"] = d.Get("customer_gateway_id")
	}

	if !d.IsNewResource() && d.HasChange("enable_tunnels_bgp") {
		update = true
		request["EnableTunnelsBgp"] = d.Get("enable_tunnels_bgp")
	}

	if !d.IsNewResource() && (d.HasChange("tunnel_options_specification") || d.HasChange("enable_tunnels_bgp")) {
		oldValue, newValue := d.GetChange("tunnel_options_specification")
		previous := oldValue.([]interface{})
		planned, restoreErr := vpnAttachmentRestoreInheritedValues(d, previous, newValue.([]interface{}))
		if restoreErr != nil {
			return restoreErr
		}
		// Persist the restored values so the final Read matches PSK priors and
		// ordering against each tunnel's own configuration, not the positional
		// values inherited from another tunnel during plan.
		if err := d.Set("tunnel_options_specification", planned); err != nil {
			return WrapError(err)
		}
		tunnels := vpnAttachmentChangedTunnels(previous, planned, d.HasChange("enable_tunnels_bgp"))
		if len(tunnels) > 0 || len(previous) != len(planned) {
			if err := vpnAttachmentValidateTunnelUpdate(previous, tunnels, d.Get("enable_tunnels_bgp").(bool)); err != nil {
				return WrapError(err)
			}
			update = true
			request["TunnelOptionsSpecification"] = tunnels
		}
	}

	if !d.IsNewResource() && d.HasChange("ike_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("ike_config"); v != nil {
			ikeAuthAlg3, _ := jsonpath.Get("$[0].ike_auth_alg", v)
			if ikeAuthAlg3 != nil && (d.HasChange("ike_config.0.ike_auth_alg") || ikeAuthAlg3 != "") {
				objectDataLocalMap1["IkeAuthAlg"] = ikeAuthAlg3
			}
			localId3, _ := jsonpath.Get("$[0].local_id", v)
			if localId3 != nil && (d.HasChange("ike_config.0.local_id") || localId3 != "") {
				objectDataLocalMap1["LocalId"] = localId3
			}
			ikeEncAlg3, _ := jsonpath.Get("$[0].ike_enc_alg", v)
			if ikeEncAlg3 != nil && (d.HasChange("ike_config.0.ike_enc_alg") || ikeEncAlg3 != "") {
				objectDataLocalMap1["IkeEncAlg"] = ikeEncAlg3
			}
			ikeVersion3, _ := jsonpath.Get("$[0].ike_version", v)
			if ikeVersion3 != nil && (d.HasChange("ike_config.0.ike_version") || ikeVersion3 != "") {
				objectDataLocalMap1["IkeVersion"] = ikeVersion3
			}
			ikeMode3, _ := jsonpath.Get("$[0].ike_mode", v)
			if ikeMode3 != nil && (d.HasChange("ike_config.0.ike_mode") || ikeMode3 != "") {
				objectDataLocalMap1["IkeMode"] = ikeMode3
			}
			ikeLifetime3, _ := jsonpath.Get("$[0].ike_lifetime", v)
			if ikeLifetime3 != nil && (d.HasChange("ike_config.0.ike_lifetime") || ikeLifetime3 != "") {
				objectDataLocalMap1["IkeLifetime"] = ikeLifetime3
			}
			psk3, _ := jsonpath.Get("$[0].psk", v)
			if psk3 != nil && (d.HasChange("ike_config.0.psk") || psk3 != "") {
				objectDataLocalMap1["Psk"] = psk3
			}
			remoteId3, _ := jsonpath.Get("$[0].remote_id", v)
			if remoteId3 != nil && (d.HasChange("ike_config.0.remote_id") || remoteId3 != "") {
				objectDataLocalMap1["RemoteId"] = remoteId3
			}
			ikePfs3, _ := jsonpath.Get("$[0].ike_pfs", v)
			if ikePfs3 != nil && (d.HasChange("ike_config.0.ike_pfs") || ikePfs3 != "") {
				objectDataLocalMap1["IkePfs"] = ikePfs3
			}

			if _, ok := objectDataLocalMap1["Psk"]; !ok {
				return WrapError(fmt.Errorf("ike_config: updating IKE config without the current pre-shared key would make the API generate a random 16-character key; the key is empty in state (for example after import), so set ike_config.psk explicitly and retry"))
			}
			request["IkeConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap1)
		}
	}

	if !d.IsNewResource() && d.HasChange("ipsec_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})

		if v := d.Get("ipsec_config"); v != nil {
			ipsecPfs3, _ := jsonpath.Get("$[0].ipsec_pfs", v)
			if ipsecPfs3 != nil && (d.HasChange("ipsec_config.0.ipsec_pfs") || ipsecPfs3 != "") {
				objectDataLocalMap2["IpsecPfs"] = ipsecPfs3
			}
			ipsecEncAlg3, _ := jsonpath.Get("$[0].ipsec_enc_alg", v)
			if ipsecEncAlg3 != nil && (d.HasChange("ipsec_config.0.ipsec_enc_alg") || ipsecEncAlg3 != "") {
				objectDataLocalMap2["IpsecEncAlg"] = ipsecEncAlg3
			}
			ipsecAuthAlg3, _ := jsonpath.Get("$[0].ipsec_auth_alg", v)
			if ipsecAuthAlg3 != nil && (d.HasChange("ipsec_config.0.ipsec_auth_alg") || ipsecAuthAlg3 != "") {
				objectDataLocalMap2["IpsecAuthAlg"] = ipsecAuthAlg3
			}
			ipsecLifetime3, _ := jsonpath.Get("$[0].ipsec_lifetime", v)
			if ipsecLifetime3 != nil && (d.HasChange("ipsec_config.0.ipsec_lifetime") || ipsecLifetime3 != "") {
				objectDataLocalMap2["IpsecLifetime"] = ipsecLifetime3
			}

			request["IpsecConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap2)
		}
	}

	if d.HasChange("health_check_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

		if v := d.Get("health_check_config"); v != nil {
			enable2, _ := jsonpath.Get("$[0].enable", v)
			if enable2 != nil && (d.HasChange("health_check_config.0.enable") || enable2 != "") {
				objectDataLocalMap3["Enable"] = enable2
			}
			dip1, _ := jsonpath.Get("$[0].dip", v)
			if dip1 != nil && (d.HasChange("health_check_config.0.dip") || dip1 != "") {
				objectDataLocalMap3["Dip"] = dip1
			}
			sip1, _ := jsonpath.Get("$[0].sip", v)
			if sip1 != nil && (d.HasChange("health_check_config.0.sip") || sip1 != "") {
				objectDataLocalMap3["Sip"] = sip1
			}
			interval1, _ := jsonpath.Get("$[0].interval", v)
			if interval1 != nil && (d.HasChange("health_check_config.0.interval") || interval1 != "") {
				objectDataLocalMap3["Interval"] = interval1
			}
			retry1, _ := jsonpath.Get("$[0].retry", v)
			if retry1 != nil && (d.HasChange("health_check_config.0.retry") || retry1 != "") {
				objectDataLocalMap3["Retry"] = retry1
			}
			policy1, _ := jsonpath.Get("$[0].policy", v)
			if policy1 != nil && (d.HasChange("health_check_config.0.policy") || policy1 != "") {
				objectDataLocalMap3["Policy"] = policy1
			}

			request["HealthCheckConfig"] = convertMapToJsonStringIgnoreError(objectDataLocalMap3)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"VpnConnection.Configuring", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnGateway.Configuring"}) || NeedRetry(err) {
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
		vpnService := VPNGatewayServiceV2{client}
		stateConf := BuildStateConf(
			[]string{"updating", "provisioning", "upgrading", "attaching", "detaching"},
			[]string{"attached", "active", "init"}, d.Timeout(schema.TimeoutUpdate), 0,
			vpnService.VpnGatewayVpnAttachmentStateRefreshFunc(d.Id(), "State", []string{"financialLocked", "deleted"}),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveVpnResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "VPNATTACHMENT"
	if update {
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
	}

	if d.HasChange("tags") {
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		if err := vPNGatewayServiceV2.SetResourceTags(d, "VPNATTACHMENT"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpnGatewayVpnAttachmentRead(d, meta)
}

func resourceAliCloudVpnGatewayVpnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpnAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VpnConnectionId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

// Empty collections are omitted by RPC serialization, so they cannot remove
// existing tunnels or reset a configuration block. Reject these plans before
// sending an update that would leave the remote values unchanged.
func vpnAttachmentValidateTunnelChanges(d *schema.ResourceDiff, meta interface{}) error {
	const key = "tunnel_options_specification"
	if !d.NewValueKnown(key) {
		return nil
	}
	previous, planned := d.GetChange(key)
	positions := make(map[int]int)
	for i, raw := range planned.([]interface{}) {
		if !d.NewValueKnown(fmt.Sprintf("%s.%d.tunnel_index", key, i)) {
			return nil
		}
		index := formatInt(raw.(map[string]interface{})["tunnel_index"])
		positions[index] = i
	}
	for oldPosition, raw := range previous.([]interface{}) {
		old := raw.(map[string]interface{})
		index := formatInt(old["tunnel_index"])
		position, exists := positions[index]
		if !exists {
			return fmt.Errorf("tunnel_options_specification cannot remove existing tunnel_index %d; retain both tunnels in the configuration", index)
		}
		next := planned.([]interface{})[position].(map[string]interface{})
		var atPosition map[string]interface{}
		if position != oldPosition && position < len(previous.([]interface{})) {
			atPosition = previous.([]interface{})[position].(map[string]interface{})
		}
		for _, block := range []string{"tunnel_bgp_config", "tunnel_ike_config", "tunnel_ipsec_config"} {
			if !d.NewValueKnown(fmt.Sprintf("%s.%d.%s", key, position, block)) || IsNil(old[block]) || !IsNil(next[block]) {
				continue
			}
			if atPosition != nil && IsNil(atPosition[block]) {
				// A moved tunnel inherits the positional predecessor's absent
				// block; this is not an explicit clear. Update restores the
				// tunnel's own block from the previous state.
				continue
			}
			return fmt.Errorf("%s for tunnel_index %d cannot be cleared with an empty list; omit the block to keep its current values", block, index)
		}
	}
	return nil
}

// vpnAttachmentOrderTunnels preserves list positions across API responses. New or
// imported tunnels use tunnel_index order so refreshes do not depend on API order.
func vpnAttachmentOrderTunnels(tunnels []map[string]interface{}, previous []interface{}) {
	positions := make(map[int]int, len(previous))
	for i, raw := range previous {
		positions[formatInt(raw.(map[string]interface{})["tunnel_index"])] = i
	}
	sort.SliceStable(tunnels, func(i, j int) bool {
		left, right := formatInt(tunnels[i]["tunnel_index"]), formatInt(tunnels[j]["tunnel_index"])
		li, lok := positions[left]
		ri, rok := positions[right]
		if lok && rok {
			return li < ri
		}
		if lok != rok {
			return lok
		}
		return left < right
	})
}

// An empty PSK in a Describe response discloses no usable secret, so retain
// the prior state value. Any non-empty value is adopted as real, including
// all-asterisk values: '*' is a legal PSK character and the API documents no
// all-asterisk mask form (the observed masked form keeps a plaintext prefix,
// e.g. "123456****"). Treating an all-asterisk value as masked would hide an
// externally rotated key and let a later update restore the stale prior.
func vpnAttachmentReadPsk(remote interface{}, prior string) string {
	value, _ := remote.(string)
	if value == "" {
		return prior
	}
	return value
}

// In the merged planned state, a reordered tunnel inherits omitted
// Optional+Computed fields from the tunnel previously at the same list
// position, not from itself. Restore each moved tunnel's own previous values
// for every leaf the plan did not present as a change, so change detection,
// the update request and the persisted state use the configuration the user
// actually wrote.
//
// The legacy positional list diff cannot present an explicit write whose new
// value equals the positional predecessor's current value: it emits no entry
// for that leaf, exactly as for an inherited value, and raw configuration is
// unreachable during apply. When the tunnel is actively modified in the same
// apply, a leaf without a plan-visible change whose planned value differs
// from the tunnel's own previous value is therefore ambiguous, and keeping
// the own value would silently defer an explicit change to a second apply.
// Such updates are rejected with guidance to split the reorder and the value
// change instead. A scope with no plan-visible change at all is an omission:
// the own values are restored, which is what keeps a pure reorder from
// writing another tunnel's values.
func vpnAttachmentRestoreInheritedValues(d *schema.ResourceData, previous, planned []interface{}) ([]interface{}, error) {
	restored := vpnAttachmentCopyTunnels(planned)
	oldByIndex := make(map[int]map[string]interface{}, len(previous))
	for _, raw := range previous {
		old := raw.(map[string]interface{})
		oldByIndex[formatInt(old["tunnel_index"])] = old
	}
	for position, raw := range restored {
		tunnel := raw.(map[string]interface{})
		own := oldByIndex[formatInt(tunnel["tunnel_index"])]
		if own == nil || position >= len(previous) {
			continue
		}
		atPosition := previous[position].(map[string]interface{})
		if formatInt(atPosition["tunnel_index"]) == formatInt(tunnel["tunnel_index"]) {
			continue
		}
		prefix := fmt.Sprintf("tunnel_options_specification.%d", position)
		conflicts := vpnAttachmentReorderConflicts(d, prefix, tunnel, own)
		if len(conflicts) > 0 {
			sort.Strings(conflicts)
			return nil, fmt.Errorf("tunnel_options_specification: reordering tunnels while changing their configuration cannot be applied reliably: the plan presents no change for %s, so its planned value would be inherited from the tunnel previously at the same list position, and the SDK cannot distinguish that inheritance from an explicit write of the same value. Split the operation into two applies (first apply the reorder without changing tunnel values, then apply the value changes) and retry; the update was not sent to the API", strings.Join(conflicts, ", "))
		}
		for _, field := range []string{"role", "enable_dpd", "enable_nat_traversal"} {
			if !d.HasChange(prefix + "." + field) {
				tunnel[field] = own[field]
			}
		}
		for _, block := range []string{"tunnel_bgp_config", "tunnel_ike_config", "tunnel_ipsec_config"} {
			vpnAttachmentRestoreInheritedBlock(d, prefix, tunnel, block, own, atPosition)
		}
	}
	return restored, nil
}

// Report the leaves of a moved tunnel whose planned value differs from the
// tunnel's own previous value without a plan-visible change, while the scope
// holding them is actively modified: a field group (role, enable_dpd,
// enable_nat_traversal) with another visible change, or a nested block with
// another visible leaf. A block the configuration leaves out entirely shows
// no visible leaf and keeps the omission semantics: it restores the tunnel's
// own values silently. A tunnel that never had a block has no own value to
// contrast, so its incomplete writes stay on the existing explicit
// validation path.
//
// Computed-only leaves (such as tunnel_bgp_config bgp_status, peer_bgp_ip and
// peer_asn) are excluded: the positional list diff inherits them from the
// tunnel previously at the same position, they can never carry an explicit
// user write, and the block restore keeps each tunnel's own values for them.
func vpnAttachmentReorderConflicts(d *schema.ResourceData, prefix string, tunnel, own map[string]interface{}) []string {
	var conflicts []string
	fields := []string{"role", "enable_dpd", "enable_nat_traversal"}
	fieldVisible := false
	for _, field := range fields {
		if d.HasChange(prefix + "." + field) {
			fieldVisible = true
		}
	}
	if fieldVisible {
		for _, field := range fields {
			if !d.HasChange(prefix+"."+field) && !reflect.DeepEqual(tunnel[field], own[field]) {
				conflicts = append(conflicts, prefix+"."+field)
			}
		}
	}
	for _, block := range []string{"tunnel_bgp_config", "tunnel_ike_config", "tunnel_ipsec_config"} {
		current, _ := tunnel[block].([]interface{})
		if len(current) == 0 {
			continue
		}
		currentMap, _ := current[0].(map[string]interface{})
		ownBlock, _ := own[block].([]interface{})
		if currentMap == nil || len(ownBlock) == 0 {
			continue
		}
		ownMap, _ := ownBlock[0].(map[string]interface{})
		if ownMap == nil {
			continue
		}
		leafPrefix := prefix + "." + block + ".0."
		leafVisible := false
		for key := range currentMap {
			if d.HasChange(leafPrefix + key) {
				leafVisible = true
				break
			}
		}
		if !leafVisible {
			continue
		}
		blockSchema := vpnAttachmentTunnelBlockSchema(block)
		for key := range currentMap {
			if ownValue, ok := ownMap[key]; ok && !d.HasChange(leafPrefix+key) && !reflect.DeepEqual(currentMap[key], ownValue) {
				if leaf := blockSchema[key]; leaf != nil && leaf.Computed && !leaf.Optional {
					continue
				}
				conflicts = append(conflicts, leafPrefix+key)
			}
		}
	}
	return conflicts
}

// vpnAttachmentTunnelBlockSchema resolves a nested tunnel block's leaf
// schemas from the declared resource schema, so the conflict scan can tell
// user-writable leaves from Computed-only ones without duplicating key lists.
func vpnAttachmentTunnelBlockSchema(block string) map[string]*schema.Schema {
	tunnel, _ := resourceAliCloudVpnGatewayVpnAttachment().Schema["tunnel_options_specification"].Elem.(*schema.Resource)
	if tunnel == nil {
		return nil
	}
	nested, _ := tunnel.Schema[block].Elem.(*schema.Resource)
	if nested == nil {
		return nil
	}
	return nested.Schema
}

func vpnAttachmentRestoreInheritedBlock(d *schema.ResourceData, prefix string, tunnel map[string]interface{}, block string, own, atPosition map[string]interface{}) {
	current, _ := tunnel[block].([]interface{})
	positional, _ := atPosition[block].([]interface{})
	ownBlock, _ := own[block].([]interface{})
	blockAddr := prefix + "." + block
	if len(current) == 0 {
		// The merged value lost this block only because the positional
		// predecessor had none to inherit; a plan-visible removal is an
		// explicit clear, rejected by vpnAttachmentValidateTunnelChanges.
		if len(ownBlock) > 0 && !d.HasChange(blockAddr) {
			tunnel[block] = vpnAttachmentCopyBlock(ownBlock)
		}
		return
	}
	currentMap, _ := current[0].(map[string]interface{})
	if currentMap == nil || len(positional) == 0 {
		return
	}
	positionalMap, _ := positional[0].(map[string]interface{})
	if positionalMap == nil {
		return
	}
	var ownMap map[string]interface{}
	if len(ownBlock) > 0 {
		ownMap, _ = ownBlock[0].(map[string]interface{})
	}
	leafPrefix := blockAddr + ".0."
	visible := false
	for key := range currentMap {
		if d.HasChange(leafPrefix + key) {
			visible = true
			break
		}
	}
	if !visible {
		// No leaf shows a plan-visible change: the whole block was inherited
		// from the positional predecessor.
		if ownMap == nil {
			tunnel[block] = []interface{}{}
		} else {
			tunnel[block] = vpnAttachmentCopyBlock(ownBlock)
		}
		return
	}
	if ownMap == nil {
		// The tunnel never had this block, so every leaf without a
		// plan-visible change was inherited from the positional predecessor
		// and must not be submitted as this tunnel's configuration. Keep
		// only the leaves the user actually wrote; an incomplete block is
		// rejected before the update by vpnAttachmentValidateTunnelUpdate.
		for key := range currentMap {
			if !d.HasChange(leafPrefix + key) {
				delete(currentMap, key)
			}
		}
		return
	}
	for key := range currentMap {
		if d.HasChange(leafPrefix + key) {
			continue
		}
		if ownValue, ok := ownMap[key]; ok {
			currentMap[key] = ownValue
		} else {
			delete(currentMap, key)
		}
	}
}

func vpnAttachmentCopyTunnels(tunnels []interface{}) []interface{} {
	copied := make([]interface{}, 0, len(tunnels))
	for _, raw := range tunnels {
		tunnel, ok := raw.(map[string]interface{})
		if !ok {
			copied = append(copied, raw)
			continue
		}
		result := make(map[string]interface{}, len(tunnel))
		for key, value := range tunnel {
			if list, ok := value.([]interface{}); ok && len(list) > 0 {
				if _, ok := list[0].(map[string]interface{}); ok {
					result[key] = vpnAttachmentCopyBlock(list)
					continue
				}
			}
			result[key] = value
		}
		copied = append(copied, result)
	}
	return copied
}

func vpnAttachmentCopyBlock(block []interface{}) []interface{} {
	copied := make([]interface{}, 0, len(block))
	for _, raw := range block {
		if m, ok := raw.(map[string]interface{}); ok {
			result := make(map[string]interface{}, len(m))
			for key, value := range m {
				result[key] = value
			}
			copied = append(copied, result)
		} else {
			copied = append(copied, raw)
		}
	}
	return copied
}

// Compare API configuration by tunnel_index. Computed status, peer addresses and
// tunnel IDs must not turn a reorder into a remote update. The planned values
// must already be normalized by vpnAttachmentRestoreInheritedValues.
func vpnAttachmentChangedTunnels(previous, planned []interface{}, includeBGP bool) []interface{} {
	oldByIndex := make(map[int]map[string]interface{}, len(previous))
	for _, raw := range previous {
		old := raw.(map[string]interface{})
		oldByIndex[formatInt(old["tunnel_index"])] = vpnAttachmentExpandTunnel(old)
	}
	result := make([]interface{}, 0, len(planned))
	for _, raw := range planned {
		tunnel := raw.(map[string]interface{})
		next := vpnAttachmentExpandTunnel(tunnel)
		old := oldByIndex[formatInt(tunnel["tunnel_index"])]
		if !includeBGP && reflect.DeepEqual(old, next) {
			continue
		}
		// Avoid resending unrelated IKE/PSK or IPsec settings for a BGP update.
		// If IKE changes, send its full block: omitting only Psk can generate a key.
		for _, group := range []string{"TunnelIkeConfig", "TunnelIpsecConfig", "TunnelBgpConfig"} {
			if !(includeBGP && group == "TunnelBgpConfig") && old != nil && reflect.DeepEqual(old[group], next[group]) {
				delete(next, group)
			}
		}
		result = append(result, next)
	}
	return result
}

// Reject tunnel update payloads the API cannot honor safely, before any RPC:
// an IKE block without Psk makes the API generate a random 16-character key,
// so it is only acceptable for a newly added tunnel, which has no key to
// reset; a BGP block without TunnelCidr is rejected by the API whenever BGP
// is enabled, which is how an incomplete block newly written during a reorder
// surfaces instead of silently inheriting the other tunnel's values.
func vpnAttachmentValidateTunnelUpdate(previous, tunnels []interface{}, enableTunnelsBgp bool) error {
	existing := make(map[int]bool, len(previous))
	for _, raw := range previous {
		existing[formatInt(raw.(map[string]interface{})["tunnel_index"])] = true
	}
	for _, raw := range tunnels {
		tunnel, _ := raw.(map[string]interface{})
		index := formatInt(tunnel["TunnelIndex"])
		if ike, ok := tunnel["TunnelIkeConfig"].(map[string]interface{}); ok && existing[index] {
			if psk, _ := ike["Psk"].(string); psk == "" {
				return fmt.Errorf("tunnel_options_specification tunnel_index %d: updating tunnel_ike_config without the current pre-shared key would make the API generate a random 16-character key; the key is empty in state (for example after import), so set tunnel_ike_config.psk explicitly and retry", index)
			}
		}
		if bgp, ok := tunnel["TunnelBgpConfig"].(map[string]interface{}); ok && enableTunnelsBgp {
			if cidr, _ := bgp["TunnelCidr"].(string); cidr == "" {
				return fmt.Errorf("tunnel_options_specification tunnel_index %d: tunnel_bgp_config.tunnel_cidr is required when enable_tunnels_bgp is true; write the complete tunnel_bgp_config block and retry", index)
			}
		}
	}
	return nil
}

func vpnAttachmentExpandTunnel(tunnel map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["CustomerGatewayId"] = tunnel["customer_gateway_id"]
	result["EnableDpd"] = tunnel["enable_dpd"]
	result["EnableNatTraversal"] = tunnel["enable_nat_traversal"]
	result["Role"] = tunnel["role"]
	result["TunnelIndex"] = tunnel["tunnel_index"]
	if !IsNil(tunnel["tunnel_bgp_config"]) {
		localData1 := make(map[string]interface{})
		localAsn3, _ := jsonpath.Get("$[0].local_asn", tunnel["tunnel_bgp_config"])
		if localAsn3 != nil && localAsn3 != "" {
			localData1["LocalAsn"] = localAsn3
		}
		localBgpIp3, _ := jsonpath.Get("$[0].local_bgp_ip", tunnel["tunnel_bgp_config"])
		if localBgpIp3 != nil && localBgpIp3 != "" {
			localData1["LocalBgpIp"] = localBgpIp3
		}
		tunnelCidr3, _ := jsonpath.Get("$[0].tunnel_cidr", tunnel["tunnel_bgp_config"])
		if tunnelCidr3 != nil && tunnelCidr3 != "" {
			localData1["TunnelCidr"] = tunnelCidr3
		}
		result["TunnelBgpConfig"] = localData1
	}
	if !IsNil(tunnel["tunnel_ike_config"]) {
		localData2 := make(map[string]interface{})
		ikeAuthAlg1, _ := jsonpath.Get("$[0].ike_auth_alg", tunnel["tunnel_ike_config"])
		if ikeAuthAlg1 != nil && ikeAuthAlg1 != "" {
			localData2["IkeAuthAlg"] = ikeAuthAlg1
		}
		ikeEncAlg1, _ := jsonpath.Get("$[0].ike_enc_alg", tunnel["tunnel_ike_config"])
		if ikeEncAlg1 != nil && ikeEncAlg1 != "" {
			localData2["IkeEncAlg"] = ikeEncAlg1
		}
		ikeLifetime1, _ := jsonpath.Get("$[0].ike_lifetime", tunnel["tunnel_ike_config"])
		if ikeLifetime1 != nil && ikeLifetime1 != "" {
			localData2["IkeLifetime"] = ikeLifetime1
		}
		ikeMode1, _ := jsonpath.Get("$[0].ike_mode", tunnel["tunnel_ike_config"])
		if ikeMode1 != nil && ikeMode1 != "" {
			localData2["IkeMode"] = ikeMode1
		}
		ikePfs1, _ := jsonpath.Get("$[0].ike_pfs", tunnel["tunnel_ike_config"])
		if ikePfs1 != nil && ikePfs1 != "" {
			localData2["IkePfs"] = ikePfs1
		}
		ikeVersion1, _ := jsonpath.Get("$[0].ike_version", tunnel["tunnel_ike_config"])
		if ikeVersion1 != nil && ikeVersion1 != "" {
			localData2["IkeVersion"] = ikeVersion1
		}
		localId1, _ := jsonpath.Get("$[0].local_id", tunnel["tunnel_ike_config"])
		if localId1 != nil && localId1 != "" {
			localData2["LocalId"] = localId1
		}
		psk1, _ := jsonpath.Get("$[0].psk", tunnel["tunnel_ike_config"])
		if psk1 != nil && psk1 != "" {
			localData2["Psk"] = psk1
		}
		remoteId1, _ := jsonpath.Get("$[0].remote_id", tunnel["tunnel_ike_config"])
		if remoteId1 != nil && remoteId1 != "" {
			localData2["RemoteId"] = remoteId1
		}
		result["TunnelIkeConfig"] = localData2
	}
	if !IsNil(tunnel["tunnel_ipsec_config"]) {
		localData3 := make(map[string]interface{})
		ipsecAuthAlg1, _ := jsonpath.Get("$[0].ipsec_auth_alg", tunnel["tunnel_ipsec_config"])
		if ipsecAuthAlg1 != nil && ipsecAuthAlg1 != "" {
			localData3["IpsecAuthAlg"] = ipsecAuthAlg1
		}
		ipsecEncAlg1, _ := jsonpath.Get("$[0].ipsec_enc_alg", tunnel["tunnel_ipsec_config"])
		if ipsecEncAlg1 != nil && ipsecEncAlg1 != "" {
			localData3["IpsecEncAlg"] = ipsecEncAlg1
		}
		ipsecLifetime1, _ := jsonpath.Get("$[0].ipsec_lifetime", tunnel["tunnel_ipsec_config"])
		if ipsecLifetime1 != nil && ipsecLifetime1 != "" {
			localData3["IpsecLifetime"] = ipsecLifetime1
		}
		ipsecPfs1, _ := jsonpath.Get("$[0].ipsec_pfs", tunnel["tunnel_ipsec_config"])
		if ipsecPfs1 != nil && ipsecPfs1 != "" {
			localData3["IpsecPfs"] = ipsecPfs1
		}
		result["TunnelIpsecConfig"] = localData3
	}
	return result
}

// resourceAliCloudVpnGatewayVpnAttachmentV0 freezes the complete version 0 state
// type, including the SDK's implicit id and timeouts attributes. Do not derive
// this from the current resource: future schema changes must not alter decoding
// of existing set-based state.
func resourceAliCloudVpnGatewayVpnAttachmentV0() cty.Type {
	ikeConfig := cty.List(cty.Object(map[string]cty.Type{
		"ike_auth_alg": cty.String,
		"ike_enc_alg":  cty.String,
		"ike_lifetime": cty.Number,
		"ike_mode":     cty.String,
		"ike_pfs":      cty.String,
		"ike_version":  cty.String,
		"local_id":     cty.String,
		"psk":          cty.String,
		"remote_id":    cty.String,
	}))
	ipsecConfig := cty.List(cty.Object(map[string]cty.Type{
		"ipsec_auth_alg": cty.String,
		"ipsec_enc_alg":  cty.String,
		"ipsec_lifetime": cty.Number,
		"ipsec_pfs":      cty.String,
	}))
	return cty.Object(map[string]cty.Type{
		"bgp_config": cty.List(cty.Object(map[string]cty.Type{
			"enable":       cty.Bool,
			"local_asn":    cty.Number,
			"local_bgp_ip": cty.String,
			"status":       cty.String,
			"tunnel_cidr":  cty.String,
		})),
		"create_time":          cty.String,
		"customer_gateway_id":  cty.String,
		"effect_immediately":   cty.Bool,
		"enable_dpd":           cty.Bool,
		"enable_nat_traversal": cty.Bool,
		"enable_tunnels_bgp":   cty.Bool,
		"health_check_config": cty.List(cty.Object(map[string]cty.Type{
			"dip":      cty.String,
			"enable":   cty.Bool,
			"interval": cty.Number,
			"policy":   cty.String,
			"retry":    cty.Number,
			"sip":      cty.String,
			"status":   cty.String,
		})),
		"id":                cty.String,
		"ike_config":        ikeConfig,
		"ipsec_config":      ipsecConfig,
		"local_subnet":      cty.String,
		"network_type":      cty.String,
		"remote_subnet":     cty.String,
		"resource_group_id": cty.String,
		"status":            cty.String,
		"tags":              cty.Map(cty.String),
		"timeouts": cty.Object(map[string]cty.Type{
			"create": cty.String,
			"delete": cty.String,
			"update": cty.String,
		}),
		"tunnel_bandwidth": cty.String,
		"tunnel_options_specification": cty.Set(cty.Object(map[string]cty.Type{
			"customer_gateway_id":  cty.String,
			"enable_dpd":           cty.Bool,
			"enable_nat_traversal": cty.Bool,
			"internet_ip":          cty.String,
			"role":                 cty.String,
			"state":                cty.String,
			"status":               cty.String,
			"tunnel_bgp_config": cty.List(cty.Object(map[string]cty.Type{
				"bgp_status":   cty.String,
				"local_asn":    cty.Number,
				"local_bgp_ip": cty.String,
				"peer_asn":     cty.String,
				"peer_bgp_ip":  cty.String,
				"tunnel_cidr":  cty.String,
			})),
			"tunnel_id":           cty.String,
			"tunnel_ike_config":   ikeConfig,
			"tunnel_index":        cty.Number,
			"tunnel_ipsec_config": ipsecConfig,
			"zone_no":             cty.String,
		})),
		"vpn_attachment_name": cty.String,
	})
}

func resourceAliCloudVpnGatewayVpnAttachmentStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState["tunnel_options_specification"] == nil {
		return rawState, nil
	}
	tunnels, ok := rawState["tunnel_options_specification"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid version 0 tunnel_options_specification: expected an array")
	}
	for _, value := range tunnels {
		tunnel, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid version 0 tunnel_options_specification: expected a tunnel object")
		}
		// StateUpgrader receives the SDK's default JSON decoding, whose numbers
		// are float64 values even when the schema declares TypeInt.
		if _, ok := tunnel["tunnel_index"].(float64); !ok {
			return nil, fmt.Errorf("invalid version 0 tunnel_options_specification: expected a numeric tunnel_index")
		}
	}
	// A set did not retain HCL block order. Establish deterministic state order
	// from tunnel_index while preserving every tunnel attribute, including PSK.
	sort.SliceStable(tunnels, func(i, j int) bool {
		return tunnels[i].(map[string]interface{})["tunnel_index"].(float64) < tunnels[j].(map[string]interface{})["tunnel_index"].(float64)
	})
	return rawState, nil
}
