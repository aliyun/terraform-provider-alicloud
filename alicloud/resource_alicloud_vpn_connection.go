package alicloud

import (
	"fmt"
	"time"

	"log"

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"local_subnet": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"remote_subnet": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"effect_immediately": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ike_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ipsec_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			/*"ike_config": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_enc_alg": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_auth_alg": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_pfs": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_lifetime": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_local_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_remote_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ipsec_config": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_auth_alg": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_pfs": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_lifetime": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},*/
		},
	}
}

func resourceAliyunVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	var vpnConn *vpc.CreateVpnConnectionResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunVpnConnectionArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateVpcRequest got an error: %#v", err))
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

	//resp, err := client.DescribeVpc(d.Id())
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
	d.Set("local_subnet", resp.LocalSubnet)
	d.Set("remote_subnet", resp.RemoteSubnet)
	d.Set("status", resp.Status)
	d.Set("ike_config", resp.IkeConfig)
	d.Set("ipsec_config", resp.IpsecConfig)
	return nil
}

func resourceAliyunVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {

	attributeUpdate := false
	request := vpc.CreateModifyVpnConnectionAttributeRequest()
	request.VpnConnectionId = d.Id()
	request.LocalSubnet = d.Get("local_subnet").(string)
	request.RemoteSubnet = d.Get("remote_subnet").(string)

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("local_subnet") || d.HasChange("remote_subnet") {
		attributeUpdate = true
	}

	if d.HasChange("effect_immediately") {

		request.EffectImmediately = requests.NewBoolean(d.Get("effect_immediately").(bool))
		attributeUpdate = true
	}

	if d.HasChange("ike_config") {
		request.IkeConfig = d.Get("ike_config").(string)
		attributeUpdate = true
	}

	if d.HasChange("ipsec_config") {
		request.IpsecConfig = d.Get("ipsec_config").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
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
	//log.Printf("Try to delete vpn connection")
	log.Printf("[ERROR]  testlog %s: Try to delete vpn connection", d.Id())
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteVpnConnection(request)

		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create vpn connnection: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Delete VPN connection timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeVpnConnection(d.Id()); err != nil {
			if IsExceptedError(err, VpnConnNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	request := vpc.CreateCreateVpnConnectionRequest()
	request.RegionId = string(getRegion(d, meta))
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = d.Get("local_subnet").(string)
	request.RemoteSubnet = d.Get("remote_subnet").(string)

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v := d.Get("ike_config").(string); v != "" {
		request.IkeConfig = v
	}

	if v := d.Get("ipsec_config").(string); v != "" {
		request.IpsecConfig = v
	}
	return request, nil
}
