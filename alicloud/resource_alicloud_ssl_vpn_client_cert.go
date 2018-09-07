package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSslVpnClientCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSslVpnClientCertCreate,
		Read:   resourceAliyunSslVpnClientCertRead,
		Update: resourceAliyunSslVpnClientCertUpdate,
		Delete: resourceAliyunSslVpnClientCertDelete,

		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSslVpnClientCertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	var sslVpnClientCert *vpc.CreateSslVpnClientCertResponse

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := buildAliyunSslVpnClientCertArgs(d, meta)

		resp, err := client.vpcconn.CreateSslVpnClientCert(args)
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		sslVpnClientCert = resp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Create ssl vpn client cert got an error :%#v", err)
	}

	d.SetId(sslVpnClientCert.SslVpnClientCertId)

	err = client.WaitForSslVpnClientCert(d.Id(), Ssl_Cert_Normal, 60)
	if err != nil {
		return fmt.Errorf("Wait Ssl Vpn Client Cert %s ready got error: %#v", sslVpnClientCert.SslVpnClientCertId, err)
	}

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	resp, err := client.DescribeSslVpnClientCert(d.Id())

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

	request := vpc.CreateModifySslVpnClientCertRequest()
	request.SslVpnClientCertId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		if _, err := meta.(*AliyunClient).vpcconn.ModifySslVpnClientCert(request); err != nil {
			return err
		}
	}

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := vpc.CreateDeleteSslVpnClientCertRequest()
	request.SslVpnClientCertId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteSslVpnClientCert(request)

		if err != nil {
			if IsExceptedError(err, SslVpnClientCertNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete SslVpnClientCert %s timeout and got an error: %#v.", request.SslVpnClientCertId, err))
		}

		if _, err := client.DescribeSslVpnClientCert(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunSslVpnClientCertArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateSslVpnClientCertRequest {
	request := vpc.CreateCreateSslVpnClientCertRequest()
	request.RegionId = string(getRegion(d, meta))
	request.SslVpnServerId = d.Get("ssl_vpn_server_id").(string)

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	return request
}
