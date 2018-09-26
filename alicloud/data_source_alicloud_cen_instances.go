package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
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

	multiFilters, err := constructCenRequestFilters(d)
	if err != nil {
		return err
	}
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
	conn := meta.(*AliyunClient).cenconn

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
		resp, err := conn.DescribeCens(args)
		if err != nil {
			return allCens, err
		}

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

func constructCenRequestFilters(d *schema.ResourceData) ([][]cbn.DescribeCensFilter, error) {
	var res [][]cbn.DescribeCensFilter

	maxQueryItem := 5
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		filters := new([]cbn.DescribeCensFilter)
		filterElem := new([]string)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			*filterElem = append(*filterElem, Trim(vv.(string)))
			if len(*filterElem) >= maxQueryItem {
				*filters = append(*filters, cbn.DescribeCensFilter{
					Key:   "CenId",
					Value: filterElem,
				})
				res =  append(res, *filters)
				filters = new([]cbn.DescribeCensFilter)
				filterElem = new([]string)
			}
		}
		if len(*filterElem) > 0 {
			*filters = append(*filters, cbn.DescribeCensFilter{
				Key:   "CenId",
				Value: filterElem,
			})
			res =  append(res, *filters)
		}

	}
	return res, nil
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
		if instanceIds != nil && len(instanceIds) > 0 {
			mapping["child_instance_ids"] = instanceIds
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

	var instanceIds []string

	childInstancesReq := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	childInstancesReq.CenId = cenId
	childInstancesReq.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		resp, err := meta.(*AliyunClient).cenconn.DescribeCenAttachedChildInstances(childInstancesReq)
		if err != nil {
			return nil, fmt.Errorf("DescribeCenAttachedChildInstances got an error: %#v.", err)
		}
		if resp != nil && len(resp.ChildInstances.ChildInstance) > 0 {

			for _, inst := range resp.ChildInstances.ChildInstance {
				instanceIds = append(instanceIds, inst.ChildInstanceId)
			}

		}

		if len(resp.ChildInstances.ChildInstance) < PageSizeLarge {
			break
		}

		childInstancesReq.PageNumber += requests.NewInteger(1)
	}

	return instanceIds, nil
}
