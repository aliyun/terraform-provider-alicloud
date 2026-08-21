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

func dataSourceAlicloudBssOpenApiBudgets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBssOpenApiBudgetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"budgets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"budget_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"budget_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cycle_end_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cycle_start_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cycle_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cycle_quota": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cycle_period": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quota": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"query_filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"select_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"warn_confs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"msc_channels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"threshold_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"msc_contacts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"event_bridge": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"warn_target": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"threshold_type": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudBssOpenApiBudgetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	budgetsMaps := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	names := make([]string, 0)

	pageNo := 1
	pageSize := PageSizeLarge
	for {
		action := "DescribeBudgets"
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["PageNo"] = pageNo
		request["PageSize"] = pageSize
		request["ProductCode"] = ""
		request["ProductType"] = ""
		if client.IsInternationalAccount() {
			request["ProductType"] = ""
		}
		var endpoint string
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2023-09-30", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
					request["ProductCode"] = ""
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bss_open_api_budgets", action, AlibabaCloudSdkGoERROR)
		}

		dataRaw, _ := jsonpath.Get("$.data", response)
		items := convertToInterfaceArray(dataRaw)
		for _, itemRaw := range items {
			item, ok := itemRaw.(map[string]interface{})
			if !ok {
				continue
			}
			budgetName := fmt.Sprint(item["budgetName"])
			if nameRegex != nil && !nameRegex.MatchString(budgetName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, exist := idsMap[budgetName]; !exist {
					continue
				}
			}
			budgetMap := make(map[string]interface{})
			budgetMap["budget_name"] = budgetName
			budgetMap["budget_type"] = item["budgetType"]
			budgetMap["comment"] = item["comment"]
			budgetMap["cycle_end_period"] = item["cycleEndPeriod"]
			budgetMap["cycle_start_period"] = item["cycleStartPeriod"]
			budgetMap["cycle_type"] = item["cycleType"]
			budgetMap["metric"] = item["metric"]
			budgetMap["quota"] = item["quota"]
			budgetMap["quota_type"] = item["quotaType"]

			cycleQuotaRaw := item["cycleQuota"]
			cycleQuotaMaps := make([]map[string]interface{}, 0)
			if cycleQuotaRaw != nil {
				for _, childRaw := range convertToInterfaceArray(cycleQuotaRaw) {
					child, ok := childRaw.(map[string]interface{})
					if !ok {
						continue
					}
					m := make(map[string]interface{})
					m["cycle_period"] = child["cyclePeriod"]
					m["quota"] = child["quota"]
					cycleQuotaMaps = append(cycleQuotaMaps, m)
				}
			}
			budgetMap["cycle_quota"] = cycleQuotaMaps

			queryFilterRaw := item["queryFilter"]
			queryFilterMaps := make([]map[string]interface{}, 0)
			if queryFilterRaw != nil {
				for _, childRaw := range convertToInterfaceArray(queryFilterRaw) {
					child, ok := childRaw.(map[string]interface{})
					if !ok {
						continue
					}
					m := make(map[string]interface{})
					m["code"] = child["code"]
					m["select_type"] = child["selectType"]
					valuesRaw := make([]interface{}, 0)
					if child["values"] != nil {
						valuesRaw = convertToInterfaceArray(child["values"])
					}
					m["values"] = valuesRaw
					queryFilterMaps = append(queryFilterMaps, m)
				}
			}
			budgetMap["query_filter"] = queryFilterMaps

			warnConfRaw := item["warnConf"]
			warnConfsMaps := make([]map[string]interface{}, 0)
			if warnConfRaw != nil {
				for _, childRaw := range convertToInterfaceArray(warnConfRaw) {
					child, ok := childRaw.(map[string]interface{})
					if !ok {
						continue
					}
					m := make(map[string]interface{})
					m["comment"] = child["comment"]
					m["event_bridge"] = child["eventBridge"]
					m["name"] = child["name"]
					m["threshold_type"] = child["thresholdType"]
					m["threshold_value"] = child["thresholdValue"]
					m["warn_target"] = child["warnTarget"]
					mscChannelsRaw := make([]interface{}, 0)
					if child["mscChannels"] != nil {
						mscChannelsRaw = convertToInterfaceArray(child["mscChannels"])
					}
					m["msc_channels"] = mscChannelsRaw
					mscContactsRaw := make([]interface{}, 0)
					if child["mscContacts"] != nil {
						mscContactsRaw = convertToInterfaceArray(child["mscContacts"])
					}
					m["msc_contacts"] = mscContactsRaw
					warnConfsMaps = append(warnConfsMaps, m)
				}
			}
			budgetMap["warn_confs"] = warnConfsMaps

			budgetsMaps = append(budgetsMaps, budgetMap)
			ids = append(ids, budgetName)
			names = append(names, budgetName)
		}

		totalCount := 0
		if tc, _ := jsonpath.Get("$.totalCount", response); tc != nil {
			totalCount = int(tc.(float64))
		}
		if pageNo*pageSize >= totalCount || len(items) == 0 {
			break
		}
		pageNo++
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("budgets", budgetsMaps); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), budgetsMaps); err != nil {
			return WrapError(err)
		}
	}

	return nil
}
