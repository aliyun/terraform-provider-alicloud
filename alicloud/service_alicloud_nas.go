package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type NasService struct {
	client *connectivity.AliyunClient
}

func (s *NasService) DescribeNasFileSystem(id string) (fs nas.FileSystem, err error) {

	request := nas.CreateDescribeFileSystemsRequest()
	request.FileSystemId = id
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidFileSystemIDNotFound, ForbiddenNasNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*nas.DescribeFileSystemsResponse)
		if response == nil || len(response.FileSystems.FileSystem) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("File System", id)), NotFoundMsg, ProviderERROR)
		}
		fs = response.FileSystems.FileSystem[0]
		return nil
	})
	return
}

func (s *NasService) DescribeNasMountTarget(id string) (fs nas.MountTarget, err error) {

	request := nas.CreateDescribeMountTargetsRequest()
	request.RegionId = string(s.client.Region)
	split := strings.Split(id, "-")
	request.FileSystemId = split[0]
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(request)
		})
		if err != nil {
			if IsExceptedErrors(err, NasNotFound) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*nas.DescribeMountTargetsResponse)
		if response == nil || len(response.MountTargets.MountTarget) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Mount Target", id)), NotFoundMsg, ProviderERROR)
		}
		fs = response.MountTargets.MountTarget[0]
		return nil
	})
	return
}

func (s *NasService) DescribeNasAccessGroup(id string) (ag nas.AccessGroup, err error) {

	request := nas.CreateDescribeAccessGroupsRequest()
	request.RegionId = string(s.client.Region)
	request.AccessGroupName = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidAccessGroupNotFound, ForbiddenNasNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)

		}
		response, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if response == nil || len(response.AccessGroups.AccessGroup) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Access Group", id)), NotFoundMsg, ProviderERROR)
		}
		ag = response.AccessGroups.AccessGroup[0]
		return nil
	})
	return
}

func (s *NasService) DescribeNasAccessRule(id string) (fs nas.AccessRule, err error) {

	request := nas.CreateDescribeAccessRulesRequest()
	request.RegionId = string(s.client.Region)
	request.AccessGroupName = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidAccessGroupNotFound, ForbiddenNasNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*nas.DescribeAccessRulesResponse)
		if response == nil || len(response.AccessRules.AccessRule) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Access Rule", id)), NotFoundMsg, ProviderERROR)
		}
		fs = response.AccessRules.AccessRule[0]
		return nil
	})
	return
}

func (s *NasService) WaitForMountTarget(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		mt, err := s.DescribeNasMountTarget(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
		if mt.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return WrapError(Error(GetTimeoutMessage("Nas MountTarget", string(status))))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
