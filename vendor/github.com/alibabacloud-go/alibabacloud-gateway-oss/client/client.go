// This file is auto-generated, don't edit it. Thanks.
package client

import (
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	array "github.com/alibabacloud-go/darabonba-array/client"
	encodeutil "github.com/alibabacloud-go/darabonba-encode-util/client"
	map_ "github.com/alibabacloud-go/darabonba-map/client"
	signatureutil "github.com/alibabacloud-go/darabonba-signature-util/client"
	string_ "github.com/alibabacloud-go/darabonba-string/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	ossutil "github.com/alibabacloud-go/tea-oss-utils/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	xml "github.com/alibabacloud-go/tea-xml/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Client struct {
	Client                spi.Client
	Default_signed_params []*string
	Except_signed_params  []*string
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
	client.Default_signed_params = []*string{tea.String("location"), tea.String("cors"), tea.String("objectMeta"), tea.String("uploadId"), tea.String("partNumber"), tea.String("security-token"), tea.String("position"), tea.String("img"), tea.String("style"), tea.String("styleName"), tea.String("replication"), tea.String("replicationProgress"), tea.String("replicationLocation"), tea.String("cname"), tea.String("qos"), tea.String("startTime"), tea.String("endTime"), tea.String("symlink"), tea.String("x-oss-process"), tea.String("response-content-type"), tea.String("response-content-language"), tea.String("response-expires"), tea.String("response-cache-control"), tea.String("response-content-disposition"), tea.String("response-content-encoding"), tea.String("udf"), tea.String("udfName"), tea.String("udfImage"), tea.String("udfId"), tea.String("udfImageDesc"), tea.String("udfApplication"), tea.String("udfApplicationLog"), tea.String("restore"), tea.String("callback"), tea.String("callback-var"), tea.String("policy"), tea.String("encryption"), tea.String("versions"), tea.String("versioning"), tea.String("versionId"), tea.String("wormId")}
	client.Except_signed_params = []*string{tea.String("list-type"), tea.String("regions")}
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

	bucketName := hostMap["bucket"]
	if tea.BoolValue(util.IsUnset(bucketName)) {
		bucketName = tea.String("")
	}

	if !tea.BoolValue(util.IsUnset(request.Headers["x-oss-meta-*"])) {
		tmp := util.ParseJSON(request.Headers["x-oss-meta-*"])
		mapData, _err := util.AssertAsMap(tmp)
		if _err != nil {
			return _err
		}

		metaData := util.StringifyMapValue(mapData)
		metaKeySet := map_.KeySet(metaData)
		request.Headers["x-oss-meta-*"] = nil
		for _, key := range metaKeySet {
			newKey := tea.String("x-oss-meta-" + tea.StringValue(key))
			request.Headers[tea.StringValue(newKey)] = metaData[tea.StringValue(key)]
		}
	}

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

	if !tea.BoolValue(util.Empty(securityToken)) {
		request.Headers["x-oss-security-token"] = securityToken
	}

	if !tea.BoolValue(util.IsUnset(request.Body)) {
		if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("xml"))) {
			reqBodyMap, _err := util.AssertAsMap(request.Body)
			if _err != nil {
				return _err
			}

			request.Stream = tea.ToReader(xml.ToXML(reqBodyMap))
			request.Headers["content-type"] = tea.String("application/xml")
		} else if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("json"))) {
			reqBodyStr := util.ToJSONString(request.Body)
			request.Stream = tea.ToReader(reqBodyStr)
			request.Headers["content-type"] = tea.String("application/json; charset=utf-8")
		} else if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("formData"))) {
			reqBodyForm, _err := util.AssertAsMap(request.Body)
			if _err != nil {
				return _err
			}

			request.Stream = tea.ToReader(openapiutil.ToForm(reqBodyForm))
			request.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
		} else if tea.BoolValue(string_.Equals(request.ReqBodyType, tea.String("binary"))) {
			attributeMap.Key = map[string]*string{
				"crc": tea.String(""),
				"md5": tea.String(""),
			}
			request.Stream = ossutil.Inject(request.Stream, attributeMap.Key)
			request.Headers["content-type"] = tea.String("application/octet-stream")
		}

	}

	host, _err := client.GetHost(config.EndpointType, bucketName, config.Endpoint)
	if _err != nil {
		return _err
	}

	request.Headers = tea.Merge(map[string]*string{
		"host":       host,
		"date":       util.GetDateUTCString(),
		"user-agent": request.UserAgent,
	}, request.Headers)
	request.Headers["authorization"], _err = client.GetAuthorization(request.SignatureVersion, bucketName, request.Pathname, request.Method, request.Query, request.Headers, accessKeyId, accessKeySecret)
	if _err != nil {
		return _err
	}

	return _err
}

func (client *Client) ModifyResponse(context *spi.InterceptorContext, attributeMap *spi.AttributeMap) (_err error) {
	request := context.Request
	response := context.Response
	var bodyStr *string
	if tea.BoolValue(util.Is4xx(response.StatusCode)) || tea.BoolValue(util.Is5xx(response.StatusCode)) {
		bodyStr, _err = util.ReadAsString(response.Body)
		if _err != nil {
			return _err
		}

		if !tea.BoolValue(util.Empty(bodyStr)) {
			respMap := xml.ParseXml(bodyStr, nil)
			err, _err := util.AssertAsMap(respMap["Error"])
			if _err != nil {
				return _err
			}

			_err = tea.NewSDKError(map[string]interface{}{
				"code":    err["Code"],
				"message": err["Message"],
				"data": map[string]interface{}{
					"statusCode": tea.IntValue(response.StatusCode),
					"requestId":  err["RequestId"],
					"hostId":     err["HostId"],
				},
			})
			return _err
		} else {
			headers := response.Headers
			requestId := headers["x-oss-request-id"]
			_err = tea.NewSDKError(map[string]interface{}{
				"code":    tea.IntValue(response.StatusCode),
				"message": nil,
				"data": map[string]interface{}{
					"statusCode": tea.IntValue(response.StatusCode),
					"requestId":  tea.StringValue(requestId),
				},
			})
			return _err
		}

	}

	ctx := attributeMap.Key
	if !tea.BoolValue(util.IsUnset(ctx)) {
		if !tea.BoolValue(util.IsUnset(ctx["crc"])) && !tea.BoolValue(util.IsUnset(response.Headers["x-oss-hash-crc64ecma"])) && !tea.BoolValue(string_.Equals(ctx["crc"], response.Headers["x-oss-hash-crc64ecma"])) {
			_err = tea.NewSDKError(map[string]interface{}{
				"code": "CrcNotMatched",
				"data": map[string]*string{
					"clientCrc": ctx["crc"],
					"serverCrc": response.Headers["x-oss-hash-crc64ecma"],
				},
			})
			return _err
		}

		if !tea.BoolValue(util.IsUnset(ctx["md5"])) && !tea.BoolValue(util.IsUnset(response.Headers["content-md5"])) && !tea.BoolValue(string_.Equals(ctx["md5"], response.Headers["content-md5"])) {
			_err = tea.NewSDKError(map[string]interface{}{
				"code": "MD5NotMatched",
				"data": map[string]*string{
					"clientMD5": ctx["md5"],
					"serverMD5": response.Headers["content-md5"],
				},
			})
			return _err
		}

	}

	if !tea.BoolValue(util.IsUnset(response.Body)) {
		if tea.BoolValue(util.EqualNumber(response.StatusCode, tea.Int(204))) {
			_, _err = util.ReadAsString(response.Body)
			if _err != nil {
				return _err
			}
		} else if tea.BoolValue(string_.Equals(request.BodyType, tea.String("xml"))) {
			bodyStr, _err = util.ReadAsString(response.Body)
			if _err != nil {
				return _err
			}

			result := xml.ParseXml(bodyStr, nil)
			list := map_.KeySet(result)
			if tea.BoolValue(util.EqualNumber(array.Size(list), tea.Int(1))) {
				tmp := list[0]
				tryErr := func() (_e error) {
					defer func() {
						if r := tea.Recover(recover()); r != nil {
							_e = r
						}
					}()
					response.DeserializedBody, _err = util.AssertAsMap(result[tea.StringValue(tmp)])
					if _err != nil {
						return _err
					}

					return nil
				}()

				if tryErr != nil {
					var error = &tea.SDKError{}
					if _t, ok := tryErr.(*tea.SDKError); ok {
						error = _t
					} else {
						error.Message = tea.String(tryErr.Error())
					}
					response.DeserializedBody = result
				}
			} else {
				response.DeserializedBody = result
			}

		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("binary"))) {
			response.DeserializedBody = response.Body
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("byte"))) {
			byt, _err := util.ReadAsBytes(response.Body)
			if _err != nil {
				return _err
			}

			response.DeserializedBody = byt
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("string"))) {
			response.DeserializedBody, _err = util.ReadAsString(response.Body)
			if _err != nil {
				return _err
			}

		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("json"))) {
			obj, _err := util.ReadAsJSON(response.Body)
			if _err != nil {
				return _err
			}

			res, _err := util.AssertAsMap(obj)
			if _err != nil {
				return _err
			}

			response.DeserializedBody = res
		} else if tea.BoolValue(util.EqualString(request.BodyType, tea.String("array"))) {
			response.DeserializedBody, _err = util.ReadAsJSON(response.Body)
			if _err != nil {
				return _err
			}

		} else {
			response.DeserializedBody, _err = util.ReadAsString(response.Body)
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
		if tea.BoolValue(string_.Contains(network, tea.String("internal"))) {
			_result = tea.String("oss-" + tea.StringValue(regionId) + "-internal.aliyuncs.com")
			return _result, _err
		} else if tea.BoolValue(string_.Contains(network, tea.String("ipv6"))) {
			_result = tea.String(tea.StringValue(regionId) + "oss.aliyuncs.com")
			return _result, _err
		} else if tea.BoolValue(string_.Contains(network, tea.String("accelerate"))) {
			_result = tea.String("oss-" + tea.StringValue(network) + ".aliyuncs.com")
			return _result, _err
		}

	}

	_result = tea.String("oss-" + tea.StringValue(regionId) + ".aliyuncs.com")
	return _result, _err
}

func (client *Client) GetHost(endpointType *string, bucketName *string, endpoint *string) (_result *string, _err error) {
	if tea.BoolValue(util.Empty(bucketName)) {
		_result = endpoint
		return _result, _err
	}

	host := tea.String(tea.StringValue(bucketName) + "." + tea.StringValue(endpoint))
	if !tea.BoolValue(util.Empty(endpointType)) {
		if tea.BoolValue(string_.Equals(endpointType, tea.String("ip"))) {
			host = tea.String(tea.StringValue(endpoint) + "/" + tea.StringValue(bucketName))
		} else if tea.BoolValue(string_.Equals(endpointType, tea.String("cname"))) {
			host = endpoint
		}

	}

	_result = host
	return _result, _err
}

func (client *Client) GetAuthorization(signatureVersion *string, bucketName *string, pathname *string, method *string, query map[string]*string, headers map[string]*string, ak *string, secret *string) (_result *string, _err error) {
	sign := tea.String("")
	if tea.BoolValue(util.IsUnset(signatureVersion)) || tea.BoolValue(string_.Equals(signatureVersion, tea.String("v1"))) {
		sign, _err = client.GetSignatureV1(bucketName, pathname, method, query, headers, secret)
		if _err != nil {
			return _result, _err
		}

		_result = tea.String("OSS " + tea.StringValue(ak) + ":" + tea.StringValue(sign))
		return _result, _err
	} else {
		sign, _err = client.GetSignatureV2(bucketName, pathname, method, query, headers, secret)
		if _err != nil {
			return _result, _err
		}

		_result = tea.String("OSS2 AccessKeyId:" + tea.StringValue(ak) + ",Signature:" + tea.StringValue(sign))
		return _result, _err
	}

}

func (client *Client) GetSignatureV1(bucketName *string, pathname *string, method *string, query map[string]*string, headers map[string]*string, secret *string) (_result *string, _err error) {
	resource := tea.String("")
	stringToSign := tea.String("")
	if !tea.BoolValue(util.Empty(bucketName)) {
		resource = tea.String("/" + tea.StringValue(bucketName))
	}

	resource = tea.String(tea.StringValue(resource) + tea.StringValue(pathname))
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
	subResourcesMap := make(map[string]*string)
	canonicalizedResource := pathname
	if !tea.BoolValue(util.Empty(pathname)) {
		paths := string_.Split(pathname, tea.String("?"), tea.Int(2))
		canonicalizedResource = paths[0]
		if tea.BoolValue(util.EqualNumber(array.Size(paths), tea.Int(2))) {
			subResources := string_.Split(paths[1], tea.String("&"), nil)
			for _, sub := range subResources {
				hasExcepts := tea.Bool(false)
				for _, excepts := range client.Except_signed_params {
					if tea.BoolValue(string_.Contains(sub, excepts)) {
						hasExcepts = tea.Bool(true)
					}

				}
				if !tea.BoolValue(hasExcepts) {
					item := string_.Split(sub, tea.String("="), nil)
					key := item[0]
					var value *string
					if tea.BoolValue(util.EqualNumber(array.Size(item), tea.Int(2))) {
						value = item[1]
					}

					subResourcesMap[tea.StringValue(key)] = value
				}

			}
		}

	}

	subResourcesArray := map_.KeySet(subResourcesMap)
	newQueryList := subResourcesArray
	if !tea.BoolValue(util.IsUnset(query)) {
		queryList := map_.KeySet(query)
		newQueryList = array.Concat(queryList, subResourcesArray)
	}

	sortedParams := array.AscSort(newQueryList)
	separator := tea.String("?")
	for _, paramName := range sortedParams {
		if tea.BoolValue(array.Contains(client.Default_signed_params, paramName)) {
			canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + tea.StringValue(separator) + tea.StringValue(paramName))
			if !tea.BoolValue(util.IsUnset(query)) && !tea.BoolValue(util.IsUnset(query[tea.StringValue(paramName)])) {
				canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + "=" + tea.StringValue(query[tea.StringValue(paramName)]))
			} else if !tea.BoolValue(util.IsUnset(subResourcesMap[tea.StringValue(paramName)])) {
				canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + "=" + tea.StringValue(subResourcesMap[tea.StringValue(paramName)]))
			}

		} else if tea.BoolValue(array.Contains(subResourcesArray, paramName)) {
			canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + tea.StringValue(separator) + tea.StringValue(paramName))
			if !tea.BoolValue(util.IsUnset(subResourcesMap[tea.StringValue(paramName)])) {
				canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + "=" + tea.StringValue(subResourcesMap[tea.StringValue(paramName)]))
			}

		}

		separator = tea.String("&")
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
		if tea.BoolValue(string_.Contains(string_.ToLower(header), tea.String("x-oss-"))) && !tea.BoolValue(util.IsUnset(headers[tea.StringValue(header)])) {
			canonicalizedHeaders = tea.String(tea.StringValue(canonicalizedHeaders) + tea.StringValue(header) + ":" + tea.StringValue(headers[tea.StringValue(header)]) + "\n")
		}

	}
	_result = canonicalizedHeaders
	return _result, _err
}

func (client *Client) GetSignatureV2(bucketName *string, pathname *string, method *string, query map[string]*string, headers map[string]*string, secret *string) (_result *string, _err error) {
	_result = tea.String("")
	return _result, _err
}
