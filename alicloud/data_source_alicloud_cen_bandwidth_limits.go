package alicloud

import (
	"fmt"

	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCenBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"bandwidth_limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
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

func dataSourceAlicloudCenBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	var allCenBwLimits []cbn.CenInterRegionBandwidthLimit

	instanceIds := make([]string, 0)
	if v, ok := d.GetOk("instance_ids"); ok {
		for _, vv := range v.([]interface{}) {
			instanceIds = append(instanceIds, Trim(vv.(string)))
		}
	} else {
		instanceIds = append(instanceIds, "")
	}

	for _, instanceId := range instanceIds {
		tmpAllCenBwLimits, err := getCenBandwidthLimits(instanceId, meta)
		if err != nil {
			return err
		} else {
			allCenBwLimits = append(allCenBwLimits, tmpAllCenBwLimits...)
		}
	}

	if len(allCenBwLimits) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_cen_bandwidth_limits - CensBwLimits found: %#v", allCenBwLimits)
	return cenInterRegionBandwidthLimitsAttributes(d, allCenBwLimits)
}

func getCenBandwidthLimits(instanceId string, meta interface{}) ([]cbn.CenInterRegionBandwidthLimit, error) {
	conn := meta.(*AliyunClient).cenconn

	args := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	if instanceId != "" {
		args.CenId = instanceId
	}

	var allCenBwLimits []cbn.CenInterRegionBandwidthLimit

	for {
		resp, err := conn.DescribeCenInterRegionBandwidthLimits(args)
		if err != nil {
			return allCenBwLimits, err
		}

		if resp == nil || len(resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < 1 {
			break
		}
		allCenBwLimits = append(allCenBwLimits, resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit...)

		if len(resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return allCenBwLimits, err
		} else {
			args.PageNumber = page
		}
	}

	return allCenBwLimits, nil
}

func cenInterRegionBandwidthLimitsAttributes(d *schema.ResourceData, allCenBwLimits []cbn.CenInterRegionBandwidthLimit) error {
	var ids []string
	var s []map[string]interface{}

	for _, cenBwLimit := range allCenBwLimits {
		mapping := map[string]interface{}{
			"instance_id":        cenBwLimit.CenId,
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
	if err := d.Set("bandwidth_limits", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
