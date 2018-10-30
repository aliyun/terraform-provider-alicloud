package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunApigatewayApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayApiCreate,
		Read:   resourceAliyunApigatewayApiRead,
		Update: resourceAliyunApigatewayApiUpdate,
		Delete: resourceAliyunApigatewayApiDelete,

		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"request_config": &schema.Schema{
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

			"service_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"http_service_config": &schema.Schema{
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
					},
				},
			},

			"http_vpc_service_config": &schema.Schema{
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
					},
				},
			},

			"mock_service_config": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"request_parameters": &schema.Schema{
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

			"constant_parameters": &schema.Schema{
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

			"system_parameters": &schema.Schema{
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
		},
	}
}

func resourceAliyunApigatewayApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildAliyunApiArgs(d, meta)
	if err != nil {
		return err
	}

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.CreateApi(request)
	})
	if err != nil {
		return fmt.Errorf("Creating apigatway api error: %#v", err)
	}

	resp, _ := raw.(*cloudapi.CreateApiResponse)
	d.SetId(resp.ApiId)
	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateDescribeApiRequest()
	request.ApiId = d.Id()
	request.GroupId = d.Get("group_id").(string)
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(request)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	resp, _ := raw.(*cloudapi.DescribeApiResponse)
	d.Set("group_id", resp.GroupId)
	d.Set("name", resp.ApiName)
	d.Set("description", resp.Description)
	d.Set("auth_type", resp.AuthType)

	requestConfig := map[string]interface{}{}
	requestConfig["protocol"] = resp.RequestConfig.RequestProtocol
	requestConfig["method"] = resp.RequestConfig.RequestHttpMethod
	requestConfig["path"] = resp.RequestConfig.RequestPath
	requestConfig["mode"] = resp.RequestConfig.RequestMode
	if resp.RequestConfig.BodyFormat != "" {
		requestConfig["body_format"] = resp.RequestConfig.BodyFormat
	}
	if err := d.Set("request_config", []map[string]interface{}{requestConfig}); err != nil {
		return err
	}

	if resp.ServiceConfig.Mock == "TRUE" {
		d.Set("service_type", "MOCK")
		MockServiceConfig := map[string]interface{}{}
		requestConfig["result"] = resp.ServiceConfig.MockResult
		if err := d.Set("mock_service_config", []map[string]interface{}{MockServiceConfig}); err != nil {
			return err
		}
	} else if resp.ServiceConfig.ServiceVpcEnable == "TRUE" {
		d.Set("service_type", "VPC")
		vpcServiceConfig := map[string]interface{}{}
		vpcServiceConfig["name"] = resp.ServiceConfig.VpcConfig.Name
		vpcServiceConfig["path"] = resp.ServiceConfig.ServicePath
		vpcServiceConfig["method"] = resp.ServiceConfig.ServiceHttpMethod
		vpcServiceConfig["timeout"] = resp.ServiceConfig.ServiceTimeout
		if err := d.Set("http_vpc_service_config", []map[string]interface{}{vpcServiceConfig}); err != nil {
			return err
		}
	} else {
		d.Set("service_type", "HTTP")
		httpServiceConfig := map[string]interface{}{}
		httpServiceConfig["address"] = resp.ServiceConfig.ServiceAddress
		httpServiceConfig["path"] = resp.ServiceConfig.ServicePath
		httpServiceConfig["method"] = resp.ServiceConfig.ServiceHttpMethod
		httpServiceConfig["timeout"] = resp.ServiceConfig.ServiceTimeout
		httpServiceConfig["aone_name"] = resp.ServiceConfig.AoneAppName
		if err := d.Set("http_service_config", []map[string]interface{}{httpServiceConfig}); err != nil {
			return err
		}
	}

	requestParams := []map[string]interface{}{}
	for _, mapParam := range resp.ServiceParametersMap.ServiceParameterMap {
		param := map[string]interface{}{}
		requestName := mapParam.RequestParameterName
		serviceName := mapParam.ServiceParameterName
		for _, serviceParam := range resp.ServiceParameters.ServiceParameter {
			if serviceParam.ServiceParameterName == serviceName {
				param["name_service"] = serviceName
				param["in_service"] = strings.ToUpper(serviceParam.Location)
				break
			}
		}
		for _, requestParam := range resp.RequestParameters.RequestParameter {
			if requestParam.ApiParameterName == requestName {
				param["name"] = requestName
				param["type"] = requestParam.ParameterType
				param["required"] = requestParam.Required
				param["in"] = requestParam.Location

				if requestParam.Description != "" {
					param["description"] = requestParam.Description
				}
				if requestParam.DefaultValue != "" {
					param["default_value"] = requestParam.DefaultValue
				}
				break
			}
		}
		requestParams = append(requestParams, param)
	}
	d.Set("request_parameters", requestParams)

	constantParams := []map[string]interface{}{}
	for _, constantParam := range resp.ConstantParameters.ConstantParameter {
		param := map[string]interface{}{}
		param["name"] = constantParam.ServiceParameterName
		param["in"] = constantParam.Location
		param["value"] = constantParam.ConstantValue
		if constantParam.Description != "" {
			param["description"] = constantParam.Description
		}
		constantParams = append(constantParams, param)

	}
	d.Set("constant_parameters", constantParams)

	SystemParams := []map[string]interface{}{}
	for _, systemParam := range resp.SystemParameters.SystemParameter {
		param := map[string]interface{}{}
		param["name"] = systemParam.ParameterName
		param["in"] = systemParam.Location
		param["name_service"] = systemParam.ServiceParameterName
	}
	d.Set("system_parameters", SystemParams)

	return nil
}

func resourceAliyunApigatewayApiUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	req := cloudapi.CreateModifyApiRequest()
	req.ApiId = d.Id()
	req.GroupId = d.Get("group_id").(string)
	update := false

	if d.HasChange("name") {
		update = true
		_, nv := d.GetChange("name")
		req.ApiName = nv.(string)
	} else {
		req.ApiName = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update = true
		_, nv := d.GetChange("description")
		req.Description = nv.(string)
	} else {
		req.Description = d.Get("description").(string)
	}

	if d.HasChange("auth_type") {
		update = true
		_, nv := d.GetChange("auth_type")
		req.AuthType = nv.(string)
	} else {
		req.AuthType = d.Get("auth_type").(string)
	}

	var paramErr error
	var paramConfig string
	if d.HasChange("request_config") {
		update = true
		_, nv := d.GetChange("auth_type")
		paramConfig, paramErr = requestConfigToJsonStr(nv.([]interface{}))
	} else {
		paramConfig, paramErr = requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	}
	if paramErr != nil {
		return paramErr
	}
	req.RequestConfig = paramConfig

	if d.HasChange("service_type") || d.HasChange("http_service_config") || d.HasChange("http_vpc_service_config") || d.HasChange("mock_service_config") {
		update = true
	}
	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return err
	}
	req.ServiceConfig = serviceConfig

	if d.HasChange("request_parameters") || d.HasChange("constant_parameters") || d.HasChange("system_parameters") {
		update = true
	}
	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return err
	}
	req.RequestParameters = string(rps)
	req.ServiceParameters = string(sps)
	req.ServiceParametersMap = string(spm)

	if update {
		req.ResultType = "JSON"
		req.ResultSample = "Result Sample"
		req.Visibility = "PRIVATE"
		req.AllowSignatureMethod = "HmacSHA256"
		req.WebSocketApiType = "COMMON"

		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApi(req)
		})
		if err != nil {
			return fmt.Errorf("Modify Api got an error: %#v", err)
		}
	}

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	req := cloudapi.CreateDeleteApiRequest()
	req.GroupId = d.Get("group_id").(string)
	req.ApiId = d.Id()
	_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApi(req)
	})
	if err != nil {
		return fmt.Errorf("Error deleting Api failed: %#v", err)
	}
	return nil
}

func buildAliyunApiArgs(d *schema.ResourceData, meta interface{}) (*cloudapi.CreateApiRequest, error) {

	request := cloudapi.CreateCreateApiRequest()
	request.GroupId = d.Get("group_id").(string)
	request.Description = d.Get("description").(string)
	request.ApiName = d.Get("name").(string)
	request.AuthType = d.Get("auth_type").(string)

	requestConfig, err := requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	if err != nil {
		return request, err
	}
	request.RequestConfig = requestConfig

	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return request, err
	}
	request.ServiceConfig = serviceConfig

	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return request, err
	}

	request.RequestParameters = string(rps)
	request.ServiceParameters = string(sps)
	request.ServiceParametersMap = string(spm)

	request.ResultType = "JSON"
	request.ResultSample = "Result Sample"
	request.Visibility = "PRIVATE"
	request.AllowSignatureMethod = "HmacSHA256"
	request.WebSocketApiType = "COMMON"

	return request, err
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

	return string(configStr), err
}

func getHttpServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	if d.HasChange("http_service_config") {
		_, nv := d.GetChange("http_service_config")
		l = nv.([]interface{})
	} else {
		v, ok := d.GetOk("http_service_config")
		if !ok {
			return []byte{}, fmt.Errorf("Creating apigatway api error: http_service_config is null")
		}
		l = v.([]interface{})
	}

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
	configStr, err := json.Marshal(serviceConfig)

	return configStr, err
}

func getHttpVpcServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	if d.HasChange("http_vpc_service_config") {
		_, nv := d.GetChange("http_vpc_service_config")
		l = nv.([]interface{})
	} else {
		v, ok := d.GetOk("http_vpc_service_config")
		if !ok {
			return []byte{}, fmt.Errorf("Creating apigatway api error: http_vpc_service_config is null")
		}
		l = v.([]interface{})
	}

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
	configStr, err := json.Marshal(serviceConfig)

	return configStr, err
}

func getMockServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	if d.HasChange("mock_service_config") {
		_, nv := d.GetChange("mock_service_config")
		l = nv.([]interface{})
	} else {
		v, ok := d.GetOk("mock_service_config")
		if !ok {
			return []byte{}, fmt.Errorf("Creating apigatway api error: mock_service_config is null")
		}
		l = v.([]interface{})
	}

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.Method = config["result"].(string)
	serviceConfig.MockEnable = "TRUE"
	serviceConfig.VpcEnable = "FALSE"
	configStr, err := json.Marshal(serviceConfig)

	return configStr, err
}

func serviceConfigToJsonStr(d *schema.ResourceData) (string, error) {
	var err error
	var configStr []byte

	serviceType := d.Get("service_type").(string)
	if d.HasChange("service_type") {
		_, nv := d.GetChange("service_type")
		serviceType = nv.(string)
	}

	switch serviceType {
	case "HTTP":
		configStr, err = getHttpServiceConfig(d)
		break
	case "HTTP-VPC":
		configStr, err = getHttpVpcServiceConfig(d)
		break
	case "MOCK":
		configStr, err = getMockServiceConfig(d)
		break
	default:
		return "", fmt.Errorf("Creating apigatway api error: unsupport service_type")
	}

	return string(configStr), err
}

func setParameters(d *schema.ResourceData) (rps []byte, sps []byte, spm []byte, err error) {
	var requestParameters []ApiGatewayRequestParam = make([]ApiGatewayRequestParam, 0)
	var serviceParameters []ApiGatewayServiceParam = make([]ApiGatewayServiceParam, 0)
	var serviceParamMaps []ApiGatewayParameterMap = make([]ApiGatewayParameterMap, 0)

	requestParameters, serviceParameters, serviceParamMaps = setRequestParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setConstantParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setSystemParameters(d, requestParameters, serviceParameters, serviceParamMaps)

	rps, err = json.Marshal(requestParameters)
	if err != nil {
		return
	}
	sps, err = json.Marshal(serviceParameters)
	if err != nil {
		return
	}
	spm, err = json.Marshal(serviceParamMaps)
	if err != nil {
		return
	}

	return rps, sps, spm, err
}

func setSystemParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("system_parameters"); ok {
		if d.HasChange("system_parameters") {
			_, nv := d.GetChange("system_parameters")
			l = nv
		}

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
			serviceParam.Catalog = catalogSystem
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
		if d.HasChange("constant_parameters") {
			_, nv := d.GetChange("constant_parameters")
			l = nv
		}

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			name := request["name"].(string)
			paramType := request["type"].(string)
			in := request["in"].(string)
			value := request["value"].(string)

			requestParam.Name = name
			requestParam.Required = "REQUIRED"
			requestParam.ApiParameterName = name
			requestParam.In = in
			requestParam.Type = "String"
			if description, ok := request["description"]; !ok {
				requestParam.Description = description.(string)
			}
			requestParam.DefualtValue = value
			requestParameters = append(requestParameters, requestParam)

			serviceParam.Type = paramType
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
		if d.HasChange("request_parameters") {
			_, nv := d.GetChange("request_parameters")
			l = nv
		}

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
