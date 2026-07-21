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

func resourceAliCloudResourceManagerDeliveryChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerDeliveryChannelCreate,
		Read:   resourceAliCloudResourceManagerDeliveryChannelRead,
		Update: resourceAliCloudResourceManagerDeliveryChannelUpdate,
		Delete: resourceAliCloudResourceManagerDeliveryChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"delivery_channel_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delivery_channel_filter": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"delivery_channel_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_change_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sls_properties": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oversized_data_oss_target_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"target_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"resource_snapshot_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delivery_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sls_properties": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oversized_data_oss_target_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"target_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudResourceManagerDeliveryChannelCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDeliveryChannel"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("resource_snapshot_delivery"); ok {
		resourceSnapshotDeliveryCustomExpressionJsonPath, err := jsonpath.Get("$[0].custom_expression", v)
		if err == nil && resourceSnapshotDeliveryCustomExpressionJsonPath != "" {
			request["ResourceSnapshotDelivery.CustomExpression"] = resourceSnapshotDeliveryCustomExpressionJsonPath
		}
	}
	if v, ok := d.GetOk("resource_snapshot_delivery"); ok {
		resourceSnapshotDeliveryDeliveryTimeJsonPath, err := jsonpath.Get("$[0].delivery_time", v)
		if err == nil && resourceSnapshotDeliveryDeliveryTimeJsonPath != "" {
			request["ResourceSnapshotDelivery.DeliveryTime"] = resourceSnapshotDeliveryDeliveryTimeJsonPath
		}
	}
	request["ResourceChangeDelivery.TargetType"] = "SLS"
	dataList := make(map[string]interface{})

	if v := d.Get("delivery_channel_filter"); v != nil {
		resourceTypes1, _ := jsonpath.Get("$[0].resource_types", v)
		if resourceTypes1 != nil && resourceTypes1 != "" {
			dataList["ResourceTypes"] = resourceTypes1
		}

		request["DeliveryChannelFilter"] = dataList
	}

	if v, ok := d.GetOk("resource_snapshot_delivery"); ok {
		resourceSnapshotDeliveryTargetTypeJsonPath, err := jsonpath.Get("$[0].target_type", v)
		if err == nil && resourceSnapshotDeliveryTargetTypeJsonPath != "" {
			request["ResourceSnapshotDelivery.TargetType"] = resourceSnapshotDeliveryTargetTypeJsonPath
		}
	}
	if v, ok := d.GetOk("delivery_channel_description"); ok {
		request["DeliveryChannelDescription"] = v
	}
	request["DeliveryChannelName"] = d.Get("delivery_channel_name")
	if v, ok := d.GetOk("resource_snapshot_delivery"); ok {
		resourceSnapshotDeliveryTargetArnJsonPath, err := jsonpath.Get("$[0].target_arn", v)
		if err == nil && resourceSnapshotDeliveryTargetArnJsonPath != "" {
			request["ResourceSnapshotDelivery.TargetArn"] = resourceSnapshotDeliveryTargetArnJsonPath
		}
	}
	dataList1 := make(map[string]interface{})

	if v := d.Get("resource_change_delivery"); !IsNil(v) {
		targetType1, _ := jsonpath.Get("$[0].target_type", v)
		if targetType1 != nil && targetType1 != "" {
			dataList1["TargetType"] = targetType1
		}
		slsProperties := make(map[string]interface{})
		oversizedDataOssTargetArn1, _ := jsonpath.Get("$[0].sls_properties[0].oversized_data_oss_target_arn", d.Get("resource_change_delivery"))
		if oversizedDataOssTargetArn1 != nil && oversizedDataOssTargetArn1 != "" {
			slsProperties["OversizedDataOssTargetArn"] = oversizedDataOssTargetArn1
		}

		dataList1["SlsProperties"] = slsProperties
		targetArn1, _ := jsonpath.Get("$[0].target_arn", v)
		if targetArn1 != nil && targetArn1 != "" {
			dataList1["TargetArn"] = targetArn1
		}

		request["ResourceChangeDelivery"] = dataList1
	}

	dataList2 := make(map[string]interface{})

	if v := d.Get("resource_snapshot_delivery"); !IsNil(v) {
		customExpression1, _ := jsonpath.Get("$[0].custom_expression", v)
		if customExpression1 != nil && customExpression1 != "" {
			dataList2["CustomExpression"] = customExpression1
		}
		deliveryTime1, _ := jsonpath.Get("$[0].delivery_time", v)
		if deliveryTime1 != nil && deliveryTime1 != "" {
			dataList2["DeliveryTime"] = deliveryTime1
		}
		targetType3, _ := jsonpath.Get("$[0].target_type", v)
		if targetType3 != nil && targetType3 != "" {
			dataList2["TargetType"] = targetType3
		}
		targetArn3, _ := jsonpath.Get("$[0].target_arn", v)
		if targetArn3 != nil && targetArn3 != "" {
			dataList2["TargetArn"] = targetArn3
		}
		slsProperties1 := make(map[string]interface{})
		oversizedDataOssTargetArn3, _ := jsonpath.Get("$[0].sls_properties[0].oversized_data_oss_target_arn", d.Get("resource_snapshot_delivery"))
		if oversizedDataOssTargetArn3 != nil && oversizedDataOssTargetArn3 != "" {
			slsProperties1["OversizedDataOssTargetArn"] = oversizedDataOssTargetArn3
		}

		dataList2["SlsProperties"] = slsProperties1

		request["ResourceSnapshotDelivery"] = dataList2
	}

	if v, ok := d.GetOk("resource_change_delivery"); ok {
		resourceChangeDeliveryTargetArnJsonPath, err := jsonpath.Get("$[0].target_arn", v)
		if err == nil && resourceChangeDeliveryTargetArnJsonPath != "" {
			request["ResourceChangeDelivery.TargetArn"] = resourceChangeDeliveryTargetArnJsonPath
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceCenter", "2022-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_delivery_channel", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeliveryChannelId"]))

	return resourceAliCloudResourceManagerDeliveryChannelRead(d, meta)
}

func resourceAliCloudResourceManagerDeliveryChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerDeliveryChannel(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_delivery_channel DescribeResourceManagerDeliveryChannel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("delivery_channel_description", objectRaw["DeliveryChannelDescription"])
	d.Set("delivery_channel_name", objectRaw["DeliveryChannelName"])

	deliveryChannelFilterMaps := make([]map[string]interface{}, 0)
	deliveryChannelFilterMap := make(map[string]interface{})
	resourceTypesRaw, _ := jsonpath.Get("$.DeliveryChannelFilter.ResourceTypes", objectRaw)

	deliveryChannelFilterMap["resource_types"] = resourceTypesRaw
	deliveryChannelFilterMaps = append(deliveryChannelFilterMaps, deliveryChannelFilterMap)
	if err := d.Set("delivery_channel_filter", deliveryChannelFilterMaps); err != nil {
		return err
	}
	resourceChangeDeliveryMaps := make([]map[string]interface{}, 0)
	resourceChangeDeliveryMap := make(map[string]interface{})
	resourceChangeDeliveryRaw := make(map[string]interface{})
	if objectRaw["ResourceChangeDelivery"] != nil {
		resourceChangeDeliveryRaw = objectRaw["ResourceChangeDelivery"].(map[string]interface{})
	}
	if len(resourceChangeDeliveryRaw) > 0 {
		resourceChangeDeliveryMap["enabled"] = formatBool(resourceChangeDeliveryRaw["Enabled"])
		resourceChangeDeliveryMap["target_arn"] = resourceChangeDeliveryRaw["TargetArn"]
		resourceChangeDeliveryMap["target_type"] = resourceChangeDeliveryRaw["TargetType"]

		slsPropertiesMaps := make([]map[string]interface{}, 0)
		slsPropertiesMap := make(map[string]interface{})
		slsPropertiesRaw := make(map[string]interface{})
		if resourceChangeDeliveryRaw["SlsProperties"] != nil {
			slsPropertiesRaw = resourceChangeDeliveryRaw["SlsProperties"].(map[string]interface{})
		}
		if len(slsPropertiesRaw) > 0 {
			slsPropertiesMap["oversized_data_oss_target_arn"] = slsPropertiesRaw["OversizedDataOssTargetArn"]

			slsPropertiesMaps = append(slsPropertiesMaps, slsPropertiesMap)
		}
		resourceChangeDeliveryMap["sls_properties"] = slsPropertiesMaps
		resourceChangeDeliveryMaps = append(resourceChangeDeliveryMaps, resourceChangeDeliveryMap)
	}
	if err := d.Set("resource_change_delivery", resourceChangeDeliveryMaps); err != nil {
		return err
	}
	resourceSnapshotDeliveryMaps := make([]map[string]interface{}, 0)
	resourceSnapshotDeliveryMap := make(map[string]interface{})
	resourceSnapshotDeliveryRaw := make(map[string]interface{})
	if objectRaw["ResourceSnapshotDelivery"] != nil {
		resourceSnapshotDeliveryRaw = objectRaw["ResourceSnapshotDelivery"].(map[string]interface{})
	}
	if len(resourceSnapshotDeliveryRaw) > 0 {
		resourceSnapshotDeliveryMap["custom_expression"] = resourceSnapshotDeliveryRaw["CustomExpression"]
		resourceSnapshotDeliveryMap["delivery_time"] = resourceSnapshotDeliveryRaw["DeliveryTime"]
		resourceSnapshotDeliveryMap["enabled"] = formatBool(resourceSnapshotDeliveryRaw["Enabled"])
		resourceSnapshotDeliveryMap["target_arn"] = resourceSnapshotDeliveryRaw["TargetArn"]
		resourceSnapshotDeliveryMap["target_type"] = resourceSnapshotDeliveryRaw["TargetType"]

		slsPropertiesMaps := make([]map[string]interface{}, 0)
		slsPropertiesMap := make(map[string]interface{})
		slsPropertiesRaw := make(map[string]interface{})
		if resourceSnapshotDeliveryRaw["SlsProperties"] != nil {
			slsPropertiesRaw = resourceSnapshotDeliveryRaw["SlsProperties"].(map[string]interface{})
		}
		if len(slsPropertiesRaw) > 0 {
			slsPropertiesMap["oversized_data_oss_target_arn"] = slsPropertiesRaw["OversizedDataOssTargetArn"]

			slsPropertiesMaps = append(slsPropertiesMaps, slsPropertiesMap)
		}
		resourceSnapshotDeliveryMap["sls_properties"] = slsPropertiesMaps
		resourceSnapshotDeliveryMaps = append(resourceSnapshotDeliveryMaps, resourceSnapshotDeliveryMap)
	}
	if err := d.Set("resource_snapshot_delivery", resourceSnapshotDeliveryMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudResourceManagerDeliveryChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateDeliveryChannel"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DeliveryChannelId"] = d.Id()

	if d.HasChange("resource_snapshot_delivery.0.custom_expression") {
		update = true
		resourceSnapshotDeliveryCustomExpressionJsonPath, err := jsonpath.Get("$[0].custom_expression", d.Get("resource_snapshot_delivery"))
		if err == nil {
			request["ResourceSnapshotDelivery.CustomExpression"] = resourceSnapshotDeliveryCustomExpressionJsonPath
		}
	}

	if d.HasChange("delivery_channel_name") {
		update = true
	}
	request["DeliveryChannelName"] = d.Get("delivery_channel_name")
	if d.HasChange("resource_change_delivery") {
		update = true
		dataList := make(map[string]interface{})

		if v := d.Get("resource_change_delivery"); v != nil {
			slsProperties := make(map[string]interface{})
			oversizedDataOssTargetArn1, _ := jsonpath.Get("$[0].sls_properties[0].oversized_data_oss_target_arn", d.Get("resource_change_delivery"))
			if oversizedDataOssTargetArn1 != nil && (d.HasChange("resource_change_delivery.0.sls_properties.0.oversized_data_oss_target_arn") || oversizedDataOssTargetArn1 != "") {
				slsProperties["OversizedDataOssTargetArn"] = oversizedDataOssTargetArn1
			}

			dataList["SlsProperties"] = slsProperties
			enabled1, _ := jsonpath.Get("$[0].enabled", v)
			if enabled1 != nil && (d.HasChange("resource_change_delivery.0.enabled") || enabled1 != "") {
				dataList["Enabled"] = enabled1
			}
			targetArn1, _ := jsonpath.Get("$[0].target_arn", v)
			if targetArn1 != nil && (d.HasChange("resource_change_delivery.0.target_arn") || targetArn1 != "") {
				dataList["TargetArn"] = targetArn1
			}
			targetType1, _ := jsonpath.Get("$[0].target_type", v)
			if targetType1 != nil && (d.HasChange("resource_change_delivery.0.target_type") || targetType1 != "") {
				dataList["TargetType"] = targetType1
			}

			request["ResourceChangeDelivery"] = dataList
		}
	}

	request["ResourceChangeDelivery.Enabled"] = "true"
	request["ResourceSnapshotDelivery.Enabled"] = "true"
	if d.HasChange("resource_snapshot_delivery") {
		update = true
		dataList1 := make(map[string]interface{})

		if v := d.Get("resource_snapshot_delivery"); v != nil {
			customExpression1, _ := jsonpath.Get("$[0].custom_expression", v)
			if customExpression1 != nil && (d.HasChange("resource_snapshot_delivery.0.custom_expression") || customExpression1 != "") {
				dataList1["CustomExpression"] = customExpression1
			}
			enabled3, _ := jsonpath.Get("$[0].enabled", v)
			if enabled3 != nil && (d.HasChange("resource_snapshot_delivery.0.enabled") || enabled3 != "") {
				dataList1["Enabled"] = enabled3
			}
			slsProperties1 := make(map[string]interface{})
			oversizedDataOssTargetArn3, _ := jsonpath.Get("$[0].sls_properties[0].oversized_data_oss_target_arn", d.Get("resource_snapshot_delivery"))
			if oversizedDataOssTargetArn3 != nil && (d.HasChange("resource_snapshot_delivery.0.sls_properties.0.oversized_data_oss_target_arn") || oversizedDataOssTargetArn3 != "") {
				slsProperties1["OversizedDataOssTargetArn"] = oversizedDataOssTargetArn3
			}

			dataList1["SlsProperties"] = slsProperties1
			deliveryTime1, _ := jsonpath.Get("$[0].delivery_time", v)
			if deliveryTime1 != nil && (d.HasChange("resource_snapshot_delivery.0.delivery_time") || deliveryTime1 != "") {
				dataList1["DeliveryTime"] = deliveryTime1
			}
			targetType3, _ := jsonpath.Get("$[0].target_type", v)
			if targetType3 != nil && (d.HasChange("resource_snapshot_delivery.0.target_type") || targetType3 != "") {
				dataList1["TargetType"] = targetType3
			}
			targetArn3, _ := jsonpath.Get("$[0].target_arn", v)
			if targetArn3 != nil && (d.HasChange("resource_snapshot_delivery.0.target_arn") || targetArn3 != "") {
				dataList1["TargetArn"] = targetArn3
			}

			request["ResourceSnapshotDelivery"] = dataList1
		}
	}

	if d.HasChange("resource_change_delivery.0.target_arn") {
		update = true
		resourceChangeDeliveryTargetArnJsonPath, err := jsonpath.Get("$[0].target_arn", d.Get("resource_change_delivery"))
		if err == nil {
			request["ResourceChangeDelivery.TargetArn"] = resourceChangeDeliveryTargetArnJsonPath
		}
	}

	if d.HasChange("resource_snapshot_delivery.0.delivery_time") {
		update = true
		resourceSnapshotDeliveryDeliveryTimeJsonPath, err := jsonpath.Get("$[0].delivery_time", d.Get("resource_snapshot_delivery"))
		if err == nil {
			request["ResourceSnapshotDelivery.DeliveryTime"] = resourceSnapshotDeliveryDeliveryTimeJsonPath
		}
	}

	request["ResourceChangeDelivery.TargetType"] = "SLS"
	if d.HasChange("delivery_channel_filter") {
		update = true
	}
	dataList2 := make(map[string]interface{})

	if v := d.Get("delivery_channel_filter"); v != nil {
		resourceTypes1, _ := jsonpath.Get("$[0].resource_types", v)
		if resourceTypes1 != nil && (d.HasChange("delivery_channel_filter.0.resource_types") || resourceTypes1 != "") {
			dataList2["ResourceTypes"] = resourceTypes1
		}

		request["DeliveryChannelFilter"] = dataList2
	}

	if d.HasChange("resource_snapshot_delivery.0.target_type") {
		update = true
		resourceSnapshotDeliveryTargetTypeJsonPath, err := jsonpath.Get("$[0].target_type", d.Get("resource_snapshot_delivery"))
		if err == nil {
			request["ResourceSnapshotDelivery.TargetType"] = resourceSnapshotDeliveryTargetTypeJsonPath
		}
	}

	if d.HasChange("delivery_channel_description") {
		update = true
		request["DeliveryChannelDescription"] = d.Get("delivery_channel_description")
	}

	if d.HasChange("resource_snapshot_delivery.0.target_arn") {
		update = true
		resourceSnapshotDeliveryTargetArnJsonPath, err := jsonpath.Get("$[0].target_arn", d.Get("resource_snapshot_delivery"))
		if err == nil {
			request["ResourceSnapshotDelivery.TargetArn"] = resourceSnapshotDeliveryTargetArnJsonPath
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceCenter", "2022-12-01", action, query, request, true)
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

	return resourceAliCloudResourceManagerDeliveryChannelRead(d, meta)
}

func resourceAliCloudResourceManagerDeliveryChannelDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeliveryChannel"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DeliveryChannelId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceCenter", "2022-12-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"NotExists.DeliveryChannelId"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
