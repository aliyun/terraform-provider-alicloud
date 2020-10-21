package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
	request.RegionId = client.RegionId
	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = zoneId.(string)
	}

	if keyword, ok := d.GetOk("keyword"); ok {
		request.Keyword = keyword.(string)
	}

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	var pvtzZoneRecords []pvtz.Record
	var ids []string

	invoker := PvtzInvoker()
	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for {
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

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
		if len(response.Records.Record) < 1 {
			break
		}

		for _, key := range response.Records.Record {
			if len(idsMap) > 0 {
				if _, ok := idsMap[strconv.FormatInt(key.RecordId, 10)]; !ok {
					continue
				}
			}
			pvtzZoneRecords = append(pvtzZoneRecords, key)
			ids = append(ids, strconv.FormatInt(key.RecordId, 10))
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	d.SetId(dataResourceIdHash(ids))
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
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), zoneRecords)
	}

	return nil
}
