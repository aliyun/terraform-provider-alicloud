package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPrivatelinkVpcEndpointServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPrivatelinkVpcEndpointServicesRead,
		Schema: map[string]*schema.Schema{
			"auto_accept_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"service_business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Normal",
				ValidateFunc: validation.StringInSlice([]string{"Normal", "FinancialLocked", "SecurityLocked"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleted", "Deleting", "Pending"}, false),
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
			"vpc_endpoint_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_accept_connection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connect_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"service_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_endpoint_service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudPrivatelinkVpcEndpointServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcEndpointServices"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("auto_accept_connection"); ok {
		request["AutoAcceptEnabled"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("service_business_status"); ok {
		request["ServiceBusinessStatus"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["ServiceStatus"] = v
	}
	if v, ok := d.GetOk("vpc_endpoint_service_name"); ok {
		request["ServiceName"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var vpcEndpointServiceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpcEndpointServiceNameRegex = r
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
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_services", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Services", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Services", response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			if vpcEndpointServiceNameRegex != nil {
				if !vpcEndpointServiceNameRegex.MatchString(item["ServiceName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServiceId"])]; !ok {
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
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"auto_accept_connection":    object["AutoAcceptEnabled"],
			"connect_bandwidth":         formatInt(object["ConnectBandwidth"]),
			"service_business_status":   object["ServiceBusinessStatus"],
			"service_description":       object["ServiceDescription"],
			"service_domain":            object["ServiceDomain"],
			"id":                        fmt.Sprint(object["ServiceId"]),
			"service_id":                fmt.Sprint(object["ServiceId"]),
			"status":                    object["ServiceStatus"],
			"vpc_endpoint_service_name": object["ServiceName"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ServiceId"]))
			names = append(names, object["ServiceName"].(string))
			s = append(s, mapping)
			continue
		}

		privatelinkService := PrivatelinkService{client}
		id := fmt.Sprint(object["ServiceId"])
		getResp, err := privatelinkService.ListVpcEndpointServiceResources(id)
		if err != nil {
			return WrapError(err)
		}

		resources := make([]map[string]interface{}, 0)
		if resourcesList, ok := getResp["Resources"].([]interface{}); ok {
			for _, v := range resourcesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"resource_id":   m1["ResourceId"],
						"resource_type": m1["ResourceType"],
					}
					resources = append(resources, temp1)
				}
			}
		}
		mapping["resources"] = resources
		ids = append(ids, fmt.Sprint(object["ServiceId"]))
		names = append(names, object["ServiceName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("services", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
