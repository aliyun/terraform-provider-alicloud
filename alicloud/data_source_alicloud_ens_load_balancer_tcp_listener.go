// Package alicloud. Data source for ENS Load Balancer TCP Listener.
package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEnsLoadBalancerTcpListener() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEnsLoadBalancerTcpListenerRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"backend_server_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduler": {
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
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_check_connect_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_check_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_check_connect_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_check_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check_http_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_transmit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudEnsLoadBalancerTcpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	d.SetId(fmt.Sprint(d.Get("load_balancer_id"), ":", d.Get("listener_port")))

	object, err := ensServiceV2.DescribeEnsLoadBalancerTcpListener(d.Id())
	if err != nil {
		return WrapError(err)
	}

	// load_balancer_id is not returned by DescribeLoadBalancerTCPListenerAttribute;
	// parse it from the data source ID so the attribute is populated.
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("load_balancer_id", parts[0])
	d.Set("listener_port", formatInt(object["ListenerPort"]))
	d.Set("backend_server_port", formatInt(object["BackendServerPort"]))
	d.Set("description", object["Description"])
	d.Set("scheduler", object["Scheduler"])
	d.Set("persistence_timeout", formatInt(object["PersistenceTimeout"]))
	d.Set("established_timeout", formatInt(object["EstablishedTimeout"]))
	d.Set("healthy_threshold", formatInt(object["HealthyThreshold"]))
	d.Set("unhealthy_threshold", formatInt(object["UnhealthyThreshold"]))
	d.Set("health_check_connect_timeout", formatInt(object["HealthCheckConnectTimeout"]))
	d.Set("health_check_interval", formatInt(object["HealthCheckInterval"]))
	d.Set("health_check_connect_port", formatInt(object["HealthCheckConnectPort"]))
	d.Set("health_check_type", object["HealthCheckType"])
	d.Set("eip_transmit", object["EipTransmit"])
	d.Set("status", object["Status"])
	// protocol is always "tcp" for this TCP listener data source; the API
	// does not return it in the Describe response.
	d.Set("protocol", "tcp")
	// health_check_domain/http_code/uri are not returned by the Describe
	// API for TCP listeners; leave them unset (computed) unless the API
	// provides one (e.g. when health_check_type is "http").
	if v, ok := object["HealthCheckDomain"]; ok {
		d.Set("health_check_domain", v)
	}
	if v, ok := object["HealthCheckHttpCode"]; ok {
		d.Set("health_check_http_code", v)
	}
	if v, ok := object["HealthCheckURI"]; ok {
		d.Set("health_check_uri", v)
	}
	return nil
}
