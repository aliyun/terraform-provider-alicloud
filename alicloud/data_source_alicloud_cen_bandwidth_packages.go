package alicloud

import (
	"regexp"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCenBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCensBandwidthPackagesRead,

		Schema: map[string]*schema.Schema{
			"cen_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cen_bandwidth_package_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"cen_bandwidth_packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_region_a_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_region_b_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCensBandwidthPackagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).cenconn

	args := cbn.CreateDescribeCenBandwidthPackagesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)

	cenIdsMap := make(map[string]string)
	cenBandwidthPackageIdsMap := make(map[string]string)

	if v, ok := d.GetOk("cen_ids"); ok && len(v.([]interface{})) > 0 {
		for _, vv := range v.([]interface{}) {
			cenIdsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	if v, ok := d.GetOk("cen_bandwidth_package_ids"); ok && len(v.([]interface{})) > 0 {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			cenBandwidthPackageIdsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}

	var allCenBwps []cbn.CenBandwidthPackage

	for {
		resp, err := conn.DescribeCenBandwidthPackages(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			break
		}

		for _, e := range resp.CenBandwidthPackages.CenBandwidthPackage {

			// filter by cenIds
			if len(cenIdsMap) > 0 {
				isContainCenId := false
				cenIds := e.CenIds.CenId
				for _, cenId := range cenIds {
					if _, ok := cenIdsMap[cenId]; ok {
						isContainCenId = true
					}
				}
				if !isContainCenId {
					continue
				}
			}

			// filter by cenBandwidthPackageIds
			if len(cenBandwidthPackageIdsMap) > 0 {
				if _, ok := cenBandwidthPackageIdsMap[e.CenBandwidthPackageId]; !ok {
					continue
				}

			}

			// filter by name
			if nameRegex != nil {
				if !nameRegex.MatchString(e.Name) {
					continue
				}
			}

			// filter by status
			if status, ok := d.GetOk("status"); ok && string(e.Status) != status.(string) {
				continue
			}

			allCenBwps = append(allCenBwps, e)
		}

		if len(resp.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		args.PageNumber += requests.NewInteger(1)
	}

	if len(allCenBwps) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	return cenBandwidthPackageAttributes(d, allCenBwps)

}

func cenBandwidthPackageAttributes(d *schema.ResourceData, allCenBwps []cbn.CenBandwidthPackage) error {
	var ids []string
	var s []map[string]interface{}

	for _, cenBwp := range allCenBwps {
		mapping := map[string]interface{}{
			"id":                            cenBwp.CenBandwidthPackageId,
			"cen_ids":                       cenBwp.CenIds.CenId,
			"name":                          cenBwp.Name,
			"description":                   cenBwp.Description,
			"business_status":               cenBwp.BusinessStatus,
			"status":                        cenBwp.Status,
			"bandwidth":                     cenBwp.Bandwidth,
			"creation_time":                 cenBwp.CreationTime,
			"bandwidth_package_charge_type": cenBwp.BandwidthPackageChargeType,
			"geographic_region_a_id":        convertGeographicRegionId(cenBwp.GeographicRegionAId),
			"geographic_region_b_id":        convertGeographicRegionId(cenBwp.GeographicRegionBId),
		}

		ids = append(ids, cenBwp.CenBandwidthPackageId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("cen_bandwidth_packages", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
