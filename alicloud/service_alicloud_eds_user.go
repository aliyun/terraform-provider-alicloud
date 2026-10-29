package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type EdsUserService struct {
	client *connectivity.AliyunClient
}

func (s *EdsUserService) DescribeEcdUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeUsers"
	// Filter by EndUserIds so DescribeUsers returns only the requested user
	// instead of paging through the whole list and matching EndUserId in Go.
	// The list traversal used to miss a freshly created user until it became
	// visible in the full list, which made Create's waiter time out with
	// "couldn't find resource".
	request := map[string]interface{}{
		"EndUserIds": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	}
	for _, user := range v.([]interface{}) {
		if fmt.Sprint(user.(map[string]interface{})["EndUserId"]) == id {
			return user.(map[string]interface{}), nil
		}
	}
	return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
}

func (s *EdsUserService) EcdUserStateRefreshFunc(id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdUser(id)
		if err != nil {
			if NotFoundError(err) {
				// A freshly created user is not always visible to DescribeUsers
				// immediately. Return a pending state instead of an empty one so
				// the waiter keeps retrying until the user becomes visible, rather
				// than failing with "couldn't find resource" before the create
				// timeout elapses.
				return nil, "Pending", nil
			}
			return nil, "", WrapError(err)
		}
		statusCode, formatErr := formatInt(object["Status"])
		if formatErr != nil {
			return nil, "", formatErr
		}
		status := convertEcdUserStatusResponse(strconv.Itoa(statusCode))
		return object, status, nil
	}
}
func (s *EdsUserService) DescribeEcdCustomProperty(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListProperty"
	request := map[string]interface{}{}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	properties, ok := v.(map[string]interface{})["Properties"]
	if !ok {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	}

	if len(properties.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	}
	for _, v := range properties.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["PropertyId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	}
	return object, nil
}
