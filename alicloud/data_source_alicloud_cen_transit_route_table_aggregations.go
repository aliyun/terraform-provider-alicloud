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

func dataSourceAlicloudCenTransitRouteTableAggregations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouteTableAggregationsRead,
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
			"transit_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_route_table_aggregation_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AllConfigured", "Configuring", "ConfigFailed", "PartialConfigured", "Deleting"}, false),
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
			"transit_route_table_aggregations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_route_table_aggregation_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_route_table_aggregation_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_route_table_aggregation_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_route_table_aggregation_description": {
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

func dataSourceAlicloudCenTransitRouteTableAggregationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeTransitRouteTableAggregation"
	request := make(map[string]interface{})
	request["ClientToken"] = buildClientToken("DescribeTransitRouteTableAggregation")
	request["MaxResults"] = PageSizeLarge
	request["TransitRouteTableId"] = d.Get("transit_route_table_id")

	if v, ok := d.GetOk("transit_route_table_aggregation_cidr"); ok {
		request["TransitRouteTableAggregationCidr"] = v
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	var transitRouteTableAggregationNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouteTableAggregationNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_route_table_aggregations", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if transitRouteTableAggregationNameRegex != nil && !transitRouteTableAggregationNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["TrRouteTableId"], item["TransitRouteTableAggregationCidr"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
			"id":                                          fmt.Sprintf("%v:%v", object["TrRouteTableId"], object["TransitRouteTableAggregationCidr"]),
			"transit_route_table_id":                      fmt.Sprint(object["TrRouteTableId"]),
			"transit_route_table_aggregation_cidr":        fmt.Sprint(object["TransitRouteTableAggregationCidr"]),
			"transit_route_table_aggregation_scope":       object["Scope"],
			"route_type":                                  object["RouteType"],
			"transit_route_table_aggregation_name":        object["Name"],
			"transit_route_table_aggregation_description": object["Description"],
			"status": object["Status"],
		}

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

	if err := d.Set("transit_route_table_aggregations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
