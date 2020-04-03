// Package alicloud common functions used by maxcompute
package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type MaxComputeService struct {
	client *connectivity.AliyunClient
}

func (s *MaxComputeService) DescribeMaxComputeProject(id string) (*maxcompute.GetProjectResponse, error) {
	response := &maxcompute.GetProjectResponse{}
	request := maxcompute.CreateGetProjectRequest()

	request.RegionName = s.client.RegionId
	request.ProjectName = id

	raw, err := s.client.WithMaxComputeClient(func(MaxComputeClient *maxcompute.Client) (interface{}, error) {
		return MaxComputeClient.GetProject(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*maxcompute.GetProjectResponse)

	if response.Code != "200" {
		if isProjectNotExistError(response.Data) {
			return response, WrapErrorf(err, NotFoundMsg, AliyunMaxComputeSdkGo)
		}

		return response, WrapError(Error("%v", response))
	}

	return response, nil
}

func (s *MaxComputeService) WaitForMaxComputeProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		response, err := s.DescribeMaxComputeProject(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, response.Data, id, ProviderERROR)
		}

	}
}

func isProjectNotExistError(data string) bool {
	if strings.Contains(data, "Project not found") {
		return true
	}

	return false
}
