package alicloud

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudVpnGatewayVpnAttachment() *schema.Resource {
	r := &schema.Resource{
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
				Type:     schema.TypeSet,
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
							// Preserve the server-reported role even when legacy configuration supplies one.
							DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool { return true },
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
	// Role is returned by the service but is not accepted by either attachment
	// write API. Retain its configuration syntax for compatibility, while the
	// hash includes every other writable field and excludes service-only values.
	tunnels := r.Schema["tunnel_options_specification"]
	fields := tunnels.Elem.(*schema.Resource).Schema
	hashFields := make(map[string]*schema.Schema, len(fields))
	for name, field := range fields {
		hashFields[name] = field
	}
	role := *fields["role"]
	role.Optional = false
	role.Computed = true
	hashFields["role"] = &role
	tunnels.Set = schema.HashResource(&schema.Resource{Schema: hashFields})
	return r
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
		for _, dataLoop1 := range vpnAttachmentSortedTunnels(v.(*schema.Set).List()) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CustomerGatewayId"] = dataLoop1Tmp["customer_gateway_id"]
			dataLoop1Map["EnableDpd"] = dataLoop1Tmp["enable_dpd"]
			dataLoop1Map["EnableNatTraversal"] = dataLoop1Tmp["enable_nat_traversal"]
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
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Appliance.Configuring", "VpnGateway.Configuring", "VpnTask.CONFLICT", "TaskConflict"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
	for _, raw := range d.Get("tunnel_options_specification").(*schema.Set).List() {
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
		previous := oldValue.(*schema.Set).List()
		planned := newValue.(*schema.Set).List()
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

			if psk, ok := objectDataLocalMap1["Psk"].(string); !ok || psk == "" {
				return WrapError(fmt.Errorf("ike_config: updating IKE config without the current pre-shared key would make the API generate a random 16-character key; the key is empty in state or configuration (for example after import or an explicit empty value), so set a non-empty ike_config.psk explicitly and retry"))
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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"VpnConnection.Configuring", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnGateway.Configuring"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
//
// The guards hold for in-place updates only: a change to a ForceNew field
// replaces the whole attachment, and the removed configuration is simply not
// part of the new resource.
// Full writable hashes detect nested edits. An omitted Optional+Computed value
// belongs to the same tunnel_index, even when another edit changes its set hash.
func vpnAttachmentValidateTunnelChanges(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// The legacy SDK turns an explicit empty Optional+Computed PSK into a
	// computed value and can restore the old key during Apply. Reject that
	// explicit request while configuration presence is still observable.
	if d.Id() != "" && !vpnAttachmentForceNewChanged(d) && d.HasChange("ike_config") && d.NewValueKnown("ike_config.0.psk") {
		if psk, exists := d.GetOkExists("ike_config.0.psk"); exists && psk == "" {
			return fmt.Errorf("ike_config: updating IKE config requires a non-empty pre-shared key; ike_config.psk is empty or unavailable")
		}
	}
	const key = "tunnel_options_specification"
	if !d.NewValueKnown(key) {
		return nil
	}
	// The SDK prefixes set hashes with ~ when configuration contains unknowns.
	// Keep that diff intact so an explicit unknown never becomes a prior value.
	if len(d.GetChangedKeysPrefix(key+".~")) != 0 {
		return nil
	}
	previousValue, plannedValue := d.GetChange(key)
	previous := previousValue.(*schema.Set).List()
	plannedSet := plannedValue.(*schema.Set)
	planned := vpnAttachmentCopyTunnels(plannedSet.List())
	oldByIndex := make(map[int]map[string]interface{}, len(previous))
	for _, raw := range previous {
		old := raw.(map[string]interface{})
		oldByIndex[formatInt(old["tunnel_index"])] = old
	}
	// DiffFieldReader rehashes set members after merging computed fields, so
	// its hashes need not equal the addresses in the original diff. Resolve
	// those addresses by the required identity fields before inheriting.
	prefixes := make(map[int]string)
	for _, address := range d.GetChangedKeysPrefix(key + ".") {
		if !strings.HasSuffix(address, ".tunnel_index") {
			continue
		}
		prefix := strings.TrimSuffix(address, ".tunnel_index")
		if gateway, ok := d.Get(prefix + ".customer_gateway_id").(string); ok && gateway != "" {
			prefixes[formatInt(d.Get(address))] = prefix
		}
	}
	seen := make(map[int]bool)
	fields := resourceAliCloudVpnGatewayVpnAttachment().Schema[key].Elem.(*schema.Resource).Schema
	for _, raw := range planned {
		next := raw.(map[string]interface{})
		index := formatInt(next["tunnel_index"])
		if seen[index] {
			return fmt.Errorf("tunnel_options_specification contains duplicate tunnel_index %d", index)
		}
		seen[index] = true
		if old := oldByIndex[index]; old != nil {
			prefix := prefixes[index]
			if prefix == "" {
				prefix = fmt.Sprintf("%s.%d", key, plannedSet.F(next))
			}
			if err := vpnAttachmentInheritComputed(d, prefix, next, old, fields, index); err != nil {
				return err
			}
		}
	}
	if !vpnAttachmentForceNewChanged(d) {
		for index := range oldByIndex {
			if !seen[index] {
				return fmt.Errorf("tunnel_options_specification cannot remove existing tunnel_index %d; retain both tunnels in the configuration", index)
			}
		}
	}
	if !reflect.DeepEqual(planned, plannedSet.List()) {
		return d.SetNew(key, planned)
	}
	return nil
}

func vpnAttachmentInheritComputed(d *schema.ResourceDiff, prefix string, next, old map[string]interface{}, fields map[string]*schema.Schema, index int) error {
	for name, field := range fields {
		address := prefix + "." + name
		if field.Computed && !field.Optional {
			next[name] = old[name]
			continue
		}
		// GetOkExists preserves explicit false, zero and empty strings. This
		// SDK mutates a set hash address to a config list position while
		// reading, so NewValueKnown alone can report an absent field as known.
		// Read each original address as well: a merged set member can contain
		// nil nested blocks even when its configuration contains a block.
		value, exists := d.GetOkExists(address)
		if field.Computed && !exists {
			next[name] = old[name]
			continue
		}
		if exists {
			next[name] = value
			if block, ok := value.([]interface{}); ok {
				next[name] = vpnAttachmentCopyBlock(block)
			}
		}
		if nested, ok := field.Elem.(*schema.Resource); ok && field.Type == schema.TypeList {
			before, _ := old[name].([]interface{})
			after, _ := next[name].([]interface{})
			if len(before) > 0 && len(after) == 0 && !vpnAttachmentForceNewChanged(d) {
				return fmt.Errorf("%s for tunnel_index %d cannot be cleared with an empty list; omit the block to keep its current values", name, index)
			}
			if len(before) > 0 && len(after) > 0 && before[0] != nil && after[0] != nil {
				if err := vpnAttachmentInheritComputed(d, address+".0", after[0].(map[string]interface{}), before[0].(map[string]interface{}), nested.Schema, index); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func vpnAttachmentSortedTunnels(tunnels []interface{}) []interface{} {
	sort.Slice(tunnels, func(i, j int) bool {
		return formatInt(tunnels[i].(map[string]interface{})["tunnel_index"]) < formatInt(tunnels[j].(map[string]interface{})["tunnel_index"])
	})
	return tunnels
}

// vpnAttachmentForceNewChanged reports whether any schema attribute marked
// ForceNew has a plan-visible change, which means the SDK plans a replacement
// instead of an in-place update.
func vpnAttachmentForceNewChanged(d *schema.ResourceDiff) bool {
	for name, s := range resourceAliCloudVpnGatewayVpnAttachment().Schema {
		if s.ForceNew && d.HasChange(name) {
			return true
		}
	}
	return false
}

// An empty Describe PSK retains the prior secret. Treat every non-empty
// response as a real value rather than guessing a mask format.
func vpnAttachmentReadPsk(remote interface{}, prior string) string {
	value, _ := remote.(string)
	if value == "" {
		return prior
	}
	return value
}

// Copy nested blocks before normalizing the planned set.
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
// are normalized by tunnel_index during planning.
func vpnAttachmentChangedTunnels(previous, planned []interface{}, includeBGP bool) []interface{} {
	oldByIndex := make(map[int]map[string]interface{}, len(previous))
	for _, raw := range previous {
		old := raw.(map[string]interface{})
		oldByIndex[formatInt(old["tunnel_index"])] = vpnAttachmentExpandTunnel(old)
	}
	result := make([]interface{}, 0, len(planned))
	for _, raw := range vpnAttachmentSortedTunnels(planned) {
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
// is enabled.
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
