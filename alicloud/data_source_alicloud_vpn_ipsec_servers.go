package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnIpsecServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnIpsecServersRead,
		Schema: map[string]*schema.Schema{
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
			"ipsec_server_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip_pool": {
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
						"idaas_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
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
										Type:     schema.TypeInt,
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
									"remote_id": {
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
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ipsec_pfs": {
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
						"ipsec_server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipsec_server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"multi_factor_auth_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"online_client_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"psk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"psk_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnIpsecServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListIpsecServers"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("ipsec_server_name"); ok {
		request["IpsecServerName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		request["VpnGatewayId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var ipsecServerNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ipsecServerNameRegex = r
	}

	if v, ok := d.GetOk("ids"); ok {
		request["IpsecServerId"] = v
	}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpn_ipsec_servers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.IpsecServers", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.IpsecServers", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if ipsecServerNameRegex != nil && !ipsecServerNameRegex.MatchString(fmt.Sprint(item["IpsecServerName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"client_ip_pool":            object["ClientIpPool"],
			"create_time":               object["CreationTime"],
			"effect_immediately":        object["EffectImmediately"],
			"idaas_instance_id":         object["IDaaSInstanceId"],
			"internet_ip":               object["InternetIp"],
			"ipsec_config":              object["IpsecConfig"],
			"id":                        fmt.Sprint(object["IpsecServerId"]),
			"ipsec_server_id":           fmt.Sprint(object["IpsecServerId"]),
			"ipsec_server_name":         object["IpsecServerName"],
			"local_subnet":              object["LocalSubnet"],
			"max_connections":           formatInt(object["MaxConnections"]),
			"multi_factor_auth_enabled": object["MultiFactorAuthEnabled"],
			"online_client_count":       formatInt(object["OnlineClientCount"]),
			"psk":                       object["Psk"],
			"psk_enabled":               object["PskEnabled"],
			"vpn_gateway_id":            object["VpnGatewayId"],
		}

		ikeConfigSli := make([]map[string]interface{}, 0)
		if len(object["IkeConfig"].(map[string]interface{})) > 0 {
			ikeConfig := object["IkeConfig"]
			ikeConfigMap := make(map[string]interface{})
			if ikeConfigArg, ok := ikeConfig.(map[string]interface{}); ok {
				ikeConfigMap["ike_auth_alg"] = ikeConfigArg["IkeAuthAlg"]
				ikeConfigMap["ike_enc_alg"] = ikeConfigArg["IkeEncAlg"]
				ikeConfigMap["ike_lifetime"] = formatInt(ikeConfigArg["IkeLifetime"])
				ikeConfigMap["ike_mode"] = ikeConfigArg["IkeMode"]
				ikeConfigMap["ike_pfs"] = ikeConfigArg["IkePfs"]
				ikeConfigMap["ike_version"] = ikeConfigArg["IkeVersion"]
				ikeConfigMap["local_id"] = ikeConfigArg["LocalId"]
				ikeConfigMap["remote_id"] = ikeConfigArg["RemoteId"]
				ikeConfigSli = append(ikeConfigSli, ikeConfigMap)
			}

		}
		mapping["ike_config"] = ikeConfigSli

		ipsecConfigSli := make([]map[string]interface{}, 0)
		if len(object["IpsecConfig"].(map[string]interface{})) > 0 {
			ipsecConfig := object["IpsecConfig"]
			ipsecConfigMap := make(map[string]interface{})
			if ipsecConfigArg, ok := ipsecConfig.(map[string]interface{}); ok {
				ipsecConfigMap["ipsec_auth_alg"] = ipsecConfigArg["IpsecAuthAlg"]
				ipsecConfigMap["ipsec_enc_alg"] = ipsecConfigArg["IpsecEncAlg"]
				ipsecConfigMap["ipsec_lifetime"] = formatInt(ipsecConfigArg["IpsecLifetime"])
				ipsecConfigMap["ipsec_pfs"] = ipsecConfigArg["IpsecPfs"]
				ipsecConfigSli = append(ipsecConfigSli, ipsecConfigMap)
			}

		}
		mapping["ipsec_config"] = ipsecConfigSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["IpsecServerName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("servers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
