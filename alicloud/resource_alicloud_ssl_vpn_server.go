package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSslVpnServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSslVpnServerCreate,
		Read:   resourceAliyunSslVpnServerRead,
		Update: resourceAliyunSslVpnServerUpdate,
		Delete: resourceAliyunSslVpnServerDelete,

		Schema: map[string]*schema.Schema{
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

			"client_ip_pool": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"local_subnet": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_UDP_PROTO,
				ValidateFunc: validateAllowedStringValue([]string{VPN_UDP_PROTO, VPN_TCP_PROTO}),
			},

			"cipher": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SSL_VPN_ENC_AES_128,
				ValidateFunc: validateAllowedStringValue([]string{SSL_VPN_ENC_AES_128, SSL_VPN_ENC_AES_192, SSL_VPN_ENC_AES_256, SSL_VPN_ENC_NONE}),
			},
			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1194,
				ValidateFunc: validateSslVpnPortValue([]int{22, 2222, 22222, 9000, 9001, 9002, 7505, 80, 443, 53, 68, 123, 4510, 4560, 500, 4500}),
			},
			"compress": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"connections": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_connections": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"internet_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSslVpnServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var sslVpnServer *vpc.CreateSslVpnServerResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := buildAliyunSslVpnServerArgs(d, meta)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSslVpnServer(args)
		})
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create SslVpnServer error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		sslVpnServer, _ = raw.(*vpc.CreateSslVpnServerResponse)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create ssl vpn server got an error :%#v", err)
	}

	d.SetId(sslVpnServer.SslVpnServerId)

	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	resp, err := vpnGatewayService.DescribeSslVpnServer(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vpn_gateway_id", resp.VpnGatewayId)
	d.Set("name", resp.Name)
	d.Set("local_subnet", resp.LocalSubnet)
	d.Set("client_ip_pool", resp.ClientIpPool)
	d.Set("cipher", resp.Cipher)
	d.Set("protocol", resp.Proto)
	d.Set("port", resp.Port)
	d.Set("compress", resp.Compress)
	d.Set("connections", resp.Connections)
	d.Set("max_connections", resp.MaxConnections)
	d.Set("internet_ip", resp.InternetIp)

	return nil
}

func resourceAliyunSslVpnServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	attributeUpdate := false
	request := vpc.CreateModifySslVpnServerRequest()
	request.SslVpnServerId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("client_ip_pool") {
		request.ClientIpPool = d.Get("client_ip_pool").(string)
		attributeUpdate = true
	}

	if d.HasChange("local_subnet") {
		request.LocalSubnet = d.Get("local_subnet").(string)
		attributeUpdate = true
	}

	if d.HasChange("protocol") {
		request.Proto = d.Get("protocol").(string)
		attributeUpdate = true
	}

	if d.HasChange("cipher") {
		request.Cipher = d.Get("cipher").(string)
		attributeUpdate = true
	}

	if d.HasChange("port") {
		request.Port = requests.NewInteger(d.Get("port").(int))
		attributeUpdate = true
	}

	if d.HasChange("compress") {
		request.Compress = requests.NewBoolean(d.Get("compress").(bool))
		attributeUpdate = true
	}

	if attributeUpdate {

		res := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ModifySslVpnServer(request)
			})

			if err != nil {
				if IsExceptedError(err, VpnConfiguring) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(fmt.Errorf("Modify SslVpnServer error: %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modify SslVpnServer timeout and got an error: %#v.", err))
			}
			return nil
		})

		if res != nil {
			return fmt.Errorf("Modify SslVpnServer failed: %#v", res)
		}
	}
	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteSslVpnServerRequest()
	request.SslVpnServerId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnServer(request)
		})

		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete SslVpnServer error: %#v", err))
			}

			if IsExceptedError(err, SslVpnServerNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Delete SslVpnServer timeout and got an error: %#v.", err))
		}

		if _, err := vpnGatewayService.DescribeSslVpnServer(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunSslVpnServerArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateSslVpnServerRequest {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateCreateSslVpnServerRequest()
	request.RegionId = string(client.Region)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.ClientIpPool = d.Get("client_ip_pool").(string)
	request.LocalSubnet = d.Get("local_subnet").(string)
	request.Name = d.Get("name").(string)

	if v := d.Get("protocol").(string); v != "" {
		request.Proto = v
	}

	if v := d.Get("cipher").(string); v != "" {
		request.Cipher = v
	}

	if v, ok := d.GetOk("port"); ok && v.(int) != 0 {
		request.Port = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("compress"); ok {
		request.Compress = requests.NewBoolean(v.(bool))
	}

	return request
}
