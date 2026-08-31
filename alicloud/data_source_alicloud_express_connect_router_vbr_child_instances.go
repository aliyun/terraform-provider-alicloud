package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudExpressConnectRouterVbrChildInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudExpressConnectRouterVbrChildInstanceRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ecr_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"child_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"child_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"child_instance_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"CREATING", "ACTIVE", "ASSOCIATING", "DISSOCIATING", "UPDATING", "DELETING"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecr_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAliCloudExpressConnectRouterVbrChildInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeExpressConnectRouterChildInstance"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	request["ClientToken"] = buildClientToken(action)

	request["EcrId"] = d.Get("ecr_id")

	if v, ok := d.GetOk("child_instance_id"); ok {
		request["ChildInstanceId"] = v
	}

	if v, ok := d.GetOk("child_instance_type"); ok {
		request["ChildInstanceType"] = v
	}

	if v, ok := d.GetOk("child_instance_region_id"); ok {
		request["ChildInstanceRegionId"] = v
	}

	status, statusOk := d.GetOk("status")

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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_router_vbr_child_instances", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ChildInstanceList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ChildInstanceList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v:%v", item["EcrId"], item["ChildInstanceId"], item["ChildInstanceType"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                       fmt.Sprintf("%v:%v:%v", object["EcrId"], object["ChildInstanceId"], object["ChildInstanceType"]),
			"ecr_id":                   fmt.Sprint(object["EcrId"]),
			"child_instance_id":        fmt.Sprint(object["ChildInstanceId"]),
			"child_instance_type":      fmt.Sprint(object["ChildInstanceType"]),
			"child_instance_owner_id":  fmt.Sprint(object["ChildInstanceOwnerId"]),
			"child_instance_region_id": object["ChildInstanceRegionId"],
			"description":              object["Description"],
			"status":                   object["Status"],
			"create_time":              object["GmtCreate"],
			"modify_time":              object["GmtModified"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
