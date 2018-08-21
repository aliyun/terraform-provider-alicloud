package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSslVpnServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSslVpnServerCreate,
		Read:   resourceAliyunSslVpnServerRead,
		Update: resourceAliyunSslVpnServerUpdate,
		Delete: resourceAliyunSslVpnServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceName,
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

			"proto": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UDP",
			},

			"cipher": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AES-128-CBC",
			},
			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateInstancePort,
				Default:      1194,
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
	client := meta.(*AliyunClient)
	var sslVpnServer *vpc.CreateSslVpnServerResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunSslVpnServerArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateSslVpnServerRequest got an error: %#v", err))
		}

		resp, err := client.vpcconn.CreateSslVpnServer(args)
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create SslVpnServer error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		sslVpnServer = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create ssl vpn server got an error :%#v", err)
	}

	d.SetId(sslVpnServer.SslVpnServerId)

	/*err = client.WaitForVpc(d.Id(), Available, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForVpcAvailable")
	}*/

	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	resp, err := client.DescribeSslVpnServers("", d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	sslVpnServer := resp.SslVpnServers.SslVpnServer[0]

	d.Set("ssl_vpn_server_id", sslVpnServer.SslVpnServerId)
	d.Set("vpn_gateway_id", sslVpnServer.VpnGatewayId)
	d.Set("name", sslVpnServer.Name)
	d.Set("local_subnet", sslVpnServer.LocalSubnet)
	d.Set("client_ip_pool", sslVpnServer.ClientIpPool)
	d.Set("cipher", sslVpnServer.Cipher)
	d.Set("proto", sslVpnServer.Proto)
	d.Set("port", sslVpnServer.Port)
	d.Set("compress", sslVpnServer.Compress)
	d.Set("connections", sslVpnServer.Connections)
	d.Set("max_connections", sslVpnServer.MaxConnections)
	d.Set("internet_ip", sslVpnServer.InternetIp)

	return nil
}

func resourceAliyunSslVpnServerUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if d.HasChange("proto") {
		request.Proto = d.Get("proto").(string)
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
			_, err := meta.(*AliyunClient).vpcconn.ModifySslVpnServer(request)

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
			return fmt.Errorf("Modify SslVpnServer failed")
		}
	}
	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := vpc.CreateDeleteSslVpnServerRequest()
	request.SslVpnServerId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteSslVpnServer(request)

		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete SslVpnServer error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Delete SslVpnServer timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeSslVpnServers("", d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunSslVpnServerArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateSslVpnServerRequest, error) {
	request := vpc.CreateCreateSslVpnServerRequest()
	request.RegionId = string(getRegion(d, meta))
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.ClientIpPool = d.Get("client_ip_pool").(string)
	request.LocalSubnet = d.Get("local_subnet").(string)
	request.Name = d.Get("name").(string)

	if v := d.Get("proto").(string); v != "" {
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

	return request, nil
}
