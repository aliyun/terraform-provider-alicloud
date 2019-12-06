package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DrdsService struct {
	client *connectivity.AliyunClient
}

func (s *DrdsService) DescribeDrdsInstance(id string) (*drds.DescribeDrdsInstanceResponse, error) {
	response := &drds.DescribeDrdsInstanceResponse{}
	request := drds.CreateDescribeDrdsInstanceRequest()
	request.RegionId = s.client.RegionId
	request.DrdsInstanceId = id
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(request)
	})

	if err != nil {
		if IsExceptedError(err, InvalidDRDSInstanceIdNotFound) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*drds.DescribeDrdsInstanceResponse)
	if response.Data.Status == "RELEASE" {
		return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *DrdsService) DrdsInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, WrapError(Error(FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}
