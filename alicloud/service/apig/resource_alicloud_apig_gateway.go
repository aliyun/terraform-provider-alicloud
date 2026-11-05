// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceGatewayCreate,
		Read:   resourceGatewayRead,
		Update: resourceGatewayUpdate,
		Delete: resourceGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gateway_edition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"load_balancers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"address_ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
							ValidateFunc: validators.StringInSlice([]string{"Internet", "Intranet", "InternetAndIntranet"}, false),
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"tags": tags.TagsSchema(),
			"target_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
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
			"version": {
				Type:     schema.TypeString,
				Computed: true,
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
							ValidateFunc: validators.StringInSlice([]string{"Auto", "Manual"}, false),
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

func resourceGatewayCreate(d *schema.ResourceData, meta interface{}) error {

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
	zoneConfig := make(map[string]interface{})

	if v := d.Get("zones"); !helper.IsNil(v) {
		if v, ok := d.GetOk("zones"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convert.ConvertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["zoneId"] = dataLoopTmp["zone_id"]
				dataLoopMap["vSwitchId"] = dataLoopTmp["vswitch_id"]
				localMaps = append(localMaps, dataLoopMap)
			}
			zoneConfig["zones"] = localMaps
		}

	}

	if v, ok := d.GetOk("zone_config"); ok {
		selectOption1, _ := jsonpath.Get("$[0].select_option", v)
		if selectOption1 != nil && selectOption1 != "" {
			zoneConfig["selectOption"] = selectOption1
		}
	}

	if v, ok := d.GetOk("vswitch"); ok {
		vSwitchId3, _ := jsonpath.Get("$[0].vswitch_id", v)
		if vSwitchId3 != nil && vSwitchId3 != "" {
			zoneConfig["vSwitchId"] = vSwitchId3
		}
	}

	request["zoneConfig"] = zoneConfig

	if v, ok := d.GetOk("gateway_edition"); ok {
		request["gatewayEdition"] = v
	}
	if v, ok := d.GetOk("gateway_name"); ok {
		request["name"] = v
	}
	logConfig := make(map[string]interface{})

	if v := d.Get("log_config"); !helper.IsNil(v) {
		sls := make(map[string]interface{})
		enable1, _ := jsonpath.Get("$[0].sls[0].enable", d.Get("log_config"))
		if enable1 != nil && enable1 != "" {
			sls["enable"] = enable1
		}

		if len(sls) > 0 {
			logConfig["sls"] = sls
		}

		request["logConfig"] = logConfig
	}

	if v, ok := d.GetOk("spec"); ok {
		request["spec"] = v
	}
	if v, ok := d.GetOk("gateway_type"); ok {
		request["gatewayType"] = v
	}
	networkAccessConfig := make(map[string]interface{})

	if v := d.Get("network_access_config"); !helper.IsNil(v) {
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			networkAccessConfig["type"] = type1
		}

		request["networkAccessConfig"] = networkAccessConfig
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_apig_gateway", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.gatewayId", response)
	d.SetId(fmt.Sprint(id))

	apigServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, apigServiceV2.ApigGatewayStateRefreshFunc(d.Id(), "status", []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceGatewayUpdate(d, meta)
}

func resourceGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := NewService(client)

	objectRaw, err := apigServiceV2.DescribeApigGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_gateway DescribeApigGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_from", objectRaw["createFrom"])
	d.Set("create_time", objectRaw["createTimestamp"])
	d.Set("expire_time", objectRaw["expireTimestamp"])
	d.Set("gateway_edition", objectRaw["gatewayEdition"])
	d.Set("gateway_name", objectRaw["name"])
	d.Set("gateway_type", objectRaw["gatewayType"])
	d.Set("payment_type", convertApigGatewaydatachargeTypeResponse(objectRaw["chargeType"]))
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("spec", objectRaw["spec"])
	d.Set("status", objectRaw["status"])
	d.Set("target_version", objectRaw["targetVersion"])
	d.Set("update_time", objectRaw["updateTimestamp"])
	d.Set("version", objectRaw["version"])

	environmentsRaw := objectRaw["environments"]
	environmentsMaps := make([]map[string]interface{}, 0)
	if environmentsRaw != nil {
		for _, environmentsChildRaw := range convert.ConvertToInterfaceArray(environmentsRaw) {
			environmentsMap := make(map[string]interface{})
			environmentsChildRaw := environmentsChildRaw.(map[string]interface{})
			environmentsMap["alias"] = environmentsChildRaw["alias"]
			environmentsMap["environment_id"] = environmentsChildRaw["environmentId"]
			environmentsMap["name"] = environmentsChildRaw["name"]

			environmentsMaps = append(environmentsMaps, environmentsMap)
		}
	}
	if err := d.Set("environments", environmentsMaps); err != nil {
		return err
	}
	loadBalancersRaw := objectRaw["loadBalancers"]
	loadBalancersMaps := make([]map[string]interface{}, 0)
	if loadBalancersRaw != nil {
		for _, loadBalancersChildRaw := range convert.ConvertToInterfaceArray(loadBalancersRaw) {
			loadBalancersMap := make(map[string]interface{})
			loadBalancersChildRaw := loadBalancersChildRaw.(map[string]interface{})
			loadBalancersMap["address"] = loadBalancersChildRaw["address"]
			loadBalancersMap["address_ip_version"] = loadBalancersChildRaw["addressIpVersion"]
			loadBalancersMap["address_type"] = loadBalancersChildRaw["addressType"]
			loadBalancersMap["gateway_default"] = loadBalancersChildRaw["gatewayDefault"]
			loadBalancersMap["load_balancer_id"] = loadBalancersChildRaw["loadBalancerId"]
			loadBalancersMap["mode"] = loadBalancersChildRaw["mode"]
			loadBalancersMap["status"] = loadBalancersChildRaw["status"]
			loadBalancersMap["type"] = loadBalancersChildRaw["type"]

			ipv4AddressesRaw := make([]interface{}, 0)
			if loadBalancersChildRaw["ipv4Addresses"] != nil {
				ipv4AddressesRaw = convert.ConvertToInterfaceArray(loadBalancersChildRaw["ipv4Addresses"])
			}

			loadBalancersMap["ipv4_addresses"] = ipv4AddressesRaw
			ipv6AddressesRaw := make([]interface{}, 0)
			if loadBalancersChildRaw["ipv6Addresses"] != nil {
				ipv6AddressesRaw = convert.ConvertToInterfaceArray(loadBalancersChildRaw["ipv6Addresses"])
			}

			loadBalancersMap["ipv6_addresses"] = ipv6AddressesRaw
			portsRaw := loadBalancersChildRaw["ports"]
			portsMaps := make([]map[string]interface{}, 0)
			if portsRaw != nil {
				for _, portsChildRaw := range convert.ConvertToInterfaceArray(portsRaw) {
					portsMap := make(map[string]interface{})
					portsChildRaw := portsChildRaw.(map[string]interface{})
					portsMap["port"] = portsChildRaw["port"]
					portsMap["protocol"] = portsChildRaw["protocol"]

					portsMaps = append(portsMaps, portsMap)
				}
			}
			loadBalancersMap["ports"] = portsMaps
			loadBalancersMaps = append(loadBalancersMaps, loadBalancersMap)
		}
	}
	if err := d.Set("load_balancers", loadBalancersMaps); err != nil {
		return err
	}
	securityGroupMaps := make([]map[string]interface{}, 0)
	securityGroupMap := make(map[string]interface{})
	securityGroupRaw := make(map[string]interface{})
	if objectRaw["securityGroup"] != nil {
		securityGroupRaw = objectRaw["securityGroup"].(map[string]interface{})
	}
	if len(securityGroupRaw) > 0 {
		securityGroupMap["name"] = securityGroupRaw["name"]
		securityGroupMap["security_group_id"] = securityGroupRaw["securityGroupId"]

		securityGroupMaps = append(securityGroupMaps, securityGroupMap)
	}
	if err := d.Set("security_group", securityGroupMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))
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
		for _, zonesChildRaw := range convert.ConvertToInterfaceArray(zonesRaw) {
			zonesMap := make(map[string]interface{})
			zonesChildRaw := zonesChildRaw.(map[string]interface{})
			zonesMap["name"] = zonesChildRaw["name"]
			zonesMap["zone_id"] = zonesChildRaw["zoneId"]

			vSwitchRawObj, _ := jsonpath.Get("$.vSwitch", zonesChildRaw)
			vSwitchRaw := make(map[string]interface{})
			if vSwitchRawObj != nil {
				vSwitchRaw = vSwitchRawObj.(map[string]interface{})
			}
			zonesMap["vswitch_id"] = vSwitchRaw["vSwitchId"]

			zonesMaps = append(zonesMaps, zonesMap)
		}
	}
	if err := d.Set("zones", zonesMaps); err != nil {
		return err
	}

	return nil
}

func resourceGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	gatewayId := d.Id()
	action := fmt.Sprintf("/v1/gateways/%s/name", gatewayId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("gateway_name") {
		update = true
	}
	if v, ok := d.GetOk("gateway_name"); ok {
		query["name"] = helper.StringPointer(v.(string))
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = helper.StringPointer(v.(string))
	}

	query["Service"] = helper.StringPointer("APIG")
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		apigServiceV2 := NewService(client)
		expectedResourceGroupId := d.Get("resource_group_id").(string)
		stateConf := helper.BuildStateConf([]string{}, []string{expectedResourceGroupId}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, apigServiceV2.ApigGatewayStateRefreshFunc(d.Id(), "resourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		apigServiceV2 := NewService(client)
		if err := apigServiceV2.SetResourceTags(d, "gateway"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceGatewayRead(d, meta)
}

func resourceGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	enableDelete := false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		gatewayId := d.Id()
		action := fmt.Sprintf("/v1/gateways/%s", gatewayId)
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		var err error
		request = make(map[string]interface{})

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotFound.GatewayNotFound", "Conflict.GatewayIsDeleted"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		apigServiceV2 := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, apigServiceV2.ApigGatewayStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

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
