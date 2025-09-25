// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudApigGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigGatewayCreate,
		Read:   resourceAliCloudApigGatewayRead,
		Update: resourceAliCloudApigGatewayUpdate,
		Delete: resourceAliCloudApigGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sls": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"network_access_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Internet", "Intranet", "InternetAndIntranet"}, false),
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"zone_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"select_option": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Auto", "Manual"}, false),
						},
					},
				},
			},
			"zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudApigGatewayCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/gateways")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("vpc"); ok {
		vpcVpcIdJsonPath, err := jsonpath.Get("$[0].vpc_id", v)
		if err == nil && vpcVpcIdJsonPath != "" {
			request["vpcId"] = vpcVpcIdJsonPath
		}
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["chargeType"] = convertApigGatewaychargeTypeRequest(d.Get("payment_type").(string))
	dataList := make(map[string]interface{})

	if v := d.Get("zones"); !IsNil(v) {
		if v, ok := d.GetOk("zones"); ok {
			localData, err := jsonpath.Get("$", v)
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
				dataLoopMap["zoneId"] = dataLoopTmp["zone_id"]
				dataLoopMap["vSwitchId"] = dataLoopTmp["vswitch_id"]
				localMaps = append(localMaps, dataLoopMap)
			}
			dataList["zones"] = localMaps
		}

	}

	if v, ok := d.GetOk("zone_config"); ok {
		selectOption1, _ := jsonpath.Get("$[0].select_option", v)
		if selectOption1 != nil && selectOption1 != "" {
			dataList["selectOption"] = selectOption1
		}
	}

	if v, ok := d.GetOk("vswitch"); ok {
		vSwitchId3, _ := jsonpath.Get("$[0].vswitch_id", v)
		if vSwitchId3 != nil && vSwitchId3 != "" {
			dataList["vSwitchId"] = vSwitchId3
		}
	}

	request["zoneConfig"] = dataList

	if v, ok := d.GetOk("gateway_name"); ok {
		request["name"] = v
	}
	dataList1 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) {
		sls := make(map[string]interface{})
		enable1, _ := jsonpath.Get("$[0].sls[0].enable", v)
		if enable1 != nil && enable1 != "" {
			sls["enable"] = enable1
		}

		dataList1["sls"] = sls

		request["logConfig"] = dataList1
	}

	if v, ok := d.GetOk("spec"); ok {
		request["spec"] = v
	}
	if v, ok := d.GetOk("gateway_type"); ok {
		request["gatewayType"] = v
	}
	dataList2 := make(map[string]interface{})

	if v := d.Get("network_access_config"); !IsNil(v) {
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			dataList2["type"] = type1
		}

		request["networkAccessConfig"] = dataList2
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_gateway", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.gatewayId", response)
	d.SetId(fmt.Sprint(id))

	apigServiceV2 := ApigServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, apigServiceV2.ApigGatewayStateRefreshFunc(d.Id(), "status", []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudApigGatewayUpdate(d, meta)
}

func resourceAliCloudApigGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_gateway DescribeApigGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTimestamp"])
	d.Set("gateway_name", objectRaw["name"])
	d.Set("gateway_type", objectRaw["gatewayType"])
	d.Set("payment_type", convertApigGatewaydatachargeTypeResponse(objectRaw["chargeType"]))
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("spec", objectRaw["spec"])
	d.Set("status", objectRaw["status"])

	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	vSwitchMaps := make([]map[string]interface{}, 0)
	vSwitchMap := make(map[string]interface{})
	vSwitchRaw := make(map[string]interface{})
	if objectRaw["vSwitch"] != nil {
		vSwitchRaw = objectRaw["vSwitch"].(map[string]interface{})
	}
	if len(vSwitchRaw) > 0 {
		vSwitchMap["name"] = vSwitchRaw["name"]
		vSwitchMap["vswitch_id"] = vSwitchRaw["vSwitchId"]

		vSwitchMaps = append(vSwitchMaps, vSwitchMap)
	}
	if err := d.Set("vswitch", vSwitchMaps); err != nil {
		return err
	}
	vpcMaps := make([]map[string]interface{}, 0)
	vpcMap := make(map[string]interface{})
	vpcRaw := make(map[string]interface{})
	if objectRaw["vpc"] != nil {
		vpcRaw = objectRaw["vpc"].(map[string]interface{})
	}
	if len(vpcRaw) > 0 {
		vpcMap["name"] = vpcRaw["name"]
		vpcMap["vpc_id"] = vpcRaw["vpcId"]

		vpcMaps = append(vpcMaps, vpcMap)
	}
	if err := d.Set("vpc", vpcMaps); err != nil {
		return err
	}

	zonesRaw := objectRaw["zones"]
	zonesMaps := make([]map[string]interface{}, 0)
	if zonesRaw != nil {
		for _, zonesChildRaw := range zonesRaw.([]interface{}) {
			zonesMap := make(map[string]interface{})
			zonesChildRaw := zonesChildRaw.(map[string]interface{})
			zonesMap["name"] = zonesChildRaw["name"]
			zonesMap["zone_id"] = zonesChildRaw["zoneId"]

			if v, ok := zonesChildRaw["vSwitch"]; ok {
				vSwitchArg := v.(map[string]interface{})

				zonesMap["vswitch_id"] = vSwitchArg["vSwitchId"]
			}

			zonesMaps = append(zonesMaps, zonesMap)
		}
		if err := d.Set("zones", zonesMaps); err != nil {
			return err
		}
	}
	return nil
}

func resourceAliCloudApigGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	gatewayId := d.Id()
	action := fmt.Sprintf("/v1/gateways/%s/name", gatewayId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["gatewayId"] = d.Id()

	if d.HasChange("gateway_name") {
		update = true
	}
	if v, ok := d.GetOk("gateway_name"); ok {
		query["name"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = StringPointer(v.(string))
	}

	query["Service"] = StringPointer("APIG")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		apigServiceV2 := ApigServiceV2{client}
		if err := apigServiceV2.SetResourceTags(d, "gateway"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudApigGatewayRead(d, meta)
}

func resourceAliCloudApigGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_apig_gateway which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	gatewayId := d.Id()
	action := fmt.Sprintf("/v1/gateways/%s", gatewayId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["gatewayId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)

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
		if IsExpectedErrors(err, []string{"NotFound.GatewayNotFound", "Conflict.GatewayIsDeleted"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	apigServiceV2 := ApigServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, apigServiceV2.ApigGatewayStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertApigGatewaydatachargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
func convertApigGatewaychargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PREPAY"
	case "PayAsYouGo":
		return "POSTPAY"
	}
	return source
}
