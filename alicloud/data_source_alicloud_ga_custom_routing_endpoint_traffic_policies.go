package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudGaCustomRoutingEndpointTrafficPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGaCustomRoutingEndpointTrafficPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_routing_endpoint_traffic_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_routing_endpoint_traffic_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"to_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudGaCustomRoutingEndpointTrafficPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListCustomRoutingEndpointTrafficPolicies"
	request := make(map[string]interface{})
	setPagingRequest(d, request, PageSizeLarge)
	request["RegionId"] = client.RegionId
	request["AcceleratorId"] = d.Get("accelerator_id")

	if v, ok := d.GetOk("listener_id"); ok {
		request["ListenerId"] = v
	}

	if v, ok := d.GetOk("endpoint_group_id"); ok {
		request["EndpointGroupId"] = v
	}

	if v, ok := d.GetOk("endpoint_id"); ok {
		request["EndpointId"] = v
	}

	if v, ok := d.GetOk("address"); ok {
		request["Address"] = v
	}

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

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_custom_routing_endpoint_traffic_policies", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Policies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["EndpointId"], item["PolicyId"])]; !ok {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprintf("%v:%v", object["EndpointId"], object["PolicyId"]),
			"endpoint_id": fmt.Sprint(object["EndpointId"]),
			"custom_routing_endpoint_traffic_policy_id": fmt.Sprint(object["PolicyId"]),
			"accelerator_id":    object["AcceleratorId"],
			"listener_id":       object["ListenerId"],
			"endpoint_group_id": object["EndpointGroupId"],
			"address":           object["Address"],
		}

		if portRangesList, ok := object["PortRanges"]; ok {
			portRangesMaps := make([]map[string]interface{}, 0)
			for _, portRanges := range portRangesList.([]interface{}) {
				portRangesArg := portRanges.(map[string]interface{})
				portRangesMap := map[string]interface{}{}

				if fromPort, ok := portRangesArg["FromPort"]; ok {
					portRangesMap["from_port"] = fromPort
				}

				if toPort, ok := portRangesArg["ToPort"]; ok {
					portRangesMap["to_port"] = toPort
				}

				portRangesMaps = append(portRangesMaps, portRangesMap)
			}

			mapping["port_ranges"] = portRangesMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("custom_routing_endpoint_traffic_policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
