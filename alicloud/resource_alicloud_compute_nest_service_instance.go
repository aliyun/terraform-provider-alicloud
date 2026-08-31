package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudComputeNestServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudComputeNestServiceInstanceCreate,
		Read:   resourceAlicloudComputeNestServiceInstanceRead,
		Update: resourceAlicloudComputeNestServiceInstanceUpdate,
		Delete: resourceAlicloudComputeNestServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_instance_ops": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"specification_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Permanent", "Subscription", "PayAsYouGo", "CustomFixTime"}, false),
			},
			"enable_user_prometheus": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"operation_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation_start_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"operation_end_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resources": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"operated_service_instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"commodity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pay_period": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pay_period_unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Year", "Month", "Day"}, false),
						},
					},
				},
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudComputeNestServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	computeNestService := ComputeNestService{client}
	var response map[string]interface{}
	action := "CreateServiceInstance"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateServiceInstance")
	request["ServiceId"] = d.Get("service_id")
	request["ServiceVersion"] = d.Get("service_version")

	if v, ok := d.GetOk("service_instance_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}

	if v, ok := d.GetOkExists("enable_instance_ops"); ok {
		request["EnableInstanceOps"] = v
	}

	if v, ok := d.GetOk("template_name"); ok {
		request["TemplateName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("specification_name"); ok {
		request["SpecificationName"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		switch v.(string) {
		case "Permanent":
			request["PayType"] = 0
		case "Subscription":
			request["PayType"] = 1
		case "PayAsYouGo":
			request["PayType"] = 2
		case "CustomFixTime":
			request["PayType"] = 3
		}
	}

	if v, ok := d.GetOkExists("enable_user_prometheus"); ok {
		request["EnableUserPrometheus"] = v
	}

	if v, ok := d.GetOk("operation_metadata"); ok {
		operationMetadataMap := map[string]interface{}{}
		for _, operationMetadataList := range v.([]interface{}) {
			operationMetadataArg := operationMetadataList.(map[string]interface{})

			if startTime, ok := operationMetadataArg["operation_start_time"]; ok && startTime.(string) != "" {
				operationMetadataMap["StartTime"] = startTime
			}

			if endTime, ok := operationMetadataArg["operation_end_time"]; ok && endTime.(string) != "" {
				operationMetadataMap["EndTime"] = endTime
			}

			if resources, ok := operationMetadataArg["resources"]; ok && resources.(string) != "" {
				operationMetadataMap["Resources"] = resources
			}

			if serviceInstanceId, ok := operationMetadataArg["operated_service_instance_id"]; ok && serviceInstanceId.(string) != "" {
				operationMetadataMap["ServiceInstanceId"] = serviceInstanceId
			}
		}

		request["OperationMetadata"] = operationMetadataMap
	}

	if v, ok := d.GetOk("commodity"); ok {
		commodityMap := map[string]interface{}{}
		for _, commodityList := range v.([]interface{}) {
			commodityArg := commodityList.(map[string]interface{})

			if payPeriod, ok := commodityArg["pay_period"]; ok {
				commodityMap["PayPeriod"] = payPeriod
			}

			if payPeriodUnit, ok := commodityArg["pay_period_unit"]; ok {
				commodityMap["PayPeriodUnit"] = payPeriodUnit
			}
		}

		request["Commodity"] = commodityMap
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ComputeNest", "2021-06-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_compute_nest_service_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceInstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"Deployed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, computeNestService.ComputeNestServiceInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudComputeNestServiceInstanceUpdate(d, meta)
}

func resourceAlicloudComputeNestServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	computeNestService := ComputeNestService{client}

	object, err := computeNestService.DescribeComputeNestServiceInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if service, ok := object["Service"]; ok {
		serviceArg := service.(map[string]interface{})

		if serviceId, ok := serviceArg["ServiceId"]; ok {
			d.Set("service_id", serviceId)
		}

		if version, ok := serviceArg["Version"]; ok {
			d.Set("service_version", version)
		}
	}

	d.Set("service_instance_name", object["Name"])
	d.Set("enable_instance_ops", object["EnableInstanceOps"])
	d.Set("template_name", object["TemplateName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("specification_name", object["PredefinedParameterName"])
	d.Set("payment_type", object["PayType"])
	d.Set("enable_user_prometheus", object["EnableUserPrometheus"])

	operationMetadataMaps := make([]map[string]interface{}, 0)
	operationMetadataMap := make(map[string]interface{})

	if startTime, ok := object["OperationStartTime"]; ok {
		startTime, err := time.Parse(time.RFC3339Nano, startTime.(string))
		if err != nil {
			return WrapError(err)
		}

		operationMetadataMap["operation_start_time"] = fmt.Sprint(startTime.UnixMilli())
	}

	if endTime, ok := object["OperationEndTime"]; ok {
		endTime, err := time.Parse(time.RFC3339Nano, endTime.(string))
		if err != nil {
			return WrapError(err)
		}

		operationMetadataMap["operation_end_time"] = fmt.Sprint(endTime.UnixMilli())
	}

	if resources, ok := object["Resources"]; ok {
		operationMetadataMap["resources"] = resources
	}

	if serviceInstanceId, ok := object["OperatedServiceInstanceId"]; ok {
		operationMetadataMap["operated_service_instance_id"] = serviceInstanceId
	}

	operationMetadataMaps = append(operationMetadataMaps, operationMetadataMap)

	d.Set("operation_metadata", operationMetadataMaps)

	d.Set("status", object["Status"])

	listTagResourcesObject, err := computeNestService.ListTagResources(d.Id(), "serviceinstance")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlicloudComputeNestServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	computeNestService := ComputeNestService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := computeNestService.SetResourceTags(d, "serviceinstance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceId":   d.Id(),
		"ResourceType": "serviceinstance",
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ComputeNest", "2021-06-01", action, nil, changeResourceGroupReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAlicloudComputeNestServiceInstanceRead(d, meta)
}

func resourceAlicloudComputeNestServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	computeNestService := ComputeNestService{client}
	action := "DeleteServiceInstances"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"ClientToken":       buildClientToken("DeleteServiceInstances"),
		"ServiceInstanceId": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("ComputeNest", "2021-06-01", action, nil, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, computeNestService.ComputeNestServiceInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
