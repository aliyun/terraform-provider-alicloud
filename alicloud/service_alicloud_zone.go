package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func zoneIdsDescriptionAttributes(d *schema.ResourceData, zones []string) error {
	var s []map[string]interface{}
	var zoneIds []string
	for _, t := range zones {
		mapping := map[string]interface{}{
			"id":             t,
			"multi_zone_ids": splitMultiZoneId(t),
		}
		s = append(s, mapping)
		zoneIds = append(zoneIds, t)
	}

	d.SetId(dataResourceIdHash(zones))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	if err := d.Set("ids", zoneIds); err != nil {
		return err
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func splitMultiZoneId(id string) (ids []string) {
	if !(strings.Contains(id, MULTI_IZ_SYMBOL) || strings.Contains(id, "(")) {
		return
	}
	firstIndex := strings.Index(id, MULTI_IZ_SYMBOL)
	secondIndex := strings.Index(id, "(")
	for _, p := range strings.Split(id[secondIndex+1:len(id)-1], COMMA_SEPARATED) {
		ids = append(ids, id[:firstIndex]+string(p))
	}
	return
}
