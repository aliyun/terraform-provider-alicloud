package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EciService struct {
	client *connectivity.AliyunClient
}

func (s *EciService) DescribeEciImageCache(id string) (object eci.DescribeImageCachesImageCache0, err error) {
	request := eci.CreateDescribeImageCachesRequest()
	request.RegionId = s.client.RegionId

	request.ImageCacheId = id

	raw, err := s.client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*eci.DescribeImageCachesResponse)

	if len(response.ImageCaches) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("EciImageCache", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.ImageCaches[0], nil
}

func (s *EciService) EciImageCacheStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciImageCache(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}
