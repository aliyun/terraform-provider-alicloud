// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: StringInSlice([]string{"enable", "disable"}, false),
									},
									"ip_whitelist": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"flow_out_bandwidth": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: IntBetween(1, 1000),
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
										Required: true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6}),
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
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: FloatBetween(.2, .5),
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
				ValidateFunc: StringInSlice([]string{"cluster_ha", "single_node"}, false),
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
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("network_info"); !IsNil(v) {
		vpcInfo := make(map[string]interface{})
		nodeNative, _ := jsonpath.Get("$[0].vpc_info[0].vpc_id", v)
		if nodeNative != "" {
			vpcInfo["vpcId"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].vpc_info[0].vswitch_id", v)
		if nodeNative1 != "" {
			vpcInfo["vSwitchId"] = nodeNative1
		}
		objectDataLocalMap["vpcInfo"] = vpcInfo
		internetInfo := make(map[string]interface{})
		nodeNative2, _ := jsonpath.Get("$[0].internet_info[0].internet_spec", v)
		if nodeNative2 != "" {
			internetInfo["internetSpec"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].internet_info[0].flow_out_type", v)
		if nodeNative3 != "" {
			internetInfo["flowOutType"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].internet_info[0].flow_out_bandwidth", v)
		if nodeNative4 != "" {
			internetInfo["flowOutBandwidth"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].internet_info[0].ip_whitelist", v)
		if nodeNative5 != "" {
			internetInfo["ipWhitelist"] = nodeNative5
		}
		objectDataLocalMap["internetInfo"] = internetInfo
	}
	request["networkInfo"] = objectDataLocalMap

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
		nodeNative6, _ := jsonpath.Get("$[0].msg_process_spec", v)
		if nodeNative6 != "" {
			objectDataLocalMap1["msgProcessSpec"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].send_receive_ratio", v)
		if nodeNative7 != "" {
			objectDataLocalMap1["sendReceiveRatio"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].auto_scaling", v)
		if nodeNative8 != "" {
			objectDataLocalMap1["autoScaling"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].message_retention_time", v)
		if nodeNative9 != "" {
			objectDataLocalMap1["messageRetentionTime"] = nodeNative9
		}
	}
	request["productInfo"] = objectDataLocalMap1

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
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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

	d.SetId(fmt.Sprint(response["body"].(map[string]interface{})["data"]))

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

	d.Set("create_time", objectRaw["createTime"])
	d.Set("instance_name", objectRaw["instanceName"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("remark", objectRaw["remark"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("series_code", objectRaw["seriesCode"])
	d.Set("service_code", objectRaw["serviceCode"])
	d.Set("status", objectRaw["status"])
	d.Set("sub_series_code", objectRaw["subSeriesCode"])

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
			vpcInfoMap["vswitch_id"] = vpcInfo1Raw["vSwitchId"]
			vpcInfoMap["vpc_id"] = vpcInfo1Raw["vpcId"]
			vpcInfoMaps = append(vpcInfoMaps, vpcInfoMap)
		}
		networkInfoMap["vpc_info"] = vpcInfoMaps
		networkInfoMaps = append(networkInfoMaps, networkInfoMap)
	}
	d.Set("network_info", networkInfoMaps)
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
	d.Set("product_info", productInfoMaps)
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
	d.Set("software", softwareMaps)
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
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)
	instanceId := d.Id()
	action := fmt.Sprintf("/instances/%s", instanceId)
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["instanceName"] = d.Get("instance_name")
	}

	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
	}
	request["remark"] = d.Get("remark")

	if !d.IsNewResource() && d.HasChange("product_info") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("product_info"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].send_receive_ratio", v)
		if nodeNative != "" {
			objectDataLocalMap["sendReceiveRatio"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].auto_scaling", v)
		if nodeNative1 != "" {
			objectDataLocalMap["autoScaling"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].message_retention_time", v)
		if nodeNative2 != "" {
			objectDataLocalMap["messageRetentionTime"] = nodeNative2
		}
	}
	request["productInfo"] = objectDataLocalMap

	if !d.IsNewResource() && d.HasChange("network_info") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("network_info"); !IsNil(v) {
		internetInfo := make(map[string]interface{})
		nodeNative3, _ := jsonpath.Get("$[0].internet_info[0].ip_whitelist", v)
		if nodeNative3 != "" {
			internetInfo["ipWhitelist"] = nodeNative3
		}
		objectDataLocalMap1["internetInfo"] = internetInfo
	}
	request["networkInfo"] = objectDataLocalMap1

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("PATCH"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		d.SetPartial("instance_name")
		d.SetPartial("remark")
	}
	update = false
	action = fmt.Sprintf("/resourceGroup/change")
	conn, err = client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["regionId"] = StringPointer(client.RegionId)
	query["resourceId"] = StringPointer(d.Id())
	if v, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		query["resourceGroupId"] = StringPointer(v.(string))
	}
	query["resourceType"] = StringPointer("instance")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		d.SetPartial("resource_group_id")
	}
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/instances/%s/software/config", instanceId)
	conn, err = client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("software") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].maintain_time", d.Get("software"))
		if err != nil {
			return WrapError(err)
		}
		request["maintainTime"] = jsonPathResult
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("PATCH"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		d.SetPartial("maintain_time")
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
		client := meta.(*connectivity.AliyunClient)
		conn, err := client.NewBssopenapiClient()
		action := "SetRenewal"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable", "SignatureDoesNotMatch"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					request["ProductCode"] = "ons"
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
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
		d.SetPartial("auto_renew_period_unit")
	}

	if d.HasChange("tags") {
		rocketmqServiceV2 := RocketmqServiceV2{client}
		if err := rocketmqServiceV2.SetResourceTags(d, ""); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
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
	body := make(map[string]interface{})
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["instanceId"] = d.Id()

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, rocketmqServiceV2.RocketmqInstanceStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
