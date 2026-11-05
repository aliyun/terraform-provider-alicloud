package ga

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAlicloudGaIpSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaIpSetsRead,
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "deleting", "init", "updating"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerate_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_address_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
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

func dataSourceAlicloudGaIpSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListIpSets"
	request := make(map[string]interface{})
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ga_ip_sets", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.IpSets", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.IpSets", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IpSetId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		bandwidthVal, err := convert.FormatInt(object["Bandwidth"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"accelerate_region_id": object["AccelerateRegionId"],
			"bandwidth":            bandwidthVal,
			"ip_address_list":      object["IpAddressList"],
			"id":                   fmt.Sprint(object["IpSetId"]),
			"ip_set_id":            fmt.Sprint(object["IpSetId"]),
			"ip_version":           object["IpVersion"],
			"status":               object["State"],
		}
		ids = append(ids, fmt.Sprint(object["IpSetId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("sets", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
