// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudEnsLoadBalancerUdpListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsLoadBalancerUdpListenerRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_server_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_transmit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"established_timeout": {
							Type:     schema.TypeInt,
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
						"health_check_exp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_req": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudEnsLoadBalancerUdpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeLoadBalancerListeners"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}
	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Listeners.Listener[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["Protocol"]) != "udp" {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["LoadBalancerId"], ":", item["ListenerPort"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(request["LoadBalancerId"], ":", objectRaw["ListenerPort"])

		mapping["description"] = objectRaw["Description"]
		mapping["status"] = objectRaw["Status"]
		mapping["listener_port"] = formatInt(objectRaw["ListenerPort"])
		mapping["load_balancer_id"] = objectRaw["LoadBalancerId"]
		mapping["protocol"] = objectRaw["Protocol"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(request["LoadBalancerId"], ":", objectRaw["ListenerPort"])
		mapping, err = dataSourceAliCloudEnsLoadBalancerUdpListenerReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("listeners", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudEnsLoadBalancerUdpListenerReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ensServiceV2 := EnsServiceV2{client}
	getResp, err := ensServiceV2.DescribeEnsLoadBalancerUdpListener(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["backend_server_port"] = objectRaw["BackendServerPort"]
	mapping["description"] = objectRaw["Description"]
	mapping["eip_transmit"] = objectRaw["EipTransmit"]
	mapping["established_timeout"] = objectRaw["EstablishedTimeout"]
	mapping["health_check_connect_port"] = objectRaw["HealthCheckConnectPort"]
	mapping["health_check_connect_timeout"] = objectRaw["HealthCheckConnectTimeout"]
	mapping["health_check_exp"] = objectRaw["HealthCheckExp"]
	mapping["health_check_interval"] = objectRaw["HealthCheckInterval"]
	mapping["health_check_req"] = objectRaw["HealthCheckReq"]
	mapping["healthy_threshold"] = objectRaw["HealthyThreshold"]
	mapping["scheduler"] = objectRaw["Scheduler"]
	mapping["status"] = objectRaw["Status"]
	mapping["unhealthy_threshold"] = objectRaw["UnhealthyThreshold"]
	mapping["listener_port"] = objectRaw["ListenerPort"]

	return mapping, nil
}
