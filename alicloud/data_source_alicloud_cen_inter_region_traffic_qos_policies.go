package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenInterRegionTrafficQosPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenInterRegionTrafficQosPoliciesRead,
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
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_qos_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"traffic_qos_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"traffic_qos_policy_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Active", "Modifying", "Deleting", "Deleted"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inter_region_traffic_qos_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inter_region_traffic_qos_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inter_region_traffic_qos_policy_description": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAlicloudCenInterRegionTrafficQosPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListCenInterRegionTrafficQosPolicies"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeLarge
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["TransitRouterAttachmentId"] = d.Get("transit_router_attachment_id")

	if v, ok := d.GetOk("traffic_qos_policy_id"); ok {
		request["TrafficQosPolicyId"] = v
	}

	if v, ok := d.GetOk("traffic_qos_policy_name"); ok {
		request["TrafficQosPolicyName"] = v
	}

	if v, ok := d.GetOk("traffic_qos_policy_description"); ok {
		request["TrafficQosPolicyDescription"] = v
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	var interRegionTrafficQosPolicyNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		interRegionTrafficQosPolicyNameRegex = r
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
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_inter_region_traffic_qos_policies", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TrafficQosPolicies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrafficQosPolicies", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if interRegionTrafficQosPolicyNameRegex != nil && !interRegionTrafficQosPolicyNameRegex.MatchString(fmt.Sprint(item["TrafficQosPolicyName"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TrafficQosPolicyId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["TrafficQosPolicyStatus"].(string) {
				continue
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                                          fmt.Sprint(object["TrafficQosPolicyId"]),
			"transit_router_id":                           fmt.Sprint(object["TransitRouterId"]),
			"transit_router_attachment_id":                fmt.Sprint(object["TransitRouterAttachmentId"]),
			"inter_region_traffic_qos_policy_id":          fmt.Sprint(object["TrafficQosPolicyId"]),
			"inter_region_traffic_qos_policy_name":        object["TrafficQosPolicyName"],
			"inter_region_traffic_qos_policy_description": object["TrafficQosPolicyDescription"],
			"status": object["TrafficQosPolicyStatus"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TrafficQosPolicyName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
