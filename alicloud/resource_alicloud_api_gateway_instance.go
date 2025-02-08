// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayInstanceCreate,
		Read:   resourceAliCloudApiGatewayInstanceRead,
		Update: resourceAliCloudApiGatewayInstanceUpdate,
		Delete: resourceAliCloudApiGatewayInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Update: schema.DefaultTimeout(25 * time.Minute),
			Delete: schema.DefaultTimeout(25 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"egress_ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"https_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Instance name"),
			},
			"instance_spec": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"normal", "vpc_connect"}, true),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, true),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"year", "month"}, true),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "User's VpcID"),
			},
			"vpc_slb_intranet_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_vswitch_security_group": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"connect_cidr_blocks": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"to_connect_vpc_ip_block": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"delete_vpc_ip_block"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Required: true,
						},
						"customized": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"delete_vpc_ip_block": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"to_connect_vpc_ip_block"},
			},
		},
	}
}

func resourceAliCloudApiGatewayInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Token"] = buildClientToken(action)

	request["InstanceName"] = d.Get("instance_name")
	request["InstanceSpec"] = d.Get("instance_spec")
	request["ChargeType"] = convertApiGatewayInstanceChargeTypeRequest(d.Get("payment_type").(string))
	request["HttpsPolicy"] = d.Get("https_policy")
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("instance_cidr"); ok {
		request["InstanceCidr"] = v
	}
	if connectConfigs, ok := d.GetOk("zone_vswitch_security_group"); ok {
		connectConfigList := connectConfigs.([]interface{})
		zoneVSwitchSecurityGroups := make([]map[string]interface{}, len(connectConfigList))
		for i, connectedConfig := range connectConfigList {
			connectConfigMap := connectedConfig.(map[string]interface{})
			zoneVSwitchSecurityGroup := make(map[string]interface{})
			zoneVSwitchSecurityGroup["ZoneId"] = connectConfigMap["zone_id"]
			zoneVSwitchSecurityGroup["VSwitchId"] = connectConfigMap["vswitch_id"]
			zoneVSwitchSecurityGroup["CidrBlock"] = connectConfigMap["cidr_block"]
			zoneVSwitchSecurityGroup["SecurityGroupId"] = connectConfigMap["security_group"]
			zoneVSwitchSecurityGroups[i] = zoneVSwitchSecurityGroup
		}
		request["ZoneVSwitchSecurityGroup"] = zoneVSwitchSecurityGroups
	}
	if v, ok := d.GetOk("user_vpc_id"); ok {
		request["UserVpcId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
		request["Token"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	apiGatewayServiceV2 := ApiGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, apiGatewayServiceV2.ApiGatewayInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudApiGatewayInstanceUpdate(d, meta)
}

func resourceAliCloudApiGatewayInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_instance DescribeApiGatewayInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("https_policy", objectRaw["HttpsPolicies"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_spec", objectRaw["InstanceSpec"])
	d.Set("instance_type", objectRaw["DedicatedInstanceType"])
	d.Set("payment_type", convertApiGatewayInstanceInstancesInstanceAttributeInstanceChargeTypeResponse(objectRaw["InstanceChargeType"]))
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("egress_ipv6_enable", objectRaw["EgressIpv6Enable"])
	d.Set("vpc_slb_intranet_enable", objectRaw["VpcSlbIntranetEnable"])
	d.Set("support_ipv6", objectRaw["SupportIpv6"])
	d.Set("instance_cidr", objectRaw["InstanceCidrBlock"])
	d.Set("user_vpc_id", objectRaw["ConnectVpcId"])
	v := convertApiGatewayInstanceZoneVSwitchSecurityGroupResponse(objectRaw["NetworkInterfaceAttributes"])
	if v != nil {
		d.Set("zone_vswitch_security_group", v)
	}
	d.Set("connect_cidr_blocks", objectRaw["ConnectCidrBlocks"])

	if vipTypeList, ok := objectRaw["VipTypeList"].(string); ok && strings.Contains(vipTypeList, "VPC_INTERNET_IPV6") {
		d.Set("ipv6_enabled", true)
	} else {
		d.Set("ipv6_enabled", false)
	}

	return nil
}

func resourceAliCloudApiGatewayInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyInstanceAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	request["InstanceName"] = d.Get("instance_name")
	if !d.IsNewResource() && d.HasChange("https_policy") {
		update = true
	}
	request["HttpsPolicy"] = d.Get("https_policy")
	if v, ok := d.GetOkExists("egress_ipv6_enable"); ok {
		request["EgressIpv6Enable"] = v
	}
	if v, ok := d.GetOkExists("vpc_slb_intranet_enable"); ok {
		request["VpcSlbIntranetEnable"] = v
	}
	if v, ok := d.GetOkExists("ipv6_enabled"); ok {
		request["IPV6Enabled"] = v
	}
	if d.HasChange("to_connect_vpc_ip_block") {
		if v, ok := d.GetOk("to_connect_vpc_ip_block"); ok {
			originalVpcIpBlock := v.(map[string]interface{})
			newVpcIpBlock := make(map[string]interface{})
			if cidrBlock, ok := originalVpcIpBlock["cidr_block"]; ok {
				newVpcIpBlock["CidrBlock"] = cidrBlock
			}
			if vswitchId, ok := originalVpcIpBlock["vswitch_id"]; ok {
				newVpcIpBlock["VSwitchId"] = vswitchId
			}
			if zoneId, ok := originalVpcIpBlock["zone_id"]; ok {
				newVpcIpBlock["ZoneId"] = zoneId
			}
			if customized, ok := originalVpcIpBlock["customized"]; ok {
				newVpcIpBlock["Customized"] = customized
			}
			request["ToConnectVpcIpBlock"], err = convertMaptoJsonString(newVpcIpBlock)
			if err != nil {
				return WrapError(err)
			}
			update = true
		}
	}
	if d.HasChange("delete_vpc_ip_block") {
		if deleteVpcIpBlock, ok := d.GetOk("delete_vpc_ip_block"); ok {
			request["DeleteVpcIpBlock"] = deleteVpcIpBlock
			update = true
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)

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
		apiGatewayServiceV2 := ApiGatewayServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, apiGatewayServiceV2.ApiGatewayInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudApiGatewayInstanceRead(d, meta)
}

func resourceAliCloudApiGatewayInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)

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

	apiGatewayServiceV2 := ApiGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, apiGatewayServiceV2.ApiGatewayInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertApiGatewayInstanceInstancesInstanceAttributeInstanceChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertApiGatewayInstanceZoneVSwitchSecurityGroupResponse(source interface{}) interface{} {
	if source == nil {
		return nil
	}
	networkAttributesMap, ok := source.(map[string]interface{})
	if !ok {
		return nil
	}
	networkInterfaces, ok := networkAttributesMap["NetworkInterfaceAttribute"].([]interface{})
	if !ok {
		return nil
	}
	var zoneVswitchSecurityGroups []map[string]interface{}
	for _, ni := range networkInterfaces {
		networkInterface, ok := ni.(map[string]interface{})
		if !ok {
			continue
		}
		zoneVSwitchSecurityGroup := map[string]interface{}{
			"zone_id":        networkInterface["ZoneId"],
			"vswitch_id":     networkInterface["VswitchId"],
			"cidr_block":     networkInterface["CidrBlock"],
			"security_group": networkInterface["SecurityGroupId"],
		}
		zoneVswitchSecurityGroups = append(zoneVswitchSecurityGroups, zoneVSwitchSecurityGroup)
	}
	return zoneVswitchSecurityGroups
}

func convertApiGatewayInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PREPAY"
	}
	return source
}
