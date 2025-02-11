package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudCenTransitRouterVpcAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCenTransitRouterVpcAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Attached", "Attaching", "Detaching"}, false),
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
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_publish_route_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"transit_router_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCenTransitRouterVpcAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterVpcAttachments"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	request["CenId"] = d.Get("cen_id")

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

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

	var transitRouterVpcAttachmentNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouterVpcAttachmentNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_vpc_attachments", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAttachments", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["CenId"], item["TransitRouterAttachmentId"])]; !ok {
					continue
				}
			}

			if transitRouterVpcAttachmentNameRegex != nil && !transitRouterVpcAttachmentNameRegex.MatchString(fmt.Sprint(item["TransitRouterAttachmentName"])) {
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
			"id":                                    fmt.Sprintf("%v:%v", object["CenId"], object["TransitRouterAttachmentId"]),
			"cen_id":                                fmt.Sprint(object["CenId"]),
			"transit_router_attachment_id":          fmt.Sprint(object["TransitRouterAttachmentId"]),
			"vpc_id":                                object["VpcId"],
			"transit_router_id":                     object["TransitRouterId"],
			"resource_type":                         object["ResourceType"],
			"payment_type":                          convertCenTransitRouterVpcAttachmentPaymentTypeResponse(object["ChargeType"].(string)),
			"vpc_owner_id":                          fmt.Sprint(object["VpcOwnerId"]),
			"auto_publish_route_enabled":            object["AutoPublishRouteEnabled"],
			"transit_router_attachment_name":        object["TransitRouterAttachmentName"],
			"transit_router_attachment_description": object["TransitRouterAttachmentDescription"],
			"status":                                object["Status"],
		}

		if zoneMappings, ok := object["ZoneMappings"]; ok {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, zoneMappingsList := range zoneMappings.([]interface{}) {
				zoneMappingsArg := zoneMappingsList.(map[string]interface{})
				zoneMappingsMap := map[string]interface{}{}

				if vSwitchId, ok := zoneMappingsArg["VSwitchId"]; ok {
					zoneMappingsMap["vswitch_id"] = vSwitchId
				}

				if zoneId, ok := zoneMappingsArg["ZoneId"]; ok {
					zoneMappingsMap["zone_id"] = zoneId
				}

				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}

			mapping["zone_mappings"] = zoneMappingsMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TransitRouterAttachmentName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
