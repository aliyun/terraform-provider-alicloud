package alicloud

import (
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type Waf_openapiService struct {
	client *connectivity.AliyunClient
}

func (s *Waf_openapiService) DescribeWafDomain(id string) (wafDomain waf_openapi.DescribeDomainResponse, err error) {
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
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*waf_openapi.DescribeDomainResponse)
	return *response, nil
}
