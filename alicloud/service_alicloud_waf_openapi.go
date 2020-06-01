package alicloud

import (
	"encoding/json"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type Waf_openapiService struct {
	client *connectivity.AliyunClient
}

func (s *Waf_openapiService) convertLogHeadersToString(v []interface{}) (string, error) {
	arrayMaps := make([]waf_openapi.LogHeader, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = waf_openapi.LogHeader{
			K: item["key"].(string),
			V: item["value"].(string),
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *Waf_openapiService) DescribeWafDomain(id string) (object waf_openapi.Domain, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := waf_openapi.CreateDescribeDomainRequest()
	request.RegionId = s.client.RegionId
	request.Domain = parts[1]
	request.InstanceId = parts[0]

	raw, err := s.client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DescribeDomain(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNotExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("WafDomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*waf_openapi.DescribeDomainResponse)
	return response.Domain, nil
}

func (s *Waf_openapiService) DescribeWafInstance(id string) (object waf_openapi.DescribeInstanceInfoResponse, err error) {
	request := waf_openapi.CreateDescribeInstanceInfoRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DescribeInstanceInfo(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*waf_openapi.DescribeInstanceInfoResponse)
	if response != nil && response.InstanceInfo.InstanceId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("WafInstance", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}
