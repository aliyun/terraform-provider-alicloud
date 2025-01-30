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

func dataSourceAlicloudExpressConnectRouterInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudExpressConnectRouterInterfacesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Optional: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"value": {
							Optional: true,
							ForceNew: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"include_reservation_data": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
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
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"interfaces": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"access_point_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"bandwidth": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"business_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"connected_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cross_border": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"has_reservation_data": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"hc_rate": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"hc_threshold": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"health_check_source_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"health_check_target_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_access_point_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_bandwidth": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"opposite_interface_business_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_interface_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_interface_owner_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_interface_spec": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_interface_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_router_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_router_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"opposite_vpc_instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"reservation_active_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"reservation_bandwidth": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"reservation_internet_charge_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"reservation_order_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"role": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"router_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"router_interface_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"router_interface_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"router_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"spec": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudExpressConnectRouterInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("filter"); ok {
		filterMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			filter := value0.(map[string]interface{})
			filterMap := make(map[string]interface{})
			filterMap["Key"] = filter["key"]

			filterMap["Value"] = filter["value"].([]interface{})
			filterMaps = append(filterMaps, filterMap)
		}
		request["Filter"] = filterMaps
	}

	if v, ok := d.GetOk("include_reservation_data"); ok {
		request["IncludeReservationData"] = v
	}

	setPagingRequest(d, request, PageSizeLarge)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var routerInterfaceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		routerInterfaceNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeRouterInterfaces"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_router_interfaces", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.RouterInterfaceSet.RouterInterfaceType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RouterInterfaceId"])]; !ok {
					continue
				}
			}

			if routerInterfaceNameRegex != nil && !routerInterfaceNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                                 fmt.Sprint(object["RouterInterfaceId"]),
			"access_point_id":                    object["AccessPointId"],
			"bandwidth":                          object["Bandwidth"],
			"business_status":                    object["BusinessStatus"],
			"connected_time":                     object["ConnectedTime"],
			"create_time":                        object["CreationTime"],
			"cross_border":                       object["CrossBorder"],
			"description":                        object["Description"],
			"end_time":                           object["EndTime"],
			"has_reservation_data":               fmt.Sprint(object["HasReservationData"]),
			"hc_rate":                            object["HcRate"],
			"hc_threshold":                       object["HcThreshold"],
			"health_check_source_ip":             object["HealthCheckSourceIp"],
			"health_check_target_ip":             object["HealthCheckTargetIp"],
			"opposite_access_point_id":           object["OppositeAccessPointId"],
			"opposite_bandwidth":                 object["OppositeBandwidth"],
			"opposite_interface_business_status": object["OppositeInterfaceBusinessStatus"],
			"opposite_interface_id":              object["OppositeInterfaceId"],
			"opposite_interface_owner_id":        object["OppositeInterfaceOwnerId"],
			"opposite_interface_spec":            object["OppositeInterfaceSpec"],
			"opposite_interface_status":          object["OppositeInterfaceStatus"],
			"opposite_region_id":                 object["OppositeRegionId"],
			"opposite_router_id":                 object["OppositeRouterId"],
			"opposite_router_type":               object["OppositeRouterType"],
			"opposite_vpc_instance_id":           object["OppositeVpcInstanceId"],
			"payment_type":                       object["ChargeType"],
			"reservation_active_time":            object["ReservationActiveTime"],
			"reservation_bandwidth":              object["ReservationBandwidth"],
			"reservation_internet_charge_type":   object["ReservationInternetChargeType"],
			"reservation_order_type":             object["ReservationOrderType"],
			"role":                               object["Role"],
			"router_id":                          object["RouterId"],
			"router_interface_id":                object["RouterInterfaceId"],
			"router_interface_name":              object["Name"],
			"router_type":                        object["RouterType"],
			"spec":                               object["Spec"],
			"status":                             object["Status"],
			"vpc_instance_id":                    object["VpcInstanceId"],
		}

		ids = append(ids, fmt.Sprint(object["RouterInterfaceId"]))
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

	if err := d.Set("interfaces", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
