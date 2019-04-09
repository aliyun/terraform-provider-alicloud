package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_VERSION_1,
							ValidateFunc: validateAllowedStringValue([]string{IKE_VERSION_1, IKE_VERSION_2}),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_MODE_MAIN,
							ValidateFunc: validateAllowedStringValue([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}),
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validateAllowedStringValue([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}),
						},
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validateAllowedStringValue([]string{VPN_AUTH_SHA, VPN_AUTH_MD5}),
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validateAllowedStringValue([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validateIntegerInRange(0, 86400),
						},
						"ike_local_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
						"ike_remote_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
					},
				},
			},

			"ipsec_config": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validateAllowedStringValue([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}),
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validateAllowedStringValue([]string{VPN_AUTH_SHA, VPN_AUTH_MD5}),
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validateAllowedStringValue([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 86400),
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
	request, err := buildAliyunVpnConnectionArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var vpnConn *vpc.CreateVpnConnectionResponse
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpnConnection(&args)
		})
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		vpnConn, _ = raw.(*vpc.CreateVpnConnectionResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(vpnConn.VpnConnectionId)

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	resp, err := vpnGatewayService.DescribeVpnConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("customer_gateway_id", resp.CustomerGatewayId)
	d.Set("vpn_gateway_id", resp.VpnGatewayId)
	d.Set("name", resp.Name)

	localSubnet := strings.Split(resp.LocalSubnet, ",")
	localSubnetSet := make([]string, 0, len(localSubnet))
	for _, e := range localSubnet {
		localSubnetSet = append(localSubnet, e)
	}
	d.Set("local_subnet", localSubnetSet)

	remoteSubnet := strings.Split(resp.RemoteSubnet, ",")
	remoteSubnetSet := make([]string, 0, len(remoteSubnet))
	for _, e := range remoteSubnet {
		remoteSubnetSet = append(remoteSubnet, e)
	}
	d.Set("remote_subnet", remoteSubnetSet)

	d.Set("effect_immediately", resp.EffectImmediately)
	d.Set("status", resp.Status)

	if err := d.Set("ike_config", vpnGatewayService.ParseIkeConfig(resp.IkeConfig)); err != nil {
		return err
	}

	if err := d.Set("ipsec_config", vpnGatewayService.ParseIpsecConfig(resp.IpsecConfig)); err != nil {
		return err
	}

	return nil
}

func resourceAliyunVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	attributeUpdate := false
	request := vpc.CreateModifyVpnConnectionAttributeRequest()
	request.VpnConnectionId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("local_subnet") {
		request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
		attributeUpdate = true
	}

	if d.HasChange("remote_subnet") {
		request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())
		attributeUpdate = true
	}

	/* If not set effect_immediately value, VPN connection will automatically set the value to false*/
	if d.HasChange("effect_immediately") {
		attributeUpdate = true
		/*The value will be read below*/
	}

	if d.HasChange("ike_config") {
		ike_config, err := vpnGatewayService.AssembleIkeConfig(d.Get("ike_config").([]interface{}))
		if err != nil {
			return fmt.Errorf("wrong ike_config: %#v", err)
		}
		request.IkeConfig = ike_config
		attributeUpdate = true
	}

	if d.HasChange("ipsec_config") {
		ipsec_config, err := vpnGatewayService.AssembleIpsecConfig(d.Get("ipsec_config").([]interface{}))
		if err != nil {
			return fmt.Errorf("wrong ipsec_config: %#v", err)
		}
		request.IpsecConfig = ipsec_config
		attributeUpdate = true
	}

	if attributeUpdate {
		if v, ok := d.GetOk("effect_immediately"); ok {
			request.EffectImmediately = requests.NewBoolean(v.(bool))
		}

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnConnectionAttribute(request)
		})
		if err != nil {
			return err
		}
	}

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteVpnConnectionRequest()
	request.VpnConnectionId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnConnection(request)
		})

		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete VPN connnection: %#v", err))
			}
			if IsExceptedError(err, VpnConnNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Delete VPN connection timeout and got an error: %#v.", err))
		}

		if _, err := vpnGatewayService.DescribeVpnConnection(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateCreateVpnConnectionRequest()
	request.RegionId = client.RegionId
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig, err := vpnGatewayService.AssembleIkeConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ike_config: %#v", err)
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
