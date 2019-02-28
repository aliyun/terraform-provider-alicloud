package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type NasService struct {
	client *connectivity.AliyunClient
}

func (s *NasService) DescribeFileSystems(fileSystemId string) (fs nas.FileSystem, err error) {

	args := nas.CreateDescribeFileSystemsRequest()
	args.FileSystemId = fileSystemId
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidFileSystemIDNotFound, ForbiddenNasNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, fileSystemId, args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*nas.DescribeFileSystemsResponse)
		if resp == nil || len(resp.FileSystems.FileSystem) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("fileSystemId", fileSystemId))
		}
		fs = resp.FileSystems.FileSystem[0]
		return nil
	})
	return
}

func (s *NasService) DescribeMountTargets(fileSystemId string) (fs nas.MountTarget, err error) {

	args := nas.CreateDescribeMountTargetsRequest()
	args.RegionId = string(s.client.Region)
	args.FileSystemId = fileSystemId
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeMountTargetsResponse)
		if resp == nil || len(resp.MountTargets.MountTarget) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("fileSystemId", fileSystemId))
		}
		fs = resp.MountTargets.MountTarget[0]
		return nil
	})
	return
}

func (s *NasService) DescribeAccessGroup(accessGroupName string) (ag nas.AccessGroup, err error) {

	args := nas.CreateDescribeAccessGroupsRequest()
	args.RegionId = string(s.client.Region)
	args.AccessGroupName = accessGroupName

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("accessGroupName", accessGroupName))
		}
		ag = resp.AccessGroups.AccessGroup[0]
		return nil
	})
	return
}

func (s *NasService) DescribeAccessRules(accessGroupName string) (fs nas.AccessRule, err error) {

	args := nas.CreateDescribeAccessRulesRequest()
	args.RegionId = string(s.client.Region)
	args.AccessGroupName = accessGroupName

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeAccessRulesResponse)
		if resp == nil || len(resp.AccessRules.AccessRule) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("accessGroupName", accessGroupName))
		}
		fs = resp.AccessRules.AccessRule[0]
		return nil
	})
	return
}
