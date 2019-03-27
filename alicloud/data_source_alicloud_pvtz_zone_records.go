package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudPvtzZoneRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZoneRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
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
	client := meta.(*connectivity.AliyunClient)
	request := pvtz.CreateDescribeZoneRecordsRequest()
	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = zoneId.(string)
	}

	if keyword, ok := d.GetOk("keyword"); ok {
		request.Keyword = keyword.(string)
	}

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	var pvtzZoneRecords []pvtz.Record
	recordIds := []string{}

	invoker := PvtzInvoker()

	for true {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZoneRecords(request)
			})
			return err
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_zone_records", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw)

		response, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
		if len(response.Records.Record) < 1 {
			break
		}

		for _, key := range response.Records.Record {
			pvtzZoneRecords = append(pvtzZoneRecords, key)
			recordIds = append(recordIds, strconv.Itoa(key.RecordId))
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
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
