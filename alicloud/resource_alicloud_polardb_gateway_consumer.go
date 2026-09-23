package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGatewayConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayConsumerCreate,
		Read:   resourceAlicloudPolarDBGatewayConsumerRead,
		Update: resourceAlicloudPolarDBGatewayConsumerUpdate,
		Delete: resourceAlicloudPolarDBGatewayConsumerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the consumer.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"name": {
				Type: schema.TypeString, Required: true,
				Description: "The name of the consumer.",
			},
			"consumer_group_id": {
				Type: schema.TypeString, Optional: true,
				Description: "The ID of the consumer group to which the consumer belongs.",
			},
			"consumer_group_name": {
				Type: schema.TypeString, Computed: true,
				Description: "The name of the consumer group.",
			},
			"nickname": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The nickname of the consumer.",
			},
			"key_type": {
				Type: schema.TypeString, Optional: true, ForceNew: true, Default: "ApiKey",
				ValidateFunc: StringInSlice([]string{"ApiKey"}, false),
				Description:  "The key type. Currently, only `ApiKey` is supported.",
			},
			"is_default": {
				Type: schema.TypeString, Optional: true, Default: "0",
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
				Description:  "Whether the consumer belongs to the default group. Valid values: `0`, `1`.",
			},
			"api_key_reset_token": {
				Type: schema.TypeString, Optional: true,
				Description: "An arbitrary value that, when changed, resets the consumer API key.",
			},
			"api_key": {
				Type: schema.TypeString, Computed: true, Sensitive: true,
				Description: "The consumer API key. It is returned only when the consumer is created or the key is reset and is stored in Terraform state.",
			},
			"allowed_models": {
				Type: schema.TypeString, Computed: true,
				Description: "The models that the consumer is allowed to access.",
			},
			"month_to_date_cost_count": {
				Type: schema.TypeInt, Computed: true,
				Description: "The month-to-date cost count.",
			},
			"lifetime_cost_count": {
				Type: schema.TypeInt, Computed: true,
				Description: "The lifetime cost count.",
			},
			"month_to_date_token_count": {
				Type: schema.TypeInt, Computed: true,
				Description: "The month-to-date token count.",
			},
			"lifetime_token_count": {
				Type: schema.TypeInt, Computed: true,
				Description: "The lifetime token count.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the consumer was created.",
			},
			"modify_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the consumer was last modified.",
			},
		},
	}
}

func resourceAlicloudPolarDBGatewayConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateConsumer"
	request := polarDBGatewayConsumerRequest(d, client.RegionId)
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway_consumer", action, AlibabaCloudSdkGoERROR)
	}
	consumerID := fmt.Sprint(response["ConsumerId"])
	if consumerID == "" || consumerID == "<nil>" {
		return WrapError(fmt.Errorf("CreateConsumer returned empty ConsumerId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), consumerID))
	if apiKey := fmt.Sprint(response["ApiKey"]); apiKey != "" && apiKey != "<nil>" {
		if err = d.Set("api_key", apiKey); err != nil {
			return WrapError(err)
		}
	}
	if d.Get("is_default").(string) == "1" {
		modifyRequest := map[string]interface{}{
			"RegionId": client.RegionId, "GwClusterId": d.Get("gateway_id"), "ConsumerId": consumerID,
			"Name": d.Get("name"), "IsDefault": d.Get("is_default"),
		}
		if value, ok := d.GetOk("consumer_group_id"); ok {
			modifyRequest["ConsumerGroupName"] = value
		}
		if _, err = callPolarDBGatewayAI(client, "ModifyConsumer", modifyRequest, d.Timeout(schema.TimeoutCreate)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyConsumer", AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudPolarDBGatewayConsumerRead(d, meta)
}

func resourceAlicloudPolarDBGatewayConsumerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	object, err := findPolarDBGatewayAIItem(client, "DescribeConsumers", gatewayID, "ConsumerId", "ConsumerId", consumerID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeConsumers", AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	values := map[string]interface{}{
		"gateway_id": gatewayID, "consumer_id": consumerID, "name": object["Name"],
		"consumer_group_id": object["ConsumerGroupId"], "consumer_group_name": object["ConsumerGroupName"],
		"allowed_models":           object["AllowedModels"],
		"month_to_date_cost_count": object["MtdCostCount"], "lifetime_cost_count": object["LifetimeCostCount"],
		"month_to_date_token_count": object["MtdTokenCount"], "lifetime_token_count": object["LifetimeTokenCount"],
		"create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
	}
	if nickname, ok := object["ConsumerTag"]; ok {
		values["nickname"] = nickname
	} else if nickname, ok = object["NickName"]; ok {
		values["nickname"] = nickname
	}
	for key, value := range values {
		if err = d.Set(key, value); err != nil {
			return WrapError(err)
		}
	}
	// ApiKey is returned only by create/reset operations. Never replace the
	// configured state value with a masked or absent value from DescribeConsumers.
	return nil
}

func resourceAlicloudPolarDBGatewayConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if d.HasChanges("name", "consumer_group_id", "nickname", "is_default") {
		request := map[string]interface{}{
			"RegionId": client.RegionId, "GwClusterId": gatewayID, "ConsumerId": consumerID,
			"Name": d.Get("name"), "IsDefault": d.Get("is_default"),
		}
		if value, ok := d.GetOk("consumer_group_id"); ok {
			request["ConsumerGroupName"] = value
		}
		if value, ok := d.GetOk("nickname"); ok {
			request["NickName"] = value
		}
		if _, err = callPolarDBGatewayAI(client, "ModifyConsumer", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyConsumer", AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("api_key_reset_token") {
		_, newValue := d.GetChange("api_key_reset_token")
		if fmt.Sprint(newValue) != "" {
			request := map[string]interface{}{
				"RegionId": client.RegionId, "GwClusterId": gatewayID, "ConsumerId": consumerID,
			}
			response, resetErr := callPolarDBGatewayAI(client, "ResetConsumerApiKey", request, d.Timeout(schema.TimeoutUpdate))
			if resetErr != nil {
				return WrapErrorf(resetErr, DefaultErrorMsg, d.Id(), "ResetConsumerApiKey", AlibabaCloudSdkGoERROR)
			}
			apiKey := fmt.Sprint(response["ApiKey"])
			if apiKey == "" || apiKey == "<nil>" {
				return WrapError(fmt.Errorf("ResetConsumerApiKey returned empty ApiKey"))
			}
			if err = d.Set("api_key", apiKey); err != nil {
				return WrapError(err)
			}
		}
	}
	return resourceAlicloudPolarDBGatewayConsumerRead(d, meta)
}

func resourceAlicloudPolarDBGatewayConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "ConsumerId": consumerID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteConsumer", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteConsumer", AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayConsumerRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"),
		"Name": d.Get("name"), "KeyType": d.Get("key_type"),
	}
	for key, requestKey := range map[string]string{"consumer_group_id": "ConsumerGroupName", "nickname": "NickName"} {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
	return request
}
