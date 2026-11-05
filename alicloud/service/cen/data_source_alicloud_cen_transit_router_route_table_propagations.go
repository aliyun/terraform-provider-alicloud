package cen

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudCenTransitRouterRouteTablePropagations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCenTransitRouterRouteTablePropagationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Active", "Enabling", "Disabling"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"propagations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCenTransitRouterRouteTablePropagationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterRouteTablePropagations"
	request := make(map[string]interface{})
	request["MaxResults"] = helper.PageSizeLarge
	request["TransitRouterRouteTableId"] = d.Get("transit_router_route_table_id")

	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cen_transit_router_route_table_propagations", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterPropagations", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.TransitRouterPropagations", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterAttachmentId"])]; !ok {
					continue
				}
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
			"id":                            fmt.Sprint(object["TransitRouterAttachmentId"]),
			"transit_router_attachment_id":  fmt.Sprint(object["TransitRouterAttachmentId"]),
			"transit_router_route_table_id": object["TransitRouterRouteTableId"],
			"resource_id":                   object["ResourceId"],
			"resource_type":                 object["ResourceType"],
			"status":                        object["Status"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("propagations", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
