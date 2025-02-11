package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenChildInstanceRouteEntryToAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenChildInstanceRouteEntryToAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"child_instance_route_table_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"service_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"cen_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"transit_router_attachment_id": {
				Required: true,
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
			"attachments": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cen_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"child_instance_route_table_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_cidr_block": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"service_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"transit_router_attachment_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenChildInstanceRouteEntryToAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("child_instance_route_table_id"); ok {
		request["ChildInstanceRouteTableId"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("service_type"); ok {
		request["ServiceType"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListCenChildInstanceRouteEntriesToAttachment"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_child_instance_route_entry_to_attachments", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.RouteEntry", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.RouteEntry", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CenId"], ":", item["ChildInstanceRouteTableId"], ":", item["TransitRouterAttachmentId"], ":", item["DestinationCidrBlock"])]; !ok {
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
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                            fmt.Sprint(object["CenId"], ":", object["ChildInstanceRouteTableId"], ":", object["TransitRouterAttachmentId"], ":", object["DestinationCidrBlock"]),
			"cen_id":                        object["CenId"],
			"child_instance_route_table_id": object["ChildInstanceRouteTableId"],
			"destination_cidr_block":        object["DestinationCidrBlock"],
			"service_type":                  object["ServiceType"],
			"status":                        object["Status"],
			"transit_router_attachment_id":  object["TransitRouterAttachmentId"],
		}

		ids = append(ids, fmt.Sprint(object["CenId"], ":", object["ChildInstanceRouteTableId"], ":", object["TransitRouterAttachmentId"], ":", object["DestinationCidrBlock"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
