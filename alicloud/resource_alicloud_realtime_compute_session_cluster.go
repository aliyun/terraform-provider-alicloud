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

func resourceAliCloudRealtimeComputeSessionCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeSessionClusterCreate,
		Read:   resourceAliCloudRealtimeComputeSessionClusterRead,
		Update: resourceAliCloudRealtimeComputeSessionClusterUpdate,
		Delete: resourceAliCloudRealtimeComputeSessionClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"basic_resource_setting": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"taskmanager_resource_setting_spec": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
						"parallelism": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"jobmanager_resource_setting_spec": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_target_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flink_conf": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log4j2_configuration_template": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"logging_profile": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"log4j_loggers": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logger_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"logger_level": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"log_reserve_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"open_history": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"expiration_days": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"modified_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"session_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"running": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"reference_deployment_ids": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"started_at": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"current_session_cluster_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failure": {
							Type:     schema.TypeList,
							Computed: true,
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
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRealtimeComputeSessionClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	action := fmt.Sprintf("/api/v2/namespaces/%s/sessionclusters", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("workspace").(string))
	// The CreateSessionCluster API validates the SessionCluster body structure
	// and requires workspace and namespace in the body, in addition to the
	// workspace header and the namespace path segment.
	request["workspace"] = d.Get("workspace")
	request["namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("session_cluster_name"); ok {
		request["name"] = v
	}

	basicResourceSetting := make(map[string]interface{})

	if v := d.Get("basic_resource_setting"); v != nil {
		taskmanagerResourceSettingSpec := make(map[string]interface{})
		cpu1, _ := jsonpath.Get("$[0].taskmanager_resource_setting_spec[0].cpu", d.Get("basic_resource_setting"))
		if cpu1 != nil && cpu1 != "" {
			taskmanagerResourceSettingSpec["cpu"] = cpu1
		}
		memory1, _ := jsonpath.Get("$[0].taskmanager_resource_setting_spec[0].memory", d.Get("basic_resource_setting"))
		if memory1 != nil && memory1 != "" {
			taskmanagerResourceSettingSpec["memory"] = memory1
		}

		if len(taskmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec
		}
		jobmanagerResourceSettingSpec := make(map[string]interface{})
		cpu3, _ := jsonpath.Get("$[0].jobmanager_resource_setting_spec[0].cpu", d.Get("basic_resource_setting"))
		if cpu3 != nil && cpu3 != "" {
			jobmanagerResourceSettingSpec["cpu"] = cpu3
		}
		memory3, _ := jsonpath.Get("$[0].jobmanager_resource_setting_spec[0].memory", d.Get("basic_resource_setting"))
		if memory3 != nil && memory3 != "" {
			jobmanagerResourceSettingSpec["memory"] = memory3
		}

		if len(jobmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec
		}
		parallelism1, _ := jsonpath.Get("$[0].parallelism", v)
		if parallelism1 != nil && parallelism1 != "" {
			basicResourceSetting["parallelism"] = parallelism1
		}

		request["basicResourceSetting"] = basicResourceSetting
	}

	if v, ok := d.GetOk("flink_conf"); ok {
		request["flinkConf"] = v
	}
	logging := make(map[string]interface{})

	if v := d.Get("logging"); !IsNil(v) {
		logReservePolicy := make(map[string]interface{})
		expirationDays1, _ := jsonpath.Get("$[0].log_reserve_policy[0].expiration_days", d.Get("logging"))
		if expirationDays1 != nil && expirationDays1 != "" {
			logReservePolicy["expirationDays"] = expirationDays1
		}
		openHistory1, _ := jsonpath.Get("$[0].log_reserve_policy[0].open_history", d.Get("logging"))
		if openHistory1 != nil && openHistory1 != "" {
			logReservePolicy["openHistory"] = openHistory1
		}

		if len(logReservePolicy) > 0 {
			logging["logReservePolicy"] = logReservePolicy
		}
		if v, ok := d.GetOk("logging"); ok {
			localData, err := jsonpath.Get("$[0].log4j_loggers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["loggerLevel"] = dataLoopTmp["logger_level"]
				dataLoopMap["loggerName"] = dataLoopTmp["logger_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			logging["log4jLoggers"] = localMaps
		}

		loggingProfile1, _ := jsonpath.Get("$[0].logging_profile", v)
		if loggingProfile1 != nil && loggingProfile1 != "" {
			logging["loggingProfile"] = loggingProfile1
		}
		log4J2ConfigurationTemplate, _ := jsonpath.Get("$[0].log4j2_configuration_template", v)
		if log4J2ConfigurationTemplate != nil && log4J2ConfigurationTemplate != "" {
			logging["log4j2ConfigurationTemplate"] = log4J2ConfigurationTemplate
		}

		request["logging"] = logging
	} else {
		// The CreateSessionCluster API requires the logging parameter,
		// fall back to the default logging profile when it is not specified.
		request["logging"] = map[string]interface{}{
			"loggingProfile": "default",
		}
	}

	if v, ok := d.GetOk("labels"); ok {
		request["labels"] = v
	}
	request["deploymentTargetName"] = d.Get("deployment_target_name")
	request["engineVersion"] = d.Get("engine_version")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_session_cluster", action, AlibabaCloudSdkGoERROR)
	}

	dataworkspaceVar, _ := jsonpath.Get("$.data.workspace", response)
	datanamespaceVar, _ := jsonpath.Get("$.data.namespace", response)
	datanameVar, _ := jsonpath.Get("$.data.name", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", dataworkspaceVar, datanamespaceVar, datanameVar))

	return resourceAliCloudRealtimeComputeSessionClusterRead(d, meta)
}

func resourceAliCloudRealtimeComputeSessionClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeSessionCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_session_cluster DescribeRealtimeComputeSessionCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("created_at", objectRaw["createdAt"])
	d.Set("creator", objectRaw["creator"])
	d.Set("creator_name", objectRaw["creatorName"])
	d.Set("deployment_target_name", objectRaw["deploymentTargetName"])
	d.Set("engine_version", objectRaw["engineVersion"])
	d.Set("flink_conf", objectRaw["flinkConf"])
	d.Set("labels", objectRaw["labels"])
	d.Set("modified_at", objectRaw["modifiedAt"])
	d.Set("modifier", objectRaw["modifier"])
	d.Set("modifier_name", objectRaw["modifierName"])
	d.Set("session_cluster_id", objectRaw["sessionClusterId"])
	d.Set("namespace", objectRaw["namespace"])
	d.Set("session_cluster_name", objectRaw["name"])
	d.Set("workspace", objectRaw["workspace"])

	basicResourceSettingMaps := make([]map[string]interface{}, 0)
	basicResourceSettingMap := make(map[string]interface{})
	basicResourceSettingRaw := make(map[string]interface{})
	if objectRaw["basicResourceSetting"] != nil {
		basicResourceSettingRaw = objectRaw["basicResourceSetting"].(map[string]interface{})
	}
	if len(basicResourceSettingRaw) > 0 {
		basicResourceSettingMap["parallelism"] = basicResourceSettingRaw["parallelism"]

		jobmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
		jobmanagerResourceSettingSpecMap := make(map[string]interface{})
		jobmanagerResourceSettingSpecRaw := make(map[string]interface{})
		if basicResourceSettingRaw["jobmanagerResourceSettingSpec"] != nil {
			jobmanagerResourceSettingSpecRaw = basicResourceSettingRaw["jobmanagerResourceSettingSpec"].(map[string]interface{})
		}
		if len(jobmanagerResourceSettingSpecRaw) > 0 {
			jobmanagerResourceSettingSpecMap["cpu"] = jobmanagerResourceSettingSpecRaw["cpu"]
			jobmanagerResourceSettingSpecMap["memory"] = jobmanagerResourceSettingSpecRaw["memory"]

			jobmanagerResourceSettingSpecMaps = append(jobmanagerResourceSettingSpecMaps, jobmanagerResourceSettingSpecMap)
		}
		basicResourceSettingMap["jobmanager_resource_setting_spec"] = jobmanagerResourceSettingSpecMaps
		taskmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
		taskmanagerResourceSettingSpecMap := make(map[string]interface{})
		taskmanagerResourceSettingSpecRaw := make(map[string]interface{})
		if basicResourceSettingRaw["taskmanagerResourceSettingSpec"] != nil {
			taskmanagerResourceSettingSpecRaw = basicResourceSettingRaw["taskmanagerResourceSettingSpec"].(map[string]interface{})
		}
		if len(taskmanagerResourceSettingSpecRaw) > 0 {
			taskmanagerResourceSettingSpecMap["cpu"] = taskmanagerResourceSettingSpecRaw["cpu"]
			taskmanagerResourceSettingSpecMap["memory"] = taskmanagerResourceSettingSpecRaw["memory"]

			taskmanagerResourceSettingSpecMaps = append(taskmanagerResourceSettingSpecMaps, taskmanagerResourceSettingSpecMap)
		}
		basicResourceSettingMap["taskmanager_resource_setting_spec"] = taskmanagerResourceSettingSpecMaps
		basicResourceSettingMaps = append(basicResourceSettingMaps, basicResourceSettingMap)
	}
	if err := d.Set("basic_resource_setting", basicResourceSettingMaps); err != nil {
		return err
	}
	loggingMaps := make([]map[string]interface{}, 0)
	loggingMap := make(map[string]interface{})
	loggingRaw := make(map[string]interface{})
	if objectRaw["logging"] != nil {
		loggingRaw = objectRaw["logging"].(map[string]interface{})
	}
	if len(loggingRaw) > 0 {
		loggingMap["log4j2_configuration_template"] = loggingRaw["log4j2ConfigurationTemplate"]
		loggingMap["logging_profile"] = loggingRaw["loggingProfile"]

		log4jLoggersRaw := loggingRaw["log4jLoggers"]
		log4JLoggersMaps := make([]map[string]interface{}, 0)
		if log4jLoggersRaw != nil {
			for _, log4jLoggersChildRaw := range convertToInterfaceArray(log4jLoggersRaw) {
				log4JLoggersMap := make(map[string]interface{})
				log4jLoggersChildRaw := log4jLoggersChildRaw.(map[string]interface{})
				log4JLoggersMap["logger_level"] = log4jLoggersChildRaw["loggerLevel"]
				log4JLoggersMap["logger_name"] = log4jLoggersChildRaw["loggerName"]

				log4JLoggersMaps = append(log4JLoggersMaps, log4JLoggersMap)
			}
		}
		loggingMap["log4j_loggers"] = log4JLoggersMaps
		logReservePolicyMaps := make([]map[string]interface{}, 0)
		logReservePolicyMap := make(map[string]interface{})
		logReservePolicyRaw := make(map[string]interface{})
		if loggingRaw["logReservePolicy"] != nil {
			logReservePolicyRaw = loggingRaw["logReservePolicy"].(map[string]interface{})
		}
		if len(logReservePolicyRaw) > 0 {
			logReservePolicyMap["expiration_days"] = logReservePolicyRaw["expirationDays"]
			logReservePolicyMap["open_history"] = logReservePolicyRaw["openHistory"]

			logReservePolicyMaps = append(logReservePolicyMaps, logReservePolicyMap)
		}
		loggingMap["log_reserve_policy"] = logReservePolicyMaps
		loggingMaps = append(loggingMaps, loggingMap)
	}
	if err := d.Set("logging", loggingMaps); err != nil {
		return err
	}
	statusMaps := make([]map[string]interface{}, 0)
	statusMap := make(map[string]interface{})
	statusRaw := make(map[string]interface{})
	if objectRaw["status"] != nil {
		statusRaw = objectRaw["status"].(map[string]interface{})
	}
	if len(statusRaw) > 0 {
		statusMap["current_session_cluster_status"] = statusRaw["currentSessionClusterStatus"]

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
			runningMap["last_update_time"] = runningRaw["lastUpdateTime"]
			runningMap["started_at"] = runningRaw["startedAt"]

			referenceDeploymentIdsRaw := make([]interface{}, 0)
			if runningRaw["referenceDeploymentIds"] != nil {
				referenceDeploymentIdsRaw = convertToInterfaceArray(runningRaw["referenceDeploymentIds"])
			}

			runningMap["reference_deployment_ids"] = referenceDeploymentIdsRaw
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

func resourceAliCloudRealtimeComputeSessionClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	sessionClusterName := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/sessionclusters/%s", namespace, sessionClusterName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	body = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])
	// Keep the UpdateSessionCluster body consistent with the Create path: the
	// backend validates the SessionCluster body structure, so include workspace
	// and namespace alongside the workspace header and namespace path segment.
	request["workspace"] = parts[0]
	request["namespace"] = namespace

	if d.HasChange("basic_resource_setting") {
		update = true
	}
	basicResourceSetting := make(map[string]interface{})

	if v := d.Get("basic_resource_setting"); v != nil {
		taskmanagerResourceSettingSpec := make(map[string]interface{})
		cpu1, _ := jsonpath.Get("$[0].taskmanager_resource_setting_spec[0].cpu", d.Get("basic_resource_setting"))
		if cpu1 != nil && cpu1 != "" {
			taskmanagerResourceSettingSpec["cpu"] = cpu1
		}
		memory1, _ := jsonpath.Get("$[0].taskmanager_resource_setting_spec[0].memory", d.Get("basic_resource_setting"))
		if memory1 != nil && memory1 != "" {
			taskmanagerResourceSettingSpec["memory"] = memory1
		}

		if len(taskmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec
		}
		jobmanagerResourceSettingSpec := make(map[string]interface{})
		cpu3, _ := jsonpath.Get("$[0].jobmanager_resource_setting_spec[0].cpu", d.Get("basic_resource_setting"))
		if cpu3 != nil && cpu3 != "" {
			jobmanagerResourceSettingSpec["cpu"] = cpu3
		}
		memory3, _ := jsonpath.Get("$[0].jobmanager_resource_setting_spec[0].memory", d.Get("basic_resource_setting"))
		if memory3 != nil && memory3 != "" {
			jobmanagerResourceSettingSpec["memory"] = memory3
		}

		if len(jobmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec
		}
		parallelism1, _ := jsonpath.Get("$[0].parallelism", v)
		if parallelism1 != nil && parallelism1 != "" {
			basicResourceSetting["parallelism"] = parallelism1
		}

		request["basicResourceSetting"] = basicResourceSetting
	}

	if d.HasChange("flink_conf") {
		update = true
	}
	if v, ok := d.GetOk("flink_conf"); ok || d.HasChange("flink_conf") {
		request["flinkConf"] = v
	}
	if d.HasChange("logging") {
		update = true
	}
	logging := make(map[string]interface{})

	if v := d.Get("logging"); !IsNil(v) || d.HasChange("logging") {
		logReservePolicy := make(map[string]interface{})
		expirationDays1, _ := jsonpath.Get("$[0].log_reserve_policy[0].expiration_days", d.Get("logging"))
		if expirationDays1 != nil && expirationDays1 != "" {
			logReservePolicy["expirationDays"] = expirationDays1
		}
		openHistory1, _ := jsonpath.Get("$[0].log_reserve_policy[0].open_history", d.Get("logging"))
		if openHistory1 != nil && openHistory1 != "" {
			logReservePolicy["openHistory"] = openHistory1
		}

		if len(logReservePolicy) > 0 {
			logging["logReservePolicy"] = logReservePolicy
		}
		if v, ok := d.GetOk("logging"); ok {
			localData, err := jsonpath.Get("$[0].log4j_loggers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["loggerLevel"] = dataLoopTmp["logger_level"]
				dataLoopMap["loggerName"] = dataLoopTmp["logger_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			logging["log4jLoggers"] = localMaps
		}

		loggingProfile1, _ := jsonpath.Get("$[0].logging_profile", v)
		if loggingProfile1 != nil && loggingProfile1 != "" {
			logging["loggingProfile"] = loggingProfile1
		}
		log4J2ConfigurationTemplate, _ := jsonpath.Get("$[0].log4j2_configuration_template", v)
		if log4J2ConfigurationTemplate != nil && log4J2ConfigurationTemplate != "" {
			logging["log4j2ConfigurationTemplate"] = log4J2ConfigurationTemplate
		}

		request["logging"] = logging
	}

	if d.HasChange("labels") {
		update = true
	}
	if v, ok := d.GetOk("labels"); ok || d.HasChange("labels") {
		request["labels"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("ververica", "2022-07-18", action, query, header, body, true)
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

	return resourceAliCloudRealtimeComputeSessionClusterRead(d, meta)
}

func resourceAliCloudRealtimeComputeSessionClusterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	sessionClusterName := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/sessionclusters/%s", namespace, sessionClusterName)
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
