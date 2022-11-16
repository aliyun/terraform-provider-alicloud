package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEventBridgeEventStreamings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEventBridgeEventStreamingsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PAUSED", "READY", "RUNNING", "RUNNING_FAILED", "STARTING", "STARTING_FAILED"}, false),
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"streamings": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"event_streaming_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"filter_pattern": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"run_options": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch_window": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count_based_window": {
													Computed: true,
													Type:     schema.TypeInt,
												},
												"time_based_window": {
													Computed: true,
													Type:     schema.TypeInt,
												},
											},
										},
									},
									"dead_letter_queue": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"arn": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"errors_tolerance": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"maximum_tasks": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"retry_strategy": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"maximum_event_age_in_seconds": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"maximum_retry_attempts": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"push_retry_strategy": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sink_fc_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"function_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"invocation_type": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"qualifier": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"service_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
									"sink_kafka_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"acks": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"instance_id": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"key": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"sasl_user": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"value": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
									"sink_mns_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"is_base64_encode": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"queue_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
									"sink_rabbit_mq_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"exchange": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"instance_id": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"message_id": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"properties": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"queue_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"routing_key": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"target_type": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"virtual_host_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
									"sink_rocket_mq_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"instance_id": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"keys": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"properties": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"tags": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
									"sink_sls_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"log_store": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"project": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"role_name": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"form": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"template": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"value": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"source": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_dts_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"broker_url": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"init_check_point": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"password": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"sid": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"task_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"username": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_kafka_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"consumer_group": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"instance_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"network": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"offset_reset": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"security_group_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"vswitch_ids": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"vpc_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_mns_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_base64_decode": {
													Computed: true,
													Type:     schema.TypeBool,
												},
												"queue_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_mqtt_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_rabbit_mq_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"queue_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"virtual_host_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_rocket_mq_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"instance_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"offset": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"tag": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"timestamp": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"topic": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"source_sls_parameters": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"consume_position": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"consumer_group": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"log_store": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"project": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"role_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEventBridgeEventStreamingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")

	var eventStreamingNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		eventStreamingNameRegex = r
	}

	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListEventStreaming"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_event_streamings", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Data.EventStreamings", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.EventStreamings", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EventStreamingName"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			if eventStreamingNameRegex != nil && !eventStreamingNameRegex.MatchString(fmt.Sprint(item["EventStreamingName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                   fmt.Sprint(object["EventStreamingName"]),
			"description":          object["Description"],
			"event_streaming_name": object["EventStreamingName"],
			"filter_pattern":       object["FilterPattern"],
			"status":               object["Status"],
		}

		runOptionsSli := make([]map[string]interface{}, 0)
		if runOptions, ok := object["RunOptions"]; ok {
			if len(runOptions.(map[string]interface{})) > 0 {
				runOptionsMap := make(map[string]interface{})
				batchWindowSli := make([]map[string]interface{}, 0)
				if batchWindow, ok := runOptions.(map[string]interface{})["BatchWindow"]; ok {
					if len(batchWindow.(map[string]interface{})) > 0 {
						batchWindowMap := make(map[string]interface{})
						batchWindowMap["count_based_window"] = batchWindow.(map[string]interface{})["CountBasedWindow"]
						batchWindowMap["time_based_window"] = batchWindow.(map[string]interface{})["TimeBasedWindow"]
						batchWindowSli = append(batchWindowSli, batchWindowMap)
					}
				}
				runOptionsMap["batch_window"] = batchWindowSli

				deadLetterQueueSli := make([]map[string]interface{}, 0)
				if deadLetterQueue, ok := runOptions.(map[string]interface{})["DeadLetterQueue"]; ok {
					if len(deadLetterQueue.(map[string]interface{})) > 0 {
						deadLetterQueueMap := make(map[string]interface{})
						deadLetterQueueMap["arn"] = deadLetterQueue.(map[string]interface{})["Arn"]
						deadLetterQueueSli = append(deadLetterQueueSli, deadLetterQueueMap)
					}
				}

				runOptionsMap["dead_letter_queue"] = deadLetterQueueSli
				runOptionsMap["errors_tolerance"] = runOptions.(map[string]interface{})["ErrorsTolerance"]
				runOptionsMap["maximum_tasks"] = runOptions.(map[string]interface{})["MaximumTasks"]

				retryStrategySli := make([]map[string]interface{}, 0)
				if retryStrategy, ok := runOptions.(map[string]interface{})["RetryStrategy"]; ok {
					if len(retryStrategy.(map[string]interface{})) > 0 {
						retryStrategyMap := make(map[string]interface{})
						retryStrategyMap["maximum_event_age_in_seconds"] = retryStrategy.(map[string]interface{})["MaximumEventAgeInSeconds"]
						retryStrategyMap["maximum_retry_attempts"] = retryStrategy.(map[string]interface{})["MaximumRetryAttempts"]
						retryStrategyMap["push_retry_strategy"] = retryStrategy.(map[string]interface{})["PushRetryStrategy"]
						retryStrategySli = append(retryStrategySli, retryStrategyMap)
					}
				}
				runOptionsMap["retry_strategy"] = retryStrategySli
				runOptionsSli = append(runOptionsSli, runOptionsMap)
			}
		}
		mapping["run_options"] = runOptionsSli

		sinkSli := make([]map[string]interface{}, 0)
		if sink, ok := object["Sink"]; ok {
			if len(sink.(map[string]interface{})) > 0 {
				sinkMap := make(map[string]interface{})
				sinkFcParametersSli := make([]map[string]interface{}, 0)
				if sinkFcParameters, ok := sink.(map[string]interface{})["SinkFcParameters"]; ok {
					if len(sinkFcParameters.(map[string]interface{})) > 0 {
						sinkFcParametersMap := make(map[string]interface{})
						bodySli := make([]map[string]interface{}, 0)
						if body, ok := sinkFcParameters.(map[string]interface{})["Body"]; ok {
							if len(body.(map[string]interface{})) > 0 {
								bodyMap := make(map[string]interface{})
								bodyMap["form"] = body.(map[string]interface{})["Form"]
								bodyMap["template"] = body.(map[string]interface{})["Template"]
								bodyMap["value"] = body.(map[string]interface{})["Value"]
								bodySli = append(bodySli, bodyMap)
							}
						}
						sinkFcParametersMap["body"] = bodySli

						functionNameSli := make([]map[string]interface{}, 0)
						if functionName, ok := sinkFcParameters.(map[string]interface{})["FunctionName"]; ok {
							if len(sinkFcParameters.(map[string]interface{})["FunctionName"].(map[string]interface{})) > 0 {
								functionNameMap := make(map[string]interface{})
								functionNameMap["form"] = functionName.(map[string]interface{})["Form"]
								functionNameMap["template"] = functionName.(map[string]interface{})["Template"]
								functionNameMap["value"] = functionName.(map[string]interface{})["Value"]
								functionNameSli = append(functionNameSli, functionNameMap)
							}
						}
						sinkFcParametersMap["function_name"] = functionNameSli

						invocationTypeSli := make([]map[string]interface{}, 0)
						if invocationType, ok := sinkFcParameters.(map[string]interface{})["InvocationType"]; ok {
							if len(sinkFcParameters.(map[string]interface{})["InvocationType"].(map[string]interface{})) > 0 {
								invocationTypeMap := make(map[string]interface{})
								invocationTypeMap["form"] = invocationType.(map[string]interface{})["Form"]
								invocationTypeMap["template"] = invocationType.(map[string]interface{})["Template"]
								invocationTypeMap["value"] = invocationType.(map[string]interface{})["Value"]
								invocationTypeSli = append(invocationTypeSli, invocationTypeMap)
							}
						}
						sinkFcParametersMap["invocation_type"] = invocationTypeSli

						qualifierSli := make([]map[string]interface{}, 0)
						if qualifier, ok := sinkFcParameters.(map[string]interface{})["Qualifier"]; ok {
							if len(sinkFcParameters.(map[string]interface{})["Qualifier"].(map[string]interface{})) > 0 {
								qualifierMap := make(map[string]interface{})
								qualifierMap["form"] = qualifier.(map[string]interface{})["Form"]
								qualifierMap["template"] = qualifier.(map[string]interface{})["Template"]
								qualifierMap["value"] = qualifier.(map[string]interface{})["Value"]
								qualifierSli = append(qualifierSli, qualifierMap)
							}
						}
						sinkFcParametersMap["qualifier"] = qualifierSli

						serviceNameSli := make([]map[string]interface{}, 0)
						if serviceName, ok := sinkFcParameters.(map[string]interface{})["ServiceName"]; ok {
							if len(sinkFcParameters.(map[string]interface{})["ServiceName"].(map[string]interface{})) > 0 {
								serviceNameMap := make(map[string]interface{})
								serviceNameMap["form"] = serviceName.(map[string]interface{})["Form"]
								serviceNameMap["template"] = serviceName.(map[string]interface{})["Template"]
								serviceNameMap["value"] = serviceName.(map[string]interface{})["Value"]
								serviceNameSli = append(serviceNameSli, serviceNameMap)
							}
						}
						sinkFcParametersMap["service_name"] = serviceNameSli
						sinkFcParametersSli = append(sinkFcParametersSli, sinkFcParametersMap)
					}
				}
				sinkMap["sink_fc_parameters"] = sinkFcParametersSli

				sinkKafkaParametersSli := make([]map[string]interface{}, 0)
				if sinkKafkaParameters, ok := sink.(map[string]interface{})["SinkKafkaParameters"]; ok {
					if len(sink.(map[string]interface{})["SinkKafkaParameters"].(map[string]interface{})) > 0 {
						sinkKafkaParametersMap := make(map[string]interface{})
						acksSli := make([]map[string]interface{}, 0)
						if acks, ok := sinkKafkaParameters.(map[string]interface{})["Acks"]; ok {
							if len(acks.(map[string]interface{})) > 0 {
								acksMap := make(map[string]interface{})
								acksMap["form"] = acks.(map[string]interface{})["Form"]
								acksMap["template"] = acks.(map[string]interface{})["Template"]
								acksMap["value"] = acks.(map[string]interface{})["Value"]
								acksSli = append(acksSli, acksMap)
							}
						}
						sinkKafkaParametersMap["acks"] = acksSli

						instanceIdSli := make([]map[string]interface{}, 0)
						if instanceId, ok := sinkKafkaParameters.(map[string]interface{})["InstanceId"]; ok {
							if len(sinkKafkaParameters.(map[string]interface{})["InstanceId"].(map[string]interface{})) > 0 {

								instanceIdMap := make(map[string]interface{})
								instanceIdMap["form"] = instanceId.(map[string]interface{})["Form"]
								instanceIdMap["template"] = instanceId.(map[string]interface{})["Template"]
								instanceIdMap["value"] = instanceId.(map[string]interface{})["Value"]
								instanceIdSli = append(instanceIdSli, instanceIdMap)
							}
						}

						sinkKafkaParametersMap["instance_id"] = instanceIdSli

						keySli := make([]map[string]interface{}, 0)
						if key, ok := sinkKafkaParameters.(map[string]interface{})["Key"]; ok {
							if len(sinkKafkaParameters.(map[string]interface{})["Key"].(map[string]interface{})) > 0 {

								keyMap := make(map[string]interface{})
								keyMap["form"] = key.(map[string]interface{})["Form"]
								keyMap["template"] = key.(map[string]interface{})["Template"]
								keyMap["value"] = key.(map[string]interface{})["Value"]
								keySli = append(keySli, keyMap)
							}
						}

						sinkKafkaParametersMap["key"] = keySli

						saslUserSli := make([]map[string]interface{}, 0)
						if saslUser, ok := sinkKafkaParameters.(map[string]interface{})["SaslUser"]; ok {
							if len(sinkKafkaParameters.(map[string]interface{})["SaslUser"].(map[string]interface{})) > 0 {

								saslUserMap := make(map[string]interface{})
								saslUserMap["form"] = saslUser.(map[string]interface{})["Form"]
								saslUserMap["template"] = saslUser.(map[string]interface{})["Template"]
								saslUserMap["value"] = saslUser.(map[string]interface{})["Value"]
								saslUserSli = append(saslUserSli, saslUserMap)
							}
						}

						sinkKafkaParametersMap["sasl_user"] = saslUserSli

						topicSli := make([]map[string]interface{}, 0)
						if topic, ok := sinkKafkaParameters.(map[string]interface{})["Topic"]; ok {
							if len(sinkKafkaParameters.(map[string]interface{})["Topic"].(map[string]interface{})) > 0 {

								topicMap := make(map[string]interface{})
								topicMap["form"] = topic.(map[string]interface{})["Form"]
								topicMap["template"] = topic.(map[string]interface{})["Template"]
								topicMap["value"] = topic.(map[string]interface{})["Value"]
								topicSli = append(topicSli, topicMap)
							}
						}

						sinkKafkaParametersMap["topic"] = topicSli

						valueSli := make([]map[string]interface{}, 0)
						if value, ok := sinkKafkaParameters.(map[string]interface{})["Value"]; ok {
							if len(sinkKafkaParameters.(map[string]interface{})["Value"].(map[string]interface{})) > 0 {

								valueMap := make(map[string]interface{})
								valueMap["form"] = value.(map[string]interface{})["Form"]
								valueMap["template"] = value.(map[string]interface{})["Template"]
								valueMap["value"] = value.(map[string]interface{})["Value"]
								valueSli = append(valueSli, valueMap)
							}
						}

						sinkKafkaParametersMap["value"] = valueSli
						sinkKafkaParametersSli = append(sinkKafkaParametersSli, sinkKafkaParametersMap)
					}
				}
				sinkMap["sink_kafka_parameters"] = sinkKafkaParametersSli

				sinkMNSParametersSli := make([]map[string]interface{}, 0)
				if sinkMNSParameters, ok := sink.(map[string]interface{})["SinkMNSParameters"]; ok {
					if len(sink.(map[string]interface{})["SinkMNSParameters"].(map[string]interface{})) > 0 {
						sinkMNSParametersMap := make(map[string]interface{})
						bodySli := make([]map[string]interface{}, 0)
						if body, ok := sinkMNSParameters.(map[string]interface{})["Body"]; ok {
							if len(sinkMNSParameters.(map[string]interface{})["Body"].(map[string]interface{})) > 0 {

								bodyMap := make(map[string]interface{})
								bodyMap["form"] = body.(map[string]interface{})["Form"]
								bodyMap["template"] = body.(map[string]interface{})["Template"]
								bodyMap["value"] = body.(map[string]interface{})["Value"]
								bodySli = append(bodySli, bodyMap)
							}
						}

						sinkMNSParametersMap["body"] = bodySli

						isBase64EncodeSli := make([]map[string]interface{}, 0)
						if isBase64Encode, ok := sinkMNSParameters.(map[string]interface{})["IsBase64Encode"]; ok {
							if len(sinkMNSParameters.(map[string]interface{})["IsBase64Encode"].(map[string]interface{})) > 0 {

								isBase64EncodeMap := make(map[string]interface{})
								isBase64EncodeMap["form"] = isBase64Encode.(map[string]interface{})["Form"]
								isBase64EncodeMap["template"] = isBase64Encode.(map[string]interface{})["Template"]
								isBase64EncodeMap["value"] = isBase64Encode.(map[string]interface{})["Value"]
								isBase64EncodeSli = append(isBase64EncodeSli, isBase64EncodeMap)
							}
						}

						sinkMNSParametersMap["is_base64_encode"] = isBase64EncodeSli

						queueNameSli := make([]map[string]interface{}, 0)
						if queueName, ok := sinkMNSParameters.(map[string]interface{})["QueueName"]; ok {
							if len(sinkMNSParameters.(map[string]interface{})["QueueName"].(map[string]interface{})) > 0 {

								queueNameMap := make(map[string]interface{})
								queueNameMap["form"] = queueName.(map[string]interface{})["Form"]
								queueNameMap["template"] = queueName.(map[string]interface{})["Template"]
								queueNameMap["value"] = queueName.(map[string]interface{})["Value"]
								queueNameSli = append(queueNameSli, queueNameMap)
							}
						}

						sinkMNSParametersMap["queue_name"] = queueNameSli
						sinkMNSParametersSli = append(sinkMNSParametersSli, sinkMNSParametersMap)
					}
				}
				sinkMap["sink_mns_parameters"] = sinkMNSParametersSli

				sinkRabbitMQParametersSli := make([]map[string]interface{}, 0)
				if sinkRabbitMQParameters, ok := sink.(map[string]interface{})["SinkRabbitMQParameters"]; ok {
					if len(sinkRabbitMQParameters.(map[string]interface{})) > 0 {

						sinkRabbitMQParametersMap := make(map[string]interface{})

						bodySli := make([]map[string]interface{}, 0)
						if body, ok := sinkRabbitMQParameters.(map[string]interface{})["Body"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["Body"].(map[string]interface{})) > 0 {

								bodyMap := make(map[string]interface{})
								bodyMap["form"] = body.(map[string]interface{})["Form"]
								bodyMap["template"] = body.(map[string]interface{})["Template"]
								bodyMap["value"] = body.(map[string]interface{})["Value"]
								bodySli = append(bodySli, bodyMap)
							}
						}

						sinkRabbitMQParametersMap["body"] = bodySli

						exchangeSli := make([]map[string]interface{}, 0)
						if exchange, ok := sinkRabbitMQParameters.(map[string]interface{})["Exchange"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["Exchange"].(map[string]interface{})) > 0 {

								exchangeMap := make(map[string]interface{})
								exchangeMap["form"] = exchange.(map[string]interface{})["Form"]
								exchangeMap["template"] = exchange.(map[string]interface{})["Template"]
								exchangeMap["value"] = exchange.(map[string]interface{})["Value"]
								exchangeSli = append(exchangeSli, exchangeMap)
							}
						}

						sinkRabbitMQParametersMap["exchange"] = exchangeSli

						instanceIdSli := make([]map[string]interface{}, 0)
						if instanceId, ok := sinkRabbitMQParameters.(map[string]interface{})["InstanceId"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["InstanceId"].(map[string]interface{})) > 0 {

								instanceIdMap := make(map[string]interface{})
								instanceIdMap["form"] = instanceId.(map[string]interface{})["Form"]
								instanceIdMap["template"] = instanceId.(map[string]interface{})["Template"]
								instanceIdMap["value"] = instanceId.(map[string]interface{})["Value"]
								instanceIdSli = append(instanceIdSli, instanceIdMap)
							}
						}

						sinkRabbitMQParametersMap["instance_id"] = instanceIdSli

						messageIdSli := make([]map[string]interface{}, 0)
						if messageId, ok := sinkRabbitMQParameters.(map[string]interface{})["MessageId"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["MessageId"].(map[string]interface{})) > 0 {

								messageIdMap := make(map[string]interface{})
								messageIdMap["form"] = messageId.(map[string]interface{})["Form"]
								messageIdMap["template"] = messageId.(map[string]interface{})["Template"]
								messageIdMap["value"] = messageId.(map[string]interface{})["Value"]
								messageIdSli = append(messageIdSli, messageIdMap)
							}
						}

						sinkRabbitMQParametersMap["message_id"] = messageIdSli

						propertiesSli := make([]map[string]interface{}, 0)
						if properties, ok := sinkRabbitMQParameters.(map[string]interface{})["Properties"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["Properties"].(map[string]interface{})) > 0 {

								propertiesMap := make(map[string]interface{})
								propertiesMap["form"] = properties.(map[string]interface{})["Form"]
								propertiesMap["template"] = properties.(map[string]interface{})["Template"]
								propertiesMap["value"] = properties.(map[string]interface{})["Value"]
								propertiesSli = append(propertiesSli, propertiesMap)
							}
						}

						sinkRabbitMQParametersMap["properties"] = propertiesSli

						queueNameSli := make([]map[string]interface{}, 0)
						if queueName, ok := sinkRabbitMQParameters.(map[string]interface{})["QueueName"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["QueueName"].(map[string]interface{})) > 0 {

								queueNameMap := make(map[string]interface{})
								queueNameMap["form"] = queueName.(map[string]interface{})["Form"]
								queueNameMap["template"] = queueName.(map[string]interface{})["Template"]
								queueNameMap["value"] = queueName.(map[string]interface{})["Value"]
								queueNameSli = append(queueNameSli, queueNameMap)
							}
						}

						sinkRabbitMQParametersMap["queue_name"] = queueNameSli

						routingKeySli := make([]map[string]interface{}, 0)
						if routingKey, ok := sinkRabbitMQParameters.(map[string]interface{})["RoutingKey"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["RoutingKey"].(map[string]interface{})) > 0 {

								routingKeyMap := make(map[string]interface{})
								routingKeyMap["form"] = routingKey.(map[string]interface{})["Form"]
								routingKeyMap["template"] = routingKey.(map[string]interface{})["Template"]
								routingKeyMap["value"] = routingKey.(map[string]interface{})["Value"]
								routingKeySli = append(routingKeySli, routingKeyMap)
							}
						}

						sinkRabbitMQParametersMap["routing_key"] = routingKeySli

						targetTypeSli := make([]map[string]interface{}, 0)
						if targetType, ok := sinkRabbitMQParameters.(map[string]interface{})["TargetType"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["TargetType"].(map[string]interface{})) > 0 {

								targetTypeMap := make(map[string]interface{})
								targetTypeMap["form"] = targetType.(map[string]interface{})["Form"]
								targetTypeMap["template"] = targetType.(map[string]interface{})["Template"]
								targetTypeMap["value"] = targetType.(map[string]interface{})["Value"]
								targetTypeSli = append(targetTypeSli, targetTypeMap)
							}
						}

						sinkRabbitMQParametersMap["target_type"] = targetTypeSli

						virtualHostNameSli := make([]map[string]interface{}, 0)
						if virtualHostName, ok := sinkRabbitMQParameters.(map[string]interface{})["VirtualHostName"]; ok {
							if len(sinkRabbitMQParameters.(map[string]interface{})["VirtualHostName"].(map[string]interface{})) > 0 {

								virtualHostNameMap := make(map[string]interface{})
								virtualHostNameMap["form"] = virtualHostName.(map[string]interface{})["Form"]
								virtualHostNameMap["template"] = virtualHostName.(map[string]interface{})["Template"]
								virtualHostNameMap["value"] = virtualHostName.(map[string]interface{})["Value"]
								virtualHostNameSli = append(virtualHostNameSli, virtualHostNameMap)
							}
						}

						sinkRabbitMQParametersMap["virtual_host_name"] = virtualHostNameSli
						sinkRabbitMQParametersSli = append(sinkRabbitMQParametersSli, sinkRabbitMQParametersMap)
					}
				}
				sinkMap["sink_rabbit_mq_parameters"] = sinkRabbitMQParametersSli

				sinkRocketMQParametersSli := make([]map[string]interface{}, 0)
				if sinkRocketMQParameters, ok := sink.(map[string]interface{})["SinkRocketMQParameters"]; ok {
					if len(sinkRocketMQParameters.(map[string]interface{})) > 0 {
						sinkRocketMQParametersMap := make(map[string]interface{})
						bodySli := make([]map[string]interface{}, 0)
						if body, ok := sinkRocketMQParameters.(map[string]interface{})["Body"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["Body"].(map[string]interface{})) > 0 {

								bodyMap := make(map[string]interface{})
								bodyMap["form"] = body.(map[string]interface{})["Form"]
								bodyMap["template"] = body.(map[string]interface{})["Template"]
								bodyMap["value"] = body.(map[string]interface{})["Value"]
								bodySli = append(bodySli, bodyMap)
							}
						}

						sinkRocketMQParametersMap["body"] = bodySli

						instanceIdSli := make([]map[string]interface{}, 0)
						if instanceId, ok := sinkRocketMQParameters.(map[string]interface{})["InstanceId"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["InstanceId"].(map[string]interface{})) > 0 {

								instanceIdMap := make(map[string]interface{})
								instanceIdMap["form"] = instanceId.(map[string]interface{})["Form"]
								instanceIdMap["template"] = instanceId.(map[string]interface{})["Template"]
								instanceIdMap["value"] = instanceId.(map[string]interface{})["Value"]
								instanceIdSli = append(instanceIdSli, instanceIdMap)
							}
						}

						sinkRocketMQParametersMap["instance_id"] = instanceIdSli

						keysSli := make([]map[string]interface{}, 0)
						if keys, ok := sinkRocketMQParameters.(map[string]interface{})["Keys"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["Keys"].(map[string]interface{})) > 0 {

								keysMap := make(map[string]interface{})
								keysMap["form"] = keys.(map[string]interface{})["Form"]
								keysMap["template"] = keys.(map[string]interface{})["Template"]
								keysMap["value"] = keys.(map[string]interface{})["Value"]
								keysSli = append(keysSli, keysMap)
							}
						}

						sinkRocketMQParametersMap["keys"] = keysSli

						propertiesSli := make([]map[string]interface{}, 0)
						if properties, ok := sinkRocketMQParameters.(map[string]interface{})["Properties"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["Properties"].(map[string]interface{})) > 0 {

								propertiesMap := make(map[string]interface{})
								propertiesMap["form"] = properties.(map[string]interface{})["Form"]
								propertiesMap["template"] = properties.(map[string]interface{})["Template"]
								propertiesMap["value"] = properties.(map[string]interface{})["Value"]
								propertiesSli = append(propertiesSli, propertiesMap)
							}
						}

						sinkRocketMQParametersMap["properties"] = propertiesSli

						tagsSli := make([]map[string]interface{}, 0)
						if tags, ok := sinkRocketMQParameters.(map[string]interface{})["Tags"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["Tags"].(map[string]interface{})) > 0 {

								tagsMap := make(map[string]interface{})
								tagsMap["form"] = tags.(map[string]interface{})["Form"]
								tagsMap["template"] = tags.(map[string]interface{})["Template"]
								tagsMap["value"] = tags.(map[string]interface{})["Value"]
								tagsSli = append(tagsSli, tagsMap)
							}
						}

						sinkRocketMQParametersMap["tags"] = tagsSli

						topicSli := make([]map[string]interface{}, 0)
						if topic, ok := sinkRocketMQParameters.(map[string]interface{})["Topic"]; ok {
							if len(sinkRocketMQParameters.(map[string]interface{})["Topic"].(map[string]interface{})) > 0 {

								topicMap := make(map[string]interface{})
								topicMap["form"] = topic.(map[string]interface{})["Form"]
								topicMap["template"] = topic.(map[string]interface{})["Template"]
								topicMap["value"] = topic.(map[string]interface{})["Value"]
								topicSli = append(topicSli, topicMap)
							}
						}

						sinkRocketMQParametersMap["topic"] = topicSli
						sinkRocketMQParametersSli = append(sinkRocketMQParametersSli, sinkRocketMQParametersMap)
					}

				}
				sinkMap["sink_rocket_mq_parameters"] = sinkRocketMQParametersSli

				sinkSLSParametersSli := make([]map[string]interface{}, 0)
				if sinkSLSParameters, ok := sink.(map[string]interface{})["SinkSLSParameters"]; ok {
					if len(sink.(map[string]interface{})["SinkSLSParameters"].(map[string]interface{})) > 0 {
						sinkSLSParametersMap := make(map[string]interface{})
						bodySli := make([]map[string]interface{}, 0)
						if body, ok := sinkSLSParameters.(map[string]interface{})["Body"]; ok {
							if len(sinkSLSParameters.(map[string]interface{})["Body"].(map[string]interface{})) > 0 {
								bodyMap := make(map[string]interface{})
								bodyMap["form"] = body.(map[string]interface{})["Form"]
								bodyMap["template"] = body.(map[string]interface{})["Template"]
								bodyMap["value"] = body.(map[string]interface{})["Value"]
								bodySli = append(bodySli, bodyMap)
							}
						}

						sinkSLSParametersMap["body"] = bodySli

						logStoreSli := make([]map[string]interface{}, 0)
						if logStore, ok := sinkSLSParameters.(map[string]interface{})["LogStore"]; ok {
							if len(sinkSLSParameters.(map[string]interface{})["LogStore"].(map[string]interface{})) > 0 {
								logStoreMap := make(map[string]interface{})
								logStoreMap["form"] = logStore.(map[string]interface{})["Form"]
								logStoreMap["template"] = logStore.(map[string]interface{})["Template"]
								logStoreMap["value"] = logStore.(map[string]interface{})["Value"]
								logStoreSli = append(logStoreSli, logStoreMap)
							}
						}
						sinkSLSParametersMap["log_store"] = logStoreSli

						projectSli := make([]map[string]interface{}, 0)
						if project, ok := sinkSLSParameters.(map[string]interface{})["Project"]; ok {
							if len(sinkSLSParameters.(map[string]interface{})["Project"].(map[string]interface{})) > 0 {
								projectMap := make(map[string]interface{})
								projectMap["form"] = project.(map[string]interface{})["Form"]
								projectMap["template"] = project.(map[string]interface{})["Template"]
								projectMap["value"] = project.(map[string]interface{})["Value"]
								projectSli = append(projectSli, projectMap)
							}
						}
						sinkSLSParametersMap["project"] = projectSli

						roleNameSli := make([]map[string]interface{}, 0)
						if roleName, ok := sinkSLSParameters.(map[string]interface{})["RoleName"]; ok {
							if len(roleName.(map[string]interface{})) > 0 {
								roleNameMap := make(map[string]interface{})
								roleNameMap["form"] = roleName.(map[string]interface{})["Form"]
								roleNameMap["template"] = roleName.(map[string]interface{})["Template"]
								roleNameMap["value"] = roleName.(map[string]interface{})["Value"]
								roleNameSli = append(roleNameSli, roleNameMap)
							}
						}
						sinkSLSParametersMap["role_name"] = roleNameSli

						topicSli := make([]map[string]interface{}, 0)
						if topic, ok := sinkSLSParameters.(map[string]interface{})["Topic"]; ok {
							if len(sinkSLSParameters.(map[string]interface{})["Topic"].(map[string]interface{})) > 0 {
								topicMap := make(map[string]interface{})
								topicMap["form"] = topic.(map[string]interface{})["Form"]
								topicMap["template"] = topic.(map[string]interface{})["Template"]
								topicMap["value"] = topic.(map[string]interface{})["Value"]
								topicSli = append(topicSli, topicMap)
							}
						}
						sinkSLSParametersMap["topic"] = topicSli
						sinkSLSParametersSli = append(sinkSLSParametersSli, sinkSLSParametersMap)
					}
				}

				sinkMap["sink_sls_parameters"] = sinkSLSParametersSli
				sinkSli = append(sinkSli, sinkMap)
			}
		}
		mapping["sink"] = sinkSli

		sourceSli := make([]map[string]interface{}, 0)
		if source, ok := object["Source"]; ok {
			if len(source.(map[string]interface{})) > 0 {
				sourceMap := make(map[string]interface{})
				sourceDTSParametersSli := make([]map[string]interface{}, 0)
				if sourceDTSParameters, ok := source.(map[string]interface{})["SourceDTSParameters"]; ok {
					if len(sourceDTSParameters.(map[string]interface{})) > 0 {
						sourceDTSParametersMap := make(map[string]interface{})
						sourceDTSParametersMap["broker_url"] = sourceDTSParameters.(map[string]interface{})["BrokerUrl"]
						sourceDTSParametersMap["init_check_point"] = sourceDTSParameters.(map[string]interface{})["InitCheckPoint"]
						sourceDTSParametersMap["password"] = sourceDTSParameters.(map[string]interface{})["Password"]
						sourceDTSParametersMap["sid"] = sourceDTSParameters.(map[string]interface{})["Sid"]
						sourceDTSParametersMap["task_id"] = sourceDTSParameters.(map[string]interface{})["TaskId"]
						sourceDTSParametersMap["topic"] = sourceDTSParameters.(map[string]interface{})["Topic"]
						sourceDTSParametersMap["username"] = sourceDTSParameters.(map[string]interface{})["Username"]
						sourceDTSParametersSli = append(sourceDTSParametersSli, sourceDTSParametersMap)
					}
				}
				sourceMap["source_dts_parameters"] = sourceDTSParametersSli

				sourceKafkaParametersSli := make([]map[string]interface{}, 0)
				if sourceKafkaParameters, ok := source.(map[string]interface{})["SourceKafkaParameters"]; ok {
					if len(source.(map[string]interface{})["SourceKafkaParameters"].(map[string]interface{})) > 0 {

						sourceKafkaParametersMap := make(map[string]interface{})
						sourceKafkaParametersMap["consumer_group"] = sourceKafkaParameters.(map[string]interface{})["ConsumerGroup"]
						sourceKafkaParametersMap["instance_id"] = sourceKafkaParameters.(map[string]interface{})["InstanceId"]
						sourceKafkaParametersMap["network"] = sourceKafkaParameters.(map[string]interface{})["Network"]
						sourceKafkaParametersMap["offset_reset"] = sourceKafkaParameters.(map[string]interface{})["OffsetReset"]
						sourceKafkaParametersMap["region_id"] = sourceKafkaParameters.(map[string]interface{})["RegionId"]
						sourceKafkaParametersMap["security_group_id"] = sourceKafkaParameters.(map[string]interface{})["SecurityGroupId"]
						sourceKafkaParametersMap["topic"] = sourceKafkaParameters.(map[string]interface{})["Topic"]
						sourceKafkaParametersMap["vswitch_ids"] = sourceKafkaParameters.(map[string]interface{})["VSwitchIds"]
						sourceKafkaParametersMap["vpc_id"] = sourceKafkaParameters.(map[string]interface{})["VpcId"]
						sourceKafkaParametersSli = append(sourceKafkaParametersSli, sourceKafkaParametersMap)
					}
				}

				sourceMap["source_kafka_parameters"] = sourceKafkaParametersSli

				sourceMNSParametersSli := make([]map[string]interface{}, 0)
				if sourceMNSParameters, ok := source.(map[string]interface{})["SourceMNSParameters"]; ok {
					if len(source.(map[string]interface{})["SourceMNSParameters"].(map[string]interface{})) > 0 {
						sourceMNSParametersMap := make(map[string]interface{})
						sourceMNSParametersMap["is_base64_decode"] = sourceMNSParameters.(map[string]interface{})["IsBase64Decode"]
						sourceMNSParametersMap["queue_name"] = sourceMNSParameters.(map[string]interface{})["QueueName"]
						sourceMNSParametersMap["region_id"] = sourceMNSParameters.(map[string]interface{})["RegionId"]
						sourceMNSParametersSli = append(sourceMNSParametersSli, sourceMNSParametersMap)
					}
				}

				sourceMap["source_mns_parameters"] = sourceMNSParametersSli

				sourceMQTTParametersSli := make([]map[string]interface{}, 0)
				if sourceMQTTParameters, ok := source.(map[string]interface{})["SourceMQTTParameters"]; ok {
					if len(sourceMQTTParameters.(map[string]interface{})) > 0 {
						sourceMQTTParametersMap := make(map[string]interface{})
						sourceMQTTParametersMap["instance_id"] = sourceMQTTParameters.(map[string]interface{})["InstanceId"]
						sourceMQTTParametersMap["region_id"] = sourceMQTTParameters.(map[string]interface{})["RegionId"]
						sourceMQTTParametersMap["topic"] = sourceMQTTParameters.(map[string]interface{})["Topic"]
						sourceMQTTParametersSli = append(sourceMQTTParametersSli, sourceMQTTParametersMap)
					}
				}

				sourceMap["source_mqtt_parameters"] = sourceMQTTParametersSli

				sourceRabbitMQParametersSli := make([]map[string]interface{}, 0)
				if sourceRabbitMQParameters, ok := source.(map[string]interface{})["SourceRabbitMQParameters"]; ok {
					if len(sourceRabbitMQParameters.(map[string]interface{})) > 0 {
						sourceRabbitMQParametersMap := make(map[string]interface{})
						sourceRabbitMQParametersMap["instance_id"] = sourceRabbitMQParameters.(map[string]interface{})["InstanceId"]
						sourceRabbitMQParametersMap["queue_name"] = sourceRabbitMQParameters.(map[string]interface{})["QueueName"]
						sourceRabbitMQParametersMap["region_id"] = sourceRabbitMQParameters.(map[string]interface{})["RegionId"]
						sourceRabbitMQParametersMap["virtual_host_name"] = sourceRabbitMQParameters.(map[string]interface{})["VirtualHostName"]
						sourceRabbitMQParametersSli = append(sourceRabbitMQParametersSli, sourceRabbitMQParametersMap)
					}
				}

				sourceMap["source_rabbit_mq_parameters"] = sourceRabbitMQParametersSli

				sourceRocketMQParametersSli := make([]map[string]interface{}, 0)
				if sourceRocketMQParameters, ok := source.(map[string]interface{})["SourceRocketMQParameters"]; ok {
					if len(sourceRocketMQParameters.(map[string]interface{})) > 0 {
						sourceRocketMQParametersMap := make(map[string]interface{})
						sourceRocketMQParametersMap["group_id"] = sourceRocketMQParameters.(map[string]interface{})["GroupID"]
						sourceRocketMQParametersMap["instance_id"] = sourceRocketMQParameters.(map[string]interface{})["InstanceId"]
						sourceRocketMQParametersMap["offset"] = sourceRocketMQParameters.(map[string]interface{})["Offset"]
						sourceRocketMQParametersMap["region_id"] = sourceRocketMQParameters.(map[string]interface{})["RegionId"]
						sourceRocketMQParametersMap["tag"] = sourceRocketMQParameters.(map[string]interface{})["Tag"]
						sourceRocketMQParametersMap["timestamp"] = sourceRocketMQParameters.(map[string]interface{})["Timestamp"]
						sourceRocketMQParametersMap["topic"] = sourceRocketMQParameters.(map[string]interface{})["Topic"]
						sourceRocketMQParametersSli = append(sourceRocketMQParametersSli, sourceRocketMQParametersMap)
					}
				}
				sourceMap["source_rocket_mq_parameters"] = sourceRocketMQParametersSli

				sourceSLSParametersSli := make([]map[string]interface{}, 0)
				if sourceSLSParameters, ok := source.(map[string]interface{})["SourceSLSParameters"]; ok {
					if len(sourceSLSParameters.(map[string]interface{})) > 0 {
						sourceSLSParametersMap := make(map[string]interface{})
						sourceSLSParametersMap["consume_position"] = sourceSLSParameters.(map[string]interface{})["ConsumePosition"]
						sourceSLSParametersMap["consumer_group"] = sourceSLSParameters.(map[string]interface{})["ConsumerGroup"]
						sourceSLSParametersMap["log_store"] = sourceSLSParameters.(map[string]interface{})["LogStore"]
						sourceSLSParametersMap["project"] = sourceSLSParameters.(map[string]interface{})["Project"]
						sourceSLSParametersMap["role_name"] = sourceSLSParameters.(map[string]interface{})["RoleName"]
						sourceSLSParametersSli = append(sourceSLSParametersSli, sourceSLSParametersMap)
					}
				}
				sourceMap["source_sls_parameters"] = sourceSLSParametersSli
				sourceSli = append(sourceSli, sourceMap)
			}
		}
		mapping["source"] = sourceSli

		ids = append(ids, fmt.Sprint(object["EventStreamingName"]))
		names = append(names, object["EventStreamingName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("streamings", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
