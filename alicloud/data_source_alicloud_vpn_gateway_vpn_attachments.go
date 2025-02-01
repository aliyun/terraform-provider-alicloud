package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnGatewayVpnAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnGatewayVpnAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"init", "active", "attaching", "attached", "detaching", "financialLocked", "provisioning", "updating", "upgrading", "deleted"}, false),
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"vpn_gateway_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "The parameter 'vpn_gateway_id' has been deprecated from 1.194.0.",
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bgp_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_asn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_bgp_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_cidr": {
										Type:     schema.TypeString,
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
							Computed: true,
						},
						"effect_immediately": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
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
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ike_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ike_auth_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_enc_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_lifetime": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_pfs": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ike_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_id": {
										Type:     schema.TypeString,
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
								},
							},
						},
						"ipsec_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_auth_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_enc_alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_lifetime": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipsec_pfs": {
										Type:     schema.TypeString,
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnGatewayVpnAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeVpnConnections"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		request["VpnGatewayId"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)
	var vpnAttachmentNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpnAttachmentNameRegex = r
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
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpn_gateway_vpn_attachments", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VpnConnections.VpnConnection", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VpnConnections.VpnConnection", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if vpnAttachmentNameRegex != nil && !vpnAttachmentNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VpnConnectionId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"customer_gateway_id": object["CustomerGatewayId"],
			"effect_immediately":  object["EffectImmediately"],
			"local_subnet":        object["LocalSubnet"],
			"network_type":        object["NetworkType"],
			"remote_subnet":       object["RemoteSubnet"],
			"status":              object["State"],
			"connection_status":   object["Status"],
			"vpn_attachment_name": object["Name"],
			"id":                  fmt.Sprint(object["VpnConnectionId"]),
			"vpn_connection_id":   fmt.Sprint(object["VpnConnectionId"]),
			"create_time":         object["CreateTime"],
			"internet_ip":         object["InternetIp"],
		}

		if ipsecConfig, ok := object["VpnBgpConfig"]; ok {
			bgpConfig := ipsecConfig.(map[string]interface{})
			bgpConfigMaps := make([]map[string]interface{}, 0)
			bgpConfigMaps = append(bgpConfigMaps, map[string]interface{}{
				"status":       bgpConfig["Status"],
				"local_asn":    bgpConfig["LocalAsn"],
				"tunnel_cidr":  bgpConfig["TunnelCidr"],
				"local_bgp_ip": bgpConfig["LocalBgpIp"],
			})
			mapping["bgp_config"] = bgpConfigMaps
		}

		if ipsecConfig, ok := object["VcoHealthCheck"]; ok {
			healthCheckConfig := ipsecConfig.(map[string]interface{})
			healthChecksMaps := make([]map[string]interface{}, 0)
			healthChecksMaps = append(healthChecksMaps,
				map[string]interface{}{
					"enable":   convertStringToBool(healthCheckConfig["Enable"].(string)),
					"dip":      healthCheckConfig["Dip"],
					"sip":      healthCheckConfig["Sip"],
					"interval": formatInt(healthCheckConfig["Interval"]),
					"retry":    formatInt(healthCheckConfig["Retry"]),
					"status":   healthCheckConfig["Status"],
					"policy":   healthCheckConfig["Policy"],
				})
			mapping["health_check_config"] = healthChecksMaps
		}

		if ipsecConfig, ok := object["IkeConfig"]; ok {
			ikeConfig := ipsecConfig.(map[string]interface{})
			ipsecConfigMaps := make([]map[string]interface{}, 0)
			ipsecConfigMaps = append(ipsecConfigMaps,
				map[string]interface{}{
					"ike_auth_alg": ikeConfig["IkeAuthAlg"],
					"ike_enc_alg":  ikeConfig["IkeEncAlg"],
					"ike_lifetime": ikeConfig["IkeLifetime"],
					"local_id":     ikeConfig["LocalId"],
					"ike_mode":     ikeConfig["IkeMode"],
					"ike_pfs":      ikeConfig["IkePfs"],
					"remote_id":    ikeConfig["RemoteId"],
					"ike_version":  ikeConfig["IkeVersion"],
					"psk":          ikeConfig["Psk"],
				})
			mapping["ike_config"] = ipsecConfigMaps
		}

		if ipsecConfig, ok := object["IpsecConfig"]; ok {
			ipsecConfigArg := ipsecConfig.(map[string]interface{})
			ipsecConfigMaps := make([]map[string]interface{}, 0)

			ipsecConfigMaps = append(ipsecConfigMaps,
				map[string]interface{}{
					"ipsec_auth_alg": ipsecConfigArg["IpsecAuthAlg"],
					"ipsec_enc_alg":  ipsecConfigArg["IpsecEncAlg"],
					"ipsec_lifetime": ipsecConfigArg["IpsecLifetime"],
					"ipsec_pfs":      ipsecConfigArg["IpsecPfs"],
				})
			mapping["ipsec_config"] = ipsecConfigMaps
		}

		ids = append(ids, fmt.Sprint(object["VpnConnectionId"].(string)))
		names = append(names, object["Name"])
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
