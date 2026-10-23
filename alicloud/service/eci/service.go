package eci

import (
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
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
