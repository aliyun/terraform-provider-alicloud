package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DdosbgpService struct {
	client *connectivity.AliyunClient
}

func (s *DdosbgpService) DescribeDdosbgpInstance(id string, region string) (v ddosbgp.Instance, err error) {
	request := ddosbgp.CreateDescribeInstanceListRequest()
	request.RegionId = region
	request.DdosRegionId = region
	request.InstanceIdList = "[\"" + id + "\"]"
	request.PageNo = "1"
	request.PageSize = "10"

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.DescribeInstanceList(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdosbgpInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		resp, _ := raw.(*ddosbgp.DescribeInstanceListResponse)
		if resp == nil || len(resp.InstanceList) == 0 || resp.InstanceList[0].InstanceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("Ddosbgp Instance", id)), NotFoundMsg, ProviderERROR)
		}

		v = resp.InstanceList[0]
		return nil
	})

	return v, WrapError(err)
}

func (s *DdosbgpService) DescribeDdosbgpInstanceSpec(id string, region string) (v ddosbgp.InstanceSpec, err error) {
	request := ddosbgp.CreateDescribeInstanceSpecsRequest()
	request.InstanceIdList = "[\"" + id + "\"]"
	request.DdosRegionId = region
	request.RegionId = region

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.DescribeInstanceSpecs(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdosbgpInstanceNotFound, InvalidDdosbgpInstance}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		resp, _ := raw.(*ddosbgp.DescribeInstanceSpecsResponse)
		if resp == nil || len(resp.InstanceSpecs) == 0 || resp.InstanceSpecs[0].InstanceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("Ddosbgp Instance", id)), NotFoundMsg, ProviderERROR)
		}

		v = resp.InstanceSpecs[0]
		return nil
	})

	return v, WrapError(err)
}

func (s *DdosbgpService) UpdateDdosbgpInstanceName(instanceId string, name string, region string) error {
	request := ddosbgp.CreateModifyRemarkRequest()
	request.InstanceId = instanceId
	request.RegionId = region
	request.ResourceRegionId = region

	request.Remark = name

	if _, err := s.client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
		return ddosbgpClient.ModifyRemark(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
