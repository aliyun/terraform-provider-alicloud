package polardb

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGatewayConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGatewayConsumerCreate,
		Read:   resourceGatewayConsumerRead,
		Update: resourceGatewayConsumerUpdate,
		Delete: resourceGatewayConsumerDelete,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the consumer.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the consumer.",
			},
			"consumer_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the consumer group to which the consumer belongs.",
			},
			"consumer_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the consumer group.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The nickname of the consumer.",
			},
			"key_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ApiKey",
				ValidateFunc: validators.StringInSlice([]string{"ApiKey"}, false),
				Description:  "The key type. Currently, only `ApiKey` is supported.",
			},
			"is_default": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "0",
				ValidateFunc: validators.StringInSlice([]string{"0", "1"}, false),
				Description:  "Whether the consumer belongs to the default group. Valid values: `0`, `1`.",
			},
			"api_key_reset_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An arbitrary value that, when changed, resets the consumer API key.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The consumer API key. It is returned only when the consumer is created or the key is reset and is stored in Terraform state.",
			},
			"allowed_models": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The models that the consumer is allowed to access.",
			},
			"month_to_date_cost_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The month-to-date cost count.",
			},
			"lifetime_cost_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The lifetime cost count.",
			},
			"month_to_date_token_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The month-to-date token count.",
			},
			"lifetime_token_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The lifetime token count.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the consumer was created.",
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the consumer was last modified.",
			},
		},
	}
}

func resourceGatewayConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateConsumer"
	request := polarDBGatewayConsumerRequest(d, client.RegionId)
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_gateway_consumer", action, errs.AlibabaCloudSdkGoERROR)
	}
	consumerID := fmt.Sprint(response["ConsumerId"])
	if consumerID == "" || consumerID == "<nil>" {
		return errs.WrapError(fmt.Errorf("CreateConsumer returned empty ConsumerId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), consumerID))
	if apiKey := fmt.Sprint(response["ApiKey"]); apiKey != "" && apiKey != "<nil>" {
		if err = d.Set("api_key", apiKey); err != nil {
			return errs.WrapError(err)
		}
	}
	if d.Get("is_default").(string) == "1" {
		modifyRequest := map[string]interface{}{
			"RegionId":    client.RegionId,
			"GwClusterId": d.Get("gateway_id"),
			"ConsumerId":  consumerID,
			"Name":        d.Get("name"),
			"IsDefault":   d.Get("is_default"),
		}
		if value, ok := d.GetOk("consumer_group_id"); ok {
			modifyRequest["ConsumerGroupName"] = value
		}
		if _, err = callPolarDBGatewayAI(client, "ModifyConsumer", modifyRequest, d.Timeout(schema.TimeoutCreate)); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyConsumer", errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceGatewayConsumerRead(d, meta)
}

func resourceGatewayConsumerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := FindPolarDBGatewayAIItem(client, "DescribeConsumers", gatewayID, "ConsumerId", "ConsumerId", consumerID)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DescribeConsumers", errs.AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("consumer_id", consumerID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("name", object["Name"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("consumer_group_id", object["ConsumerGroupId"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("consumer_group_name", object["ConsumerGroupName"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("allowed_models", object["AllowedModels"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("month_to_date_cost_count", object["MtdCostCount"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("lifetime_cost_count", object["LifetimeCostCount"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("month_to_date_token_count", object["MtdTokenCount"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("lifetime_token_count", object["LifetimeTokenCount"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("create_time", object["GmtCreated"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("modify_time", object["GmtModified"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("nickname", object["NickName"]); err != nil {
		return errs.WrapError(err)
	}
	// ApiKey is returned only by create/reset operations. Never replace the
	// configured state value with a masked or absent value from DescribeConsumers.
	return nil
}

func resourceGatewayConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	if d.HasChange("api_key_reset_token") {
		_, newValue := d.GetChange("api_key_reset_token")
		if fmt.Sprint(newValue) != "" {
			request := map[string]interface{}{
				"RegionId":    client.RegionId,
				"GwClusterId": gatewayID,
				"ConsumerId":  consumerID,
			}
			response, resetErr := callPolarDBGatewayAI(client, "ResetConsumerApiKey", request, d.Timeout(schema.TimeoutUpdate))
			if resetErr != nil {
				return errs.WrapErrorf(resetErr, errs.DefaultErrorMsg, d.Id(), "ResetConsumerApiKey", errs.AlibabaCloudSdkGoERROR)
			}
			apiKey := fmt.Sprint(response["ApiKey"])
			if apiKey == "" || apiKey == "<nil>" {
				return errs.WrapError(fmt.Errorf("ResetConsumerApiKey returned empty ApiKey"))
			}
			if err = d.Set("api_key", apiKey); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	return resourceGatewayConsumerRead(d, meta)
}

func resourceGatewayConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"GwClusterId": gatewayID,
		"ConsumerId":  consumerID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteConsumer", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteConsumer", errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayConsumerRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":    regionID,
		"GwClusterId": d.Get("gateway_id"),
		"Name":        d.Get("name"),
		"KeyType":     d.Get("key_type"),
	}
	for key, requestKey := range map[string]string{
		"consumer_group_id": "ConsumerGroupName",
		"nickname":          "NickName",
	} {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
	return request
}
