package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenTransitRouterMulticastDomainPeerMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterMulticastDomainPeerMembersRead,
		Schema: map[string]*schema.Schema{
			"peer_transit_router_multicast_domains": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"resource_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"transit_router_attachment_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"transit_router_multicast_domain_id": {
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
			"members": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_ip_address": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peer_transit_router_multicast_domain_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"transit_router_multicast_domain_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenTransitRouterMulticastDomainPeerMembersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("peer_transit_router_multicast_domains"); ok {
		request["PeerTransitRouterMulticastDomains.1"] = v.([]interface{})
	}
	if v, ok := d.GetOk("resource_id"); ok {
		request["ResourceId"] = v
	}
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_id"); ok {
		request["TransitRouterMulticastDomainId"] = v
	}
	request["MaxResults"] = PageSizeMedium

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
		action := "ListTransitRouterMulticastGroups"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_peer_members", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TransitRouterMulticastGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterMulticastGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterMulticastDomainId"], ":", item["GroupIpAddress"], ":", item["PeerTransitRouterMulticastDomainId"])]; !ok {
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
			"id":               fmt.Sprint(object["TransitRouterMulticastDomainId"], ":", object["GroupIpAddress"], ":", object["PeerTransitRouterMulticastDomainId"]),
			"group_ip_address": object["GroupIpAddress"],
			"peer_transit_router_multicast_domain_id": object["PeerTransitRouterMulticastDomainId"],
			"status":                             object["Status"],
			"transit_router_multicast_domain_id": object["TransitRouterMulticastDomainId"],
		}

		ids = append(ids, fmt.Sprint(object["TransitRouterMulticastDomainId"], ":", object["GroupIpAddress"], ":", object["PeerTransitRouterMulticastDomainId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("members", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
