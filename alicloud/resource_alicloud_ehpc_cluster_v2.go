package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}
	clusterCredentials := make(map[string]interface{})

	if v := d.Get("cluster_credentials"); v != nil {
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
	d.Set("cluster_vswitch_id", objectRaw["ClusterVSwitchId"])
	d.Set("cluster_vpc_id", objectRaw["ClusterVpcId"])
	d.Set("create_time", objectRaw["ClusterCreateTime"])
	deletionProtection := fmt.Sprint(objectRaw["DeleteProtection"])
	d.Set("deletion_protection", formatBool(deletionProtection))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("security_group_id", objectRaw["SecurityGroupId"])

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

	if d.HasChange("client_version") {
		update = true
		request["ClientVersion"] = d.Get("client_version")
	}

	if d.HasChange("deletion_protection") {
		update = true
		request["DeletionProtection"] = d.Get("deletion_protection")
	}

	if d.HasChange("cluster_name") {
		update = true
		request["ClusterName"] = d.Get("cluster_name")
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
