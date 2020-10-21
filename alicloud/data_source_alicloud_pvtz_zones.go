package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPvtzZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZonesRead,

		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_ptr": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputShortVpcsSchema()},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPvtzZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	request := pvtz.CreateDescribeZonesRequest()
	request.RegionId = client.RegionId
	if keyword, ok := d.GetOk("keyword"); ok {
		request.Keyword = keyword.(string)
	}

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	var pvtzZones []pvtz.Zone
	pvtzZoneBindVpcs := make(map[string][]map[string]interface{})

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	invoker := PvtzInvoker()

	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZones(request)
			})
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return err
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*pvtz.DescribeZonesResponse)
		if response == nil || len(response.Zones.Zone) < 1 {
			break
		}

		for _, key := range response.Zones.Zone {
			if len(idsMap) > 0 {
				if _, ok := idsMap[key.ZoneId]; !ok {
					continue
				}
			}
			pvtzZones = append(pvtzZones, key)

			object, err := pvtzService.DescribePvtzZone(key.ZoneId)

			if err != nil {
				return WrapError(err)
			}

			var vpcs []map[string]interface{}
			for _, vpc := range object.BindVpcs.Vpc {
				mapping := map[string]interface{}{
					"region_id": vpc.RegionId,
					"vpc_id":    vpc.VpcId,
					"vpc_name":  vpc.VpcName,
				}

				vpcs = append(vpcs, mapping)
			}
			pvtzZoneBindVpcs[key.ZoneId] = vpcs
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	return pvtzZoneDescriptionAttributes(d, pvtzZones, pvtzZoneBindVpcs)
}

func pvtzZoneDescriptionAttributes(d *schema.ResourceData, pvtzZones []pvtz.Zone, pvtzZoneBindVpcs map[string][]map[string]interface{}) error {
	var names []string
	var ids []string
	var s []map[string]interface{}
	for _, key := range pvtzZones {
		mapping := map[string]interface{}{
			"id":            key.ZoneId,
			"name":          key.ZoneName,
			"remark":        key.Remark,
			"record_count":  key.RecordCount,
			"is_ptr":        key.IsPtr,
			"creation_time": key.CreateTime,
			"update_time":   key.UpdateTime,
			"bind_vpcs":     pvtzZoneBindVpcs[key.ZoneId],
		}

		names = append(names, key.ZoneName)
		ids = append(ids, key.ZoneId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
