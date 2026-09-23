// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigSources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigSourceRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"association_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_source_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"nacos_source_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudApigSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListSources
	action := fmt.Sprintf("/v1/sources")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		query["type"] = StringPointer(v.(string))
	}

	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["sourceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["sourceId"]

		mapping["association_reason"] = objectRaw["associationReason"]
		mapping["association_status"] = objectRaw["associationStatus"]
		mapping["create_time"] = objectRaw["createTimestamp"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["source_name"] = objectRaw["name"]
		mapping["update_time"] = objectRaw["updateTimestamp"]
		mapping["source_id"] = objectRaw["sourceId"]

		k8SSourceInfoMaps := make([]map[string]interface{}, 0)
		k8SSourceInfoMap := make(map[string]interface{})
		k8sSourceInfoRaw := make(map[string]interface{})
		if objectRaw["k8sSourceInfo"] != nil {
			k8sSourceInfoRaw = objectRaw["k8sSourceInfo"].(map[string]interface{})
		}
		if len(k8sSourceInfoRaw) > 0 {
			k8SSourceInfoMap["cluster_id"] = k8sSourceInfoRaw["clusterId"]

			k8SSourceInfoMaps = append(k8SSourceInfoMaps, k8SSourceInfoMap)
		}
		mapping["k8s_source_info"] = k8SSourceInfoMaps
		nacosSourceInfoMaps := make([]map[string]interface{}, 0)
		nacosSourceInfoMap := make(map[string]interface{})
		nacosSourceInfoRaw := make(map[string]interface{})
		if objectRaw["nacosSourceInfo"] != nil {
			nacosSourceInfoRaw = objectRaw["nacosSourceInfo"].(map[string]interface{})
		}
		if len(nacosSourceInfoRaw) > 0 {
			nacosSourceInfoMap["address"] = nacosSourceInfoRaw["address"]
			nacosSourceInfoMap["cluster_id"] = nacosSourceInfoRaw["clusterId"]
			nacosSourceInfoMap["instance_id"] = nacosSourceInfoRaw["instanceId"]

			nacosSourceInfoMaps = append(nacosSourceInfoMaps, nacosSourceInfoMap)
		}
		mapping["nacos_source_info"] = nacosSourceInfoMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["sourceId"])
		mapping, err = dataSourceAliCloudApigSourceReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("sources", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudApigSourceReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	apigServiceV2 := ApigServiceV2{client}
	getResp, err := apigServiceV2.DescribeApigSource(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["association_reason"] = objectRaw["associationReason"]
	mapping["association_status"] = objectRaw["associationStatus"]
	mapping["create_time"] = objectRaw["createTimestamp"]
	mapping["gateway_id"] = objectRaw["gatewayId"]
	mapping["resource_group_id"] = objectRaw["resourceGroupId"]
	mapping["source_name"] = objectRaw["name"]
	mapping["type"] = objectRaw["type"]
	mapping["update_time"] = objectRaw["updateTimestamp"]
	mapping["source_id"] = objectRaw["sourceId"]

	k8SSourceInfoMaps := make([]map[string]interface{}, 0)
	k8SSourceInfoMap := make(map[string]interface{})
	k8SSourceInfoRaw := make(map[string]interface{})
	if objectRaw["k8SSourceInfo"] != nil {
		k8SSourceInfoRaw = objectRaw["k8SSourceInfo"].(map[string]interface{})
	}
	if len(k8SSourceInfoRaw) > 0 {
		k8SSourceInfoMap["cluster_id"] = k8SSourceInfoRaw["clusterId"]

		k8SSourceInfoMaps = append(k8SSourceInfoMaps, k8SSourceInfoMap)
	}
	mapping["k8s_source_info"] = k8SSourceInfoMaps
	nacosSourceInfoMaps := make([]map[string]interface{}, 0)
	nacosSourceInfoMap := make(map[string]interface{})
	nacosSourceInfoRaw := make(map[string]interface{})
	if objectRaw["nacosSourceInfo"] != nil {
		nacosSourceInfoRaw = objectRaw["nacosSourceInfo"].(map[string]interface{})
	}
	if len(nacosSourceInfoRaw) > 0 {
		nacosSourceInfoMap["address"] = nacosSourceInfoRaw["address"]
		nacosSourceInfoMap["cluster_id"] = nacosSourceInfoRaw["clusterId"]
		nacosSourceInfoMap["instance_id"] = nacosSourceInfoRaw["instanceId"]

		nacosSourceInfoMaps = append(nacosSourceInfoMaps, nacosSourceInfoMap)
	}
	mapping["nacos_source_info"] = nacosSourceInfoMaps

	return mapping, nil
}
