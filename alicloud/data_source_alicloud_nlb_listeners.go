package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudNlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNlbListenersRead,
		Schema: map[string]*schema.Schema{
			"listener_protocol": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"load_balancer_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"listeners": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"alpn_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"alpn_policy": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ca_certificate_ids": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ca_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"certificate_ids": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cps": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"end_port": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"idle_timeout": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"listener_description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"listener_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"listener_port": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"listener_protocol": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"mss": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"proxy_protocol_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"sec_sensor_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"security_policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"server_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_port": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("listener_protocol"); ok {
		request["ListenerProtocol"] = v
	}
	if v, ok := d.GetOk("load_balancer_ids"); ok {
		request["LoadBalancerIds"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListListeners"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nlb_listeners", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Listeners", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Listeners", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ListenerId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                     fmt.Sprint(object["ListenerId"]),
			"alpn_enabled":           object["AlpnEnabled"],
			"alpn_policy":            object["AlpnPolicy"],
			"ca_enabled":             object["CaEnabled"],
			"cps":                    object["Cps"],
			"end_port":               object["EndPort"],
			"idle_timeout":           object["IdleTimeout"],
			"listener_description":   object["ListenerDescription"],
			"listener_id":            object["ListenerId"],
			"listener_port":          object["ListenerPort"],
			"listener_protocol":      object["ListenerProtocol"],
			"load_balancer_id":       object["LoadBalancerId"],
			"mss":                    object["Mss"],
			"proxy_protocol_enabled": object["ProxyProtocolEnabled"],
			"sec_sensor_enabled":     object["SecSensorEnabled"],
			"security_policy_id":     object["SecurityPolicyId"],
			"server_group_id":        object["ServerGroupId"],
			"start_port":             object["StartPort"],
			"status":                 object["ListenerStatus"],
		}
		if v, ok := object["CaCertificateIds"]; ok {
			mapping["ca_certificate_ids"] = v.([]interface{})
		}
		if v, ok := object["CertificateIds"]; ok {
			mapping["certificate_ids"] = v.([]interface{})
		}

		ids = append(ids, fmt.Sprint(object["ListenerId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("listeners", s); err != nil {
		return WrapError(err)
	}
	return nil
}
