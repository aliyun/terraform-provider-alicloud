package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSslVpnClientCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSslVpnClientCertCreate,
		Read:   resourceAliyunSslVpnClientCertRead,
		Update: resourceAliyunSslVpnClientCertUpdate,
		Delete: resourceAliyunSslVpnClientCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSslVpnClientCertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateSslVpnClientCertRequest()
	request.RegionId = string(client.Region)
	request.SslVpnServerId = d.Get("ssl_vpn_server_id").(string)
	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	var sslVpnClientCert *vpc.CreateSslVpnClientCertResponse

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSslVpnClientCert(request)
		})
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		sslVpnClientCert, _ = raw.(*vpc.CreateSslVpnClientCertResponse)
		return nil
	})

	if err != nil {
		return fmt.Errorf("Create ssl vpn client cert got an error :%#v", err)
	}

	d.SetId(sslVpnClientCert.SslVpnClientCertId)

	err = vpnGatewayService.WaitForSslVpnClientCert(d.Id(), Ssl_Cert_Normal, 60)
	if err != nil {
		return fmt.Errorf("Wait Ssl Vpn Client Cert %s ready got error: %#v", sslVpnClientCert.SslVpnClientCertId, err)
	}

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	resp, err := vpnGatewayService.DescribeSslVpnClientCert(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", resp.Name)
	d.Set("status", resp.Status)
	d.Set("ssl_vpn_server_id", resp.SslVpnServerId)

	return nil
}

func resourceAliyunSslVpnClientCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateModifySslVpnClientCertRequest()
	request.SslVpnClientCertId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifySslVpnClientCert(request)
		})
		if err != nil {
			return err
		}
	}

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteSslVpnClientCertRequest()
	request.SslVpnClientCertId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnClientCert(request)
		})

		if err != nil {
			if IsExceptedError(err, SslVpnClientCertNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete SslVpnClientCert %s timeout and got an error: %#v.", request.SslVpnClientCertId, err))
		}

		_, err = vpnGatewayService.DescribeSslVpnClientCert(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}