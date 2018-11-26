package alicloud

import (
	"regexp"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDRDSInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDRDSInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drdsInstanceId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"createTime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"netWorkType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zoneId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vips": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDRDSInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := drds.CreateDescribeDrdsInstancesRequest()

	args.RegionId = d.Get("regionId").(string)
	args.Type = "1"

	var dbi []drds.Instance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(args)
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(*drds.DescribeDrdsInstancesResponse)
	if resp == nil || len(resp.Data.Instance) < 1 {
		return fmt.Errorf("Not found instances regionId : %s", args.RegionId)
	}

	for _, item := range resp.Data.Instance {
		if nameRegex != nil {
			if !nameRegex.MatchString(item.Description) {
				continue
			}
		}
		dbi = append(dbi, item)
	}

	return drdsInstancesDescription(d, dbi)
}

func drdsInstancesDescription(d *schema.ResourceData, dbi []drds.Instance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"drdsInstanceId": item.DrdsInstanceId,
			"description":    item.Description,
			"type":           item.Type,
			"regionId":       item.RegionId,
			"createTime":     item.CreateTime,
			"status":         item.Status,
			"netWorkType":    item.NetworkType,
			"zoneId":         item.ZoneId,
			"version":        item.Version,
			"vips":           item.Vips,
		}

		ids = append(ids, item.DrdsInstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
