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
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
						"send_receive_ratio": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"message_retention_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"support_auto_scaling": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"msg_process_spec": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
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
	var err error
	query := make(map[string]*string)
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("network_info"); v != nil {
		vpcInfo := make(map[string]interface{})
		vpcId1, _ := jsonpath.Get("$[0].vpc_info[0].vpc_id", d.Get("network_info"))
		if vpcId1 != nil && vpcId1 != "" {
			vpcInfo["vpcId"] = vpcId1
		}
		vSwitchId1, _ := jsonpath.Get("$[0].vpc_info[0].vswitch_id", d.Get("network_info"))
		if vSwitchId1 != nil && vSwitchId1 != "" {
			vpcInfo["vSwitchId"] = vSwitchId1
		}
		securityGroupIds1, _ := jsonpath.Get("$[0].vpc_info[0].security_group_ids", d.Get("network_info"))
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
		internetSpec1, _ := jsonpath.Get("$[0].internet_info[0].internet_spec", d.Get("network_info"))
		if internetSpec1 != nil && internetSpec1 != "" {
			internetInfo["internetSpec"] = internetSpec1
		}
		flowOutType1, _ := jsonpath.Get("$[0].internet_info[0].flow_out_type", d.Get("network_info"))
		if flowOutType1 != nil && flowOutType1 != "" {
			internetInfo["flowOutType"] = flowOutType1
		}
		flowOutBandwidth1, _ := jsonpath.Get("$[0].internet_info[0].flow_out_bandwidth", d.Get("network_info"))
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
		msgProcessSpec1, _ := jsonpath.Get("$[0].msg_process_spec", d.Get("product_info"))
		if msgProcessSpec1 != nil && msgProcessSpec1 != "" {
			objectDataLocalMap1["msgProcessSpec"] = msgProcessSpec1
		}
		sendReceiveRatio1, _ := jsonpath.Get("$[0].send_receive_ratio", d.Get("product_info"))
		if sendReceiveRatio1 != nil && sendReceiveRatio1 != "" {
			objectDataLocalMap1["sendReceiveRatio"] = sendReceiveRatio1
		}
		autoScaling1, _ := jsonpath.Get("$[0].auto_scaling", d.Get("product_info"))
		if autoScaling1 != nil && autoScaling1 != "" {
			objectDataLocalMap1["autoScaling"] = autoScaling1
		}
		messageRetentionTime1, _ := jsonpath.Get("$[0].message_retention_time", d.Get("product_info"))
		if messageRetentionTime1 != nil && messageRetentionTime1 != "" {
			objectDataLocalMap1["messageRetentionTime"] = messageRetentionTime1
		}

		request["productInfo"] = objectDataLocalMap1
	}

	if v, ok := d.GetOk("period"); ok {
		request["period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["periodUnit"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.body.data", response)
	d.SetId(fmt.Sprint(id))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 6*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "status", []string{}))
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

	if objectRaw["commodityCode"] != nil {
		d.Set("commodity_code", objectRaw["commodityCode"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["instanceName"] != nil {
		d.Set("instance_name", objectRaw["instanceName"])
	}
	if objectRaw["paymentType"] != nil {
		d.Set("payment_type", objectRaw["paymentType"])
	}
	if objectRaw["remark"] != nil {
		d.Set("remark", objectRaw["remark"])
	}
	if objectRaw["resourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["resourceGroupId"])
	}
	if objectRaw["seriesCode"] != nil {
		d.Set("series_code", objectRaw["seriesCode"])
	}
	if objectRaw["serviceCode"] != nil {
		d.Set("service_code", objectRaw["serviceCode"])
	}
	if objectRaw["status"] != nil {
		d.Set("status", objectRaw["status"])
	}
	if objectRaw["subSeriesCode"] != nil {
		d.Set("sub_series_code", objectRaw["subSeriesCode"])
	}

	networkInfoMaps := make([]map[string]interface{}, 0)
	networkInfoMap := make(map[string]interface{})
	networkInfo1Raw := make(map[string]interface{})
	if objectRaw["networkInfo"] != nil {
		networkInfo1Raw = objectRaw["networkInfo"].(map[string]interface{})
	}
	if len(networkInfo1Raw) > 0 {

		endpoints1Raw := networkInfo1Raw["endpoints"]
		endpointsMaps := make([]map[string]interface{}, 0)
		if endpoints1Raw != nil {
			for _, endpointsChild1Raw := range endpoints1Raw.([]interface{}) {
				endpointsMap := make(map[string]interface{})
				endpointsChild1Raw := endpointsChild1Raw.(map[string]interface{})
				endpointsMap["endpoint_type"] = endpointsChild1Raw["endpointType"]
				endpointsMap["endpoint_url"] = endpointsChild1Raw["endpointUrl"]

				ipWhitelist2Raw := make([]interface{}, 0)
				if endpointsChild1Raw["ipWhitelist"] != nil {
					ipWhitelist2Raw = endpointsChild1Raw["ipWhitelist"].([]interface{})
				}

				endpointsMap["ip_white_list"] = ipWhitelist2Raw
				endpointsMaps = append(endpointsMaps, endpointsMap)
			}
		}
		networkInfoMap["endpoints"] = endpointsMaps
		internetInfoMaps := make([]map[string]interface{}, 0)
		internetInfoMap := make(map[string]interface{})
		internetInfo1Raw := make(map[string]interface{})
		if networkInfo1Raw["internetInfo"] != nil {
			internetInfo1Raw = networkInfo1Raw["internetInfo"].(map[string]interface{})
		}
		if len(internetInfo1Raw) > 0 {
			internetInfoMap["flow_out_bandwidth"] = internetInfo1Raw["flowOutBandwidth"]
			internetInfoMap["flow_out_type"] = internetInfo1Raw["flowOutType"]
			internetInfoMap["internet_spec"] = internetInfo1Raw["internetSpec"]

			ipWhitelist3Raw := make([]interface{}, 0)
			if internetInfo1Raw["ipWhitelist"] != nil {
				ipWhitelist3Raw = internetInfo1Raw["ipWhitelist"].([]interface{})
			}

			internetInfoMap["ip_whitelist"] = ipWhitelist3Raw
			internetInfoMaps = append(internetInfoMaps, internetInfoMap)
		}
		networkInfoMap["internet_info"] = internetInfoMaps
		vpcInfoMaps := make([]map[string]interface{}, 0)
		vpcInfoMap := make(map[string]interface{})
		vpcInfo1Raw := make(map[string]interface{})
		if networkInfo1Raw["vpcInfo"] != nil {
			vpcInfo1Raw = networkInfo1Raw["vpcInfo"].(map[string]interface{})
		}
		if len(vpcInfo1Raw) > 0 {
			vpcInfoMap["security_group_ids"] = vpcInfo1Raw["securityGroupIds"]
			vpcInfoMap["vswitch_id"] = vpcInfo1Raw["vSwitchId"]
			vpcInfoMap["vpc_id"] = vpcInfo1Raw["vpcId"]

			vSwitches1Raw := vpcInfo1Raw["vSwitches"]
			vSwitchesMaps := make([]map[string]interface{}, 0)
			if vSwitches1Raw != nil {
				for _, vSwitchesChild1Raw := range vSwitches1Raw.([]interface{}) {
					vSwitchesMap := make(map[string]interface{})
					vSwitchesChild1Raw := vSwitchesChild1Raw.(map[string]interface{})
					vSwitchesMap["vswitch_id"] = vSwitchesChild1Raw["vSwitchId"]

					vSwitchesMaps = append(vSwitchesMaps, vSwitchesMap)
				}
			}
			vpcInfoMap["vswitches"] = vSwitchesMaps
			vpcInfoMaps = append(vpcInfoMaps, vpcInfoMap)
		}
		networkInfoMap["vpc_info"] = vpcInfoMaps
		networkInfoMaps = append(networkInfoMaps, networkInfoMap)
	}
	if objectRaw["networkInfo"] != nil {
		if err := d.Set("network_info", networkInfoMaps); err != nil {
			return err
		}
	}
	productInfoMaps := make([]map[string]interface{}, 0)
	productInfoMap := make(map[string]interface{})
	productInfo1Raw := make(map[string]interface{})
	if objectRaw["productInfo"] != nil {
		productInfo1Raw = objectRaw["productInfo"].(map[string]interface{})
	}
	if len(productInfo1Raw) > 0 {
		productInfoMap["auto_scaling"] = productInfo1Raw["autoScaling"]
		productInfoMap["message_retention_time"] = productInfo1Raw["messageRetentionTime"]
		productInfoMap["msg_process_spec"] = productInfo1Raw["msgProcessSpec"]
		productInfoMap["send_receive_ratio"] = productInfo1Raw["sendReceiveRatio"]
		productInfoMap["support_auto_scaling"] = productInfo1Raw["supportAutoScaling"]

		productInfoMaps = append(productInfoMaps, productInfoMap)
	}
	if objectRaw["productInfo"] != nil {
		if err := d.Set("product_info", productInfoMaps); err != nil {
			return err
		}
	}
	softwareMaps := make([]map[string]interface{}, 0)
	softwareMap := make(map[string]interface{})
	software1Raw := make(map[string]interface{})
	if objectRaw["software"] != nil {
		software1Raw = objectRaw["software"].(map[string]interface{})
	}
	if len(software1Raw) > 0 {
		softwareMap["maintain_time"] = software1Raw["maintainTime"]
		softwareMap["software_version"] = software1Raw["softwareVersion"]
		softwareMap["upgrade_method"] = software1Raw["upgradeMethod"]

		softwareMaps = append(softwareMaps, softwareMap)
	}
	if objectRaw["software"] != nil {
		if err := d.Set("software", softwareMaps); err != nil {
			return err
		}
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
	return nil
}

func resourceAliCloudRocketmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var query map[string]*string
	update := false
	d.Partial(true)
	instanceId := d.Id()
	action := fmt.Sprintf("/instances/%s", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
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

	if v := d.Get("product_info"); !IsNil(v) || d.HasChange("product_info") {
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

		request["productInfo"] = objectDataLocalMap
	}

	if !d.IsNewResource() && d.HasChange("network_info") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("network_info"); v != nil {
		internetInfo := make(map[string]interface{})
		ipWhitelist1, _ := jsonpath.Get("$[0].internet_info[0].ip_whitelist", d.Get("network_info"))
		if ipWhitelist1 != nil && (d.HasChange("network_info.0.internet_info.0.ip_whitelist") || ipWhitelist1 != "") {
			internetInfo["ipWhitelist"] = ipWhitelist1
		}

		objectDataLocalMap1["internetInfo"] = internetInfo

		request["networkInfo"] = objectDataLocalMap1
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, request, false)
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
	query["resourceId"] = StringPointer(d.Id())
	query["regionId"] = StringPointer(client.RegionId)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}
	query["resourceType"] = StringPointer("instance")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, request, false)
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
	}
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/instances/%s/software/config", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
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
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, request, false)
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

	if d.HasChange("tags") {
		rocketmqServiceV2 := RocketmqServiceV2{client}
		if err := rocketmqServiceV2.SetResourceTags(d, ""); err != nil {
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
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
