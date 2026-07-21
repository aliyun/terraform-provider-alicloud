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

func resourceAliCloudMaxComputeQuotaSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeQuotaScheduleCreate,
		Read:   resourceAliCloudMaxComputeQuotaScheduleRead,
		Update: resourceAliCloudMaxComputeQuotaScheduleUpdate,
		Delete: resourceAliCloudMaxComputeQuotaScheduleDelete,
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
			"schedule_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"at": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"plan": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMaxComputeQuotaScheduleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	nickname := d.Get("nickname")
	action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaSchedule", nickname)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["scheduleTimezone"] = StringPointer(d.Get("timezone").(string))

	if v, ok := d.GetOk("schedule_list"); ok {
		bodyMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			localData1 := make(map[string]interface{})
			at1, _ := jsonpath.Get("$[0].at", dataLoopTmp["condition"])
			if at1 != nil && at1 != "" {
				localData1["at"] = at1
			}
			dataLoopMap["condition"] = localData1
			dataLoopMap["plan"] = dataLoopTmp["plan"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			bodyMapsArray = append(bodyMapsArray, dataLoopMap)
		}
		request["body"] = bodyMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body["body"], true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_quota_schedule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", nickname, *query["scheduleTimezone"]))

	return resourceAliCloudMaxComputeQuotaScheduleRead(d, meta)
}

func resourceAliCloudMaxComputeQuotaScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeQuotaSchedule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_quota_schedule DescribeMaxComputeQuotaSchedule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	data1Raw, _ := jsonpath.Get("$.data", objectRaw)

	scheduleListMaps := make([]map[string]interface{}, 0)
	if data1Raw != nil {
		for _, dataChild1Raw := range data1Raw.([]interface{}) {
			scheduleListMap := make(map[string]interface{})
			dataChild1Raw := dataChild1Raw.(map[string]interface{})
			scheduleListMap["plan"] = dataChild1Raw["plan"]
			scheduleListMap["type"] = dataChild1Raw["type"]

			conditionMaps := make([]map[string]interface{}, 0)
			conditionMap := make(map[string]interface{})
			condition1Raw := make(map[string]interface{})
			if dataChild1Raw["condition"] != nil {
				condition1Raw = dataChild1Raw["condition"].(map[string]interface{})
			}
			if len(condition1Raw) > 0 {
				conditionMap["at"] = condition1Raw["at"]

				conditionMaps = append(conditionMaps, conditionMap)
			}
			scheduleListMap["condition"] = conditionMaps
			scheduleListMaps = append(scheduleListMaps, scheduleListMap)
		}
	}
	if objectRaw["data"] != nil {
		if err := d.Set("schedule_list", scheduleListMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("nickname", parts[0])
	d.Set("timezone", parts[1])

	return nil
}

func resourceAliCloudMaxComputeQuotaScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	nickname := parts[0]
	action := fmt.Sprintf("/api/v1/quotas/%s/computeQuotaSchedule", nickname)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["scheduleTimezone"] = StringPointer(parts[1])

	if d.HasChange("schedule_list") {
		update = true
	}
	if v, ok := d.GetOk("schedule_list"); ok || d.HasChange("schedule_list") {
		bodyMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			if !IsNil(dataLoopTmp["condition"]) {
				localData1 := make(map[string]interface{})
				at1, _ := jsonpath.Get("$[0].at", dataLoopTmp["condition"])
				if at1 != nil && at1 != "" {
					localData1["at"] = at1
				}
				dataLoopMap["condition"] = localData1
			}
			dataLoopMap["plan"] = dataLoopTmp["plan"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			bodyMapsArray = append(bodyMapsArray, dataLoopMap)
		}
		request["body"] = bodyMapsArray
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

	return resourceAliCloudMaxComputeQuotaScheduleRead(d, meta)
}

func resourceAliCloudMaxComputeQuotaScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Quota Schedule. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
