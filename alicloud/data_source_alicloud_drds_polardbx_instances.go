package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudDrdsPolardbxInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDrdsPolardbxInstancesRead,
		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
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
						"polardbx_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cn_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cn_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dn_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dn_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tertiary_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
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
	}
}

func dataSourceAliCloudDrdsPolardbxInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstances"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
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

	statusFilter := ""
	if v, ok := d.GetOk("status"); ok {
		statusFilter = v.(string)
	}

	var objects []interface{}
	var response map[string]interface{}
	var err error
	for {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_drds_polardbx_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.DBInstances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBInstances", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if descriptionRegex != nil {
				if !descriptionRegex.MatchString(fmt.Sprint(item["Description"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceName"])]; !ok {
					continue
				}
			}
			if statusFilter != "" && fmt.Sprint(item["Status"]) != statusFilter {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	descriptions := make([]string, 0)
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		id := fmt.Sprint(object["DBInstanceName"])
		mapping := map[string]interface{}{
			"id":                   id,
			"polardbx_instance_id": id,
			"cn_class":             object["CnNodeClassCode"],
			"cn_node_count":        object["CnNodeCount"],
			"create_time":          object["CreateTime"],
			"description":          object["Description"],
			"dn_class":             object["DnNodeClassCode"],
			"dn_node_count":        object["DnNodeCount"],
			"engine_version":       object["EngineVersion"],
			"network_type":         object["Network"],
			"payment_type":         convertDrdsPolardbxInstancePaymentTypeResponse(object["PayType"]),
			"primary_zone":         object["PrimaryZone"],
			"region_id":            object["RegionId"],
			"resource_group_id":    object["ResourceGroupId"],
			"secondary_zone":       object["SecondaryZone"],
			"status":               object["Status"],
			"storage_type":         object["StorageType"],
			"tertiary_zone":        object["TertiaryZone"],
			"topology_type":        object["TopologyType"],
			"vpc_id":               object["VPCId"],
			"zone_id":              object["ZoneId"],
		}
		if desc, ok := object["Description"].(string); ok {
			descriptions = append(descriptions, desc)
		}
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", formatInt(response["TotalNumber"])); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func convertDrdsPolardbxInstancePaymentTypeResponse(source interface{}) interface{} {
	switch fmt.Sprint(source) {
	case "POSTPAY":
		return "Postpaid"
	case "PREPAY":
		return "Prepaid"
	}
	return source
}
