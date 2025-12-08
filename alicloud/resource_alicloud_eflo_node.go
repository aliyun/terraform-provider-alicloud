// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
	"github.com/tidwall/sjson"
)

func resourceAliCloudEfloNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloNodeCreate,
		Read:   resourceAliCloudEfloNodeRead,
		Update: resourceAliCloudEfloNodeUpdate,
		Delete: resourceAliCloudEfloNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"billing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"classify": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"computing_server": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'computing_server' has been deprecated from provider version 1.261.0. New field 'machine_type' instead.",
				ConflictsWith: []string{"machine_type"},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discount_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hpn_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
									"hostname": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"login_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"payment_ratio": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"product_form": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_arch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stage_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Unused", "Using"}, false),
			},
			"tags": tagsSchema(),
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"install_pai": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudEfloNodeCreate(d *schema.ResourceData, meta interface{}) error {

	installPai := false
	if v, ok := d.GetOk("install_pai"); ok && v.(bool) {
		installPai = true
	}

	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("payment_type"); ok && InArray(fmt.Sprint(v), []string{"Subscription"}) {
		action := "CreateInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})

		request["ClientToken"] = buildClientToken(action)

		request["SubscriptionType"] = d.Get("payment_type")
		parameterMapList := make([]map[string]interface{}, 0)
		if v, ok := d.GetOk("server_arch"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "ServerArch",
				"Value": v,
			})
		}
		if v, ok := d.GetOk("hpn_zone"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "HpnZone",
				"Value": v,
			})
		}
		if v, ok := d.GetOk("stage_num"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "StageNum",
				"Value": v,
			})
		}
		if v, ok := d.GetOk("payment_ratio"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "PaymentRatio",
				"Value": v,
			})
		}
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "RegionId",
			"Value": client.RegionId,
		})
		if v, ok := d.GetOk("classify"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "Classify",
				"Value": v,
			})
		}
		discountlevelCode := "discountlevel"
		if installPai {
			discountlevelCode = "DiscountLevel"
		}
		if v, ok := d.GetOk("discount_level"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  discountlevelCode,
				"Value": v,
			})
		}
		if v, ok := d.GetOk("billing_cycle"); ok {
			if v.(string) == "1month" && installPai {
				v = "1m"
			}
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "BillingCycle",
				"Value": v,
			})
		}
		computingServerCode := "computingserver"
		if installPai {
			computingServerCode = "ComputingServer"
		}
		if v, ok := d.GetOk("machine_type"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  computingServerCode,
				"Value": v,
			})
		}
		if v, ok := d.GetOk("computing_server"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  computingServerCode,
				"Value": v,
			})
		}
		if v, ok := d.GetOk("zone"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "Zone",
				"Value": v,
			})
		}
		if v, ok := d.GetOk("product_form"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "ProductForm",
				"Value": v,
			})
		}
		request["Parameter"] = parameterMapList

		if v, ok := d.GetOk("renewal_status"); ok {
			request["RenewalStatus"] = v
		}
		if v, ok := d.GetOkExists("period"); ok {
			request["Period"] = v
		}
		if v, ok := d.GetOkExists("renew_period"); ok {
			request["RenewPeriod"] = v
		}
		var endpoint string
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_cn"
		if installPai {
			request["ProductCode"] = "learn"
			request["ProductType"] = "learn_eflocomputing_public_cn"
		}
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_cn"
			if installPai {
				return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
			}
		}
		if client.IsInternationalAccount() {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_eflocomputing_public_intl"
			if installPai {
				request["ProductCode"] = "learn"
				request["ProductType"] = "learn_eflocomputing_public_intl"
			}
			if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_computinginstance_public_intl"
				if installPai {
					return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
				}
			}
		}
		if request["SubscriptionType"] == "" {
			request["SubscriptionType"] = "Subscription"
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if IsExpectedErrors(err, []string{"CSS_CHECK_ORDER_ERROR", "InternalError", "SYSTEM.CONCURRENT_OPERATE"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_eflocomputing_public_intl"
					if installPai {
						request["ProductCode"] = "learn"
						request["ProductType"] = "learn_eflocomputing_public_intl"
					}
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "bccluster"
						request["ProductType"] = "bccluster_computinginstance_public_intl"
						if installPai {
							return resource.RetryableError(err)
						}
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_node", action, AlibabaCloudSdkGoERROR)
		}

		id, _ := jsonpath.Get("$.Data.InstanceId", response)
		d.SetId(fmt.Sprint(id))

		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Unused"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if v, ok := d.GetOk("payment_type"); ok && InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
		action := "ExtendCluster"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId

		if v, ok := d.GetOk("ip_allocation_policy"); ok {
			ipAllocationPolicyMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				localData1 := make(map[string]interface{})
				if v, ok := dataLoopTmp["bond_policy"]; ok {
					localData2, err := jsonpath.Get("$[0].bonds", v)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps1 := make([]interface{}, 0)
					for _, dataLoop2 := range convertToInterfaceArray(localData2) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["Subnet"] = dataLoop2Tmp["subnet"]
						dataLoop2Map["Name"] = dataLoop2Tmp["name"]
						localMaps1 = append(localMaps1, dataLoop2Map)
					}
					localData1["Bonds"] = localMaps1
				}

				bondDefaultSubnet1, _ := jsonpath.Get("$[0].bond_default_subnet", dataLoopTmp["bond_policy"])
				if bondDefaultSubnet1 != nil && bondDefaultSubnet1 != "" {
					localData1["BondDefaultSubnet"] = bondDefaultSubnet1
				}
				if len(localData1) > 0 {
					dataLoopMap["BondPolicy"] = localData1
				}
				localMaps2 := make([]interface{}, 0)
				localData3 := dataLoopTmp["node_policy"]
				for _, dataLoop3 := range convertToInterfaceArray(localData3) {
					dataLoop3Tmp := dataLoop3.(map[string]interface{})
					dataLoop3Map := make(map[string]interface{})
					localMaps3 := make([]interface{}, 0)
					localData4 := dataLoop3Tmp["bonds"]
					for _, dataLoop4 := range convertToInterfaceArray(localData4) {
						dataLoop4Tmp := dataLoop4.(map[string]interface{})
						dataLoop4Map := make(map[string]interface{})
						dataLoop4Map["Name"] = dataLoop4Tmp["name"]
						dataLoop4Map["Subnet"] = dataLoop4Tmp["subnet"]
						localMaps3 = append(localMaps3, dataLoop4Map)
					}
					dataLoop3Map["Bonds"] = localMaps3
					dataLoop3Map["Hostname"] = dataLoop3Tmp["hostname"]
					dataLoop3Map["NodeId"] = dataLoop3Tmp["node_id"]
					localMaps2 = append(localMaps2, dataLoop3Map)
				}
				dataLoopMap["NodePolicy"] = localMaps2
				localMaps4 := make([]interface{}, 0)
				localData5 := dataLoopTmp["machine_type_policy"]
				for _, dataLoop5 := range convertToInterfaceArray(localData5) {
					dataLoop5Tmp := dataLoop5.(map[string]interface{})
					dataLoop5Map := make(map[string]interface{})
					localMaps5 := make([]interface{}, 0)
					localData6 := dataLoop5Tmp["bonds"]
					for _, dataLoop6 := range convertToInterfaceArray(localData6) {
						dataLoop6Tmp := dataLoop6.(map[string]interface{})
						dataLoop6Map := make(map[string]interface{})
						dataLoop6Map["Subnet"] = dataLoop6Tmp["subnet"]
						dataLoop6Map["Name"] = dataLoop6Tmp["name"]
						localMaps5 = append(localMaps5, dataLoop6Map)
					}
					dataLoop5Map["Bonds"] = localMaps5
					dataLoop5Map["MachineType"] = dataLoop5Tmp["machine_type"]
					localMaps4 = append(localMaps4, dataLoop5Map)
				}
				dataLoopMap["MachineTypePolicy"] = localMaps4
				ipAllocationPolicyMapsArray = append(ipAllocationPolicyMapsArray, dataLoopMap)
			}
			ipAllocationPolicyMapsJson, err := json.Marshal(ipAllocationPolicyMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["IpAllocationPolicy"] = string(ipAllocationPolicyMapsJson)
		}

		nodeGroupsDataList := make(map[string]interface{})

		if v, ok := d.GetOk("payment_type"); ok {
			nodeGroupsDataList["ChargeType"] = convertEfloNodeNodeGroupsArrayChargeTypeRequest(v)
		}

		if v, ok := d.GetOk("hostname"); ok {
			nodeGroupsDataList["Hostnames"] = []interface{}{v}
		}

		if v, ok := d.GetOk("vswitch_id"); ok {
			nodeGroupsDataList["VSwitchId"] = v
		}

		if v, ok := d.GetOk("node_group_id"); ok {
			nodeGroupsDataList["NodeGroupId"] = v
		}

		nodeGroupsDataList["Amount"] = "1"
		nodeGroupsDataList["Period"] = "1"

		if v, ok := d.GetOk("login_password"); ok {
			nodeGroupsDataList["LoginPassword"] = v
		}

		if v, ok := d.GetOk("user_data"); ok {
			nodeGroupsDataList["UserData"] = v
		}

		if v, ok := d.GetOk("vpc_id"); ok {
			nodeGroupsDataList["VpcId"] = v
		}

		if v, ok := d.GetOk("zone"); ok {
			nodeGroupsDataList["ZoneId"] = v
		}

		NodeGroupsMap := make([]interface{}, 0)
		NodeGroupsMap = append(NodeGroupsMap, nodeGroupsDataList)
		nodeGroupsDataListJson, err := json.Marshal(NodeGroupsMap)
		if err != nil {
			return WrapError(err)
		}
		request["NodeGroups"] = string(nodeGroupsDataListJson)

		if v, ok := d.GetOk("cluster_id"); ok {
			request["ClusterId"] = v
		}

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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_node", action, AlibabaCloudSdkGoERROR)
		}

		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Using"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEfloNodeUpdate(d, meta)
}

func resourceAliCloudEfloNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloNode(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_node DescribeEfloNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_id", objectRaw["ClusterId"])
	d.Set("computing_server", objectRaw["MachineType"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("hostname", objectRaw["Hostname"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("machine_type", objectRaw["MachineType"])
	d.Set("node_group_id", objectRaw["NodeGroupId"])
	d.Set("node_type", objectRaw["NodeType"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["OperatingState"])
	d.Set("user_data", objectRaw["UserData"])
	d.Set("zone", objectRaw["ZoneId"])
	if objectRaw["Disks"] != nil {
		disksMaps := make([]interface{}, 0)
		for _, v := range objectRaw["Disks"].([]interface{}) {
			disksMap := make(map[string]interface{})
			if fmt.Sprint(v.(map[string]interface{})["DiskType"]) == "system" {
				continue
			}
			disksMap["category"] = v.(map[string]interface{})["Category"]
			disksMap["performance_level"] = v.(map[string]interface{})["PerformanceLevel"]
			disksMap["size"] = v.(map[string]interface{})["Size"]
			disksMaps = append(disksMaps, disksMap)
		}
		d.Set("data_disk", disksMaps)
	}

	objectRaw, err = efloServiceV2.DescribeNodeListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = efloServiceV2.DescribeNodeQueryAvailableInstances(d)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("region_id", objectRaw["Region"])
	if fmt.Sprint(objectRaw["RenewalDurationUnit"]) == "Y" {
		d.Set("renew_period", formatInt(objectRaw["RenewalDuration"])*12)
	} else {
		d.Set("renew_period", objectRaw["RenewalDuration"])
	}
	d.Set("renewal_status", objectRaw["RenewStatus"])

	if d.Get("status") == "Using" {
		objectRaw, err = efloServiceV2.DescribeNodeListClusterNodes(d)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		d.Set("vswitch_id", objectRaw["VSwitchId"])
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	return nil
}

func resourceAliCloudEfloNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "Node"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NeedRetry(err) {
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
	update = false
	action = "SetRenewal"
	installPai := false
	if v, ok := d.GetOk("install_pai"); ok && v.(bool) {
		installPai = true
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		request["SubscriptionType"] = d.Get("payment_type")
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renewal_status")
	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true
		request["RenewalPeriod"] = d.Get("renew_period")
	}

	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if installPai {
		request["ProductCode"] = "learn"
		request["ProductType"] = "learn_eflocomputing_public_cn"
	}
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
		if installPai {
			return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
		}
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if installPai {
			request["ProductCode"] = "learn"
			request["ProductType"] = "learn_eflocomputing_public_intl"
		}
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
			if installPai {
				return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
			}
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v
	}
	if request["SubscriptionType"] == "" {
		request["SubscriptionType"] = "Subscription"
	}
	if request["SubscriptionType"] == "Subscription" && request["RenewalStatus"] == "AutoRenewal" {
		v, ok := d.GetOk("renew_period")
		if !ok {
			return WrapError(Error("renew_period is required when renewal_status is set to AutoRenewal."))
		}
		request["RenewalPeriod"] = v
		if v.(int) < 12 {
			request["RenewalPeriod"] = v
			request["RenewalPeriodUnit"] = "M"
		} else {
			if v.(int)%12 != 0 {
				return WrapError(Error("renew_period must be a multiple of 12 when renew_period more than 12."))
			}
			renewPeriod := v.(int) / 12
			request["RenewalPeriod"] = renewPeriod
			request["RenewalPeriodUnit"] = "Y"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_eflocomputing_public_intl"
					if installPai {
						request["ProductCode"] = "learn"
						request["ProductType"] = "learn_eflocomputing_public_intl"
					}
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "bccluster"
						request["ProductType"] = "bccluster_computinginstance_public_intl"
						if installPai {
							return resource.RetryableError(err)
						}
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
	update = false
	action = "ChangeNodeTypes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if d.HasChange("node_type") {
		update = true
		request["NodeType"] = d.Get("node_type")
	}

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "NodeIds.0", d.Id())
	_ = json.Unmarshal([]byte(jsonString), &request)

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
		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Using"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		efloServiceV2 := EfloServiceV2{client}
		if err := efloServiceV2.SetResourceTags(d, "Node"); err != nil {
			return WrapError(err)
		}
	}
	if d.IsNewResource() && d.Get("payment_type").(string) == "PayAsYouGo" {
		d.Partial(false)
		return resourceAliCloudEfloNodeRead(d, meta)
	}
	if d.HasChanges("status", "cluster_id", "node_group_id") {
		oldEntry, newEntry := d.GetChange("cluster_id")
		oldClusterId := oldEntry.(string)
		newClusterId := newEntry.(string)
		oldEntry, newEntry = d.GetChange("node_group_id")
		oldNodeGroupId := oldEntry.(string)
		newNodeGroupId := newEntry.(string)

		if fmt.Sprint(d.Get("status")) == "Unused" && oldNodeGroupId != "" {
			action := "ShrinkCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			request["ClusterId"] = oldClusterId
			nodeGroupsDataList := make(map[string]interface{})
			nodeGroupsDataList["NodeGroupId"] = oldNodeGroupId

			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, nodeGroupsDataList)
			request["NodeGroups"] = NodeGroupsMap

			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.NodeId", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)
			request["NodeGroups"] = convertObjectToJsonString(request["NodeGroups"])

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
			stateConf := BuildStateConf([]string{}, []string{"Unused"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.Set("status", "Unused")
		}

		if fmt.Sprint(d.Get("status")) == "Using" && newNodeGroupId != "" {
			action := "ExtendCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			request["ClusterId"] = newClusterId
			if v, ok := d.GetOk("ip_allocation_policy"); ok {
				ipAllocationPolicyMapsArray := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(v) {
					dataLoopTmp := dataLoop.(map[string]interface{})
					dataLoopMap := make(map[string]interface{})
					localData1 := make(map[string]interface{})
					if v, ok := dataLoopTmp["bond_policy"]; ok {
						localData2, err := jsonpath.Get("$[0].bonds", v)
						if err != nil {
							localData2 = make([]interface{}, 0)
						}
						localMaps1 := make([]interface{}, 0)
						for _, dataLoop2 := range convertToInterfaceArray(localData2) {
							dataLoop2Tmp := make(map[string]interface{})
							if dataLoop2 != nil {
								dataLoop2Tmp = dataLoop2.(map[string]interface{})
							}
							dataLoop2Map := make(map[string]interface{})
							dataLoop2Map["Subnet"] = dataLoop2Tmp["subnet"]
							dataLoop2Map["Name"] = dataLoop2Tmp["name"]
							localMaps1 = append(localMaps1, dataLoop2Map)
						}
						localData1["Bonds"] = localMaps1
					}

					bondDefaultSubnet1, _ := jsonpath.Get("$[0].bond_default_subnet", dataLoopTmp["bond_policy"])
					if bondDefaultSubnet1 != nil && bondDefaultSubnet1 != "" {
						localData1["BondDefaultSubnet"] = bondDefaultSubnet1
					}
					if len(localData1) > 0 {
						dataLoopMap["BondPolicy"] = localData1
					}
					localMaps2 := make([]interface{}, 0)
					localData3 := dataLoopTmp["node_policy"]
					for _, dataLoop3 := range convertToInterfaceArray(localData3) {
						dataLoop3Tmp := dataLoop3.(map[string]interface{})
						dataLoop3Map := make(map[string]interface{})
						localMaps3 := make([]interface{}, 0)
						localData4 := dataLoop3Tmp["bonds"]
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := dataLoop4.(map[string]interface{})
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["Name"] = dataLoop4Tmp["name"]
							dataLoop4Map["Subnet"] = dataLoop4Tmp["subnet"]
							localMaps3 = append(localMaps3, dataLoop4Map)
						}
						dataLoop3Map["Bonds"] = localMaps3
						dataLoop3Map["Hostname"] = dataLoop3Tmp["hostname"]
						dataLoop3Map["NodeId"] = dataLoop3Tmp["node_id"]
						localMaps2 = append(localMaps2, dataLoop3Map)
					}
					dataLoopMap["NodePolicy"] = localMaps2
					localMaps4 := make([]interface{}, 0)
					localData5 := dataLoopTmp["machine_type_policy"]
					for _, dataLoop5 := range convertToInterfaceArray(localData5) {
						dataLoop5Tmp := dataLoop5.(map[string]interface{})
						dataLoop5Map := make(map[string]interface{})
						localMaps5 := make([]interface{}, 0)
						localData6 := dataLoop5Tmp["bonds"]
						for _, dataLoop6 := range convertToInterfaceArray(localData6) {
							dataLoop6Tmp := dataLoop6.(map[string]interface{})
							dataLoop6Map := make(map[string]interface{})
							dataLoop6Map["Subnet"] = dataLoop6Tmp["subnet"]
							dataLoop6Map["Name"] = dataLoop6Tmp["name"]
							localMaps5 = append(localMaps5, dataLoop6Map)
						}
						dataLoop5Map["Bonds"] = localMaps5
						dataLoop5Map["MachineType"] = dataLoop5Tmp["machine_type"]
						localMaps4 = append(localMaps4, dataLoop5Map)
					}
					dataLoopMap["MachineTypePolicy"] = localMaps4
					ipAllocationPolicyMapsArray = append(ipAllocationPolicyMapsArray, dataLoopMap)
				}
				ipAllocationPolicyMapsJson, err := json.Marshal(ipAllocationPolicyMapsArray)
				if err != nil {
					return WrapError(err)
				}
				request["IpAllocationPolicy"] = string(ipAllocationPolicyMapsJson)
			}
			nodeGroupsDataList := make(map[string]interface{})
			nodeGroupsDataList["NodeGroupId"] = newNodeGroupId
			if v, ok := d.GetOk("user_data"); ok {
				nodeGroupsDataList["UserData"] = v
			}
			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, nodeGroupsDataList)
			request["NodeGroups"] = NodeGroupsMap

			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.NodeId", d.Id())
			if v, ok := d.GetOk("hostname"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.Hostname", d.Id())
			}
			if v, ok := d.GetOk("login_password"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.LoginPassword", v)
			}
			if v, ok := d.GetOk("vpc_id"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.VpcId", v)
			}
			if v, ok := d.GetOk("vswitch_id"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.VSwitchId", v)
			}
			if v, ok := d.GetOk("data_disk"); ok {
				dataDisks := make([]map[string]interface{}, 0)
				for _, disk := range v.([]interface{}) {
					diskMap := disk.(map[string]interface{})
					dataDisk := make(map[string]interface{})
					if val, ok := diskMap["category"]; ok {
						dataDisk["Category"] = val
					}
					if val, ok := diskMap["size"]; ok {
						dataDisk["Size"] = val
					}
					if val, ok := diskMap["performance_level"]; ok {
						dataDisk["PerformanceLevel"] = val
					}

					dataDisks = append(dataDisks, dataDisk)
				}
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.Nodes.0.DataDisk", dataDisks)
			}
			_ = json.Unmarshal([]byte(jsonString), &request)
			request["NodeGroups"] = convertObjectToJsonString(request["NodeGroups"])
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
			stateConf := BuildStateConf([]string{}, []string{"Using"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.Set("status", "Using")
		}
	}
	d.Partial(false)
	return resourceAliCloudEfloNodeRead(d, meta)
}

func resourceAliCloudEfloNodeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	enableDelete := false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"Subscription"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		request["ClientToken"] = buildClientToken(action)

		request["ImmediatelyRelease"] = "1"
		installPai := false
		if v, ok := d.GetOk("install_pai"); ok && v.(bool) {
			installPai = true
		}

		if !installPai {
			installPai, err = isInstallPai(d.Id(), d.Timeout(schema.TimeoutDelete), client)
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		var endpoint string
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_cn"
		if installPai {
			request["ProductCode"] = "learn"
			request["ProductType"] = "learn_eflocomputing_public_cn"
		}
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_cn"
			if installPai {
				return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
			}
		}
		if client.IsInternationalAccount() {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_eflocomputing_public_intl"
			if installPai {
				request["ProductCode"] = "learn"
				request["ProductType"] = "learn_eflocomputing_public_intl"
			}
			if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_computinginstance_public_intl"
				if installPai {
					return WrapError(Error("InstallPai currently does not support pay-as-you-go products."))
				}
			}
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_eflocomputing_public_intl"
					if installPai {
						request["ProductCode"] = "learn"
						request["ProductType"] = "learn_eflocomputing_public_intl"
					}
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "bccluster"
						request["ProductType"] = "bccluster_computinginstance_public_intl"
						if installPai {
							return resource.RetryableError(err)
						}
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"RESOURCE_NOT_FOUND"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "$.NodeId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	enableDelete = false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "DeleteNode"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["NodeId"] = d.Id()
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
			if IsExpectedErrors(err, []string{"RESOURCE_NOT_FOUND", "InvalidNodeId.NotFound"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "$.NodeId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return nil
}

func getProductCodeAndType(instanceID string, timeout time.Duration, client *connectivity.AliyunClient) (string, string, error) {
	var productCode, productType string

	request := map[string]interface{}{"InstanceIDs": instanceID}
	wait := incrementalWait(3*time.Second, 5*time.Second)

	err := resource.Retry(timeout, func() *resource.RetryError {
		response, err := client.RpcPostWithEndpoint(
			"BssOpenApi", "2017-12-14", "QueryAvailableInstances",
			map[string]interface{}{}, request, true, "",
		)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		instances, ok := getInstanceList(response)
		if !ok || len(instances) == 0 {
			return resource.NonRetryableError(WrapErrorf(NotFoundErr("Node", instanceID), NotFoundMsg, response))
		}

		instance, ok := instances[0].(map[string]interface{})
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("unexpected instance format for id %s", instanceID))
		}
		productCode, _ = instance["ProductCode"].(string)
		productType, _ = instance["ProductType"].(string)
		return nil
	})

	return productCode, productType, err
}

// getInstanceList extracts the instance list from the response.
// Returns (list, true) if found, otherwise (nil, false).
func getInstanceList(response interface{}) ([]interface{}, bool) {
	instancesI, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return nil, false
	}
	instances, ok := instancesI.([]interface{})
	return instances, ok
}

func isInstallPai(instanceId string, timeout time.Duration, client *connectivity.AliyunClient) (bool, error) {
	productCode, _, err := getProductCodeAndType(instanceId, timeout, client)
	if err != nil {
		return false, err
	}
	return productCode == "learn", nil
}

func convertEfloNodeNodeGroupsArrayChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PREPAY"
	case "PayAsYouGo":
		return "POSTPAY"
	}
	return source
}
