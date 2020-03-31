package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EciService struct {
	client *connectivity.AliyunClient
}

func (s *EciService) DescribeEciopenapiImageCache(id string) (object eci.DescribeImageCachesImageCache0, err error) {
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
	if len(response.ImageCaches)<1 {
		err=WrapErrorf(Error(GetNotFoundMessage("", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.ImageCaches[0], nil
}
/*
func (s *EciService) DescribeEciopenapiContainerGroup(id string) (object eci.ContainerGroups, err error) {
request:=eci.CreateDescribeContainerGroupsRequest()
request.RegionId = s.client.RegionId

request. = id

raw, err := s.client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
	return eciClient.DescribeContainerGroups(request)
})
if err != nil {
	err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	return
}
addDebug(request.GetActionName(), raw, request.RpcRequest, request)
response, _ := raw.(*eci.DescribeContainerGroupsResponse)

if len(response.ContainerGroups) < 1 {
	err = WrapErrorf(Error(GetNotFoundMessage("EciopenapiContainerGroup", id)), NotFoundMsg, ProviderERROR)
	return
}
return response.ContainerGroups[0], nil
}
*/
