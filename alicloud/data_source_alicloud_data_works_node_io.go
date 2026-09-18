package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksNodeIo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksNodeIoRead,
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_env": {
				Type:     schema.TypeString,
				Required: true,
			},
			"io_type": {
				Type:     schema.TypeString,
				Required: true,
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
			"node_ios": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksNodeIoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNodeIO"
	request := make(map[string]interface{})
	request["NodeId"] = d.Get("node_id").(string)
	request["ProjectEnv"] = d.Get("project_env").(string)
	request["IoType"] = d.Get("io_type").(string)

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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_node_io", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	result, _ := resp.([]interface{})

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		tableName := fmt.Sprint(item["TableName"])
		if len(idsMap) > 0 {
			if _, ok := idsMap[tableName]; !ok {
				continue
			}
		}
		nodeId := fmt.Sprint(item["NodeId"])
		mapping := map[string]interface{}{
			"id":         tableName,
			"table_name": tableName,
			"data":       fmt.Sprint(item["Data"]),
			"node_id":    nodeId,
		}
		ids = append(ids, tableName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("node_ios", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
