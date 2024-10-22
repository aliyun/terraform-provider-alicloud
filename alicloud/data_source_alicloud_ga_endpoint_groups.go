package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudGaEndpointGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGaEndpointGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
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
			"endpoint_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "default",
				ValidateFunc: StringInSlice([]string{"default", "virtual"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"active", "configuring", "creating", "init"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_group_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"traffic_percentage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"endpoint_group_ip_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_overrides": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"listener_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"endpoint_configurations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"probe_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"probe_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable_clientip_preservation": {
										Type:     schema.TypeBool,
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

func dataSourceAliCloudGaEndpointGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListEndpointGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["AcceleratorId"] = d.Get("accelerator_id")

	if v, ok := d.GetOk("listener_id"); ok {
		request["ListenerId"] = v
	}

	if v, ok := d.GetOk("endpoint_group_type"); ok {
		request["EndpointGroupType"] = v
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	var endpointGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		endpointGroupNameRegex = r
	}

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_endpoint_groups", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.EndpointGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EndpointGroups", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if endpointGroupNameRegex != nil {
				if !endpointGroupNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EndpointGroupId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                            fmt.Sprint(object["EndpointGroupId"]),
			"endpoint_group_id":             fmt.Sprint(object["EndpointGroupId"]),
			"listener_id":                   object["ListenerId"],
			"endpoint_group_region":         object["EndpointGroupRegion"],
			"name":                          object["Name"],
			"description":                   object["Description"],
			"health_check_interval_seconds": formatInt(object["HealthCheckIntervalSeconds"]),
			"health_check_path":             object["HealthCheckPath"],
			"health_check_port":             formatInt(object["HealthCheckPort"]),
			"health_check_protocol":         object["HealthCheckProtocol"],
			"threshold_count":               formatInt(object["ThresholdCount"]),
			"traffic_percentage":            formatInt(object["TrafficPercentage"]),
			"endpoint_group_ip_list":        object["EndpointGroupIpList"],
			"status":                        object["State"],
		}

		portOverrides := make([]map[string]interface{}, 0)
		if portOverridesList, ok := object["PortOverrides"].([]interface{}); ok {
			for _, v := range portOverridesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"endpoint_port": m1["EndpointPort"],
						"listener_port": m1["ListenerPort"],
					}
					portOverrides = append(portOverrides, temp1)
				}
			}
		}
		mapping["port_overrides"] = portOverrides

		endpointConfigurations := make([]map[string]interface{}, 0)
		if endpointConfigurationsList, ok := object["EndpointConfigurations"].([]interface{}); ok {
			for _, v := range endpointConfigurationsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"enable_clientip_preservation": m1["EnableClientIPPreservation"],
						"endpoint":                     m1["Endpoint"],
						"probe_port":                   m1["ProbePort"],
						"probe_protocol":               m1["ProbeProtocol"],
						"type":                         m1["Type"],
						"weight":                       m1["Weight"],
					}
					endpointConfigurations = append(endpointConfigurations, temp1)
				}
			}
		}
		mapping["endpoint_configurations"] = endpointConfigurations

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
