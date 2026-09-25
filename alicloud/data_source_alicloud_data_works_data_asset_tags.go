package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksDataAssetTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksDataAssetTagsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "System"}, false),
			},
			"value_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Boolean", "Int", "String", "Double"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"managers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksDataAssetTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDataAssetTags"
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	query["RegionId"] = client.RegionId
	query["PageSize"] = PageSizeLarge
	query["PageNumber"] = 1
	if v, ok := d.GetOk("key"); ok {
		query["Key"] = v
	}
	if v, ok := d.GetOk("category"); ok {
		query["Category"] = v
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
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_data_asset_tags", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DataAssetTags", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DataAssetTags", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Key"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		if vt, ok := d.GetOk("value_type"); ok {
			if fmt.Sprint(object["ValueType"]) != vt.(string) {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["Key"]),
			"key":         fmt.Sprint(object["Key"]),
			"value_type":  fmt.Sprint(object["ValueType"]),
			"category":    fmt.Sprint(object["Category"]),
			"description": fmt.Sprint(object["Description"]),
			"create_time": fmt.Sprint(object["CreateTime"]),
			"modify_time": fmt.Sprint(object["ModifyTime"]),
		}
		if v, ok := object["Values"]; ok && v != nil {
			mapping["values"] = v
		} else {
			mapping["values"] = make([]interface{}, 0)
		}
		if v, ok := object["Managers"]; ok && v != nil {
			mapping["managers"] = v
		} else {
			mapping["managers"] = make([]interface{}, 0)
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tags", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
