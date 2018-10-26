package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCenRegionRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenRegionDomainRouteEntriesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenRegionDomainRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).cenconn

	args := cbn.CreateDescribeCenRegionDomainRouteEntriesRequest()
	args.CenId = d.Get("instance_id").(string)
	args.CenRegionId = d.Get("region_id").(string)

	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allCenRouteEntries []cbn.CenRouteEntry
	for pageNumber := 1; ; pageNumber++ {
		resp, err := conn.DescribeCenRegionDomainRouteEntries(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.CenRouteEntries.CenRouteEntry) < 1 {
			break
		}
		allCenRouteEntries = append(allCenRouteEntries, resp.CenRouteEntries.CenRouteEntry...)

		if len(resp.CenRouteEntries.CenRouteEntry) < PageSizeLarge {
			break
		}

		args.PageNumber = requests.NewInteger(pageNumber)
	}

	if len(allCenRouteEntries) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_cen_region_route_entries - CenRouteEntries found: %#v", allCenRouteEntries)

	return cenRegionDomainRouteEntriesAttributes(d, allCenRouteEntries)

}

func cenRegionDomainRouteEntriesAttributes(d *schema.ResourceData, allCenRouteEntries []cbn.CenRouteEntry) error {
	var s []map[string]interface{}

	for _, cenRouteEntry := range allCenRouteEntries {
		mapping := map[string]interface{}{
			"cidr_block":         cenRouteEntry.DestinationCidrBlock,
			"type":               cenRouteEntry.Type,
			"next_hop_id":        cenRouteEntry.NextHopInstanceId,
			"next_hop_type":      cenRouteEntry.NextHopType,
			"next_hop_region_id": cenRouteEntry.NextHopRegionId,
		}

		s = append(s, mapping)
	}
	id := d.Get("instance_id").(string) + COLON_SEPARATED + d.Get("region_id").(string)
	d.SetId(id)
	if err := d.Set("entries", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
