package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudPvtzZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZonesRead,

		Schema: map[string]*schema.Schema{
			"keyword": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zones": &schema.Schema{
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

	args := pvtz.CreateDescribeZonesRequest()
	if keyword, ok := d.GetOk("keyword"); ok {
		args.Keyword = keyword.(string)
	}

	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var pvtzZones []pvtz.Zone
	pvtzZoneBindVpcs := make(map[string][]map[string]interface{})

	for {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZones(args)
		})
		if err != nil {
			return fmt.Errorf("Error DescribeZones: %#v", err)
		}
		results, _ := raw.(*pvtz.DescribeZonesResponse)
		if results == nil || len(results.Zones.Zone) < 1 {
			break
		}

		for _, key := range results.Zones.Zone {
			pvtzZones = append(pvtzZones, key)

			request := pvtz.CreateDescribeZoneInfoRequest()
			request.ZoneId = key.ZoneId

			raw, errZoneInfo := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZoneInfo(request)
			})

			if errZoneInfo != nil {
				return fmt.Errorf("Error DescribeZoneInfo: %#v", errZoneInfo)
			}
			response, _ := raw.(*pvtz.DescribeZoneInfoResponse)

			var vpcs []map[string]interface{}
			for _, vpc := range response.BindVpcs.Vpc {
				mapping := map[string]interface{}{
					"region_id": vpc.RegionId,
					"vpc_id":    vpc.VpcId,
					"vpc_name":  vpc.VpcName,
				}

				vpcs = append(vpcs, mapping)
			}
			pvtzZoneBindVpcs[key.ZoneId] = vpcs
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	if len(pvtzZones) < 1 {
		return fmt.Errorf("Your query private zone return no results. Please change your keyword and try again.")
	}

	return pvtzZoneDescriptionAttributes(d, pvtzZones, pvtzZoneBindVpcs)
}

func pvtzZoneDescriptionAttributes(d *schema.ResourceData, pvtzZones []pvtz.Zone, pvtzZoneBindVpcs map[string][]map[string]interface{}) error {
	var names []string
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

		log.Printf("[DEBUG] alicloud_pvtz_zones - adding pvtzZone mapping: %v", mapping)
		names = append(names, string(key.ZoneName))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
