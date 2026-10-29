package ess

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func dataSourceAlicloudEssLifecycleHooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssLifecycleHooksRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"hooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"heartbeat_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lifecycle_transition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_metadata": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssLifecycleHooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDescribeLifecycleHooksRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}
	var allLifecycleHooks []ess.LifecycleHook

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeLifecycleHooks(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ess_lifecycle_hooks", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeLifecycleHooksResponse)
		if len(response.LifecycleHooks.LifecycleHook) < 1 {
			break
		}
		allLifecycleHooks = append(allLifecycleHooks, response.LifecycleHooks.LifecycleHook...)
		if len(response.LifecycleHooks.LifecycleHook) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredLifecycleHooks = make([]ess.LifecycleHook, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, hook := range allLifecycleHooks {
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return errs.WrapError(err)
				}
				if r != nil && !r.MatchString(hook.LifecycleHookName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[hook.LifecycleHookId]; !ok {
					continue
				}
			}
			filteredLifecycleHooks = append(filteredLifecycleHooks, hook)
		}
	} else {
		filteredLifecycleHooks = allLifecycleHooks
	}

	return lifecycleHooksDescriptionAttribute(d, filteredLifecycleHooks, meta)
}

func lifecycleHooksDescriptionAttribute(d *schema.ResourceData, lifecycleHooks []ess.LifecycleHook, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, hook := range lifecycleHooks {
		mapping := map[string]interface{}{
			"id":                    hook.LifecycleHookId,
			"name":                  hook.LifecycleHookName,
			"scaling_group_id":      hook.ScalingGroupId,
			"default_result":        hook.DefaultResult,
			"heartbeat_timeout":     hook.HeartbeatTimeout,
			"lifecycle_transition":  hook.LifecycleTransition,
			"notification_arn":      hook.NotificationArn,
			"notification_metadata": hook.NotificationMetadata,
		}
		ids = append(ids, hook.LifecycleHookId)
		names = append(names, hook.LifecycleHookName)
		s = append(s, mapping)
	}
	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("hooks", s); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
