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

func dataSourceAlicloudVpcIpv4Gateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcIpv4GatewaysRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipv4_gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Created", "Deleting", "Pending", "Deleted"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ipv4_gateway_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_gateway_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcIpv4GatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListIpv4Gateways"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("ipv4_gateway_name"); ok {
		request["Ipv4GatewayName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var ipv4GatewayNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ipv4GatewayNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_ipv4_gateways", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Ipv4GatewayModels", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Ipv4GatewayModels", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if ipv4GatewayNameRegex != nil && !ipv4GatewayNameRegex.MatchString(fmt.Sprint(item["Ipv4GatewayName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Ipv4GatewayId"])]; !ok {
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
			"create_time":                 object["GmtCreate"],
			"enabled":                     object["Enabled"],
			"ipv4_gateway_description":    object["Ipv4GatewayDescription"],
			"id":                          fmt.Sprint(object["Ipv4GatewayId"]),
			"ipv4_gateway_id":             fmt.Sprint(object["Ipv4GatewayId"]),
			"ipv4_gateway_name":           object["Ipv4GatewayName"],
			"ipv4_gateway_route_table_id": object["Ipv4GatewayRouteTableId"],
			"status":                      object["Status"],
			"vpc_id":                      object["VpcId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Ipv4GatewayName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
