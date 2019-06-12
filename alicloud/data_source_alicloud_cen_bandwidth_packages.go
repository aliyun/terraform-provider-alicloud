package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCenBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCensBandwidthPackagesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
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
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packages": {
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

	multiFilters := constructCenBwpRequestFilters(d)
	if len(multiFilters) <= 0 {
		multiFilters = append(multiFilters, nil)
	}
	for _, filters := range multiFilters {
		tmpCenBwps, err := doRequestCenBandwidthPackages(filters, d, meta)
		if err != nil {
			return err
		}
		if len(tmpCenBwps) > 0 {
			allCenBwps = append(allCenBwps, tmpCenBwps...)
		}
	}

	return cenBandwidthPackageAttributes(d, allCenBwps)
}

func doRequestCenBandwidthPackages(filters []cbn.DescribeCenBandwidthPackagesFilter, d *schema.ResourceData, meta interface{}) ([]cbn.CenBandwidthPackage, error) {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	// filter by instance_id
	if instanceId, ok := d.GetOk("instance_id"); ok {
		filters = append(filters, cbn.DescribeCenBandwidthPackagesFilter{
			Key:   "CenId",
			Value: &[]string{instanceId.(string)},
		})
	}

	if filters != nil {
		request.Filter = &filters
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}

	var allCenBwps []cbn.CenBandwidthPackage

	for {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
		})
		if err != nil {
			return allCenBwps, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_bandwidth_packages", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)

		if len(response.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			break
		}

		for _, e := range response.CenBandwidthPackages.CenBandwidthPackage {

			// filter by name
			if nameRegex != nil {
				if !nameRegex.MatchString(e.Name) {
					continue
				}
			}
			allCenBwps = append(allCenBwps, e)
		}

		if len(response.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return allCenBwps, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return allCenBwps, nil
}

func constructCenBwpRequestFilters(d *schema.ResourceData) [][]cbn.DescribeCenBandwidthPackagesFilter {
	var res [][]cbn.DescribeCenBandwidthPackagesFilter
	maxQueryItem := 5
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
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
				Key:   "CenBandwidthPackageId",
				Value: &filtersArr[k],
			},
			}
			res = append(res, filters)
		}
	}

	return res
}

func cenBandwidthPackageAttributes(d *schema.ResourceData, allCenBwps []cbn.CenBandwidthPackage) error {
	var ids []string
	var names []string
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
		names = append(names, cenBwp.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("packages", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
