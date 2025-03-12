package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRocketmqInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqInstanceCreate,
		Read:   resourceAliCloudRocketmqInstanceRead,
		Update: resourceAliCloudRocketmqInstanceUpdate,
		Delete: resourceAliCloudRocketmqInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(35 * time.Minute),
			Update: schema.DefaultTimeout(35 * time.Minute),
			Delete: schema.DefaultTimeout(35 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_vpc_auth_free": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"acl_types": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 6, 12}),
			},
			"auto_renew_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"commodity_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ons_rmqsub_public_cn", "ons_rmqpost_public_cn", "ons_rmqsrvlesspost_public_cn"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_whitelists": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_white_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"endpoint_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_info": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flow_out_type": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: StringInSlice([]string{"payByBandwidth", "payByTraffic", "uninvolved"}, false),
									},
									"internet_spec": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"ip_whitelist": {
										Type:       schema.TypeList,
										Optional:   true,
										Deprecated: "Field 'ip_whitelist' has been deprecated from provider version 1.245.0. New field 'ip_whitelists' instead.",
										Elem:       &schema.Schema{Type: schema.TypeString},
									},
									"flow_out_bandwidth": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"vpc_info": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"vswitch_id": {
										Type:       schema.TypeString,
										Optional:   true,
										Computed:   true,
										ForceNew:   true,
										Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.231.0. New field 'vswitches' instead.",
									},
									"vswitches": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vswitch_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
													ForceNew: true,
												},
											},
										},
									},
									"security_group_ids": {
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
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"product_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_secret_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"send_receive_ratio": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"message_retention_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"storage_encryption": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"support_auto_scaling": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"trace_on": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"msg_process_spec": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"series_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"standard", "professional", "ultimate"}, false),
			},
			"service_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"software": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"software_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgrade_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_series_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cluster_ha", "single_node", "serverless"}, false),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudRocketmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/instances")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("network_info"); v != nil {
		vpcInfo := make(map[string]interface{})
		vpcId1, _ := jsonpath.Get("$[0].vpc_info[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			vpcInfo["vpcId"] = vpcId1
		}
		vSwitchId1, _ := jsonpath.Get("$[0].vpc_info[0].vswitch_id", v)
		if vSwitchId1 != nil && vSwitchId1 != "" {
			vpcInfo["vSwitchId"] = vSwitchId1
		}
		securityGroupIds1, _ := jsonpath.Get("$[0].vpc_info[0].security_group_ids", v)
		if securityGroupIds1 != nil && securityGroupIds1 != "" {
			vpcInfo["securityGroupIds"] = securityGroupIds1
		}
		if v, ok := d.GetOk("network_info"); ok {
			localData, err := jsonpath.Get("$[0].vpc_info[0].vswitches", v)
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
				dataLoopMap["vSwitchId"] = dataLoopTmp["vswitch_id"]
				localMaps = append(localMaps, dataLoopMap)
			}
			vpcInfo["vSwitches"] = localMaps
		}

		objectDataLocalMap["vpcInfo"] = vpcInfo
		internetInfo := make(map[string]interface{})
		internetSpec1, _ := jsonpath.Get("$[0].internet_info[0].internet_spec", v)
		if internetSpec1 != nil && internetSpec1 != "" {
			internetInfo["internetSpec"] = internetSpec1
		}
		flowOutType1, _ := jsonpath.Get("$[0].internet_info[0].flow_out_type", v)
		if flowOutType1 != nil && flowOutType1 != "" {
			internetInfo["flowOutType"] = flowOutType1
		}
		flowOutBandwidth1, _ := jsonpath.Get("$[0].internet_info[0].flow_out_bandwidth", v)
		// todo: property dependent
		if flowOutBandwidth1 != nil && flowOutBandwidth1 != "" && internetInfo["internetSpec"] == "enable" {
			internetInfo["flowOutBandwidth"] = flowOutBandwidth1
		}
		ipWhitelist1, _ := jsonpath.Get("$[0].internet_info[0].ip_whitelist", v)
		if ipWhitelist1 != nil && ipWhitelist1 != "" {
			internetInfo["ipWhitelist"] = ipWhitelist1
		}

		objectDataLocalMap["internetInfo"] = internetInfo

		request["networkInfo"] = objectDataLocalMap
	}

	request["subSeriesCode"] = d.Get("sub_series_code")
	if v, ok := d.GetOk("remark"); ok {
		request["remark"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["instanceName"] = v
	}
	request["seriesCode"] = d.Get("series_code")
	request["paymentType"] = d.Get("payment_type")
	request["serviceCode"] = d.Get("service_code")
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("product_info"); !IsNil(v) {
		msgProcessSpec1, _ := jsonpath.Get("$[0].msg_process_spec", v)
		if msgProcessSpec1 != nil && msgProcessSpec1 != "" {
			objectDataLocalMap1["msgProcessSpec"] = msgProcessSpec1
		}
		sendReceiveRatio1, _ := jsonpath.Get("$[0].send_receive_ratio", v)
		if sendReceiveRatio1 != nil && sendReceiveRatio1 != "" {
			objectDataLocalMap1["sendReceiveRatio"] = sendReceiveRatio1
		}
		autoScaling1, _ := jsonpath.Get("$[0].auto_scaling", v)
		if autoScaling1 != nil && autoScaling1 != "" {
			objectDataLocalMap1["autoScaling"] = autoScaling1
		}
		messageRetentionTime1, _ := jsonpath.Get("$[0].message_retention_time", v)
		if messageRetentionTime1 != nil && messageRetentionTime1 != "" {
			objectDataLocalMap1["messageRetentionTime"] = messageRetentionTime1
		}
		storageSecretKey1, _ := jsonpath.Get("$[0].storage_secret_key", v)
		if storageSecretKey1 != nil && storageSecretKey1 != "" {
			objectDataLocalMap1["storageSecretKey"] = storageSecretKey1
		}
		storageEncryption1, _ := jsonpath.Get("$[0].storage_encryption", v)
		if storageEncryption1 != nil && storageEncryption1 != "" {
			objectDataLocalMap1["storageEncryption"] = storageEncryption1
		}

		request["productInfo"] = objectDataLocalMap1
	}

	if v, ok := d.GetOkExists("period"); ok {
		request["period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["periodUnit"] = v
	}
	if v, ok := d.GetOkExists("auto_renew_period"); ok {
		request["autoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["autoRenew"] = v
	}
	if v, ok := d.GetOk("commodity_code"); ok {
		request["commodityCode"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data", response)
	d.SetId(fmt.Sprint(id))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRocketmqInstanceUpdate(d, meta)
}

func resourceAliCloudRocketmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	objectRaw, err := rocketmqServiceV2.DescribeRocketmqInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_instance DescribeRocketmqInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("commodity_code", objectRaw["commodityCode"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("instance_name", objectRaw["instanceName"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("remark", objectRaw["remark"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("series_code", objectRaw["seriesCode"])
	d.Set("service_code", objectRaw["serviceCode"])
	d.Set("status", objectRaw["status"])
	d.Set("sub_series_code", objectRaw["subSeriesCode"])

	aclInfoMaps := make([]map[string]interface{}, 0)
	aclInfoMap := make(map[string]interface{})
	aclInfoRaw := make(map[string]interface{})
	if objectRaw["aclInfo"] != nil {
		aclInfoRaw = objectRaw["aclInfo"].(map[string]interface{})
	}
	if len(aclInfoRaw) > 0 {
		aclInfoMap["default_vpc_auth_free"] = aclInfoRaw["defaultVpcAuthFree"]

		aclTypesRaw := make([]interface{}, 0)
		if aclInfoRaw["aclTypes"] != nil {
			aclTypesRaw = aclInfoRaw["aclTypes"].([]interface{})
		}

		aclInfoMap["acl_types"] = aclTypesRaw
		aclInfoMaps = append(aclInfoMaps, aclInfoMap)
	}
	if err := d.Set("acl_info", aclInfoMaps); err != nil {
		return err
	}
	networkInfoMaps := make([]map[string]interface{}, 0)
	networkInfoMap := make(map[string]interface{})
	networkInfoRaw := make(map[string]interface{})
	if objectRaw["networkInfo"] != nil {
		networkInfoRaw = objectRaw["networkInfo"].(map[string]interface{})
	}
	if len(networkInfoRaw) > 0 {

		endpointsRaw := networkInfoRaw["endpoints"]
		endpointsMaps := make([]map[string]interface{}, 0)
		if endpointsRaw != nil {
			for _, endpointsChildRaw := range endpointsRaw.([]interface{}) {
				endpointsMap := make(map[string]interface{})
				endpointsChildRaw := endpointsChildRaw.(map[string]interface{})
				endpointsMap["endpoint_type"] = endpointsChildRaw["endpointType"]
				endpointsMap["endpoint_url"] = endpointsChildRaw["endpointUrl"]

				ipWhitelistRaw := make([]interface{}, 0)
				if endpointsChildRaw["ipWhitelist"] != nil {
					ipWhitelistRaw = endpointsChildRaw["ipWhitelist"].([]interface{})
				}

				endpointsMap["ip_white_list"] = ipWhitelistRaw
				endpointsMaps = append(endpointsMaps, endpointsMap)
			}
		}
		networkInfoMap["endpoints"] = endpointsMaps
		internetInfoMaps := make([]map[string]interface{}, 0)
		internetInfoMap := make(map[string]interface{})
		internetInfoRaw := make(map[string]interface{})
		if networkInfoRaw["internetInfo"] != nil {
			internetInfoRaw = networkInfoRaw["internetInfo"].(map[string]interface{})
		}
		if len(internetInfoRaw) > 0 {
			internetInfoMap["flow_out_bandwidth"] = internetInfoRaw["flowOutBandwidth"]
			internetInfoMap["flow_out_type"] = internetInfoRaw["flowOutType"]
			internetInfoMap["internet_spec"] = internetInfoRaw["internetSpec"]

			ipWhitelistRaw := make([]interface{}, 0)
			if internetInfoRaw["ipWhitelist"] != nil {
				ipWhitelistRaw = internetInfoRaw["ipWhitelist"].([]interface{})
			}

			internetInfoMap["ip_whitelist"] = ipWhitelistRaw
			internetInfoMaps = append(internetInfoMaps, internetInfoMap)
		}
		networkInfoMap["internet_info"] = internetInfoMaps
		vpcInfoMaps := make([]map[string]interface{}, 0)
		vpcInfoMap := make(map[string]interface{})
		vpcInfoRaw := make(map[string]interface{})
		if networkInfoRaw["vpcInfo"] != nil {
			vpcInfoRaw = networkInfoRaw["vpcInfo"].(map[string]interface{})
		}
		if len(vpcInfoRaw) > 0 {
			vpcInfoMap["security_group_ids"] = vpcInfoRaw["securityGroupIds"]
			vpcInfoMap["vswitch_id"] = vpcInfoRaw["vSwitchId"]
			vpcInfoMap["vpc_id"] = vpcInfoRaw["vpcId"]

			vSwitchesRaw := vpcInfoRaw["vSwitches"]
			vSwitchesMaps := make([]map[string]interface{}, 0)
			if vSwitchesRaw != nil {
				for _, vSwitchesChildRaw := range vSwitchesRaw.([]interface{}) {
					vSwitchesMap := make(map[string]interface{})
					vSwitchesChildRaw := vSwitchesChildRaw.(map[string]interface{})
					vSwitchesMap["vswitch_id"] = vSwitchesChildRaw["vSwitchId"]

					vSwitchesMaps = append(vSwitchesMaps, vSwitchesMap)
				}
			}
			vpcInfoMap["vswitches"] = vSwitchesMaps
			vpcInfoMaps = append(vpcInfoMaps, vpcInfoMap)
		}
		networkInfoMap["vpc_info"] = vpcInfoMaps
		networkInfoMaps = append(networkInfoMaps, networkInfoMap)
	}
	if err := d.Set("network_info", networkInfoMaps); err != nil {
		return err
	}
	productInfoMaps := make([]map[string]interface{}, 0)
	productInfoMap := make(map[string]interface{})
	productInfoRaw := make(map[string]interface{})
	if objectRaw["productInfo"] != nil {
		productInfoRaw = objectRaw["productInfo"].(map[string]interface{})
	}
	if len(productInfoRaw) > 0 {
		productInfoMap["auto_scaling"] = productInfoRaw["autoScaling"]
		productInfoMap["message_retention_time"] = productInfoRaw["messageRetentionTime"]
		productInfoMap["msg_process_spec"] = productInfoRaw["msgProcessSpec"]
		productInfoMap["send_receive_ratio"] = productInfoRaw["sendReceiveRatio"]
		productInfoMap["storage_encryption"] = productInfoRaw["storageEncryption"]
		productInfoMap["storage_secret_key"] = productInfoRaw["storageSecretKey"]
		productInfoMap["support_auto_scaling"] = productInfoRaw["supportAutoScaling"]
		productInfoMap["trace_on"] = productInfoRaw["traceOn"]

		productInfoMaps = append(productInfoMaps, productInfoMap)
	}
	if err := d.Set("product_info", productInfoMaps); err != nil {
		return err
	}
	softwareMaps := make([]map[string]interface{}, 0)
	softwareMap := make(map[string]interface{})
	softwareRaw := make(map[string]interface{})
	if objectRaw["software"] != nil {
		softwareRaw = objectRaw["software"].(map[string]interface{})
	}
	if len(softwareRaw) > 0 {
		softwareMap["maintain_time"] = softwareRaw["maintainTime"]
		softwareMap["software_version"] = softwareRaw["softwareVersion"]
		softwareMap["upgrade_method"] = softwareRaw["upgradeMethod"]

		softwareMaps = append(softwareMaps, softwareMap)
	}
	if err := d.Set("software", softwareMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = rocketmqServiceV2.DescribeGetInstanceAccount(d.Id())
	if err != nil {
		return WrapError(err)
	}

	bssOpenApiService := BssOpenApiService{client}
	queryAvailableInstancesObject, err := bssOpenApiService.QueryAvailableInstancesWithoutProductType(d.Id(), client.RegionId, "ons", "ons")
	if err != nil {
		return WrapError(err)
	}
	d.Set("auto_renew", queryAvailableInstancesObject["RenewStatus"] == "AutoRenewal")
	if v, ok := queryAvailableInstancesObject["RenewalDuration"]; ok && fmt.Sprint(v) != "0" {
		d.Set("auto_renew_period", formatInt(v))
	}
	d.Set("auto_renew_period_unit", convertAmqpInstanceRenewalDurationUnitResponse(queryAvailableInstancesObject["RenewalDurationUnit"]))

	objectRaw, err = rocketmqServiceV2.DescribeInstanceGetInstanceIpWhitelist(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	ipWhitelistsRaw := make([]interface{}, 0)
	if objectRaw["ipWhitelists"] != nil {
		ipWhitelistsRaw = objectRaw["ipWhitelists"].([]interface{})
	}

	d.Set("ip_whitelists", ipWhitelistsRaw)

	return nil
}

func resourceAliCloudRocketmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	instanceId := d.Id()
	action := fmt.Sprintf("/instances/%s", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok || d.HasChange("instance_name") {
		request["instanceName"] = v
	}
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
	}
	if v, ok := d.GetOk("remark"); ok || d.HasChange("remark") {
		request["remark"] = v
	}
	if !d.IsNewResource() && d.HasChange("product_info") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("product_info"); v != nil {
		sendReceiveRatio1, _ := jsonpath.Get("$[0].send_receive_ratio", v)
		if sendReceiveRatio1 != nil && (d.HasChange("product_info.0.send_receive_ratio") || sendReceiveRatio1 != "") {
			objectDataLocalMap["sendReceiveRatio"] = sendReceiveRatio1
		}
		autoScaling1, _ := jsonpath.Get("$[0].auto_scaling", v)
		if autoScaling1 != nil && (d.HasChange("product_info.0.auto_scaling") || autoScaling1 != "") {
			objectDataLocalMap["autoScaling"] = autoScaling1
		}
		messageRetentionTime1, _ := jsonpath.Get("$[0].message_retention_time", v)
		if messageRetentionTime1 != nil && (d.HasChange("product_info.0.message_retention_time") || messageRetentionTime1 != "") {
			objectDataLocalMap["messageRetentionTime"] = messageRetentionTime1
		}
		traceOn1, _ := jsonpath.Get("$[0].trace_on", v)
		if traceOn1 != nil && (d.HasChange("product_info.0.trace_on") || traceOn1 != "") {
			update = true
			objectDataLocalMap["traceOn"] = traceOn1
		}

		request["productInfo"] = objectDataLocalMap
	}

	if !d.IsNewResource() && d.HasChange("network_info") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("network_info"); v != nil {
		internetInfo := make(map[string]interface{})
		ipWhitelist1, _ := jsonpath.Get("$[0].internet_info[0].ip_whitelist", v)
		if ipWhitelist1 != nil && (d.HasChange("network_info.0.internet_info.0.ip_whitelist") || ipWhitelist1 != "") {
			internetInfo["ipWhitelist"] = ipWhitelist1
		}

		objectDataLocalMap1["internetInfo"] = internetInfo

		request["networkInfo"] = objectDataLocalMap1
	}

	if d.HasChange("acl_info") {
		update = true
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("acl_info"); v != nil {
		aclTypes1, _ := jsonpath.Get("$[0].acl_types", d.Get("acl_info"))
		if aclTypes1 != nil && (d.HasChange("acl_info.0.acl_types") || aclTypes1 != "") {
			objectDataLocalMap2["aclTypes"] = aclTypes1
		}
		defaultVpcAuthFree1, _ := jsonpath.Get("$[0].default_vpc_auth_free", v)
		if defaultVpcAuthFree1 != nil && (d.HasChange("acl_info.0.default_vpc_auth_free") || defaultVpcAuthFree1 != "") {
			objectDataLocalMap2["defaultVpcAuthFree"] = defaultVpcAuthFree1
		}

		request["aclInfo"] = objectDataLocalMap2
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
		rocketmqServiceV2 := RocketmqServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = fmt.Sprintf("/resourceGroup/change")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["resourceId"] = StringPointer(d.Id())
	query["regionId"] = StringPointer(client.RegionId)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	query["resourceType"] = StringPointer("instance")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/instances/%s/software/config", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("software.0.maintain_time") {
		update = true
	}
	if v, ok := d.GetOk("software"); ok || d.HasChange("software") {
		jsonPathResult, err := jsonpath.Get("$[0].maintain_time", v)
		if err == nil && jsonPathResult != "" {
			request["maintainTime"] = jsonPathResult
		}
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, body, true)
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

	update = false
	request = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if !d.IsNewResource() && d.HasChange("auto_renew") {
		update = true
	}
	request["RenewalStatus"] = "ManualRenewal"
	if v, ok := d.GetOk("auto_renew"); ok {
		if v.(bool) {
			request["RenewalStatus"] = "AutoRenewal"
		}
	}
	if !d.IsNewResource() && d.HasChange("auto_renew_period") {
		update = true
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["RenewalPeriod"] = v
	}
	if d.HasChange("auto_renew_period_unit") {
		update = true
	}
	if v, ok := d.GetOk("auto_renew_period_unit"); ok {
		request["RenewalPeriodUnit"] = convertAmqpInstanceRenewalDurationUnitRequest(v.(string))
	}

	request["ProductCode"] = "ons"
	if update {
		action := "SetRenewal"
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, false, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("ip_whitelists") {
		oldEntry, newEntry := d.GetChange("ip_whitelists")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			instanceId := d.Id()
			action := fmt.Sprintf("/instances/%s/ip/whitelist", instanceId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["instanceId"] = StringPointer(d.Id())

			localData := removed.List()
			ipWhitelistsMapsArray := localData
			query["ipWhitelists"] = StringPointer(convertListToCommaSeparate(ipWhitelistsMapsArray))

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, true)
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

		if added.Len() > 0 {
			instanceId := d.Id()
			action := fmt.Sprintf("/instances/%s/ip/whitelist", instanceId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["instanceId"] = d.Id()

			localData := added.List()
			ipWhitelistsMapsArray := localData
			request["ipWhitelists"] = ipWhitelistsMapsArray

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
	}

	if d.HasChange("tags") {
		rocketmqServiceV2 := RocketmqServiceV2{client}
		if err := rocketmqServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudRocketmqInstanceRead(d, meta)
}

func resourceAliCloudRocketmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if d.Get("payment_type").(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_rocketmq_instance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	instanceId := d.Id()
	action := fmt.Sprintf("/instances/%s", instanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["instanceId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, true)

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

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "$.instanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
