package privatezone

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceZoneRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRecordsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
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

func dataSourceZoneRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZoneRecords"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	request["ZoneId"] = d.Get("zone_id")

	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}

	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}

	if v, ok := d.GetOk("tag"); ok {
		request["Tag"] = v
	}

	if v, ok := d.GetOk("search_mode"); ok {
		request["SearchMode"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"System.Busy", "ServiceUnavailable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_pvtz_zone_records", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Records.Record", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Records.Record", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RecordId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		priorityVal, err := convert.FormatInt(object["Priority"])
		if err != nil {
			return err
		}

		ttlVal, err := convert.FormatInt(object["Ttl"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"id":              fmt.Sprint(object["RecordId"], ":", request["ZoneId"]),
			"record_id":       fmt.Sprint(object["RecordId"]),
			"priority":        priorityVal,
			"remark":          object["Remark"],
			"rr":              object["Rr"],
			"resource_record": object["Rr"],
			"ttl":             ttlVal,
			"type":            object["Type"],
			"value":           object["Value"],
			"status":          object["Status"],
		}

		ids = append(ids, fmt.Sprint(object["RecordId"]))

		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("records", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
