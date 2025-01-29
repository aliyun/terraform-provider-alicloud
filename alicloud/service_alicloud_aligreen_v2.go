package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AligreenServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAligreenAuditCallback <<< Encapsulated get interface for Aligreen AuditCallback.

func (s *AligreenServiceV2) DescribeAligreenAuditCallback(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeAuditCallbackList"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CallbackList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CallbackList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AuditCallback", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Name"]) != id {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("AuditCallback", id)), NotFoundMsg, response)
}

func (s *AligreenServiceV2) AligreenAuditCallbackStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenAuditCallback(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenAuditCallback >>> Encapsulated.

// DescribeAligreenCallback <<< Encapsulated get interface for Aligreen Callback.

func (s *AligreenServiceV2) DescribeAligreenCallback(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeCallbackList"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CallbackList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CallbackList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Callback", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Id"]) != id {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("Callback", id)), NotFoundMsg, response)
}

func (s *AligreenServiceV2) AligreenCallbackStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenCallback(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenCallback >>> Encapsulated.
// DescribeAligreenBizType <<< Encapsulated get interface for Aligreen BizType.

func (s *AligreenServiceV2) DescribeAligreenBizType(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeUserBizTypes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.BizTypeList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BizTypeList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("BizType", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["BizType"]) != id {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("BizType", id)), NotFoundMsg, response)
}

func (s *AligreenServiceV2) AligreenBizTypeStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenBizType(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenBizType >>> Encapsulated.

// DescribeAligreenImageLib <<< Encapsulated get interface for Aligreen ImageLib.

func (s *AligreenServiceV2) DescribeAligreenImageLib(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeImageLib"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["ServiceModule"] = "open_api"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ImageLibList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ImageLibList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ImageLib", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Id"]) != id {
			continue
		}
		if fmt.Sprint(item["ServiceModule"]) != "open_api" {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("ImageLib", id)), NotFoundMsg, response)
}

func (s *AligreenServiceV2) AligreenImageLibStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenImageLib(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenImageLib >>> Encapsulated.
// DescribeAligreenKeywordLib <<< Encapsulated get interface for Aligreen KeywordLib.

func (s *AligreenServiceV2) DescribeAligreenKeywordLib(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeKeywordLib"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["ServiceModule"] = "open_api"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.KeywordLibList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.KeywordLibList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("KeywordLib", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Id"]) != id {
			continue
		}
		if fmt.Sprint(item["ServiceModule"]) != "open_api" {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("KeywordLib", id)), NotFoundMsg, response)
}

func (s *AligreenServiceV2) AligreenKeywordLibStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenKeywordLib(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenKeywordLib >>> Encapsulated.

// DescribeAligreenOssStockTask <<< Encapsulated get interface for Aligreen OssStockTask.

func (s *AligreenServiceV2) DescribeAligreenOssStockTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeOssStockTaskSetting"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AligreenServiceV2) AligreenOssStockTaskStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAligreenOssStockTask(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.AudioAntispamFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.AudioAntispamFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.AudioScenes[*]" {
			e := jsonata.MustCompile("ApiOutput.AudioScenes")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.Buckets[*]" {
			e := jsonata.MustCompile("ApiOutput.Buckets")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.ImageAdFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.ImageAdFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.ImageLiveFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.ImageLiveFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.ImagePornFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.ImagePornFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.ImageTerrorismFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.ImageTerrorismFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoAdFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.VideoAdFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoLiveFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.VideoLiveFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoPornFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.VideoPornFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoScenes[*]" {
			e := jsonata.MustCompile("ApiOutput.VideoScenes")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoTerrorismFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.VideoTerrorismFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.VideoVoiceAntispamFreezeConfig.Type" {
			e := jsonata.MustCompile("ApiOutput.VideoVoiceAntispamFreezeConfig")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAligreenOssStockTask >>> Encapsulated.
