package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAdbResourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAdbResourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"group_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
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
			"groups": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"db_cluster_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_num": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"user": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAdbResourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
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

	var objects []interface{}
	var response map[string]interface{}
	var err error
	action := "DescribeDBResourceGroup"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_resource_groups", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.GroupsInfo", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.GroupsInfo", response)
	}
	dbClusterId, err := jsonpath.Get("$.DBClusterId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBClusterId", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(dbClusterId, ":", item["GroupName"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":            fmt.Sprint(dbClusterId, ":", object["GroupName"]),
			"create_time":   object["CreateTime"],
			"db_cluster_id": object["DBClusterId"],
			"group_name":    object["GroupName"],
			"group_type":    object["GroupType"],
			"node_num":      object["NodeNum"],
			"user":          object["GroupUsers"],
		}

		ids = append(ids, fmt.Sprint(dbClusterId, ":", object["GroupName"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	return nil
}
