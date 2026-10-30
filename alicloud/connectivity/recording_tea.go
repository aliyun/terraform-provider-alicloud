// Copyright (c) 2009-present Alibaba Cloud. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Adapted from tea-rpc/client and tea-roa/client v1.3.4 DoRequest methods.
// The upstream methods do not expose runtime["httpClient"]. These adapters add
// that one injection point, preserving serialization, signing, retry and SDK
// errors at the HTTP boundary. They are used only for opt-in recording; ordinary
// execution continues through the upstream methods. See recording_tea_LICENSE.
package connectivity

import (
	roautil "github.com/alibabacloud-go/tea-roa-utils/service"
	roa "github.com/alibabacloud-go/tea-roa/client"
	rpcutil "github.com/alibabacloud-go/tea-rpc-utils/service"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"time"
)

func recordingRPCRequest(client *rpc.Client, action *string, protocol *string, method *string, version *string, authType *string, query map[string]interface{}, body map[string]interface{}, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"httpClient":     recording.HTTPClient{Timeout: time.Duration(tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout))+tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout))) * time.Millisecond},
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		var serializationErr error
		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = tea.String("/")
			request_.Query, _err = rpcutil.FlattenParams(tea.ToMap(map[string]interface{}{
				"Action":         tea.StringValue(action),
				"Format":         "json",
				"Timestamp":      tea.StringValue(rpcutil.GetTimestamp()),
				"Version":        tea.StringValue(version),
				"SignatureNonce": tea.StringValue(util.GetNonce()),
			}, query))
			if _err != nil {
				serializationErr = _err
				return _result, serializationErr
			}
			if !tea.BoolValue(util.IsUnset(client.SourceIp)) {
				request_.Query["SourceIp"] = client.SourceIp
			}

			if !tea.BoolValue(util.IsUnset(client.SecureTransport)) {
				request_.Query["SecureTransport"] = client.SecureTransport
			}
			headers, _err := client.GetHeaders()
			if tea.BoolValue(util.IsUnset(headers)) {
				// endpoint is setted in product client
				request_.Headers = map[string]*string{
					"host":          client.Endpoint,
					"x-acs-version": version,
					"x-acs-action":  action,
					"user-agent":    client.GetUserAgent(),
				}
			} else {
				request_.Headers = tea.Merge(map[string]*string{
					"host":          client.Endpoint,
					"x-acs-version": version,
					"x-acs-action":  action,
					"user-agent":    client.GetUserAgent(),
				}, headers)
			}
			if !tea.BoolValue(util.IsUnset(body)) {
				bodyParams, err := rpcutil.FlattenParams(body)
				if err != nil {
					serializationErr = err
					return _result, serializationErr
				}
				tmp := util.AnyifyMapValue(bodyParams)
				request_.Body = tea.ToReader(util.ToFormString(tmp))
				request_.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				accessKeyId, _err := client.GetAccessKeyId()
				if _err != nil {
					return _result, _err
				}

				accessKeySecret, _err := client.GetAccessKeySecret()
				if _err != nil {
					return _result, _err
				}

				securityToken, _err := client.GetSecurityToken()
				if _err != nil {
					return _result, _err
				}

				if !tea.BoolValue(util.Empty(securityToken)) {
					request_.Query["SecurityToken"] = securityToken
				}

				request_.Query["SignatureMethod"] = tea.String("HMAC-SHA1")
				request_.Query["SignatureVersion"] = tea.String("1.0")
				request_.Query["AccessKeyId"] = accessKeyId
				bodyParams, err := rpcutil.FlattenParams(body)
				if err != nil {
					serializationErr = err
					return _result, serializationErr
				}
				signedParam := tea.Merge(request_.Query, bodyParams)
				request_.Query["Signature"] = rpcutil.GetSignatureV1(signedParam, request_.Method, accessKeySecret)
			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			obj, _err := util.ReadAsJSON(response_.Body)
			if _err != nil {
				return _result, _err
			}
			res := util.AssertAsMap(obj)
			res["_headers"] = response_.Headers
			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				_err = tea.NewSDKError(map[string]interface{}{
					"code":       tea.ToString(rpc.DefaultAny(res["Code"], res["code"])),
					"statusCode": tea.IntValue(response_.StatusCode),
					"message":    "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(rpc.DefaultAny(res["Message"], res["message"])) + " request id: " + tea.ToString(rpc.DefaultAny(res["RequestId"], res["requestId"])),
					"data":       res,
					"_headers":   response_.Headers,
				})
				return _result, _err
			}

			_result = res
			_result["_headers"] = response_.Headers
			return _result, _err
		}()
		if serializationErr != nil {
			return _resp, serializationErr
		}
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

func recordingROARequest(client *roa.Client, action *string, version *string, protocol *string, method *string, authType *string, pathname *string, query map[string]*string, headers map[string]*string, body interface{}, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"httpClient":     recording.HTTPClient{Timeout: time.Duration(tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout))+tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout))) * time.Millisecond},
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = pathname
			defaultHeaders := map[string]*string{
				"date":                    util.GetDateUTCString(),
				"host":                    client.EndpointHost,
				"accept":                  tea.String("application/json"),
				"x-acs-signature-nonce":   util.GetNonce(),
				"x-acs-signature-method":  tea.String("HMAC-SHA1"),
				"x-acs-signature-version": tea.String("1.0"),
				"x-acs-version":           version,
				"user-agent":              util.GetUserAgent(client.UserAgent),
				// x-sdk-client': helper.DEFAULT_CLIENT
			}
			if tea.StringValue(action) != "" {
				defaultHeaders["x-acs-action"] = action
			}
			request_.Headers = tea.Merge(defaultHeaders, headers)
			if !tea.BoolValue(util.IsUnset(client.SourceIp)) {
				request_.Headers["x-acs-source-ip"] = client.SourceIp
			}

			if !tea.BoolValue(util.IsUnset(client.SecureTransport)) {
				request_.Headers["x-acs-secure-transport"] = client.SecureTransport
			}

			if !tea.BoolValue(util.IsUnset(body)) {
				request_.Body = tea.ToReader(util.ToJSONString(body))
				request_.Headers["content-type"] = tea.String("application/json; charset=utf-8")
			}

			if !tea.BoolValue(util.IsUnset(query)) {
				request_.Query = query
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				accessKeyId, _err := client.Credential.GetAccessKeyId()
				if _err != nil {
					return _result, _err
				}

				accessKeySecret, _err := client.Credential.GetAccessKeySecret()
				if _err != nil {
					return _result, _err
				}

				securityToken, _err := client.Credential.GetSecurityToken()
				if _err != nil {
					return _result, _err
				}

				if !tea.BoolValue(util.Empty(securityToken)) {
					request_.Headers["x-acs-accesskey-id"] = accessKeyId
					request_.Headers["x-acs-security-token"] = securityToken
				}

				stringToSign := roautil.GetStringToSign(request_)
				request_.Headers["authorization"] = tea.String("acs " + tea.StringValue(accessKeyId) + ":" + tea.StringValue(roautil.GetSignature(stringToSign, accessKeySecret)))
			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			if tea.BoolValue(util.EqualNumber(response_.StatusCode, tea.Int(204))) {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]map[string]*string{
					"headers": response_.Headers,
				}, &_result)
				return _result, _err
			}

			result, _err := util.ReadAsJSON(response_.Body)
			if _err != nil {
				return _result, _err
			}

			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				err := util.AssertAsMap(result)
				err["_headers"] = response_.Headers
				_err = tea.NewSDKError(map[string]interface{}{
					"code":       tea.ToString(roa.DefaultAny(err["Code"], err["code"])),
					"statusCode": tea.IntValue(response_.StatusCode),
					"message":    "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(roa.DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(roa.DefaultAny(err["RequestId"], err["requestId"])),
					"data":       err,
				})
				return _result, _err
			}

			_result = make(map[string]interface{})
			_err = tea.Convert(map[string]interface{}{
				"headers": response_.Headers,
				"body":    result,
			}, &_result)
			_result["_headers"] = response_.Headers
			return _result, _err
		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}
