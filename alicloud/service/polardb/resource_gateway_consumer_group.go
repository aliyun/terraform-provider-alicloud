package polardb

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGatewayConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGatewayConsumerGroupCreate,
		Read:   resourceGatewayConsumerGroupRead,
		Update: resourceGatewayConsumerGroupUpdate,
		Delete: resourceGatewayConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consumer_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the consumer group.",
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
				Description: "The name of the consumer group.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The nickname of the consumer group.",
			},
			"is_default": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the consumer group is the default group.",
			},
			"allowed_models": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The models that the consumer group is allowed to access.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the consumer group was created.",
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the consumer group was last modified.",
			},
		},
	}
}

func resourceGatewayConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateConsumerGroup"
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"GwClusterId":       d.Get("gateway_id"),
		"ConsumerGroupName": d.Get("name"),
	}
	if value, ok := d.GetOk("nickname"); ok {
		request["NickName"] = value
	}
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_gateway_consumer_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	consumerGroupID := fmt.Sprint(response["ConsumerGroupId"])
	if consumerGroupID == "" || consumerGroupID == "<nil>" {
		return errs.WrapError(fmt.Errorf("CreateConsumerGroup returned empty ConsumerGroupId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), consumerGroupID))
	return resourceGatewayConsumerGroupRead(d, meta)
}

func resourceGatewayConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerGroupID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := FindPolarDBGatewayAIItem(client, "DescribeConsumerGroups", gatewayID, "ConsumerGroupId", "ConsumerGroupId", consumerGroupID)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DescribeConsumerGroups", errs.AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("consumer_group_id", consumerGroupID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("name", object["ConsumerGroupName"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("is_default", object["IsDefault"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("allowed_models", object["AllowedModels"]); err != nil {
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
	return nil
}

func resourceGatewayConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerGroupID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"GwClusterId":       gatewayID,
		"ConsumerGroupName": consumerGroupID,
	}
	if value, ok := d.GetOk("nickname"); ok {
		request["NickName"] = value
	}
	if _, err = callPolarDBGatewayAI(client, "ModifyConsumerGroup", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyConsumerGroup", errs.AlibabaCloudSdkGoERROR)
	}
	return resourceGatewayConsumerGroupRead(d, meta)
}

func resourceGatewayConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, consumerGroupID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"GwClusterId":       gatewayID,
		"ConsumerGroupName": consumerGroupID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteConsumerGroup", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteConsumerGroup", errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}
