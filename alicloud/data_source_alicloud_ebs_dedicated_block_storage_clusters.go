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

func dataSourceAlicloudEbsDedicatedBlockStorageClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEbsDedicatedBlockStorageClustersRead,
		Schema: map[string]*schema.Schema{
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
			"clusters": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"available_capacity": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"category": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dedicated_block_storage_cluster_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dedicated_block_storage_cluster_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"delivery_capacity": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"expired_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"performance_level": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"supported_category": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"total_capacity": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"used_capacity": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEbsDedicatedBlockStorageClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("dedicated_block_storage_cluster_id"); ok {
		request["DedicatedBlockStorageClusterId"] = v
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

	var dedicatedBlockStorageClusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dedicatedBlockStorageClusterNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDedicatedBlockStorageClusters"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ebs_dedicated_block_storage_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DedicatedBlockStorageClusters", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DedicatedBlockStorageClusters", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DedicatedBlockStorageClusterId"])]; !ok {
					continue
				}
			}

			if dedicatedBlockStorageClusterNameRegex != nil && !dedicatedBlockStorageClusterNameRegex.MatchString(fmt.Sprint(item["DedicatedBlockStorageClusterName"])) {
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
			"id":                                   fmt.Sprint(object["DedicatedBlockStorageClusterId"]),
			"available_capacity":                   object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["AvailableCapacity"],
			"category":                             object["Category"],
			"create_time":                          object["CreateTime"],
			"dedicated_block_storage_cluster_id":   object["DedicatedBlockStorageClusterId"],
			"dedicated_block_storage_cluster_name": object["DedicatedBlockStorageClusterName"],
			"delivery_capacity":                    object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["DeliveryCapacity"],
			"description":                          object["Description"],
			"expired_time":                         object["ExpiredTime"],
			"performance_level":                    object["PerformanceLevel"],
			"resource_group_id":                    object["ResourceGroupId"],
			"status":                               object["Status"],
			"supported_category":                   object["SupportedCategory"],
			"total_capacity":                       object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["TotalCapacity"],
			"type":                                 object["Type"],
			"used_capacity":                        object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["UsedCapacity"],
			"zone_id":                              object["ZoneId"],
		}

		ids = append(ids, fmt.Sprint(object["DedicatedBlockStorageClusterId"]))
		names = append(names, object["DedicatedBlockStorageClusterName"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
