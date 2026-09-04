// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type EcdServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeEcdDesktopGroup <<< Encapsulated get interface for Ecd DesktopGroup.

func (s *EcdServiceV2) DescribeEcdDesktopGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DesktopGroupId"] = id
	request["RegionId"] = client.RegionId
	action := "GetDesktopGroupDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Desktops", response)
	if err != nil || v == nil {
		// GetDesktopGroupDetail returns a response without the Desktops key
		// when the desktop group does not exist.
		return object, WrapErrorf(NotFoundErr("DesktopGroup", id), NotFoundMsg, response)
	}

	desktops, ok := v.(map[string]interface{})
	if !ok || len(desktops) == 0 {
		return object, WrapErrorf(NotFoundErr("DesktopGroup", id), NotFoundMsg, response)
	}

	return desktops, nil
}
func (s *EcdServiceV2) DescribeDesktopGroupDescribeUsersInGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeUsersInGroup"
	endUserIds := make([]interface{}, 0)
	nextToken := ""

	for {
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["DesktopGroupId"] = id
		request["RegionId"] = client.RegionId
		request["MaxResults"] = PageSizeLarge
		if nextToken != "" {
			request["NextToken"] = nextToken
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		v, _ := jsonpath.Get("$.EndUsers[*]", response)
		if users, ok := v.([]interface{}); ok {
			for _, user := range users {
				if userMap, ok := user.(map[string]interface{}); ok {
					if userId, ok := userMap["EndUserId"]; ok && userId != nil {
						endUserIds = append(endUserIds, userId)
					}
				}
			}
		}

		if token, ok := response["NextToken"].(string); ok && token != "" {
			nextToken = token
		} else {
			break
		}
	}

	sort.Slice(endUserIds, func(i, j int) bool {
		return fmt.Sprint(endUserIds[i]) < fmt.Sprint(endUserIds[j])
	})
	object = map[string]interface{}{
		"EndUserIds": endUserIds,
	}
	return object, nil
}

func (s *EcdServiceV2) EcdDesktopGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.EcdDesktopGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeEcdDesktopGroup)
}

func (s *EcdServiceV2) EcdDesktopGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeEcdDesktopGroup >>> Encapsulated.
