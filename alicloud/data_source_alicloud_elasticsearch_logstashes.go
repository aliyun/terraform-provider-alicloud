// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudElasticsearchLogstashes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudElasticsearchLogstashRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logstashes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vs_area": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
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
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
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
		},
	}
}

func dataSourceAliCloudElasticsearchLogstashRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListLogstash
	action := fmt.Sprintf("/openapi/logstashes")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	query["instanceId"] = StringPointer(d.Get("instance_id").(string))
	if v, ok := d.GetOk("description"); ok {
		query["description"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		query["instanceId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		query["version"] = StringPointer(v.(string))
	}

	query["size"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["page"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("elasticsearch", "2017-06-13", action, query, nil, nil)

			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
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

		resp, _ := jsonpath.Get("$.Result[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["instanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["page"])
		query["page"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["instanceId"]

		mapping["create_time"] = objectRaw["createdAt"]
		mapping["description"] = objectRaw["description"]
		mapping["node_amount"] = objectRaw["nodeAmount"]
		mapping["payment_type"] = objectRaw["paymentType"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["status"] = objectRaw["status"]
		mapping["updated_at"] = objectRaw["updatedAt"]
		mapping["version"] = objectRaw["version"]
		mapping["instance_id"] = objectRaw["instanceId"]

		networkConfigMaps := make([]map[string]interface{}, 0)
		networkConfigMap := make(map[string]interface{})
		networkConfigRaw := make(map[string]interface{})
		if objectRaw["networkConfig"] != nil {
			networkConfigRaw = objectRaw["networkConfig"].(map[string]interface{})
		}
		if len(networkConfigRaw) > 0 {
			networkConfigMap["type"] = networkConfigRaw["type"]
			networkConfigMap["vswitch_id"] = networkConfigRaw["vswitchId"]
			networkConfigMap["vpc_id"] = networkConfigRaw["vpcId"]
			networkConfigMap["vs_area"] = networkConfigRaw["vsArea"]

			networkConfigMaps = append(networkConfigMaps, networkConfigMap)
		}
		mapping["network_config"] = networkConfigMaps
		nodeSpecMaps := make([]map[string]interface{}, 0)
		nodeSpecMap := make(map[string]interface{})
		nodeSpecRaw := make(map[string]interface{})
		if objectRaw["nodeSpec"] != nil {
			nodeSpecRaw = objectRaw["nodeSpec"].(map[string]interface{})
		}
		if len(nodeSpecRaw) > 0 {
			nodeSpecMap["disk"] = nodeSpecRaw["disk"]
			nodeSpecMap["disk_type"] = nodeSpecRaw["diskType"]
			nodeSpecMap["spec"] = nodeSpecRaw["spec"]

			nodeSpecMaps = append(nodeSpecMaps, nodeSpecMap)
		}
		mapping["node_spec"] = nodeSpecMaps

		tagsMaps := objectRaw["tags"]
		mapping["tags"] = tagsToMap(tagsMaps)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("logstashes", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
