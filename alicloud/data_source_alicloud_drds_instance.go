package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDRDSInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDRDSInstancesRead,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 129),
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(PrivateType)}),
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"specification": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pay_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DRDSInstancePostPayType)}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_series": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudDRDSInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := drds.CreateDescribeDrdsInstanceRequest()

	args.RegionId = client.RegionId
	args.DrdsInstanceId = d.Get("drdsInstanceId").(string)

	var dbi []drds.Instance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.DescribeDrdsInstance(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*drds.DescribeDrdsInstancesResponse)
		if resp == nil || len(resp.Data.Instance) < 1 {
			break
		}

		for _, item := range resp.Data.Instance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DrdsInstanceId) {
					continue
				}
			}
			dbi = append(dbi, item)
		}

	}

	return drdsInstancesDescription(d, dbi)
}

func drdsInstancesDescription(d *schema.ResourceData, dbi []drds.Instance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                 item.DrdsInstanceId,
			"regionId":           item.RegionId,
			"zoneId":             item.ZoneId,
			"type":               item.Type,
			"description":        item.Description,
			"netWorkType":        item.NetworkType,
			"status":             item.Status,
			"version":            item.Version,
			"vips":               item.Vips,
			"vpcCloudInstanceId": item.VpcCloudInstanceId,
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
