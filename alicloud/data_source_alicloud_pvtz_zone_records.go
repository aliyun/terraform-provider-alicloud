package alicloud

import (
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudPvtzZoneRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZoneRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPvtzZoneRecordsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).pvtzconn
	args := pvtz.CreateDescribeZoneRecordsRequest()
	if zoneId, ok := d.GetOk("zone_id"); ok {
		args.ZoneId = zoneId.(string)
	}

	if keyword, ok := d.GetOk("keyword"); ok {
		args.Keyword = keyword.(string)
	}

	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var pvtzZoneRecords []pvtz.Record
	recordIds := []string{}

	for true {
		results, err := conn.DescribeZoneRecords(args)

		if err != nil {
			return fmt.Errorf("Error DescribeZoneRecords: %#v", err)
		}

		if results == nil || len(results.Records.Record) < 1 {
			break
		}

		for _, key := range results.Records.Record {
			pvtzZoneRecords = append(pvtzZoneRecords, key)
			recordIds = append(recordIds, strconv.Itoa(key.RecordId))
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	d.SetId(dataResourceIdHash(recordIds))
	var zoneRecords []map[string]interface{}

	for _, r := range pvtzZoneRecords {
		mapping := map[string]interface{}{
			"id":              r.RecordId,
			"resource_record": r.Rr,
			"type":            r.Type,
			"status":          r.Status,
			"value":           r.Value,
			"ttl":             r.Ttl,
			"priority":        r.Priority,
		}

		zoneRecords = append(zoneRecords, mapping)
	}
	if err := d.Set("records", zoneRecords); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), zoneRecords)
	}

	return nil
}
