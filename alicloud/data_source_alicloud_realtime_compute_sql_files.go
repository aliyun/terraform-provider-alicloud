// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudRealtimeComputeSqlFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRealtimeComputeSqlFileRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_file_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_script": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
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

func dataSourceAliCloudRealtimeComputeSqlFileRead(d *schema.ResourceData, meta interface{}) error {
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
	var header map[string]*string
	// ListSqlFiles
	namespace := d.Get("namespace")
	action := fmt.Sprintf("/api/v2/namespaces/%s/sql-file", namespace)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = StringPointer(d.Get("workspace").(string))
	request["namespace"] = d.Get("namespace")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

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

	resp, _ := jsonpath.Get("$.data[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["workspace"], ":", item["namespace"], ":", item["sqlFileId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["workspace"], ":", objectRaw["namespace"], ":", objectRaw["sqlFileId"])

		if v, ok := objectRaw["batchMode"]; ok && v != nil {
			mapping["batch_mode"] = fmt.Sprint(v)
		}
		mapping["description"] = objectRaw["description"]
		mapping["name"] = objectRaw["name"]
		mapping["parent_id"] = objectRaw["parentId"]
		mapping["session_cluster_name"] = objectRaw["sessionClusterName"]
		mapping["sql_script"] = objectRaw["sqlScript"]
		mapping["namespace"] = objectRaw["namespace"]
		mapping["sql_file_id"] = objectRaw["sqlFileId"]
		mapping["workspace"] = objectRaw["workspace"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("files", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
