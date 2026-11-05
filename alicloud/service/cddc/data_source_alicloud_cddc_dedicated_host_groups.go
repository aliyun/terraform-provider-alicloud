package cddc

import (
	"fmt"
	"regexp"
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

func DataSourceDedicatedHostGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedHostGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Redis", "SQLServer", "MySQL", "PostgreSQL", "MongoDB"}, false),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bastion_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cpu_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cpu_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_count_group_by_host_type": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"place_holder": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dedicated_host_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_used_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_utility": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_replace_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_used_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_utility": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedHostGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDedicatedHostGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var userNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		userNameRegex = r
	}

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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cddc_dedicated_host_groups", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.DedicatedHostGroups.DedicatedHostGroups", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.DedicatedHostGroups.DedicatedHostGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if userNameRegex != nil && !userNameRegex.MatchString(fmt.Sprint(item["DedicatedHostGroupDesc"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DedicatedHostGroupId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		cpu_allocation_ratioVal, err := convert.FormatInt(object["CpuAllocationRatio"])
		if err != nil {
			return err
		}

		disk_allocation_ratioVal, err := convert.FormatInt(object["DiskAllocationRatio"])
		if err != nil {
			return err
		}

		host_numberVal, err := convert.FormatInt(object["HostNumber"])
		if err != nil {
			return err
		}

		instance_numberVal, err := convert.FormatInt(object["InstanceNumber"])
		if err != nil {
			return err
		}

		mem_allocation_ratioVal, err := convert.FormatInt(object["MemAllocationRatio"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"allocation_policy":         object["AllocationPolicy"],
			"bastion_instance_id":       object["BastionInstanceId"],
			"cpu_allocate_ration":       object["CpuAllocateRation"],
			"cpu_allocated_amount":      object["CpuAllocatedAmount"],
			"cpu_allocation_ratio":      cpu_allocation_ratioVal,
			"create_time":               object["CreateTime"],
			"dedicated_host_group_desc": object["DedicatedHostGroupDesc"],
			"id":                        fmt.Sprint(object["DedicatedHostGroupId"]),
			"dedicated_host_group_id":   fmt.Sprint(object["DedicatedHostGroupId"]),
			"deploy_type":               object["DeployType"],
			"disk_allocate_ration":      object["DiskAllocateRation"],
			"disk_allocated_amount":     object["DiskAllocatedAmount"],
			"disk_allocation_ratio":     disk_allocation_ratioVal,
			"disk_used_amount":          object["DiskUsedAmount"],
			"disk_utility":              object["DiskUtility"],
			"engine":                    convertDedicatedHostGroupEngineResponse(object["Engine"]),
			"host_number":               host_numberVal,
			"host_replace_policy":       object["HostReplacePolicy"],
			"instance_number":           instance_numberVal,
			"mem_allocate_ration":       object["MemAllocateRation"],
			"mem_allocated_amount":      object["MemAllocatedAmount"],
			"mem_allocation_ratio":      mem_allocation_ratioVal,
			"mem_used_amount":           object["MemUsedAmount"],
			"mem_utility":               object["MemUtility"],
			"text":                      object["Text"],
			"vpc_id":                    object["VPCId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DedicatedHostGroupDesc"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
