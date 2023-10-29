// This file is auto-generated, don't edit it. Thanks.
package client

import (
	sls_util "github.com/alibabacloud-go/alibabacloud-gateway-sls-util/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	array "github.com/alibabacloud-go/darabonba-array/client"
	encodeutil "github.com/alibabacloud-go/darabonba-encode-util/client"
	map_ "github.com/alibabacloud-go/darabonba-map/client"
	signatureutil "github.com/alibabacloud-go/darabonba-signature-util/client"
	string_ "github.com/alibabacloud-go/darabonba-string/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Client struct {
	spi.Client
}

func NewClient() (*Client, error) {
	client := new(Client)
	err := client.Init()
	return client, err
}

func (client *Client) Init() (_err error) {
	_err = client.Client.Init()
	if _err != nil {
		return _err
	}
	return nil
}

func (client *Client) ModifyConfiguration(context *spi.InterceptorContext, attributeMap *spi.AttributeMap) (_err error) {
	config := context.Configuration
	config.Endpoint, _err = client.GetEndpoint(config.RegionId, config.Network, config.Endpoint)
	if _err != nil {
		return _err
	}

	return _err
}

func (client *Client) ModifyRequest(context *spi.InterceptorContext, attributeMap *spi.AttributeMap) (_err error) {
	request := context.Request
	hostMap := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(request.HostMap)) {
		hostMap = request.HostMap
	}

	project := hostMap["project"]
	config := context.Configuration
	credential := request.Credential
	accessKeyId, _err := credential.GetAccessKeyId()
	if _err != nil {
		return _err
	}

	accessKeySecret, _err := credential.GetAccessKeySecret()
	if _err != nil {
		return _err
	}

	securityToken, _err := credential.GetSecurityToken()
	if _err != nil {
		return _err
	}

	if !tea.BoolValue(util.Empty(accessKeyId)) {
		request.Headers["x-log-signaturemethod"] = tea.String("hmac-sha1")
	}

	if !tea.BoolValue(util.Empty(securityToken)) {
		request.Headers["x-acs-security-token"] = securityToken
	}

	if !tea.BoolValue(util.IsUnset(request.Body)) {
		if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("protobuf"))) {
			// var bodyMap = Util.assertAsMap(request.body);
			// 缺少body的Content-MD5计算，以及protobuf处理
			request.Headers["content-type"] = tea.String("application/x-protobuf")
		} else if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("json"))) {
			bodyStr := util.ToJSONString(request.Body)
			request.Headers["content-md5"] = string_.ToUpper(encodeutil.HexEncode(signatureutil.MD5Sign(bodyStr)))
			request.Stream = tea.ToReader(bodyStr)
			request.Headers["content-type"] = tea.String("application/json")
		} else if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("formData"))) {
			str := util.ToJSONString(request.Body)
			request.Headers["content-md5"] = string_.ToUpper(encodeutil.HexEncode(signatureutil.MD5Sign(str)))
			request.Stream = tea.ToReader(str)
			request.Headers["content-type"] = tea.String("application/json")
		}

	}

	host, _err := client.GetHost(config.Network, project, config.Endpoint)
	if _err != nil {
		return _err
	}

	request.Headers = tea.Merge(map[string]*string{
		"accept":            tea.String("application/json"),
		"host":              host,
		"date":              util.GetDateUTCString(),
		"user-agent":        request.UserAgent,
		"x-log-apiversion":  tea.String("0.6.0"),
		"x-log-bodyrawsize": tea.String("0"),
	}, request.Headers)
	request.Headers["authorization"], _err = client.GetAuthorization(request.Pathname, request.Method, request.Query, request.Headers, accessKeyId, accessKeySecret)
	if _err != nil {
		return _err
	}

	_err = client.BuildRequest(context)
	if _err != nil {
		return _err
	}
	return _err
}

func (client *Client) ModifyResponse(context *spi.InterceptorContext, attributeMap *spi.AttributeMap) (_err error) {
	request := context.Request
	response := context.Response
	if tea.BoolValue(util.Is4xx(response.StatusCode)) || tea.BoolValue(util.Is5xx(response.StatusCode)) {
		error, _err := util.ReadAsJSON(response.Body)
		if _err != nil {
			return _err
		}

		resMap, _err := util.AssertAsMap(error)
		if _err != nil {
			return _err
		}

		_err = tea.NewSDKError(map[string]interface{}{
			"code":    resMap["errorCode"],
			"message": resMap["errorMessage"],
			"data": map[string]interface{}{
				"httpCode":   tea.IntValue(response.StatusCode),
				"requestId":  tea.StringValue(response.Headers["x-log-requestid"]),
				"statusCode": tea.IntValue(response.StatusCode),
			},
		})
		return _err
	}

	if !tea.BoolValue(util.IsUnset(response.Body)) {
		bodyrawSize := response.Headers["x-log-bodyrawsize"]
		compressType := response.Headers["x-log-compresstype"]
		uncompressedData := response.Body
		if !tea.BoolValue(util.IsUnset(bodyrawSize)) && !tea.BoolValue(util.IsUnset(compressType)) {
			uncompressedData, _err = sls_util.ReadAndUncompressBlock(response.Body, compressType, bodyrawSize)
			if _err != nil {
				return _err
			}

		}

		if tea.BoolValue(util.EqualString(request.BodyType, tea.String("binary"))) {
			response.DeserializedBody = uncompressedData
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("byte"))) {
			byt, _err := util.ReadAsBytes(uncompressedData)
			if _err != nil {
				return _err
			}

			response.DeserializedBody = byt
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("string"))) {
			response.DeserializedBody, _err = util.ReadAsString(uncompressedData)
			if _err != nil {
				return _err
			}

		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("json"))) {
			obj, _err := util.ReadAsJSON(uncompressedData)
			if _err != nil {
				return _err
			}

			// var res = Util.assertAsMap(obj);
			response.DeserializedBody = obj
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("array"))) {
			response.DeserializedBody, _err = util.ReadAsJSON(uncompressedData)
			if _err != nil {
				return _err
			}

		} else {
			response.DeserializedBody, _err = util.ReadAsString(uncompressedData)
			if _err != nil {
				return _err
			}

		}

	}

	return _err
}

func (client *Client) GetEndpoint(regionId *string, network *string, endpoint *string) (_result *string, _err error) {
	if !tea.BoolValue(util.Empty(endpoint)) {
		_result = endpoint
		return _result, _err
	}

	if tea.BoolValue(util.Empty(regionId)) {
		regionId = tea.String("cn-hangzhou")
	}

	if !tea.BoolValue(util.Empty(network)) {
		if tea.BoolValue(string_.Equals(network, tea.String("intranet"))) {
			_result = tea.String(tea.StringValue(regionId) + "-intranet.log.aliyuncs.com")
			return _result, _err
		} else if tea.BoolValue(string_.Equals(network, tea.String("accelerate"))) {
			_result = tea.String("log-global.aliyuncs.com")
			return _result, _err
		} else if tea.BoolValue(string_.Equals(network, tea.String("share"))) {
			if tea.BoolValue(string_.Equals(regionId, tea.String("cn-hangzhou-corp"))) || tea.BoolValue(string_.Equals(regionId, tea.String("cn-shanghai-corp"))) {
				_result = tea.String(tea.StringValue(regionId) + ".sls.aliyuncs.com")
				return _result, _err
			} else if tea.BoolValue(string_.Equals(regionId, tea.String("cn-zhangbei-corp"))) {
				_result = tea.String("zhangbei-corp-share.log.aliyuncs.com")
				return _result, _err
			}

			_result = tea.String(tea.StringValue(regionId) + "-share.log.aliyuncs.com")
			return _result, _err
		}

	}

	_result = tea.String(tea.StringValue(regionId) + ".log.aliyuncs.com")
	return _result, _err
}

func (client *Client) GetHost(network *string, project *string, endpoint *string) (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(project)) {
		_result = endpoint
		return _result, _err
	}

	_result = tea.String(tea.StringValue(project) + "." + tea.StringValue(endpoint))
	return _result, _err
}

func (client *Client) GetAuthorization(pathname *string, method *string, query map[string]*string, headers map[string]*string, ak *string, secret *string) (_result *string, _err error) {
	sign, _err := client.GetSignature(pathname, method, query, headers, secret)
	if _err != nil {
		return _result, _err
	}

	_result = tea.String("LOG " + tea.StringValue(ak) + ":" + tea.StringValue(sign))
	return _result, _err
}

func (client *Client) GetSignature(pathname *string, method *string, query map[string]*string, headers map[string]*string, secret *string) (_result *string, _err error) {
	resource := pathname
	stringToSign := tea.String("")
	canonicalizedResource, _err := client.BuildCanonicalizedResource(resource, query)
	if _err != nil {
		return _result, _err
	}

	canonicalizedHeaders, _err := client.BuildCanonicalizedHeaders(headers)
	if _err != nil {
		return _result, _err
	}

	stringToSign = tea.String(tea.StringValue(method) + "\n" + tea.StringValue(canonicalizedHeaders) + tea.StringValue(canonicalizedResource))
	_body := encodeutil.Base64EncodeToString(signatureutil.HmacSHA1Sign(stringToSign, secret))
	_result = _body
	return _result, _err
}

func (client *Client) BuildCanonicalizedResource(pathname *string, query map[string]*string) (_result *string, _err error) {
	canonicalizedResource := pathname
	paramsMap := tea.Merge(query)
	if !tea.BoolValue(util.Empty(pathname)) {
		paths := string_.Split(pathname, tea.String("?"), tea.Int(2))
		canonicalizedResource = paths[0]
		if tea.BoolValue(util.EqualNumber(array.Size(paths), tea.Int(2))) {
			params := string_.Split(paths[1], tea.String("&"), nil)
			for _, sub := range params {
				item := string_.Split(sub, tea.String("="), nil)
				key := item[0]
				var value *string
				if tea.BoolValue(util.EqualNumber(array.Size(item), tea.Int(2))) {
					value = item[1]
				}

				paramsMap[tea.StringValue(key)] = value
			}
		}

	}

	if !tea.BoolValue(util.IsUnset(paramsMap)) {
		queryList := map_.KeySet(paramsMap)
		sortedParams := array.AscSort(queryList)
		separator := tea.String("?")
		for _, paramName := range sortedParams {
			canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + tea.StringValue(separator) + tea.StringValue(paramName))
			paramValue := paramsMap[tea.StringValue(paramName)]
			if !tea.BoolValue(util.IsUnset(paramValue)) {
				canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + "=" + tea.StringValue(paramValue))
			}

			separator = tea.String("&")
		}
	}

	_result = canonicalizedResource
	return _result, _err
}

func (client *Client) BuildCanonicalizedHeaders(headers map[string]*string) (_result *string, _err error) {
	canonicalizedHeaders := tea.String("")
	contentType := headers["content-type"]
	if tea.BoolValue(util.IsUnset(contentType)) {
		contentType = tea.String("")
	}

	contentMd5 := headers["content-md5"]
	if tea.BoolValue(util.IsUnset(contentMd5)) {
		contentMd5 = tea.String("")
	}

	canonicalizedHeaders = tea.String(tea.StringValue(canonicalizedHeaders) + tea.StringValue(contentMd5) + "\n" + tea.StringValue(contentType) + "\n" + tea.StringValue(headers["date"]) + "\n")
	keys := map_.KeySet(headers)
	sortedHeaders := array.AscSort(keys)
	for _, header := range sortedHeaders {
		if tea.BoolValue(string_.Contains(string_.ToLower(header), tea.String("x-log-"))) || tea.BoolValue(string_.Contains(string_.ToLower(header), tea.String("x-acs-"))) {
			canonicalizedHeaders = tea.String(tea.StringValue(canonicalizedHeaders) + tea.StringValue(header) + ":" + tea.StringValue(headers[tea.StringValue(header)]) + "\n")
		}

	}
	_result = canonicalizedHeaders
	return _result, _err
}

func (client *Client) BuildRequest(context *spi.InterceptorContext) (_err error) {
	request := context.Request
	resource := request.Pathname
	if !tea.BoolValue(util.Empty(resource)) {
		paths := string_.Split(resource, tea.String("?"), tea.Int(2))
		resource = paths[0]
		if tea.BoolValue(util.EqualNumber(array.Size(paths), tea.Int(2))) {
			params := string_.Split(paths[1], tea.String("&"), nil)
			for _, sub := range params {
				item := string_.Split(sub, tea.String("="), nil)
				key := item[0]
				var value *string
				if tea.BoolValue(util.EqualNumber(array.Size(item), tea.Int(2))) {
					value = item[1]
				}

				request.Query[tea.StringValue(key)] = value
			}
		}

	}

	request.Pathname = resource
	return _err
}
