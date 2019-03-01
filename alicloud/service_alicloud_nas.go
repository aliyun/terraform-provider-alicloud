package alicloud

import (
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
		resp, _ := raw.(*nas.DescribeFileSystemsResponse)
		if resp == nil || len(resp.FileSystems.FileSystem) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("File System", id)), NotFoundMsg, ProviderERROR)
		}
		fs = resp.FileSystems.FileSystem[0]
		return nil
	})
	return
}

func (s *NasService) DescribeNasMountTarget(id string) (fs nas.MountTarget, err error) {

	request := nas.CreateDescribeMountTargetsRequest()
	request.RegionId = string(s.client.Region)
	request.FileSystemId = id
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*nas.DescribeMountTargetsResponse)
		if resp == nil || len(resp.MountTargets.MountTarget) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Mount Target", id)), NotFoundMsg, ProviderERROR)
		}
		fs = resp.MountTargets.MountTarget[0]
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
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)

		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Access Group", id)), NotFoundMsg, ProviderERROR)
		}
		ag = resp.AccessGroups.AccessGroup[0]
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
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)

		}
		resp, _ := raw.(*nas.DescribeAccessRulesResponse)
		if resp == nil || len(resp.AccessRules.AccessRule) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("Access Rule", id)), NotFoundMsg, ProviderERROR)
		}
		fs = resp.AccessRules.AccessRule[0]
		return nil
	})
	return
}
