package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/oos"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type OosService struct {
	client *connectivity.AliyunClient
}

func (s *OosService) DescribeOosTemplate(id string) (object oos.Template, err error) {
	request := oos.CreateGetTemplateRequest()
	request.RegionId = s.client.RegionId

	request.TemplateName = id

	raw, err := s.client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
		return oosClient.GetTemplate(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OosTemplate", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*oos.GetTemplateResponse)
	return response.Template, nil
}
