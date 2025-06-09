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

func dataSourceAliCloudCenTransitRouterRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCenTransitRouterRouteTablesRead,
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
			"transit_router_route_table_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"System", "Custom"}, false),
			},
			"transit_router_route_table_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Creating", "Updating", "Deleting"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Creating", "Updating", "Deleting"}, false),
			},
			"transit_router_route_table_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_route_table_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_description": {
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

func dataSourceAliCloudCenTransitRouterRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterRouteTables"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeLarge
	request["TransitRouterId"] = d.Get("transit_router_id")

	if v, ok := d.GetOk("transit_router_route_table_type"); ok {
		request["TransitRouterRouteTableType"] = v
	}

	if v, ok := d.GetOk("transit_router_route_table_status"); ok {
		request["TransitRouterRouteTableStatus"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["TransitRouterRouteTableStatus"] = v
	}

	if v, ok := d.GetOk("transit_router_route_table_ids"); ok {
		request["TransitRouterRouteTableIds"] = v
	}

	if v, ok := d.GetOk("transit_router_route_table_names"); ok {
		request["TransitRouterRouteTableNames"] = v
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

	var transitRouterRouteTableNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouterRouteTableNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_route_tables", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterRouteTables", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterRouteTables", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterRouteTableId"])]; !ok {
					continue
				}
			}

			if transitRouterRouteTableNameRegex != nil {
				if !transitRouterRouteTableNameRegex.MatchString(fmt.Sprint(item["TransitRouterRouteTableName"])) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                                     fmt.Sprint(object["TransitRouterRouteTableId"]),
			"transit_router_route_table_id":          fmt.Sprint(object["TransitRouterRouteTableId"]),
			"transit_router_route_table_type":        object["TransitRouterRouteTableType"],
			"transit_router_route_table_name":        object["TransitRouterRouteTableName"],
			"transit_router_route_table_description": object["TransitRouterRouteTableDescription"],
			"status":                                 object["TransitRouterRouteTableStatus"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TransitRouterRouteTableName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tables", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
