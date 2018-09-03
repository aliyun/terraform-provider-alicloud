package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnConnectionCreate,
		Read:   resourceAliyunVpnConnectionRead,
		Update: resourceAliyunVpnConnectionUpdate,
		Delete: resourceAliyunVpnConnectionDelete,

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},

			"local_subnet": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"remote_subnet": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"effect_immediately": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ike_config": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
						"ike_version": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_VERSION_1,
							ValidateFunc: validateAllowedStringValue([]string{IKE_VERSION_1, IKE_VERSION_2}),
						},
						"ike_mode": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_MODE_MAIN,
							ValidateFunc: validateAllowedStringValue([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}),
						},
						"ike_enc_alg": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validateAllowedStringValue([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}),
						},
						"ike_auth_alg": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validateAllowedStringValue([]string{VPN_AUTH_SHA, VPN_AUTH_MD5}),
						},
						"ike_pfs": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validateAllowedStringValue([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}),
						},
						"ike_lifetime": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validateIntegerInRange(0, 86400),
						},
						"ike_local_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
						"ike_remote_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
						},
					},
				},
			},

			"ipsec_config": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validateAllowedStringValue([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}),
						},
						"ipsec_auth_alg": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validateAllowedStringValue([]string{VPN_AUTH_SHA, VPN_AUTH_MD5}),
						},
						"ipsec_pfs": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validateAllowedStringValue([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}),
						},
						"ipsec_lifetime": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 86400),
						},
					},
				},
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	var vpnConn *vpc.CreateVpnConnectionResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunVpnConnectionArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building buildAliyunVpnConnectionArgs got an error: %#v", err))
		}

		resp, err := client.vpcconn.CreateVpnConnection(args)
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create vpn connnection: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		vpnConn = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create VPN Connection got an error :%#v", err)
	}

	d.SetId(vpnConn.VpnConnectionId)

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	resp, err := client.DescribeVpnConnection(d.Id())
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

	if err := d.Set("ike_config", ParseIkeConfig(resp.IkeConfig)); err != nil {
		return err
	}

	if err := d.Set("ipsec_config", ParseIpsecConfig(resp.IpsecConfig)); err != nil {
		return err
	}

	return nil
}

func resourceAliyunVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {

	attributeUpdate := false
	request := vpc.CreateModifyVpnConnectionAttributeRequest()
	request.VpnConnectionId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("local_subnet") {
		request.LocalSubnet = AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
		attributeUpdate = true
	}

	if d.HasChange("remote_subnet") {
		request.RemoteSubnet = AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())
		attributeUpdate = true
	}

	/* If not set effect_immediately value, VPN connection will automatically set the value to false*/
	if d.HasChange("effect_immediately") {
		attributeUpdate = true
		/*The value will be read below*/
	}

	if d.HasChange("ike_config") {
		ike_config, err := AssembleIkeConfig(d.Get("ike_config").([]interface{}))
		if err != nil {
			return fmt.Errorf("wrong ike_config: %#v", err)
		}
		request.IkeConfig = ike_config
		attributeUpdate = true
	}

	if d.HasChange("ipsec_config") {
		ipsec_config, err := AssembleIpsecConfig(d.Get("ipsec_config").([]interface{}))
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

		if _, err := meta.(*AliyunClient).vpcconn.ModifyVpnConnectionAttribute(request); err != nil {
			return err
		}
	}

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := vpc.CreateDeleteVpnConnectionRequest()
	request.VpnConnectionId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteVpnConnection(request)

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

		if _, err := client.DescribeVpnConnection(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	request := vpc.CreateCreateVpnConnectionRequest()
	request.RegionId = getRegionId(d, meta)
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig, err := AssembleIkeConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ike_config: %#v", err)
		}
		request.IkeConfig = ikeConfig
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig, err := AssembleIpsecConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ipsec_config: %#v", err)
		}
		request.IpsecConfig = ipsecConfig
	}

	return request, nil
}
