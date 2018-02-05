package alicloud

import (
	"fmt"
	"log"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEipsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"in_use": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"ip_addresses": {
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
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.DescribeEipAddressesArgs{
		RegionId: getRegion(d, meta),
		Status:   ecs.EipStatusAvailable,
	}
	if v, ok := d.GetOk("in_use"); ok && v.(bool) {
		args.Status = ecs.EipStatusInUse
	}
	idsMap := make(map[string]string)
	ipsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	if v, ok := d.GetOk("ip_addresses"); ok {
		for _, vv := range v.([]interface{}) {
			ipsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allEips []ecs.EipAddressSetType

	for {
		eips, paginationResult, err := conn.DescribeEipAddresses(args)
		if err != nil {
			return err
		}

		for _, e := range eips {
			if len(idsMap) > 0 {
				if _, ok := idsMap[e.AllocationId]; !ok {
					continue
				}
			}
			if len(ipsMap) > 0 {
				if _, ok := ipsMap[e.IpAddress]; !ok {
					continue
				}
			}
			allEips = append(allEips, e)
		}

		pagination := paginationResult.NextPage()
		if pagination == nil {
			break
		}

		args.Pagination = *pagination
	}

	if len(allEips) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_eips - EIPs found: %#v", allEips)

	return eipsDecriptionAttributes(d, allEips, meta)
}

func eipsDecriptionAttributes(d *schema.ResourceData, eipSetTypes []ecs.EipAddressSetType, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, eip := range eipSetTypes {
		mapping := map[string]interface{}{
			"id":                   eip.AllocationId,
			"status":               eip.Status,
			"ip_address":           eip.IpAddress,
			"bandwidth":            eip.Bandwidth,
			"instance_id":          eip.InstanceId,
			"instance_type":        eip.InstanceType,
			"internet_charge_type": eip.InternetChargeType,
			"creation_time":        eip.AllocationTime.String(),
		}
		log.Printf("[DEBUG] alicloud_eip - adding eip: %v", mapping)
		ids = append(ids, eip.AllocationId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("eips", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
