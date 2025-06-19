// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudEfloNodeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloNodeGroupCreate,
		Read:   resourceAliCloudEfloNodeGroupRead,
		Update: resourceAliCloudEfloNodeGroupUpdate,
		Delete: resourceAliCloudEfloNodeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"az": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ignore_failed_node_tasks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"login_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_group_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					if v, ok := m["node_id"]; ok {
						buf.WriteString(fmt.Sprint(v))
					}
					return hashcode.String(buf.String())
				},
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
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"login_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpd_subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEfloNodeGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNodeGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("node_group_description"); ok {
		objectDataLocalMap["NodeGroupDescription"] = v
	}

	if v, ok := d.GetOk("key_pair_name"); ok {
		objectDataLocalMap["KeyPairName"] = v
	}

	if v, ok := d.GetOk("login_password"); ok {
		objectDataLocalMap["LoginPassword"] = v
	}

	if v, ok := d.GetOk("az"); ok {
		objectDataLocalMap["Az"] = v
	}

	if v, ok := d.GetOk("machine_type"); ok {
		objectDataLocalMap["MachineType"] = v
	}

	if v, ok := d.GetOk("image_id"); ok {
		objectDataLocalMap["ImageId"] = v
	}

	if v, ok := d.GetOk("node_group_name"); ok {
		objectDataLocalMap["NodeGroupName"] = v
	}

	objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
	if err != nil {
		return WrapError(err)
	}
	request["NodeGroup"] = string(objectDataLocalMapJson)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_node_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ClusterId"], response["NodeGroupId"]))

	return resourceAliCloudEfloNodeGroupUpdate(d, meta)
}

func resourceAliCloudEfloNodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloNodeGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_node_group DescribeEfloNodeGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("az", objectRaw["ZoneId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("image_id", objectRaw["ImageId"])
	d.Set("machine_type", objectRaw["MachineType"])
	d.Set("node_group_description", objectRaw["Description"])
	d.Set("node_group_name", objectRaw["GroupName"])
	d.Set("cluster_id", objectRaw["ClusterId"])
	d.Set("node_group_id", objectRaw["GroupId"])

	objectRaw, err = efloServiceV2.DescribeNodeGroupListClusterNodes(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	loginPasswordMap := make(map[string]interface{})
	if v := d.Get("nodes"); !IsNil(v) {
		if v, ok := d.GetOk("nodes"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				loginPasswordMap[dataLoopTmp["node_id"].(string)] = dataLoopTmp["login_password"]
			}
		}
	}

	nodesRaw, _ := jsonpath.Get("$.Nodes", objectRaw)

	nodesMaps := make([]map[string]interface{}, 0)
	if nodesRaw != nil {
		for _, nodesChildRaw := range nodesRaw.([]interface{}) {
			nodesMap := make(map[string]interface{})
			nodesChildRaw := nodesChildRaw.(map[string]interface{})
			nodesMap["hostname"] = nodesChildRaw["Hostname"]
			nodesMap["node_id"] = nodesChildRaw["NodeId"]
			nodesMap["vswitch_id"] = nodesChildRaw["VSwitchId"]
			nodesMap["vpc_id"] = nodesChildRaw["VpcId"]
			if password, ok := loginPasswordMap[fmt.Sprint(nodesMap["node_id"])]; ok {
				nodesMap["login_password"] = password
			}

			nodesMaps = append(nodesMaps, nodesMap)
		}
	}
	if err := d.Set("nodes", nodesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEfloNodeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateNodeGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NodeGroupId"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("key_pair_name") {
		update = true
		request["KeyPairName"] = d.Get("key_pair_name")
	}

	if !d.IsNewResource() && d.HasChange("image_id") {
		update = true
	}
	request["ImageId"] = d.Get("image_id")
	if !d.IsNewResource() && d.HasChange("login_password") {
		update = true
		request["LoginPassword"] = d.Get("login_password")
	}

	if !d.IsNewResource() && d.HasChange("node_group_name") {
		update = true
	}
	request["NewNodeGroupName"] = d.Get("node_group_name")
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

	if d.HasChange("nodes") {
		oldEntry, newEntry := d.GetChange("nodes")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ShrinkCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ClusterId"] = parts[0]
			request["RegionId"] = client.RegionId
			objectDataLocalMap := make(map[string]interface{})

			localMaps := make([]interface{}, 0)
			for _, dataLoop := range removed.List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["NodeId"] = dataLoopTmp["node_id"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["Nodes"] = localMaps

			if v, ok := d.GetOk("node_group_id"); ok {
				objectDataLocalMap["NodeGroupId"] = v
			}

			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, objectDataLocalMap)
			objectDataLocalMapJson, err := json.Marshal(NodeGroupsMap)
			if err != nil {
				return WrapError(err)
			}
			request["NodeGroups"] = string(objectDataLocalMapJson)

			if v, ok := d.GetOkExists("ignore_failed_node_tasks"); ok {
				request["IgnoreFailedNodeTasks"] = v
			}
			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.NodeGroupId", parts[1])
			_ = json.Unmarshal([]byte(jsonString), &request)

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
			efloServiceV2 := EfloServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, efloServiceV2.DescribeAsyncEfloNodeGroupStateRefreshFunc(d, response, "$.TaskState", []string{}))
			if jobDetail, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
			}

		}

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ExtendCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ClusterId"] = parts[0]
			request["RegionId"] = client.RegionId
			objectDataLocalMap := make(map[string]interface{})

			localMaps := make([]interface{}, 0)
			for _, dataLoop := range added.List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["NodeId"] = dataLoopTmp["node_id"]
				dataLoopMap["Hostname"] = dataLoopTmp["hostname"]
				dataLoopMap["LoginPassword"] = dataLoopTmp["login_password"]
				dataLoopMap["VpcId"] = dataLoopTmp["vpc_id"]
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["Nodes"] = localMaps

			if v, ok := d.GetOk("node_group_id"); ok {
				objectDataLocalMap["NodeGroupId"] = v
			}

			if v, ok := d.GetOk("user_data"); ok {
				objectDataLocalMap["UserData"] = v
			}

			if v, ok := d.GetOk("zone_id"); ok {
				objectDataLocalMap["ZoneId"] = v
			}

			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, objectDataLocalMap)
			objectDataLocalMapJson, err := json.Marshal(NodeGroupsMap)
			if err != nil {
				return WrapError(err)
			}
			request["NodeGroups"] = string(objectDataLocalMapJson)

			localData1 := d.Get("ip_allocation_policy").([]interface{})
			ipAllocationPolicyMapsArray := make([]interface{}, 0)
			for _, dataLoop1 := range localData1 {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})

				localData3 := make(map[string]interface{})
				bondDefaultSubnet1, _ := jsonpath.Get("$[0].bond_default_subnet", dataLoop1Tmp["bond_policy"])
				if bondDefaultSubnet1 != nil && bondDefaultSubnet1 != "" {
					localData3["BondDefaultSubnet"] = bondDefaultSubnet1
				}
				if v, ok := dataLoop1Tmp["bond_policy"]; ok {
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
					localData3["Bonds"] = localMaps2
				}
				dataLoop1Map["BondPolicy"] = localData3

				localMaps3 := make([]interface{}, 0)
				localData4 := dataLoop1Tmp["node_policy"]
				for _, dataLoop4 := range localData4.([]interface{}) {
					dataLoop4Tmp := dataLoop4.(map[string]interface{})
					dataLoop4Map := make(map[string]interface{})
					dataLoop4Map["NodeId"] = dataLoop4Tmp["node_id"]
					localData2, err := jsonpath.Get("$.bonds", dataLoop4Tmp)
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
					dataLoop4Map["Bonds"] = localMaps2
					localMaps3 = append(localMaps3, dataLoop4Map)
				}
				dataLoop1Map["NodePolicy"] = localMaps3
				localMaps5 := make([]interface{}, 0)
				localData6 := dataLoop1Tmp["machine_type_policy"]
				for _, dataLoop6 := range localData6.([]interface{}) {
					dataLoop6Tmp := dataLoop6.(map[string]interface{})
					dataLoop6Map := make(map[string]interface{})
					dataLoop6Map["MachineType"] = dataLoop6Tmp["machine_type"]
					localData2, err := jsonpath.Get("$.bonds", dataLoop6Tmp)
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
					dataLoop6Map["Bonds"] = localMaps2
					localMaps5 = append(localMaps5, dataLoop6Map)
				}
				dataLoop1Map["MachineTypePolicy"] = localMaps5
				ipAllocationPolicyMapsArray = append(ipAllocationPolicyMapsArray, dataLoop1Map)
			}
			ipAllocationPolicyMapsJson, err := json.Marshal(ipAllocationPolicyMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["IpAllocationPolicy"] = string(ipAllocationPolicyMapsJson)

			localData8 := d.Get("vpd_subnets").([]interface{})
			vpdSubnetsMapsArray := localData8
			vpdSubnetsMapsJson, err := json.Marshal(vpdSubnetsMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["VpdSubnets"] = string(vpdSubnetsMapsJson)

			if v, ok := d.GetOkExists("ignore_failed_node_tasks"); ok {
				request["IgnoreFailedNodeTasks"] = v
			}
			if v, ok := d.GetOk("vswitch_zone_id"); ok {
				request["VSwitchZoneId"] = v
			}
			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.NodeGroupId", parts[1])
			_ = json.Unmarshal([]byte(jsonString), &request)

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
			efloServiceV2 := EfloServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, efloServiceV2.DescribeAsyncEfloNodeGroupStateRefreshFunc(d, response, "$.TaskState", []string{}))
			if jobDetail, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
			}

		}

	}
	d.Partial(false)
	return resourceAliCloudEfloNodeGroupRead(d, meta)
}

func resourceAliCloudEfloNodeGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteNodeGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = parts[0]
	request["NodeGroupId"] = parts[1]
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
