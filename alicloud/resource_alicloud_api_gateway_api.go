package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliyunApigatewayApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayApiCreate,
		Read:   resourceAliyunApigatewayApiRead,
		Update: resourceAliyunApigatewayApiUpdate,
		Delete: resourceAliyunApigatewayApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},

			"auth_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"request_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MOCK", "HTTP-VPC", "FunctionCompute", "HTTP"}, false),
			},

			"http_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type_category": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"content_type_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"http_vpc_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_scheme": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type_category": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"content_type_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"fc_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "2.0",
						},
						"function_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"FCEvent", "HttpTrigger"}, false),
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"function_base_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"only_business_path": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"arn_role": {
							Type:     schema.TypeString,
							Required: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"mock_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:     schema.TypeString,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"request_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"required": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"constant_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"system_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"stage_names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"PRE", "RELEASE", "TEST"}, false),
				},
				Optional: true,
			},

			"api_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_nonce_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewayApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildAliyunApiArgs(d, meta)
	request.RegionId = client.RegionId
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.CreateApi(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apigateway_api", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cloudapi.CreateApiResponse)

	d.SetId(fmt.Sprintf("%s%s%s", request.GroupId, COLON_SEPARATED, response.ApiId))

	if l, ok := d.GetOk("stage_names"); ok {
		err = updateApiStages(d, l.(*schema.Set), meta)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}
	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayApi(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_api DescribeApiGatewayApi Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	stageNames, err := getStageNameList(d, cloudApiService)
	if err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
	}
	if err := d.Set("stage_names", stageNames); err != nil {
		return WrapError(err)
	}

	d.Set("api_id", objectRaw["ApiId"])
	d.Set("group_id", objectRaw["GroupId"])
	d.Set("name", objectRaw["ApiName"])
	d.Set("description", objectRaw["Description"])
	d.Set("auth_type", objectRaw["AuthType"])
	d.Set("force_nonce_check", objectRaw["ForceNonceCheck"])

	v := convertApiGatewayApiRequestConfigResponse(objectRaw["RequestConfig"])
	if err := d.Set("request_config", []map[string]interface{}{v}); err != nil {
		return WrapError(err)
	}

	serviceConfig, ok := objectRaw["ServiceConfig"].(map[string]interface{})
	if !ok {
		return WrapError(Error("ApiGateway resource Api service_config is not valid."))
	}
	if mock, ok := serviceConfig["Mock"]; ok && mock == "TRUE" {
		d.Set("service_type", "MOCK")
		v := convertApiGatewayApiServiceConfigMockServiceConfigResponse(serviceConfig)
		if err := d.Set("mock_service_config", []map[string]interface{}{v}); err != nil {
			return WrapError(err)
		}
	} else if vpcEnable, ok := serviceConfig["ServiceVpcEnable"]; ok && vpcEnable == "TRUE" {
		d.Set("service_type", "HTTP-VPC")
		v := convertApiGatewayApiServiceConfigVpcServiceConfigResponse(serviceConfig)
		if err := d.Set("http_vpc_service_config", []map[string]interface{}{v}); err != nil {
			return WrapError(err)
		}
	} else if serviceProtocol, ok := serviceConfig["ServiceProtocol"]; ok && serviceProtocol == "FunctionCompute" {
		d.Set("service_type", "FunctionCompute")
		v := convertApiGatewayApiServiceConfigFcServiceConfigResponse(serviceConfig)
		if err := d.Set("fc_service_config", []map[string]interface{}{v}); err != nil {
			return WrapError(err)
		}
	} else {
		d.Set("service_type", "HTTP")
		v := convertApiGatewayApiServiceConfigHttpServiceConfigResponse(serviceConfig)
		if err := d.Set("http_service_config", []map[string]interface{}{v}); err != nil {
			return WrapError(err)
		}
	}

	d.Set("request_parameters", convertApiGatewayApiRequestParamsResponse(objectRaw))

	d.Set("constant_parameters", convertApiGatewayApiConstantParamsResponse(objectRaw["ConstantParameters"]))

	d.Set("system_parameters", convertApiGatewayApiSystemParamsResponse(objectRaw["SystemParameters"]))

	return nil
}

func resourceAliyunApigatewayApiUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateModifyApiRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	update := false

	d.Partial(true)

	if d.HasChanges("name", "description", "auth_type") {
		update = true
	}
	request.ApiName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.AuthType = d.Get("auth_type").(string)

	if d.HasChange("force_nonce_check") {
		update = true
	}
	if v, exist := d.GetOk("force_nonce_check"); exist {
		request.ForceNonceCheck = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	var paramErr error
	var paramConfig string
	if d.HasChange("request_config") {
		update = true
	}
	paramConfig, paramErr = requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	if paramErr != nil {
		return paramErr
	}
	request.RequestConfig = paramConfig

	if d.HasChanges("service_type", "http_service_config", "http_vpc_service_config", "mock_service_config", "fc_service_config") {
		update = true
	}
	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return WrapError(err)
	}
	request.ServiceConfig = serviceConfig

	if d.HasChanges("request_parameters", "constant_parameters", "system_parameters") {
		update = true
	}
	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return WrapError(err)
	}
	request.RequestParameters = string(rps)
	request.ServiceParameters = string(sps)
	request.ServiceParametersMap = string(spm)

	if update {
		request.ResultType = ResultType
		request.ResultSample = ResultSample
		request.Visibility = Visibility
		request.AllowSignatureMethod = AllowSignatureMethod
		request.WebSocketApiType = WebSocketApiType

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApi(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("auth_type")
		d.SetPartial("service_type")
		d.SetPartial("http_service_config")
		d.SetPartial("http_vpc_service_config")
		d.SetPartial("fc_service_config")
		d.SetPartial("mock_service_config")
		d.SetPartial("request_parameters")
		d.SetPartial("constant_parameters")
		d.SetPartial("system_parameters")

	}

	if update || d.HasChange("stage_names") {
		if l, ok := d.GetOk("stage_names"); ok {
			err = updateApiStages(d, l.(*schema.Set), meta)
			if err != nil {
				return WrapError(err)
			}
		}
		d.SetPartial("stage_names")
	}

	d.Partial(false)

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDeleteApiRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	for _, stageName := range ApiGatewayStageNames {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			err := cloudApiService.AbolishApi(d.Id(), stageName)
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrencyLockTimeout"}) {
					time.Sleep(3 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}
		_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
		if err != nil {
			if !NotFoundError(err) {
				return WrapError(err)
			}
		}
	}

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApi(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApi"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayApi(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunApiArgs(d *schema.ResourceData, meta interface{}) (*cloudapi.CreateApiRequest, error) {

	request := cloudapi.CreateCreateApiRequest()
	request.GroupId = d.Get("group_id").(string)
	request.Description = d.Get("description").(string)
	request.ApiName = d.Get("name").(string)
	request.AuthType = d.Get("auth_type").(string)
	if v, exist := d.GetOk("force_nonce_check"); exist {
		request.ForceNonceCheck = requests.Boolean(strconv.FormatBool(v.(bool)))
	}
	requestConfig, err := requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	if err != nil {
		return request, WrapError(err)
	}
	request.RequestConfig = requestConfig

	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return request, WrapError(err)
	}
	request.ServiceConfig = serviceConfig

	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return request, WrapError(err)
	}

	request.RequestParameters = string(rps)
	request.ServiceParameters = string(sps)
	request.ServiceParametersMap = string(spm)

	request.ResultType = ResultType
	request.ResultSample = ResultSample
	request.Visibility = Visibility
	request.AllowSignatureMethod = AllowSignatureMethod
	request.WebSocketApiType = WebSocketApiType

	return request, WrapError(err)
}

func requestConfigToJsonStr(l []interface{}) (string, error) {
	config := l[0].(map[string]interface{})

	var requestConfig ApiGatewayRequestConfig
	requestConfig.Protocol = config["protocol"].(string)
	requestConfig.Path = config["path"].(string)
	requestConfig.Method = config["method"].(string)
	requestConfig.Mode = config["mode"].(string)

	if v, ok := config["body_format"]; ok {
		requestConfig.BodyFormat = v.(string)
	}
	configStr, err := json.Marshal(requestConfig)

	return string(configStr), WrapError(err)
}

func getHttpServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("http_service_config")
	if !ok {
		return []byte{}, fmt.Errorf("Creating apigatway api error: http_service_config is null")
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.Address = config["address"].(string)
	serviceConfig.Path = config["path"].(string)
	serviceConfig.Method = config["method"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	if v, ok := config["content_type_category"]; ok {
		serviceConfig.ContentTypeCategory = v.(string)
	}
	if v, ok := config["content_type_value"]; ok {
		serviceConfig.ContentTypeValue = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getHttpVpcServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("http_vpc_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: http_vpc_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.VpcConfig.Name = config["name"].(string)
	serviceConfig.Path = config["path"].(string)
	serviceConfig.Method = config["method"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.VpcEnable = "TRUE"
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	if v, ok := config["vpc_scheme"]; ok {
		serviceConfig.VpcConfig.VpcScheme = v.(string)
	}
	if v, ok := config["content_type_category"]; ok {
		serviceConfig.ContentTypeCategory = v.(string)
	}
	if v, ok := config["content_type_value"]; ok {
		serviceConfig.ContentTypeValue = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getFcServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("fc_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: fc_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "FunctionCompute"
	serviceConfig.FcConfig.FunctionVersion = config["function_version"].(string)
	serviceConfig.FcConfig.FunctionType = config["function_type"].(string)
	serviceConfig.FcConfig.FunctionBaseUrl = config["function_base_url"].(string)
	serviceConfig.FcConfig.Path = config["path"].(string)
	serviceConfig.FcConfig.Method = config["method"].(string)
	serviceConfig.FcConfig.OnlyBusinessPath = config["only_business_path"].(bool)
	serviceConfig.FcConfig.Qualifier = config["qualifier"].(string)
	serviceConfig.FcConfig.Region = config["region"].(string)
	serviceConfig.FcConfig.FunctionName = config["function_name"].(string)
	serviceConfig.FcConfig.ServiceName = config["service_name"].(string)
	serviceConfig.FcConfig.Arn = config["arn_role"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getMockServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("mock_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: mock_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.Method = "GET"
	serviceConfig.MockResult = config["result"].(string)
	serviceConfig.MockEnable = "TRUE"
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.Timeout = ApigatewayDefaultTimeout
	serviceConfig.Address = ApigatewayDefaultAddress
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func serviceConfigToJsonStr(d *schema.ResourceData) (string, error) {
	var err error
	var configStr []byte

	serviceType := d.Get("service_type").(string)

	switch serviceType {
	case "HTTP":
		configStr, err = getHttpServiceConfig(d)
		break
	case "HTTP-VPC":
		configStr, err = getHttpVpcServiceConfig(d)
		break
	case "FunctionCompute":
		configStr, err = getFcServiceConfig(d)
		break
	case "MOCK":
		configStr, err = getMockServiceConfig(d)
		break
	}
	if err != nil {
		return "", WrapError(err)
	}
	return string(configStr), nil
}

func setParameters(d *schema.ResourceData) (rps []byte, sps []byte, spm []byte, err error) {
	requestParameters := make([]ApiGatewayRequestParam, 0)
	serviceParameters := make([]ApiGatewayServiceParam, 0)
	serviceParamMaps := make([]ApiGatewayParameterMap, 0)

	requestParameters, serviceParameters, serviceParamMaps = setRequestParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setConstantParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setSystemParameters(d, requestParameters, serviceParameters, serviceParamMaps)

	rps, err = json.Marshal(requestParameters)
	if err != nil {
		err = WrapError(err)
		return
	}
	sps, err = json.Marshal(serviceParameters)
	if err != nil {
		err = WrapError(err)
		return
	}
	spm, err = json.Marshal(serviceParamMaps)
	if err != nil {
		err = WrapError(err)
		return
	}

	return rps, sps, spm, WrapError(err)
}

func setSystemParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("system_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			nameRequest := request["name"].(string)
			nameService := request["name_service"].(string)
			in := request["in"].(string)

			requestParam.Name = nameRequest
			requestParam.ApiParameterName = nameRequest
			requestParam.In = in
			requestParam.Required = "REQUIRED"
			requestParam.Type = "String"
			requestParameters = append(requestParameters, requestParam)

			serviceParam.Type = "String"
			serviceParam.In = in
			serviceParam.Name = nameService
			serviceParam.Catalog = CatalogSystem
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = nameRequest
			serviceParamMap.ServiceParamName = nameService
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func setConstantParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("constant_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			name := request["name"].(string)
			in := request["in"].(string)
			value := request["value"].(string)

			requestParam.Name = name
			requestParam.Required = "REQUIRED"
			requestParam.ApiParameterName = name
			requestParam.In = in
			requestParam.Type = "String"
			if description, ok := request["description"]; ok {
				requestParam.Description = description.(string)
			}
			requestParam.DefualtValue = value
			requestParameters = append(requestParameters, requestParam)

			serviceParam.In = in
			serviceParam.Name = name
			serviceParam.Catalog = CatalogConstant
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = name
			serviceParamMap.ServiceParamName = name
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func setRequestParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("request_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			nameRequest := request["name"].(string)
			paramType := request["type"].(string)
			required := request["required"].(string)
			in := request["in"].(string)

			inService := request["in_service"].(string)
			nameService := request["name_service"].(string)

			if description, ok := request["description"]; ok {
				requestParam.Description = description.(string)
			}
			if defaultValue, ok := request["default_value"]; ok {
				requestParam.DefualtValue = defaultValue.(string)
			}

			requestParam.Name = nameRequest
			requestParam.Required = required
			requestParam.ApiParameterName = nameRequest
			requestParam.In = in
			requestParam.Type = paramType

			requestParameters = append(requestParameters, requestParam)

			serviceParam.Type = paramType
			serviceParam.In = inService
			serviceParam.Name = nameService
			serviceParam.Catalog = CatalogRequest
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = nameRequest
			serviceParamMap.ServiceParamName = nameService
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func getStageNameList(d *schema.ResourceData, cloudApiService CloudApiService) ([]string, error) {
	var stageNames []string

	for _, stageName := range ApiGatewayStageNames {
		_, err := cloudApiService.DescribeDeployedApi(d.Id(), stageName)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return nil, WrapError(err)
		}
		stageNames = append(stageNames, stageName)
	}
	return stageNames, nil
}

func updateApiStages(d *schema.ResourceData, stageNames *schema.Set, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, stageName := range ApiGatewayStageNames {
		if stageNames.Contains(stageName) {
			err := cloudApiService.DeployedApi(d.Id(), stageName)

			if err != nil {
				return WrapError(err)
			}

			_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
			if err != nil {
				return WrapError(err)
			}

		} else {
			err := cloudApiService.AbolishApi(d.Id(), stageName)
			if err != nil {
				return WrapError(err)
			}
			_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
			if err != nil {
				if !NotFoundError(err) {
					return WrapError(err)
				}
			}
		}
	}
	return nil
}

func convertApiGatewayApiRequestConfigResponse(source interface{}) map[string]interface{} {
	requestConfig := map[string]interface{}{}
	if source == nil {
		return requestConfig
	}
	requestConfigMap, ok := source.(map[string]interface{})
	if !ok {
		return requestConfig
	}

	requestConfig["protocol"] = requestConfigMap["RequestProtocol"]
	requestConfig["method"] = requestConfigMap["RequestHttpMethod"]
	requestConfig["path"] = requestConfigMap["RequestPath"]
	requestConfig["mode"] = requestConfigMap["RequestMode"]

	if bodyFormat, ok := requestConfigMap["BodyFormat"]; ok && bodyFormat != "" {
		requestConfig["body_format"] = bodyFormat
	}

	return requestConfig
}

func convertApiGatewayApiServiceConfigMockServiceConfigResponse(serviceConfig map[string]interface{}) map[string]interface{} {
	mockServiceConfig := map[string]interface{}{}
	if serviceConfig == nil {
		return mockServiceConfig
	}
	mockServiceConfig["result"] = serviceConfig["MockResult"]
	mockServiceConfig["aone_name"] = serviceConfig["AoneAppName"]

	return mockServiceConfig
}

func convertApiGatewayApiServiceConfigVpcServiceConfigResponse(serviceConfig map[string]interface{}) map[string]interface{} {
	vpcServiceConfig := map[string]interface{}{}
	if serviceConfig == nil {
		return vpcServiceConfig
	}

	vpcServiceConfig["path"] = serviceConfig["ServicePath"]
	vpcServiceConfig["method"] = serviceConfig["ServiceHttpMethod"]
	vpcServiceConfig["timeout"] = serviceConfig["ServiceTimeout"]
	vpcServiceConfig["aone_name"] = serviceConfig["AoneAppName"]
	vpcServiceConfig["content_type_category"] = serviceConfig["ContentTypeCatagory"]
	vpcServiceConfig["content_type_value"] = serviceConfig["ContentTypeValue"]

	if vpcConfig, ok := serviceConfig["VpcConfig"].(map[string]interface{}); ok {
		vpcServiceConfig["name"] = vpcConfig["Name"]
		vpcServiceConfig["vpc_scheme"] = vpcConfig["VpcScheme"]
	}

	return vpcServiceConfig
}

func convertApiGatewayApiServiceConfigFcServiceConfigResponse(serviceConfig map[string]interface{}) map[string]interface{} {
	fcServiceConfig := map[string]interface{}{}
	if serviceConfig == nil {
		return fcServiceConfig
	}

	fcServiceConfig["timeout"] = serviceConfig["ServiceTimeout"]

	if fcConfig, ok := serviceConfig["FunctionComputeConfig"].(map[string]interface{}); ok {
		fcServiceConfig["region"] = fcConfig["RegionId"]
		fcServiceConfig["function_version"] = fcConfig["FcVersion"]
		fcServiceConfig["function_type"] = fcConfig["FcType"]
		fcServiceConfig["function_base_url"] = fcConfig["FcBaseUrl"]
		fcServiceConfig["path"] = fcConfig["Path"]
		fcServiceConfig["method"] = fcConfig["Method"]
		fcServiceConfig["only_business_path"] = fcConfig["OnlyBusinessPath"]
		fcServiceConfig["qualifier"] = fcConfig["Qualifier"]
		fcServiceConfig["function_name"] = fcConfig["FunctionName"]
		fcServiceConfig["service_name"] = fcConfig["ServiceName"]
		fcServiceConfig["arn_role"] = fcConfig["RoleArn"]
	}

	return fcServiceConfig
}

func convertApiGatewayApiServiceConfigHttpServiceConfigResponse(serviceConfig map[string]interface{}) map[string]interface{} {
	httpServiceConfig := map[string]interface{}{}
	if serviceConfig == nil {
		return httpServiceConfig
	}

	httpServiceConfig["address"] = serviceConfig["ServiceAddress"]
	httpServiceConfig["path"] = serviceConfig["ServicePath"]
	httpServiceConfig["method"] = serviceConfig["ServiceHttpMethod"]
	httpServiceConfig["timeout"] = serviceConfig["ServiceTimeout"]
	httpServiceConfig["aone_name"] = serviceConfig["AoneAppName"]
	httpServiceConfig["content_type_category"] = serviceConfig["ContentTypeCatagory"]
	httpServiceConfig["content_type_value"] = serviceConfig["ContentTypeValue"]

	return httpServiceConfig
}

func convertApiGatewayApiRequestParamsResponse(objectRaw map[string]interface{}) []map[string]interface{} {
	var requestParams []map[string]interface{}
	if objectRaw == nil {
		return requestParams
	}
	serviceParametersMap, ok := objectRaw["ServiceParametersMap"].(map[string]interface{})
	if !ok {
		return requestParams
	}
	serviceParameterMap, ok := serviceParametersMap["ServiceParameterMap"].([]interface{})
	if !ok {
		return requestParams
	}
	for _, mapParam := range serviceParameterMap {
		param := map[string]interface{}{}
		paramMap, ok := mapParam.(map[string]interface{})
		if !ok {
			continue
		}
		requestName := paramMap["RequestParameterName"]
		serviceName := paramMap["ServiceParameterName"]
		serviceParameters, ok := objectRaw["ServiceParameters"].(map[string]interface{})
		if !ok {
			continue
		}
		serviceParameter, ok := serviceParameters["ServiceParameter"].([]interface{})
		if !ok {
			continue
		}
		for _, serviceParam := range serviceParameter {
			serviceParamMap, ok := serviceParam.(map[string]interface{})
			if !ok {
				continue
			}
			if serviceParamMap["ServiceParameterName"] == serviceName {
				param["name_service"] = serviceName
				param["in_service"] = strings.ToUpper(serviceParamMap["Location"].(string))
				break
			}
		}

		requestParameters, ok := objectRaw["RequestParameters"].(map[string]interface{})
		if !ok {
			continue
		}
		requestParameter, ok := requestParameters["RequestParameter"].([]interface{})
		if !ok {
			continue
		}
		for _, requestParam := range requestParameter {
			requestParamMap, ok := requestParam.(map[string]interface{})
			if !ok {
				continue
			}
			if requestParamMap["ApiParameterName"] == requestName {
				param["name"] = requestName
				param["type"] = requestParamMap["ParameterType"]
				param["required"] = requestParamMap["Required"]
				param["in"] = requestParamMap["Location"]

				if description, ok := requestParamMap["Description"]; ok && description != "" {
					param["description"] = description
				}
				if defaultValue, ok := requestParamMap["DefaultValue"]; ok && defaultValue != "" {
					param["default_value"] = defaultValue
				}
				break
			}
		}
		requestParams = append(requestParams, param)
	}

	return requestParams
}

func convertApiGatewayApiConstantParamsResponse(source interface{}) []map[string]interface{} {
	var constantParams []map[string]interface{}
	if source == nil {
		return constantParams
	}
	constantParametersMap, ok := source.(map[string]interface{})
	if !ok {
		return constantParams
	}
	constantParameter, ok := constantParametersMap["ConstantParameter"].([]interface{})
	if !ok {
		return constantParams
	}
	for _, constantParam := range constantParameter {
		param := map[string]interface{}{}
		constantParamMap, ok := constantParam.(map[string]interface{})
		if !ok {
			continue
		}
		param["name"] = constantParamMap["ServiceParameterName"]
		param["in"] = constantParamMap["Location"]
		param["value"] = constantParamMap["ConstantValue"]
		if description, ok := constantParamMap["Description"]; ok && description != "" {
			param["description"] = description
		}
		constantParams = append(constantParams, param)
	}

	return constantParams
}

func convertApiGatewayApiSystemParamsResponse(source interface{}) []map[string]interface{} {
	var systemParams []map[string]interface{}
	if source == nil {
		return systemParams
	}
	systemParametersMap, ok := source.(map[string]interface{})
	if !ok {
		return systemParams
	}
	systemParameter, ok := systemParametersMap["SystemParameter"].([]interface{})
	if !ok {
		return systemParams
	}
	for _, systemParam := range systemParameter {
		param := map[string]interface{}{}
		systemParamMap, ok := systemParam.(map[string]interface{})
		if !ok {
			continue
		}
		param["name"] = systemParamMap["ParameterName"]
		param["in"] = systemParamMap["Location"]
		param["name_service"] = systemParamMap["ServiceParameterName"]
		systemParams = append(systemParams, param)
	}

	return systemParams
}
