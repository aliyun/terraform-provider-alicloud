// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudRealtimeComputeSessionClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRealtimeComputeSessionClusterRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_resource_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"taskmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"parallelism": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"jobmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Computed: true,
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
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flink_conf": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"logging": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log4j2_configuration_template": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"logging_profile": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log4j_loggers": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"logger_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"logger_level": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"log_reserve_policy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"open_history": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"expiration_days": {
													Type:     schema.TypeInt,
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
							Computed: true,
						},
						"session_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
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
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudRealtimeComputeSessionClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	namespace := d.Get("namespace")
	// ListSessionClusters
	action := fmt.Sprintf("/api/v2/namespaces/%s/sessionclusters", namespace)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = StringPointer(d.Get("workspace").(string))
	request["namespace"] = namespace
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.data[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["workspace"], ":", item["namespace"], ":", item["name"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["workspace"], ":", objectRaw["namespace"], ":", objectRaw["name"])

		mapping["created_at"] = objectRaw["createdAt"]
		mapping["creator"] = objectRaw["creator"]
		mapping["creator_name"] = objectRaw["creatorName"]
		mapping["deployment_target_name"] = objectRaw["deploymentTargetName"]
		mapping["engine_version"] = objectRaw["engineVersion"]
		mapping["flink_conf"] = objectRaw["flinkConf"]
		mapping["labels"] = objectRaw["labels"]
		mapping["modified_at"] = objectRaw["modifiedAt"]
		mapping["modifier"] = objectRaw["modifier"]
		mapping["modifier_name"] = objectRaw["modifierName"]
		mapping["session_cluster_id"] = objectRaw["sessionClusterId"]
		mapping["namespace"] = objectRaw["namespace"]
		mapping["session_cluster_name"] = objectRaw["name"]
		mapping["workspace"] = objectRaw["workspace"]

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
		mapping["basic_resource_setting"] = basicResourceSettingMaps
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
		mapping["logging"] = loggingMaps
		statusMaps := make([]map[string]interface{}, 0)
		statusMap := make(map[string]interface{})
		statusRaw := make(map[string]interface{})
		if objectRaw["status"] != nil {
			statusRaw = objectRaw["status"].(map[string]interface{})
		}
		if len(statusRaw) > 0 {
			statusMap["current_session_cluster_status"] = statusRaw["currentSessionClusterStatus"]

			statusMaps = append(statusMaps, statusMap)
		}
		mapping["status"] = statusMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["workspace"], ":", objectRaw["namespace"], ":", objectRaw["name"])
		mapping, err = dataSourceAliCloudRealtimeComputeSessionClusterReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudRealtimeComputeSessionClusterReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	getResp, err := realtimeComputeServiceV2.DescribeRealtimeComputeSessionCluster(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["created_at"] = objectRaw["createdAt"]
	mapping["creator"] = objectRaw["creator"]
	mapping["creator_name"] = objectRaw["creatorName"]
	mapping["deployment_target_name"] = objectRaw["deploymentTargetName"]
	mapping["engine_version"] = objectRaw["engineVersion"]
	mapping["flink_conf"] = objectRaw["flinkConf"]
	mapping["labels"] = objectRaw["labels"]
	mapping["modified_at"] = objectRaw["modifiedAt"]
	mapping["modifier"] = objectRaw["modifier"]
	mapping["modifier_name"] = objectRaw["modifierName"]
	mapping["session_cluster_id"] = objectRaw["sessionClusterId"]
	mapping["namespace"] = objectRaw["namespace"]
	mapping["session_cluster_name"] = objectRaw["name"]
	mapping["workspace"] = objectRaw["workspace"]

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
	mapping["basic_resource_setting"] = basicResourceSettingMaps
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
	mapping["logging"] = loggingMaps
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
	mapping["status"] = statusMaps

	return mapping, nil
}
