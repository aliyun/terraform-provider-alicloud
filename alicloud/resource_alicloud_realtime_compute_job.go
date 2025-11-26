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

func resourceAliCloudRealtimeComputeJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeJobCreate,
		Read:   resourceAliCloudRealtimeComputeJobRead,
		Update: resourceAliCloudRealtimeComputeJobUpdate,
		Delete: resourceAliCloudRealtimeComputeJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"deployment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_variables": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"savepoint_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"kind": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"NONE", "LATEST_SAVEPOINT", "FROM_SAVEPOINT", "LATEST_STATE"}, false),
						},
						"allow_non_restored_state": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"job_start_time_in_ms": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"current_job_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"running": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"observed_flink_job_restarts": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"observed_flink_job_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"risk_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failure": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"failed_at": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"stop_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"NONE", "STOP_WITH_SAVEPOINT", "STOP_WITH_DRAIN"}, false),
			},
		},
	}
}

func resourceAliCloudRealtimeComputeJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	action := fmt.Sprintf("/api/v2/namespaces/%s/jobs:start", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("resource_id").(string))

	if v, ok := d.GetOk("local_variables"); ok {
		localVariablesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["value"] = dataLoopTmp["value"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			localVariablesMapsArray = append(localVariablesMapsArray, dataLoopMap)
		}
		request["localVariables"] = localVariablesMapsArray
	}

	restoreStrategy := make(map[string]interface{})

	if v := d.Get("restore_strategy"); !IsNil(v) {
		kind1, _ := jsonpath.Get("$[0].kind", v)
		if kind1 != nil && kind1 != "" {
			restoreStrategy["kind"] = kind1
		}
		savepointId1, _ := jsonpath.Get("$[0].savepoint_id", v)
		if savepointId1 != nil && savepointId1 != "" {
			restoreStrategy["savepointId"] = savepointId1
		}
		allowNonRestoredState1, _ := jsonpath.Get("$[0].allow_non_restored_state", v)
		if allowNonRestoredState1 != nil && allowNonRestoredState1 != "" {
			restoreStrategy["allowNonRestoredState"] = allowNonRestoredState1
		}
		jobStartTimeInMs1, _ := jsonpath.Get("$[0].job_start_time_in_ms", v)
		if jobStartTimeInMs1 != nil && jobStartTimeInMs1 != "" {
			restoreStrategy["jobStartTimeInMs"] = jobStartTimeInMs1
		}

		request["restoreStrategy"] = restoreStrategy
	}

	if v, ok := d.GetOk("deployment_id"); ok {
		request["deploymentId"] = v
	}
	if v, ok := d.GetOk("resource_queue_name"); ok {
		request["resourceQueueName"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_job", action, AlibabaCloudSdkGoERROR)
	}

	dataworkspaceVar, _ := jsonpath.Get("$.data.workspace", response)
	datanamespaceVar, _ := jsonpath.Get("$.data.namespace", response)
	datajobIdVar, _ := jsonpath.Get("$.data.jobId", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", dataworkspaceVar, datanamespaceVar, datajobIdVar))

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, realtimeComputeServiceV2.RealtimeComputeJobStateRefreshFunc(d.Id(), "$.status.currentJobStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRealtimeComputeJobUpdate(d, meta)
}

func resourceAliCloudRealtimeComputeJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_job DescribeRealtimeComputeJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("deployment_id", objectRaw["deploymentId"])
	d.Set("job_id", objectRaw["jobId"])
	d.Set("namespace", objectRaw["namespace"])
	d.Set("resource_id", objectRaw["workspace"])

	localVariablesRaw := objectRaw["localVariables"]
	localVariablesMaps := make([]map[string]interface{}, 0)
	if localVariablesRaw != nil {
		for _, localVariablesChildRaw := range convertToInterfaceArray(localVariablesRaw) {
			localVariablesMap := make(map[string]interface{})
			localVariablesChildRaw := localVariablesChildRaw.(map[string]interface{})
			localVariablesMap["name"] = localVariablesChildRaw["name"]
			localVariablesMap["value"] = localVariablesChildRaw["value"]

			localVariablesMaps = append(localVariablesMaps, localVariablesMap)
		}
	}
	if err := d.Set("local_variables", localVariablesMaps); err != nil {
		return err
	}
	restoreStrategyMaps := make([]map[string]interface{}, 0)
	restoreStrategyMap := make(map[string]interface{})
	restoreStrategyRaw := make(map[string]interface{})
	if objectRaw["restoreStrategy"] != nil {
		restoreStrategyRaw = objectRaw["restoreStrategy"].(map[string]interface{})
	}
	if len(restoreStrategyRaw) > 0 {
		restoreStrategyMap["allow_non_restored_state"] = restoreStrategyRaw["allowNonRestoredState"]
		restoreStrategyMap["job_start_time_in_ms"] = restoreStrategyRaw["jobStartTimeInMs"]
		restoreStrategyMap["kind"] = restoreStrategyRaw["kind"]
		restoreStrategyMap["savepoint_id"] = restoreStrategyRaw["savepointId"]

		restoreStrategyMaps = append(restoreStrategyMaps, restoreStrategyMap)
	}
	if err := d.Set("restore_strategy", restoreStrategyMaps); err != nil {
		return err
	}
	statusMaps := make([]map[string]interface{}, 0)
	statusMap := make(map[string]interface{})
	statusRaw := make(map[string]interface{})
	if objectRaw["status"] != nil {
		statusRaw = objectRaw["status"].(map[string]interface{})
	}
	if len(statusRaw) > 0 {
		statusMap["current_job_status"] = statusRaw["currentJobStatus"]
		statusMap["health_score"] = statusRaw["healthScore"]
		statusMap["risk_level"] = statusRaw["riskLevel"]

		failureMaps := make([]map[string]interface{}, 0)
		failureMap := make(map[string]interface{})
		failureRaw := make(map[string]interface{})
		if statusRaw["failure"] != nil {
			failureRaw = statusRaw["failure"].(map[string]interface{})
		}
		if len(failureRaw) > 0 {
			failureMap["failed_at"] = failureRaw["failedAt"]
			failureMap["message"] = failureRaw["message"]
			failureMap["reason"] = failureRaw["reason"]

			failureMaps = append(failureMaps, failureMap)
		}
		statusMap["failure"] = failureMaps
		runningMaps := make([]map[string]interface{}, 0)
		runningMap := make(map[string]interface{})
		runningRaw := make(map[string]interface{})
		if statusRaw["running"] != nil {
			runningRaw = statusRaw["running"].(map[string]interface{})
		}
		if len(runningRaw) > 0 {
			runningMap["observed_flink_job_restarts"] = runningRaw["observedFlinkJobRestarts"]
			runningMap["observed_flink_job_status"] = runningRaw["observedFlinkJobStatus"]

			runningMaps = append(runningMaps, runningMap)
		}
		statusMap["running"] = runningMaps
		statusMaps = append(statusMaps, statusMap)
	}
	if err := d.Set("status", statusMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudRealtimeComputeJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	objectRaw, _ := realtimeComputeServiceV2.DescribeRealtimeComputeJob(d.Id())

	if d.HasChange("status.0.current_job_status") {
		var err error
		target := d.Get("status.0.current_job_status").(string)

		currentStatus, err := jsonpath.Get("$.status.currentJobStatus", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.status.currentJobStatus", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "CANCELLED" {
				parts := strings.Split(d.Id(), ":")
				namespace := parts[1]
				jobId := parts[2]
				action := fmt.Sprintf("/api/v2/namespaces/%s/jobs/%s:stop", namespace, jobId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				header = make(map[string]*string)
				body = make(map[string]interface{})
				header["workspace"] = StringPointer(parts[0])

				request["stopStrategy"] = d.Get("stop_strategy")
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
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
				realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"CANCELLED"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, realtimeComputeServiceV2.RealtimeComputeJobStateRefreshFunc(d.Id(), "$.status.currentJobStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	return resourceAliCloudRealtimeComputeJobRead(d, meta)
}

func resourceAliCloudRealtimeComputeJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	jobId := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/jobs/%s", namespace, jobId)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
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
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
