package alicloud

import (
	"fmt"
	"strings"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"is_outdated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	zoneId, validZones, err := client.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}

	mapInstanceTypes := make(map[string][]string)
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}

					zones, _ := mapInstanceTypes[t.Value]
					zones = append(zones, zone.ZoneId)
					mapInstanceTypes[t.Value] = zones
				}
			}
		}
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := d.Get("memory_size").(float64)
	family := strings.TrimSpace(d.Get("instance_type_family").(string))

	resp, err := client.ecsconn.DescribeInstanceTypesNew(&ecs.DescribeInstanceTypesArgs{})
	if err != nil {
		return err
	}

	var instanceTypes []ecs.InstanceTypeItemType
	for _, types := range resp {
		if _, ok := mapInstanceTypes[types.InstanceTypeId]; !ok {
			continue
		}

		if cpu > 0 && types.CpuCoreCount != cpu {
			continue
		}

		if mem > 0 && types.MemorySize != mem {
			continue
		}

		if family != "" && types.InstanceTypeFamily != family {
			continue
		}

		instanceTypes = append(instanceTypes, types)
	}

	if len(instanceTypes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	return instanceTypesDescriptionAttributes(d, instanceTypes, mapInstanceTypes)
}

func instanceTypesDescriptionAttributes(d *schema.ResourceData, types []ecs.InstanceTypeItemType, mapTypes map[string][]string) error {
	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":             t.InstanceTypeId,
			"cpu_core_count": t.CpuCoreCount,
			"memory_size":    t.MemorySize,
			"family":         t.InstanceTypeFamily,
		}

		ids = append(ids, t.InstanceTypeId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instance_types", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
