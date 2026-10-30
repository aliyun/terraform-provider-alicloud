package swasopen

import (
	"fmt"
	"strings"
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

func DataSourcePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"core": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flow": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Linux", "Windows"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"flow": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"support_platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPlans"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
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

	bandwidth, bandwidthOk := d.GetOk("bandwidth")
	core, coreOk := d.GetOk("core")
	diskSize, diskSizeOk := d.GetOk("disk_size")
	memory, memoryOk := d.GetOk("memory")
	flow, flowOk := d.GetOk("flow")

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_simple_application_server_plans", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Plans", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Plans", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["PlanId"])]; !ok {
				continue
			}
		}

		if bandwidthOk && bandwidth.(int) != 0 {
			bandwidthVal, bandwidthErr := convert.FormatInt(item["Bandwidth"])
			if bandwidthErr != nil {
				return bandwidthErr
			}
			if bandwidth.(int) != bandwidthVal {
				continue
			}
		}

		if coreOk && core.(int) != 0 {
			coreVal, coreErr := convert.FormatInt(item["Core"])
			if coreErr != nil {
				return coreErr
			}
			if core.(int) != coreVal {
				continue
			}
		}

		if diskSizeOk && diskSize.(int) != 0 {
			diskSizeVal, diskSizeErr := convert.FormatInt(item["DiskSize"])
			if diskSizeErr != nil {
				return diskSizeErr
			}
			if diskSize.(int) != diskSizeVal {
				continue
			}
		}

		if memoryOk && memory.(float64) != 0 {
			memoryVal, memoryErr := convert.FormatFloat64(item["Memory"])
			if memoryErr != nil {
				return memoryErr
			}
			if memory.(float64) != memoryVal {
				continue
			}
		}

		if flowOk && flow.(int) != 0 {
			flowVal, flowErr := convert.FormatInt(item["Flow"])
			if flowErr != nil {
				return flowErr
			}
			if flow.(int) != flowVal {
				continue
			}
		}

		if v, ok := d.GetOk("platform"); ok && v.(string) != "" && !strings.Contains(fmt.Sprint(item["SupportPlatform"]), v.(string)) {
			continue
		}

		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		bandwidthVal, err := convert.FormatInt(object["Bandwidth"])
		if err != nil {
			return err
		}
		coreVal, err := convert.FormatInt(object["Core"])
		if err != nil {
			return err
		}
		disk_sizeVal, err := convert.FormatInt(object["DiskSize"])
		if err != nil {
			return err
		}
		flowVal, err := convert.FormatInt(object["Flow"])
		if err != nil {
			return err
		}

		memoryVal, err := convert.FormatFloat64(object["Memory"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"bandwidth":        bandwidthVal,
			"core":             coreVal,
			"disk_size":        disk_sizeVal,
			"flow":             flowVal,
			"id":               fmt.Sprint(object["PlanId"]),
			"plan_id":          fmt.Sprint(object["PlanId"]),
			"memory":           memoryVal,
			"support_platform": fmt.Sprint(object["SupportPlatform"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
