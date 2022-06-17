package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnConnectionCreate,
		Read:   resourceAliyunVpnConnectionRead,
		Update: resourceAliyunVpnConnectionUpdate,
		Delete: resourceAliyunVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
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

			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ike_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_VERSION_1,
							ValidateFunc: validation.StringInSlice([]string{IKE_VERSION_1, IKE_VERSION_2}, false),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_MODE_MAIN,
							ValidateFunc: validation.StringInSlice([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}, false),
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}, false),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"ike_local_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_remote_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
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
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24, VPN_PFS_DISABLED}, false),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
					},
				},
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
					},
				},
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
							Type:     schema.TypeString,
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

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateVpnConnection"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	request["CustomerGatewayId"] = d.Get("customer_gateway_id")
	request["VpnGatewayId"] = d.Get("vpn_gateway_id")

	if v, ok := d.GetOk("local_subnet"); ok {
		request["LocalSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("remote_subnet"); ok {
		request["RemoteSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOkExists("effect_immediately"); ok {
		request["EffectImmediately"] = v
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfigsArg := v.([]interface{})[0].(map[string]interface{})
		ikeConfigsMap := map[string]interface{}{
			"IkeAuthAlg":  ikeConfigsArg["ike_auth_alg"],
			"IkeEncAlg":   ikeConfigsArg["ike_enc_alg"],
			"IkeLifetime": ikeConfigsArg["ike_lifetime"],
			"LocalId":     ikeConfigsArg["ike_local_id"],
			"IkeMode":     ikeConfigsArg["ike_mode"],
			"IkePfs":      ikeConfigsArg["ike_pfs"],
			"RemoteId":    ikeConfigsArg["ike_remote_id"],
			"IkeVersion":  ikeConfigsArg["ike_version"],
			"Psk":         ikeConfigsArg["psk"],
		}
		ikeConfigsMapsStrting, _ := convertMaptoJsonString(ikeConfigsMap)
		request["IkeConfig"] = ikeConfigsMapsStrting
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecsArg := v.([]interface{})[0].(map[string]interface{})
		ipsecsMap := map[string]interface{}{
			"IpsecAuthAlg":  ipsecsArg["ipsec_auth_alg"],
			"IpsecEncAlg":   ipsecsArg["ipsec_enc_alg"],
			"IpsecLifetime": ipsecsArg["ipsec_lifetime"],
			"IpsecPfs":      ipsecsArg["ipsec_pfs"],
		}
		ipsecsMapsStrting, _ := convertMaptoJsonString(ipsecsMap)
		request["IpsecConfig"] = ipsecsMapsStrting
	}

	if v, ok := d.GetOk("bgp_config"); ok {
		bgpsArg := v.([]interface{})[0].(map[string]interface{})
		bgpsMap := map[string]interface{}{
			"EnableBgp":  bgpsArg["enable"],
			"LocalAsn":   bgpsArg["local_asn"],
			"TunnelCidr": bgpsArg["tunnel_cidr"],
			"LocalBgpIp": bgpsArg["local_bgp_ip"],
		}
		bgpsMapsStrting, _ := convertMaptoJsonString(bgpsMap)
		request["BgpConfig "] = bgpsMapsStrting
	}

	if v, ok := d.GetOk("health_check_config"); ok {
		healthChecksArg := v.([]interface{})[0].(map[string]interface{})
		healthChecksMap := map[string]interface{}{
			"enable":   healthChecksArg["enable"],
			"dip":      healthChecksArg["dip"],
			"sip":      healthChecksArg["sip"],
			"interval": formatInt(healthChecksArg["interval"]),
			"retry":    formatInt(healthChecksArg["retry"]),
		}

		healthChecksMapsStrting, _ := convertMaptoJsonString(healthChecksMap)
		request["HealthCheckConfig "] = healthChecksMapsStrting
	}

	if v, ok := d.GetOkExists("enable_dpd"); ok {
		request["EnableDpd"] = v
	}

	if v, ok := d.GetOkExists("enable_nat_traversal"); ok {
		request["EnableNatTraversal"] = v
	}

	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
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
	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway vpcService.DescribeVpnConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("customer_gateway_id", object["CustomerGatewayId"])
	d.Set("vpn_gateway_id", object["VpnGatewayId"])
	d.Set("name", object["Name"])

	localSubnet := strings.Split(object["LocalSubnet"].(string), ",")
	d.Set("local_subnet", localSubnet)

	remoteSubnet := strings.Split(object["RemoteSubnet"].(string), ",")
	d.Set("remote_subnet", remoteSubnet)

	d.Set("effect_immediately", object["EffectImmediately"])
	d.Set("status", object["Status"])

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

	if ipsecConfig, ok := object["IkeConfig"]; ok {
		ikeConfig := ipsecConfig.(map[string]interface{})
		ipsecConfigMaps := make([]map[string]interface{}, 0)
		ipsecConfigMaps = append(ipsecConfigMaps,
			map[string]interface{}{
				"ike_auth_alg":  ikeConfig["IkeAuthAlg"],
				"ike_enc_alg":   ikeConfig["IkeEncAlg"],
				"ike_lifetime":  ikeConfig["IkeLifetime"],
				"ike_local_id":  ikeConfig["LocalId"],
				"ike_mode":      ikeConfig["IkeMode"],
				"ike_pfs":       ikeConfig["IkePfs"],
				"ike_remote_id": ikeConfig["RemoteId"],
				"ike_version":   ikeConfig["IkeVersion"],
				"psk":           ikeConfig["Psk"],
			})
		d.Set("ike_config", ipsecConfigMaps)
	}

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
			})
		d.Set("health_check_config", healthChecksMaps)
	}

	d.Set("enable_dpd", object["EnableDpd"])
	d.Set("enable_nat_traversal", object["EnableNatTraversal"])

	return nil
}

func resourceAliyunVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ModifyVpnConnectionAttribute"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"VpnConnectionId": d.Id(),
	}

	update := false
	if d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("local_subnet") {
		update = true
	}
	if v, ok := d.GetOk("local_subnet"); ok {
		request["LocalSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if d.HasChange("remote_subnet") {
		update = true
	}
	if v, ok := d.GetOk("remote_subnet"); ok {
		request["RemoteSubnet"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if d.HasChange("effect_immediately") {
		update = true
		if v, ok := d.GetOkExists("effect_immediately"); ok {
			request["EffectImmediately"] = v
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
				"LocalId":     ikeConfigsArg["ike_local_id"],
				"IkeMode":     ikeConfigsArg["ike_mode"],
				"IkePfs":      ikeConfigsArg["ike_pfs"],
				"RemoteId":    ikeConfigsArg["ike_remote_id"],
				"IkeVersion":  ikeConfigsArg["ike_version"],
				"Psk":         ikeConfigsArg["psk"],
			}
			ikeConfigsMapsStrting, _ := convertMaptoJsonString(ikeConfigsMap)
			request["IkeConfig"] = ikeConfigsMapsStrting
		}
	}

	if d.HasChange("ipsec_config") {
		update = true
		if v, ok := d.GetOk("ipsec_config"); ok {
			ipsecsArg := v.([]interface{})[0].(map[string]interface{})
			ipsecsMap := map[string]interface{}{
				"IpsecAuthAlg":  ipsecsArg["ipsec_auth_alg"],
				"IpsecEncAlg":   ipsecsArg["ipsec_enc_alg"],
				"IpsecLifetime": ipsecsArg["ipsec_lifetime"],
				"IpsecPfs":      ipsecsArg["ipsec_pfs"],
			}
			ipsecsMapsStrting, _ := convertMaptoJsonString(ipsecsMap)
			request["IpsecConfig"] = ipsecsMapsStrting
		}
	}

	if d.HasChange("bgp_config") {
		update = true
		if v, ok := d.GetOk("bgp_config"); ok {
			bgpsArg := v.([]interface{})[0].(map[string]interface{})
			bgpsMap := map[string]interface{}{
				"EnableBgp":  bgpsArg["enable"],
				"LocalAsn":   bgpsArg["local_asn"],
				"TunnelCidr": bgpsArg["tunnel_cidr"],
				"LocalBgpIp": bgpsArg["local_bgp_ip"],
			}
			bgpsMapsStrting, _ := convertMaptoJsonString(bgpsMap)
			request["BgpConfig "] = bgpsMapsStrting
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
				"interval": healthChecksArg["interval"],
				"retry":    healthChecksArg["retry"],
			}
			healthChecksMapsStrting, _ := convertMaptoJsonString(healthChecksMap)
			request["HealthCheckConfig "] = healthChecksMapsStrting
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

	if update {
		request["ClientToken"] = buildClientToken(action)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
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
	}

	return resourceAliyunVpnConnectionRead(d, meta)

}

func resourceAliyunVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpnConnection"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"VpnConnectionId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func buildAliyunVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateCreateVpnConnectionRequest()
	request.RegionId = client.RegionId
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").([]interface{}))
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").([]interface{}))

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig, err := vpnGatewayService.AssembleIkeConfig(v.([]interface{}))
		if err != nil {
			return nil, WrapError(err)
		}
		request.IkeConfig = ikeConfig
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig, err := vpnGatewayService.AssembleIpsecConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ipsec_config: %#v", err)
		}
		request.IpsecConfig = ipsecConfig
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
