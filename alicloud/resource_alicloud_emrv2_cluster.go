package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEmrV2Cluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEmrV2ClusterCreate,
		Read:   resourceAlicloudEmrV2ClusterRead,
		Update: resourceAlicloudEmrV2ClusterUpdate,
		Delete: resourceAlicloudEmrV2ClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"subscription_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_duration_unit": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
						},
						"payment_duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_pay_order": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_renew_duration_unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
						},
						"auto_renew_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
						},
					},
				},
				MaxItems: 1,
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DATALAKE", "OLAP", "DATAFLOW", "DATASERVING", "CUSTOM"}, false),
			},
			"release_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_collect_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"deploy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "HA"}, false),
			},
			"security_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "KERBEROS"}, false),
			},
			"applications": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"application_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_file_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_item_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_item_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_scope": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"CLUSTER", "NODE_GROUP"}, false),
						},
						"config_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"node_attributes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ram_role": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"data_disk_encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"data_disk_kms_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"system_disk_encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"system_disk_kms_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
				ForceNew: true,
			},
			"node_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_group_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"MASTER", "CORE", "TASK", "GATEWAY", "MASTER-EXTEND"}, false),
						},
						"node_group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"payment_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
						},
						"deployment_set_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "CLUSTER", "NODE_GROUP"}, false),
						},
						"node_resize_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"PRIORITY", "COST_OPTIMIZED"}, false),
						},
						"subscription_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"payment_duration_unit": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
									},
									"payment_duration": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
									},
									"auto_renew": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"auto_pay_order": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"auto_renew_duration_unit": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
									},
									"auto_renew_duration": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
									},
								},
							},
							MaxItems: 1,
						},
						"spot_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"NoSpot", "SpotWithPriceLimit", "SpotAsPriceGo"}, false),
						},
						"spot_bid_prices": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"bid_price": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"vswitch_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"with_public_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"additional_security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_types": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"node_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"system_disk": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency", "cloud_ssd"}, false),
									},
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"performance_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
									},
									"count": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
						"data_disks": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud", "local_hdd_pro", "local_disk", "local_ssd_pro"}, false),
									},
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"performance_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
									},
									"count": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"graceful_shutdown": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"spot_instance_remedy": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"cost_optimized_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"on_demand_base_capacity": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"on_demand_percentage_above_base_capacity": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"spot_instance_pools": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
							MaxItems: 1,
						},
						"auto_scaling_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scaling_rules": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"trigger_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: StringInSlice([]string{"TIME_TRIGGER", "METRICS_TRIGGER"}, false),
												},
												"activity_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: StringInSlice([]string{"SCALE_OUT", "SCALE_IN"}, false),
												},
												"adjustment_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringInSlice([]string{"CHANGE_IN_CAPACITY", "EXACT_CAPACITY"}, false),
												},
												"adjustment_value": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: IntBetween(1, 5000),
												},
												"min_adjustment_value": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"time_trigger": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"launch_time": {
																Type:     schema.TypeString,
																Required: true,
															},
															"start_time": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"end_time": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"launch_expiration_time": {
																Type:         schema.TypeInt,
																Optional:     true,
																ValidateFunc: IntBetween(0, 3600),
															},
															"recurrence_type": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: StringInSlice([]string{"MINUTELY", "HOURLY", "DAILY", "WEEKLY", "MONTHLY"}, false),
															},
															"recurrence_value": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
													MaxItems: 1,
												},
												"metrics_trigger": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"time_window": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: IntBetween(30, 1800),
															},
															"evaluation_count": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: IntBetween(1, 5),
															},
															"cool_down_interval": {
																Type:         schema.TypeInt,
																Optional:     true,
																ValidateFunc: IntBetween(0, 10800),
															},
															"condition_logic_operator": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: StringInSlice([]string{"And", "Or"}, false),
															},
															"time_constraints": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_time": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"end_time": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"conditions": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_name": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"statistics": {
																			Type:         schema.TypeString,
																			Required:     true,
																			ValidateFunc: StringInSlice([]string{"MAX", "MIN", "AVG"}, false),
																		},
																		"comparison_operator": {
																			Type:         schema.TypeString,
																			Required:     true,
																			ValidateFunc: StringInSlice([]string{"EQ", "NE", "GT", "LT", "GE", "LE"}, false),
																		},
																		"threshold": {
																			Type:     schema.TypeFloat,
																			Required: true,
																		},
																		"tags": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"key": {
																						Type:     schema.TypeString,
																						Required: true,
																					},
																					"value": {
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
													MaxItems: 1,
												},
											},
										},
									},
									"constraints": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_capacity": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: IntBetween(0, 2000),
												},
												"min_capacity": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: IntBetween(0, 2000),
												},
											},
										},
										MaxItems: 1,
									},
								},
							},
							MaxItems: 1,
						},
						"ack_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ack_instance_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"node_selectors": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"tolerations": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"effect": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"namespace": {
										Type:     schema.TypeString,
										Required: true,
									},
									"request_cpu": {
										Type:     schema.TypeFloat,
										Required: true,
									},
									"request_memory": {
										Type:     schema.TypeFloat,
										Required: true,
									},
									"limit_cpu": {
										Type:     schema.TypeFloat,
										Required: true,
									},
									"limit_memory": {
										Type:     schema.TypeFloat,
										Required: true,
									},
									"custom_labels": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"custom_annotations": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"pvcs": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"data_disk_storage_class": {
													Type:     schema.TypeString,
													Required: true,
												},
												"data_disk_size": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"volumes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"volume_mounts": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"pre_start_command": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"pod_affinity": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"pod_anti_affinity": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"node_affinity": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
				MinItems: 1,
			},
			"bootstrap_scripts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"script_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"script_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"script_args": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:       schema.TypeInt,
							Optional:   true,
							Deprecated: "Field 'priority' has been deprecated from provider version 1.227.0.",
						},
						"execution_moment": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"BEFORE_INSTALL", "AFTER_STARTED", "BEFORE_START"}, false),
						},
						"execution_fail_strategy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"FAILED_CONTINUE", "FAILED_BLOCK"}, false),
						},
						"node_selector": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_select_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"CLUSTER", "NODE_GROUP", "NODE"}, false),
									},
									"node_names": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"node_group_id": {
										Type:       schema.TypeString,
										Optional:   true,
										Deprecated: "Field 'node_group_id' has been deprecated from provider version 1.227.0. New field 'node_group_ids' replaces it.",
									},
									"node_group_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"node_group_types": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"node_group_name": {
										Type:       schema.TypeString,
										Optional:   true,
										Deprecated: "Field 'node_group_name' has been deprecated from provider version 1.227.0. New field 'node_group_names' replaces it.",
									},
									"node_group_names": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudEmrV2ClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RunCluster"
	var err error
	createClusterRequest := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		createClusterRequest["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		createClusterRequest["PaymentType"] = v
	}

	if v, ok := d.GetOk("subscription_config"); ok {
		subscriptionConfig := v.(*schema.Set).List()
		if len(subscriptionConfig) == 1 {
			subscriptionConfigSource := subscriptionConfig[0].(map[string]interface{})
			subscriptionConfigMap := map[string]interface{}{
				"PaymentDurationUnit":   subscriptionConfigSource["payment_duration_unit"],
				"PaymentDuration":       subscriptionConfigSource["payment_duration"],
				"AutoRenew":             subscriptionConfigSource["auto_renew"],
				"AutoRenewDurationUnit": subscriptionConfigSource["auto_renew_duration_unit"],
				"AutoRenewDuration":     subscriptionConfigSource["auto_renew_duration"],
			}
			if value, exists := subscriptionConfigSource["auto_pay_order"]; exists {
				subscriptionConfigMap["AutoPayOrder"] = value.(bool)
			} else {
				subscriptionConfigMap["AutoPayOrder"] = true
			}

			createClusterRequest["SubscriptionConfig"] = convertMapToJsonStringIgnoreError(subscriptionConfigMap)
		}
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		createClusterRequest["ClusterType"] = v
	}

	if v, ok := d.GetOk("release_version"); ok {
		createClusterRequest["ReleaseVersion"] = v
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		createClusterRequest["ClusterName"] = v
	}

	if v, ok := d.GetOk("deploy_mode"); ok {
		createClusterRequest["DeployMode"] = v
	}

	if v, ok := d.GetOk("security_mode"); ok {
		createClusterRequest["SecurityMode"] = v
	}

	if v, ok := d.GetOk("log_collect_strategy"); ok && v.(string) != "" {
		createClusterRequest["LogCollectStrategy"] = v.(string)
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		createClusterRequest["DeletionProtection"] = v
	}

	applications := make([]map[string]interface{}, 0)
	if apps, ok := d.GetOk("applications"); ok {
		for _, application := range apps.(*schema.Set).List() {
			applications = append(applications, map[string]interface{}{"ApplicationName": application.(string)})
		}
	}
	createClusterRequest["Applications"], _ = convertListMapToJsonString(applications)

	applicationConfigs := make([]map[string]interface{}, 0)
	if appConfigs, ok := d.GetOk("application_configs"); ok {
		for _, appConfig := range appConfigs.(*schema.Set).List() {
			applicationConfig := map[string]interface{}{}
			kv := appConfig.(map[string]interface{})
			if v, ok := kv["application_name"]; ok {
				applicationConfig["ApplicationName"] = v
			}
			if v, ok := kv["config_file_name"]; ok {
				applicationConfig["ConfigFileName"] = v
			}
			if v, ok := kv["config_item_key"]; ok {
				applicationConfig["ConfigItemKey"] = v
			}
			if v, ok := kv["config_item_value"]; ok {
				applicationConfig["ConfigItemValue"] = v
			}
			if v, ok := kv["config_scope"]; ok {
				applicationConfig["ConfigScope"] = v
			}
			if v, ok := kv["node_group_name"]; ok {
				applicationConfig["NodeGroupName"] = v
			}
			if v, ok := kv["node_group_id"]; ok {
				applicationConfig["NodeGroupId"] = v
			}
			applicationConfigs = append(applicationConfigs, applicationConfig)
		}
	}
	createClusterRequest["ApplicationConfigs"], _ = convertListMapToJsonString(applicationConfigs)

	if v, ok := d.GetOk("node_attributes"); ok {
		nodeAttributes := v.(*schema.Set).List()
		if len(nodeAttributes) == 1 {
			nodeAttributesSource := nodeAttributes[0].(map[string]interface{})
			nodeAttributesSourceMap := map[string]interface{}{
				"VpcId":             nodeAttributesSource["vpc_id"],
				"ZoneId":            nodeAttributesSource["zone_id"],
				"SecurityGroupId":   nodeAttributesSource["security_group_id"],
				"RamRole":           nodeAttributesSource["ram_role"],
				"KeyPairName":       nodeAttributesSource["key_pair_name"],
				"DataDiskEncrypted": nodeAttributesSource["data_disk_encrypted"],
				"DataDiskKMSKeyId":  nodeAttributesSource["data_disk_kms_key_id"],
			}
			if value, exists := nodeAttributesSource["system_disk_encrypted"]; exists {
				nodeAttributesSourceMap["SystemDiskEncrypted"] = value.(bool)
			}
			if value, exists := nodeAttributesSource["system_disk_kms_key_id"]; exists && value.(string) != "" {
				nodeAttributesSourceMap["SystemDiskKMSKeyId"] = value.(string)
			}
			createClusterRequest["NodeAttributes"] = convertMapToJsonStringIgnoreError(nodeAttributesSourceMap)
		}
	}

	nodeGroups := make([]map[string]interface{}, 0)
	if nodeGroupsList, ok := d.GetOk("node_groups"); ok {
		for _, nodeGroupItem := range nodeGroupsList.([]interface{}) {
			nodeGroup := map[string]interface{}{}
			kv := nodeGroupItem.(map[string]interface{})
			if v, ok := kv["node_group_type"]; ok {
				nodeGroup["NodeGroupType"] = v
			}
			if v, ok := kv["node_group_name"]; ok {
				nodeGroup["NodeGroupName"] = v
			}
			if v, ok := kv["payment_type"]; ok {
				nodeGroup["PaymentType"] = v
			}
			if v, ok := kv["deployment_set_strategy"]; ok && v.(string) != "" {
				nodeGroup["DeploymentSetStrategy"] = v.(string)
			}
			if v, ok := kv["node_resize_strategy"]; ok && v.(string) != "" {
				nodeGroup["NodeResizeStrategy"] = v.(string)
			}
			if v, ok := kv["spot_strategy"]; ok && v.(string) != "" {
				nodeGroup["SpotStrategy"] = v.(string)
			}
			if v, ok := kv["subscription_config"]; ok {
				subscriptionConfigs := v.(*schema.Set).List()
				if len(subscriptionConfigs) == 1 {
					subscriptionConfig := map[string]interface{}{}
					subscriptionConfigMap := subscriptionConfigs[0].(map[string]interface{})
					if value, exists := subscriptionConfigMap["payment_duration_unit"]; exists {
						subscriptionConfig["PaymentDurationUnit"] = value
					}
					if value, exists := subscriptionConfigMap["payment_duration"]; exists {
						subscriptionConfig["PaymentDuration"] = value
					}
					if value, exists := subscriptionConfigMap["auto_renew"]; exists {
						subscriptionConfig["AutoRenew"] = value
					}
					if value, exists := subscriptionConfigMap["auto_renew_duration_unit"]; exists {
						subscriptionConfig["AutoRenewDurationUnit"] = value
					}
					if value, exists := subscriptionConfigMap["auto_renew_duration"]; exists {
						subscriptionConfig["AutoRenewDuration"] = value
					}
					if value, exists := subscriptionConfigMap["auto_pay_order"]; exists {
						subscriptionConfig["AutoPayOrder"] = value.(bool)
					} else {
						subscriptionConfig["AutoPayOrder"] = true
					}
					nodeGroup["SubscriptionConfig"] = subscriptionConfig
				}
			}
			if v, ok := kv["spot_bid_prices"]; ok {
				spotBidPriceList := v.(*schema.Set).List()
				if len(spotBidPriceList) > 0 {
					spotBidPrices := make([]map[string]interface{}, 0)
					for _, spotBidPriceSource := range spotBidPriceList {
						spotBidPrice := map[string]interface{}{}
						spotBidPriceMap := spotBidPriceSource.(map[string]interface{})
						if value, exists := spotBidPriceMap["instance_type"]; exists {
							spotBidPrice["InstanceType"] = value
						}
						if value, exists := spotBidPriceMap["bid_price"]; exists {
							spotBidPrice["BidPrice"] = value
						}
						spotBidPrices = append(spotBidPrices, spotBidPrice)
					}
					nodeGroup["SpotBidPrices"] = spotBidPrices
				}
			}
			if v, ok := kv["vswitch_ids"]; ok {
				var vSwitchIds []string
				for _, vSwitchId := range v.(*schema.Set).List() {
					vSwitchIds = append(vSwitchIds, vSwitchId.(string))
				}
				nodeGroup["VSwitchIds"] = vSwitchIds
			}
			if v, ok := kv["with_public_ip"]; ok {
				nodeGroup["WithPublicIp"] = v
			}
			if v, ok := kv["additional_security_group_ids"]; ok {
				var additionalSecurityGroupIds []string
				for _, additionalSecurityGroupId := range v.(*schema.Set).List() {
					additionalSecurityGroupIds = append(additionalSecurityGroupIds, additionalSecurityGroupId.(string))
				}
				nodeGroup["AdditionalSecurityGroupIds"] = additionalSecurityGroupIds
			}
			if v, ok := kv["instance_types"]; ok {
				var instanceTypes []string
				for _, instanceType := range v.(*schema.Set).List() {
					instanceTypes = append(instanceTypes, instanceType.(string))
				}
				nodeGroup["InstanceTypes"] = instanceTypes
			}
			if v, ok := kv["node_count"]; ok {
				nodeGroup["NodeCount"] = v
			}
			if v, ok := kv["system_disk"]; ok {
				systemDisks := v.(*schema.Set).List()
				if len(systemDisks) == 1 {
					systemDisk := map[string]interface{}{}
					systemDiskMap := systemDisks[0].(map[string]interface{})
					if value, exists := systemDiskMap["category"]; exists {
						systemDisk["Category"] = value
					}
					if value, exists := systemDiskMap["size"]; exists {
						systemDisk["Size"] = value
					}
					if value, exists := systemDiskMap["performance_level"]; exists && value.(string) != "" {
						systemDisk["PerformanceLevel"] = value
					}
					if value, exists := systemDiskMap["count"]; exists {
						systemDisk["Count"] = value
					}
					nodeGroup["SystemDisk"] = systemDisk
				}
			}
			if v, ok := kv["data_disks"]; ok {
				dataDiskList := v.(*schema.Set).List()
				if len(dataDiskList) > 0 {
					dataDisks := make([]map[string]interface{}, 0)
					for _, dataDiskSource := range dataDiskList {
						dataDisk := map[string]interface{}{}
						dataDiskMap := dataDiskSource.(map[string]interface{})
						if value, exists := dataDiskMap["category"]; exists {
							dataDisk["Category"] = value
						}
						if value, exists := dataDiskMap["size"]; exists {
							dataDisk["Size"] = value
						}
						if value, exists := dataDiskMap["performance_level"]; exists && value.(string) != "" {
							dataDisk["PerformanceLevel"] = value
						}
						if value, exists := dataDiskMap["count"]; exists {
							dataDisk["Count"] = value
						}
						dataDisks = append(dataDisks, dataDisk)
					}
					nodeGroup["DataDisks"] = dataDisks
				}
			}
			if v, ok := kv["graceful_shutdown"]; ok {
				nodeGroup["GracefulShutdown"] = v
			}
			if v, ok := kv["spot_instance_remedy"]; ok {
				nodeGroup["SpotInstanceRemedy"] = v
			}
			if v, ok := kv["cost_optimized_config"]; ok {
				costOptimizedConfigs := v.(*schema.Set).List()
				if len(costOptimizedConfigs) == 1 {
					costOptimizedConfig := map[string]interface{}{}
					costOptimizedConfigMap := costOptimizedConfigs[0].(map[string]interface{})
					if value, exists := costOptimizedConfigMap["on_demand_base_capacity"]; exists {
						costOptimizedConfig["OnDemandBaseCapacity"] = value
					}
					if value, exists := costOptimizedConfigMap["on_demand_percentage_above_base_capacity"]; exists {
						costOptimizedConfig["OnDemandPercentageAboveBaseCapacity"] = value
					}
					if value, exists := costOptimizedConfigMap["spot_instance_pools"]; exists {
						costOptimizedConfig["SpotInstancePools"] = value
					}
					nodeGroup["CostOptimizedConfig"] = costOptimizedConfig
				}
			}
			if v, ok := kv["auto_scaling_policy"]; ok {
				scalingPolicies := v.([]interface{})
				if len(scalingPolicies) == 1 {
					nodeGroup["AutoScalingPolicy"] = adaptAutoScalingPolicyRequest(scalingPolicies[0].(map[string]interface{}))
				}
			}
			if v, ok := kv["ack_config"]; ok {
				ackConfig := v.([]interface{})
				if len(ackConfig) == 1 {
					nodeGroup["AckConfig"] = adaptAckConfigRequest(ackConfig[0].(map[string]interface{}))
					nodeGroup["IaaSType"] = "K8S"
					delete(nodeGroup, "InstanceTypes")
				}
			}
			nodeGroups = append(nodeGroups, nodeGroup)
		}
	}
	createClusterRequest["NodeGroups"], _ = convertListMapToJsonString(nodeGroups)

	if scripts, ok := d.GetOk("bootstrap_scripts"); ok {
		bootstrapScripts := make([]map[string]interface{}, 0)
		for _, script := range scripts.([]interface{}) {
			kv := script.(map[string]interface{})
			bootstrapScript := map[string]interface{}{}
			if v, ok := kv["script_name"]; ok {
				bootstrapScript["ScriptName"] = v
			}
			if v, ok := kv["script_path"]; ok {
				bootstrapScript["ScriptPath"] = v
			}
			if v, ok := kv["script_args"]; ok {
				bootstrapScript["ScriptArgs"] = v
			}
			if v, ok := kv["priority"]; ok {
				bootstrapScript["Priority"] = v
			}
			if v, ok := kv["execution_moment"]; ok {
				bootstrapScript["ExecutionMoment"] = v
			}
			if v, ok := kv["execution_fail_strategy"]; ok {
				bootstrapScript["ExecutionFailStrategy"] = v
			}
			if v, ok := kv["node_selector"]; ok {
				nodeSelectors := v.(*schema.Set).List()
				if len(nodeSelectors) == 1 {
					nodeSelector := map[string]interface{}{}
					nodeSelectorMap := nodeSelectors[0].(map[string]interface{})
					if value, exists := nodeSelectorMap["node_select_type"]; exists {
						nodeSelector["NodeSelectType"] = value
					}
					if value, exists := nodeSelectorMap["node_names"]; exists {
						var nodeNames []string
						for _, nodeName := range value.([]interface{}) {
							nodeNames = append(nodeNames, nodeName.(string))
						}
						nodeSelector["NodeNames"] = nodeNames
					}
					if value, exists := nodeSelectorMap["node_group_id"]; exists {
						nodeSelector["NodeGroupId"] = value
					}
					if value, exists := nodeSelectorMap["node_group_ids"]; exists {
						var nodeGroupIds []string
						for _, ngId := range value.([]interface{}) {
							nodeGroupIds = append(nodeGroupIds, ngId.(string))
						}
						nodeSelector["NodeGroupIds"] = nodeGroupIds
					}
					if value, exists := nodeSelectorMap["node_group_types"]; exists {
						var nodeGroupTypes []string
						for _, nodeGroupType := range value.([]interface{}) {
							nodeGroupTypes = append(nodeGroupTypes, nodeGroupType.(string))
						}
						nodeSelector["NodeGroupTypes"] = nodeGroupTypes
					}
					if value, exists := nodeSelectorMap["node_group_name"]; exists {
						nodeSelector["NodeGroupName"] = value
					}
					if value, exists := nodeSelectorMap["node_group_names"]; exists {
						var nodeGroupNames []string
						for _, ngName := range value.([]interface{}) {
							nodeGroupNames = append(nodeGroupNames, ngName.(string))
						}
						nodeSelector["NodeGroupNames"] = nodeGroupNames
					}
					bootstrapScript["NodeSelector"] = nodeSelector
				}
			}
			bootstrapScripts = append(bootstrapScripts, bootstrapScript)
		}
		createClusterRequest["BootstrapScripts"], _ = convertListMapToJsonString(bootstrapScripts)
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		}
		createClusterRequest["Tags"], _ = convertListMapToJsonString(tags)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, createClusterRequest, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, createClusterRequest)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_emrv2_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))

	emrService := EmrService{client}
	stateConf := BuildStateConf([]string{"STARTING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate),
		90*time.Second, emrService.EmrV2ClusterStateRefreshFunc(d.Id(), []string{"START_FAILED", "TERMINATED_WITH_ERRORS", "TERMINATED"}))
	stateConf.PollInterval = 10 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEmrV2ClusterRead(d, meta)
}

func resourceAlicloudEmrV2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := EmrService{client}

	object, err := emrService.GetEmrV2Cluster(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_name", object["ClusterName"])
	d.Set("cluster_type", object["ClusterType"])
	d.Set("payment_type", object["PaymentType"])
	d.Set("release_version", object["ReleaseVersion"])
	d.Set("deploy_mode", object["DeployMode"])
	d.Set("security_mode", object["SecurityMode"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("log_collect_strategy", object["LogCollectStrategy"])
	d.Set("deletion_protection", object["DeletionProtection"])

	if _, ok := object["SubscriptionConfig"]; ok {
		sc := d.Get("subscription_config").(*schema.Set).List()
		if len(sc) > 0 {
			d.Set("subscription_config", []map[string]interface{}{sc[0].(map[string]interface{})})
		}
	}

	var nodeAttributes []map[string]interface{}
	if v, ok := object["NodeAttributes"]; ok {
		nodeAttributesMap := v.(map[string]interface{})
		nodeAttribute := map[string]interface{}{
			"vpc_id":                nodeAttributesMap["VpcId"],
			"zone_id":               nodeAttributesMap["ZoneId"],
			"security_group_id":     nodeAttributesMap["SecurityGroupId"],
			"ram_role":              nodeAttributesMap["RamRole"],
			"key_pair_name":         nodeAttributesMap["KeyPairName"],
			"data_disk_encrypted":   nodeAttributesMap["DataDiskEncrypted"],
			"system_disk_encrypted": nodeAttributesMap["SystemDiskEncrypted"],
		}
		if v, exists := nodeAttributesMap["DataDiskKMSKeyId"]; exists && v.(string) != "" {
			nodeAttribute["data_disk_kms_key_id"] = v
		}
		if v, exists := nodeAttributesMap["SystemDiskKMSKeyId"]; exists && v.(string) != "" {
			nodeAttribute["system_disk_kms_key_id"] = v
		}

		oldNodeAttributes := d.Get("node_attributes")
		if oldNodeAttributes != nil && oldNodeAttributes.(*schema.Set).Len() > 0 {
			oldNodeAttributesMap := d.Get("node_attributes").(*schema.Set).List()[0].(map[string]interface{})
			if _, exists := oldNodeAttributesMap["system_disk_encrypted"]; exists {
				nodeAttribute["system_disk_encrypted"] = nodeAttributesMap["SystemDiskEncrypted"]
			}
			if value, exists := oldNodeAttributesMap["system_disk_kms_key_id"]; exists && value != "" {
				nodeAttribute["system_disk_kms_key_id"] = nodeAttributesMap["SystemDiskKMSKeyId"]
			}
		}

		nodeAttributes = append(nodeAttributes, nodeAttribute)
		d.Set("node_attributes", nodeAttributes)
	}
	var response map[string]interface{}
	action := "ListApplications"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Applications", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Applications", response)
	}
	if v != nil && len(v.([]interface{})) > 0 {
		emrInternalApp := map[string]struct{}{
			"TRE-EXTENSION":        {},
			"KERBEROS":             {},
			"METASTORE":            {},
			"MYSQL":                {},
			"SPARK-EXTENSION":      {},
			"SPARK-NATIVE":         {},
			"TAIHAODOCTOR":         {},
			"EMRHOOK":              {},
			"JINDOSDK":             {},
			"TRE":                  {},
			"EMRRUNTIME":           {},
			"EMRRUNTIME-EXTENSION": {},
		}
		var applications []string
		for _, item := range v.([]interface{}) {
			app := strings.ToUpper(item.(map[string]interface{})["ApplicationName"].(string))
			if _, ok := emrInternalApp[app]; ok {
				continue
			}
			applications = append(applications, app)
		}
		d.Set("applications", applications)
	}

	if v, ok := d.GetOk("application_configs"); ok && len(v.(*schema.Set).List()) > 0 {
		var applicationConfigs []map[string]interface{}
		for _, ac := range v.(*schema.Set).List() {
			applicationConfigs = append(applicationConfigs, ac.(map[string]interface{}))
		}
		d.Set("application_configs", applicationConfigs)
	}

	action = "ListScripts"
	listScriptsRequest := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ClusterId":  d.Id(),
		"ScriptType": "BOOTSTRAP",
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, listScriptsRequest, true)
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
	v, err = jsonpath.Get("$.Scripts", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Scripts", response)
	}
	if v != nil && len(v.([]interface{})) > 0 {
		scriptsMaps := v.([]interface{})
		var scripts []map[string]interface{}
		for _, item := range scriptsMaps {
			scriptMap := item.(map[string]interface{})
			script := map[string]interface{}{
				"script_name":             scriptMap["ScriptName"],
				"script_path":             scriptMap["ScriptPath"],
				"script_args":             scriptMap["ScriptArgs"],
				"execution_moment":        scriptMap["ExecutionMoment"],
				"execution_fail_strategy": scriptMap["ExecutionFailStrategy"],
			}

			if v, ok := scriptMap["NodeSelector"]; ok && len(v.(map[string]interface{})) > 0 {
				nodeSelectorMap := v.(map[string]interface{})
				nodeSelector := map[string]interface{}{
					"node_select_type": nodeSelectorMap["NodeSelectType"],
				}

				if nodeGroupId, exists := nodeSelectorMap["NodeGroupId"]; exists && nodeGroupId.(string) != "" {
					nodeSelector["node_group_id"] = nodeGroupId
				}

				if ngIDs, exists := nodeSelectorMap["NodeGroupIds"]; exists && len(ngIDs.([]interface{})) > 0 {
					var nodeGroupIDs []string
					for _, ngID := range ngIDs.([]interface{}) {
						nodeGroupIDs = append(nodeGroupIDs, ngID.(string))
					}
					nodeSelector["node_group_ids"] = nodeGroupIDs
				}

				if nodeGroupName, exists := nodeSelectorMap["NodeGroupName"]; exists && nodeGroupName.(string) != "" {
					nodeSelector["node_group_name"] = nodeGroupName
				}

				if ngNames, exists := nodeSelectorMap["NodeGroupNames"]; exists && len(ngNames.([]interface{})) > 0 {
					var nodeGroupNames []string
					for _, ngName := range ngNames.([]interface{}) {
						nodeGroupNames = append(nodeGroupNames, ngName.(string))
					}
					nodeSelector["node_group_names"] = nodeGroupNames
				}

				if ngTypes, exists := nodeSelectorMap["NodeGroupTypes"]; exists && len(ngTypes.([]interface{})) > 0 {
					var nodeGroupTypes []string
					for _, ngType := range ngTypes.([]interface{}) {
						nodeGroupTypes = append(nodeGroupTypes, ngType.(string))
					}
					nodeSelector["node_group_types"] = nodeGroupTypes
				}

				if nn, exists := nodeSelectorMap["NodeNames"]; exists && len(nn.([]interface{})) > 0 {
					var nodeNames []string
					for _, nodeName := range nn.([]interface{}) {
						nodeNames = append(nodeNames, nodeName.(string))
					}
					nodeSelector["node_names"] = nodeNames
				}
				script["node_selector"] = []map[string]interface{}{nodeSelector}
			}
			scripts = append(scripts, script)
		}
		d.Set("bootstrap_scripts", scripts)
	}

	action = "ListNodeGroups"
	request["MaxResults"] = PageSizeLarge
	request["NodeGroupStates"] = []string{"RUNNING"}
	var nodeGroupObjects []interface{}

	for {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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

		nodeGroupResp, err := jsonpath.Get("$.NodeGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.NodeGroups", response)
		}
		if nodeGroupResp != nil && len(nodeGroupResp.([]interface{})) > 0 {
			for _, ngr := range nodeGroupResp.([]interface{}) {
				nodeGroupObjects = append(nodeGroupObjects, ngr)
			}
			_, nextTokenExists := response["NextToken"]
			if len(nodeGroupResp.([]interface{})) < PageSizeLarge {
				break
			} else if len(nodeGroupResp.([]interface{})) == PageSizeLarge && !nextTokenExists {
				break
			}
		} else {
			break
		}

		nextToken, err := jsonpath.Get("$.NextToken", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NextToken", response)
		}
		request["NextToken"] = nextToken
	}
	if len(nodeGroupObjects) > 0 {
		v = nodeGroupObjects
	}
	delete(request, "NextToken")

	if v != nil && len(v.([]interface{})) > 0 {
		oldNodeGroupsMap := map[string]map[string]interface{}{}
		indexNodeGroupMap := map[string]int{}
		if oldNodeGroups, ok := d.GetOk("node_groups"); ok {
			for index, item := range oldNodeGroups.([]interface{}) {
				oldNodeGroup := item.(map[string]interface{})
				oldNodeGroupsMap[oldNodeGroup["node_group_name"].(string)] = oldNodeGroup
				indexNodeGroupMap[oldNodeGroup["node_group_name"].(string)] = index
			}
		}

		nodeGroupMaps := v.([]interface{})
		var nodeGroups []map[string]interface{}
		for _, item := range nodeGroupMaps {
			nodeGroupMap := item.(map[string]interface{})
			nodeGroup := map[string]interface{}{
				"node_group_type":      nodeGroupMap["NodeGroupType"],
				"node_group_name":      nodeGroupMap["NodeGroupName"],
				"payment_type":         nodeGroupMap["PaymentType"],
				"with_public_ip":       nodeGroupMap["WithPublicIp"],
				"graceful_shutdown":    nodeGroupMap["GracefulShutdown"],
				"spot_instance_remedy": nodeGroupMap["SpotInstanceRemedy"],
				"node_count":           nodeGroupMap["RunningNodeCount"],
			}

			if oldNodeGroup, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
				if v, ok := oldNodeGroup["subscription_config"]; ok {
					var subscriptionConfig []map[string]interface{}
					for _, item := range v.(*schema.Set).List() {
						subscriptionConfig = append(subscriptionConfig, item.(map[string]interface{}))
					}
					nodeGroup["subscription_config"] = subscriptionConfig
				}

				if v, ok := oldNodeGroup["spot_bid_prices"]; ok {
					var spotBidPrices []map[string]interface{}
					for _, item := range v.(*schema.Set).List() {
						spotBidPrices = append(spotBidPrices, item.(map[string]interface{}))
					}
					nodeGroup["spot_bid_prices"] = spotBidPrices
				}

				if v, ok := oldNodeGroup["auto_scaling_policy"]; ok && len(v.([]interface{})) == 1 {
					var scalingPolicies []map[string]interface{}
					scalingPolicy := map[string]interface{}{}
					scalingPolicyMap := v.([]interface{})[0].(map[string]interface{})
					if scalingPolicyValue, scalingPolicyExists := scalingPolicyMap["scaling_rules"]; scalingPolicyExists && len(scalingPolicyValue.([]interface{})) > 0 {
						var scalingRules []map[string]interface{}
						for _, sr := range scalingPolicyValue.([]interface{}) {
							scalingRule := map[string]interface{}{}
							scalingRuleMap := sr.(map[string]interface{})
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["rule_name"]; scalingRuleExists && scalingRuleValue.(string) != "" {
								scalingRule["rule_name"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["trigger_type"]; scalingRuleExists && scalingRuleValue.(string) != "" {
								scalingRule["trigger_type"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["activity_type"]; scalingRuleExists && scalingRuleValue.(string) != "" {
								scalingRule["activity_type"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["adjustment_type"]; scalingRuleExists && scalingRuleValue.(string) != "" {
								scalingRule["adjustment_type"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["adjustment_value"]; scalingRuleExists && scalingRuleValue.(int) != 0 {
								scalingRule["adjustment_value"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["min_adjustment_value"]; scalingRuleExists {
								scalingRule["min_adjustment_value"] = scalingRuleValue
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["time_trigger"]; scalingRuleExists && len(scalingRuleValue.([]interface{})) > 0 {
								var timeTriggers []map[string]interface{}
								for _, tt := range scalingRuleValue.([]interface{}) {
									timeTrigger := map[string]interface{}{}
									timeTriggerMap := tt.(map[string]interface{})
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["launch_time"]; timeTriggerExists && timeTriggerValue.(string) != "" {
										timeTrigger["launch_time"] = timeTriggerValue
									}
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["start_time"]; timeTriggerExists && timeTriggerValue.(string) != "" {
										timeTrigger["start_time"] = timeTriggerValue
									}
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["end_time"]; timeTriggerExists && timeTriggerValue.(string) != "" {
										timeTrigger["end_time"] = timeTriggerValue
									}
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["launch_expiration_time"]; timeTriggerExists && timeTriggerValue.(int) != 0 {
										timeTrigger["launch_expiration_time"] = timeTriggerValue
									}
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["recurrence_type"]; timeTriggerExists && timeTriggerValue.(string) != "" {
										timeTrigger["recurrence_type"] = timeTriggerValue
									}
									if timeTriggerValue, timeTriggerExists := timeTriggerMap["recurrence_value"]; timeTriggerExists && timeTriggerValue.(string) != "" {
										timeTrigger["recurrence_value"] = timeTriggerValue
									}
									timeTriggers = append(timeTriggers, timeTrigger)
								}
								scalingRule["time_trigger"] = timeTriggers
							}
							if scalingRuleValue, scalingRuleExists := scalingRuleMap["metrics_trigger"]; scalingRuleExists && len(scalingRuleValue.([]interface{})) > 0 {
								var metricsTriggers []map[string]interface{}
								for _, mt := range scalingRuleValue.([]interface{}) {
									metricsTrigger := map[string]interface{}{}
									metricsTriggerMap := mt.(map[string]interface{})
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["time_window"]; metricsTriggerExists && metricsTriggerValue.(int) != 0 {
										metricsTrigger["time_window"] = metricsTriggerValue
									}
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["evaluation_count"]; metricsTriggerExists && metricsTriggerValue.(int) != 0 {
										metricsTrigger["evaluation_count"] = metricsTriggerValue
									}
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["cool_down_interval"]; metricsTriggerExists && metricsTriggerValue.(int) != 0 {
										metricsTrigger["cool_down_interval"] = metricsTriggerValue
									}
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["condition_logic_operator"]; metricsTriggerExists && metricsTriggerValue.(string) != "" {
										metricsTrigger["condition_logic_operator"] = metricsTriggerValue
									}
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["time_constraints"]; metricsTriggerExists && len(metricsTriggerValue.([]interface{})) > 0 {
										var timeConstraints []map[string]interface{}
										for _, tc := range metricsTriggerValue.([]interface{}) {
											timeConstraint := map[string]interface{}{}
											timeConstraintMap := tc.(map[string]interface{})
											if timeConstraintValue, timeConstraintExists := timeConstraintMap["start_time"]; timeConstraintExists && timeConstraintValue.(string) != "" {
												timeConstraint["start_time"] = timeConstraintValue
											}
											if timeConstraintValue, timeConstraintExists := timeConstraintMap["end_time"]; timeConstraintExists && timeConstraintValue.(string) != "" {
												timeConstraint["end_time"] = timeConstraintValue
											}
											timeConstraints = append(timeConstraints, timeConstraint)
										}
										metricsTrigger["time_constraints"] = timeConstraints
									}
									if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["conditions"]; metricsTriggerExists && len(metricsTriggerValue.([]interface{})) > 0 {
										var conditions []map[string]interface{}
										for _, cd := range metricsTriggerValue.([]interface{}) {
											condition := map[string]interface{}{}
											conditionMap := cd.(map[string]interface{})
											if conditionValue, conditionExists := conditionMap["metric_name"]; conditionExists && conditionValue.(string) != "" {
												condition["metric_name"] = conditionValue
											}
											if conditionValue, conditionExists := conditionMap["statistics"]; conditionExists && conditionValue.(string) != "" {
												condition["statistics"] = conditionValue
											}
											if conditionValue, conditionExists := conditionMap["comparison_operator"]; conditionExists && conditionValue.(string) != "" {
												condition["comparison_operator"] = conditionValue
											}
											if conditionValue, conditionExists := conditionMap["threshold"]; conditionExists && conditionValue.(float64) != 0.0 {
												condition["threshold"] = conditionValue
											}
											if conditionValue, conditionExists := conditionMap["tags"]; conditionExists && len(conditionValue.([]interface{})) > 0 {
												var tags []map[string]interface{}
												for _, tg := range conditionValue.([]interface{}) {
													tag := map[string]interface{}{}
													tagMap := tg.(map[string]interface{})
													if tagValue, tagExists := tagMap["key"]; tagExists && tagValue.(string) != "" {
														tag["key"] = tagValue
													}
													if tagValue, tagExists := tagMap["value"]; tagExists && tagValue.(string) != "" {
														tag["value"] = tagValue
													}
													tags = append(tags, tag)
												}
												condition["tags"] = tags
											}
											conditions = append(conditions, condition)
										}
										metricsTrigger["conditions"] = conditions
									}
									metricsTriggers = append(metricsTriggers, metricsTrigger)
								}
								scalingRule["metrics_trigger"] = metricsTriggers
							}
							scalingRules = append(scalingRules, scalingRule)
						}
						scalingPolicy["scaling_rules"] = scalingRules
					}
					if scalingPolicyValue, scalingPolicyExists := scalingPolicyMap["constraints"]; scalingPolicyExists && len(scalingPolicyValue.([]interface{})) > 0 {
						var constraints []map[string]interface{}
						for _, ct := range scalingPolicyValue.([]interface{}) {
							constraint := map[string]interface{}{}
							constraintMap := ct.(map[string]interface{})
							if constraintValue, constraintExists := constraintMap["max_capacity"]; constraintExists && constraintValue.(int) != 0 {
								constraint["max_capacity"] = constraintValue
							}
							if constraintValue, constraintExists := constraintMap["min_capacity"]; constraintExists && constraintValue.(int) != 0 {
								constraint["min_capacity"] = constraintValue
							}
							constraints = append(constraints, constraint)
						}
						scalingPolicy["constraints"] = constraints
					}
					scalingPolicies = append(scalingPolicies, scalingPolicy)
					nodeGroup["auto_scaling_policy"] = scalingPolicies
				}
			}

			if v, ok := nodeGroupMap["VSwitchIds"]; ok && len(v.([]interface{})) > 0 {
				var vSwitchIDs []string
				for _, item := range v.([]interface{}) {
					vSwitchIDs = append(vSwitchIDs, item.(string))
				}
				nodeGroup["vswitch_ids"] = vSwitchIDs
			}

			if v, ok := nodeGroupMap["AdditionalSecurityGroupIds"]; ok && len(v.([]interface{})) > 0 {
				var additionalSecurityGroupIDs []string
				for _, item := range v.([]interface{}) {
					additionalSecurityGroupIDs = append(additionalSecurityGroupIDs, item.(string))
				}
				nodeGroup["additional_security_group_ids"] = additionalSecurityGroupIDs
			}

			if v, ok := nodeGroupMap["InstanceTypes"]; ok && len(v.([]interface{})) > 0 {
				var instanceTypes []string
				for _, item := range v.([]interface{}) {
					instanceTypes = append(instanceTypes, item.(string))
				}
				nodeGroup["instance_types"] = instanceTypes
			}

			if v, ok := nodeGroupMap["AckConfig"]; ok && len(v.(map[string]interface{})) > 0 {
				var ackConfigs []map[string]interface{}
				ackConfig := map[string]interface{}{}
				m := v.(map[string]interface{})
				ackConfig["ack_instance_id"] = m["AckInstanceId"]
				if value, exists := m["NodeSelectors"]; exists && len(value.([]interface{})) > 0 {
					var nodeSelectors []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						nodeSelectors = append(nodeSelectors, map[string]interface{}{
							"key":   ackConfigValue.(map[string]interface{})["Key"],
							"value": ackConfigValue.(map[string]interface{})["Value"],
						})
					}
					ackConfig["node_selectors"] = nodeSelectors
				}
				if value, exists := m["Tolerations"]; exists && len(value.([]interface{})) > 0 {
					var tolerations []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						tolerations = append(tolerations, map[string]interface{}{
							"key":      ackConfigValue.(map[string]interface{})["Key"],
							"value":    ackConfigValue.(map[string]interface{})["Value"],
							"operator": ackConfigValue.(map[string]interface{})["Operator"],
							"effect":   ackConfigValue.(map[string]interface{})["Effect"],
						})
					}
					ackConfig["tolerations"] = tolerations
				}
				if value, exists := m["Namespace"]; exists && value.(string) != "" {
					ackConfig["namespace"] = value
				}
				if value, exists := m["RequestCpu"]; exists {
					ackConfig["request_cpu"] = value
				}
				if value, exists := m["RequestMemory"]; exists {
					ackConfig["request_memory"] = value
				}
				if value, exists := m["LimitCpu"]; exists {
					ackConfig["limit_cpu"] = value
				}
				if value, exists := m["LimitMemory"]; exists {
					ackConfig["limit_memory"] = value
				}
				if value, exists := m["CustomLabels"]; exists && len(value.([]interface{})) > 0 {
					var customLabels []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						customLabels = append(customLabels, map[string]interface{}{
							"key":   ackConfigValue.(map[string]interface{})["Key"],
							"value": ackConfigValue.(map[string]interface{})["Value"],
						})
					}
					ackConfig["custom_labels"] = customLabels
				}
				if value, exists := m["CustomAnnotations"]; exists && len(value.([]interface{})) > 0 {
					var customAnnotations []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						customAnnotations = append(customAnnotations, map[string]interface{}{
							"key":   ackConfigValue.(map[string]interface{})["Key"],
							"value": ackConfigValue.(map[string]interface{})["Value"],
						})
					}
					ackConfig["custom_annotations"] = customAnnotations
				}
				if value, exists := m["Pvcs"]; exists && len(value.([]interface{})) > 0 {
					var pvcs []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						pvcs = append(pvcs, map[string]interface{}{
							"name":                    ackConfigValue.(map[string]interface{})["Name"],
							"path":                    ackConfigValue.(map[string]interface{})["Path"],
							"data_disk_storage_class": ackConfigValue.(map[string]interface{})["DataDiskStorageClass"],
							"data_disk_size":          ackConfigValue.(map[string]interface{})["DataDiskSize"],
						})
					}
					ackConfig["pvcs"] = pvcs
				}
				if value, exists := m["Volumes"]; exists && len(value.([]interface{})) > 0 {
					var volumes []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						volumes = append(volumes, map[string]interface{}{
							"name": ackConfigValue.(map[string]interface{})["Name"],
							"path": ackConfigValue.(map[string]interface{})["Path"],
							"type": ackConfigValue.(map[string]interface{})["Type"],
						})
					}
					ackConfig["volumes"] = volumes
				}
				if value, exists := m["VolumeMounts"]; exists && len(value.([]interface{})) > 0 {
					var volumeMounts []map[string]interface{}
					for _, ackConfigValue := range value.([]interface{}) {
						volumeMounts = append(volumeMounts, map[string]interface{}{
							"name": ackConfigValue.(map[string]interface{})["Name"],
							"path": ackConfigValue.(map[string]interface{})["Path"],
						})
					}
					ackConfig["volume_mounts"] = volumeMounts
				}
				if value, exists := m["PreStartCommand"]; exists && len(value.([]interface{})) > 0 {
					var preStartCommands []string
					for _, ackConfigValue := range value.([]interface{}) {
						preStartCommands = append(preStartCommands, ackConfigValue.(string))
					}
					ackConfig["pre_start_command"] = preStartCommands
				}
				if value, exists := m["PodAffinity"]; exists && value.(string) != "" {
					ackConfig["pod_affinity"] = value
				}
				if value, exists := m["PodAntiAffinity"]; exists && value.(string) != "" {
					ackConfig["pod_anti_affinity"] = value
				}
				if value, exists := m["NodeAffinity"]; exists && value.(string) != "" {
					ackConfig["node_affinity"] = value
				}
				ackConfigs = append(ackConfigs, ackConfig)
				nodeGroup["ack_config"] = ackConfigs

				if ong, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
					if instValue, instExists := ong["instance_types"]; instExists {
						var instanceTypes []string
						for _, item := range instValue.(*schema.Set).List() {
							instanceTypes = append(instanceTypes, item.(string))
						}
						nodeGroup["instance_types"] = instanceTypes
					}
				}
			}

			if oldNodeGroup, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
				if v, ok := oldNodeGroup["deployment_set_strategy"]; ok && v.(string) != "" {
					deploymentSetStrategy := v.(string)
					if st, stExists := nodeGroupMap["DeploymentSetStrategy"]; stExists && st.(string) != "" {
						deploymentSetStrategy = st.(string)
					}
					nodeGroup["deployment_set_strategy"] = deploymentSetStrategy
				}
			}

			if oldNodeGroup, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
				if v, ok := oldNodeGroup["node_resize_strategy"]; ok && v.(string) != "" {
					nodeGroup["node_resize_strategy"] = v.(string)
				}

				if v, ok := oldNodeGroup["spot_strategy"]; ok && v.(string) != "" {
					nodeGroup["spot_strategy"] = v.(string)
				}

				if v, ok := oldNodeGroup["cost_optimized_config"]; ok && len(v.(*schema.Set).List()) > 0 {
					var costOptimizedConfig []map[string]interface{}
					for _, coc := range v.(*schema.Set).List() {
						costOptimizedConfig = append(costOptimizedConfig, coc.(map[string]interface{}))
					}
					nodeGroup["cost_optimized_config"] = costOptimizedConfig
				}
			}

			if v, ok := nodeGroupMap["SystemDisk"]; ok {
				systemDiskMap := v.(map[string]interface{})
				systemDisk := map[string]interface{}{
					"category": systemDiskMap["Category"],
					"size":     formatInt(systemDiskMap["Size"]),
				}
				if oldNodeGroup, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
					oldSystemDisk := oldNodeGroup["system_disk"].(*schema.Set).List()[0].(map[string]interface{})
					if v, exists := oldSystemDisk["performance_level"]; exists && v != nil && v.(string) != "" {
						systemDisk["performance_level"] = v
					}
					if v, exists := oldSystemDisk["count"]; exists && v != nil {
						if count := formatInt(v); count > 0 {
							systemDisk["count"] = count
						}
					}
				}
				nodeGroup["system_disk"] = []map[string]interface{}{systemDisk}
			}

			if v, ok := nodeGroupMap["DataDisks"]; ok && len(v.([]interface{})) > 0 {
				var dataDisks []map[string]interface{}
				for _, item := range v.([]interface{}) {
					dataDisksMap := item.(map[string]interface{})
					dataDisk := map[string]interface{}{
						"category": dataDisksMap["Category"],
						"size":     formatInt(dataDisksMap["Size"]),
						"count":    formatInt(dataDisksMap["Count"]),
					}
					if oldNodeGroup, exists := oldNodeGroupsMap[nodeGroupMap["NodeGroupName"].(string)]; exists {
						oldDataDisk := oldNodeGroup["data_disks"].(*schema.Set).List()[0].(map[string]interface{})
						if v, exists := oldDataDisk["performance_level"]; exists && v != nil && v.(string) != "" {
							dataDisk["performance_level"] = v
						}
						if v, exists := oldDataDisk["count"]; exists && v != nil && v.(int) != 0 {
							if count := formatInt(v); count > 0 {
								dataDisk["count"] = formatInt(v)
							}
						}
					}
					dataDisks = append(dataDisks, dataDisk)
				}
				nodeGroup["data_disks"] = dataDisks
			}

			nodeGroups = append(nodeGroups, nodeGroup)
		}
		sort.Slice(nodeGroups, func(i, j int) bool {
			return indexNodeGroupMap[nodeGroups[i]["node_group_name"].(string)] < indexNodeGroupMap[nodeGroups[j]["node_group_name"].(string)]
		})

		d.Set("node_groups", nodeGroups)
	}

	tags, err := emrService.ListTagResourcesNew(d.Id(), string(TagResourceCluster))

	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(tags))

	return nil
}

func resourceAlicloudEmrV2ClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	d.Partial(true)
	emrService := EmrService{client}
	if err := emrService.SetEmrClusterTagsNew(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("cluster_name") || d.HasChange("log_collect_strategy") || d.HasChange("deletion_protection") {
		action := "UpdateClusterAttribute"
		request := map[string]interface{}{
			"ClusterId": d.Id(),
			"RegionId":  client.RegionId,
		}
		if d.HasChange("cluster_name") && d.Get("cluster_name").(string) != "" {
			request["ClusterName"] = d.Get("cluster_name")
		}
		if d.HasChange("log_collect_strategy") && d.Get("log_collect_strategy").(string) != "" {
			request["LogCollectStrategy"] = d.Get("log_collect_strategy")
		}
		if d.HasChange("deletion_protection") {
			request["DeletionProtection"] = d.Get("deletion_protection")
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, false)
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
		d.SetPartial("cluster_name")
	}

	if d.HasChange("payment_type") {
		_, newPaymentType := d.GetChange("payment_type")
		if "PayAsYouGo" == newPaymentType.(string) {
			return WrapError(Error("EMR cluster can only change paymentType from PayAsYouGo to Subscription."))
		}
		if !d.HasChange("node_groups") {
			return WrapError(Error("Subscription paymentType of emr cluster can not contains PayAsYouGo node group with 'MASTER' or 'CORE'."))
		}
	}

	if d.HasChange("node_groups") {
		oldNodeGroupsList, newNodeGroupsList := d.GetChange("node_groups")

		listNodeGroupsRequest := map[string]interface{}{
			"ClusterId":  d.Id(),
			"RegionId":   client.RegionId,
			"MaxResults": PageSizeLarge,
		}
		action := "ListNodeGroups"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Emr", "2021-03-20", action, nil, listNodeGroupsRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, listNodeGroupsRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NodeGroups", response)

		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.NodeGroups", response)
		}

		oldNodeGroupMap := map[string]map[string]interface{}{}
		for _, nodeGroupItem := range resp.([]interface{}) {
			oldNodeGroup := nodeGroupItem.(map[string]interface{})
			oldNodeGroupMap[oldNodeGroup["NodeGroupName"].(string)] = oldNodeGroup
		}

		originNodeGroupMap := map[string]map[string]interface{}{}
		for _, nodeGroupItem := range oldNodeGroupsList.([]interface{}) {
			originNodeGroup := nodeGroupItem.(map[string]interface{})
			originNodeGroupMap[originNodeGroup["node_group_name"].(string)] = originNodeGroup
		}

		newNodeGroupMap := map[string]map[string]interface{}{}
		for _, nodeGroupItem := range newNodeGroupsList.([]interface{}) {
			newNodeGroup := nodeGroupItem.(map[string]interface{})
			newNodeGroupMap[newNodeGroup["node_group_name"].(string)] = newNodeGroup
		}

		isUpdateClusterPaymentType := false
		if d.HasChange("payment_type") {
			action = "UpdateClusterPaymentType"
			updateClusterPaymentTypeRequest := map[string]interface{}{
				"ClusterId": d.Id(),
				"RegionId":  client.RegionId,
			}
			if sc, ok := d.GetOk("subscription_config"); ok {
				if autoPay, exists := sc.(*schema.Set).List()[0].(map[string]interface{})["auto_pay_order"]; exists {
					updateClusterPaymentTypeRequest["AutoPayOrder"] = autoPay.(bool)
				} else {
					updateClusterPaymentTypeRequest["AutoPayOrder"] = true
				}
				if autoRenew, exists := sc.(*schema.Set).List()[0].(map[string]interface{})["auto_renew"]; exists {
					updateClusterPaymentTypeRequest["AutoRenew"] = autoRenew.(bool)
				} else {
					updateClusterPaymentTypeRequest["AutoRenew"] = false
				}
			}
			var convertNodeGroups []map[string]interface{}
			for originNodeGroupName := range originNodeGroupMap {
				convertNodeGroup := map[string]interface{}{
					"NodeGroupId": oldNodeGroupMap[originNodeGroupName]["NodeGroupId"],
				}
				if newNodeGroup, ok := newNodeGroupMap[originNodeGroupName]; ok {
					if newNodeGroupValue, newNodeGroupExists := newNodeGroup["payment_type"]; newNodeGroupExists && "Subscription" == newNodeGroupValue {
						convertNodeGroup["PaymentType"] = newNodeGroupValue
					} else if "PayAsYouGo" == newNodeGroupValue && ("MASTER" == newNodeGroup["node_group_type"] || "CORE" == newNodeGroup["node_group_type"]) {
						return WrapError(Error("Subscription paymentType of emr cluster can not contains PayAsYouGo node group with 'MASTER' or 'CORE'."))
					} else {
						continue
					}
					if subscriptionConfig, exists := newNodeGroup["subscription_config"]; exists && len(subscriptionConfig.(*schema.Set).List()) > 0 {
						subscriptionConfigMap := subscriptionConfig.(*schema.Set).List()[0].(map[string]interface{})
						if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["payment_duration"]; subscriptionConfigExists {
							convertNodeGroup["PaymentDuration"] = subscriptionConfigValue
						}
						if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["payment_duration_unit"]; subscriptionConfigExists {
							convertNodeGroup["PaymentDurationUnit"] = subscriptionConfigValue
						}
					} else {
						return WrapError(Error("The '%s' nodeGroup: '%s' is needed parameter 'subscription_config' for changing paymentType.",
							newNodeGroup["node_group_type"], newNodeGroup["node_group_name"]))
					}
					convertNodeGroups = append(convertNodeGroups, convertNodeGroup)
				}
			}
			updateClusterPaymentTypeRequest["NodeGroups"] = convertNodeGroups
			wait = incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Emr", "2021-03-20", action, nil, updateClusterPaymentTypeRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, updateClusterPaymentTypeRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// Wait for cluster payment type has been changed
			if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				if cluster, err := emrService.GetEmrV2Cluster(d.Id()); err != nil {
					return resource.NonRetryableError(err)
				} else if cluster["PaymentType"].(string) == "Subscription" {
					return nil
				}
				return resource.RetryableError(Error("Waiting for cluster %s payment type to be changed.", d.Id()))
			}); err != nil {
				return WrapError(err)
			}
			isUpdateClusterPaymentType = true
		}

		var increaseNodesGroups []map[string]interface{}
		var decreaseNodesGroups []map[string]interface{}

		for nodeGroupName, newNodeGroup := range newNodeGroupMap {
			if oldNodeGroup, ok := oldNodeGroupMap[nodeGroupName]; ok {
				if !reflect.DeepEqual(originNodeGroupMap[nodeGroupName]["auto_scaling_policy"], newNodeGroup["auto_scaling_policy"]) {
					removeScalingPolicy := false
					if len(newNodeGroup["auto_scaling_policy"].([]interface{})) > 0 {
						adaptedScalingPolicy := adaptAutoScalingPolicyRequest(newNodeGroup["auto_scaling_policy"].([]interface{})[0].(map[string]interface{}))
						rulesExists := false
						constraintsExists := false
						if aspValue, aspExists := adaptedScalingPolicy["scalingRules"]; aspExists {
							rulesExists = aspExists
							adaptedScalingPolicy["ScalingRules"] = aspValue
							delete(adaptedScalingPolicy, "scalingRules")
						}
						if aspValue, aspExists := adaptedScalingPolicy["constraints"]; aspExists {
							constraintsExists = aspExists
							adaptedScalingPolicy["Constraints"] = aspValue
							delete(adaptedScalingPolicy, "constraints")
						}
						if rulesExists || constraintsExists {
							adaptedScalingPolicy["RegionId"] = client.RegionId
							adaptedScalingPolicy["ClusterId"] = d.Id()
							adaptedScalingPolicy["NodeGroupId"] = oldNodeGroup["NodeGroupId"]

							action = "PutAutoScalingPolicy"
							err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
								response, err = client.RpcPost("Emr", "2021-03-20", action, nil, adaptedScalingPolicy, false)
								if err != nil {
									if NeedRetry(err) {
										wait()
										return resource.RetryableError(err)
									}
									return resource.NonRetryableError(err)
								}
								return nil
							})
							addDebug(action, response, adaptedScalingPolicy)
							if err != nil {
								return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
							}
						} else {
							removeScalingPolicy = true
						}
					} else {
						removeScalingPolicy = true
					}
					if removeScalingPolicy {
						removeScalingPolicyRequest := map[string]interface{}{
							"RegionId":    client.RegionId,
							"ClusterId":   d.Id(),
							"NodeGroupId": oldNodeGroup["NodeGroupId"],
						}
						action = "RemoveAutoScalingPolicy"
						err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
							response, err = client.RpcPost("Emr", "2021-03-20", action, nil, removeScalingPolicyRequest, false)
							if err != nil {
								if NeedRetry(err) {
									wait()
									return resource.RetryableError(err)
								}
								return resource.NonRetryableError(err)
							}
							return nil
						})
						addDebug(action, response, removeScalingPolicyRequest)
						if err != nil {
							return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
						}
					}
				}
				if !reflect.DeepEqual(originNodeGroupMap[nodeGroupName]["ack_config"], newNodeGroup["ack_config"]) && "K8S" == oldNodeGroup["IaaSType"] {
					ackConfigs := newNodeGroup["ack_config"].([]interface{})
					if len(ackConfigs) > 0 {
						updateNodeGroupAttributesRequest := map[string]interface{}{
							"RegionId":    client.RegionId,
							"ClusterId":   d.Id(),
							"NodeGroupId": oldNodeGroup["NodeGroupId"],
							"AckConfig":   adaptAckConfigRequest(ackConfigs[0].(map[string]interface{})),
						}
						action = "UpdateNodeGroupAttributes"
						err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
							response, err = client.RpcPost("Emr", "2021-03-20", action, nil, updateNodeGroupAttributesRequest, false)
							if err != nil {
								if NeedRetry(err) {
									wait()
									return resource.RetryableError(err)
								}
								return resource.NonRetryableError(err)
							}
							return nil
						})
						addDebug(action, response, updateNodeGroupAttributesRequest)
						if err != nil {
							return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
						}
					}
				}

				newNodeCount := formatInt(newNodeGroup["node_count"])
				oldNodeCount := formatInt(oldNodeGroup["RunningNodeCount"])

				// increase nodes
				if oldNodeCount < newNodeCount {
					if "MASTER" == newNodeGroup["node_group_type"].(string) {
						return WrapError(Error("EMR cluster can not increase the node group type of ['MASTER']."))
					}
					count := newNodeCount - oldNodeCount
					increaseNodesGroup := map[string]interface{}{}
					increaseNodesGroup["RegionId"] = client.RegionId
					increaseNodesGroup["ClusterId"] = d.Id()
					increaseNodesGroup["NodeGroupId"] = oldNodeGroup["NodeGroupId"]
					increaseNodesGroup["IncreaseNodeCount"] = count
					increaseNodesGroup["AutoRenew"] = false
					if "Subscription" == newNodeGroup["payment_type"].(string) {
						subscriptionConfig := newNodeGroup["subscription_config"].(*schema.Set).List()
						if len(subscriptionConfig) == 1 {
							configMap := subscriptionConfig[0].(map[string]interface{})
							increaseNodesGroup["PaymentDuration"] = configMap["payment_duration"]
							increaseNodesGroup["PaymentDurationUnit"] = configMap["payment_duration_unit"]
							if value, exists := configMap["auto_pay_order"]; exists {
								increaseNodesGroup["AutoPayOrder"] = value.(bool)
							} else {
								increaseNodesGroup["AutoPayOrder"] = true
							}
							if value, exists := configMap["auto_renew"]; exists {
								increaseNodesGroup["AutoRenew"] = value.(bool)
							}
						}
					}
					increaseNodesGroups = append(increaseNodesGroups, increaseNodesGroup)
				} else if oldNodeCount > newNodeCount { // decrease nodes
					// EMR cluster can only decrease 'TASK, GATEWAY' node group.
					nodeGroupType := newNodeGroup["node_group_type"].(string)
					if "TASK" != nodeGroupType && "GATEWAY" != nodeGroupType {
						return WrapError(Error("EMR cluster can only decrease the node group type of ['TASK', 'GATEWAY']."))
					}
					decreaseNodesGroup := map[string]interface{}{
						"ClusterId":         d.Id(),
						"RegionId":          client.RegionId,
						"DecreaseNodeCount": oldNodeCount - newNodeCount,
						"NodeGroupId":       oldNodeGroup["NodeGroupId"],
					}
					decreaseNodesGroups = append(decreaseNodesGroups, decreaseNodesGroup)
				}

				// increase node disk size, we can only support single disk type.
				currDataDisk := oldNodeGroup["DataDisks"].([]interface{})[0].(map[string]interface{})
				targetDataDisk := newNodeGroup["data_disks"].(*schema.Set).List()[0].(map[string]interface{})
				if formatInt(targetDataDisk["size"]) < formatInt(currDataDisk["Size"]) {
					return WrapError(Error("EMR cluster can only increase node disk, decrease node disk is not supported."))
				} else if formatInt(targetDataDisk["size"]) > formatInt(currDataDisk["Size"]) {
					if currDataDisk["Category"].(string) == "local_hdd_pro" {
						return WrapError(Error("EMR cluster can not support increase node disk with 'local_hdd_pro' disk type."))
					}
					action := "IncreaseNodesDiskSize"
					increaseNodeDiskSizeRequest := map[string]interface{}{
						"ClusterId":   d.Id(),
						"RegionId":    client.RegionId,
						"NodeGroupId": oldNodeGroup["NodeGroupId"],
						"DataDiskSizes": []map[string]interface{}{
							{
								"Category": currDataDisk["Category"].(string),
								"Size":     formatInt(targetDataDisk["size"]),
							},
						},
					}
					wait := incrementalWait(3*time.Second, 5*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
						response, err = client.RpcPost("Emr", "2021-03-20", action, nil, increaseNodeDiskSizeRequest, false)
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					addDebug(action, response, increaseNodeDiskSizeRequest)
					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}
				}
			} else { // 'Task' NodeGroupType may not exist when create emr_cluster
				subscriptionConfig := map[string]interface{}{}
				if "Subscription" == newNodeGroup["payment_type"] {
					subscriptionMap := newNodeGroup["subscription_config"].(*schema.Set).List()[0].(map[string]interface{})
					subscriptionConfig["PaymentDurationUnit"] = subscriptionMap["payment_duration_unit"]
					subscriptionConfig["PaymentDuration"] = subscriptionMap["payment_duration"]
					subscriptionConfig["AutoRenewDurationUnit"] = subscriptionMap["auto_renew_duration_unit"]
					subscriptionConfig["AutoRenewDuration"] = subscriptionMap["auto_renew_duration"]
					subscriptionConfig["AutoRenew"] = false
					if value, exists := subscriptionMap["auto_pay_order"]; exists {
						subscriptionConfig["AutoPayOrder"] = value.(bool)
					} else {
						subscriptionConfig["AutoPayOrder"] = true
					}
					if value, exists := subscriptionMap["auto_renew"]; exists {
						subscriptionConfig["AutoRenew"] = value.(bool)
					}
				}
				var spotBidPrices []map[string]interface{}
				for _, v := range newNodeGroup["spot_bid_prices"].(*schema.Set).List() {
					sbpMap := v.(map[string]interface{})
					spotBidPrices = append(spotBidPrices, map[string]interface{}{
						"InstanceType": sbpMap["instance_type"],
						"BidPrice":     sbpMap["bid_price"],
					})
				}
				systemDiskMap := newNodeGroup["system_disk"].(*schema.Set).List()[0].(map[string]interface{})
				var dataDisks []map[string]interface{}
				for _, v := range newNodeGroup["data_disks"].(*schema.Set).List() {
					dataDiskMap := v.(map[string]interface{})
					dataDisk := map[string]interface{}{
						"Category": dataDiskMap["category"],
						"Size":     dataDiskMap["size"],
						"Count":    dataDiskMap["count"],
					}
					if value, exists := dataDiskMap["performance_level"]; exists && value.(string) != "" {
						dataDisk["PerformanceLevel"] = value.(string)
					}
					dataDisks = append(dataDisks, dataDisk)
				}
				nodeGroupParam := map[string]interface{}{
					"NodeGroupType":      newNodeGroup["node_group_type"],
					"NodeGroupName":      nodeGroupName,
					"PaymentType":        newNodeGroup["payment_type"],
					"SubscriptionConfig": subscriptionConfig,
					"SpotBidPrices":      spotBidPrices,
					"WithPublicIp":       newNodeGroup["with_public_ip"],
					"NodeCount":          newNodeGroup["node_count"],
					"DataDisks":          dataDisks,
					"GracefulShutdown":   newNodeGroup["graceful_shutdown"],
					"SpotInstanceRemedy": newNodeGroup["spot_instance_remedy"],
				}
				if value, exists := newNodeGroup["auto_scaling_policy"]; exists && len(value.([]interface{})) > 0 {
					nodeGroupParam["AutoScalingPolicy"] = adaptAutoScalingPolicyRequest(value.([]interface{})[0].(map[string]interface{}))
				}
				if value, exists := newNodeGroup["deployment_set_strategy"]; exists && value.(string) != "" {
					nodeGroupParam["DeploymentSetStrategy"] = value.(string)
				}
				if value, exists := newNodeGroup["node_resize_strategy"]; exists && value.(string) != "" {
					nodeGroupParam["NodeResizeStrategy"] = value.(string)
				}
				if value, exists := newNodeGroup["spot_strategy"]; exists && value.(string) != "" {
					nodeGroupParam["SpotStrategy"] = value.(string)
				}
				vSwitchIDList := newNodeGroup["vswitch_ids"].(*schema.Set).List()
				if len(vSwitchIDList) > 0 {
					var vSwitchIDs []string
					for _, vSwitchID := range vSwitchIDList {
						vSwitchIDs = append(vSwitchIDs, vSwitchID.(string))
					}
					nodeGroupParam["VSwitchIds"] = vSwitchIDs
				}
				systemDisk := map[string]interface{}{
					"Category": systemDiskMap["category"],
					"Size":     systemDiskMap["size"],
					"Count":    systemDiskMap["count"],
				}
				if value, exists := systemDiskMap["performance_level"]; exists && value.(string) != "" {
					systemDisk["PerformanceLevel"] = value.(string)
				}
				nodeGroupParam["SystemDisk"] = systemDisk

				instanceTypeList := newNodeGroup["instance_types"].(*schema.Set).List()
				if len(instanceTypeList) > 0 {
					var instanceTypes []string
					for _, instanceType := range instanceTypeList {
						instanceTypes = append(instanceTypes, instanceType.(string))
					}
					nodeGroupParam["InstanceTypes"] = instanceTypes
				}
				if value, exists := newNodeGroup["ack_config"]; exists && len(value.([]interface{})) > 0 {
					nodeGroupParam["AckConfig"] = adaptAckConfigRequest(value.([]interface{})[0].(map[string]interface{}))
					nodeGroupParam["IaaSType"] = "K8S"
					delete(nodeGroupParam, "InstanceTypes")
				}

				addSecurityGroupIDList := newNodeGroup["additional_security_group_ids"].(*schema.Set).List()
				if len(addSecurityGroupIDList) > 0 {
					var addSecurityGroupIDs []string
					for _, addSecurityGroupID := range addSecurityGroupIDList {
						addSecurityGroupIDs = append(addSecurityGroupIDs, addSecurityGroupID.(string))
					}
					nodeGroupParam["AdditionalSecurityGroupIds"] = addSecurityGroupIDs
				}

				costOptimizedConfigList := newNodeGroup["cost_optimized_config"].(*schema.Set).List()
				if len(costOptimizedConfigList) > 0 {
					costOptimizedConfig := costOptimizedConfigList[0].(map[string]interface{})
					nodeGroupParam["CostOptimizedConfig"] = map[string]interface{}{
						"OnDemandBaseCapacity":                costOptimizedConfig["on_demand_base_capacity"],
						"OnDemandPercentageAboveBaseCapacity": costOptimizedConfig["on_demand_percentage_above_base_capacity"],
						"SpotInstancePools":                   costOptimizedConfig["spot_instance_pools"],
					}
				}
				createNodeGroupRequest := map[string]interface{}{
					"ClusterId": d.Id(),
					"RegionId":  client.RegionId,
					"NodeGroup": nodeGroupParam,
				}

				action = "CreateNodeGroup"
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Emr", "2021-03-20", action, nil, createNodeGroupRequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, createNodeGroupRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

				listNodeGroupsRequest := map[string]interface{}{
					"ClusterId":      d.Id(),
					"NodeGroupNames": []string{nodeGroupName},
					"RegionId":       client.RegionId,
				}
				action = "ListNodeGroups"
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Emr", "2021-03-20", action, nil, listNodeGroupsRequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, listNodeGroupsRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				resp, err := jsonpath.Get("$.NodeGroups", response)
				if err != nil {
					return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.NodeGroups", response)
				}

				if len(resp.([]interface{})) == 0 {
					continue
				}

				nodeGroupId := resp.([]interface{})[0].(map[string]interface{})["NodeGroupId"].(string)

				newNodeCount := formatInt(newNodeGroup["node_count"])
				if newNodeCount > 0 {
					increaseNodesGroup := map[string]interface{}{}
					increaseNodesGroup["RegionId"] = client.RegionId
					increaseNodesGroup["ClusterId"] = d.Id()
					increaseNodesGroup["NodeGroupId"] = nodeGroupId
					increaseNodesGroup["IncreaseNodeCount"] = newNodeCount
					increaseNodesGroup["AutoPayOrder"] = true
					increaseNodesGroup["AutoRenew"] = false
					if "Subscription" == newNodeGroup["payment_type"].(string) {
						subscriptionConfig := newNodeGroup["subscription_config"].(*schema.Set).List()
						if len(subscriptionConfig) == 1 {
							configMap := subscriptionConfig[0].(map[string]interface{})
							increaseNodesGroup["PaymentDuration"] = configMap["payment_duration"]
							increaseNodesGroup["PaymentDurationUnit"] = configMap["payment_duration_unit"]
							if value, exists := configMap["auto_pay_order"]; exists {
								increaseNodesGroup["AutoPayOrder"] = value.(bool)
							}
							if value, exists := configMap["auto_renew"]; exists {
								increaseNodesGroup["AutoRenew"] = value.(bool)
							}
						}
					}
					increaseNodesGroups = append(increaseNodesGroups, increaseNodesGroup)
				}
			}
		}

		var deleteNodeGroups []map[string]interface{}
		for nodeGroupName, oldNodeGroup := range oldNodeGroupMap { // Delete empty nodeGroup
			if newNodeGroup, ok := newNodeGroupMap[nodeGroupName]; !ok {
				oldNodeCount := formatInt(oldNodeGroup["RunningNodeCount"])
				if oldNodeCount > 0 {
					return WrapError(Error("The [nodeGroup: %v, nodeGroupType: %v] can not delete cause exists running nodes", nodeGroupName, oldNodeGroup["NodeGroupType"].(string)))
				}
				deleteNodeGroups = append(deleteNodeGroups, map[string]interface{}{
					"ClusterId":   d.Id(),
					"NodeGroupId": oldNodeGroup["NodeGroupId"],
					"RegionId":    client.RegionId,
				})
			} else if newNodeGroup["payment_type"] == "Subscription" && oldNodeGroup["PaymentType"] == "PayAsYouGo" &&
				!isUpdateClusterPaymentType && d.Get("payment_type") == "Subscription" {
				action = "UpdateNodeGroupPaymentType"
				UpdateNodeGroupPaymentTypeRequest := map[string]interface{}{
					"ClusterId": d.Id(),
					"RegionId":  client.RegionId,
				}
				updateNodeGroupPaymentType := map[string]interface{}{
					"NodeGroupId": oldNodeGroup["NodeGroupId"],
				}
				if newNodeGroupValue, newNodeGroupExists := newNodeGroup["payment_type"]; newNodeGroupExists {
					updateNodeGroupPaymentType["PaymentType"] = newNodeGroupValue
				}
				if subscriptionConfig, exists := newNodeGroup["subscription_config"]; exists && len(subscriptionConfig.(*schema.Set).List()) > 0 {
					subscriptionConfigMap := subscriptionConfig.(*schema.Set).List()[0].(map[string]interface{})
					if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["payment_duration"]; subscriptionConfigExists {
						updateNodeGroupPaymentType["PaymentDuration"] = subscriptionConfigValue
					}
					if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["payment_duration_unit"]; subscriptionConfigExists {
						updateNodeGroupPaymentType["PaymentDurationUnit"] = subscriptionConfigValue
					}
					if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["auto_pay_order"]; subscriptionConfigExists {
						UpdateNodeGroupPaymentTypeRequest["AutoPayOrder"] = subscriptionConfigValue
					} else {
						UpdateNodeGroupPaymentTypeRequest["AutoPayOrder"] = true
					}
					if subscriptionConfigValue, subscriptionConfigExists := subscriptionConfigMap["auto_renew"]; subscriptionConfigExists {
						UpdateNodeGroupPaymentTypeRequest["AutoRenew"] = subscriptionConfigValue
					} else {
						UpdateNodeGroupPaymentTypeRequest["AutoRenew"] = false
					}
				} else {
					return WrapError(Error("The '%s' nodeGroup: '%s' is needed parameter 'subscription_config' for changing paymentType.",
						newNodeGroup["node_group_type"], newNodeGroup["node_group_name"]))
				}
				UpdateNodeGroupPaymentTypeRequest["NodeGroup"] = updateNodeGroupPaymentType
				wait = incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Emr", "2021-03-20", action, nil, UpdateNodeGroupPaymentTypeRequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, UpdateNodeGroupPaymentTypeRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				// Wait for node group payment type has been changed
				if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
					nodeGroupId := updateNodeGroupPaymentType["NodeGroupId"].(string)
					if nodeGroups, err := emrService.ListEmrV2NodeGroups(d.Id(), []string{nodeGroupId}); err != nil {
						return resource.NonRetryableError(err)
					} else if len(nodeGroups) > 0 && "Subscription" == nodeGroups[0].(map[string]interface{})["PaymentType"].(string) {
						return nil
					}
					return resource.RetryableError(Error("Waiting for node group %s payment type to be changed.", nodeGroupId))
				}); err != nil {
					return WrapError(err)
				}
			} else if !isUpdateClusterPaymentType && d.Get("payment_type") == "Subscription" &&
				oldNodeGroup["PaymentType"] == "Subscription" && newNodeGroup["payment_type"] == "PayAsYouGo" {
				return WrapError(Error("EMR cluster can only change paymentType from PayAsYouGo to Subscription."))
			}
		}

		var wg sync.WaitGroup
		var cm sync.Map
		waitFlag := false
		for _, increaseNodesGroupRequest := range increaseNodesGroups {
			waitFlag = true
			action := "IncreaseNodes"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Emr", "2021-03-20", action, nil, increaseNodesGroupRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, increaseNodesGroupRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.OperationId", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.OperationId", response)
			}
			wg.Add(1)
			go emrService.WaitForEmrV2Operation(d.Id(), increaseNodesGroupRequest["NodeGroupId"].(string), resp.(string), 2*Timeout5Minute, &wg, &cm)
		}

		for _, decreaseNodesGroupRequest := range decreaseNodesGroups {
			waitFlag = true
			action := "DecreaseNodes"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Emr", "2021-03-20", action, nil, decreaseNodesGroupRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, decreaseNodesGroupRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.OperationId", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.OperationId", response)
			}
			wg.Add(1)
			go emrService.WaitForEmrV2Operation(d.Id(), decreaseNodesGroupRequest["NodeGroupId"].(string), resp.(string), 2*Timeout5Minute, &wg, &cm)
		}
		if waitFlag {
			wg.Wait()
			var failedNodeGroupId []string
			cm.Range(func(k, v interface{}) bool {
				if v != nil && v.(bool) == false {
					failedNodeGroupId = append(failedNodeGroupId, k.(string))
				}
				return true
			})
			if len(failedNodeGroupId) > 0 {
				return WrapError(Error("EMR cluster resize found error result with the failed nodeGroupIds: [%s].", strings.Join(failedNodeGroupId, ",")))
			}
		}

		for _, deleteNodeGroupRequest := range deleteNodeGroups {
			action := "DeleteNodeGroup"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Emr", "2021-03-20", action, nil, deleteNodeGroupRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, deleteNodeGroupRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("node_groups")
	}

	if d.HasChange("bootstrap_scripts") {
		_, newBootstrapScripts := d.GetChange("bootstrap_scripts")
		listScriptsRequest := map[string]interface{}{
			"ClusterId":  d.Id(),
			"RegionId":   client.RegionId,
			"ScriptType": "BOOTSTRAP",
		}
		action := "ListScripts"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Emr", "2021-03-20", action, nil, listScriptsRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, listScriptsRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Scripts", response)

		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Scripts", response)
		}

		if resp != nil && len(resp.([]interface{})) > 0 {
			deleteScriptRequest := map[string]interface{}{
				"ClusterId":  d.Id(),
				"RegionId":   client.RegionId,
				"ScriptType": "BOOTSTRAP",
			}
			action = "DeleteScript"
			for _, v := range resp.([]interface{}) {
				scriptMap := v.(map[string]interface{})
				deleteScriptRequest["ScriptId"] = scriptMap["ScriptId"]

				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Emr", "2021-03-20", action, nil, deleteScriptRequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, deleteScriptRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if newBootstrapScripts != nil && len(newBootstrapScripts.([]interface{})) > 0 {
			var newScripts []map[string]interface{}
			createScriptRequest := map[string]interface{}{
				"ClusterId":  d.Id(),
				"RegionId":   client.RegionId,
				"ScriptType": "BOOTSTRAP",
			}
			for _, bs := range newBootstrapScripts.([]interface{}) {
				newScriptMap := bs.(map[string]interface{})
				newScript := map[string]interface{}{
					"ScriptName":            newScriptMap["script_name"],
					"ScriptPath":            newScriptMap["script_path"],
					"ScriptArgs":            newScriptMap["script_args"],
					"ExecutionMoment":       newScriptMap["execution_moment"],
					"ExecutionFailStrategy": newScriptMap["execution_fail_strategy"],
				}
				if value, exists := newScriptMap["priority"]; exists && value.(int) > 0 {
					newScript["Priority"] = value.(int)
				}
				if value, exists := newScriptMap["node_selector"]; exists && len(value.(*schema.Set).List()) == 1 {
					nodeSelectorMap := value.(*schema.Set).List()[0].(map[string]interface{})
					nodeSelector := map[string]interface{}{
						"NodeSelectType": nodeSelectorMap["node_select_type"],
						"NodeGroupId":    nodeSelectorMap["node_group_id"],
						"NodeGroupName":  nodeSelectorMap["node_group_name"],
					}
					if v, ok := nodeSelectorMap["node_names"]; ok && len(v.([]interface{})) > 0 {
						var nodeNames []string
						for _, nn := range v.([]interface{}) {
							nodeNames = append(nodeNames, nn.(string))
						}
						nodeSelector["NodeNames"] = nodeNames
					}
					if v, ok := nodeSelectorMap["node_group_ids"]; ok && len(v.([]interface{})) > 0 {
						var nodeGroupIds []string
						for _, ngId := range v.([]interface{}) {
							nodeGroupIds = append(nodeGroupIds, ngId.(string))
						}
						nodeSelector["NodeGroupIds"] = nodeGroupIds
					}
					if v, ok := nodeSelectorMap["node_group_types"]; ok && len(v.([]interface{})) > 0 {
						var nodeGroupTypes []string
						for _, ngType := range v.([]interface{}) {
							nodeGroupTypes = append(nodeGroupTypes, ngType.(string))
						}
						nodeSelector["NodeGroupTypes"] = nodeGroupTypes
					}
					if v, ok := nodeSelectorMap["node_group_names"]; ok && len(v.([]interface{})) > 0 {
						var nodeGroupNames []string
						for _, ngName := range v.([]interface{}) {
							nodeGroupNames = append(nodeGroupNames, ngName.(string))
						}
						nodeSelector["NodeGroupNames"] = nodeGroupNames
					}
					newScript["NodeSelector"] = nodeSelector
				}
				newScripts = append(newScripts, newScript)
			}
			createScriptRequest["Scripts"] = newScripts
			action = "CreateScript"
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Emr", "2021-03-20", action, nil, createScriptRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, createScriptRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}

	d.Partial(false)

	return nil
}

func resourceAlicloudEmrV2ClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := EmrService{client}
	var response map[string]interface{}
	var err error
	if paymentType, ok := d.GetOk("payment_type"); ok && paymentType.(string) == "Subscription" {
		v, err := emrService.ListEmrV2NodeGroups(d.Id(), []string{})
		if err != nil {
			return WrapError(err)
		}
		if v != nil && len(v) > 0 {
			action := "ListNodes"
			request := map[string]interface{}{
				"ClusterId": d.Id(),
				"RegionId":  client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			for _, item := range v {
				nodeGroupMap := item.(map[string]interface{})
				if value, exists := nodeGroupMap["PaymentType"]; exists && value.(string) == "Subscription" {
					request["MaxResults"] = 100
					request["NodeGroupIds"] = []string{nodeGroupMap["NodeGroupId"].(string)}
					err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
						response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
					nodes, err := jsonpath.Get("$.Nodes", response)
					if err != nil {
						return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Nodes", response)
					}
					if nodes != nil && len(nodes.([]interface{})) > 0 {
						if err = deleteSubscriptionInstances(d, meta, nodes.([]interface{})); err != nil {
							return WrapError(err)
						}
					}
				}
			}
		}
	}

	action := "DeleteCluster"
	request := map[string]interface{}{
		"ClusterId": d.Id(),
		"RegionId":  client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || strings.Contains(err.Error(), "cluster exists nonempty pre-paid nodeGroups") {
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

	stateConf := BuildStateConf([]string{"TERMINATING"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Millisecond, emrService.EmrV2ClusterStateRefreshFunc(d.Id(), []string{"TERMINATE_FAILED"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return WrapError(emrService.WaitForEmrV2Cluster(d.Id(), Deleted, DefaultTimeoutMedium))
}

func deleteSubscriptionInstances(d *schema.ResourceData, meta interface{}, instances []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreateModifyInstanceChargeTypeRequest()
	var instanceIds []interface{}
	for _, item := range instances {
		instanceIds = append(instanceIds, item.(map[string]interface{})["NodeId"])
	}
	request.InstanceIds = convertListToJsonString(instanceIds)
	request.AutoPay = requests.NewBoolean(true)
	request.DryRun = requests.NewBoolean(false)
	request.InstanceChargeType = string(PostPaid)
	if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceChargeType(request)
		})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	for _, instanceId := range instanceIds {
		deleteRequest := ecs.CreateDeleteInstanceRequest()
		deleteRequest.InstanceId = instanceId.(string)
		deleteRequest.Force = requests.NewBoolean(true)

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DeleteInstance(deleteRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "DependencyViolation.RouteEntry", "IncorrectInstanceStatus.Initializing"}) {
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{Throttling, "LastTokenProcessing"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteRequest.GetActionName(), raw)
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, EcsNotFound) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"Pending", "Running", "Stopped", "Stopping"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

		if _, err = stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return nil
}

func adaptAutoScalingPolicyRequest(r map[string]interface{}) map[string]interface{} {
	scalingPolicy := map[string]interface{}{}
	if value, exists := r["constraints"]; exists {
		constraints := value.([]interface{})
		if len(constraints) == 1 {
			constraint := map[string]interface{}{}
			constraintMap := constraints[0].(map[string]interface{})
			if constraintValue, constraintExists := constraintMap["max_capacity"]; constraintExists {
				constraint["maxCapacity"] = constraintValue
				constraint["MaxCapacity"] = constraintValue
			}
			if constraintValue, constraintExists := constraintMap["min_capacity"]; constraintExists {
				constraint["minCapacity"] = constraintValue
				constraint["MinCapacity"] = constraintValue
			}
			scalingPolicy["constraints"] = constraint
			scalingPolicy["Constraints"] = constraint
		}
	}
	if value, exists := r["scaling_rules"]; exists {
		var scalingRules []map[string]interface{}
		for _, sr := range value.([]interface{}) {
			scalingRule := map[string]interface{}{}
			scalingRuleMap := sr.(map[string]interface{})
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["rule_name"]; scalingRuleExists {
				scalingRule["RuleName"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["trigger_type"]; scalingRuleExists {
				scalingRule["TriggerType"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["activity_type"]; scalingRuleExists {
				scalingRule["ActivityType"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["adjustment_type"]; scalingRuleExists {
				scalingRule["AdjustmentType"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["adjustment_value"]; scalingRuleExists {
				scalingRule["AdjustmentValue"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["min_adjustment_value"]; scalingRuleExists {
				scalingRule["MinAdjustmentValue"] = scalingRuleValue
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["time_trigger"]; scalingRuleExists {
				timeTriggers := scalingRuleValue.([]interface{})
				if len(timeTriggers) == 1 {
					timeTrigger := map[string]interface{}{}
					timeTriggerMap := timeTriggers[0].(map[string]interface{})
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["launch_time"]; timeTriggerExists {
						timeTrigger["LaunchTime"] = timeTriggerValue
					}
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["start_time"]; timeTriggerExists {
						timeTrigger["StartTime"] = timeTriggerValue
					}
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["end_time"]; timeTriggerExists {
						timeTrigger["EndTime"] = timeTriggerValue
					}
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["launch_expiration_time"]; timeTriggerExists {
						timeTrigger["LaunchExpirationTime"] = timeTriggerValue
					}
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["recurrence_type"]; timeTriggerExists {
						timeTrigger["RecurrenceType"] = timeTriggerValue
					}
					if timeTriggerValue, timeTriggerExists := timeTriggerMap["recurrence_value"]; timeTriggerExists {
						timeTrigger["RecurrenceValue"] = timeTriggerValue
					}
					scalingRule["TimeTrigger"] = timeTrigger
				}
			}
			if scalingRuleValue, scalingRuleExists := scalingRuleMap["metrics_trigger"]; scalingRuleExists {
				metricsTriggers := scalingRuleValue.([]interface{})
				if len(metricsTriggers) == 1 {
					metricsTrigger := map[string]interface{}{}
					metricsTriggerMap := metricsTriggers[0].(map[string]interface{})
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["time_window"]; metricsTriggerExists {
						metricsTrigger["TimeWindow"] = metricsTriggerValue
					}
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["evaluation_count"]; metricsTriggerExists {
						metricsTrigger["EvaluationCount"] = metricsTriggerValue
					}
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["cool_down_interval"]; metricsTriggerExists {
						metricsTrigger["CoolDownInterval"] = metricsTriggerValue
					}
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["condition_logic_operator"]; metricsTriggerExists {
						metricsTrigger["ConditionLogicOperator"] = metricsTriggerValue
					}
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["time_constraints"]; metricsTriggerExists {
						var timeConstraints []map[string]interface{}
						for _, tc := range metricsTriggerValue.([]interface{}) {
							timeConstraint := map[string]interface{}{}
							timeConstraintMap := tc.(map[string]interface{})
							if timeConstraintValue, timeConstraintExists := timeConstraintMap["start_time"]; timeConstraintExists {
								timeConstraint["StartTime"] = timeConstraintValue
							}
							if timeConstraintValue, timeConstraintExists := timeConstraintMap["end_time"]; timeConstraintExists {
								timeConstraint["EndTime"] = timeConstraintValue
							}
							timeConstraints = append(timeConstraints, timeConstraint)
						}
						metricsTrigger["TimeConstraints"] = timeConstraints
					}
					if metricsTriggerValue, metricsTriggerExists := metricsTriggerMap["conditions"]; metricsTriggerExists {
						var conditions []map[string]interface{}
						for _, cd := range metricsTriggerValue.([]interface{}) {
							condition := map[string]interface{}{}
							conditionMap := cd.(map[string]interface{})
							if conditionValue, conditionExists := conditionMap["metric_name"]; conditionExists {
								condition["MetricName"] = conditionValue
							}
							if conditionValue, conditionExists := conditionMap["statistics"]; conditionExists {
								condition["Statistics"] = conditionValue
							}
							if conditionValue, conditionExists := conditionMap["comparison_operator"]; conditionExists {
								condition["ComparisonOperator"] = conditionValue
							}
							if conditionValue, conditionExists := conditionMap["threshold"]; conditionExists {
								condition["Threshold"] = conditionValue
							}
							if conditionValue, conditionExists := conditionMap["tags"]; conditionExists {
								var tags []map[string]interface{}
								for _, t := range conditionValue.([]interface{}) {
									tag := map[string]interface{}{}
									tagMap := t.(map[string]interface{})
									if tagValue, tagExists := tagMap["key"]; tagExists {
										tag["Key"] = tagValue
									}
									if tagValue, tagExists := tagMap["value"]; tagExists {
										tag["Value"] = tagValue
									}
									tags = append(tags, tag)
								}
								condition["Tags"] = tags
							}
							conditions = append(conditions, condition)
						}
						metricsTrigger["Conditions"] = conditions
					}
					scalingRule["MetricsTrigger"] = metricsTrigger
				}
			}
			scalingRules = append(scalingRules, scalingRule)
		}
		scalingPolicy["scalingRules"] = scalingRules
		scalingPolicy["ScalingRules"] = scalingRules
	}
	return scalingPolicy
}

func adaptAckConfigRequest(r map[string]interface{}) map[string]interface{} {
	ackConfig := map[string]interface{}{}
	if value, exists := r["ack_instance_id"]; exists {
		ackConfig["AckInstanceId"] = value
	}
	if value, exists := r["node_selectors"]; exists {
		nodeSelectors := value.(*schema.Set).List()
		if len(nodeSelectors) > 0 {
			var nodeSelectorsReq []map[string]interface{}
			for _, ns := range nodeSelectors {
				nodeSelectorsReq = append(nodeSelectorsReq, map[string]interface{}{
					"Key":   ns.(map[string]interface{})["key"],
					"Value": ns.(map[string]interface{})["value"],
				})
			}
			ackConfig["NodeSelectors"] = nodeSelectorsReq
		}
	}
	if value, exists := r["tolerations"]; exists {
		tolerations := value.([]interface{})
		if len(tolerations) > 0 {
			var tolerationsReq []map[string]interface{}
			for _, t := range tolerations {
				toleration := t.(map[string]interface{})
				tolerationsReq = append(tolerationsReq, map[string]interface{}{
					"Key":      toleration["key"],
					"Value":    toleration["value"],
					"Operator": toleration["operator"],
					"Effect":   toleration["effect"],
				})
			}
			ackConfig["Tolerations"] = tolerationsReq
		}
	}
	if value, exists := r["namespace"]; exists {
		ackConfig["Namespace"] = value
	}
	if value, exists := r["request_cpu"]; exists {
		ackConfig["RequestCpu"] = value
	}
	if value, exists := r["request_memory"]; exists {
		ackConfig["RequestMemory"] = value
	}
	if value, exists := r["limit_cpu"]; exists {
		ackConfig["LimitCpu"] = value
	}
	if value, exists := r["limit_memory"]; exists {
		ackConfig["LimitMemory"] = value
	}
	if value, exists := r["custom_labels"]; exists {
		customLabels := value.(*schema.Set).List()
		if len(customLabels) > 0 {
			var customLabelsReq []map[string]interface{}
			for _, cl := range customLabels {
				customLabel := cl.(map[string]interface{})
				customLabelsReq = append(customLabelsReq, map[string]interface{}{
					"Key":   customLabel["key"],
					"Value": customLabel["value"],
				})
			}
			ackConfig["CustomLabels"] = customLabelsReq
		}
	}
	if value, exists := r["custom_annotations"]; exists {
		customAnnotations := value.(*schema.Set).List()
		if len(customAnnotations) > 0 {
			var customAnnotationsReq []map[string]interface{}
			for _, ca := range customAnnotations {
				customAnnotation := ca.(map[string]interface{})
				customAnnotationsReq = append(customAnnotationsReq, map[string]interface{}{
					"Key":   customAnnotation["key"],
					"Value": customAnnotation["value"],
				})
			}
			ackConfig["CustomAnnotations"] = customAnnotationsReq
		}
	}
	if value, exists := r["pvcs"]; exists {
		pvcs := value.([]interface{})
		if len(pvcs) > 0 {
			var pvcsReq []map[string]interface{}
			for _, pvc := range pvcs {
				pvcMap := pvc.(map[string]interface{})
				pvcsReq = append(pvcsReq, map[string]interface{}{
					"DataDiskStorageClass": pvcMap["data_disk_storage_class"],
					"DataDiskSize":         pvcMap["data_disk_size"],
					"Path":                 pvcMap["path"],
					"Name":                 pvcMap["name"],
				})
			}
			ackConfig["Pvcs"] = pvcsReq
		}
	}
	if value, exists := r["volumes"]; exists {
		volumes := value.([]interface{})
		if len(volumes) > 0 {
			var volumesReq []map[string]interface{}
			for _, vl := range volumes {
				volume := vl.(map[string]interface{})
				volumesReq = append(volumesReq, map[string]interface{}{
					"Name": volume["name"],
					"Path": volume["path"],
					"Type": volume["type"],
				})
			}
			ackConfig["Volumes"] = volumesReq
		}
	}
	if value, exists := r["volume_mounts"]; exists {
		volumeMounts := value.([]interface{})
		if len(volumeMounts) > 0 {
			var volumeMountsReq []map[string]interface{}
			for _, vm := range volumeMounts {
				volumeMount := vm.(map[string]interface{})
				volumeMountsReq = append(volumeMountsReq, map[string]interface{}{
					"Name": volumeMount["name"],
					"Path": volumeMount["path"],
				})
			}
			ackConfig["VolumeMounts"] = volumeMountsReq
		}
	}
	if value, exists := r["pre_start_command"]; exists {
		preStartCommand := value.([]interface{})
		if len(preStartCommand) > 0 {
			var preStartCommandReq []string
			for _, prc := range preStartCommand {
				preStartCommandReq = append(preStartCommandReq, prc.(string))
			}
			ackConfig["PreStartCommand"] = preStartCommandReq
		}
	}
	if value, exists := r["pod_affinity"]; exists {
		ackConfig["PodAffinity"] = value
	}
	if value, exists := r["pod_anti_affinity"]; exists {
		ackConfig["PodAntiAffinity"] = value
	}
	if value, exists := r["node_affinity"]; exists {
		ackConfig["NodeAffinity"] = value
	}

	return ackConfig
}
