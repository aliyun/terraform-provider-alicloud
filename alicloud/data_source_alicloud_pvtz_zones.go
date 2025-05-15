package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudPvtzZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPvtzZonesRead,
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
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"query_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"query_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"LIKE", "EXACT"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"en", "zh", "jp"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_ptr": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slave_dns": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bind_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field 'creation_time' has been removed from provider version 1.107.0",
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field 'update_time' has been removed from provider version 1.107.0",
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudPvtzZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZones"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("query_vpc_id"); ok {
		request["QueryVpcId"] = v
	}

	if v, ok := d.GetOk("query_region_id"); ok {
		request["QueryRegionId"] = v
	}

	if v, ok := d.GetOk("search_mode"); ok {
		request["SearchMode"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
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

	var zoneNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		zoneNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_zones", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Zones.Zone", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Zones.Zone", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ZoneId"])]; !ok {
					continue
				}
			}

			if zoneNameRegex != nil {
				if !zoneNameRegex.MatchString(fmt.Sprint(item["ZoneName"])) {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["ZoneId"]),
			"zone_id":           fmt.Sprint(object["ZoneId"]),
			"zone_name":         object["ZoneName"],
			"name":              object["ZoneName"],
			"proxy_pattern":     object["ProxyPattern"],
			"record_count":      formatInt(object["RecordCount"]),
			"resource_group_id": object["ResourceGroupId"],
			"remark":            object["Remark"],
			"is_ptr":            object["IsPtr"],
			"create_timestamp":  formatInt(object["CreateTimestamp"]),
			"update_timestamp":  formatInt(object["UpdateTimestamp"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ZoneName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["ZoneId"])
		pvtzService := PvtzService{client}

		pvtzZoneDetail, err := pvtzService.DescribePvtzZone(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["slave_dns"] = pvtzZoneDetail["SlaveDns"]

		if bindVpcs, ok := pvtzZoneDetail["BindVpcs"]; ok {
			if vpcList, ok := bindVpcs.(map[string]interface{})["Vpc"]; ok {
				vpcMaps := make([]map[string]interface{}, 0)
				for _, vpcs := range vpcList.([]interface{}) {
					vpcArg := vpcs.(map[string]interface{})
					vpcMap := map[string]interface{}{}

					if vpcId, ok := vpcArg["VpcId"]; ok {
						vpcMap["vpc_id"] = vpcId
					}

					if vpcName, ok := vpcArg["VpcName"]; ok {
						vpcMap["vpc_name"] = vpcName
					}

					if regionId, ok := vpcArg["RegionId"]; ok {
						vpcMap["region_id"] = regionId
					}

					if regionName, ok := vpcArg["RegionName"]; ok {
						vpcMap["region_name"] = regionName
					}

					vpcMaps = append(vpcMaps, vpcMap)
				}

				mapping["bind_vpcs"] = vpcMaps
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
