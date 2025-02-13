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

func dataSourceAlicloudEcsCapacityReservations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsCapacityReservationsRead,
		Schema: map[string]*schema.Schema{
			"capacity_reservation_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"platform": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"windows", "linux", "all"}, false),
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"All", "Pending", "Preparing", "Prepared", "Active", "Released"}, false),
			},
			"tags": tagsSchema(),
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
			"reservations": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"capacity_reservation_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"capacity_reservation_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"match_criteria": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"platform": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_time_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"tags": tagsSchema(),
						"time_slot": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_amount": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zone_ids": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsCapacityReservationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("capacity_reservation_ids"); ok {
		request["PrivatePoolOptions.Ids"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
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

	var capacityReservationNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		capacityReservationNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeCapacityReservations"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_capacity_reservations", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.CapacityReservationSet.CapacityReservationItem", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CapacityReservationSet.CapacityReservationItem", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PrivatePoolOptionsId"])]; !ok {
					continue
				}
			}

			if capacityReservationNameRegex != nil && !capacityReservationNameRegex.MatchString(fmt.Sprint(item["PrivatePoolOptionsName"])) {
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
			"id":                        fmt.Sprint(object["PrivatePoolOptionsId"]),
			"capacity_reservation_id":   object["PrivatePoolOptionsId"],
			"capacity_reservation_name": object["PrivatePoolOptionsName"],
			"description":               object["Description"],
			"end_time":                  object["EndTime"],
			"end_time_type":             object["EndTimeType"],
			"match_criteria":            object["PrivatePoolOptionsMatchCriteria"],
			"payment_type":              object["InstanceChargeType"],
			"platform":                  object["Platform"],
			"resource_group_id":         object["ResourceGroupId"],
			"start_time":                object["StartTime"],
			"start_time_type":           object["StartTimeType"],
			"status":                    object["Status"],
			"time_slot":                 object["TimeSlot"],
		}

		tagsMap := make(map[string]interface{})
		tagsRaw, _ := jsonpath.Get("$.Tags.Tag", object)
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

		if v, ok := object["AllocatedResources"]; ok {
			allocatedResources := v.(map[string]interface{})
			if v, ok := allocatedResources["AllocatedResource"]; ok && len(v.([]interface{})) > 0 {
				allocatedResourceMap := v.([]interface{})[0].(map[string]interface{})
				mapping["instance_type"] = allocatedResourceMap["InstanceType"]
				mapping["instance_amount"] = allocatedResourceMap["TotalAmount"]
				mapping["zone_ids"] = []string{fmt.Sprint(allocatedResourceMap["zoneId"])}
			}
		}

		ids = append(ids, fmt.Sprint(object["PrivatePoolOptionsId"]))
		names = append(names, object["PrivatePoolOptionsName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("reservations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
