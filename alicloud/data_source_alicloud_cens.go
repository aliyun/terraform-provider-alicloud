package alicloud

import (
	"fmt"

	"regexp"

	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCens() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCensRead,

		Schema: map[string]*schema.Schema{
			"cen_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"cen_bandwidth_package_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"cens": {
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
						"cen_bandwidth_package_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_ids": {
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

func dataSourceAlicloudCensRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).cenconn

	args := cbn.CreateDescribeCensRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	pageNumber := 1
	args.PageNumber = requests.NewInteger(pageNumber)

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

	var allCens []cbn.Cen

	for {
		resp, err := conn.DescribeCens(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Cens.Cen) < 1 {
			break
		}

		for _, e := range resp.Cens.Cen {

			// filter by cenIds
			if len(cenIdsMap) > 0 {
				if _, ok := cenIdsMap[e.CenId]; !ok {
					continue
				}
			}

			// filter by cenBandwidthPackageIds
			if len(cenBandwidthPackageIdsMap) > 0 {
				isContainCenBandwidthPackageId := false
				cenBandwidthPackageIds := e.CenBandwidthPackageIds.CenBandwidthPackageId
				for _, cenBandwidthPackageId := range cenBandwidthPackageIds {
					if _, ok := cenBandwidthPackageIdsMap[cenBandwidthPackageId]; ok {
						isContainCenBandwidthPackageId = true
					}
				}
				if !isContainCenBandwidthPackageId {
					continue
				}
			}

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

		pageNumber++
		args.PageNumber = requests.NewInteger(pageNumber)
	}

	if len(allCens) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_cens - Cens found: %#v", allCens)

	return censDescriptionAttributes(d, allCens, meta)
}

func censDescriptionAttributes(d *schema.ResourceData, cenSetTypes []cbn.Cen, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, cen := range cenSetTypes {
		mapping := map[string]interface{}{
			"id":                        cen.CenId,
			"name":                      cen.Name,
			"status":                    cen.Status,
			"cen_bandwidth_package_ids": cen.CenBandwidthPackageIds.CenBandwidthPackageId,
			"description":               cen.Description,
		}

		// get child instances
		instanceIds, err := censDescribeCenAttachedChildInstances(d, cen.CenId, meta)
		if err != nil {
			return err
		}
		if instanceIds != nil && len(instanceIds) > 0 {
			mapping["instance_ids"] = instanceIds
		}

		//log.Printf("[DEBUG] alicloud_cen - adding cen: %v", mapping)
		ids = append(ids, cen.CenId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("cens", s); err != nil {
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
