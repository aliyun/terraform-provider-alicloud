package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnGatewayCreate,
		Read:   resourceAliyunVpnGatewayRead,
		Update: resourceAliyunVpnGatewayUpdate,
		Delete: resourceAliyunVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnName,
				Default:      resource.PrefixedUniqueId("tf-vpn-"),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateInstanceChargeType,
			},

			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateVpnPeriod,
			},

			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateVpnBandwidth([]int{5, 10, 20, 50, 100, 200, 500, 1000}),
			},

			"enable_ipsec": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_connections": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: vpnSslConnectionsDiffSuppressFunc,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnDescription,
			},

			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateVpnGatewayRequest()
	request.RegionId = client.RegionId

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = d.Get("name").(string)
	}

	request.VpcId = d.Get("vpc_id").(string)

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		if v.(string) == string(PostPaid) {
			request.InstanceChargeType = string("POSTPAY")
		} else {
			request.InstanceChargeType = string("PREPAY")
		}
	}

	if v, ok := d.GetOk("period"); ok && v.(int) != 0 && request.InstanceChargeType == string("PREPAY") {
		request.Period = requests.NewInteger(v.(int))
	}

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	if v, ok := d.GetOk("enable_ipsec"); ok {
		request.EnableIpsec = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("enable_ssl"); ok {
		request.EnableSsl = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ssl_connections"); ok && v.(int) != 0 {
		request.SslConnections = requests.NewInteger(v.(int))
	}

	request.AutoPay = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateVpnGateway(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateVpnGatewayResponse)
	d.SetId(response.VpnGatewayId)

	time.Sleep(10 * time.Second)
	if err := vpnGatewayService.WaitForVpnGateway(d.Id(), Active, 2*DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpnGatewayUpdate(d, meta)
}

func resourceAliyunVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeVpnGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("internet_ip", object.InternetIp)
	d.Set("status", object.Status)
	if strings.ToLower(VpnEnable) == strings.ToLower(object.IpsecVpn) {
		d.Set("enable_ipsec", true)
	} else {
		d.Set("enable_ipsec", false)
	}

	if strings.ToLower(VpnEnable) == strings.ToLower(object.SslVpn) {
		d.Set("enable_ssl", true)
	} else {
		d.Set("enable_ssl", false)
	}

	d.Set("ssl_connections", object.SslMaxConnections)
	d.Set("business_status", object.BusinessStatus)

	spec := strings.Split(object.Spec, "M")[0]
	bandwidth, err := strconv.Atoi(spec)

	if err == nil {
		d.Set("bandwidth", bandwidth)
	} else {
		return WrapError(err)
	}

	if string("PostpayByFlow") == object.ChargeType {
		d.Set("instance_charge_type", string(PostPaid))
	} else {
		d.Set("instance_charge_type", string(PrePaid))
	}

	return nil
}

func resourceAliyunVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateModifyVpnGatewayAttributeRequest()
	request.RegionId = client.RegionId
	request.VpnGatewayId = d.Id()
	update := false
	d.Partial(true)
	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnGatewayAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
		d.SetPartial("description")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunVpnGatewayRead(d, meta)
	}

	if d.HasChange("bandwidth") {
		return fmt.Errorf("Now Cann't Support modify vpn gateway bandwidth, try to modify on the web console")
	}

	if d.HasChange("enable_ipsec") || d.HasChange("enable_ssl") {
		return fmt.Errorf("Now Cann't Support modify ipsec/ssl switch, try to modify on the web console")
	}

	d.Partial(false)
	return resourceAliyunVpnGatewayRead(d, meta)
}

func resourceAliyunVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateDeleteVpnGatewayRequest()
	request.RegionId = client.RegionId
	request.VpnGatewayId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnGateway(&args)
		})
		if err != nil {
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			/*Vpn known issue: while the vpn is configuring, it will return unknown error*/
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedError(err, VpnNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(vpnGatewayService.WaitForVpnGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}
