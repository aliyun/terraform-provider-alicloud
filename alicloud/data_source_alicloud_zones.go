package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"reflect"
	"sort"
)

func dataSourceAlicloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesRead,

		Schema: map[string]*schema.Schema{
			"available_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"available_resource_creation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"available_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudSSD),
					string(ecs.DiskCategoryCloudEfficiency),
				}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_resource_creation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_disk_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	insType, _ := d.Get("available_instance_type").(string)
	resType, _ := d.Get("available_resource_creation").(string)
	diskType, _ := d.Get("available_disk_category").(string)
	validData, err := meta.(*AliyunClient).CheckParameterValidity(d, meta)
	if err != nil {
		return err
	}
	zones := make(map[string]ecs.ZoneType)
	if val, ok := validData[ZoneKey]; ok {
		zones = val.(map[string]ecs.ZoneType)
	}

	zoneTypes := make(map[string]ecs.ZoneType)
	var zoneIds []string
	for _, zone := range zones {

		if len(zone.AvailableInstanceTypes.InstanceTypes) == 0 {
			continue
		}

		if insType != "" && !constraints(zone.AvailableInstanceTypes.InstanceTypes, insType) {
			continue
		}

		if len(zone.AvailableResourceCreation.ResourceTypes) == 0 || (resType != "" && !constraints(zone.AvailableResourceCreation.ResourceTypes, resType)) {
			continue
		}

		if len(zone.AvailableDiskCategories.DiskCategories) == 0 || (diskType != "" && !constraints(zone.AvailableDiskCategories.DiskCategories, diskType)) {
			continue
		}
		zoneTypes[zone.ZoneId] = zone
		zoneIds = append(zoneIds, zone.ZoneId)
	}

	if len(zoneTypes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	// Sort zones before reading
	sort.Strings(zoneIds)

	var newZoneTypes []ecs.ZoneType
	for _, id := range zoneIds {
		newZoneTypes = append(newZoneTypes, zoneTypes[id])
	}

	log.Printf("[DEBUG] alicloud_zones - Zones found: %#v", newZoneTypes)
	return zonesDescriptionAttributes(d, newZoneTypes)
}

// check array constraints str
func constraints(arr interface{}, v string) bool {
	arrs := reflect.ValueOf(arr)
	len := arrs.Len()
	for i := 0; i < len; i++ {
		if arrs.Index(i).String() == v {
			return true
		}
	}
	return false
}

func zonesDescriptionAttributes(d *schema.ResourceData, types []ecs.ZoneType) error {
	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":                          t.ZoneId,
			"local_name":                  t.LocalName,
			"available_instance_types":    t.AvailableInstanceTypes.InstanceTypes,
			"available_resource_creation": t.AvailableResourceCreation.ResourceTypes,
			"available_disk_categories":   t.AvailableDiskCategories.DiskCategories,
		}

		log.Printf("[DEBUG] alicloud_zones - adding zone mapping: %v", mapping)
		ids = append(ids, t.ZoneId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
