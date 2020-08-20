package alicloud

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DcdnService struct {
	client *connectivity.AliyunClient
}

func (s *DcdnService) convertSourcesToString(v []interface{}) (string, error) {
	arrayMaps := make([]dcdn.Source, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = dcdn.Source{
			Content:  item["content"].(string),
			Port:     item["port"].(int),
			Priority: item["priority"].(string),
			Type:     item["type"].(string),
			Weight:   item["weight"].(string),
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *DcdnService) DescribeDcdnDomainCertificateInfo(id string) (object dcdn.CertInfo, err error) {
	request := dcdn.CreateDescribeDcdnDomainCertificateInfoRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
		return dcdnClient.DescribeDcdnDomainCertificateInfo(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*dcdn.DescribeDcdnDomainCertificateInfoResponse)

	if len(response.CertInfos.CertInfo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DcdnDomain", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.CertInfos.CertInfo[0], nil
}

func (s *DcdnService) DescribeDcdnDomain(id string) (object dcdn.DomainDetail, err error) {
	request := dcdn.CreateDescribeDcdnDomainDetailRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
		return dcdnClient.DescribeDcdnDomainDetail(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DcdnDomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*dcdn.DescribeDcdnDomainDetailResponse)
	return response.DomainDetail, nil
}

func (s *DcdnService) DcdnDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DomainStatus == failState {
				return object, object.DomainStatus, WrapError(Error(FailedToReachTargetStatus, object.DomainStatus))
			}
		}
		return object, object.DomainStatus, nil
	}
}
