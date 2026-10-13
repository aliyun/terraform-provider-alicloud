package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudOpenSearchAppGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOpenSearchAppGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{`standard`, `enhanced`}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"app_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_way": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"commodity_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_rank_algo_deployment_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"has_pending_quota_review_task": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locked_by_expiration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pending_second_rank_algo_deployment_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processing_order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"produced": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compute_resource": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"doc_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_rank_algo_deployment_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"switched_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAlicloudOpenSearchAppGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/v4/openapi/app-groups"
	request := make(map[string]*string)
	if v, ok := d.GetOk("instance_id"); ok {
		request["instanceId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		request["type"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		request["name"] = StringPointer(v.(string))
	}
	request["PageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	request["PageNumber"] = StringPointer("1")
	var objects []map[string]interface{}
	var appGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		appGroupNameRegex = r
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("OpenSearch", "2017-12-25", action, request, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug("GET "+action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_open_search_app_groups", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.result", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.result", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if appGroupNameRegex != nil && !appGroupNameRegex.MatchString(fmt.Sprint(item["name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["instanceId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		charge_wayVal, err274 := formatInt(object["chargingWay"])
		if err274 != nil {
			return err274
		}

		create_timeVal, err276 := formatInt(object["created"])
		if err276 != nil {
			return err276
		}

		first_rank_algo_deployment_idVal, err281 := formatInt(object["firstRankAlgoDeploymentId"])
		if err281 != nil {
			return err281
		}
		has_pending_quota_review_taskVal, err282 := formatInt(object["hasPendingQuotaReviewTask"])
		if err282 != nil {
			return err282
		}

		locked_by_expirationVal, err285 := formatInt(object["lockedByExpiration"])
		if err285 != nil {
			return err285
		}
		pending_second_rank_algo_deployment_idVal, err286 := formatInt(object["pendingSecondRankAlgoDeploymentId"])
		if err286 != nil {
			return err286
		}

		producedVal, err288 := formatInt(object["produced"])
		if err288 != nil {
			return err288
		}

		second_rank_algo_deployment_idVal, err290 := formatInt(object["secondRankAlgoDeploymentId"])
		if err290 != nil {
			return err290
		}

		switched_timeVal, err292 := formatInt(object["switchedTime"])
		if err292 != nil {
			return err292
		}
		mapping := map[string]interface{}{
			"payment_type":                           convertOpenSearchAppGroupPaymentTypeResponse(object["chargeType"].(string)),
			"charge_way":                             charge_wayVal,
			"commodity_code":                         object["commodityCode"],
			"create_time":                            create_timeVal,
			"current_version":                        object["currentVersion"],
			"description":                            object["description"],
			"domain":                                 object["domain"],
			"expire_on":                              object["expireOn"],
			"first_rank_algo_deployment_id":          first_rank_algo_deployment_idVal,
			"has_pending_quota_review_task":          has_pending_quota_review_taskVal,
			"instance_id":                            object["instanceId"],
			"lock_mode":                              object["lockMode"],
			"locked_by_expiration":                   locked_by_expirationVal,
			"pending_second_rank_algo_deployment_id": pending_second_rank_algo_deployment_idVal,
			"processing_order_id":                    object["processingOrderId"],
			"produced":                               producedVal,
			"project_id":                             object["projectId"],
			"second_rank_algo_deployment_id":         second_rank_algo_deployment_idVal,
			"status":                                 object["status"],
			"switched_time":                          switched_timeVal,
			"type":                                   object["type"],
		}
		quotaSli := make([]map[string]interface{}, 0)
		if _, exist := object["quota"]; exist {
			quotaval := object["quota"].(map[string]interface{})
			quotaSli = append(quotaSli, map[string]interface{}{
				"doc_size":         quotaval["docSize"].(json.Number).String(),
				"compute_resource": quotaval["computeResource"].(json.Number).String(),
				"spec":             quotaval["spec"],
			})
		}
		mapping["quota"] = quotaSli
		ids = append(ids, fmt.Sprint(mapping["name"]))
		names = append(names, fmt.Sprint(mapping["name"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["name"])
		openSearchService := OpenSearchService{client}
		getResp, err := openSearchService.DescribeOpenSearchAppGroup(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["app_group_id"] = getResp["id"]
		mapping["app_group_name"] = getResp["name"]
		mapping["resource_group_id"] = getResp["resourceGroupId"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
