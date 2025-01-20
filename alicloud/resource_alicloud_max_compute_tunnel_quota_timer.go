// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMaxComputeTunnelQuotaTimer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeTunnelQuotaTimerCreate,
		Read:   resourceAliCloudMaxComputeTunnelQuotaTimerRead,
		Update: resourceAliCloudMaxComputeTunnelQuotaTimerUpdate,
		Delete: resourceAliCloudMaxComputeTunnelQuotaTimerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"nickname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_timer": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tunnel_quota_parameter": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"elastic_reserved_slot_num": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"slot_num": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"begin_time": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudMaxComputeTunnelQuotaTimerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	nickname := d.Get("nickname")
	action := fmt.Sprintf("/api/v1/tunnel/%s/timers", nickname)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("nickname"); ok {
		request["nickname"] = v
	}

	if v, ok := d.GetOk("quota_timer"); ok {
		bodyMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["beginTime"] = dataLoopTmp["begin_time"]
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			localData1 := make(map[string]interface{})
			slotNum1, _ := jsonpath.Get("$[0].slot_num", dataLoopTmp["tunnel_quota_parameter"])
			if slotNum1 != nil && slotNum1 != "" {
				localData1["slotNum"] = slotNum1
			}
			elasticReservedSlotNum1, _ := jsonpath.Get("$[0].elastic_reserved_slot_num", dataLoopTmp["tunnel_quota_parameter"])
			if elasticReservedSlotNum1 != nil && elasticReservedSlotNum1 != "" {
				localData1["elasticReservedSlotNum"] = elasticReservedSlotNum1
			}
			dataLoopMap["tunnelQuotaParameter"] = localData1
			bodyMapsArray = append(bodyMapsArray, dataLoopMap)
		}
		request["body"] = bodyMapsArray
	}

	if v, ok := d.GetOk("time_zone"); ok {
		query["timezone"] = StringPointer(v.(string))
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body["body"], &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_tunnel_quota_timer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(nickname))

	return resourceAliCloudMaxComputeTunnelQuotaTimerRead(d, meta)
}

func resourceAliCloudMaxComputeTunnelQuotaTimerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeTunnelQuotaTimer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_tunnel_quota_timer DescribeMaxComputeTunnelQuotaTimer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	data1Raw, _ := jsonpath.Get("$.data", objectRaw)

	quotaTimerMaps := make([]map[string]interface{}, 0)
	if data1Raw != nil {
		for _, dataChild1Raw := range data1Raw.([]interface{}) {
			quotaTimerMap := make(map[string]interface{})
			dataChild1Raw := dataChild1Raw.(map[string]interface{})
			quotaTimerMap["begin_time"] = dataChild1Raw["beginTime"]
			quotaTimerMap["end_time"] = dataChild1Raw["endTime"]

			tunnelQuotaParameterMaps := make([]map[string]interface{}, 0)
			tunnelQuotaParameterMap := make(map[string]interface{})
			tunnelQuotaParameter1Raw := make(map[string]interface{})
			if dataChild1Raw["tunnelQuotaParameter"] != nil {
				tunnelQuotaParameter1Raw = dataChild1Raw["tunnelQuotaParameter"].(map[string]interface{})
			}
			if len(tunnelQuotaParameter1Raw) > 0 {
				tunnelQuotaParameterMap["elastic_reserved_slot_num"] = tunnelQuotaParameter1Raw["elasticReservedSlotNum"]
				tunnelQuotaParameterMap["slot_num"] = tunnelQuotaParameter1Raw["slotNum"]

				tunnelQuotaParameterMaps = append(tunnelQuotaParameterMaps, tunnelQuotaParameterMap)
			}
			quotaTimerMap["tunnel_quota_parameter"] = tunnelQuotaParameterMaps
			quotaTimerMaps = append(quotaTimerMaps, quotaTimerMap)
		}
	}
	if objectRaw["data"] != nil {
		if err := d.Set("quota_timer", quotaTimerMaps); err != nil {
			return err
		}
	}

	d.Set("nickname", d.Id())

	return nil
}

func resourceAliCloudMaxComputeTunnelQuotaTimerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	nickname := d.Id()
	action := fmt.Sprintf("/api/v1/tunnel/%s/timers", nickname)
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["nickname"] = d.Id()

	if d.HasChange("quota_timer") {
		update = true
	}
	if v, ok := d.GetOk("quota_timer"); ok || d.HasChange("quota_timer") {
		bodyMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["beginTime"] = dataLoopTmp["begin_time"]
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			if !IsNil(dataLoopTmp["tunnel_quota_parameter"]) {
				localData1 := make(map[string]interface{})
				slotNum1, _ := jsonpath.Get("$[0].slot_num", dataLoopTmp["tunnel_quota_parameter"])
				if slotNum1 != nil && slotNum1 != "" {
					localData1["slotNum"] = slotNum1
				}
				elasticReservedSlotNum1, _ := jsonpath.Get("$[0].elastic_reserved_slot_num", dataLoopTmp["tunnel_quota_parameter"])
				if elasticReservedSlotNum1 != nil && elasticReservedSlotNum1 != "" {
					localData1["elasticReservedSlotNum"] = elasticReservedSlotNum1
				}
				dataLoopMap["tunnelQuotaParameter"] = localData1
			}
			bodyMapsArray = append(bodyMapsArray, dataLoopMap)
		}
		request["body"] = bodyMapsArray
	}

	if v, ok := d.GetOk("time_zone"); ok {
		query["timezone"] = StringPointer(v.(string))
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body["body"], &runtime)
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

	return resourceAliCloudMaxComputeTunnelQuotaTimerRead(d, meta)
}

func resourceAliCloudMaxComputeTunnelQuotaTimerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Tunnel Quota Timer. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
