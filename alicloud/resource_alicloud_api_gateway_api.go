package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// apiGatewayErrorCodeSample defines the error code sample structure for CreateApi/ModifyApi.
type apiGatewayErrorCodeSample struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Model       string `json:"model,omitempty"`
	Description string `json:"description,omitempty"`
}

// apiGatewayOpenIdConnectConfig defines the OpenID Connect config structure for CreateApi/ModifyApi.
type apiGatewayOpenIdConnectConfig struct {
	IdTokenParamName string `json:"idTokenParamName,omitempty"`
	PublicKeyId      string `json:"publicKeyId,omitempty"`
	OpenIdApiType    string `json:"openIdApiType,omitempty"`
}

// apiGatewayRequestConfigExt extends the request config with governance fields.
type apiGatewayRequestConfigExt struct {
	ApiGatewayRequestConfig
	PostBodyDescription string `json:"postBodyDescription,omitempty"`
	BodyModel           string `json:"bodyModel,omitempty"`
}

// apiGatewayRequestParamExt extends the request param with governance fields.
type apiGatewayRequestParamExt struct {
	ApiGatewayRequestParam
	DemoValue         string `json:"demoValue,omitempty"`
	MaxValue          int64  `json:"maxValue,omitempty"`
	MinValue          int64  `json:"minValue,omitempty"`
	MaxLength         int64  `json:"maxLength,omitempty"`
	MinLength         int64  `json:"minLength,omitempty"`
	RegularExpression string `json:"regularExpression,omitempty"`
	JsonScheme        string `json:"jsonScheme,omitempty"`
	EnumValue         string `json:"enumValue,omitempty"`
	DocShow           string `json:"docShow,omitempty"`
	DocOrder          int    `json:"docOrder,omitempty"`
	ArrayItemsType    string `json:"arrayItemsType,omitempty"`
}

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

			"app_code_auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "HEADER", "HEADER_QUERY"}, false),
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
						"post_body_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"body_model": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"arn_role": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
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
						"qualifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"content_type_category": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "CUSTOM", "CLIENT"}, false),
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
				Set:      apiGatewayApiRequestParamHash,
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
						"demo_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"max_value": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"min_value": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"min_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"regular_expression": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"json_scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enum_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"doc_show": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"doc_order": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"array_items_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
				Set:      apiGatewayApiSystemParamHash,
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
						"demo_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"stage_names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"backend_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backend_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Governance fields: hardcoded constants converted to Optional+Computed with fallback.
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allow_signature_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"result_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"result_sample": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"web_socket_api_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Governance fields: new writable properties.
			"fail_result_sample": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_body_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disable_internet": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"open_id_connect_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id_token_param_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"open_id_api_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"error_code_samples": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"message": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": tagsSchemaForceNew(),

			// Governance fields: readOnly Computed from DescribeApi response.
			"model_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_valid_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployed_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployed_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effective_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stage_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"result_descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"custom_system_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_parameter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"traffic_control_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_control_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_backend_service": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"signature_id": {
				Type:     schema.TypeString,
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
	d.Set("app_code_auth_type", objectRaw["AppCodeAuthType"])
	d.Set("force_nonce_check", objectRaw["ForceNonceCheck"])

	// Governance fields: converted constants read back from DescribeApi.
	d.Set("visibility", objectRaw["Visibility"])
	d.Set("allow_signature_method", objectRaw["AllowSignatureMethod"])
	d.Set("result_type", objectRaw["ResultType"])
	d.Set("result_sample", objectRaw["ResultSample"])
	d.Set("web_socket_api_type", objectRaw["WebSocketApiType"])
	d.Set("fail_result_sample", objectRaw["FailResultSample"])
	d.Set("result_body_model", objectRaw["ResultBodyModel"])
	d.Set("disable_internet", objectRaw["DisableInternet"])
	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("group_name", objectRaw["GroupName"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("traffic_control_id", objectRaw["TrafficControlId"])
	d.Set("ip_control_id", objectRaw["IpControlId"])
	d.Set("signature_id", objectRaw["SignatureId"])
	d.Set("model_name", objectRaw["ModelName"])
	d.Set("auth_valid_time", objectRaw["AuthValidTime"])
	d.Set("stage_id", objectRaw["StageId"])

	// AppId is returned as float64 in JSON map.
	if v, ok := objectRaw["AppId"]; ok {
		d.Set("app_id", v)
	}

	// OpenIdConnectConfig read back.
	d.Set("open_id_connect_config", convertApiGatewayApiOpenIdConnectConfigResponse(objectRaw["OpenIdConnectConfig"]))

	// ErrorCodeSamples read back.
	d.Set("error_code_samples", convertApiGatewayApiErrorCodeSamplesResponse(objectRaw["ErrorCodeSamples"]))

	// ResultDescriptions read back.
	d.Set("result_descriptions", convertApiGatewayApiResultDescriptionsResponse(objectRaw["ResultDescriptions"]))

	// DeployedInfos read back.
	d.Set("deployed_infos", convertApiGatewayApiDeployedInfosResponse(objectRaw["DeployedInfos"]))

	// CustomSystemParameters read back.
	d.Set("custom_system_parameters", convertApiGatewayApiCustomSystemParametersResponse(objectRaw["CustomSystemParameters"]))

	// Backend config additional fields.
	if backendConfig, ok := objectRaw["BackendConfig"].(map[string]interface{}); ok {
		if v, ok := backendConfig["BackendName"].(string); ok {
			d.Set("backend_name", v)
		}
		if v, ok := backendConfig["Type"].(string); ok {
			d.Set("use_backend_service", v == "HTTP")
		}
	}

	// Tags are read back via ListTagResources per the CloudSpec model get
	// mapping ($.TagResources.TagResource[*]); DescribeApi does not return tags.
	// ListTagResources expects the raw ApiId, not the composite resource id.
	apiId, _ := objectRaw["ApiId"].(string)
	if apiId == "" {
		if parts := strings.Split(d.Id(), COLON_SEPARATED); len(parts) > 1 {
			apiId = parts[1]
		} else {
			apiId = d.Id()
		}
	}
	if tagResources, e := cloudApiService.DescribeTags(apiId, nil, TagResourceType("api")); e != nil {
		if NotFoundError(e) {
			d.Set("tags", cloudApiService.tagsToMap(nil))
		} else {
			// Tags are ForceNew/create-only; on a read failure do not overwrite
			// the configured value to avoid a spurious ForceNew diff.
			addDebug("ListTagResources", e)
		}
	} else {
		d.Set("tags", cloudApiService.tagsToMap(tagResources))
	}

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

	if backendConfig, ok := objectRaw["BackendConfig"].(map[string]interface{}); ok {
		if v, ok := backendConfig["BackendId"].(string); ok && v != "" {
			d.Set("backend_id", v)
		}
	}
	if v, ok := objectRaw["BackendEnable"]; ok {
		d.Set("backend_enabled", v)
	}

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

	if d.HasChanges("name", "description", "auth_type", "app_code_auth_type") {
		update = true
	}
	request.ApiName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.AuthType = d.Get("auth_type").(string)
	if d.HasChange("app_code_auth_type") || d.Get("auth_type").(string) == "APP" {
		if v, exist := d.GetOk("app_code_auth_type"); exist {
			request.AppCodeAuthType = v.(string)
		}
	}

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

	if d.HasChanges("service_type", "http_service_config", "http_vpc_service_config", "mock_service_config", "fc_service_config", "backend_id", "backend_enabled") {
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

	if d.HasChanges("visibility", "allow_signature_method", "result_type", "result_sample", "web_socket_api_type",
		"fail_result_sample", "result_body_model", "disable_internet", "open_id_connect_config", "error_code_samples") {
		update = true
	}
	if update {
		request.ResultType = getApiGatewayApiFieldWithDefault(d, "result_type", ResultType)
		request.ResultSample = getApiGatewayApiFieldWithDefault(d, "result_sample", ResultSample)
		request.Visibility = getApiGatewayApiFieldWithDefault(d, "visibility", Visibility)
		request.AllowSignatureMethod = getApiGatewayApiFieldWithDefault(d, "allow_signature_method", AllowSignatureMethod)
		request.WebSocketApiType = getApiGatewayApiFieldWithDefault(d, "web_socket_api_type", WebSocketApiType)

		if v, ok := d.GetOk("fail_result_sample"); ok && v.(string) != "" {
			request.FailResultSample = v.(string)
		}
		if v, ok := d.GetOk("result_body_model"); ok && v.(string) != "" {
			request.ResultBodyModel = v.(string)
		}
		if v, ok := d.GetOk("disable_internet"); ok {
			request.DisableInternet = requests.Boolean(strconv.FormatBool(v.(bool)))
		}
		if v, ok := d.GetOk("open_id_connect_config"); ok {
			if openIdJson, e := openIdConnectConfigToJsonStr(v.([]interface{})); e == nil && openIdJson != "" {
				request.OpenIdConnectConfig = openIdJson
			}
		}
		if v, ok := d.GetOk("error_code_samples"); ok {
			if ecJson, e := errorCodeSamplesToJsonStr(v.(*schema.Set)); e == nil && ecJson != "" {
				request.ErrorCodeSamples = ecJson
			}
		}

		backendEnabled := d.Get("backend_enabled").(bool)
		request.QueryParams["BackendEnable"] = strconv.FormatBool(backendEnabled)
		if backendEnabled {
			if v, ok := d.GetOk("backend_id"); ok && v.(string) != "" {
				request.QueryParams["BackendId"] = v.(string)
			}
		}

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApi(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	}

	if update || d.HasChange("stage_names") {
		stageNames := d.Get("stage_names").(*schema.Set)
		err = updateApiStages(d, stageNames, meta)
		if err != nil {
			return WrapError(err)
		}
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

	if v, ok := d.GetOk("stage_names"); ok {
		for _, sn := range v.(*schema.Set).List() {
			stageName := sn.(string)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				err := cloudApiService.AbolishApi(d.Id(), stageName)
				if err != nil {
					if IsExpectedErrors(err, []string{"ConcurrencyLockTimeout"}) {
						time.Sleep(3 * time.Second)
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
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
	if v, exist := d.GetOk("app_code_auth_type"); exist {
		request.AppCodeAuthType = v.(string)
	}
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

	request.ResultType = getApiGatewayApiFieldWithDefault(d, "result_type", ResultType)
	request.ResultSample = getApiGatewayApiFieldWithDefault(d, "result_sample", ResultSample)
	request.Visibility = getApiGatewayApiFieldWithDefault(d, "visibility", Visibility)
	request.AllowSignatureMethod = getApiGatewayApiFieldWithDefault(d, "allow_signature_method", AllowSignatureMethod)
	request.WebSocketApiType = getApiGatewayApiFieldWithDefault(d, "web_socket_api_type", WebSocketApiType)

	if v, ok := d.GetOk("fail_result_sample"); ok && v.(string) != "" {
		request.FailResultSample = v.(string)
	}
	if v, ok := d.GetOk("result_body_model"); ok && v.(string) != "" {
		request.ResultBodyModel = v.(string)
	}
	if v, ok := d.GetOk("disable_internet"); ok {
		request.DisableInternet = requests.Boolean(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("open_id_connect_config"); ok {
		if openIdJson, e := openIdConnectConfigToJsonStr(v.([]interface{})); e == nil && openIdJson != "" {
			request.OpenIdConnectConfig = openIdJson
		}
	}
	if v, ok := d.GetOk("error_code_samples"); ok {
		if ecJson, e := errorCodeSamplesToJsonStr(v.(*schema.Set)); e == nil && ecJson != "" {
			request.ErrorCodeSamples = ecJson
		}
	}
	// Tags are mapped to the CreateApi Tag[*] repeatList parameters per the
	// CloudSpec operationMapping (CreateApi Tags→Tag[*].Key/Tag[*].Value).
	// ModifyApi does not accept Tag, so tags are create-time only (ForceNew).
	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		count := 1
		for key, value := range tags {
			request.QueryParams[fmt.Sprintf("Tag.%d.Key", count)] = key
			request.QueryParams[fmt.Sprintf("Tag.%d.Value", count)] = value.(string)
			count++
		}
	}

	backendEnabled := d.Get("backend_enabled").(bool)
	request.QueryParams["BackendEnable"] = strconv.FormatBool(backendEnabled)
	if backendEnabled {
		if v, ok := d.GetOk("backend_id"); ok && v.(string) != "" {
			request.QueryParams["BackendId"] = v.(string)
		}
	}

	return request, WrapError(err)
}

func requestConfigToJsonStr(l []interface{}) (string, error) {
	config := l[0].(map[string]interface{})

	var requestConfig apiGatewayRequestConfigExt
	requestConfig.Protocol = config["protocol"].(string)
	requestConfig.Path = config["path"].(string)
	requestConfig.Method = config["method"].(string)
	requestConfig.Mode = config["mode"].(string)

	if v, ok := config["body_format"]; ok {
		requestConfig.BodyFormat = v.(string)
	}
	if v, ok := config["post_body_description"]; ok {
		requestConfig.PostBodyDescription = v.(string)
	}
	if v, ok := config["body_model"]; ok {
		requestConfig.BodyModel = v.(string)
	}
	configStr, err := json.Marshal(requestConfig)

	return string(configStr), WrapError(err)
}

// getApiGatewayApiFieldWithDefault returns the field value if set, otherwise the default constant.
func getApiGatewayApiFieldWithDefault(d *schema.ResourceData, key, defaultValue string) string {
	if v, ok := d.GetOk(key); ok && v.(string) != "" {
		return v.(string)
	}
	return defaultValue
}

// openIdConnectConfigToJsonStr converts the open_id_connect_config schema block to a JSON string.
func openIdConnectConfigToJsonStr(l []interface{}) (string, error) {
	if len(l) == 0 {
		return "", nil
	}
	config := l[0].(map[string]interface{})
	var cfg apiGatewayOpenIdConnectConfig
	cfg.IdTokenParamName = config["id_token_param_name"].(string)
	cfg.PublicKeyId = config["public_key_id"].(string)
	cfg.OpenIdApiType = config["open_id_api_type"].(string)
	configStr, err := json.Marshal(cfg)
	return string(configStr), WrapError(err)
}

// errorCodeSamplesToJsonStr converts the error_code_samples schema set to a JSON string.
func errorCodeSamplesToJsonStr(s *schema.Set) (string, error) {
	samples := make([]apiGatewayErrorCodeSample, 0)
	for _, element := range s.List() {
		req := element.(map[string]interface{})
		var sample apiGatewayErrorCodeSample
		sample.Code = req["code"].(string)
		sample.Message = req["message"].(string)
		if v, ok := req["description"]; ok {
			sample.Description = v.(string)
		}
		if v, ok := req["model"]; ok {
			sample.Model = v.(string)
		}
		samples = append(samples, sample)
	}
	configStr, err := json.Marshal(samples)
	return string(configStr), WrapError(err)
}

// apiGatewayApiRequestParamHash computes the set identity hash for an API
// Gateway request parameter. It intentionally excludes the Optional+Computed
// governance fields (demo_value, max_value, min_value, max_length,
// min_length, regular_expression, json_scheme, enum_value, doc_show,
// doc_order, array_items_type) that the API populates with defaults when the
// user leaves them unset, so that a parameter written in config without those
// fields still hashes identically to the same parameter read back from the
// API. Without this the TypeSet would show a perpetual diff, because the
// default schema.HashResource hashes Computed fields that are empty in config
// but populated in state. The remaining identity fields are sufficient to
// uniquely distinguish a request parameter within one API.
func apiGatewayApiRequestParamHash(i interface{}) int {
	m, ok := i.(map[string]interface{})
	if !ok {
		return 0
	}
	return schema.HashString(fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v",
		m["name"], m["type"], m["required"], m["in"], m["in_service"],
		m["name_service"], m["description"], m["default_value"]))
}

// apiGatewayApiSystemParamHash computes the set identity hash for an API
// Gateway system parameter. It intentionally excludes demo_value and
// description (Optional+Computed, populated with server defaults by the API
// when the user leaves them unset) so that a system parameter written in
// config without those fields still hashes identically to the same parameter
// read back from the API. The remaining identity fields are sufficient to
// uniquely distinguish a system parameter within one API.
func apiGatewayApiSystemParamHash(i interface{}) int {
	m, ok := i.(map[string]interface{})
	if !ok {
		return 0
	}
	return schema.HashString(fmt.Sprintf("%v|%v|%v",
		m["name"], m["name_service"], m["in"]))
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

	serviceConfig.FcConfig.ContentTypeCategory = config["content_type_category"].(string)
	contentTypeValue := config["content_type_value"].(string)
	if !strings.EqualFold(serviceConfig.FcConfig.ContentTypeCategory, "CLIENT") {
		serviceConfig.FcConfig.ContentTypeValue = contentTypeValue
	}

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
	requestParameters := make([]apiGatewayRequestParamExt, 0)
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

func setSystemParameters(d *schema.ResourceData, requestParameters []apiGatewayRequestParamExt, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]apiGatewayRequestParamExt, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("system_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam apiGatewayRequestParamExt
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
			if v, ok := request["description"]; ok && v.(string) != "" {
				requestParam.Description = v.(string)
			}
			if v, ok := request["demo_value"]; ok && v.(string) != "" {
				requestParam.DemoValue = v.(string)
			}
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

func setConstantParameters(d *schema.ResourceData, requestParameters []apiGatewayRequestParamExt, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]apiGatewayRequestParamExt, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("constant_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam apiGatewayRequestParamExt
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

func setRequestParameters(d *schema.ResourceData, requestParameters []apiGatewayRequestParamExt, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]apiGatewayRequestParamExt, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("request_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam apiGatewayRequestParamExt
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

			if v, ok := request["demo_value"]; ok && v.(string) != "" {
				requestParam.DemoValue = v.(string)
			}
			if v, ok := request["max_value"]; ok {
				requestParam.MaxValue = int64(v.(int))
			}
			if v, ok := request["min_value"]; ok {
				requestParam.MinValue = int64(v.(int))
			}
			if v, ok := request["max_length"]; ok {
				requestParam.MaxLength = int64(v.(int))
			}
			if v, ok := request["min_length"]; ok {
				requestParam.MinLength = int64(v.(int))
			}
			if v, ok := request["regular_expression"]; ok && v.(string) != "" {
				requestParam.RegularExpression = v.(string)
			}
			if v, ok := request["json_scheme"]; ok && v.(string) != "" {
				requestParam.JsonScheme = v.(string)
			}
			if v, ok := request["enum_value"]; ok && v.(string) != "" {
				requestParam.EnumValue = v.(string)
			}
			if v, ok := request["doc_show"]; ok && v.(string) != "" {
				requestParam.DocShow = v.(string)
			}
			if v, ok := request["doc_order"]; ok {
				requestParam.DocOrder = v.(int)
			}
			if v, ok := request["array_items_type"]; ok && v.(string) != "" {
				requestParam.ArrayItemsType = v.(string)
			}

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
	var checkList []string
	if v, ok := d.GetOk("stage_names"); ok {
		for _, sn := range v.(*schema.Set).List() {
			checkList = append(checkList, sn.(string))
		}
	} else {
		// Fallback for ImportState: no config available, check system stages
		checkList = ApiGatewayStageNames
	}
	for _, stageName := range checkList {
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

	// Deploy to all stages in the user's stage_names list
	for _, stageName := range stageNames.List() {
		err := cloudApiService.DeployedApi(d.Id(), stageName.(string))
		if err != nil {
			return WrapError(err)
		}
		_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName.(string))
		if err != nil {
			return WrapError(err)
		}
	}

	// Abolish from old stages that are no longer in the new list
	o, _ := d.GetChange("stage_names")
	if oldSet, ok := o.(*schema.Set); ok {
		for _, stageName := range oldSet.Difference(stageNames).List() {
			err := cloudApiService.AbolishApi(d.Id(), stageName.(string))
			if err != nil {
				return WrapError(err)
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
	if v, ok := requestConfigMap["PostBodyDescription"]; ok {
		requestConfig["post_body_description"] = v
	}
	if v, ok := requestConfigMap["BodyModel"]; ok {
		requestConfig["body_model"] = v
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
		fcServiceConfig["content_type_value"] = fcConfig["ContentTypeValue"]
		fcServiceConfig["content_type_category"] = fcConfig["ContentTypeCatagory"]
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
				if v, ok := requestParamMap["DemoValue"]; ok && v != "" {
					param["demo_value"] = v
				}
				if v, ok := requestParamMap["MaxValue"]; ok {
					param["max_value"] = v
				}
				if v, ok := requestParamMap["MinValue"]; ok {
					param["min_value"] = v
				}
				if v, ok := requestParamMap["MaxLength"]; ok {
					param["max_length"] = v
				}
				if v, ok := requestParamMap["MinLength"]; ok {
					param["min_length"] = v
				}
				if v, ok := requestParamMap["RegularExpression"]; ok && v != "" {
					param["regular_expression"] = v
				}
				if v, ok := requestParamMap["JsonScheme"]; ok && v != "" {
					param["json_scheme"] = v
				}
				if v, ok := requestParamMap["EnumValue"]; ok && v != "" {
					param["enum_value"] = v
				}
				if v, ok := requestParamMap["DocShow"]; ok && v != "" {
					param["doc_show"] = v
				}
				if v, ok := requestParamMap["DocOrder"]; ok {
					param["doc_order"] = v
				}
				if v, ok := requestParamMap["ArrayItemsType"]; ok && v != "" {
					param["array_items_type"] = v
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
		if v, ok := systemParamMap["DemoValue"]; ok && v != "" {
			param["demo_value"] = v
		}
		if v, ok := systemParamMap["Description"]; ok && v != "" {
			param["description"] = v
		}
		systemParams = append(systemParams, param)
	}

	return systemParams
}

func convertApiGatewayApiOpenIdConnectConfigResponse(source interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if source == nil {
		return result
	}
	var cfg map[string]interface{}
	switch s := source.(type) {
	case map[string]interface{}:
		cfg = s
	case string:
		// DescribeApi may return OpenIdConnectConfig as the raw JSON string
		// submitted to CreateApi (camelCase keys) rather than a parsed object.
		if s == "" {
			return result
		}
		if err := json.Unmarshal([]byte(s), &cfg); err != nil {
			return result
		}
	default:
		return result
	}
	if len(cfg) == 0 {
		return result
	}
	m := map[string]interface{}{}
	for _, key := range []string{"IdTokenParamName", "idTokenParamName"} {
		if v, ok := cfg[key]; ok && v != nil {
			m["id_token_param_name"] = v
			break
		}
	}
	for _, key := range []string{"PublicKeyId", "publicKeyId"} {
		if v, ok := cfg[key]; ok && v != nil {
			m["public_key_id"] = v
			break
		}
	}
	for _, key := range []string{"OpenIdApiType", "openIdApiType"} {
		if v, ok := cfg[key]; ok && v != nil {
			m["open_id_api_type"] = v
			break
		}
	}
	result = append(result, m)
	return result
}

func convertApiGatewayApiErrorCodeSamplesResponse(source interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if source == nil {
		return result
	}
	ecsMap, ok := source.(map[string]interface{})
	if !ok {
		return result
	}
	ecs, ok := ecsMap["ErrorCodeSample"].([]interface{})
	if !ok {
		return result
	}
	for _, ec := range ecs {
		ecMap, ok := ec.(map[string]interface{})
		if !ok {
			continue
		}
		m := map[string]interface{}{}
		m["code"] = ecMap["Code"]
		m["message"] = ecMap["Message"]
		if v, ok := ecMap["Description"]; ok && v != "" {
			m["description"] = v
		}
		if v, ok := ecMap["Model"]; ok && v != "" {
			m["model"] = v
		}
		result = append(result, m)
	}
	return result
}

func convertApiGatewayApiResultDescriptionsResponse(source interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if source == nil {
		return result
	}
	rdMap, ok := source.(map[string]interface{})
	if !ok {
		return result
	}
	rds, ok := rdMap["ResultDescription"].([]interface{})
	if !ok {
		return result
	}
	for _, rd := range rds {
		rdM, ok := rd.(map[string]interface{})
		if !ok {
			continue
		}
		m := map[string]interface{}{}
		if v, ok := rdM["Key"]; ok {
			m["key"] = v
		}
		if v, ok := rdM["Name"]; ok {
			m["name"] = v
		}
		if v, ok := rdM["Type"]; ok {
			m["type"] = v
		}
		if v, ok := rdM["Description"]; ok {
			m["description"] = v
		}
		result = append(result, m)
	}
	return result
}

func convertApiGatewayApiDeployedInfosResponse(source interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if source == nil {
		return result
	}
	diMap, ok := source.(map[string]interface{})
	if !ok {
		return result
	}
	dis, ok := diMap["DeployedInfo"].([]interface{})
	if !ok {
		return result
	}
	for _, di := range dis {
		diM, ok := di.(map[string]interface{})
		if !ok {
			continue
		}
		m := map[string]interface{}{}
		if v, ok := diM["DeployedStatus"]; ok {
			m["deployed_status"] = v
		}
		if v, ok := diM["EffectiveVersion"]; ok {
			m["effective_version"] = v
		}
		if v, ok := diM["StageName"]; ok {
			m["stage_name"] = v
		}
		result = append(result, m)
	}
	return result
}

func convertApiGatewayApiCustomSystemParametersResponse(source interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if source == nil {
		return result
	}
	cspMap, ok := source.(map[string]interface{})
	if !ok {
		return result
	}
	csps, ok := cspMap["CustomSystemParameter"].([]interface{})
	if !ok {
		return result
	}
	for _, csp := range csps {
		cspM, ok := csp.(map[string]interface{})
		if !ok {
			continue
		}
		m := map[string]interface{}{}
		if v, ok := cspM["ParameterName"]; ok {
			m["parameter_name"] = v
		}
		if v, ok := cspM["ServiceParameterName"]; ok {
			m["service_parameter_name"] = v
		}
		if v, ok := cspM["Location"]; ok {
			m["location"] = v
		}
		result = append(result, m)
	}
	return result
}
