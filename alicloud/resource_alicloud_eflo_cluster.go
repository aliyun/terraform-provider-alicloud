package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEfloCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloClusterCreate,
		Read:   resourceAliCloudEfloClusterRead,
		Update: resourceAliCloudEfloClusterUpdate,
		Delete: resourceAliCloudEfloClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"component_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"basic_args": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"node_units": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hpn_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_failed_node_tasks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpd_info": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpd_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vpd_subnets": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"ip_allocation_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"machine_type_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bonds": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"subnet": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"machine_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"bond_policy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bonds": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"subnet": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"bond_default_subnet": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"node_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bonds": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"subnet": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"node_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tail_ip_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"new_vpd_info": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_link_cidr": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"monitor_vpc_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cen_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"11111"}, false),
									},
									"vpd_cidr": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"monitor_vswitch_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vpd_subnets": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subnet_cidr": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"zone_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"subnet_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"cloud_link_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"nimiz_vswitches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_group_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"node_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"login_password": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"open_eni_jumbo_frame": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudEfloClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("networks"); !IsNil(v) {
		vSwitchId1, _ := jsonpath.Get("$[0].vswitch_id", v)
		if vSwitchId1 != nil && vSwitchId1 != "" {
			objectDataLocalMap["VSwitchId"] = vSwitchId1
		}
		if v, ok := d.GetOk("networks"); ok {
			localData, err := jsonpath.Get("$[0].ip_allocation_policy", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				localData1 := make(map[string]interface{})
				bondDefaultSubnet1, _ := jsonpath.Get("$[0].bond_default_subnet", dataLoopTmp["bond_policy"])
				if bondDefaultSubnet1 != nil && bondDefaultSubnet1 != "" {
					localData1["BondDefaultSubnet"] = bondDefaultSubnet1
				}
				if v, ok := dataLoopTmp["bond_policy"]; ok {
					localData2, err := jsonpath.Get("$[0].bonds", v)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps2 := make([]interface{}, 0)
					for _, dataLoop2 := range localData2.([]interface{}) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["Subnet"] = dataLoop2Tmp["subnet"]
						dataLoop2Map["Name"] = dataLoop2Tmp["name"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					localData1["Bonds"] = localMaps2
				}

				dataLoopMap["BondPolicy"] = localData1
				localMaps3 := make([]interface{}, 0)
				localData3 := dataLoopTmp["node_policy"]
				for _, dataLoop3 := range localData3.([]interface{}) {
					dataLoop3Tmp := dataLoop3.(map[string]interface{})
					dataLoop3Map := make(map[string]interface{})
					dataLoop3Map["NodeId"] = dataLoop3Tmp["node_id"]
					localData2, err := jsonpath.Get("$.bonds", dataLoop3Tmp)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps2 := make([]interface{}, 0)
					for _, dataLoop2 := range localData2.([]interface{}) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["Subnet"] = dataLoop2Tmp["subnet"]
						dataLoop2Map["Name"] = dataLoop2Tmp["name"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					dataLoop3Map["Bonds"] = localMaps2
					localMaps3 = append(localMaps3, dataLoop3Map)
				}
				dataLoopMap["NodePolicy"] = localMaps3
				localMaps5 := make([]interface{}, 0)
				localData5 := dataLoopTmp["machine_type_policy"]
				for _, dataLoop5 := range localData5.([]interface{}) {
					dataLoop5Tmp := dataLoop5.(map[string]interface{})
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["MachineType"] = dataLoop5Tmp["machine_type"]
					localData2, err := jsonpath.Get("$.bonds", dataLoop5Tmp)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps2 := make([]interface{}, 0)
					for _, dataLoop2 := range localData2.([]interface{}) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["Subnet"] = dataLoop2Tmp["subnet"]
						dataLoop2Map["Name"] = dataLoop2Tmp["name"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					dataLoop5Map["Bonds"] = localMaps2
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				dataLoopMap["MachineTypePolicy"] = localMaps5
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["IpAllocationPolicy"] = localMaps
		}

		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			objectDataLocalMap["VpcId"] = vpcId1
		}
		newVpdInfo := make(map[string]interface{})
		if v, ok := d.GetOk("networks"); ok {
			localData7, err := jsonpath.Get("$[0].new_vpd_info[0].vpd_subnets", v)
			if err != nil {
				localData7 = make([]interface{}, 0)
			}
			localMaps7 := make([]interface{}, 0)
			for _, dataLoop7 := range localData7.([]interface{}) {
				dataLoop7Tmp := make(map[string]interface{})
				if dataLoop7 != nil {
					dataLoop7Tmp = dataLoop7.(map[string]interface{})
				}
				dataLoop7Map := make(map[string]interface{})
				dataLoop7Map["SubnetCidr"] = dataLoop7Tmp["subnet_cidr"]
				dataLoop7Map["ZoneId"] = dataLoop7Tmp["zone_id"]
				dataLoop7Map["SubnetType"] = dataLoop7Tmp["subnet_type"]
				localMaps7 = append(localMaps7, dataLoop7Map)
			}
			newVpdInfo["VpdSubnets"] = localMaps7
		}

		cloudLinkId1, _ := jsonpath.Get("$[0].new_vpd_info[0].cloud_link_id", v)
		if cloudLinkId1 != nil && cloudLinkId1 != "" {
			newVpdInfo["CloudLinkId"] = cloudLinkId1
		}
		cloudLinkCidr1, _ := jsonpath.Get("$[0].new_vpd_info[0].cloud_link_cidr", v)
		if cloudLinkCidr1 != nil && cloudLinkCidr1 != "" {
			newVpdInfo["CloudLinkCidr"] = cloudLinkCidr1
		}
		monitorVswitchId1, _ := jsonpath.Get("$[0].new_vpd_info[0].monitor_vswitch_id", v)
		if monitorVswitchId1 != nil && monitorVswitchId1 != "" {
			newVpdInfo["MonitorVswitchId"] = monitorVswitchId1
		}
		vpdCidr1, _ := jsonpath.Get("$[0].new_vpd_info[0].vpd_cidr", v)
		if vpdCidr1 != nil && vpdCidr1 != "" {
			newVpdInfo["VpdCidr"] = vpdCidr1
		}
		cenId1, _ := jsonpath.Get("$[0].new_vpd_info[0].cen_id", v)
		if cenId1 != nil && cenId1 != "" {
			newVpdInfo["CenId"] = cenId1
		}
		monitorVpcId1, _ := jsonpath.Get("$[0].new_vpd_info[0].monitor_vpc_id", v)
		if monitorVpcId1 != nil && monitorVpcId1 != "" {
			newVpdInfo["MonitorVpcId"] = monitorVpcId1
		}

		objectDataLocalMap["NewVpdInfo"] = newVpdInfo
		securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
		if securityGroupId1 != nil && securityGroupId1 != "" {
			objectDataLocalMap["SecurityGroupId"] = securityGroupId1
		}
		vpdInfo := make(map[string]interface{})
		vpdSubnets2, _ := jsonpath.Get("$[0].vpd_info[0].vpd_subnets", v)
		if vpdSubnets2 != nil && vpdSubnets2 != "" {
			vpdInfo["VpdSubnets"] = vpdSubnets2
		}
		vpdId1, _ := jsonpath.Get("$[0].vpd_info[0].vpd_id", v)
		if vpdId1 != nil && vpdId1 != "" {
			vpdInfo["VpdId"] = vpdId1
		}

		objectDataLocalMap["VpdInfo"] = vpdInfo
		vSwitchZoneId1, _ := jsonpath.Get("$[0].vswitch_zone_id", v)
		if vSwitchZoneId1 != nil && vSwitchZoneId1 != "" {
			objectDataLocalMap["VSwitchZoneId"] = vSwitchZoneId1
		}
		tailIpVersion1, _ := jsonpath.Get("$[0].tail_ip_version", v)
		if tailIpVersion1 != nil && tailIpVersion1 != "" {
			objectDataLocalMap["TailIpVersion"] = tailIpVersion1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Networks"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("node_groups"); ok {
		nodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop8 := range v.([]interface{}) {
			dataLoop8Tmp := dataLoop8.(map[string]interface{})
			dataLoop8Map := make(map[string]interface{})
			dataLoop8Map["UserData"] = dataLoop8Tmp["user_data"]
			dataLoop8Map["ZoneId"] = dataLoop8Tmp["zone_id"]
			dataLoop8Map["NodeGroupDescription"] = dataLoop8Tmp["node_group_description"]
			localMaps8 := make([]interface{}, 0)
			localData9 := dataLoop8Tmp["nodes"]
			for _, dataLoop9 := range localData9.([]interface{}) {
				dataLoop9Tmp := dataLoop9.(map[string]interface{})
				dataLoop9Map := make(map[string]interface{})
				dataLoop9Map["Hostname"] = dataLoop9Tmp["hostname"]
				dataLoop9Map["VSwitchId"] = dataLoop9Tmp["vswitch_id"]
				dataLoop9Map["VpcId"] = dataLoop9Tmp["vpc_id"]
				dataLoop9Map["NodeId"] = dataLoop9Tmp["node_id"]
				dataLoop9Map["LoginPassword"] = dataLoop9Tmp["login_password"]
				localMaps8 = append(localMaps8, dataLoop9Map)
			}
			dataLoop8Map["Nodes"] = localMaps8
			dataLoop8Map["NodeGroupName"] = dataLoop8Tmp["node_group_name"]
			dataLoop8Map["MachineType"] = dataLoop8Tmp["machine_type"]
			dataLoop8Map["ImageId"] = dataLoop8Tmp["image_id"]
			nodeGroupsMapsArray = append(nodeGroupsMapsArray, dataLoop8Map)
		}
		nodeGroupsMapsJson, err := json.Marshal(nodeGroupsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["NodeGroups"] = string(nodeGroupsMapsJson)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request["ClusterName"] = v
	}
	if v, ok := d.GetOk("cluster_description"); ok {
		request["ClusterDescription"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("components"); ok {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop11 := range v.([]interface{}) {
			dataLoop11Tmp := dataLoop11.(map[string]interface{})
			dataLoop11Map := make(map[string]interface{})
			localData12 := make(map[string]interface{})
			basicArgs1, _ := jsonpath.Get("$[0].basic_args", dataLoop11Tmp["component_config"])
			if basicArgs1 != nil && basicArgs1 != "" {
				localData12["BasicArgs"] = basicArgs1
			}
			nodeUnits1, _ := jsonpath.Get("$[0].node_units", dataLoop11Tmp["component_config"])
			if nodeUnits1 != nil && nodeUnits1 != "" {
				localData12["NodeUnits"] = nodeUnits1
			}
			dataLoop11Map["ComponentConfig"] = localData12
			dataLoop11Map["ComponentType"] = dataLoop11Tmp["component_type"]
			componentsMapsArray = append(componentsMapsArray, dataLoop11Map)
		}
		componentsMapsJson, err := json.Marshal(componentsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Components"] = string(componentsMapsJson)
	}

	if v, ok := d.GetOkExists("ignore_failed_node_tasks"); ok {
		request["IgnoreFailedNodeTasks"] = v
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		request["ClusterType"] = v
	}
	if v, ok := d.GetOk("nimiz_vswitches"); ok {
		nimizVSwitchesMapsArray := v.([]interface{})
		nimizVSwitchesMapsJson, err := json.Marshal(nimizVSwitchesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["NimizVSwitches"] = string(nimizVSwitchesMapsJson)
	}

	if v, ok := d.GetOk("hpn_zone"); ok {
		request["HpnZone"] = v
	}
	if v, ok := d.GetOkExists("open_eni_jumbo_frame"); ok {
		request["OpenEniJumboFrame"] = v
	}
	wait := incrementalWait(10*time.Second, 60*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
		if err != nil {
			addDebug(action, err, request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))

	return resourceAliCloudEfloClusterUpdate(d, meta)
}

func resourceAliCloudEfloClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_cluster DescribeEfloCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_description", objectRaw["ClusterDescription"])
	d.Set("cluster_name", objectRaw["ClusterName"])
	d.Set("cluster_type", objectRaw["ClusterType"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["OperatingState"])

	objectRaw, err = efloServiceV2.DescribeClusterListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEfloClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "Cluster"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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

	if d.HasChange("tags") {
		efloServiceV2 := EfloServiceV2{client}
		if err := efloServiceV2.SetResourceTags(d, "Cluster"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudEfloClusterRead(d, meta)
}

func resourceAliCloudEfloClusterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
