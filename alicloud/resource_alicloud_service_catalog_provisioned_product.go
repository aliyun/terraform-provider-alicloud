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

func resourceAlicloudServiceCatalogProvisionedProduct() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudServiceCatalogProvisionedProductCreate,
		Read:   resourceAlicloudServiceCatalogProvisionedProductRead,
		Update: resourceAlicloudServiceCatalogProvisionedProductUpdate,
		Delete: resourceAlicloudServiceCatalogProvisionedProductDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(24 * time.Minute),
			Delete: schema.DefaultTimeout(24 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"last_provisioning_task_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"last_successful_provisioning_task_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"last_task_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"outputs": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"output_key": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"output_value": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"owner_principal_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"owner_principal_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"parameters": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"parameter_value": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"portfolio_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"product_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"product_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"product_version_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"product_version_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"provisioned_product_arn": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"provisioned_product_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"provisioned_product_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"provisioned_product_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"stack_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"stack_region_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status_message": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudServiceCatalogProvisionedProductCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicecatalogService := ServicecatalogService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("portfolio_id"); ok {
		request["PortfolioId"] = v
	}
	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
	}
	if v, ok := d.GetOk("product_version_id"); ok {
		request["ProductVersionId"] = v
	}
	if v, ok := d.GetOk("provisioned_product_name"); ok {
		request["ProvisionedProductName"] = v
	}
	if v, ok := d.GetOk("stack_region_id"); ok {
		request["StackRegionId"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parametersMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			parameters := value0.(map[string]interface{})
			parametersMap := make(map[string]interface{})
			parametersMap["ParameterKey"] = parameters["parameter_key"]
			parametersMap["ParameterValue"] = parameters["parameter_value"]
			parametersMaps = append(parametersMaps, parametersMap)
		}
		request["Parameters"] = parametersMaps
	}

	if v, ok := d.GetOk("tags"); ok {
		request["Tags"] = tagsFromMap(v.(map[string]interface{}))
	}

	var response map[string]interface{}
	action := "LaunchProduct"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("servicecatalog", "2021-09-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"undefined"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_catalog_provisioned_product", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ProvisionedProductId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_service_catalog_provisioned_product")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, servicecatalogService.ServiceCatalogProvisionedProductStateRefreshFunc(d, []string{"Error"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudServiceCatalogProvisionedProductRead(d, meta)
}

func resourceAlicloudServiceCatalogProvisionedProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicecatalogService := ServicecatalogService{client}

	object, err := servicecatalogService.DescribeServiceCatalogProvisionedProduct(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_catalog_provisioned_product servicecatalogService.DescribeServiceCatalogProvisionedProduct Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("last_provisioning_task_id", object["LastProvisioningTaskId"])
	d.Set("last_successful_provisioning_task_id", object["LastSuccessfulProvisioningTaskId"])
	d.Set("last_task_id", object["LastTaskId"])
	outputs67Maps := make([]map[string]interface{}, 0)
	outputs67Raw := object["Outputs"]
	for _, value0 := range outputs67Raw.([]interface{}) {
		outputs67 := value0.(map[string]interface{})
		outputs67Map := make(map[string]interface{})
		outputs67Map["description"] = outputs67["Description"]
		outputs67Map["output_key"] = outputs67["OutputKey"]
		outputs67Map["output_value"] = outputs67["OutputValue"]
		outputs67Maps = append(outputs67Maps, outputs67Map)
	}
	d.Set("outputs", outputs67Maps)
	d.Set("owner_principal_id", object["OwnerPrincipalId"])
	d.Set("owner_principal_type", object["OwnerPrincipalType"])
	parameters31Maps := make([]map[string]interface{}, 0)
	parameters31Raw := object["Parameters"]
	for _, value0 := range parameters31Raw.([]interface{}) {
		parameters31 := value0.(map[string]interface{})
		parameters31Map := make(map[string]interface{})
		parameters31Map["parameter_key"] = parameters31["ParameterKey"]
		parameters31Map["parameter_value"] = parameters31["ParameterValue"]
		parameters31Maps = append(parameters31Maps, parameters31Map)
	}
	d.Set("parameters", parameters31Maps)
	d.Set("portfolio_id", object["PortfolioId"])
	d.Set("product_id", object["ProductId"])
	d.Set("product_name", object["ProductName"])
	d.Set("product_version_id", object["ProductVersionId"])
	d.Set("product_version_name", object["ProductVersionName"])
	d.Set("provisioned_product_arn", object["ProvisionedProductArn"])
	d.Set("provisioned_product_name", object["ProvisionedProductName"])
	d.Set("provisioned_product_type", object["ProvisionedProductType"])
	d.Set("stack_id", object["StackId"])
	d.Set("stack_region_id", object["StackRegionId"])
	d.Set("status", object["Status"])
	d.Set("status_message", object["StatusMessage"])
	d.Set("provisioned_product_id", object["ProvisionedProductId"])
	tagsRaw, err := jsonpath.Get("$.TaskTags", object)
	if err == nil {
		d.Set("tags", tagsToMap(tagsRaw))
	}

	return nil
}

func resourceAlicloudServiceCatalogProvisionedProductUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	servicecatalogService := ServicecatalogService{client}
	var err error
	update := false
	request := map[string]interface{}{
		"ProvisionedProductId": d.Id(),
	}

	if d.HasChange("parameters") {
		update = true
		if v, ok := d.GetOk("parameters"); ok {
			parametersMaps := make([]map[string]interface{}, 0)
			for _, value0 := range v.(*schema.Set).List() {
				parameters := value0.(map[string]interface{})
				parametersMap := make(map[string]interface{})
				parametersMap["ParameterKey"] = parameters["parameter_key"]
				parametersMap["ParameterValue"] = parameters["parameter_value"]
				parametersMaps = append(parametersMaps, parametersMap)
			}
			request["Parameters"] = parametersMaps
		}
	}
	if d.HasChange("portfolio_id") {
		update = true
	}
	request["PortfolioId"] = d.Get("portfolio_id")
	if d.HasChange("product_id") {
		update = true
	}
	request["ProductId"] = d.Get("product_id")
	if d.HasChange("product_version_id") {
		update = true
	}
	request["ProductVersionId"] = d.Get("product_version_id")
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			request["Tags"] = tagsFromMap(v.(map[string]interface{}))
		}
	}

	if update {
		action := "UpdateProvisionedProduct"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("servicecatalog", "2021-09-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, servicecatalogService.ServiceCatalogProvisionedProductStateRefreshFunc(d, []string{"Error"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudServiceCatalogProvisionedProductRead(d, meta)
}

func resourceAlicloudServiceCatalogProvisionedProductDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicecatalogService := ServicecatalogService{client}
	var err error

	request := map[string]interface{}{
		"ProvisionedProductId": d.Id(),
	}

	action := "TerminateProvisionedProduct"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("servicecatalog", "2021-09-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidProvisionedProduct.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, servicecatalogService.ServiceCatalogProvisionedProductStateRefreshFunc(d, []string{"Error"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
