package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEhpcClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEhpcClusterV2Create,
		Read:   resourceAliCloudEhpcClusterV2Read,
		Update: resourceAliCloudEhpcClusterV2Update,
		Delete: resourceAliCloudEhpcClusterV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(8 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"additional_packages": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"addons": {
				Type:      schema.TypeList,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"services_spec": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"resources_spec": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"name": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"client_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_credentials": {
				Type:      schema.TypeList,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
				MaxItems:  1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_pair_name": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"cluster_custom_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"script": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"args": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cluster_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ehpc_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_scale_in": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_scale_out": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"grow_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"idle_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"manager": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"manager_node": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_renew_period": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_charge_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"auto_renew": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"duration": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"system_disk": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"category": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"size": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
												},
												"level": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
									"enable_ht": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"expired_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"spot_price_limit": {
										Type:     schema.TypeFloat,
										Optional: true,
										ForceNew: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"spot_strategy": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"period_unit": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"scheduler": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"dns": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"directory_service": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"max_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_compute_load_monitor": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"queues": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enable_scale_out": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"enable_scale_in": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"min_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"max_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"initial_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"inter_connect": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"compute_nodes": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"instance_charge_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"period_unit": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"auto_renew": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"auto_renew_period": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"spot_strategy": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"spot_price_limit": {
										Type:     schema.TypeFloat,
										Optional: true,
										ForceNew: true,
									},
									"duration": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"enable_ht": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"system_disk": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"category": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"size": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
												},
												"level": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"allocation_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ram_role": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"hostname_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"hostname_suffix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"keep_alive_nodes": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_count_per_cycle": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"reserved_node_pool_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scheduler_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_topology_awareness": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"shared_storages": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_directory": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"nas_directory": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"mount_target_domain": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"mount_options": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"tags": tagsSchemaForceNew(),
		},
	}
}

func resourceAliCloudEhpcClusterV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	manager := make(map[string]interface{})

	if v := d.Get("manager"); !IsNil(v) {
		managerNode := make(map[string]interface{})
		enableHt, _ := jsonpath.Get("$[0].manager_node[0].enable_ht", d.Get("manager"))
		if enableHt != nil && enableHt != "" {
			managerNode["EnableHT"] = enableHt
		}
		imageId1, _ := jsonpath.Get("$[0].manager_node[0].image_id", d.Get("manager"))
		if imageId1 != nil && imageId1 != "" {
			managerNode["ImageId"] = imageId1
		}
		spotPriceLimit1, _ := jsonpath.Get("$[0].manager_node[0].spot_price_limit", d.Get("manager"))
		if spotPriceLimit1 != nil && spotPriceLimit1 != "" {
			managerNode["SpotPriceLimit"] = spotPriceLimit1
		}
		systemDisk := make(map[string]interface{})
		size1, _ := jsonpath.Get("$[0].manager_node[0].system_disk[0].size", d.Get("manager"))
		if size1 != nil && size1 != "" {
			systemDisk["Size"] = size1
		}
		level1, _ := jsonpath.Get("$[0].manager_node[0].system_disk[0].level", d.Get("manager"))
		if level1 != nil && level1 != "" {
			systemDisk["Level"] = level1
		}
		category1, _ := jsonpath.Get("$[0].manager_node[0].system_disk[0].category", d.Get("manager"))
		if category1 != nil && category1 != "" {
			systemDisk["Category"] = category1
		}

		if len(systemDisk) > 0 {
			managerNode["SystemDisk"] = systemDisk
		}
		periodUnit1, _ := jsonpath.Get("$[0].manager_node[0].period_unit", d.Get("manager"))
		if periodUnit1 != nil && periodUnit1 != "" {
			managerNode["PeriodUnit"] = periodUnit1
		}
		autoRenew1, _ := jsonpath.Get("$[0].manager_node[0].auto_renew", d.Get("manager"))
		if autoRenew1 != nil && autoRenew1 != "" {
			managerNode["AutoRenew"] = autoRenew1
		}
		instanceType1, _ := jsonpath.Get("$[0].manager_node[0].instance_type", d.Get("manager"))
		if instanceType1 != nil && instanceType1 != "" {
			managerNode["InstanceType"] = instanceType1
		}
		duration1, _ := jsonpath.Get("$[0].manager_node[0].duration", d.Get("manager"))
		if duration1 != nil && duration1 != "" {
			managerNode["Duration"] = duration1
		}
		spotStrategy1, _ := jsonpath.Get("$[0].manager_node[0].spot_strategy", d.Get("manager"))
		if spotStrategy1 != nil && spotStrategy1 != "" {
			managerNode["SpotStrategy"] = spotStrategy1
		}
		period1, _ := jsonpath.Get("$[0].manager_node[0].period", d.Get("manager"))
		if period1 != nil && period1 != "" {
			managerNode["Period"] = period1
		}
		autoRenewPeriod1, _ := jsonpath.Get("$[0].manager_node[0].auto_renew_period", d.Get("manager"))
		if autoRenewPeriod1 != nil && autoRenewPeriod1 != "" {
			managerNode["AutoRenewPeriod"] = autoRenewPeriod1
		}
		instanceChargeType1, _ := jsonpath.Get("$[0].manager_node[0].instance_charge_type", d.Get("manager"))
		if instanceChargeType1 != nil && instanceChargeType1 != "" {
			managerNode["InstanceChargeType"] = instanceChargeType1
		}

		if len(managerNode) > 0 {
			manager["ManagerNode"] = managerNode
		}
		directoryService := make(map[string]interface{})
		type1, _ := jsonpath.Get("$[0].directory_service[0].type", d.Get("manager"))
		if type1 != nil && type1 != "" {
			directoryService["Type"] = type1
		}
		version1, _ := jsonpath.Get("$[0].directory_service[0].version", d.Get("manager"))
		if version1 != nil && version1 != "" {
			directoryService["Version"] = version1
		}

		if len(directoryService) > 0 {
			manager["DirectoryService"] = directoryService
		}
		scheduler := make(map[string]interface{})
		type3, _ := jsonpath.Get("$[0].scheduler[0].type", d.Get("manager"))
		if type3 != nil && type3 != "" {
			scheduler["Type"] = type3
		}
		version3, _ := jsonpath.Get("$[0].scheduler[0].version", d.Get("manager"))
		if version3 != nil && version3 != "" {
			scheduler["Version"] = version3
		}

		if len(scheduler) > 0 {
			manager["Scheduler"] = scheduler
		}
		dNS := make(map[string]interface{})
		type5, _ := jsonpath.Get("$[0].dns[0].type", d.Get("manager"))
		if type5 != nil && type5 != "" {
			dNS["Type"] = type5
		}
		version5, _ := jsonpath.Get("$[0].dns[0].version", d.Get("manager"))
		if version5 != nil && version5 != "" {
			dNS["Version"] = version5
		}

		if len(dNS) > 0 {
			manager["DNS"] = dNS
		}

		managerJson, err := json.Marshal(manager)
		if err != nil {
			return WrapError(err)
		}
		request["Manager"] = string(managerJson)
	}

	if v, ok := d.GetOk("shared_storages"); ok {
		sharedStoragesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["MountTargetDomain"] = dataLoopTmp["mount_target_domain"]
			dataLoopMap["NASDirectory"] = dataLoopTmp["nas_directory"]
			dataLoopMap["MountDirectory"] = dataLoopTmp["mount_directory"]
			dataLoopMap["MountOptions"] = dataLoopTmp["mount_options"]
			dataLoopMap["ProtocolType"] = dataLoopTmp["protocol_type"]
			dataLoopMap["FileSystemId"] = dataLoopTmp["file_system_id"]
			sharedStoragesMapsArray = append(sharedStoragesMapsArray, dataLoopMap)
		}
		sharedStoragesMapsJson, err := json.Marshal(sharedStoragesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["SharedStorages"] = string(sharedStoragesMapsJson)
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		request["ClusterName"] = v
	}
	if v, ok := d.GetOk("addons"); ok {
		addonsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Name"] = dataLoop1Tmp["name"]
			dataLoop1Map["ServicesSpec"] = dataLoop1Tmp["services_spec"]
			dataLoop1Map["ResourcesSpec"] = dataLoop1Tmp["resources_spec"]
			dataLoop1Map["Version"] = dataLoop1Tmp["version"]
			addonsMapsArray = append(addonsMapsArray, dataLoop1Map)
		}
		addonsMapsJson, err := json.Marshal(addonsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Addons"] = string(addonsMapsJson)
	}

	if v, ok := d.GetOk("cluster_category"); ok {
		request["ClusterCategory"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	clusterCredentials := make(map[string]interface{})

	if v := d.Get("cluster_credentials"); v != nil {
		keyPairName1, _ := jsonpath.Get("$[0].key_pair_name", v)
		if keyPairName1 != nil && keyPairName1 != "" {
			clusterCredentials["KeyPairName"] = keyPairName1
		}
		password1, _ := jsonpath.Get("$[0].password", v)
		if password1 != nil && password1 != "" {
			clusterCredentials["Password"] = password1
		}

		clusterCredentialsJson, err := json.Marshal(clusterCredentials)
		if err != nil {
			return WrapError(err)
		}
		request["ClusterCredentials"] = string(clusterCredentialsJson)
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}
	if v, ok := d.GetOk("cluster_vswitch_id"); ok {
		request["ClusterVSwitchId"] = v
	}
	if v, ok := d.GetOk("cluster_mode"); ok {
		request["ClusterMode"] = v
	}
	if v, ok := d.GetOk("cluster_vpc_id"); ok {
		request["ClusterVpcId"] = v
	}
	if v, ok := d.GetOk("client_version"); ok {
		request["ClientVersion"] = v
	}
	if v, ok := d.GetOk("cluster_description"); ok {
		request["ClusterDescription"] = v
	}
	if v, ok := d.GetOkExists("is_enterprise_security_group"); ok {
		request["IsEnterpriseSecurityGroup"] = v
	}
	if v := d.Get("cluster_custom_configuration"); len(convertToInterfaceArray(v)) > 0 {
		clusterCustomConfiguration := make(map[string]interface{})
		script1, _ := jsonpath.Get("$[0].script", v)
		if script1 != nil && script1 != "" {
			clusterCustomConfiguration["Script"] = script1
		}
		args1, _ := jsonpath.Get("$[0].args", v)
		if args1 != nil && args1 != "" {
			clusterCustomConfiguration["Args"] = args1
		}

		if len(clusterCustomConfiguration) > 0 {
			clusterCustomConfigurationJson, err := json.Marshal(clusterCustomConfiguration)
			if err != nil {
				return WrapError(err)
			}
			request["ClusterCustomConfiguration"] = string(clusterCustomConfigurationJson)
		}
	}
	if v, ok := d.GetOkExists("max_count"); ok {
		request["MaxCount"] = v
	}
	if v, ok := d.GetOkExists("max_core_count"); ok {
		request["MaxCoreCount"] = v
	}
	if v, ok := d.GetOkExists("grow_interval"); ok {
		request["GrowInterval"] = v
	}
	if v, ok := d.GetOkExists("idle_interval"); ok {
		request["IdleInterval"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMapsArray := make([]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			// CreateCluster expects the lowercase item fields "key"/"value";
			// PascalCase fields are rejected with InvalidParams "Tags.0.key".
			tagsMapsArray = append(tagsMapsArray, map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
		tagsMapsJson, err := json.Marshal(tagsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = string(tagsMapsJson)
	}
	if v, ok := d.GetOk("additional_packages"); ok {
		additionalPackagesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			if rawValue, ok := dataLoopTmp["name"]; ok && rawValue != "" {
				dataLoopMap["Name"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["version"]; ok && rawValue != "" {
				dataLoopMap["Version"] = rawValue
			}
			additionalPackagesMapsArray = append(additionalPackagesMapsArray, dataLoopMap)
		}
		additionalPackagesMapsJson, err := json.Marshal(additionalPackagesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["AdditionalPackages"] = string(additionalPackagesMapsJson)
	}
	if v, ok := d.GetOk("queues"); ok {
		queuesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			if rawValue, ok := dataLoopTmp["queue_name"]; ok && rawValue != "" {
				dataLoopMap["QueueName"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["enable_scale_out"]; ok && rawValue.(bool) {
				dataLoopMap["EnableScaleOut"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["enable_scale_in"]; ok && rawValue.(bool) {
				dataLoopMap["EnableScaleIn"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["min_count"]; ok && rawValue.(int) != 0 {
				dataLoopMap["MinCount"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["max_count"]; ok && rawValue.(int) != 0 {
				dataLoopMap["MaxCount"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["initial_count"]; ok && rawValue.(int) != 0 {
				dataLoopMap["InitialCount"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["inter_connect"]; ok && rawValue != "" {
				dataLoopMap["InterConnect"] = rawValue
			}
			if rawValue := convertToInterfaceArray(dataLoopTmp["vswitch_ids"]); len(rawValue) > 0 {
				dataLoopMap["VSwitchIds"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["allocation_strategy"]; ok && rawValue != "" {
				dataLoopMap["AllocationStrategy"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["ram_role"]; ok && rawValue != "" {
				dataLoopMap["RamRole"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["hostname_prefix"]; ok && rawValue != "" {
				dataLoopMap["HostnamePrefix"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["hostname_suffix"]; ok && rawValue != "" {
				dataLoopMap["HostnameSuffix"] = rawValue
			}
			if rawValue := convertToInterfaceArray(dataLoopTmp["keep_alive_nodes"]); len(rawValue) > 0 {
				dataLoopMap["KeepAliveNodes"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["max_count_per_cycle"]; ok && rawValue.(int) != 0 {
				dataLoopMap["MaxCountPerCycle"] = rawValue
			}
			if rawValue, ok := dataLoopTmp["reserved_node_pool_id"]; ok && rawValue != "" {
				dataLoopMap["ReservedNodePoolId"] = rawValue
			}
			computeNodesMapsArray := make([]interface{}, 0)
			for _, computeNodeLoop := range convertToInterfaceArray(dataLoopTmp["compute_nodes"]) {
				computeNodeLoopTmp := computeNodeLoop.(map[string]interface{})
				computeNodeMap := make(map[string]interface{})
				if rawValue, ok := computeNodeLoopTmp["instance_type"]; ok && rawValue != "" {
					computeNodeMap["InstanceType"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["image_id"]; ok && rawValue != "" {
					computeNodeMap["ImageId"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["instance_charge_type"]; ok && rawValue != "" {
					computeNodeMap["InstanceChargeType"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["period_unit"]; ok && rawValue != "" {
					computeNodeMap["PeriodUnit"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["period"]; ok && rawValue.(int) != 0 {
					computeNodeMap["Period"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["auto_renew"]; ok && rawValue.(bool) {
					computeNodeMap["AutoRenew"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["auto_renew_period"]; ok && rawValue.(int) != 0 {
					computeNodeMap["AutoRenewPeriod"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["spot_strategy"]; ok && rawValue != "" {
					computeNodeMap["SpotStrategy"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["spot_price_limit"]; ok && rawValue.(float64) != 0 {
					computeNodeMap["SpotPriceLimit"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["duration"]; ok && rawValue.(int) != 0 {
					computeNodeMap["Duration"] = rawValue
				}
				if rawValue, ok := computeNodeLoopTmp["enable_ht"]; ok && rawValue.(bool) {
					computeNodeMap["EnableHT"] = rawValue
				}
				systemDisk := make(map[string]interface{})
				if systemDiskLoop := convertToInterfaceArray(computeNodeLoopTmp["system_disk"]); len(systemDiskLoop) > 0 {
					systemDiskTmp := systemDiskLoop[0].(map[string]interface{})
					if rawValue, ok := systemDiskTmp["category"]; ok && rawValue != "" {
						systemDisk["Category"] = rawValue
					}
					if rawValue, ok := systemDiskTmp["size"]; ok && rawValue.(int) != 0 {
						systemDisk["Size"] = rawValue
					}
					if rawValue, ok := systemDiskTmp["level"]; ok && rawValue != "" {
						systemDisk["Level"] = rawValue
					}
				}
				if len(systemDisk) > 0 {
					computeNodeMap["SystemDisk"] = systemDisk
				}
				computeNodesMapsArray = append(computeNodesMapsArray, computeNodeMap)
			}
			if len(computeNodesMapsArray) > 0 {
				dataLoopMap["ComputeNodes"] = computeNodesMapsArray
			}
			queuesMapsArray = append(queuesMapsArray, dataLoopMap)
		}
		queuesMapsJson, err := json.Marshal(queuesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Queues"] = string(queuesMapsJson)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ehpc_cluster_v2", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))

	ehpcServiceV2 := EhpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ehpcServiceV2.EhpcClusterV2StateRefreshFunc(d.Id(), "$.ClusterStatus", []string{"exception"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// CreateCluster ignores GrowInterval/IdleInterval and does not accept
	// EnableScaleOut/EnableScaleIn/MonitorSpec/SchedulerSpec; apply the full
	// updatable state through UpdateCluster right after the cluster is running.
	postCreateUpdate := false
	if _, ok := d.GetOkExists("grow_interval"); ok {
		postCreateUpdate = true
	}
	if _, ok := d.GetOkExists("idle_interval"); ok {
		postCreateUpdate = true
	}
	if _, ok := d.GetOkExists("enable_scale_out"); ok {
		postCreateUpdate = true
	}
	if _, ok := d.GetOkExists("enable_scale_in"); ok {
		postCreateUpdate = true
	}
	if len(convertToInterfaceArray(d.Get("monitor_spec"))) > 0 {
		postCreateUpdate = true
	}
	if len(convertToInterfaceArray(d.Get("scheduler_spec"))) > 0 {
		postCreateUpdate = true
	}
	if postCreateUpdate {
		updateRequest, err := ehpcClusterV2UpdatableRequest(d)
		if err != nil {
			return WrapError(err)
		}
		updateRequest["ClusterId"] = d.Id()
		action = "UpdateCluster"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("EHPC", "2024-07-30", action, query, updateRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEhpcClusterV2Read(d, meta)
}

func resourceAliCloudEhpcClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ehpcServiceV2 := EhpcServiceV2{client}

	objectRaw, err := ehpcServiceV2.DescribeEhpcClusterV2(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ehpc_cluster_v2 DescribeEhpcClusterV2 Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("client_version", objectRaw["ClientVersion"])
	d.Set("cluster_category", objectRaw["ClusterCategory"])
	d.Set("cluster_mode", objectRaw["ClusterMode"])
	d.Set("cluster_name", objectRaw["ClusterName"])
	d.Set("cluster_status", objectRaw["ClusterStatus"])
	d.Set("cluster_vswitch_id", objectRaw["ClusterVSwitchId"])
	d.Set("cluster_vpc_id", objectRaw["ClusterVpcId"])
	d.Set("create_time", objectRaw["ClusterCreateTime"])
	deletionProtection := fmt.Sprint(objectRaw["DeleteProtection"])
	d.Set("deletion_protection", formatBool(deletionProtection))
	d.Set("ehpc_version", objectRaw["EhpcVersion"])
	d.Set("enable_scale_in", objectRaw["EnableScaleIn"])
	d.Set("enable_scale_out", objectRaw["EnableScaleOut"])
	d.Set("grow_interval", formatInt(objectRaw["GrowInterval"]))
	d.Set("idle_interval", formatInt(objectRaw["IdleInterval"]))
	d.Set("max_core_count", formatInt(objectRaw["MaxCoreCount"]))
	d.Set("max_count", formatInt(objectRaw["MaxCount"]))
	d.Set("modify_time", objectRaw["ClusterModifyTime"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("security_group_id", objectRaw["SecurityGroupId"])

	clusterCustomConfigurationMaps := make([]map[string]interface{}, 0)
	clusterCustomConfigurationMap := make(map[string]interface{})
	clusterCustomConfigurationRaw := make(map[string]interface{})
	if objectRaw["ClusterCustomConfiguration"] != nil {
		clusterCustomConfigurationRaw = objectRaw["ClusterCustomConfiguration"].(map[string]interface{})
	}
	if len(clusterCustomConfigurationRaw) > 0 {
		clusterCustomConfigurationMap["script"] = clusterCustomConfigurationRaw["Script"]
		clusterCustomConfigurationMap["args"] = clusterCustomConfigurationRaw["Args"]

		clusterCustomConfigurationMaps = append(clusterCustomConfigurationMaps, clusterCustomConfigurationMap)
	}
	if err := d.Set("cluster_custom_configuration", clusterCustomConfigurationMaps); err != nil {
		return err
	}

	monitorSpecMaps := make([]map[string]interface{}, 0)
	monitorSpecMap := make(map[string]interface{})
	monitorSpecRaw := make(map[string]interface{})
	if objectRaw["MonitorSpec"] != nil {
		monitorSpecRaw = objectRaw["MonitorSpec"].(map[string]interface{})
	}
	if len(monitorSpecRaw) > 0 {
		monitorSpecMap["enable_compute_load_monitor"] = monitorSpecRaw["EnableComputeLoadMonitor"]

		monitorSpecMaps = append(monitorSpecMaps, monitorSpecMap)
	}
	if err := d.Set("monitor_spec", monitorSpecMaps); err != nil {
		return err
	}

	schedulerSpecMaps := make([]map[string]interface{}, 0)
	schedulerSpecMap := make(map[string]interface{})
	schedulerSpecRaw := make(map[string]interface{})
	if objectRaw["SchedulerSpec"] != nil {
		schedulerSpecRaw = objectRaw["SchedulerSpec"].(map[string]interface{})
	}
	if len(schedulerSpecRaw) > 0 {
		schedulerSpecMap["enable_topology_awareness"] = schedulerSpecRaw["EnableTopologyAwareness"]

		schedulerSpecMaps = append(schedulerSpecMaps, schedulerSpecMap)
	}
	if err := d.Set("scheduler_spec", schedulerSpecMaps); err != nil {
		return err
	}

	managerMaps := make([]map[string]interface{}, 0)
	managerMap := make(map[string]interface{})
	managerRaw := make(map[string]interface{})
	if objectRaw["Manager"] != nil {
		managerRaw = objectRaw["Manager"].(map[string]interface{})
	}
	if len(managerRaw) > 0 {

		directoryServiceMaps := make([]map[string]interface{}, 0)
		directoryServiceMap := make(map[string]interface{})
		directoryServiceRaw := make(map[string]interface{})
		if managerRaw["DirectoryService"] != nil {
			directoryServiceRaw = managerRaw["DirectoryService"].(map[string]interface{})
		}
		if len(directoryServiceRaw) > 0 {
			directoryServiceMap["type"] = directoryServiceRaw["Type"]
			directoryServiceMap["version"] = convertEhpcClusterV2ManagerDirectoryServiceVersionResponse(directoryServiceRaw["Version"])

			directoryServiceMaps = append(directoryServiceMaps, directoryServiceMap)
		}
		managerMap["directory_service"] = directoryServiceMaps
		dnsMaps := make([]map[string]interface{}, 0)
		dnsMap := make(map[string]interface{})
		dNSRaw := make(map[string]interface{})
		if managerRaw["DNS"] != nil {
			dNSRaw = managerRaw["DNS"].(map[string]interface{})
		}
		if len(dNSRaw) > 0 {
			dnsMap["type"] = dNSRaw["Type"]
			dnsMap["version"] = convertEhpcClusterV2ManagerDNSVersionResponse(dNSRaw["Version"])

			dnsMaps = append(dnsMaps, dnsMap)
		}
		managerMap["dns"] = dnsMaps
		managerNodeMaps := make([]map[string]interface{}, 0)
		managerNodeMap := make(map[string]interface{})
		managerNodeRaw := make(map[string]interface{})
		if managerRaw["ManagerNode"] != nil {
			managerNodeRaw = managerRaw["ManagerNode"].(map[string]interface{})
		}
		if len(managerNodeRaw) > 0 {
			managerNodeMap["auto_renew"] = managerNodeRaw["AutoRenew"]
			managerNodeMap["auto_renew_period"] = managerNodeRaw["AutoRenewPeriod"]
			managerNodeMap["duration"] = managerNodeRaw["Duration"]
			managerNodeMap["enable_ht"] = managerNodeRaw["EnableHt"]
			managerNodeMap["expired_time"] = managerNodeRaw["ExpiredTime"]
			managerNodeMap["image_id"] = managerNodeRaw["ImageId"]
			managerNodeMap["instance_charge_type"] = managerNodeRaw["InstanceChargeType"]
			managerNodeMap["instance_id"] = managerNodeRaw["InstanceId"]
			managerNodeMap["instance_type"] = managerNodeRaw["InstanceType"]
			managerNodeMap["period"] = managerNodeRaw["Period"]
			managerNodeMap["period_unit"] = managerNodeRaw["PeriodUnit"]
			managerNodeMap["spot_price_limit"] = managerNodeRaw["SpotPriceLimit"]
			managerNodeMap["spot_strategy"] = managerNodeRaw["SpotStrategy"]

			systemDiskMaps := make([]map[string]interface{}, 0)
			systemDiskMap := make(map[string]interface{})
			systemDiskRaw := make(map[string]interface{})
			if managerNodeRaw["SystemDisk"] != nil {
				systemDiskRaw = managerNodeRaw["SystemDisk"].(map[string]interface{})
			}
			if len(systemDiskRaw) > 0 {
				systemDiskMap["category"] = systemDiskRaw["Category"]
				systemDiskMap["level"] = systemDiskRaw["Level"]
				systemDiskMap["size"] = systemDiskRaw["Size"]

				systemDiskMaps = append(systemDiskMaps, systemDiskMap)
			}
			managerNodeMap["system_disk"] = systemDiskMaps
			managerNodeMaps = append(managerNodeMaps, managerNodeMap)
		}
		managerMap["manager_node"] = managerNodeMaps
		schedulerMaps := make([]map[string]interface{}, 0)
		schedulerMap := make(map[string]interface{})
		schedulerRaw := make(map[string]interface{})
		if managerRaw["Scheduler"] != nil {
			schedulerRaw = managerRaw["Scheduler"].(map[string]interface{})
		}
		if len(schedulerRaw) > 0 {
			schedulerMap["type"] = convertEhpcClusterV2ManagerSchedulerTypeResponse(schedulerRaw["Type"])
			schedulerMap["version"] = convertEhpcClusterV2ManagerSchedulerVersionResponse(schedulerRaw["Version"])

			schedulerMaps = append(schedulerMaps, schedulerMap)
		}
		managerMap["scheduler"] = schedulerMaps
		managerMaps = append(managerMaps, managerMap)
	}
	if err := d.Set("manager", managerMaps); err != nil {
		return err
	}

	objectRaw, err = ehpcServiceV2.DescribeClusterV2ListSharedStorages(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	// 从 ListSharedStorages API 获取 SharedStorages 数组
	sharedStoragesRawObj, _ := jsonpath.Get("$.SharedStorages[*]", objectRaw)
	sharedStoragesMaps := make([]map[string]interface{}, 0)

	if sharedStoragesRawObj != nil {
		sharedStoragesRaw := convertToInterfaceArray(sharedStoragesRawObj)
		// 遍历每个 SharedStorage，将其中的 MountInfo 数组铺平
		for _, sharedStorageRaw := range sharedStoragesRaw {
			sharedStorage := sharedStorageRaw.(map[string]interface{})
			fileSystemId := sharedStorage["FileSystemId"]

			// 获取当前 SharedStorage 下的 MountInfo 数组
			mountInfoRawObj, _ := jsonpath.Get("$.MountInfo[*]", sharedStorage)
			if mountInfoRawObj != nil {
				mountInfoArray := convertToInterfaceArray(mountInfoRawObj)

				// 将 MountInfo 数组中的每一项展开到结果数组中
				for _, mountInfoRaw := range mountInfoArray {
					mountInfo := mountInfoRaw.(map[string]interface{})
					sharedStoragesMap := make(map[string]interface{})

					sharedStoragesMap["file_system_id"] = fileSystemId
					sharedStoragesMap["mount_directory"] = mountInfo["MountDirectory"]
					sharedStoragesMap["mount_options"] = mountInfo["MountOptions"]
					sharedStoragesMap["mount_target_domain"] = mountInfo["MountTarget"]
					sharedStoragesMap["protocol_type"] = mountInfo["ProtocolType"]

					// 处理 StorageDirectory：如果 MountDirectory 是 /home 或 /opt，需要去掉 /ehpc 后缀
					if storageDir, ok := mountInfo["StorageDirectory"].(string); ok {
						mountDir, _ := mountInfo["MountDirectory"].(string)
						nasDirectory := storageDir
						if (mountDir == "/home" || mountDir == "/opt") && len(storageDir) > 0 {
							// 去掉 /ehpc 后缀部分，使用 strings.Split 分割
							parts := strings.Split(storageDir, "/ehpc")
							if len(parts) > 0 {
								nasDirectory = parts[0]
							}
						}
						sharedStoragesMap["nas_directory"] = nasDirectory
					}

					sharedStoragesMaps = append(sharedStoragesMaps, sharedStoragesMap)
				}
			}
		}
	}

	if err := d.Set("shared_storages", sharedStoragesMaps); err != nil {
		return err
	}

	return nil
}

// ehpcClusterV2UpdatableRequest builds the complete set of updatable parameters for
// UpdateCluster. UpdateCluster applies a full state: parameters that are not carried
// in the request are reset to their defaults (DeletionProtection is turned off,
// GrowInterval/IdleInterval fall back to 2/6), so every updatable attribute must be
// resent together to preserve the values applied at create time.
func ehpcClusterV2UpdatableRequest(d *schema.ResourceData) (map[string]interface{}, error) {
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}
	if v := d.Get("client_version"); v.(string) != "" {
		request["ClientVersion"] = v
	}
	if v := d.Get("cluster_name"); v.(string) != "" {
		request["ClusterName"] = v
	}
	if v := d.Get("cluster_description"); v.(string) != "" {
		request["ClusterDescription"] = v
	}
	request["MaxCount"] = d.Get("max_count")
	request["MaxCoreCount"] = d.Get("max_core_count")
	request["GrowInterval"] = d.Get("grow_interval")
	request["IdleInterval"] = d.Get("idle_interval")
	request["EnableScaleOut"] = d.Get("enable_scale_out")
	request["EnableScaleIn"] = d.Get("enable_scale_in")

	clusterCustomConfiguration := make(map[string]interface{})
	script1, _ := jsonpath.Get("$[0].script", d.Get("cluster_custom_configuration"))
	if script1 != nil && script1 != "" {
		clusterCustomConfiguration["Script"] = script1
	}
	args1, _ := jsonpath.Get("$[0].args", d.Get("cluster_custom_configuration"))
	if args1 != nil && args1 != "" {
		clusterCustomConfiguration["Args"] = args1
	}
	if len(clusterCustomConfiguration) > 0 {
		clusterCustomConfigurationJson, err := json.Marshal(clusterCustomConfiguration)
		if err != nil {
			return nil, WrapError(err)
		}
		request["ClusterCustomConfiguration"] = string(clusterCustomConfigurationJson)
	}

	monitorSpec := make(map[string]interface{})
	enableComputeLoadMonitor1, _ := jsonpath.Get("$[0].enable_compute_load_monitor", d.Get("monitor_spec"))
	if enableComputeLoadMonitor1 != nil && enableComputeLoadMonitor1 != "" {
		monitorSpec["EnableComputeLoadMonitor"] = enableComputeLoadMonitor1
	}
	if len(monitorSpec) > 0 {
		monitorSpecJson, err := json.Marshal(monitorSpec)
		if err != nil {
			return nil, WrapError(err)
		}
		request["MonitorSpec"] = string(monitorSpecJson)
	}

	schedulerSpec := make(map[string]interface{})
	enableTopologyAwareness1, _ := jsonpath.Get("$[0].enable_topology_awareness", d.Get("scheduler_spec"))
	if enableTopologyAwareness1 != nil && enableTopologyAwareness1 != "" {
		schedulerSpec["EnableTopologyAwareness"] = enableTopologyAwareness1
	}
	if len(schedulerSpec) > 0 {
		schedulerSpecJson, err := json.Marshal(schedulerSpec)
		if err != nil {
			return nil, WrapError(err)
		}
		request["SchedulerSpec"] = string(schedulerSpecJson)
	}

	return request, nil
}

func resourceAliCloudEhpcClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateCluster"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = d.Id()

	if d.HasChange("client_version") || d.HasChange("deletion_protection") || d.HasChange("cluster_name") || d.HasChange("cluster_description") || d.HasChange("max_count") || d.HasChange("max_core_count") ||
		d.HasChange("grow_interval") || d.HasChange("idle_interval") || d.HasChange("enable_scale_out") || d.HasChange("enable_scale_in") ||
		d.HasChange("cluster_custom_configuration") || d.HasChange("monitor_spec") || d.HasChange("scheduler_spec") {
		update = true
	}

	if update {
		fullRequest, err := ehpcClusterV2UpdatableRequest(d)
		if err != nil {
			return WrapError(err)
		}
		for k, v := range fullRequest {
			request[k] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
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

	return resourceAliCloudEhpcClusterV2Read(d, meta)
}

func resourceAliCloudEhpcClusterV2Delete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidClusterStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ClusterNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertEhpcClusterV2ManagerDirectoryServiceVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
func convertEhpcClusterV2ManagerDNSVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
func convertEhpcClusterV2ManagerSchedulerTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "slurm22":
		return "SLURM"
	}
	return source
}
func convertEhpcClusterV2ManagerSchedulerVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
