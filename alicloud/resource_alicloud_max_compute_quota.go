package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMaxComputeQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeQuotaCreate,
		Read:   resourceAliCloudMaxComputeQuotaRead,
		Update: resourceAliCloudMaxComputeQuotaUpdate,
		Delete: resourceAliCloudMaxComputeQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"commodity_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commodity_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"part_nick_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"sub_quota_info_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"parameter": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"force_reserved_min": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"single_job_cu_limit": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"enable_priority": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"min_cu": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"scheduler_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudMaxComputeQuotaCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/quotas")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("payment_type"); ok {
		query["chargeType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("commodity_code"); ok {
		query["commodityCode"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("commodity_data"); ok {
		query["commodityData"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("part_nick_name"); ok {
		query["partNickName"] = StringPointer(v.(string))
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_quota", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.nickName", response)
	if id == nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_quota", action, AlibabaCloudSdkGoERROR, response)
	}
	d.SetId(fmt.Sprint(id))

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(id)}, d.Timeout(schema.TimeoutCreate), 30*time.Second, maxComputeServiceV2.MaxComputeQuotaStateRefreshFunc(d.Id(), "$.nickName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudMaxComputeQuotaUpdate(d, meta)
}

func resourceAliCloudMaxComputeQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeQuota(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_quota DescribeMaxComputeQuota Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	billingPolicyRawObj, _ := jsonpath.Get("$.billingPolicy", objectRaw)
	billingPolicyRaw := make(map[string]interface{})
	if billingPolicyRawObj != nil {
		billingPolicyRaw = billingPolicyRawObj.(map[string]interface{})
	}
	d.Set("payment_type", convertMaxComputeQuotabillingPolicybillingMethodResponse(billingPolicyRaw["billingMethod"]))

	subQuotaInfoListRaw := objectRaw["subQuotaInfoList"]
	subQuotaInfoListMaps := make([]map[string]interface{}, 0)
	if subQuotaInfoListRaw != nil {
		for _, subQuotaInfoListChildRaw := range subQuotaInfoListRaw.([]interface{}) {
			subQuotaInfoListMap := make(map[string]interface{})
			subQuotaInfoListChildRaw := subQuotaInfoListChildRaw.(map[string]interface{})
			subQuotaInfoListMap["nick_name"] = subQuotaInfoListChildRaw["nickName"]
			subQuotaInfoListMap["type"] = subQuotaInfoListChildRaw["type"]

			parameterMaps := make([]map[string]interface{}, 0)
			parameterMap := make(map[string]interface{})
			parameterRaw := make(map[string]interface{})
			if subQuotaInfoListChildRaw["parameter"] != nil {
				parameterRaw = subQuotaInfoListChildRaw["parameter"].(map[string]interface{})
			}
			if len(parameterRaw) > 0 {
				parameterMap["enable_priority"] = parameterRaw["enablePriority"]
				parameterMap["force_reserved_min"] = parameterRaw["forceReservedMin"]
				parameterMap["max_cu"] = parameterRaw["maxCU"]
				parameterMap["min_cu"] = parameterRaw["minCU"]
				parameterMap["scheduler_type"] = parameterRaw["schedulerType"]
				parameterMap["single_job_cu_limit"] = parameterRaw["singleJobCULimit"]

				parameterMaps = append(parameterMaps, parameterMap)
			}
			subQuotaInfoListMap["parameter"] = parameterMaps
			subQuotaInfoListMaps = append(subQuotaInfoListMaps, subQuotaInfoListMap)
		}
	}
	if err := d.Set("sub_quota_info_list", subQuotaInfoListMaps); err != nil {
		return err
	}

	objectRaw, err = maxComputeServiceV2.DescribeQuotaListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudMaxComputeQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	nickname := d.Id()
	action := fmt.Sprintf("/api/v1/quotas/%s/computeSubQuota", nickname)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header := make(map[string]*string)
	body = make(map[string]interface{})
	request["nickname"] = d.Id()

	if d.HasChange("sub_quota_info_list") {
		update = true
	}
	if v, ok := d.GetOk("sub_quota_info_list"); ok || d.HasChange("sub_quota_info_list") {
		subQuotaInfoListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["nickName"] = dataLoopTmp["nick_name"]
			if !IsNil(dataLoopTmp["parameter"]) {
				localData1 := make(map[string]interface{})
				minCu, _ := jsonpath.Get("$[0].min_cu", dataLoopTmp["parameter"])
				if minCu != nil && minCu != "" {
					localData1["minCU"] = minCu
				}
				singleJobCuLimit, _ := jsonpath.Get("$[0].single_job_cu_limit", dataLoopTmp["parameter"])
				if singleJobCuLimit != nil && singleJobCuLimit != "" && singleJobCuLimit.(int) > 0 {
					localData1["singleJobCULimit"] = singleJobCuLimit
				}
				enablePriority1, _ := jsonpath.Get("$[0].enable_priority", dataLoopTmp["parameter"])
				if enablePriority1 != nil && enablePriority1 != "" {
					localData1["enablePriority"] = enablePriority1
				}
				maxCu, _ := jsonpath.Get("$[0].max_cu", dataLoopTmp["parameter"])
				if maxCu != nil && maxCu != "" {
					localData1["maxCU"] = maxCu
				}
				schedulerType1, _ := jsonpath.Get("$[0].scheduler_type", dataLoopTmp["parameter"])
				if schedulerType1 != nil && schedulerType1 != "" {
					localData1["schedulerType"] = schedulerType1
				}
				forceReservedMin1, _ := jsonpath.Get("$[0].force_reserved_min", dataLoopTmp["parameter"])
				if forceReservedMin1 != nil && forceReservedMin1 != "" {
					localData1["forceReservedMin"] = forceReservedMin1
				}
				dataLoopMap["parameter"] = localData1
			}
			if dataLoopTmp["type"] != "" {
				dataLoopMap["type"] = dataLoopTmp["type"]
			}
			subQuotaInfoListMapsArray = append(subQuotaInfoListMapsArray, dataLoopMap)
		}
		request["subQuotaInfoList"] = subQuotaInfoListMapsArray
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, header, body, true)
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

	return resourceAliCloudMaxComputeQuotaRead(d, meta)
}

func resourceAliCloudMaxComputeQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Quota. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertMaxComputeQuotabillingPolicybillingMethodResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "payasyougo":
		return "PayAsYouGo"
	case "subscription":
		return "Subscription"
	}
	return source
}
