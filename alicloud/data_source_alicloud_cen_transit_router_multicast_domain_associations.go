package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenTransitRouterMulticastDomainAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterMulticastDomainAssociationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_multicast_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Associated", "Associating", "Dissociating"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_multicast_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_owner_id": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudCenTransitRouterMulticastDomainAssociationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterMulticastDomainAssociations"
	request := make(map[string]interface{})
	request["ClientToken"] = buildClientToken("ListTransitRouterMulticastDomainAssociations")
	request["MaxResults"] = PageSizeLarge
	request["TransitRouterMulticastDomainId"] = d.Get("transit_router_multicast_domain_id")

	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchIds"] = []string{v.(string)}
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request["ResourceId"] = v
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_associations", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterMulticastAssociations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterMulticastAssociations", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v:%v", item["TransitRouterMulticastDomainId"], item["TransitRouterAttachmentId"], item["VSwitchId"])]; !ok {
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
			"id":                                 fmt.Sprintf("%v:%v:%v", object["TransitRouterMulticastDomainId"], object["TransitRouterAttachmentId"], object["VSwitchId"]),
			"transit_router_multicast_domain_id": fmt.Sprint(object["TransitRouterMulticastDomainId"]),
			"transit_router_attachment_id":       fmt.Sprint(object["TransitRouterAttachmentId"]),
			"vswitch_id":                         fmt.Sprint(object["VSwitchId"]),
			"resource_id":                        object["ResourceId"],
			"resource_owner_id":                  formatInt(object["ResourceOwnerId"]),
			"resource_type":                      object["ResourceType"],
			"status":                             object["Status"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("associations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
