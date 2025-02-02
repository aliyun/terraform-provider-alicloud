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

func dataSourceAlicloudDtsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDtsInstancesRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
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
				Default:  20,
			},
			"instances": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_endpoint_engine_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dts_instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_class": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_endpoint_engine_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_region": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"tags": tagsSchema(),
						"type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_region": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDtsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNum"] = v.(int)
	} else {
		request["PageNum"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDtsInstances"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dts_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DtsInstanceStatusList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DtsInstanceStatusList", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DtsInstanceId"])]; !ok {
					continue
				}
			}

			if instanceNameRegex != nil && !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                               fmt.Sprint(object["DtsInstanceId"]),
			"create_time":                      object["CreateTime"],
			"destination_endpoint_engine_name": object["DestEndpointEngineType"],
			"dts_instance_id":                  object["DtsInstanceId"],
			"instance_class":                   object["InstanceClass"],
			"payment_type":                     object["PayType"],
			"resource_group_id":                object["ResourceGroupId"],
			"source_endpoint_engine_name":      object["SourceEndpointEngineType"],
			"source_region":                    object["SourceEndpointRegion"],
			"status":                           object["Status"],
			"type":                             object["Type"],
			"instance_name":                    object["InstanceName"],
			"destination_region":               object["DestEndpointRegion"],
		}

		tagsMap := make(map[string]interface{})
		tagsRaw, _ := jsonpath.Get("$.Tags", object)
		if tagsRaw != nil {
			for _, value0 := range tagsRaw.([]interface{}) {
				tags := value0.(map[string]interface{})
				key := tags["TagKey"].(string)
				value := tags["TagValue"]
				if !ignoredTags(key, value) {
					tagsMap[key] = value
				}
			}
		}
		if len(tagsMap) > 0 {
			mapping["tags"] = tagsMap
		}

		ids = append(ids, fmt.Sprint(object["DtsInstanceId"]))
		names = append(names, object["InstanceName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
