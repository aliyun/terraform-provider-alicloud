package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
							ValidateFunc: validation.StringInSlice([]string{"MASTER", "CORE", "TASK", "GATEWAY"}, false),
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
							ValidateFunc: validation.StringInSlice([]string{"BEFORE_INSTALL", "AFTER_STARTED"}, false),
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
	action := "CreateCluster"
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
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

			createClusterRequest["SubscriptionConfig"] = subscriptionConfigMap
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

	applications := make([]map[string]interface{}, 0)
	if apps, ok := d.GetOk("applications"); ok {
		for _, application := range apps.(*schema.Set).List() {
			applications = append(applications, map[string]interface{}{"ApplicationName": application.(string)})
		}
	}
	createClusterRequest["Applications"] = applications

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
	createClusterRequest["ApplicationConfigs"] = applicationConfigs

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
			createClusterRequest["NodeAttributes"] = nodeAttributesSourceMap
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
			nodeGroups = append(nodeGroups, nodeGroup)
		}
	}
	createClusterRequest["NodeGroups"] = nodeGroups

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
		createClusterRequest["BootstrapScripts"] = bootstrapScripts
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		}
		createClusterRequest["Tags"] = tags
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, createClusterRequest, &runtime)
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

	if v, ok := d.GetOk("log_collect_strategy"); ok && v.(string) != "" {
		logCollectStrategy := v.(string)
		if value, exists := object["LogCollectStrategy"]; exists && value.(string) != "" {
			logCollectStrategy = value.(string)
		}
		d.Set("log_collect_strategy", logCollectStrategy)
	}

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
			"vpc_id":              nodeAttributesMap["VpcId"],
			"zone_id":             nodeAttributesMap["ZoneId"],
			"security_group_id":   nodeAttributesMap["SecurityGroupId"],
			"ram_role":            nodeAttributesMap["RamRole"],
			"key_pair_name":       nodeAttributesMap["KeyPairName"],
			"data_disk_encrypted": nodeAttributesMap["DataDiskEncrypted"],
		}
		if v, exists := nodeAttributesMap["DataDiskKMSKeyId"]; exists && v.(string) != "" {
			nodeAttribute["data_disk_kms_key_id"] = v
		}
		nodeAttributes = append(nodeAttributes, nodeAttribute)
		d.Set("node_attributes", nodeAttributes)
	}
	var response map[string]interface{}
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ListApplications"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, request, &runtime)
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

	action = "ListScripts"
	listScriptsRequest := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ClusterId":  d.Id(),
		"ScriptType": "BOOTSTRAP",
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, listScriptsRequest, &runtime)
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
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, request, &runtime)
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
	v, err = jsonpath.Get("$.NodeGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.NodeGroups", response)
	}
	if v != nil && len(v.([]interface{})) > 0 {
		oldNodeGroupsMap := map[string]map[string]interface{}{}
		if oldNodeGroups, ok := d.GetOk("node_groups"); ok {
			for _, item := range oldNodeGroups.([]interface{}) {
				oldNodeGroup := item.(map[string]interface{})
				oldNodeGroupsMap[oldNodeGroup["node_group_name"].(string)] = oldNodeGroup
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

			if v, ok := nodeGroupMap["SpotBidPrices"]; ok && len(v.([]interface{})) > 0 {
				var spotBidPrices []map[string]interface{}
				for _, item := range v.([]interface{}) {
					spotBidPricesMap := item.(map[string]interface{})
					spotBidPrices = append(spotBidPrices, map[string]interface{}{
						"instance_type": spotBidPricesMap["InstanceType"],
						"bid_price":     formatInt(spotBidPricesMap["BidPrice"]),
					})
				}
				nodeGroup["spot_bid_prices"] = spotBidPrices
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
					nodeResizeStrategy := v.(string)
					if nrs, nrsExists := nodeGroupMap["NodeResizeStrategy"]; nrsExists && nrs.(string) != "" {
						nodeResizeStrategy = nrs.(string)
					}
					nodeGroup["node_resize_strategy"] = nodeResizeStrategy
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

			if v, ok := nodeGroupMap["CostOptimizedConfig"]; ok {
				costOptimizedConfigMap := v.(map[string]interface{})
				nodeGroup["cost_optimized_config"] = []map[string]interface{}{
					{
						"on_demand_base_capacity":                  formatInt(costOptimizedConfigMap["OnDemandBaseCapacity"]),
						"on_demand_percentage_above_base_capacity": formatInt(costOptimizedConfigMap["OnDemandPercentageAboveBaseCapacity"]),
						"spot_instance_pools":                      formatInt(costOptimizedConfigMap["SpotInstancePools"]),
					},
				}
			}

			nodeGroups = append(nodeGroups, nodeGroup)
		}

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
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)
	emrService := EmrService{client}
	if err := emrService.SetEmrClusterTagsNew(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("cluster_name") || d.HasChange("log_collect_strategy") {
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			"ClusterId": d.Id(),
			"RegionId":  client.RegionId,
		}
		action := "ListNodeGroups"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, listNodeGroupsRequest, &runtime)
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
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, updateClusterPaymentTypeRequest, &runtime)
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
					if len(newNodeGroup["auto_scaling_policy"].([]interface{})) > 0 {
						adaptedScalingPolicy := adaptAutoScalingPolicyRequest(newNodeGroup["auto_scaling_policy"].([]interface{})[0].(map[string]interface{}))
						if aspValue, aspExists := adaptedScalingPolicy["scalingRules"]; aspExists {
							adaptedScalingPolicy["ScalingRules"] = aspValue
							delete(adaptedScalingPolicy, "scalingRules")
						}
						if aspValue, aspExists := adaptedScalingPolicy["constraints"]; aspExists {
							adaptedScalingPolicy["Constraints"] = aspValue
							delete(adaptedScalingPolicy, "constraints")
						}
						adaptedScalingPolicy["RegionId"] = client.RegionId
						adaptedScalingPolicy["ClusterId"] = d.Id()
						adaptedScalingPolicy["NodeGroupId"] = oldNodeGroup["NodeGroupId"]

						action = "PutAutoScalingPolicy"
						err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
							response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, adaptedScalingPolicy, &runtime)
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
					}
				}

				newNodeCount := formatInt(newNodeGroup["node_count"])
				oldNodeCount := formatInt(oldNodeGroup["RunningNodeCount"])

				// increase nodes
				if oldNodeCount < newNodeCount {
					count := newNodeCount - oldNodeCount
					increaseNodesGroup := map[string]interface{}{}
					increaseNodesGroup["RegionId"] = client.RegionId
					increaseNodesGroup["ClusterId"] = d.Id()
					increaseNodesGroup["NodeGroupId"] = oldNodeGroup["NodeGroupId"]
					increaseNodesGroup["IncreaseNodeCount"] = count
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
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 5*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, increaseNodeDiskSizeRequest, &runtime)
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
					subscriptionConfig["AutoRenew"] = subscriptionMap["auto_renew"]
					subscriptionConfig["AutoRenewDurationUnit"] = subscriptionMap["auto_renew_duration_unit"]
					subscriptionConfig["AutoRenewDuration"] = subscriptionMap["auto_renew_duration"]
					if value, exists := subscriptionMap["auto_pay_order"]; exists {
						subscriptionConfig["AutoPayOrder"] = value.(bool)
					} else {
						subscriptionConfig["AutoPayOrder"] = true
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, createNodeGroupRequest, &runtime)
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, listNodeGroupsRequest, &runtime)
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
					if "Subscription" == newNodeGroup["payment_type"].(string) {
						subscriptionConfig := newNodeGroup["subscription_config"].(*schema.Set).List()
						if len(subscriptionConfig) == 1 {
							configMap := subscriptionConfig[0].(map[string]interface{})
							increaseNodesGroup["PaymentDuration"] = configMap["payment_duration"]
							increaseNodesGroup["PaymentDurationUnit"] = configMap["payment_duration_unit"]
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
					return WrapError(Error(fmt.Sprintf("The [nodeGroup: %v, nodeGroupType: %v] can not delete cause exists running nodes", nodeGroupName, oldNodeGroup["NodeGroupType"].(string))))
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, UpdateNodeGroupPaymentTypeRequest, &runtime)
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

		for _, increaseNodesGroupRequest := range increaseNodesGroups {
			action := "IncreaseNodes"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"),
					StringPointer("2021-03-20"), StringPointer("AK"), nil, increaseNodesGroupRequest, &runtime)
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
		}

		for _, decreaseNodesGroupRequest := range decreaseNodesGroups {
			action := "DecreaseNodes"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"),
					StringPointer("2021-03-20"), StringPointer("AK"), nil, decreaseNodesGroupRequest, &runtime)
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
		}

		for _, deleteNodeGroupRequest := range deleteNodeGroups {
			action := "DeleteNodeGroup"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"),
					StringPointer("2021-03-20"), StringPointer("AK"), nil, deleteNodeGroupRequest, &runtime)
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, listScriptsRequest, &runtime)
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, deleteScriptRequest, &runtime)
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
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, createScriptRequest, &runtime)
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
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
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
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			for _, item := range v {
				nodeGroupMap := item.(map[string]interface{})
				if value, exists := nodeGroupMap["PaymentType"]; exists && value.(string) == "Subscription" {
					request["MaxResults"] = 100
					request["NodeGroupIds"] = []string{nodeGroupMap["NodeGroupId"].(string)}
					err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, request, &runtime)
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, request, &runtime)
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
	}
	return scalingPolicy
}
