package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEcsElasticityAssurances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsElasticityAssurancesRead,
		Schema: map[string]*schema.Schema{
			"private_pool_options_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				ValidateFunc: validation.StringInSlice([]string{"All", "Preparing", "Prepared", "Active", "Released"}, false),
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
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"assurances": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"allocated_resources": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"total_amount": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"used_amount": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"zone_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"elasticity_assurance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_charge_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"private_pool_options_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"private_pool_options_match_criteria": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"private_pool_options_name": {
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
						"total_assurance_times": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"used_assurance_times": {
							Computed: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsElasticityAssurancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if v, ok := d.GetOk("private_pool_options_ids"); ok {
		request["PrivatePoolOptions.Ids"] = convertListToJsonString(v.([]interface{}))
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeElasticityAssurances"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_elasticity_assurances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ElasticityAssuranceSet.ElasticityAssuranceItem", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ElasticityAssuranceSet.ElasticityAssuranceItem", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PrivatePoolOptionsId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                                  fmt.Sprint(object["PrivatePoolOptionsId"]),
			"description":                         object["Description"],
			"elasticity_assurance_id":             object["PrivatePoolOptionsId"],
			"private_pool_options_id":             object["PrivatePoolOptionsId"],
			"end_time":                            object["EndTime"],
			"instance_charge_type":                object["InstanceChargeType"],
			"private_pool_options_match_criteria": object["PrivatePoolOptionsMatchCriteria"],
			"resource_group_id":                   object["ResourceGroupId"],
			"start_time":                          object["StartTime"],
			"start_time_type":                     object["StartTimeType"],
			"status":                              object["Status"],
			"total_assurance_times":               object["TotalAssuranceTimes"],
			"used_assurance_times":                object["UsedAssuranceTimes"],
			"private_pool_options_name":           object["PrivatePoolOptionsName"],
		}

		allocatedResources20Maps := make([]map[string]interface{}, 0)
		allocatedResources20Raw, _ := jsonpath.Get("$.AllocatedResources.AllocatedResource", object)
		for _, value0 := range allocatedResources20Raw.([]interface{}) {
			allocatedResources20 := value0.(map[string]interface{})
			allocatedResources20Map := make(map[string]interface{})
			allocatedResources20Map["zone_id"] = allocatedResources20["zoneId"]
			allocatedResources20Map["instance_type"] = allocatedResources20["InstanceType"]
			allocatedResources20Map["total_amount"] = allocatedResources20["TotalAmount"]
			allocatedResources20Map["used_amount"] = allocatedResources20["UsedAmount"]
			allocatedResources20Maps = append(allocatedResources20Maps, allocatedResources20Map)
		}
		mapping["allocated_resources"] = allocatedResources20Maps
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

		ids = append(ids, fmt.Sprint(object["PrivatePoolOptionsId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("assurances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
