package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpnGatewayVpnAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpnGatewayVpnAttachmentCreate,
		Read:   resourceAlicloudVpnGatewayVpnAttachmentRead,
		Update: resourceAlicloudVpnGatewayVpnAttachmentUpdate,
		Delete: resourceAlicloudVpnGatewayVpnAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_config": {
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
						"local_bgp_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_dpd": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"enable_nat_traversal": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
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
						"retry": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ike_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}, false),
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}, false),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{IKE_VERSION_1, IKE_VERSION_2}, false),
						},
						"local_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"psk": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"remote_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
					},
				},
			},
			"ipsec_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24, VPN_PFS_DISABLED}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"remote_subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpnGatewayVpnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVpnAttachment"
	request := make(map[string]interface{})
	var err error
	request["CustomerGatewayId"] = d.Get("customer_gateway_id")
	if v, ok := d.GetOkExists("effect_immediately"); ok {
		request["EffectImmediately"] = v
	}
	if v, ok := d.GetOkExists("enable_dpd"); ok {
		request["EnableDpd"] = v
	}
	if v, ok := d.GetOkExists("enable_nat_traversal"); ok {
		request["EnableNatTraversal"] = v
	}
	request["LocalSubnet"] = d.Get("local_subnet")
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	request["RegionId"] = client.RegionId
	request["RemoteSubnet"] = d.Get("remote_subnet")
	if v, ok := d.GetOk("vpn_attachment_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfigsArg := v.([]interface{})[0].(map[string]interface{})
		ikeConfigsMap := map[string]interface{}{
			"IkeAuthAlg":  ikeConfigsArg["ike_auth_alg"],
			"IkeEncAlg":   ikeConfigsArg["ike_enc_alg"],
			"IkeLifetime": ikeConfigsArg["ike_lifetime"],
			"LocalId":     ikeConfigsArg["local_id"],
			"IkeMode":     ikeConfigsArg["ike_mode"],
			"IkePfs":      ikeConfigsArg["ike_pfs"],
			"RemoteId":    ikeConfigsArg["remote_id"],
			"IkeVersion":  ikeConfigsArg["ike_version"],
			"Psk":         ikeConfigsArg["psk"],
		}
		ikeConfigsMapsString, _ := convertMaptoJsonString(ikeConfigsMap)
		request["IkeConfig"] = ikeConfigsMapsString
	}
	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecArg := v.([]interface{})[0].(map[string]interface{})
		ipsecMap := map[string]interface{}{
			"IpsecAuthAlg":  ipsecArg["ipsec_auth_alg"],
			"IpsecEncAlg":   ipsecArg["ipsec_enc_alg"],
			"IpsecLifetime": ipsecArg["ipsec_lifetime"],
			"IpsecPfs":      ipsecArg["ipsec_pfs"],
		}
		ipsecMapsString, _ := convertMaptoJsonString(ipsecMap)
		request["IpsecConfig"] = ipsecMapsString
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		healthChecksArg := v.([]interface{})[0].(map[string]interface{})
		healthChecksMap := map[string]interface{}{
			"enable":   healthChecksArg["enable"],
			"dip":      healthChecksArg["dip"],
			"sip":      healthChecksArg["sip"],
			"interval": formatInt(healthChecksArg["interval"]),
			"retry":    formatInt(healthChecksArg["retry"]),
			"Policy":   healthChecksArg["policy"],
		}

		healthChecksMapsString, _ := convertMaptoJsonString(healthChecksMap)
		request["HealthCheckConfig"] = healthChecksMapsString
	}
	if v, ok := d.GetOk("bgp_config"); ok {
		bgpArg := v.([]interface{})[0].(map[string]interface{})
		bgpMap := map[string]interface{}{
			"EnableBgp":  bgpArg["enable"],
			"LocalAsn":   bgpArg["local_asn"],
			"TunnelCidr": bgpArg["tunnel_cidr"],
			"LocalBgpIp": bgpArg["local_bgp_ip"],
		}
		bgpMapsString, _ := convertMaptoJsonString(bgpMap)
		request["BgpConfig"] = bgpMapsString
	}
	request["ClientToken"] = buildClientToken("CreateVpnAttachment")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway_vpn_attachment", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["VpnConnectionId"]))

	return resourceAlicloudVpnGatewayVpnAttachmentRead(d, meta)
}
func resourceAlicloudVpnGatewayVpnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnGatewayVpnAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway_vpn_attachment vpcService.DescribeVpnGatewayVpnAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("effect_immediately", object["EffectImmediately"])
	d.Set("customer_gateway_id", object["CustomerGatewayId"])
	d.Set("enable_dpd", object["EnableDpd"])
	d.Set("enable_nat_traversal", object["EnableNatTraversal"])

	if ipsecConfig, ok := object["VpnBgpConfig"]; ok {
		bgpConfig := ipsecConfig.(map[string]interface{})
		bgpConfigMaps := make([]map[string]interface{}, 0)
		bgpConfigMaps = append(bgpConfigMaps, map[string]interface{}{
			"enable":       convertStringToBool(bgpConfig["EnableBgp"].(string)),
			"local_asn":    bgpConfig["LocalAsn"],
			"tunnel_cidr":  bgpConfig["TunnelCidr"],
			"local_bgp_ip": bgpConfig["LocalBgpIp"],
		})
		d.Set("bgp_config", bgpConfigMaps)
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
				"policy":   healthCheckConfig["Policy"],
			})
		d.Set("health_check_config", healthChecksMaps)
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
		d.Set("ike_config", ipsecConfigMaps)
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
		d.Set("ipsec_config", ipsecConfigMaps)
	}
	d.Set("local_subnet", object["LocalSubnet"])
	d.Set("network_type", object["NetworkType"])
	d.Set("remote_subnet", object["RemoteSubnet"])
	d.Set("status", object["State"])
	d.Set("vpn_attachment_name", object["Name"])
	d.Set("internet_ip", object["InternetIp"])
	return nil
}
func resourceAlicloudVpnGatewayVpnAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"VpnConnectionId": d.Id(),
		"RegionId":        client.RegionId,
	}
	if d.HasChange("effect_immediately") {
		update = true
		if v, ok := d.GetOkExists("effect_immediately"); ok {
			request["EffectImmediately"] = v
		}
	}
	if d.HasChange("enable_dpd") {
		update = true
		if v, ok := d.GetOkExists("enable_dpd"); ok {
			request["EnableDpd"] = v
		}
	}
	if d.HasChange("enable_nat_traversal") {
		update = true
		if v, ok := d.GetOkExists("enable_nat_traversal"); ok {
			request["EnableNatTraversal"] = v
		}
	}
	if d.HasChange("local_subnet") {
		update = true
	}
	request["LocalSubnet"] = d.Get("local_subnet")
	if d.HasChange("remote_subnet") {
		update = true
	}
	request["RemoteSubnet"] = d.Get("remote_subnet")
	if d.HasChange("vpn_attachment_name") {
		update = true
		if v, ok := d.GetOk("vpn_attachment_name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("ike_config") {
		update = true
		if v, ok := d.GetOk("ike_config"); ok {
			ikeConfigsArg := v.([]interface{})[0].(map[string]interface{})
			ikeConfigsMap := map[string]interface{}{
				"IkeAuthAlg":  ikeConfigsArg["ike_auth_alg"],
				"IkeEncAlg":   ikeConfigsArg["ike_enc_alg"],
				"IkeLifetime": ikeConfigsArg["ike_lifetime"],
				"LocalId":     ikeConfigsArg["local_id"],
				"IkeMode":     ikeConfigsArg["ike_mode"],
				"IkePfs":      ikeConfigsArg["ike_pfs"],
				"RemoteId":    ikeConfigsArg["remote_id"],
				"IkeVersion":  ikeConfigsArg["ike_version"],
				"Psk":         ikeConfigsArg["psk"],
			}
			ikeConfigsMapsString, _ := convertMaptoJsonString(ikeConfigsMap)
			request["IkeConfig"] = ikeConfigsMapsString
		}
	}
	if d.HasChange("ipsec_config") {
		update = true
		if v, ok := d.GetOk("ipsec_config"); ok {
			ipsecArg := v.([]interface{})[0].(map[string]interface{})
			ipsecMap := map[string]interface{}{
				"IpsecAuthAlg":  ipsecArg["ipsec_auth_alg"],
				"IpsecEncAlg":   ipsecArg["ipsec_enc_alg"],
				"IpsecLifetime": ipsecArg["ipsec_lifetime"],
				"IpsecPfs":      ipsecArg["ipsec_pfs"],
			}
			ipsecMapsString, _ := convertMaptoJsonString(ipsecMap)
			request["IpsecConfig"] = ipsecMapsString
		}
	}
	if d.HasChange("health_check_config") {
		update = true
		if v, ok := d.GetOk("health_check_config"); ok {
			healthChecksArg := v.([]interface{})[0].(map[string]interface{})
			healthChecksMap := map[string]interface{}{
				"enable":   healthChecksArg["enable"],
				"dip":      healthChecksArg["dip"],
				"sip":      healthChecksArg["sip"],
				"interval": formatInt(healthChecksArg["interval"]),
				"retry":    formatInt(healthChecksArg["retry"]),
				"Policy":   healthChecksArg["policy"],
			}

			healthChecksMapsString, _ := convertMaptoJsonString(healthChecksMap)
			request["HealthCheckConfig"] = healthChecksMapsString
		}
	}
	if d.HasChange("bgp_config") {
		update = true
		if v, ok := d.GetOk("bgp_config"); ok {
			bgpArg := v.([]interface{})[0].(map[string]interface{})
			bgpMap := map[string]interface{}{
				"EnableBgp":  bgpArg["enable"],
				"LocalAsn":   bgpArg["local_asn"],
				"TunnelCidr": bgpArg["tunnel_cidr"],
				"LocalBgpIp": bgpArg["local_bgp_ip"],
			}
			bgpMapsString, _ := convertMaptoJsonString(bgpMap)
			request["BgpConfig"] = bgpMapsString
		}
	}

	if d.HasChange("customer_gateway_id") {
		update = true
		if v, ok := d.GetOk("customer_gateway_id"); ok {
			request["CustomerGatewayId"] = v
		}
	}

	if update {
		action := "ModifyVpnAttachmentAttribute"
		request["ClientToken"] = buildClientToken("ModifyVpnAttachmentAttribute")
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"init", "attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpnGatewayVpnAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudVpnGatewayVpnAttachmentRead(d, meta)
}
func resourceAlicloudVpnGatewayVpnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpnAttachment"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{}

	request["RegionId"] = client.RegionId
	request["VpnConnectionId"] = d.Id()
	request["ClientToken"] = buildClientToken("DeleteVpnAttachment")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
