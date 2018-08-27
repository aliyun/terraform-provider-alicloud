package alicloud

import (
	"fmt"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnName,
				Default:      resource.PrefixedUniqueId("tf-vpn-"),
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateVpnInstanceChargeType,
			},

			"period": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateVpnPeriod,
			},

			"bandwidth": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateVpnBandwidth,
			},

			"enable_ipsec": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_connections": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnDescription,
			},

			"ssl_max_connections": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"internet_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"specification": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"business_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	args := vpc.CreateCreateVpnGatewayRequest()
	args.RegionId = getRegionId(d, meta)

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.Name = d.Get("name").(string)
	}

	args.VpcId = d.Get("vpc_id").(string)

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		if strings.ToLower(v.(string)) == strings.ToLower(string(PrePaid)) {
			args.InstanceChargeType = string("PREPAY")
		} else {
			args.InstanceChargeType = string("POSTPAY")
		}
	}

	if v, ok := d.GetOk("period"); ok && v.(int) != 0 {
		args.Period = requests.NewInteger(v.(int))
	}

	args.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	if v, ok := d.GetOk("enable_ipsec"); ok {
		args.EnableIpsec = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("enable_ssl"); ok {
		args.EnableSsl = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ssl_connections"); ok && v.(int) != 0 {
		args.SslConnections = requests.NewInteger(v.(int))
	}

	args.AutoPay = requests.NewBoolean(true)

	vpn, err := client.vpcconn.CreateVpnGateway(args)

	if err != nil {
		return fmt.Errorf("Create Vpn got an error: %#v", err)
	}

	d.SetId(vpn.VpnGatewayId)

	//time.Sleep(time.Duration(5) * time.Second)
	time.Sleep(10 * time.Second)
	if err := client.WaitForVpn(vpn.VpnGatewayId, Active, 2*DefaultTimeout); err != nil {
		return fmt.Errorf("WaitVpnGateway %s got error: %#v, %s", Active, err, vpn.VpnGatewayId)
	}

	return resourceAliyunVpnGatewayUpdate(d, meta)
}

func resourceAliyunVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	resp, err := client.DescribeVpnGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId(resp.VpnGatewayId)

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("vpc_id", resp.VpcId)
	d.Set("internet_ip", resp.InternetIp)
	d.Set("status", resp.Status)
	if strings.ToLower(VpnEnable) == strings.ToLower(resp.IpsecVpn) {
		d.Set("enable_ipsec", true)
	} else {
		d.Set("enable_ipsec", false)
	}

	if strings.ToLower(VpnEnable) == strings.ToLower(resp.SslVpn) {
		d.Set("enable_ssl", true)
	} else {
		d.Set("enable_ssl", false)
	}

	d.Set("ssl_max_connections", resp.SslMaxConnections)
	d.Set("business_status", resp.BusinessStatus)
	d.Set("specification", resp.Spec)
	return nil
}

func resourceAliyunVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	vpcconn := meta.(*AliyunClient).vpcconn
	req := vpc.CreateModifyVpnGatewayAttributeRequest()
	req.VpnGatewayId = d.Id()
	update := false
	d.Partial(true)
	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		req.Description = d.Get("description").(string)
		update = true
	}

	if update {
		if _, err := vpcconn.ModifyVpnGatewayAttribute(req); err != nil {
			return fmt.Errorf("ModifyVpnGatewayAttribute got an error: %#v", err)
		}
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

	if d.HasChange("instance_charge_type") {
		return fmt.Errorf("Now Cann't Support modify vpn gateway instance_charge_type")
	}

	if d.HasChange("enable_ipsec") || d.HasChange("enable_ssl") {
		return fmt.Errorf("Now Cann't Support modify ipsec/ssl switch, try to modify on the web console")
	}

	d.Partial(false)
	return resourceAliyunVpnGatewayRead(d, meta)
}

func resourceAliyunVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	req := vpc.CreateDeleteVpnGatewayRequest()
	req.VpnGatewayId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.DeleteVpnGateway(req); err != nil {
			if IsExceptedError(err, VpnNotFound) {
				return nil
			}
			if IsExceptedError(err, VpnConfiguring) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(fmt.Errorf("Error deleting vpn failed: %#v", err))
			}
			/*Vpn known issue: while the vpn is configuring, it will return unknown error*/
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(fmt.Errorf("Error deleting vpn failed: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting vpn failed: %#v", err))
		}

		if _, err := client.DescribeVpnGateway(d.Id()); err != nil {
			//if IsExceptedError(err, VpnNotFound) || IsExceptedError(err, InstanceNotFound) {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing vpn failed when deleting VPN: %#v", err))
		}
		return nil
	})
}
