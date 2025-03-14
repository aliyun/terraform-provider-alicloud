// This file is auto-generated, don't edit it. Thanks.
package client

import (
	oss_util "github.com/alibabacloud-go/alibabacloud-gateway-oss-util/client"
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
	client.Default_signed_params = []*string{tea.String("response-content-type"), tea.String("response-content-language"), tea.String("response-cache-control"), tea.String("logging"), tea.String("response-content-encoding"), tea.String("acl"), tea.String("uploadId"), tea.String("uploads"), tea.String("partNumber"), tea.String("group"), tea.String("link"), tea.String("delete"), tea.String("website"), tea.String("location"), tea.String("objectInfo"), tea.String("objectMeta"), tea.String("response-expires"), tea.String("response-content-disposition"), tea.String("cors"), tea.String("lifecycle"), tea.String("restore"), tea.String("qos"), tea.String("referer"), tea.String("stat"), tea.String("bucketInfo"), tea.String("append"), tea.String("position"), tea.String("security-token"), tea.String("live"), tea.String("comp"), tea.String("status"), tea.String("vod"), tea.String("startTime"), tea.String("endTime"), tea.String("x-oss-process"), tea.String("symlink"), tea.String("callback"), tea.String("callback-var"), tea.String("tagging"), tea.String("encryption"), tea.String("versions"), tea.String("versioning"), tea.String("versionId"), tea.String("policy"), tea.String("requestPayment"), tea.String("x-oss-traffic-limit"), tea.String("qosInfo"), tea.String("asyncFetch"), tea.String("x-oss-request-payer"), tea.String("sequential"), tea.String("inventory"), tea.String("inventoryId"), tea.String("continuation-token"), tea.String("callback"), tea.String("callback-var"), tea.String("worm"), tea.String("wormId"), tea.String("wormExtend"), tea.String("replication"), tea.String("replicationLocation"), tea.String("replicationProgress"), tea.String("transferAcceleration"), tea.String("cname"), tea.String("metaQuery"), tea.String("x-oss-ac-source-ip"), tea.String("x-oss-ac-subnet-mask"), tea.String("x-oss-ac-vpc-id"), tea.String("x-oss-ac-forward-allow"), tea.String("resourceGroup"), tea.String("style"), tea.String("styleName"), tea.String("x-oss-async-process"), tea.String("rtc"), tea.String("accessPoint"), tea.String("accessPointPolicy"), tea.String("httpsConfig"), tea.String("regionsV2"), tea.String("publicAccessBlock"), tea.String("policyStatus"), tea.String("redundancyTransition"), tea.String("redundancyType"), tea.String("redundancyProgress"), tea.String("dataAccelerator"), tea.String("verbose"), tea.String("accessPointForObjectProcess"), tea.String("accessPointConfigForObjectProcess"), tea.String("accessPointPolicyForObjectProcess"), tea.String("bucketArchiveDirectRead"), tea.String("responseHeader"), tea.String("userDefinedLogFieldsConfig"), tea.String("reservedcapacity"), tea.String("requesterQosInfo"), tea.String("qosRequester"), tea.String("resourcePool"), tea.String("resourcePoolInfo"), tea.String("resourcePoolBuckets"), tea.String("processConfiguration"), tea.String("img"), tea.String("asyncFetch"), tea.String("virtualBucket"), tea.String("copy"), tea.String("userRegion"), tea.String("partSize"), tea.String("chunkSize"), tea.String("partUploadId"), tea.String("chunkNumber"), tea.String("userRegion"), tea.String("regionList"), tea.String("eventnotification"), tea.String("cacheConfiguration"), tea.String("dfs"), tea.String("dfsadmin"), tea.String("dfssecurity")}
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
	regionId := config.RegionId
	if tea.BoolValue(util.IsUnset(regionId)) || tea.BoolValue(util.Empty(regionId)) {
		regionId, _err = client.GetRegionIdFromEndpoint(config.Endpoint)
		if _err != nil {
			return _err
		}

	}

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

			// for python:
			// xml_str = OSS_UtilClient.to_xml(req_body_map)
			xmlStr := xml.ToXML(reqBodyMap)
			request.Stream = tea.ToReader(xmlStr)
			request.Headers["content-type"] = tea.String("application/xml")
			request.Headers["content-md5"] = encodeutil.Base64EncodeToString(signatureutil.MD5Sign(xmlStr))
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

	host, _err := client.GetHost(config.EndpointType, bucketName, config.Endpoint, context)
	if _err != nil {
		return _err
	}

	request.Headers = tea.Merge(map[string]*string{
		"host":       host,
		"date":       util.GetDateUTCString(),
		"user-agent": request.UserAgent,
	}, request.Headers)
	originPath := request.Pathname
	originQuery := request.Query
	if !tea.BoolValue(util.Empty(originPath)) {
		pathAndQueries := string_.Split(originPath, tea.String("?"), tea.Int(2))
		request.Pathname = pathAndQueries[0]
		if tea.BoolValue(util.EqualNumber(array.Size(pathAndQueries), tea.Int(2))) {
			pathQueries := string_.Split(pathAndQueries[1], tea.String("&"), nil)
			for _, sub := range pathQueries {
				item := string_.Split(sub, tea.String("="), nil)
				queryKey := item[0]
				queryValue := tea.String("")
				if tea.BoolValue(util.EqualNumber(array.Size(item), tea.Int(2))) {
					queryValue = item[1]
				}

				if tea.BoolValue(util.Empty(originQuery[tea.StringValue(queryKey)])) {
					request.Query[tea.StringValue(queryKey)] = queryValue
				}

			}
		}

	}

	signatureVersion := util.DefaultString(request.SignatureVersion, tea.String("v4"))
	request.Headers["authorization"], _err = client.GetAuthorization(signatureVersion, bucketName, request.Pathname, request.Method, request.Query, request.Headers, accessKeyId, accessKeySecret, regionId)
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
					"statusCode":         tea.IntValue(response.StatusCode),
					"requestId":          err["RequestId"],
					"ecCode":             err["EC"],
					"Recommend":          err["RecommendDoc"],
					"hostId":             err["HostId"],
					"AccessDeniedDetail": err["AccessDeniedDetail"],
				},
			})
			return _err
		} else {
			headers := response.Headers
			requestId := headers["x-oss-request-id"]
			ecCode := headers["x-oss-ec-code"]
			_err = tea.NewSDKError(map[string]interface{}{
				"code":    tea.IntValue(response.StatusCode),
				"message": nil,
				"data": map[string]interface{}{
					"statusCode": tea.IntValue(response.StatusCode),
					"requestId":  tea.StringValue(requestId),
					"ecCode":     tea.StringValue(ecCode),
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

			response.DeserializedBody = bodyStr
			if !tea.BoolValue(util.Empty(bodyStr)) {
				result, _err := oss_util.ParseXml(bodyStr, request.Action)
				if _err != nil {
					return _err
				}

				// for no util language
				// var result : any = XML.parseXml(bodyStr, null);
				tryErr := func() (_e error) {
					defer func() {
						if r := tea.Recover(recover()); r != nil {
							_e = r
						}
					}()
					response.DeserializedBody, _err = util.AssertAsMap(result)
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

func (client *Client) GetRegionIdFromEndpoint(endpoint *string) (_result *string, _err error) {
	if !tea.BoolValue(util.Empty(endpoint)) {
		idx := tea.Int(-1)
		if tea.BoolValue(string_.HasPrefix(endpoint, tea.String("oss-"))) && tea.BoolValue(string_.HasSuffix(endpoint, tea.String(".aliyuncs.com"))) {
			idx = string_.Index(endpoint, tea.String(".aliyuncs.com"))
			_body := string_.SubString(endpoint, tea.Int(4), idx)
			_result = _body
			return _result, _err
		}

		if tea.BoolValue(string_.HasSuffix(endpoint, tea.String(".mgw.aliyuncs.com"))) {
			idx = string_.Index(endpoint, tea.String(".mgw.aliyuncs.com"))
			_body := string_.SubString(endpoint, tea.Int(0), idx)
			_result = _body
			return _result, _err
		}

		if tea.BoolValue(string_.HasSuffix(endpoint, tea.String("-internal.oss-data-acc.aliyuncs.com"))) {
			idx = string_.Index(endpoint, tea.String("-internal.oss-data-acc.aliyuncs.com"))
			_body := string_.SubString(endpoint, tea.Int(0), idx)
			_result = _body
			return _result, _err
		}

		if tea.BoolValue(string_.HasSuffix(endpoint, tea.String(".oss-dls.aliyuncs.com"))) {
			idx = string_.Index(endpoint, tea.String(".oss-dls.aliyuncs.com"))
			_body := string_.SubString(endpoint, tea.Int(0), idx)
			_result = _body
			return _result, _err
		}

	}

	_result = tea.String("cn-hangzhou")
	return _result, _err
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

func (client *Client) GetHost(endpointType *string, bucketName *string, endpoint *string, context *spi.InterceptorContext) (_result *string, _err error) {
	if tea.BoolValue(string_.Contains(endpoint, tea.String(".mgw.aliyuncs.com"))) && !tea.BoolValue(util.IsUnset(context.Request.HostMap["userid"])) {
		_result = tea.String(tea.StringValue(context.Request.HostMap["userid"]) + "." + tea.StringValue(endpoint))
		return _result, _err
	}

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

func (client *Client) GetAuthorization(signatureVersion *string, bucketName *string, pathname *string, method *string, query map[string]*string, headers map[string]*string, ak *string, secret *string, regionId *string) (_result *string, _err error) {
	sign := tea.String("")
	if !tea.BoolValue(util.IsUnset(signatureVersion)) {
		if tea.BoolValue(string_.Equals(signatureVersion, tea.String("v1"))) {
			sign, _err = client.GetSignatureV1(bucketName, pathname, method, query, headers, secret)
			if _err != nil {
				return _result, _err
			}

			_result = tea.String("OSS " + tea.StringValue(ak) + ":" + tea.StringValue(sign))
			return _result, _err
		}

		if tea.BoolValue(string_.Equals(signatureVersion, tea.String("v2"))) {
			sign, _err = client.GetSignatureV2(bucketName, pathname, method, query, headers, secret)
			if _err != nil {
				return _result, _err
			}

			_result = tea.String("OSS2 AccessKeyId:" + tea.StringValue(ak) + ",Signature:" + tea.StringValue(sign))
			return _result, _err
		}

	}

	dateTime := openapiutil.GetTimestamp()
	dateTime = string_.Replace(dateTime, tea.String("-"), tea.String(""), nil)
	dateTime = string_.Replace(dateTime, tea.String(":"), tea.String(""), nil)
	headers["x-oss-date"] = dateTime
	headers["x-oss-content-sha256"] = tea.String("UNSIGNED-PAYLOAD")
	onlyDate := string_.SubString(dateTime, tea.Int(0), tea.Int(8))
	cred := tea.String(tea.StringValue(ak) + "/" + tea.StringValue(onlyDate) + "/" + tea.StringValue(regionId) + "/oss/aliyun_v4_request")
	sign, _err = client.GetSignatureV4(bucketName, pathname, method, query, headers, onlyDate, regionId, secret)
	if _err != nil {
		return _result, _err
	}

	_result = tea.String("OSS4-HMAC-SHA256 Credential=" + tea.StringValue(cred) + ", Signature=" + tea.StringValue(sign))
	return _result, _err
}

func (client *Client) GetSignKey(secret *string, onlyDate *string, regionId *string) (_result []byte, _err error) {
	temp := tea.String("aliyun_v4" + tea.StringValue(secret))
	res := signatureutil.HmacSHA256Sign(onlyDate, temp)
	res = signatureutil.HmacSHA256SignByBytes(regionId, res)
	res = signatureutil.HmacSHA256SignByBytes(tea.String("oss"), res)
	res = signatureutil.HmacSHA256SignByBytes(tea.String("aliyun_v4_request"), res)
	_result = res
	return _result, _err
}

func (client *Client) GetSignatureV4(bucketName *string, pathname *string, method *string, query map[string]*string, headers map[string]*string, onlyDate *string, regionId *string, secret *string) (_result *string, _err error) {
	signingkey, _err := client.GetSignKey(secret, onlyDate, regionId)
	if _err != nil {
		return _result, _err
	}

	canonicalizedUri := pathname
	if !tea.BoolValue(util.Empty(pathname)) {
		if !tea.BoolValue(util.Empty(bucketName)) {
			canonicalizedUri = tea.String("/" + tea.StringValue(bucketName) + tea.StringValue(canonicalizedUri))
		}

	} else {
		if !tea.BoolValue(util.Empty(bucketName)) {
			canonicalizedUri = tea.String("/" + tea.StringValue(bucketName) + "/")
		} else {
			canonicalizedUri = tea.String("/")
		}

	}

	// for java:
	// String suffix = (!canonicalizedUri.equals("/") && canonicalizedUri.endsWith("/"))? "/" : "";
	// canonicalizedUri = com.aliyun.openapiutil.Client.getEncodePath(canonicalizedUri) + suffix;
	canonicalizedUri = openapiutil.GetEncodePath(canonicalizedUri)
	queryMap := make(map[string]*string)
	for _, queryKey := range map_.KeySet(query) {
		var queryValue *string
		if !tea.BoolValue(util.Empty(query[tea.StringValue(queryKey)])) {
			queryValue = encodeutil.PercentEncode(query[tea.StringValue(queryKey)])
			queryValue = string_.Replace(queryValue, tea.String("+"), tea.String("%20"), nil)
		}

		queryKey = encodeutil.PercentEncode(queryKey)
		queryKey = string_.Replace(queryKey, tea.String("+"), tea.String("%20"), nil)
		// for go : queryMap[tea.StringValue(queryKey)] = queryValue
		queryMap[tea.StringValue(queryKey)] = queryValue
	}
	canonicalizedQueryString, _err := client.BuildCanonicalizedQueryStringV4(queryMap)
	if _err != nil {
		return _result, _err
	}

	canonicalizedHeaders, _err := client.BuildCanonicalizedHeadersV4(headers)
	if _err != nil {
		return _result, _err
	}

	payload := tea.String("UNSIGNED-PAYLOAD")
	canonicalRequest := tea.String(tea.StringValue(method) + "\n" + tea.StringValue(canonicalizedUri) + "\n" + tea.StringValue(canonicalizedQueryString) + "\n" + tea.StringValue(canonicalizedHeaders) + "\n\n" + tea.StringValue(payload))
	hex := encodeutil.HexEncode(encodeutil.Hash(util.ToBytes(canonicalRequest), tea.String("ACS4-HMAC-SHA256")))
	scope := tea.String(tea.StringValue(onlyDate) + "/" + tea.StringValue(regionId) + "/oss/aliyun_v4_request")
	stringToSign := tea.String("OSS4-HMAC-SHA256\n" + tea.StringValue(headers["x-oss-date"]) + "\n" + tea.StringValue(scope) + "\n" + tea.StringValue(hex))
	signature := signatureutil.HmacSHA256SignByBytes(stringToSign, signingkey)
	_body := encodeutil.HexEncode(signature)
	_result = _body
	return _result, _err
}

func (client *Client) BuildCanonicalizedQueryStringV4(queryMap map[string]*string) (_result *string, _err error) {
	canonicalizedQueryString := tea.String("")
	if !tea.BoolValue(util.IsUnset(queryMap)) {
		queryArray := map_.KeySet(queryMap)
		sortedQueryArray := array.AscSort(queryArray)
		separator := tea.String("")
		for _, key := range sortedQueryArray {
			canonicalizedQueryString = tea.String(tea.StringValue(canonicalizedQueryString) + tea.StringValue(separator) + tea.StringValue(key))
			if !tea.BoolValue(util.Empty(queryMap[tea.StringValue(key)])) {
				canonicalizedQueryString = tea.String(tea.StringValue(canonicalizedQueryString) + "=" + tea.StringValue(queryMap[tea.StringValue(key)]))
			}

			separator = tea.String("&")
		}
	}

	_result = canonicalizedQueryString
	return _result, _err
}

func (client *Client) BuildCanonicalizedHeadersV4(headers map[string]*string) (_result *string, _err error) {
	canonicalizedHeaders := tea.String("")
	headersArray := map_.KeySet(headers)
	sortedHeadersArray := array.AscSort(headersArray)
	for _, key := range sortedHeadersArray {
		lowerKey := string_.ToLower(key)
		if tea.BoolValue(string_.HasPrefix(lowerKey, tea.String("x-oss-"))) || tea.BoolValue(string_.Equals(lowerKey, tea.String("content-type"))) || tea.BoolValue(string_.Equals(lowerKey, tea.String("content-md5"))) {
			canonicalizedHeaders = tea.String(tea.StringValue(canonicalizedHeaders) + tea.StringValue(key) + ":" + tea.StringValue(string_.Trim(headers[tea.StringValue(key)])) + "\n")
		}

	}
	_result = canonicalizedHeaders
	return _result, _err
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
	canonicalizedResource := pathname
	queryKeys := map_.KeySet(query)
	sortedParams := array.AscSort(queryKeys)
	separator := tea.String("?")
	for _, paramName := range sortedParams {
		if tea.BoolValue(array.Contains(client.Default_signed_params, paramName)) || tea.BoolValue(string_.HasPrefix(paramName, tea.String("x-oss-"))) {
			canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + tea.StringValue(separator) + tea.StringValue(paramName))
			if !tea.BoolValue(util.Empty(query[tea.StringValue(paramName)])) {
				canonicalizedResource = tea.String(tea.StringValue(canonicalizedResource) + "=" + tea.StringValue(query[tea.StringValue(paramName)]))
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
