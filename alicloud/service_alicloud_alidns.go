package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type AlidnsService struct {
	client *connectivity.AliyunClient
}

func (s *AlidnsService) DescribeAlidnsDomainGroup(id string) (object alidns.DomainGroup, err error) {
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)

		if len(response.DomainGroups.DomainGroup) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("AlidnsDomainGroup", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		for _, object := range response.DomainGroups.DomainGroup {
			if object.GroupId == id {
				return object, nil
			}
		}
		if len(response.DomainGroups.DomainGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("AlidnsDomainGroup", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *AlidnsService) DescribeAlidnsRecord(id string) (object alidns.DescribeDomainRecordInfoResponse, err error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RegionId = s.client.RegionId

	request.RecordId = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("AlidnsRecord", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainRecordInfoResponse)
	return *response, nil
}
