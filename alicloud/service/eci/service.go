package eci

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeEciImageCache(id string) (object eci.DescribeImageCachesImageCache0, err error) {
	request := eci.CreateDescribeImageCachesRequest()
	request.RegionId = s.client.RegionId

	request.ImageCacheId = id

	raw, err := s.client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*eci.DescribeImageCachesResponse)

	if len(response.ImageCaches) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("EciImageCache", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}
	return response.ImageCaches[0], nil
}

func (s *Service) EciImageCacheStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciImageCache(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *Service) DescribeEciContainerGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeContainerGroups"
	jsonId, err := json.Marshal([]string{id})
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"ContainerGroupIds": string(jsonId),
	}
	response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("EciContainerGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.ContainerGroups", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ContainerGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ECIOpenAPI", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ContainerGroupId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ECIOpenAPI", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) EciContainerGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciContainerGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *Service) DescribeEciVirtualNode(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVirtualNodes"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
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
		v, err := jsonpath.Get("$.VirtualNodes", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VirtualNodes", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("ECI", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["VirtualNodeId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ECI", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *Service) EciVirtualNodeStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciVirtualNode(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
