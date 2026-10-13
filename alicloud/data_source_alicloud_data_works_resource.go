package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksResourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"modify_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {Type: schema.TypeString, Computed: true},
									"type": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"script": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path":            {Type: schema.TypeString, Computed: true},
									"runtime_command": {Type: schema.TypeString, Computed: true},
									"script_id":       {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListResources"
	query := make(map[string]interface{})
	query["RegionId"] = client.RegionId
	query["ProjectId"] = d.Get("project_id")
	if v, ok := d.GetOk("owner"); ok {
		query["Owner"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		query["Type"] = v
	}
	query["PageSize"] = PageSizeLarge
	query["PageNumber"] = 1

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
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
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_resource", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PagingInfo.Resources", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PagingInfo.Resources", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				itemId := fmt.Sprint(item["Id"])
				if _, ok := idsMap[itemId]; !ok {
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
	projectId := d.Get("project_id").(string)
	for _, object := range objects {
		resourceId := fmt.Sprint(object["Id"])
		dataSourceMap := make(map[string]interface{})
		if dsRaw, ok := object["DataSource"].(map[string]interface{}); ok {
			dataSourceMap["name"] = dsRaw["Name"]
			dataSourceMap["type"] = dsRaw["Type"]
		}
		scriptMap := make(map[string]interface{})
		if scriptRaw, ok := object["Script"].(map[string]interface{}); ok {
			scriptMap["path"] = scriptRaw["Path"]
			scriptMap["script_id"] = scriptRaw["Id"]
			if runtimeRaw, ok := scriptRaw["Runtime"].(map[string]interface{}); ok {
				scriptMap["runtime_command"] = runtimeRaw["Command"]
			}
		}
		mapping := map[string]interface{}{
			"id":            fmt.Sprint(projectId, ":", resourceId),
			"resource_id":   resourceId,
			"resource_name": object["Name"],
			"project_id":    projectId,
			"owner":         object["Owner"],
			"type":          object["Type"],
			"create_time":   object["CreateTime"],
			"modify_time":   object["ModifyTime"],
			"source_type":   object["SourceType"],
			"source_path":   object["SourcePath"],
			"target_type":   object["TargetType"],
			"target_path":   object["TargetPath"],
			"data_source":   []map[string]interface{}{dataSourceMap},
			"script":        []map[string]interface{}{scriptMap},
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("resources", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
