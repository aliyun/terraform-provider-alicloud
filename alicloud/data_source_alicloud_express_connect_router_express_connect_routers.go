// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAliCloudExpressConnectRouterExpressConnectRouters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudExpressConnectRouterExpressConnectRouterRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ecr_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecr_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ACTIVE", "UPDATING", "ASSOCIATING", "DISSOCIATING", "LOCKED_ATTACHING", "LOCKED_DETACHING", "RECLAIMING", "DELETING"}, false),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routers": {
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
						"ecr_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alibaba_side_asn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_status": {
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudExpressConnectRouterExpressConnectRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeExpressConnectRouter"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ecr_id"); ok {
		request["EcrId"] = v
	}

	if v, ok := d.GetOk("ecr_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
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

	var expressConnectRouterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		expressConnectRouterNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_router_express_connect_routers", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.EcrList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EcrList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EcrId"])]; !ok {
					continue
				}
			}

			if expressConnectRouterNameRegex != nil {
				if !expressConnectRouterNameRegex.MatchString(fmt.Sprint(item["Name"])) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["EcrId"]),
			"ecr_id":            fmt.Sprint(object["EcrId"]),
			"ecr_name":          object["Name"],
			"description":       object["Description"],
			"owner_id":          object["OwnerId"],
			"resource_group_id": object["ResourceGroupId"],
			"alibaba_side_asn":  fmt.Sprint(object["AlibabaSideAsn"]),
			"biz_status":        object["BizStatus"],
			"status":            object["Status"],
			"create_time":       fmt.Sprint(object["GmtCreate"]),
			"modify_time":       fmt.Sprint(object["GmtModified"]),
			"tags":              tagsToMap(object["Tags"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("routers", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
