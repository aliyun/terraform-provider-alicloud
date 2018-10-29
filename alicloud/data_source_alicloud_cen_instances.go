package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCenInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"child_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenInstancesRead(d *schema.ResourceData, meta interface{}) error {

	multiFilters := constructCenRequestFilters(d)
	if len(multiFilters) <= 0 {
		multiFilters = append(multiFilters, nil)
	}

	var allCens []cbn.Cen
	for _, filters := range multiFilters {
		tmpCens, err := getCenInstances(filters, d, meta)
		if err != nil {
			return err
		}
		if len(tmpCens) > 0 {
			allCens = append(allCens, tmpCens...)
		}
	}

	if len(allCens) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_cen_instances - Cens found: %#v", allCens)
	return censDescriptionAttributes(d, allCens, meta)
}

func getCenInstances(filters []cbn.DescribeCensFilter, d *schema.ResourceData, meta interface{}) ([]cbn.Cen, error) {
	client := meta.(*connectivity.AliyunClient)

	args := cbn.CreateDescribeCensRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	if filters != nil {
		args.Filter = &filters
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}

	var allCens []cbn.Cen
	for {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(args)
		})
		if err != nil {
			return allCens, err
		}
		resp, _ := raw.(*cbn.DescribeCensResponse)

		if resp == nil || len(resp.Cens.Cen) < 1 {
			break
		}

		for _, e := range resp.Cens.Cen {
			// filter by name
			if nameRegex != nil {
				if !nameRegex.MatchString(e.Name) {
					continue
				}
			}
			allCens = append(allCens, e)
		}

		if len(resp.Cens.Cen) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return allCens, err
		} else {
			args.PageNumber = page
		}
	}
	return allCens, nil
}

func constructCenRequestFilters(d *schema.ResourceData) [][]cbn.DescribeCensFilter {
	var res [][]cbn.DescribeCensFilter
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
			filters := []cbn.DescribeCensFilter{{
				Key:   "CenId",
				Value: &filtersArr[k],
			},
			}
			res = append(res, filters)
		}
	}
	return res
}

func censDescriptionAttributes(d *schema.ResourceData, cenSetTypes []cbn.Cen, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, cen := range cenSetTypes {
		mapping := map[string]interface{}{
			"id":                    cen.CenId,
			"name":                  cen.Name,
			"status":                cen.Status,
			"bandwidth_package_ids": cen.CenBandwidthPackageIds.CenBandwidthPackageId,
			"description":           cen.Description,
		}

		// get child instances
		instanceIds, err := censDescribeCenAttachedChildInstances(d, cen.CenId, meta)
		if err != nil {
			return err
		}

		if instanceIds != nil {
			mapping["child_instance_ids"] = instanceIds
		} else {
			mapping["child_instance_ids"] = []string{}
		}

		ids = append(ids, cen.CenId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func censDescribeCenAttachedChildInstances(d *schema.ResourceData, cenId string, meta interface{}) ([]string, error) {
	client := meta.(*connectivity.AliyunClient)
	var instanceIds []string

	args := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	args.CenId = cenId
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	for {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstances(args)
		})
		if err != nil {
			return nil, fmt.Errorf("DescribeCenAttachedChildInstances got an error: %#v.", err)
		}
		resp, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)
		if resp != nil && len(resp.ChildInstances.ChildInstance) > 0 {

			for _, inst := range resp.ChildInstances.ChildInstance {
				instanceIds = append(instanceIds, inst.ChildInstanceId)
			}

		}

		if len(resp.ChildInstances.ChildInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return instanceIds, err
		} else {
			args.PageNumber = page
		}
	}

	return instanceIds, nil
}
