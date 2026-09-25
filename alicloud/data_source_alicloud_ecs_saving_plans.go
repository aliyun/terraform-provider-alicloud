package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcsSavingPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsSavingPlansRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"plans": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"saving_plan_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"committed_amount": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_family": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"offering_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"period": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"plan_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsSavingPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QuerySavingsPlansInstance"
	request := map[string]interface{}{
		"PageSize": "100",
		"PageNum":  "1",
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	var response map[string]interface{}
	var items []interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		resp, e := client.RpcPost("BssOpenApi", "2017-12-14", action, nil, request, true)
		if e != nil {
			if NeedRetry(e) {
				wait()
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		response = resp
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "data.alicloud_ecs_saving_plans", action, AlibabaCloudSdkGoERROR)
	}

	v, e := jsonpath.Get("$.Data.Items", response)
	if e == nil && v != nil {
		if rawItems, ok := v.([]interface{}); ok {
			items = rawItems
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, e := regexp.Compile(v.(string)); e == nil {
			nameRegex = r
		}
	}

	var ids, names []string
	var plans []map[string]interface{}
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		var savingPlanId string
		if v, ok := item["InstanceId"]; ok {
			savingPlanId = fmt.Sprint(v)
		}
		if savingPlanId == "" {
			continue
		}

		var name string
		if v, ok := item["Name"]; ok {
			name = fmt.Sprint(v)
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		plan := map[string]interface{}{
			"id":             savingPlanId,
			"saving_plan_id": savingPlanId,
		}
		if v, ok := item["PoolValue"]; ok {
			plan["committed_amount"] = fmt.Sprint(v)
		}
		if v, ok := item["StartTimestamp"]; ok {
			plan["create_time"] = fmt.Sprint(v)
		}
		if v, ok := item["InstanceFamily"]; ok {
			plan["instance_family"] = fmt.Sprint(v)
		}
		if v, ok := item["PayMode"]; ok {
			plan["offering_type"] = fmt.Sprint(v)
		}
		if v, ok := item["PaymentType"]; ok {
			plan["payment_type"] = fmt.Sprint(v)
		}
		if v, ok := item["SavingsType"]; ok {
			plan["plan_type"] = fmt.Sprint(v)
		}
		if v, ok := item["Region"]; ok {
			plan["region_id"] = fmt.Sprint(v)
		}
		if v, ok := item["StartTime"]; ok {
			plan["start_time"] = fmt.Sprint(v)
		}
		if v, ok := item["Cycle"]; ok {
			if cycle, ok := v.(string); ok && cycle != "" {
				if period, err := parseSavingPlanCycle(cycle); err == nil {
					plan["period"] = period
				}
			}
		}

		ids = append(ids, savingPlanId)
		names = append(names, name)
		plans = append(plans, plan)
	}

	d.SetId(dataSourceAlicloudEcsSavingPlanID(d, client))
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("plans", plans)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), plans); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func dataSourceAlicloudEcsSavingPlanID(d *schema.ResourceData, client *connectivity.AliyunClient) string {
	return client.RegionId
}
