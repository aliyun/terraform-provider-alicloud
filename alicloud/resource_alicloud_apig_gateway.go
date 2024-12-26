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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	if v, ok := d.GetOk("spec"); ok {
		request["spec"] = v
	}
	if v, ok := d.GetOk("gateway_name"); ok {
		request["name"] = v
	}
	if v, ok := d.GetOk("vpc"); ok {
		jsonPathResult2, err := jsonpath.Get("$[0].vpc_id", v)
		if err == nil && jsonPathResult2 != "" {
			request["vpcId"] = jsonPathResult2
		}
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("network_access_config"); !IsNil(v) {
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			objectDataLocalMap["type"] = type1
		}

		request["networkAccessConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) {
		sls := make(map[string]interface{})
		enable1, _ := jsonpath.Get("$[0].sls[0].enable", v)
		if enable1 != nil && enable1 != "" {
			sls["enable"] = enable1
		}

		objectDataLocalMap1["sls"] = sls

		request["logConfig"] = objectDataLocalMap1
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if v, ok := d.GetOk("zone_config"); ok {
		selectOption1, _ := jsonpath.Get("$[0].select_option", v)
		if selectOption1 != nil && selectOption1 != "" {
			objectDataLocalMap2["selectOption"] = selectOption1
		}
	}

	if v, ok := d.GetOk("vswitch"); ok {
		vSwitchId1, _ := jsonpath.Get("$[0].vswitch_id", v)
		if vSwitchId1 != nil && vSwitchId1 != "" {
			objectDataLocalMap2["vSwitchId"] = vSwitchId1
		}
	}

	request["zoneConfig"] = objectDataLocalMap2
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["chargeType"] = convertApigGatewaychargeTypeRequest(d.Get("payment_type").(string))
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	id, _ := jsonpath.Get("$.body.data.gatewayId", response)
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

	if objectRaw["createTimestamp"] != nil {
		d.Set("create_time", objectRaw["createTimestamp"])
	}
	if objectRaw["name"] != nil {
		d.Set("gateway_name", objectRaw["name"])
	}
	if objectRaw["chargeType"] != nil {
		d.Set("payment_type", convertApigGatewaydatachargeTypeResponse(objectRaw["chargeType"]))
	}
	if objectRaw["resourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["resourceGroupId"])
	}
	if objectRaw["spec"] != nil {
		d.Set("spec", objectRaw["spec"])
	}
	if objectRaw["status"] != nil {
		d.Set("status", objectRaw["status"])
	}

	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	vSwitchMaps := make([]map[string]interface{}, 0)
	vSwitchMap := make(map[string]interface{})
	vSwitch1Raw := make(map[string]interface{})
	if objectRaw["vSwitch"] != nil {
		vSwitch1Raw = objectRaw["vSwitch"].(map[string]interface{})
	}
	if len(vSwitch1Raw) > 0 {
		vSwitchMap["name"] = vSwitch1Raw["name"]
		vSwitchMap["vswitch_id"] = vSwitch1Raw["vSwitchId"]

		vSwitchMaps = append(vSwitchMaps, vSwitchMap)
	}
	if objectRaw["vSwitch"] != nil {
		if err := d.Set("vswitch", vSwitchMaps); err != nil {
			return err
		}
	}
	vpcMaps := make([]map[string]interface{}, 0)
	vpcMap := make(map[string]interface{})
	vpc1Raw := make(map[string]interface{})
	if objectRaw["vpc"] != nil {
		vpc1Raw = objectRaw["vpc"].(map[string]interface{})
	}
	if len(vpc1Raw) > 0 {
		vpcMap["name"] = vpc1Raw["name"]
		vpcMap["vpc_id"] = vpc1Raw["vpcId"]

		vpcMaps = append(vpcMaps, vpcMap)
	}
	if objectRaw["vpc"] != nil {
		if err := d.Set("vpc", vpcMaps); err != nil {
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

	gatewayId := d.Id()
	action := fmt.Sprintf("/v1/gateways/%s/name", gatewayId)
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err = client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
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

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["gatewayId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
