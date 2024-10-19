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

func dataSourceAlicloudGaBasicEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaBasicEndpointsRead,
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
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENI", "SLB", "ECS", "NLB"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"init", "active", "updating", "binding", "unbinding", "deleting", "bound"}, false),
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
			"endpoints": {
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
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_sub_address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_sub_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_endpoint_name": {
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

func dataSourceAlicloudGaBasicEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListBasicEndpoints"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	request["EndpointGroupId"] = d.Get("endpoint_group_id")

	if v, ok := d.GetOk("endpoint_id"); ok {
		request["EndpointId"] = v
	}

	if v, ok := d.GetOk("endpoint_type"); ok {
		request["EndpointType"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	var basicEndpointNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		basicEndpointNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_basic_endpoints", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Endpoints", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Endpoints", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if basicEndpointNameRegex != nil && !basicEndpointNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["EndpointGroupId"], item["EndpointId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
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
			"id":                        fmt.Sprintf("%v:%v", object["EndpointGroupId"], object["EndpointId"]),
			"endpoint_group_id":         fmt.Sprint(object["EndpointGroupId"]),
			"endpoint_id":               fmt.Sprint(object["EndpointId"]),
			"accelerator_id":            object["AcceleratorId"],
			"endpoint_type":             object["EndpointType"],
			"endpoint_address":          object["EndpointAddress"],
			"endpoint_sub_address_type": object["EndpointSubAddressType"],
			"endpoint_sub_address":      object["EndpointSubAddress"],
			"endpoint_zone_id":          object["EndpointZoneId"],
			"basic_endpoint_name":       object["Name"],
			"status":                    object["State"],
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

	if err := d.Set("endpoints", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
