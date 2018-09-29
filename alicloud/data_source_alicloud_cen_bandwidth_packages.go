package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCenBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCensBandwidthPackagesRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(InUse),
					string(Idle),
				}),
			},
			"bandwidth_package_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
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
			"bandwidth_packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
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
	var allCenBwps []cbn.CenBandwidthPackage

	// There are three cases:
	//		1. If instance_ids input exists, first request with instance_ids filters,
	// 		then filter results by possible bandwidth_package_ids.
	//		2. Else if bandwidth_package_ids input exists, request with bandwidth_package_ids filters
	//		3. Else request with normal parameters.
	if v, ok := d.GetOk("instance_ids"); ok && len(v.([]interface{})) > 0 {
		tmpAllCenBwps, err := getCenBandwidthPackagesForKey(d, "instance_ids", meta)
		if err != nil {
			return err
		}
		// Filter tmpAllCenBwps by bandwidth_package_ids
		bwpIdsMap := make(map[string]string)
		if v, ok := d.GetOk("bandwidth_package_ids"); ok && len(v.([]interface{})) > 0 {
			for _, vv := range v.([]interface{}) {
				bwpIdsMap[Trim(vv.(string))] = Trim(vv.(string))
			}
			for _, bwp := range tmpAllCenBwps {
				if _, ok := bwpIdsMap[bwp.CenBandwidthPackageId]; !ok {
					continue
				}
				allCenBwps = append(allCenBwps, bwp)
			}
		} else {
			allCenBwps = tmpAllCenBwps
		}
	} else if v, ok := d.GetOk("bandwidth_package_ids"); ok && len(v.([]interface{})) > 0 {
		var err error
		allCenBwps, err = getCenBandwidthPackagesForKey(d, "bandwidth_package_ids", meta)
		if err != nil {
			return err
		}
	} else {
		var err error
		allCenBwps, err = doRequestCenBandwidthPackages(nil, d, meta)
		if err != nil {
			return err
		}
	}

	if len(allCenBwps) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	return cenBandwidthPackageAttributes(d, allCenBwps)
}

func getCenBandwidthPackagesForKey(d *schema.ResourceData, key string, meta interface{}) ([]cbn.CenBandwidthPackage, error) {
	var allCenBwps []cbn.CenBandwidthPackage
	multiFilters := constructCenBwpRequestFiltersForKey(d, key)
	if len(multiFilters) <= 0 {
		multiFilters = append(multiFilters, nil)
	}
	for _, filters := range multiFilters {
		tmpCenBwps, err := doRequestCenBandwidthPackages(filters, d, meta)
		if err != nil {
			return allCenBwps, err
		}
		if len(tmpCenBwps) > 0 {
			allCenBwps = append(allCenBwps, tmpCenBwps...)
		}
	}
	return allCenBwps, nil
}

func doRequestCenBandwidthPackages(filters []cbn.DescribeCenBandwidthPackagesFilter, d *schema.ResourceData, meta interface{}) ([]cbn.CenBandwidthPackage, error) {
	conn := meta.(*AliyunClient).cenconn

	args := cbn.CreateDescribeCenBandwidthPackagesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	// filter by status
	if status, ok := d.GetOk("status"); ok {
		filters = append(filters, cbn.DescribeCenBandwidthPackagesFilter{
			Key:   "Status",
			Value: &[]string{status.(string)},
		})
	}
	if filters != nil {
		args.Filter = &filters
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
			return allCenBwps, err
		}

		if resp == nil || len(resp.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			break
		}

		for _, e := range resp.CenBandwidthPackages.CenBandwidthPackage {

			// filter by name
			if nameRegex != nil {
				if !nameRegex.MatchString(e.Name) {
					continue
				}
			}
			allCenBwps = append(allCenBwps, e)
		}

		if len(resp.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return allCenBwps, err
		} else {
			args.PageNumber = page
		}
	}

	return allCenBwps, nil
}

func constructCenBwpRequestFiltersForKey(d *schema.ResourceData, key string) [][]cbn.DescribeCenBandwidthPackagesFilter {
	var res [][]cbn.DescribeCenBandwidthPackagesFilter
	maxQueryItem := 5
	if v, ok := d.GetOk(key); ok && len(v.([]interface{})) > 0 {
		// split ids
		requestTimes := len(v.([]interface{})) / maxQueryItem
		if (len(v.([]interface{})) % maxQueryItem) > 0 {
			requestTimes += 1
		}
		filtersArr := make([][]string, requestTimes)
		for k, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}

			index := k / maxQueryItem
			filtersArr[index] = append(filtersArr[index], Trim(vv.(string)))
		}

		for k := range filtersArr {
			filters := []cbn.DescribeCenBandwidthPackagesFilter{{
				Key:   convertFilterKey(key),
				Value: &filtersArr[k],
			},
			}
			res = append(res, filters)
		}
	}

	return res
}

func convertFilterKey(inputKey string) (retKey string) {
	switch inputKey {
	case "instance_ids":
		retKey = "CenId"
	case "bandwidth_package_ids":
		retKey = "CenBandwidthPackageId"
	default:
		return ""
	}
	return retKey
}

func cenBandwidthPackageAttributes(d *schema.ResourceData, allCenBwps []cbn.CenBandwidthPackage) error {
	var ids []string
	var s []map[string]interface{}

	for _, cenBwp := range allCenBwps {
		mapping := map[string]interface{}{
			"id":                            cenBwp.CenBandwidthPackageId,
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
		if len(cenBwp.CenIds.CenId) > 0 {
			// one bandwidth package attached to only one cen instance
			mapping["instance_id"] = cenBwp.CenIds.CenId[0]
		} else {
			mapping["instance_id"] = ""
		}

		ids = append(ids, cenBwp.CenBandwidthPackageId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("bandwidth_packages", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
