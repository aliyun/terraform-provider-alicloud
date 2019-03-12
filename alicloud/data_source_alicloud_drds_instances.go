package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
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
	var dbi []drds.Instance
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	if v, ok := d.GetOk("region_id"); ok && v.(string) != "" {
		args.RegionId = v.(string)
	}
	instIds := ""
	vswitchId := ""
	vsws := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && v.(string) != "" {
		instIds = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		vswitchId = v.(string)
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

		if instIds != "" && !strings.Contains(instIds, item.DrdsInstanceId) {
			continue
		}

		for _, vsw := range item.Vips.Vip {
			vsws[vsw.VswitchId] = vsw.VswitchId
		}
		if _, ok := vsws[vswitchId]; !ok && vswitchId != "" {
			continue
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
			"id":           item.DrdsInstanceId,
			"description":  item.Description,
			"type":         item.Type,
			"create_time":  item.CreateTime,
			"status":       item.Status,
			"network_type": item.NetworkType,
			"zone_id":      item.ZoneId,
			"version":      item.Version,
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
