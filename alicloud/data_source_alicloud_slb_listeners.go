package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbListenersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"slb_listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frontend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_slave_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"established_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sticky_session": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sticky_session_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cookie_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cookie": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_http_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gzip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// http https
						"idle_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// http https
						"request_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// https
						"enable_http2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// https
						"tls_cipher_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := slb.CreateDescribeLoadBalancerAttributeRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(args)
	})
	if err != nil {
		return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
	}
	resp, _ := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	var filteredListenersTemp []slb.ListenerPortAndProtocol
	port := -1
	if v, ok := d.GetOk("frontend_port"); ok && v.(int) != 0 {
		port = v.(int)
	}
	protocol := ""
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		protocol = v.(string)
	}
	if port != -1 && protocol != "" {
		for _, listener := range resp.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if port != -1 && listener.ListenerPort != port {
				continue
			}
			if protocol != "" && listener.ListenerProtocol != protocol {
				continue
			}

			filteredListenersTemp = append(filteredListenersTemp, listener)
		}
	} else {
		filteredListenersTemp = resp.ListenerPortsAndProtocol.ListenerPortAndProtocol
	}

	if len(filteredListenersTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_slb_listeners - Slb listeners found: %#v", filteredListenersTemp)

	return slbListenersDescriptionAttributes(d, filteredListenersTemp, meta)
}

func slbListenersDescriptionAttributes(d *schema.ResourceData, listeners []slb.ListenerPortAndProtocol, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var s []map[string]interface{}

	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"frontend_port": listener.ListenerPort,
			"protocol":      listener.ListenerProtocol,
		}

		loadBalancerId := d.Get("load_balancer_id").(string)
		switch Protocol(listener.ListenerProtocol) {
		case Http:
			args := slb.CreateDescribeLoadBalancerHTTPListenerAttributeRequest()
			args.LoadBalancerId = loadBalancerId
			args.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPListenerAttribute(args)
			})
			if err == nil {
				resp, _ := raw.(*slb.DescribeLoadBalancerHTTPListenerAttributeResponse)
				mapping["backend_port"] = resp.BackendServerPort
				mapping["status"] = resp.Status
				mapping["bandwidth"] = resp.Bandwidth
				mapping["scheduler"] = resp.Scheduler
				mapping["server_group_id"] = resp.VServerGroupId
				mapping["sticky_session"] = resp.StickySession
				mapping["sticky_session_type"] = resp.StickySessionType
				mapping["cookie_timeout"] = resp.CookieTimeout
				mapping["cookie"] = resp.Cookie
				mapping["health_check"] = resp.HealthCheck
				mapping["health_check_domain"] = resp.HealthCheckDomain
				mapping["health_check_uri"] = resp.HealthCheckURI
				mapping["health_check_connect_port"] = resp.HealthCheckConnectPort
				mapping["healthy_threshold"] = resp.HealthyThreshold
				mapping["unhealthy_threshold"] = resp.UnhealthyThreshold
				mapping["health_check_timeout"] = resp.HealthCheckTimeout
				mapping["health_check_interval"] = resp.HealthCheckInterval
				mapping["health_check_http_code"] = resp.HealthCheckHttpCode
				mapping["gzip"] = resp.Gzip
				mapping["x_forwarded_for"] = resp.XForwardedFor
				mapping["x_forwarded_for_slb_ip"] = resp.XForwardedForSLBIP
				mapping["x_forwarded_for_slb_id"] = resp.XForwardedForSLBID
				mapping["x_forwarded_for_slb_proto"] = resp.XForwardedForProto
				mapping["idle_timeout"] = resp.IdleTimeout
				mapping["request_timeout"] = resp.RequestTimeout
			} else {
				log.Printf("[WARN] alicloud_slb_listeners - DescribeLoadBalancerHTTPListenerAttribute error: %v", err)
			}
		case Https:
			args := slb.CreateDescribeLoadBalancerHTTPSListenerAttributeRequest()
			args.LoadBalancerId = loadBalancerId
			args.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPSListenerAttribute(args)
			})
			if err == nil {
				resp, _ := raw.(*slb.DescribeLoadBalancerHTTPSListenerAttributeResponse)
				mapping["backend_port"] = resp.BackendServerPort
				mapping["status"] = resp.Status
				mapping["security_status"] = resp.SecurityStatus
				mapping["bandwidth"] = resp.Bandwidth
				mapping["scheduler"] = resp.Scheduler
				mapping["server_group_id"] = resp.VServerGroupId
				mapping["sticky_session"] = resp.StickySession
				mapping["sticky_session_type"] = resp.StickySessionType
				mapping["cookie_timeout"] = resp.CookieTimeout
				mapping["cookie"] = resp.Cookie
				mapping["health_check"] = resp.HealthCheck
				mapping["health_check_domain"] = resp.HealthCheckDomain
				mapping["health_check_uri"] = resp.HealthCheckURI
				mapping["health_check_connect_port"] = resp.HealthCheckConnectPort
				mapping["healthy_threshold"] = resp.HealthyThreshold
				mapping["unhealthy_threshold"] = resp.UnhealthyThreshold
				mapping["health_check_timeout"] = resp.HealthCheckTimeout
				mapping["health_check_interval"] = resp.HealthCheckInterval
				mapping["health_check_http_code"] = resp.HealthCheckHttpCode
				mapping["gzip"] = resp.Gzip
				mapping["ssl_certificate_id"] = resp.ServerCertificateId
				mapping["ca_certificate_id"] = resp.CACertificateId
				mapping["x_forwarded_for"] = resp.XForwardedFor
				mapping["x_forwarded_for_slb_ip"] = resp.XForwardedForSLBIP
				mapping["x_forwarded_for_slb_id"] = resp.XForwardedForSLBID
				mapping["x_forwarded_for_slb_proto"] = resp.XForwardedForProto
				mapping["idle_timeout"] = resp.IdleTimeout
				mapping["request_timeout"] = resp.RequestTimeout
				mapping["enable_http2"] = resp.EnableHttp2
				mapping["tls_cipher_policy"] = resp.TLSCipherPolicy
			} else {
				log.Printf("[WARN] alicloud_slb_listeners - DescribeLoadBalancerHTTPSListenerAttribute error: %v", err)
			}
		case Tcp:
			args := slb.CreateDescribeLoadBalancerTCPListenerAttributeRequest()
			args.LoadBalancerId = loadBalancerId
			args.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerTCPListenerAttribute(args)
			})
			if err == nil {
				resp, _ := raw.(*slb.DescribeLoadBalancerTCPListenerAttributeResponse)
				mapping["backend_port"] = resp.BackendServerPort
				mapping["status"] = resp.Status
				mapping["bandwidth"] = resp.Bandwidth
				mapping["scheduler"] = resp.Scheduler
				mapping["server_group_id"] = resp.VServerGroupId
				mapping["master_slave_server_group_id"] = resp.MasterSlaveServerGroupId
				mapping["persistence_timeout"] = resp.PersistenceTimeout
				mapping["established_timeout"] = resp.EstablishedTimeout
				mapping["health_check"] = resp.HealthCheck
				mapping["health_check_type"] = resp.HealthCheckType
				mapping["health_check_domain"] = resp.HealthCheckDomain
				mapping["health_check_uri"] = resp.HealthCheckURI
				mapping["health_check_connect_port"] = resp.HealthCheckConnectPort
				mapping["health_check_connect_timeout"] = resp.HealthCheckConnectTimeout
				mapping["healthy_threshold"] = resp.HealthyThreshold
				mapping["unhealthy_threshold"] = resp.UnhealthyThreshold
				mapping["health_check_interval"] = resp.HealthCheckInterval
				mapping["health_check_http_code"] = resp.HealthCheckHttpCode
			} else {
				log.Printf("[WARN] alicloud_slb_listeners - DescribeLoadBalancerTCPListenerAttribute error: %v", err)
			}
		case Udp:
			args := slb.CreateDescribeLoadBalancerUDPListenerAttributeRequest()
			args.LoadBalancerId = loadBalancerId
			args.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerUDPListenerAttribute(args)
			})
			if err == nil {
				resp, _ := raw.(*slb.DescribeLoadBalancerUDPListenerAttributeResponse)
				mapping["backend_port"] = resp.BackendServerPort
				mapping["status"] = resp.Status
				mapping["bandwidth"] = resp.Bandwidth
				mapping["scheduler"] = resp.Scheduler
				mapping["server_group_id"] = resp.VServerGroupId
				mapping["master_slave_server_group_id"] = resp.MasterSlaveServerGroupId
				mapping["persistence_timeout"] = resp.PersistenceTimeout
				mapping["health_check"] = resp.HealthCheck
				mapping["health_check_connect_port"] = resp.HealthCheckConnectPort
				mapping["health_check_connect_timeout"] = resp.HealthCheckConnectTimeout
				mapping["healthy_threshold"] = resp.HealthyThreshold
				mapping["unhealthy_threshold"] = resp.UnhealthyThreshold
				mapping["health_check_interval"] = resp.HealthCheckInterval
			} else {
				log.Printf("[WARN] alicloud_slb_listeners - DescribeLoadBalancerUDPListenerAttribute error: %v", err)
			}
		}

		log.Printf("[DEBUG] alicloud_slb_listeners - adding slb_listener mapping: %v", mapping)
		ids = append(ids, strconv.Itoa(listener.ListenerPort))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_listeners", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
