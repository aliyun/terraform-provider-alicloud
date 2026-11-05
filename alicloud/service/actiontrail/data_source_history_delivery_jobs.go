package actiontrail

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceHistoryDeliveryJobs() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceHistoryDeliveryJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"history_delivery_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"home_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func DataSourceHistoryDeliveryJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDeliveryHistoryJobs"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
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
	status, statusOk := d.GetOkExists("status")
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_actiontrail_history_delivery_jobs", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DeliveryHistoryJobs", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.DeliveryHistoryJobs", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["JobId"])]; !ok {
					continue
				}
			}
			if statusOk {
				jobStatusVal, jobStatusErr := convert.FormatInt(item["JobStatus"])
				if jobStatusErr != nil {
					return jobStatusErr
				}
				if status.(int) != jobStatusVal {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		statusVal, err := convert.FormatInt(object["JobStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"create_time":             object["CreatedTime"],
			"end_time":                object["EndTime"],
			"id":                      fmt.Sprint(object["JobId"]),
			"history_delivery_job_id": fmt.Sprint(object["JobId"]),
			"home_region":             object["HomeRegion"],
			"start_time":              object["StartTime"],
			"status":                  statusVal,
			"trail_name":              object["TrailName"],
			"updated_time":            object["UpdatedTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["JobId"])
		actiontrailService := NewActiontrailService(client)
		getResp, err := actiontrailService.DescribeActiontrailHistoryDeliveryJob(id)
		if err != nil {
			return errs.WrapError(err)
		}
		jobStatus := make([]map[string]interface{}, 0)
		if jobStatusList, ok := getResp["Status"]; ok {
			jobStatusListArg, _ := jobStatusList.([]interface{})
			for _, v := range jobStatusListArg {
				if m1, ok := v.(map[string]interface{}); ok {

					deliveryJobStatus, err := convert.FormatInt(m1["Status"])
					if err != nil {
						return err
					}
					temp1 := map[string]interface{}{
						"region": m1["Region"],
						"status": deliveryJobStatus,
					}
					jobStatus = append(jobStatus, temp1)
				}
			}
		}
		mapping["job_status"] = jobStatus
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("jobs", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
