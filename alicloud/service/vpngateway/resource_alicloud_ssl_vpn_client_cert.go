package vpngateway

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
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
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_config": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
	request.ClientToken = helper.BuildClientToken(request.GetActionName())

	var response *vpc.CreateSslVpnClientCertResponse
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSslVpnClientCert(&args)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateSslVpnClientCertResponse)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ssl_vpn_client_cert", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.SslVpnClientCertId)

	err = vpnGatewayService.WaitForSslVpnClientCert(d.Id(), Ssl_Cert_Normal, DefaultTimeout)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeSslVpnClientCert(d.Id())

	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("status", object.Status)
	d.Set("ssl_vpn_server_id", object.SslVpnServerId)
	d.Set("ca_cert", object.CaCert)
	d.Set("client_cert", object.ClientCert)
	d.Set("client_key", object.ClientKey)
	d.Set("client_config", object.ClientConfig)

	return nil
}

func resourceAliyunSslVpnClientCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateModifySslVpnClientCertRequest()
	request.RegionId = client.RegionId
	request.SslVpnClientCertId = d.Id()

	request.Name = d.Get("name").(string)
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifySslVpnClientCert(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAliyunSslVpnClientCertRead(d, meta)
}

func resourceAliyunSslVpnClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteSslVpnClientCertRequest()
	request.RegionId = client.RegionId
	request.SslVpnClientCertId = d.Id()

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnClientCert(request)
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				return retry.RetryableError(err)
			} else {
				return retry.NonRetryableError(err)
			}
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidSslVpnClientCertId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return errs.WrapError(vpnGatewayService.WaitForSslVpnClientCert(d.Id(), Deleted, DefaultTimeout))
}
