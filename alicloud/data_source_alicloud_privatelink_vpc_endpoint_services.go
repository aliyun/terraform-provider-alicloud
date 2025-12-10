package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudPrivateLinkVpcEndpointServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPrivateLinkVpcEndpointServicesRead,
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
			"vpc_endpoint_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auto_accept_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"service_business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Normal",
				ValidateFunc: StringInSlice([]string{"Normal", "FinancialLocked", "SecurityLocked"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Creating", "Deleted", "Deleting", "Pending"}, false),
			},
			"tags": tagsSchemaForceNew(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_endpoint_service_name": {
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
						"connect_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_accept_connection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"service_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudPrivateLinkVpcEndpointServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcEndpointServices"
	request := make(map[string]interface{})

	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge

	if v, ok := d.GetOk("vpc_endpoint_service_name"); ok {
		request["ServiceName"] = v
	}

	if v, ok := d.GetOkExists("auto_accept_connection"); ok {
		request["AutoAcceptEnabled"] = v
	}

	if v, ok := d.GetOk("service_business_status"); ok {
		request["ServiceBusinessStatus"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["ServiceStatus"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMaps := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMaps)
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

	var vpcEndpointServiceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpcEndpointServiceNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_services", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Services", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Services", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServiceId"])]; !ok {
					continue
				}
			}

			if vpcEndpointServiceNameRegex != nil {
				if !vpcEndpointServiceNameRegex.MatchString(fmt.Sprint(item["ServiceName"])) {
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
			"id":                        fmt.Sprint(object["ServiceId"]),
			"service_id":                fmt.Sprint(object["ServiceId"]),
			"vpc_endpoint_service_name": object["ServiceName"],
			"service_description":       object["ServiceDescription"],
			"service_domain":            object["ServiceDomain"],
			"connect_bandwidth":         formatInt(object["ConnectBandwidth"]),
			"auto_accept_connection":    object["AutoAcceptEnabled"],
			"service_business_status":   object["ServiceBusinessStatus"],
			"status":                    object["ServiceStatus"],
			"tags":                      tagsToMap(object["Tags"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ServiceName"])
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
