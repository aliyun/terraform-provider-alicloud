package oss

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/tidwall/sjson"
)

type OssServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewOssServiceV2(client *connectivity.AliyunClient) *OssServiceV2 {
	return &OssServiceV2{client: client}
}

func ParseOssBucketAclResponse(response map[string]interface{}) (map[string]interface{}, error) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	var normalizedResponse map[string]interface{}
	if err := json.Unmarshal(responseJSON, &normalizedResponse); err != nil {
		return nil, err
	}

	v, err := jsonpath.Get("$.AccessControlPolicy.AccessControlList", normalizedResponse)
	if err != nil {
		return nil, err
	}
	acl, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected AccessControlList type %T", v)
	}
	return acl, nil
}

// DescribeOssBucketAcl <<< Encapsulated get interface for Oss BucketAcl.

func (s *OssServiceV2) DescribeOssBucketAcl(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?acl")

	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketAcl", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketAcl", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketAcl", id), errs.NotFoundMsg, response)
	}

	v, err := ParseOssBucketAclResponse(response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccessControlPolicy.AccessControlList", response)
	}

	return v, nil
}

func (s *OssServiceV2) OssBucketAclStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketAcl(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketAcl >>> Encapsulated.

// DescribeOssBucketReferer <<< Encapsulated get interface for Oss BucketReferer.

func (s *OssServiceV2) DescribeOssBucketReferer(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?referer")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketReferer", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketReferer", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketReferer", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.RefererConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RefererConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketRefererStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketReferer(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketReferer >>> Encapsulated.

// DescribeOssBucketHttpsConfig <<< Encapsulated get interface for Oss BucketHttpsConfig.

func (s *OssServiceV2) DescribeOssBucketHttpsConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?httpsConfig")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketHttpsConfig", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchHttpsConfig", "NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketHttpsConfig", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketHttpsConfig", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.HttpsConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.HttpsConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketHttpsConfigStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketHttpsConfig(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketHttpsConfig >>> Encapsulated.

// DescribeOssBucketCors <<< Encapsulated get interface for Oss BucketCors.

func (s *OssServiceV2) DescribeOssBucketCors(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?cors")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketCors", action), query, nil, nil, hostMap, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchCORSConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketCors", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketCors", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.CORSConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CORSConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketCorsStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketCors(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketCors >>> Encapsulated.

// DescribeOssBucketPolicy
func (s *OssServiceV2) DescribeOssBucketPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?policy")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlJsonParam("GET", "2019-05-17", "GetBucketPolicy", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketPolicy", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return response, nil
}

func (s *OssServiceV2) OssBucketPolicyStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketPolicy(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketPolicy >>> Encapsulated.

// DescribeOssBucketVersioning <<< Encapsulated get interface for Oss BucketVersioning.

func (s *OssServiceV2) DescribeOssBucketVersioning(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?versioning")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketVersioning", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketVersioning", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.VersioningConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketVersioning", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Status"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketVersioning", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketVersioningStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketVersioning(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketVersioning >>> Encapsulated.

// DescribeOssBucketRequestPayment <<< Encapsulated get interface for Oss BucketRequestPayment.
func (s *OssServiceV2) DescribeOssBucketRequestPayment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?requestPayment")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketRequestPayment", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketRequestPayment", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.RequestPaymentConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RequestPaymentConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketRequestPaymentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketRequestPayment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketRequestPayment >>> Encapsulated.

// DescribeOssBucketTransferAcceleration <<< Encapsulated get interface for Oss BucketTransferAcceleration.
func (s *OssServiceV2) DescribeOssBucketTransferAcceleration(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?transferAcceleration")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketTransferAcceleration", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"404"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketTransferAcceleration", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketTransferAcceleration", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.TransferAccelerationConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransferAccelerationConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketTransferAccelerationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketTransferAcceleration(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketTransferAcceleration >>> Encapsulated.

// DescribeOssBucketAccessMonitor <<< Encapsulated get interface for Oss BucketAccessMonitor.
func (s *OssServiceV2) DescribeOssBucketAccessMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?accessmonitor")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketAccessMonitor", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketAccessMonitor", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketAccessMonitor", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.AccessMonitorConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccessMonitorConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketAccessMonitorStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketAccessMonitor(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketAccessMonitor >>> Encapsulated.

// DescribeOssBucketServerSideEncryption <<< Encapsulated get interface for Oss BucketServerSideEncryption.
func (s *OssServiceV2) DescribeOssBucketServerSideEncryption(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?encryption")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketEncryption", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchServerSideEncryptionRule", "NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketServerSideEncryption", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketServerSideEncryption", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketServerSideEncryptionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketServerSideEncryption(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketServerSideEncryption >>> Encapsulated.

// DescribeOssBucketUserDefinedLogFields <<< Encapsulated get interface for Oss BucketUserDefinedLogFields.
func (s *OssServiceV2) DescribeOssBucketUserDefinedLogFields(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?userDefinedLogFieldsConfig")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketUserDefinedLogFieldsConfig", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchUserDefinedLogFieldsConfig"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketUserDefinedLogFields", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketUserDefinedLogFields", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *OssServiceV2) OssBucketUserDefinedLogFieldsStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketUserDefinedLogFields(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketUserDefinedLogFields >>> Encapsulated.

// DescribeOssBucketMetaQuery <<< Encapsulated get interface for Oss BucketMetaQuery.
func (s *OssServiceV2) DescribeOssBucketMetaQuery(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?metaQuery")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketMetaQueryStatus", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "MetaQueryNotExist", "AccessDenied"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketMetaQuery", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketMetaQuery", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.MetaQueryStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.MetaQueryStatus", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketMetaQueryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketMetaQuery(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketMetaQuery >>> Encapsulated.

// DescribeOssBucketDataRedundancyTransition <<< Encapsulated get interface for Oss BucketDataRedundancyTransition.
func (s *OssServiceV2) DescribeOssBucketDataRedundancyTransition(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := fmt.Sprintf("/?redundancyTransition")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["x-oss-redundancy-transition-taskid"] = helper.StringPointer(parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketDataRedundancyTransition", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "BucketDataRedundancyTransitionTaskNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketDataRedundancyTransition", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketDataRedundancyTransition", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *OssServiceV2) OssBucketDataRedundancyTransitionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketDataRedundancyTransition(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketDataRedundancyTransition >>> Encapsulated.

// DescribeOssAccountPublicAccessBlock <<< Encapsulated get interface for Oss AccountPublicAccessBlock.
func (s *OssServiceV2) DescribeOssAccountPublicAccessBlock(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 0 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 0, len(parts)))
	}
	action := fmt.Sprintf("/?publicAccessBlock")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetPublicAccessBlock", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AccountPublicAccessBlock", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("AccountPublicAccessBlock", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.PublicAccessBlockConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PublicAccessBlockConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssAccountPublicAccessBlockStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssAccountPublicAccessBlock(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssAccountPublicAccessBlock >>> Encapsulated.

// DescribeOssBucketPublicAccessBlock <<< Encapsulated get interface for Oss BucketPublicAccessBlock.
func (s *OssServiceV2) DescribeOssBucketPublicAccessBlock(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?publicAccessBlock")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketPublicAccessBlock", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketPublicAccessBlock", id), errs.NotFoundMsg, response)
	}

	currentStatus, err := jsonpath.Get("$.PublicAccessBlockConfiguration.BlockPublicAccess", response)
	if currentStatus == "" {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketPublicAccessBlock", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *OssServiceV2) OssBucketPublicAccessBlockStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketPublicAccessBlock(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketPublicAccessBlock >>> Encapsulated.

// DescribeOssBucketCname <<< Encapsulated get interface for Oss BucketCname.
func (s *OssServiceV2) DescribeOssBucketCname(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := fmt.Sprintf("/?cname")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "ListCname", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketCname", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketCname", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.ListCnameResult", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ListCnameResult", response)
	}

	item := v.(map[string]interface{})
	domains, err := jsonpath.Get("$.Cname[*]", item)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Cname[*]", item)
	}
	for _, vv := range domains.([]interface{}) {
		if vv.(map[string]interface{})["Domain"] == parts[1] {
			return vv.(map[string]interface{}), nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("BucketCname", id), errs.NotFoundMsg, response)
}

func (s *OssServiceV2) OssBucketCnameStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketCname(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		if err != nil {
			return object, "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, field, object)
		}
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketCname >>> Encapsulated.

// DescribeOssBucketCnameToken <<< Encapsulated get interface for Oss BucketCnameToken.
func (s *OssServiceV2) DescribeOssBucketCnameToken(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := fmt.Sprintf("/?cname&comp=token")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])

	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "BucketCnameConfiguration.Cname.Domain", parts[1])
	json.Unmarshal([]byte(jsonString), &request)
	body := request

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("POST", "2019-05-17", "CreateCnameToken", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "CnameTokenNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketCnameToken", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketCnameToken", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.CnameToken", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CnameToken", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketCnameTokenStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketCnameToken(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketCnameToken >>> Encapsulated.

// DescribeOssBucketWebsite <<< Encapsulated get interface for Oss BucketWebsite.
func (s *OssServiceV2) DescribeOssBucketWebsite(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?website")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketWebsite", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchWebsiteConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketWebsite", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketWebsite", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.WebsiteConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.WebsiteConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketWebsiteStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketWebsite(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketWebsite >>> Encapsulated.

// DescribeOssBucketLifecycle <<< Encapsulated get interface for Oss BucketLifecycle.
func (s *OssServiceV2) DescribeOssBucketLifecycle(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/?lifecycle")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketLifecycle", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchLifecycle"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketLifecycle", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketLifecycle", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.LifecycleConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.LifecycleConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketLifecycleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketLifecycle(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketLifecycle >>> Encapsulated.

// DescribeOssBucketWorm <<< Encapsulated get interface for Oss BucketWorm.
func (s *OssServiceV2) DescribeOssBucketWorm(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := fmt.Sprintf("/?worm")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketWorm", action), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchWORMConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketWorm", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketWorm", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.WormConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.WormConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketWormStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketWorm(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketWorm >>> Encapsulated.

// DescribeOssBucketStyle <<< Encapsulated get interface for Oss BucketStyle.

func (s *OssServiceV2) DescribeOssBucketStyle(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["styleName"] = helper.StringPointer(parts[1])

	action := fmt.Sprintf("/?style")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetStyle", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchStyle", "NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketStyle", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Style", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Style", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketStyleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketStyle(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketStyle >>> Encapsulated.

// DescribeOssBucketLogging <<< Encapsulated get interface for Oss BucketLogging.

func (s *OssServiceV2) DescribeOssBucketLogging(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?logging")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketLogging", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.BucketLoggingStatus.LoggingEnabled", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketLogging", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["TargetBucket"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("BucketLogging", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketLoggingStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketLoggingStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketLogging)
}

func (s *OssServiceV2) OssBucketLoggingStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketLogging >>> Encapsulated.
// DescribeOssBucketArchiveDirectRead <<< Encapsulated get interface for Oss BucketArchiveDirectRead.

func (s *OssServiceV2) DescribeOssBucketArchiveDirectRead(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?bucketArchiveDirectRead")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketArchiveDirectRead", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketArchiveDirectRead", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ArchiveDirectReadConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ArchiveDirectReadConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketArchiveDirectReadStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketArchiveDirectReadStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketArchiveDirectRead)
}

func (s *OssServiceV2) OssBucketArchiveDirectReadStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketArchiveDirectRead >>> Encapsulated.
// DescribeOssBucketOverwriteConfig <<< Encapsulated get interface for Oss BucketOverwriteConfig.

func (s *OssServiceV2) DescribeOssBucketOverwriteConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?overwriteConfig")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketOverwriteConfig", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchBucketOverwriteConfig"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketOverwriteConfig", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.OverwriteConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.OverwriteConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketOverwriteConfigStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketOverwriteConfigStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketOverwriteConfig)
}

func (s *OssServiceV2) OssBucketOverwriteConfigStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketOverwriteConfig >>> Encapsulated.
// DescribeOssBucketResponseHeader <<< Encapsulated get interface for Oss BucketResponseHeader.

func (s *OssServiceV2) DescribeOssBucketResponseHeader(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?responseHeader")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketResponseHeader", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchResponseHeaderConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketResponseHeader", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ResponseHeaderConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ResponseHeaderConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketResponseHeaderStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketResponseHeaderStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketResponseHeader)
}

func (s *OssServiceV2) OssBucketResponseHeaderStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketResponseHeader >>> Encapsulated.
// DescribeOssBucketObjectWormConfiguration <<< Encapsulated get interface for Oss BucketObjectWormConfiguration.

func (s *OssServiceV2) DescribeOssBucketObjectWormConfiguration(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(id)

	action := fmt.Sprintf("/?objectWorm")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketObjectWormConfiguration", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ObjectWormConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ObjectWormConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketObjectWormConfigurationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketObjectWormConfigurationStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketObjectWormConfiguration)
}

func (s *OssServiceV2) OssBucketObjectWormConfigurationStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketObjectWormConfiguration >>> Encapsulated.
// DescribeOssBucketInventory <<< Encapsulated get interface for Oss BucketInventory.

func (s *OssServiceV2) DescribeOssBucketInventory(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["inventoryId"] = helper.StringPointer(parts[1])

	action := fmt.Sprintf("/?inventory&inventoryId=%s", parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketInventory", action), query, nil, nil, hostMap, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchInventoryConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("BucketInventory", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.InventoryConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InventoryConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketInventoryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OssBucketInventoryStateRefreshFuncWithApi(id, field, failStates, s.DescribeOssBucketInventory)
}

func (s *OssServiceV2) OssBucketInventoryStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketInventory >>> Encapsulated.
