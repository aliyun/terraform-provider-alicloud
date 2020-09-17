package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type Dms_enterpriseService struct {
	client *connectivity.AliyunClient
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseInstance(id string) (object dms_enterprise.Instance, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := dms_enterprise.CreateGetInstanceRequest()
	request.RegionId = s.client.RegionId
	request.Host = parts[0]
	if v, err := strconv.Atoi(parts[1]); err == nil {
		request.Port = requests.NewInteger(v)
	} else {
		err = WrapError(err)
		return object, err
	}

	raw, err := s.client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
		return dms_enterpriseClient.GetInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DmsEnterpriseInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*dms_enterprise.GetInstanceResponse)
	return response.Instance, nil
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseUser(id string) (object dms_enterprise.User, err error) {
	request := dms_enterprise.CreateGetUserRequest()
	request.RegionId = s.client.RegionId

	if v, err := strconv.Atoi(id); err == nil {
		request.Uid = requests.NewInteger(v)
	} else {
		err = WrapError(err)
		return object, err
	}

	raw, err := s.client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
		return dms_enterpriseClient.GetUser(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*dms_enterprise.GetUserResponse)
	return response.User, nil
}
