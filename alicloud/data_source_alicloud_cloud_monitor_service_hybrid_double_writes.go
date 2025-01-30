package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudMonitorServiceHybridDoubleWrites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudMonitorServiceHybridDoubleWritesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hybrid_double_writes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCloudMonitorServiceHybridDoubleWritesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeHybridDoubleWrite"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("source_namespace"); ok {
		request["SourceNamespace"] = v
	}

	if v, ok := d.GetOk("source_user_id"); ok {
		request["SourceUserId"] = v
	}

	if v, ok := d.GetOk("namespace"); ok {
		request["TargetNamespace"] = v
	}

	if v, ok := d.GetOk("user_id"); ok {
		request["TargetUserId"] = v
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2018-03-08", action, nil, request, false)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_monitor_service_hybrid_double_writes", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	resp, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Result", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%v:%v", item["SourceNamespace"], item["SourceUserId"])]; !ok {
				continue
			}
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprintf("%v:%v", object["SourceNamespace"], object["SourceUserId"]),
			"source_namespace": fmt.Sprint(object["SourceNamespace"]),
			"source_user_id":   fmt.Sprint(object["SourceUserId"]),
			"namespace":        object["Namespace"],
			"user_id":          object["UserId"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("hybrid_double_writes", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
