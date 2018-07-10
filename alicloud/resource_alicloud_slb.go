package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSlbName,
				Default:      resource.PrefixedUniqueId("tf-lb-"),
			},

			"internet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"vswitch_id": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbInternetDiffSuppressFunc,
			},

			"internet_charge_type": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          slb.PayByTraffic,
				ValidateFunc:     validateSlbInternetChargeType,
				DiffSuppressFunc: slbInternetChargeTypeDiffSuppressFunc,
			},

			"specification": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validateSlbInstanceSpecType,
				DiffSuppressFunc: slbInstanceSpecDiffSuppressFunc,
			},

			"bandwidth": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validateIntegerInRange(1, 1000),
				Default:          1,
				DiffSuppressFunc: slbBandwidthDiffSuppressFunc,
			},

			"listener": &schema.Schema{
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'listener' has been deprecated, and using new resource 'alicloud_slb_listener' to replace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"lb_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true, Computed: true,
						},

						"lb_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"bandwidth": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"scheduler": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"sticky_session": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"sticky_session_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https
						"cookie_timeout": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https
						"cookie": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//tcp & udp
						"persistence_timeout": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https
						"health_check": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//tcp
						"health_check_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_domain": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_uri": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_connect_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"healthy_threshold": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"unhealthy_threshold": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"health_check_timeout": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"health_check_interval": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						//http & https & tcp
						"health_check_http_code": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						//https
						"ssl_certificate_id": &schema.Schema{
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
			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instances' has been deprecated from provider version 1.6.0. New resource 'alicloud_slb_attachment' replaces it.",
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSlbCreate(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn
	args := &slb.CreateLoadBalancerArgs{
		RegionId:           getRegion(d, meta),
		LoadBalancerName:   d.Get("name").(string),
		AddressType:        slb.IntranetAddressType,
		InternetChargeType: slb.PayByTraffic,
	}
	if d.Get("internet").(bool) {
		args.AddressType = slb.InternetAddressType
	}

	if v, ok := d.GetOk("internet_charge_type"); ok && v.(string) != "" {
		args.InternetChargeType = slb.InternetChargeType(v.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		args.VSwitchId = v.(string)
	}

	if v, ok := d.GetOk("bandwidth"); ok && v.(int) != 0 {
		args.Bandwidth = v.(int)
	}

	if v, ok := d.GetOk("specification"); ok && v.(string) != "" {
		args.LoadBalancerSpec = slb.LoadBalancerSpecType(v.(string))
	}

	lb, err := slbconn.CreateLoadBalancer(args)

	if err != nil {
		if IsExceptedError(err, SlbOrderFailed) {
			return fmt.Errorf("Your account may not support to create 'paybybandwidth' load balancer. Please change it to 'paybytraffic' and try again.")
		}
		return fmt.Errorf("Create load balancer got an error: %#v", err)
	}

	d.SetId(lb.LoadBalancerId)

	if err := slbconn.WaitForLoadBalancerAsyn(lb.LoadBalancerId, slb.ActiveStatus, DefaultTimeout); err != nil {
		return fmt.Errorf("WaitForLoadbalancer %s got error: %#v", slb.ActiveStatus, err)
	}

	return resourceAliyunSlbUpdate(d, meta)
}

func resourceAliyunSlbRead(d *schema.ResourceData, meta interface{}) error {
	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", loadBalancer.LoadBalancerName)

	if loadBalancer.AddressType == slb.InternetAddressType {
		d.Set("internet", true)
	} else {
		d.Set("internet", false)
	}
	d.Set("internet_charge_type", loadBalancer.InternetChargeType)
	d.Set("bandwidth", loadBalancer.Bandwidth)
	d.Set("vswitch_id", loadBalancer.VSwitchId)
	d.Set("address", loadBalancer.Address)
	d.Set("specification", loadBalancer.LoadBalancerSpec)

	return nil
}

func resourceAliyunSlbUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		if err := slbconn.SetLoadBalancerName(d.Id(), d.Get("name").(string)); err != nil {
			return fmt.Errorf("SetLoadBalancerName got an error: %#v", err)
		}

		d.SetPartial("name")
	}

	update := false
	args := &slb.ModifyLoadBalancerInternetSpecArgs{
		LoadBalancerId: d.Id(),
	}
	if d.HasChange("internet_charge_type") && !d.IsNewResource() {
		args.InternetChargeType = slb.InternetChargeType(d.Get("internet_charge_type").(string))
		update = true
		d.SetPartial("internet_charge_type")

	}
	if d.HasChange("bandwidth") && !d.IsNewResource() {
		args.Bandwidth = d.Get("bandwidth").(int)
		update = true
		d.SetPartial("bandwidth")

	}
	if update {
		if err := slbconn.ModifyLoadBalancerInternetSpec(args); err != nil {
			return fmt.Errorf("ModifyLoadBalancerInternetSpec got an error: %#v", err)
		}

	}

	if d.HasChange("specification") && !d.IsNewResource() {
		if err := slbconn.ModifyLoadBalancerInstanceSpec(&slb.ModifyLoadBalancerInstanceSpecArgs{
			RegionId:         getRegion(d, meta),
			LoadBalancerId:   d.Id(),
			LoadBalancerSpec: slb.LoadBalancerSpecType(d.Get("specification").(string)),
		}); err != nil {
			return fmt.Errorf("ModifyLoadBalancerInstanceSpec got an error: %#v", err)
		}
		d.SetPartial("specification")
	}

	d.Partial(false)

	return resourceAliyunSlbRead(d, meta)
}

func resourceAliyunSlbDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).slbconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteLoadBalancer(d.Id())

		if err != nil {
			if IsExceptedError(err, LoadBalancerNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting slb failed: %#v", err))
		}

		loadBalancer, err := conn.DescribeLoadBalancerAttribute(d.Id())
		if err != nil {
			if IsExceptedError(err, LoadBalancerNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing slb failed when deleting SLB: %#v", err))
		}
		if loadBalancer != nil {
			return resource.RetryableError(fmt.Errorf("Delete load balancer timeout and got an error: %#v.", err))
		}
		return nil
	})
}
