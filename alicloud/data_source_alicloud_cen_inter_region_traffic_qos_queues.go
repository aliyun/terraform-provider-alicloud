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

func dataSourceAlicloudCenInterRegionTrafficQosQueues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenInterRegionTrafficQosQueuesRead,
		Schema: map[string]*schema.Schema{
			"traffic_qos_policy_id": {
				Required: true,
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
			"queues": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dscps": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"inter_region_traffic_qos_queue_description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"inter_region_traffic_qos_queue_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"inter_region_traffic_qos_queue_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"remain_bandwidth_percent": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"traffic_qos_policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenInterRegionTrafficQosQueuesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("traffic_qos_policy_id"); ok {
		request["TrafficQosPolicyId"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var interRegionTrafficQosQueueNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		interRegionTrafficQosQueueNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListCenInterRegionTrafficQosQueues"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_inter_region_traffic_qos_queues", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TrafficQosQueues", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrafficQosQueues", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TrafficQosQueueId"])]; !ok {
					continue
				}
			}

			if interRegionTrafficQosQueueNameRegex != nil && !interRegionTrafficQosQueueNameRegex.MatchString(fmt.Sprint(item["TrafficQosQueueName"])) {
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
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":    fmt.Sprint(object["TrafficQosQueueId"]),
			"dscps": convertJsonStringToStringList(object["Dscps"].([]interface{})),
			"inter_region_traffic_qos_queue_description": object["TrafficQosQueueDescription"],
			"inter_region_traffic_qos_queue_id":          object["TrafficQosQueueId"],
			"inter_region_traffic_qos_queue_name":        object["TrafficQosQueueName"],
			"remain_bandwidth_percent":                   object["RemainBandwidthPercent"],
			"status":                                     object["Status"],
			"traffic_qos_policy_id":                      object["TrafficQosPolicyId"],
		}

		ids = append(ids, fmt.Sprint(object["TrafficQosQueueId"]))
		names = append(names, object["TrafficQosQueueName"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("queues", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
