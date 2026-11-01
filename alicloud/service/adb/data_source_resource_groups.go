package adb

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAdbResourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAdbResourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"group_name": {
				Optional: true,
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

func DataSourceAdbResourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_adb_resource_groups", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.GroupsInfo", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.GroupsInfo", response)
	}
	dbClusterId, err := jsonpath.Get("$.DBClusterId", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.DBClusterId", response)
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

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}
	return nil
}
