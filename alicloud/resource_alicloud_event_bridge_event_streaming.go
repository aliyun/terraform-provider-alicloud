package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEventBridgeEventStreaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeEventStreamingCreate,
		Read:   resourceAlicloudEventBridgeEventStreamingRead,
		Update: resourceAlicloudEventBridgeEventStreamingUpdate,
		Delete: resourceAlicloudEventBridgeEventStreamingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Required: true,
				Type:     schema.TypeString,
			},
			"event_streaming_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"filter_pattern": {
				Required: true,
				Type:     schema.TypeString,
			},
			"run_options": {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_window": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count_based_window": {
										Optional: true,
										Type:     schema.TypeInt,
									},
									"time_based_window": {
										Optional: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
						"dead_letter_queue": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"arn": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"errors_tolerance": {
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"ALL", "NONE"}, false),
						},
						"maximum_tasks": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeInt,
						},
						"retry_strategy": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum_event_age_in_seconds": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeInt,
									},
									"maximum_retry_attempts": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeInt,
									},
									"push_retry_strategy": {
										Optional:     true,
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"BACKOFF_RETRY", "EXPONENTIAL_DECAY_RETRY"}, false),
									},
								},
							},
						},
					},
				},
			},
			"sink": {
				Required: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sink_fc_parameters": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"function_name": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"invocation_type": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"qualifier": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"service_name": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink_kafka_parameters": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"acks": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"instance_id": {
										Optional: true,
										MaxItems: 1,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"key": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"sasl_user": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"topic": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"value": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink_mns_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"is_base64_encode": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"queue_name": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink_rabbit_mq_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"exchange": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"instance_id": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"message_id": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"properties": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"queue_name": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"routing_key": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"target_type": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"Exchange", "Queue"}, false),
												},
											},
										},
									},
									"virtual_host_name": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink_rocket_mq_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"instance_id": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"keys": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"properties": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"tags": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"topic": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"sink_sls_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"log_store": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"project": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"role_name": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"topic": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"form": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"template": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"value": {
													Optional: true,
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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_dts_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"broker_url": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"init_check_point": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"password": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"sid": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"task_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"topic": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"username": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_kafka_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"consumer_group": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"instance_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"network": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"PublicNetwork", "Default"}, false),
									},
									"offset_reset": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"security_group_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"topic": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"vswitch_ids": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"vpc_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_mns_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_base64_decode": {
										Optional: true,
										Type:     schema.TypeBool,
									},
									"queue_name": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_mqtt_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"topic": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_rabbit_mq_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"queue_name": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"virtual_host_name": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_rocket_mq_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"offset": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"CONSUME_FROM_LAST_OFFSET", "CONSUME_FROM_FIRST_OFFSET", "CONSUME_FROM_TIMESTAMP"}, false),
									},
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"tag": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"timestamp": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"topic": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"source_sls_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"consume_position": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"consumer_group": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"log_store": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"project": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"role_name": {
										Optional: true,
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
	}
}

func resourceAlicloudEventBridgeEventStreamingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEventStreaming"
	request := make(map[string]interface{})
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request["Description"] = d.Get("description")
	request["EventStreamingName"] = d.Get("event_streaming_name")
	request["FilterPattern"] = d.Get("filter_pattern")

	if v, ok := d.GetOk("run_options"); ok {
		runOptionsMap := make(map[string]interface{})
		for _, runOptions := range v.([]interface{}) {
			runOptionsArg := runOptions.(map[string]interface{})
			if runOptionsArg["batch_window"] != nil {
				batchWindowMap := make(map[string]interface{})
				for _, critical := range runOptionsArg["batch_window"].([]interface{}) {
					batchWindowArg := critical.(map[string]interface{})
					batchWindowMap["CountBasedWindow"] = batchWindowArg["count_based_window"].(int)
					batchWindowMap["TimeBasedWindow"] = batchWindowArg["time_based_window"].(int)
				}
				runOptionsMap["BatchWindow"] = batchWindowMap
			}
			if runOptionsArg["dead_letter_queue"] != nil {
				deadLetterMap := make(map[string]interface{})
				for _, deadLetter := range runOptionsArg["dead_letter_queue"].([]interface{}) {
					deadLetterArg := deadLetter.(map[string]interface{})
					deadLetterMap["Arn"] = deadLetterArg["arn"].(string)
				}
				runOptionsMap["DeadLetterQueue"] = deadLetterMap
			}
			if runOptionsArg["retry_strategy"] != nil {
				retryStrategyMap := make(map[string]interface{})
				for _, retryStrategy := range runOptionsArg["retry_strategy"].([]interface{}) {
					retryStrategyArg := retryStrategy.(map[string]interface{})
					retryStrategyMap["MaximumEventAgeInSeconds"] = retryStrategyArg["maximum_event_age_in_seconds"].(int)
					retryStrategyMap["MaximumRetryAttempts"] = retryStrategyArg["maximum_retry_attempts"].(int)
					retryStrategyMap["PushRetryStrategy"] = retryStrategyArg["push_retry_strategy"].(string)
				}
				runOptionsMap["RetryStrategy"] = retryStrategyMap
			}
			if runOptionsArg["maximum_tasks"] != nil {
				runOptionsMap["MaximumTasks"] = runOptionsArg["maximum_tasks"]
			}
			if runOptionsArg["errors_tolerance"] != nil {
				runOptionsMap["ErrorsTolerance"] = runOptionsArg["errors_tolerance"]
			}
		}
		request["RunOptions"], _ = convertArrayObjectToJsonString(runOptionsMap)
	}
	if v, ok := d.GetOk("source"); ok {
		sourceMap := make(map[string]interface{})
		for _, source := range v.([]interface{}) {
			sourceArg := source.(map[string]interface{})
			if sourceArg["source_dts_parameters"] != nil {
				sourceDtsParametersMap := make(map[string]interface{})
				for _, sourceDtsParameters := range sourceArg["source_dts_parameters"].([]interface{}) {
					sourceDtsParametersArg := sourceDtsParameters.(map[string]interface{})
					sourceDtsParametersMap["BrokerUrl"] = sourceDtsParametersArg["broker_url"].(string)
					sourceDtsParametersMap["InitCheckPoint"] = sourceDtsParametersArg["init_check_point"].(string)
					sourceDtsParametersMap["Password"] = sourceDtsParametersArg["password"].(string)
					sourceDtsParametersMap["Sid"] = sourceDtsParametersArg["sid"].(string)
					sourceDtsParametersMap["TaskId"] = sourceDtsParametersArg["task_id"].(string)
					sourceDtsParametersMap["Topic"] = sourceDtsParametersArg["topic"].(string)
					sourceDtsParametersMap["Username"] = sourceDtsParametersArg["username"].(string)
				}
				sourceMap["SourceDtsParameters"] = sourceDtsParametersMap
			}
			if sourceArg["source_kafka_parameters"] != nil {
				sourceKafkaParametersMap := make(map[string]interface{})
				for _, sourceKafkaParameters := range sourceArg["source_kafka_parameters"].([]interface{}) {
					sourceKafkaParametersArg := sourceKafkaParameters.(map[string]interface{})
					sourceKafkaParametersMap["ConsumerGroup"] = sourceKafkaParametersArg["consumer_group"].(string)
					sourceKafkaParametersMap["InstanceId"] = sourceKafkaParametersArg["instance_id"].(string)
					sourceKafkaParametersMap["Network"] = sourceKafkaParametersArg["network"].(string)
					sourceKafkaParametersMap["OffsetReset"] = sourceKafkaParametersArg["offset_reset"].(string)
					sourceKafkaParametersMap["RegionId"] = sourceKafkaParametersArg["region_id"].(string)
					sourceKafkaParametersMap["SecurityGroupId"] = sourceKafkaParametersArg["security_group_id"].(string)
					sourceKafkaParametersMap["VSwitchIds"] = sourceKafkaParametersArg["vswitch_ids"].(string)
					sourceKafkaParametersMap["Topic"] = sourceKafkaParametersArg["topic"].(string)
					sourceKafkaParametersMap["VpcId"] = sourceKafkaParametersArg["vpc_id"].(string)
				}
				sourceMap["SourceKafkaParameters"] = sourceKafkaParametersMap
			}
			if sourceArg["source_mns_parameters"] != nil {
				sourceMnsParametersMap := make(map[string]interface{})
				for _, sourceMnsParameters := range sourceArg["source_mns_parameters"].([]interface{}) {
					sourceMnsParametersArg := sourceMnsParameters.(map[string]interface{})
					sourceMnsParametersMap["IsBase64Decode"] = sourceMnsParametersArg["is_base64_decode"].(bool)
					sourceMnsParametersMap["QueueName"] = sourceMnsParametersArg["queue_name"].(string)
					sourceMnsParametersMap["RegionId"] = sourceMnsParametersArg["region_id"].(string)
				}
				sourceMap["SourceMNSParameters"] = sourceMnsParametersMap
			}
			if sourceArg["source_mqtt_parameters"] != nil {
				sourceMqttParametersMap := make(map[string]interface{})
				for _, sourceMqttParameters := range sourceArg["source_mqtt_parameters"].([]interface{}) {
					sourceMqttParametersArg := sourceMqttParameters.(map[string]interface{})
					sourceMqttParametersMap["InstanceId"] = sourceMqttParametersArg["instance_id"].(string)
					sourceMqttParametersMap["RegionId"] = sourceMqttParametersArg["region_id"].(string)
					sourceMqttParametersMap["Topic"] = sourceMqttParametersArg["topic"].(string)
				}
				sourceMap["SourceMQTTParameters"] = sourceMqttParametersMap
			}
			if sourceArg["source_rabbit_mq_parameters"] != nil {
				sourceRabbitMqParametersMap := make(map[string]interface{})
				for _, sourceRabbitMqParameters := range sourceArg["source_rabbit_mq_parameters"].([]interface{}) {
					sourceRabbitMqParametersArg := sourceRabbitMqParameters.(map[string]interface{})
					sourceRabbitMqParametersMap["InstanceId"] = sourceRabbitMqParametersArg["instance_id"].(string)
					sourceRabbitMqParametersMap["QueueName"] = sourceRabbitMqParametersArg["queue_name"].(string)
					sourceRabbitMqParametersMap["RegionId"] = sourceRabbitMqParametersArg["region_id"].(string)
					sourceRabbitMqParametersMap["VirtualHostName"] = sourceRabbitMqParametersArg["virtual_host_name"].(string)
				}
				sourceMap["SourceRabbitMQParameters"] = sourceRabbitMqParametersMap
			}
			if sourceArg["source_rocket_mq_parameters"] != nil {
				sourceRocketMqParametersMap := make(map[string]interface{})
				for _, sourceRocketMqParameters := range sourceArg["source_rocket_mq_parameters"].([]interface{}) {
					sourceRocketMqParametersArg := sourceRocketMqParameters.(map[string]interface{})
					sourceRocketMqParametersMap["GroupId"] = sourceRocketMqParametersArg["group_id"].(string)
					sourceRocketMqParametersMap["InstanceId"] = sourceRocketMqParametersArg["instance_id"].(string)
					sourceRocketMqParametersMap["Offset"] = sourceRocketMqParametersArg["offset"].(string)
					sourceRocketMqParametersMap["RegionId"] = sourceRocketMqParametersArg["region_id"].(string)
					sourceRocketMqParametersMap["Tag"] = sourceRocketMqParametersArg["tag"].(string)
					sourceRocketMqParametersMap["Timestamp"] = sourceRocketMqParametersArg["timestamp"].(int)
					sourceRocketMqParametersMap["Topic"] = sourceRocketMqParametersArg["topic"].(string)
				}
				sourceMap["SourceRocketMQParameters"] = sourceRocketMqParametersMap
			}
			if sourceArg["source_sls_parameters"] != nil {
				sourceSlsParametersMap := make(map[string]interface{})
				for _, sourceSlsParameters := range sourceArg["source_sls_parameters"].([]interface{}) {
					sourceSlsParametersArg := sourceSlsParameters.(map[string]interface{})
					sourceSlsParametersMap["ConsumePosition"] = sourceSlsParametersArg["consume_position"].(string)
					sourceSlsParametersMap["LogStore"] = sourceSlsParametersArg["log_store"].(string)
					sourceSlsParametersMap["Project"] = sourceSlsParametersArg["project"].(string)
					sourceSlsParametersMap["RoleName"] = sourceSlsParametersArg["role_name"].(string)
				}
				sourceMap["SourceSLSParameters"] = sourceSlsParametersMap
			}
		}
		request["Source"], _ = convertArrayObjectToJsonString(sourceMap)
	}
	if v, ok := d.GetOk("sink"); ok {
		sinkMap := make(map[string]interface{})
		for _, sink := range v.([]interface{}) {
			sinkArg := sink.(map[string]interface{})
			if sinkArg["sink_fc_parameters"] != nil {
				sinkFcParametersMap := make(map[string]interface{})
				for _, critical := range sinkArg["sink_fc_parameters"].([]interface{}) {
					sinkFcParametersArg := critical.(map[string]interface{})
					if sinkFcParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkFcParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
						}
						sinkFcParametersMap["Body"] = bodyMap
					}
					if sinkFcParametersArg["function_name"] != nil {
						functionNameMap := make(map[string]interface{})
						for _, functionName := range sinkFcParametersArg["function_name"].([]interface{}) {
							functionNameArg := functionName.(map[string]interface{})
							functionNameMap["Form"] = functionNameArg["form"].(string)
							functionNameMap["Value"] = functionNameArg["value"].(string)
							functionNameMap["Template"] = functionNameArg["template"].(string)
						}
						sinkFcParametersMap["FunctionName"] = functionNameMap
					}
					if sinkFcParametersArg["qualifier"] != nil {
						qualifierMap := make(map[string]interface{})
						for _, qualifier := range sinkFcParametersArg["qualifier"].([]interface{}) {
							qualifierArg := qualifier.(map[string]interface{})
							qualifierMap["Form"] = qualifierArg["form"].(string)
							qualifierMap["Value"] = qualifierArg["value"].(string)
							qualifierMap["Template"] = qualifierArg["template"].(string)
						}
						sinkFcParametersMap["Qualifier"] = qualifierMap
					}
					if sinkFcParametersArg["invocation_type"] != nil {
						invocationTypeMap := make(map[string]interface{})
						for _, invocationType := range sinkFcParametersArg["invocation_type"].([]interface{}) {
							invocationTypeArg := invocationType.(map[string]interface{})
							invocationTypeMap["Form"] = invocationTypeArg["form"].(string)
							invocationTypeMap["Value"] = invocationTypeArg["value"].(string)
							invocationTypeMap["Template"] = invocationTypeArg["template"].(string)
						}
						sinkFcParametersMap["InvocationType"] = invocationTypeMap
					}
					if sinkFcParametersArg["service_name"] != nil {
						serviceNameMap := make(map[string]interface{})
						for _, serviceName := range sinkFcParametersArg["service_name"].([]interface{}) {
							serviceNameArg := serviceName.(map[string]interface{})
							serviceNameMap["Form"] = serviceNameArg["form"].(string)
							serviceNameMap["Value"] = serviceNameArg["value"].(string)
							serviceNameMap["Template"] = serviceNameArg["template"].(string)
						}
						sinkFcParametersMap["ServiceName"] = serviceNameMap
					}
				}
				sinkMap["SinkFcParameters"] = sinkFcParametersMap
			}

			if sinkArg["sink_kafka_parameters"] != nil {
				sinkKafkaParametersMap := make(map[string]interface{})
				for _, sinkKafkaParameters := range sinkArg["sink_kafka_parameters"].([]interface{}) {
					sinkKafkaParametersArg := sinkKafkaParameters.(map[string]interface{})
					if sinkKafkaParametersArg["acks"] != nil {
						acksMap := make(map[string]interface{})
						for _, acks := range sinkKafkaParametersArg["acks"].([]interface{}) {
							acksArg := acks.(map[string]interface{})
							acksMap["Form"] = acksArg["form"].(string)
							acksMap["Value"] = acksArg["value"].(string)
							acksMap["Template"] = acksArg["template"].(string)
						}
						sinkKafkaParametersMap["Acks"] = acksMap
					}
					if sinkKafkaParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkKafkaParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkKafkaParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkKafkaParametersArg["key"] != nil {
						keyMap := make(map[string]interface{})
						for _, key := range sinkKafkaParametersArg["key"].([]interface{}) {
							keyArg := key.(map[string]interface{})
							keyMap["Form"] = keyArg["form"].(string)
							keyMap["Value"] = keyArg["value"].(string)
							keyMap["Template"] = keyArg["template"].(string)
						}
						sinkKafkaParametersMap["Key"] = keyMap
					}
					if sinkKafkaParametersArg["sasl_user"] != nil {
						saslUserMap := make(map[string]interface{})
						for _, saslUser := range sinkKafkaParametersArg["sasl_user"].([]interface{}) {
							saslUserArg := saslUser.(map[string]interface{})
							saslUserMap["Form"] = saslUserArg["form"].(string)
							saslUserMap["Value"] = saslUserArg["value"].(string)
							saslUserMap["Template"] = saslUserArg["template"].(string)
						}
						sinkKafkaParametersMap["SaslUser"] = saslUserMap
					}
					if sinkKafkaParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkKafkaParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkKafkaParametersMap["Topic"] = topicMap
					}
					if sinkKafkaParametersArg["value"] != nil {
						valueMap := make(map[string]interface{})
						for _, value := range sinkKafkaParametersArg["value"].([]interface{}) {
							valueArg := value.(map[string]interface{})
							valueMap["Form"] = valueArg["form"].(string)
							valueMap["Value"] = valueArg["value"].(string)
							valueMap["Template"] = valueArg["template"].(string)
						}
						sinkKafkaParametersMap["Value"] = valueMap
					}
				}
				sinkMap["SinkKafkaParameters"] = sinkKafkaParametersMap
			}

			if sinkArg["sink_mns_parameters"] != nil {
				sinkMnsParametersMap := make(map[string]interface{})
				for _, sinkMnsParameters := range sinkArg["sink_mns_parameters"].([]interface{}) {
					sinkMnsParametersArg := sinkMnsParameters.(map[string]interface{})
					if sinkMnsParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkMnsParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkMnsParametersMap["Body"] = bodyMap
					}
					if sinkMnsParametersArg["is_base64_encode"] != nil {
						isBase64EncodeMap := make(map[string]interface{})
						for _, isBase64Encode := range sinkMnsParametersArg["is_base64_encode"].([]interface{}) {
							isBase64EncodeArg := isBase64Encode.(map[string]interface{})
							isBase64EncodeMap["Form"] = isBase64EncodeArg["form"].(string)
							isBase64EncodeMap["Value"] = isBase64EncodeArg["value"].(string)
							isBase64EncodeMap["Template"] = isBase64EncodeArg["template"].(string)
						}
						sinkMnsParametersMap["IsBase64Encode"] = isBase64EncodeMap
					}
					if sinkMnsParametersArg["queue_name"] != nil {
						queueNameMap := make(map[string]interface{})
						for _, queueName := range sinkMnsParametersArg["queue_name"].([]interface{}) {
							queueNameArg := queueName.(map[string]interface{})
							queueNameMap["Form"] = queueNameArg["form"].(string)
							queueNameMap["Value"] = queueNameArg["value"].(string)
							queueNameMap["Template"] = queueNameArg["template"].(string)
						}
						sinkMnsParametersMap["QueueName"] = queueNameMap
					}
				}
				sinkMap["SinkMNSParameters"] = sinkMnsParametersMap
			}

			if sinkArg["sink_rabbit_mq_parameters"] != nil {
				sinkRabbitMqParametersMap := make(map[string]interface{})
				for _, sinkRabbitMqParameters := range sinkArg["sink_rabbit_mq_parameters"].([]interface{}) {
					sinkRabbitMqParametersArg := sinkRabbitMqParameters.(map[string]interface{})
					if sinkRabbitMqParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkRabbitMqParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Body"] = bodyMap
					}
					if sinkRabbitMqParametersArg["exchange"] != nil {
						exchangeMap := make(map[string]interface{})
						for _, exchange := range sinkRabbitMqParametersArg["exchange"].([]interface{}) {
							exchangeArg := exchange.(map[string]interface{})
							exchangeMap["Form"] = exchangeArg["form"].(string)
							exchangeMap["Value"] = exchangeArg["value"].(string)
							exchangeMap["Template"] = exchangeArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Exchange"] = exchangeMap
					}
					if sinkRabbitMqParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkRabbitMqParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkRabbitMqParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkRabbitMqParametersArg["message_id"] != nil {
						messageIdMap := make(map[string]interface{})
						for _, messageId := range sinkRabbitMqParametersArg["message_id"].([]interface{}) {
							messageIdArg := messageId.(map[string]interface{})
							messageIdMap["Form"] = messageIdArg["form"].(string)
							messageIdMap["Value"] = messageIdArg["value"].(string)
							messageIdMap["Template"] = messageIdArg["template"].(string)
						}
						sinkRabbitMqParametersMap["MessageId"] = messageIdMap
					}
					if sinkRabbitMqParametersArg["properties"] != nil {
						propertiesMap := make(map[string]interface{})
						for _, properties := range sinkRabbitMqParametersArg["properties"].([]interface{}) {
							propertiesArg := properties.(map[string]interface{})
							propertiesMap["Form"] = propertiesArg["form"].(string)
							propertiesMap["Value"] = propertiesArg["value"].(string)
							propertiesMap["Template"] = propertiesArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Properties"] = propertiesMap
					}
					if sinkRabbitMqParametersArg["queue_name"] != nil {
						queueNameMap := make(map[string]interface{})
						for _, queueName := range sinkRabbitMqParametersArg["queue_name"].([]interface{}) {
							queueNameArg := queueName.(map[string]interface{})
							queueNameMap["Form"] = queueNameArg["form"].(string)
							queueNameMap["Value"] = queueNameArg["value"].(string)
							queueNameMap["Template"] = queueNameArg["template"].(string)
						}
						sinkRabbitMqParametersMap["QueueName"] = queueNameMap
					}
					if sinkRabbitMqParametersArg["routing_key"] != nil {
						routingKeyMap := make(map[string]interface{})
						for _, routingKey := range sinkRabbitMqParametersArg["routing_key"].([]interface{}) {
							routingKeyArg := routingKey.(map[string]interface{})
							routingKeyMap["Form"] = routingKeyArg["form"].(string)
							routingKeyMap["Value"] = routingKeyArg["value"].(string)
							routingKeyMap["Template"] = routingKeyArg["template"].(string)
						}
						sinkRabbitMqParametersMap["RoutingKey"] = routingKeyMap
					}
					if sinkRabbitMqParametersArg["target_type"] != nil {
						targetTypeMap := make(map[string]interface{})
						for _, targetType := range sinkRabbitMqParametersArg["target_type"].([]interface{}) {
							targetTypeArg := targetType.(map[string]interface{})
							targetTypeMap["Form"] = targetTypeArg["form"].(string)
							targetTypeMap["Value"] = targetTypeArg["value"].(string)
							targetTypeMap["Template"] = targetTypeArg["template"].(string)
						}
						sinkRabbitMqParametersMap["TargetType"] = targetTypeMap
					}
					if sinkRabbitMqParametersArg["virtual_host_name"] != nil {
						virtualHostNameMap := make(map[string]interface{})
						for _, virtualHostName := range sinkRabbitMqParametersArg["virtual_host_name"].([]interface{}) {
							virtualHostNameArg := virtualHostName.(map[string]interface{})
							virtualHostNameMap["Form"] = virtualHostNameArg["form"].(string)
							virtualHostNameMap["Value"] = virtualHostNameArg["value"].(string)
							virtualHostNameMap["Template"] = virtualHostNameArg["template"].(string)
						}
						sinkRabbitMqParametersMap["VirtualHostName"] = virtualHostNameMap
					}
				}
				sinkMap["SinkRabbitMqParameters"] = sinkRabbitMqParametersMap
			}

			if sinkArg["sink_rocket_mq_parameters"] != nil {
				sinkRocketMqParametersMap := make(map[string]interface{})
				for _, sinkRocketMqParameters := range sinkArg["sink_rocket_mq_parameters"].([]interface{}) {
					sinkRocketMqParametersArg := sinkRocketMqParameters.(map[string]interface{})
					if sinkRocketMqParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkRocketMqParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkRocketMqParametersMap["Body"] = bodyMap
					}
					if sinkRocketMqParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkRocketMqParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkRocketMqParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkRocketMqParametersArg["keys"] != nil {
						keysMap := make(map[string]interface{})
						for _, keys := range sinkRocketMqParametersArg["keys"].([]interface{}) {
							keysArg := keys.(map[string]interface{})
							keysMap["Form"] = keysArg["form"].(string)
							keysMap["Value"] = keysArg["value"].(string)
							keysMap["Template"] = keysArg["template"].(string)
						}
						sinkRocketMqParametersMap["Keys"] = keysMap
					}
					if sinkRocketMqParametersArg["properties"] != nil {
						propertiesMap := make(map[string]interface{})
						for _, properties := range sinkRocketMqParametersArg["properties"].([]interface{}) {
							propertiesArg := properties.(map[string]interface{})
							propertiesMap["Form"] = propertiesArg["form"].(string)
							propertiesMap["Value"] = propertiesArg["value"].(string)
							propertiesMap["Template"] = propertiesArg["template"].(string)
						}
						sinkRocketMqParametersMap["Properties"] = propertiesMap
					}
					if sinkRocketMqParametersArg["tags"] != nil {
						tagsMap := make(map[string]interface{})
						for _, tags := range sinkRocketMqParametersArg["tags"].([]interface{}) {
							tagsArg := tags.(map[string]interface{})
							tagsMap["Form"] = tagsArg["form"].(string)
							tagsMap["Value"] = tagsArg["value"].(string)
							tagsMap["Template"] = tagsArg["template"].(string)
						}
						sinkRocketMqParametersMap["Tags"] = tagsMap
					}
					if sinkRocketMqParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkRocketMqParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkRocketMqParametersMap["Topic"] = topicMap
					}
				}
				sinkMap["SinkRocketMqParameters"] = sinkRocketMqParametersMap
			}

			if sinkArg["sink_sls_parameters"] != nil {
				sinkSlsParametersMap := make(map[string]interface{})
				for _, sinkSlsParameters := range sinkArg["sink_sls_parameters"].([]interface{}) {
					sinkSlsParametersArg := sinkSlsParameters.(map[string]interface{})
					if sinkSlsParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkSlsParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkSlsParametersMap["Body"] = bodyMap
					}
					if sinkSlsParametersArg["log_store"] != nil {
						logStoreMap := make(map[string]interface{})
						for _, logStore := range sinkSlsParametersArg["log_store"].([]interface{}) {
							logStoreArg := logStore.(map[string]interface{})
							logStoreMap["Form"] = logStoreArg["form"].(string)
							logStoreMap["Value"] = logStoreArg["value"].(string)
							logStoreMap["Template"] = logStoreArg["template"].(string)
						}
						sinkSlsParametersMap["LogStore"] = logStoreMap
					}
					if sinkSlsParametersArg["project"] != nil {
						projectMap := make(map[string]interface{})
						for _, project := range sinkSlsParametersArg["project"].([]interface{}) {
							projectArg := project.(map[string]interface{})
							projectMap["Form"] = projectArg["form"].(string)
							projectMap["Value"] = projectArg["value"].(string)
							projectMap["Template"] = projectArg["template"].(string)
						}
						sinkSlsParametersMap["Project"] = projectMap
					}
					if sinkSlsParametersArg["role_name"] != nil {
						roleNameMap := make(map[string]interface{})
						for _, roleName := range sinkSlsParametersArg["role_name"].([]interface{}) {
							roleNameArg := roleName.(map[string]interface{})
							roleNameMap["Form"] = roleNameArg["form"].(string)
							roleNameMap["Value"] = roleNameArg["value"].(string)
							roleNameMap["Template"] = roleNameArg["template"].(string)
						}
						sinkSlsParametersMap["RoleName"] = roleNameMap
					}
					if sinkSlsParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkSlsParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkSlsParametersMap["Topic"] = topicMap
					}
				}
				sinkMap["SinkSLSParameters"] = sinkSlsParametersMap
			}

		}
		request["Sink"], _ = convertArrayObjectToJsonString(sinkMap)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_event_streaming", action, AlibabaCloudSdkGoERROR)
	}
	if v, ok := response["Success"]; !ok || fmt.Sprint(v) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["EventStreamingName"]))

	return resourceAlicloudEventBridgeEventStreamingRead(d, meta)
}

func resourceAlicloudEventBridgeEventStreamingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}

	object, err := eventbridgeService.DescribeEventBridgeEventStreaming(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_streaming eventbridgeService.DescribeEventBridgeEventStreaming Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_streaming_name", d.Id())
	d.Set("description", object["Description"])
	d.Set("filter_pattern", object["FilterPattern"])

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
	d.Set("run_options", runOptionsSli)

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
	d.Set("sink", sinkSli)

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
	d.Set("source", sourceSli)
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudEventBridgeEventStreamingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"EventStreamingName": d.Id(),
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if d.HasChange("filter_pattern") {
		update = true
	}
	request["FilterPattern"] = d.Get("filter_pattern")

	if d.HasChange("run_options") {
		update = true
	}
	if v, ok := d.GetOk("run_options"); ok {
		runOptionsMap := make(map[string]interface{})
		for _, runOptions := range v.([]interface{}) {
			runOptionsArg := runOptions.(map[string]interface{})
			if runOptionsArg["batch_window"] != nil {
				batchWindowMap := make(map[string]interface{})
				for _, critical := range runOptionsArg["batch_window"].([]interface{}) {
					batchWindowArg := critical.(map[string]interface{})
					batchWindowMap["CountBasedWindow"] = batchWindowArg["count_based_window"].(int)
					batchWindowMap["TimeBasedWindow"] = batchWindowArg["time_based_window"].(int)
				}
				runOptionsMap["BatchWindow"] = batchWindowMap
			}
			if runOptionsArg["dead_letter_queue"] != nil {
				deadLetterMap := make(map[string]interface{})
				for _, deadLetter := range runOptionsArg["dead_letter_queue"].([]interface{}) {
					deadLetterArg := deadLetter.(map[string]interface{})
					deadLetterMap["Arn"] = deadLetterArg["arn"].(string)
				}
				runOptionsMap["DeadLetterQueue"] = deadLetterMap
			}
			if runOptionsArg["retry_strategy"] != nil {
				retryStrategyMap := make(map[string]interface{})
				for _, retryStrategy := range runOptionsArg["retry_strategy"].([]interface{}) {
					retryStrategyArg := retryStrategy.(map[string]interface{})
					retryStrategyMap["MaximumEventAgeInSeconds"] = retryStrategyArg["maximum_event_age_in_seconds"].(int)
					retryStrategyMap["MaximumRetryAttempts"] = retryStrategyArg["maximum_retry_attempts"].(int)
					retryStrategyMap["PushRetryStrategy"] = retryStrategyArg["push_retry_strategy"].(string)
				}
				runOptionsMap["RetryStrategy"] = retryStrategyMap
			}
			if runOptionsArg["maximum_tasks"] != nil {
				runOptionsMap["MaximumTasks"] = runOptionsArg["maximum_tasks"]
			}
			if runOptionsArg["errors_tolerance"] != nil {
				runOptionsMap["ErrorsTolerance"] = runOptionsArg["errors_tolerance"]
			}
		}
		request["RunOptions"], _ = convertArrayObjectToJsonString(runOptionsMap)
	}
	if d.HasChange("source") {
		update = true
	}
	if v, ok := d.GetOk("source"); ok {
		sourceMap := make(map[string]interface{})
		for _, source := range v.([]interface{}) {
			sourceArg := source.(map[string]interface{})
			if sourceArg["source_dts_parameters"] != nil {
				sourceDtsParametersMap := make(map[string]interface{})
				for _, sourceDtsParameters := range sourceArg["source_dts_parameters"].([]interface{}) {
					sourceDtsParametersArg := sourceDtsParameters.(map[string]interface{})
					sourceDtsParametersMap["BrokerUrl"] = sourceDtsParametersArg["broker_url"].(string)
					sourceDtsParametersMap["InitCheckPoint"] = sourceDtsParametersArg["init_check_point"].(string)
					sourceDtsParametersMap["Password"] = sourceDtsParametersArg["password"].(string)
					sourceDtsParametersMap["Sid"] = sourceDtsParametersArg["sid"].(string)
					sourceDtsParametersMap["TaskId"] = sourceDtsParametersArg["task_id"].(string)
					sourceDtsParametersMap["Topic"] = sourceDtsParametersArg["topic"].(string)
					sourceDtsParametersMap["Username"] = sourceDtsParametersArg["username"].(string)
				}
				sourceMap["SourceDtsParameters"] = sourceDtsParametersMap
			}
			if sourceArg["source_kafka_parameters"] != nil {
				sourceKafkaParametersMap := make(map[string]interface{})
				for _, sourceKafkaParameters := range sourceArg["source_kafka_parameters"].([]interface{}) {
					sourceKafkaParametersArg := sourceKafkaParameters.(map[string]interface{})
					sourceKafkaParametersMap["ConsumerGroup"] = sourceKafkaParametersArg["consumer_group"].(string)
					sourceKafkaParametersMap["InstanceId"] = sourceKafkaParametersArg["instance_id"].(string)
					sourceKafkaParametersMap["Network"] = sourceKafkaParametersArg["network"].(string)
					sourceKafkaParametersMap["OffsetReset"] = sourceKafkaParametersArg["offset_reset"].(string)
					sourceKafkaParametersMap["RegionId"] = sourceKafkaParametersArg["region_id"].(string)
					sourceKafkaParametersMap["SecurityGroupId"] = sourceKafkaParametersArg["security_group_id"].(string)
					sourceKafkaParametersMap["VSwitchIds"] = sourceKafkaParametersArg["vswitch_ids"].(string)
					sourceKafkaParametersMap["Topic"] = sourceKafkaParametersArg["topic"].(string)
					sourceKafkaParametersMap["VpcId"] = sourceKafkaParametersArg["vpc_id"].(string)
				}
				sourceMap["SourceKafkaParameters"] = sourceKafkaParametersMap
			}
			if sourceArg["source_mns_parameters"] != nil {
				sourceMnsParametersMap := make(map[string]interface{})
				for _, sourceMnsParameters := range sourceArg["source_mns_parameters"].([]interface{}) {
					sourceMnsParametersArg := sourceMnsParameters.(map[string]interface{})
					sourceMnsParametersMap["IsBase64Decode"] = sourceMnsParametersArg["is_base64_decode"].(bool)
					sourceMnsParametersMap["QueueName"] = sourceMnsParametersArg["queue_name"].(string)
					sourceMnsParametersMap["RegionId"] = sourceMnsParametersArg["region_id"].(string)
				}
				sourceMap["SourceMNSParameters"] = sourceMnsParametersMap
			}
			if sourceArg["source_mqtt_parameters"] != nil {
				sourceMqttParametersMap := make(map[string]interface{})
				for _, sourceMqttParameters := range sourceArg["source_mqtt_parameters"].([]interface{}) {
					sourceMqttParametersArg := sourceMqttParameters.(map[string]interface{})
					sourceMqttParametersMap["InstanceId"] = sourceMqttParametersArg["instance_id"].(string)
					sourceMqttParametersMap["RegionId"] = sourceMqttParametersArg["region_id"].(string)
					sourceMqttParametersMap["Topic"] = sourceMqttParametersArg["topic"].(string)
				}
				sourceMap["SourceMQTTParameters"] = sourceMqttParametersMap
			}
			if sourceArg["source_rabbit_mq_parameters"] != nil {
				sourceRabbitMqParametersMap := make(map[string]interface{})
				for _, sourceRabbitMqParameters := range sourceArg["source_rabbit_mq_parameters"].([]interface{}) {
					sourceRabbitMqParametersArg := sourceRabbitMqParameters.(map[string]interface{})
					sourceRabbitMqParametersMap["InstanceId"] = sourceRabbitMqParametersArg["instance_id"].(string)
					sourceRabbitMqParametersMap["QueueName"] = sourceRabbitMqParametersArg["queue_name"].(string)
					sourceRabbitMqParametersMap["RegionId"] = sourceRabbitMqParametersArg["region_id"].(string)
					sourceRabbitMqParametersMap["VirtualHostName"] = sourceRabbitMqParametersArg["virtual_host_name"].(string)
				}
				sourceMap["SourceRabbitMQParameters"] = sourceRabbitMqParametersMap
			}
			if sourceArg["source_rocket_mq_parameters"] != nil {
				sourceRocketMqParametersMap := make(map[string]interface{})
				for _, sourceRocketMqParameters := range sourceArg["source_rocket_mq_parameters"].([]interface{}) {
					sourceRocketMqParametersArg := sourceRocketMqParameters.(map[string]interface{})
					sourceRocketMqParametersMap["GroupId"] = sourceRocketMqParametersArg["group_id"].(string)
					sourceRocketMqParametersMap["InstanceId"] = sourceRocketMqParametersArg["instance_id"].(string)
					sourceRocketMqParametersMap["Offset"] = sourceRocketMqParametersArg["offset"].(string)
					sourceRocketMqParametersMap["RegionId"] = sourceRocketMqParametersArg["region_id"].(string)
					sourceRocketMqParametersMap["Tag"] = sourceRocketMqParametersArg["tag"].(string)
					sourceRocketMqParametersMap["Timestamp"] = sourceRocketMqParametersArg["timestamp"].(int)
					sourceRocketMqParametersMap["Topic"] = sourceRocketMqParametersArg["topic"].(string)
				}
				sourceMap["SourceRocketMQParameters"] = sourceRocketMqParametersMap
			}
			if sourceArg["source_sls_parameters"] != nil {
				sourceSlsParametersMap := make(map[string]interface{})
				for _, sourceSlsParameters := range sourceArg["source_sls_parameters"].([]interface{}) {
					sourceSlsParametersArg := sourceSlsParameters.(map[string]interface{})
					sourceSlsParametersMap["ConsumePosition"] = sourceSlsParametersArg["consume_position"].(string)
					sourceSlsParametersMap["LogStore"] = sourceSlsParametersArg["log_store"].(string)
					sourceSlsParametersMap["Project"] = sourceSlsParametersArg["project"].(string)
					sourceSlsParametersMap["RoleName"] = sourceSlsParametersArg["role_name"].(string)
				}
				sourceMap["SourceSLSParameters"] = sourceSlsParametersMap
			}
		}
		request["Source"], _ = convertArrayObjectToJsonString(sourceMap)
	}
	if d.HasChange("sink") {
		update = true
	}
	if v, ok := d.GetOk("sink"); ok {
		sinkMap := make(map[string]interface{})
		for _, sink := range v.([]interface{}) {
			sinkArg := sink.(map[string]interface{})
			if sinkArg["sink_fc_parameters"] != nil {
				sinkFcParametersMap := make(map[string]interface{})
				for _, critical := range sinkArg["sink_fc_parameters"].([]interface{}) {
					sinkFcParametersArg := critical.(map[string]interface{})
					if sinkFcParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkFcParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
						}
						sinkFcParametersMap["Body"] = bodyMap
					}
					if sinkFcParametersArg["function_name"] != nil {
						functionNameMap := make(map[string]interface{})
						for _, functionName := range sinkFcParametersArg["function_name"].([]interface{}) {
							functionNameArg := functionName.(map[string]interface{})
							functionNameMap["Form"] = functionNameArg["form"].(string)
							functionNameMap["Value"] = functionNameArg["value"].(string)
							functionNameMap["Template"] = functionNameArg["template"].(string)
						}
						sinkFcParametersMap["FunctionName"] = functionNameMap
					}
					if sinkFcParametersArg["qualifier"] != nil {
						qualifierMap := make(map[string]interface{})
						for _, qualifier := range sinkFcParametersArg["qualifier"].([]interface{}) {
							qualifierArg := qualifier.(map[string]interface{})
							qualifierMap["Form"] = qualifierArg["form"].(string)
							qualifierMap["Value"] = qualifierArg["value"].(string)
							qualifierMap["Template"] = qualifierArg["template"].(string)
						}
						sinkFcParametersMap["Qualifier"] = qualifierMap
					}
					if sinkFcParametersArg["invocation_type"] != nil {
						invocationTypeMap := make(map[string]interface{})
						for _, invocationType := range sinkFcParametersArg["invocation_type"].([]interface{}) {
							invocationTypeArg := invocationType.(map[string]interface{})
							invocationTypeMap["Form"] = invocationTypeArg["form"].(string)
							invocationTypeMap["Value"] = invocationTypeArg["value"].(string)
							invocationTypeMap["Template"] = invocationTypeArg["template"].(string)
						}
						sinkFcParametersMap["InvocationType"] = invocationTypeMap
					}
					if sinkFcParametersArg["service_name"] != nil {
						serviceNameMap := make(map[string]interface{})
						for _, serviceName := range sinkFcParametersArg["service_name"].([]interface{}) {
							serviceNameArg := serviceName.(map[string]interface{})
							serviceNameMap["Form"] = serviceNameArg["form"].(string)
							serviceNameMap["Value"] = serviceNameArg["value"].(string)
							serviceNameMap["Template"] = serviceNameArg["template"].(string)
						}
						sinkFcParametersMap["ServiceName"] = serviceNameMap
					}
				}
				sinkMap["SinkFcParameters"] = sinkFcParametersMap
			}

			if sinkArg["sink_kafka_parameters"] != nil {
				sinkKafkaParametersMap := make(map[string]interface{})
				for _, sinkKafkaParameters := range sinkArg["sink_kafka_parameters"].([]interface{}) {
					sinkKafkaParametersArg := sinkKafkaParameters.(map[string]interface{})
					if sinkKafkaParametersArg["acks"] != nil {
						acksMap := make(map[string]interface{})
						for _, acks := range sinkKafkaParametersArg["acks"].([]interface{}) {
							acksArg := acks.(map[string]interface{})
							acksMap["Form"] = acksArg["form"].(string)
							acksMap["Value"] = acksArg["value"].(string)
							acksMap["Template"] = acksArg["template"].(string)
						}
						sinkKafkaParametersMap["Acks"] = acksMap
					}
					if sinkKafkaParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkKafkaParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkKafkaParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkKafkaParametersArg["key"] != nil {
						keyMap := make(map[string]interface{})
						for _, key := range sinkKafkaParametersArg["key"].([]interface{}) {
							keyArg := key.(map[string]interface{})
							keyMap["Form"] = keyArg["form"].(string)
							keyMap["Value"] = keyArg["value"].(string)
							keyMap["Template"] = keyArg["template"].(string)
						}
						sinkKafkaParametersMap["Key"] = keyMap
					}
					if sinkKafkaParametersArg["sasl_user"] != nil {
						saslUserMap := make(map[string]interface{})
						for _, saslUser := range sinkKafkaParametersArg["sasl_user"].([]interface{}) {
							saslUserArg := saslUser.(map[string]interface{})
							saslUserMap["Form"] = saslUserArg["form"].(string)
							saslUserMap["Value"] = saslUserArg["value"].(string)
							saslUserMap["Template"] = saslUserArg["template"].(string)
						}
						sinkKafkaParametersMap["SaslUser"] = saslUserMap
					}
					if sinkKafkaParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkKafkaParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkKafkaParametersMap["Topic"] = topicMap
					}
					if sinkKafkaParametersArg["value"] != nil {
						valueMap := make(map[string]interface{})
						for _, value := range sinkKafkaParametersArg["value"].([]interface{}) {
							valueArg := value.(map[string]interface{})
							valueMap["Form"] = valueArg["form"].(string)
							valueMap["Value"] = valueArg["value"].(string)
							valueMap["Template"] = valueArg["template"].(string)
						}
						sinkKafkaParametersMap["Value"] = valueMap
					}
				}
				sinkMap["SinkKafkaParameters"] = sinkKafkaParametersMap
			}

			if sinkArg["sink_mns_parameters"] != nil {
				sinkMnsParametersMap := make(map[string]interface{})
				for _, sinkMnsParameters := range sinkArg["sink_mns_parameters"].([]interface{}) {
					sinkMnsParametersArg := sinkMnsParameters.(map[string]interface{})
					if sinkMnsParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkMnsParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkMnsParametersMap["Body"] = bodyMap
					}
					if sinkMnsParametersArg["is_base64_encode"] != nil {
						isBase64EncodeMap := make(map[string]interface{})
						for _, isBase64Encode := range sinkMnsParametersArg["is_base64_encode"].([]interface{}) {
							isBase64EncodeArg := isBase64Encode.(map[string]interface{})
							isBase64EncodeMap["Form"] = isBase64EncodeArg["form"].(string)
							isBase64EncodeMap["Value"] = isBase64EncodeArg["value"].(string)
							isBase64EncodeMap["Template"] = isBase64EncodeArg["template"].(string)
						}
						sinkMnsParametersMap["IsBase64Encode"] = isBase64EncodeMap
					}
					if sinkMnsParametersArg["queue_name"] != nil {
						queueNameMap := make(map[string]interface{})
						for _, queueName := range sinkMnsParametersArg["queue_name"].([]interface{}) {
							queueNameArg := queueName.(map[string]interface{})
							queueNameMap["Form"] = queueNameArg["form"].(string)
							queueNameMap["Value"] = queueNameArg["value"].(string)
							queueNameMap["Template"] = queueNameArg["template"].(string)
						}
						sinkMnsParametersMap["QueueName"] = queueNameMap
					}
				}
				sinkMap["SinkMNSParameters"] = sinkMnsParametersMap
			}

			if sinkArg["sink_rabbit_mq_parameters"] != nil {
				sinkRabbitMqParametersMap := make(map[string]interface{})
				for _, sinkRabbitMqParameters := range sinkArg["sink_rabbit_mq_parameters"].([]interface{}) {
					sinkRabbitMqParametersArg := sinkRabbitMqParameters.(map[string]interface{})
					if sinkRabbitMqParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkRabbitMqParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Body"] = bodyMap
					}
					if sinkRabbitMqParametersArg["exchange"] != nil {
						exchangeMap := make(map[string]interface{})
						for _, exchange := range sinkRabbitMqParametersArg["exchange"].([]interface{}) {
							exchangeArg := exchange.(map[string]interface{})
							exchangeMap["Form"] = exchangeArg["form"].(string)
							exchangeMap["Value"] = exchangeArg["value"].(string)
							exchangeMap["Template"] = exchangeArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Exchange"] = exchangeMap
					}
					if sinkRabbitMqParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkRabbitMqParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkRabbitMqParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkRabbitMqParametersArg["message_id"] != nil {
						messageIdMap := make(map[string]interface{})
						for _, messageId := range sinkRabbitMqParametersArg["message_id"].([]interface{}) {
							messageIdArg := messageId.(map[string]interface{})
							messageIdMap["Form"] = messageIdArg["form"].(string)
							messageIdMap["Value"] = messageIdArg["value"].(string)
							messageIdMap["Template"] = messageIdArg["template"].(string)
						}
						sinkRabbitMqParametersMap["MessageId"] = messageIdMap
					}
					if sinkRabbitMqParametersArg["properties"] != nil {
						propertiesMap := make(map[string]interface{})
						for _, properties := range sinkRabbitMqParametersArg["properties"].([]interface{}) {
							propertiesArg := properties.(map[string]interface{})
							propertiesMap["Form"] = propertiesArg["form"].(string)
							propertiesMap["Value"] = propertiesArg["value"].(string)
							propertiesMap["Template"] = propertiesArg["template"].(string)
						}
						sinkRabbitMqParametersMap["Properties"] = propertiesMap
					}
					if sinkRabbitMqParametersArg["queue_name"] != nil {
						queueNameMap := make(map[string]interface{})
						for _, queueName := range sinkRabbitMqParametersArg["queue_name"].([]interface{}) {
							queueNameArg := queueName.(map[string]interface{})
							queueNameMap["Form"] = queueNameArg["form"].(string)
							queueNameMap["Value"] = queueNameArg["value"].(string)
							queueNameMap["Template"] = queueNameArg["template"].(string)
						}
						sinkRabbitMqParametersMap["QueueName"] = queueNameMap
					}
					if sinkRabbitMqParametersArg["routing_key"] != nil {
						routingKeyMap := make(map[string]interface{})
						for _, routingKey := range sinkRabbitMqParametersArg["routing_key"].([]interface{}) {
							routingKeyArg := routingKey.(map[string]interface{})
							routingKeyMap["Form"] = routingKeyArg["form"].(string)
							routingKeyMap["Value"] = routingKeyArg["value"].(string)
							routingKeyMap["Template"] = routingKeyArg["template"].(string)
						}
						sinkRabbitMqParametersMap["RoutingKey"] = routingKeyMap
					}
					if sinkRabbitMqParametersArg["target_type"] != nil {
						targetTypeMap := make(map[string]interface{})
						for _, targetType := range sinkRabbitMqParametersArg["target_type"].([]interface{}) {
							targetTypeArg := targetType.(map[string]interface{})
							targetTypeMap["Form"] = targetTypeArg["form"].(string)
							targetTypeMap["Value"] = targetTypeArg["value"].(string)
							targetTypeMap["Template"] = targetTypeArg["template"].(string)
						}
						sinkRabbitMqParametersMap["TargetType"] = targetTypeMap
					}
					if sinkRabbitMqParametersArg["virtual_host_name"] != nil {
						virtualHostNameMap := make(map[string]interface{})
						for _, virtualHostName := range sinkRabbitMqParametersArg["virtual_host_name"].([]interface{}) {
							virtualHostNameArg := virtualHostName.(map[string]interface{})
							virtualHostNameMap["Form"] = virtualHostNameArg["form"].(string)
							virtualHostNameMap["Value"] = virtualHostNameArg["value"].(string)
							virtualHostNameMap["Template"] = virtualHostNameArg["template"].(string)
						}
						sinkRabbitMqParametersMap["VirtualHostName"] = virtualHostNameMap
					}
				}
				sinkMap["SinkRabbitMqParameters"] = sinkRabbitMqParametersMap
			}

			if sinkArg["sink_rocket_mq_parameters"] != nil {
				sinkRocketMqParametersMap := make(map[string]interface{})
				for _, sinkRocketMqParameters := range sinkArg["sink_rocket_mq_parameters"].([]interface{}) {
					sinkRocketMqParametersArg := sinkRocketMqParameters.(map[string]interface{})
					if sinkRocketMqParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkRocketMqParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkRocketMqParametersMap["Body"] = bodyMap
					}
					if sinkRocketMqParametersArg["instance_id"] != nil {
						instanceIdMap := make(map[string]interface{})
						for _, instanceId := range sinkRocketMqParametersArg["instance_id"].([]interface{}) {
							instanceIdArg := instanceId.(map[string]interface{})
							instanceIdMap["Form"] = instanceIdArg["form"].(string)
							instanceIdMap["Value"] = instanceIdArg["value"].(string)
							instanceIdMap["Template"] = instanceIdArg["template"].(string)
						}
						sinkRocketMqParametersMap["InstanceId"] = instanceIdMap
					}
					if sinkRocketMqParametersArg["keys"] != nil {
						keysMap := make(map[string]interface{})
						for _, keys := range sinkRocketMqParametersArg["keys"].([]interface{}) {
							keysArg := keys.(map[string]interface{})
							keysMap["Form"] = keysArg["form"].(string)
							keysMap["Value"] = keysArg["value"].(string)
							keysMap["Template"] = keysArg["template"].(string)
						}
						sinkRocketMqParametersMap["Keys"] = keysMap
					}
					if sinkRocketMqParametersArg["properties"] != nil {
						propertiesMap := make(map[string]interface{})
						for _, properties := range sinkRocketMqParametersArg["properties"].([]interface{}) {
							propertiesArg := properties.(map[string]interface{})
							propertiesMap["Form"] = propertiesArg["form"].(string)
							propertiesMap["Value"] = propertiesArg["value"].(string)
							propertiesMap["Template"] = propertiesArg["template"].(string)
						}
						sinkRocketMqParametersMap["Properties"] = propertiesMap
					}
					if sinkRocketMqParametersArg["tags"] != nil {
						tagsMap := make(map[string]interface{})
						for _, tags := range sinkRocketMqParametersArg["tags"].([]interface{}) {
							tagsArg := tags.(map[string]interface{})
							tagsMap["Form"] = tagsArg["form"].(string)
							tagsMap["Value"] = tagsArg["value"].(string)
							tagsMap["Template"] = tagsArg["template"].(string)
						}
						sinkRocketMqParametersMap["Tags"] = tagsMap
					}
					if sinkRocketMqParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkRocketMqParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkRocketMqParametersMap["Topic"] = topicMap
					}
				}
				sinkMap["SinkRocketMqParameters"] = sinkRocketMqParametersMap
			}

			if sinkArg["sink_sls_parameters"] != nil {
				sinkSlsParametersMap := make(map[string]interface{})
				for _, sinkSlsParameters := range sinkArg["sink_sls_parameters"].([]interface{}) {
					sinkSlsParametersArg := sinkSlsParameters.(map[string]interface{})
					if sinkSlsParametersArg["body"] != nil {
						bodyMap := make(map[string]interface{})
						for _, body := range sinkSlsParametersArg["body"].([]interface{}) {
							bodyArg := body.(map[string]interface{})
							bodyMap["Form"] = bodyArg["form"].(string)
							bodyMap["Value"] = bodyArg["value"].(string)
							bodyMap["Template"] = bodyArg["template"].(string)
						}
						sinkSlsParametersMap["Body"] = bodyMap
					}
					if sinkSlsParametersArg["log_store"] != nil {
						logStoreMap := make(map[string]interface{})
						for _, logStore := range sinkSlsParametersArg["log_store"].([]interface{}) {
							logStoreArg := logStore.(map[string]interface{})
							logStoreMap["Form"] = logStoreArg["form"].(string)
							logStoreMap["Value"] = logStoreArg["value"].(string)
							logStoreMap["Template"] = logStoreArg["template"].(string)
						}
						sinkSlsParametersMap["LogStore"] = logStoreMap
					}
					if sinkSlsParametersArg["project"] != nil {
						projectMap := make(map[string]interface{})
						for _, project := range sinkSlsParametersArg["project"].([]interface{}) {
							projectArg := project.(map[string]interface{})
							projectMap["Form"] = projectArg["form"].(string)
							projectMap["Value"] = projectArg["value"].(string)
							projectMap["Template"] = projectArg["template"].(string)
						}
						sinkSlsParametersMap["Project"] = projectMap
					}
					if sinkSlsParametersArg["role_name"] != nil {
						roleNameMap := make(map[string]interface{})
						for _, roleName := range sinkSlsParametersArg["role_name"].([]interface{}) {
							roleNameArg := roleName.(map[string]interface{})
							roleNameMap["Form"] = roleNameArg["form"].(string)
							roleNameMap["Value"] = roleNameArg["value"].(string)
							roleNameMap["Template"] = roleNameArg["template"].(string)
						}
						sinkSlsParametersMap["RoleName"] = roleNameMap
					}
					if sinkSlsParametersArg["topic"] != nil {
						topicMap := make(map[string]interface{})
						for _, topic := range sinkSlsParametersArg["topic"].([]interface{}) {
							topicArg := topic.(map[string]interface{})
							topicMap["Form"] = topicArg["form"].(string)
							topicMap["Value"] = topicArg["value"].(string)
							topicMap["Template"] = topicArg["template"].(string)
						}
						sinkSlsParametersMap["Topic"] = topicMap
					}
				}
				sinkMap["SinkSLSParameters"] = sinkSlsParametersMap
			}

		}
		request["Sink"], _ = convertArrayObjectToJsonString(sinkMap)
	}

	if update {
		action := "UpdateEventStreaming"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudEventBridgeEventStreamingRead(d, meta)
}

func resourceAlicloudEventBridgeEventStreamingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"EventStreamingName": d.Id(),
	}

	action := "DeleteEventStreaming"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
