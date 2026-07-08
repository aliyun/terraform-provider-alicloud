package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudExpressConnectRouterTrAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudExpressConnectRouterTrAssociationRead,
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
			"association_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"association_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"CREATING", "ACTIVE", "INACTIVE", "ASSOCIATING", "DISSOCIATING", "UPDATING", "DELETING"}, false),
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
						"ecr_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association_node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allowed_prefixes_mode": {
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
						"allowed_prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudExpressConnectRouterTrAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeExpressConnectRouterAssociation"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["AssociationNodeType"] = "TR"
	request["MaxResults"] = PageSizeLarge
	request["ClientToken"] = buildClientToken(action)

	request["EcrId"] = d.Get("ecr_id")

	if v, ok := d.GetOk("association_id"); ok {
		request["AssociationId"] = v
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	if v, ok := d.GetOk("association_region_id"); ok {
		request["AssociationRegionId"] = v
	}

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_router_tr_associations", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.AssociationList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AssociationList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v:%v", item["EcrId"], item["AssociationId"], item["TransitRouterId"])]; !ok {
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
			"id":                      fmt.Sprintf("%v:%v:%v", object["EcrId"], object["AssociationId"], object["TransitRouterId"]),
			"ecr_id":                  fmt.Sprint(object["EcrId"]),
			"association_id":          fmt.Sprint(object["AssociationId"]),
			"transit_router_id":       fmt.Sprint(object["TransitRouterId"]),
			"association_node_type":   object["AssociationNodeType"],
			"transit_router_owner_id": fmt.Sprint(object["TransitRouterOwnerId"]),
			"cen_id":                  object["CenId"],
			"allowed_prefixes_mode":   object["AllowedPrefixesMode"],
			"status":                  object["Status"],
			"create_time":             object["GmtCreate"],
			"modify_time":             object["GmtModified"],
			"allowed_prefixes":        object["AllowedPrefixes"],
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
