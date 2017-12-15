package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
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
				Computed:     true,
			},

			"internet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "paybytraffic",
				ValidateFunc: validateSlbInternetChargeType,
			},

			"bandwidth": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateSlbBandwidth,
				Computed:     true,
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
				Set:      schema.HashString,
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

	var slbName string
	if v, ok := d.GetOk("name"); ok {
		slbName = v.(string)
	} else {
		slbName = resource.PrefixedUniqueId("tf-lb-")
		d.Set("name", slbName)
	}

	slbArgs := &slb.CreateLoadBalancerArgs{
		RegionId:         getRegion(d, meta),
		LoadBalancerName: slbName,
		AddressType:      slb.IntranetAddressType,
	}

	if internet, ok := d.GetOk("internet"); ok && internet.(bool) {
		slbArgs.AddressType = slb.InternetAddressType
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		slbArgs.InternetChargeType = slb.InternetChargeType(v.(string))
	}

	if v, ok := d.GetOk("bandwidth"); ok && v.(int) != 0 {
		slbArgs.Bandwidth = v.(int)
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		slbArgs.VSwitchId = v.(string)
	}
	slb, err := slbconn.CreateLoadBalancer(slbArgs)
	if err != nil {
		return err
	}

	d.SetId(slb.LoadBalancerId)

	return resourceAliyunSlbUpdate(d, meta)
}

func resourceAliyunSlbRead(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn
	loadBalancer, err := slbconn.DescribeLoadBalancerAttribute(d.Id())
	if err != nil {
		if IsExceptedError(err, LoadBalancerNotFound) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error describing load balancer failed: %#v", err)
	}

	if loadBalancer == nil {
		d.SetId("")
		return nil
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

	return nil
}

func resourceAliyunSlbUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn

	d.Partial(true)

	if d.HasChange("name") {
		err := slbconn.SetLoadBalancerName(d.Id(), d.Get("name").(string))
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.Get("internet") == true && d.Get("internet_charge_type") == "paybybandwidth" {
		//don't intranet web and paybybandwidth, then can modify bandwidth
		if d.HasChange("bandwidth") {
			args := &slb.ModifyLoadBalancerInternetSpecArgs{
				LoadBalancerId: d.Id(),
				Bandwidth:      d.Get("bandwidth").(int),
			}
			err := slbconn.ModifyLoadBalancerInternetSpec(args)
			if err != nil {
				return err
			}

			d.SetPartial("bandwidth")
		}
	}

	// If we currently have instances, or did have instances,
	// we want to figure out what to add and remove from the load
	// balancer
	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandBackendServers(os.Difference(ns).List())
		add := expandBackendServers(ns.Difference(os).List())

		if len(add) > 0 {
			_, err := slbconn.AddBackendServers(d.Id(), add)
			if err != nil {
				return err
			}
		}
		if len(remove) > 0 {
			removeBackendServers := make([]string, 0, len(remove))
			for _, e := range remove {
				removeBackendServers = append(removeBackendServers, e.ServerId)
			}
			_, err := slbconn.RemoveBackendServers(d.Id(), removeBackendServers)
			if err != nil {
				return err
			}
		}

		d.SetPartial("instances")
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
			return resource.RetryableError(fmt.Errorf("LoadBalancer in use - trying again while it deleted."))
		}
		return nil
	})
}
