package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayCreate,
		Read:   resourceAlicloudPolarDBGatewayRead,
		Delete: resourceAlicloudPolarDBGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The region ID of the PolarDB gateway.",
			},
			"zone_id": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The zone ID of the PolarDB gateway.",
			},
			"db_cluster_class": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The specifications of the PolarDB gateway.",
			},
			"pay_type": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"Postpaid", "Prepaid"}, false),
				Description:  "The billing method. Valid values: `Postpaid`, `Prepaid`.",
			},
			"auto_renew": {
				Type: schema.TypeBool, Optional: true, ForceNew: true, Default: false,
				Description: "Whether to enable automatic renewal for a subscription gateway.",
			},
			"period": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
				Description:  "The unit of the subscription duration. Valid values: `Month`, `Year`.",
			},
			"used_time": {
				Type: schema.TypeInt, Optional: true, ForceNew: true,
				ValidateFunc: IntBetween(1, 9),
				Description:  "The subscription duration. Valid values are 1 to 9 for Month and 1 to 3 for Year.",
			},
			"vpc_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the VPC.",
			},
			"vswitch_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the vSwitch.",
			},
			"security_group_id": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The ID of the security group.",
			},
			"db_type": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"MySQL", "PostgreSQL"}, false),
				Description:  "The database engine. Valid values: `MySQL`, `PostgreSQL`.",
			},
			"status": {
				Type: schema.TypeString, Computed: true,
				Description: "The status of the PolarDB gateway.",
			},
			"description": {
				Type: schema.TypeString, Computed: true,
				Description: "The description of the PolarDB gateway.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the PolarDB gateway was created.",
			},
			"modify_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the PolarDB gateway was last modified.",
			},
			"expire_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The expiration time of the PolarDB gateway.",
			},
			"expired": {
				Type: schema.TypeBool, Computed: true,
				Description: "Indicates whether the PolarDB gateway has expired.",
			},
			"latest_version": {
				Type: schema.TypeString, Computed: true,
				Description: "The latest available gateway version.",
			},
			"current_version": {
				Type: schema.TypeString, Computed: true,
				Description: "The current gateway version.",
			},
			"running_version": {
				Type: schema.TypeString, Computed: true,
				Description: "The running gateway version.",
			},
			"endpoints": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"address":      {Type: schema.TypeString, Computed: true},
					"endpoint_id":  {Type: schema.TypeString, Computed: true},
					"gateway_id":   {Type: schema.TypeString, Computed: true},
					"port":         {Type: schema.TypeInt, Computed: true},
					"tunnel_id":    {Type: schema.TypeString, Computed: true},
					"vpc_id":       {Type: schema.TypeString, Computed: true},
					"network_type": {Type: schema.TypeString, Computed: true},
				}},
			},
			"security_ip_arrays": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name":    {Type: schema.TypeString, Computed: true},
					"ip_list": {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func resourceAlicloudPolarDBGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := PolarDBService{client}

	if d.Get("pay_type").(string) == "Prepaid" {
		period, periodOk := d.GetOk("period")
		usedTime, usedTimeOk := d.GetOk("used_time")
		if !periodOk || !usedTimeOk {
			return WrapError(fmt.Errorf("'period' and 'used_time' are required when 'pay_type' is 'Prepaid'"))
		}
		if period.(string) == "Year" && usedTime.(int) > 3 {
			return WrapError(fmt.Errorf("'used_time' must be between 1 and 3 when 'period' is 'Year'"))
		}
	}

	action := "CreateGateway"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"PayType":   d.Get("pay_type"),
		"VPCId":     d.Get("vpc_id"),
		"VSwitchId": d.Get("vswitch_id"),
	}
	for key, requestKey := range map[string]string{
		"zone_id": "ZoneId", "db_cluster_class": "DBClusterClass",
		"security_group_id": "SecurityGroupId", "db_type": "DBType",
	} {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
	if d.Get("pay_type").(string) == "Prepaid" {
		request["AutoRenew"] = d.Get("auto_renew")
		request["Period"] = d.Get("period")
		request["UsedTime"] = d.Get("used_time")
	}

	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway", action, AlibabaCloudSdkGoERROR)
	}

	gatewayId := fmt.Sprint(response["GwClusterId"])
	if gatewayId == "" || gatewayId == "<nil>" {
		return WrapError(fmt.Errorf("CreateGateway returned empty GwClusterId"))
	}
	d.SetId(gatewayId)

	stateConf := BuildStateConf([]string{"CREATE", "CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, service.PolarDBGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudPolarDBGatewayRead(d, meta)
}

func resourceAlicloudPolarDBGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := PolarDBService{client}
	object, err := service.DescribePolarDBGatewayAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if err = d.Set("region_id", client.RegionId); err != nil {
		return WrapError(err)
	}
	for key, responseKey := range map[string]string{
		"status": "Status", "description": "GwDescription", "create_time": "CreateTime",
		"modify_time": "ModifyTime", "expire_time": "ExpireTime", "expired": "Expired",
		"latest_version": "LatestVersion", "current_version": "CurrentVersion", "running_version": "RunningVersion",
	} {
		if value, ok := object[responseKey]; ok && value != nil {
			if err = d.Set(key, value); err != nil {
				return WrapError(err)
			}
		}
	}
	if err = d.Set("endpoints", flattenPolarDBGatewayEndpoints(object["Endpoints"])); err != nil {
		return WrapError(err)
	}
	if err = d.Set("security_ip_arrays", flattenPolarDBGatewaySecurityIPArrays(object["SecurityIPArrays"])); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudPolarDBGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := PolarDBService{client}
	action := "DeleteGateway"
	request := map[string]interface{}{"GwClusterId": d.Id(), "RegionId": client.RegionId}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
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

	stateConf := BuildStateConf([]string{"CREATE", "ACTIVATION", "DELETING"}, []string{"DELETED"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, service.PolarDBGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId("")
	return nil
}

func flattenPolarDBGatewayEndpoints(raw interface{}) []map[string]interface{} {
	items := polarDBGatewayList(raw, "Endpoint")
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"address": item["Address"], "endpoint_id": item["EndpointId"], "gateway_id": item["GwClusterId"],
			"port": formatInt(item["Port"]), "tunnel_id": item["TunnelId"], "vpc_id": item["VpcId"], "network_type": item["NetType"],
		})
	}
	return result
}

func flattenPolarDBGatewaySecurityIPArrays(raw interface{}) []map[string]interface{} {
	items := polarDBGatewayList(raw, "SecurityIPArray")
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{"name": item["SecurityIPArrayName"], "ip_list": item["SecurityIPList"]})
	}
	return result
}

func polarDBGatewayList(raw interface{}, nestedKey string) []map[string]interface{} {
	if wrapper, ok := raw.(map[string]interface{}); ok {
		raw = wrapper[nestedKey]
	}
	values, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		if item, ok := value.(map[string]interface{}); ok {
			result = append(result, item)
		}
	}
	return result
}
