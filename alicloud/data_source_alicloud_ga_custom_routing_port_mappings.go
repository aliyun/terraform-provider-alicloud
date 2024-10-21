package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudGaCustomRoutingPortMappings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaCustomRoutingPortMappingsRead,
		Schema: map[string]*schema.Schema{
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, false),
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
			"custom_routing_port_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vswitch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_group_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocols": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_socket_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGaCustomRoutingPortMappingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListCustomRoutingPortMappings"
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

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_custom_routing_port_mappings", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.PortMappings", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PortMappings", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if statusOk && status.(string) != "" && status.(string) != item["DestinationTrafficState"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"accelerator_id":        object["AcceleratorId"],
			"listener_id":           object["ListenerId"],
			"endpoint_group_id":     object["EndpointGroupId"],
			"endpoint_id":           object["EndpointId"],
			"accelerator_port":      object["AcceleratorPort"],
			"vswitch":               object["Vswitch"],
			"endpoint_group_region": object["EndpointGroupRegion"],
			"protocols":             object["Protocols"],
			"status":                object["DestinationTrafficState"],
		}

		if v, ok := object["DestinationSocketAddress"]; ok {
			destinationSocketAddressMaps := make([]map[string]interface{}, 0)
			destinationSocketAddressArg := v.(map[string]interface{})
			destinationSocketAddressMap := map[string]interface{}{}

			if ipAddress, ok := destinationSocketAddressArg["IpAddress"]; ok {
				destinationSocketAddressMap["ip_address"] = ipAddress
			}

			if port, ok := destinationSocketAddressArg["Port"]; ok {
				destinationSocketAddressMap["port"] = port
			}

			destinationSocketAddressMaps = append(destinationSocketAddressMaps, destinationSocketAddressMap)
			mapping["destination_socket_address"] = destinationSocketAddressMaps
		}

		s = append(s, mapping)
	}

	d.SetId(tea.ToString(hashcode.String("GaCustomRoutingPortMappings")))

	if err := d.Set("custom_routing_port_mappings", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
