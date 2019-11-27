package alicloud

import (
	"strings"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSlb() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbCreate,
		Read:   resourceAliyunSlbRead,
		Update: resourceAliyunSlbUpdate,
		Delete: resourceAliyunSlbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 80),
				Default:      resource.PrefixedUniqueId("tf-lb-"),
			},

			"internet": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Field 'internet' has been deprecated from provider version 1.55.3. Use 'address_type' replaces it.",
			},

			"address_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"internet"},
				ValidateFunc:  validation.StringInSlice([]string{"internet", "intranet"}, false),
			},

			"vswitch_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbInternetDiffSuppressFunc,
			},

			"internet_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          PayByTraffic,
				ValidateFunc:     validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, true),
				DiffSuppressFunc: slbInternetChargeTypeDiffSuppressFunc,
			},

			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{S1Small, S2Small, S2Medium, S3Small, S3Medium, S3Large, S4Large}, false),
			},

			"bandwidth": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 1000),
				Default:          1,
				DiffSuppressFunc: slbBandwidthDiffSuppressFunc,
			},

			"listener": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'listener' has been deprecated, and using new resource 'alicloud_slb_listener' to replace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"lb_port": {
							Type:     schema.TypeInt,
							Optional: true, Computed: true,
						},

						"lb_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"sticky_session": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"sticky_session_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"cookie_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https
						"cookie": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//tcp & udp
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https
						"health_check": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//tcp
						"health_check_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_uri": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"health_check_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_http_code": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//https
						"ssl_certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},

			//deprecated
			"instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instances' has been deprecated from provider version 1.6.0. New resource 'alicloud_slb_attachment' replaces it.",
			},

			"address": {
				Type:             schema.TypeString,
				Computed:         true,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validation.SingleIP(),
				DiffSuppressFunc: slbAddressDiffSuppressFunc,
			},

			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10),
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},

			"master_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"slave_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"delete_protection": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				DiffSuppressFunc: slbDeleteProtectionSuppressFunc,
				Default:          string(OffFlag),
			},

			"address_ip_version": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
				Default:          string(IPV4),
				ForceNew:         true,
				DiffSuppressFunc: slbAddressIpVersionSuppressFunc,
			},
		},
	}
}

func resourceAliyunSlbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateCreateLoadBalancerRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerName = d.Get("name").(string)
	request.AddressType = strings.ToLower(string(Intranet))
	request.InternetChargeType = strings.ToLower(string(PayByTraffic))
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("address_type"); ok && v.(string) != "" {
		request.AddressType = strings.ToLower(v.(string))
	} else if v, ok := d.GetOkExists("internet"); ok {
		request.AddressType = strings.ToLower(string(Intranet))
		if v.(bool) {
			request.AddressType = strings.ToLower(string(Internet))
		}
	}

	if v, ok := d.GetOk("internet_charge_type"); ok && v.(string) != "" {
		request.InternetChargeType = strings.ToLower((v.(string)))
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}

	if v, ok := d.GetOk("bandwidth"); ok && v.(int) != 0 {
		request.Bandwidth = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("specification"); ok && v.(string) != "" {
		request.LoadBalancerSpec = v.(string)
	}

	if v, ok := d.GetOk("master_zone_id"); ok && v.(string) != "" {
		request.MasterZoneId = v.(string)
	}

	if v, ok := d.GetOk("slave_zone_id"); ok && v.(string) != "" {
		request.SlaveZoneId = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		request.PayType = v.(string)
		if request.PayType == string(PrePaid) {
			request.PayType = "PrePay"
		} else {
			request.PayType = "PayOnDemand"
		}
	}
	if v, ok := d.GetOk("delete_protection"); ok && v.(string) != "" {
		request.DeleteProtection = d.Get("delete_protection").(string)
	}

	if v, ok := d.GetOk("address_ip_version"); ok && v.(string) != "" {
		request.AddressIPVersion = v.(string)
	}

	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		request.Address = v.(string)
	}

	if request.PayType == string("PrePay") {
		period := d.Get("period").(int)
		request.Duration = requests.NewInteger(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Duration = requests.NewInteger(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	var raw interface{}

	invoker := Invoker{}
	invoker.AddCatcher(Catcher{SlbTokenIsProcessing, 10, 5})

	if err := invoker.Run(func() error {
		resp, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateLoadBalancer(request)
		})
		raw = resp
		return err
	}); err != nil {
		if IsExceptedError(err, SlbOrderFailed) {
			return WrapError(err)
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateLoadBalancerResponse)
	d.SetId(response.LoadBalancerId)

	if err := slbService.WaitForSlb(response.LoadBalancerId, Active, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSlbUpdate(d, meta)
}

func resourceAliyunSlbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.LoadBalancerName)

	d.Set("internet", object.AddressType == strings.ToLower(string(Internet)))
	d.Set("address_type", object.AddressType)

	if object.InternetChargeType == strings.ToLower(string(PayByTraffic)) {
		d.Set("internet_charge_type", PayByTraffic)
	} else {
		d.Set("internet_charge_type", PayByBandwidth)
	}
	d.Set("bandwidth", object.Bandwidth)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("address", object.Address)
	d.Set("specification", object.LoadBalancerSpec)
	d.Set("instance_charge_type", object.PayType)
	d.Set("master_zone_id", object.MasterZoneId)
	d.Set("slave_zone_id", object.SlaveZoneId)
	d.Set("address_ip_version", object.AddressIPVersion)
	d.Set("resource_group_id", object.ResourceGroupId)
	if object.PayType == "PrePay" {
		d.Set("instance_charge_type", PrePaid)
	} else {
		d.Set("instance_charge_type", PostPaid)
	}
	d.Set("delete_protection", object.DeleteProtection)
	tags, _ := slbService.describeTags(d.Id())
	if len(tags) > 0 {
		if err := d.Set("tags", slbService.slbTagsToMap(tags)); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAliyunSlbUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	d.Partial(true)

	// set instance tags
	if err := slbService.setSlbInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunSlbRead(d, meta)
	}

	if d.HasChange("name") {
		request := slb.CreateSetLoadBalancerNameRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.LoadBalancerName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetLoadBalancerName(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}

	if d.HasChange("specification") {
		request := slb.CreateModifyLoadBalancerInstanceSpecRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		if _, ok := d.GetOk("specification"); !ok {
			return WrapError(Error(`'specification': required field is not set when updating it'.`))
		}
		request.LoadBalancerSpec = d.Get("specification").(string)

		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInstanceSpec(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("specification")
	}

	if d.HasChange("delete_protection") {
		request := slb.CreateSetLoadBalancerDeleteProtectionRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.DeleteProtection = d.Get("delete_protection").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetLoadBalancerDeleteProtection(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("delete_protection")
	}
	update := false
	modifyLoadBalancerInternetSpecRequest := slb.CreateModifyLoadBalancerInternetSpecRequest()
	modifyLoadBalancerInternetSpecRequest.RegionId = client.RegionId
	modifyLoadBalancerInternetSpecRequest.LoadBalancerId = d.Id()
	if d.HasChange("internet_charge_type") {
		modifyLoadBalancerInternetSpecRequest.InternetChargeType = strings.ToLower(d.Get("internet_charge_type").(string))
		update = true
		d.SetPartial("internet_charge_type")

	}
	if d.HasChange("bandwidth") {
		modifyLoadBalancerInternetSpecRequest.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
		update = true
		d.SetPartial("bandwidth")

	}
	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInternetSpec(modifyLoadBalancerInternetSpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerInternetSpecRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyLoadBalancerInternetSpecRequest.GetActionName(), raw, modifyLoadBalancerInternetSpecRequest.RpcRequest, modifyLoadBalancerInternetSpecRequest)
	}

	update = false
	modifyLoadBalancerPayTypeRequest := slb.CreateModifyLoadBalancerPayTypeRequest()
	modifyLoadBalancerPayTypeRequest.RegionId = client.RegionId
	modifyLoadBalancerPayTypeRequest.LoadBalancerId = d.Id()
	if d.HasChange("instance_charge_type") {
		modifyLoadBalancerPayTypeRequest.PayType = d.Get("instance_charge_type").(string)
		if modifyLoadBalancerPayTypeRequest.PayType == string(PrePaid) {
			modifyLoadBalancerPayTypeRequest.PayType = "PrePay"
		} else {
			modifyLoadBalancerPayTypeRequest.PayType = "PayOnDemand"
		}
		if modifyLoadBalancerPayTypeRequest.PayType == "PrePay" {
			period := d.Get("period").(int)
			modifyLoadBalancerPayTypeRequest.Duration = requests.NewInteger(period)
			modifyLoadBalancerPayTypeRequest.PricingCycle = string(Month)
			if period > 9 {
				modifyLoadBalancerPayTypeRequest.Duration = requests.NewInteger(period / 12)
				modifyLoadBalancerPayTypeRequest.PricingCycle = string(Year)
			}
			modifyLoadBalancerPayTypeRequest.AutoPay = requests.NewBoolean(true)
		}
		update = true
		d.SetPartial("instance_charge_type")
	}

	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerPayType(modifyLoadBalancerPayTypeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerPayTypeRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyLoadBalancerPayTypeRequest.GetActionName(), raw, modifyLoadBalancerPayTypeRequest.RpcRequest, modifyLoadBalancerPayTypeRequest)
	}
	d.Partial(false)

	return resourceAliyunSlbRead(d, meta)
}

func resourceAliyunSlbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteLoadBalancerRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Id()

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(slbService.WaitForSlb(d.Id(), Deleted, DefaultTimeoutMedium))
}
