package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudSimpleApplicationServerPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSimpleApplicationServerPlansRead,
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

func dataSourceAlicloudSimpleApplicationServerPlansRead(d *schema.ResourceData, meta interface{}) error {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_simple_application_server_plans", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Plans", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Plans", response)
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
			bandwidthVal, bandwidthErr := formatInt(item["Bandwidth"])
			if bandwidthErr != nil {
				return bandwidthErr
			}
			if bandwidth.(int) != bandwidthVal {
				continue
			}
		}

		if coreOk && core.(int) != 0 {
			coreVal, coreErr := formatInt(item["Core"])
			if coreErr != nil {
				return coreErr
			}
			if core.(int) != coreVal {
				continue
			}
		}

		if diskSizeOk && diskSize.(int) != 0 {
			diskSizeVal, diskSizeErr := formatInt(item["DiskSize"])
			if diskSizeErr != nil {
				return diskSizeErr
			}
			if diskSize.(int) != diskSizeVal {
				continue
			}
		}

		if memoryOk && memory.(float64) != 0 {
			memoryVal, memoryErr := formatFloat64(item["Memory"])
			if memoryErr != nil {
				return memoryErr
			}
			if memory.(float64) != memoryVal {
				continue
			}
		}

		if flowOk && flow.(int) != 0 {
			flowVal, flowErr := formatInt(item["Flow"])
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

		bandwidthVal, err183 := formatInt(object["Bandwidth"])
		if err183 != nil {
			return err183
		}
		coreVal, err184 := formatInt(object["Core"])
		if err184 != nil {
			return err184
		}
		disk_sizeVal, err185 := formatInt(object["DiskSize"])
		if err185 != nil {
			return err185
		}
		flowVal, err186 := formatInt(object["Flow"])
		if err186 != nil {
			return err186
		}

		memoryVal, err189 := formatFloat64(object["Memory"])
		if err189 != nil {
			return err189
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

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
