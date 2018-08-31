package alicloud

import (
	"fmt"

	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCenInterRegionBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenInterRegionBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"cen_inter_region_bandwidth_limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenInterRegionBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).cenconn

	cenId := d.Get("cen_id")

	args := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.CenId = cenId.(string)

	pageNumber := 1
	args.PageNumber = requests.NewInteger(pageNumber)

	var allcenBwLimits []cbn.CenInterRegionBandwidthLimit

	for {
		resp, err := conn.DescribeCenInterRegionBandwidthLimits(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < 1 {
			break
		}
		allcenBwLimits = append(allcenBwLimits, resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit...)

		if len(resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < PageSizeLarge {
			break
		}

		pageNumber++
		args.PageNumber = requests.NewInteger(pageNumber)
	}

	if len(allcenBwLimits) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_cen_inter_region_bandwidth_limits - CensBwLimits found: %#v", allcenBwLimits)

	return cenInterRegionBandwidthLimitsAttributes(d, allcenBwLimits)

}

func cenInterRegionBandwidthLimitsAttributes(d *schema.ResourceData, allcenBwLimits []cbn.CenInterRegionBandwidthLimit) error {
	var ids []string
	var s []map[string]interface{}

	for _, cenBwLimit := range allcenBwLimits {
		mapping := map[string]interface{}{
			"cen_id":             cenBwLimit.CenId,
			"local_region_id":    cenBwLimit.LocalRegionId,
			"opposite_region_id": cenBwLimit.OppositeRegionId,
			"status":             cenBwLimit.Status,
			"bandwidth_limit":    cenBwLimit.BandwidthLimit,
		}

		id := cenBwLimit.CenId + ":" + cenBwLimit.LocalRegionId + ":" + cenBwLimit.OppositeRegionId
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("cen_inter_region_bandwidth_limits", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
