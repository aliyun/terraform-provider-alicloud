package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenTransitRouterPrefixListAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterPrefixListAssociationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix_list_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Updating"}, false),
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
						"prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_uid": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudCenTransitRouterPrefixListAssociationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterPrefixListAssociation"
	request := make(map[string]interface{})
	setPagingRequest(d, request, PageSizeLarge)
	request["RegionId"] = client.RegionId
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["TransitRouterTableId"] = d.Get("transit_router_table_id")

	if v, ok := d.GetOk("prefix_list_id"); ok {
		request["PrefixListId"] = v
	}

	if v, ok := d.GetOk("owner_uid"); ok {
		request["OwnerUid"] = v
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_prefix_list_associations", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.PrefixLists", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PrefixLists", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v:%v:%v", item["PrefixListId"], item["TransitRouterId"], item["TransitRouterTableId"], item["NextHop"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		ownerUid, _ := strconv.Atoi(fmt.Sprint(object["OwnerUid"]))
		mapping := map[string]interface{}{
			"id":                      fmt.Sprintf("%v:%v:%v:%v", object["PrefixListId"], object["TransitRouterId"], object["TransitRouterTableId"], object["NextHop"]),
			"prefix_list_id":          fmt.Sprint(object["PrefixListId"]),
			"transit_router_id":       object["TransitRouterId"],
			"transit_router_table_id": object["TransitRouterTableId"],
			"next_hop":                object["NextHop"],
			"next_hop_type":           object["NextHopType"],
			"next_hop_instance_id":    object["NextHopInstanceId"],
			"owner_uid":               ownerUid,
			"status":                  object["Status"],
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
