// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMaxComputeQuotaPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeQuotaPlanCreate,
		Read:   resourceAliCloudMaxComputeQuotaPlanRead,
		Update: resourceAliCloudMaxComputeQuotaPlanUpdate,
		Delete: resourceAliCloudMaxComputeQuotaPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"is_effective": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sub_quota_info_list": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"elastic_reserved_cu": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"min_cu": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"max_cu": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"nick_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"parameter": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"elastic_reserved_cu": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"min_cu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_cu": {
										Type:     schema.TypeInt,
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

func resourceAliCloudMaxComputeQuotaPlanCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	nickname := d.Get("nickname")
	action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaPlan", nickname)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("plan_name"); ok {
		request["name"] = v
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("quota"); !IsNil(v) {
		parameter := make(map[string]interface{})
		elasticReservedCu, _ := jsonpath.Get("$[0].parameter[0].elastic_reserved_cu", v)
		if elasticReservedCu != nil && elasticReservedCu != "" {
			parameter["elasticReservedCU"] = elasticReservedCu
		}

		objectDataLocalMap["parameter"] = parameter
		if v, ok := d.GetOk("quota"); ok {
			localData, err := jsonpath.Get("$[0].sub_quota_info_list", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				localData1 := make(map[string]interface{})
				minCu, _ := jsonpath.Get("$[0].min_cu", dataLoopTmp["parameter"])
				if minCu != nil && minCu != "" {
					localData1["minCU"] = minCu
				}
				maxCu, _ := jsonpath.Get("$[0].max_cu", dataLoopTmp["parameter"])
				if maxCu != nil && maxCu != "" {
					localData1["maxCU"] = maxCu
				}
				elasticReservedCu1, _ := jsonpath.Get("$[0].elastic_reserved_cu", dataLoopTmp["parameter"])
				if elasticReservedCu1 != nil && elasticReservedCu1 != "" {
					localData1["elasticReservedCU"] = elasticReservedCu1
				}
				dataLoopMap["parameter"] = localData1
				dataLoopMap["nickName"] = dataLoopTmp["nick_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["subQuotaInfoList"] = localMaps
		}

		request["quota"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("MaxCompute", "2022-01-04", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_quota_plan", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", nickname, request["name"]))

	return resourceAliCloudMaxComputeQuotaPlanUpdate(d, meta)
}

func resourceAliCloudMaxComputeQuotaPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeQuotaPlan(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_quota_plan DescribeMaxComputeQuotaPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["isEffective"] != nil {
		d.Set("is_effective", objectRaw["isEffective"])
	}
	if objectRaw["name"] != nil {
		d.Set("plan_name", objectRaw["name"])
	}

	quotaMaps := make([]map[string]interface{}, 0)
	quotaMap := make(map[string]interface{})
	quota1Raw := make(map[string]interface{})
	if objectRaw["quota"] != nil {
		quota1Raw = objectRaw["quota"].(map[string]interface{})
	}
	if len(quota1Raw) > 0 {

		parameterMaps := make([]map[string]interface{}, 0)
		parameterMap := make(map[string]interface{})
		parameter2Raw := make(map[string]interface{})
		if quota1Raw["parameter"] != nil {
			parameter2Raw = quota1Raw["parameter"].(map[string]interface{})
		}
		if len(parameter2Raw) > 0 {
			parameterMap["elastic_reserved_cu"] = parameter2Raw["elasticReservedCU"]
			parameterMap["max_cu"] = parameter2Raw["maxCU"]
			parameterMap["min_cu"] = parameter2Raw["minCU"]

			parameterMaps = append(parameterMaps, parameterMap)
		}
		quotaMap["parameter"] = parameterMaps
		subQuotaInfoList1Raw := quota1Raw["subQuotaInfoList"]
		subQuotaInfoListMaps := make([]map[string]interface{}, 0)
		if subQuotaInfoList1Raw != nil {
			for _, subQuotaInfoListChild1Raw := range subQuotaInfoList1Raw.([]interface{}) {
				subQuotaInfoListMap := make(map[string]interface{})
				subQuotaInfoListChild1Raw := subQuotaInfoListChild1Raw.(map[string]interface{})
				subQuotaInfoListMap["nick_name"] = subQuotaInfoListChild1Raw["nickName"]

				parameterMaps := make([]map[string]interface{}, 0)
				parameterMap := make(map[string]interface{})
				parameter3Raw := make(map[string]interface{})
				if subQuotaInfoListChild1Raw["parameter"] != nil {
					parameter3Raw = subQuotaInfoListChild1Raw["parameter"].(map[string]interface{})
				}
				if len(parameter3Raw) > 0 {
					parameterMap["elastic_reserved_cu"] = parameter3Raw["elasticReservedCU"]
					parameterMap["max_cu"] = parameter3Raw["maxCU"]
					parameterMap["min_cu"] = parameter3Raw["minCU"]

					parameterMaps = append(parameterMaps, parameterMap)
				}
				subQuotaInfoListMap["parameter"] = parameterMaps
				subQuotaInfoListMaps = append(subQuotaInfoListMaps, subQuotaInfoListMap)
			}
		}
		quotaMap["sub_quota_info_list"] = subQuotaInfoListMaps
		quotaMaps = append(quotaMaps, quotaMap)
	}
	if objectRaw["quota"] != nil {
		if err := d.Set("quota", quotaMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("nickname", parts[0])

	return nil
}

func resourceAliCloudMaxComputeQuotaPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	if d.HasChange("is_effective") {
		maxComputeServiceV2 := MaxComputeServiceV2{client}
		object, err := maxComputeServiceV2.DescribeMaxComputeQuotaPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("is_effective").(bool)
		if object["isEffective"].(bool) != target {
			if target == true {
				parts := strings.Split(d.Id(), ":")
				planName := parts[1]
				nickname := parts[0]
				action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaPlan/%s/apply", nickname, planName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	parts := strings.Split(d.Id(), ":")
	nickname := parts[0]
	action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaPlan", nickname)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["name"] = parts[1]

	if d.HasChange("quota") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("quota"); v != nil {
		parameter := make(map[string]interface{})
		elasticReservedCu, _ := jsonpath.Get("$[0].parameter[0].elastic_reserved_cu", v)
		if elasticReservedCu != nil && (d.HasChange("quota.0.parameter.0.elastic_reserved_cu") || elasticReservedCu != "") {
			parameter["elasticReservedCU"] = elasticReservedCu
		}

		objectDataLocalMap["parameter"] = parameter
		if v, ok := d.GetOk("quota"); ok {
			localData, err := jsonpath.Get("$[0].sub_quota_info_list", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				if !IsNil(dataLoopTmp["parameter"]) {
					localData1 := make(map[string]interface{})
					elasticReservedCu1, _ := jsonpath.Get("$[0].elastic_reserved_cu", dataLoopTmp["parameter"])
					if elasticReservedCu1 != nil && elasticReservedCu1 != "" {
						localData1["elasticReservedCU"] = elasticReservedCu1
					}
					maxCu, _ := jsonpath.Get("$[0].max_cu", dataLoopTmp["parameter"])
					if maxCu != nil && maxCu != "" {
						localData1["maxCU"] = maxCu
					}
					minCu, _ := jsonpath.Get("$[0].min_cu", dataLoopTmp["parameter"])
					if minCu != nil && minCu != "" {
						localData1["minCU"] = minCu
					}
					dataLoopMap["parameter"] = localData1
				}
				dataLoopMap["nickName"] = dataLoopTmp["nick_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["subQuotaInfoList"] = localMaps
		}

		request["quota"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudMaxComputeQuotaPlanRead(d, meta)
}

func resourceAliCloudMaxComputeQuotaPlanDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	planName := parts[1]
	nickname := parts[0]
	action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaPlan/%s", nickname, planName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("MaxCompute", "2022-01-04", action, query, nil, nil, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
